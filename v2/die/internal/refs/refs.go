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

// Package refs does the work of a die's methods, and of the functions that
// take a die, with functions that aren't generic, so that it is compiled once
// rather than for each type of die: a generic type's methods are compiled for
// each struct type they're instantiated with.
package refs

import (
	"encoding/json"
	"reflect"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	utiljson "k8s.io/apimachinery/pkg/util/json"
	"sigs.k8s.io/yaml"

	"reconciler.io/dies/v2/die/internal/core"
)

// Ref is a die's state: its holder, its location within it, and the type of
// its value.
type Ref struct {
	H  *core.Holder
	At core.Path
	T  reflect.Type
}

// Of returns the state of a die of any type.
func Of(d any) (Ref, bool) {
	h, at, ok := core.StateOf(d)
	if !ok {
		return Ref{}, false
	}
	return Ref{H: h, At: at, T: core.ValueTypeOf(d)}, true
}

// NewHolder is the holder of a zero die, mutable, for a decoder.
func NewHolder(t reflect.Type) *core.Holder {
	return &core.Holder{Base: reflect.Zero(t), Mutable: true}
}

// Root returns the holder of a root die that isn't a jig, for the methods
// only it supports.
func (r Ref) Root(method string) (*core.Holder, error) {
	if len(r.At) > 0 {
		return nil, core.Errorf("die: %s is only supported on a root die, not at %s", method, r.At)
	}
	if r.H.View != nil {
		return nil, core.Errorf("die: %s is not supported on a jig; use it on the jig's source", method)
	}
	return r.H, nil
}

// CopyOf returns a new root holder with a copy of the die: a root's value
// and intent, or the current value of a die within another or presented by a
// jig, without intent, so the copy is unlinked from its source.
func (r Ref) CopyOf(mutable bool) (*core.Holder, error) {
	h := r.H
	if len(r.At) > 0 || h.View != nil {
		v, err := r.Value()
		if err != nil {
			return nil, err
		}
		return &core.Holder{Base: v, Mutable: mutable, Schema: h.Schema.At(r.At)}, nil
	}
	return &core.Holder{Base: h.Base, Ov: h.Ov, Wire: h.Wire, Mutable: mutable, Schema: h.Schema}, nil
}

// Has reports whether the value is present. A root always is.
func (r Ref) Has() (bool, error) {
	if len(r.At) == 0 {
		return true, nil
	}
	return r.H.Has(r.At)
}

// Value returns a copy of the current value, as a value of the die's type.
func (r Ref) Value() (reflect.Value, error) {
	v, err := r.H.ValueAt(r.At)
	if err != nil {
		return reflect.Value{}, err
	}
	// a value within an unstructured one is held as an any
	if v.Kind() == reflect.Interface && v.Type() != r.T {
		if v.IsNil() {
			return reflect.Zero(r.T), nil
		}
		v = v.Elem()
	}
	if v.Kind() == reflect.Pointer && v.Type().Elem() == r.T {
		if v.IsNil() {
			return reflect.Zero(r.T), nil
		}
		v = v.Elem()
	}
	return v, nil
}

// MarshalJSON encodes the current value as its presence says.
func (r Ref) MarshalJSON() ([]byte, error) {
	v, err := r.H.EncodeAt(r.At)
	if err != nil {
		return nil, err
	}
	return json.Marshal(v)
}

func (r Ref) ReleaseJSON() ([]byte, error) {
	data, err := r.MarshalJSON()
	if err != nil {
		return nil, core.Errorf("die: encoding %s: %w", r.T, err)
	}
	return data, nil
}

func (r Ref) ReleaseYAML() ([]byte, error) {
	j, err := r.ReleaseJSON()
	if err != nil {
		return nil, err
	}
	y, err := yaml.JSONToYAML(j)
	if err != nil {
		return nil, core.Errorf("die: encoding %s as YAML: %w", r.T, err)
	}
	return y, nil
}

func (r Ref) ReleaseRawExtension() (runtime.RawExtension, error) {
	data, err := r.ReleaseJSON()
	return runtime.RawExtension{Raw: data}, err
}

