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
	"slices"
	"strings"

	"k8s.io/client-go/util/jsonpath"
)

// JSONPath expressions, as kubectl and client-go write them
// (https://kubernetes.io/docs/reference/kubectl/jsonpath/), resolve to
// locations within a die, as accessors find them, rather than to copies of
// values:
//
//   - an element of a list map is found by its keys, so containers[0] and
//     containers[?(@.name=='app')] are both .containers[name=app], and a write
//     through either merges by key
//   - an element of an atomic list is found by its index
//   - in an unstructured list, a filter that matches one element by one field
//     with == finds it by that field, as Item does; other elements, by index
//   - a definite path, with no wildcard, slice, filter or recursive descent
//     before a step, resolves by the schema whether or not the value
//     is set, so that setting it creates it; other steps only match values
//     that are present
//   - in a typed die, a field the schema doesn't have is an error
//
// $ is the root die, and an expression that starts with it resolves from
// there; any other expression, such as .name, ..name or [n], resolves from the
// die queried. Either way, the paths found are within the root, as every path
// is.
//
// Expressions are parsed by client-go: one expression, with or without
// {braces}. Filters compare a field with a literal: ==, !=, <, >, <= and >=
// (on numbers), or test that it exists.

type jpSel struct {
	name     string
	wildcard bool
	array    *[3]jsonpath.ParamsEntry
	filter   *jpFilter
	union    [][]jpSeg
}

type jpFilter struct {
	rel []string
	op  string
	lit any
}

type jpSeg struct {
	recursive bool
	sel       jpSel
}

// parseJSONPath parses an expression, and reports whether it starts at the
// root.
func parseJSONPath(expr string) ([]jpSeg, bool, error) {
	s := strings.TrimSpace(expr)
	if !strings.Contains(s, "{") {
		s = "{" + s + "}"
	}
	inner := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(s), "{"))
	// client-go's parser drops a leading $, which means the root here
	fromRoot := strings.HasPrefix(inner, "$")
	p, err := jsonpath.Parse("die", s)
	if err != nil {
		return nil, false, Errorf("die: JSONPath %q: %w", expr, err)
	}
	var lists []*jsonpath.ListNode
	for _, n := range p.Root.Nodes {
		switch n := n.(type) {
		case *jsonpath.ListNode:
			lists = append(lists, n)
		case *jsonpath.TextNode:
			if strings.TrimSpace(n.Text) != "" {
				return nil, false, Errorf("die: JSONPath %q: text outside of an expression, %q", expr, n.Text)
			}
		default:
			return nil, false, Errorf("die: JSONPath %q: unsupported %s", expr, n)
		}
	}
	var segs []jpSeg
	for _, l := range lists {
		if segs, err = segmentsOf(l); err != nil {
			return nil, false, Errorf("die: JSONPath %q: %w", expr, err)
		}
	}
	if len(lists) != 1 {
		return nil, false, Errorf("die: JSONPath %q: want one expression, not %d", expr, len(lists))
	}
	return segs, fromRoot, nil
}

func segmentsOf(list *jsonpath.ListNode) ([]jpSeg, error) {
	var segs []jpSeg
	recursive := false
	for _, n := range list.Nodes {
		var sel jpSel
		switch n := n.(type) {
		case *jsonpath.RecursiveNode:
			recursive = true
			continue
		case *jsonpath.FieldNode:
			if n.Value == "" {
				// the current node, as in {.}
				continue
			}
			sel.name = n.Value
		case *jsonpath.WildcardNode:
			sel.wildcard = true
		case *jsonpath.ArrayNode:
			if !n.Params[0].Known && !n.Params[1].Known && !n.Params[2].Known {
				sel.wildcard = true
			} else {
				params := n.Params
				sel.array = &params
			}
		case *jsonpath.FilterNode:
			f, err := filterOf(n)
			if err != nil {
				return nil, err
			}
			sel.filter = f
		case *jsonpath.UnionNode:
			for _, alt := range n.Nodes {
				segs, err := segmentsOf(alt)
				if err != nil {
					return nil, err
				}
				sel.union = append(sel.union, segs)
			}
		default:
			return nil, fmt.Errorf("unsupported %s", n)
		}
		segs = append(segs, jpSeg{recursive: recursive, sel: sel})
		recursive = false
	}
	if recursive {
		// a trailing .. is every descendant
		segs = append(segs, jpSeg{recursive: true, sel: jpSel{wildcard: true}})
	}
	return segs, nil
}

