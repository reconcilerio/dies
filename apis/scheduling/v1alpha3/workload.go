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

package v1alpha3

import (
	schedulingv1alpha3 "k8s.io/api/scheduling/v1alpha3"
)

// +die:object=true,apiVersion=scheduling.k8s.io/v1alpha3,kind=Workload
type _ = schedulingv1alpha3.Workload

// +die
// +die:field:name=ControllerRef,die=TypedLocalObjectReferenceDie,pointer=true
// +die:field:name=PodGroupTemplates,die=PodGroupTemplateDie,listType=map
// +die:field:name=CompositePodGroupTemplates,die=CompositePodGroupTemplateDie,listType=map
type _ = schedulingv1alpha3.WorkloadSpec

// +die
type _ = schedulingv1alpha3.TypedLocalObjectReference

// +die
// +die:field:name=SchedulingPolicy,die=PodGroupSchedulingPolicyDie
// +die:field:name=SchedulingConstraints,die=PodGroupSchedulingConstraintsDie,pointer=true
// +die:field:name=ResourceClaims,die=PodGroupResourceClaimDie,listType=map
// +die:field:name=DisruptionMode,die=DisruptionModeDie,pointer=true
type _ = schedulingv1alpha3.PodGroupTemplate

// +die
// +die:field:name=SchedulingPolicy,die=CompositePodGroupSchedulingPolicyDie
type _ = schedulingv1alpha3.CompositePodGroupTemplate

// +die
// +die:field:name=Basic,die=CompositeBasicSchedulingPolicyDie,pointer=true
// +die:field:name=Gang,die=CompositeGangSchedulingPolicyDie,pointer=true
type _ = schedulingv1alpha3.CompositePodGroupSchedulingPolicy

// +die
type _ = schedulingv1alpha3.CompositeBasicSchedulingPolicy

// +die
type _ = schedulingv1alpha3.CompositeGangSchedulingPolicy
