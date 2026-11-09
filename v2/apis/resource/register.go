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

// Package resource registers the dies of every version of the resource API group.
package resource

import (
	resourcev1 "reconciler.io/dies/v2/apis/resource/v1"
	resourcev1alpha3 "reconciler.io/dies/v2/apis/resource/v1alpha3"
	resourcev1beta1 "reconciler.io/dies/v2/apis/resource/v1beta1"
	resourcev1beta2 "reconciler.io/dies/v2/apis/resource/v1beta2"
	"reconciler.io/dies/v2/die/schemes"
)

const GroupName = resourcev1.GroupName

// SchemeBuilder registers the types and dies of every version of the group,
// of each category, in order
// of version priority.
var SchemeBuilder = schemes.NewBuilder(
	resourcev1.SchemeBuilder.Add,
	resourcev1beta2.SchemeBuilder.Add,
	resourcev1beta1.SchemeBuilder.Add,
	resourcev1alpha3.SchemeBuilder.Add,
)

// AddToScheme registers the mutable dies of every version of the group
// with a scheme. SchemeBuilder.AddToSchemeAs registers another category, and
// SchemeBuilder.AddToComposite each.
var AddToScheme = SchemeBuilder.AddToScheme
