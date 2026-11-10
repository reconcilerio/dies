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

// Package v2 is the newer version of a resource with two versions, and its
// hub: the other version, [reconciler.io/dies/v2/examples/sprocket/v1], is
// converted to and from it, by the conversion in the group's package,
// [reconciler.io/dies/v2/examples/sprocket].
package v2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +die:object:apiVersion=example.dev/v2,kind=Sprocket
// +kubebuilder:subresource:status

// Sprocket is a toothed wheel.
type Sprocket struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SprocketSpec   `json:"spec,omitempty"`
	Status SprocketStatus `json:"status,omitempty"`
}

// +die
type SprocketSpec struct {
	// Teeth was size in v1.
	Teeth int32 `json:"teeth,omitempty"`
	// Image was one string in v1.
	Image SprocketImage `json:"image,omitempty"`
	// Paused is new in v2.
	Paused bool `json:"paused,omitempty"`
}

// +die
type SprocketImage struct {
	Repository string `json:"repository,omitempty"`
	Tag        string `json:"tag,omitempty"`
}

// +die
type SprocketStatus struct {
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}
