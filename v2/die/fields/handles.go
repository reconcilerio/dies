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
// A leaf's mutators return the parent die so that calls chain at the parent's
// level. On an immutable root they return a new root.

// leaf is the state shared by leaf handles: the die's location, and the
// leaf's within it, usually a field of the die.
type leaf[D any] struct {
	h    *core.Holder
	at   base.Path
	rel  base.Path
	self func(*core.Holder) (D, error)
}

func (l leaf[D]) path() base.Path { return core.JoinPath(l.at, l.rel) }

func (l leaf[D]) write(o core.Op, v reflect.Value, union bool) (D, error) {
	h, err := l.h.Write(l.at, l.rel, core.SetOp(o, v), union)
	if err != nil {
		var zero D
		return zero, err
	}
	return l.self(h)
}

// selfOf is the die a write through a handle of d returns.
func selfOf[T any, D base.DieOf[T, D]](d D) func(*core.Holder) (D, error) {
	self := dies.SelfOf[T](d)
	return func(h *core.Holder) (D, error) { return self(h), nil }
}

// Value is a handle to a leaf field of type V on a D die.
type Value[V, D any] struct{ l leaf[D] }

// DieValue marks a handle as a location of a V, for LocatedAs. It does
// nothing.
func (Value[V, D]) DieValue(*V) {}

// ValueAt returns the handle for the value of type V at a path within d's
// root, such as one a JSONPath expression found, whose writes return d. The
// path may be outside d itself, such as the namespace of the object a
// subject is part of; then, as for any write through a child of an immutable
// root, writing through it fails unless the root is mutable. It panics if
// the value at the path isn't a V.
func ValueAt[V any, D base.Die](d D, p base.Path) Value[V, D] { return core.Must(ValueAtE[V](d, p)) }

// ValueAtE is ValueAt, returning an error rather than panicking.
func ValueAtE[V any, D base.Die](d D, p base.Path) (Value[V, D], error) {
	h, at := dies.MustStateOf(d)
	if t := core.TypeAt(h.Base.Type(), p); t != reflect.TypeFor[V]() {
		return Value[V, D]{}, core.Errorf("die: the value at %s is a %s, not a %s", p, t, reflect.TypeFor[V]())
	}
	if len(p) >= len(at) && p[:len(at)].String() == at.String() {
		self := func(nh *core.Holder) (D, error) {
			if nh != h {
				return dies.Wrap[D](nh, nil), nil
			}
			return dies.Wrap[D](nh, at), nil
		}
		return Value[V, D]{leaf[D]{h: h, at: at, rel: append(base.Path{}, p[len(at):]...), self: self}}, nil
	}
	// outside d: written from the root, returning d
	self := func(nh *core.Holder) (D, error) {
		if nh != h {
			// a write on an immutable root made a new one, which d is not in
			var zero D
			return zero, core.ErrImmutableChild
		}
		return dies.Wrap[D](nh, at), nil
	}
	return Value[V, D]{leaf[D]{h: h, rel: append(base.Path{}, p...), self: self}}, nil
}

// ValueOf returns the handle for a leaf field of d.
func ValueOf[T any, D base.DieOf[T, D], V any](d D, f Field[T, V]) Value[V, D] {
	h, at := dies.MustStateOf(d)
	return Value[V, D]{leaf[D]{h: h, at: at, rel: base.Path{core.FieldStep(f.fd.Wire)}, self: selfOf[T](d)}}
}

// Get returns a copy of the current value.
func (v Value[V, D]) Get() V { return core.Must(v.GetE()) }

// GetE is Get, returning an error rather than panicking.
func (v Value[V, D]) GetE() (V, error) {
	val, err := v.l.h.ValueAt(v.l.path())
	if err != nil {
		var zero V
		return zero, err
	}
	return val.Interface().(V), nil
}

// Has reports whether the field is present.
func (v Value[V, D]) Has() bool { return core.Must(v.HasE()) }

// HasE is Has, returning an error rather than panicking.
func (v Value[V, D]) HasE() (bool, error) { return v.l.h.Has(v.l.path()) }

