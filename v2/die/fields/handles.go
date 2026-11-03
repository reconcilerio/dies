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

package fields

import (
	"reflect"

	"reconciler.io/dies/v2/die/base"
	"reconciler.io/dies/v2/die/internal/core"
	"reconciler.io/dies/v2/die/internal/dies"
)

// Every field accessor returns a handle that is acted on: a child die for a
// struct, a collection die for a list map or map, and a Value or Optional for
// a leaf. Handles have no fields of their own, so their methods need no Die
// prefix.
//
// An immutable die's handles read. A mutable die's handles, ValueMut,
// OptionalMut, ListMapMut and MapMut, embed them and also write, in place,
// returning only an error.

// leaf is the state shared by leaf handles: the die's location, and the
// leaf's within it, usually a field of the die.
type leaf struct {
	h   *core.Holder
	at  base.Path
	rel base.Path
}

func (l leaf) path() base.Path { return core.JoinPath(l.at, l.rel) }

func leafOf[T, V any](d base.LocatedAs[T], f Field[T, V]) leaf {
	h, at := dies.MustStateOf(d)
	return leaf{h: h, at: at, rel: base.Path{core.FieldStep(f.fd.Wire)}}
}

// write writes through a leaf.
func write(l leaf, o core.Op, v reflect.Value, union bool) error {
	_, err := l.h.Write(l.at, l.rel, core.SetOp(o, v), union)
	return err
}

// Value is a handle to a leaf field of type V, of an immutable die, which
// reads it.
type Value[V any] struct{ l leaf }

// ValueMut is a handle to a leaf field of type V of a mutable die, which
// writes it.
type ValueMut[V any] struct{ Value[V] }

// DieValue marks a handle as a location of a V, for [base.LocatedAs]. It does
// nothing.
func (Value[V]) DieValue(*V) {}

// LocateValueAt returns the handle for the value of type V at a path within d's
// root, such as one a JSONPath expression found. It returns an error if the
// value at the path isn't a V.
func LocateValueAt[V any](d base.Die, p base.Path) (Value[V], error) {
	l, err := leafAt[V](d, p)
	return Value[V]{l}, err
}

// MustLocateValueAt is [LocateValueAt], panicking rather than returning an error.
func MustLocateValueAt[V any](d base.Die, p base.Path) Value[V] {
	return core.Must(LocateValueAt[V](d, p))
}

// LocateValueMutAt returns the handle for the value of type V at a path within a
// mutable die's root. The path may be outside d itself, such as the namespace
// of the object a subject is part of. It panics if the value at the path isn't
// a V.
func LocateValueMutAt[V any](d base.MutableDie, p base.Path) (ValueMut[V], error) {
	l, err := leafAt[V](d, p)
	return ValueMut[V]{Value[V]{l}}, err
}

// MustLocateValueMutAt is [LocateValueMutAt], panicking rather than returning an error.
func MustLocateValueMutAt[V any](d base.MutableDie, p base.Path) ValueMut[V] {
	return core.Must(LocateValueMutAt[V](d, p))
}

func leafAt[V any](d base.Die, p base.Path) (leaf, error) {
	h, at := dies.MustStateOf(d)
	if t := core.TypeAt(h.Base.Type(), p); t != reflect.TypeFor[V]() {
		return leaf{}, core.Errorf("die: the value at %s is a %s, not a %s", p, t, reflect.TypeFor[V]())
	}
	if len(p) >= len(at) && p[:len(at)].String() == at.String() {
		return leaf{h: h, at: at, rel: append(base.Path{}, p[len(at):]...)}, nil
	}
	// outside d: written from the root
	return leaf{h: h, rel: append(base.Path{}, p...)}, nil
}

// ValueOf returns the handle for a leaf field of an immutable die.
func ValueOf[T, V any](d base.LocatedAs[T], f Field[T, V]) Value[V] {
	return Value[V]{leafOf(d, f)}
}

// ValueMutOf returns the handle for a leaf field of a mutable die.
func ValueMutOf[T, V any](d base.LocatedAs[T], f Field[T, V]) ValueMut[V] {
	return ValueMut[V]{Value[V]{leafOf(d, f)}}
}

// Get returns a copy of the current value.
func (v Value[V]) Get() V {
	val := core.Must(v.l.h.ValueAt(v.l.path()))
	return val.Interface().(V)
}

// Has reports whether the field is present.
func (v Value[V]) Has() bool {
	return core.Must(v.l.h.Has(v.l.path()))
}

func (v Value[V]) DiePath() base.Path { return v.l.path() }

// DieRoot returns the root die this is within: see [base.DieBase.DieRoot].
func (v Value[V]) DieRoot() base.Die { return base.RootOf(v.l.h) }

// DieRootType returns the type of the root value this is within.
func (v Value[V]) DieRootType() reflect.Type { return v.l.h.Base.Type() }

// Set records the value as intent, even when it is the zero value.
func (v ValueMut[V]) Set(x V) error {
	return write(v.l, core.OpSet, core.DeepCopy(reflect.ValueOf(&x).Elem()), true)
}

// Clear records that the field is absent: null in patches.
func (v ValueMut[V]) Clear() error {
	return write(v.l, core.OpClear, reflect.Value{}, false)
}

// Empty records that a list or map field is present and empty.
func (v ValueMut[V]) Empty() error {
	return write(v.l, core.OpEmpty, reflect.Value{}, true)
}

// Optional is a handle to a pointer leaf field of type *V of an immutable
// die, which reads V, with presence reported separately.
type Optional[V any] struct{ l leaf }

// OptionalMut is a handle to a pointer leaf field of type *V of a mutable
// die, which writes V.
type OptionalMut[V any] struct{ Optional[V] }

// OptionalOf returns the handle for a pointer leaf field of an immutable die.
func OptionalOf[T, V any](d base.LocatedAs[T], f Field[T, *V]) Optional[V] {
	return Optional[V]{leafOf(d, f)}
}

// OptionalMutOf returns the handle for a pointer leaf field of a mutable die.
func OptionalMutOf[T, V any](d base.LocatedAs[T], f Field[T, *V]) OptionalMut[V] {
	return OptionalMut[V]{Optional[V]{leafOf(d, f)}}
}

// Get returns the current value, or the zero value when absent.
func (o Optional[V]) Get() V {
	if p := o.Ptr(); p != nil {
		return *p
	}
	var zero V
	return zero
}

// Ptr returns a copy of the current value, or nil when absent.
func (o Optional[V]) Ptr() *V {
	val := core.Must(o.l.h.ValueAt(o.l.path()))
	return val.Interface().(*V)
}

// Has reports whether the field is present.
func (o Optional[V]) Has() bool {
	return core.Must(o.l.h.Has(o.l.path()))
}

func (o Optional[V]) DiePath() base.Path { return o.l.path() }

// DieRoot returns the root die this is within: see [base.DieBase.DieRoot].
func (o Optional[V]) DieRoot() base.Die { return base.RootOf(o.l.h) }

// DieRootType returns the type of the root value this is within.
func (o Optional[V]) DieRootType() reflect.Type { return o.l.h.Base.Type() }

// Set records the value as intent.
func (o OptionalMut[V]) Set(x V) error {
	return write(o.l, core.OpSet, reflect.ValueOf(&x), true)
}

// Clear records that the field is absent: null in patches.
func (o OptionalMut[V]) Clear() error {
	return write(o.l, core.OpClear, reflect.Value{}, false)
}
