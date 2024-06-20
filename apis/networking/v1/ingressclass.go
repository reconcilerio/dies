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
	networkingv1 "k8s.io/api/networking/v1"
)

// +die:object=true,apiVersion=networking.k8s.io/v1,kind=IngressClass
type _ = networkingv1.IngressClass

// +die
// +die:field:name=Parameters,die=IngressClassParametersReferenceDie,pointer=true
type _ = networkingv1.IngressClassSpec

// +die
type _ = networkingv1.IngressClassParametersReference
