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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// +die:object=true
type StatefulSpec = appsv1.StatefulSet

// +die
type StatefulSetSpec = appsv1.StatefulSetSpec

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

// +die
type StatefulSetStatus = appsv1.StatefulSetStatus

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
