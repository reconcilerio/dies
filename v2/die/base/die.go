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

// Package base has the bases dies embed, which give every die its methods,
// and blanks.
//
// A die is an immutable base value plus a persistent overlay of intent. Writes
// record intent (set, clear, empty) by wire path and copy only the overlay
// nodes along that path; the base is never copied. Reads apply the overlay
// along the path being read. Patches are encoded from the overlay, so they
// contain exactly what was stamped.
//
// Every die type comes in two: an immutable die, such as DeploymentDie, which
// only reads, and a mutable die, such as DeploymentDieMut, which also writes,
// in place. Blanks are immutable; [DieBase.DieStamp] passes a mutable copy to
// functions and returns an immutable die of the result. [DieBase.DieMutable]
// and [DieMutBase.DieImmutable] convert, always by copying, so an immutable
// die never sees a write.
//
// A child die returned by an accessor is attached to its root: it addresses
// into the root by path, reads never create anything, and the first write
// creates whatever the path needs. The children of an immutable die are
// immutable, and those of a mutable die mutable.
package base

import (
	"encoding/json"
	"fmt"
	"reflect"

	"reconciler.io/dies/v2/die/internal/core"
	"reconciler.io/dies/v2/die/internal/refs"
)

// dieOf is what immutable and mutable dies of a T share: reading.
type dieOf[T any] interface {
	Die
	LocatedAs[T]

	DieRelease() *T
	DieReleaseE() (*T, error)
}

// DieOf carries an immutable die's types for generic code: T is the die's
// value type, D its own pointer type, so its methods return the concrete
// die, and M the pointer type of its mutable counterpart. It has every method
// of an immutable die. Generic code constrains a die by it,
// D DieOf[T, D, M]; code that takes a die of any type takes a [Die].
type DieOf[T any, D DieOf[T, D, M], M DieMutOf[T, D, M]] interface {
	dieOf[T]
	ImmutableDie
	// DieFrom wraps state in a new die. It is called on a nil receiver.
	DieFrom(DieBase[T, D, M]) D

	DieStamp(fns ...func(m M)) D
	DieStampE(fns ...func(m M) error) (D, error)
	DieMutable() M
	DieMutableE() (M, error)
	DeepCopy() D
	DeepCopyE() (D, error)
	DieTare() D
	DieTareE() (D, error)
	Equal(o D) bool
	EqualE(o D) (bool, error)
}

// DieMutOf carries a mutable die's types for generic code, as [DieOf] does an
// immutable die's. It has every method of a mutable die.
type DieMutOf[T any, D DieOf[T, D, M], M DieMutOf[T, D, M]] interface {
	dieOf[T]
	MutableDie
	json.Unmarshaler
	// DieFrom wraps state in a new die. It is called on a nil receiver.
	DieFrom(DieMutBase[T, D, M]) M

	DieClear() error
	DieFeed(r *T) error
	DieStamp(fns ...func(m M)) M
	DieStampE(fns ...func(m M) error) (M, error)
	DieImmutable() D
	DieImmutableE() (D, error)
	DeepCopy() M
	DeepCopyE() (M, error)
	DieTare() error
	Equal(o M) bool
	EqualE(o M) (bool, error)
}

// Die is a die of any type, immutable or mutable: the methods every die has
// that don't name its types, for functions that take any die. Handles,
// collections and nodes have some of them, but aren't dies. [DieOf] and
// [DieMutOf] have the rest, for generic code.
type Die interface {
	// MarshalJSON encodes the die's value as its presence says.
	json.Marshaler

	DiePath() Path
	DieRootType() reflect.Type
	DieRoot() Die
	DieIntent() string
	DieFields() []FieldInfo
	DieHas() bool
	DieHasE() (bool, error)
}

// ImmutableDie is an immutable die of any type, for functions that take only
// immutable dies, such as jigs.Jig.Of. Only the dies of [DieBase], and the
// typed and object bases that embed it, are ImmutableDies.
type ImmutableDie interface {
	Die
	dieImmutableMarker()
}

// MutableDie is a mutable die of any type, for functions that take only
// mutable dies, such as jigs.Jig.OfMut. Only the dies of [DieMutBase], and the
// typed and object bases that embed it, are MutableDies.
type MutableDie interface {
	Die
	dieMutableMarker()
}

