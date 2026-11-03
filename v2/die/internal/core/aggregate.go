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
	"strings"
)

type Member struct {
	Src  Path
	Root reflect.Type
	// single members are one value, present when set, rather than a list
	Single bool
	// view presents source elements of another type as X
	View *ViewSource
	// Err is why the member couldn't be made, reported by its binding
	Err error
}

// ViewSource is a view, erased: its bindings for a source type, and the type
// it presents.
type ViewSource struct {
	BindingsFor func(t reflect.Type) ([]Binding, bool)
	Type        reflect.Type
}

type Aggregate struct {
	Members []Member
}

// element is an element of an aggregate, located in its source.
//
// An element of a keyed source list, or a single value, is located by its
// path. An element of an atomic list has no key of its own, so it is located
// by its index: it is read and written as a detached copy, and a change is
// written back as the whole list, as befits an atomic list.
type element struct {
	member *Member
	// at is the element's location in the source, or for an element of an
	// atomic list, the list's
	at Path
	// index and detached locate an element of an atomic list
	index    int
	detached *Holder
	// view presents the element, for members with an element view
	view *Viewing
	key  []any
}

func (e element) valueAt(v *Viewing, rest Path, t reflect.Type) (reflect.Value, error) {
	switch {
	case e.view != nil:
		val, err := e.view.valueAt(rest)
		return adapt(val, t), err
	case e.detached != nil:
		val, err := e.detached.ValueAt(rest)
		return adapt(val, t), err
	}
	return v.Src.valueAt(e.at.join(rest), t)
}

func (e element) has(v *Viewing, rest Path) (bool, error) {
	switch {
	case e.view != nil:
		return e.view.has(rest)
	case e.detached != nil:
		return e.detached.Has(rest)
	}
	return v.Src.has(e.at.join(rest))
}

func (e element) write(v *Viewing, rest Path, o Op, val reflect.Value, t reflect.Type) error {
	switch {
	case e.view != nil:
		return e.view.write(rest, o, val)
	case e.detached == nil:
		return v.writeSource(e.at.join(rest), o, val, t)
	}
	if _, err := e.detached.Write(nil, rest, SetOp(o, val), o != OpClear); err != nil {
		return err
	}
	list, err := v.Src.valueAt(e.at, nil)
	if err != nil {
		return err
	}
	el, err := e.detached.Materialize()
	if err != nil {
		return err
	}
	list = DeepCopy(list)
	list.Index(e.index).Set(adapt(el, list.Type().Elem()))
	return v.writeSource(e.at, OpSet, list, nil)
}

// remove removes the element from its list, or clears a single value.
func (e element) remove(v *Viewing) error {
	if e.detached == nil {
		return v.writeSource(e.at, OpClear, reflect.Value{}, nil)
	}
	list, err := v.Src.valueAt(e.at, nil)
	if err != nil {
		return err
	}
	out := reflect.MakeSlice(list.Type(), 0, list.Len()-1)
	out = reflect.AppendSlice(out, list.Slice(0, e.index))
	out = reflect.AppendSlice(out, list.Slice(e.index+1, list.Len()))
	return v.writeSource(e.at, OpSet, DeepCopy(out), nil)
}

func describeMembers(ms []*Member) string {
	var names []string
	for _, m := range ms {
		names = append(names, m.Src.String())
	}
	return strings.Join(names, " and ")
}

// fieldAt returns the descriptor of the field at a path.
func fieldAt(root reflect.Type, p Path) *FieldDesc {
	return Describe(Deref(TypeAt(root, p[:len(p)-1]))).Field(p[len(p)-1].name)
}

// keys are the key fields of the view's list map.
func (v *Viewing) keys(b Binding) []string {
	return Describe(Deref(TypeAt(v.Vtype, b.view[:len(b.view)-1]))).Field(b.view[len(b.view)-1].name).Keys
}

