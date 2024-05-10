/*
Copyright 2022 the original author or authors.

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
	corev1 "k8s.io/api/core/v1"
)

// +die:object=true,apiVersion=v1,kind=Binding
type _ = corev1.Binding

func (d *BindingDie) TargetDie(fn func(d *ObjectReferenceDie)) *BindingDie {
	return d.DieStamp(func(r *corev1.Binding) {
		d := ObjectReferenceBlank.DieImmutable(false).DieFeed(r.Target)
		fn(d)
		r.Target = d.DieRelease()
	})
}