// IsMutable reports whether a die is a mutable die.
func IsMutable(d any) bool {
	_, ok := d.(MutableDie)
	return ok
}

// state is the state of a die, immutable or mutable, with the methods both
// have, which read. [DieBase] and [DieMutBase] embed it, rather than one the
// other, so that each has only its own methods: [DieBase.DieMutable] only on
// an immutable die, and [DieMutBase.DieImmutable] only on a mutable one.
type state[T any] struct {
	h  *core.Holder
	at Path
}

// DieBase holds the state of an immutable die and provides its methods,
// which read. Immutable dies embed it by value.
type DieBase[T any, D DieOf[T, D, M], M DieMutOf[T, D, M]] struct {
	state[T]
}

// DieMutBase holds the state of a mutable die and provides its methods. Its
// writes change the die in place, and return only an error. Mutable dies
// embed it by value.
//
// The zero value is a mutable, empty die, so a die allocated by a decoder can
// be unmarshaled into.
type DieMutBase[T any, D DieOf[T, D, M], M DieMutOf[T, D, M]] struct {
	state[T]
}

func wrap[T any, D DieOf[T, D, M], M DieMutOf[T, D, M]](h *core.Holder, at Path) D {
	if len(at) == 0 {
		h.SetRoot(func(h *core.Holder) any { return wrap[T, D, M](h, nil) })
	}
	var d D
	return d.DieFrom(DieBase[T, D, M]{state[T]{h: h, at: at}})
}

func wrapMut[T any, D DieOf[T, D, M], M DieMutOf[T, D, M]](h *core.Holder, at Path) M {
	if len(at) == 0 {
		h.SetRoot(func(h *core.Holder) any { return wrapMut[T, D, M](h, nil) })
	}
	var m M
	return m.DieFrom(DieMutBase[T, D, M]{state[T]{h: h, at: at}})
}

// BlankOf returns an immutable die of a value, such as an object with only
// its type metadata, for a blank.
func BlankOf[T any, D DieOf[T, D, M], M DieMutOf[T, D, M]](v T) D {
	return wrap[T, D, M](&core.Holder{Base: reflect.ValueOf(v)}, nil)
}

// Blank returns an immutable, empty die.
func Blank[T any, D DieOf[T, D, M], M DieMutOf[T, D, M]]() D {
	return wrap[T, D, M](&core.Holder{Base: reflect.Zero(reflect.TypeFor[T]())}, nil)
}

// DieValue marks a die as a location of a T, for [LocatedAs]. It does nothing.
func (*state[T]) DieValue(*T) {}

// dieState is how the packages of handles, jigs and unstructured dies reach
// a die's state, through the functions installed below.
type dieState interface {
	hold() *core.Holder
	dieAt() Path
	setDieState(h *core.Holder, at Path)
	dieValueType() reflect.Type
}

func (b *state[T]) dieAt() Path { return b.at }

func (b *state[T]) setDieState(h *core.Holder, at Path) { b.h, b.at = h, at }

func (b *state[T]) dieValueType() reflect.Type { return reflect.TypeFor[T]() }

func (*DieBase[T, D, M]) dieImmutableMarker() {}

func (*DieMutBase[T, D, M]) dieMutableMarker() {}

func init() {
	core.StateOf = func(d any) (*core.Holder, Path, bool) {
		s, ok := d.(dieState)
		if !ok {
			return nil, nil, false
		}
		return s.hold(), s.dieAt(), true
	}
	core.SetState = func(d any, h *core.Holder, at Path) {
		if len(at) == 0 {
			// the die of the root, of whatever type it is
			t := reflect.TypeOf(d).Elem()
			h.SetRoot(func(h *core.Holder) any {
				root := reflect.New(t).Interface()
				root.(dieState).setDieState(h, nil)
				return root
			})
		}
		d.(dieState).setDieState(h, at)
	}
	core.ValueTypeOf = func(d any) reflect.Type { return d.(dieState).dieValueType() }
	core.IsMutable = IsMutable
}

// hold returns the die's holder. A zero die's is made on first use, mutable,
// for a decoder; an immutable die never writes to it.
func (b *state[T]) hold() *core.Holder {
	if b.h == nil {
		b.h = refs.NewHolder(reflect.TypeFor[T]())
	}
	return b.h
}

