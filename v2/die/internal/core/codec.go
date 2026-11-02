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
	"bytes"
	"encoding/json"
	"reflect"
)

// A Go value can't say which of its fields are present: a zero int is
// indistinguishable from an unset one. So the codec keeps presence beside
// the base, as a tree of the paths present on the wire, recorded when a die
// is decoded. Where there is no tree, as for a value fed or set from Go,
// presence follows the Go rule: a nil pointer, slice, map or interface is
// absent, as is an omitempty or omitzero field with an empty or zero value.
// Either way, an absent value is never encoded, not even as null.
//
// The tree also keeps the fields that a value had on the wire but its Go type
// doesn't, such as those of a newer API, so that encoding the die sends them
// back.

// pnode is a node of the presence tree. A node is present. Its kids are the
// present values below it, by step id, and extra is the fields its Go type
// doesn't have. A loose node is one whose values came from Go, not the wire:
// below it, a value with no node of its own follows the Go rule. Nodes are
// never modified after they are shared.
type pnode struct {
	kids  map[string]*pnode
	extra map[string]any
	loose bool
}

func (n *pnode) kid(s Step) *pnode {
	if n == nil {
		return nil
	}
	return n.kids[s.id()]
}

func (n *pnode) clone() *pnode {
	c := &pnode{loose: n.loose, extra: n.extra}
	if len(n.kids) > 0 {
		c.kids = make(map[string]*pnode, len(n.kids))
		for k, v := range n.kids {
			c.kids[k] = v
		}
	}
	return c
}

func (n *pnode) setKid(s Step, k *pnode) {
	if k == nil {
		delete(n.kids, s.id())
		return
	}
	if n.kids == nil {
		n.kids = map[string]*pnode{}
	}
	n.kids[s.id()] = k
}

// decodeWire records the presence of the values in data, decoded as t.
func DecodeWire(t reflect.Type, data []byte) (*pnode, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var raw any
	if err := dec.Decode(&raw); err != nil {
		return nil, err
	}
	d := &wireDecoder{}
	w := d.wireOf(t, raw)
	return w, d.err
}

// wireDecoder records presence, and the first error it meets: the presence
// tree is complete for what decoded, and the error says what didn't.
type wireDecoder struct{ err error }

// wireOf records the presence of a decoded value of type t. null is absent.
func (d *wireDecoder) wireOf(t reflect.Type, raw any) *pnode {
	t = Deref(t)
	m, ok := raw.(map[string]any)
	if !ok || t.Kind() != reflect.Struct || customJSON(t) {
		return &pnode{}
	}
	n := &pnode{}
	sd := Describe(t)
	for k, v := range m {
		if v == nil {
			continue
		}
		fd, ok := sd.byWire[k]
		if !ok {
			if n.extra == nil {
				n.extra = map[string]any{}
			}
			n.extra[k] = v
			continue
		}
		n.setKid(FieldStep(k), d.wireOfField(fd, v))
	}
	return n
}

func (d *wireDecoder) wireOfField(fd *FieldDesc, raw any) *pnode {
	switch fd.Kind {
	case KindStruct:
		return d.wireOf(fd.Typ, raw)
	case KindListMap:
		n := &pnode{}
		elem := Deref(fd.Typ).Elem()
		sd := Describe(elem)
		for _, e := range asList(raw) {
			m, _ := e.(map[string]any)
			vals := make([]any, len(fd.Keys))
			for i, k := range fd.Keys {
				kv := reflect.Zero(sd.Field(k).Typ)
				if v, ok := m[k]; ok && v != nil {
					var err error
					if kv, err = fromUnstructured(v, kv.Type()); err != nil {
						if d.err == nil {
							d.err = err
						}
						return n
					}
				}
				vals[i] = kv.Interface()
			}
			n.setKid(KeyStep(elem, fd.Keys, vals), d.wireOf(elem, e))
		}
		return n
	case KindMap:
		n := &pnode{}
		elem := Deref(fd.Typ).Elem()
		m, _ := raw.(map[string]any)
		for k, v := range m {
			if v != nil {
				n.setKid(MapKeyStep(k), d.wireOf(elem, v))
			}
		}
		return n
	}
	// an atomic list of structs keeps its elements' unknown fields, by index
	if t := Deref(fd.Typ); t.Kind() == reflect.Slice && structured(t.Elem()) {
		n := &pnode{}
		for i, e := range asList(raw) {
			n.setKid(IndexStep(i), d.wireOf(t.Elem(), e))
		}
		return n
	}
	return &pnode{}
}

