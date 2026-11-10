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

// Package corev1 holds a subset of core/v1 types, shaped as upstream, and
// hand-written dies that stand in for generated ones. The `die` struct tags
// stand in for the upstream +listType, +listMapKey and union markers. As
// these aren't k8s.io types, an apply sends null only for fields marked
// nullable, as for a custom resource, so the fields the tests clear are.
package corev1

import "reconciler.io/dies/v2/spike/metav1"

type PodTemplateSpec struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PodSpec `json:"spec,omitempty"`
}

type PodSpec struct {
	Volumes        []Volume    `json:"volumes,omitempty" die:"listMap=name"`
	InitContainers []Container `json:"initContainers,omitempty" die:"listMap=name"`
	// not omitempty, as upstream: encoding/json writes nil as null
	Containers                   []Container       `json:"containers" die:"listMap=name"`
	NodeSelector                 map[string]string `json:"nodeSelector,omitempty" die:"nullable"`
	ServiceAccountName           string            `json:"serviceAccountName,omitempty"`
	AutomountServiceAccountToken *bool             `json:"automountServiceAccountToken,omitempty" die:"nullable"`
}

type Container struct {
	Name         string          `json:"name"`
	Image        string          `json:"image,omitempty"`
	Args         []string        `json:"args,omitempty" die:"atomic"`
	Ports        []ContainerPort `json:"ports,omitempty" die:"listMap=containerPort,protocol"`
	Env          []EnvVar        `json:"env,omitempty" die:"listMap=name"`
	VolumeMounts []VolumeMount   `json:"volumeMounts,omitempty" die:"listMap=mountPath"`
}

type ContainerPort struct {
	Name          string `json:"name,omitempty"`
	ContainerPort int32  `json:"containerPort"`
	// +default="TCP"
	Protocol string `json:"protocol,omitempty" die:"default=\"TCP\""`
}

type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
}

type VolumeMount struct {
	Name      string `json:"name"`
	ReadOnly  bool   `json:"readOnly,omitempty"`
	MountPath string `json:"mountPath"`
}

type Volume struct {
	Name string `json:"name"`
	// inlined with an empty json name, as upstream
	VolumeSource `json:""`
}

type VolumeSource struct {
	Secret   *SecretVolumeSource   `json:"secret,omitempty" die:"union=source;nullable"`
	EmptyDir *EmptyDirVolumeSource `json:"emptyDir,omitempty" die:"union=source"`
}

type SecretVolumeSource struct {
	SecretName string `json:"secretName,omitempty"`
}

type EmptyDirVolumeSource struct {
	Medium string `json:"medium,omitempty"`
}
