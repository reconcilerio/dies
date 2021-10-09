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
	diemetav1 "github.com/scothis/dies/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
)

// +die:target=k8s.io/api/core/v1.Node,object=true

// +die:target=k8s.io/api/core/v1.NodeSpec
// +die:field:receiver=NodeSpecDie,name=PodCIDR,type=string
// +die:field:receiver=NodeSpecDie,name=PodCIDRs,type=[]string
// +die:field:receiver=NodeSpecDie,name=ProviderID,type=string
// +die:field:receiver=NodeSpecDie,name=Unschedulable,type=bool
// +die:field:receiver=NodeSpecDie,name=Taints,type=[]k8s.io/api/core/v1.Taint
// +die:field:receiver=NodeSpecDie,name=ConfigSource,type=*k8s.io/api/core/v1.NodeConfigSource
// +die:field:receiver=NodeSpecDie,name=DoNotUseExternalID,type=string

// +die:target=k8s.io/api/core/v1.NodeStatus
// +die:field:receiver=NodeStatusDie,name=Capacity,type=k8s.io/api/core/v1.ResourceList
// +die:field:receiver=NodeStatusDie,name=Allocatable,type=k8s.io/api/core/v1.ResourceList
// +die:field:receiver=NodeStatusDie,name=Phase,type=k8s.io/api/core/v1.NodePhase
// +die:field:receiver=NodeStatusDie,name=Conditions,type=[]k8s.io/api/core/v1.NodeCondition
// +die:field:receiver=NodeStatusDie,name=Addresses,type=[]k8s.io/api/core/v1.NodeAddress
// +die:field:receiver=NodeStatusDie,name=DaemonEndpoints,type=k8s.io/api/core/v1.NodeDaemonEndpoints
// +die:field:receiver=NodeStatusDie,name=NodeInfo,type=k8s.io/api/core/v1.NodeSystemInfo
// +die:field:receiver=NodeStatusDie,name=Images,type=[]k8s.io/api/core/v1.ContainerImage
// +die:field:receiver=NodeStatusDie,name=VolumesInUse,type=[]k8s.io/api/core/v1.UniqueVolumeName
// +die:field:receiver=NodeStatusDie,name=VolumesAttached,type=[]k8s.io/api/core/v1.AttachedVolume
// +die:field:receiver=NodeStatusDie,name=Config,type=*k8s.io/api/core/v1.NodeConfigStatus

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
