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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +die:object=true
type _ = corev1.Service

// +die
type _ = corev1.ServiceSpec

// +die
type _ = corev1.ServiceStatus

type serviceStatus interface {
	LoadBalancerDie(fn func(d LoadBalancerStatusDie)) ServiceStatusDie
	ConditionsDie(conditions ...diemetav1.ConditionDie) ServiceStatusDie
}

func (d *serviceStatusDie) LoadBalancerDie(fn func(d LoadBalancerStatusDie)) ServiceStatusDie {
	return d.DieStamp(func(r *corev1.ServiceStatus) {
		d := LoadBalancerStatusBlank.DieImmutable(false).DieFeed(r.LoadBalancer)
		fn(d)
		r.LoadBalancer = d.DieRelease()
	})
}

func (d *serviceStatusDie) ConditionsDie(conditions ...diemetav1.ConditionDie) ServiceStatusDie {
	return d.DieStamp(func(r *corev1.ServiceStatus) {
		r.Conditions = make([]metav1.Condition, len(conditions))
		for i, c := range conditions {
			r.Conditions[i] = c.DieRelease()
		}
	})
}

// +die
type _ = corev1.LoadBalancerStatus

type loadBalancerStatus interface {
	LoadBalancerDie(ingress ...LoadBalancerIngressDie) LoadBalancerStatusDie
}

func (d *loadBalancerStatusDie) LoadBalancerDie(ingress ...LoadBalancerIngressDie) LoadBalancerStatusDie {
	return d.DieStamp(func(r *corev1.LoadBalancerStatus) {
		r.Ingress = make([]corev1.LoadBalancerIngress, len(ingress))
		for i := range ingress {
			r.Ingress[i] = ingress[i].DieRelease()
		}
	})
}

// +die
type _ = corev1.LoadBalancerIngress

type loadBalancerIngress interface {
	PortsDie(ports ...PortStatusDie) LoadBalancerIngressDie
}

func (d *loadBalancerIngressDie) PortsDie(ports ...PortStatusDie) LoadBalancerIngressDie {
	return d.DieStamp(func(r *corev1.LoadBalancerIngress) {
		r.Ports = make([]corev1.PortStatus, len(ports))
		for i := range ports {
			r.Ports[i] = ports[i].DieRelease()
		}
	})
}

// +die
type _ = corev1.PortStatus
