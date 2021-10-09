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
	diecorev1 "github.com/scothis/dies/apis/core/v1"
	diemetav1 "github.com/scothis/dies/apis/meta/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

// +die:target=k8s.io/api/batch/v1.Job,object=true

// +die:target=k8s.io/api/batch/v1.JobSpec
// +die:field:receiver=JobSpecDie,name=Parallelism,type=*int32
// +die:field:receiver=JobSpecDie,name=Completions,type=*int32
// +die:field:receiver=JobSpecDie,name=ActiveDeadlineSeconds,type=*int64
// +die:field:receiver=JobSpecDie,name=BackoffLimit,type=*int32
// +die:field:receiver=JobSpecDie,name=Selector,type=*k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector
// +die:field:receiver=JobSpecDie,name=ManualSelector,type=*bool
// +die:field:receiver=JobSpecDie,name=Template,type=k8s.io/api/core/v1.PodTemplateSpec
// +die:field:receiver=JobSpecDie,name=TTLSecondsAfterFinished,type=*int32
// +die:field:receiver=JobSpecDie,name=CompletionMode,type=*k8s.io/api/batch/v1.CompletionMode
// +die:field:receiver=JobSpecDie,name=Suspend,type=*bool

func (d *JobSpecDie) TemplateDie(fn func(d *diecorev1.PodTemplateSpecDie)) *JobSpecDie {
	return d.DieStamp(func(r *batchv1.JobSpec) {
		d := diecorev1.PodTemplateSpecBlank.DieImmutable(false).DieFeed(r.Template)
		fn(d)
		r.Template = d.DieRelease()
	})
}

// +die:target=k8s.io/api/batch/v1.JobStatus
// +die:field:receiver=JobStatusDie,name=Conditions,type=[]k8s.io/api/batch/v1.JobCondition
// +die:field:receiver=JobStatusDie,name=StartTime,type=*k8s.io/apimachinery/pkg/apis/meta/v1.Time
// +die:field:receiver=JobStatusDie,name=CompletionTime,type=*k8s.io/apimachinery/pkg/apis/meta/v1.Time
// +die:field:receiver=JobStatusDie,name=Active,type=int32
// +die:field:receiver=JobStatusDie,name=Succeeded,type=int32
// +die:field:receiver=JobStatusDie,name=Failed,type=int32
// +die:field:receiver=JobStatusDie,name=CompletedIndexes,type=string
// +die:field:receiver=JobStatusDie,name=UncountedTerminatedPods,type=*k8s.io/api/batch/v1.UncountedTerminatedPods

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
