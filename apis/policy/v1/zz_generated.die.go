//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2021-2022 the original author or authors.

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
	policyv1 "k8s.io/api/policy/v1"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
)

var PodDisruptionBudgetBlank = (&PodDisruptionBudgetDie{}).DieFeed(policyv1.PodDisruptionBudget{})

type PodDisruptionBudgetDie struct {
	metav1.FrozenObjectMeta
	mutable bool
	r       policyv1.PodDisruptionBudget
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *PodDisruptionBudgetDie) DieImmutable(immutable bool) *PodDisruptionBudgetDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *PodDisruptionBudgetDie) DieFeed(r policyv1.PodDisruptionBudget) *PodDisruptionBudgetDie {
	if d.mutable {
		d.FrozenObjectMeta = metav1.FreezeObjectMeta(r.ObjectMeta)
		d.r = r
		return d
	}
	return &PodDisruptionBudgetDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *PodDisruptionBudgetDie) DieFeedPtr(r *policyv1.PodDisruptionBudget) *PodDisruptionBudgetDie {
	if r == nil {
		r = &policyv1.PodDisruptionBudget{}
	}
	return d.DieFeed(*r)
}

// DieRelease returns the resource managed by the die.
func (d *PodDisruptionBudgetDie) DieRelease() policyv1.PodDisruptionBudget {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *PodDisruptionBudgetDie) DieReleasePtr() *policyv1.PodDisruptionBudget {
	r := d.DieRelease()
	return &r
}

// DieReleaseUnstructured returns the resource managed by the die as an unstructured object.
func (d *PodDisruptionBudgetDie) DieReleaseUnstructured() runtime.Unstructured {
	r := d.DieReleasePtr()
	u, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(r)
	return &unstructured.Unstructured{
		Object: u,
	}
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *PodDisruptionBudgetDie) DieStamp(fn func(r *policyv1.PodDisruptionBudget)) *PodDisruptionBudgetDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *PodDisruptionBudgetDie) DeepCopy() *PodDisruptionBudgetDie {
	r := *d.r.DeepCopy()
	return &PodDisruptionBudgetDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

var _ runtime.Object = (*PodDisruptionBudgetDie)(nil)

func (d *PodDisruptionBudgetDie) DeepCopyObject() runtime.Object {
	return d.r.DeepCopy()
}

func (d *PodDisruptionBudgetDie) GetObjectKind() schema.ObjectKind {
	r := d.DieRelease()
	return r.GetObjectKind()
}

func (d *PodDisruptionBudgetDie) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.r)
}

func (d *PodDisruptionBudgetDie) UnmarshalJSON(b []byte) error {
	if d == PodDisruptionBudgetBlank {
		return fmtx.Errorf("cannot unmarshal into the blank die, create a copy first")
	}
	if !d.mutable {
		return fmtx.Errorf("cannot unmarshal into immutable dies, create a mutable version first")
	}
	r := &policyv1.PodDisruptionBudget{}
	err := json.Unmarshal(b, r)
	*d = *d.DieFeed(*r)
	return err
}

// MetadataDie stamps the resource's ObjectMeta field with a mutable die.
func (d *PodDisruptionBudgetDie) MetadataDie(fn func(d *metav1.ObjectMetaDie)) *PodDisruptionBudgetDie {
	return d.DieStamp(func(r *policyv1.PodDisruptionBudget) {
		d := metav1.ObjectMetaBlank.DieImmutable(false).DieFeed(r.ObjectMeta)
		fn(d)
		r.ObjectMeta = d.DieRelease()
	})
}

// SpecDie stamps the resource's spec field with a mutable die.
func (d *PodDisruptionBudgetDie) SpecDie(fn func(d *PodDisruptionBudgetSpecDie)) *PodDisruptionBudgetDie {
	return d.DieStamp(func(r *policyv1.PodDisruptionBudget) {
		d := PodDisruptionBudgetSpecBlank.DieImmutable(false).DieFeed(r.Spec)
		fn(d)
		r.Spec = d.DieRelease()
	})
}

// StatusDie stamps the resource's status field with a mutable die.
func (d *PodDisruptionBudgetDie) StatusDie(fn func(d *PodDisruptionBudgetStatusDie)) *PodDisruptionBudgetDie {
	return d.DieStamp(func(r *policyv1.PodDisruptionBudget) {
		d := PodDisruptionBudgetStatusBlank.DieImmutable(false).DieFeed(r.Status)
		fn(d)
		r.Status = d.DieRelease()
	})
}

