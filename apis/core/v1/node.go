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
	diemetav1 "github.com/scothis/dies/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
)

// +die:object=true
type Node = corev1.Node

// +die
type NodeSpec = corev1.NodeSpec

// +die
type NodeStatus = corev1.NodeStatus

func (d *NodeStatusDie) ConditionsDie(conditions ...*diemetav1.ConditionDie) *NodeStatusDie {
	return d.DieStamp(func(r *corev1.NodeStatus) {
		r.Conditions = make([]corev1.NodeCondition, len(conditions))
		for i := range conditions {
			c := conditions[i].DieRelease()
			r.Conditions[i] = corev1.NodeCondition{
				Type:               corev1.NodeConditionType(c.Type),
				Status:             corev1.ConditionStatus(c.Status),
				Reason:             c.Reason,
				Message:            c.Message,
				LastTransitionTime: c.LastTransitionTime,
			}
		}
	})
}
