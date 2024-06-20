/*
Copyright 2024 the original author or authors.

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

package v2

import (
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	diemetav1 "reconciler.io/dies/apis/meta/v1"
)

// +die:object=true,apiVersion=autoscaling/v2,kind=HorizontalPodAutoscaler
type _ = autoscalingv2.HorizontalPodAutoscaler

// +die
// +die:field:name=ScaleTargetRef,die=CrossVersionObjectReferenceDie
// +die:field:name=Behavior,die=HorizontalPodAutoscalerBehaviorDie,pointer=true
// +die:field:name=Metrics,die=MetricSpecDie,listType=atomic
type _ = autoscalingv2.HorizontalPodAutoscalerSpec

// +die
type _ = autoscalingv2.CrossVersionObjectReference

// +die
// +die:field:name=Object,die=ObjectMetricSourceDie,pointer=true
// +die:field:name=Pods,die=PodsMetricSourceDie,pointer=true
// +die:field:name=Resource,die=ResourceMetricSourceDie,pointer=true
// +die:field:name=ContainerResource,die=ContainerResourceMetricSourceDie,pointer=true
// +die:field:name=External,die=ExternalMetricSourceDie,pointer=true
type _ = autoscalingv2.MetricSpec

// +die
// +die:field:name=DescribedObject,die=CrossVersionObjectReferenceDie
// +die:field:name=Target,die=MetricTargetDie
// +die:field:name=Metric,die=MetricIdentifierDie
type _ = autoscalingv2.ObjectMetricSource

// +die
type _ = autoscalingv2.MetricTarget

// +die
// +die:field:name=Selector,package=_/meta/v1,die=LabelSelectorDie,pointer=true
type _ = autoscalingv2.MetricIdentifier

// +die
// +die:field:name=Metric,die=MetricIdentifierDie
// +die:field:name=Target,die=MetricTargetDie
type _ = autoscalingv2.PodsMetricSource

// +die
// +die:field:name=Target,die=MetricTargetDie
type _ = autoscalingv2.ResourceMetricSource

// +die
// +die:field:name=Target,die=MetricTargetDie
type _ = autoscalingv2.ContainerResourceMetricSource

// +die
// +die:field:name=Metric,die=MetricIdentifierDie
// +die:field:name=Target,die=MetricTargetDie
type _ = autoscalingv2.ExternalMetricSource

// +die
// +die:field:name=ScaleUp,die=HPAScalingRulesDie,pointer=true
// +die:field:name=ScaleDown,die=HPAScalingRulesDie,pointer=true
type _ = autoscalingv2.HorizontalPodAutoscalerBehavior

// +die
// +die:field:name=Policies,die=HPAScalingPolicyDie,listType=atomic
type _ = autoscalingv2.HPAScalingRules

// +die
type _ = autoscalingv2.HPAScalingPolicy

// +die
// +die:field:name=CurrentMetrics,die=MetricStatusDie,listType=atomic
type _ = autoscalingv2.HorizontalPodAutoscalerStatus

func (d *HorizontalPodAutoscalerStatusDie) ConditionsDie(conditions ...*diemetav1.ConditionDie) *HorizontalPodAutoscalerStatusDie {
	return d.DieStamp(func(r *autoscalingv2.HorizontalPodAutoscalerStatus) {
		r.Conditions = make([]autoscalingv2.HorizontalPodAutoscalerCondition, len(conditions))
		for i := range conditions {
			c := conditions[i].DieRelease()
			r.Conditions[i] = autoscalingv2.HorizontalPodAutoscalerCondition{
				Type:               autoscalingv2.HorizontalPodAutoscalerConditionType(c.Type),
				Status:             corev1.ConditionStatus(c.Status),
				LastTransitionTime: c.LastTransitionTime,
				Reason:             c.Reason,
				Message:            c.Message,
			}
		}
	})
}

// +die
// +die:field:name=Object,die=ObjectMetricStatusDie,pointer=true
// +die:field:name=Pods,die=PodsMetricStatusDie,pointer=true
// +die:field:name=Resource,die=ResourceMetricStatusDie,pointer=true
// +die:field:name=ContainerResource,die=ContainerResourceMetricStatusDie,pointer=true
// +die:field:name=External,die=ExternalMetricStatusDie,pointer=true
type _ = autoscalingv2.MetricStatus

// +die
// +die:field:name=Metric,die=MetricIdentifierDie
// +die:field:name=Current,die=MetricValueStatusDie
// +die:field:name=DescribedObject,die=CrossVersionObjectReferenceDie
type _ = autoscalingv2.ObjectMetricStatus

// +die
type _ = autoscalingv2.MetricValueStatus

// +die
// +die:field:name=Metric,die=MetricIdentifierDie
// +die:field:name=Current,die=MetricValueStatusDie
type _ = autoscalingv2.PodsMetricStatus

// +die
// +die:field:name=Current,die=MetricValueStatusDie
type _ = autoscalingv2.ResourceMetricStatus

// +die
// +die:field:name=Current,die=MetricValueStatusDie
type _ = autoscalingv2.ContainerResourceMetricStatus

// +die
// +die:field:name=Metric,die=MetricIdentifierDie
// +die:field:name=Current,die=MetricValueStatusDie
type _ = autoscalingv2.ExternalMetricStatus