func asList(raw any) []any {
	l, _ := raw.([]any)
	return l
}

// structured reports whether values of a type are encoded field by field.
func structured(t reflect.Type) bool {
	t = Deref(t)
	return t.Kind() == reflect.Struct && !customJSON(t)
}

// commitWire returns the presence once an overlay is folded into the base.
// What the overlay set came from Go, so it is loose; what it cleared is gone.
func CommitWire(w *pnode, n *Onode) *pnode {
	if n == nil {
		return w
	}
	if n.removes() {
		return nil
	}
	var out *pnode
	switch {
	case n.Op == OpSet || n.Op == OpEmpty:
		out = &pnode{loose: true}
	case w == nil:
		out = &pnode{loose: true}
	default:
		out = w.clone()
	}
	for _, id := range n.Order {
		s := n.Steps[id]
		var k *pnode
		if n.Op == opNone {
			k = w.kid(s)
		}
		out.setKid(s, CommitWire(k, n.Kids[id]))
	}
	return out
}

// goPresent reports whether a field's Go value is present by the Go rule.
func goPresent(fd *FieldDesc, v reflect.Value) bool {
	if !v.IsValid() {
		return false
	}
	switch v.Kind() {
	case reflect.Pointer, reflect.Interface, reflect.Map, reflect.Slice:
		if v.IsNil() {
			return false
		}
	}
	if !fd.omit {
		return true
	}
	switch v.Kind() {
	case reflect.Map, reflect.Slice, reflect.String, reflect.Array:
		return v.Len() > 0
	}
	return !v.IsZero()
}

// wirePresent reports whether the value at p of the base is present.
func (h *Holder) wirePresent(p Path) (bool, error) {
	n := h.Wire
	parent := h.Base.Type()
	for i, s := range p {
		if n == nil || (n.loose && n.kid(s) == nil) {
			// by the Go rule, from the parent of what's left
			present, err := h.goStepPresent(parent, p, i)
			if err != nil || !present {
				return false, err
			}
			n = nil
			continue
		}
		n = n.kid(s)
		if n == nil {
			return false, nil
		}
	}
	return true, nil
}

// goStepPresent reports whether the i-th step of p is present by the Go rule,
// given its parent is.
func (h *Holder) goStepPresent(root reflect.Type, p Path, i int) (bool, error) {
	s := p[i]
	parent, err := h.ValueAt(p[:i])
	if err != nil {
		return false, err
	}
	for parent.Kind() == reflect.Interface || parent.Kind() == reflect.Pointer {
		if parent.IsNil() {
			return false, nil
		}
		parent = parent.Elem()
	}
	switch {
	case s.kind == stepKey:
		i, err := FindElem(parent, s)
		return i >= 0, err
	case parent.Kind() == reflect.Map:
		if parent.IsNil() {
			return false, nil
		}
		v := parent.MapIndex(reflect.ValueOf(s.name).Convert(parent.Type().Key()))
		return v.IsValid() && !(v.Kind() == reflect.Interface && v.IsNil()), nil
	case parent.Kind() == reflect.Struct:
		fd := Describe(parent.Type()).Field(s.name)
		fv, _ := parent.FieldByIndexErr(fd.index)
		return goPresent(fd, fv), nil
	}
	v, err := h.ValueAt(p[:i+1])
	if err != nil {
		return false, err
	}
	return !v.IsZero(), nil
}

