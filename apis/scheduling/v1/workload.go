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
	schedulingv1 "k8s.io/api/scheduling/v1"
)

// +die
type _ = schedulingv1.WorkloadTopologyConstraint

// +die
// +die:field:name=Basic,die=WorkloadPodGroupBasicSchedulingPolicyDie,pointer=true
// +die:field:name=Gang,die=WorkloadPodGroupGangSchedulingPolicyDie,pointer=true
type _ = schedulingv1.WorkloadPodGroupSchedulingPolicy

// +die
type _ = schedulingv1.WorkloadPodGroupBasicSchedulingPolicy

// +die
type _ = schedulingv1.WorkloadPodGroupGangSchedulingPolicy

// +die
// +die:field:name=Topology,die=WorkloadTopologyConstraintDie,listType=atomic
type _ = schedulingv1.WorkloadPodGroupSchedulingConstraints

// +die
type _ = schedulingv1.WorkloadPodGroupResourceClaim

// +die
// +die:field:name=Single,die=WorkloadPodGroupSingleDisruptionModeDie,pointer=true
// +die:field:name=All,die=WorkloadPodGroupAllDisruptionModeDie,pointer=true
type _ = schedulingv1.WorkloadPodGroupDisruptionMode

// +die
type _ = schedulingv1.WorkloadPodGroupSingleDisruptionMode

// +die
type _ = schedulingv1.WorkloadPodGroupAllDisruptionMode

// +die
// +die:field:name=Basic,die=WorkloadCompositePodGroupBasicSchedulingPolicyDie,pointer=true
// +die:field:name=Gang,die=WorkloadCompositePodGroupGangSchedulingPolicyDie,pointer=true
type _ = schedulingv1.WorkloadCompositePodGroupSchedulingPolicy

// +die
type _ = schedulingv1.WorkloadCompositePodGroupBasicSchedulingPolicy

// +die
type _ = schedulingv1.WorkloadCompositePodGroupGangSchedulingPolicy

// +die
// +die:field:name=Topology,die=WorkloadTopologyConstraintDie,listType=atomic
type _ = schedulingv1.WorkloadCompositePodGroupSchedulingConstraints

// +die
// +die:field:name=Single,die=WorkloadCompositePodGroupSingleDisruptionModeDie,pointer=true
// +die:field:name=All,die=WorkloadCompositePodGroupAllDisruptionModeDie,pointer=true
type _ = schedulingv1.WorkloadCompositePodGroupDisruptionMode

// +die
type _ = schedulingv1.WorkloadCompositePodGroupSingleDisruptionMode

// +die
type _ = schedulingv1.WorkloadCompositePodGroupAllDisruptionMode
