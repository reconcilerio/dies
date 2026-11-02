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

// Package fields has the handles a die's accessors return, Value, Optional,
// ListMap and Map, and the field tables generated dies declare to make them.
package fields

import (
	"reflect"

	"reconciler.io/dies/v2/die/base"
	"reconciler.io/dies/v2/die/internal/core"
	"reconciler.io/dies/v2/die/internal/dies"
)

// Typed field references. A generator would emit one per field; they are
// checked against the schema when created, so a mismatch fails at init.

// Field is a leaf field of T with value type V.
type Field[T, V any] struct{ fd *core.FieldDesc }

func NewField[T, V any](wire string) Field[T, V] {
	return Field[T, V]{fd: lookup[T](wire, reflect.TypeFor[V]())}
}

// Struct is a struct field of T, held by value or pointer, whose die is CD.
type Struct[T, C any, CD base.DieOf[C, CD]] struct{ fd *core.FieldDesc }

func NewStruct[T, C any, CD base.DieOf[C, CD]](wire string) Struct[T, C, CD] {
	fd := lookup[T](wire, nil)
	if core.Deref(fd.Typ) != reflect.TypeFor[C]() {
		panic(core.Errorf("die: field %q of %s is %s, not %s", wire, reflect.TypeFor[T](), fd.Typ, reflect.TypeFor[C]()))
	}
	return Struct[T, C, CD]{fd: fd}
}

// ListMapField is a list-map field of T with elements C keyed by K.
type ListMapField[T, C any, CD base.DieOf[C, CD], K any] struct {
	fd    *core.FieldDesc
	key   func(K) []any
	keyOf func(C) K
}

// NewListMap declares a list-map field. key returns the key values in the
// order of the field's listMap keys; keyOf returns an element's key.
func NewListMap[T, C any, CD base.DieOf[C, CD], K any](wire string, key func(K) []any, keyOf func(C) K) ListMapField[T, C, CD, K] {
	// the list may be a named type, such as ValidationRules []ValidationRule
	fd := lookup[T](wire, nil)
	if fd.Typ.Kind() != reflect.Slice || fd.Typ.Elem() != reflect.TypeFor[C]() {
		panic(core.Errorf("die: field %q of %s is %s, not a list of %s", wire, reflect.TypeFor[T](), fd.Typ, reflect.TypeFor[C]()))
	}
	if fd.Kind != core.KindListMap {
		panic(core.Errorf("die: field %q of %s is not a list map", wire, reflect.TypeFor[T]()))
	}
	return ListMapField[T, C, CD, K]{fd: fd, key: key, keyOf: keyOf}
}

// MapField is a granular map field of T with string-kinded keys K and values V.
type MapField[T any, K ~string, V any] struct{ fd *core.FieldDesc }

func NewMap[T any, K ~string, V any](wire string) MapField[T, K, V] {
	fd := lookup[T](wire, nil)
	if fd.Kind != core.KindMap || fd.Typ.Key() != reflect.TypeFor[K]() || fd.Typ.Elem() != reflect.TypeFor[V]() {
		panic(core.Errorf("die: field %q of %s is %s, not a map of %s", wire, reflect.TypeFor[T](), fd.Typ, reflect.TypeFor[V]()))
	}
	return MapField[T, K, V]{fd: fd}
}

// lookup finds a field for a field table, recording that T's die manages it.
func lookup[T any](wire string, want reflect.Type) *core.FieldDesc {
	t := reflect.TypeFor[T]()
	fd := core.Describe(t).Field(wire)
	if want != nil && fd.Typ != want {
		panic(core.Errorf("die: field %q of %s is %s, not %s", wire, t, fd.Typ, want))
	}
	core.Manage(t, fd)
	return fd
}

// Child returns the die for a struct field, attached to d.
func Child[T any, D base.DieOf[T, D], C any, CD base.DieOf[C, CD]](d D, f Struct[T, C, CD]) CD {
	h, at := dies.MustStateOf(d)
	return dies.WrapDie[C, CD](h, core.ChildPath(at, core.FieldStep(f.fd.Wire)))
}

// ListMapOf returns the collection die for a list-map field, attached to d.
func ListMapOf[T any, D base.DieOf[T, D], C any, CD base.DieOf[C, CD], K any](d D, f ListMapField[T, C, CD, K]) *ListMap[C, CD, K] {
	h, at := dies.MustStateOf(d)
	return &ListMap[C, CD, K]{h: h, at: core.ChildPath(at, core.FieldStep(f.fd.Wire)), keys: f.fd.Keys, key: f.key, keyOf: f.keyOf}
}

// MapOf returns the collection die for a map field, attached to d.
func MapOf[T any, D base.DieOf[T, D], K ~string, V any](d D, f MapField[T, K, V]) *Map[K, V] {
	h, at := dies.MustStateOf(d)
	return &Map[K, V]{h: h, at: core.ChildPath(at, core.FieldStep(f.fd.Wire))}
}
