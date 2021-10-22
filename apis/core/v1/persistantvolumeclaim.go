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
	"k8s.io/apimachinery/pkg/api/resource"
)

// +die:object=true
type _ = corev1.PersistentVolumeClaim

// +die
type _ = corev1.PersistentVolumeClaimSpec

type persistentVolumeClaimSpecDieExtension interface {
	SelectorDie(fn func(d diemetav1.LabelSelectorDie)) PersistentVolumeClaimSpecDie
	ResourcesDie(fn func(d ResourceRequirementsDie)) PersistentVolumeClaimSpecDie
	DataSourceDie(fn func(d TypedLocalObjectReferenceDie)) PersistentVolumeClaimSpecDie
	DataSourceRefDie(fn func(d TypedLocalObjectReferenceDie)) PersistentVolumeClaimSpecDie
}

func (d *persistentVolumeClaimSpecDie) SelectorDie(fn func(d diemetav1.LabelSelectorDie)) PersistentVolumeClaimSpecDie {
	return d.DieStamp(func(r *corev1.PersistentVolumeClaimSpec) {
		d := diemetav1.LabelSelectorBlank.DieImmutable(false).DieFeedPtr(r.Selector)
		fn(d)
		r.Selector = d.DieReleasePtr()
	})
}

func (d *persistentVolumeClaimSpecDie) ResourcesDie(fn func(d ResourceRequirementsDie)) PersistentVolumeClaimSpecDie {
	return d.DieStamp(func(r *corev1.PersistentVolumeClaimSpec) {
		d := ResourceRequirementsBlank.DieImmutable(false).DieFeed(r.Resources)
		fn(d)
		r.Resources = d.DieRelease()
	})
}

func (d *persistentVolumeClaimSpecDie) DataSourceDie(fn func(d TypedLocalObjectReferenceDie)) PersistentVolumeClaimSpecDie {
	return d.DieStamp(func(r *corev1.PersistentVolumeClaimSpec) {
		d := TypedLocalObjectReferenceBlank.DieImmutable(false).DieFeedPtr(r.DataSource)
		fn(d)
		r.DataSource = d.DieReleasePtr()
	})
}

func (d *persistentVolumeClaimSpecDie) DataSourceRefDie(fn func(d TypedLocalObjectReferenceDie)) PersistentVolumeClaimSpecDie {
	return d.DieStamp(func(r *corev1.PersistentVolumeClaimSpec) {
		d := TypedLocalObjectReferenceBlank.DieImmutable(false).DieFeedPtr(r.DataSourceRef)
		fn(d)
		r.DataSourceRef = d.DieReleasePtr()
	})
}

// +die
type _ = corev1.PersistentVolumeClaimStatus

type persistentVolumeClaimStatusDieExtension interface {
	AddCapacity(name corev1.ResourceName, quantity resource.Quantity) PersistentVolumeClaimStatusDie
	AddCapacityString(name corev1.ResourceName, quantity string) PersistentVolumeClaimStatusDie
	ConditionsDie(conditions ...diemetav1.ConditionDie) PersistentVolumeClaimStatusDie
}

func (d *persistentVolumeClaimStatusDie) AddCapacity(name corev1.ResourceName, quantity resource.Quantity) PersistentVolumeClaimStatusDie {
	return d.DieStamp(func(r *corev1.PersistentVolumeClaimStatus) {
		r.Capacity[name] = quantity
	})
}

func (d *persistentVolumeClaimStatusDie) AddCapacityString(name corev1.ResourceName, quantity string) PersistentVolumeClaimStatusDie {
	return d.AddCapacity(name, resource.MustParse(quantity))
}

func (d *persistentVolumeClaimStatusDie) ConditionsDie(conditions ...diemetav1.ConditionDie) PersistentVolumeClaimStatusDie {
	return d.DieStamp(func(r *corev1.PersistentVolumeClaimStatus) {
		r.Conditions = make([]corev1.PersistentVolumeClaimCondition, len(conditions))
		for i := range conditions {
			c := conditions[i].DieRelease()
			r.Conditions[i] = corev1.PersistentVolumeClaimCondition{
				Type:               corev1.PersistentVolumeClaimConditionType(c.Type),
				Status:             corev1.ConditionStatus(c.Status),
				Reason:             c.Reason,
				Message:            c.Message,
				LastTransitionTime: c.LastTransitionTime,
			}
		}
	})
}

// +die
type _ corev1.PersistentVolumeClaimTemplate

type persistentVolumeClaimTemplateDieExtension interface {
	MetadataDie(fn func(d diemetav1.ObjectMetaDie)) PersistentVolumeClaimTemplateDie
	SpecDie(fn func(d PersistentVolumeClaimSpecDie)) PersistentVolumeClaimTemplateDie
}

func (d *persistentVolumeClaimTemplateDie) MetadataDie(fn func(d diemetav1.ObjectMetaDie)) PersistentVolumeClaimTemplateDie {
	return d.DieStamp(func(r *corev1.PersistentVolumeClaimTemplate) {
		d := diemetav1.ObjectMetaBlank.DieImmutable(false).DieFeed(r.ObjectMeta)
		fn(d)
		r.ObjectMeta = d.DieRelease()
	})
}

func (d *persistentVolumeClaimTemplateDie) SpecDie(fn func(d PersistentVolumeClaimSpecDie)) PersistentVolumeClaimTemplateDie {
	return d.DieStamp(func(r *corev1.PersistentVolumeClaimTemplate) {
		d := PersistentVolumeClaimSpecBlank.DieImmutable(false).DieFeed(r.Spec)
		fn(d)
		r.Spec = d.DieRelease()
	})
}
