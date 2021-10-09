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

// +die:target=k8s.io/api/core/v1.Event,object=true
// +die:field:receiver=EventDie,name=InvolvedObject,type=k8s.io/api/core/v1.ObjectReference
// +die:field:receiver=EventDie,name=Reason,type=string
// +die:field:receiver=EventDie,name=Message,type=string
// +die:field:receiver=EventDie,name=Source,type=k8s.io/api/core/v1.EventSource
// +die:field:receiver=EventDie,name=FirstTimestamp,type=k8s.io/apimachinery/pkg/apis/meta/v1.Time
// +die:field:receiver=EventDie,name=LastTimestamp,type=k8s.io/apimachinery/pkg/apis/meta/v1.Time
// +die:field:receiver=EventDie,name=Count,type=int32
// +die:field:receiver=EventDie,name=Type,type=string
// +die:field:receiver=EventDie,name=EventTime,type=k8s.io/apimachinery/pkg/apis/meta/v1.MicroTime
// +die:field:receiver=EventDie,name=Series,type=*k8s.io/api/core/v1.EventSeries
// +die:field:receiver=EventDie,name=Action,type=string
// +die:field:receiver=EventDie,name=Related,type=*k8s.io/api/core/v1.ObjectReference
// +die:field:receiver=EventDie,name=ReportingController,type=string
// +die:field:receiver=EventDie,name=ReportingInstance,type=string
