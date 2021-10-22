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
	diecorev1 "dies.dev/apis/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
)

// +die:object=true
type _ = networkingv1.Ingress

// +die
type _ = networkingv1.IngressSpec

type ingressSpecDieExtension interface {
	DefaultBackendDie(fn func(d IngressBackendDie)) IngressSpecDie
	TLSDie(tls ...IngressTLSDie) IngressSpecDie
	RulesDie(rules ...IngressRuleDie) IngressSpecDie
}

func (d *ingressSpecDie) DefaultBackendDie(fn func(d IngressBackendDie)) IngressSpecDie {
	return d.DieStamp(func(r *networkingv1.IngressSpec) {
		d := IngressBackendBlank.DieImmutable(false).DieFeedPtr(r.DefaultBackend)
		fn(d)
		r.DefaultBackend = d.DieReleasePtr()
	})
}

func (d *ingressSpecDie) TLSDie(tls ...IngressTLSDie) IngressSpecDie {
	return d.DieStamp(func(r *networkingv1.IngressSpec) {
		r.TLS = make([]networkingv1.IngressTLS, len(tls))
		for i := range tls {
			r.TLS[i] = tls[i].DieRelease()
		}
	})
}

func (d *ingressSpecDie) RulesDie(rules ...IngressRuleDie) IngressSpecDie {
	return d.DieStamp(func(r *networkingv1.IngressSpec) {
		r.Rules = make([]networkingv1.IngressRule, len(rules))
		for i := range rules {
			r.Rules[i] = rules[i].DieRelease()
		}
	})
}

// +die
type _ = networkingv1.IngressBackend

type ingressBackendDieExtension interface {
	ServiceDie(fn func(d IngressServiceBackendDie)) IngressBackendDie
	ResourceDie(fn func(d diecorev1.TypedLocalObjectReferenceDie)) IngressBackendDie
}

func (d *ingressBackendDie) ServiceDie(fn func(d IngressServiceBackendDie)) IngressBackendDie {
	return d.DieStamp(func(r *networkingv1.IngressBackend) {
		d := IngressServiceBackendBlank.DieImmutable(false).DieFeedPtr(r.Service)
		fn(d)
		r.Service = d.DieReleasePtr()
	})
}

func (d *ingressBackendDie) ResourceDie(fn func(d diecorev1.TypedLocalObjectReferenceDie)) IngressBackendDie {
	return d.DieStamp(func(r *networkingv1.IngressBackend) {
		d := diecorev1.TypedLocalObjectReferenceBlank.DieImmutable(false).DieFeedPtr(r.Resource)
		fn(d)
		r.Resource = d.DieReleasePtr()
	})
}

// +die
type _ = networkingv1.IngressServiceBackend

type ingressServiceBackendDieExtension interface {
	PortDie(fn func(d ServiceBackendPortDie)) IngressServiceBackendDie
}

func (d *ingressServiceBackendDie) PortDie(fn func(d ServiceBackendPortDie)) IngressServiceBackendDie {
	return d.DieStamp(func(r *networkingv1.IngressServiceBackend) {
		d := ServiceBackendPortBlank.DieImmutable(false).DieFeed(r.Port)
		fn(d)
		r.Port = d.DieRelease()
	})
}

// +die
type _ = networkingv1.ServiceBackendPort

// +die
type _ = networkingv1.IngressTLS

// +die
type _ = networkingv1.IngressRule

type ingressRuleDieExtension interface {
	HTTPDie(fn func(d HTTPIngressRuleValueDie)) IngressRuleDie
}

func (d *ingressRuleDie) HTTPDie(fn func(d HTTPIngressRuleValueDie)) IngressRuleDie {
	return d.DieStamp(func(r *networkingv1.IngressRule) {
		d := HTTPIngressRuleValueBlank.DieImmutable(false).DieFeedPtr(r.HTTP)
		fn(d)
		r.HTTP = d.DieReleasePtr()
	})
}

// +die
type _ = networkingv1.HTTPIngressRuleValue

type httpIngressRuleValueDieExtension interface {
	PathsDie(paths ...HTTPIngressPathDie) HTTPIngressRuleValueDie
}

func (d *httpIngressRuleValueDie) PathsDie(paths ...HTTPIngressPathDie) HTTPIngressRuleValueDie {
	return d.DieStamp(func(r *networkingv1.HTTPIngressRuleValue) {
		r.Paths = make([]networkingv1.HTTPIngressPath, len(paths))
		for i := range paths {
			r.Paths[i] = paths[i].DieRelease()
		}
	})
}

// +die
type _ = networkingv1.HTTPIngressPath

type httpIngressPathDieExtension interface {
	BackendDie(fn func(d IngressBackendDie)) HTTPIngressPathDie
}

func (d *httpIngressPathDie) BackendDie(fn func(d IngressBackendDie)) HTTPIngressPathDie {
	return d.DieStamp(func(r *networkingv1.HTTPIngressPath) {
		d := IngressBackendBlank.DieImmutable(false).DieFeed(r.Backend)
		fn(d)
		r.Backend = d.DieRelease()
	})
}

// +die
type IngressStatus = networkingv1.IngressStatus

type ingressStatusDieExtension interface {
	LoadBalancerDie(fn func(d diecorev1.LoadBalancerStatusDie)) IngressStatusDie
}

func (d *ingressStatusDie) LoadBalancerDie(fn func(d diecorev1.LoadBalancerStatusDie)) IngressStatusDie {
	return d.DieStamp(func(r *networkingv1.IngressStatus) {
		d := diecorev1.LoadBalancerStatusBlank.DieImmutable(false).DieFeed(r.LoadBalancer)
		fn(d)
		r.LoadBalancer = d.DieRelease()
	})
}
