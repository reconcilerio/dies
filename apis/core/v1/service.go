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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +die:object=true
type Service = corev1.Service

// +die
type ServiceSpec = corev1.ServiceSpec

// +die
type ServiceStatus = corev1.ServiceStatus

func (d *ServiceStatusDie) ConditionsDie(conditions ...diemetav1.ConditionDie) *ServiceStatusDie {
	return d.DieStamp(func(r *corev1.ServiceStatus) {
		r.Conditions = make([]metav1.Condition, len(conditions))
		for i, c := range conditions {
			r.Conditions[i] = c.DieRelease()
		}
	})
}
