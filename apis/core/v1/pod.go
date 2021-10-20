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
	diemetav1 "dies.dev/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
)

// +die:object=true
type _ = corev1.Pod

// +die
type _ = corev1.PodSpec

type podSpec interface {
	InitContainerDie(name string, fn func(d ContainerDie)) PodSpecDie
	ContainerDie(name string, fn func(d ContainerDie)) PodSpecDie
	VolumeDie(name string, fn func(d VolumeDie)) PodSpecDie
}

func (d *podSpecDie) InitContainerDie(name string, fn func(d ContainerDie)) PodSpecDie {
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

func (d *podSpecDie) ContainerDie(name string, fn func(d ContainerDie)) PodSpecDie {
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

func (d *podSpecDie) VolumeDie(name string, fn func(d VolumeDie)) PodSpecDie {
	return d.DieStamp(func(r *corev1.PodSpec) {
		for i := range r.Volumes {
			if name == r.Volumes[i].Name {
				d := VolumeBlank.DieImmutable(false).DieFeed(r.Volumes[i])
				fn(d)
				r.Volumes[i] = d.DieRelease()
				return
			}
		}

		d := VolumeBlank.DieImmutable(false).DieFeed(corev1.Volume{Name: name})
		fn(d)
		r.Volumes = append(r.Volumes, d.DieRelease())
	})
}

// +die
type _ = corev1.PodStatus

type podStatus interface {
	ConditionsDie(conditions ...diemetav1.ConditionDie) PodStatusDie
}

func (d *podStatusDie) ConditionsDie(conditions ...diemetav1.ConditionDie) PodStatusDie {
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
