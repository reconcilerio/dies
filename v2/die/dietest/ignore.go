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

package dietest

import (
	"reflect"
	"strings"

	"github.com/google/go-cmp/cmp"

	"reconciler.io/dies/v2/die/base"
	"reconciler.io/dies/v2/die/internal/core"
)

// IgnoreFields is a go-cmp option that ignores locations within dies, named
// as a jig binds them, by navigating from a blank:
//
//	dietest.Diff(x, y, dietest.IgnoreFields(
//		dieappsv1.DeploymentBlank.Metadata().ResourceVersion(),
//		dieappsv1.DeploymentBlank.Status(),
//	))
//
// It works with both [Transform] and [TransformRelease], so the same
// locations serve [Diff] and [DiffRelease]. A location within a blank
// applies to dies of the blank's type, wherever they are within the values
// compared: a location within a DeploymentBlank to Deployment dies, not to
// dies of a Deployment's parts or of other types. A [base.Path] applies to
// dies of every type.
//
// With [Transform], it ignores the value at and below each location, and
// the intent there: each line of intent for the location or below it. A
// line that sets an enclosing value whole, such as one feeding the die, is
// not ignored, as part of it is outside the location.
//
// A location within a list map ignores the element with its key, such as
// Containers().At("sidecar"), wherever the element is in the list, as does
// one of an unstructured list's element by its key, such as
// Field("containers").Item("name", "sidecar"). A location of an element by
// its index, such as Index(0), ignores its value but not its intent, as
// intent locates elements by key.
func IgnoreFields(at ...base.Located) cmp.Option {
	ig := ignoring{}
	for _, l := range at {
		root, err := core.RootOf(l)
		if err != nil {
			root = nil
		}
		ig.paths = append(ig.paths, ignored{root: root, path: l.DiePath()})
	}
	return cmp.FilterPath(ig.filter, cmp.Ignore())
}

type ignored struct {
	// root is the type of dies the path is within, or nil for any
	root reflect.Type
	path core.Path
}

type ignoring struct{ paths []ignored }

func (ig ignoring) filter(p cmp.Path) (ignore bool) {
	// a value the walk doesn't understand isn't ignored
	defer func() {
		if recover() != nil {
			ignore = false
		}
	}()
	// the location within the nearest die
	at := -1
	release := false
	for i := len(p) - 1; i > 0 && at < 0; i-- {
		if t, ok := p[i].(cmp.Transform); ok {
			switch t.Name() {
			case "die.Content":
				at = i
			case "die.Release":
				at, release = i, true
			}
		}
	}
	if at < 0 {
		return false
	}
	t := dieValueType(p[at-1])
	if t == nil {
		return false
	}
	var paths []core.Path
	for _, i := range ig.paths {
		if i.root == nil || core.Deref(i.root) == core.Deref(t) {
			paths = append(paths, i.path)
		}
	}
	if len(paths) == 0 {
		return false
	}
	steps := p[at+1:]

	if !release {
		// the content's value or intent
		for len(steps) > 0 {
			if _, ok := steps[0].(cmp.Indirect); !ok {
				break
			}
			steps = steps[1:]
		}
		if len(steps) == 0 {
			return false
		}
		f, ok := steps[0].(cmp.StructField)
		if !ok {
			return false
		}
		steps = steps[1:]
		if f.Name() == "Intent" {
			return ignoresIntent(steps, paths)
		}
		if f.Name() != "Value" {
			return false
		}
	}

	// ignored when each side present is at or below an ignored path
	present := false
	for side := range 2 {
		sp, ok, exists := walk(t, steps, side, paths)
		if !exists {
			continue
		}
		present = true
		if !ok || !within(sp, paths) {
			return false
		}
	}
	return present
}

// dieValueType is the value type of the die a step reaches.
func dieValueType(s cmp.PathStep) reflect.Type {
	vx, vy := s.Values()
	for _, v := range []reflect.Value{vx, vy} {
		if !v.IsValid() || ((v.Kind() == reflect.Interface || v.Kind() == reflect.Pointer) && v.IsNil()) {
			continue
		}
		return core.ValueTypeOf(v.Interface())
	}
	return nil
}

func within(p core.Path, paths []core.Path) bool {
	for _, q := range paths {
		if core.HasPrefix(p, q) {
			return true
		}
	}
	return false
}

