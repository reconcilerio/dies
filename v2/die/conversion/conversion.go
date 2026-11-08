/*
Copyright 2026 the original author or authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package conversion converts a resource between its versions, as a CRD's
// conversion webhook does, with jigs. Package
// [reconciler.io/dies/v2/die/conversion/webhook] serves a [Converter] as one.
//
// One version is the hub, usually the storage version, and every other
// version is a spoke, converted to and from the hub: a conversion between two
// spokes goes through it. A spoke is a jig of the spoke's type over the hub's,
// so one set of bindings converts both ways. Holding a hub as a spoke reads
// it as the spoke, and writing a spoke through the jig into an empty hub
// writes the hub. Each binding is copied only where its source has a value,
// so what is unset stays unset. The metadata is bound as it is, as the API
// server allows conversion to change only labels and annotations.
//
// The jig's transforms are lossless, but its bindings needn't cover each
// field: a field of one version may have no place in the other. What a
// conversion drops, it stashes, in an annotation of the converted object
// named for the version it came from (see [StashAnnotation] and
// [Converter.WithStashAnnotation]), as the part of
// the original that a conversion back wouldn't restore. Converting back
// restores the stash, filling only what the conversion leaves unset, and
// removes the annotation, so a round trip changes nothing. A stash keeps
// whole values: a list whose elements lose a field is stashed whole, and
// isn't restored over the list the conversion back writes.
package conversion

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"reflect"
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation"

	"reconciler.io/dies/v2/die/base"
	"reconciler.io/dies/v2/die/fields"
	"reconciler.io/dies/v2/die/internal/core"
	"reconciler.io/dies/v2/die/jigs"
)

// StashAnnotation is the prefix of the annotations that keep what a
// conversion from a version dropped, followed by the version, such as
// conversion.dies.reconciler.io/v1 for what a conversion from v1 dropped. A
// converter uses it unless given another, with
// [Converter.WithStashAnnotation].
const StashAnnotation = "conversion.dies.reconciler.io/"

// spoke converts a version of a resource to and from its hub version.
type spoke struct {
	gvk      schema.GroupVersionKind
	newSpoke func() base.MutableDie
	newHub   func() base.MutableDie
	hold     func(hub base.MutableDie) (base.MutableDie, error)
	paths    []core.Path
}

// toHub writes a spoke into an empty hub, through the jig.
func (s *spoke) toHub(d base.MutableDie) (base.MutableDie, error) {
	hub := s.newHub()
	view, err := s.hold(hub)
	if err != nil {
		return nil, err
	}
	return hub, s.copy(d, view)
}

// fromHub reads a hub as a new spoke, through the jig.
func (s *spoke) fromHub(hub base.MutableDie) (base.MutableDie, error) {
	view, err := s.hold(hub)
	if err != nil {
		return nil, err
	}
	d := s.newSpoke()
	return d, s.copy(view, d)
}

// copy copies the value of each binding, of the spoke's type, that from has
// to to.
func (s *spoke) copy(from, to base.Die) error {
	fh, fat, ok := core.StateOf(from)
	th, tat, ok2 := core.StateOf(to)
	if !ok || !ok2 {
		return core.Errorf("conversion: %T or %T is not a die", from, to)
	}
	for _, p := range s.paths {
		has, err := fh.Has(core.JoinPath(fat, p))
		if err != nil {
			return err
		}
		if !has {
			continue
		}
		v, err := fh.ValueAt(core.JoinPath(fat, p))
		if err != nil {
			return err
		}
		if _, err := th.Write(tat, p, core.SetOp(core.OpSet, v), false); err != nil {
			return fmt.Errorf("conversion: writing %s: %w", p, err)
		}
	}
	return nil
}

// Converter converts a resource between its hub version and its spokes. It
// is defined once, typically as a package variable, from the hub's blank,
// with a spoke for each other version:
//
//	var Conversion = conversion.New(v2.SprocketBlank).
//		AddSpoke(v1.SprocketBlank,
//			jigs.Bind(v1.SprocketBlank.Spec().Size(), v2.SprocketBlank.Spec().Teeth()),
//		)
//
// A definition that can't convert, such as a binding of the wrong types, is a
// programming error: the methods that define a converter panic rather than
// return it. Each returns a new converter, leaving the one it's called on as
// it was.
type Converter struct {
	gvk    schema.GroupVersionKind
	newHub func() base.MutableDie
	// hub is the hub's blank, which a spoke's jig is clamped to, and
	hub    base.Located
	spokes map[schema.GroupVersion]*spoke
	// stash is the prefix of the stash annotations
	stash string
}

// New returns the converter for a resource whose hub version is given as its
// blank, with no spokes. The hub's group, version and kind are those its die
// package records from its +die:object marker, not those of the blank's
// value, and converted objects start from an empty die of its type, so a
// value stamped on the blank isn't part of them.
func New[H any, HD fields.DieObjectOf[H, HD, HDM], HDM fields.DieObjectMutOf[H, HD, HDM]](hub HD) *Converter {
	gvk := core.Must(identity[H]())
	return &Converter{
		gvk:    gvk,
		newHub: empty[H, HD, HDM](gvk),
		hub:    hub,
		spokes: map[schema.GroupVersion]*spoke{},
		stash:  StashAnnotation,
	}
}

// AddSpoke returns the converter with another version, given as its blank,
// converted to and from the hub by bindings of the spoke's locations to the
// hub's, as a jig binds them: [jigs.Bind] for a value that's the same in each
// version, [jigs.BindVia] for one that is changed. The metadata is bound
// already. The spoke's group, version and kind are its type's, as the hub's
// are.
//
// It panics if the spoke isn't another version of the hub's group and kind,
// already has a conversion, or a binding can't be made.
func (c *Converter) AddSpoke[S any, SD fields.DieObjectOf[S, SD, SDM], SDM fields.DieObjectMutOf[S, SD, SDM]](blank SD, bindings ...jigs.Binding) *Converter {
	gvk := core.Must(identity[S]())
	gv := gvk.GroupVersion()
	if gvk.GroupKind() != c.gvk.GroupKind() || gv == c.gvk.GroupVersion() {
		panic(core.Errorf("conversion: %s is not another version of %s", gvk, c.gvk))
	}
	if c.spokes[gv] != nil {
		panic(core.Errorf("conversion: %s has more than one conversion", gv))
	}
	core.Check(checkStash(c.stash, gv))
	bindings = append([]jigs.Binding{bindMetadata(blank, c.hub)}, bindings...)
	jig := jigs.NewJig(blank).MustClamp(c.hub, bindings...)
	s := &spoke{
		gvk:      gvk,
		newSpoke: empty[S, SD, SDM](gvk),
		newHub:   c.newHub,
		hold: func(h base.MutableDie) (base.MutableDie, error) {
			v, ok, err := jig.HoldMut(h)
			if err == nil && !ok {
				err = core.Errorf("conversion: the jig has no bindings for %T", h)
			}
			return v, err
		},
	}
	for _, b := range bindings {
		view, _, _, _ := core.BindingInfo(b)
		s.paths = append(s.paths, view)
	}
	out := c.clone()
	out.spokes[gv] = s
	return out
}

// WithStashAnnotation returns the converter with stash annotations of another
// prefix, followed by the version, as [StashAnnotation] is:
// "sprockets.example.dev/dropped-" names the stash of v1
// sprockets.example.dev/dropped-v1. It panics if that isn't a valid
// annotation for each version.
//
// The prefix is part of the resource's stored data: a stash under another
// prefix, such as one made before it changed, isn't restored.
func (c *Converter) WithStashAnnotation(prefix string) *Converter {
	versions := []schema.GroupVersion{c.gvk.GroupVersion()}
	for gv := range c.spokes {
		versions = append(versions, gv)
	}
	core.Check(checkStash(prefix, versions...))
	out := c.clone()
	out.stash = prefix
	return out
}

// checkStash returns an error unless a stash prefix makes a valid annotation
// for each version.
func checkStash(prefix string, versions ...schema.GroupVersion) error {
	for _, gv := range versions {
		key := prefix + gv.Version
		if errs := validation.IsQualifiedName(key); len(errs) > 0 {
			return core.Errorf("conversion: the stash annotation %q is not valid: %s", key, strings.Join(errs, "; "))
		}
	}
	return nil
}

func (c *Converter) clone() *Converter {
	out := *c
	out.spokes = maps.Clone(c.spokes)
	return &out
}

// StashAnnotation is the annotation that keeps what a conversion from a
// version dropped.
func (c *Converter) StashAnnotation(gv schema.GroupVersion) string {
	return c.stash + gv.Version
}

// GroupKind is the group and kind of the resource the converter converts.
func (c *Converter) GroupKind() schema.GroupKind { return c.gvk.GroupKind() }

// Convert converts an object, as JSON, to a version of its group.
func (c *Converter) Convert(obj []byte, to schema.GroupVersion) ([]byte, error) {
	m, err := decode(obj)
	if err != nil {
		return nil, err
	}
	wireGV, _ := m["apiVersion"].(string)
	kind, _ := m["kind"].(string)
	from, err := schema.ParseGroupVersion(wireGV)
	if err != nil {
		return nil, core.Errorf("conversion: %w", err)
	}
	if gk := from.WithKind(kind).GroupKind(); gk != c.gvk.GroupKind() {
		return nil, core.Errorf("conversion: can't convert a %s, only a %s", gk, c.gvk.GroupKind())
	}
	for _, v := range []schema.GroupVersion{from, to} {
		if !c.has(v) {
			return nil, core.Errorf("conversion: %s is not a version of %s", v, c.gvk.GroupKind())
		}
	}
	if from == to {
		return obj, nil
	}
	hub := c.gvk.GroupVersion()
	if from != hub {
		if m, err = c.toHub(c.spokes[from], m); err != nil {
			return nil, fmt.Errorf("conversion: from %s: %w", from, err)
		}
	}
	if to != hub {
		if m, err = c.fromHub(c.spokes[to], m); err != nil {
			return nil, fmt.Errorf("conversion: to %s: %w", to, err)
		}
	}
	return json.Marshal(m)
}

// has reports whether a group version is the hub's or a spoke's.
func (c *Converter) has(gv schema.GroupVersion) bool {
	return gv == c.gvk.GroupVersion() || c.spokes[gv] != nil
}

// toHub converts a spoke to the hub, stashing what the hub can't hold, and
// restoring what the hub stashed when it was converted to the spoke.
func (c *Converter) toHub(s *spoke, obj map[string]any) (map[string]any, error) {
	stashed, err := unstash(obj, c.StashAnnotation(c.gvk.GroupVersion()))
	if err != nil {
		return nil, err
	}
	return convert(obj, s.newSpoke, s.toHub, s.fromHub, stashed, c.StashAnnotation(s.gvk.GroupVersion()))
}

// fromHub converts the hub to a spoke, the other way.
func (c *Converter) fromHub(s *spoke, hub map[string]any) (map[string]any, error) {
	stashed, err := unstash(hub, c.StashAnnotation(s.gvk.GroupVersion()))
	if err != nil {
		return nil, err
	}
	return convert(hub, c.newHub, s.fromHub, s.toHub, stashed, c.StashAnnotation(c.gvk.GroupVersion()))
}

// convert converts an object one way, restoring what was stashed from the
// version it's converted to, and stashing what converting it back would
// drop, in the annotation key.
func convert(obj map[string]any, blank func() base.MutableDie, there, back func(base.MutableDie) (base.MutableDie, error), stashed map[string]any, key string) (map[string]any, error) {
	d := blank()
	if err := feed(d, obj); err != nil {
		return nil, err
	}
	converted, err := there(d)
	if err != nil {
		return nil, err
	}
	returned, err := back(converted)
	if err != nil {
		return nil, err
	}
	out, err := release(converted)
	if err != nil {
		return nil, err
	}
	again, err := release(returned)
	if err != nil {
		return nil, err
	}
	fill(out, stashed)
	if dropped := drops(obj, again); len(dropped) > 0 {
		if err := stash(out, key, dropped); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// drops returns what of want isn't in got: the fields whose values got
// doesn't have, recursively within objects.
func drops(want, got map[string]any) map[string]any {
	out := map[string]any{}
	for k, w := range want {
		g, ok := got[k]
		if !ok {
			out[k] = w
			continue
		}
		wm, wok := w.(map[string]any)
		gm, gok := g.(map[string]any)
		if wok && gok {
			if d := drops(wm, gm); len(d) > 0 {
				out[k] = d
			}
		}
	}
	return out
}

// fill sets what of from isn't in into, recursively within objects.
func fill(into, from map[string]any) {
	for k, f := range from {
		i, ok := into[k]
		if !ok {
			into[k] = f
			continue
		}
		im, iok := i.(map[string]any)
		fm, fok := f.(map[string]any)
		if iok && fok {
			fill(im, fm)
		}
	}
}

// stash records what a conversion dropped, as the annotation key of the
// converted object.
func stash(obj map[string]any, key string, dropped map[string]any) error {
	b, err := json.Marshal(dropped)
	if err != nil {
		return err
	}
	md, _ := obj["metadata"].(map[string]any)
	if md == nil {
		md = map[string]any{}
		obj["metadata"] = md
	}
	ann, _ := md["annotations"].(map[string]any)
	if ann == nil {
		ann = map[string]any{}
		md["annotations"] = ann
	}
	ann[key] = string(b)
	return nil
}

// unstash removes and returns what a conversion stashed in an object, as the
// annotation key.
func unstash(obj map[string]any, key string) (map[string]any, error) {
	md, _ := obj["metadata"].(map[string]any)
	ann, _ := md["annotations"].(map[string]any)
	raw, ok := ann[key].(string)
	if !ok {
		return nil, nil
	}
	delete(ann, key)
	if len(ann) == 0 {
		delete(md, "annotations")
	}
	stashed, err := decode([]byte(raw))
	if err != nil {
		return nil, fmt.Errorf("conversion: the annotation %s: %w", key, err)
	}
	return stashed, nil
}

func feed(d base.MutableDie, obj map[string]any) error {
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, d)
}

func release(d base.Die) (map[string]any, error) {
	b, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	return decode(b)
}

// decode decodes a JSON object, keeping numbers as they're written.
func decode(b []byte) (map[string]any, error) {
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	var m map[string]any
	if err := dec.Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}

// bindMetadata binds a spoke's metadata to the hub's. An object die's
// metadata is a metav1.ObjectMeta at metadata, as its base requires, so it's
// bound by its path, without naming the die of meta/v1 it has.
func bindMetadata(spoke, hub base.Located) jigs.Binding {
	vroot, verr := core.RootOf(spoke)
	root, err := core.RootOf(hub)
	path := core.Path{core.FieldStep("metadata")}
	return core.WithErr(core.NewBinding(path, path, vroot, root, nil, nil), errors.Join(verr, err))
}

// identity returns the group, version and kind an object type's die package
// recorded from its +die:object marker.
func identity[T any]() (schema.GroupVersionKind, error) {
	t := reflect.TypeFor[T]()
	gvk, ok, err := core.IdentityOf(t)
	if err != nil {
		return schema.GroupVersionKind{}, err
	}
	if !ok {
		return schema.GroupVersionKind{}, core.Errorf("conversion: %s has no group, version and kind: mark it +die:object:apiVersion=…,kind=…", t)
	}
	return gvk, nil
}

// empty returns a function that makes an empty mutable die of an object type,
// with only its group, version and kind.
func empty[T any, D fields.DieObjectOf[T, D, M], M fields.DieObjectMutOf[T, D, M]](gvk schema.GroupVersionKind) func() base.MutableDie {
	return func() base.MutableDie {
		d := base.Blank[T, D, M]().DieMutable()
		d.GetObjectKind().SetGroupVersionKind(gvk)
		return d
	}
}
