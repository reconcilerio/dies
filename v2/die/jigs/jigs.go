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

// Package jigs holds dies and unstructured objects in the shape of another
// type, so the same tools work on each of them. As a jig in a workshop holds
// material in place for an operation, a jig here holds a source die so that
// it can be worked through a die of the jig's type: the locations of the
// jig's type are bound to those of each source type or kind, directly,
// through transforms, or as an aggregate of several sources.
package jigs

import (
	"errors"
	"fmt"
	"reflect"

	"reconciler.io/dies/v2/die/base"
	"reconciler.io/dies/v2/die/internal/core"
	"reconciler.io/dies/v2/die/internal/dies"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// A jig holds a source die as a die of another type, V. Each location of V
// is bound to a location of the source, separately for each source type, so
// one jig can hold different types the same way: the template of a
// Deployment is at .spec.template, of a CronJob at
// .spec.jobTemplate.spec.template.
//
// A jig holds no state of its own. Reads and writes go to the source, so a
// write through a jig is recorded as intent on the source, and the source's
// patches carry it. The jig of an immutable die is an immutable die, which
// reads; of a mutable die, a mutable die, which also writes.

// Jig holds sources as dies of type V. Typed sources are bound by their
// type; unstructured objects by their kind, with a fallback for kinds without
// bindings of their own.
type Jig[V any, VD base.DieOf[V, VD, VDM], VDM base.DieMutOf[V, VD, VDM]] struct {
	sources map[reflect.Type][]Binding
	kinds   map[schema.GroupVersionKind][]Binding
	unknown []Binding
}

// Binding binds a location of a jig to a location of a source: made by
// [Bind], [BindVia], [BindPath], [BindPathVia] and [BindAll].
type Binding = core.Binding

// NewJig returns a jig with no sources.
func NewJig[V any, VD base.DieOf[V, VD, VDM], VDM base.DieMutOf[V, VD, VDM]]() *Jig[V, VD, VDM] {
	return &Jig[V, VD, VDM]{sources: map[reflect.Type][]Binding{}, kinds: map[schema.GroupVersionKind][]Binding{}}
}

// Bind binds a location of a jig to a location of a source. Both are
// locations navigated from blanks, and must hold the same type of value, so
// a binding of mismatched types does not compile:
//
//	Bind(PodSpecableBlank.Template(), dieappsv1.DeploymentBlank.Spec().Template())
func Bind[X any](jig, source base.LocatedAs[X]) Binding { return bind(jig, source) }

// bind binds a jig location to a source location, recording why it can't, if
// it can't, for the jig to report.
func bind(jig, source base.Located) Binding {
	vroot, verr := core.RootOf(jig)
	root, serr := core.RootOf(source)
	return core.WithErr(core.NewBinding(jig.DiePath(), source.DiePath(), vroot, root, nil, nil), errors.Join(verr, serr))
}

// BindPath binds a location of a jig to a path in an unstructured object,
// for [Jig.ForKind] and [Jig.ForUnknownKinds].
func BindPath[X any](jig base.LocatedAs[X], path string) Binding {
	vroot, err := core.RootOf(jig)
	return core.WithErr(core.NewBinding(jig.DiePath(), core.ParsePath(path), vroot, nil, nil, nil), err)
}

// BindPathVia binds a location of a jig to a path in an unstructured object
// through a transform.
func BindPathVia[S, V any](jig base.LocatedAs[V], path string, t Transform[S, V]) Binding {
	b := BindPath[V](jig, path)
	b = core.WithTransform(b, newTransform(t))
	return b
}

// Transform converts between a source value S and a jig value V. It must be
// lossless both ways:
//
//   - get-put: writing back what was read changes nothing,
//     From(To(s), s) == s
//   - put-get: what was written reads back unchanged, To(From(v, s)) == v
//
// [Transform.From] receives the current source value, so that what the jig
// does not show is kept. Put-get is checked on every write through a
// transform: a jig value the source cannot represent is an error, rather than
// being changed. Check get-put with [CheckTransform].
type Transform[S, V any] struct {
	// To presents a source value as a jig value.
	To func(s S) V
	// From updates a source value from a jig value.
	From func(v V, s S) S
}

// CheckTransform checks the laws of a transform over sample source values,
// and jig values written to each of them.
func CheckTransform[S, V any](t Transform[S, V], sources []S, jigs []V) error {
	var errs []error
	for _, s := range sources {
		if got := t.From(t.To(s), s); !reflect.DeepEqual(got, s) {
			errs = append(errs, fmt.Errorf("get-put: writing back %+v read from %+v changed it to %+v", t.To(s), s, got))
		}
		for _, v := range jigs {
			if got := t.To(t.From(v, s)); !reflect.DeepEqual(got, v) {
				errs = append(errs, fmt.Errorf("put-get: %+v written to %+v reads back as %+v", v, s, got))
			}
		}
	}
	return errors.Join(errs...)
}

// BindVia binds a location of a jig to a location of a source through a
// transform between their values:
//
//	BindVia(ContainerImageBlank.Image(), diecorev1.ContainerBlank.Image(), imageReference)
func BindVia[S, V any](jig base.LocatedAs[V], source base.LocatedAs[S], t Transform[S, V]) Binding {
	return core.WithTransform(bind(jig, source), newTransform(t))
}

func newTransform[S, V any](t Transform[S, V]) *core.Transform {
	return &core.Transform{
		Stype: reflect.TypeFor[S](),
		To:    func(s reflect.Value) reflect.Value { return reflect.ValueOf(t.To(s.Interface().(S))) },
		From: func(v, s reflect.Value) reflect.Value {
			return reflect.ValueOf(t.From(v.Interface().(V), s.Interface().(S)))
		},
	}
}

// For binds the jig to a source type, given as its blank. It panics if a
// binding is from another type, or couldn't be made.
func (v *Jig[V, VD, VDM]) For(source base.Located, bindings ...Binding) *Jig[V, VD, VDM] {
	return core.Must(v.ForE(source, bindings...))
}

// ForE is [Jig.For], returning an error rather than panicking.
func (v *Jig[V, VD, VDM]) ForE(source base.Located, bindings ...Binding) (*Jig[V, VD, VDM], error) {
	st, err := core.RootOf(source)
	if err != nil {
		return nil, err
	}
	vt := reflect.TypeFor[V]()
	for _, b := range bindings {
		if err := core.BindingErr(b); err != nil {
			return nil, err
		}
		if jig, src, vroot, root := core.BindingInfo(b); vroot != vt || root != st {
			return nil, core.Errorf("die: binding %s to %s is from %s to %s, not %s to %s", jig, src, vroot, root, vt, st)
		}
	}
	v.sources[st] = append(v.sources[st], bindings...)
	return v, nil
}

// ForKind binds the jig to unstructured objects of a kind, by path.
func (v *Jig[V, VD, VDM]) ForKind(gvk schema.GroupVersionKind, bindings ...Binding) *Jig[V, VD, VDM] {
	return core.Must(v.ForKindE(gvk, bindings...))
}

// ForKindE is [Jig.ForKind], returning an error rather than panicking.
func (v *Jig[V, VD, VDM]) ForKindE(gvk schema.GroupVersionKind, bindings ...Binding) (*Jig[V, VD, VDM], error) {
	if err := v.untyped(gvk.String(), bindings); err != nil {
		return nil, err
	}
	v.kinds[gvk] = append(v.kinds[gvk], bindings...)
	return v, nil
}

// ForUnknownKinds binds the jig to unstructured objects of any kind without
// bindings of its own, by path.
func (v *Jig[V, VD, VDM]) ForUnknownKinds(bindings ...Binding) *Jig[V, VD, VDM] {
	return core.Must(v.ForUnknownKindsE(bindings...))
}

// ForUnknownKindsE is [Jig.ForUnknownKinds], returning an error rather than
// panicking.
func (v *Jig[V, VD, VDM]) ForUnknownKindsE(bindings ...Binding) (*Jig[V, VD, VDM], error) {
	if err := v.untyped("unknown kinds", bindings); err != nil {
		return nil, err
	}
	v.unknown = append(v.unknown, bindings...)
	return v, nil
}

func (v *Jig[V, VD, VDM]) untyped(what string, bindings []Binding) error {
	vt := reflect.TypeFor[V]()
	for _, b := range bindings {
		if err := core.BindingErr(b); err != nil {
			return err
		}
		if jig, _, vroot, root := core.BindingInfo(b); vroot != vt || root != nil {
			return core.Errorf("die: binding %s for %s must be from %s to a path", jig, what, vt)
		}
	}
	return nil
}

// Of returns the jig of an immutable die, an immutable die of the jig's
// type, if the jig has bindings for the source's type, or for an
// unstructured die, its kind. A mutable die has [Jig.OfMut].
func (v *Jig[V, VD, VDM]) Of(source base.ImmutableDie) (VD, bool) {
	d, ok, err := v.OfE(source)
	core.Check(err)
	return d, ok
}

// OfE is [Jig.Of], returning an error rather than panicking, such as when an
// unstructured die's kind can't be read.
func (v *Jig[V, VD, VDM]) OfE(source base.ImmutableDie) (VD, bool, error) {
	var zero VD
	src, ok, err := v.sourceOf(source)
	if !ok || err != nil {
		return zero, ok, err
	}
	return dies.WrapDie[V, VD](v.holder(src), nil), true, nil
}

// OfMut returns the jig of a mutable die, a mutable die of the jig's type
// whose writes are writes to the source, if the jig has bindings for the
// source's type, or for an unstructured die, its kind. An immutable die has
// [Jig.Of].
func (v *Jig[V, VD, VDM]) OfMut(source base.MutableDie) (VDM, bool) {
	d, ok, err := v.OfMutE(source)
	core.Check(err)
	return d, ok
}

// OfMutE is [Jig.OfMut], returning an error rather than panicking.
func (v *Jig[V, VD, VDM]) OfMutE(source base.MutableDie) (VDM, bool, error) {
	var zero VDM
	src, ok, err := v.sourceOf(source)
	if !ok || err != nil {
		return zero, ok, err
	}
	return dies.WrapDieMut[V, VD](v.holder(src), nil), true, nil
}

func (v *Jig[V, VD, VDM]) sourceOf(source base.Die) (core.SourceAt, bool, error) {
	h, at, ok := core.StateOf(source)
	if !ok {
		return core.SourceAt{}, false, nil
	}
	if len(at) == 0 && core.UnstructuredValue(h.Base.Type()) {
		ds := core.DieSource{H: h}
		gvk, err := ds.GVK()
		if err != nil {
			return core.SourceAt{}, false, err
		}
		binds, ok := v.forKind(gvk)
		return core.SourceAt{Src: ds, Binds: binds}, ok, nil
	}
	binds, ok := v.sources[core.Deref(core.TypeAt(h.Base.Type(), at))]
	return core.SourceAt{Src: core.HolderSource{H: h}, At: at, Binds: binds}, ok, nil
}

// OfUnstructured returns the jig of an unstructured object, if the jig has
// bindings for its kind: a mutable die, which changes the object in place, as
// an unstructured object has no intent to record, unlike an unstructured die.
func (v *Jig[V, VD, VDM]) OfUnstructured(u *unstructured.Unstructured) (VDM, bool) {
	binds, ok := v.forKind(u.GroupVersionKind())
	if !ok {
		var zero VDM
		return zero, false
	}
	return dies.WrapDieMut[V, VD](v.holder(core.SourceAt{Src: core.UnstructuredSource{Obj: u}, Binds: binds}), nil), true
}

// holder returns the holder of the jig of a source.
func (v *Jig[V, VD, VDM]) holder(src core.SourceAt) *core.Holder {
	var base V
	return &core.Holder{
		Base:    reflect.ValueOf(base),
		Mutable: src.Src.Mutable(),
		View:    &core.Viewing{Src: src.Src, At: src.At, Binds: src.Binds, Vtype: reflect.TypeFor[V]()},
	}
}

func (v *Jig[V, VD, VDM]) forKind(gvk schema.GroupVersionKind) ([]Binding, bool) {
	if binds, ok := v.kinds[gvk]; ok {
		return binds, true
	}
	return v.unknown, v.unknown != nil
}
