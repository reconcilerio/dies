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
	"fmt"
	"reflect"
	"strings"
)

// A list map's elements are found by their keys, so two elements with the
// same key can't be told apart: reading or writing either is ambiguous. Such
// content is refused where it enters a die, when it is fed, decoded or
// written, so that a die always has unique keys and its reads can't fail.

// CheckKeys returns an error if a value of a die's type has a list map, at
// any depth, with two elements of the same key. A key left unset and the
// same key set to its default are the same key.
func CheckKeys(v reflect.Value) error {
	return checkKeys(v, nil)
}

func checkKeys(v reflect.Value, p Path) error {
	for v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct || customJSON(v.Type()) {
		return nil
	}
	for _, fd := range Describe(v.Type()).Fields {
		f := v.FieldByIndex(fd.index)
		switch fd.Kind {
		case KindStruct:
			if err := checkKeys(f, p.child(FieldStep(fd.Wire))); err != nil {
				return err
			}
		case KindListMap:
			if err := checkList(f, fd.Keys, p.child(FieldStep(fd.Wire))); err != nil {
				return err
			}
		case KindMap:
			if f.IsNil() || Deref(f.Type().Elem()).Kind() != reflect.Struct {
				continue
			}
			iter := f.MapRange()
			for iter.Next() {
				if err := checkKeys(iter.Value(), p.child(MapKeyStep(iter.Key().String()))); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// checkList checks a list map's keys, and its elements' content.
func checkList(l reflect.Value, keys []string, p Path) error {
	if l.Kind() != reflect.Slice || l.Len() == 0 {
		return nil
	}
	elem := Deref(l.Type().Elem())
	sd := Describe(elem)
	seen := map[string]bool{}
	for i := range l.Len() {
		e := l.Index(i)
		for e.Kind() == reflect.Pointer {
			e = e.Elem()
		}
		if !e.IsValid() {
			continue
		}
		vals := make([]any, len(keys))
		var id strings.Builder
		for j, k := range keys {
			fd := sd.Field(k)
			vals[j] = fd.keyValue(e.FieldByIndex(fd.index)).Interface()
			fmt.Fprintf(&id, "%#v\x00", vals[j])
		}
		if seen[id.String()] {
			return Errorf("die: %s has more than one element %s", p, KeyStep(elem, keys, vals))
		}
		seen[id.String()] = true
		if err := checkKeys(e, p.child(KeyStep(elem, keys, vals))); err != nil {
			return err
		}
	}
	return nil
}

// checkKeysAt checks a value written at a path within a root of type root:
// its content, and when it is a list map itself, its keys.
func checkKeysAt(root reflect.Type, at Path, v reflect.Value) error {
	if !v.IsValid() {
		return nil
	}
	if fd := listFieldAt(root, at); fd != nil && fd.Kind == KindListMap {
		return checkList(v, fd.Keys, at)
	}
	return checkKeys(v, at)
}

// listFieldAt returns the field a path ends at, or nil if it ends elsewhere,
// such as at a list element.
func listFieldAt(root reflect.Type, at Path) *FieldDesc {
	if len(at) == 0 || at[len(at)-1].kind != stepField {
		return nil
	}
	parent := Deref(TypeAt(root, at[:len(at)-1]))
	if parent.Kind() != reflect.Struct {
		return nil
	}
	fd, ok := Describe(parent).byWire[at[len(at)-1].name]
	if !ok {
		return nil
	}
	return fd
}

// CheckRekeyed checks the lists a write within list-map elements could have
// given a duplicate key, as a holder would have them after the write: a
// write to an element whole, or to one of its key fields, can make its key
// another element's.
func CheckRekeyed(h *Holder, full Path) error {
	if !full.keyed() {
		return nil
	}
	return checkRekeyed(h, full)
}

func checkRekeyed(h *Holder, full Path) error {
	for j, s := range full {
		if s.kind != stepKey {
			continue
		}
		if j+1 < len(full) && !s.hasKey(full[j+1]) {
			continue
		}
		l, err := h.ValueAt(full[:j])
		if err != nil {
			return err
		}
		fd := listFieldAt(h.Base.Type(), full[:j])
		if fd == nil || fd.Kind != KindListMap {
			continue
		}
		if err := checkList(l, fd.Keys, full[:j]); err != nil {
			return err
		}
	}
	return nil
}

// hasKey reports whether a step names one of a list-map step's key fields.
func (s Step) hasKey(next Step) bool {
	if next.kind != stepField {
		return false
	}
	for _, k := range s.keys {
		if k == next.name {
			return true
		}
	}
	return false
}
