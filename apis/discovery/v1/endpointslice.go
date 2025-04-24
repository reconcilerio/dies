/*
Copyright 2024 the original author or authors.

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
	discoveryv1 "k8s.io/api/discovery/v1"
)

// +die:object=true,apiVersion=discovery.k8s.io/v1,kind=EndpointSlice
// +die:field:name=Endpoints,die=EndpointDie,listType=atomic
// +die:field:name=Ports,die=EndpointPortDie,listType=atomic
type _ = discoveryv1.EndpointSlice

// +die
// +die:field:name=Conditions,die=EndpointConditionsDie
// +die:field:name=TargetRef,package=_/core/v1,die=ObjectReferenceDie,pointer=true
// +die:field:name=Hints,die=EndpointHintsDie,pointer=true
type _ = discoveryv1.Endpoint

func (d *EndpointDie) AddDeprecatedTopology(key, value string) *EndpointDie {
	return d.DieStamp(func(r *discoveryv1.Endpoint) {
		if r.DeprecatedTopology == nil {
			r.DeprecatedTopology = map[string]string{}
		}
		r.DeprecatedTopology[key] = value
	})
}

// +die
type _ = discoveryv1.EndpointConditions

// +die
// +die:field:name=ForZones,die=ForZoneDie,listType=atomic
// +die:field:name=ForNodes,die=ForNodeDie,listType=atomic
type _ = discoveryv1.EndpointHints

// +die
type _ = discoveryv1.ForZone

// +die
type _ = discoveryv1.ForNode

// +die
type _ = discoveryv1.EndpointPort
