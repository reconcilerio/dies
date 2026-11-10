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

	"reconciler.io/dies/v2/die/base"
	"reconciler.io/dies/v2/die/internal/core"
	"reconciler.io/dies/v2/die/internal/dies"
)

// An atomic list, such as a role's rules or a container's args, is written
// whole, as the API server replaces it whole, but it needn't be built whole:
// a [List] reaches its elements by index, as dies of their own where they
// have them and as [Value] handles where they don't, and a [ListMut] sets,
// appends and removes them. Each write is of the whole list, so patches are
// as they were.

// ListField is a list field of T, other than a list map, of E elements, whose
// handles are ED and EDM: dies of E where E has them, and [Value] and
// [ValueMut] where it doesn't.
type ListField[T, E, ED, EDM any] struct {
	fd      *core.FieldDesc
	wrap    func(h *core.Holder, list base.Path, i int) ED
	wrapMut func(h *core.Holder, list base.Path, i int) EDM
}

// NewList declares a list field of T whose E elements have dies, ED and EDM.
func NewList[T, E any, ED base.DieOf[E, ED, EDM], EDM base.DieMutOf[E, ED, EDM]](wire string) ListField[T, E, ED, EDM] {
	fd := listField[T, E](wire)
	// for code that walks dies by their fields, such as validation
	fd.Dies = childDies[E, ED, EDM]
	return ListField[T, E, ED, EDM]{
		fd: fd,
		wrap: func(h *core.Holder, list base.Path, i int) ED {
			return dies.WrapDie[E, ED](h, core.ChildPath(list, core.IndexStep(i)))
		},
		wrapMut: func(h *core.Holder, list base.Path, i int) EDM {
			return dies.WrapDieMut[E, ED](h, core.ChildPath(list, core.IndexStep(i)))
		},
	}
}

// NewValueList declares a list field of T whose E elements have no dies, such
// as strings, whose handles are [Value] and [ValueMut].
func NewValueList[T, E any](wire string) ListField[T, E, Value[E], ValueMut[E]] {
	return ListField[T, E, Value[E], ValueMut[E]]{
		fd:   listField[T, E](wire),
		wrap: func(h *core.Holder, list base.Path, i int) Value[E] { return Value[E]{elementLeaf(h, list, i)} },
		wrapMut: func(h *core.Holder, list base.Path, i int) ValueMut[E] {
			return ValueMut[E]{Value[E]{elementLeaf(h, list, i)}}
		},
	}
}

func elementLeaf(h *core.Holder, list base.Path, i int) leaf {
	return leaf{h: h, at: list, rel: base.Path{core.IndexStep(i)}}
}

// listField finds a list field of E elements, of any slice type of them.
func listField[T, E any](wire string) *core.FieldDesc {
	fd := lookup[T](wire, nil)
	if fd.Typ.Kind() != reflect.Slice || fd.Typ.Elem() != reflect.TypeFor[E]() {
		panic(core.Errorf("die: field %q of %s is %s, not a list of %s", wire, reflect.TypeFor[T](), fd.Typ, reflect.TypeFor[E]()))
	}
	return fd
}

// ListOf returns the handle of a list field of an immutable die.
func ListOf[T, E, ED, EDM any](d base.LocatedAs[T], f ListField[T, E, ED, EDM]) *List[E, ED] {
	h, at := dies.MustStateOf(d)
	return &List[E, ED]{h: h, at: core.ChildPath(at, core.FieldStep(f.fd.Wire)), t: f.fd.Typ, wrap: f.wrap}
}

// ListMutOf returns the handle of a list field of a mutable die.
func ListMutOf[T, E, ED, EDM any](d base.LocatedAs[T], f ListField[T, E, ED, EDM]) *ListMut[E, ED, EDM] {
	h, at := dies.MustStateOf(d)
	return &ListMut[E, ED, EDM]{List: List[E, EDM]{h: h, at: core.ChildPath(at, core.FieldStep(f.fd.Wire)), t: f.fd.Typ, wrap: f.wrapMut}, readOnly: f.wrap}
}

// List is the handle of a list of E elements of an immutable die, other
// than a list map, whose elements' handles are ED.
type List[E, ED any] struct {
	h  *core.Holder
	at base.Path
	// t is the list's type, a slice of E, which may be named
	t    reflect.Type
	wrap func(h *core.Holder, list base.Path, i int) ED
}

// ListMut is the handle of a list of a mutable die, whose elements'
// handles, EDM, are mutable. ED is the immutable handle of an element, of
// its [ListMut.DieReadOnly].
type ListMut[E, ED, EDM any] struct {
	List[E, EDM]
	readOnly func(h *core.Holder, list base.Path, i int) ED
}

// DieReadOnly returns the immutable list of the same location, sharing this
// one's state: it reads what is written through this one.
func (l *ListMut[E, ED, EDM]) DieReadOnly() *List[E, ED] {
	return &List[E, ED]{h: l.h, at: l.at, t: l.t, wrap: l.readOnly}
}

func (l *List[E, ED]) DiePath() base.Path { return l.at }

func (l *List[E, ED]) handleState() (*core.Holder, base.Path) { return l.h, l.at }

// DieValue marks a list as a location of a []E, for [base.LocatedAs]. It
// does nothing.
func (*List[E, ED]) DieValue(*[]E) {}

// DieRoot returns the root die this is within: see [base.DieBase.DieRoot].
func (l *List[E, ED]) DieRoot() base.Die { return base.RootOf(l.h) }

// DieRootType returns the type of the root value this is within.
func (l *List[E, ED]) DieRootType() reflect.Type { return l.h.Base.Type() }

