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
	resource "k8s.io/apimachinery/pkg/api/resource"
	diemetav1 "reconciler.io/dies/apis/meta/v1"
)

// +die:object=true,apiVersion=autoscaling/v2,kind=HorizontalPodAutoscaler
type _ = autoscalingv2.HorizontalPodAutoscaler

// +die
type _ = autoscalingv2.HorizontalPodAutoscalerSpec

func (d *HorizontalPodAutoscalerSpecDie) ScaleTargetRefDie(fn func(d *CrossVersionObjectReferenceDie)) *HorizontalPodAutoscalerSpecDie {
	return d.DieStamp(func(r *autoscalingv2.HorizontalPodAutoscalerSpec) {
		d := CrossVersionObjectReferenceBlank.DieImmutable(false).DieFeed(r.ScaleTargetRef)
		fn(d)
		r.ScaleTargetRef = d.DieRelease()
	})
}

func (d *HorizontalPodAutoscalerSpecDie) ConditionsDie(metrics ...*MetricSpecDie) *HorizontalPodAutoscalerSpecDie {
	return d.DieStamp(func(r *autoscalingv2.HorizontalPodAutoscalerSpec) {
		r.Metrics = make([]autoscalingv2.MetricSpec, len(metrics))
		for i := range metrics {
			r.Metrics[i] = metrics[i].DieRelease()
		}
	})
}

func (d *HorizontalPodAutoscalerSpecDie) BehaviorDie(fn func(d *HorizontalPodAutoscalerBehaviorDie)) *HorizontalPodAutoscalerSpecDie {
	return d.DieStamp(func(r *autoscalingv2.HorizontalPodAutoscalerSpec) {
		d := HorizontalPodAutoscalerBehaviorBlank.DieImmutable(false).DieFeedPtr(r.Behavior)
		fn(d)
		r.Behavior = d.DieReleasePtr()
	})
}

// +die
type _ = autoscalingv2.CrossVersionObjectReference

// +die
type _ = autoscalingv2.MetricSpec

func (d *MetricSpecDie) BehaviorDie(fn func(d *ObjectMetricSourceDie)) *MetricSpecDie {
	return d.DieStamp(func(r *autoscalingv2.MetricSpec) {
		d := ObjectMetricSourceBlank.DieImmutable(false).DieFeedPtr(r.Object)
		fn(d)
		r.Object = d.DieReleasePtr()
	})
}

func (d *MetricSpecDie) PodsDie(fn func(d *PodsMetricSourceDie)) *MetricSpecDie {
	return d.DieStamp(func(r *autoscalingv2.MetricSpec) {
		d := PodsMetricSourceBlank.DieImmutable(false).DieFeedPtr(r.Pods)
		fn(d)
		r.Pods = d.DieReleasePtr()
	})
}

func (d *MetricSpecDie) ResourceDie(fn func(d *ResourceMetricSourceDie)) *MetricSpecDie {
	return d.DieStamp(func(r *autoscalingv2.MetricSpec) {
		d := ResourceMetricSourceBlank.DieImmutable(false).DieFeedPtr(r.Resource)
		fn(d)
		r.Resource = d.DieReleasePtr()
	})
}

func (d *MetricSpecDie) ContainerResourceDie(fn func(d *ContainerResourceMetricSourceDie)) *MetricSpecDie {
	return d.DieStamp(func(r *autoscalingv2.MetricSpec) {
		d := ContainerResourceMetricSourceBlank.DieImmutable(false).DieFeedPtr(r.ContainerResource)
		fn(d)
		r.ContainerResource = d.DieReleasePtr()
	})
}

func (d *MetricSpecDie) ExternalDie(fn func(d *ExternalMetricSourceDie)) *MetricSpecDie {
	return d.DieStamp(func(r *autoscalingv2.MetricSpec) {
		d := ExternalMetricSourceBlank.DieImmutable(false).DieFeedPtr(r.External)
		fn(d)
		r.External = d.DieReleasePtr()
	})
}

// +die
type _ = autoscalingv2.ObjectMetricSource

func (d *ObjectMetricSourceDie) DescribedObjectDie(fn func(d *CrossVersionObjectReferenceDie)) *ObjectMetricSourceDie {
	return d.DieStamp(func(r *autoscalingv2.ObjectMetricSource) {
		d := CrossVersionObjectReferenceBlank.DieImmutable(false).DieFeed(r.DescribedObject)
		fn(d)
		r.DescribedObject = d.DieRelease()
	})
}

