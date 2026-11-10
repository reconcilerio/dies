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

// Package jsonpath finds locations within dies with JSONPath expressions
// (https://goessner.net/articles/JsonPath/). A match is a location within
// the die, found as its accessors find it, not a copy of a value: reads see
// the die as it is, and writes are recorded as intent, so patches carry them.
//
// $ is the root die and . the die queried, so a function written for a part
// of an object, such as a PodSpecDie, finds what's within that part with
// .containers[*], and what's about the object with $.metadata.namespace,
// whatever object it is part of. Any expression that doesn't start with $
// starts at the die queried. The paths found are within the root either way,
// as every path is.
//
//   - an element of a list map is found by its keys, so containers[0] and
//     containers[?(@.name=='app')] are both .containers[name=app], and a write
//     through either merges by key
//   - an element of an atomic list is found by its index
//   - in an unstructured list, a filter that matches one element by one field
//     finds it by that field; other elements, by index
//   - a definite path, with no wildcard, slice, filter or recursive descent
//     before a step, resolves by the schema whether or not the value is set, so that
//     setting it creates it; other steps only match values that are present
//   - in a typed die, a field the schema doesn't have is an error, not a
//     silent miss
//
// Expressions are kubectl's (https://kubernetes.io/docs/reference/kubectl/jsonpath/),
// parsed by client-go: names, indexes, slices, unions, wildcards, recursive
// descent, and filters that compare a field with a literal (==, !=, and <,
// >, <= and >= on numbers) or test that it exists. One expression, with or
// without {braces}; templates, such as range, aren't expressions.
//
// Unstructured dies and their nodes find nodes with Node.Find.
package jsonpath

import (
	"reflect"

	"reconciler.io/dies/v2/die/base"
	"reconciler.io/dies/v2/die/fields"
	"reconciler.io/dies/v2/die/internal/core"
	"reconciler.io/dies/v2/die/internal/dies"
)

// Find returns the locations within a die that an expression matches.
func Find(d base.Die, expr string) ([]base.Path, error) {
	h, at, ok := core.StateOf(d)
	if !ok {
		return nil, core.Errorf("die: %T is not a die", d)
	}
	return core.ResolveJSONPath(h, at, expr)
}

// MustFind is [Find], panicking rather than returning an error.
func MustFind(d base.Die, expr string) []base.Path { return core.Must(Find(d, expr)) }

// FindDies returns the matches of an expression within an immutable die as
// immutable dies of type CD, such as *ContainerDie, each a die within d. It
// fails if a match isn't a value of CD's type. An unstructured die is a die
// for any unstructured value.
func FindDies[CD base.ImmutableDie](d base.ImmutableDie, expr string) ([]CD, error) {
	return matches[CD](d, expr)
}

// MustFindDies is [FindDies], panicking rather than returning an error.
func MustFindDies[CD base.ImmutableDie](d base.ImmutableDie, expr string) []CD {
	return core.Must(FindDies[CD](d, expr))
}

// FindDiesMut returns the matches of an expression within a mutable die as mutable
// dies of type CDM, such as *ContainerDieMut, whose writes are writes to d.
func FindDiesMut[CDM base.MutableDie](d base.MutableDie, expr string) ([]CDM, error) {
	return matches[CDM](d, expr)
}

// MustFindDiesMut is [FindDiesMut], panicking rather than returning an error.
func MustFindDiesMut[CDM base.MutableDie](d base.MutableDie, expr string) []CDM {
	return core.Must(FindDiesMut[CDM](d, expr))
}

func matches[CD base.Die](d base.Die, expr string) ([]CD, error) {
	h, at, ok := core.StateOf(d)
	if !ok {
		return nil, core.Errorf("die: %T is not a die", d)
	}
	paths, err := core.ResolveJSONPath(h, at, expr)
	if err != nil {
		return nil, err
	}
	out := make([]CD, len(paths))
	for i, p := range paths {
		cd := dies.Wrap[CD](h, p)
		have, want := core.TypeAt(h.Base.Type(), p), core.ValueTypeOf(cd)
		if core.Deref(have) != want && (have.Kind() != reflect.Interface || want != reflect.TypeFor[map[string]any]()) {
			return nil, core.Errorf("die: JSONPath %q: the value at %s is a %s, not a %s", expr, p, have, want)
		}
		out[i] = cd
	}
	return out, nil
}

// FindValues returns the matches of an expression as handles to values of type V
// within d. It fails if a match isn't a V.
func FindValues[V any](d base.Die, expr string) ([]fields.Value[V], error) {
	paths, err := valuePaths[V](d, expr)
	if err != nil {
		return nil, err
	}
	out := make([]fields.Value[V], len(paths))
	for i, p := range paths {
		if out[i], err = fields.LocateValueAt[V](d, p); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// MustFindValues is [FindValues], panicking rather than returning an error.
func MustFindValues[V any](d base.Die, expr string) []fields.Value[V] {
	return core.Must(FindValues[V](d, expr))
}

// FindValuesMut returns the matches of an expression as handles to values of type
// V within a mutable die. It fails if a match isn't a V.
func FindValuesMut[V any](d base.MutableDie, expr string) ([]fields.ValueMut[V], error) {
	paths, err := valuePaths[V](d, expr)
	if err != nil {
		return nil, err
	}
	out := make([]fields.ValueMut[V], len(paths))
	for i, p := range paths {
		if out[i], err = fields.LocateValueMutAt[V](d, p); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// MustFindValuesMut is [FindValuesMut], panicking rather than returning an error.
func MustFindValuesMut[V any](d base.MutableDie, expr string) []fields.ValueMut[V] {
	return core.Must(FindValuesMut[V](d, expr))
}

func valuePaths[V any](d base.Die, expr string) ([]base.Path, error) {
	h, at, ok := core.StateOf(d)
	if !ok {
		return nil, core.Errorf("die: %T is not a die", d)
	}
	paths, err := core.ResolveJSONPath(h, at, expr)
	if err != nil {
		return nil, err
	}
	for _, p := range paths {
		if t := core.TypeAt(h.Base.Type(), p); t != reflect.TypeFor[V]() {
			return nil, core.Errorf("die: JSONPath %q: the value at %s is a %s, not a %s", expr, p, t, reflect.TypeFor[V]())
		}
	}
	return paths, nil
}
