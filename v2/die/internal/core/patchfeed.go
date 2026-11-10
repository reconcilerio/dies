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

package core

import (
	"maps"
	"reflect"
	"slices"
)

// Patches are fed as intent, the reverse of encoding them: what a patch says
// becomes intent over the base, so that releasing the same kind of patch
// again sends it back.
//
// An apply patch is a manager's opinion: each field present is set, null
// clears a field, [] and {} are an empty list or map, list-map elements are
// found by their keys, and an atomic struct is set whole.
//
// A merge patch, RFC 7386, is a change to the base: null clears a field,
// objects and maps are merged field by field, and anything else, lists
// included, is set whole.
//
// Fields a typed die's type doesn't have aren't fed. An unstructured die
// follows its mold, if it has one; without, objects merge field by field and
// lists are set whole.

// PatchKind is the kind of patch to feed.
type PatchKind int

const (
	ApplyPatch PatchKind = iota
	MergePatch
)

// FeedPatch replaces the intent of a root die with what a patch, decoded as
// unstructured JSON, says over its base. It returns the holder written, as
// Replace does.
func FeedPatch(h *Holder, kind PatchKind, patch map[string]any) (*Holder, error) {
	f := &patchFeeder{kind: kind}
	n := f.intent(h.Base.Type(), patch, h.Schema)
	if f.err != nil {
		return nil, f.err
	}
	return h.Replace(nil, n)
}

// patchFeeder builds the intent a patch says, and records the first value it
// can't convert, such as a string where the type has a number.
type patchFeeder struct {
	kind PatchKind
	err  error
}

// leaf is the intent that sets a value as it is.
func (f *patchFeeder) leaf(t reflect.Type, raw any) *Onode {
	n, err := setLeaf(t, raw)
	if err != nil && f.err == nil {
		f.err = err
	}
	return n
}

func clearNode() *Onode { return SetOp(OpClear, reflect.Value{})(&Onode{}) }
func emptyNode() *Onode { return SetOp(OpEmpty, reflect.Value{})(&Onode{}) }

// intent is the intent a patch's value says for a value of type t, with the
// unstructured schema sc, if any.
func (f *patchFeeder) intent(t reflect.Type, raw any, sc *Schema) *Onode {
	if raw == nil {
		return clearNode()
	}
	// an unstructured die's object, or a value within it
	if t := Deref(t); t.Kind() == reflect.Interface || t == reflect.TypeFor[map[string]any]() {
		return f.unstructured(raw, sc)
	}
	m, ok := raw.(map[string]any)
	if !ok || !structured(t) {
		return f.leaf(t, raw)
	}
	n := &Onode{}
	sd := Describe(Deref(t))
	for _, k := range slices.Sorted(maps.Keys(m)) {
		fd, ok := sd.byWire[k]
		if !ok {
			continue
		}
		n.setKid(FieldStep(k), f.field(fd, m[k]))
	}
	return n
}

func (f *patchFeeder) field(fd *FieldDesc, raw any) *Onode {
	if raw == nil {
		return clearNode()
	}
	switch fd.Kind {
	case KindStruct:
		if f.kind == ApplyPatch && fd.Atomic {
			// an apply sends an atomic struct whole
			return f.leaf(fd.Typ, raw)
		}
		return f.intent(fd.Typ, raw, nil)
	case KindListMap:
		list, ok := raw.([]any)
		if !ok || f.kind == MergePatch {
			// a merge patch replaces a list
			return f.leaf(fd.Typ, raw)
		}
		if len(list) == 0 {
			return emptyNode()
		}
		elem := Deref(fd.Typ).Elem()
		sd := Describe(elem)
		n := &Onode{}
		for _, e := range list {
			m, _ := e.(map[string]any)
			vals := make([]any, len(fd.Keys))
			for i, k := range fd.Keys {
				kv := reflect.Zero(sd.Field(k).Typ)
				if v, ok := m[k]; ok && v != nil {
					var err error
					if kv, err = fromUnstructured(v, kv.Type()); err != nil {
						if f.err == nil {
							f.err = err
						}
						return n
					}
				}
				vals[i] = kv.Interface()
			}
			en := f.intent(elem, e, nil)
			keyed(en, fd.Keys, vals)
			n.setKid(KeyStep(elem, fd.Keys, vals), en)
		}
		return n
	case KindMap:
		m, ok := raw.(map[string]any)
		if !ok {
			return f.leaf(fd.Typ, raw)
		}
		if len(m) == 0 && f.kind == ApplyPatch {
			return emptyNode()
		}
		elem := Deref(fd.Typ).Elem()
		n := &Onode{}
		for _, k := range slices.Sorted(maps.Keys(m)) {
			switch {
			case m[k] == nil:
				n.setKid(MapKeyStep(k), clearNode())
			case structured(elem):
				n.setKid(MapKeyStep(k), f.intent(elem, m[k], nil))
			default:
				n.setKid(MapKeyStep(k), f.leaf(elem, m[k]))
			}
		}
		return n
	}
	if list, ok := raw.([]any); ok && len(list) == 0 && f.kind == ApplyPatch {
		return emptyNode()
	}
	return f.leaf(fd.Typ, raw)
}

// unstructured is the intent for an unstructured value, which follows its
// schema where it has one.
func (f *patchFeeder) unstructured(raw any, sc *Schema) *Onode {
	if sc.atomic() && f.kind == ApplyPatch {
		return f.leaf(reflect.TypeFor[any](), raw)
	}
	switch v := raw.(type) {
	case map[string]any:
		if len(v) == 0 && f.kind == ApplyPatch {
			return emptyNode()
		}
		n := &Onode{}
		object := sc != nil && sc.Kind == SchemaObject
		for _, k := range slices.Sorted(maps.Keys(v)) {
			ksc := sc.elem()
			if object {
				ksc = sc.field(k)
			}
			n.setKid(FieldStep(k), f.unstructured(v[k], ksc))
		}
		return n
	case []any:
		if f.kind == ApplyPatch && len(v) == 0 {
			return emptyNode()
		}
		if f.kind == MergePatch || sc == nil || sc.Kind != SchemaListMap {
			// without a list map's keys, a list is set whole
			return f.leaf(reflect.TypeFor[any](), raw)
		}
		n := &Onode{}
		for _, e := range v {
			m, _ := e.(map[string]any)
			vals := make([]any, len(sc.Keys))
			for i, k := range sc.Keys {
				vals[i] = m[k]
			}
			en := f.unstructured(e, sc.elem())
			keyed(en, sc.Keys, vals)
			n.setKid(UnstructuredKeyStep(sc.Keys, vals), en)
		}
		return n
	case nil:
		return clearNode()
	}
	return f.leaf(reflect.TypeFor[any](), raw)
}

// keyed sets an element's keys, so that it exists however few of its other
// fields are set.
func keyed(n *Onode, keys []string, vals []any) {
	for i, k := range keys {
		if vals[i] != nil && n.kid(FieldStep(k)) == nil {
			n.setKid(FieldStep(k), SetOp(OpSet, reflect.ValueOf(vals[i]))(&Onode{}))
		}
	}
}
