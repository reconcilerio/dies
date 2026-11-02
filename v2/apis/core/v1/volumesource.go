/*
Copyright 2026 the original author or authors.

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
)

// VolumeSource's members form a union, which upstream only documents: "Only
// one of its members may be specified." This adds that to the generated
// alias's metadata.
//
// +die:union:name=source,members={hostPath,emptyDir,gcePersistentDisk,awsElasticBlockStore,gitRepo,secret,nfs,iscsi,glusterfs,persistentVolumeClaim,rbd,flexVolume,cinder,cephfs,flocker,downwardAPI,fc,azureFile,configMap,vsphereVolume,quobyte,azureDisk,photonPersistentDisk,projected,portworxVolume,scaleIO,storageos,csi,ephemeral,image}
type _ = corev1.VolumeSource
