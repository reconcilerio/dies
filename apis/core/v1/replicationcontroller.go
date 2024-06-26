/*
Copyright 2022 the original author or authors.

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
	corev1 "k8s.io/api/core/v1"
	diemetav1 "reconciler.io/dies/apis/meta/v1"
)

// +die:object=true,apiVersion=v1,kind=ReplicationController
type _ = corev1.ReplicationController

// +die
// +die:field:name=Template,die=PodTemplateSpecDie,pointer=true
type _ = corev1.ReplicationControllerSpec

// +die
type _ = corev1.ReplicationControllerStatus

func (d *ReplicationControllerStatusDie) ConditionsDie(conditions ...*diemetav1.ConditionDie) *ReplicationControllerStatusDie {
	return d.DieStamp(func(r *corev1.ReplicationControllerStatus) {
		r.Conditions = make([]corev1.ReplicationControllerCondition, len(conditions))
		for i := range conditions {
			c := conditions[i].DieRelease()
			r.Conditions[i] = corev1.ReplicationControllerCondition{
				Type:               corev1.ReplicationControllerConditionType(c.Type),
				Status:             corev1.ConditionStatus(c.Status),
				Reason:             c.Reason,
				Message:            c.Message,
				LastTransitionTime: c.LastTransitionTime,
			}
		}
	})
}