// Specification of the desired behavior of the PodDisruptionBudget.
func (d *PodDisruptionBudgetDie) Spec(v policyv1.PodDisruptionBudgetSpec) *PodDisruptionBudgetDie {
	return d.DieStamp(func(r *policyv1.PodDisruptionBudget) {
		r.Spec = v
	})
}

// Most recently observed status of the PodDisruptionBudget.
func (d *PodDisruptionBudgetDie) Status(v policyv1.PodDisruptionBudgetStatus) *PodDisruptionBudgetDie {
	return d.DieStamp(func(r *policyv1.PodDisruptionBudget) {
		r.Status = v
	})
}

var PodDisruptionBudgetSpecBlank = (&PodDisruptionBudgetSpecDie{}).DieFeed(policyv1.PodDisruptionBudgetSpec{})

type PodDisruptionBudgetSpecDie struct {
	mutable bool
	r       policyv1.PodDisruptionBudgetSpec
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *PodDisruptionBudgetSpecDie) DieImmutable(immutable bool) *PodDisruptionBudgetSpecDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *PodDisruptionBudgetSpecDie) DieFeed(r policyv1.PodDisruptionBudgetSpec) *PodDisruptionBudgetSpecDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &PodDisruptionBudgetSpecDie{
		mutable: d.mutable,
		r:       r,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *PodDisruptionBudgetSpecDie) DieFeedPtr(r *policyv1.PodDisruptionBudgetSpec) *PodDisruptionBudgetSpecDie {
	if r == nil {
		r = &policyv1.PodDisruptionBudgetSpec{}
	}
	return d.DieFeed(*r)
}

// DieRelease returns the resource managed by the die.
func (d *PodDisruptionBudgetSpecDie) DieRelease() policyv1.PodDisruptionBudgetSpec {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *PodDisruptionBudgetSpecDie) DieReleasePtr() *policyv1.PodDisruptionBudgetSpec {
	r := d.DieRelease()
	return &r
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *PodDisruptionBudgetSpecDie) DieStamp(fn func(r *policyv1.PodDisruptionBudgetSpec)) *PodDisruptionBudgetSpecDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *PodDisruptionBudgetSpecDie) DeepCopy() *PodDisruptionBudgetSpecDie {
	r := *d.r.DeepCopy()
	return &PodDisruptionBudgetSpecDie{
		mutable: d.mutable,
		r:       r,
	}
}

// An eviction is allowed if at least "minAvailable" pods selected by "selector" will still be available after the eviction, i.e. even in the absence of the evicted pod.  So for example you can prevent all voluntary evictions by specifying "100%".
func (d *PodDisruptionBudgetSpecDie) MinAvailable(v *intstr.IntOrString) *PodDisruptionBudgetSpecDie {
	return d.DieStamp(func(r *policyv1.PodDisruptionBudgetSpec) {
		r.MinAvailable = v
	})
}

func (d *PodDisruptionBudgetSpecDie) MinAvailableInt(i int) *PodDisruptionBudgetSpecDie {
	return d.DieStamp(func(r *policyv1.PodDisruptionBudgetSpec) {
		v := intstr.FromInt(i)
		r.MinAvailable = &v
	})
}

func (d *PodDisruptionBudgetSpecDie) MinAvailableString(s string) *PodDisruptionBudgetSpecDie {
	return d.DieStamp(func(r *policyv1.PodDisruptionBudgetSpec) {
		v := intstr.FromString(s)
		r.MinAvailable = &v
	})
}

// Label query over pods whose evictions are managed by the disruption budget. A null selector will match no pods, while an empty ({}) selector will select all pods within the namespace.
func (d *PodDisruptionBudgetSpecDie) Selector(v *apismetav1.LabelSelector) *PodDisruptionBudgetSpecDie {
	return d.DieStamp(func(r *policyv1.PodDisruptionBudgetSpec) {
		r.Selector = v
	})
}

// An eviction is allowed if at most "maxUnavailable" pods selected by "selector" are unavailable after the eviction, i.e. even in absence of the evicted pod. For example, one can prevent all voluntary evictions by specifying 0. This is a mutually exclusive setting with "minAvailable".
func (d *PodDisruptionBudgetSpecDie) MaxUnavailable(v *intstr.IntOrString) *PodDisruptionBudgetSpecDie {
	return d.DieStamp(func(r *policyv1.PodDisruptionBudgetSpec) {
		r.MaxUnavailable = v
	})
}

func (d *PodDisruptionBudgetSpecDie) MaxUnavailableInt(i int) *PodDisruptionBudgetSpecDie {
	return d.DieStamp(func(r *policyv1.PodDisruptionBudgetSpec) {
		v := intstr.FromInt(i)
		r.MaxUnavailable = &v
	})
}

