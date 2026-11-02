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

// Package authentication registers the dies of every version of the authentication API group.
package authentication

import (
	"k8s.io/apimachinery/pkg/runtime"

	authenticationv1 "reconciler.io/dies/v2/apis/authentication/v1"
	authenticationv1alpha1 "reconciler.io/dies/v2/apis/authentication/v1alpha1"
	authenticationv1beta1 "reconciler.io/dies/v2/apis/authentication/v1beta1"
)

const GroupName = authenticationv1.GroupName

// SchemeBuilder registers the dies of every version of the group, in order
// of version priority.
var SchemeBuilder = runtime.NewSchemeBuilder(
	authenticationv1.AddToScheme,
	authenticationv1beta1.AddToScheme,
	authenticationv1alpha1.AddToScheme,
)

// AddToScheme registers the dies of every version of the group with a
// scheme.
var AddToScheme = SchemeBuilder.AddToScheme
