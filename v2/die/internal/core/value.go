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

// Package core is the engine of dies: the state a die holds, the overlay of
// intent, paths, schemas, the codec, patches, diffs and views' runtime. Its
// names are exported for the die packages, not for users.
package core

import (
	"reflect"
)

var ErrImmutableChild = Errorf("die: cannot mutate through a child of an immutable die, mutate inside DieWith instead")

// holder is the state of a root object. Child dies share their root's holder
// and address into it by path.
type Holder struct {
	// base is the fed or decoded value. It is never mutated.
	Base reflect.Value
	// ov is the intent recorded on top of base.
	Ov *Onode
	// wire is the presence of the base's values, when it was decoded, and the
	// fields its type doesn't have
	Wire    *pnode
	Mutable bool
	// view is set for a view, whose reads and writes go to its source
	View *Viewing
	// Schema is an unstructured value's structure, from a server's OpenAPI
	// schema, if it's known. It's part of the die, not its value: feeding
	// and copying the die keep it.
	Schema *Schema
	// root makes the root die over this holder, as the first die made at its
	// root, which knows its type, records it
	root func(*Holder) any
}

// SetRoot records how to make the root die over the holder, unless it's
// recorded.
func (h *Holder) SetRoot(root func(*Holder) any) {
	if h.root == nil {
		h.root = root
	}
}

// Root returns the root die over the holder, or nil if none was made.
func (h *Holder) Root() any {
	if h.root == nil {
		return nil
	}
	return h.root(h)
}

// valueAt returns a private copy of the value at p: base with the overlay
// applied. Only the overlay along p and below it is consulted.
func (h *Holder) ValueAt(p Path) (reflect.Value, error) {
	if h.View != nil {
		return h.View.valueAt(p)
	}
	if k := p.indexed(); k >= 0 {
		e, err := h.element(p[:k+1])
		if err != nil {
			return reflect.Value{}, err
		}
		if e != nil {
			return e.ValueAt(p[k+1:])
		}
		return reflect.Zero(TypeAt(h.Base.Type(), p)), nil
	}
	src := h.Base
	n := h.Ov
	for _, s := range p {
		src = readOp(src, n)
		var err error
		if src, err = descend(src, s); err != nil {
			return reflect.Value{}, err
		}
		n = n.kid(s)
	}
	src = readOp(src, n)
	out := DeepCopy(src)
	if n != nil {
		if err := applyKids(out, n); err != nil {
			return reflect.Value{}, err
		}
	}
	return out, nil
}

func (h *Holder) Has(p Path) (bool, error) {
	if h.View != nil {
		return h.View.has(p)
	}
	if k := p.indexed(); k >= 0 {
		e, err := h.element(p[:k+1])
		if err != nil || e == nil {
			return false, err
		}
		return e.Has(p[k+1:])
	}
	n := h.Ov
	for _, s := range p {
		if n == nil {
			break
		}
		if n.removes() {
			return false, nil
		}
		n = n.kid(s)
	}
	if n != nil {
		switch {
		case n.removes():
			return false, nil
		case n.Op == OpSet || n.Op == OpEmpty, creates(n):
			return true, nil
		}
	}
	if len(p) > 0 {
		last := p[len(p)-1]
		switch last.kind {
		case stepKey:
			l, err := h.ValueAt(p[:len(p)-1])
			if err != nil {
				return false, err
			}
			i, err := FindElem(l, last)
			return i >= 0, err
		case stepMapKey, stepField:
			m, err := h.ValueAt(p[:len(p)-1])
			if err != nil {
				return false, err
			}
			for m.Kind() == reflect.Interface && !m.IsNil() {
				m = m.Elem()
			}
			if m.Kind() == reflect.Map {
				if m.IsNil() {
					return false, nil
				}
				v := m.MapIndex(reflect.ValueOf(last.name).Convert(m.Type().Key()))
				return v.IsValid() && !(v.Kind() == reflect.Interface && v.IsNil()), nil
			}
		}
	}
	return h.wirePresent(p)
}

// write records intent at at+rel. Writes to a mutable holder are made in
// place. An immutable root returns a new holder, while an immutable child
// cannot be written.
func (h *Holder) Write(at, rel Path, fn func(*Onode) *Onode, union bool) (*Holder, error) {
	if h.View != nil {
		n := fn(&Onode{})
		return h, h.View.write(at.join(rel), n.Op, n.Val)
	}
	if !h.Mutable && len(at) > 0 {
		return nil, ErrImmutableChild
	}
	full := at.join(rel)
	if k := full.indexed(); k >= 0 {
		return h.writeElement(at, full, k, fn, union)
	}
	ov := h.Ov
	if union {
		var err error
		if ov, err = h.clearUnionSiblings(ov, full); err != nil {
			return nil, err
		}
	}
	if TypeAt(h.Base.Type(), full).Kind() == reflect.Interface {
		// values written into an unstructured value take unstructured form
		var convErr error
		ov = with(ov, full, func(n *Onode) *Onode {
			n = fn(n)
			if n.Op == OpSet && n.Val.IsValid() {
				u, err := ToUnstructured(n.Val)
				if err != nil {
					convErr = err
					return n
				}
				n.Val = reflect.ValueOf(u)
			}
			return n
		})
		if convErr != nil {
			return nil, convErr
		}
		return h.Replace(at, ov)
	}
	return h.Replace(at, with(ov, full, fn))
}

