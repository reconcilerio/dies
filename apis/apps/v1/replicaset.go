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
	diemetav1 "reconciler.io/dies/apis/meta/v1"
)

// +die:object=true,apiVersion=apps/v1,kind=ReplicaSet
type _ = appsv1.ReplicaSet

// +die
// +die:field:name=Selector,package=_/meta/v1,die=LabelSelectorDie,pointer=true
// +die:field:name=Template,package=_/core/v1,die=PodTemplateSpecDie
type _ = appsv1.ReplicaSetSpec

// +die
type _ = appsv1.ReplicaSetStatus

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
