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

package base

import (
	"strings"

	"reconciler.io/dies/v2/die/internal/core"
	"reconciler.io/dies/v2/die/internal/refs"
)

// The work of a die's methods is done by functions that aren't generic, so
// that it is compiled once rather than for each type of die: a generic type's
// methods are compiled for each struct type they're instantiated with. The
// methods of the generic bases are wrappers that convert to and from the
// die's types. Most of the work is in [refs], which the die package's
// functions share; what needs this package's helpers is here.

func intent(r refs.Ref) string {
	var lines []string
	describeIntent(core.NodeAt(core.IntentOf(r.H), r.At), r.At, &lines)
	return strings.Join(lines, "\n")
}

// with returns the holder of an immutable die of what a mutable copy, made
// for functions, holds after them, with the copy's observers.
func with(h *core.Holder, m any) (*core.Holder, error) {
	err := keepIn(h, m)
	out := &core.Holder{Base: h.Base, Ov: h.Ov, Wire: h.Wire, Schema: h.Schema}
	out.ObserversOf(h, nil)
	return out, err
}

// withIn keeps what functions passed a mutable die left it holding in the
// die's own holder, which its children share, as a client zeroes the object
// it decodes into, holder and all.
func withIn(r refs.Ref, m any) error {
	err := keepIn(r.H, m)
	core.SetState(m, r.H, r.At)
	return err
}

// keepIn keeps what m holds in h, if it holds another holder, notifying h's
// observers: the other holder, made when its state was zeroed, has none.
func keepIn(h *core.Holder, m any) error {
	// what m holds now may not be the holder it was given: a client zeroes
	// the object it decodes into, holder and all
	r := held(m, h)
	if r == h {
		return nil
	}
	return h.Commit(r.Base, r.Ov, r.Wire, nil)
}
