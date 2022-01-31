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
	diemetav1 "dies.dev/apis/meta/v1"
	certificatesv1 "k8s.io/api/certificates/v1"
	corev1 "k8s.io/api/core/v1"
)

// +die:object=true
type _ = certificatesv1.CertificateSigningRequest

// TODO(scothis) fix import for maps with struct values, ignore the 'Extra' field until then

// +die:ignore={Extra}
type _ = certificatesv1.CertificateSigningRequestSpec

// Extra corresponds to the user.Info.GetExtra() method from the authenticator.  Since that is input to the authorizer it needs a reflection here.
func (d *CertificateSigningRequestSpecDie) Extra(v map[string]certificatesv1.ExtraValue) *CertificateSigningRequestSpecDie {
	return d.DieStamp(func(r *certificatesv1.CertificateSigningRequestSpec) {
		r.Extra = v
	})
}

func (d *CertificateSigningRequestSpecDie) AddExtra(key string, value certificatesv1.ExtraValue) *CertificateSigningRequestSpecDie {
	return d.DieStamp(func(r *certificatesv1.CertificateSigningRequestSpec) {
		if r.Extra == nil {
			r.Extra = map[string]certificatesv1.ExtraValue{}
		}
		r.Extra[key] = value
	})
}

// +die
type _ = certificatesv1.CertificateSigningRequestStatus

func (d *CertificateSigningRequestStatusDie) ConditionsDie(conditions ...*diemetav1.ConditionDie) *CertificateSigningRequestStatusDie {
	return d.DieStamp(func(r *certificatesv1.CertificateSigningRequestStatus) {
		r.Conditions = make([]certificatesv1.CertificateSigningRequestCondition, len(conditions))
		for i := range conditions {
			c := conditions[i].DieRelease()
			r.Conditions[i] = certificatesv1.CertificateSigningRequestCondition{
				Type:               certificatesv1.RequestConditionType(c.Type),
				Status:             corev1.ConditionStatus(c.Status),
				Reason:             c.Reason,
				Message:            c.Message,
				LastTransitionTime: c.LastTransitionTime,
			}
		}
	})
}
