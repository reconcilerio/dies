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
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type stepKind uint8

const (
	stepField stepKind = iota
	stepKey
	stepMapKey
	// stepIndex is a list element by its index. It is resolved against the
	// list as it is, and never recorded in the overlay: a write through it
	// sets the whole list, as befits a list without keys.
	stepIndex
)

// Step is one element of a wire path: a field, a list-map element selected by
// its keys, or a map entry.
type Step struct {
	kind stepKind
	name string
	keys []string
	vals []any
	// elem is the element type of a list-map step, whose keys may have
	// defaults
	elem reflect.Type
	// index is the index of a stepIndex
	index int
}

func FieldStep(name string) Step { return Step{kind: stepField, name: name} }

func MapKeyStep(key string) Step { return Step{kind: stepMapKey, name: key} }

func IndexStep(i int) Step {
	if i < 0 {
		panic(Errorf("die: index %d is negative", i))
	}
	return Step{kind: stepIndex, index: i}
}

// keyStep selects a list-map element. Key values are converted to the
// element's field types and unset keys are defaulted, so that the same
// element always has the same id.
func KeyStep(elem reflect.Type, keys []string, vals []any) Step {
	if len(keys) != len(vals) {
		panic(Errorf("die: %s is keyed by %v, got %d values", elem, keys, len(vals)))
	}
	sd := Describe(elem)
	conv := make([]any, len(vals))
	for i, k := range keys {
		fd := sd.Field(k)
		conv[i] = fd.keyValue(reflect.ValueOf(vals[i]).Convert(fd.Typ)).Interface()
	}
	return Step{kind: stepKey, keys: keys, vals: conv, elem: elem}
}

func (s Step) id() string {
	switch s.kind {
	case stepKey:
		b, _ := json.Marshal(s.vals)
		return "k" + string(b)
	case stepMapKey:
		return "m" + s.name
	case stepIndex:
		return "i" + strconv.Itoa(s.index)
	}
	return "f" + s.name
}

func (s Step) String() string {
	switch s.kind {
	case stepKey:
		parts := make([]string, len(s.keys))
		for i, k := range s.keys {
			parts[i] = fmt.Sprintf("%s=%v", k, s.vals[i])
		}
		return "[" + strings.Join(parts, ",") + "]"
	case stepMapKey:
		return "[" + strconv.Quote(s.name) + "]"
	case stepIndex:
		return "[" + strconv.Itoa(s.index) + "]"
	}
	return "." + s.name
}

// childPath, joinPath and firstIndex are what other packages do with paths;
// as methods, they'd be part of Path's API.
func ChildPath(p Path, s Step) Path { return p.child(s) }

func JoinPath(p, q Path) Path { return p.join(q) }

func FirstIndex(p Path) int { return p.indexed() }

// unstructuredKeyStep selects an element of an unstructured list by key.
func UnstructuredKeyStep(keys []string, vals []any) Step {
	return Step{kind: stepKey, keys: keys, vals: vals}
}

// DiePath returns the path, so that a [Path] is a [Located].
func (p Path) DiePath() Path { return p }

// Path is a location in an object's wire schema. Embedded structs are
// flattened, so paths match JSON, field sets and patches.
type Path []Step

// ParsePath parses a path of field steps, such as ".status".
func ParsePath(s string) Path {
	var p Path
	for name := range strings.SplitSeq(strings.TrimPrefix(s, "."), ".") {
		if name != "" {
			p = append(p, FieldStep(name))
		}
	}
	return p
}

func (p Path) String() string {
	if len(p) == 0 {
		return "."
	}
	var sb strings.Builder
	for _, s := range p {
		sb.WriteString(s.String())
	}
	return sb.String()
}

func (p Path) child(s Step) Path {
	out := make(Path, len(p)+1)
	copy(out, p)
	out[len(p)] = s
	return out
}

func (p Path) join(q Path) Path {
	out := make(Path, 0, len(p)+len(q))
	return append(append(out, p...), q...)
}

// indexed returns the position of the first index step in a path, or -1.
// keyed reports whether a path passes through a list-map element by key.
func (p Path) keyed() bool {
	for _, s := range p {
		if s.kind == stepKey {
			return true
		}
	}
	return false
}

func (p Path) indexed() int {
	for i, s := range p {
		if s.kind == stepIndex {
			return i
		}
	}
	return -1
}
