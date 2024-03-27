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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	diecorev1 "reconciler.io/dies/apis/core/v1"
	diemetav1 "reconciler.io/dies/apis/meta/v1"
)

// +die:object=true
type _ = appsv1.StatefulSet

// +die
type _ = appsv1.StatefulSetSpec

func (d *StatefulSetSpecDie) SelectorDie(fn func(d *diemetav1.LabelSelectorDie)) *StatefulSetSpecDie {
	return d.DieStamp(func(r *appsv1.StatefulSetSpec) {
		d := diemetav1.LabelSelectorBlank.DieImmutable(false).DieFeedPtr(r.Selector)
		fn(d)
		r.Selector = d.DieReleasePtr()
	})
}

func (d *StatefulSetSpecDie) TemplateDie(fn func(d *diecorev1.PodTemplateSpecDie)) *StatefulSetSpecDie {
	return d.DieStamp(func(r *appsv1.StatefulSetSpec) {
		d := diecorev1.PodTemplateSpecBlank.DieImmutable(false).DieFeed(r.Template)
		fn(d)
		r.Template = d.DieRelease()
	})
}

func (d *StatefulSetSpecDie) VolumeClaimTemplatesDie(volumeClaimTemplates ...*diecorev1.PersistentVolumeClaimDie) *StatefulSetSpecDie {
	return d.DieStamp(func(r *appsv1.StatefulSetSpec) {
		r.VolumeClaimTemplates = make([]corev1.PersistentVolumeClaim, len(volumeClaimTemplates))
		for i, v := range volumeClaimTemplates {
			r.VolumeClaimTemplates[i] = v.DieRelease()
		}
	})
}

func (d *StatefulSetSpecDie) UpdateStrategyDie(fn func(d *StatefulSetUpdateStrategyDie)) *StatefulSetSpecDie {
	return d.DieStamp(func(r *appsv1.StatefulSetSpec) {
		d := StatefulSetUpdateStrategyBlank.DieImmutable(false).DieFeed(r.UpdateStrategy)
		fn(d)
		r.UpdateStrategy = d.DieRelease()
	})
}

func (d *StatefulSetSpecDie) PersistentVolumeClaimRetentionPolicyDie(fn func(d *StatefulSetPersistentVolumeClaimRetentionPolicyDie)) *StatefulSetSpecDie {
	return d.DieStamp(func(r *appsv1.StatefulSetSpec) {
		d := StatefulSetPersistentVolumeClaimRetentionPolicyBlank.DieImmutable(false).DieFeedPtr(r.PersistentVolumeClaimRetentionPolicy)
		fn(d)
		r.PersistentVolumeClaimRetentionPolicy = d.DieReleasePtr()
	})
}

func (d *StatefulSetSpecDie) OrdinalsDie(fn func(d *StatefulSetOrdinalsDie)) *StatefulSetSpecDie {
	return d.DieStamp(func(r *appsv1.StatefulSetSpec) {
		d := StatefulSetOrdinalsBlank.DieImmutable(false).DieFeedPtr(r.Ordinals)
		fn(d)
		r.Ordinals = d.DieReleasePtr()
	})
}

// +die
type _ = appsv1.StatefulSetUpdateStrategy

func (d *StatefulSetUpdateStrategyDie) OnDelete() *StatefulSetUpdateStrategyDie {
	return d.DieStamp(func(r *appsv1.StatefulSetUpdateStrategy) {
		r.Type = appsv1.OnDeleteStatefulSetStrategyType
		r.RollingUpdate = nil
	})
}

func (d *StatefulSetUpdateStrategyDie) RollingUpdateDie(fn func(d *RollingUpdateStatefulSetStrategyDie)) *StatefulSetUpdateStrategyDie {
	return d.DieStamp(func(r *appsv1.StatefulSetUpdateStrategy) {
		r.Type = appsv1.RollingUpdateStatefulSetStrategyType
		d := RollingUpdateStatefulSetStrategyBlank.DieImmutable(false).DieFeedPtr(r.RollingUpdate)
		fn(d)
		r.RollingUpdate = d.DieReleasePtr()
	})
}

// +die
type _ = appsv1.RollingUpdateStatefulSetStrategy

// +die
type _ = appsv1.StatefulSetPersistentVolumeClaimRetentionPolicy

// +die
type _ = appsv1.StatefulSetOrdinals

// +die
type _ = appsv1.StatefulSetStatus

func (d *StatefulSetStatusDie) ConditionsDie(conditions ...*diemetav1.ConditionDie) *StatefulSetStatusDie {
	return d.DieStamp(func(r *appsv1.StatefulSetStatus) {
		r.Conditions = make([]appsv1.StatefulSetCondition, len(conditions))
		for i := range conditions {
			c := conditions[i].DieRelease()
			r.Conditions[i] = appsv1.StatefulSetCondition{
				Type:               appsv1.StatefulSetConditionType(c.Type),
				Status:             corev1.ConditionStatus(c.Status),
				Reason:             c.Reason,
				Message:            c.Message,
				LastTransitionTime: c.LastTransitionTime,
			}
		}
	})
}
