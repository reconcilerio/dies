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

type CronJobDie interface {
	// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
	DieStamp(fn func(r *batchv1.CronJob)) CronJobDie
	// DieFeed returns a new die with the provided resource.
	DieFeed(r batchv1.CronJob) CronJobDie
	// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
	DieFeedPtr(r *batchv1.CronJob) CronJobDie
	// DieRelease returns the resource managed by the die.
	DieRelease() batchv1.CronJob
	// DieReleasePtr returns a pointer to the resource managed by the die.
	DieReleasePtr() *batchv1.CronJob
	// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
	DieImmutable(immutable bool) CronJobDie
	// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
	DeepCopy() CronJobDie

	// MetadataDie stamps the resource's ObjectMeta field with a mutable die.
	MetadataDie(fn func(d metav1.ObjectMetaDie)) CronJobDie
	// SpecDie stamps the resource's spec field with a mutable die.
	SpecDie(fn func(d CronJobSpecDie)) CronJobDie
	// StatusDie stamps the resource's status field with a mutable die.
	StatusDie(fn func(d CronJobStatusDie)) CronJobDie
	// Specification of the desired behavior of a cron job, including the schedule. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	Spec(Spec batchv1.CronJobSpec) CronJobDie
	// Current status of a cron job. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	Status(Status batchv1.CronJobStatus) CronJobDie

	runtime.Object
	apismetav1.Object
	apismetav1.ObjectMetaAccessor
}

var _ CronJobDie = (*cronJobDie)(nil)
var CronJobBlank = (&cronJobDie{}).DieFeed(batchv1.CronJob{})

type cronJobDie struct {
	metav1.FrozenObjectMeta
	mutable bool
	r       batchv1.CronJob
}

func (d *cronJobDie) DieImmutable(immutable bool) CronJobDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy().(*cronJobDie)
	d.mutable = !immutable
	return d
}

