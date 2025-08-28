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
	"k8s.io/apimachinery/pkg/api/resource"
)

// +die:object=true,apiVersion=resource.k8s.io/v1,kind=ResourceClaim
type _ = resourcev1.ResourceClaim

// +die
// +die:field:name=Devices,die=DeviceClaimDie
type _ = resourcev1.ResourceClaimSpec

// +die
// +die:field:name=Requests,die=DeviceRequestDie,listType=atomic
// +die:field:name=Constraints,die=DeviceConstraintDie,listType=atomic
// +die:field:name=Config,die=DeviceClaimConfigurationDie,listType=atomic
type _ = resourcev1.DeviceClaim

// +die
// +die:field:name=Exactly,die=ExactDeviceRequestDie,pointer=true
// +die:field:name=FirstAvailable,die=DeviceSubRequestDie,listType=atomic
type _ = resourcev1.DeviceRequest

// +die
// +die:field:name=Selectors,die=DeviceSelectorDie,listType=atomic
// +die:field:name=Tolerations,die=DeviceTolerationDie,listType=atomic
// +die:field:name=Capacity,die=CapacityRequirementsDie,pointer=true
type _ = resourcev1.DeviceSubRequest

// +die:ignore=Requests
type _ = resourcev1.CapacityRequirements

// Requests represent individual device resource requests for distinct resources,
// all of which must be provided by the device.
//
// This value is used as an additional filtering condition against the available capacity on the device.
// This is semantically equivalent to a CEL selector with
// `device.capacity[<domain>].<name>.compareTo(quantity(<request quantity>)) >= 0`.
// For example, device.capacity['test-driver.cdi.k8s.io'].counters.compareTo(quantity('2')) >= 0.
//
// When a requestPolicy is defined, the requested amount is adjusted upward
// to the nearest valid value based on the policy.
// If the requested amount cannot be adjusted to a valid value—because it exceeds what the requestPolicy allows—
// the device is considered ineligible for allocation.
//
// For any capacity that is not explicitly requested:
//   - If no requestPolicy is set, the default consumed capacity is equal to the full device capacity
//     (i.e., the whole device is claimed).
//   - If a requestPolicy is set, the default consumed capacity is determined according to that policy.
//
// If the device allows multiple allocation,
// the aggregated amount across all requests must not exceed the capacity value.
// The consumed capacity, which may be adjusted based on the requestPolicy if defined,
// is recorded in the resource claim’s status.devices[*].consumedCapacity field.
//
// +optional
func (d *CapacityRequirementsDie) Requests(v map[resourcev1.QualifiedName]resource.Quantity) *CapacityRequirementsDie {
	return d.DieStamp(func(r *resourcev1.CapacityRequirements) {
		r.Requests = v
	})
}

// AddRequest sets a single quantity on the Requests map.
//
// Requests represent individual device resource requests for distinct resources,
// all of which must be provided by the device.
//
// This value is used as an additional filtering condition against the available capacity on the device.
// This is semantically equivalent to a CEL selector with
// `device.capacity[<domain>].<name>.compareTo(quantity(<request quantity>)) >= 0`.
// For example, device.capacity['test-driver.cdi.k8s.io'].counters.compareTo(quantity('2')) >= 0.
//
// When a requestPolicy is defined, the requested amount is adjusted upward
// to the nearest valid value based on the policy.
// If the requested amount cannot be adjusted to a valid value—because it exceeds what the requestPolicy allows—
// the device is considered ineligible for allocation.
//
// For any capacity that is not explicitly requested:
//   - If no requestPolicy is set, the default consumed capacity is equal to the full device capacity
//     (i.e., the whole device is claimed).
//   - If a requestPolicy is set, the default consumed capacity is determined according to that policy.
//
// If the device allows multiple allocation,
// the aggregated amount across all requests must not exceed the capacity value.
// The consumed capacity, which may be adjusted based on the requestPolicy if defined,
// is recorded in the resource claim’s status.devices[*].consumedCapacity field.
//
// +optional
func (d *CapacityRequirementsDie) AddRequest(name resourcev1.QualifiedName, quantity resource.Quantity) *CapacityRequirementsDie {
	return d.DieStamp(func(r *resourcev1.CapacityRequirements) {
		if r.Requests == nil {
			r.Requests = map[resourcev1.QualifiedName]resource.Quantity{}
		}
		r.Requests[name] = quantity
	})
}

