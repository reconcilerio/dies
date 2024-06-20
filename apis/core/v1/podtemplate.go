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

package v1

import (
	corev1 "k8s.io/api/core/v1"
)

// +die:object=true,spec=-
// +die:field:name=Template,die=PodTemplateSpecDie
type _ = corev1.PodTemplate

// +die
// +die:field:name=ObjectMeta,method=MetadataDie,package=_/meta/v1,die=ObjectMetaDie
// +die:field:name=Spec,die=PodSpecDie
type _ = corev1.PodTemplateSpec
