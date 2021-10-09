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

// +die:target=k8s.io/api/apps/v1.DaemonSet,object=true

// +die:target=k8s.io/api/apps/v1.DaemonSetSpec
// +die:field:receiver=DaemonSetSpecDie,name=Selector,type=*k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector
// +die:field:receiver=DaemonSetSpecDie,name=Template,type=k8s.io/api/core/v1.PodTemplateSpec
// +die:field:receiver=DaemonSetSpecDie,name=UpdateStrategy,type=k8s.io/api/apps/v1.DaemonSetUpdateStrategy
// +die:field:receiver=DaemonSetSpecDie,name=MinReadySeconds,type=int32
// +die:field:receiver=DaemonSetSpecDie,name=RevisionHistoryLimit,type=*int32

func (d *DaemonSetSpecDie) TemplateDie(fn func(d *diecorev1.PodTemplateSpecDie)) *DaemonSetSpecDie {
	return d.DieStamp(func(r *appsv1.DaemonSetSpec) {
		d := diecorev1.PodTemplateSpecBlank.DieImmutable(false).DieFeed(r.Template)
		fn(d)
		r.Template = d.DieRelease()
	})
}

// +die:target=k8s.io/api/apps/v1.DaemonSetStatus
// +die:field:receiver=DaemonSetStatusDie,name=CurrentNumberScheduled,type=int32
// +die:field:receiver=DaemonSetStatusDie,name=NumberMisscheduled,type=int32
// +die:field:receiver=DaemonSetStatusDie,name=DesiredNumberScheduled,type=int32
// +die:field:receiver=DaemonSetStatusDie,name=NumberReady,type=int32
// +die:field:receiver=DaemonSetStatusDie,name=ObservedGeneration,type=int64
// +die:field:receiver=DaemonSetStatusDie,name=UpdatedNumberScheduled,type=int32
// +die:field:receiver=DaemonSetStatusDie,name=NumberAvailable,type=int32
// +die:field:receiver=DaemonSetStatusDie,name=NumberUnavailable,type=int32
// +die:field:receiver=DaemonSetStatusDie,name=CollisionCount,type=*int32
// +die:field:receiver=DaemonSetStatusDie,name=Conditions,type=[]k8s.io/api/apps/v1.DaemonSetCondition

func (d *DaemonSetStatusDie) ConditionsDie(conditions ...*diemetav1.ConditionDie) *DaemonSetStatusDie {
	return d.DieStamp(func(r *appsv1.DaemonSetStatus) {
		r.Conditions = make([]appsv1.DaemonSetCondition, len(conditions))
		for i := range conditions {
			c := conditions[i].DieRelease()
			r.Conditions[i] = appsv1.DaemonSetCondition{
				Type:               appsv1.DaemonSetConditionType(c.Type),
				Status:             corev1.ConditionStatus(c.Status),
				Reason:             c.Reason,
				Message:            c.Message,
				LastTransitionTime: c.LastTransitionTime,
			}
		}
	})
}
