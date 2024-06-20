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

// +die:object=true,apiVersion=networking.k8s.io/v1,kind=Ingress
type _ = networkingv1.Ingress

// +die
// +die:field:name=DefaultBackend,die=IngressBackendDie,pointer=true
// +die:field:name=TLS,die=IngressTLSDie,listType=atomic
// +die:field:name=Rules,die=IngressRuleDie,listType=atomic
type _ = networkingv1.IngressSpec

// +die
// +die:field:name=Service,die=IngressServiceBackendDie,pointer=true
// +die:field:name=Resource,package=_/core/v1,die=TypedLocalObjectReferenceDie,pointer=true
type _ = networkingv1.IngressBackend

// +die
// +die:field:name=Port,die=ServiceBackendPortDie
type _ = networkingv1.IngressServiceBackend

// +die
type _ = networkingv1.ServiceBackendPort

// +die
type _ = networkingv1.IngressTLS

// +die
// +die:field:name=HTTP,die=HTTPIngressRuleValueDie,pointer=true
type _ = networkingv1.IngressRule

// +die
// +die:field:name=Paths,die=HTTPIngressPathDie,listType=atomic
type _ = networkingv1.HTTPIngressRuleValue

// +die
// +die:field:name=Backend,die=IngressBackendDie
type _ = networkingv1.HTTPIngressPath

// +die
// +die:field:name=LoadBalancer,die=IngressLoadBalancerStatusDie
type IngressStatus = networkingv1.IngressStatus

// +die
// +die:field:name=Ingress,die=IngressLoadBalancerIngressDie,listType=atomic
type IngressLoadBalancerStatus = networkingv1.IngressLoadBalancerStatus

// +die
// +die:field:name=Ports,die=IngressPortStatusDie,listType=atomic
type IngressLoadBalancerIngress = networkingv1.IngressLoadBalancerIngress

// +die
type IngressPortStatus = networkingv1.IngressPortStatus
