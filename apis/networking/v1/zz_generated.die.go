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
	json "encoding/json"
	fmtx "fmt"
	metav1 "github.com/scothis/dies/apis/meta/v1"
	util "github.com/scothis/dies/util"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	GroupVersion  = schema.GroupVersion{Group: "networking.k8s.io", Version: "v1"}
	SchemeBuilder = runtime.NewSchemeBuilder()
	AddToScheme   = SchemeBuilder.AddToScheme
)

type IngressDie struct {
	metav1.FrozenObjectMeta
	mutable bool
	r       networkingv1.Ingress
}

var IngressBlank = (&IngressDie{}).DieFeed(networkingv1.Ingress{})

func (d *IngressDie) DieImmutable(immutable bool) *IngressDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

func (d *IngressDie) DieFeed(r networkingv1.Ingress) *IngressDie {
	if d.mutable {
		d.FrozenObjectMeta = metav1.FreezeObjectMeta(r.ObjectMeta)
		d.r = r
		return d
	}
	return &IngressDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

func (d *IngressDie) DieRelease() networkingv1.Ingress {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

func (d *IngressDie) DieStamp(fn func(r *networkingv1.Ingress)) *IngressDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

func (d *IngressDie) DeepCopy() *IngressDie {
	r := *d.r.DeepCopy()
	return &IngressDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

func (d *IngressDie) DeepCopyObject() runtime.Object {
	return d.DeepCopy()
}

func (d *IngressDie) GetObjectKind() schema.ObjectKind {
	r := d.DieRelease()
	return r.GetObjectKind()
}

func (d *IngressDie) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.r)
}

func (d *IngressDie) UnmarshalJSON(b []byte) error {
	if d == IngressBlank {
		return fmtx.Errorf("cannot unmarshing into the root object, create a copy first")
	}
	r := &networkingv1.Ingress{}
	err := json.Unmarshal(b, r)
	*d = *d.DieFeed(*r)
	return err
}

func (d *IngressDie) MetadataDie(fn func(d *metav1.ObjectMetaDie)) *IngressDie {
	return d.DieStamp(func(r *networkingv1.Ingress) {
		d := metav1.ObjectMetaBlank.DieImmutable(false).DieFeed(r.ObjectMeta)
		fn(d)
		r.ObjectMeta = d.DieRelease()
	})
}

func (d *IngressDie) Spec(v networkingv1.IngressSpec) *IngressDie {
	return d.DieStamp(func(r *networkingv1.Ingress) {
		r.Spec = v
	})
}

func (d *IngressDie) SpecDie(fn func(d *IngressSpecDie)) *IngressDie {
	return d.DieStamp(func(r *networkingv1.Ingress) {
		d := IngressSpecBlank.DieImmutable(false).DieFeed(r.Spec)
		fn(d)
		r.Spec = d.DieRelease()
	})
}

func (d *IngressDie) Status(v networkingv1.IngressStatus) *IngressDie {
	return d.DieStamp(func(r *networkingv1.Ingress) {
		r.Status = v
	})
}

func (d *IngressDie) StatusDie(fn func(d *IngressStatusDie)) *IngressDie {
	return d.DieStamp(func(r *networkingv1.Ingress) {
		d := IngressStatusBlank.DieImmutable(false).DieFeed(r.Status)
		fn(d)
		r.Status = d.DieRelease()
	})
}

var _ apismetav1.Object = (*IngressDie)(nil)
var _ apismetav1.ObjectMetaAccessor = (*IngressDie)(nil)
var _ runtime.Object = (*IngressDie)(nil)

func init() {
	gvk := GroupVersion.WithKind("Ingress")
	obj := &IngressDie{}
	util.Register(SchemeBuilder, gvk, obj)
}

type IngressSpecDie struct {
	mutable bool
	r       networkingv1.IngressSpec
}

var IngressSpecBlank = (&IngressSpecDie{}).DieFeed(networkingv1.IngressSpec{})

func (d *IngressSpecDie) DieImmutable(immutable bool) *IngressSpecDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

func (d *IngressSpecDie) DieFeed(r networkingv1.IngressSpec) *IngressSpecDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &IngressSpecDie{
		mutable: d.mutable,
		r:       r,
	}
}

