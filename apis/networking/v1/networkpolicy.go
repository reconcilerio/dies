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
	networkingv1 "k8s.io/api/networking/v1"
)

// +die:object=true
type _ = networkingv1.NetworkPolicy

// +die
type _ = networkingv1.NetworkPolicySpec

func (d *NetworkPolicySpecDie) PodSelectorDie(fn func(d *diemetav1.LabelSelectorDie)) *NetworkPolicySpecDie {
	return d.DieStamp(func(r *networkingv1.NetworkPolicySpec) {
		d := diemetav1.LabelSelectorBlank.DieImmutable(false).DieFeed(r.PodSelector)
		fn(d)
		r.PodSelector = d.DieRelease()
	})
}

func (d *NetworkPolicySpecDie) IngressDie(ingress ...*NetworkPolicyIngressRuleDie) *NetworkPolicySpecDie {
	return d.DieStamp(func(r *networkingv1.NetworkPolicySpec) {
		r.Ingress = make([]networkingv1.NetworkPolicyIngressRule, len(ingress))
		for i := range ingress {
			r.Ingress[i] = ingress[i].DieRelease()
		}
	})
}

func (d *NetworkPolicySpecDie) EgressDie(egress ...*NetworkPolicyEgressRuleDie) *NetworkPolicySpecDie {
	return d.DieStamp(func(r *networkingv1.NetworkPolicySpec) {
		r.Egress = make([]networkingv1.NetworkPolicyEgressRule, len(egress))
		for i := range egress {
			r.Egress[i] = egress[i].DieRelease()
		}
	})
}

// +die
type _ = networkingv1.NetworkPolicyIngressRule

func (d *NetworkPolicyIngressRuleDie) PortsDie(ports ...*NetworkPolicyPortDie) *NetworkPolicyIngressRuleDie {
	return d.DieStamp(func(r *networkingv1.NetworkPolicyIngressRule) {
		r.Ports = make([]networkingv1.NetworkPolicyPort, len(ports))
		for i := range ports {
			r.Ports[i] = ports[i].DieRelease()
		}
	})
}

func (d *NetworkPolicyIngressRuleDie) FromDie(from ...*NetworkPolicyPeerDie) *NetworkPolicyIngressRuleDie {
	return d.DieStamp(func(r *networkingv1.NetworkPolicyIngressRule) {
		r.From = make([]networkingv1.NetworkPolicyPeer, len(from))
		for i := range from {
			r.From[i] = from[i].DieRelease()
		}
	})
}

// +die
type _ = networkingv1.NetworkPolicyEgressRule

func (d *NetworkPolicyEgressRuleDie) PortsDie(ports ...*NetworkPolicyPortDie) *NetworkPolicyEgressRuleDie {
	return d.DieStamp(func(r *networkingv1.NetworkPolicyEgressRule) {
		r.Ports = make([]networkingv1.NetworkPolicyPort, len(ports))
		for i := range ports {
			r.Ports[i] = ports[i].DieRelease()
		}
	})
}

func (d *NetworkPolicyEgressRuleDie) ToDie(to ...*NetworkPolicyPeerDie) *NetworkPolicyEgressRuleDie {
	return d.DieStamp(func(r *networkingv1.NetworkPolicyEgressRule) {
		r.To = make([]networkingv1.NetworkPolicyPeer, len(to))
		for i := range to {
			r.To[i] = to[i].DieRelease()
		}
	})
}

// +die
type _ = networkingv1.NetworkPolicyPort

// +die
type _ = networkingv1.NetworkPolicyPeer

func (d *NetworkPolicyPeerDie) PodSelectorDie(fn func(d *diemetav1.LabelSelectorDie)) *NetworkPolicyPeerDie {
	return d.DieStamp(func(r *networkingv1.NetworkPolicyPeer) {
		d := diemetav1.LabelSelectorBlank.DieImmutable(false).DieFeedPtr(r.PodSelector)
		fn(d)
		r.PodSelector = d.DieReleasePtr()
	})
}

func (d *NetworkPolicyPeerDie) NamespaceSelectorDie(fn func(d *diemetav1.LabelSelectorDie)) *NetworkPolicyPeerDie {
	return d.DieStamp(func(r *networkingv1.NetworkPolicyPeer) {
		d := diemetav1.LabelSelectorBlank.DieImmutable(false).DieFeedPtr(r.NamespaceSelector)
		fn(d)
		r.NamespaceSelector = d.DieReleasePtr()
	})
}

func (d *NetworkPolicyPeerDie) IPBlockDie(fn func(d *IPBlockDie)) *NetworkPolicyPeerDie {
	return d.DieStamp(func(r *networkingv1.NetworkPolicyPeer) {
		d := IPBlockBlank.DieImmutable(false).DieFeedPtr(r.IPBlock)
		fn(d)
		r.IPBlock = d.DieReleasePtr()
	})
}

// +die
type _ = networkingv1.IPBlock
