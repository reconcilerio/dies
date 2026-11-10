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

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/utils/ptr"
	"reconciler.io/dies/v2/die/base"
	"reconciler.io/dies/v2/die/fields"
)

// DieObjectOf carries an immutable object die's types for generic code, as
// die.DieObjectOf does, with the accessor of its metadata's die, from this
// package. The die API's interface doesn't name a die of meta/v1, so that it
// doesn't depend on one; code that reaches into an object's metadata through
// its die names this one instead.
type DieObjectOf[T any, D base.DieOf[T, D, M], M base.DieMutOf[T, D, M]] interface {
	fields.DieObjectOf[T, D, M]
	Metadata() *ObjectMetaDie
}

// DieObjectMutOf carries a mutable object die's types for generic code, as
// [DieObjectOf] does an immutable one's.
type DieObjectMutOf[T any, D base.DieOf[T, D, M], M base.DieMutOf[T, D, M]] interface {
	fields.DieObjectMutOf[T, D, M]
	Metadata() *ObjectMetaDieMut
}

func (m *ObjectMetaDieMut) ControlledBy(obj runtime.Object, scheme *runtime.Scheme) error {
	gvks, _, err := scheme.ObjectKinds(obj)
	if err != nil {
		if gvk := obj.GetObjectKind().GroupVersionKind(); gvk != (schema.GroupVersionKind{}) {
			gvks = []schema.GroupVersionKind{gvk}
		} else {
			return err
		}
	}
	ref := &OwnerReference{
		APIVersion:         gvks[0].GroupVersion().String(),
		Kind:               gvks[0].Kind,
		Name:               obj.(Object).GetName(),
		UID:                obj.(Object).GetUID(),
		BlockOwnerDeletion: ptr.To(true),
		Controller:         ptr.To(true),
	}
	return m.OwnerReferences().At(ref.UID).DieFeed(ref)
}
