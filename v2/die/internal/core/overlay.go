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
	"maps"
	"reflect"
	"slices"
)

// op is the intent recorded at a node of the overlay.
type Op uint8

const (
	// opNone nodes only lead to deeper intent.
	opNone Op = iota
	// opSet replaces the value; deeper nodes apply on top of it.
	OpSet
	// opClear removes the value (nil, absent, list element removed).
	OpClear
	// opEmpty sets a present, empty value ([] or {}).
	OpEmpty
)

// onode is a node of the intent overlay: a persistent tree keyed by wire
// path. Nodes are never modified after they are shared, so dies can share
// overlays freely and a write only copies the nodes along its path.
type Onode struct {
	Op    Op
	Val   reflect.Value
	Kids  map[string]*Onode
	Steps map[string]Step
	Order []string
}

func (n *Onode) clone() *Onode {
	if n == nil {
		return &Onode{}
	}
	c := &Onode{Op: n.Op, Val: n.Val}
	if len(n.Order) > 0 {
		c.Kids = maps.Clone(n.Kids)
		c.Steps = maps.Clone(n.Steps)
		c.Order = slices.Clone(n.Order)
	}
	return c
}

func (n *Onode) kid(s Step) *Onode {
	if n == nil {
		return nil
	}
	return n.Kids[s.id()]
}

// setKid must only be called on a node that is not yet shared.
func (n *Onode) setKid(s Step, k *Onode) {
	id := s.id()
	if k == nil {
		if _, ok := n.Kids[id]; ok {
			delete(n.Kids, id)
			delete(n.Steps, id)
			n.Order = slices.DeleteFunc(n.Order, func(o string) bool { return o == id })
		}
		return
	}
	if n.Kids == nil {
		n.Kids, n.Steps = map[string]*Onode{}, map[string]Step{}
	}
	if _, ok := n.Kids[id]; !ok {
		n.Order = append(n.Order, id)
	}
	n.Kids[id], n.Steps[id] = k, s
}

// removes reports whether the node only removes its value.
func (n *Onode) removes() bool { return n.Op == OpClear && len(n.Order) == 0 }

// replaces reports whether the node replaces its whole value, as opposed to
// only carrying deeper intent.
func (n *Onode) replaces() bool { return n.Op != opNone && !n.removes() }

// with returns a copy of n with fn applied to a private copy of the node at p.
func with(n *Onode, p Path, fn func(*Onode) *Onode) *Onode {
	c := n.clone()
	if len(p) == 0 {
		return fn(c)
	}
	c.setKid(p[0], with(n.kid(p[0]), p[1:], fn))
	return c
}

func SetOp(o Op, v reflect.Value) func(*Onode) *Onode {
	return func(n *Onode) *Onode {
		n.Op, n.Val = o, v
		n.Kids, n.Steps, n.Order = nil, nil, nil
		return n
	}
}

func NodeAt(n *Onode, p Path) *Onode {
	for _, s := range p {
		n = n.kid(s)
	}
	return n
}

// only keeps the intent at and below the given paths.
func Only(n *Onode, paths []Path) *Onode {
	var out *Onode
	for _, p := range paths {
		if k := NodeAt(n, p); k != nil {
			out = with(out, p, func(*Onode) *Onode { return k })
		}
	}
	return out
}

// except drops the intent at and below the given paths.
func Except(n *Onode, paths []Path) *Onode {
	for _, p := range paths {
		if len(p) == 0 {
			return nil
		}
		if NodeAt(n, p) == nil {
			continue
		}
		last := p[len(p)-1]
		n = with(n, p[:len(p)-1], func(c *Onode) *Onode {
			c.setKid(last, nil)
			return c
		})
	}
	return n
}
