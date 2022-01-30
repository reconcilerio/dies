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

	dieauthenticationv1 "dies.dev/apis/authentication/v1"
	diemetav1 "dies.dev/apis/meta/v1"
	"github.com/google/go-cmp/cmp"
	authenticationv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestTokenReview(t *testing.T) {
	tests := []struct {
		name     string
		die      *dieauthenticationv1.TokenReviewDie
		expected authenticationv1.TokenReview
	}{
		{
			name:     "empty",
			die:      dieauthenticationv1.TokenReviewBlank,
			expected: authenticationv1.TokenReview{},
		},
		{
			name: "object metadata",
			die: dieauthenticationv1.TokenReviewBlank.
				MetadataDie(func(d *diemetav1.ObjectMetaDie) {
					d.Name("my-name")
				}),
			expected: authenticationv1.TokenReview{
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
