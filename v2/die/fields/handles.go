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
// OptionalMut, ListMapMut and MapMut, embed them and also write, in place: a
// leaf's writes return the die it is a field of, so calls chain at the die's
// level.

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

// write writes through a leaf of d, returning d.
func write[M any](d M, l leaf, o core.Op, v reflect.Value, union bool) (M, error) {
	if _, err := l.h.Write(l.at, l.rel, core.SetOp(o, v), union); err != nil {
		var zero M
		return zero, err
	}
	return d, nil
}

// Value is a handle to a leaf field of type V, of an immutable die, which
// reads it.
type Value[V any] struct{ l leaf }

// ValueMut is a handle to a leaf field of type V of a mutable M die, which
// writes it.
type ValueMut[V, M any] struct {
	Value[V]
	d M
}

// DieValue marks a handle as a location of a V, for [base.LocatedAs]. It does
// nothing.
func (Value[V]) DieValue(*V) {}

// ValueAt returns the handle for the value of type V at a path within d's
// root, such as one a JSONPath expression found. It panics if the value at
// the path isn't a V.
func ValueAt[V any](d base.Die, p base.Path) Value[V] { return core.Must(ValueAtE[V](d, p)) }

// ValueAtE is [ValueAt], returning an error rather than panicking.
func ValueAtE[V any](d base.Die, p base.Path) (Value[V], error) {
	l, err := leafAt[V](d, p)
	return Value[V]{l}, err
}

// ValueMutAt returns the handle for the value of type V at a path within a
// mutable die's root, whose writes return d. The path may be outside d
// itself, such as the namespace of the object a subject is part of. It
// panics if the value at the path isn't a V.
func ValueMutAt[V any, M base.MutableDie](d M, p base.Path) ValueMut[V, M] {
	return core.Must(ValueMutAtE[V](d, p))
}

// ValueMutAtE is [ValueMutAt], returning an error rather than panicking.
func ValueMutAtE[V any, M base.MutableDie](d M, p base.Path) (ValueMut[V, M], error) {
	l, err := leafAt[V](d, p)
	return ValueMut[V, M]{Value[V]{l}, d}, err
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
func ValueMutOf[T, V any, M base.LocatedAs[T]](d M, f Field[T, V]) ValueMut[V, M] {
	return ValueMut[V, M]{Value[V]{leafOf[T](d, f)}, d}
}

// Get returns a copy of the current value.
func (v Value[V]) Get() V { return core.Must(v.GetE()) }

// GetE is [Value.Get], returning an error rather than panicking.
func (v Value[V]) GetE() (V, error) {
	val, err := v.l.h.ValueAt(v.l.path())
	if err != nil {
		var zero V
		return zero, err
	}
	return val.Interface().(V), nil
}

// Has reports whether the field is present.
func (v Value[V]) Has() bool { return core.Must(v.HasE()) }

// HasE is [Value.Has], returning an error rather than panicking.
func (v Value[V]) HasE() (bool, error) { return v.l.h.Has(v.l.path()) }

func (v Value[V]) DiePath() base.Path { return v.l.path() }

// DieRoot returns the root die this is within: see [base.DieBase.DieRoot].
func (v Value[V]) DieRoot() base.Die { return base.RootOf(v.l.h) }

// DieRootType returns the type of the root value this is within.
func (v Value[V]) DieRootType() reflect.Type { return v.l.h.Base.Type() }

// Set records the value as intent, even when it is the zero value.
func (v ValueMut[V, M]) Set(x V) M { return core.Must(v.SetE(x)) }

// SetE is [ValueMut.Set], returning an error rather than panicking.
func (v ValueMut[V, M]) SetE(x V) (M, error) {
	return write(v.d, v.l, core.OpSet, core.DeepCopy(reflect.ValueOf(&x).Elem()), true)
}

// Clear records that the field is absent: null in patches.
func (v ValueMut[V, M]) Clear() M { return core.Must(v.ClearE()) }

// ClearE is [ValueMut.Clear], returning an error rather than panicking.
func (v ValueMut[V, M]) ClearE() (M, error) {
	return write(v.d, v.l, core.OpClear, reflect.Value{}, false)
}

// Empty records that a list or map field is present and empty.
func (v ValueMut[V, M]) Empty() M { return core.Must(v.EmptyE()) }

// EmptyE is [ValueMut.Empty], returning an error rather than panicking.
func (v ValueMut[V, M]) EmptyE() (M, error) {
	return write(v.d, v.l, core.OpEmpty, reflect.Value{}, true)
}

// Optional is a handle to a pointer leaf field of type *V of an immutable
// die, which reads V, with presence reported separately.
type Optional[V any] struct{ l leaf }

// OptionalMut is a handle to a pointer leaf field of type *V of a mutable M
// die, which writes V.
type OptionalMut[V, M any] struct {
	Optional[V]
	d M
}

// OptionalOf returns the handle for a pointer leaf field of an immutable die.
func OptionalOf[T, V any](d base.LocatedAs[T], f Field[T, *V]) Optional[V] {
	return Optional[V]{leafOf(d, f)}
}

// OptionalMutOf returns the handle for a pointer leaf field of a mutable die.
func OptionalMutOf[T, V any, M base.LocatedAs[T]](d M, f Field[T, *V]) OptionalMut[V, M] {
	return OptionalMut[V, M]{Optional[V]{leafOf[T](d, f)}, d}
}

// Get returns the current value, or the zero value when absent.
func (o Optional[V]) Get() V { return core.Must(o.GetE()) }

// GetE is [Optional.Get], returning an error rather than panicking.
func (o Optional[V]) GetE() (V, error) {
	p, err := o.PtrE()
	if err != nil || p == nil {
		var zero V
		return zero, err
	}
	return *p, nil
}

// Ptr returns a copy of the current value, or nil when absent.
func (o Optional[V]) Ptr() *V { return core.Must(o.PtrE()) }

// PtrE is [Optional.Ptr], returning an error rather than panicking.
func (o Optional[V]) PtrE() (*V, error) {
	val, err := o.l.h.ValueAt(o.l.path())
	if err != nil {
		return nil, err
	}
	return val.Interface().(*V), nil
}

// Has reports whether the field is present.
func (o Optional[V]) Has() bool { return core.Must(o.HasE()) }

// HasE is [Optional.Has], returning an error rather than panicking.
func (o Optional[V]) HasE() (bool, error) { return o.l.h.Has(o.l.path()) }

func (o Optional[V]) DiePath() base.Path { return o.l.path() }

// DieRoot returns the root die this is within: see [base.DieBase.DieRoot].
func (o Optional[V]) DieRoot() base.Die { return base.RootOf(o.l.h) }

// DieRootType returns the type of the root value this is within.
func (o Optional[V]) DieRootType() reflect.Type { return o.l.h.Base.Type() }

// Set records the value as intent.
func (o OptionalMut[V, M]) Set(x V) M { return core.Must(o.SetE(x)) }

// SetE is [OptionalMut.Set], returning an error rather than panicking.
func (o OptionalMut[V, M]) SetE(x V) (M, error) {
	return write(o.d, o.l, core.OpSet, reflect.ValueOf(&x), true)
}

// Clear records that the field is absent: null in patches.
func (o OptionalMut[V, M]) Clear() M { return core.Must(o.ClearE()) }

// ClearE is [OptionalMut.Clear], returning an error rather than panicking.
func (o OptionalMut[V, M]) ClearE() (M, error) {
	return write(o.d, o.l, core.OpClear, reflect.Value{}, false)
}