// elements lists the elements of an aggregate, with where each is in the
// source.
func (v *Viewing) elements(b Binding) ([]element, error) {
	keys := v.keys(b)
	elemType := TypeAt(v.Vtype, b.view).Elem()
	var out []element
	add := func(m *Member, at Path, index int, detached *Holder) error {
		e := element{member: m, at: at, index: index, detached: detached}
		if m.View != nil {
			src := v.Src
			sat := at
			if detached != nil {
				src, sat = HolderSource{H: detached}, nil
			}
			if src.typeAt(sat) == nil {
				return Errorf("die: the elements of %s are presented by a jig, which needs a typed source", b.view)
			}
			st := Deref(src.typeAt(sat))
			binds, ok := m.View.BindingsFor(st)
			if !ok {
				return Errorf("die: the jig for %s has no bindings for %s", b.view, st)
			}
			e.view = &Viewing{Src: src, At: sat, Binds: binds, Vtype: m.View.Type}
		}
		val, err := e.valueAt(v, nil, elemType)
		if err != nil {
			return err
		}
		sd := Describe(elemType)
		for _, k := range keys {
			e.key = append(e.key, sd.Field(k).keyValue(val.FieldByIndex(sd.Field(k).index)).Interface())
		}
		out = append(out, e)
		return nil
	}
	for i := range b.agg.Members {
		m := &b.agg.Members[i]
		at := v.At.join(m.Src)
		if m.Single {
			has, err := v.Src.has(at)
			if err != nil {
				return nil, err
			}
			if has {
				if err := add(m, at, -1, nil); err != nil {
					return nil, err
				}
			}
			continue
		}
		lt := v.Src.typeAt(at)
		if lt == nil {
			lt = reflect.SliceOf(elemType)
		}
		list, err := v.Src.valueAt(at, lt)
		if err != nil {
			return nil, err
		}
		listKeys := v.Src.listKeys(at, keys)
		if len(listKeys) == 0 {
			// an atomic list's elements are found by index
			for j := range list.Len() {
				if err := add(m, at, j, &Holder{Base: DeepCopy(list.Index(j)), Mutable: true}); err != nil {
					return nil, err
				}
			}
			continue
		}
		for j := range list.Len() {
			el := list.Index(j)
			sd := Describe(el.Type())
			vals := make([]any, len(listKeys))
			for k, key := range listKeys {
				vals[k] = el.FieldByIndex(sd.Field(key).index).Interface()
			}
			if err := add(m, at.child(KeyStep(el.Type(), listKeys, vals)), -1, nil); err != nil {
				return nil, err
			}
		}
	}
	return out, nil
}

// find returns the elements whose key matches a view key step.
func (v *Viewing) find(b Binding, s Step) ([]element, error) {
	elements, err := v.elements(b)
	if err != nil {
		return nil, err
	}
	var out []element
	for _, e := range elements {
		if reflect.DeepEqual(e.key, s.vals) {
			out = append(out, e)
		}
	}
	return out, nil
}

func (v *Viewing) readAggregate(b Binding, rest Path) (reflect.Value, error) {
	t := TypeAt(v.Vtype, b.view.join(rest))
	if len(rest) == 0 {
		elements, err := v.elements(b)
		if err != nil {
			return reflect.Value{}, err
		}
		elemType := t.Elem()
		list := reflect.MakeSlice(t, 0, 0)
		for _, e := range elements {
			val, err := e.valueAt(v, nil, elemType)
			if err != nil {
				return reflect.Value{}, err
			}
			list = reflect.Append(list, val)
		}
		if list.Len() == 0 {
			return reflect.Zero(t), nil
		}
		return list, nil
	}
	if rest[0].kind != stepKey {
		panic(Errorf("die: an aggregate list is addressed by key, not %s", rest[0]))
	}
	found, err := v.find(b, rest[0])
	if err != nil {
		return reflect.Value{}, err
	}
	switch len(found) {
	case 0:
		val := newElem(TypeAt(v.Vtype, b.view).Elem(), rest[0])
		for _, s := range rest[1:] {
			if val, err = descend(val, s); err != nil {
				return reflect.Value{}, err
			}
		}
		return adapt(DeepCopy(val), t), nil
	case 1:
		return found[0].valueAt(v, rest[1:], t)
	default:
		return reflect.Value{}, v.duplicate(b, rest[0], found)
	}
}

