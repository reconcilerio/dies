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
	diecorev1 "dies.dev/apis/core/v1"
	diemetav1 "dies.dev/apis/meta/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

// +die:object=true
type _ = batchv1.Job

// +die
type _ = batchv1.JobSpec

func (d *JobSpecDie) PodFailurePolicyDie(fn func(d *PodFailurePolicyDie)) *JobSpecDie {
	return d.DieStamp(func(r *batchv1.JobSpec) {
		d := PodFailurePolicyBlank.DieImmutable(false).DieFeedPtr(r.PodFailurePolicy)
		fn(d)
		r.PodFailurePolicy = d.DieReleasePtr()
	})
}

func (d *JobSpecDie) SelectorDie(fn func(d *diemetav1.LabelSelectorDie)) *JobSpecDie {
	return d.DieStamp(func(r *batchv1.JobSpec) {
		d := diemetav1.LabelSelectorBlank.DieImmutable(false).DieFeedPtr(r.Selector)
		fn(d)
		r.Selector = d.DieReleasePtr()
	})
}

func (d *JobSpecDie) TemplateDie(fn func(d *diecorev1.PodTemplateSpecDie)) *JobSpecDie {
	return d.DieStamp(func(r *batchv1.JobSpec) {
		d := diecorev1.PodTemplateSpecBlank.DieImmutable(false).DieFeed(r.Template)
		fn(d)
		r.Template = d.DieRelease()
	})
}

// +die
type _ = batchv1.PodFailurePolicy

func (d *PodFailurePolicyDie) RulesDie(rules ...*PodFailurePolicyRuleDie) *PodFailurePolicyDie {
	return d.DieStamp(func(r *batchv1.PodFailurePolicy) {
		r.Rules = make([]batchv1.PodFailurePolicyRule, len(rules))
		for i := range rules {
			r.Rules[i] = rules[i].DieRelease()
		}
	})
}

// +die
type _ = batchv1.PodFailurePolicyRule

func (d *PodFailurePolicyRuleDie) OnExitCodesDie(fn func(d *PodFailurePolicyOnExitCodesRequirementDie)) *PodFailurePolicyRuleDie {
	return d.DieStamp(func(r *batchv1.PodFailurePolicyRule) {
		d := PodFailurePolicyOnExitCodesRequirementBlank.DieImmutable(false).DieFeedPtr(r.OnExitCodes)
		fn(d)
		r.OnExitCodes = d.DieReleasePtr()
	})
}

func (d *PodFailurePolicyRuleDie) OnPodConditionsDie(patterns ...*PodFailurePolicyOnPodConditionsPatternDie) *PodFailurePolicyRuleDie {
	return d.DieStamp(func(r *batchv1.PodFailurePolicyRule) {
		r.OnPodConditions = make([]batchv1.PodFailurePolicyOnPodConditionsPattern, len(patterns))
		for i := range patterns {
			r.OnPodConditions[i] = patterns[i].DieRelease()
		}
	})
}

// +die
type _ = batchv1.PodFailurePolicyOnExitCodesRequirement

// +die
type _ = batchv1.PodFailurePolicyOnPodConditionsPattern

// +die
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

func (d *JobStatusDie) UncountedTerminatedPodsDie(fn func(d *UncountedTerminatedPodsDie)) *JobStatusDie {
	return d.DieStamp(func(r *batchv1.JobStatus) {
		d := UncountedTerminatedPodsBlank.DieImmutable(false).DieFeedPtr(r.UncountedTerminatedPods)
		fn(d)
		r.UncountedTerminatedPods = d.DieReleasePtr()
	})
}

// +die
type _ = batchv1.UncountedTerminatedPods
