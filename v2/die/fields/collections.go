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
	"iter"
	"reflect"
	"slices"

	"reconciler.io/dies/v2/die/base"
	"reconciler.io/dies/v2/die/internal/core"
)

// Collections are always attached to a die. An immutable die's read, and
// return immutable dies of their elements; a mutable die's, ListMapMut and
// MapMut, embed them and also write, in place.

// ListMap is a list of C elements, each identified by a key K, whose dies
// are CD.
type ListMap[C, CD, K any] struct {
	h     *core.Holder
	at    base.Path
	keys  []string
	key   func(K) []any
	keyOf func(C) K
	wrap  func(*core.Holder, core.Path) CD
}

// ListMapMut is the list map of a mutable die, whose element dies, CDM, are
// mutable. CD is the immutable die of an element, of its [ListMapMut.DieReadOnly].
type ListMapMut[C, CD, CDM, K any] struct {
	ListMap[C, CDM, K]
	readOnly func(*core.Holder, core.Path) CD
}

// DieReadOnly returns the immutable list map of the same location, sharing
// this one's state, as [base.DieMutBase.DieReadOnly] does a die's: it reads
// what is written through this one.
func (l *ListMapMut[C, CD, CDM, K]) DieReadOnly() *ListMap[C, CD, K] {
	return &ListMap[C, CD, K]{h: l.h, at: l.at, keys: l.keys, key: l.key, keyOf: l.keyOf, wrap: l.readOnly}
}

func (l *ListMap[C, CD, K]) DiePath() base.Path { return l.at }

func (l *ListMap[C, CD, K]) handleState() (*core.Holder, base.Path) { return l.h, l.at }

// DieValue marks a list map as a location of a []C, for [base.LocatedAs]. It
// does nothing.
func (*ListMap[C, CD, K]) DieValue(*[]C) {}

// DieRoot returns the root die this is within: see [base.DieBase.DieRoot].
func (l *ListMap[C, CD, K]) DieRoot() base.Die { return base.RootOf(l.h) }

// DieRootType returns the type of the root value this is within.
func (l *ListMap[C, CD, K]) DieRootType() reflect.Type { return l.h.Base.Type() }

func (l *ListMap[C, CD, K]) step(k K) core.Step {
	return core.KeyStep(reflect.TypeFor[C](), l.keys, l.key(k))
}

// At returns the element with key k. Reading a missing element returns zero
// values; the first write creates it with its key fields set.
func (l *ListMap[C, CD, K]) At(k K) CD {
	return l.wrap(l.h, core.ChildPath(l.at, l.step(k)))
}

// Has reports whether an element with key k exists.
func (l *ListMap[C, CD, K]) Has(k K) bool {
	list := core.Must(l.h.ValueAt(l.at))
	i := core.Must(core.FindElem(list, l.step(k)))
	return i >= 0
}

// write writes through the list, returning it.
func (l *ListMapMut[C, CD, CDM, K]) write(rel base.Path, o core.Op, union bool) error {
	_, err := l.h.Write(l.at, rel, core.SetOp(o, reflect.Value{}), union)
	return err
}

// Remove removes the element with key k.
func (l *ListMapMut[C, CD, CDM, K]) Remove(k K) error {
	return l.write(base.Path{l.step(k)}, core.OpClear, false)
}

// Set sets the element with c's key to c: it is appended if the list
// doesn't have one, and replaced if it does, as intent, as [ListMapMut.Replace]
// replaces a list: what differs is written, field by field, so fields c
// doesn't have are cleared, and nothing is written if they're the same. The
// element is a copy of c, so later changes to c aren't seen here.
//
//	d.Spec().Containers().Set(corev1.Container{Name: "app", Image: "app:2"})
func (l *ListMapMut[C, CD, CDM, K]) Set(c C) error {
	v := core.DeepCopy(reflect.ValueOf(&c).Elem())
	if err := core.CheckKeys(v); err != nil {
		return err
	}
	// keyed as the list keys its elements, with its key fields' defaults
	k := l.keyOf(core.WithKeyDefaults(v, l.keys).Interface().(C))
	return core.ReplaceValueAt(l.h, l.at, base.Path{l.step(k)}, true, v)
}

// Clear removes the list: null in patches.
func (l *ListMapMut[C, CD, CDM, K]) Clear() error {
	return l.write(nil, core.OpClear, false)
}

// Empty sets the list to present and empty: [] in patches.
func (l *ListMapMut[C, CD, CDM, K]) Empty() error {
	return l.write(nil, core.OpEmpty, true)
}