// ref returns the die's state for the functions behind its methods.
func (b *state[T]) ref() refs.Ref { return refs.Ref{H: b.hold(), At: b.at, T: reflect.TypeFor[T]()} }

// DiePath returns the location of this die within its root.
func (b *state[T]) DiePath() Path { return b.at }

// DieRoot returns the root die this die is within, sharing its state. A root
// returns itself. Its type is the root's die type, such as *DeploymentDie
// for a die within an immutable Deployment, or *DeploymentDieMut within a
// mutable one, whose writes are writes to this die's tree. Assert it as a
// type, or as metav1.Object for an object's identity whatever its type:
// d.DieRoot().(metav1.Object).GetNamespace().
func (b *DieBase[T, D, M]) DieRoot() Die {
	if len(b.at) == 0 {
		return wrap[T, D, M](b.hold(), nil)
	}
	return RootOf(b.hold())
}

// DieRoot returns the root die this die is within: see [DieBase.DieRoot].
func (b *DieMutBase[T, D, M]) DieRoot() Die {
	if len(b.at) == 0 {
		return wrapMut[T, D, M](b.hold(), nil)
	}
	return RootOf(b.hold())
}

// RootOf returns the root die over a holder, or nil if none was made.
func RootOf(h *core.Holder) Die {
	root, _ := h.Root().(Die)
	return root
}

// DieRootType returns the type of the root value this die is within.
func (b *state[T]) DieRootType() reflect.Type { return b.hold().Base.Type() }

// DieHas reports whether the value is present. A root is always present.
func (b *state[T]) DieHas() bool { return core.Must(b.ref().Has()) }

// DieHasE is [DieBase.DieHas], returning an error rather than panicking.
func (b *state[T]) DieHasE() (bool, error) { return b.ref().Has() }

// DieRelease returns a pointer to a copy of the current value.
func (b *state[T]) DieRelease() *T { return core.Must(b.DieReleaseE()) }

// DieReleaseE is [DieBase.DieRelease], returning an error rather than
// panicking.
func (b *state[T]) DieReleaseE() (*T, error) {
	v, err := b.ref().Value()
	if err != nil {
		return nil, err
	}
	r := v.Interface().(T)
	return &r, nil
}

// DieStamp passes a mutable copy of the die to each fn, and returns an
// immutable die of the result. The copy of a root has its intent; the copy
// of a die within another, or presented by a jig, is a new root holding its
// current value, as [DieBase.DeepCopy] makes it.
func (b *DieBase[T, D, M]) DieStamp(fns ...func(m M)) D {
	return core.Must(b.DieStampE(withE(fns)...))
}

// DieStampE is [DieBase.DieStamp] for functions that can fail, such as a client
// reading into a die: it passes a mutable copy to each fn, as
// [DieBase.DieStamp] does, until one returns an error, and returns that error
// with an immutable die of the copy, as the functions left it.
//
//	current, err := DeploymentBlank.DieStampE(func(d *DeploymentDieMut) error {
//		return c.Get(ctx, key, d)
//	})
func (b *DieBase[T, D, M]) DieStampE(fns ...func(m M) error) (D, error) {
	h, err := b.ref().CopyOf(true)
	if err != nil {
		var zero D
		return zero, err
	}
	m := wrapMut[T, D, M](h, nil)
	err = runWith(m, fns)
	return wrap[T, D, M](with(h, m), nil), err
}

func withE[M any](fns []func(m M)) []func(M) error {
	efns := make([]func(M) error, len(fns))
	for i, fn := range fns {
		efns[i] = func(m M) error { fn(m); return nil }
	}
	return efns
}

func runWith[M any](m M, fns []func(M) error) error {
	for _, fn := range fns {
		if err := fn(m); err != nil {
			return err
		}
	}
	return nil
}

// held returns the holder of a die passed to functions, restoring the schema
// of the holder it was given if its state was replaced, as a client zeroes
// the object it decodes into.
func held(d any, was *core.Holder) *core.Holder {
	m, _, _ := core.StateOf(d)
	if m != was && m.Schema == nil {
		m.Schema = was.Schema
	}
	return m
}

// DieMutable returns a mutable copy of the die: of a root, its value and
// intent; of a die within another, or presented by a jig, a new root holding
// its current value, as [DieBase.DeepCopy] makes it. Writes to the copy never
// reach the die.
func (b *DieBase[T, D, M]) DieMutable() M { return core.Must(b.DieMutableE()) }

