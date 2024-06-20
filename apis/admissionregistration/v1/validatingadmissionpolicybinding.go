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
)

// +die:object=true,apiVersion=admissionregistration.k8s.io/v1,kind=ValidatingAdmissionPolicyBinding
type _ = admissionregistrationv1.ValidatingAdmissionPolicyBinding

// +die
// +die:field:name=ParamRef,die=ParamRefDie,pointer=true
// +die:field:name=MatchResources,die=MatchResourcesDie,pointer=true
type _ = admissionregistrationv1.ValidatingAdmissionPolicyBindingSpec

// +die
// +die:field:name=Selector,package=_/meta/v1,die=LabelSelectorDie,pointer=true
type _ = admissionregistrationv1.ParamRef
