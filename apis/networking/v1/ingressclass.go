/*
Copyright 2021 the original author or authors.

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
	networkingv1 "k8s.io/api/networking/v1"
)

// +die:object=true
type IngressClass = networkingv1.IngressClass

// +die
type IngressClassSpec = networkingv1.IngressClassSpec

type ingressClassSpec interface {
	ParametersDie(fn func(d IngressClassParametersReferenceDie)) IngressClassSpecDie
}

func (d *ingressClassSpecDie) ParametersDie(fn func(d IngressClassParametersReferenceDie)) IngressClassSpecDie {
	return d.DieStamp(func(r *networkingv1.IngressClassSpec) {
		d := IngressClassParametersReferenceBlank.DieImmutable(false).DieFeedPtr(r.Parameters)
		fn(d)
		r.Parameters = d.DieReleasePtr()
	})
}

// +die
type IngressClassParametersReference = networkingv1.IngressClassParametersReference
