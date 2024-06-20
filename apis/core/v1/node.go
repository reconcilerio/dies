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

import (
	corev1 "k8s.io/api/core/v1"
	diemetav1 "reconciler.io/dies/apis/meta/v1"
)

// +die:object=true,apiVersion=v1,kind=Node
type _ = corev1.Node

// +die
// +die:field:name=ConfigSource,die=NodeConfigSourceDie,pointer=true
// +die:field:name=Taints,die=TaintDie,listType=map,listMapKey=Key
type _ = corev1.NodeSpec

// +die
type _ = corev1.Taint

// +die
// +die:field:name=ConfigMap,die=ConfigMapNodeConfigSourceDie,pointer=true
type _ = corev1.NodeConfigSource

// +die
type _ = corev1.ConfigMapNodeConfigSource

// +die
// +die:field:name=DaemonEndpoints,die=NodeDaemonEndpointsDie
// +die:field:name=NodeInfo,die=NodeSystemInfoDie
// +die:field:name=Config,die=NodeConfigStatusDie,pointer=true
// +die:field:name=Addresses,die=NodeAddressDie,listType=atomic
// +die:field:name=Images,die=ContainerImageDie,listType=atomic
// +die:field:name=VolumesAttached,method=VolumeAttachedDie,die=AttachedVolumeDie,listType=map,listMapKeyPackage=k8s.io/api/core/v1,listMapKeyType=UniqueVolumeName
// +die:field:name=RuntimeHandlers,die=NodeRuntimeHandlerDie,listType=atomic
type _ = corev1.NodeStatus

func (d *NodeStatusDie) ConditionsDie(conditions ...*diemetav1.ConditionDie) *NodeStatusDie {
	return d.DieStamp(func(r *corev1.NodeStatus) {
		r.Conditions = make([]corev1.NodeCondition, len(conditions))
		for i := range conditions {
			c := conditions[i].DieRelease()
			r.Conditions[i] = corev1.NodeCondition{
				Type:               corev1.NodeConditionType(c.Type),
				Status:             corev1.ConditionStatus(c.Status),
				Reason:             c.Reason,
				Message:            c.Message,
				LastTransitionTime: c.LastTransitionTime,
			}
		}
	})
}

// deprecated: use AddressesDie
func (d *NodeStatusDie) AddresssDie(addresses ...*NodeAddressDie) *NodeStatusDie {
	return d.AddressesDie(addresses...)
}

// +die
type _ = corev1.NodeAddress

// +die
// +die:field:name=KubeletEndpoint,die=DaemonEndpointDie
type _ = corev1.NodeDaemonEndpoints

// +die
type _ = corev1.DaemonEndpoint

// +die
type _ = corev1.NodeSystemInfo

// +die
type _ = corev1.ContainerImage

// +die
type _ = corev1.AttachedVolume

// +die
// +die:field:name=Assigned,die=NodeConfigSourceDie,pointer=true
// +die:field:name=Active,die=NodeConfigSourceDie,pointer=true
// +die:field:name=LastKnownGood,die=NodeConfigSourceDie,pointer=true
type _ = corev1.NodeConfigStatus

// +die
// +die:field:name=Features,die=NodeRuntimeHandlerFeaturesDie,pointer=true
type _ = corev1.NodeRuntimeHandler

// +die
type _ = corev1.NodeRuntimeHandlerFeatures
