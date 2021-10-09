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

// +die:target=k8s.io/api/apps/v1.Deployment,object=true

// +die:target=k8s.io/api/apps/v1.DeploymentSpec
// +die:field:receiver=DeploymentSpecDie,name=Replicas,type=*int32
// +die:field:receiver=DeploymentSpecDie,name=Selector,type=*k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector
// +die:field:receiver=DeploymentSpecDie,name=Template,type=k8s.io/api/core/v1.PodTemplateSpec
// +die:field:receiver=DeploymentSpecDie,name=Strategy,type=k8s.io/api/apps/v1.DeploymentStrategy
// +die:field:receiver=DeploymentSpecDie,name=MinReadySeconds,type=int32
// +die:field:receiver=DeploymentSpecDie,name=RevisionHistoryLimit,type=*int32
// +die:field:receiver=DeploymentSpecDie,name=Paused,type=bool
// +die:field:receiver=DeploymentSpecDie,name=ProgressDeadlineSeconds,type=*int32

func (d *DeploymentSpecDie) TemplateDie(fn func(d *diecorev1.PodTemplateSpecDie)) *DeploymentSpecDie {
	return d.DieStamp(func(r *appsv1.DeploymentSpec) {
		d := diecorev1.PodTemplateSpecBlank.DieImmutable(false).DieFeed(r.Template)
		fn(d)
		r.Template = d.DieRelease()
	})
}

// +die:target=k8s.io/api/apps/v1.DeploymentStatus
// +die:field:receiver=DeploymentStatusDie,name=ObservedGeneration,type=int64
// +die:field:receiver=DeploymentStatusDie,name=Replicas,type=int32
// +die:field:receiver=DeploymentStatusDie,name=UpdatedReplicas,type=int32
// +die:field:receiver=DeploymentStatusDie,name=ReadyReplicas,type=int32
// +die:field:receiver=DeploymentStatusDie,name=AvailableReplicas,type=int32
// +die:field:receiver=DeploymentStatusDie,name=UnavailableReplicas,type=int32
// +die:field:receiver=DeploymentStatusDie,name=Conditions,type=[]k8s.io/api/apps/v1.DeploymentCondition
// +die:field:receiver=DeploymentStatusDie,name=CollisionCount,type=*int32

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
