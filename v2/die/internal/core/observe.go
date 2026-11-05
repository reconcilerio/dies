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

// An observer is called when a location it observes, its target, changes.
// Observers are held by a holder, and copied with it, so the dies derived
// from a die observe as it does.
//
// Each observer remembers the value it last saw at its target. After a
// mutation of a mutable holder, every observer whose target the mutation
// may have reached reads its target again, and is called with what it saw
// and what is there now, if they differ. What an observer writes notifies
// the others in turn, so observers cascade; while an observer runs, it is
// not called again, but once it returns it is called again if its target
// changed since it was called, until its target settles.

// maxObserving is how deep observers may cascade, and how many times one may
// be called again for its own changes, before they are taken not to settle.
const maxObserving = 32

// Observer is an observer of a holder, as a die registers it.
type Observer struct {
	// At is the location of the die it was registered on.
	At Path
	// Target is the location it observes.
	Target Path
	// Read reads the value at a location of a holder, as the observer sees
	// it.
	Read func(h *Holder, at Path) (reflect.Value, error)
	// Fn is called with the holder, the die's location, and the change.
	Fn func(h *Holder, at Path, d Delta) error
}

// Delta is a change an observer sees: the target's value, and whether it was
// present, when the observer last saw it and now, and the die's state then
// and now, frozen, for immutable dies of them.
type Delta struct {
	Before, After reflect.Value
	Had, Has      bool
	Was, Now      *Holder
}

// ReadOnly returns an immutable die of a die's state, sharing it: the
// immutable die of a mutable die. Package base installs it.
var ReadOnly = func(d any) any { return d }

// frozen returns an immutable holder of the holder's state now, sharing it,
// as it is never changed in place: writes to the holder replace its state.
// A jig's reads go to its source, so its value is copied. The root die of
// the frozen holder is the read-only form of the holder's.
func (h *Holder) frozen() (*Holder, error) {
	f := &Holder{Base: h.Base, Ov: h.Ov, Wire: h.Wire, Schema: h.Schema}
	if h.View != nil {
		v, err := h.ValueAt(nil)
		if err != nil {
			return nil, err
		}
		f = &Holder{Base: v, Schema: h.Schema}
	}
	f.root = func(f *Holder) any {
		if h.root == nil {
			return nil
		}
		return ReadOnly(h.root(f))
	}
	return f, nil
}

// observing is an observer of one holder, with what it last saw.
type observing struct {
	o          *Observer
	at, target Path
	seen       reflect.Value
	seenHas    bool
	// seenAt is the die's state when the observer last saw it, frozen
	seenAt  *Holder
	running bool
}

// Observe adds an observer, which has seen the value now at its target.
func (h *Holder) Observe(o *Observer) error {
	seen, has, err := h.observe(o, o.Target)
	if err != nil {
		return err
	}
	at, err := h.frozen()
	if err != nil {
		return err
	}
	h.observers = append(h.observers[:len(h.observers):len(h.observers)], &observing{o: o, at: o.At, target: o.Target, seen: seen, seenHas: has, seenAt: at})
	return nil
}

// observe reads an observer's target, and whether it is present: a root
// always is.
func (h *Holder) observe(o *Observer, target Path) (reflect.Value, bool, error) {
	v, err := o.Read(h, target)
	if err != nil || len(target) == 0 {
		return v, true, err
	}
	has, err := h.Has(target)
	return v, has, err
}

// HasObservers reports whether the holder has observers.
func (h *Holder) HasObservers() bool { return len(h.observers) > 0 }

// ObserversOf gives the holder copies of the observers of another, the
// source of a copy of its die at within: those registered at or below
// within, rebased onto the copy, whose root is the value there.
func (h *Holder) ObserversOf(src *Holder, within Path) {
	h.observers = nil
	for _, s := range src.observers {
		if !s.at.hasPrefix(within) || !s.target.hasPrefix(within) {
			continue
		}
		seenAt := s.seenAt
		if len(within) > 0 {
			// the copy's root is the value within; so is what it saw
			v, err := seenAt.ValueAt(within)
			if err != nil {
				continue
			}
			seenAt = &Holder{Base: v, Schema: seenAt.Schema.At(within)}
			seenAt.root = func(f *Holder) any {
				if h.root == nil {
					return nil
				}
				return ReadOnly(h.root(f))
			}
		}
		h.observers = append(h.observers, &observing{o: s.o, at: s.at[len(within):], target: s.target[len(within):], seen: s.seen, seenHas: s.seenHas, seenAt: seenAt})
	}
}

// Notify calls the observers whose targets a mutation at p may have reached,
// if their targets changed. It does nothing for an immutable holder, which a
// mutation doesn't change in place.
//
// While an atomic stamp holds notifications, it records p, to notify of when
// the stamp is done.
func (h *Holder) Notify(p Path) error {
	if len(h.observers) == 0 || !h.Mutable {
		return nil
	}
	if h.holding > 0 {
		h.pending = append(h.pending, p)
		return nil
	}
	return h.notify([]Path{p})
}

