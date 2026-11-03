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
	"slices"
)

// Binding binds a location of a view to a location of a source.
type Binding struct {
	view, src   Path
	vroot, root reflect.Type
	// tf transforms the source value into the view value, and back
	tf *Transform
	// agg aggregates the elements of several source locations
	agg *Aggregate
	// err is why the binding couldn't be made, reported where it's used
	err error
}

// RootOf returns the type of the root die a location is within.
func RootOf(l Located) (reflect.Type, error) {
	r, ok := l.(Rooted)
	if !ok {
		return nil, Errorf("die: %s is not a location within a die", l.DiePath())
	}
	return r.DieRootType(), nil
}

// BindingErr returns why a binding couldn't be made, if it couldn't.
func BindingErr(b Binding) error { return b.err }

// WithErr returns a binding that couldn't be made, for err, if it's set.
func WithErr(b Binding, err error) Binding {
	if err != nil && b.err == nil {
		b.err = err
	}
	return b
}

type Transform struct {
	// stype is the source type, for sources with no types of their own
	Stype reflect.Type
	To    func(s reflect.Value) reflect.Value
	From  func(v, s reflect.Value) reflect.Value
}

// NewBinding binds a location of a view to a location of a source, through
// a transform or an aggregate if either is given. The source's root is nil
// for a path in an unstructured object.
func NewBinding(view, src Path, vroot, root reflect.Type, tf *Transform, agg *Aggregate) Binding {
	return Binding{view: view, src: src, vroot: vroot, root: root, tf: tf, agg: agg}
}

// WithTransform returns a binding through a transform.
func WithTransform(b Binding, tf *Transform) Binding {
	b.tf = tf
	return b
}

// BindingInfo returns where a binding binds, and the types of the roots of
// its locations.
func BindingInfo(b Binding) (view, src Path, vroot, root reflect.Type) {
	return b.view, b.src, b.vroot, b.root
}

// wholeView is a view of the whole of a location in a source, as a type t.
func WholeView(src source, at Path, t reflect.Type) *Viewing {
	return &Viewing{Src: src, At: at, Binds: []Binding{{}}, Vtype: t}
}

type SourceAt struct {
	Src   source
	At    Path
	Binds []Binding
}

// viewing is the state of a view: where its source is, and how it is bound.
type Viewing struct {
	Src   source
	At    Path
	Binds []Binding
	Vtype reflect.Type
}

var ErrReadOnlyJig = Errorf("die: a jig of an immutable die is read-only")

// within finds the binding that contains a view path, the longest if they
// nest, with the rest of the path below it.
func (v *Viewing) within(p Path) (Binding, Path, bool) {
	var best Binding
	found := false
	for _, b := range v.Binds {
		if p.hasPrefix(b.view) && (!found || len(b.view) > len(best.view)) {
			best, found = b, true
		}
	}
	if !found {
		return Binding{}, nil, false
	}
	return best, p[len(best.view):], true
}

// below finds the bindings for locations under a view path.
func (v *Viewing) below(p Path) []Binding {
	var out []Binding
	for _, b := range v.Binds {
		if len(b.view) > len(p) && b.view.hasPrefix(p) {
			out = append(out, b)
		}
	}
	return out
}

func (v *Viewing) source(b Binding, rest Path) Path { return v.At.join(b.src).join(rest) }

func (v *Viewing) valueAt(p Path) (reflect.Value, error) {
	if b, rest, ok := v.within(p); ok {
		return v.read(b, rest)
	}
	// a location above the bindings is assembled from them
	var ov *Onode
	for _, b := range v.below(p) {
		val, err := v.read(b, nil)
		if err != nil {
			return reflect.Value{}, err
		}
		ov = with(ov, b.view[len(p):], SetOp(OpSet, val))
	}
	out := reflect.New(TypeAt(v.Vtype, p)).Elem()
	if ov != nil {
		if err := applyKids(out, ov); err != nil {
			return reflect.Value{}, err
		}
	}
	return out, nil
}

