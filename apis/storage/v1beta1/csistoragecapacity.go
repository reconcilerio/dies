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

package v1beta1

import (
	diemetav1 "dies.dev/apis/meta/v1"
	storagev1beta1 "k8s.io/api/storage/v1beta1"
	resource "k8s.io/apimachinery/pkg/api/resource"
)

// +die:object=true
type _ = storagev1beta1.CSIStorageCapacity

func (d *CSIStorageCapacityDie) NodeTopologyDie(fn func(d *diemetav1.LabelSelectorDie)) *CSIStorageCapacityDie {
	return d.DieStamp(func(r *storagev1beta1.CSIStorageCapacity) {
		d := diemetav1.LabelSelectorBlank.DieImmutable(false).DieFeedPtr(r.NodeTopology)
		fn(d)
		r.NodeTopology = d.DieReleasePtr()
	})
}

func (d *CSIStorageCapacityDie) CapacityString(quantity string) *CSIStorageCapacityDie {
	q := resource.MustParse(quantity)
	return d.Capacity(&q)
}

func (d *CSIStorageCapacityDie) MaximumVolumeSizeString(quantity string) *CSIStorageCapacityDie {
	q := resource.MustParse(quantity)
	return d.MaximumVolumeSize(&q)
}
