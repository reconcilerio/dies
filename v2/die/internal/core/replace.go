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

import "reflect"

// ReplaceAt replaces the value at at+rel with a copy of src's, a die or a
// handle of any root, as intent: what differs is written, as writes of each
// change would be, so the intent is the diff from the value to src's. A
// field src doesn't have is cleared, as are a map's keys and a list map's
// elements, and the value is cleared when src is absent. Nothing is written
// where they're the same. The value shares nothing with src's, so later
// writes to either aren't seen by the other.
func ReplaceAt(h *Holder, at, rel Path, src Located) error {
	present, after, err := ReadLocated(src, TypeAt(h.Base.Type(), at.join(rel)))
	if err != nil {
		return err
	}
	return ReplaceValueAt(h, at, rel, present, after)
}

// ReplaceValueAt replaces the value at at+rel with after, a value of its
// type that the caller doesn't keep, as intent, as [ReplaceAt] does: or
// clears it, unless present.
func ReplaceValueAt(h *Holder, at, rel Path, present bool, after reflect.Value) error {
	full := at.join(rel)
	t := TypeAt(h.Base.Type(), full)
	has, err := h.Has(full)
	if err != nil {
		return err
	}
	switch {
	case !present:
		if !has {
			return nil
		}
		_, err = h.Write(at, rel, SetOp(OpClear, reflect.Value{}), false)
		return err
	case !has, t.Kind() == reflect.Interface, FirstIndex(full) >= 0:
		// set whole: there's nothing to diff against, an unstructured
		// value takes unstructured form as it is written, and a list element
		// found by its index is written whole
		_, err = h.Write(at, rel, SetOp(OpSet, after), true)
		return err
	}
	before, err := h.ValueAt(full)
	if err != nil {
		return err
	}
	if reflect.DeepEqual(before.Interface(), after.Interface()) {
		return nil
	}
	ov, err := diffAt(h, full, t, before, after)
	if err != nil {
		return err
	}
	if h.View == nil {
		if err := CheckKeys(after); err != nil {
			return err
		}
		if err := CheckRekeyed(&Holder{Base: h.Base, Ov: ov}, full); err != nil {
			return err
		}
	}
	_, err = h.replace(at, full, ov)
	return err
}

// diffAt returns the holder's intent with the change at p from before to
// after. A field is diffed as part of its struct, as a view diffs its
// source, so that a map or list map field's change is recorded by key.
func diffAt(h *Holder, p Path, t reflect.Type, before, after reflect.Value) (*Onode, error) {
	ov := IntentOf(h)
	if h.View != nil {
		// a view's writes go to its source, replayed from this diff alone
		ov = nil
	}
	if len(p) == 0 || p[len(p)-1].kind != stepField {
		return DiffInto(ov, p, t, before, after)
	}
	pp := p[:len(p)-1]
	pt := TypeAt(h.Base.Type(), pp)
	if Deref(pt).Kind() != reflect.Struct {
		return DiffInto(ov, p, t, before, after)
	}
	pbefore, err := h.ValueAt(pp)
	if err != nil {
		return nil, err
	}
	pafter := DeepCopy(pbefore)
	target := pafter
	if target.Kind() == reflect.Pointer {
		if target.IsNil() {
			return DiffInto(ov, p, t, before, after)
		}
		target = target.Elem()
	}
	target.FieldByIndex(Describe(target.Type()).Field(p[len(p)-1].name).index).Set(after)
	return DiffInto(ov, pp, pt, pbefore, pafter)
}

// ReadLocated returns whether src, a die or a handle, is present and, if it
// is, a copy of its value, as a t: a value and a pointer to one are the same
// location's, as a die is of a struct field or a pointer to one, and a list
// is of any slice type of its elements, as a named list type is.
func ReadLocated(src Located, t reflect.Type) (bool, reflect.Value, error) {
	h, p, ok := StateOf(src)
	if !ok && HandleStateOf != nil {
		h, p, ok = HandleStateOf(src)
	}
	if !ok {
		return false, reflect.Value{}, Errorf("die: replacing from %s: it is not a die or a handle", src.DiePath())
	}
	has, err := h.Has(p)
	if err != nil || !has {
		return false, reflect.Value{}, err
	}
	v, err := h.ValueAt(p)
	if err != nil {
		return false, reflect.Value{}, err
	}
	switch {
	case v.Type() == t:
	case t.Kind() == reflect.Interface:
		// unstructured: what a value of any type is, written
	case v.Kind() == reflect.Pointer && v.Type().Elem() == t:
		if v.IsNil() {
			return false, reflect.Value{}, nil
		}
		v = v.Elem()
	case t.Kind() == reflect.Pointer && t.Elem() == v.Type():
		ptr := reflect.New(v.Type())
		ptr.Elem().Set(v)
		v = ptr
	case t.Kind() == reflect.Slice && v.Kind() == reflect.Slice && t.Elem() == v.Type().Elem():
		v = v.Convert(t)
	default:
		return false, reflect.Value{}, Errorf("die: replacing a %s from %s, a %s", t, src.DiePath(), v.Type())
	}
	return true, DeepCopy(v), nil
}