// list returns a private copy of the list at p, or an invalid value when
// there is none.
func (h *Holder) list(p Path) (reflect.Value, error) {
	l, err := h.ValueAt(p)
	if err != nil {
		return reflect.Value{}, err
	}
	for l.Kind() == reflect.Interface || l.Kind() == reflect.Pointer {
		if l.IsNil() {
			return reflect.Value{}, nil
		}
		l = l.Elem()
	}
	if l.Kind() != reflect.Slice {
		return reflect.Value{}, nil
	}
	return l, nil
}

// element returns a detached holder of the list element that p, ending in an
// index step, selects, or nil when the list has no such element.
func (h *Holder) element(p Path) (*Holder, error) {
	l, err := h.list(p[:len(p)-1])
	if err != nil {
		return nil, err
	}
	i := p[len(p)-1].index
	if !l.IsValid() || i >= l.Len() {
		return nil, nil
	}
	return &Holder{Base: l.Index(i), Mutable: true}, nil
}

// writeElement writes through the k-th step of full, an index step, by
// writing to a detached copy of the element and setting the whole list.
// Clearing the element itself removes it from the list.
func (h *Holder) writeElement(at, full Path, k int, fn func(*Onode) *Onode, union bool) (*Holder, error) {
	lp, i, rest := full[:k], full[k].index, full[k+1:]
	l, err := h.list(lp)
	if err != nil {
		return nil, err
	}
	if !l.IsValid() || i >= l.Len() {
		return nil, Errorf("die: %s is out of range", full[:k+1])
	}
	if len(rest) == 0 && fn(&Onode{}).removes() {
		out := reflect.MakeSlice(l.Type(), 0, l.Len()-1)
		l = reflect.AppendSlice(reflect.AppendSlice(out, l.Slice(0, i)), l.Slice(i+1, l.Len()))
	} else {
		e := &Holder{Base: l.Index(i), Mutable: true}
		if _, err := e.Write(nil, rest, fn, union); err != nil {
			return nil, err
		}
		v, err := e.Materialize()
		if err != nil {
			return nil, err
		}
		l.Index(i).Set(v)
	}
	// the list is written at the same depth of at, so that writing through
	// a child of an immutable die is still refused
	if len(at) > k {
		at = lp
	}
	return h.Write(at, lp[len(at):], SetOp(OpSet, l), false)
}

func (h *Holder) Replace(at Path, ov *Onode) (*Holder, error) {
	if h.View != nil {
		return h, h.View.replay(NodeAt(ov, at), at)
	}
	if h.Mutable {
		h.Ov = ov
		return h, nil
	}
	if len(at) > 0 {
		return nil, ErrImmutableChild
	}
	return &Holder{Base: h.Base, Ov: ov, Wire: h.Wire, Schema: h.Schema}, nil
}

// clearUnionSiblings clears the present members of every union the path
// writes into, so setting one member removes the others.
func (h *Holder) clearUnionSiblings(ov *Onode, p Path) (*Onode, error) {
	t := h.Base.Type()
	for i, s := range p {
		t = Deref(t)
		if t.Kind() != reflect.Struct {
			// unions are only declared on structs
			t = TypeAt(t, Path{s})
			continue
		}
		switch s.kind {
		case stepField:
			sd := Describe(t)
			fd := sd.Field(s.name)
			if fd.Union != "" {
				for _, sib := range sd.Fields {
					if sib == fd || sib.Union != fd.Union {
						continue
					}
					sp := p[:i].child(FieldStep(sib.Wire))
					has, err := (&Holder{Base: h.Base, Ov: ov}).Has(sp)
					if err != nil {
						return nil, err
					}
					if has {
						ov = with(ov, sp, SetOp(OpClear, reflect.Value{}))
					}
				}
			}
			t = fd.Typ
		default:
			t = t.Elem()
		}
	}
	return ov, nil
}

func (h *Holder) Materialize() (reflect.Value, error) {
	return h.ValueAt(nil)
}

