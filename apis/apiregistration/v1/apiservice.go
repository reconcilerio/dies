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
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration"
	diemetav1 "reconciler.io/dies/apis/meta/v1"
)

// +die:object=true,apiVersion=apiregistration.k8s.io/v1,kind=APIService
type _ = apiregistrationv1.APIService

// +die
// +die:field:name=Service,die=ServiceReferenceDie,pointer=true
type _ = apiregistrationv1.APIServiceSpec

// +die
type _ = apiregistrationv1.ServiceReference

// +die
type _ = apiregistrationv1.APIServiceStatus

func (d *APIServiceStatusDie) ConditionsDie(conditions ...*diemetav1.ConditionDie) *APIServiceStatusDie {
	return d.DieStamp(func(r *apiregistrationv1.APIServiceStatus) {
		r.Conditions = make([]apiregistrationv1.APIServiceCondition, len(conditions))
		for i := range conditions {
			c := conditions[i].DieRelease()
			r.Conditions[i] = apiregistrationv1.APIServiceCondition{
				Type:               apiregistrationv1.APIServiceConditionType(c.Type),
				Status:             apiregistrationv1.ConditionStatus(c.Status),
				Reason:             c.Reason,
				Message:            c.Message,
				LastTransitionTime: c.LastTransitionTime,
			}
		}
	})
}
