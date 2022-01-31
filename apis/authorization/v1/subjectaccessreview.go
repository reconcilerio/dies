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

// +die:object=true
type _ = authorizationv1.SubjectAccessReview

// TODO(scothis) fix import for maps with struct values, ignore the 'Extra' field until then

// +die:ignore={Extra}
type _ = authorizationv1.SubjectAccessReviewSpec

func (d *SubjectAccessReviewSpecDie) ResourceAttributesDie(fn func(d *ResourceAttributesDie)) *SubjectAccessReviewSpecDie {
	return d.DieStamp(func(r *authorizationv1.SubjectAccessReviewSpec) {
		d := ResourceAttributesBlank.DieImmutable(false).DieFeedPtr(r.ResourceAttributes)
		fn(d)
		r.ResourceAttributes = d.DieReleasePtr()
	})
}

func (d *SubjectAccessReviewSpecDie) NonResourceAttributesDie(fn func(d *NonResourceAttributesDie)) *SubjectAccessReviewSpecDie {
	return d.DieStamp(func(r *authorizationv1.SubjectAccessReviewSpec) {
		d := NonResourceAttributesBlank.DieImmutable(false).DieFeedPtr(r.NonResourceAttributes)
		fn(d)
		r.NonResourceAttributes = d.DieReleasePtr()
	})
}

// Extra corresponds to the user.Info.GetExtra() method from the authenticator.  Since that is input to the authorizer it needs a reflection here.
func (d *SubjectAccessReviewSpecDie) Extra(v map[string]authorizationv1.ExtraValue) *SubjectAccessReviewSpecDie {
	return d.DieStamp(func(r *authorizationv1.SubjectAccessReviewSpec) {
		r.Extra = v
	})
}

func (d *SubjectAccessReviewSpecDie) AddExtra(key string, value authorizationv1.ExtraValue) *SubjectAccessReviewSpecDie {
	return d.DieStamp(func(r *authorizationv1.SubjectAccessReviewSpec) {
		if r.Extra == nil {
			r.Extra = map[string]authorizationv1.ExtraValue{}
		}
		r.Extra[key] = value
	})
}

// +die
type _ = authorizationv1.ResourceAttributes

// +die
type _ = authorizationv1.NonResourceAttributes

// +die
type _ = authorizationv1.SubjectAccessReviewStatus
