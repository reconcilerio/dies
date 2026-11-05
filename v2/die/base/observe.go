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

package base

import (
	"reflect"

	"reconciler.io/dies/v2/die/internal/core"
	"reconciler.io/dies/v2/die/internal/refs"
)

// Observers: a die can be observed, with DieObserve, at a location within
// it, its target, with die.ObserveAt, or a list map's elements, with its
// DieObserveEach. The observer is called after each mutation that changes
// the target, with a mutable handle to the target, within the die as the
// mutation left it, and the change, a [Change]: immutable forms of the
// target as the observer last saw it, and now. The handle's DieRoot is the
// mutated die, so the observer may write anywhere in it, as well as at its
// target.
//
// What an observer writes notifies the die's other observers, so observers
// cascade: an observer of A writing B calls the observers of B, and so on.
// While an observer runs it isn't called again for its own writes, but once
// it returns it is called again if its target changed since it was called,
// so an observer that writes its own target is called until the target
// settles. An observer must settle: one that changes its target every time
// it is called makes the write that called it fail.
//
// Observers are notified as each write is made, so they see a change written
// in parts, such as a condition's status and then its reason, part by part.
// Within [DieMutBase.DieStampAtomic] they are notified as each of its
// functions returns, of what it wrote.
//
// A write whose observers fail returns their error, and is kept, with
// whatever they wrote: nothing is rolled back implicitly. To make writes and
// what they cause commit or roll back together, write them within
// DieStampAtomic, whose fns draw the boundary: a fn that returns nil commits,
// and one that returns an error, or whose observers fail, rolls back the
// stamp.
//
// Observers are kept by the die and every die derived from it: its stamps,
// mutable and immutable copies, deep copies and tares. A copy of a die
// within another keeps the observers registered at or below it. A jig's
// observers are notified of writes through the jig, not of those made to its
// source directly.

// Change is a change an observer observes: immutable forms of its target, a
// die or a handle, as the observer last saw it, and now, each within its
// root as it was then, so what it is within reads as it was too. Neither is
// copied, as a die's state is never changed in place. Presence is theirs: a
// die's DieHas, or a handle's Has, says whether it was, and is, present; a
// root always is.
type Change[D any] struct {
	Before, After D
}

// DieObserve returns a copy of the die observed, as [DieMutBase.DieObserve]
// observes a mutable die. An immutable die doesn't change, so its observers
// are called within the dies derived from it that do, such as within
// [DieBase.DieStamp].
//
// It panics on a die within another, or a jig, whose copy would be detached
// from what it is within: observe a part of an immutable root within a stamp
// of the root, root.DieStamp(func(m M) { m.Status().DieObserve(fn) }).
func (b *DieBase[T, D, M]) DieObserve(fn func(m M, c Change[D]) error) D {
	return observeCopy[T, D, M](b, "DieObserve", func(h *core.Holder) error { return observeSelf[T, D, M](h, nil, fn) })
}

// DieObserve observes the die: fn is called after each mutation that changes
// its value, with the die, mutable, as the mutation left it, and the change,
// as immutable dies. It is called on this die and the dies derived from it.
// A die within another, such as m.Status(), shares its root's state, so it
// observes that part of the root.
func (b *DieMutBase[T, D, M]) DieObserve(fn func(m M, c Change[D]) error) error {
	return observeSelf[T, D, M](b.hold(), b.at, fn)
}

func observeSelf[T any, D DieOf[T, D, M], M DieMutOf[T, D, M]](h *core.Holder, at Path, fn func(m M, c Change[D]) error) error {
	t := reflect.TypeFor[T]()
	return h.Observe(&core.Observer{
		At:     at,
		Target: at,
		Read: func(h *core.Holder, at Path) (reflect.Value, error) {
			return refs.Ref{H: h, At: at, T: t}.Value()
		},
		Fn: func(h *core.Holder, at Path, d core.Delta) error {
			return fn(wrapMut[T, D, M](h, at), Change[D]{Before: wrap[T, D, M](d.Was, at), After: wrap[T, D, M](d.Now, at)})
		},
	})
}

// observeCopy returns a copy of an immutable root, observed by register.
func observeCopy[T any, D DieOf[T, D, M], M DieMutOf[T, D, M]](b *DieBase[T, D, M], method string, register func(h *core.Holder) error) D {
	r := b.ref()
	switch {
	case len(r.At) > 0:
		panic(core.Errorf("die: %s on an immutable die at %s would observe a detached copy; observe it within a stamp of its root, as a mutable die", method, r.At))
	case r.H.View != nil:
		panic(core.Errorf("die: %s on an immutable jig would observe a detached copy; observe the jig's source, or a mutable jig", method))
	}
	h := core.Must(r.CopyOf(false))
	d := wrap[T, D, M](h, nil)
	core.Check(register(h))
	return d
}

// dieReadOnly is the hook the engine makes an immutable die of a mutable
// root with, for the frozen state an observer's change holds.
func (b *DieMutBase[T, D, M]) dieReadOnly() any { return wrap[T, D, M](b.hold(), b.at) }

func init() {
	core.ReadOnly = func(d any) any {
		if m, ok := d.(interface{ dieReadOnly() any }); ok {
			return m.dieReadOnly()
		}
		return d
	}
}
