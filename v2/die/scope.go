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
	"reconciler.io/dies/v2/die/internal/refs"
)

// A Scope narrows the intent a patch is released from, for that patch: the
// die keeps all of its intent. Intent a scope leaves out is folded into the
// value the patch is made against, so the patch sends only what the scope
// keeps. A scope can be kept in a variable and passed to several releases:
//
//	status := die.Subresource("status")
//	body, err := die.ReleaseApplyPatchE(d, status)
//
// Scopes passed together each narrow what the one before kept.
type Scope struct {
	keep func(h *core.Holder, t reflect.Type) (*core.Holder, error)
}

// Only keeps the intent at and below the locations, such as
// d.Spec().Replicas(), which are within the die released.
func Only(at ...Located) Scope {
	return Scope{func(h *core.Holder, _ reflect.Type) (*core.Holder, error) { return refs.Keep(h, at, true) }}
}

// Except keeps the intent but at and below the locations.
func Except(at ...Located) Scope {
	return Scope{func(h *core.Holder, _ reflect.Type) (*core.Holder, error) { return refs.Keep(h, at, false) }}
}

// Subresource keeps the intent for a subresource, such as "status", as the
// fields its type marks as served through it.
func Subresource(name string) Scope {
	return Scope{func(h *core.Holder, t reflect.Type) (*core.Holder, error) {
		return refs.Keep(h, refs.Subresources(t, func(s string) bool { return s == name }), true)
	}}
}

// MainResource keeps the intent for the main resource, leaving out every
// subresource's.
func MainResource() Scope {
	return Scope{func(h *core.Holder, t reflect.Type) (*core.Holder, error) {
		return refs.Keep(h, refs.Subresources(t, func(s string) bool { return s != "" }), false)
	}}
}

// scoped returns the state of a die with only the intent scopes keep, for a
// release.
func scoped(d Die, method string, scopes []Scope) (refs.Ref, error) {
	r, err := refOf(d)
	if err != nil || len(scopes) == 0 {
		return r, err
	}
	h := r.Root(method)
	for _, s := range scopes {
		if s.keep == nil {
			continue
		}
		if h, err = s.keep(h, r.T); err != nil {
			return refs.Ref{}, err
		}
	}
	return refs.Ref{H: h, T: r.T}, nil
}
