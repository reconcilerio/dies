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
	dieschedulingv1 "reconciler.io/dies/apis/scheduling/v1"
)

// The workload-aware scheduling types graduated to scheduling.k8s.io/v1 and
// remain in k8s.io/api/scheduling/v1alpha3 as type aliases. The dies below
// mirror those aliases.

// Deprecated WorkloadTopologyConstraintDie, moved to package reconciler.io/dies/apis/scheduling/v1
type WorkloadTopologyConstraintDie = dieschedulingv1.WorkloadTopologyConstraintDie

// Deprecated WorkloadTopologyConstraintBlank, moved to package reconciler.io/dies/apis/scheduling/v1
var WorkloadTopologyConstraintBlank = dieschedulingv1.WorkloadTopologyConstraintBlank

// Deprecated WorkloadPodGroupSchedulingPolicyDie, moved to package reconciler.io/dies/apis/scheduling/v1
type WorkloadPodGroupSchedulingPolicyDie = dieschedulingv1.WorkloadPodGroupSchedulingPolicyDie

// Deprecated WorkloadPodGroupSchedulingPolicyBlank, moved to package reconciler.io/dies/apis/scheduling/v1
var WorkloadPodGroupSchedulingPolicyBlank = dieschedulingv1.WorkloadPodGroupSchedulingPolicyBlank

// Deprecated WorkloadPodGroupBasicSchedulingPolicyDie, moved to package reconciler.io/dies/apis/scheduling/v1
type WorkloadPodGroupBasicSchedulingPolicyDie = dieschedulingv1.WorkloadPodGroupBasicSchedulingPolicyDie

// Deprecated WorkloadPodGroupBasicSchedulingPolicyBlank, moved to package reconciler.io/dies/apis/scheduling/v1
var WorkloadPodGroupBasicSchedulingPolicyBlank = dieschedulingv1.WorkloadPodGroupBasicSchedulingPolicyBlank

// Deprecated WorkloadPodGroupGangSchedulingPolicyDie, moved to package reconciler.io/dies/apis/scheduling/v1
type WorkloadPodGroupGangSchedulingPolicyDie = dieschedulingv1.WorkloadPodGroupGangSchedulingPolicyDie

// Deprecated WorkloadPodGroupGangSchedulingPolicyBlank, moved to package reconciler.io/dies/apis/scheduling/v1
var WorkloadPodGroupGangSchedulingPolicyBlank = dieschedulingv1.WorkloadPodGroupGangSchedulingPolicyBlank

// Deprecated WorkloadPodGroupSchedulingConstraintsDie, moved to package reconciler.io/dies/apis/scheduling/v1
type WorkloadPodGroupSchedulingConstraintsDie = dieschedulingv1.WorkloadPodGroupSchedulingConstraintsDie

// Deprecated WorkloadPodGroupSchedulingConstraintsBlank, moved to package reconciler.io/dies/apis/scheduling/v1
var WorkloadPodGroupSchedulingConstraintsBlank = dieschedulingv1.WorkloadPodGroupSchedulingConstraintsBlank

// Deprecated WorkloadPodGroupResourceClaimDie, moved to package reconciler.io/dies/apis/scheduling/v1
type WorkloadPodGroupResourceClaimDie = dieschedulingv1.WorkloadPodGroupResourceClaimDie

// Deprecated WorkloadPodGroupResourceClaimBlank, moved to package reconciler.io/dies/apis/scheduling/v1
var WorkloadPodGroupResourceClaimBlank = dieschedulingv1.WorkloadPodGroupResourceClaimBlank

// Deprecated WorkloadPodGroupDisruptionModeDie, moved to package reconciler.io/dies/apis/scheduling/v1
type WorkloadPodGroupDisruptionModeDie = dieschedulingv1.WorkloadPodGroupDisruptionModeDie

// Deprecated WorkloadPodGroupDisruptionModeBlank, moved to package reconciler.io/dies/apis/scheduling/v1
var WorkloadPodGroupDisruptionModeBlank = dieschedulingv1.WorkloadPodGroupDisruptionModeBlank

// Deprecated WorkloadPodGroupSingleDisruptionModeDie, moved to package reconciler.io/dies/apis/scheduling/v1
type WorkloadPodGroupSingleDisruptionModeDie = dieschedulingv1.WorkloadPodGroupSingleDisruptionModeDie

// Deprecated WorkloadPodGroupSingleDisruptionModeBlank, moved to package reconciler.io/dies/apis/scheduling/v1
var WorkloadPodGroupSingleDisruptionModeBlank = dieschedulingv1.WorkloadPodGroupSingleDisruptionModeBlank