// walk returns the path within a die of t that steps reach on one side,
// whether it could be found, and whether the side has a value there. It
// walks the die's value as its Go type has it, or as its JSON decodes,
// following t's structure for the names of fields and the keys of lists. An
// element of a list t doesn't describe, such as an unstructured die's, is
// located by the key of an ignored path there, if it has that key, and
// otherwise by its index.
func walk(t reflect.Type, steps []cmp.PathStep, side int, paths []core.Path) (path core.Path, ok, exists bool) {
	if len(steps) > 0 {
		vx, vy := steps[len(steps)-1].Values()
		exists = []reflect.Value{vx, vy}[side].IsValid()
	} else {
		exists = true
	}
	cur := t
	var field *core.FieldDesc
	// a struct, and the Go index within it, of an inline struct's fields
	var outer reflect.Type
	var inline []int
	for _, s := range steps {
		switch s := s.(type) {
		case cmp.Indirect, cmp.TypeAssertion:
		case cmp.StructField:
			idx := []int{s.Index()}
			if inline != nil {
				idx = append(inline, s.Index())
			} else {
				outer = core.Deref(cur)
			}
			fd, partial := core.WireFieldAt(outer, idx)
			switch {
			case fd != nil:
				path = append(path, core.FieldStep(fd.Wire))
				cur, field, inline = fd.Typ, fd, nil
			case partial:
				inline = idx
			default:
				return nil, false, exists
			}
		case cmp.SliceIndex:
			ct := deref(cur)
			vx, vy := s.Values()
			el := []reflect.Value{vx, vy}[side]
			ix, iy := s.SplitKeys()
			if ct != nil && ct.Kind() == reflect.Slice && field != nil && len(field.Keys) > 0 {
				if !el.IsValid() {
					return nil, false, false
				}
				path = append(path, core.KeyStepOf(core.Deref(ct.Elem()), field.Keys, el))
			} else if st, keyed := keyedLike(path, paths, el); keyed {
				path = append(path, st)
			} else {
				i := []int{ix, iy}[side]
				if i < 0 {
					return nil, false, false
				}
				path = append(path, core.IndexStep(i))
			}
			field = nil
			if ct != nil && ct.Kind() == reflect.Slice {
				cur = ct.Elem()
			} else {
				cur = nil
			}
		case cmp.MapIndex:
			if s.Key().Kind() != reflect.String {
				return nil, false, exists
			}
			k := s.Key().String()
			ct := deref(cur)
			switch {
			case ct != nil && ct.Kind() == reflect.Struct:
				path = append(path, core.FieldStep(k))
				field = core.WireField(ct, k)
				cur = nil
				if field != nil {
					cur = field.Typ
				}
			case ct != nil && ct.Kind() == reflect.Map && ct.Elem().Kind() != reflect.Interface:
				path = append(path, core.MapKeyStep(k))
				cur, field = ct.Elem(), nil
			default:
				// an object of no known type, such as an unstructured die's
				path = append(path, core.FieldStep(k))
				cur, field = nil, nil
			}
		default:
			return nil, false, exists
		}
	}
	return path, true, exists
}

// keyedLike returns the step of a list element by the key of an ignored path
// at the element's location, if the element has that key.
func keyedLike(at core.Path, paths []core.Path, el reflect.Value) (core.Step, bool) {
	if !el.IsValid() {
		return core.Step{}, false
	}
	for _, q := range paths {
		if len(q) <= len(at) || !core.HasPrefix(at, q[:len(at)]) {
			continue
		}
		if st, ok := core.UnstructuredKeyStepOf(q[len(at)], el); ok && core.HasPrefix(core.Path{st}, q[len(at):len(at)+1]) {
			return st, true
		}
	}
	return core.Step{}, false
}

func deref(t reflect.Type) reflect.Type {
	if t == nil {
		return nil
	}
	return core.Deref(t)
}

// ignoresIntent reports whether steps reach a line of intent for an ignored
// path or below it.
func ignoresIntent(steps []cmp.PathStep, paths []core.Path) bool {
	if len(steps) != 1 {
		return false
	}
	s, ok := steps[0].(cmp.SliceIndex)
	if !ok {
		return false
	}
	present := false
	vx, vy := s.Values()
	for _, v := range []reflect.Value{vx, vy} {
		if !v.IsValid() {
			continue
		}
		present = true
		if !intentWithin(v.String(), paths) {
			return false
		}
	}
	return present
}

// intentWithin reports whether a line of intent, such as
// `set .spec.replicas = 1`, is for an ignored path or below it.
func intentWithin(line string, paths []core.Path) bool {
	_, lp, _ := strings.Cut(line, " ")
	lp, _, _ = strings.Cut(lp, " = ")
	for _, q := range paths {
		qs := q.String()
		if qs == "." || lp == qs || strings.HasPrefix(lp, qs) && (lp[len(qs)] == '.' || lp[len(qs)] == '[') {
			return true
		}
	}
	return false
}