func (d *IngressSpecDie) DieRelease() networkingv1.IngressSpec {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

func (d *IngressSpecDie) DieStamp(fn func(r *networkingv1.IngressSpec)) *IngressSpecDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

func (d *IngressSpecDie) DeepCopy() *IngressSpecDie {
	r := *d.r.DeepCopy()
	return &IngressSpecDie{
		mutable: d.mutable,
		r:       r,
	}
}

func (d *IngressSpecDie) IngressClassName(v *string) *IngressSpecDie {
	return d.DieStamp(func(r *networkingv1.IngressSpec) {
		r.IngressClassName = v
	})
}

func (d *IngressSpecDie) DefaultBackend(v *networkingv1.IngressBackend) *IngressSpecDie {
	return d.DieStamp(func(r *networkingv1.IngressSpec) {
		r.DefaultBackend = v
	})
}

func (d *IngressSpecDie) TLS(v ...networkingv1.IngressTLS) *IngressSpecDie {
	return d.DieStamp(func(r *networkingv1.IngressSpec) {
		r.TLS = v
	})
}

func (d *IngressSpecDie) Rules(v ...networkingv1.IngressRule) *IngressSpecDie {
	return d.DieStamp(func(r *networkingv1.IngressSpec) {
		r.Rules = v
	})
}

type IngressStatusDie struct {
	mutable bool
	r       networkingv1.IngressStatus
}

var IngressStatusBlank = (&IngressStatusDie{}).DieFeed(networkingv1.IngressStatus{})

func (d *IngressStatusDie) DieImmutable(immutable bool) *IngressStatusDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

func (d *IngressStatusDie) DieFeed(r networkingv1.IngressStatus) *IngressStatusDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &IngressStatusDie{
		mutable: d.mutable,
		r:       r,
	}
}

