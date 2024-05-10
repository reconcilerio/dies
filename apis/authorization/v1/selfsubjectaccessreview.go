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
	authorizationv1 "k8s.io/api/authorization/v1"
)

// +die:object=true,apiVersion=authorization.k8s.io/v1,kind=SelfSubjectAccessReview,status=SubjectAccessReviewStatus
type _ = authorizationv1.SelfSubjectAccessReview

// +die
type _ = authorizationv1.SelfSubjectAccessReviewSpec

func (d *SelfSubjectAccessReviewSpecDie) ResourceAttributesDie(fn func(d *ResourceAttributesDie)) *SelfSubjectAccessReviewSpecDie {
	return d.DieStamp(func(r *authorizationv1.SelfSubjectAccessReviewSpec) {
		d := ResourceAttributesBlank.DieImmutable(false).DieFeedPtr(r.ResourceAttributes)
		fn(d)
		r.ResourceAttributes = d.DieReleasePtr()
	})
}

func (d *SelfSubjectAccessReviewSpecDie) NonResourceAttributesDie(fn func(d *NonResourceAttributesDie)) *SelfSubjectAccessReviewSpecDie {
	return d.DieStamp(func(r *authorizationv1.SelfSubjectAccessReviewSpec) {
		d := NonResourceAttributesBlank.DieImmutable(false).DieFeedPtr(r.NonResourceAttributes)
		fn(d)
		r.NonResourceAttributes = d.DieReleasePtr()
	})
}
