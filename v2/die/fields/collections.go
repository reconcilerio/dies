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
// mutable.
type ListMapMut[C, CDM, K any] struct {
	ListMap[C, CDM, K]
}

func (l *ListMap[C, CD, K]) DiePath() base.Path { return l.at }

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
func (l *ListMap[C, CD, K]) Has(k K) bool { return core.Must(l.HasE(k)) }

// HasE is [ListMap.Has], returning an error rather than panicking.
func (l *ListMap[C, CD, K]) HasE(k K) (bool, error) {
	list, err := l.h.ValueAt(l.at)
	if err != nil {
		return false, err
	}
	i, err := core.FindElem(list, l.step(k))
	return i >= 0, err
}

// write writes through the list, returning it.
func (l *ListMapMut[C, CDM, K]) write(rel base.Path, o core.Op, union bool) (*ListMapMut[C, CDM, K], error) {
	if _, err := l.h.Write(l.at, rel, core.SetOp(o, reflect.Value{}), union); err != nil {
		return nil, err
	}
	return l, nil
}

// Remove removes the element with key k.
func (l *ListMapMut[C, CDM, K]) Remove(k K) *ListMapMut[C, CDM, K] { return core.Must(l.RemoveE(k)) }

// RemoveE is [ListMapMut.Remove], returning an error rather than panicking.
func (l *ListMapMut[C, CDM, K]) RemoveE(k K) (*ListMapMut[C, CDM, K], error) {
	return l.write(base.Path{l.step(k)}, core.OpClear, false)
}

// Clear removes the list: null in patches.
func (l *ListMapMut[C, CDM, K]) Clear() *ListMapMut[C, CDM, K] { return core.Must(l.ClearE()) }

// ClearE is [ListMapMut.Clear], returning an error rather than panicking.
func (l *ListMapMut[C, CDM, K]) ClearE() (*ListMapMut[C, CDM, K], error) {
	return l.write(nil, core.OpClear, false)
}

// Empty sets the list to present and empty: [] in patches.
func (l *ListMapMut[C, CDM, K]) Empty() *ListMapMut[C, CDM, K] { return core.Must(l.EmptyE()) }

// EmptyE is [ListMapMut.Empty], returning an error rather than panicking.
func (l *ListMapMut[C, CDM, K]) EmptyE() (*ListMapMut[C, CDM, K], error) {
	return l.write(nil, core.OpEmpty, true)
}

// Len returns the number of elements.
func (l *ListMap[C, CD, K]) Len() int { return core.Must(l.LenE()) }

// LenE is [ListMap.Len], returning an error rather than panicking.
func (l *ListMap[C, CD, K]) LenE() (int, error) {
	list, err := l.h.ValueAt(l.at)
	if err != nil {
		return 0, err
	}
	return list.Len(), nil
}

// Keys returns the element keys in list order, with defaulted key fields
// filled in.
func (l *ListMap[C, CD, K]) Keys() []K { return core.Must(l.KeysE()) }

// KeysE is [ListMap.Keys], returning an error rather than panicking.
func (l *ListMap[C, CD, K]) KeysE() ([]K, error) {
	list, err := l.h.ValueAt(l.at)
	if err != nil {
		return nil, err
	}
	keys := make([]K, list.Len())
	for i := range list.Len() {
		keys[i] = l.keyOf(core.WithKeyDefaults(list.Index(i), l.keys).Interface().(C))
	}
	return keys, nil
}