func (d *PodDisruptionBudgetSpecDie) MaxUnavailableString(s string) *PodDisruptionBudgetSpecDie {
	return d.DieStamp(func(r *policyv1.PodDisruptionBudgetSpec) {
		v := intstr.FromString(s)
		r.MaxUnavailable = &v
	})
}

var PodDisruptionBudgetStatusBlank = (&PodDisruptionBudgetStatusDie{}).DieFeed(policyv1.PodDisruptionBudgetStatus{})

type PodDisruptionBudgetStatusDie struct {
	mutable bool
	r       policyv1.PodDisruptionBudgetStatus
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *PodDisruptionBudgetStatusDie) DieImmutable(immutable bool) *PodDisruptionBudgetStatusDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *PodDisruptionBudgetStatusDie) DieFeed(r policyv1.PodDisruptionBudgetStatus) *PodDisruptionBudgetStatusDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &PodDisruptionBudgetStatusDie{
		mutable: d.mutable,
		r:       r,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *PodDisruptionBudgetStatusDie) DieFeedPtr(r *policyv1.PodDisruptionBudgetStatus) *PodDisruptionBudgetStatusDie {
	if r == nil {
		r = &policyv1.PodDisruptionBudgetStatus{}
	}
	return d.DieFeed(*r)
}

// DieRelease returns the resource managed by the die.
func (d *PodDisruptionBudgetStatusDie) DieRelease() policyv1.PodDisruptionBudgetStatus {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *PodDisruptionBudgetStatusDie) DieReleasePtr() *policyv1.PodDisruptionBudgetStatus {
	r := d.DieRelease()
	return &r
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *PodDisruptionBudgetStatusDie) DieStamp(fn func(r *policyv1.PodDisruptionBudgetStatus)) *PodDisruptionBudgetStatusDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *PodDisruptionBudgetStatusDie) DeepCopy() *PodDisruptionBudgetStatusDie {
	r := *d.r.DeepCopy()
	return &PodDisruptionBudgetStatusDie{
		mutable: d.mutable,
		r:       r,
	}
}

// Most recent generation observed when updating this PDB status. DisruptionsAllowed and other status information is valid only if observedGeneration equals to PDB's object generation.
func (d *PodDisruptionBudgetStatusDie) ObservedGeneration(v int64) *PodDisruptionBudgetStatusDie {
	return d.DieStamp(func(r *policyv1.PodDisruptionBudgetStatus) {
		r.ObservedGeneration = v
	})
}

// Number of pod disruptions that are currently allowed.
func (d *PodDisruptionBudgetStatusDie) DisruptionsAllowed(v int32) *PodDisruptionBudgetStatusDie {
	return d.DieStamp(func(r *policyv1.PodDisruptionBudgetStatus) {
		r.DisruptionsAllowed = v
	})
}

// current number of healthy pods
func (d *PodDisruptionBudgetStatusDie) CurrentHealthy(v int32) *PodDisruptionBudgetStatusDie {
	return d.DieStamp(func(r *policyv1.PodDisruptionBudgetStatus) {
		r.CurrentHealthy = v
	})
}

// minimum desired number of healthy pods
func (d *PodDisruptionBudgetStatusDie) DesiredHealthy(v int32) *PodDisruptionBudgetStatusDie {
	return d.DieStamp(func(r *policyv1.PodDisruptionBudgetStatus) {
		r.DesiredHealthy = v
	})
}

// total number of pods counted by this disruption budget
func (d *PodDisruptionBudgetStatusDie) ExpectedPods(v int32) *PodDisruptionBudgetStatusDie {
	return d.DieStamp(func(r *policyv1.PodDisruptionBudgetStatus) {
		r.ExpectedPods = v
	})
}

// Conditions contain conditions for PDB. The disruption controller sets the DisruptionAllowed condition. The following are known values for the reason field (additional reasons could be added in the future): - SyncFailed: The controller encountered an error and wasn't able to compute the number of allowed disruptions. Therefore no disruptions are allowed and the status of the condition will be False. - InsufficientPods: The number of pods are either at or below the number required by the PodDisruptionBudget. No disruptions are allowed and the status of the condition will be False. - SufficientPods: There are more pods than required by the PodDisruptionBudget. The condition will be True, and the number of allowed disruptions are provided by the disruptionsAllowed property.
func (d *PodDisruptionBudgetStatusDie) Conditions(v ...apismetav1.Condition) *PodDisruptionBudgetStatusDie {
	return d.DieStamp(func(r *policyv1.PodDisruptionBudgetStatus) {
		r.Conditions = v
	})
}