func readOp(src reflect.Value, n *Onode) reflect.Value {
	if n == nil {
		return src
	}
	switch n.Op {
	case OpSet:
		if !n.Val.IsValid() {
			return reflect.Zero(src.Type())
		}
		return n.Val
	case OpClear:
		return reflect.Zero(src.Type())
	case OpEmpty:
		if n.Val.IsValid() {
			return n.Val
		}
		return emptyValue(src.Type())
	}
	return src
}

// descend reads one step down without copying. Missing values read as zero,
// and a missing list-map element reads as a zero element with its keys set.
//
// Within an unstructured value, a field step reads an object's key, and a key
// step reads the element of a list with those key fields.
func descend(src reflect.Value, s Step) (reflect.Value, error) {
	for src.Kind() == reflect.Pointer || src.Kind() == reflect.Interface {
		if src.IsNil() {
			if src.Kind() == reflect.Interface {
				return reflect.Zero(anyType), nil
			}
			src = reflect.Zero(src.Type().Elem())
		} else {
			src = src.Elem()
		}
	}
	if s.kind == stepField && src.Kind() == reflect.Map {
		s.kind = stepMapKey
	}
	switch s.kind {
	case stepKey:
		i, err := FindElem(src, s)
		if err != nil {
			return reflect.Value{}, err
		}
		if i >= 0 {
			return src.Index(i), nil
		}
		return newElem(src.Type().Elem(), s), nil
	case stepMapKey:
		if !src.IsNil() {
			if v := src.MapIndex(reflect.ValueOf(s.name).Convert(src.Type().Key())); v.IsValid() {
				return v, nil
			}
		}
		return reflect.Zero(src.Type().Elem()), nil
	}
	return src.FieldByIndex(Describe(src.Type()).Field(s.name).index), nil
}

func apply(dst reflect.Value, n *Onode) error {
	switch {
	case n.Op == OpSet && n.Val.IsValid():
		dst.Set(DeepCopy(n.Val))
	case n.Op == OpSet, n.Op == OpClear:
		dst.SetZero()
	case n.Op == OpEmpty && n.Val.IsValid():
		// an unstructured empty value says whether it is an object or a list
		dst.Set(DeepCopy(n.Val))
	case n.Op == OpEmpty:
		dst.Set(emptyValue(dst.Type()))
	}
	return applyKids(dst, n)
}

// applyKids applies the overlay below n onto dst, which must be a private,
// settable value.
func applyKids(dst reflect.Value, n *Onode) error {
	if len(n.Order) == 0 {
		return nil
	}
	if dst.Kind() == reflect.Interface {
		// an unstructured value: an object or a list, created when missing
		cur := dst.Elem()
		if !cur.IsValid() {
			if !creates(n) {
				return nil
			}
			cur = reflect.ValueOf(map[string]any{})
			if n.Steps[n.Order[0]].kind == stepKey {
				cur = reflect.ValueOf([]any{})
			}
		}
		tmp := reflect.New(cur.Type()).Elem()
		tmp.Set(cur)
		if err := applyKids(tmp, n); err != nil {
			return err
		}
		dst.Set(tmp)
		return nil
	}
	if dst.Kind() == reflect.Pointer {
		if dst.IsNil() {
			if !creates(n) {
				return nil
			}
			dst.Set(reflect.New(dst.Type().Elem()))
		}
		dst = dst.Elem()
	}
	for _, id := range n.Order {
		k, s := n.Kids[id], n.Steps[id]
		if s.kind == stepField && dst.Kind() == reflect.Map {
			s.kind = stepMapKey
		}
		switch s.kind {
		case stepField:
			if err := apply(dst.FieldByIndex(Describe(dst.Type()).Field(s.name).index), k); err != nil {
				return err
			}
		case stepKey:
			i, err := FindElem(dst, s)
			if err != nil {
				return err
			}
			if k.removes() {
				if i >= 0 {
					dst.Set(reflect.AppendSlice(dst.Slice(0, i), dst.Slice(i+1, dst.Len())))
				}
				continue
			}
			if i < 0 {
				dst.Set(reflect.Append(dst, newElem(dst.Type().Elem(), s)))
				i = dst.Len() - 1
			}
			if err := apply(dst.Index(i), k); err != nil {
				return err
			}
		case stepMapKey:
			key := reflect.ValueOf(s.name).Convert(dst.Type().Key())
			if k.removes() {
				if !dst.IsNil() {
					dst.SetMapIndex(key, reflect.Value{})
				}
				continue
			}
			if dst.IsNil() {
				dst.Set(reflect.MakeMap(dst.Type()))
			}
			v := reflect.New(dst.Type().Elem()).Elem()
			if cur := dst.MapIndex(key); cur.IsValid() {
				v.Set(cur)
			}
			if err := apply(v, k); err != nil {
				return err
			}
			dst.SetMapIndex(key, v)
		}
	}
	return nil
}

