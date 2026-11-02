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

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// source is what a view reads and writes: a typed die, or an unstructured
// object whose types are only known where a view describes them.
type source interface {
	// valueAt reads the value at a path as type t; a typed source may read
	// its own type when t is nil
	valueAt(p Path, t reflect.Type) (reflect.Value, error)
	has(p Path) (bool, error)
	// write writes a value at a path; t is the view's type there, for an
	// untyped source to know what an empty value is
	write(p Path, o Op, val reflect.Value, t reflect.Type) error
	Mutable() bool
	// typeAt is the source's own type at a path, or nil for an untyped source
	typeAt(p Path) reflect.Type
	// listKeys are the keys of the list at a path, nil for an atomic list. An
	// untyped source has no keys of its own and uses the view's.
	listKeys(p Path, view []string) []string
	// intent is the source's recorded intent, at paths within it
	intent() *Onode
	String() string
}

// holderSource is a typed die.
type HolderSource struct{ H *Holder }

func (s HolderSource) valueAt(p Path, t reflect.Type) (reflect.Value, error) {
	v, err := s.H.ValueAt(p)
	if err != nil || t == nil {
		return v, err
	}
	return adapt(v, t), nil
}

func (s HolderSource) has(p Path) (bool, error) { return s.H.Has(p) }

func (s HolderSource) write(p Path, o Op, val reflect.Value, _ reflect.Type) error {
	if o == OpSet {
		val = adapt(val, s.typeAt(p))
	}
	_, err := s.H.Write(nil, p, SetOp(o, val), o != OpClear)
	return err
}

func (s HolderSource) Mutable() bool { return s.H.Mutable }

func (s HolderSource) typeAt(p Path) reflect.Type { return TypeAt(s.H.Base.Type(), p) }

func (s HolderSource) listKeys(p Path, _ []string) []string { return fieldAt(s.H.Base.Type(), p).Keys }

func (s HolderSource) intent() *Onode { return IntentOf(s.H) }

func (s HolderSource) String() string { return s.H.Base.Type().String() }

// unstructuredSource is an unstructured object. Reads convert the part a view
// describes to the view's types; writes change only the leaf written, so
// what the view's types don't describe, such as fields newer than they are,
// is kept. A view changes the object in place: it has no intent to record.
type UnstructuredSource struct{ Obj *unstructured.Unstructured }

func (u UnstructuredSource) String() string { return u.Obj.GroupVersionKind().String() }

func (u UnstructuredSource) Mutable() bool { return true }

func (u UnstructuredSource) intent() *Onode { return nil }

func (u UnstructuredSource) typeAt(Path) reflect.Type { return nil }

func (u UnstructuredSource) listKeys(_ Path, view []string) []string { return view }

func (u UnstructuredSource) get(p Path) (any, bool, error) {
	var cur any = u.Obj.Object
	for _, s := range p {
		switch s.kind {
		case stepKey:
			l, _ := cur.([]any)
			i, err := findUnstructured(l, s)
			if err != nil || i < 0 {
				return nil, false, err
			}
			cur = l[i]
		default:
			m, _ := cur.(map[string]any)
			v, ok := m[s.name]
			if !ok {
				return nil, false, nil
			}
			cur = v
		}
	}
	return cur, true, nil
}

func (u UnstructuredSource) valueAt(p Path, t reflect.Type) (reflect.Value, error) {
	raw, ok, err := u.get(p)
	if err != nil {
		return reflect.Value{}, err
	}
	if !ok || raw == nil {
		if n := len(p); n > 0 && p[n-1].kind == stepKey && p[n-1].elem == t {
			return newElem(t, p[n-1]), nil
		}
		return reflect.Zero(t), nil
	}
	return fromUnstructured(raw, t)
}

func (u UnstructuredSource) has(p Path) (bool, error) {
	raw, ok, err := u.get(p)
	return ok && raw != nil, err
}

func (u UnstructuredSource) write(p Path, o Op, val reflect.Value, t reflect.Type) error {
	if o == OpClear {
		has, err := u.has(p)
		if err != nil || !has {
			return err
		}
		obj, err := update(u.Obj.Object, p, func(any) (any, bool) { return nil, true })
		if err != nil {
			return err
		}
		u.Obj.Object = obj.(map[string]any)
		return nil
	}
	var leaf any
	switch o {
	case OpSet:
		var err error
		if leaf, err = ToUnstructured(val); err != nil {
			return err
		}
	case OpEmpty:
		if t != nil && Deref(t).Kind() == reflect.Slice {
			leaf = []any{}
		} else {
			leaf = map[string]any{}
		}
	}
	obj, err := update(u.Obj.Object, p, func(any) (any, bool) { return leaf, false })
	if err != nil {
		return err
	}
	u.Obj.Object = obj.(map[string]any)
	return nil
}

// update returns cur with the value at p replaced by what fn returns for it,
// or removed. Objects, lists and list elements along the path are created as
// needed, with their key fields set.
func update(cur any, p Path, fn func(old any) (v any, remove bool)) (any, error) {
	s := p[0]
	if s.kind == stepKey {
		l, _ := cur.([]any)
		i, err := findUnstructured(l, s)
		if err != nil {
			return nil, err
		}
		if len(p) == 1 {
			var old any
			if i >= 0 {
				old = l[i]
			}
			v, remove := fn(old)
			switch {
			case remove && i >= 0:
				return append(l[:i:i], l[i+1:]...), nil
			case remove:
				return l, nil
			case i < 0:
				return append(l, v), nil
			}
			l[i] = v
			return l, nil
		}
		if i < 0 {
			l = append(l, keyedElement(s))
			i = len(l) - 1
		}
		if l[i], err = update(l[i], p[1:], fn); err != nil {
			return nil, err
		}
		return l, nil
	}
	m, _ := cur.(map[string]any)
	if m == nil {
		m = map[string]any{}
	}
	if len(p) == 1 {
		v, remove := fn(m[s.name])
		if remove {
			delete(m, s.name)
		} else {
			m[s.name] = v
		}
		return m, nil
	}
	v, err := update(m[s.name], p[1:], fn)
	if err != nil {
		return nil, err
	}
	m[s.name] = v
	return m, nil
}