func filterOf(n *jsonpath.FilterNode) (*jpFilter, error) {
	f := &jpFilter{op: n.Operator}
	for _, l := range n.Left.Nodes {
		fn, ok := l.(*jsonpath.FieldNode)
		if !ok {
			return nil, fmt.Errorf("a filter compares a field, not %s", l)
		}
		if fn.Value != "" {
			f.rel = append(f.rel, fn.Value)
		}
	}
	if len(f.rel) == 0 {
		return nil, fmt.Errorf("a filter must name a field")
	}
	if f.op == "exists" {
		return f, nil
	}
	if n.Right == nil || len(n.Right.Nodes) != 1 {
		return nil, fmt.Errorf("a filter compares with one literal")
	}
	switch r := n.Right.Nodes[0].(type) {
	case *jsonpath.TextNode:
		f.lit = r.Text
	case *jsonpath.IntNode:
		f.lit = int64(r.Value)
	case *jsonpath.FloatNode:
		f.lit = r.Value
	case *jsonpath.BoolNode:
		f.lit = r.Value
	default:
		return nil, fmt.Errorf("a filter compares with a literal, not %s", r)
	}
	switch f.op {
	case "==", "!=", "<", ">", "<=", ">=":
	case "=":
		f.op = "=="
	default:
		return nil, fmt.Errorf("unsupported filter operator %q", f.op)
	}
	return f, nil
}

// jpNode is a match: where it is, its current value and its static type, and
// for a list, the field it is, whose keys find its elements.
type jpNode struct {
	at Path
	v  reflect.Value
	t  reflect.Type
	fd *FieldDesc
}

// ResolveJSONPath returns the locations within a holder, from at, that a
// JSONPath expression matches.
func ResolveJSONPath(h *Holder, at Path, expr string) ([]Path, error) {
	segs, fromRoot, err := parseJSONPath(expr)
	if err != nil {
		return nil, err
	}
	if fromRoot {
		at = nil
	}
	v, err := h.ValueAt(at)
	if err != nil {
		return nil, err
	}
	start := []jpNode{{at: at, v: v, t: TypeAt(h.Base.Type(), at)}}
	nodes, err := evaluate(h, start, segs, true, expr)
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	var out []Path
	for _, n := range nodes {
		if k := n.at.String(); !seen[k] {
			seen[k] = true
			out = append(out, n.at)
		}
	}
	return out, nil
}

func evaluate(h *Holder, nodes []jpNode, segs []jpSeg, definite bool, expr string) ([]jpNode, error) {
	for _, seg := range segs {
		// a union's alternatives are as definite as their own steps
		if seg.recursive || seg.sel.wildcard || seg.sel.filter != nil || (seg.sel.array != nil && !single(*seg.sel.array)) {
			definite = false
		}
		var next []jpNode
		for _, n := range nodes {
			from := []jpNode{n}
			if seg.recursive {
				var err error
				if from, err = descendants(h, n); err != nil {
					return nil, err
				}
			}
			for _, f := range from {
				if seg.sel.union != nil {
					for _, alt := range seg.sel.union {
						m, err := evaluate(h, []jpNode{f}, alt, definite, expr)
						if err != nil {
							return nil, err
						}
						next = append(next, m...)
					}
					continue
				}
				m, err := selectJSONPath(h, f, seg.sel, definite, expr)
				if err != nil {
					return nil, err
				}
				next = append(next, m...)
			}
		}
		nodes = next
	}
	return nodes, nil
}

// single reports whether array params select one index.
func single(p [3]jsonpath.ParamsEntry) bool { return p[0].Known && p[1].Derived }

// indexes are the indexes array params select in a list of length n, as
// client-go selects them.
func indexes(p [3]jsonpath.ParamsEntry, n int) []int {
	if single(p) {
		i := p[0].Value
		if i < 0 {
			i += n
		}
		if i < 0 || i >= n {
			return nil
		}
		return []int{i}
	}
	start, end, step := 0, n, 1
	if p[0].Known {
		start = p[0].Value
	}
	if p[1].Known {
		end = p[1].Value
	}
	if p[2].Known && p[2].Value > 0 {
		step = p[2].Value
	}
	if start < 0 {
		start += n
	}
	if end < 0 {
		end += n
	}
	start, end = max(0, min(start, n)), max(0, min(end, n))
	var out []int
	for i := start; i < end; i += step {
		out = append(out, i)
	}
	return out
}

