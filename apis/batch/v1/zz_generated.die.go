//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by diegen. DO NOT EDIT.

package v1

import (
	metav1 "dies.dev/apis/meta/v1"
	json "encoding/json"
	fmtx "fmt"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
)

type CronJobDie struct {
	metav1.FrozenObjectMeta
	mutable bool
	r       batchv1.CronJob
}

var CronJobBlank = (&CronJobDie{}).DieFeed(batchv1.CronJob{})

func (d *CronJobDie) DieImmutable(immutable bool) *CronJobDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

func (d *CronJobDie) DieFeed(r batchv1.CronJob) *CronJobDie {
	if d.mutable {
		d.FrozenObjectMeta = metav1.FreezeObjectMeta(r.ObjectMeta)
		d.r = r
		return d
	}
	return &CronJobDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

func (d *CronJobDie) DieFeedPtr(r *batchv1.CronJob) *CronJobDie {
	if r == nil {
		r = &batchv1.CronJob{}
	}
	return d.DieFeed(*r)
}

func (d *CronJobDie) DieRelease() batchv1.CronJob {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

func (d *CronJobDie) DieReleasePtr() *batchv1.CronJob {
	r := d.DieRelease()
	return &r
}

func (d *CronJobDie) DieStamp(fn func(r *batchv1.CronJob)) *CronJobDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

func (d *CronJobDie) DeepCopy() *CronJobDie {
	r := *d.r.DeepCopy()
	return &CronJobDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

func (d *CronJobDie) DeepCopyObject() runtime.Object {
	return d.r.DeepCopy()
}

func (d *CronJobDie) GetObjectKind() schema.ObjectKind {
	r := d.DieRelease()
	return r.GetObjectKind()
}

func (d *CronJobDie) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.r)
}

func (d *CronJobDie) UnmarshalJSON(b []byte) error {
	if d == CronJobBlank {
		return fmtx.Errorf("cannot unmarshal into the root object, create a copy first")
	}
	r := &batchv1.CronJob{}
	err := json.Unmarshal(b, r)
	*d = *d.DieFeed(*r)
	return err
}

func (d *CronJobDie) MetadataDie(fn func(d *metav1.ObjectMetaDie)) *CronJobDie {
	return d.DieStamp(func(r *batchv1.CronJob) {
		d := metav1.ObjectMetaBlank.DieImmutable(false).DieFeed(r.ObjectMeta)
		fn(d)
		r.ObjectMeta = d.DieRelease()
	})
}

func (d *CronJobDie) SpecDie(fn func(d *CronJobSpecDie)) *CronJobDie {
	return d.DieStamp(func(r *batchv1.CronJob) {
		d := CronJobSpecBlank.DieImmutable(false).DieFeed(r.Spec)
		fn(d)
		r.Spec = d.DieRelease()
	})
}

func (d *CronJobDie) StatusDie(fn func(d *CronJobStatusDie)) *CronJobDie {
	return d.DieStamp(func(r *batchv1.CronJob) {
		d := CronJobStatusBlank.DieImmutable(false).DieFeed(r.Status)
		fn(d)
		r.Status = d.DieRelease()
	})
}

var _ apismetav1.Object = (*CronJobDie)(nil)
var _ apismetav1.ObjectMetaAccessor = (*CronJobDie)(nil)
var _ runtime.Object = (*CronJobDie)(nil)

// Specification of the desired behavior of a cron job, including the schedule. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
func (d *CronJobDie) Spec(v batchv1.CronJobSpec) *CronJobDie {
	return d.DieStamp(func(r *batchv1.CronJob) {
		r.Spec = v
	})
}

// Current status of a cron job. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
func (d *CronJobDie) Status(v batchv1.CronJobStatus) *CronJobDie {
	return d.DieStamp(func(r *batchv1.CronJob) {
		r.Status = v
	})
}

type CronJobSpecDie struct {
	mutable bool
	r       batchv1.CronJobSpec
}

var CronJobSpecBlank = (&CronJobSpecDie{}).DieFeed(batchv1.CronJobSpec{})

func (d *CronJobSpecDie) DieImmutable(immutable bool) *CronJobSpecDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

func (d *CronJobSpecDie) DieFeed(r batchv1.CronJobSpec) *CronJobSpecDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &CronJobSpecDie{
		mutable: d.mutable,
		r:       r,
	}
}

