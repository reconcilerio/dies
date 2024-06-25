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

// +die:object=true,apiVersion=apps/v1,kind=Deployment
type _ = appsv1.Deployment

// +die
// +die:field:name=Selector,package=_/meta/v1,die=LabelSelectorDie,pointer=true
// +die:field:name=Template,package=_/core/v1,die=PodTemplateSpecDie
// +die:field:name=Strategy,die=DeploymentStrategyDie
type _ = appsv1.DeploymentSpec

// +die
type _ = appsv1.DeploymentStrategy

func (d *DeploymentStrategyDie) Recreate() *DeploymentStrategyDie {
	return d.DieStamp(func(r *appsv1.DeploymentStrategy) {
		r.Type = appsv1.RecreateDeploymentStrategyType
		r.RollingUpdate = nil
	})
}

func (d *DeploymentStrategyDie) RollingUpdateDie(fn func(d *RollingUpdateDeploymentDie)) *DeploymentStrategyDie {
	return d.DieStamp(func(r *appsv1.DeploymentStrategy) {
		r.Type = appsv1.RollingUpdateDeploymentStrategyType
		d := RollingUpdateDeploymentBlank.DieImmutable(false).DieFeedPtr(r.RollingUpdate)
		fn(d)
		r.RollingUpdate = d.DieReleasePtr()
	})
}

// +die
type _ = appsv1.RollingUpdateDeployment

// +die
type _ = appsv1.DeploymentStatus

func (d *DeploymentStatusDie) ConditionsDie(conditions ...*diemetav1.ConditionDie) *DeploymentStatusDie {
	return d.DieStamp(func(r *appsv1.DeploymentStatus) {
		r.Conditions = make([]appsv1.DeploymentCondition, len(conditions))
		for i := range conditions {
			c := conditions[i].DieRelease()
			r.Conditions[i] = appsv1.DeploymentCondition{
				Type:               appsv1.DeploymentConditionType(c.Type),
				Status:             corev1.ConditionStatus(c.Status),
				Reason:             c.Reason,
				Message:            c.Message,
				LastTransitionTime: c.LastTransitionTime,
			}
		}
	})
}