// read returns the view value at a path within a binding.
func (v *Viewing) read(b Binding, rest Path) (reflect.Value, error) {
	if b.agg != nil {
		return v.readAggregate(b, rest)
	}
	t := TypeAt(v.Vtype, b.view.join(rest))
	if b.tf == nil {
		return v.Src.valueAt(v.source(b, rest), t)
	}
	sv, err := v.sourceValue(b)
	if err != nil {
		return reflect.Value{}, err
	}
	val := adapt(b.tf.To(sv), TypeAt(v.Vtype, b.view))
	for _, s := range rest {
		if val, err = descend(val, s); err != nil {
			return reflect.Value{}, err
		}
	}
	return adapt(DeepCopy(val), t), nil
}

// sourceValue is the source value of a binding with a transform.
func (v *Viewing) sourceValue(b Binding) (reflect.Value, error) {
	sp := v.source(b, nil)
	t := v.Src.typeAt(sp)
	if t == nil {
		t = b.tf.Stype
	}
	return v.Src.valueAt(sp, t)
}

func (v *Viewing) has(p Path) (bool, error) {
	if b, rest, ok := v.within(p); ok {
		if b.agg != nil {
			return v.hasAggregate(b, rest)
		}
		if b.tf != nil && len(rest) > 0 {
			val, err := v.read(b, rest)
			return err == nil && !val.IsZero(), err
		}
		return v.Src.has(v.source(b, rest))
	}
	for _, b := range v.below(p) {
		has, err := v.Src.has(v.source(b, nil))
		if err != nil || has {
			return has, err
		}
	}
	return false, nil
}

func (v *Viewing) write(p Path, o Op, val reflect.Value) error {
	if !v.Src.Mutable() {
		return ErrReadOnlyJig
	}
	if b, rest, ok := v.within(p); ok {
		return v.writeBinding(b, rest, o, val)
	}
	// a write above the bindings is split across them
	bs := v.below(p)
	if len(bs) == 0 {
		return Errorf("die: %s is not bound in this jig of %s", p, v.Src)
	}
	for _, b := range bs {
		sub := val
		if o == OpSet {
			for _, s := range b.view[len(p):] {
				var err error
				if sub, err = descend(sub, s); err != nil {
					return err
				}
			}
		}
		if err := v.writeBinding(b, nil, o, sub); err != nil {
			return err
		}
	}
	return nil
}

func (v *Viewing) writeBinding(b Binding, rest Path, o Op, val reflect.Value) error {
	if b.agg != nil {
		return v.writeAggregate(b, rest, o, val)
	}
	if b.tf == nil {
		return v.writeSource(v.source(b, rest), o, val, TypeAt(v.Vtype, b.view.join(rest)))
	}
	// apply the write to the view value, transform it back, and record only
	// what changed in the source as intent
	sp := v.source(b, nil)
	before, err := v.sourceValue(b)
	if err != nil {
		return err
	}
	st := before.Type()
	view := DeepCopy(adapt(b.tf.To(before), TypeAt(v.Vtype, b.view)))
	if err := apply(view, with(nil, rest, SetOp(o, val))); err != nil {
		return err
	}
	after := adapt(b.tf.From(view, before), st)
	if got := adapt(b.tf.To(after), view.Type()); !reflect.DeepEqual(got.Interface(), view.Interface()) {
		return Errorf("die: the transform for %s cannot represent %+v in its source; it would read back as %+v", b.view, view.Interface(), got.Interface())
	}
	n, err := v.diffSource(sp, st, before, after)
	if err != nil {
		return err
	}
	return v.replaySource(n, nil)
}

// diffSource records the change to a source value as intent. A field is
// diffed as part of its struct, so that a change within a list map is
// recorded on its elements rather than as the whole list.
func (v *Viewing) diffSource(sp Path, st reflect.Type, before, after reflect.Value) (*Onode, error) {
	if len(sp) == 0 || sp[len(sp)-1].kind != stepField || v.Src.typeAt(sp) == nil {
		return DiffInto(nil, sp, st, before, after)
	}
	pp := sp[:len(sp)-1]
	pt := v.Src.typeAt(pp)
	pbefore, err := v.Src.valueAt(pp, nil)
	if err != nil {
		return nil, err
	}
	pafter := DeepCopy(pbefore)
	target := pafter
	if target.Kind() == reflect.Pointer {
		if target.IsNil() {
			return DiffInto(nil, sp, st, before, after)
		}
		target = target.Elem()
	}
	target.FieldByIndex(Describe(target.Type()).Field(sp[len(sp)-1].name).index).Set(after)
	return DiffInto(nil, pp, pt, pbefore, pafter)
}

