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

// +die:target=k8s.io/api/core/v1.Pod,object=true

// +die:target=k8s.io/api/core/v1.PodSpec
// +die:field:receiver=PodSpecDie,name=Volumes,type=[]k8s.io/api/core/v1.Volume
// +die:field:receiver=PodSpecDie,name=InitContainers,type=[]k8s.io/api/core/v1.Container
// +die:field:receiver=PodSpecDie,name=Containers,type=[]k8s.io/api/core/v1.Container
// +die:field:receiver=PodSpecDie,name=EphemeralContainers,type=[]k8s.io/api/core/v1.EphemeralContainer
// +die:field:receiver=PodSpecDie,name=RestartPolicy,type=k8s.io/api/core/v1.RestartPolicy
// +die:field:receiver=PodSpecDie,name=TerminationGracePeriodSeconds,type=*int64
// +die:field:receiver=PodSpecDie,name=ActiveDeadlineSeconds,type=*int64
// +die:field:receiver=PodSpecDie,name=DNSPolicy,type=k8s.io/api/core/v1.DNSPolicy
// +die:field:receiver=PodSpecDie,name=NodeSelector,type=map[string]string
// +die:field:receiver=PodSpecDie,name=ServiceAccountName,type=string
// +die:field:receiver=PodSpecDie,name=AutomountServiceAccountToken,type=*bool
// +die:field:receiver=PodSpecDie,name=NodeName,type=string
// +die:field:receiver=PodSpecDie,name=HostNetwork,type=bool
// +die:field:receiver=PodSpecDie,name=HostPID,type=bool
// +die:field:receiver=PodSpecDie,name=HostIPC,type=bool
// +die:field:receiver=PodSpecDie,name=ShareProcessNamespace,type=*bool
// +die:field:receiver=PodSpecDie,name=SecurityContext,type=*k8s.io/api/core/v1.PodSecurityContext
// +die:field:receiver=PodSpecDie,name=ImagePullSecrets,type=[]k8s.io/api/core/v1.LocalObjectReference
// +die:field:receiver=PodSpecDie,name=Hostname,type=string
// +die:field:receiver=PodSpecDie,name=Subdomain,type=string
// +die:field:receiver=PodSpecDie,name=Affinity,type=*k8s.io/api/core/v1.Affinity
// +die:field:receiver=PodSpecDie,name=SchedulerName,type=string
// +die:field:receiver=PodSpecDie,name=Tolerations,type=[]k8s.io/api/core/v1.Toleration
// +die:field:receiver=PodSpecDie,name=HostAliases,type=[]k8s.io/api/core/v1.HostAlias
// +die:field:receiver=PodSpecDie,name=PriorityClassName,type=string
// +die:field:receiver=PodSpecDie,name=Priority,type=*int32
// +die:field:receiver=PodSpecDie,name=DNSConfig,type=*k8s.io/api/core/v1.PodDNSConfig
// +die:field:receiver=PodSpecDie,name=ReadinessGates,type=[]k8s.io/api/core/v1.PodReadinessGate
// +die:field:receiver=PodSpecDie,name=RuntimeClassName,type=*string
// +die:field:receiver=PodSpecDie,name=EnableServiceLinks,type=*bool
// +die:field:receiver=PodSpecDie,name=PreemptionPolicy,type=*k8s.io/api/core/v1.PreemptionPolicy
// +die:field:receiver=PodSpecDie,name=Overhead,type=k8s.io/api/core/v1.ResourceList
// +die:field:receiver=PodSpecDie,name=TopologySpreadConstraints,type=[]k8s.io/api/core/v1.TopologySpreadConstraint
// +die:field:receiver=PodSpecDie,name=SetHostnameAsFQDN,type=*bool

func (d *PodSpecDie) InitContainerDie(name string, fn func(d *ContainerDie)) *PodSpecDie {
	return d.DieStamp(func(r *corev1.PodSpec) {
		for i := range r.InitContainers {
			if name == r.InitContainers[i].Name {
				d := ContainerBlank.DieImmutable(false).DieFeed(r.InitContainers[i])
				fn(d)
				r.InitContainers[i] = d.DieRelease()
				return
			}
		}

		d := ContainerBlank.DieImmutable(false).DieFeed(corev1.Container{Name: name})
		fn(d)
		r.InitContainers = append(r.InitContainers, d.DieRelease())
	})
}

func (d *PodSpecDie) ContainerDie(name string, fn func(d *ContainerDie)) *PodSpecDie {
	return d.DieStamp(func(r *corev1.PodSpec) {
		for i := range r.Containers {
			if name == r.Containers[i].Name {
				d := ContainerBlank.DieImmutable(false).DieFeed(r.Containers[i])
				fn(d)
				r.Containers[i] = d.DieRelease()
				return
			}
		}

		d := ContainerBlank.DieImmutable(false).DieFeed(corev1.Container{Name: name})
		fn(d)
		r.Containers = append(r.Containers, d.DieRelease())
	})
}

func (d *PodSpecDie) AddVolumes(volumes ...corev1.Volume) *PodSpecDie {
	return d.DieStamp(func(r *corev1.PodSpec) {
		for _, m := range volumes {
			found := false
			for i := range r.Volumes {
				if m.Name == r.Volumes[i].Name {
					found = true
					r.Volumes[i] = m
				}
			}
			if !found {
				r.Volumes = append(r.Volumes, m)
			}
		}
	})
}

// +die:target=k8s.io/api/core/v1.PodStatus
// +die:field:receiver=PodStatusDie,name=Phase,type=k8s.io/api/core/v1.PodPhase
// +die:field:receiver=PodStatusDie,name=Conditions,type=[]k8s.io/api/core/v1.PodCondition
// +die:field:receiver=PodStatusDie,name=Message,type=string
// +die:field:receiver=PodStatusDie,name=Reason,type=string
// +die:field:receiver=PodStatusDie,name=NominatedNodeName,type=string
// +die:field:receiver=PodStatusDie,name=HostIP,type=string
// +die:field:receiver=PodStatusDie,name=PodIP,type=string
// +die:field:receiver=PodStatusDie,name=PodIPs,type=[]k8s.io/api/core/v1.PodIP
// +die:field:receiver=PodStatusDie,name=StartTime,type=*k8s.io/apimachinery/pkg/apis/meta/v1.Time
// +die:field:receiver=PodStatusDie,name=InitContainerStatuses,type=[]k8s.io/api/core/v1.ContainerStatus
// +die:field:receiver=PodStatusDie,name=ContainerStatuses,type=[]k8s.io/api/core/v1.ContainerStatus
// +die:field:receiver=PodStatusDie,name=QOSClass,type=k8s.io/api/core/v1.PodQOSClass
// +die:field:receiver=PodStatusDie,name=EphemeralContainerStatuses,type=[]k8s.io/api/core/v1.ContainerStatus

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