// Deprecated WorkloadPodGroupAllDisruptionModeDie, moved to package reconciler.io/dies/apis/scheduling/v1
type WorkloadPodGroupAllDisruptionModeDie = dieschedulingv1.WorkloadPodGroupAllDisruptionModeDie

// Deprecated WorkloadPodGroupAllDisruptionModeBlank, moved to package reconciler.io/dies/apis/scheduling/v1
var WorkloadPodGroupAllDisruptionModeBlank = dieschedulingv1.WorkloadPodGroupAllDisruptionModeBlank

// Deprecated WorkloadCompositePodGroupSchedulingPolicyDie, moved to package reconciler.io/dies/apis/scheduling/v1
type WorkloadCompositePodGroupSchedulingPolicyDie = dieschedulingv1.WorkloadCompositePodGroupSchedulingPolicyDie

// Deprecated WorkloadCompositePodGroupSchedulingPolicyBlank, moved to package reconciler.io/dies/apis/scheduling/v1
var WorkloadCompositePodGroupSchedulingPolicyBlank = dieschedulingv1.WorkloadCompositePodGroupSchedulingPolicyBlank

// Deprecated WorkloadCompositePodGroupBasicSchedulingPolicyDie, moved to package reconciler.io/dies/apis/scheduling/v1
type WorkloadCompositePodGroupBasicSchedulingPolicyDie = dieschedulingv1.WorkloadCompositePodGroupBasicSchedulingPolicyDie

// Deprecated WorkloadCompositePodGroupBasicSchedulingPolicyBlank, moved to package reconciler.io/dies/apis/scheduling/v1
var WorkloadCompositePodGroupBasicSchedulingPolicyBlank = dieschedulingv1.WorkloadCompositePodGroupBasicSchedulingPolicyBlank

// Deprecated WorkloadCompositePodGroupGangSchedulingPolicyDie, moved to package reconciler.io/dies/apis/scheduling/v1
type WorkloadCompositePodGroupGangSchedulingPolicyDie = dieschedulingv1.WorkloadCompositePodGroupGangSchedulingPolicyDie

// Deprecated WorkloadCompositePodGroupGangSchedulingPolicyBlank, moved to package reconciler.io/dies/apis/scheduling/v1
var WorkloadCompositePodGroupGangSchedulingPolicyBlank = dieschedulingv1.WorkloadCompositePodGroupGangSchedulingPolicyBlank

// Deprecated WorkloadCompositePodGroupSchedulingConstraintsDie, moved to package reconciler.io/dies/apis/scheduling/v1
type WorkloadCompositePodGroupSchedulingConstraintsDie = dieschedulingv1.WorkloadCompositePodGroupSchedulingConstraintsDie

// Deprecated WorkloadCompositePodGroupSchedulingConstraintsBlank, moved to package reconciler.io/dies/apis/scheduling/v1
var WorkloadCompositePodGroupSchedulingConstraintsBlank = dieschedulingv1.WorkloadCompositePodGroupSchedulingConstraintsBlank

// Deprecated WorkloadCompositePodGroupDisruptionModeDie, moved to package reconciler.io/dies/apis/scheduling/v1
type WorkloadCompositePodGroupDisruptionModeDie = dieschedulingv1.WorkloadCompositePodGroupDisruptionModeDie

// Deprecated WorkloadCompositePodGroupDisruptionModeBlank, moved to package reconciler.io/dies/apis/scheduling/v1
var WorkloadCompositePodGroupDisruptionModeBlank = dieschedulingv1.WorkloadCompositePodGroupDisruptionModeBlank

// Deprecated WorkloadCompositePodGroupSingleDisruptionModeDie, moved to package reconciler.io/dies/apis/scheduling/v1
type WorkloadCompositePodGroupSingleDisruptionModeDie = dieschedulingv1.WorkloadCompositePodGroupSingleDisruptionModeDie

// Deprecated WorkloadCompositePodGroupSingleDisruptionModeBlank, moved to package reconciler.io/dies/apis/scheduling/v1
var WorkloadCompositePodGroupSingleDisruptionModeBlank = dieschedulingv1.WorkloadCompositePodGroupSingleDisruptionModeBlank

// Deprecated WorkloadCompositePodGroupAllDisruptionModeDie, moved to package reconciler.io/dies/apis/scheduling/v1
type WorkloadCompositePodGroupAllDisruptionModeDie = dieschedulingv1.WorkloadCompositePodGroupAllDisruptionModeDie

// Deprecated WorkloadCompositePodGroupAllDisruptionModeBlank, moved to package reconciler.io/dies/apis/scheduling/v1
var WorkloadCompositePodGroupAllDisruptionModeBlank = dieschedulingv1.WorkloadCompositePodGroupAllDisruptionModeBlank