func (d *ObjectMetricSourceDie) TargetDie(fn func(d *MetricTargetDie)) *ObjectMetricSourceDie {
	return d.DieStamp(func(r *autoscalingv2.ObjectMetricSource) {
		d := MetricTargetBlank.DieImmutable(false).DieFeed(r.Target)
		fn(d)
		r.Target = d.DieRelease()
	})
}

func (d *ObjectMetricSourceDie) MetricDie(fn func(d *MetricIdentifierDie)) *ObjectMetricSourceDie {
	return d.DieStamp(func(r *autoscalingv2.ObjectMetricSource) {
		d := MetricIdentifierBlank.DieImmutable(false).DieFeed(r.Metric)
		fn(d)
		r.Metric = d.DieRelease()
	})
}

// +die
type _ = autoscalingv2.MetricTarget

func (d *MetricTargetDie) ValueString(quantity string) *MetricTargetDie {
	q := resource.MustParse(quantity)
	return d.Value(&q)
}

func (d *MetricTargetDie) AverageValueString(quantity string) *MetricTargetDie {
	q := resource.MustParse(quantity)
	return d.AverageValue(&q)
}

// +die
type _ = autoscalingv2.MetricIdentifier

func (d *MetricIdentifierDie) SelectorDie(fn func(d *diemetav1.LabelSelectorDie)) *MetricIdentifierDie {
	return d.DieStamp(func(r *autoscalingv2.MetricIdentifier) {
		d := diemetav1.LabelSelectorBlank.DieImmutable(false).DieFeedPtr(r.Selector)
		fn(d)
		r.Selector = d.DieReleasePtr()
	})
}

// +die
type _ = autoscalingv2.PodsMetricSource

func (d *PodsMetricSourceDie) MetricDie(fn func(d *MetricIdentifierDie)) *PodsMetricSourceDie {
	return d.DieStamp(func(r *autoscalingv2.PodsMetricSource) {
		d := MetricIdentifierBlank.DieImmutable(false).DieFeed(r.Metric)
		fn(d)
		r.Metric = d.DieRelease()
	})
}

func (d *PodsMetricSourceDie) TargetDie(fn func(d *MetricTargetDie)) *PodsMetricSourceDie {
	return d.DieStamp(func(r *autoscalingv2.PodsMetricSource) {
		d := MetricTargetBlank.DieImmutable(false).DieFeed(r.Target)
		fn(d)
		r.Target = d.DieRelease()
	})
}

// +die
type _ = autoscalingv2.ResourceMetricSource

func (d *ResourceMetricSourceDie) TargetDie(fn func(d *MetricTargetDie)) *ResourceMetricSourceDie {
	return d.DieStamp(func(r *autoscalingv2.ResourceMetricSource) {
		d := MetricTargetBlank.DieImmutable(false).DieFeed(r.Target)
		fn(d)
		r.Target = d.DieRelease()
	})
}

// +die
type _ = autoscalingv2.ContainerResourceMetricSource

func (d *ContainerResourceMetricSourceDie) TargetDie(fn func(d *MetricTargetDie)) *ContainerResourceMetricSourceDie {
	return d.DieStamp(func(r *autoscalingv2.ContainerResourceMetricSource) {
		d := MetricTargetBlank.DieImmutable(false).DieFeed(r.Target)
		fn(d)
		r.Target = d.DieRelease()
	})
}

// +die
type _ = autoscalingv2.ExternalMetricSource

func (d *ExternalMetricSourceDie) MetricDie(fn func(d *MetricIdentifierDie)) *ExternalMetricSourceDie {
	return d.DieStamp(func(r *autoscalingv2.ExternalMetricSource) {
		d := MetricIdentifierBlank.DieImmutable(false).DieFeed(r.Metric)
		fn(d)
		r.Metric = d.DieRelease()
	})
}

func (d *ExternalMetricSourceDie) TargetDie(fn func(d *MetricTargetDie)) *ExternalMetricSourceDie {
	return d.DieStamp(func(r *autoscalingv2.ExternalMetricSource) {
		d := MetricTargetBlank.DieImmutable(false).DieFeed(r.Target)
		fn(d)
		r.Target = d.DieRelease()
	})
}

