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

// +die:target=k8s.io/api/apps/v1.ReplicaSet,object=true

// +die:target=k8s.io/api/apps/v1.ReplicaSetSpec
// +die:field:receiver=ReplicaSetSpecDie,name=Replicas,type=*int32
// +die:field:receiver=ReplicaSetSpecDie,name=MinReadySeconds,type=int32
// +die:field:receiver=ReplicaSetSpecDie,name=Selector,type=*k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector
// +die:field:receiver=ReplicaSetSpecDie,name=Template,type=k8s.io/api/core/v1.PodTemplateSpec

func (d *ReplicaSetSpecDie) TemplateDie(fn func(d *diecorev1.PodTemplateSpecDie)) *ReplicaSetSpecDie {
	return d.DieStamp(func(r *appsv1.ReplicaSetSpec) {
		d := diecorev1.PodTemplateSpecBlank.DieImmutable(false).DieFeed(r.Template)
		fn(d)
		r.Template = d.DieRelease()
	})
}

// +die:target=k8s.io/api/apps/v1.ReplicaSetStatus
// +die:field:receiver=ReplicaSetStatusDie,name=Replicas,type=int32
// +die:field:receiver=ReplicaSetStatusDie,name=FullyLabeledReplicas,type=int32
// +die:field:receiver=ReplicaSetStatusDie,name=ReadyReplicas,type=int32
// +die:field:receiver=ReplicaSetStatusDie,name=AvailableReplicas,type=int32
// +die:field:receiver=ReplicaSetStatusDie,name=ObservedGeneration,type=int64
// +die:field:receiver=ReplicaSetStatusDie,name=Conditions,type=[]k8s.io/api/apps/v1.ReplicaSetCondition

func (d *ReplicaSetStatusDie) ConditionsDie(conditions ...*diemetav1.ConditionDie) *ReplicaSetStatusDie {
	return d.DieStamp(func(r *appsv1.ReplicaSetStatus) {
		r.Conditions = make([]appsv1.ReplicaSetCondition, len(conditions))
		for i := range conditions {
			c := conditions[i].DieRelease()
			r.Conditions[i] = appsv1.ReplicaSetCondition{
				Type:               appsv1.ReplicaSetConditionType(c.Type),
				Status:             corev1.ConditionStatus(c.Status),
				Reason:             c.Reason,
				Message:            c.Message,
				LastTransitionTime: c.LastTransitionTime,
			}
		}
	})
}
