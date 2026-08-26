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
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
)

// +die:object=true,apiVersion=admissionregistration.k8s.io/v1,kind=MutatingAdmissionPolicy
type _ = admissionregistrationv1.MutatingAdmissionPolicy

// +die
// +die:field:name=ParamKind,die=ParamKindDie,pointer=true
// +die:field:name=MatchConstraints,die=MatchResourcesDie,pointer=true
// +die:field:name=Variables,die=VariableDie,listType=atomic
// +die:field:name=Mutations,die=MutationDie,listType=atomic
// +die:field:name=MatchConditions,die=MatchConditionDie,listType=map
type _ = admissionregistrationv1.MutatingAdmissionPolicySpec

// +die
// +die:field:name=ApplyConfiguration,die=ApplyConfigurationDie,pointer=true
// +die:field:name=JSONPatch,die=JSONPatchDie,pointer=true
type _ = admissionregistrationv1.Mutation

// +die
type _ = admissionregistrationv1.ApplyConfiguration

// +die
type _ = admissionregistrationv1.JSONPatch
