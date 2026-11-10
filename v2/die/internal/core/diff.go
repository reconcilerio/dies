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
)

// diffInto records the changes from before to after, both of type t at p, as
// intent. die.StampValue uses it so raw mutations are tracked like setter calls. A
// raw assignment that does not change the value cannot be seen, so explicit
// zero values must be written with setters.
func DiffInto(ov *Onode, p Path, t reflect.Type, before, after reflect.Value) (*Onode, error) {
	if reflect.DeepEqual(before.Interface(), after.Interface()) {
		return ov, nil
	}
	set := func() (*Onode, error) { return with(ov, p, SetOp(OpSet, DeepCopy(after))), nil }
	clear := func() (*Onode, error) { return with(ov, p, SetOp(OpClear, reflect.Value{})), nil }

	switch t.Kind() {
	case reflect.Pointer:
		if after.IsNil() {
			return clear()
		}
		if before.IsNil() || Deref(t).Kind() != reflect.Struct || customJSON(t.Elem()) {
			return set()
		}
		return DiffInto(ov, p, t.Elem(), before.Elem(), after.Elem())
	case reflect.Struct:
		if customJSON(t) {
			return set()
		}
		for _, fd := range Describe(t).Fields {
			fp := p.child(FieldStep(fd.Wire))
			bv, av := before.FieldByIndex(fd.index), after.FieldByIndex(fd.index)
			var err error
			switch fd.Kind {
			case KindStruct:
				ov, err = DiffInto(ov, fp, fd.Typ, bv, av)
			case KindListMap:
				ov, err = diffList(ov, fp, fd, bv, av)
			case KindMap:
				ov = diffMap(ov, fp, bv, av)
			default:
				if reflect.DeepEqual(bv.Interface(), av.Interface()) {
					continue
				}
				if nilable(av) && av.IsNil() {
					ov = with(ov, fp, SetOp(OpClear, reflect.Value{}))
				} else {
					ov = with(ov, fp, SetOp(OpSet, DeepCopy(av)))
				}
			}
			if err != nil {
				return nil, err
			}
		}
		return ov, nil
	}
	if nilable(after) && after.IsNil() {
		return clear()
	}
	return set()
}

func diffList(ov *Onode, p Path, fd *FieldDesc, before, after reflect.Value) (*Onode, error) {
	if after.IsNil() && !before.IsNil() {
		return with(ov, p, SetOp(OpClear, reflect.Value{})), nil
	}
	if before.IsNil() && after.IsNil() {
		return ov, nil
	}
	elem := fd.Typ.Elem()
	sd := Describe(elem)
	stepOf := func(e reflect.Value) Step {
		vals := make([]any, len(fd.Keys))
		for i, k := range fd.Keys {
			vals[i] = e.FieldByIndex(sd.Field(k).index).Interface()
		}
		return KeyStep(elem, fd.Keys, vals)
	}
	seen := map[string]bool{}
	for i := range after.Len() {
		ae := after.Index(i)
		s := stepOf(ae)
		seen[s.id()] = true
		j, err := FindElem(before, s)
		if err != nil {
			return nil, err
		}
		if j >= 0 {
			if ov, err = DiffInto(ov, p.child(s), elem, before.Index(j), ae); err != nil {
				return nil, err
			}
		} else {
			ov = with(ov, p.child(s), SetOp(OpSet, DeepCopy(ae)))
		}
	}
	for i := range before.Len() {
		if s := stepOf(before.Index(i)); !seen[s.id()] {
			ov = with(ov, p.child(s), SetOp(OpClear, reflect.Value{}))
		}
	}
	return ov, nil
}

func diffMap(ov *Onode, p Path, before, after reflect.Value) *Onode {
	if after.IsNil() && !before.IsNil() {
		return with(ov, p, SetOp(OpClear, reflect.Value{}))
	}
	for it := after.MapRange(); it.Next(); {
		bv := reflect.Value{}
		if !before.IsNil() {
			bv = before.MapIndex(it.Key())
		}
		if bv.IsValid() && reflect.DeepEqual(bv.Interface(), it.Value().Interface()) {
			continue
		}
		ov = with(ov, p.child(MapKeyStep(it.Key().String())), SetOp(OpSet, DeepCopy(it.Value())))
	}
	if !before.IsNil() {
		for it := before.MapRange(); it.Next(); {
			if after.IsNil() || !after.MapIndex(it.Key()).IsValid() {
				ov = with(ov, p.child(MapKeyStep(it.Key().String())), SetOp(OpClear, reflect.Value{}))
			}
		}
	}
	return ov
}

func nilable(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Pointer, reflect.Slice, reflect.Map, reflect.Interface:
		return true
	}
	return false
}
