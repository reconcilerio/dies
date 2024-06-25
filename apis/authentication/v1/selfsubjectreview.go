/*
Copyright 2024 the original author or authors.

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

// +die:object=true,apiVersion=authentication.k8s.io/v1,kind=SelfSubjectReview
type _ = authenticationv1.SelfSubjectReview

// +die
// +die:field:name=UserInfo,die=UserInfoDie
type _ = authenticationv1.SelfSubjectReviewStatus

// +die:ignore={Extra}
type _ = authenticationv1.UserInfo

// Any additional information provided by the authenticator.
func (d *UserInfoDie) Extra(v map[string]authenticationv1.ExtraValue) *UserInfoDie {
	return d.DieStamp(func(r *authenticationv1.UserInfo) {
		r.Extra = v
	})
}

func (d *UserInfoDie) AddExtra(key string, value ...string) *UserInfoDie {
	return d.DieStamp(func(r *authenticationv1.UserInfo) {
		if r.Extra == nil {
			r.Extra = map[string]authenticationv1.ExtraValue{}
		}
		r.Extra[key] = value
	})
}