// findUnstructured finds the element of an unstructured list with a key,
// matching unset key fields by their defaults, as a typed list does. A key
// that is not unique is ambiguous.
func findUnstructured(l []any, s Step) (int, error) {
	found := -1
	for i, el := range l {
		match, err := matchesUnstructured(el, s)
		if err != nil {
			return -1, err
		}
		if match {
			if found >= 0 {
				return -1, Errorf("die: %s is ambiguous: the list has it more than once", s)
			}
			found = i
		}
	}
	return found, nil
}

func matchesUnstructured(el any, s Step) (bool, error) {
	m, ok := el.(map[string]any)
	if !ok {
		return false, nil
	}
	sd := Describe(s.elem)
	for j, k := range s.keys {
		fd := sd.Field(k)
		v := reflect.Zero(fd.Typ)
		if raw, ok := m[k]; ok && raw != nil {
			var err error
			if v, err = fromUnstructured(raw, fd.Typ); err != nil {
				return false, err
			}
		}
		if fd.keyValue(v).Interface() != s.vals[j] {
			return false, nil
		}
	}
	return true, nil
}

func keyedElement(s Step) map[string]any {
	m := map[string]any{}
	for j, k := range s.keys {
		// keys are strings and numbers, which always convert
		v, err := ToUnstructured(reflect.ValueOf(s.vals[j]))
		if err != nil {
			panic(err)
		}
		m[k] = v
	}
	return m
}

func fromUnstructured(raw any, t reflect.Type) (reflect.Value, error) {
	b, err := json.Marshal(raw)
	if err != nil {
		return reflect.Value{}, Wrap(err)
	}
	v := reflect.New(t)
	if err := json.Unmarshal(b, v.Interface()); err != nil {
		return reflect.Value{}, Errorf("die: reading %s from an unstructured object: %w", t, err)
	}
	return v.Elem(), nil
}

// toUnstructured converts a value to the form unstructured objects hold:
// maps, lists, strings, booleans, int64 and float64.
func ToUnstructured(val reflect.Value) (any, error) {
	if !val.IsValid() {
		return nil, nil
	}
	b, err := json.Marshal(val.Interface())
	if err != nil {
		return nil, Wrap(err)
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	var out any
	if err := dec.Decode(&out); err != nil {
		return nil, Wrap(err)
	}
	return normalize(out), nil
}

func normalize(v any) any {
	switch v := v.(type) {
	case map[string]any:
		for k, e := range v {
			v[k] = normalize(e)
		}
	case []any:
		for i, e := range v {
			v[i] = normalize(e)
		}
	case json.Number:
		if i, err := v.Int64(); err == nil {
			return i
		}
		f, _ := v.Float64()
		return f
	}
	return v
}

// dieSource is an unstructured die. Unlike an unstructured object, it records
// what a view writes as intent, so patches carry it.
type DieSource struct{ H *Holder }

func (s DieSource) valueAt(p Path, t reflect.Type) (reflect.Value, error) {
	v, err := s.H.ValueAt(p)
	if err != nil || t == nil {
		return v, err
	}
	if !v.IsValid() || (v.Kind() == reflect.Interface && v.IsNil()) {
		if n := len(p); n > 0 && p[n-1].kind == stepKey && p[n-1].elem == t {
			return newElem(t, p[n-1]), nil
		}
		return reflect.Zero(t), nil
	}
	return fromUnstructured(v.Interface(), t)
}

func (s DieSource) has(p Path) (bool, error) { return s.H.Has(p) }

func (s DieSource) write(p Path, o Op, val reflect.Value, t reflect.Type) error {
	if o == OpEmpty {
		val = reflect.ValueOf(map[string]any{})
		if t != nil && Deref(t).Kind() == reflect.Slice {
			val = reflect.ValueOf([]any{})
		}
	}
	_, err := s.H.Write(nil, p, SetOp(o, val), false)
	return err
}

func (s DieSource) Mutable() bool { return s.H.Mutable }

func (s DieSource) intent() *Onode { return IntentOf(s.H) }

func (s DieSource) typeAt(Path) reflect.Type { return nil }

func (s DieSource) listKeys(_ Path, view []string) []string { return view }

func (s DieSource) String() string {
	gvk, _ := s.GVK()
	return gvk.String()
}

// GVK is the kind of the unstructured die.
func (s DieSource) GVK() (schema.GroupVersionKind, error) {
	apiVersion, err := GetAt[string](s.H, nil, Path{FieldStep("apiVersion")})
	if err != nil {
		return schema.GroupVersionKind{}, err
	}
	kind, err := GetAt[string](s.H, nil, Path{FieldStep("kind")})
	if err != nil {
		return schema.GroupVersionKind{}, err
	}
	return schema.FromAPIVersionAndKind(apiVersion, kind), nil
}

// unstructuredValue reports whether a type is an unstructured object.
func UnstructuredValue(t reflect.Type) bool {
	return t.Kind() == reflect.Map && t.Key().Kind() == reflect.String && t.Elem().Kind() == reflect.Interface
}
