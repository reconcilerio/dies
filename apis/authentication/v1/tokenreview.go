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
	authenticationv1 "k8s.io/api/authentication/v1"
)

// +die:object=true,apiVersion=authentication.k8s.io/v1,kind=TokenReview
type _ = authenticationv1.TokenReview

// +die
// +die:field:name=BoundObjectRef,die=BoundObjectReferenceDie,pointer=true
type _ = authenticationv1.TokenRequestSpec

// +die
type _ = authenticationv1.BoundObjectReference

// +die
type _ = authenticationv1.TokenRequestStatus
