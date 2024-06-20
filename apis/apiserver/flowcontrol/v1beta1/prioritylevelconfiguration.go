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

package v1beta1

import (
	flowcontrolv1beta1 "k8s.io/api/flowcontrol/v1beta1"
	diemetav1 "reconciler.io/dies/apis/meta/v1"
)

// +die:object=true,apiVersion=flowcontrol.apiserver.k8s.io/v1beta1,kind=PriorityLevelConfiguration
type _ = flowcontrolv1beta1.PriorityLevelConfiguration

// +die
// +die:field:name=Limited,die=LimitedPriorityLevelConfigurationDie,pointer=true
// +die:field:name=Exempt,die=ExemptPriorityLevelConfigurationDie,pointer=true
type _ = flowcontrolv1beta1.PriorityLevelConfigurationSpec

// +die
// +die:field:name=LimitResponse,die=LimitResponseDie
type _ = flowcontrolv1beta1.LimitedPriorityLevelConfiguration

// +die
// +die:field:name=Queuing,die=QueuingConfigurationDie,pointer=true
type _ = flowcontrolv1beta1.LimitResponse

// +die
type _ = flowcontrolv1beta1.ExemptPriorityLevelConfiguration

// +die
type _ = flowcontrolv1beta1.QueuingConfiguration

// +die
type _ = flowcontrolv1beta1.PriorityLevelConfigurationStatus

func (d *PriorityLevelConfigurationStatusDie) ConditionsDie(conditions ...*diemetav1.ConditionDie) *PriorityLevelConfigurationStatusDie {
	return d.DieStamp(func(r *flowcontrolv1beta1.PriorityLevelConfigurationStatus) {
		r.Conditions = make([]flowcontrolv1beta1.PriorityLevelConfigurationCondition, len(conditions))
		for i := range conditions {
			c := conditions[i].DieRelease()
			r.Conditions[i] = flowcontrolv1beta1.PriorityLevelConfigurationCondition{
				Type:               flowcontrolv1beta1.PriorityLevelConfigurationConditionType(c.Type),
				Status:             flowcontrolv1beta1.ConditionStatus(c.Status),
				Reason:             c.Reason,
				Message:            c.Message,
				LastTransitionTime: c.LastTransitionTime,
			}
		}
	})
}