// Replace replaces the list with a copy of src's, of a die or handle of any
// root, or removes it, as Clear does, when src is absent. It records the
// diff as intent, as writes of each change would: elements src doesn't have
// are removed, its new elements are set, and those of both are diffed by
// field. The list is a copy, so later writes to src aren't seen here.
//
//	current.Status().Conditions().Replace(desired.Status().Conditions())
func (l *ListMapMut[C, CD, CDM, K]) Replace(src base.LocatedAs[[]C]) error {
	return core.ReplaceAt(l.h, l.at, nil, src)
}

// Len returns the number of elements.
func (l *ListMap[C, CD, K]) Len() int {
	list := core.Must(l.h.ValueAt(l.at))
	return list.Len()
}

// Keys returns the element keys in list order, with defaulted key fields
// filled in.
func (l *ListMap[C, CD, K]) Keys() []K {
	list := core.Must(l.h.ValueAt(l.at))
	keys := make([]K, list.Len())
	for i := range list.Len() {
		keys[i] = l.keyOf(core.WithKeyDefaults(list.Index(i), l.keys).Interface().(C))
	}
	return keys
}

// All yields each element's key and attached die, in list order.
func (l *ListMap[C, CD, K]) All() iter.Seq2[K, CD] {
	return func(yield func(K, CD) bool) {
		for _, k := range l.Keys() {
			if !yield(k, l.At(k)) {
				return
			}
		}
	}
}

// DieObserveEach observes each element of the list: after each mutation of
// the list, fn is called once for each element it changed, added or
// removed, with the element's die, mutable, and the element's change, as
// immutable dies of the element as the observer last saw it and now, so the
// elements come matched by key. Changed and added elements are called in
// list order, then removed ones, in the order they were. The die of an
// element that was absent, or is, reads as absent, so DieHas tells an added
// or removed element; writing through a removed element's die adds it
// again.
//
// It is an observer of the list, as die.ObserveAt registers one: it is kept
// by the dies derived from the die, cascades, settles and is held within an
// atomic stamp, as other observers are. An element's die reads it as it is,
// which another element's observer may have changed since the change, as
// the observer is then called again.
//
// The list of an immutable die is observed within a stamp of its root, as
// a die within one is.
func (l *ListMapMut[C, CD, CDM, K]) DieObserveEach(fn func(e CDM, c base.Change[CD]) error) error {
	proto := *l
	on := func(h *core.Holder, at core.Path) *ListMapMut[C, CD, CDM, K] {
		list := proto
		list.h, list.at = h, at
		return &list
	}
	return l.h.Observe(&core.Observer{
		At:     l.at,
		Target: l.at,
		Read:   func(h *core.Holder, at core.Path) (reflect.Value, error) { return h.ValueAt(at) },
		Fn: func(h *core.Holder, at core.Path, d core.Delta) error {
			was, now := on(d.Was, at).DieReadOnly(), on(d.Now, at).DieReadOnly()
			return on(h, at).each(d, func(k K, e CDM) error {
				return fn(e, base.Change[CD]{Before: was.At(k), After: now.At(k)})
			})
		},
	})
}

// each calls fn for each element that differs between a change's values of
// the list.
func (l *ListMap[C, CD, K]) each(d core.Delta, fn func(k K, e CD) error) error {
	was, wasOrder := l.byKey(d.Before)
	now, order := l.byKey(d.After)
	for _, k := range order {
		if b, had := was[any(k)]; had && reflect.DeepEqual(now[any(k)], b) {
			continue
		}
		if err := fn(k, l.At(k)); err != nil {
			return err
		}
	}
	for _, k := range wasOrder {
		if _, has := now[any(k)]; has {
			continue
		}
		if err := fn(k, l.At(k)); err != nil {
			return err
		}
	}
	return nil
}

// byKey returns a value of the list's elements by key, with the keys in list
// order.
func (l *ListMap[C, CD, K]) byKey(v reflect.Value) (map[any]C, []K) {
	for v.IsValid() && (v.Kind() == reflect.Interface || v.Kind() == reflect.Pointer) && !v.IsNil() {
		v = v.Elem()
	}
	out := map[any]C{}
	var order []K
	if !v.IsValid() || v.Kind() != reflect.Slice {
		return out, order
	}
	for i := range v.Len() {
		e := v.Index(i).Interface().(C)
		k := l.keyOf(core.WithKeyDefaults(v.Index(i), l.keys).Interface().(C))
		out[any(k)] = e
		order = append(order, k)
	}
	return out, order
}

// Map is a map of string-kinded keys K to V values.
type Map[K ~string, V any] struct {
	h  *core.Holder
	at base.Path
}

// MapMut is the map of a mutable die.
type MapMut[K ~string, V any] struct {
	Map[K, V]
}

// DieReadOnly returns the immutable map of the same location, sharing this
// one's state: it reads what is written through this one.
func (m *MapMut[K, V]) DieReadOnly() *Map[K, V] {
	ro := m.Map
	return &ro
}