// Set records the value as intent, even when it is the zero value.
func (v Value[V, D]) Set(x V) D { return core.Must(v.SetE(x)) }

// SetE is Set, returning an error rather than panicking.
func (v Value[V, D]) SetE(x V) (D, error) {
	return v.l.write(core.OpSet, core.DeepCopy(reflect.ValueOf(&x).Elem()), true)
}

// Clear records that the field is absent: null in patches.
func (v Value[V, D]) Clear() D { return core.Must(v.ClearE()) }

// ClearE is Clear, returning an error rather than panicking.
func (v Value[V, D]) ClearE() (D, error) { return v.l.write(core.OpClear, reflect.Value{}, false) }

// Empty records that a list or map field is present and empty.
func (v Value[V, D]) Empty() D { return core.Must(v.EmptyE()) }

// EmptyE is Empty, returning an error rather than panicking.
func (v Value[V, D]) EmptyE() (D, error) { return v.l.write(core.OpEmpty, reflect.Value{}, true) }

func (v Value[V, D]) DiePath() base.Path { return v.l.path() }

// DieRoot returns the root die this is within: see base.DieBase.DieRoot.
func (v Value[V, D]) DieRoot() base.Die { return base.RootOf(v.l.h) }

// DieRootType returns the type of the root value this is within.
func (v Value[V, D]) DieRootType() reflect.Type { return v.l.h.Base.Type() }

// Optional is a handle to a pointer leaf field of type *V on a D die. It reads
// and writes V, with presence reported separately.
type Optional[V, D any] struct{ l leaf[D] }

// OptionalOf returns the handle for a pointer leaf field of d.
func OptionalOf[T any, D base.DieOf[T, D], V any](d D, f Field[T, *V]) Optional[V, D] {
	h, at := dies.MustStateOf(d)
	return Optional[V, D]{leaf[D]{h: h, at: at, rel: base.Path{core.FieldStep(f.fd.Wire)}, self: selfOf[T](d)}}
}

// Get returns the current value, or the zero value when absent.
func (o Optional[V, D]) Get() V { return core.Must(o.GetE()) }

// GetE is Get, returning an error rather than panicking.
func (o Optional[V, D]) GetE() (V, error) {
	p, err := o.PtrE()
	if err != nil || p == nil {
		var zero V
		return zero, err
	}
	return *p, nil
}

// Ptr returns a copy of the current value, or nil when absent.
func (o Optional[V, D]) Ptr() *V { return core.Must(o.PtrE()) }

// PtrE is Ptr, returning an error rather than panicking.
func (o Optional[V, D]) PtrE() (*V, error) {
	val, err := o.l.h.ValueAt(o.l.path())
	if err != nil {
		return nil, err
	}
	return val.Interface().(*V), nil
}

// Has reports whether the field is present.
func (o Optional[V, D]) Has() bool { return core.Must(o.HasE()) }

// HasE is Has, returning an error rather than panicking.
func (o Optional[V, D]) HasE() (bool, error) { return o.l.h.Has(o.l.path()) }

// Set records the value as intent.
func (o Optional[V, D]) Set(x V) D { return core.Must(o.SetE(x)) }

// SetE is Set, returning an error rather than panicking.
func (o Optional[V, D]) SetE(x V) (D, error) {
	return o.l.write(core.OpSet, reflect.ValueOf(&x), true)
}

// Clear records that the field is absent: null in patches.
func (o Optional[V, D]) Clear() D { return core.Must(o.ClearE()) }

// ClearE is Clear, returning an error rather than panicking.
func (o Optional[V, D]) ClearE() (D, error) { return o.l.write(core.OpClear, reflect.Value{}, false) }

func (o Optional[V, D]) DiePath() base.Path { return o.l.path() }

// DieRoot returns the root die this is within: see base.DieBase.DieRoot.
func (o Optional[V, D]) DieRoot() base.Die { return base.RootOf(o.l.h) }

// DieRootType returns the type of the root value this is within.
func (o Optional[V, D]) DieRootType() reflect.Type { return o.l.h.Base.Type() }
