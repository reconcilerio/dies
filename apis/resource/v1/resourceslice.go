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

package v1

import (
	resourcev1 "k8s.io/api/resource/v1"
)

// +die:object=true
type _ = resourcev1.ResourceSlice

// +die
// +die:field:name=Pool,die=ResourcePoolDie
// +die:field:name=NodeSelector,package=_/core/v1,die=NodeSelectorDie,pointer=true
// +die:field:name=Devices,die=DeviceDie,listType=atomic
// +die:field:name=SharedCounters,die=CounterSetDie,listType=atomic
type _ = resourcev1.ResourceSliceSpec

// +die
type _ = resourcev1.ResourcePool

// +die:ignore=Attributes;Capacity
// +die:field:name=ConsumesCounters,die=DeviceCounterConsumptionDie,listType=atomic
// +die:field:name=NodeSelector,package=_/core/v1,die=NodeSelectorDie,pointer=true
// +die:field:name=Taints,die=DeviceTaintDie,listType=atomic
type _ = resourcev1.Device

// Attributes defines the set of attributes for this device.
// The name of each attribute must be unique in that set.
//
// The maximum number of attributes and capacities combined is 32.
//
// +optional
func (d *DeviceDie) Attributes(v map[resourcev1.QualifiedName]resourcev1.DeviceAttribute) *DeviceDie {
	return d.DieStamp(func(r *resourcev1.Device) {
		r.Attributes = v
	})
}

// Attributes defines the set of attributes for this device.
// The name of each attribute must be unique in that set.
//
// The maximum number of attributes and capacities combined is 32.
//
// +optional
func (d *DeviceDie) AttributesDie(v resourcev1.QualifiedName, fn func(d *DeviceAttributeDie)) *DeviceDie {
	return d.DieStamp(func(r *resourcev1.Device) {
		d := DeviceAttributeBlank.DieImmutable(false).DieFeed(r.Attributes[v])
		fn(d)
		r.Attributes[v] = d.DieRelease()
	})
}

// Capacity defines the set of capacities for this device.
// The name of each capacity must be unique in that set.
//
// The maximum number of attributes and capacities combined is 32.
//
// +optional
func (d *DeviceDie) Capacity(v map[resourcev1.QualifiedName]resourcev1.DeviceAttribute) *DeviceDie {
	return d.DieStamp(func(r *resourcev1.Device) {
		r.Attributes = v
	})
}

// Capacity defines the set of capacities for this device.
// The name of each capacity must be unique in that set.
//
// The maximum number of attributes and capacities combined is 32.
//
// +optional
func (d *DeviceDie) CapacityDie(v resourcev1.QualifiedName, fn func(d *DeviceCapacityDie)) *DeviceDie {
	return d.DieStamp(func(r *resourcev1.Device) {
		d := DeviceCapacityBlank.DieImmutable(false).DieFeed(r.Capacity[v])
		fn(d)
		r.Capacity[v] = d.DieRelease()
	})
}

// +die
type _ = resourcev1.DeviceAttribute

// +die
// +die:field:name=RequestPolicy,die=CapacityRequestPolicyDie,pointer=true
type _ = resourcev1.DeviceCapacity

// +die
// +die:field:name=ValidRange,die=CapacityRequestPolicyRangeDie,pointer=true
type _ = resourcev1.CapacityRequestPolicy

// +die
type _ = resourcev1.CapacityRequestPolicyRange

// +die:ignore=Counters
type _ = resourcev1.DeviceCounterConsumption

// Counters defines the counters that will be consumed by the device.
//
// The maximum number counters in a device is 32.
// In addition, the maximum number of all counters
// in all devices is 1024 (for example, 64 devices with
// 16 counters each).
//
// +required
func (d *DeviceCounterConsumptionDie) Counters(v map[string]resourcev1.Counter) *DeviceCounterConsumptionDie {
	return d.DieStamp(func(r *resourcev1.DeviceCounterConsumption) {
		r.Counters = v
	})
}

// Counters defines the counters that will be consumed by the device.
//
// The maximum number counters in a device is 32.
// In addition, the maximum number of all counters
// in all devices is 1024 (for example, 64 devices with
// 16 counters each).
//
// +required
func (d *DeviceCounterConsumptionDie) CounterDie(v string, fn func(d *CounterDie)) *DeviceCounterConsumptionDie {
	return d.DieStamp(func(r *resourcev1.DeviceCounterConsumption) {
		d := CounterBlank.DieImmutable(false).DieFeed(r.Counters[v])
		fn(d)
		r.Counters[v] = d.DieRelease()
	})
}

// +die
type _ = resourcev1.Counter

// +die:ignore=Counters
type _ = resourcev1.CounterSet

// Counters defines the set of counters for this CounterSet
// The name of each counter must be unique in that set and must be a DNS label.
//
// The maximum number of counters in all sets is 32.
//
// +required
func (d *CounterSetDie) Counters(v map[string]resourcev1.Counter) *CounterSetDie {
	return d.DieStamp(func(r *resourcev1.CounterSet) {
		r.Counters = v
	})
}

// Counters defines the set of counters for this CounterSet
// The name of each counter must be unique in that set and must be a DNS label.
//
// The maximum number of counters in all sets is 32.
//
// +required
func (d *CounterSetDie) CounterDie(v string, fn func(d *CounterDie)) *CounterSetDie {
	return d.DieStamp(func(r *resourcev1.CounterSet) {
		d := CounterBlank.DieImmutable(false).DieFeed(r.Counters[v])
		fn(d)
		r.Counters[v] = d.DieRelease()
	})
}

// +die
type _ = resourcev1.DeviceTaint
