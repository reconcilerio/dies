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

// Package apps registers the dies of every version of the apps API group.
package apps

import (
	"k8s.io/apimachinery/pkg/runtime"

	appsv1 "reconciler.io/dies/v2/apis/apps/v1"
	appsv1beta1 "reconciler.io/dies/v2/apis/apps/v1beta1"
	appsv1beta2 "reconciler.io/dies/v2/apis/apps/v1beta2"
)

const GroupName = appsv1.GroupName

// SchemeBuilder registers the dies of every version of the group, in order
// of version priority.
var SchemeBuilder = runtime.NewSchemeBuilder(
	appsv1.AddToScheme,
	appsv1beta2.AddToScheme,
	appsv1beta1.AddToScheme,
)

// AddToScheme registers the dies of every version of the group with a
// scheme.
var AddToScheme = SchemeBuilder.AddToScheme
