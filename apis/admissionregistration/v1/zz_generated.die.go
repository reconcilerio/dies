//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by diegen. DO NOT EDIT.

package v1

import (
	metav1 "dies.dev/apis/meta/v1"
	json "encoding/json"
	fmtx "fmt"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
)

type MutatingWebhookConfigurationDie struct {
	metav1.FrozenObjectMeta
	mutable bool
	r       admissionregistrationv1.MutatingWebhookConfiguration
}

var MutatingWebhookConfigurationBlank = (&MutatingWebhookConfigurationDie{}).DieFeed(admissionregistrationv1.MutatingWebhookConfiguration{})

func (d *MutatingWebhookConfigurationDie) DieImmutable(immutable bool) *MutatingWebhookConfigurationDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

func (d *MutatingWebhookConfigurationDie) DieFeed(r admissionregistrationv1.MutatingWebhookConfiguration) *MutatingWebhookConfigurationDie {
	if d.mutable {
		d.FrozenObjectMeta = metav1.FreezeObjectMeta(r.ObjectMeta)
		d.r = r
		return d
	}
	return &MutatingWebhookConfigurationDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

func (d *MutatingWebhookConfigurationDie) DieFeedPtr(r *admissionregistrationv1.MutatingWebhookConfiguration) *MutatingWebhookConfigurationDie {
	if r == nil {
		r = &admissionregistrationv1.MutatingWebhookConfiguration{}
	}
	return d.DieFeed(*r)
}

func (d *MutatingWebhookConfigurationDie) DieRelease() admissionregistrationv1.MutatingWebhookConfiguration {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

func (d *MutatingWebhookConfigurationDie) DieReleasePtr() *admissionregistrationv1.MutatingWebhookConfiguration {
	r := d.DieRelease()
	return &r
}

func (d *MutatingWebhookConfigurationDie) DieStamp(fn func(r *admissionregistrationv1.MutatingWebhookConfiguration)) *MutatingWebhookConfigurationDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

func (d *MutatingWebhookConfigurationDie) DeepCopy() *MutatingWebhookConfigurationDie {
	r := *d.r.DeepCopy()
	return &MutatingWebhookConfigurationDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

func (d *MutatingWebhookConfigurationDie) DeepCopyObject() runtime.Object {
	return d.r.DeepCopy()
}

func (d *MutatingWebhookConfigurationDie) GetObjectKind() schema.ObjectKind {
	r := d.DieRelease()
	return r.GetObjectKind()
}

func (d *MutatingWebhookConfigurationDie) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.r)
}

func (d *MutatingWebhookConfigurationDie) UnmarshalJSON(b []byte) error {
	if d == MutatingWebhookConfigurationBlank {
		return fmtx.Errorf("cannot unmarshal into the root object, create a copy first")
	}
	r := &admissionregistrationv1.MutatingWebhookConfiguration{}
	err := json.Unmarshal(b, r)
	*d = *d.DieFeed(*r)
	return err
}

func (d *MutatingWebhookConfigurationDie) MetadataDie(fn func(d *metav1.ObjectMetaDie)) *MutatingWebhookConfigurationDie {
	return d.DieStamp(func(r *admissionregistrationv1.MutatingWebhookConfiguration) {
		d := metav1.ObjectMetaBlank.DieImmutable(false).DieFeed(r.ObjectMeta)
		fn(d)
		r.ObjectMeta = d.DieRelease()
	})
}

var _ apismetav1.Object = (*MutatingWebhookConfigurationDie)(nil)
var _ apismetav1.ObjectMetaAccessor = (*MutatingWebhookConfigurationDie)(nil)
var _ runtime.Object = (*MutatingWebhookConfigurationDie)(nil)

// Webhooks is a list of webhooks and the affected resources and operations.
func (d *MutatingWebhookConfigurationDie) Webhooks(v ...admissionregistrationv1.MutatingWebhook) *MutatingWebhookConfigurationDie {
	return d.DieStamp(func(r *admissionregistrationv1.MutatingWebhookConfiguration) {
		r.Webhooks = v
	})
}

type ValidatingWebhookConfigurationDie struct {
	metav1.FrozenObjectMeta
	mutable bool
	r       admissionregistrationv1.ValidatingWebhookConfiguration
}

var ValidatingWebhookConfigurationBlank = (&ValidatingWebhookConfigurationDie{}).DieFeed(admissionregistrationv1.ValidatingWebhookConfiguration{})

func (d *ValidatingWebhookConfigurationDie) DieImmutable(immutable bool) *ValidatingWebhookConfigurationDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

func (d *ValidatingWebhookConfigurationDie) DieFeed(r admissionregistrationv1.ValidatingWebhookConfiguration) *ValidatingWebhookConfigurationDie {
	if d.mutable {
		d.FrozenObjectMeta = metav1.FreezeObjectMeta(r.ObjectMeta)
		d.r = r
		return d
	}
	return &ValidatingWebhookConfigurationDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

func (d *ValidatingWebhookConfigurationDie) DieFeedPtr(r *admissionregistrationv1.ValidatingWebhookConfiguration) *ValidatingWebhookConfigurationDie {
	if r == nil {
		r = &admissionregistrationv1.ValidatingWebhookConfiguration{}
	}
	return d.DieFeed(*r)
}

func (d *ValidatingWebhookConfigurationDie) DieRelease() admissionregistrationv1.ValidatingWebhookConfiguration {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

func (d *ValidatingWebhookConfigurationDie) DieReleasePtr() *admissionregistrationv1.ValidatingWebhookConfiguration {
	r := d.DieRelease()
	return &r
}

func (d *ValidatingWebhookConfigurationDie) DieStamp(fn func(r *admissionregistrationv1.ValidatingWebhookConfiguration)) *ValidatingWebhookConfigurationDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

func (d *ValidatingWebhookConfigurationDie) DeepCopy() *ValidatingWebhookConfigurationDie {
	r := *d.r.DeepCopy()
	return &ValidatingWebhookConfigurationDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

func (d *ValidatingWebhookConfigurationDie) DeepCopyObject() runtime.Object {
	return d.r.DeepCopy()
}

func (d *ValidatingWebhookConfigurationDie) GetObjectKind() schema.ObjectKind {
	r := d.DieRelease()
	return r.GetObjectKind()
}

func (d *ValidatingWebhookConfigurationDie) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.r)
}

func (d *ValidatingWebhookConfigurationDie) UnmarshalJSON(b []byte) error {
	if d == ValidatingWebhookConfigurationBlank {
		return fmtx.Errorf("cannot unmarshal into the root object, create a copy first")
	}
	r := &admissionregistrationv1.ValidatingWebhookConfiguration{}
	err := json.Unmarshal(b, r)
	*d = *d.DieFeed(*r)
	return err
}

func (d *ValidatingWebhookConfigurationDie) MetadataDie(fn func(d *metav1.ObjectMetaDie)) *ValidatingWebhookConfigurationDie {
	return d.DieStamp(func(r *admissionregistrationv1.ValidatingWebhookConfiguration) {
		d := metav1.ObjectMetaBlank.DieImmutable(false).DieFeed(r.ObjectMeta)
		fn(d)
		r.ObjectMeta = d.DieRelease()
	})
}

var _ apismetav1.Object = (*ValidatingWebhookConfigurationDie)(nil)
var _ apismetav1.ObjectMetaAccessor = (*ValidatingWebhookConfigurationDie)(nil)
var _ runtime.Object = (*ValidatingWebhookConfigurationDie)(nil)

// Webhooks is a list of webhooks and the affected resources and operations.
func (d *ValidatingWebhookConfigurationDie) Webhooks(v ...admissionregistrationv1.ValidatingWebhook) *ValidatingWebhookConfigurationDie {
	return d.DieStamp(func(r *admissionregistrationv1.ValidatingWebhookConfiguration) {
		r.Webhooks = v
	})
}