// ReleaseUnstructured releases the value as an unstructured object: an
// unstructured die's content as it is, and another die's as it encodes.
func (r Ref) ReleaseUnstructured() (*unstructured.Unstructured, error) {
	if r.T == reflect.TypeFor[map[string]any]() {
		v, err := r.Value()
		if err != nil {
			return nil, err
		}
		return &unstructured.Unstructured{Object: v.Interface().(map[string]any)}, nil
	}
	data, err := r.ReleaseJSON()
	if err != nil {
		return nil, err
	}
	u := &unstructured.Unstructured{Object: map[string]any{}}
	// numbers decode as int64 where they can, as unstructured objects hold them
	if err := utiljson.Unmarshal(data, &u.Object); err != nil {
		return nil, core.Errorf("die: decoding %s as an unstructured object: %w", r.T, err)
	}
	return u, nil
}

func (r Ref) ReleaseApplyPatch() ([]byte, error) {
	h, err := r.Root("ReleaseApplyPatch")
	if err != nil {
		return nil, err
	}
	patch, err := h.ApplyPatch()
	if err != nil {
		return nil, core.Errorf("die: encoding the intent of %s as an apply: %w", r.T, err)
	}
	return patch, nil
}

func (r Ref) ReleaseMergePatch() ([]byte, error) {
	h, err := r.Root("ReleaseMergePatch")
	if err != nil {
		return nil, err
	}
	patch, err := h.MergePatch()
	if err != nil {
		return nil, core.Errorf("die: encoding the intent of %s as a merge patch: %w", r.T, err)
	}
	return patch, nil
}

// Tare returns a new holder of the die's value with its intent folded in.
func (r Ref) Tare() (*core.Holder, error) {
	h, err := r.Root("DieTare")
	if err != nil {
		return nil, err
	}
	v, err := h.Materialize()
	if err != nil {
		return nil, err
	}
	return &core.Holder{Base: v, Wire: core.CommitWire(h.Wire, h.Ov), Mutable: h.Mutable, Schema: h.Schema}, nil
}

// Replace replaces the die's state with a holder's, in place.
func (r Ref) Replace(h *core.Holder, err error) error {
	if err != nil {
		return err
	}
	r.H.Base, r.H.Ov, r.H.Wire = h.Base, h.Ov, h.Wire
	return nil
}

// Clear records that the value is absent, or resets a root to the zero value.
func (r Ref) Clear() error {
	if len(r.At) == 0 {
		return r.Feed(reflect.Zero(r.T))
	}
	_, err := r.H.Write(r.At, nil, core.SetOp(core.OpClear, reflect.Value{}), false)
	return err
}

// Feed replaces the value: of a root, its base, discarding its intent; of a
// die within another, as intent.
func (r Ref) Feed(v reflect.Value) error {
	v = core.DeepCopy(v)
	if len(r.At) == 0 && r.H.View == nil {
		r.H.Base, r.H.Ov, r.H.Wire = v, nil, nil
		return nil
	}
	_, err := r.H.Write(r.At, nil, core.SetOp(core.OpSet, fit(core.TypeAt(r.H.Base.Type(), r.At), v)), true)
	return err
}

// fit adapts a value to the static type at its location, which may be a
// pointer to it.
func fit(t reflect.Type, v reflect.Value) reflect.Value {
	if t != v.Type() && t.Kind() == reflect.Pointer && t.Elem() == v.Type() {
		p := reflect.New(v.Type())
		p.Elem().Set(v)
		return p
	}
	return v
}

// decode decodes JSON as a value of the die's type. Numbers in unstructured
// values are int64 where they can be, as apimachinery holds them; structs
// decode as encoding/json decodes them.
func (r Ref) decode(data []byte) (reflect.Value, error) {
	p := reflect.New(r.T)
	if err := utiljson.Unmarshal(data, p.Interface()); err != nil {
		return reflect.Value{}, err
	}
	return p.Elem(), nil
}

func (r Ref) FeedJSON(data []byte) error {
	v, err := r.decode(data)
	if err != nil {
		return core.Errorf("die: feeding %s from JSON: %w", r.T, err)
	}
	if len(r.At) > 0 || r.H.View != nil {
		return r.Feed(v)
	}
	w, err := core.DecodeWire(r.T, data)
	if err != nil {
		return core.Errorf("die: feeding %s from JSON: %w", r.T, err)
	}
	r.H.Base, r.H.Ov, r.H.Wire = v, nil, w
	return nil
}

func (r Ref) FeedYAML(data []byte) error {
	j, err := yaml.YAMLToJSON(data)
	if err != nil {
		return core.Errorf("die: feeding %s from YAML: %w", r.T, err)
	}
	return r.FeedJSON(j)
}

