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

// WrapDie returns an immutable D die over state.
func WrapDie[T any, D base.DieOf[T, D, M], M base.DieMutOf[T, D, M]](h *core.Holder, at core.Path) D {
	var d D
	d = d.DieFrom(base.DieBase[T, D, M]{})
	core.SetState(d, h, at)
	return d
}

// WrapDieMut returns a mutable M die over state.
func WrapDieMut[T any, D base.DieOf[T, D, M], M base.DieMutOf[T, D, M]](h *core.Holder, at core.Path) M {
	var m M
	m = m.DieFrom(base.DieMutBase[T, D, M]{})
	core.SetState(m, h, at)
	return m
}

// Wrap returns a die of type D, which may be a die of any type, immutable or
// mutable, over state.
func Wrap[D any](h *core.Holder, at core.Path) D {
	d := reflect.New(reflect.TypeFor[D]().Elem()).Interface().(D)
	core.SetState(d, h, at)
	return d
}

// MustStateOf returns a die's state.
func MustStateOf(d any) (*core.Holder, core.Path) {
	h, at, ok := core.StateOf(d)
	if !ok {
		panic(core.Errorf("die: not a die"))
	}
	return h, at
}
