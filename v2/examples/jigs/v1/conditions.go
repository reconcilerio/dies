/*
Copyright 2026 the original author or authors.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	"reconciler.io/dies/v2/die/jigs"
)

// DeploymentConditions presents a Deployment's legacy conditions as
// [metav1.Condition] values. It is lossless: lastUpdateTime, which a
// [metav1.Condition] doesn't have, is kept from the condition of the same
// type, and observedGeneration, which an [appsv1.DeploymentCondition] doesn't
// have, is always zero, so writing anything else is refused.
var DeploymentConditions = jigs.Transform[[]appsv1.DeploymentCondition, []metav1.Condition]{
	To: func(legacy []appsv1.DeploymentCondition) []metav1.Condition {
		if legacy == nil {
			return nil
		}
		out := make([]metav1.Condition, len(legacy))
		for i, c := range legacy {
			out[i] = metav1.Condition{
				Type:               string(c.Type),
				Status:             metav1.ConditionStatus(c.Status),
				LastTransitionTime: c.LastTransitionTime,
				Reason:             c.Reason,
				Message:            c.Message,
			}
		}
		return out
	},
	From: func(conditions []metav1.Condition, legacy []appsv1.DeploymentCondition) []appsv1.DeploymentCondition {
		if conditions == nil {
			return nil
		}
		out := make([]appsv1.DeploymentCondition, len(conditions))
		for i, c := range conditions {
			out[i] = appsv1.DeploymentCondition{
				Type:               appsv1.DeploymentConditionType(c.Type),
				Status:             corev1.ConditionStatus(c.Status),
				LastTransitionTime: c.LastTransitionTime,
				Reason:             c.Reason,
				Message:            c.Message,
			}
			for _, l := range legacy {
				if string(l.Type) == c.Type {
					out[i].LastUpdateTime = l.LastUpdateTime
				}
			}
		}
		return out
	},
}

// ConditionedJig presents a status's conditions in their harmonized form.
var ConditionedJig = jigs.NewJig[Conditioned, *ConditionedDie]().
	MustClamp(dieappsv1.DeploymentStatusBlank,
		jigs.BindVia(ConditionedBlank.Conditions(), dieappsv1.DeploymentStatusBlank.Conditions(), DeploymentConditions),
	)