func (r Ref) FeedRawExtension(raw runtime.RawExtension) error {
	data, err := raw.MarshalJSON()
	if err != nil {
		return core.Errorf("die: feeding %s from a raw extension: %w", r.T, err)
	}
	return r.FeedJSON(data)
}

// FeedUnstructured feeds an unstructured object: an unstructured die copies
// its content, so its numbers stay int64 where they are, and another die
// decodes it.
func (r Ref) FeedUnstructured(u *unstructured.Unstructured) error {
	if r.T == reflect.TypeFor[map[string]any]() {
		return r.Feed(reflect.ValueOf(u.DeepCopy().Object))
	}
	data, err := u.MarshalJSON()
	if err != nil {
		return core.Errorf("die: encoding an unstructured object: %w", err)
	}
	return r.FeedJSON(data)
}

// Stamp passes a pointer to a copy of the value to fn, and records what fn
// changed as intent.
func (r Ref) Stamp(fn func(p any)) error {
	before, err := r.Value()
	if err != nil {
		return err
	}
	p := reflect.New(r.T)
	p.Elem().Set(core.DeepCopy(before))
	fn(p.Interface())
	after := p.Elem()
	if core.FirstIndex(r.At) >= 0 {
		// within a list element found by its index, which a write sets whole
		if reflect.DeepEqual(before.Interface(), after.Interface()) {
			return nil
		}
		_, err := r.H.Write(r.At, nil, core.SetOp(core.OpSet, after), false)
		return err
	}
	ov, err := core.DiffInto(r.H.Ov, r.At, r.T, before, after)
	if err != nil {
		return err
	}
	_, err = r.H.Replace(r.At, ov)
	return err
}

// FeedPatch feeds a patch as intent over the value.
func (r Ref) FeedPatch(method string, kind core.PatchKind, patch []byte) error {
	h, err := r.Root(method)
	if err != nil {
		return err
	}
	var raw map[string]any
	// numbers are int64 where they can be, as unstructured values hold them
	if err := utiljson.Unmarshal(patch, &raw); err != nil {
		return core.Errorf("die: %s: %w", method, err)
	}
	_, err = core.FeedPatch(h, kind, raw)
	return err
}

// Unmarshal replaces the base and keeps the intent, so the intent is
// reapplied on top of the decoded value.
func (r Ref) Unmarshal(data []byte) error {
	if len(r.At) > 0 || r.H.View != nil {
		return core.Errorf("die: can only unmarshal into a root die")
	}
	v, err := r.decode(data)
	if err != nil {
		return err
	}
	w, err := core.DecodeWire(r.T, data)
	if err != nil {
		return err
	}
	r.H.Base, r.H.Wire = v, w
	return nil
}

// Keep returns a new holder of a root's value keeping the intent at the
// locations, or all but it, with the rest folded into the base.
func Keep(h *core.Holder, at []core.Located, only bool) (*core.Holder, error) {
	paths, err := pathsOf(h.Base.Type(), at)
	if err != nil {
		return nil, err
	}
	kept, fold := core.Only(h.Ov, paths), core.Except(h.Ov, paths)
	if !only {
		kept, fold = fold, kept
	}
	base, err := (&core.Holder{Base: h.Base, Ov: fold}).Materialize()
	if err != nil {
		return nil, err
	}
	return &core.Holder{Base: base, Ov: kept, Wire: core.CommitWire(h.Wire, fold), Mutable: h.Mutable, Schema: h.Schema}, nil
}

func pathsOf(root reflect.Type, at []core.Located) ([]core.Path, error) {
	paths := make([]core.Path, len(at))
	for i, l := range at {
		if r, ok := l.(core.Rooted); ok && r.DieRootType() != root {
			return nil, core.Errorf("die: location %s is within %s, not %s", l.DiePath(), r.DieRootType(), root)
		}
		paths[i] = l.DiePath()
	}
	return paths, nil
}

// Subresources returns the locations of a type's fields served through
// subresources that match, such as status.
func Subresources(t reflect.Type, match func(string) bool) []core.Located {
	var paths []core.Located
	for _, fd := range core.Describe(t).Fields {
		if match(fd.Subresource) {
			paths = append(paths, core.Path{core.FieldStep(fd.Wire)})
		}
	}
	return paths
}
