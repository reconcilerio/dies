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
	storagev1 "k8s.io/api/storage/v1"
)

// +die:object=true,apiVersion=storage.k8s.io/v1,kind=CSINode
type _ = storagev1.CSINode

// +die
type _ = storagev1.CSINodeSpec

func (d *CSINodeSpecDie) DriversDie(drivers ...*CSINodeDriverDie) *CSINodeSpecDie {
	return d.DieStamp(func(r *storagev1.CSINodeSpec) {
		r.Drivers = make([]storagev1.CSINodeDriver, len(drivers))
		for i := range drivers {
			r.Drivers[i] = drivers[i].DieRelease()
		}
	})
}

// +die
type _ = storagev1.CSINodeDriver

func (d *CSINodeDriverDie) AllocatableDie(fn func(d *VolumeNodeResourcesDie)) *CSINodeDriverDie {
	return d.DieStamp(func(r *storagev1.CSINodeDriver) {
		d := VolumeNodeResourcesBlank.DieImmutable(false).DieFeedPtr(r.Allocatable)
		fn(d)
		r.Allocatable = d.DieReleasePtr()
	})
}

// +die
type _ = storagev1.VolumeNodeResources