func (v *Viewing) writeSource(sp Path, o Op, val reflect.Value, t reflect.Type) error {
	return v.Src.write(sp, o, val, t)
}

// replaySource writes each operation of an overlay to the source.
func (v *Viewing) replaySource(n *Onode, p Path) error {
	if n == nil {
		return nil
	}
	if n.Op != opNone {
		if err := v.writeSource(p, n.Op, n.Val, nil); err != nil {
			return err
		}
	}
	for _, id := range n.Order {
		if err := v.replaySource(n.Kids[id], p.child(n.Steps[id])); err != nil {
			return err
		}
	}
	return nil
}

// replay writes each operation of an overlay, such as the changes die.StampValue
// finds, through the view.
func (v *Viewing) replay(n *Onode, p Path) error {
	if n == nil {
		return nil
	}
	if n.Op != opNone {
		if err := v.write(p, n.Op, n.Val); err != nil {
			return err
		}
	}
	for _, id := range n.Order {
		if err := v.replay(n.Kids[id], p.child(n.Steps[id])); err != nil {
			return err
		}
	}
	return nil
}

// adapt converts between a value and a pointer to it, as a view and its
// source may hold the same type either way.
func adapt(val reflect.Value, t reflect.Type) reflect.Value {
	switch {
	case !val.IsValid() || val.Type() == t:
		return val
	case val.Kind() == reflect.Pointer && val.Type().Elem() == t:
		if val.IsNil() {
			return reflect.Zero(t)
		}
		return val.Elem()
	case t.Kind() == reflect.Pointer && t.Elem() == val.Type():
		p := reflect.New(val.Type())
		p.Elem().Set(val)
		return p
	}
	return val
}

func (p Path) hasPrefix(q Path) bool {
	if len(q) > len(p) {
		return false
	}
	for i := range q {
		if p[i].id() != q[i].id() {
			return false
		}
	}
	return true
}

// IntentOf returns a holder's intent. A view records none of its own, as its
// writes go to its source: its intent is its source's, through its bindings.
func IntentOf(h *Holder) *Onode {
	if h.View != nil {
		return h.View.intent()
	}
	return h.Ov
}

// intent is the view's intent, derived from its source's through its
// bindings. Through a binding without a transform or an aggregate, the
// source's intent at the bound location is the view's at its own, operation
// for operation. A transform or an aggregate can't map operations, so where
// the source has any intent below such a binding, the view's intent sets the
// binding's whole value, as the view reads it. A binding nested in another
// replaces the outer binding's intent below it, as its reads do.
func (v *Viewing) intent() *Onode {
	src := v.Src.intent()
	if src == nil {
		return nil
	}
	binds := slices.Clone(v.Binds)
	slices.SortStableFunc(binds, func(a, b Binding) int { return len(a.view) - len(b.view) })
	var out *Onode
	for _, b := range binds {
		var n *Onode
		switch {
		case b.agg != nil:
			for _, m := range b.agg.Members {
				if NodeAt(src, v.At.join(m.Src)) != nil {
					n = v.setsBinding(b)
					break
				}
			}
		case b.tf != nil:
			if NodeAt(src, v.source(b, nil)) != nil {
				n = v.setsBinding(b)
			}
		default:
			n = NodeAt(src, v.source(b, nil))
		}
		out = graft(out, b.view, n)
	}
	return out
}

// setsBinding is intent that sets a binding's whole value as the view reads
// it, or that sets nothing if it can't be read.
func (v *Viewing) setsBinding(b Binding) *Onode {
	val, err := v.read(b, nil)
	if err != nil {
		return nil
	}
	return &Onode{Op: OpSet, Val: val}
}

// graft returns a copy of n with the node at p replaced by k, or removed if k
// is nil.
func graft(n *Onode, p Path, k *Onode) *Onode {
	if len(p) == 0 {
		return k
	}
	if k == nil && NodeAt(n, p) == nil {
		return n
	}
	last := p[len(p)-1]
	return with(n, p[:len(p)-1], func(c *Onode) *Onode {
		c.setKid(last, k)
		return c
	})
}