func (d *IngressStatusDie) DieRelease() networkingv1.IngressStatus {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

func (d *IngressStatusDie) DieStamp(fn func(r *networkingv1.IngressStatus)) *IngressStatusDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

func (d *IngressStatusDie) DeepCopy() *IngressStatusDie {
	r := *d.r.DeepCopy()
	return &IngressStatusDie{
		mutable: d.mutable,
		r:       r,
	}
}

func (d *IngressStatusDie) LoadBalancer(v corev1.LoadBalancerStatus) *IngressStatusDie {
	return d.DieStamp(func(r *networkingv1.IngressStatus) {
		r.LoadBalancer = v
	})
}

type IngressClassDie struct {
	metav1.FrozenObjectMeta
	mutable bool
	r       networkingv1.IngressClass
}

var IngressClassBlank = (&IngressClassDie{}).DieFeed(networkingv1.IngressClass{})

func (d *IngressClassDie) DieImmutable(immutable bool) *IngressClassDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

func (d *IngressClassDie) DieFeed(r networkingv1.IngressClass) *IngressClassDie {
	if d.mutable {
		d.FrozenObjectMeta = metav1.FreezeObjectMeta(r.ObjectMeta)
		d.r = r
		return d
	}
	return &IngressClassDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

func (d *IngressClassDie) DieRelease() networkingv1.IngressClass {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

func (d *IngressClassDie) DieStamp(fn func(r *networkingv1.IngressClass)) *IngressClassDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

func (d *IngressClassDie) DeepCopy() *IngressClassDie {
	r := *d.r.DeepCopy()
	return &IngressClassDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

func (d *IngressClassDie) DeepCopyObject() runtime.Object {
	return d.DeepCopy()
}

func (d *IngressClassDie) GetObjectKind() schema.ObjectKind {
	r := d.DieRelease()
	return r.GetObjectKind()
}

func (d *IngressClassDie) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.r)
}

func (d *IngressClassDie) UnmarshalJSON(b []byte) error {
	if d == IngressClassBlank {
		return fmtx.Errorf("cannot unmarshing into the root object, create a copy first")
	}
	r := &networkingv1.IngressClass{}
	err := json.Unmarshal(b, r)
	*d = *d.DieFeed(*r)
	return err
}

func (d *IngressClassDie) MetadataDie(fn func(d *metav1.ObjectMetaDie)) *IngressClassDie {
	return d.DieStamp(func(r *networkingv1.IngressClass) {
		d := metav1.ObjectMetaBlank.DieImmutable(false).DieFeed(r.ObjectMeta)
		fn(d)
		r.ObjectMeta = d.DieRelease()
	})
}

func (d *IngressClassDie) Spec(v networkingv1.IngressClassSpec) *IngressClassDie {
	return d.DieStamp(func(r *networkingv1.IngressClass) {
		r.Spec = v
	})
}

func (d *IngressClassDie) SpecDie(fn func(d *IngressClassSpecDie)) *IngressClassDie {
	return d.DieStamp(func(r *networkingv1.IngressClass) {
		d := IngressClassSpecBlank.DieImmutable(false).DieFeed(r.Spec)
		fn(d)
		r.Spec = d.DieRelease()
	})
}

var _ apismetav1.Object = (*IngressClassDie)(nil)
var _ apismetav1.ObjectMetaAccessor = (*IngressClassDie)(nil)
var _ runtime.Object = (*IngressClassDie)(nil)

func init() {
	gvk := GroupVersion.WithKind("IngressClass")
	obj := &IngressClassDie{}
	util.Register(SchemeBuilder, gvk, obj)
}

type IngressClassSpecDie struct {
	mutable bool
	r       networkingv1.IngressClassSpec
}

var IngressClassSpecBlank = (&IngressClassSpecDie{}).DieFeed(networkingv1.IngressClassSpec{})

func (d *IngressClassSpecDie) DieImmutable(immutable bool) *IngressClassSpecDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

func (d *IngressClassSpecDie) DieFeed(r networkingv1.IngressClassSpec) *IngressClassSpecDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &IngressClassSpecDie{
		mutable: d.mutable,
		r:       r,
	}
}

func (d *IngressClassSpecDie) DieRelease() networkingv1.IngressClassSpec {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

func (d *IngressClassSpecDie) DieStamp(fn func(r *networkingv1.IngressClassSpec)) *IngressClassSpecDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

func (d *IngressClassSpecDie) DeepCopy() *IngressClassSpecDie {
	r := *d.r.DeepCopy()
	return &IngressClassSpecDie{
		mutable: d.mutable,
		r:       r,
	}
}

func (d *IngressClassSpecDie) Controller(v string) *IngressClassSpecDie {
	return d.DieStamp(func(r *networkingv1.IngressClassSpec) {
		r.Controller = v
	})
}

func (d *IngressClassSpecDie) Parameters(v *networkingv1.IngressClassParametersReference) *IngressClassSpecDie {
	return d.DieStamp(func(r *networkingv1.IngressClassSpec) {
		r.Parameters = v
	})
}

type NetworkPolicyDie struct {
	metav1.FrozenObjectMeta
	mutable bool
	r       networkingv1.NetworkPolicy
}

var NetworkPolicyBlank = (&NetworkPolicyDie{}).DieFeed(networkingv1.NetworkPolicy{})

func (d *NetworkPolicyDie) DieImmutable(immutable bool) *NetworkPolicyDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

func (d *NetworkPolicyDie) DieFeed(r networkingv1.NetworkPolicy) *NetworkPolicyDie {
	if d.mutable {
		d.FrozenObjectMeta = metav1.FreezeObjectMeta(r.ObjectMeta)
		d.r = r
		return d
	}
	return &NetworkPolicyDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

func (d *NetworkPolicyDie) DieRelease() networkingv1.NetworkPolicy {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

func (d *NetworkPolicyDie) DieStamp(fn func(r *networkingv1.NetworkPolicy)) *NetworkPolicyDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

func (d *NetworkPolicyDie) DeepCopy() *NetworkPolicyDie {
	r := *d.r.DeepCopy()
	return &NetworkPolicyDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

func (d *NetworkPolicyDie) DeepCopyObject() runtime.Object {
	return d.DeepCopy()
}

func (d *NetworkPolicyDie) GetObjectKind() schema.ObjectKind {
	r := d.DieRelease()
	return r.GetObjectKind()
}

func (d *NetworkPolicyDie) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.r)
}

func (d *NetworkPolicyDie) UnmarshalJSON(b []byte) error {
	if d == NetworkPolicyBlank {
		return fmtx.Errorf("cannot unmarshing into the root object, create a copy first")
	}
	r := &networkingv1.NetworkPolicy{}
	err := json.Unmarshal(b, r)
	*d = *d.DieFeed(*r)
	return err
}

func (d *NetworkPolicyDie) MetadataDie(fn func(d *metav1.ObjectMetaDie)) *NetworkPolicyDie {
	return d.DieStamp(func(r *networkingv1.NetworkPolicy) {
		d := metav1.ObjectMetaBlank.DieImmutable(false).DieFeed(r.ObjectMeta)
		fn(d)
		r.ObjectMeta = d.DieRelease()
	})
}

func (d *NetworkPolicyDie) Spec(v networkingv1.NetworkPolicySpec) *NetworkPolicyDie {
	return d.DieStamp(func(r *networkingv1.NetworkPolicy) {
		r.Spec = v
	})
}

func (d *NetworkPolicyDie) SpecDie(fn func(d *NetworkPolicySpecDie)) *NetworkPolicyDie {
	return d.DieStamp(func(r *networkingv1.NetworkPolicy) {
		d := NetworkPolicySpecBlank.DieImmutable(false).DieFeed(r.Spec)
		fn(d)
		r.Spec = d.DieRelease()
	})
}

var _ apismetav1.Object = (*NetworkPolicyDie)(nil)
var _ apismetav1.ObjectMetaAccessor = (*NetworkPolicyDie)(nil)
var _ runtime.Object = (*NetworkPolicyDie)(nil)

func init() {
	gvk := GroupVersion.WithKind("NetworkPolicy")
	obj := &NetworkPolicyDie{}
	util.Register(SchemeBuilder, gvk, obj)
}

type NetworkPolicySpecDie struct {
	mutable bool
	r       networkingv1.NetworkPolicySpec
}

var NetworkPolicySpecBlank = (&NetworkPolicySpecDie{}).DieFeed(networkingv1.NetworkPolicySpec{})

func (d *NetworkPolicySpecDie) DieImmutable(immutable bool) *NetworkPolicySpecDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

func (d *NetworkPolicySpecDie) DieFeed(r networkingv1.NetworkPolicySpec) *NetworkPolicySpecDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &NetworkPolicySpecDie{
		mutable: d.mutable,
		r:       r,
	}
}

func (d *NetworkPolicySpecDie) DieRelease() networkingv1.NetworkPolicySpec {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

func (d *NetworkPolicySpecDie) DieStamp(fn func(r *networkingv1.NetworkPolicySpec)) *NetworkPolicySpecDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

func (d *NetworkPolicySpecDie) DeepCopy() *NetworkPolicySpecDie {
	r := *d.r.DeepCopy()
	return &NetworkPolicySpecDie{
		mutable: d.mutable,
		r:       r,
	}
}

func (d *NetworkPolicySpecDie) PodSelector(v apismetav1.LabelSelector) *NetworkPolicySpecDie {
	return d.DieStamp(func(r *networkingv1.NetworkPolicySpec) {
		r.PodSelector = v
	})
}

func (d *NetworkPolicySpecDie) Ingress(v ...networkingv1.NetworkPolicyIngressRule) *NetworkPolicySpecDie {
	return d.DieStamp(func(r *networkingv1.NetworkPolicySpec) {
		r.Ingress = v
	})
}

func (d *NetworkPolicySpecDie) Egress(v ...networkingv1.NetworkPolicyEgressRule) *NetworkPolicySpecDie {
	return d.DieStamp(func(r *networkingv1.NetworkPolicySpec) {
		r.Egress = v
	})
}

func (d *NetworkPolicySpecDie) PolicyTypes(v ...networkingv1.PolicyType) *NetworkPolicySpecDie {
	return d.DieStamp(func(r *networkingv1.NetworkPolicySpec) {
		r.PolicyTypes = v
	})
}