// +die
type _ = autoscalingv2.HorizontalPodAutoscalerBehavior

func (d *HorizontalPodAutoscalerBehaviorDie) ScaleUpDie(fn func(d *HPAScalingRulesDie)) *HorizontalPodAutoscalerBehaviorDie {
	return d.DieStamp(func(r *autoscalingv2.HorizontalPodAutoscalerBehavior) {
		d := HPAScalingRulesBlank.DieImmutable(false).DieFeedPtr(r.ScaleUp)
		fn(d)
		r.ScaleUp = d.DieReleasePtr()
	})
}

func (d *HorizontalPodAutoscalerBehaviorDie) ScaleDownDie(fn func(d *HPAScalingRulesDie)) *HorizontalPodAutoscalerBehaviorDie {
	return d.DieStamp(func(r *autoscalingv2.HorizontalPodAutoscalerBehavior) {
		d := HPAScalingRulesBlank.DieImmutable(false).DieFeedPtr(r.ScaleDown)
		fn(d)
		r.ScaleDown = d.DieReleasePtr()
	})
}

// +die
type _ = autoscalingv2.HPAScalingRules

func (d *HPAScalingRulesDie) PoliciesDie(policies ...*HPAScalingPolicyDie) *HPAScalingRulesDie {
	return d.DieStamp(func(r *autoscalingv2.HPAScalingRules) {
		r.Policies = make([]autoscalingv2.HPAScalingPolicy, len(policies))
		for i := range policies {
			r.Policies[i] = policies[i].DieRelease()
		}
	})
}

// +die
type _ = autoscalingv2.HPAScalingPolicy

// +die
type _ = autoscalingv2.HorizontalPodAutoscalerStatus

func (d *HorizontalPodAutoscalerStatusDie) CurrentMetricsDie(metrics ...*MetricStatusDie) *HorizontalPodAutoscalerStatusDie {
	return d.DieStamp(func(r *autoscalingv2.HorizontalPodAutoscalerStatus) {
		r.CurrentMetrics = make([]autoscalingv2.MetricStatus, len(metrics))
		for i := range metrics {
			r.CurrentMetrics[i] = metrics[i].DieRelease()
		}
	})
}

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
type _ = autoscalingv2.MetricStatus

func (d *MetricStatusDie) ObjectDie(fn func(d *ObjectMetricStatusDie)) *MetricStatusDie {
	return d.DieStamp(func(r *autoscalingv2.MetricStatus) {
		d := ObjectMetricStatusBlank.DieImmutable(false).DieFeedPtr(r.Object)
		fn(d)
		r.Object = d.DieReleasePtr()
	})
}

func (d *MetricStatusDie) PodDie(fn func(d *PodsMetricStatusDie)) *MetricStatusDie {
	return d.DieStamp(func(r *autoscalingv2.MetricStatus) {
		d := PodsMetricStatusBlank.DieImmutable(false).DieFeedPtr(r.Pods)
		fn(d)
		r.Pods = d.DieReleasePtr()
	})
}

func (d *MetricStatusDie) ResourceDie(fn func(d *ResourceMetricStatusDie)) *MetricStatusDie {
	return d.DieStamp(func(r *autoscalingv2.MetricStatus) {
		d := ResourceMetricStatusBlank.DieImmutable(false).DieFeedPtr(r.Resource)
		fn(d)
		r.Resource = d.DieReleasePtr()
	})
}

func (d *MetricStatusDie) ContainerResourceDie(fn func(d *ContainerResourceMetricStatusDie)) *MetricStatusDie {
	return d.DieStamp(func(r *autoscalingv2.MetricStatus) {
		d := ContainerResourceMetricStatusBlank.DieImmutable(false).DieFeedPtr(r.ContainerResource)
		fn(d)
		r.ContainerResource = d.DieReleasePtr()
	})
}

func (d *MetricStatusDie) ExternalDie(fn func(d *ExternalMetricStatusDie)) *MetricStatusDie {
	return d.DieStamp(func(r *autoscalingv2.MetricStatus) {
		d := ExternalMetricStatusBlank.DieImmutable(false).DieFeedPtr(r.External)
		fn(d)
		r.External = d.DieReleasePtr()
	})
}

// +die
type _ = autoscalingv2.ObjectMetricStatus