// DieMutableE is [DieBase.DieMutable], returning an error rather than
// panicking.
func (b *DieBase[T, D, M]) DieMutableE() (M, error) {
	h, err := b.ref().CopyOf(true)
	return wrapped(h, err, wrapMut[T, D, M])
}

// DeepCopy returns an independent copy. The base and overlay are shared, as
// neither is ever mutated, so this is cheap. A child die, or a die presented
// by another, such as a jig or a list item, is copied as a new root holding
// its current value, without intent: the copy is unlinked from its source.
func (b *DieBase[T, D, M]) DeepCopy() D { return core.Must(b.DeepCopyE()) }

// DeepCopyE is [DieBase.DeepCopy], returning an error rather than panicking.
func (b *DieBase[T, D, M]) DeepCopyE() (D, error) {
	h, err := b.ref().CopyOf(false)
	return wrapped(h, err, wrap[T, D, M])
}

// DieTare returns the die with its intent folded into its base: its current
// value becomes the value later writes are measured from, as taring a scale
// makes what is on it the new zero. Tare after a successful write, or before
// changes that should be intent of their own.
func (b *DieBase[T, D, M]) DieTare() D { return core.Must(b.DieTareE()) }

// DieTareE is [DieBase.DieTare], returning an error rather than panicking.
func (b *DieBase[T, D, M]) DieTareE() (D, error) {
	h, err := b.ref().Tare()
	return wrapped(h, err, wrap[T, D, M])
}

func wrapped[X any](h *core.Holder, err error, wrap func(*core.Holder, Path) X) (X, error) {
	if err != nil {
		var zero X
		return zero, err
	}
	return wrap(h, nil), nil
}

// DieIntent describes the recorded intent, one operation per line. A jig
// records none of its own, as its writes go to its sources: its intent is
// theirs, through its bindings, operation for operation where a location is
// bound as is, and as a set of the whole value where it is bound through a
// transform or an aggregate.
func (b *state[T]) DieIntent() string { return intent(b.ref()) }

func describeIntent(n *core.Onode, p Path, lines *[]string) {
	if n == nil {
		return
	}
	switch n.Op {
	case core.OpSet:
		// an unstructured value set to nil has no value: null
		v := []byte("null")
		if n.Val.IsValid() {
			v, _ = json.Marshal(n.Val.Interface())
		}
		*lines = append(*lines, fmt.Sprintf("set %s = %s", p, v))
	case core.OpClear:
		*lines = append(*lines, fmt.Sprintf("clear %s", p))
	case core.OpEmpty:
		*lines = append(*lines, fmt.Sprintf("empty %s", p))
	}
	for _, id := range n.Order {
		describeIntent(n.Kids[id], core.ChildPath(p, n.Steps[id]), lines)
	}
}

// Equal reports whether two dies have the same [Content]: the same current
// value, as encoded, fields its type doesn't have included, and the same
// intent. Values that encode alike, such as times that differ only within a
// second, are equal.
func (b *DieBase[T, D, M]) Equal(o D) bool { return core.Must(equal(b.ref(), o)) }

// EqualE is [DieBase.Equal], returning an error rather than panicking.
func (b *DieBase[T, D, M]) EqualE(o D) (bool, error) { return equal(b.ref(), o) }

// MarshalJSON encodes the current value as its presence says: values that
// were decoded as they were on the wire, fields its type doesn't have
// included, and written values as Go has them, nil and omitempty zero values
// omitted.
func (b *state[T]) MarshalJSON() ([]byte, error) { return b.ref().MarshalJSON() }

// UnmarshalJSON fails with ErrImmutable: an immutable die can't be decoded
// into. Decode into its mutable die, such as within [DieBase.DieStampE].
func (b *DieBase[T, D, M]) UnmarshalJSON([]byte) error {
	return core.Errorf("die: decoding into an immutable %s die: %w", reflect.TypeFor[T](), core.ErrImmutable)
}

// The methods of mutable dies: those that write, which change the die in
// place and return only an error, and those of DieBase that return a die,
// returning a mutable one.

// self returns the mutable die at this die's location.
func (b *DieMutBase[T, D, M]) self() M { return wrapMut[T, D, M](b.hold(), b.at) }

