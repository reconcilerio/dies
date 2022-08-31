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
	diecorev1 "dies.dev/apis/core/v1"
	corev1 "k8s.io/api/core/v1"
	nodev1 "k8s.io/api/node/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// +die:object=true
type _ = nodev1.RuntimeClass

// +die
type _ = nodev1.Overhead

func (d *OverheadDie) AddPodFixed(name corev1.ResourceName, quantity resource.Quantity) *OverheadDie {
	return d.DieStamp(func(r *nodev1.Overhead) {
		if r.PodFixed == nil {
			r.PodFixed = corev1.ResourceList{}
		}
		r.PodFixed[name] = quantity
	})
}

func (d *OverheadDie) AddPodFixedString(name corev1.ResourceName, quantity string) *OverheadDie {
	return d.AddPodFixed(name, resource.MustParse(quantity))
}

// +die
type _ = nodev1.Scheduling

func (d *SchedulingDie) TolerationsDie(key string, fn func(d *diecorev1.TolerationDie)) *SchedulingDie {
	return d.DieStamp(func(r *nodev1.Scheduling) {
		for i := range r.Tolerations {
			if key == r.Tolerations[i].Key {
				d := diecorev1.TolerationBlank.DieImmutable(false).DieFeed(r.Tolerations[i])
				fn(d)
				r.Tolerations[i] = d.DieRelease()
				return
			}
		}

		d := diecorev1.TolerationBlank.DieImmutable(false).DieFeed(corev1.Toleration{Key: key})
		fn(d)
		r.Tolerations = append(r.Tolerations, d.DieRelease())
	})
}
