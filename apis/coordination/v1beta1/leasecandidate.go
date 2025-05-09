/*
Copyright 2025 the original author or authors.

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

package v1beta1

import (
	coordinationv1beta1 "k8s.io/api/coordination/v1beta1"
)

// +die:object=true,apiVersion=coordination.k8s.io/v1beta1,kind=LeaseCandidate
type _ = coordinationv1beta1.LeaseCandidate

// +die
type _ = coordinationv1beta1.LeaseCandidateSpec
