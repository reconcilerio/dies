/*
Copyright 2025 the original author or authors.

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

package v2

import (
	apidiscoveryv2 "k8s.io/api/apidiscovery/v2"
)

// +die:object=true,apiVersion=apidiscovery.k8s.io/v2,kind=APIGroupDiscovery
// +die:field:name=Versions,die=APIVersionDiscoveryDie,listType=map,listMapKey=Version
type _ = apidiscoveryv2.APIGroupDiscovery

// +die
// +die:field:name=Resources,die=APIResourceDiscoveryDie,listType=map,listMapKey=Resource
type _ = apidiscoveryv2.APIVersionDiscovery

// +die
// +die:field:name=ResponseKind,package=_/meta/v1,die=GroupVersionKindDie,pointer=true
// +die:field:name=Subresources,die=APISubresourceDiscoveryDie,listType=map,listMapKey=Subresource
type _ = apidiscoveryv2.APIResourceDiscovery

// +die
// +die:field:name=ResponseKind,package=_/meta/v1,die=GroupVersionKindDie,pointer=true
// +die:field:name=AcceptedTypes,package=_/meta/v1,die=GroupVersionKindDie,listType=atomic
type _ = apidiscoveryv2.APISubresourceDiscovery
