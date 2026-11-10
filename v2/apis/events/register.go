/*
Copyright 2026 the original author or authors.

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

// Package events registers the dies of every version of the events API group.
package events

import (
	eventsv1 "reconciler.io/dies/v2/apis/events/v1"
	eventsv1beta1 "reconciler.io/dies/v2/apis/events/v1beta1"
	"reconciler.io/dies/v2/die/schemes"
)

const GroupName = eventsv1.GroupName

// SchemeBuilder registers the types and dies of every version of the group,
// of each category, in order
// of version priority.
var SchemeBuilder = schemes.NewBuilder(
	eventsv1.SchemeBuilder.Add,
	eventsv1beta1.SchemeBuilder.Add,
)

// AddToScheme registers the mutable dies of every version of the group
// with a scheme. SchemeBuilder.AddToSchemeAs registers another category, and
// SchemeBuilder.AddToComposite each.
var AddToScheme = SchemeBuilder.AddToScheme