// DieClear records that the value is absent: null in patches, or a removed
// element. On a root it resets to the zero value.
func (b *DieMutBase[T, D, M]) DieClear() error { return b.ref().Clear() }

// DieFeed replaces the value, with a copy of what r points to, or the zero
// value for nil. On a root die the value becomes the new base and the intent
// is discarded. On a child die the value is set as intent.
func (b *DieMutBase[T, D, M]) DieFeed(r *T) error {
	if r == nil {
		return b.ref().Feed(reflect.Zero(reflect.TypeFor[T]()))
	}
	return b.ref().Feed(reflect.ValueOf(r).Elem())
}

// DieStamp passes the die itself to each fn, so their writes change it, and
// returns it.
func (b *DieMutBase[T, D, M]) DieStamp(fns ...func(m M)) M {
	return core.Must(b.DieStampE(withE(fns)...))
}

// DieStampE is [DieMutBase.DieStamp] for functions that can fail: it passes the
// die to each fn until one returns an error, and returns that error with the
// die, as the functions left it.
func (b *DieMutBase[T, D, M]) DieStampE(fns ...func(m M) error) (M, error) {
	r := b.ref()
	m := b.self()
	err := runWith(m, fns)
	withIn(r, m)
	return m, err
}

// DieImmutable returns an immutable copy of the die: of a root, its value and
// intent; of a die within another, or presented by a jig, a new root holding
// its current value, as [DieMutBase.DeepCopy] makes it. Later writes to the
// die never reach the copy.
func (b *DieMutBase[T, D, M]) DieImmutable() D { return core.Must(b.DieImmutableE()) }

// DieImmutableE is [DieMutBase.DieImmutable], returning an error rather than
// panicking.
func (b *DieMutBase[T, D, M]) DieImmutableE() (D, error) {
	h, err := b.ref().CopyOf(false)
	return wrapped(h, err, wrap[T, D, M])
}

// DeepCopy returns an independent, mutable copy: see [DieBase.DeepCopy].
func (b *DieMutBase[T, D, M]) DeepCopy() M { return core.Must(b.DeepCopyE()) }

// DeepCopyE is [DieMutBase.DeepCopy], returning an error rather than
// panicking.
func (b *DieMutBase[T, D, M]) DeepCopyE() (M, error) {
	h, err := b.ref().CopyOf(true)
	return wrapped(h, err, wrapMut[T, D, M])
}

// DieTare folds the intent into the base, in place: see [DieBase.DieTare].
func (b *DieMutBase[T, D, M]) DieTare() error {
	r := b.ref()
	return r.Replace(r.Tare())
}

// Equal reports whether two dies have the same [Content], as
// [DieBase.Equal] does.
func (b *DieMutBase[T, D, M]) Equal(o M) bool { return core.Must(equal(b.ref(), o)) }

// EqualE is [DieMutBase.Equal], returning an error rather than panicking.
func (b *DieMutBase[T, D, M]) EqualE(o M) (bool, error) { return equal(b.ref(), o) }

// UnmarshalJSON replaces the base and keeps the intent, so the intent is
// reapplied on top of the decoded value.
func (b *DieMutBase[T, D, M]) UnmarshalJSON(data []byte) error { return b.ref().Unmarshal(data) }

// Error is a die's error. Writes return one, and so do the E variants of
// reads that can fail, whose fluent forms panic with it: the fluent method for
// content that is known, as in tests and fixtures, and the E variant
// otherwise. Use [errors.As] to tell a die's error from another, and
// [errors.Is] with [ErrImmutable], [ErrImmutableChild] or [ErrReadOnlyJig]
// for those.
//
// Methods that implement Kubernetes' and Go's interfaces, such as
// metav1.Object's setters, can't return an error: they panic with one, such
// as on an immutable die.
type Error = core.Error

var (
	// ErrImmutable is the error of a write to an immutable die through an
	// interface it implements, such as metav1.Object's setters: write
	// within [DieBase.DieStamp], or to a [DieBase.DieMutable] copy.
	ErrImmutable error = core.ErrImmutable
	// ErrImmutableChild is the error of a write through a child of an
	// immutable die, which would have nowhere to go: write within DieStamp.
	ErrImmutableChild error = core.ErrImmutableChild
	// ErrReadOnlyJig is the error of a write through a jig of an immutable
	// die.
	ErrReadOnlyJig error = core.ErrReadOnlyJig
)