func (d *cronJobDie) DieFeed(r batchv1.CronJob) CronJobDie {
	if d.mutable {
		d.FrozenObjectMeta = metav1.FreezeObjectMeta(r.ObjectMeta)
		d.r = r
		return d
	}
	return &cronJobDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

func (d *cronJobDie) DieFeedPtr(r *batchv1.CronJob) CronJobDie {
	if r == nil {
		r = &batchv1.CronJob{}
	}
	return d.DieFeed(*r)
}

func (d *cronJobDie) DieRelease() batchv1.CronJob {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

func (d *cronJobDie) DieReleasePtr() *batchv1.CronJob {
	r := d.DieRelease()
	return &r
}

func (d *cronJobDie) DieStamp(fn func(r *batchv1.CronJob)) CronJobDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

func (d *cronJobDie) DeepCopy() CronJobDie {
	r := *d.r.DeepCopy()
	return &cronJobDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

func (d *cronJobDie) DeepCopyObject() runtime.Object {
	return d.r.DeepCopy()
}

func (d *cronJobDie) GetObjectKind() schema.ObjectKind {
	r := d.DieRelease()
	return r.GetObjectKind()
}

func (d *cronJobDie) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.r)
}

func (d *cronJobDie) UnmarshalJSON(b []byte) error {
	if d == CronJobBlank {
		return fmtx.Errorf("cannot unmarshal into the root object, create a copy first")
	}
	r := &batchv1.CronJob{}
	err := json.Unmarshal(b, r)
	*d = *d.DieFeed(*r).(*cronJobDie)
	return err
}

func (d *cronJobDie) MetadataDie(fn func(d metav1.ObjectMetaDie)) CronJobDie {
	return d.DieStamp(func(r *batchv1.CronJob) {
		d := metav1.ObjectMetaBlank.DieImmutable(false).DieFeed(r.ObjectMeta)
		fn(d)
		r.ObjectMeta = d.DieRelease()
	})
}

func (d *cronJobDie) SpecDie(fn func(d CronJobSpecDie)) CronJobDie {
	return d.DieStamp(func(r *batchv1.CronJob) {
		d := CronJobSpecBlank.DieImmutable(false).DieFeed(r.Spec)
		fn(d)
		r.Spec = d.DieRelease()
	})
}

func (d *cronJobDie) StatusDie(fn func(d CronJobStatusDie)) CronJobDie {
	return d.DieStamp(func(r *batchv1.CronJob) {
		d := CronJobStatusBlank.DieImmutable(false).DieFeed(r.Status)
		fn(d)
		r.Status = d.DieRelease()
	})
}

func (d *cronJobDie) Spec(v batchv1.CronJobSpec) CronJobDie {
	return d.DieStamp(func(r *batchv1.CronJob) {
		r.Spec = v
	})
}

func (d *cronJobDie) Status(v batchv1.CronJobStatus) CronJobDie {
	return d.DieStamp(func(r *batchv1.CronJob) {
		r.Status = v
	})
}

type CronJobSpecDie interface {
	// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
	DieStamp(fn func(r *batchv1.CronJobSpec)) CronJobSpecDie
	// DieFeed returns a new die with the provided resource.
	DieFeed(r batchv1.CronJobSpec) CronJobSpecDie
	// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
	DieFeedPtr(r *batchv1.CronJobSpec) CronJobSpecDie
	// DieRelease returns the resource managed by the die.
	DieRelease() batchv1.CronJobSpec
	// DieReleasePtr returns a pointer to the resource managed by the die.
	DieReleasePtr() *batchv1.CronJobSpec
	// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
	DieImmutable(immutable bool) CronJobSpecDie
	// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
	DeepCopy() CronJobSpecDie

	cronJobSpec
	// The schedule in Cron format, see https://en.wikipedia.org/wiki/Cron.
	Schedule(Schedule string) CronJobSpecDie
	// Optional deadline in seconds for starting the job if it misses scheduled time for any reason.  Missed jobs executions will be counted as failed ones.
	StartingDeadlineSeconds(StartingDeadlineSeconds *int64) CronJobSpecDie
	// Specifies how to treat concurrent executions of a Job. Valid values are: - "Allow" (default): allows CronJobs to run concurrently; - "Forbid": forbids concurrent runs, skipping next run if previous run hasn't finished yet; - "Replace": cancels currently running job and replaces it with a new one
	ConcurrencyPolicy(ConcurrencyPolicy batchv1.ConcurrencyPolicy) CronJobSpecDie
	// This flag tells the controller to suspend subsequent executions, it does not apply to already started executions.  Defaults to false.
	Suspend(Suspend *bool) CronJobSpecDie
	// Specifies the job that will be created when executing a CronJob.
	JobTemplate(JobTemplate batchv1.JobTemplateSpec) CronJobSpecDie
	// The number of successful finished jobs to retain. Value must be non-negative integer. Defaults to 3.
	SuccessfulJobsHistoryLimit(SuccessfulJobsHistoryLimit *int32) CronJobSpecDie
	// The number of failed finished jobs to retain. Value must be non-negative integer. Defaults to 1.
	FailedJobsHistoryLimit(FailedJobsHistoryLimit *int32) CronJobSpecDie
}

var _ CronJobSpecDie = (*cronJobSpecDie)(nil)
var CronJobSpecBlank = (&cronJobSpecDie{}).DieFeed(batchv1.CronJobSpec{})

type cronJobSpecDie struct {
	mutable bool
	r       batchv1.CronJobSpec
}

func (d *cronJobSpecDie) DieImmutable(immutable bool) CronJobSpecDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy().(*cronJobSpecDie)
	d.mutable = !immutable
	return d
}

func (d *cronJobSpecDie) DieFeed(r batchv1.CronJobSpec) CronJobSpecDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &cronJobSpecDie{
		mutable: d.mutable,
		r:       r,
	}
}

func (d *cronJobSpecDie) DieFeedPtr(r *batchv1.CronJobSpec) CronJobSpecDie {
	if r == nil {
		r = &batchv1.CronJobSpec{}
	}
	return d.DieFeed(*r)
}