// AddRequestString parses the quantity setting a single value on the Requests map. Panics if the string is not parsable.
//
// Requests represent individual device resource requests for distinct resources,
// all of which must be provided by the device.
//
// This value is used as an additional filtering condition against the available capacity on the device.
// This is semantically equivalent to a CEL selector with
// `device.capacity[<domain>].<name>.compareTo(quantity(<request quantity>)) >= 0`.
// For example, device.capacity['test-driver.cdi.k8s.io'].counters.compareTo(quantity('2')) >= 0.
//
// When a requestPolicy is defined, the requested amount is adjusted upward
// to the nearest valid value based on the policy.
// If the requested amount cannot be adjusted to a valid value—because it exceeds what the requestPolicy allows—
// the device is considered ineligible for allocation.
//
// For any capacity that is not explicitly requested:
//   - If no requestPolicy is set, the default consumed capacity is equal to the full device capacity
//     (i.e., the whole device is claimed).
//   - If a requestPolicy is set, the default consumed capacity is determined according to that policy.
//
// If the device allows multiple allocation,
// the aggregated amount across all requests must not exceed the capacity value.
// The consumed capacity, which may be adjusted based on the requestPolicy if defined,
// is recorded in the resource claim’s status.devices[*].consumedCapacity field.
//
// +optional
func (d *CapacityRequirementsDie) AddRequestString(name resourcev1.QualifiedName, quantity string) *CapacityRequirementsDie {
	q := resource.MustParse(quantity)
	return d.AddRequest(name, q)
}

// +die
// +die:field:name=Selectors,die=DeviceSelectorDie,listType=atomic
// +die:field:name=Tolerations,die=DeviceTolerationDie,listType=atomic
type _ = resourcev1.ExactDeviceRequest

// +die
type _ = resourcev1.DeviceToleration

// +die
type _ = resourcev1.DeviceConstraint

// +die
// +die:field:name=DeviceConfiguration,die=DeviceConfigurationDie
type _ = resourcev1.DeviceClaimConfiguration

// +die
// +die:field:name=Allocation,die=AllocationResultDie,pointer=true
// +die:field:name=ReservedFor,die=ResourceClaimConsumerReferenceDie,listType=map,listMapKey=UID,listMapKeyPackage=k8s.io/apimachinery/pkg/types,listMapKeyType=UID
type _ = resourcev1.ResourceClaimStatus

// DevicesDie mutates a single item in Devices matched by the nested fields Driver, Device, Pool, and ShareID, appending a new item if no match is found.
//
// Devices contains the status of each device allocated for this
// claim, as reported by the driver. This can include driver-specific
// information. Entries are owned by their respective drivers.
//
// +optional
// +listType=map
// +listMapKey=driver
// +listMapKey=device
// +listMapKey=pool
// +listMapKey=shareID
// +featureGate=DRAResourceClaimDeviceStatus
func (d *ResourceClaimStatusDie) DevicesDie(driver string, device string, pool string, shareID *string, fn func(d *AllocatedDeviceStatusDie)) *ResourceClaimStatusDie {
	return d.DieStamp(func(r *resourcev1.ResourceClaimStatus) {
		for i := range r.Devices {
			if driver == r.Devices[i].Driver && device == r.Devices[i].Device && pool == r.Devices[i].Pool && shareID == r.Devices[i].ShareID {
				d := AllocatedDeviceStatusBlank.DieImmutable(false).DieFeed(r.Devices[i])
				fn(d)
				r.Devices[i] = d.DieRelease()
				return
			}
		}

		d := AllocatedDeviceStatusBlank.DieImmutable(false).DieFeed(resourcev1.AllocatedDeviceStatus{
			Driver:  driver,
			Device:  device,
			Pool:    pool,
			ShareID: shareID,
		})
		fn(d)
		r.Devices = append(r.Devices, d.DieRelease())
	})
}

// +die
// +die:field:name=Conditions,package=_/meta/v1,die=ConditionDie,listType=map,listMapKey=Type
// +die:field:name=NetworkData,die=NetworkDeviceDataDie,pointer=true
type _ = resourcev1.AllocatedDeviceStatus

// +die
type _ = resourcev1.NetworkDeviceData

// +die
// +die:field:name=Devices,die=DeviceAllocationResultDie
// +die:field:name=NodeSelector,package=_/core/v1,die=NodeSelectorDie,pointer=true
type _ = resourcev1.AllocationResult

// +die
// +die:field:name=Results,die=DeviceRequestAllocationResultDie,listType=atomic
// +die:field:name=Config,die=DeviceAllocationConfigurationDie,listType=atomic
type _ = resourcev1.DeviceAllocationResult

// +die:ignore=ConsumedCapacity
// +die:field:name=Tolerations,die=DeviceTolerationDie,listType=atomic
type _ = resourcev1.DeviceRequestAllocationResult

