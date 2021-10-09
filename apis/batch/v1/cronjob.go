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
)

// +die:target=k8s.io/api/batch/v1.CronJob,object=true

// +die:target=k8s.io/api/batch/v1.CronJobSpec
// +die:field:receiver=CronJobSpecDie,name=Schedule,type=string
// +die:field:receiver=CronJobSpecDie,name=StartingDeadlineSeconds,type=*int64
// +die:field:receiver=CronJobSpecDie,name=ConcurrencyPolicy,type=k8s.io/api/batch/v1.ConcurrencyPolicy
// +die:field:receiver=CronJobSpecDie,name=Suspend,type=*bool
// +die:field:receiver=CronJobSpecDie,name=JobTemplate,type=k8s.io/api/batch/v1.JobTemplateSpec
// +die:field:receiver=CronJobSpecDie,name=SuccessfulJobsHistoryLimit,type=*int32
// +die:field:receiver=CronJobSpecDie,name=FailedJobsHistoryLimit,type=*int32

func (d *CronJobSpecDie) JobTemplateDie(fn func(d *JobDie)) *CronJobSpecDie {
	return d.DieStamp(func(r *batchv1.CronJobSpec) {
		d := JobBlank.DieImmutable(false).DieFeed(batchv1.Job{
			ObjectMeta: r.JobTemplate.ObjectMeta,
			Spec:       r.JobTemplate.Spec,
		})
		fn(d)
		r.JobTemplate.ObjectMeta = d.DieRelease().ObjectMeta
		r.JobTemplate.Spec = d.DieRelease().Spec
	})
}

// +die:target=k8s.io/api/batch/v1.CronJobStatus
// +die:field:receiver=CronJobStatusDie,name=Active,type=[]k8s.io/api/core/v1.ObjectReference
// +die:field:receiver=CronJobStatusDie,name=LastScheduleTime,type=*k8s.io/apimachinery/pkg/apis/meta/v1.Time
// +die:field:receiver=CronJobStatusDie,name=LastSuccessfulTime,type=*k8s.io/apimachinery/pkg/apis/meta/v1.Time
