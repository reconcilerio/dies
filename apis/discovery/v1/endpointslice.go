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
	diecorev1 "reconciler.io/dies/apis/core/v1"
)

// +die:object=true,apiVersion=discovery.k8s.io/v1,kind=EndpointSlice
type _ = discoveryv1.EndpointSlice

func (d *EndpointSliceDie) EndpointsDie(endpoints ...*EndpointDie) *EndpointSliceDie {
	return d.DieStamp(func(r *discoveryv1.EndpointSlice) {
		r.Endpoints = make([]discoveryv1.Endpoint, len(endpoints))
		for i := range endpoints {
			r.Endpoints[i] = endpoints[i].DieRelease()
		}
	})
}

func (d *EndpointSliceDie) PortsDie(ports ...*EndpointPortDie) *EndpointSliceDie {
	return d.DieStamp(func(r *discoveryv1.EndpointSlice) {
		r.Ports = make([]discoveryv1.EndpointPort, len(ports))
		for i := range ports {
			r.Ports[i] = ports[i].DieRelease()
		}
	})
}

// +die
type _ = discoveryv1.Endpoint

func (d *EndpointDie) ConditionsDie(fn func(d *EndpointConditionsDie)) *EndpointDie {
	return d.DieStamp(func(r *discoveryv1.Endpoint) {
		d := EndpointConditionsBlank.DieImmutable(false).DieFeed(r.Conditions)
		fn(d)
		r.Conditions = d.DieRelease()
	})
}

func (d *EndpointDie) TargetRefDie(fn func(d *diecorev1.ObjectReferenceDie)) *EndpointDie {
	return d.DieStamp(func(r *discoveryv1.Endpoint) {
		d := diecorev1.ObjectReferenceBlank.DieImmutable(false).DieFeedPtr(r.TargetRef)
		fn(d)
		r.TargetRef = d.DieReleasePtr()
	})
}

func (d *EndpointDie) AddDeprecatedTopology(key, value string) *EndpointDie {
	return d.DieStamp(func(r *discoveryv1.Endpoint) {
		if r.DeprecatedTopology == nil {
			r.DeprecatedTopology = map[string]string{}
		}
		r.DeprecatedTopology[key] = value
	})
}

func (d *EndpointDie) HintsDie(fn func(d *EndpointHintsDie)) *EndpointDie {
	return d.DieStamp(func(r *discoveryv1.Endpoint) {
		d := EndpointHintsBlank.DieImmutable(false).DieFeedPtr(r.Hints)
		fn(d)
		r.Hints = d.DieReleasePtr()
	})
}

// +die
type _ = discoveryv1.EndpointConditions

// +die
type _ = discoveryv1.EndpointHints

func (d *EndpointHintsDie) ForZonesDie(zones ...*ForZoneDie) *EndpointHintsDie {
	return d.DieStamp(func(r *discoveryv1.EndpointHints) {
		r.ForZones = make([]discoveryv1.ForZone, len(zones))
		for i := range zones {
			r.ForZones[i] = zones[i].DieRelease()
		}
	})
}

// +die
type _ = discoveryv1.ForZone

// +die
type _ = discoveryv1.EndpointPort
