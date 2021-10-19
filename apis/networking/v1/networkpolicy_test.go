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

package v1_test

import (
	"testing"

	diemetav1 "dies.dev/apis/meta/v1"
	dienetworkingv1 "dies.dev/apis/networking/v1"
	"github.com/google/go-cmp/cmp"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNetworkPolicy(t *testing.T) {
	tests := []struct {
		name     string
		die      dienetworkingv1.NetworkPolicyDie
		expected networkingv1.NetworkPolicy
	}{
		{
			name:     "empty",
			die:      dienetworkingv1.NetworkPolicyBlank,
			expected: networkingv1.NetworkPolicy{},
		},
		{
			name: "object metadata",
			die: dienetworkingv1.NetworkPolicyBlank.
				MetadataDie(func(d diemetav1.ObjectMetaDie) {
					d.Namespace("my-namespace")
					d.Name("my-name")
				}),
			expected: networkingv1.NetworkPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "my-namespace",
					Name:      "my-name",
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
