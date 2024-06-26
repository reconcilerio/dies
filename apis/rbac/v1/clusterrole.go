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
	dierbacv1 "reconciler.io/dies/apis/authorization/rbac/v1"
)

// Deprecated ClusterRoleDie, moved to package reconciler.io/dies/apis/authorization/rbac/v1
type ClusterRoleDie = dierbacv1.ClusterRoleDie

// Deprecated AggregationRuleDie, moved to package reconciler.io/dies/apis/authorization/rbac/v1
type AggregationRuleDie = dierbacv1.AggregationRuleDie