func (d *CronJobSpecDie) DieFeedPtr(r *batchv1.CronJobSpec) *CronJobSpecDie {
	if r == nil {
		r = &batchv1.CronJobSpec{}
	}
	return d.DieFeed(*r)
}

func (d *CronJobSpecDie) DieRelease() batchv1.CronJobSpec {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

func (d *CronJobSpecDie) DieReleasePtr() *batchv1.CronJobSpec {
	r := d.DieRelease()
	return &r
}

func (d *CronJobSpecDie) DieStamp(fn func(r *batchv1.CronJobSpec)) *CronJobSpecDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

func (d *CronJobSpecDie) DeepCopy() *CronJobSpecDie {
	r := *d.r.DeepCopy()
	return &CronJobSpecDie{
		mutable: d.mutable,
		r:       r,
	}
}

// The schedule in Cron format, see https://en.wikipedia.org/wiki/Cron.
func (d *CronJobSpecDie) Schedule(v string) *CronJobSpecDie {
	return d.DieStamp(func(r *batchv1.CronJobSpec) {
		r.Schedule = v
	})
}

// Optional deadline in seconds for starting the job if it misses scheduled time for any reason.  Missed jobs executions will be counted as failed ones.
func (d *CronJobSpecDie) StartingDeadlineSeconds(v *int64) *CronJobSpecDie {
	return d.DieStamp(func(r *batchv1.CronJobSpec) {
		r.StartingDeadlineSeconds = v
	})
}

// Specifies how to treat concurrent executions of a Job. Valid values are: - "Allow" (default): allows CronJobs to run concurrently; - "Forbid": forbids concurrent runs, skipping next run if previous run hasn't finished yet; - "Replace": cancels currently running job and replaces it with a new one
func (d *CronJobSpecDie) ConcurrencyPolicy(v batchv1.ConcurrencyPolicy) *CronJobSpecDie {
	return d.DieStamp(func(r *batchv1.CronJobSpec) {
		r.ConcurrencyPolicy = v
	})
}

// This flag tells the controller to suspend subsequent executions, it does not apply to already started executions.  Defaults to false.
func (d *CronJobSpecDie) Suspend(v *bool) *CronJobSpecDie {
	return d.DieStamp(func(r *batchv1.CronJobSpec) {
		r.Suspend = v
	})
}

// Specifies the job that will be created when executing a CronJob.
func (d *CronJobSpecDie) JobTemplate(v batchv1.JobTemplateSpec) *CronJobSpecDie {
	return d.DieStamp(func(r *batchv1.CronJobSpec) {
		r.JobTemplate = v
	})
}

// The number of successful finished jobs to retain. Value must be non-negative integer. Defaults to 3.
func (d *CronJobSpecDie) SuccessfulJobsHistoryLimit(v *int32) *CronJobSpecDie {
	return d.DieStamp(func(r *batchv1.CronJobSpec) {
		r.SuccessfulJobsHistoryLimit = v
	})
}

// The number of failed finished jobs to retain. Value must be non-negative integer. Defaults to 1.
func (d *CronJobSpecDie) FailedJobsHistoryLimit(v *int32) *CronJobSpecDie {
	return d.DieStamp(func(r *batchv1.CronJobSpec) {
		r.FailedJobsHistoryLimit = v
	})
}

type CronJobStatusDie struct {
	mutable bool
	r       batchv1.CronJobStatus
}

var CronJobStatusBlank = (&CronJobStatusDie{}).DieFeed(batchv1.CronJobStatus{})

func (d *CronJobStatusDie) DieImmutable(immutable bool) *CronJobStatusDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

func (d *CronJobStatusDie) DieFeed(r batchv1.CronJobStatus) *CronJobStatusDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &CronJobStatusDie{
		mutable: d.mutable,
		r:       r,
	}
}

