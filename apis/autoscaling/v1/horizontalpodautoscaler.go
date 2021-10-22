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
	autoscalingv1 "k8s.io/api/autoscaling/v1"
)

// +die:object=true
type _ = autoscalingv1.HorizontalPodAutoscaler

// +die
type _ = autoscalingv1.HorizontalPodAutoscalerSpec

type horizontalPodAutoscalerSpecDieExtension interface {
	ScaleTargetRefDie(fn func(d CrossVersionObjectReferenceDie)) HorizontalPodAutoscalerSpecDie
}

func (d *horizontalPodAutoscalerSpecDie) ScaleTargetRefDie(fn func(d CrossVersionObjectReferenceDie)) HorizontalPodAutoscalerSpecDie {
	return d.DieStamp(func(r *autoscalingv1.HorizontalPodAutoscalerSpec) {
		d := CrossVersionObjectReferenceBlank.DieImmutable(false).DieFeed(r.ScaleTargetRef)
		fn(d)
		r.ScaleTargetRef = d.DieRelease()
	})
}

// +die
type _ = autoscalingv1.CrossVersionObjectReference

// +die
type _ = autoscalingv1.HorizontalPodAutoscalerStatus