func creates(n *Onode) bool {
	for _, k := range n.Kids {
		if !k.removes() {
			return true
		}
	}
	return false
}

// findElem finds the element of a list with a key. A key that is not unique
// is ambiguous, rather than matching its first element.
func FindElem(list reflect.Value, s Step) (int, error) {
	for list.Kind() == reflect.Interface && !list.IsNil() {
		list = list.Elem()
	}
	if list.Kind() != reflect.Slice {
		return -1, nil
	}
	dynamic := list.Type().Elem().Kind() == reflect.Interface
	var sd *structDesc
	if !dynamic {
		sd = Describe(list.Type().Elem())
	}
	found := -1
	for i := range list.Len() {
		e := list.Index(i)
		match := true
		if dynamic {
			var err error
			if match, err = matchesDynamic(e.Interface(), s); err != nil {
				return -1, err
			}
		}
		for j, k := range s.keys {
			if dynamic {
				break
			}
			fd := sd.Field(k)
			if fd.keyValue(e.FieldByIndex(fd.index)).Interface() != s.vals[j] {
				match = false
				break
			}
		}
		if match {
			if found >= 0 {
				return -1, Errorf("die: %s is ambiguous: the list has it more than once", s)
			}
			found = i
		}
	}
	return found, nil
}

// withKeyDefaults returns a copy of a list-map element with unset key fields
// defaulted.
func WithKeyDefaults(e reflect.Value, keys []string) reflect.Value {
	out := DeepCopy(e)
	sd := Describe(e.Type())
	for _, k := range keys {
		fd := sd.Field(k)
		f := out.FieldByIndex(fd.index)
		f.Set(fd.keyValue(f))
	}
	return out
}

// matchesDynamic reports whether an unstructured list element has a key.
// With the element type of a view, unset key fields match their defaults.
func matchesDynamic(el any, s Step) (bool, error) {
	if s.elem != nil && Deref(s.elem).Kind() == reflect.Struct {
		return matchesUnstructured(el, s)
	}
	m, ok := el.(map[string]any)
	if !ok {
		return false, nil
	}
	for j, k := range s.keys {
		u, err := ToUnstructured(reflect.ValueOf(s.vals[j]))
		if err != nil {
			return false, err
		}
		if !reflect.DeepEqual(m[k], u) {
			return false, nil
		}
	}
	return true, nil
}

func newElem(t reflect.Type, s Step) reflect.Value {
	if t.Kind() == reflect.Interface {
		e := reflect.New(t).Elem()
		e.Set(reflect.ValueOf(keyedElement(s)))
		return e
	}
	e := reflect.New(t).Elem()
	sd := Describe(t)
	for j, k := range s.keys {
		e.FieldByIndex(sd.Field(k).index).Set(reflect.ValueOf(s.vals[j]))
	}
	return e
}

func emptyValue(t reflect.Type) reflect.Value {
	switch t.Kind() {
	case reflect.Slice:
		return reflect.MakeSlice(t, 0, 0)
	case reflect.Map:
		return reflect.MakeMap(t)
	case reflect.Pointer:
		return reflect.New(t.Elem())
	}
	return reflect.Zero(t)
}

// deepCopy returns a settable copy of v that shares no mutable state.
func DeepCopy(v reflect.Value) reflect.Value {
	out := reflect.New(v.Type()).Elem()
	copyInto(out, v)
	return out
}

func copyInto(dst, src reflect.Value) {
	switch src.Kind() {
	case reflect.Pointer:
		if src.IsNil() {
			dst.SetZero()
			return
		}
		p := reflect.New(src.Type().Elem())
		copyInto(p.Elem(), src.Elem())
		dst.Set(p)
	case reflect.Slice:
		if src.IsNil() {
			dst.SetZero()
			return
		}
		s := reflect.MakeSlice(src.Type(), src.Len(), src.Len())
		for i := range src.Len() {
			copyInto(s.Index(i), src.Index(i))
		}
		dst.Set(s)
	case reflect.Map:
		if src.IsNil() {
			dst.SetZero()
			return
		}
		m := reflect.MakeMapWithSize(src.Type(), src.Len())
		for it := src.MapRange(); it.Next(); {
			m.SetMapIndex(it.Key(), DeepCopy(it.Value()))
		}
		dst.Set(m)
	case reflect.Struct:
		// copies unexported state too, then deepens the exported fields
		dst.Set(src)
		for i := range src.NumField() {
			if dst.Field(i).CanSet() {
				copyInto(dst.Field(i), src.Field(i))
			}
		}
	case reflect.Interface:
		if src.IsNil() {
			dst.SetZero()
			return
		}
		dst.Set(DeepCopy(src.Elem()))
	default:
		dst.Set(src)
	}
}