func (d *CronJobStatusDie) DieFeedPtr(r *batchv1.CronJobStatus) *CronJobStatusDie {
	if r == nil {
		r = &batchv1.CronJobStatus{}
	}
	return d.DieFeed(*r)
}

func (d *CronJobStatusDie) DieRelease() batchv1.CronJobStatus {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

func (d *CronJobStatusDie) DieReleasePtr() *batchv1.CronJobStatus {
	r := d.DieRelease()
	return &r
}

func (d *CronJobStatusDie) DieStamp(fn func(r *batchv1.CronJobStatus)) *CronJobStatusDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

func (d *CronJobStatusDie) DeepCopy() *CronJobStatusDie {
	r := *d.r.DeepCopy()
	return &CronJobStatusDie{
		mutable: d.mutable,
		r:       r,
	}
}

// A list of pointers to currently running jobs.
func (d *CronJobStatusDie) Active(v ...corev1.ObjectReference) *CronJobStatusDie {
	return d.DieStamp(func(r *batchv1.CronJobStatus) {
		r.Active = v
	})
}

// Information when was the last time the job was successfully scheduled.
func (d *CronJobStatusDie) LastScheduleTime(v *apismetav1.Time) *CronJobStatusDie {
	return d.DieStamp(func(r *batchv1.CronJobStatus) {
		r.LastScheduleTime = v
	})
}

// Information when was the last time the job successfully completed.
func (d *CronJobStatusDie) LastSuccessfulTime(v *apismetav1.Time) *CronJobStatusDie {
	return d.DieStamp(func(r *batchv1.CronJobStatus) {
		r.LastSuccessfulTime = v
	})
}

type JobDie struct {
	metav1.FrozenObjectMeta
	mutable bool
	r       batchv1.Job
}

var JobBlank = (&JobDie{}).DieFeed(batchv1.Job{})

func (d *JobDie) DieImmutable(immutable bool) *JobDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

func (d *JobDie) DieFeed(r batchv1.Job) *JobDie {
	if d.mutable {
		d.FrozenObjectMeta = metav1.FreezeObjectMeta(r.ObjectMeta)
		d.r = r
		return d
	}
	return &JobDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

func (d *JobDie) DieFeedPtr(r *batchv1.Job) *JobDie {
	if r == nil {
		r = &batchv1.Job{}
	}
	return d.DieFeed(*r)
}

func (d *JobDie) DieRelease() batchv1.Job {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

func (d *JobDie) DieReleasePtr() *batchv1.Job {
	r := d.DieRelease()
	return &r
}

func (d *JobDie) DieStamp(fn func(r *batchv1.Job)) *JobDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

func (d *JobDie) DeepCopy() *JobDie {
	r := *d.r.DeepCopy()
	return &JobDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

func (d *JobDie) DeepCopyObject() runtime.Object {
	return d.r.DeepCopy()
}

func (d *JobDie) GetObjectKind() schema.ObjectKind {
	r := d.DieRelease()
	return r.GetObjectKind()
}

func (d *JobDie) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.r)
}

func (d *JobDie) UnmarshalJSON(b []byte) error {
	if d == JobBlank {
		return fmtx.Errorf("cannot unmarshal into the root object, create a copy first")
	}
	r := &batchv1.Job{}
	err := json.Unmarshal(b, r)
	*d = *d.DieFeed(*r)
	return err
}

func (d *JobDie) MetadataDie(fn func(d *metav1.ObjectMetaDie)) *JobDie {
	return d.DieStamp(func(r *batchv1.Job) {
		d := metav1.ObjectMetaBlank.DieImmutable(false).DieFeed(r.ObjectMeta)
		fn(d)
		r.ObjectMeta = d.DieRelease()
	})
}

func (d *JobDie) SpecDie(fn func(d *JobSpecDie)) *JobDie {
	return d.DieStamp(func(r *batchv1.Job) {
		d := JobSpecBlank.DieImmutable(false).DieFeed(r.Spec)
		fn(d)
		r.Spec = d.DieRelease()
	})
}

func (d *JobDie) StatusDie(fn func(d *JobStatusDie)) *JobDie {
	return d.DieStamp(func(r *batchv1.Job) {
		d := JobStatusBlank.DieImmutable(false).DieFeed(r.Status)
		fn(d)
		r.Status = d.DieRelease()
	})
}

var _ apismetav1.Object = (*JobDie)(nil)
var _ apismetav1.ObjectMetaAccessor = (*JobDie)(nil)
var _ runtime.Object = (*JobDie)(nil)

// Specification of the desired behavior of a job. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
func (d *JobDie) Spec(v batchv1.JobSpec) *JobDie {
	return d.DieStamp(func(r *batchv1.Job) {
		r.Spec = v
	})
}

// Current status of a job. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
func (d *JobDie) Status(v batchv1.JobStatus) *JobDie {
	return d.DieStamp(func(r *batchv1.Job) {
		r.Status = v
	})
}

type JobSpecDie struct {
	mutable bool
	r       batchv1.JobSpec
}

var JobSpecBlank = (&JobSpecDie{}).DieFeed(batchv1.JobSpec{})

func (d *JobSpecDie) DieImmutable(immutable bool) *JobSpecDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

func (d *JobSpecDie) DieFeed(r batchv1.JobSpec) *JobSpecDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &JobSpecDie{
		mutable: d.mutable,
		r:       r,
	}
}

func (d *JobSpecDie) DieFeedPtr(r *batchv1.JobSpec) *JobSpecDie {
	if r == nil {
		r = &batchv1.JobSpec{}
	}
	return d.DieFeed(*r)
}

func (d *JobSpecDie) DieRelease() batchv1.JobSpec {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

func (d *JobSpecDie) DieReleasePtr() *batchv1.JobSpec {
	r := d.DieRelease()
	return &r
}

func (d *JobSpecDie) DieStamp(fn func(r *batchv1.JobSpec)) *JobSpecDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

func (d *JobSpecDie) DeepCopy() *JobSpecDie {
	r := *d.r.DeepCopy()
	return &JobSpecDie{
		mutable: d.mutable,
		r:       r,
	}
}

// Specifies the maximum desired number of pods the job should run at any given time. The actual number of pods running in steady state will be less than this number when ((.spec.completions - .status.successful) < .spec.parallelism), i.e. when the work left to do is less than max parallelism. More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/
func (d *JobSpecDie) Parallelism(v *int32) *JobSpecDie {
	return d.DieStamp(func(r *batchv1.JobSpec) {
		r.Parallelism = v
	})
}

// Specifies the desired number of successfully finished pods the job should be run with.  Setting to nil means that the success of any pod signals the success of all pods, and allows parallelism to have any positive value.  Setting to 1 means that parallelism is limited to 1 and the success of that pod signals the success of the job. More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/
func (d *JobSpecDie) Completions(v *int32) *JobSpecDie {
	return d.DieStamp(func(r *batchv1.JobSpec) {
		r.Completions = v
	})
}

// Specifies the duration in seconds relative to the startTime that the job may be continuously active before the system tries to terminate it; value must be positive integer. If a Job is suspended (at creation or through an update), this timer will effectively be stopped and reset when the Job is resumed again.
func (d *JobSpecDie) ActiveDeadlineSeconds(v *int64) *JobSpecDie {
	return d.DieStamp(func(r *batchv1.JobSpec) {
		r.ActiveDeadlineSeconds = v
	})
}

// Specifies the number of retries before marking this job failed. Defaults to 6
func (d *JobSpecDie) BackoffLimit(v *int32) *JobSpecDie {
	return d.DieStamp(func(r *batchv1.JobSpec) {
		r.BackoffLimit = v
	})
}

// A label query over pods that should match the pod count. Normally, the system sets this field for you. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors
func (d *JobSpecDie) Selector(v *apismetav1.LabelSelector) *JobSpecDie {
	return d.DieStamp(func(r *batchv1.JobSpec) {
		r.Selector = v
	})
}

// manualSelector controls generation of pod labels and pod selectors. Leave `manualSelector` unset unless you are certain what you are doing. When false or unset, the system pick labels unique to this job and appends those labels to the pod template.  When true, the user is responsible for picking unique labels and specifying the selector.  Failure to pick a unique label may cause this and other jobs to not function correctly.  However, You may see `manualSelector=true` in jobs that were created with the old `extensions/v1beta1` API. More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/#specifying-your-own-pod-selector
func (d *JobSpecDie) ManualSelector(v *bool) *JobSpecDie {
	return d.DieStamp(func(r *batchv1.JobSpec) {
		r.ManualSelector = v
	})
}

// Describes the pod that will be created when executing a job. More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/
func (d *JobSpecDie) Template(v corev1.PodTemplateSpec) *JobSpecDie {
	return d.DieStamp(func(r *batchv1.JobSpec) {
		r.Template = v
	})
}

// ttlSecondsAfterFinished limits the lifetime of a Job that has finished execution (either Complete or Failed). If this field is set, ttlSecondsAfterFinished after the Job finishes, it is eligible to be automatically deleted. When the Job is being deleted, its lifecycle guarantees (e.g. finalizers) will be honored. If this field is unset, the Job won't be automatically deleted. If this field is set to zero, the Job becomes eligible to be deleted immediately after it finishes. This field is alpha-level and is only honored by servers that enable the TTLAfterFinished feature.
func (d *JobSpecDie) TTLSecondsAfterFinished(v *int32) *JobSpecDie {
	return d.DieStamp(func(r *batchv1.JobSpec) {
		r.TTLSecondsAfterFinished = v
	})
}

// CompletionMode specifies how Pod completions are tracked. It can be `NonIndexed` (default) or `Indexed`.
//
// `NonIndexed` means that the Job is considered complete when there have been .spec.completions successfully completed Pods. Each Pod completion is homologous to each other.
//
// `Indexed` means that the Pods of a Job get an associated completion index from 0 to (.spec.completions - 1), available in the annotation batch.kubernetes.io/job-completion-index. The Job is considered complete when there is one successfully completed Pod for each index. When value is `Indexed`, .spec.completions must be specified and `.spec.parallelism` must be less than or equal to 10^5. In addition, The Pod name takes the form `$(job-name)-$(index)-$(random-string)`, the Pod hostname takes the form `$(job-name)-$(index)`.
//
// This field is beta-level. More completion modes can be added in the future. If the Job controller observes a mode that it doesn't recognize, the controller skips updates for the Job.
func (d *JobSpecDie) CompletionMode(v *batchv1.CompletionMode) *JobSpecDie {
	return d.DieStamp(func(r *batchv1.JobSpec) {
		r.CompletionMode = v
	})
}

// Suspend specifies whether the Job controller should create Pods or not. If a Job is created with suspend set to true, no Pods are created by the Job controller. If a Job is suspended after creation (i.e. the flag goes from false to true), the Job controller will delete all active Pods associated with this Job. Users must design their workload to gracefully handle this. Suspending a Job will reset the StartTime field of the Job, effectively resetting the ActiveDeadlineSeconds timer too. Defaults to false.
//
// This field is beta-level, gated by SuspendJob feature flag (enabled by default).
func (d *JobSpecDie) Suspend(v *bool) *JobSpecDie {
	return d.DieStamp(func(r *batchv1.JobSpec) {
		r.Suspend = v
	})
}

type JobStatusDie struct {
	mutable bool
	r       batchv1.JobStatus
}

var JobStatusBlank = (&JobStatusDie{}).DieFeed(batchv1.JobStatus{})

func (d *JobStatusDie) DieImmutable(immutable bool) *JobStatusDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

func (d *JobStatusDie) DieFeed(r batchv1.JobStatus) *JobStatusDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &JobStatusDie{
		mutable: d.mutable,
		r:       r,
	}
}

func (d *JobStatusDie) DieFeedPtr(r *batchv1.JobStatus) *JobStatusDie {
	if r == nil {
		r = &batchv1.JobStatus{}
	}
	return d.DieFeed(*r)
}

func (d *JobStatusDie) DieRelease() batchv1.JobStatus {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

func (d *JobStatusDie) DieReleasePtr() *batchv1.JobStatus {
	r := d.DieRelease()
	return &r
}

func (d *JobStatusDie) DieStamp(fn func(r *batchv1.JobStatus)) *JobStatusDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

func (d *JobStatusDie) DeepCopy() *JobStatusDie {
	r := *d.r.DeepCopy()
	return &JobStatusDie{
		mutable: d.mutable,
		r:       r,
	}
}

// The latest available observations of an object's current state. When a Job fails, one of the conditions will have type "Failed" and status true. When a Job is suspended, one of the conditions will have type "Suspended" and status true; when the Job is resumed, the status of this condition will become false. When a Job is completed, one of the conditions will have type "Complete" and status true. More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/
func (d *JobStatusDie) Conditions(v ...batchv1.JobCondition) *JobStatusDie {
	return d.DieStamp(func(r *batchv1.JobStatus) {
		r.Conditions = v
	})
}

// Represents time when the job controller started processing a job. When a Job is created in the suspended state, this field is not set until the first time it is resumed. This field is reset every time a Job is resumed from suspension. It is represented in RFC3339 form and is in UTC.
func (d *JobStatusDie) StartTime(v *apismetav1.Time) *JobStatusDie {
	return d.DieStamp(func(r *batchv1.JobStatus) {
		r.StartTime = v
	})
}

// Represents time when the job was completed. It is not guaranteed to be set in happens-before order across separate operations. It is represented in RFC3339 form and is in UTC. The completion time is only set when the job finishes successfully.
func (d *JobStatusDie) CompletionTime(v *apismetav1.Time) *JobStatusDie {
	return d.DieStamp(func(r *batchv1.JobStatus) {
		r.CompletionTime = v
	})
}

// The number of actively running pods.
func (d *JobStatusDie) Active(v int32) *JobStatusDie {
	return d.DieStamp(func(r *batchv1.JobStatus) {
		r.Active = v
	})
}

// The number of pods which reached phase Succeeded.
func (d *JobStatusDie) Succeeded(v int32) *JobStatusDie {
	return d.DieStamp(func(r *batchv1.JobStatus) {
		r.Succeeded = v
	})
}

// The number of pods which reached phase Failed.
func (d *JobStatusDie) Failed(v int32) *JobStatusDie {
	return d.DieStamp(func(r *batchv1.JobStatus) {
		r.Failed = v
	})
}

// CompletedIndexes holds the completed indexes when .spec.completionMode = "Indexed" in a text format. The indexes are represented as decimal integers separated by commas. The numbers are listed in increasing order. Three or more consecutive numbers are compressed and represented by the first and last element of the series, separated by a hyphen. For example, if the completed indexes are 1, 3, 4, 5 and 7, they are represented as "1,3-5,7".
func (d *JobStatusDie) CompletedIndexes(v string) *JobStatusDie {
	return d.DieStamp(func(r *batchv1.JobStatus) {
		r.CompletedIndexes = v
	})
}

// UncountedTerminatedPods holds the UIDs of Pods that have terminated but the job controller hasn't yet accounted for in the status counters.
//
// The job controller creates pods with a finalizer. When a pod terminates (succeeded or failed), the controller does three steps to account for it in the job status: (1) Add the pod UID to the arrays in this field. (2) Remove the pod finalizer. (3) Remove the pod UID from the arrays while increasing the corresponding     counter.
//
// This field is alpha-level. The job controller only makes use of this field when the feature gate PodTrackingWithFinalizers is enabled. Old jobs might not be tracked using this field, in which case the field remains null.
func (d *JobStatusDie) UncountedTerminatedPods(v *batchv1.UncountedTerminatedPods) *JobStatusDie {
	return d.DieStamp(func(r *batchv1.JobStatus) {
		r.UncountedTerminatedPods = v
	})
}