func (d *cronJobSpecDie) DieRelease() batchv1.CronJobSpec {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

func (d *cronJobSpecDie) DieReleasePtr() *batchv1.CronJobSpec {
	r := d.DieRelease()
	return &r
}

func (d *cronJobSpecDie) DieStamp(fn func(r *batchv1.CronJobSpec)) CronJobSpecDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

func (d *cronJobSpecDie) DeepCopy() CronJobSpecDie {
	r := *d.r.DeepCopy()
	return &cronJobSpecDie{
		mutable: d.mutable,
		r:       r,
	}
}

func (d *cronJobSpecDie) Schedule(v string) CronJobSpecDie {
	return d.DieStamp(func(r *batchv1.CronJobSpec) {
		r.Schedule = v
	})
}

func (d *cronJobSpecDie) StartingDeadlineSeconds(v *int64) CronJobSpecDie {
	return d.DieStamp(func(r *batchv1.CronJobSpec) {
		r.StartingDeadlineSeconds = v
	})
}

func (d *cronJobSpecDie) ConcurrencyPolicy(v batchv1.ConcurrencyPolicy) CronJobSpecDie {
	return d.DieStamp(func(r *batchv1.CronJobSpec) {
		r.ConcurrencyPolicy = v
	})
}

func (d *cronJobSpecDie) Suspend(v *bool) CronJobSpecDie {
	return d.DieStamp(func(r *batchv1.CronJobSpec) {
		r.Suspend = v
	})
}

func (d *cronJobSpecDie) JobTemplate(v batchv1.JobTemplateSpec) CronJobSpecDie {
	return d.DieStamp(func(r *batchv1.CronJobSpec) {
		r.JobTemplate = v
	})
}

func (d *cronJobSpecDie) SuccessfulJobsHistoryLimit(v *int32) CronJobSpecDie {
	return d.DieStamp(func(r *batchv1.CronJobSpec) {
		r.SuccessfulJobsHistoryLimit = v
	})
}

func (d *cronJobSpecDie) FailedJobsHistoryLimit(v *int32) CronJobSpecDie {
	return d.DieStamp(func(r *batchv1.CronJobSpec) {
		r.FailedJobsHistoryLimit = v
	})
}

type CronJobStatusDie interface {
	// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
	DieStamp(fn func(r *batchv1.CronJobStatus)) CronJobStatusDie
	// DieFeed returns a new die with the provided resource.
	DieFeed(r batchv1.CronJobStatus) CronJobStatusDie
	// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
	DieFeedPtr(r *batchv1.CronJobStatus) CronJobStatusDie
	// DieRelease returns the resource managed by the die.
	DieRelease() batchv1.CronJobStatus
	// DieReleasePtr returns a pointer to the resource managed by the die.
	DieReleasePtr() *batchv1.CronJobStatus
	// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
	DieImmutable(immutable bool) CronJobStatusDie
	// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
	DeepCopy() CronJobStatusDie

	// A list of pointers to currently running jobs.
	Active(Active ...corev1.ObjectReference) CronJobStatusDie
	// Information when was the last time the job was successfully scheduled.
	LastScheduleTime(LastScheduleTime *apismetav1.Time) CronJobStatusDie
	// Information when was the last time the job successfully completed.
	LastSuccessfulTime(LastSuccessfulTime *apismetav1.Time) CronJobStatusDie
}

var _ CronJobStatusDie = (*cronJobStatusDie)(nil)
var CronJobStatusBlank = (&cronJobStatusDie{}).DieFeed(batchv1.CronJobStatus{})

type cronJobStatusDie struct {
	mutable bool
	r       batchv1.CronJobStatus
}

func (d *cronJobStatusDie) DieImmutable(immutable bool) CronJobStatusDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy().(*cronJobStatusDie)
	d.mutable = !immutable
	return d
}

func (d *cronJobStatusDie) DieFeed(r batchv1.CronJobStatus) CronJobStatusDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &cronJobStatusDie{
		mutable: d.mutable,
		r:       r,
	}
}

