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
	dieauthenticationv1 "dies.dev/apis/authentication/v1"
	diemetav1 "dies.dev/apis/meta/v1"
	admissionv1 "k8s.io/api/admission/v1"
)

// +die
type _ = admissionv1.AdmissionRequest

func (d *AdmissionRequestDie) KindDie(fn func(d *diemetav1.GroupVersionKindDie)) *AdmissionRequestDie {
	return d.DieStamp(func(r *admissionv1.AdmissionRequest) {
		d := diemetav1.GroupVersionKindBlank.DieImmutable(false).DieFeed(r.Kind)
		fn(d)
		r.Kind = d.DieRelease()
	})
}

func (d *AdmissionRequestDie) ResourceDie(fn func(d *diemetav1.GroupVersionResourceDie)) *AdmissionRequestDie {
	return d.DieStamp(func(r *admissionv1.AdmissionRequest) {
		d := diemetav1.GroupVersionResourceBlank.DieImmutable(false).DieFeed(r.Resource)
		fn(d)
		r.Resource = d.DieRelease()
	})
}

func (d *AdmissionRequestDie) RequestKindDie(fn func(d *diemetav1.GroupVersionKindDie)) *AdmissionRequestDie {
	return d.DieStamp(func(r *admissionv1.AdmissionRequest) {
		d := diemetav1.GroupVersionKindBlank.DieImmutable(false).DieFeedPtr(r.RequestKind)
		fn(d)
		r.RequestKind = d.DieReleasePtr()
	})
}

func (d *AdmissionRequestDie) RequestResourceDie(fn func(d *diemetav1.GroupVersionResourceDie)) *AdmissionRequestDie {
	return d.DieStamp(func(r *admissionv1.AdmissionRequest) {
		d := diemetav1.GroupVersionResourceBlank.DieImmutable(false).DieFeedPtr(r.RequestResource)
		fn(d)
		r.RequestResource = d.DieReleasePtr()
	})
}

func (d *AdmissionRequestDie) UserInfoDie(fn func(d *dieauthenticationv1.UserInfoDie)) *AdmissionRequestDie {
	return d.DieStamp(func(r *admissionv1.AdmissionRequest) {
		d := dieauthenticationv1.UserInfoBlank.DieImmutable(false).DieFeed(r.UserInfo)
		fn(d)
		r.UserInfo = d.DieRelease()
	})
}