// settle dereferences a value, to a zero value of its type when nil.
func settle(v reflect.Value, t reflect.Type) (reflect.Value, reflect.Type) {
	t = Deref(t)
	for v.IsValid() && (v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface) {
		if v.IsNil() {
			if v.Kind() == reflect.Pointer {
				return reflect.Zero(t), t
			}
			return reflect.Value{}, t
		}
		v = v.Elem()
	}
	return v, t
}

// dynamic reports whether a type is unstructured: any, or an object of any.
func dynamic(t reflect.Type) bool {
	t = Deref(t)
	return t.Kind() == reflect.Interface || (t.Kind() == reflect.Map && t.Elem().Kind() == reflect.Interface)
}

func selectJSONPath(h *Holder, n jpNode, sel jpSel, definite bool, expr string) ([]jpNode, error) {
	v, t := settle(n.v, n.t)
	switch {
	case sel.name != "":
		c, ok, err := childNamed(h, n, v, t, sel.name, definite)
		if err != nil {
			return nil, Errorf("die: JSONPath %q: %w", expr, err)
		}
		if !ok {
			return nil, nil
		}
		return []jpNode{c}, nil
	case sel.wildcard:
		return children(h, n)
	case sel.array != nil:
		if !v.IsValid() || v.Kind() != reflect.Slice {
			return nil, nil
		}
		var out []jpNode
		for _, i := range indexes(*sel.array, v.Len()) {
			out = append(out, jpElement(n, v, t, i, nil))
		}
		return out, nil
	case sel.filter != nil:
		if !v.IsValid() || v.Kind() != reflect.Slice {
			return nil, nil
		}
		var idx []int
		for i := range v.Len() {
			match, err := sel.filter.matches(v.Index(i))
			if err != nil {
				return nil, err
			}
			if match {
				idx = append(idx, i)
			}
		}
		var out []jpNode
		for _, i := range idx {
			var by *jpFilter
			if len(idx) == 1 {
				by = sel.filter
			}
			out = append(out, jpElement(n, v, t, i, by))
		}
		return out, nil
	}
	return nil, nil
}

// childNamed is a field of a struct, or an entry of an object or map.
func childNamed(h *Holder, n jpNode, v reflect.Value, t reflect.Type, name string, definite bool) (jpNode, bool, error) {
	switch {
	case t.Kind() == reflect.Struct && structured(t):
		fd, ok := Describe(t).byWire[name]
		if !ok {
			if definite {
				return jpNode{}, false, fmt.Errorf("%s has no field %q", t, name)
			}
			return jpNode{}, false, nil
		}
		at := n.at.child(FieldStep(name))
		if !definite {
			has, err := h.Has(at)
			if err != nil || !has {
				return jpNode{}, false, err
			}
		}
		var fv reflect.Value
		if v.IsValid() {
			fv, _ = v.FieldByIndexErr(fd.index)
		}
		return jpNode{at: at, v: fv, t: fd.Typ, fd: fd}, true, nil
	case dynamic(t):
		var m map[string]any
		if v.IsValid() {
			m, _ = v.Interface().(map[string]any)
		}
		e, present := m[name]
		// an absent entry is found in a definite path, so that setting it
		// creates it, unless the value is not an object
		if !present && (!definite || (v.IsValid() && m == nil)) {
			return jpNode{}, false, nil
		}
		return jpNode{at: n.at.child(FieldStep(name)), v: reflect.ValueOf(e), t: anyType}, true, nil
	case t.Kind() == reflect.Map && t.Key().Kind() == reflect.String:
		at := n.at.child(MapKeyStep(name))
		var e reflect.Value
		if v.IsValid() && !v.IsNil() {
			e = v.MapIndex(reflect.ValueOf(name).Convert(t.Key()))
		}
		if !e.IsValid() {
			if !definite {
				return jpNode{}, false, nil
			}
			e = reflect.Zero(t.Elem())
		}
		return jpNode{at: at, v: e, t: t.Elem()}, true, nil
	}
	return jpNode{}, false, nil
}

