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

// Package apiregistration registers the dies of every version of the apiregistration API group.
package apiregistration

import (
	"k8s.io/apimachinery/pkg/runtime"

	apiregistrationv1 "reconciler.io/dies/v2/apis/apiregistration/v1"
)

const GroupName = apiregistrationv1.GroupName

// SchemeBuilder registers the dies of every version of the group, in order
// of version priority.
var SchemeBuilder = runtime.NewSchemeBuilder(
	apiregistrationv1.AddToScheme,
)

// AddToScheme registers the dies of every version of the group with a
// scheme.
var AddToScheme = SchemeBuilder.AddToScheme
