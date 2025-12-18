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

package v1alpha1

import (
	imagepolicyv1alpha1 "k8s.io/api/imagepolicy/v1alpha1"
)

// +die:object=true,apiVersion=imagepolicy.k8s.io/v1alpha1,kind=ImageReview
type _ = imagepolicyv1alpha1.ImageReview

// +die
// +die:field:name=Containers,die=ImageReviewContainerSpecDie,listType=atomic
type _ = imagepolicyv1alpha1.ImageReviewSpec

func (d *ImageReviewSpecDie) AddAnnotation(key, value string) *ImageReviewSpecDie {
	return d.DieStamp(func(r *imagepolicyv1alpha1.ImageReviewSpec) {
		if r.Annotations == nil {
			r.Annotations = map[string]string{}
		}
		r.Annotations[key] = value
	})
}

// +die
type _ = imagepolicyv1alpha1.ImageReviewContainerSpec

// +die
type _ = imagepolicyv1alpha1.ImageReviewStatus

func (d *ImageReviewStatusDie) AddAuditAnnotation(key, value string) *ImageReviewStatusDie {
	return d.DieStamp(func(r *imagepolicyv1alpha1.ImageReviewStatus) {
		if r.AuditAnnotations == nil {
			r.AuditAnnotations = map[string]string{}
		}
		r.AuditAnnotations[key] = value
	})
}
