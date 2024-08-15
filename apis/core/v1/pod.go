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

// +die:object=true,apiVersion=v1,kind=Pod
type _ = corev1.Pod

// +die
// +die:field:name=SecurityContext,die=PodSecurityContextDie,pointer=true
// +die:field:name=DNSConfig,die=PodDNSConfigDie,pointer=true
// +die:field:name=OS,die=PodOSDie,pointer=true
// +die:field:name=Volumes,die=VolumeDie,listType=map
// +die:field:name=InitContainers,die=ContainerDie,listType=map
// +die:field:name=Containers,die=ContainerDie,listType=map
// +die:field:name=Tolerations,die=TolerationDie,listMapKey=Key
// +die:field:name=HostAliases,die=HostAliasDie,listType=atomic
// +die:field:name=ReadinessGates,die=PodReadinessGateDie,listType=atomic
// +die:field:name=TopologySpreadConstraints,die=TopologySpreadConstraintDie,listMapKey=TopologyKey
// +die:field:name=SchedulingGates,die=PodSchedulingGateDie,listType=atomic
// +die:field:name=ResourceClaims,die=PodResourceClaimDie,listType=atomic
type _ = corev1.PodSpec

// +die
type _ = corev1.PodSchedulingGate

// +die
type _ = corev1.PodResourceClaim

// +die
// +die:field:name=SELinuxOptions,die=SELinuxOptionsDie,pointer=true
// +die:field:name=WindowsOptions,die=WindowsSecurityContextOptionsDie,pointer=true
// +die:field:name=AppArmorProfile,die=AppArmorProfileDie,pointer=true
// +die:field:name=Sysctls,die=SysctlDie,listType=atomic
// +die:field:name=SeccompProfile,die=SeccompProfileDie,pointer=true

type _ = corev1.PodSecurityContext

// +die
type _ = corev1.Sysctl

// +die
type _ = corev1.Toleration

// +die
type _ = corev1.HostAlias

// +die
// +die:field:name=Options,die=PodDNSConfigOptionDie,listType=atomic
type _ = corev1.PodDNSConfig

// +die
type _ = corev1.PodDNSConfigOption

// +die
type _ = corev1.PodReadinessGate

// +die
// +die:field:name=LabelSelector,package=_/meta/v1,die=LabelSelectorDie,pointer=true
type _ = corev1.TopologySpreadConstraint

// +die
type _ = corev1.PodOS

// +die
// +die:field:name=InitContainerStatuses,method=InitContainerStatusDie,die=ContainerStatusDie,listType=map
// +die:field:name=ContainerStatuses,method=ContainerStatusDie,die=ContainerStatusDie,listType=map
// +die:field:name=EphemeralContainerStatuses,method=EphemeralContainerStatusDie,die=ContainerStatusDie,listType=map
type _ = corev1.PodStatus

func (d *PodStatusDie) ConditionsDie(conditions ...*diemetav1.ConditionDie) *PodStatusDie {
	return d.DieStamp(func(r *corev1.PodStatus) {
		r.Conditions = make([]corev1.PodCondition, len(conditions))
		for i := range conditions {
			c := conditions[i].DieRelease()
			r.Conditions[i] = corev1.PodCondition{
				Type:               corev1.PodConditionType(c.Type),
				Status:             corev1.ConditionStatus(c.Status),
				Reason:             c.Reason,
				Message:            c.Message,
				LastTransitionTime: c.LastTransitionTime,
			}
		}
	})
}
