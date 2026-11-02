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
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
)

// Patches are encoded from the overlay, never from the Go value: an untouched
// field is omitted even when its Go type would encode it (non-omitempty
// fields encode nil as null), and a touched zero value is sent even when the
// field is omitempty.

// applyPatch encodes the intent as a server-side apply configuration:
//
//   - untouched fields are omitted (no opinion)
//   - Clear on a field sends null (owned and absent), where the server
//     accepts null: for every field of a built-in type, whose server converts
//     it through its Go type, and for a nullable field of any other, such as
//     a custom resource, whose validation rejects null otherwise. Elsewhere
//     it's omitted, and the server removes it if this manager owned it.
//   - Empty sends [] or {} (owned and empty)
//   - removing a list-map element or a map entry omits it: an apply can't
//     remove what another manager owns, and the server removes what this
//     manager owned and no longer applies (null in a map sets an entry to
//     its zero value, rather than removing it)
//   - a list or map with only removals, or a struct with nothing else to
//     send, is omitted, so that the manager doesn't own it empty
//   - an atomic struct with intent is sent whole, as the server replaces it
//
// An unstructured die with a schema follows the schema as a typed die
// follows its type: Clear on an object's field sends null, and atomic
// values, lists included, are sent whole. Where its schema is unknown, an
// object is a map, and a list addressed by keys is a list map.
func (h *Holder) ApplyPatch() ([]byte, error) {
	v, err := encodeApply(h, h.Base.Type(), h.Ov, nil, h.Schema)
	if err != nil {
		return nil, err
	}
	obj, _ := v.(map[string]any)
	if obj == nil {
		obj = map[string]any{}
	}
	// an apply request always identifies its object, even when its value
	// doesn't, as when a client decoded it
	if id, ok := identityOf(h.Base.Type()); ok {
		obj["apiVersion"], obj["kind"] = id.apiVersion, id.kind
	}
	current, err := valueJSON(h, nil)
	if err != nil {
		return nil, err
	}
	if cur, ok := current.(map[string]any); ok {
		for _, k := range []string{"apiVersion", "kind"} {
			if v, ok := cur[k]; ok && v != "" {
				obj[k] = v
			}
		}
		if md, ok := cur["metadata"].(map[string]any); ok {
			omd, _ := obj["metadata"].(map[string]any)
			if omd == nil {
				omd = map[string]any{}
			}
			for _, k := range []string{"name", "namespace"} {
				if v, ok := md[k]; ok {
					omd[k] = v
				}
			}
			if len(omd) > 0 {
				obj["metadata"] = omd
			}
		}
	}
	return json.Marshal(obj)
}

// acceptsNull reports whether the server takes null in an apply as absent
// for every field of the die's object: a server that converts the object
// through Go types, as for built-in types, while a custom resource's
// validation rejects null for a field that isn't nullable. A typed die's
// object is taken to be served from Go types if its type is from k8s.io; an
// unstructured die's, if its schema says so.
func (h *Holder) acceptsNull() bool {
	if h.Schema != nil {
		return h.Schema.GoTypes
	}
	return strings.HasPrefix(h.Base.Type().PkgPath(), "k8s.io/")
}

// omitted is what encodeApply returns for intent that an apply leaves out.
var omitted any = new(struct{ omitted bool })

func encodeApply(h *Holder, t reflect.Type, n *Onode, p Path, sc *Schema) (any, error) {
	if n == nil {
		return omitted, nil
	}
	switch {
	case n.removes():
		return nil, nil
	case n.Op == OpEmpty && len(n.Order) == 0 && n.Val.IsValid():
		return toJSON(n.Val)
	case n.Op == OpEmpty && len(n.Order) == 0:
		return emptyJSON(t), nil
	case n.replaces():
		return h.EncodeAt(p)
	case sc.atomic():
		// intent within a value the server replaces whole
		return h.EncodeAt(p)
	}
	t = shapeOf(Deref(t), n, sc)
	switch t.Kind() {
	case reflect.Struct:
		obj := map[string]any{}
		sd := Describe(t)
		for _, id := range n.Order {
			s := n.Steps[id]
			fd := sd.Field(s.name)
			k := n.Kids[id]
			if k.removes() && !fd.Nullable && !h.acceptsNull() {
				continue
			}
			var v any
			var err error
			if fd.Atomic && !k.removes() {
				v, err = h.EncodeAt(p.child(s))
			} else {
				v, err = encodeApply(h, fd.Typ, k, p.child(s), nil)
			}
			if err != nil {
				return nil, err
			}
			if v != omitted {
				obj[fd.Wire] = v
			}
		}
		if len(obj) == 0 && len(p) > 0 {
			return omitted, nil
		}
		return obj, nil
	case reflect.Slice:
		arr := []any{}
		for _, id := range n.Order {
			k, s := n.Kids[id], n.Steps[id]
			if k.removes() {
				continue
			}
			e, err := encodeApply(h, t.Elem(), k, p.child(s), sc.elem())
			if err != nil {
				return nil, err
			}
			if e == omitted {
				// an element is identified by its keys, even with nothing
				// else to send
				e = map[string]any{}
			}
			if m, ok := e.(map[string]any); ok {
				// elements are identified by their keys
				for j, key := range s.keys {
					m[key] = s.vals[j]
				}
			}
			arr = append(arr, e)
		}
		if len(arr) == 0 {
			return omitted, nil
		}
		return arr, nil
	case reflect.Map:
		obj := map[string]any{}
		// an unstructured object whose schema is known has fields, as a
		// struct does, rather than entries
		object := sc != nil && sc.Kind == SchemaObject
		for _, id := range n.Order {
			k, s := n.Kids[id], n.Steps[id]
			ksc := sc.elem()
			if object {
				ksc = sc.field(s.name)
			}
			if k.removes() {
				if object && (h.acceptsNull() || (ksc != nil && ksc.Nullable)) {
					obj[s.name] = nil
				}
				continue
			}
			v, err := encodeApply(h, t.Elem(), k, p.child(s), ksc)
			if err != nil {
				return nil, err
			}
			if v != omitted {
				obj[s.name] = v
			}
		}
		if len(obj) == 0 && (!object || len(p) > 0) {
			return omitted, nil
		}
		return obj, nil
	}
	return valueJSON(h, p)
}

