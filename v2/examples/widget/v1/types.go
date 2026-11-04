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

// Package v1 is an example of dies for types declared locally. They are
// marked +die just as aliased upstream types are, and their fields of
// upstream types use the dies in [reconciler.io/dies/v2/apis].
package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	gadgetv1 "reconciler.io/dies/v2/examples/gadget/v1"
)

// +die:object:apiVersion=example.dev/v1,kind=Widget
// +kubebuilder:subresource:status

// Widget runs a pod template with named ports.
type Widget struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WidgetSpec   `json:"spec,omitempty"`
	Status WidgetStatus `json:"status,omitempty"`
}

// +die:object

// Templated is a duck type of the resources with a pod template, as a Widget
// has. It has no kind, so it is not registered: it is a view of other kinds,
// not a resource of its own.
type Templated struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec TemplatedSpec `json:"spec,omitempty"`
}

// +die
type TemplatedSpec struct {
	// Template for the resource's pods.
	Template corev1.PodTemplateSpec `json:"template"`
}

// +die
type WidgetSpec struct {
	// Ports exposed by the widget.
	// +listType=map
	// +listMapKey=name
	Ports []WidgetPort `json:"ports,omitempty"`
	// Template for the widget's pods.
	Template corev1.PodTemplateSpec `json:"template"`
	// Gadget embeds a type from another die package.
	Gadget WidgetGadget `json:"gadget,omitempty"`
	// Match has a field whose accessor is renamed.
	Match WidgetMatch `json:"match,omitempty"`
}

// WidgetMatch's equal field would have an accessor named Equal, which
// collides with the die method that compares dies, so it is renamed.
//
// +die
// +die:field:name=equal,accessor=Exact
type WidgetMatch struct {
	Equal bool `json:"equal,omitempty"`
	// Note accepts null, so an apply clears it with null, as a custom
	// resource's validation rejects null for fields that don't.
	// +nullable
	Note *string `json:"note,omitempty"`
}

// +die
type WidgetGadget struct {
	Label               string `json:"label,omitempty"`
	gadgetv1.GadgetSpec `json:",inline"`
}

// +die
type WidgetPort struct {
	Name string `json:"name"`
	Port int32  `json:"port"`
}

// +die
type WidgetStatus struct {
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}
