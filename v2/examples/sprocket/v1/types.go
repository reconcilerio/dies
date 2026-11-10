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

// Package v1 is the older version of a resource with two versions, which is
// converted to and from the hub,
// [reconciler.io/dies/v2/examples/sprocket/v2], by the conversion in the
// group's package, [reconciler.io/dies/v2/examples/sprocket].
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +die:object:apiVersion=example.dev/v1,kind=Sprocket
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
	// Size is the number of teeth: v2 calls it teeth.
	Size int32 `json:"size,omitempty"`
	// Image is a reference, repository:tag: v2 splits it.
	Image string `json:"image,omitempty"`
	// Lubricant has no place in v2.
	Lubricant string `json:"lubricant,omitempty"`
}

// +die
type SprocketStatus struct {
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}
