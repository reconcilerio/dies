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
	networkingv1 "k8s.io/api/networking/v1"
)

// +die:object=true,apiVersion=networking.k8s.io/v1,kind=NetworkPolicy
type _ = networkingv1.NetworkPolicy

// +die
// +die:field:name=PodSelector,package=_/meta/v1,die=LabelSelectorDie
// +die:field:name=Ingress,die=NetworkPolicyIngressRuleDie,listType=atomic
// +die:field:name=Egress,die=NetworkPolicyEgressRuleDie,listType=atomic
type _ = networkingv1.NetworkPolicySpec

// +die
// +die:field:name=Ports,die=NetworkPolicyPortDie,listType=atomic
// +die:field:name=From,die=NetworkPolicyPeerDie,listType=atomic
type _ = networkingv1.NetworkPolicyIngressRule

// +die
// +die:field:name=Ports,die=NetworkPolicyPortDie,listType=atomic
// +die:field:name=To,die=NetworkPolicyPeerDie,listType=atomic
type _ = networkingv1.NetworkPolicyEgressRule

// +die
type _ = networkingv1.NetworkPolicyPort

// +die
// +die:field:name=PodSelector,package=_/meta/v1,die=LabelSelectorDie,pointer=true
// +die:field:name=NamespaceSelector,package=_/meta/v1,die=LabelSelectorDie,pointer=true
// +die:field:name=IPBlock,die=IPBlockDie,pointer=true
type _ = networkingv1.NetworkPolicyPeer

// +die
type _ = networkingv1.IPBlock