// encodeAt encodes the value at p as its presence says: base values as the
// wire had them, written values as Go has them, unknown fields kept, and
// absent values omitted.
func (h *Holder) EncodeAt(p Path) (any, error) {
	t := TypeAt(h.Base.Type(), p)
	v, err := h.ValueAt(p)
	if err != nil {
		return nil, err
	}
	if !structured(t) {
		return toJSON(v)
	}
	// a view has no presence of its own: its values follow the Go rule,
	// with its intent, which is its source's, as a written die's do
	n, w := IntentOf(h), h.Wire
	if h.View != nil {
		w = nil
	}
	for _, s := range p {
		if n != nil && (n.Op == OpSet || n.Op == OpEmpty) {
			w = nil
		}
		n, w = n.kid(s), w.kid(s)
	}
	if n != nil && (n.Op == OpSet || n.Op == OpEmpty) {
		w = nil
	}
	return encode(v, t, n, w)
}

// encode encodes a value of type t, with the intent n and the presence w
// that apply to it. A nil w is loose: its values follow the Go rule.
func encode(v reflect.Value, t reflect.Type, n *Onode, w *pnode) (any, error) {
	t = Deref(t)
	for v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface {
		if v.IsNil() {
			v = reflect.Zero(t)
			break
		}
		v = v.Elem()
	}
	if !structured(t) {
		return toJSON(v)
	}
	obj := map[string]any{}
	strict := w != nil && !w.loose
	for _, fd := range Describe(t).Fields {
		s := FieldStep(fd.Wire)
		k, kw := n.kid(s), w.kid(s)
		fv, _ := v.FieldByIndexErr(fd.index)
		switch {
		case k != nil && k.removes():
			continue
		case k != nil && (k.Op == OpSet || k.Op == OpEmpty):
			// written whole, from Go
			kw = nil
		case k != nil && creates(k):
		case strict:
			if kw == nil {
				continue
			}
		case kw == nil && !goPresent(fd, fv):
			continue
		}
		e, err := encodeField(fd, fv, k, kw)
		if err != nil {
			return nil, err
		}
		obj[fd.Wire] = e
	}
	if w != nil {
		for k, ev := range w.extra {
			if _, ok := obj[k]; !ok {
				obj[k] = ev
			}
		}
	}
	return obj, nil
}

func encodeField(fd *FieldDesc, v reflect.Value, n *Onode, w *pnode) (any, error) {
	if !v.IsValid() {
		return nil, nil
	}
	switch fd.Kind {
	case KindStruct:
		return encode(v, fd.Typ, n, w)
	case KindListMap:
		list := reflect.Indirect(v)
		arr := make([]any, 0, list.Len())
		elem := list.Type().Elem()
		sd := Describe(elem)
		for i := range list.Len() {
			e := list.Index(i)
			vals := make([]any, len(fd.Keys))
			for j, k := range fd.Keys {
				vals[j] = e.FieldByIndex(sd.Field(k).index).Interface()
			}
			s := KeyStep(elem, fd.Keys, vals)
			ev, err := encode(e, elem, n.kid(s), w.kid(s))
			if err != nil {
				return nil, err
			}
			arr = append(arr, ev)
		}
		return arr, nil
	case KindMap:
		m := reflect.Indirect(v)
		if !structured(m.Type().Elem()) {
			return toJSON(v)
		}
		obj := make(map[string]any, m.Len())
		for it := m.MapRange(); it.Next(); {
			s := MapKeyStep(it.Key().String())
			ev, err := encode(it.Value(), m.Type().Elem(), n.kid(s), w.kid(s))
			if err != nil {
				return nil, err
			}
			obj[s.name] = ev
		}
		return obj, nil
	}
	if list := reflect.Indirect(v); list.Kind() == reflect.Slice && structured(list.Type().Elem()) && n == nil && w != nil {
		arr := make([]any, list.Len())
		for i := range list.Len() {
			ev, err := encode(list.Index(i), list.Type().Elem(), nil, w.kid(IndexStep(i)))
			if err != nil {
				return nil, err
			}
			arr[i] = ev
		}
		return arr, nil
	}
	return toJSON(v)
}
