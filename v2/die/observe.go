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

package die

import (
	"reflect"

	"reconciler.io/dies/v2/die/internal/core"
	"reconciler.io/dies/v2/die/internal/dies"
)

// ObserveAt observes a mutable die at a location within it, its target,
// which target returns given the die: a die within it, or a handle, such as
// a field's or a list map's. After each mutation that changes the target, fn
// is called with a mutable handle to it, within the die as the mutation left
// it, and the change, as LI, the target's immutable form, its DieReadOnly:
// as it was when the observer last saw it, and now. The handle's DieRoot is
// the changed die, so fn may write anywhere in it.
//
// A function of a handle and its change, such as one that keeps a resource's
// conditions, observes dies of any type that have that location:
//
//	err := die.ObserveAt(widget, func(m *WidgetDieMut) *Conditions {
//		return m.Status().Conditions()
//	}, readyConditions.Observe)
//
// It is the die's own [base.DieMutBase.DieObserve] for a location within
// it, and behaves as it does: the observer is kept by the dies derived from
// the die, cascades, settles and is held within an atomic stamp. Observe an
// immutable die within a stamp of its root.
func ObserveAt[M MutableDie, L interface {
	Located
	DieReadOnly() LI
}, LI any](m M, target func(m M) L, fn func(l L, c Change[LI]) error) error {
	h, at := dies.MustStateOf(m)
	on := func(h *core.Holder, at core.Path) L { return target(dies.Wrap[M](h, at)) }
	probe := on(h, at)
	if r, ok := any(probe).(core.Rooted); ok && r.DieRootType() != h.Base.Type() {
		return core.Errorf("die: the target of an observer is within %s, not %s", r.DieRootType(), h.Base.Type())
	}
	return h.Observe(&core.Observer{
		At:     at,
		Target: probe.DiePath(),
		Read:   func(h *core.Holder, at core.Path) (reflect.Value, error) { return h.ValueAt(at) },
		Fn: func(h *core.Holder, at core.Path, d core.Delta) error {
			// the target over the frozen states, only to read
			return fn(on(h, at), Change[LI]{Before: on(d.Was, at).DieReadOnly(), After: on(d.Now, at).DieReadOnly()})
		},
	})
}
