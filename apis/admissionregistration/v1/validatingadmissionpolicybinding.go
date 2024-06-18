/*
Copyright 2024 the original author or authors.

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
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	diemetav1 "reconciler.io/dies/apis/meta/v1"
)

// +die:object=true,apiVersion=admissionregistration.k8s.io/v1,kind=ValidatingAdmissionPolicyBinding
type _ = admissionregistrationv1.ValidatingAdmissionPolicyBinding

// +die
type _ = admissionregistrationv1.ValidatingAdmissionPolicyBindingSpec

func (d *ValidatingAdmissionPolicyBindingSpecDie) ParamRefDie(fn func(d *ParamRefDie)) *ValidatingAdmissionPolicyBindingSpecDie {
	return d.DieStamp(func(r *admissionregistrationv1.ValidatingAdmissionPolicyBindingSpec) {
		d := ParamRefBlank.DieImmutable(false).DieFeedPtr(r.ParamRef)
		fn(d)
		r.ParamRef = d.DieReleasePtr()
	})
}

func (d *ValidatingAdmissionPolicyBindingSpecDie) MatchResourcesDie(fn func(d *MatchResourcesDie)) *ValidatingAdmissionPolicyBindingSpecDie {
	return d.DieStamp(func(r *admissionregistrationv1.ValidatingAdmissionPolicyBindingSpec) {
		d := MatchResourcesBlank.DieImmutable(false).DieFeedPtr(r.MatchResources)
		fn(d)
		r.MatchResources = d.DieReleasePtr()
	})
}

// +die
type _ = admissionregistrationv1.ParamRef

func (d *ParamRefDie) SelectorDie(fn func(d *diemetav1.LabelSelectorDie)) *ParamRefDie {
	return d.DieStamp(func(r *admissionregistrationv1.ParamRef) {
		d := diemetav1.LabelSelectorBlank.DieImmutable(false).DieFeedPtr(r.Selector)
		fn(d)
		r.Selector = d.DieReleasePtr()
	})
}
