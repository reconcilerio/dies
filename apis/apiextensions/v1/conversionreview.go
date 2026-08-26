/*
Copyright 2021 the original author or authors.

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
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

// +die
// +die:field:name=TypeMeta,die=TypeMetaDie,package=_/meta/v1
// +die:field:name=Request,die=ConversionRequestDie,pointer=true
// +die:field:name=Response,die=ConversionResponseDie,pointer=true
type _ = apiextensionsv1.ConversionReview

// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
func (d *ConversionReviewDie) APIVersion(v string) *ConversionReviewDie {
	return d.DieStamp(func(r *apiextensionsv1.ConversionReview) {
		r.APIVersion = v
	})
}

// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
func (d *ConversionReviewDie) Kind(v string) *ConversionReviewDie {
	return d.DieStamp(func(r *apiextensionsv1.ConversionReview) {
		r.Kind = v
	})
}

// +die
type _ = apiextensionsv1.ConversionRequest

// +die
// +die:field:name=Result,die=StatusDie,package=_/meta/v1
type _ = apiextensionsv1.ConversionResponse
