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
	"reconciler.io/dies/v2/die/internal/core"
	"reconciler.io/dies/v2/die/internal/refs"
)

// StampValue passes fn a pointer to a copy of a mutable die's Go value, and
// records what fn changes as intent, as writes through the die's accessors
// would. The value's type is the die's, so fn's is checked against it:
//
//	err := die.StampValue(d, func(r *appsv1.Deployment) {
//		r.Spec.Paused = true
//	})
//
// Writing through the die, within its DieStamp, records intent as it is
// written; StampValue suits code written for the Go type.
func StampValue[T any](d interface {
	MutableDie
	LocatedAs[T]
}, fn func(r *T)) error {
	r, ok := refs.Of(d)
	if !ok {
		return core.Errorf("die: %T is not a die", d)
	}
	return r.Stamp(func(p any) { fn(p.(*T)) })
}
