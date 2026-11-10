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

package v1_test

import (
	appsv1 "k8s.io/api/apps/v1"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	diemetav1 "reconciler.io/dies/v2/apis/meta/v1"
)

// An object die whose metadata's die is from this package is a DieObjectOf
// of this package, with Metadata(), as well as the die API's.
var (
	_ diemetav1.DieObjectOf[appsv1.Deployment, *dieappsv1.DeploymentDie, *dieappsv1.DeploymentDieMut]    = (*dieappsv1.DeploymentDie)(nil)
	_ diemetav1.DieObjectMutOf[appsv1.Deployment, *dieappsv1.DeploymentDie, *dieappsv1.DeploymentDieMut] = (*dieappsv1.DeploymentDieMut)(nil)
)
