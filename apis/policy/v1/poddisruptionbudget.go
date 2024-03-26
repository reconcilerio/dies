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
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	diemetav1 "reconciler.io/dies/apis/meta/v1"
)

// +die:object=true
type _ = policyv1.PodDisruptionBudget

// +die
type _ = policyv1.PodDisruptionBudgetSpec

// TODO(scothis) fix import for maps with struct values, ignore the 'DisruptedPods' field until then

// +die:ignore={DisruptedPods}
type _ = policyv1.PodDisruptionBudgetStatus

// DisruptedPods contains information about pods whose eviction was processed by the API server eviction subresource handler but has not yet been observed by the PodDisruptionBudget controller. A pod will be in this map from the time when the API server processed the eviction request to the time when the pod is seen by PDB controller as having been marked for deletion (or after a timeout). The key in the map is the name of the pod and the value is the time when the API server processed the eviction request. If the deletion didn't occur and a pod is still there it will be removed from the list automatically by PodDisruptionBudget controller after some time. If everything goes smooth this map should be empty for the most of the time. Large number of entries in the map may indicate problems with pod deletions.
func (d *PodDisruptionBudgetStatusDie) DisruptedPods(v map[string]metav1.Time) *PodDisruptionBudgetStatusDie {
	return d.DieStamp(func(r *policyv1.PodDisruptionBudgetStatus) {
		r.DisruptedPods = v
	})
}

func (d *PodDisruptionBudgetStatusDie) DisruptedPodDie(key string, value metav1.Time) *PodDisruptionBudgetStatusDie {
	return d.DieStamp(func(r *policyv1.PodDisruptionBudgetStatus) {
		if r.DisruptedPods == nil {
			r.DisruptedPods = map[string]metav1.Time{}
		}
		r.DisruptedPods[key] = value
	})
}

func (d *PodDisruptionBudgetStatusDie) ConditionsDie(conditions ...*diemetav1.ConditionDie) *PodDisruptionBudgetStatusDie {
	return d.DieStamp(func(r *policyv1.PodDisruptionBudgetStatus) {
		r.Conditions = make([]metav1.Condition, len(conditions))
		for i := range conditions {
			r.Conditions[i] = conditions[i].DieRelease()
		}
	})
}