func (d *cronJobStatusDie) DieFeedPtr(r *batchv1.CronJobStatus) CronJobStatusDie {
	if r == nil {
		r = &batchv1.CronJobStatus{}
	}
	return d.DieFeed(*r)
}

func (d *cronJobStatusDie) DieRelease() batchv1.CronJobStatus {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

func (d *cronJobStatusDie) DieReleasePtr() *batchv1.CronJobStatus {
	r := d.DieRelease()
	return &r
}

func (d *cronJobStatusDie) DieStamp(fn func(r *batchv1.CronJobStatus)) CronJobStatusDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

func (d *cronJobStatusDie) DeepCopy() CronJobStatusDie {
	r := *d.r.DeepCopy()
	return &cronJobStatusDie{
		mutable: d.mutable,
		r:       r,
	}
}

func (d *cronJobStatusDie) Active(v ...corev1.ObjectReference) CronJobStatusDie {
	return d.DieStamp(func(r *batchv1.CronJobStatus) {
		r.Active = v
	})
}

func (d *cronJobStatusDie) LastScheduleTime(v *apismetav1.Time) CronJobStatusDie {
	return d.DieStamp(func(r *batchv1.CronJobStatus) {
		r.LastScheduleTime = v
	})
}

func (d *cronJobStatusDie) LastSuccessfulTime(v *apismetav1.Time) CronJobStatusDie {
	return d.DieStamp(func(r *batchv1.CronJobStatus) {
		r.LastSuccessfulTime = v
	})
}

type JobDie interface {
	// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
	DieStamp(fn func(r *batchv1.Job)) JobDie
	// DieFeed returns a new die with the provided resource.
	DieFeed(r batchv1.Job) JobDie
	// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
	DieFeedPtr(r *batchv1.Job) JobDie
	// DieRelease returns the resource managed by the die.
	DieRelease() batchv1.Job
	// DieReleasePtr returns a pointer to the resource managed by the die.
	DieReleasePtr() *batchv1.Job
	// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
	DieImmutable(immutable bool) JobDie
	// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
	DeepCopy() JobDie

	// MetadataDie stamps the resource's ObjectMeta field with a mutable die.
	MetadataDie(fn func(d metav1.ObjectMetaDie)) JobDie
	// SpecDie stamps the resource's spec field with a mutable die.
	SpecDie(fn func(d JobSpecDie)) JobDie
	// StatusDie stamps the resource's status field with a mutable die.
	StatusDie(fn func(d JobStatusDie)) JobDie
	// Specification of the desired behavior of a job. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	Spec(Spec batchv1.JobSpec) JobDie
	// Current status of a job. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	Status(Status batchv1.JobStatus) JobDie

	runtime.Object
	apismetav1.Object
	apismetav1.ObjectMetaAccessor
}

var _ JobDie = (*jobDie)(nil)
var JobBlank = (&jobDie{}).DieFeed(batchv1.Job{})

type jobDie struct {
	metav1.FrozenObjectMeta
	mutable bool
	r       batchv1.Job
}

func (d *jobDie) DieImmutable(immutable bool) JobDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy().(*jobDie)
	d.mutable = !immutable
	return d
}

