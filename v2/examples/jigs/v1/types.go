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

// Package v1 is an example of jigs: dies that present different types the
// same way.
package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PodSpecable is anything that runs pods: its metadata and its pod template.
//
// +die
type PodSpecable struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Template          corev1.PodTemplateSpec `json:"template"`
}

// ContainerImage presents a container's image as its parts.
//
// +die
type ContainerImage struct {
	Name  string         `json:"name"`
	Image ImageReference `json:"image"`
}

// ImageReference is an image reference in parts: registry.example/repo:tag
// or registry.example/repo@sha256:digest.
//
// +die
type ImageReference struct {
	Repository string `json:"repository,omitempty"`
	Tag        string `json:"tag,omitempty"`
	Digest     string `json:"digest,omitempty"`
}

// Conditioned presents the conditions of a status in their harmonized form.
//
// +die
type Conditioned struct {
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// Workload presents what a ServiceBinding projects into a workload: every
// container, wherever it is, and the volumes.
//
// +die
type Workload struct {
	// +listType=map
	// +listMapKey=name
	Containers []corev1.Container `json:"containers,omitempty"`
	// +listType=map
	// +listMapKey=name
	Volumes []corev1.Volume `json:"volumes,omitempty"`
}

// Pipeline has containers in three shapes: a list of steps, an optional
// sidecar, and a builder that isn't a container but can be presented as one.
//
// +die
type Pipeline struct {
	// +listType=map
	// +listMapKey=name
	Steps   []corev1.Container `json:"steps,omitempty"`
	Sidecar *corev1.Container  `json:"sidecar,omitempty"`
	Builder Builder            `json:"builder"`
}

// +die
type Builder struct {
	Name  string `json:"name"`
	Image string `json:"image,omitempty"`
	// +listType=map
	// +listMapKey=name
	Env []corev1.EnvVar `json:"env,omitempty"`
}

// Batch keeps its containers in an atomic list, whose elements have no keys
// of their own, and has an optional init container.
//
// +die
type Batch struct {
	// +listType=atomic
	Containers []corev1.Container `json:"containers,omitempty"`
	Init       *corev1.Container  `json:"init,omitempty"`
}
