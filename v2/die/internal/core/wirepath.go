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
	"reflect"
	"slices"
)

// These map locations in a value, as a Go type has it or as its JSON decodes,
// to the steps of a path, for code that walks values rather than dies, such
// as go-cmp options.

// WireFieldAt returns the wire field of a struct type at a Go field index
// path, such as [0 1] for the APIVersion of an inline TypeMeta. When idx is
// the start of one, as an inline struct's index is, it returns nil and true.
func WireFieldAt(t reflect.Type, idx []int) (fd *FieldDesc, partial bool) {
	for _, f := range Describe(t).Fields {
		if slices.Equal(f.index, idx) {
			return f, false
		}
		if len(f.index) > len(idx) && slices.Equal(f.index[:len(idx)], idx) {
			partial = true
		}
	}
	return nil, partial
}

// WireField returns the wire field of a struct type by its wire name, or nil
// if it has none.
func WireField(t reflect.Type, wire string) *FieldDesc {
	return Describe(t).byWire[wire]
}

// KeyStepOf returns the step of a list-map element, from its value: a struct
// of the element type, or a map[string]any as its JSON decodes. A key it
// doesn't have is the key field's zero value, or default.
func KeyStepOf(elem reflect.Type, keys []string, v reflect.Value) Step {
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	sd := Describe(elem)
	vals := make([]any, len(keys))
	for i, k := range keys {
		fd := sd.Field(k)
		val := reflect.Zero(fd.Typ)
		switch v.Kind() {
		case reflect.Struct:
			val = v.FieldByIndex(fd.index)
		case reflect.Map:
			if e := v.MapIndex(reflect.ValueOf(k)); e.IsValid() && !e.IsNil() {
				val = e.Elem()
			}
		}
		vals[i] = val.Interface()
	}
	return KeyStep(elem, keys, vals)
}

// HasPrefix reports whether a path is at or below another.
func HasPrefix(p, q Path) bool { return p.hasPrefix(q) }

// UnstructuredKeyStepOf returns the step of an element of an unstructured
// list, as its JSON decodes, by the keys of another such step, and whether
// the element has every key. It is how an element is found by a key that
// only the step knows, as the list has no type to say what its keys are.
func UnstructuredKeyStepOf(like Step, v reflect.Value) (Step, bool) {
	if like.kind != stepKey || like.elem != nil {
		return Step{}, false
	}
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if v.Kind() != reflect.Map {
		return Step{}, false
	}
	vals := make([]any, len(like.keys))
	for i, k := range like.keys {
		e := v.MapIndex(reflect.ValueOf(k))
		if !e.IsValid() {
			return Step{}, false
		}
		vals[i] = e.Interface()
	}
	return UnstructuredKeyStep(like.keys, vals), true
}
