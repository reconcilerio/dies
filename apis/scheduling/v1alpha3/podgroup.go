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

// +die:object=true,apiVersion=scheduling.k8s.io/v1alpha3,kind=PodGroup
type _ = schedulingv1alpha3.PodGroup

// +die
// +die:field:name=WorkloadRef,die=WorkloadReferenceDie,pointer=true
// +die:field:name=SchedulingPolicy,die=PodGroupSchedulingPolicyDie
// +die:field:name=SchedulingConstraints,die=PodGroupSchedulingConstraintsDie,pointer=true
// +die:field:name=ResourceClaims,die=PodGroupResourceClaimDie,listType=map
// +die:field:name=DisruptionMode,die=DisruptionModeDie,pointer=true
type _ = schedulingv1alpha3.PodGroupSpec

// +die
type _ = schedulingv1alpha3.WorkloadReference

// +die
// +die:field:name=Basic,die=BasicSchedulingPolicyDie,pointer=true
// +die:field:name=Gang,die=GangSchedulingPolicyDie,pointer=true
type _ = schedulingv1alpha3.PodGroupSchedulingPolicy

// +die
type _ = schedulingv1alpha3.BasicSchedulingPolicy

// +die
type _ = schedulingv1alpha3.GangSchedulingPolicy

// +die
// +die:field:name=Topology,die=TopologyConstraintDie,listType=atomic
type _ = schedulingv1alpha3.PodGroupSchedulingConstraints

// +die
type _ = schedulingv1alpha3.TopologyConstraint

// +die
type _ = schedulingv1alpha3.PodGroupResourceClaim

// +die
// +die:field:name=Single,die=SingleDisruptionModeDie,pointer=true
// +die:field:name=All,die=AllDisruptionModeDie,pointer=true
type _ = schedulingv1alpha3.DisruptionMode

// +die
type _ = schedulingv1alpha3.SingleDisruptionMode

// +die
type _ = schedulingv1alpha3.AllDisruptionMode

// +die
// +die:field:name=Conditions,package=_/meta/v1,die=ConditionDie,listType=map,listMapKey=Type
// +die:field:name=ResourceClaimStatuses,die=PodGroupResourceClaimStatusDie,listType=map
type _ = schedulingv1alpha3.PodGroupStatus

// +die
type _ = schedulingv1alpha3.PodGroupResourceClaimStatus
