/*
Copyright 2022 the original author or authors.

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
	eventsv1 "k8s.io/api/events/v1"
)

// +die:object=true,apiVersion=events.k8s.io/v1,kind=Event
// +die:field:name=Series,die=EventSeriesDie,pointer=true
// +die:field:name=Regarding,package=_/core/v1,die=ObjectReferenceDie
// +die:field:name=Related,package=_/core/v1,die=ObjectReferenceDie,pointer=true
// +die:field:name=DeprecatedSource,package=_/core/v1,die=EventSourceDie
type _ = eventsv1.Event

// +die
type _ = eventsv1.EventSeries