func (d *ObjectMetricStatusDie) MetricDie(fn func(d *MetricIdentifierDie)) *ObjectMetricStatusDie {
	return d.DieStamp(func(r *autoscalingv2.ObjectMetricStatus) {
		d := MetricIdentifierBlank.DieImmutable(false).DieFeed(r.Metric)
		fn(d)
		r.Metric = d.DieRelease()
	})
}

func (d *ObjectMetricStatusDie) CurrentDie(fn func(d *MetricValueStatusDie)) *ObjectMetricStatusDie {
	return d.DieStamp(func(r *autoscalingv2.ObjectMetricStatus) {
		d := MetricValueStatusBlank.DieImmutable(false).DieFeed(r.Current)
		fn(d)
		r.Current = d.DieRelease()
	})
}

func (d *ObjectMetricStatusDie) DescribedObjectDie(fn func(d *CrossVersionObjectReferenceDie)) *ObjectMetricStatusDie {
	return d.DieStamp(func(r *autoscalingv2.ObjectMetricStatus) {
		d := CrossVersionObjectReferenceBlank.DieImmutable(false).DieFeed(r.DescribedObject)
		fn(d)
		r.DescribedObject = d.DieRelease()
	})
}

// +die
type _ = autoscalingv2.MetricValueStatus

func (d *MetricValueStatusDie) ValueString(quantity string) *MetricValueStatusDie {
	q := resource.MustParse(quantity)
	return d.Value(&q)
}

func (d *MetricValueStatusDie) AverageValueString(quantity string) *MetricValueStatusDie {
	q := resource.MustParse(quantity)
	return d.AverageValue(&q)
}

// +die
type _ = autoscalingv2.PodsMetricStatus

func (d *PodsMetricStatusDie) MetricDie(fn func(d *MetricIdentifierDie)) *PodsMetricStatusDie {
	return d.DieStamp(func(r *autoscalingv2.PodsMetricStatus) {
		d := MetricIdentifierBlank.DieImmutable(false).DieFeed(r.Metric)
		fn(d)
		r.Metric = d.DieRelease()
	})
}

func (d *PodsMetricStatusDie) CurrentDie(fn func(d *MetricValueStatusDie)) *PodsMetricStatusDie {
	return d.DieStamp(func(r *autoscalingv2.PodsMetricStatus) {
		d := MetricValueStatusBlank.DieImmutable(false).DieFeed(r.Current)
		fn(d)
		r.Current = d.DieRelease()
	})
}

// +die
type _ = autoscalingv2.ResourceMetricStatus

func (d *ResourceMetricStatusDie) CurrentDie(fn func(d *MetricValueStatusDie)) *ResourceMetricStatusDie {
	return d.DieStamp(func(r *autoscalingv2.ResourceMetricStatus) {
		d := MetricValueStatusBlank.DieImmutable(false).DieFeed(r.Current)
		fn(d)
		r.Current = d.DieRelease()
	})
}

// +die
type _ = autoscalingv2.ContainerResourceMetricStatus

func (d *ContainerResourceMetricStatusDie) CurrentDie(fn func(d *MetricValueStatusDie)) *ContainerResourceMetricStatusDie {
	return d.DieStamp(func(r *autoscalingv2.ContainerResourceMetricStatus) {
		d := MetricValueStatusBlank.DieImmutable(false).DieFeed(r.Current)
		fn(d)
		r.Current = d.DieRelease()
	})
}

// +die
type _ = autoscalingv2.ExternalMetricStatus

func (d *ExternalMetricStatusDie) MetricDie(fn func(d *MetricIdentifierDie)) *ExternalMetricStatusDie {
	return d.DieStamp(func(r *autoscalingv2.ExternalMetricStatus) {
		d := MetricIdentifierBlank.DieImmutable(false).DieFeed(r.Metric)
		fn(d)
		r.Metric = d.DieRelease()
	})
}

func (d *ExternalMetricStatusDie) CurrentDie(fn func(d *MetricValueStatusDie)) *ExternalMetricStatusDie {
	return d.DieStamp(func(r *autoscalingv2.ExternalMetricStatus) {
		d := MetricValueStatusBlank.DieImmutable(false).DieFeed(r.Current)
		fn(d)
		r.Current = d.DieRelease()
	})
}