// Has reports whether the list is present.
func (l *List[E, ED]) Has() bool { return core.Must(l.h.Has(l.at)) }

// Get returns a copy of the list, or nil when it is absent.
func (l *List[E, ED]) Get() []E {
	v := core.Must(l.h.ValueAt(l.at))
	if !v.IsValid() || v.IsNil() {
		return nil
	}
	return v.Convert(reflect.TypeFor[[]E]()).Interface().([]E)
}

// Len returns the number of elements.
func (l *List[E, ED]) Len() int {
	v := core.Must(l.h.ValueAt(l.at))
	if !v.IsValid() {
		return 0
	}
	return v.Len()
}

// At returns the handle of the element at index i: its die, if it has one,
// at its index in the list, such as rules[0] in a validation error. Reading
// an index the list doesn't have returns zero values; writing one panics, as
// indexing a slice does: [ListMut.Append] adds elements.
func (l *List[E, ED]) At(i int) ED { return l.wrap(l.h, l.at, i) }

// All returns the index and handle of each element, in order.
func (l *List[E, ED]) All() iter.Seq2[int, ED] {
	return func(yield func(int, ED) bool) {
		for i := range l.Len() {
			if !yield(i, l.At(i)) {
				return
			}
		}
	}
}

// Set sets the list to copies of the items, dies or handles of E of any
// root, in order. An absent item is the zero E. Setting no items sets the
// list present and empty, as [ListMut.Empty] does.
//
//	d.Rules().Set(readRule, writeRule.DieMutable())
func (l *ListMut[E, ED, EDM]) Set(items ...base.LocatedAs[E]) error {
	out, err := l.located(reflect.MakeSlice(l.t, 0, len(items)), items)
	if err != nil {
		return err
	}
	return l.replace(out)
}

// SetValue sets the list to copies of the values, in order, as [ListMut.Set]
// sets it to dies'.
//
//	d.Args().SetValue("--verbose", "--port=8080")
func (l *ListMut[E, ED, EDM]) SetValue(items ...E) error {
	return l.replace(values(reflect.MakeSlice(l.t, 0, len(items)), items))
}

// Append appends copies of the items, dies or handles of E of any root, to
// the list, creating it if it's absent. An absent item is the zero E.
func (l *ListMut[E, ED, EDM]) Append(items ...base.LocatedAs[E]) error {
	out, err := l.located(l.current(len(items)), items)
	if err != nil {
		return err
	}
	return l.replace(out)
}

// AppendValue appends copies of the values to the list, creating it if it's
// absent, as [ListMut.Append] appends dies'.
func (l *ListMut[E, ED, EDM]) AppendValue(items ...E) error {
	return l.replace(values(l.current(len(items)), items))
}

// Remove removes the element at index i, moving those after it down one. It
// returns an error if the list has no such element.
func (l *ListMut[E, ED, EDM]) Remove(i int) error {
	if n := l.Len(); i < 0 || i >= n {
		return core.Errorf("die: removing %s: the list has %d elements", core.ChildPath(l.at, core.IndexStep(i)), n)
	}
	_, err := l.h.Write(l.at, base.Path{core.IndexStep(i)}, core.SetOp(core.OpClear, reflect.Value{}), false)
	return err
}

// Clear removes the list: null in patches.
func (l *ListMut[E, ED, EDM]) Clear() error {
	_, err := l.h.Write(l.at, nil, core.SetOp(core.OpClear, reflect.Value{}), false)
	return err
}

// Empty sets the list to present and empty: [] in patches.
func (l *ListMut[E, ED, EDM]) Empty() error {
	_, err := l.h.Write(l.at, nil, core.SetOp(core.OpEmpty, reflect.Value{}), true)
	return err
}

// Replace replaces the list with a copy of src's, of a die or handle of any
// root, or removes it, as Clear does, when src is absent, as intent: it
// writes nothing when the lists are the same. The list is a copy, so later
// writes to src aren't seen here.
func (l *ListMut[E, ED, EDM]) Replace(src base.LocatedAs[[]E]) error {
	return core.ReplaceAt(l.h, l.at, nil, src)
}

// current returns a copy of the list, of its type, with room for n more
// elements, or an empty list of its type when it is absent.
func (l *ListMut[E, ED, EDM]) current(n int) reflect.Value {
	v := core.Must(l.h.ValueAt(l.at))
	out := reflect.MakeSlice(l.t, 0, n)
	if v.IsValid() {
		out = reflect.MakeSlice(l.t, 0, v.Len()+n)
		out = reflect.AppendSlice(out, v.Convert(l.t))
	}
	return out
}

// located appends copies of the items' values to out.
func (l *ListMut[E, ED, EDM]) located(out reflect.Value, items []base.LocatedAs[E]) (reflect.Value, error) {
	t := reflect.TypeFor[E]()
	for _, item := range items {
		ok, v, err := core.ReadLocated(item, t)
		if err != nil {
			return reflect.Value{}, err
		}
		if !ok {
			v = reflect.Zero(t)
		}
		out = reflect.Append(out, v)
	}
	return out, nil
}

// values appends copies of the items to out.
func values[E any](out reflect.Value, items []E) reflect.Value {
	for i := range items {
		out = reflect.Append(out, core.DeepCopy(reflect.ValueOf(&items[i]).Elem()))
	}
	return out
}

// replace sets the list to v, a list of its type that nothing else holds,
// as intent: nothing is written if the list is v already.
func (l *ListMut[E, ED, EDM]) replace(v reflect.Value) error {
	if err := core.CheckKeys(v); err != nil {
		return err
	}
	return core.ReplaceValueAt(l.h, l.at, nil, true, v)
}