// children are the present values directly within a value.
func children(h *Holder, n jpNode) ([]jpNode, error) {
	v, t := settle(n.v, n.t)
	if !v.IsValid() {
		return nil, nil
	}
	var out []jpNode
	switch {
	case v.Kind() == reflect.Slice:
		for i := range v.Len() {
			out = append(out, jpElement(n, v, t, i, nil))
		}
	case t.Kind() == reflect.Struct && structured(t):
		for _, fd := range Describe(t).Fields {
			at := n.at.child(FieldStep(fd.Wire))
			has, err := h.Has(at)
			if err != nil {
				return nil, err
			}
			if !has {
				continue
			}
			fv, _ := v.FieldByIndexErr(fd.index)
			out = append(out, jpNode{at: at, v: fv, t: fd.Typ, fd: fd})
		}
	case v.Kind() == reflect.Map && v.Type().Key().Kind() == reflect.String:
		keys := make([]string, 0, v.Len())
		for _, k := range v.MapKeys() {
			keys = append(keys, k.String())
		}
		slices.Sort(keys)
		for _, k := range keys {
			e := v.MapIndex(reflect.ValueOf(k).Convert(v.Type().Key()))
			if dynamic(t) {
				out = append(out, jpNode{at: n.at.child(FieldStep(k)), v: e, t: anyType})
			} else {
				out = append(out, jpNode{at: n.at.child(MapKeyStep(k)), v: e, t: t.Elem()})
			}
		}
	}
	return out, nil
}

// descendants are a value and every present value within it.
func descendants(h *Holder, n jpNode) ([]jpNode, error) {
	kids, err := children(h, n)
	if err != nil {
		return nil, err
	}
	out := []jpNode{n}
	for _, c := range kids {
		d, err := descendants(h, c)
		if err != nil {
			return nil, err
		}
		out = append(out, d...)
	}
	return out, nil
}

// jpElement is the i-th element of a list, found as accessors find it: by its
// keys in a list map, by the field a filter matched it by in an unstructured
// list, and otherwise by its index.
func jpElement(n jpNode, list reflect.Value, t reflect.Type, i int, by *jpFilter) jpNode {
	e := list.Index(i)
	elemType := anyType
	if !dynamic(t) {
		elemType = t.Elem()
	}
	switch {
	case n.fd != nil && n.fd.Kind == KindListMap:
		sd := Describe(elemType)
		vals := make([]any, len(n.fd.Keys))
		for j, k := range n.fd.Keys {
			vals[j] = e.FieldByIndex(sd.Field(k).index).Interface()
		}
		return jpNode{at: n.at.child(KeyStep(elemType, n.fd.Keys, vals)), v: e, t: elemType}
	case dynamic(t) && by != nil && by.op == "==" && len(by.rel) == 1:
		if m, ok := e.Interface().(map[string]any); ok {
			if kv, ok := m[by.rel[0]]; ok && kv != nil {
				return jpNode{at: n.at.child(UnstructuredKeyStep([]string{by.rel[0]}, []any{kv})), v: e, t: anyType}
			}
		}
	}
	return jpNode{at: n.at.child(IndexStep(i)), v: e, t: elemType}
}

// matches reports whether an element satisfies the filter, comparing their
// JSON forms.
func (f *jpFilter) matches(e reflect.Value) (bool, error) {
	cur, err := ToUnstructured(e)
	if err != nil {
		return false, err
	}
	for _, name := range f.rel {
		m, ok := cur.(map[string]any)
		if !ok {
			cur = nil
			break
		}
		cur = m[name]
	}
	switch f.op {
	case "exists":
		return cur != nil, nil
	case "==":
		return jsonEqual(cur, f.lit), nil
	case "!=":
		return !jsonEqual(cur, f.lit), nil
	}
	a, aok := number(cur)
	b, bok := number(f.lit)
	if !aok || !bok {
		return false, nil
	}
	switch f.op {
	case "<":
		return a < b, nil
	case ">":
		return a > b, nil
	case "<=":
		return a <= b, nil
	}
	return a >= b, nil
}

func jsonEqual(a, b any) bool {
	fa, aok := number(a)
	fb, bok := number(b)
	if aok && bok {
		return fa == fb
	}
	return reflect.DeepEqual(a, b)
}

func number(v any) (float64, bool) {
	switch v := v.(type) {
	case int64:
		return float64(v), true
	case float64:
		return v, true
	}
	return 0, false
}
