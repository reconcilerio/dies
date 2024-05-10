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
	storagev1 "k8s.io/api/storage/v1"
	diecorev1 "reconciler.io/dies/apis/core/v1"
)

// +die:object=true,apiVersion=storage.k8s.io/v1,kind=StorageClass
type _ = storagev1.StorageClass

func (d *StorageClassDie) AddParameter(key, value string) *StorageClassDie {
	return d.DieStamp(func(r *storagev1.StorageClass) {
		if r.Parameters == nil {
			r.Parameters = map[string]string{}
		}
		r.Parameters[key] = value
	})
}

func (d *StorageClassDie) AllowedTopologiesDie(topologies ...*diecorev1.TopologySelectorTermDie) *StorageClassDie {
	return d.DieStamp(func(r *storagev1.StorageClass) {
		r.AllowedTopologies = make([]corev1.TopologySelectorTerm, len(topologies))
		for i := range topologies {
			r.AllowedTopologies[i] = topologies[i].DieRelease()
		}
	})
}