// notify calls the observers whose targets a mutation at any of paths may
// have reached, if their targets changed.
func (h *Holder) notify(paths []Path) error {
	if h.notifying >= maxObserving {
		return Errorf("die: observers cascaded more than %d deep at %s; they don't settle", maxObserving, paths[0])
	}
	h.notifying++
	defer func() { h.notifying-- }()
	// observers added while notifying aren't notified of this mutation
	for _, s := range h.observers[:len(h.observers):len(h.observers)] {
		if s.running || !slices.ContainsFunc(paths, s.target.overlaps) {
			continue
		}
		if err := h.call(s); err != nil {
			return err
		}
	}
	return nil
}

// Snapshot is a holder's state, to restore if an atomic stamp fails.
type Snapshot struct {
	base      reflect.Value
	ov        *Onode
	wire      *pnode
	observers []*observing
	seen      []reflect.Value
	seenHas   []bool
	seenAt    []*Holder
}

func (h *Holder) snapshot() *Snapshot {
	s := &Snapshot{base: h.Base, ov: h.Ov, wire: h.Wire, observers: h.observers}
	for _, o := range h.observers {
		s.seen, s.seenHas, s.seenAt = append(s.seen, o.seen), append(s.seenHas, o.seenHas), append(s.seenAt, o.seenAt)
	}
	return s
}

// restore restores a holder to a snapshot, its observers included.
func (h *Holder) restore(s *Snapshot) {
	h.Base, h.Ov, h.Wire, h.observers = s.base, s.ov, s.wire, s.observers
	for i, o := range s.observers {
		o.seen, o.seenHas, o.seenAt, o.running = s.seen[i], s.seenHas[i], s.seenAt[i], false
	}
}

// Commit replaces a mutable holder's state in place, as a mutation at p
// does, and notifies its observers. If they fail, the error is returned, and
// the mutation, with whatever they wrote, is kept: only an atomic stamp
// rolls back, at the boundary its caller draws.
func (h *Holder) Commit(base reflect.Value, ov *Onode, wire *pnode, p Path) error {
	h.Base, h.Ov, h.Wire = base, ov, wire
	return h.Notify(p)
}

// Hold holds back notifications until [Holder.Flush] or [Holder.Release],
// returning the holder's state to restore if the stamp fails. Holds nest.
func (h *Holder) Hold() *Snapshot {
	s := h.snapshot()
	h.holding++
	return s
}

// Flush notifies the observers of what was written while held, if this is
// the outermost hold: a hold within another flushes with it. While the
// observers run, nothing is held, so what they write notifies as it is
// written, unless they hold it themselves.
func (h *Holder) Flush() error {
	if h.holding != 1 || len(h.pending) == 0 {
		return nil
	}
	paths := h.pending
	h.pending = nil
	h.holding = 0
	defer func() { h.holding = 1 }()
	return h.notify(paths)
}

// Release releases a hold, flushing it. If err isn't nil, or the observers
// fail, the holder is restored to the snapshot, so whatever was written
// while it was held, by the stamp or its observers, is undone, and the error
// returned.
func (h *Holder) Release(s *Snapshot, err error) error {
	if err == nil {
		err = h.Flush()
	}
	h.holding--
	if err != nil {
		h.restore(s)
		if h.holding == 0 {
			h.pending = nil
		}
	}
	return err
}

// call calls an observer while its target differs from what it saw.
func (h *Holder) call(s *observing) error {
	for range maxObserving {
		now, has, err := h.observe(s.o, s.target)
		if err != nil {
			return err
		}
		if has == s.seenHas && reflect.DeepEqual(s.seen.Interface(), now.Interface()) {
			return nil
		}
		frozen, err := h.frozen()
		if err != nil {
			return err
		}
		d := Delta{Before: s.seen, After: now, Had: s.seenHas, Has: has, Was: s.seenAt, Now: frozen}
		s.seen, s.seenHas, s.seenAt = now, has, frozen
		s.running = true
		err = s.o.Fn(h, s.at, d)
		s.running = false
		if err != nil {
			return err
		}
	}
	return Errorf("die: the observer of %s changed it %d times; it doesn't settle", s.target, maxObserving)
}

func (s Step) named() bool { return s.kind == stepField || s.kind == stepMapKey }

// overlaps reports whether a change at one path may change a value at the
// other: one is within the other. An element found by its index and one
// found by its key may be the same element, so they overlap, as do a field
// and a map entry of the same name.
func (p Path) overlaps(q Path) bool {
	for i := range min(len(p), len(q)) {
		a, b := p[i], q[i]
		switch {
		case a.kind == stepIndex || b.kind == stepIndex:
			// the same element, perhaps
		case a.named() && b.named():
			// a field and a map entry of the same name, as typed and
			// unstructured content address the same value
			if a.name != b.name {
				return false
			}
		case a.id() != b.id():
			return false
		}
	}
	return true
}
