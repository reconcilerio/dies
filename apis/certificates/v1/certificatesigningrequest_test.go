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

package v1_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	certificatesv1 "k8s.io/api/certificates/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	diecertificatesv1 "reconciler.io/dies/apis/certificates/v1"
	diemetav1 "reconciler.io/dies/apis/meta/v1"
)

func TestCertificatSigningRequest(t *testing.T) {
	tests := []struct {
		name     string
		die      *diecertificatesv1.CertificateSigningRequestDie
		expected certificatesv1.CertificateSigningRequest
	}{
		{
			name:     "empty",
			die:      diecertificatesv1.CertificateSigningRequestBlank,
			expected: certificatesv1.CertificateSigningRequest{},
		},
		{
			name: "object metadata",
			die: diecertificatesv1.CertificateSigningRequestBlank.
				MetadataDie(func(d *diemetav1.ObjectMetaDie) {
					d.Name("my-name")
				}),
			expected: certificatesv1.CertificateSigningRequest{
				ObjectMeta: metav1.ObjectMeta{
					Name: "my-name",
				},
			},
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			actual := c.die.DieRelease()
			if diff := cmp.Diff(c.expected, actual); diff != "" {
				t.Errorf("(-expected, +actual): %s", diff)
			}
		})
	}
}
