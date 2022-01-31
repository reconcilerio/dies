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

package v1beta1_test

import (
	"testing"

	diemetav1 "dies.dev/apis/meta/v1"
	diepolicyv1beta1 "dies.dev/apis/policy/v1beta1"
	"github.com/google/go-cmp/cmp"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPodSecurityPolicy(t *testing.T) {
	tests := []struct {
		name     string
		die      *diepolicyv1beta1.PodSecurityPolicyDie
		expected policyv1beta1.PodSecurityPolicy
	}{
		{
			name:     "empty",
			die:      diepolicyv1beta1.PodSecurityPolicyBlank,
			expected: policyv1beta1.PodSecurityPolicy{},
		},
		{
			name: "object metadata",
			die: diepolicyv1beta1.PodSecurityPolicyBlank.
				MetadataDie(func(d *diemetav1.ObjectMetaDie) {
					d.Name("my-name")
				}),
			expected: policyv1beta1.PodSecurityPolicy{
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