func (m *Map[K, V]) DiePath() base.Path { return m.at }

func (m *Map[K, V]) handleState() (*core.Holder, base.Path) { return m.h, m.at }

// DieValue marks a map as a location of a map[K]V, for [base.LocatedAs]. It
// does nothing.
func (*Map[K, V]) DieValue(*map[K]V) {}

// DieRoot returns the root die this is within: see [base.DieBase.DieRoot].
func (m *Map[K, V]) DieRoot() base.Die { return base.RootOf(m.h) }

// DieRootType returns the type of the root value this is within.
func (m *Map[K, V]) DieRootType() reflect.Type { return m.h.Base.Type() }

func (m *Map[K, V]) entry(k K) base.Path { return core.ChildPath(m.at, core.MapKeyStep(string(k))) }

// At returns the value for key k, or the zero value.
func (m *Map[K, V]) At(k K) V {
	v := core.Must(m.h.ValueAt(m.entry(k)))
	return v.Interface().(V)
}

// Entry returns the handle of the value for key k, whether or not the map
// has it: a location of a V, as a field's handle is, for what takes one, such
// as a validation error about the value rather than the map, a scope, a jig's
// binding or [ValueMut.Replace]. Its path is the map's with the key.
//
//	for k, v := range d.Metadata().Labels().All() {
//		if !valid(v) {
//			errs = append(errs, validation.Invalid(d.Metadata().Labels().Entry(k), "must be a valid label value"))
//		}
//	}
func (m *Map[K, V]) Entry(k K) Value[V] {
	return Value[V]{m.entryLeaf(k)}
}

// Entry returns the mutable handle of the value for key k, whether or not
// the map has it, as [Map.Entry] does: its Set sets the key, and its Clear
// removes it.
func (m *MapMut[K, V]) Entry(k K) ValueMut[V] {
	return ValueMut[V]{Value[V]{m.entryLeaf(k)}}
}

func (m *Map[K, V]) entryLeaf(k K) leaf {
	return leaf{h: m.h, at: m.at, rel: base.Path{core.MapKeyStep(string(k))}}
}

// Has reports whether the map has key k.
func (m *Map[K, V]) Has(k K) bool {
	return core.Must(m.h.Has(m.entry(k)))
}

// write writes through the map, returning it.
func (m *MapMut[K, V]) write(rel base.Path, o core.Op, v reflect.Value, union bool) error {
	_, err := m.h.Write(m.at, rel, core.SetOp(o, v), union)
	return err
}

// Set sets the value for key k.
func (m *MapMut[K, V]) Set(k K, v V) error {
	return m.write(base.Path{core.MapKeyStep(string(k))}, core.OpSet, core.DeepCopy(reflect.ValueOf(&v).Elem()), true)
}

// Remove removes key k.
func (m *MapMut[K, V]) Remove(k K) error {
	return m.write(base.Path{core.MapKeyStep(string(k))}, core.OpClear, reflect.Value{}, false)
}

// Clear removes the map: null in patches.
func (m *MapMut[K, V]) Clear() error {
	return m.write(nil, core.OpClear, reflect.Value{}, false)
}

// Empty sets the map to present and empty: {} in patches.
func (m *MapMut[K, V]) Empty() error {
	return m.write(nil, core.OpEmpty, reflect.Value{}, true)
}

// Replace replaces the map with a copy of src's, of a die or handle of any
// root, or removes it, as Clear does, when src is absent. It records the
// diff as intent, as writes of each change would: keys src doesn't have are
// removed, and those whose values differ are set. The map is a copy, so
// later writes to src aren't seen here.
//
//	current.Metadata().Labels().Replace(desired.Metadata().Labels())
func (m *MapMut[K, V]) Replace(src base.LocatedAs[map[K]V]) error {
	return core.ReplaceAt(m.h, m.at, nil, src)
}

// Len returns the number of entries.
func (m *Map[K, V]) Len() int {
	v := core.Must(m.h.ValueAt(m.at))
	return v.Len()
}

// Keys returns the keys in sorted order.
func (m *Map[K, V]) Keys() []K {
	v := core.Must(m.h.ValueAt(m.at))
	keys := make([]K, 0, v.Len())
	for _, k := range v.MapKeys() {
		keys = append(keys, K(k.String()))
	}
	slices.Sort(keys)
	return keys
}

// All yields each entry in sorted key order.
func (m *Map[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		v := core.Must(m.h.ValueAt(m.at))
		for _, k := range m.Keys() {
			if !yield(k, v.MapIndex(reflect.ValueOf(k).Convert(v.Type().Key())).Interface().(V)) {
				return
			}
		}
	}
}
