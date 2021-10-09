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

// +die:target=k8s.io/api/coordination/v1.Lease,object=true

// +die:target=k8s.io/api/coordination/v1.LeaseSpec
// +die:field:receiver=LeaseSpecDie,name=HolderIdentity,type=*string
// +die:field:receiver=LeaseSpecDie,name=LeaseDurationSeconds,type=*int32
// +die:field:receiver=LeaseSpecDie,name=AcquireTime,type=*k8s.io/apimachinery/pkg/apis/meta/v1.MicroTime
// +die:field:receiver=LeaseSpecDie,name=RenewTime,type=*k8s.io/apimachinery/pkg/apis/meta/v1.MicroTime
// +die:field:receiver=LeaseSpecDie,name=LeaseTransitions,type=*int32