func (d *jobDie) DieFeed(r batchv1.Job) JobDie {
	if d.mutable {
		d.FrozenObjectMeta = metav1.FreezeObjectMeta(r.ObjectMeta)
		d.r = r
		return d
	}
	return &jobDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

func (d *jobDie) DieFeedPtr(r *batchv1.Job) JobDie {
	if r == nil {
		r = &batchv1.Job{}
	}
	return d.DieFeed(*r)
}

func (d *jobDie) DieRelease() batchv1.Job {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

func (d *jobDie) DieReleasePtr() *batchv1.Job {
	r := d.DieRelease()
	return &r
}

func (d *jobDie) DieStamp(fn func(r *batchv1.Job)) JobDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

func (d *jobDie) DeepCopy() JobDie {
	r := *d.r.DeepCopy()
	return &jobDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

func (d *jobDie) DeepCopyObject() runtime.Object {
	return d.r.DeepCopy()
}

func (d *jobDie) GetObjectKind() schema.ObjectKind {
	r := d.DieRelease()
	return r.GetObjectKind()
}

func (d *jobDie) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.r)
}

func (d *jobDie) UnmarshalJSON(b []byte) error {
	if d == JobBlank {
		return fmtx.Errorf("cannot unmarshal into the root object, create a copy first")
	}
	r := &batchv1.Job{}
	err := json.Unmarshal(b, r)
	*d = *d.DieFeed(*r).(*jobDie)
	return err
}

func (d *jobDie) MetadataDie(fn func(d metav1.ObjectMetaDie)) JobDie {
	return d.DieStamp(func(r *batchv1.Job) {
		d := metav1.ObjectMetaBlank.DieImmutable(false).DieFeed(r.ObjectMeta)
		fn(d)
		r.ObjectMeta = d.DieRelease()
	})
}

func (d *jobDie) SpecDie(fn func(d JobSpecDie)) JobDie {
	return d.DieStamp(func(r *batchv1.Job) {
		d := JobSpecBlank.DieImmutable(false).DieFeed(r.Spec)
		fn(d)
		r.Spec = d.DieRelease()
	})
}

func (d *jobDie) StatusDie(fn func(d JobStatusDie)) JobDie {
	return d.DieStamp(func(r *batchv1.Job) {
		d := JobStatusBlank.DieImmutable(false).DieFeed(r.Status)
		fn(d)
		r.Status = d.DieRelease()
	})
}

func (d *jobDie) Spec(v batchv1.JobSpec) JobDie {
	return d.DieStamp(func(r *batchv1.Job) {
		r.Spec = v
	})
}

func (d *jobDie) Status(v batchv1.JobStatus) JobDie {
	return d.DieStamp(func(r *batchv1.Job) {
		r.Status = v
	})
}

type JobSpecDie interface {
	// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
	DieStamp(fn func(r *batchv1.JobSpec)) JobSpecDie
	// DieFeed returns a new die with the provided resource.
	DieFeed(r batchv1.JobSpec) JobSpecDie
	// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
	DieFeedPtr(r *batchv1.JobSpec) JobSpecDie
	// DieRelease returns the resource managed by the die.
	DieRelease() batchv1.JobSpec
	// DieReleasePtr returns a pointer to the resource managed by the die.
	DieReleasePtr() *batchv1.JobSpec
	// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
	DieImmutable(immutable bool) JobSpecDie
	// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
	DeepCopy() JobSpecDie

	jobSpec
	// Specifies the maximum desired number of pods the job should run at any given time. The actual number of pods running in steady state will be less than this number when ((.spec.completions - .status.successful) < .spec.parallelism), i.e. when the work left to do is less than max parallelism. More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/
	Parallelism(Parallelism *int32) JobSpecDie
	// Specifies the desired number of successfully finished pods the job should be run with.  Setting to nil means that the success of any pod signals the success of all pods, and allows parallelism to have any positive value.  Setting to 1 means that parallelism is limited to 1 and the success of that pod signals the success of the job. More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/
	Completions(Completions *int32) JobSpecDie
	// Specifies the duration in seconds relative to the startTime that the job may be continuously active before the system tries to terminate it; value must be positive integer. If a Job is suspended (at creation or through an update), this timer will effectively be stopped and reset when the Job is resumed again.
	ActiveDeadlineSeconds(ActiveDeadlineSeconds *int64) JobSpecDie
	// Specifies the number of retries before marking this job failed. Defaults to 6
	BackoffLimit(BackoffLimit *int32) JobSpecDie
	// A label query over pods that should match the pod count. Normally, the system sets this field for you. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors
	Selector(Selector *apismetav1.LabelSelector) JobSpecDie
	// manualSelector controls generation of pod labels and pod selectors. Leave `manualSelector` unset unless you are certain what you are doing. When false or unset, the system pick labels unique to this job and appends those labels to the pod template.  When true, the user is responsible for picking unique labels and specifying the selector.  Failure to pick a unique label may cause this and other jobs to not function correctly.  However, You may see `manualSelector=true` in jobs that were created with the old `extensions/v1beta1` API. More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/#specifying-your-own-pod-selector
	ManualSelector(ManualSelector *bool) JobSpecDie
	// Describes the pod that will be created when executing a job. More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/
	Template(Template corev1.PodTemplateSpec) JobSpecDie
	// ttlSecondsAfterFinished limits the lifetime of a Job that has finished execution (either Complete or Failed). If this field is set, ttlSecondsAfterFinished after the Job finishes, it is eligible to be automatically deleted. When the Job is being deleted, its lifecycle guarantees (e.g. finalizers) will be honored. If this field is unset, the Job won't be automatically deleted. If this field is set to zero, the Job becomes eligible to be deleted immediately after it finishes. This field is alpha-level and is only honored by servers that enable the TTLAfterFinished feature.
	TTLSecondsAfterFinished(TTLSecondsAfterFinished *int32) JobSpecDie
	// CompletionMode specifies how Pod completions are tracked. It can be `NonIndexed` (default) or `Indexed`.
	//
	// `NonIndexed` means that the Job is considered complete when there have been .spec.completions successfully completed Pods. Each Pod completion is homologous to each other.
	//
	// `Indexed` means that the Pods of a Job get an associated completion index from 0 to (.spec.completions - 1), available in the annotation batch.kubernetes.io/job-completion-index. The Job is considered complete when there is one successfully completed Pod for each index. When value is `Indexed`, .spec.completions must be specified and `.spec.parallelism` must be less than or equal to 10^5. In addition, The Pod name takes the form `$(job-name)-$(index)-$(random-string)`, the Pod hostname takes the form `$(job-name)-$(index)`.
	//
	// This field is beta-level. More completion modes can be added in the future. If the Job controller observes a mode that it doesn't recognize, the controller skips updates for the Job.
	CompletionMode(CompletionMode *batchv1.CompletionMode) JobSpecDie
	// Suspend specifies whether the Job controller should create Pods or not. If a Job is created with suspend set to true, no Pods are created by the Job controller. If a Job is suspended after creation (i.e. the flag goes from false to true), the Job controller will delete all active Pods associated with this Job. Users must design their workload to gracefully handle this. Suspending a Job will reset the StartTime field of the Job, effectively resetting the ActiveDeadlineSeconds timer too. Defaults to false.
	//
	// This field is beta-level, gated by SuspendJob feature flag (enabled by default).
	Suspend(Suspend *bool) JobSpecDie
}

var _ JobSpecDie = (*jobSpecDie)(nil)
var JobSpecBlank = (&jobSpecDie{}).DieFeed(batchv1.JobSpec{})

type jobSpecDie struct {
	mutable bool
	r       batchv1.JobSpec
}

func (d *jobSpecDie) DieImmutable(immutable bool) JobSpecDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy().(*jobSpecDie)
	d.mutable = !immutable
	return d
}

