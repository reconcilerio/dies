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
	"k8s.io/apimachinery/pkg/api/resource"
)

// +die:object=true
type _ = corev1.LimitRange

// +die
type _ = corev1.LimitRangeSpec

type limitRangeSpec interface {
	LimitsDie(limits ...LimitRangeItemDie) LimitRangeSpecDie
}

func (d *limitRangeSpecDie) LimitsDie(limits ...LimitRangeItemDie) LimitRangeSpecDie {
	return d.DieStamp(func(r *corev1.LimitRangeSpec) {
		r.Limits = make([]corev1.LimitRangeItem, len(limits))
		for i := range r.Limits {
			r.Limits[i] = limits[i].DieRelease()
		}
	})
}

// +die
type _ = corev1.LimitRangeItem

type limitRangeItem interface {
	AddMax(name corev1.ResourceName, quantity resource.Quantity) LimitRangeItemDie
	AddMaxString(name corev1.ResourceName, quantity string) LimitRangeItemDie
	AddMin(name corev1.ResourceName, quantity resource.Quantity) LimitRangeItemDie
	AddMinString(name corev1.ResourceName, quantity string) LimitRangeItemDie
	AddDefault(name corev1.ResourceName, quantity resource.Quantity) LimitRangeItemDie
	AddDefaultString(name corev1.ResourceName, quantity string) LimitRangeItemDie
	AddDefaultRequest(name corev1.ResourceName, quantity resource.Quantity) LimitRangeItemDie
	AddDefaultRequestString(name corev1.ResourceName, quantity string) LimitRangeItemDie
	AddMaxLimitRequestRatio(name corev1.ResourceName, quantity resource.Quantity) LimitRangeItemDie
	AddMaxLimitRequestRatioString(name corev1.ResourceName, quantity string) LimitRangeItemDie
}

func (d *limitRangeItemDie) AddMax(name corev1.ResourceName, quantity resource.Quantity) LimitRangeItemDie {
	return d.DieStamp(func(r *corev1.LimitRangeItem) {
		r.Max[name] = quantity
	})
}

func (d *limitRangeItemDie) AddMaxString(name corev1.ResourceName, quantity string) LimitRangeItemDie {
	return d.AddMax(name, resource.MustParse(quantity))
}

func (d *limitRangeItemDie) AddMin(name corev1.ResourceName, quantity resource.Quantity) LimitRangeItemDie {
	return d.DieStamp(func(r *corev1.LimitRangeItem) {
		r.Min[name] = quantity
	})
}

func (d *limitRangeItemDie) AddMinString(name corev1.ResourceName, quantity string) LimitRangeItemDie {
	return d.AddMin(name, resource.MustParse(quantity))
}

func (d *limitRangeItemDie) AddDefault(name corev1.ResourceName, quantity resource.Quantity) LimitRangeItemDie {
	return d.DieStamp(func(r *corev1.LimitRangeItem) {
		r.Default[name] = quantity
	})
}

func (d *limitRangeItemDie) AddDefaultString(name corev1.ResourceName, quantity string) LimitRangeItemDie {
	return d.AddDefault(name, resource.MustParse(quantity))
}

func (d *limitRangeItemDie) AddDefaultRequest(name corev1.ResourceName, quantity resource.Quantity) LimitRangeItemDie {
	return d.DieStamp(func(r *corev1.LimitRangeItem) {
		r.DefaultRequest[name] = quantity
	})
}

func (d *limitRangeItemDie) AddDefaultRequestString(name corev1.ResourceName, quantity string) LimitRangeItemDie {
	return d.AddDefaultRequest(name, resource.MustParse(quantity))
}

func (d *limitRangeItemDie) AddMaxLimitRequestRatio(name corev1.ResourceName, quantity resource.Quantity) LimitRangeItemDie {
	return d.DieStamp(func(r *corev1.LimitRangeItem) {
		r.MaxLimitRequestRatio[name] = quantity
	})
}

func (d *limitRangeItemDie) AddMaxLimitRequestRatioString(name corev1.ResourceName, quantity string) LimitRangeItemDie {
	return d.AddMaxLimitRequestRatio(name, resource.MustParse(quantity))
}
