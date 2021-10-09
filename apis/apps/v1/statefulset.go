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

// +die:target=k8s.io/api/apps/v1.StatefulSet,object=true

// +die:target=k8s.io/api/apps/v1.StatefulSetSpec
// +die:field:receiver=StatefulSetSpecDie,name=Replicas,type=*int32
// +die:field:receiver=StatefulSetSpecDie,name=Selector,type=*k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector
// +die:field:receiver=StatefulSetSpecDie,name=Template,type=k8s.io/api/core/v1.PodTemplateSpec
// +die:field:receiver=StatefulSetSpecDie,name=VolumeClaimTemplates,type=[]k8s.io/api/core/v1.PersistentVolumeClaim
// +die:field:receiver=StatefulSetSpecDie,name=ServiceName,type=string
// +die:field:receiver=StatefulSetSpecDie,name=PodManagementPolicy,type=k8s.io/api/apps/v1.PodManagementPolicyType
// +die:field:receiver=StatefulSetSpecDie,name=UpdateStrategy,type=k8s.io/api/apps/v1.StatefulSetUpdateStrategy
// +die:field:receiver=StatefulSetSpecDie,name=RevisionHistoryLimit,type=*int32
// +die:field:receiver=StatefulSetSpecDie,name=MinReadySeconds,type=int32

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

// +die:target=k8s.io/api/apps/v1.StatefulSetStatus
// +die:field:receiver=StatefulSetStatusDie,name=ObservedGeneration,type=int64
// +die:field:receiver=StatefulSetStatusDie,name=Replicas,type=int32
// +die:field:receiver=StatefulSetStatusDie,name=ReadyReplicas,type=int32
// +die:field:receiver=StatefulSetStatusDie,name=CurrentReplicas,type=int32
// +die:field:receiver=StatefulSetStatusDie,name=UpdatedReplicas,type=int32
// +die:field:receiver=StatefulSetStatusDie,name=CurrentRevision,type=string
// +die:field:receiver=StatefulSetStatusDie,name=UpdateRevision,type=string
// +die:field:receiver=StatefulSetStatusDie,name=CollisionCount,type=*int32
// +die:field:receiver=StatefulSetStatusDie,name=Conditions,type=[]k8s.io/api/apps/v1.StatefulSetCondition
// +die:field:receiver=StatefulSetStatusDie,name=AvailableReplicas,type=int32

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