// All yields each element's key and attached die, in list order. It panics,
// as [ListMap.Keys] does, if the keys can't be read; use [ListMap.KeysE] and
// [ListMap.At] to handle that.
func (l *ListMap[C, CD, K]) All() iter.Seq2[K, CD] {
	return func(yield func(K, CD) bool) {
		for _, k := range l.Keys() {
			if !yield(k, l.At(k)) {
				return
			}
		}
	}
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

func (m *Map[K, V]) DiePath() base.Path { return m.at }

// DieRoot returns the root die this is within: see [base.DieBase.DieRoot].
func (m *Map[K, V]) DieRoot() base.Die { return base.RootOf(m.h) }

// DieRootType returns the type of the root value this is within.
func (m *Map[K, V]) DieRootType() reflect.Type { return m.h.Base.Type() }

func (m *Map[K, V]) entry(k K) base.Path { return core.ChildPath(m.at, core.MapKeyStep(string(k))) }

// At returns the value for key k, or the zero value.
func (m *Map[K, V]) At(k K) V { return core.Must(m.AtE(k)) }

// AtE is [Map.At], returning an error rather than panicking.
func (m *Map[K, V]) AtE(k K) (V, error) {
	v, err := m.h.ValueAt(m.entry(k))
	if err != nil {
		var zero V
		return zero, err
	}
	return v.Interface().(V), nil
}

// Has reports whether the map has key k.
func (m *Map[K, V]) Has(k K) bool { return core.Must(m.HasE(k)) }

// HasE is [Map.Has], returning an error rather than panicking.
func (m *Map[K, V]) HasE(k K) (bool, error) { return m.h.Has(m.entry(k)) }

// write writes through the map, returning it.
func (m *MapMut[K, V]) write(rel base.Path, o core.Op, v reflect.Value, union bool) (*MapMut[K, V], error) {
	if _, err := m.h.Write(m.at, rel, core.SetOp(o, v), union); err != nil {
		return nil, err
	}
	return m, nil
}

// Set sets the value for key k.
func (m *MapMut[K, V]) Set(k K, v V) *MapMut[K, V] { return core.Must(m.SetE(k, v)) }

// SetE is [MapMut.Set], returning an error rather than panicking.
func (m *MapMut[K, V]) SetE(k K, v V) (*MapMut[K, V], error) {
	return m.write(base.Path{core.MapKeyStep(string(k))}, core.OpSet, core.DeepCopy(reflect.ValueOf(&v).Elem()), true)
}

// Remove removes key k.
func (m *MapMut[K, V]) Remove(k K) *MapMut[K, V] { return core.Must(m.RemoveE(k)) }

// RemoveE is [MapMut.Remove], returning an error rather than panicking.
func (m *MapMut[K, V]) RemoveE(k K) (*MapMut[K, V], error) {
	return m.write(base.Path{core.MapKeyStep(string(k))}, core.OpClear, reflect.Value{}, false)
}

// Clear removes the map: null in patches.
func (m *MapMut[K, V]) Clear() *MapMut[K, V] { return core.Must(m.ClearE()) }

// ClearE is [MapMut.Clear], returning an error rather than panicking.
func (m *MapMut[K, V]) ClearE() (*MapMut[K, V], error) {
	return m.write(nil, core.OpClear, reflect.Value{}, false)
}

// Empty sets the map to present and empty: {} in patches.
func (m *MapMut[K, V]) Empty() *MapMut[K, V] { return core.Must(m.EmptyE()) }

// EmptyE is [MapMut.Empty], returning an error rather than panicking.
func (m *MapMut[K, V]) EmptyE() (*MapMut[K, V], error) {
	return m.write(nil, core.OpEmpty, reflect.Value{}, true)
}

// Len returns the number of entries.
func (m *Map[K, V]) Len() int { return core.Must(m.LenE()) }

// LenE is [Map.Len], returning an error rather than panicking.
func (m *Map[K, V]) LenE() (int, error) {
	v, err := m.h.ValueAt(m.at)
	if err != nil {
		return 0, err
	}
	return v.Len(), nil
}

// Keys returns the keys in sorted order.
func (m *Map[K, V]) Keys() []K { return core.Must(m.KeysE()) }

// KeysE is [Map.Keys], returning an error rather than panicking.
func (m *Map[K, V]) KeysE() ([]K, error) {
	v, err := m.h.ValueAt(m.at)
	if err != nil {
		return nil, err
	}
	keys := make([]K, 0, v.Len())
	for _, k := range v.MapKeys() {
		keys = append(keys, K(k.String()))
	}
	slices.Sort(keys)
	return keys, nil
}

// All yields each entry in sorted key order. It panics, as [Map.Keys] does,
// if the map can't be read; use [Map.KeysE] and [Map.AtE] to handle that.
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
