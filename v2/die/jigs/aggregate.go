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

package jigs

import (
	"reflect"

	"reconciler.io/dies/v2/die/base"
	"reconciler.io/dies/v2/die/internal/core"
)

// An aggregate binds a list map of a jig to several locations of a source,
// such as a pod spec's containers and its init containers. Each element of
// the jig knows its source, so a change within an element goes back to it,
// and removing an element removes it there. Adding an element is allowed
// only when exactly one source list could take it; otherwise it is
// ambiguous, an error. An element whose key is in more than one source is
// ambiguous too.

// Member is a source of the elements of an aggregate, each an X.
type Member[X any] struct {
	m core.Member
}

// erased is the jig for the engine.
func (v *Jig[V, VD, VDM]) erased() *core.ViewSource {
	return &core.ViewSource{BindingsFor: v.bindingsFor, Type: reflect.TypeFor[V]()}
}

func (v *Jig[V, VD, VDM]) bindingsFor(t reflect.Type) ([]Binding, bool) {
	b, ok := v.sources[t]
	return b, ok
}

// Each is a source list whose elements are the aggregate's elements.
func Each[X any](source base.LocatedAs[[]X]) Member[X] {
	root, err := core.RootOf(source)
	return Member[X]{core.Member{Src: source.DiePath(), Root: root, Err: err}}
}

// One is a single source value that is an element of the aggregate when set.
func One[X any](source base.LocatedAs[X]) Member[X] {
	root, err := core.RootOf(source)
	return Member[X]{core.Member{Src: source.DiePath(), Root: root, Single: true, Err: err}}
}

// EachVia is a source list whose elements are presented as the aggregate's
// elements by a jig.
func EachVia[S, X any, XD base.DieOf[X, XD, XDM], XDM base.DieMutOf[X, XD, XDM]](source base.LocatedAs[[]S], jig *Jig[X, XD, XDM]) Member[X] {
	root, err := core.RootOf(source)
	return Member[X]{core.Member{Src: source.DiePath(), Root: root, View: jig.erased(), Err: err}}
}

// EachAt is a list at a path in an unstructured object, whose elements are
// the aggregate's elements.
func EachAt[X any](path string) Member[X] {
	return Member[X]{core.Member{Src: core.ParsePath(path)}}
}

// OneAt is a single value at a path in an unstructured object.
func OneAt[X any](path string) Member[X] {
	return Member[X]{core.Member{Src: core.ParsePath(path), Single: true}}
}

// OneVia is a single source value presented as an element of the aggregate
// by a jig.
func OneVia[S, X any, XD base.DieOf[X, XD, XDM], XDM base.DieMutOf[X, XD, XDM]](source base.LocatedAs[S], jig *Jig[X, XD, XDM]) Member[X] {
	root, err := core.RootOf(source)
	return Member[X]{core.Member{Src: source.DiePath(), Root: root, Single: true, View: jig.erased(), Err: err}}
}

// BindAll binds a list map of a jig to the elements of several sources.
func BindAll[X any](jig base.LocatedAs[[]X], members ...Member[X]) Binding {
	agg := &core.Aggregate{}
	root, err := core.RootOf(jig)
	var sroot reflect.Type
	for i, m := range members {
		if err == nil && m.m.Err != nil {
			err = m.m.Err
		}
		if err == nil && i > 0 && m.m.Root != sroot {
			err = core.Errorf("die: members of an aggregate are from %v and %v", sroot, m.m.Root)
		}
		sroot = m.m.Root
		agg.Members = append(agg.Members, m.m)
	}
	return core.WithErr(core.NewBinding(jig.DiePath(), nil, root, sroot, nil, agg), err)
}
