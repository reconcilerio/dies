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

// Package appsv1 holds a subset of apps/v1 types, shaped as upstream, and
// hand-written dies that stand in for generated ones.
package appsv1

import (
	"reconciler.io/dies/v2/spike/corev1"
	"reconciler.io/dies/v2/spike/metav1"
)

type Deployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              DeploymentSpec   `json:"spec,omitempty"`
	Status            DeploymentStatus `json:"status,omitempty" die:"subresource=status"`
}

type DeploymentSpec struct {
	Replicas *int32                 `json:"replicas,omitempty"`
	Paused   bool                   `json:"paused,omitempty"`
	Template corev1.PodTemplateSpec `json:"template"`
}

type DeploymentStatus struct {
	ObservedGeneration int64              `json:"observedGeneration,omitempty"`
	Replicas           int32              `json:"replicas,omitempty"`
	Conditions         []metav1.Condition `json:"conditions,omitempty" die:"listMap=type"`
}