func (v *Viewing) hasAggregate(b Binding, rest Path) (bool, error) {
	if len(rest) == 0 {
		elements, err := v.elements(b)
		return len(elements) > 0, err
	}
	found, err := v.find(b, rest[0])
	if err != nil {
		return false, err
	}
	switch len(found) {
	case 0:
		return false, nil
	case 1:
		return found[0].has(v, rest[1:])
	default:
		return false, v.duplicate(b, rest[0], found)
	}
}

func (v *Viewing) writeAggregate(b Binding, rest Path, o Op, val reflect.Value) error {
	if len(rest) == 0 {
		panic(Errorf("die: %s is an aggregate of %s; write its elements, not the whole list", b.view, describeMembers(v.lists(b, true))))
	}
	if rest[0].kind != stepKey {
		panic(Errorf("die: an aggregate list is addressed by key, not %s", rest[0]))
	}
	found, err := v.find(b, rest[0])
	if err != nil {
		return err
	}
	switch len(found) {
	case 0:
		return v.create(b, rest, o, val)
	case 1:
		e := found[0]
		if len(rest) == 1 && o == OpClear {
			return e.remove(v)
		}
		return e.write(v, rest[1:], o, val, TypeAt(v.Vtype, b.view.join(rest)))
	default:
		return v.duplicate(b, rest[0], found)
	}
}

// lists are the members that are lists, optionally including those whose
// elements are presented by a view.
func (v *Viewing) lists(b Binding, viewed bool) []*Member {
	var out []*Member
	for i := range b.agg.Members {
		m := &b.agg.Members[i]
		if !m.Single && (viewed || m.View == nil) {
			out = append(out, m)
		}
	}
	return out
}

// create adds an element, which is only possible when exactly one source
// list could take it.
func (v *Viewing) create(b Binding, rest Path, o Op, val reflect.Value) error {
	if o == OpClear && len(rest) == 1 {
		// removing an element that isn't there
		return nil
	}
	candidates := v.lists(b, true)
	if len(candidates) != 1 {
		panic(Errorf("die: adding %s%s is ambiguous: it could go to %s", b.view, rest[0], describeMembers(candidates)))
	}
	m := candidates[0]
	if m.View != nil {
		panic(Errorf("die: adding %s%s would go to %s, whose elements are presented by a jig, which can't create them", b.view, rest[0], m.Src))
	}
	at := v.At.join(m.Src)
	if len(v.Src.listKeys(at, v.keys(b))) > 0 {
		return v.writeSource(at.child(rest[0]).join(rest[1:]), o, val, TypeAt(v.Vtype, b.view.join(rest)))
	}
	// a new element of an atomic list is appended to it
	list, err := v.Src.valueAt(at, nil)
	if err != nil {
		return err
	}
	el := &Holder{Base: newElem(list.Type().Elem(), rest[0]), Mutable: true}
	if _, err := el.Write(nil, rest[1:], SetOp(o, val), o != OpClear); err != nil {
		return err
	}
	ev, err := el.Materialize()
	if err != nil {
		return err
	}
	out := reflect.Append(DeepCopy(list), adapt(ev, list.Type().Elem()))
	return v.writeSource(at, OpSet, out, nil)
}

func (v *Viewing) duplicate(b Binding, s Step, found []element) error {
	var ms []*Member
	for _, e := range found {
		ms = append(ms, e.member)
	}
	return Errorf("die: %s%s is ambiguous: it is in %s", b.view, s, describeMembers(ms))
}
