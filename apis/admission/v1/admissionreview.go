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
	admissionv1 "k8s.io/api/admission/v1"
)

// +die
type _ = admissionv1.AdmissionReview

func (d *AdmissionReviewDie) RequestDie(fn func(d *AdmissionRequestDie)) *AdmissionReviewDie {
	return d.DieStamp(func(r *admissionv1.AdmissionReview) {
		d := AdmissionRequestBlank.DieImmutable(false).DieFeedPtr(r.Request)
		fn(d)
		r.Request = d.DieReleasePtr()
	})
}

func (d *AdmissionReviewDie) ResponseDie(fn func(d *AdmissionResponseDie)) *AdmissionReviewDie {
	return d.DieStamp(func(r *admissionv1.AdmissionReview) {
		d := AdmissionResponseBlank.DieImmutable(false).DieFeedPtr(r.Response)
		fn(d)
		r.Response = d.DieReleasePtr()
	})
}
