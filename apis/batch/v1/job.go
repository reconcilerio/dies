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
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	diemetav1 "reconciler.io/dies/apis/meta/v1"
)

// +die:object=true,apiVersion=batch/v1,kind=Job
type _ = batchv1.Job

// +die
// +die:field:name=PodFailurePolicy,die=PodFailurePolicyDie,pointer=true
// +die:field:name=SuccessPolicy,die=SuccessPolicyDie,pointer=true
// +die:field:name=Selector,package=_/meta/v1,die=LabelSelectorDie,pointer=true
// +die:field:name=Template,package=_/core/v1,die=PodTemplateSpecDie
type _ = batchv1.JobSpec

// +die
// +die:field:name=Rules,die=PodFailurePolicyRuleDie,listType=atomic
type _ = batchv1.PodFailurePolicy

// +die
// +die:field:name=OnExitCodes,die=PodFailurePolicyOnExitCodesRequirementDie,pointer=true
// +die:field:name=OnPodConditions,die=PodFailurePolicyOnPodConditionsPatternDie,listType=atomic
type _ = batchv1.PodFailurePolicyRule

// +die
type _ = batchv1.PodFailurePolicyOnExitCodesRequirement

// +die
type _ = batchv1.PodFailurePolicyOnPodConditionsPattern

// +die
// +die:field:name=Rules,die=SuccessPolicyRuleDie,listType=atomic
type _ = batchv1.SuccessPolicy

// +die
type _ = batchv1.SuccessPolicyRule

// +die
// +die:field:name=UncountedTerminatedPods,die=UncountedTerminatedPodsDie,pointer=true
type _ = batchv1.JobStatus

func (d *JobStatusDie) ConditionsDie(conditions ...*diemetav1.ConditionDie) *JobStatusDie {
	return d.DieStamp(func(r *batchv1.JobStatus) {
		r.Conditions = make([]batchv1.JobCondition, len(conditions))
		for i := range conditions {
			c := conditions[i].DieRelease()
			r.Conditions[i] = batchv1.JobCondition{
				Type:               batchv1.JobConditionType(c.Type),
				Status:             corev1.ConditionStatus(c.Status),
				Reason:             c.Reason,
				Message:            c.Message,
				LastTransitionTime: c.LastTransitionTime,
			}
		}
	})
}

// +die
type _ = batchv1.UncountedTerminatedPods
