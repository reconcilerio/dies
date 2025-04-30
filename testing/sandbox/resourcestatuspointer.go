/*
Copyright 2022 the original author or authors.

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

package sandbox

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// StatusPointer tests generating an object die when the status field is a pointer
// +die:object=true,apiVersion=internal.dies.reconciler.io,kind=ResourceStatusPointer
// +kubebuilder:object:root=true
type ResourceStatusPointer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ResourceStatusPointerSpec    `json:"spec,omitempty"`
	Status            *ResourceStatusPointerStatus `json:"status,omitempty"`
}

// +die
type ResourceStatusPointerSpec struct {
}

// +die
// +die:field:name=Conditions,die=ConditionDie,package=_/meta/v1,listMapKey=Type
type ResourceStatusPointerStatus struct {
	ObservedGeneration int64              `json:"observedGeneration"`
	Conditions         []metav1.Condition `json:"conditions"`
}
