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

// Package die has the types most code that uses dies names: dies, the
// locations within them, and the handles their accessors return. The rest is
// in packages by theme:
//
//   - base: the bases dies embed, and blanks
//   - fields: field handles, and the field tables generated dies declare
//   - lists: list dies, lists of dies as Kubernetes lists are
//   - objects: the interface of object dies for generic code
//   - jigs: dies and unstructured objects held in the shape of other types
//   - unstructured: dies for objects with no Go type
//   - jsonpath: locations within dies found by JSONPath expressions
//   - codegen: what generated dies register, and the checks their tests run
//   - dietest: helpers for tests and fixtures
package die

import (
	"reconciler.io/dies/v2/die/base"
	"reconciler.io/dies/v2/die/fields"
	"reconciler.io/dies/v2/die/lists"
	"reconciler.io/dies/v2/die/objects"
)

// Path is a location in an object's wire schema, such as
// .spec.template.spec.containers[name=app].
type Path = base.Path

// Die is a die of any type, for functions that take any die.
type Die = base.Die

// ImmutableDie is an immutable die of any type, for functions that take only
// immutable dies, such as a jig's Hold.
type ImmutableDie = base.ImmutableDie

// MutableDie is a mutable die of any type, for functions that take only
// mutable dies, such as a jig's HoldMut.
type MutableDie = base.MutableDie

// DieTyped is a die of any type with type metadata: a [Die] that is a
// runtime.Object.
type DieTyped = base.DieTyped

// DieOf carries an immutable die's types for generic code: its value type
// T, its own pointer type D, and its mutable counterpart's, M.
type DieOf[T any, D base.DieOf[T, D, M], M base.DieMutOf[T, D, M]] = base.DieOf[T, D, M]

// DieMutOf carries a mutable die's types for generic code, as [DieOf] does an
// immutable die's.
type DieMutOf[T any, D base.DieOf[T, D, M], M base.DieMutOf[T, D, M]] = base.DieMutOf[T, D, M]

// DieTypedOf carries an immutable typed die's types for generic code, with
// its type metadata's accessors.
type DieTypedOf[T any, D base.DieOf[T, D, M], M base.DieMutOf[T, D, M]] = fields.DieTypedOf[T, D, M]

// DieTypedMutOf carries a mutable typed die's types for generic code, with
// its type metadata's accessors.
type DieTypedMutOf[T any, D base.DieOf[T, D, M], M base.DieMutOf[T, D, M]] = fields.DieTypedMutOf[T, D, M]

// DieObject is an object die of any type: a [DieTyped] that is a
// metav1.Object.
type DieObject = base.DieObject

// DieObjectOf carries an immutable object die's types for generic code, with
// its metadata's accessor. An unstructured die is a [DieObject], but not a
// DieObjectOf: its fields are reached through its nodes, not accessors.
type DieObjectOf[T any, D base.DieOf[T, D, M], M base.DieMutOf[T, D, M]] = objects.DieObjectOf[T, D, M]

// DieObjectMutOf carries a mutable object die's types for generic code, with
// its metadata's accessor and the unstructured feed.
type DieObjectMutOf[T any, D base.DieOf[T, D, M], M base.DieMutOf[T, D, M]] = objects.DieObjectMutOf[T, D, M]

// DieList is a list die of any type, as the Kubernetes list interfaces see a
// list.
type DieList = lists.DieList

// DieListOf is a list die of D items, for generic code.
type DieListOf[D lists.ListItem[D]] = lists.DieListOf[D]

// Located is anything that identifies a location within a root die: a [Path],
// a die or a handle.
type Located = base.Located

// LocatedAs is a location of a value of type X: a die for X, or a handle to
// a leaf of type X.
type LocatedAs[X any] = base.LocatedAs[X]

// Value is a handle to a leaf field of type V of an immutable die.
type Value[V any] = fields.Value[V]

// ValueMut is a handle to a leaf field of type V of a mutable die, which
// writes it.
type ValueMut[V any] = fields.ValueMut[V]

// Optional is a handle to a leaf field held by pointer, of type V, of an
// immutable die.
type Optional[V any] = fields.Optional[V]

// OptionalMut is a handle to a leaf field held by pointer, of type V, of a
// mutable die, which writes it.
type OptionalMut[V any] = fields.OptionalMut[V]

// ListMap is the collection for a list-map field of an immutable die, of C
// elements with CD dies, keyed by K.
type ListMap[C, CD, K any] = fields.ListMap[C, CD, K]

// ListMapMut is the collection for a list-map field of a mutable die, of C
// elements with mutable CDM dies, keyed by K.
type ListMapMut[C, CDM, K any] = fields.ListMapMut[C, CDM, K]

// Map is the collection for a map field of an immutable die, with K keys and
// V values.
type Map[K ~string, V any] = fields.Map[K, V]

// MapMut is the collection for a map field of a mutable die.
type MapMut[K ~string, V any] = fields.MapMut[K, V]

// Error is a die's error. Writes return one, for content a die can't hold, as
// do functions that encode or take content from elsewhere, whose Must forms
// panic with it instead. A mistake in the code panics with one.
type Error = base.Error

var (
	// ErrImmutable is the error of a write to an immutable die through an
	// interface it implements, such as metav1.Object's setters.
	ErrImmutable = base.ErrImmutable
	// ErrImmutableChild is the error of a write through a child of an
	// immutable die.
	ErrImmutableChild = base.ErrImmutableChild
	// ErrReadOnlyJig is the error of a write through a jig of an immutable
	// die.
	ErrReadOnlyJig = base.ErrReadOnlyJig
)