func (d *jobSpecDie) DieFeed(r batchv1.JobSpec) JobSpecDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &jobSpecDie{
		mutable: d.mutable,
		r:       r,
	}
}

func (d *jobSpecDie) DieFeedPtr(r *batchv1.JobSpec) JobSpecDie {
	if r == nil {
		r = &batchv1.JobSpec{}
	}
	return d.DieFeed(*r)
}

func (d *jobSpecDie) DieRelease() batchv1.JobSpec {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

func (d *jobSpecDie) DieReleasePtr() *batchv1.JobSpec {
	r := d.DieRelease()
	return &r
}

func (d *jobSpecDie) DieStamp(fn func(r *batchv1.JobSpec)) JobSpecDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

func (d *jobSpecDie) DeepCopy() JobSpecDie {
	r := *d.r.DeepCopy()
	return &jobSpecDie{
		mutable: d.mutable,
		r:       r,
	}
}

func (d *jobSpecDie) Parallelism(v *int32) JobSpecDie {
	return d.DieStamp(func(r *batchv1.JobSpec) {
		r.Parallelism = v
	})
}

func (d *jobSpecDie) Completions(v *int32) JobSpecDie {
	return d.DieStamp(func(r *batchv1.JobSpec) {
		r.Completions = v
	})
}

func (d *jobSpecDie) ActiveDeadlineSeconds(v *int64) JobSpecDie {
	return d.DieStamp(func(r *batchv1.JobSpec) {
		r.ActiveDeadlineSeconds = v
	})
}

func (d *jobSpecDie) BackoffLimit(v *int32) JobSpecDie {
	return d.DieStamp(func(r *batchv1.JobSpec) {
		r.BackoffLimit = v
	})
}

func (d *jobSpecDie) Selector(v *apismetav1.LabelSelector) JobSpecDie {
	return d.DieStamp(func(r *batchv1.JobSpec) {
		r.Selector = v
	})
}

func (d *jobSpecDie) ManualSelector(v *bool) JobSpecDie {
	return d.DieStamp(func(r *batchv1.JobSpec) {
		r.ManualSelector = v
	})
}

func (d *jobSpecDie) Template(v corev1.PodTemplateSpec) JobSpecDie {
	return d.DieStamp(func(r *batchv1.JobSpec) {
		r.Template = v
	})
}

func (d *jobSpecDie) TTLSecondsAfterFinished(v *int32) JobSpecDie {
	return d.DieStamp(func(r *batchv1.JobSpec) {
		r.TTLSecondsAfterFinished = v
	})
}

func (d *jobSpecDie) CompletionMode(v *batchv1.CompletionMode) JobSpecDie {
	return d.DieStamp(func(r *batchv1.JobSpec) {
		r.CompletionMode = v
	})
}

func (d *jobSpecDie) Suspend(v *bool) JobSpecDie {
	return d.DieStamp(func(r *batchv1.JobSpec) {
		r.Suspend = v
	})
}

type JobStatusDie interface {
	// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
	DieStamp(fn func(r *batchv1.JobStatus)) JobStatusDie
	// DieFeed returns a new die with the provided resource.
	DieFeed(r batchv1.JobStatus) JobStatusDie
	// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
	DieFeedPtr(r *batchv1.JobStatus) JobStatusDie
	// DieRelease returns the resource managed by the die.
	DieRelease() batchv1.JobStatus
	// DieReleasePtr returns a pointer to the resource managed by the die.
	DieReleasePtr() *batchv1.JobStatus
	// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
	DieImmutable(immutable bool) JobStatusDie
	// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
	DeepCopy() JobStatusDie

	jobStatus
	// The latest available observations of an object's current state. When a Job fails, one of the conditions will have type "Failed" and status true. When a Job is suspended, one of the conditions will have type "Suspended" and status true; when the Job is resumed, the status of this condition will become false. When a Job is completed, one of the conditions will have type "Complete" and status true. More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/
	Conditions(Conditions ...batchv1.JobCondition) JobStatusDie
	// Represents time when the job controller started processing a job. When a Job is created in the suspended state, this field is not set until the first time it is resumed. This field is reset every time a Job is resumed from suspension. It is represented in RFC3339 form and is in UTC.
	StartTime(StartTime *apismetav1.Time) JobStatusDie
	// Represents time when the job was completed. It is not guaranteed to be set in happens-before order across separate operations. It is represented in RFC3339 form and is in UTC. The completion time is only set when the job finishes successfully.
	CompletionTime(CompletionTime *apismetav1.Time) JobStatusDie
	// The number of actively running pods.
	Active(Active int32) JobStatusDie
	// The number of pods which reached phase Succeeded.
	Succeeded(Succeeded int32) JobStatusDie
	// The number of pods which reached phase Failed.
	Failed(Failed int32) JobStatusDie
	// CompletedIndexes holds the completed indexes when .spec.completionMode = "Indexed" in a text format. The indexes are represented as decimal integers separated by commas. The numbers are listed in increasing order. Three or more consecutive numbers are compressed and represented by the first and last element of the series, separated by a hyphen. For example, if the completed indexes are 1, 3, 4, 5 and 7, they are represented as "1,3-5,7".
	CompletedIndexes(CompletedIndexes string) JobStatusDie
	// UncountedTerminatedPods holds the UIDs of Pods that have terminated but the job controller hasn't yet accounted for in the status counters.
	//
	// The job controller creates pods with a finalizer. When a pod terminates (succeeded or failed), the controller does three steps to account for it in the job status: (1) Add the pod UID to the arrays in this field. (2) Remove the pod finalizer. (3) Remove the pod UID from the arrays while increasing the corresponding     counter.
	//
	// This field is alpha-level. The job controller only makes use of this field when the feature gate PodTrackingWithFinalizers is enabled. Old jobs might not be tracked using this field, in which case the field remains null.
	UncountedTerminatedPods(UncountedTerminatedPods *batchv1.UncountedTerminatedPods) JobStatusDie
}

var _ JobStatusDie = (*jobStatusDie)(nil)
var JobStatusBlank = (&jobStatusDie{}).DieFeed(batchv1.JobStatus{})

type jobStatusDie struct {
	mutable bool
	r       batchv1.JobStatus
}

func (d *jobStatusDie) DieImmutable(immutable bool) JobStatusDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy().(*jobStatusDie)
	d.mutable = !immutable
	return d
}

