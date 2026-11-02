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

// Package dies makes dies of a type over state, for the die packages that
// can't reach a die's state themselves.
package dies

import (
	"reflect"

	"reconciler.io/dies/v2/die/base"
	"reconciler.io/dies/v2/die/internal/core"
)

// wrapDie returns a D die over state.
func WrapDie[T any, D base.DieOf[T, D]](h *core.Holder, at core.Path) D {
	var d D
	d = d.DieFrom(base.DieBase[T, D]{})
	core.SetState(d, h, at)
	return d
}

// Wrap returns a die of type D, which may be a die of any type, over state.
func Wrap[D any](h *core.Holder, at core.Path) D {
	d := reflect.New(reflect.TypeFor[D]().Elem()).Interface().(D)
	core.SetState(d, h, at)
	return d
}

// mustStateOf returns a die's state.
func MustStateOf(d any) (*core.Holder, core.Path) {
	h, at, ok := core.StateOf(d)
	if !ok {
		panic(core.Errorf("die: not a die"))
	}
	return h, at
}

// selfOf returns how a write through d returns a die: d's own location, or a
// new root when an immutable root's write made one.
func SelfOf[T any, D base.DieOf[T, D]](d D) func(*core.Holder) D {
	dh, at := MustStateOf(d)
	return func(h *core.Holder) D {
		if h != dh {
			return WrapDie[T, D](h, nil)
		}
		return WrapDie[T, D](h, at)
	}
}