// ConsumedCapacity tracks the amount of capacity consumed per device as part of the claim request.
// The consumed amount may differ from the requested amount: it is rounded up to the nearest valid
// value based on the device’s requestPolicy if applicable (i.e., may not be less than the requested amount).
//
// The total consumed capacity for each device must not exceed the DeviceCapacity's Value.
//
// This field is populated only for devices that allow multiple allocations.
// All capacity entries are included, even if the consumed amount is zero.
//
// +optional
// +featureGate=DRAConsumableCapacity
func (d *DeviceRequestAllocationResultDie) ConsumedCapacity(v map[resourcev1.QualifiedName]resource.Quantity) *DeviceRequestAllocationResultDie {
	return d.DieStamp(func(r *resourcev1.DeviceRequestAllocationResult) {
		r.ConsumedCapacity = v
	})
}

// AddRequest sets a single quantity on the Requests map.
//
// ConsumedCapacity tracks the amount of capacity consumed per device as part of the claim request.
// The consumed amount may differ from the requested amount: it is rounded up to the nearest valid
// value based on the device’s requestPolicy if applicable (i.e., may not be less than the requested amount).
//
// The total consumed capacity for each device must not exceed the DeviceCapacity's Value.
//
// This field is populated only for devices that allow multiple allocations.
// All capacity entries are included, even if the consumed amount is zero.
//
// +optional
// +featureGate=DRAConsumableCapacity
func (d *DeviceRequestAllocationResultDie) AddConsumedCapacity(name resourcev1.QualifiedName, quantity resource.Quantity) *DeviceRequestAllocationResultDie {
	return d.DieStamp(func(r *resourcev1.DeviceRequestAllocationResult) {
		if r.ConsumedCapacity == nil {
			r.ConsumedCapacity = map[resourcev1.QualifiedName]resource.Quantity{}
		}
		r.ConsumedCapacity[name] = quantity
	})
}

// AddConsumedCapacityString parses the quantity setting a single value on the ConsumedCapacity map. Panics if the string is not parsable.
//
// ConsumedCapacity tracks the amount of capacity consumed per device as part of the claim request.
// The consumed amount may differ from the requested amount: it is rounded up to the nearest valid
// value based on the device’s requestPolicy if applicable (i.e., may not be less than the requested amount).
//
// The total consumed capacity for each device must not exceed the DeviceCapacity's Value.
//
// This field is populated only for devices that allow multiple allocations.
// All capacity entries are included, even if the consumed amount is zero.
//
// +optional
// +featureGate=DRAConsumableCapacity
func (d *DeviceRequestAllocationResultDie) AddConsumedCapacityString(name resourcev1.QualifiedName, quantity string) *DeviceRequestAllocationResultDie {
	q := resource.MustParse(quantity)
	return d.AddConsumedCapacity(name, q)
}

// +die
// +die:field:name=DeviceConfiguration,die=DeviceConfigurationDie
type _ = resourcev1.DeviceAllocationConfiguration

// +die
type _ = resourcev1.ResourceClaimConsumerReference

// +die:object=true
type _ = resourcev1.ResourceClaimTemplate

// +die
// +die:field:name=Spec,die=ResourceClaimSpecDie
type _ = resourcev1.ResourceClaimTemplateSpec

func (d *ResourceClaimTemplateSpecDie) AddLabel(key, value string) *ResourceClaimTemplateSpecDie {
	return d.DieStamp(func(r *resourcev1.ResourceClaimTemplateSpec) {
		if r.Labels == nil {
			r.Labels = map[string]string{}
		}
		r.Labels[key] = value
	})
}

func (d *ResourceClaimTemplateSpecDie) AddAnnotation(key, value string) *ResourceClaimTemplateSpecDie {
	return d.DieStamp(func(r *resourcev1.ResourceClaimTemplateSpec) {
		if r.Annotations == nil {
			r.Annotations = map[string]string{}
		}
		r.Annotations[key] = value
	})
}

// Map of string keys and values that can be used to organize and categorize
//
// (scope and select) objects. May match selectors of replication controllers
//
// and services.
//
// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels
func (d *ResourceClaimTemplateSpecDie) Labels(v map[string]string) *ResourceClaimTemplateSpecDie {
	return d.DieStamp(func(r *resourcev1.ResourceClaimTemplateSpec) {
		r.Labels = v
	})
}

// Annotations is an unstructured key value map stored with a resource that may be
//
// set by external tools to store and retrieve arbitrary metadata. They are not
//
// queryable and should be preserved when modifying objects.
//
// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations
func (d *ResourceClaimTemplateSpecDie) Annotations(v map[string]string) *ResourceClaimTemplateSpecDie {
	return d.DieStamp(func(r *resourcev1.ResourceClaimTemplateSpec) {
		r.Annotations = v
	})
}
