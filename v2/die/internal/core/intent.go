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
	"encoding/json"
	"maps"
	"reflect"
	"slices"
)

// An apply configuration, such as client-go's, is what intent is: the values
// a writer has an opinion on, each present or not. These convert between the
// two, through JSON, which is the form apply configurations share with the
// types they configure.

// IntentJSON encodes the intent at a location as an apply configuration: on
// a root, what an apply sends, the object's identity included; within it, the
// intent below the location.
func IntentJSON(h *Holder, at Path) ([]byte, error) {
	if len(at) == 0 {
		return h.ApplyPatch()
	}
	v, err := encodeApply(h, TypeAt(h.Base.Type(), at), NodeAt(h.Ov, at), at, h.Schema.At(at))
	if err != nil {
		return nil, err
	}
	if v == nil || v == omitted {
		v = map[string]any{}
	}
	return json.Marshal(v)
}

// FeedIntent replaces the intent at a location with the values present in
// raw, an apply configuration in generic JSON form: each becomes intent, set
// as it is, over the base. Elements of list maps are set by their keys, and
// atomic lists whole. It returns the holder written, as Replace does.
func FeedIntent(h *Holder, at Path, raw any) (*Holder, error) {
	b := &intentBuilder{}
	n := b.intentOf(TypeAt(h.Base.Type(), at), raw)
	if b.err != nil {
		return nil, b.err
	}
	ov := with(h.Ov, at, func(*Onode) *Onode { return n })
	return h.Replace(at, ov)
}

// intentBuilder builds intent, and records the first value it can't
// convert, such as a string where the type has a number.
type intentBuilder struct{ err error }

// leaf is the intent that sets a value as it is.
func (b *intentBuilder) leaf(t reflect.Type, raw any) *Onode {
	n, err := setLeaf(t, raw)
	if err != nil && b.err == nil {
		b.err = err
	}
	return n
}

// intentOf is the intent that sets what raw has, as a value of type t.
func (b *intentBuilder) intentOf(t reflect.Type, raw any) *Onode {
	if raw == nil {
		return &Onode{}
	}
	if !structured(t) {
		return b.leaf(t, raw)
	}
	m, ok := raw.(map[string]any)
	if !ok {
		return b.leaf(t, raw)
	}
	n := &Onode{}
	sd := Describe(t)
	for _, k := range slices.Sorted(maps.Keys(m)) {
		fd, ok := sd.byWire[k]
		if !ok || m[k] == nil {
			continue
		}
		n.setKid(FieldStep(k), b.intentOfField(fd, m[k]))
	}
	return n
}

func (b *intentBuilder) intentOfField(fd *FieldDesc, raw any) *Onode {
	switch fd.Kind {
	case KindStruct:
		return b.intentOf(fd.Typ, raw)
	case KindListMap:
		list, _ := raw.([]any)
		elem := Deref(fd.Typ).Elem()
		sd := Describe(elem)
		n := &Onode{}
		for _, e := range list {
			m, _ := e.(map[string]any)
			vals := make([]any, len(fd.Keys))
			for i, k := range fd.Keys {
				kv := reflect.Zero(sd.Field(k).Typ)
				if v, ok := m[k]; ok && v != nil {
					var err error
					if kv, err = fromUnstructured(v, kv.Type()); err != nil {
						if b.err == nil {
							b.err = err
						}
						return n
					}
				}
				vals[i] = kv.Interface()
			}
			en := b.intentOf(elem, e)
			// an element's keys are set, so that it exists however few of
			// its other fields are
			for i, k := range fd.Keys {
				if en.kid(FieldStep(k)) == nil {
					en.setKid(FieldStep(k), SetOp(OpSet, reflect.ValueOf(vals[i]))(&Onode{}))
				}
			}
			n.setKid(KeyStep(elem, fd.Keys, vals), en)
		}
		return n
	case KindMap:
		m, _ := raw.(map[string]any)
		elem := Deref(fd.Typ).Elem()
		n := &Onode{}
		for _, k := range slices.Sorted(maps.Keys(m)) {
			if m[k] == nil {
				continue
			}
			if structured(elem) {
				n.setKid(MapKeyStep(k), b.intentOf(elem, m[k]))
			} else {
				n.setKid(MapKeyStep(k), b.leaf(elem, m[k]))
			}
		}
		return n
	}
	return b.leaf(fd.Typ, raw)
}

// setLeaf is the intent that sets a value as it is.
func setLeaf(t reflect.Type, raw any) (*Onode, error) {
	v := reflect.ValueOf(raw)
	if Deref(t).Kind() != reflect.Interface {
		var err error
		if v, err = fromUnstructured(raw, t); err != nil {
			return &Onode{}, err
		}
	}
	return SetOp(OpSet, v)(&Onode{}), nil
}