func (d *jobStatusDie) DieFeed(r batchv1.JobStatus) JobStatusDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &jobStatusDie{
		mutable: d.mutable,
		r:       r,
	}
}

func (d *jobStatusDie) DieFeedPtr(r *batchv1.JobStatus) JobStatusDie {
	if r == nil {
		r = &batchv1.JobStatus{}
	}
	return d.DieFeed(*r)
}

func (d *jobStatusDie) DieRelease() batchv1.JobStatus {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

func (d *jobStatusDie) DieReleasePtr() *batchv1.JobStatus {
	r := d.DieRelease()
	return &r
}

func (d *jobStatusDie) DieStamp(fn func(r *batchv1.JobStatus)) JobStatusDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

func (d *jobStatusDie) DeepCopy() JobStatusDie {
	r := *d.r.DeepCopy()
	return &jobStatusDie{
		mutable: d.mutable,
		r:       r,
	}
}

func (d *jobStatusDie) Conditions(v ...batchv1.JobCondition) JobStatusDie {
	return d.DieStamp(func(r *batchv1.JobStatus) {
		r.Conditions = v
	})
}

func (d *jobStatusDie) StartTime(v *apismetav1.Time) JobStatusDie {
	return d.DieStamp(func(r *batchv1.JobStatus) {
		r.StartTime = v
	})
}

func (d *jobStatusDie) CompletionTime(v *apismetav1.Time) JobStatusDie {
	return d.DieStamp(func(r *batchv1.JobStatus) {
		r.CompletionTime = v
	})
}

func (d *jobStatusDie) Active(v int32) JobStatusDie {
	return d.DieStamp(func(r *batchv1.JobStatus) {
		r.Active = v
	})
}

func (d *jobStatusDie) Succeeded(v int32) JobStatusDie {
	return d.DieStamp(func(r *batchv1.JobStatus) {
		r.Succeeded = v
	})
}

func (d *jobStatusDie) Failed(v int32) JobStatusDie {
	return d.DieStamp(func(r *batchv1.JobStatus) {
		r.Failed = v
	})
}

func (d *jobStatusDie) CompletedIndexes(v string) JobStatusDie {
	return d.DieStamp(func(r *batchv1.JobStatus) {
		r.CompletedIndexes = v
	})
}

func (d *jobStatusDie) UncountedTerminatedPods(v *batchv1.UncountedTerminatedPods) JobStatusDie {
	return d.DieStamp(func(r *batchv1.JobStatus) {
		r.UncountedTerminatedPods = v
	})
}