// valueJSON is the value at p in its generic JSON form.
func valueJSON(h *Holder, p Path) (any, error) {
	v, err := h.ValueAt(p)
	if err != nil {
		return nil, err
	}
	return toJSON(v)
}

// mergePatch encodes the intent as an RFC 7386 JSON merge patch against base.
// Lists are always sent whole, which is all a merge patch can express.
func (h *Holder) MergePatch() ([]byte, error) {
	base := &Holder{Base: h.Base, Wire: h.Wire}
	v, ok, err := encodeMerge(h, base, h.Base.Type(), h.Ov, nil)
	if err != nil {
		return nil, err
	}
	if !ok {
		return []byte("{}"), nil
	}
	return json.Marshal(v)
}

func encodeMerge(h, base *Holder, t reflect.Type, n *Onode, p Path) (any, bool, error) {
	if n == nil {
		return nil, false, nil
	}
	if n.removes() {
		has, err := base.Has(p)
		return nil, has, err
	}
	if n.Op != opNone {
		a, err := base.EncodeAt(p)
		if err != nil {
			return nil, false, err
		}
		b, err := h.EncodeAt(p)
		if err != nil {
			return nil, false, err
		}
		v, ok := mergeDiff(a, b)
		return v, ok, nil
	}
	t = shapeOf(Deref(t), n, nil)
	switch t.Kind() {
	case reflect.Struct, reflect.Map:
		obj := map[string]any{}
		for _, id := range n.Order {
			s := n.Steps[id]
			var kt reflect.Type
			if t.Kind() == reflect.Struct {
				kt = Describe(t).Field(s.name).Typ
			} else {
				kt = t.Elem()
			}
			v, ok, err := encodeMerge(h, base, kt, n.Kids[id], p.child(s))
			if err != nil {
				return nil, false, err
			}
			if ok {
				obj[s.name] = v
			}
		}
		return obj, len(obj) > 0, nil
	}
	a, err := valueJSON(base, p)
	if err != nil {
		return nil, false, err
	}
	b, err := valueJSON(h, p)
	if err != nil {
		return nil, false, err
	}
	v, ok := mergeDiff(a, b)
	return v, ok, nil
}

// shapeOf is the type to encode an unstructured value's intent as: a list
// when it is addressed by keys, an object otherwise.
func shapeOf(t reflect.Type, n *Onode, sc *Schema) reflect.Type {
	if t.Kind() != reflect.Interface {
		return t
	}
	switch {
	case sc == nil:
	case sc.Kind == SchemaListMap:
		return reflect.TypeFor[[]any]()
	case sc.Kind == SchemaObject || sc.Kind == SchemaMap:
		return reflect.TypeFor[map[string]any]()
	}
	for _, s := range n.Steps {
		if s.kind == stepKey {
			return reflect.TypeFor[[]any]()
		}
	}
	return reflect.TypeFor[map[string]any]()
}

// mergeDiff returns the merge patch that turns a into b.
func mergeDiff(a, b any) (any, bool) {
	am, aok := a.(map[string]any)
	bm, bok := b.(map[string]any)
	if aok && bok {
		obj := map[string]any{}
		for k, bv := range bm {
			if av, ok := am[k]; ok {
				if d, changed := mergeDiff(av, bv); changed {
					obj[k] = d
				}
			} else {
				obj[k] = bv
			}
		}
		for k := range am {
			if _, ok := bm[k]; !ok {
				obj[k] = nil
			}
		}
		return obj, len(obj) > 0
	}
	if reflect.DeepEqual(a, b) {
		return nil, false
	}
	return b, true
}

func emptyJSON(t reflect.Type) any {
	if Deref(t).Kind() == reflect.Slice {
		return []any{}
	}
	return map[string]any{}
}

// toJSON converts a value to its generic JSON form. A value that's encoded
// field by field goes through encodeAt instead, for its presence.
func toJSON(v reflect.Value) (any, error) {
	b, err := json.Marshal(v.Interface())
	if err != nil {
		return nil, Wrap(err)
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	var out any
	if err := dec.Decode(&out); err != nil {
		return nil, Wrap(err)
	}
	return out, nil
}
