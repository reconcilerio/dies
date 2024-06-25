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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	diemetav1 "reconciler.io/dies/apis/meta/v1"
)

// +die:object=true,apiVersion=v1,kind=Service
type _ = corev1.Service

// +die
// +die:field:name=SessionAffinityConfig,die=SessionAffinityConfigDie,pointer=true
type _ = corev1.ServiceSpec

func (d *ServiceSpecDie) PortDie(protocol corev1.Protocol, port int32, fn func(d *ServicePortDie)) *ServiceSpecDie {
	return d.DieStamp(func(r *corev1.ServiceSpec) {
		for i := range r.Ports {
			if protocol == r.Ports[i].Protocol && port == r.Ports[i].Port {
				d := ServicePortBlank.DieImmutable(false).DieFeed(r.Ports[i])
				fn(d)
				r.Ports[i] = d.DieRelease()
				return
			}
		}

		d := ServicePortBlank.DieImmutable(false).DieFeed(corev1.ServicePort{Protocol: protocol, Port: port})
		fn(d)
		r.Ports = append(r.Ports, d.DieRelease())
	})
}

func (d *ServiceSpecDie) AddSelector(key, value string) *ServiceSpecDie {
	return d.DieStamp(func(r *corev1.ServiceSpec) {
		if r.Selector == nil {
			r.Selector = map[string]string{}
		}
		r.Selector[key] = value
	})
}

// +die
type _ = corev1.ServicePort

// +die
// +die:field:name=ClientIP,die=ClientIPConfigDie,pointer=true
type _ = corev1.SessionAffinityConfig

// +die
type _ = corev1.ClientIPConfig

// +die
// +die:field:name=LoadBalancer,die=LoadBalancerStatusDie
type _ = corev1.ServiceStatus

func (d *ServiceStatusDie) ConditionsDie(conditions ...*diemetav1.ConditionDie) *ServiceStatusDie {
	return d.DieStamp(func(r *corev1.ServiceStatus) {
		r.Conditions = make([]metav1.Condition, len(conditions))
		for i, c := range conditions {
			r.Conditions[i] = c.DieRelease()
		}
	})
}

// +die
// +die:field:name=Ingress,die=LoadBalancerIngressDie,listType=atomic
type _ = corev1.LoadBalancerStatus

// deprecated: use IngressDie
func (d *LoadBalancerStatusDie) LoadBalancerDie(ingress ...*LoadBalancerIngressDie) *LoadBalancerStatusDie {
	return d.IngressDie(ingress...)
}

// +die
// +die:field:name=Ports,die=PortStatusDie,listType=atomic
type _ = corev1.LoadBalancerIngress

// +die
type _ = corev1.PortStatus
