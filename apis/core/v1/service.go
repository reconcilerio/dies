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
	diemetav1 "github.com/scothis/dies/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +die:target=k8s.io/api/core/v1.Service,object=true

// +die:target=k8s.io/api/core/v1.ServiceSpec
// +die:field:receiver=ServiceSpecDie,name=Ports,type=[]k8s.io/api/core/v1.ServicePort
// +die:field:receiver=ServiceSpecDie,name=Selector,type=map[string]string
// +die:field:receiver=ServiceSpecDie,name=ClusterIP,type=string
// +die:field:receiver=ServiceSpecDie,name=ClusterIPs,type=[]string
// +die:field:receiver=ServiceSpecDie,name=Type,type=k8s.io/api/core/v1.ServiceType
// +die:field:receiver=ServiceSpecDie,name=ExternalIPs,type=[]string
// +die:field:receiver=ServiceSpecDie,name=SessionAffinity,type=k8s.io/api/core/v1.ServiceAffinity
// +die:field:receiver=ServiceSpecDie,name=LoadBalancerIP,type=string
// +die:field:receiver=ServiceSpecDie,name=LoadBalancerSourceRanges,type=[]string
// +die:field:receiver=ServiceSpecDie,name=ExternalName,type=string
// +die:field:receiver=ServiceSpecDie,name=ExternalTrafficPolicy,type=k8s.io/api/core/v1.ServiceExternalTrafficPolicyType
// +die:field:receiver=ServiceSpecDie,name=HealthCheckNodePort,type=int32
// +die:field:receiver=ServiceSpecDie,name=PublishNotReadyAddresses,type=bool
// +die:field:receiver=ServiceSpecDie,name=SessionAffinityConfig,type=*k8s.io/api/core/v1.SessionAffinityConfig
// +die:field:receiver=ServiceSpecDie,name=IPFamilies,type=[]k8s.io/api/core/v1.IPFamily
// +die:field:receiver=ServiceSpecDie,name=IPFamilyPolicy,type=*k8s.io/api/core/v1.IPFamilyPolicyType
// +die:field:receiver=ServiceSpecDie,name=AllocateLoadBalancerNodePorts,type=*bool
// +die:field:receiver=ServiceSpecDie,name=LoadBalancerClass,type=*string
// +die:field:receiver=ServiceSpecDie,name=InternalTrafficPolicy,type=*k8s.io/api/core/v1.ServiceInternalTrafficPolicyType

// +die:target=k8s.io/api/core/v1.ServiceStatus
// +die:field:receiver=ServiceStatusDie,name=LoadBalancer,type=k8s.io/api/core/v1.LoadBalancerStatus
// +die:field:receiver=ServiceStatusDie,name=Conditions,type=[]k8s.io/apimachinery/pkg/apis/meta/v1.Condition

func (d *ServiceStatusDie) ConditionsDie(conditions ...diemetav1.ConditionDie) *ServiceStatusDie {
	return d.DieStamp(func(r *corev1.ServiceStatus) {
		r.Conditions = make([]metav1.Condition, len(conditions))
		for i, c := range conditions {
			r.Conditions[i] = c.DieRelease()
		}
	})
}
