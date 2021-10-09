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

// +die:target=k8s.io/api/core/v1.PersistentVolume,object=true

// +die:target=k8s.io/api/core/v1.PersistentVolumeSpec
// +die:field:receiver=PersistentVolumeSpecDie,name=Capacity,type=k8s.io/api/core/v1.ResourceList
// +die:field:receiver=PersistentVolumeSpecDie,name=PersistentVolumeSource,type=k8s.io/api/core/v1.PersistentVolumeSource
// +die:field:receiver=PersistentVolumeSpecDie,name=AccessModes,type=[]k8s.io/api/core/v1.PersistentVolumeAccessMode
// +die:field:receiver=PersistentVolumeSpecDie,name=ClaimRef,type=*k8s.io/api/core/v1.ObjectReference
// +die:field:receiver=PersistentVolumeSpecDie,name=PersistentVolumeReclaimPolicy,type=k8s.io/api/core/v1.PersistentVolumeReclaimPolicy
// +die:field:receiver=PersistentVolumeSpecDie,name=StorageClassName,type=string
// +die:field:receiver=PersistentVolumeSpecDie,name=MountOptions,type=[]string
// +die:field:receiver=PersistentVolumeSpecDie,name=VolumeMode,type=*k8s.io/api/core/v1.PersistentVolumeMode
// +die:field:receiver=PersistentVolumeSpecDie,name=NodeAffinity,type=*k8s.io/api/core/v1.VolumeNodeAffinity

// +die:target=k8s.io/api/core/v1.PersistentVolumeStatus
// +die:field:receiver=PersistentVolumeStatusDie,name=Phase,type=k8s.io/api/core/v1.PersistentVolumePhase
// +die:field:receiver=PersistentVolumeStatusDie,name=Message,type=string
// +die:field:receiver=PersistentVolumeStatusDie,name=Reason,type=string
