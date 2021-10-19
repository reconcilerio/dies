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
	"time"

	dieapiextensionsv1 "dies.dev/apis/apiextensions/v1"
	diemetav1 "dies.dev/apis/meta/v1"
	"github.com/google/go-cmp/cmp"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCustomResourceDefinition(t *testing.T) {
	now := metav1.Time{
		Time: time.Now().Round(time.Second),
	}

	tests := []struct {
		name     string
		die      *dieapiextensionsv1.CustomResourceDefinitionDie
		expected apiextensionsv1.CustomResourceDefinition
	}{
		{
			name:     "empty",
			die:      dieapiextensionsv1.CustomResourceDefinitionBlank,
			expected: apiextensionsv1.CustomResourceDefinition{},
		},
		{
			name: "object metadata",
			die: dieapiextensionsv1.CustomResourceDefinitionBlank.
				MetadataDie(func(d *diemetav1.ObjectMetaDie) {
					d.Name("my-name")
				}),
			expected: apiextensionsv1.CustomResourceDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name: "my-name",
				},
			},
		},
		{
			name: "status conditions die",
			die: dieapiextensionsv1.CustomResourceDefinitionBlank.
				StatusDie(func(d *dieapiextensionsv1.CustomResourceDefinitionStatusDie) {
					d.ConditionsDie(
						diemetav1.ConditionBlank.
							Type(string(apiextensionsv1.Established)).
							Status(metav1.ConditionTrue).
							Reason("TheReason").
							Message("a message.").
							LastTransitionTime(now),
					)
				}),
			expected: apiextensionsv1.CustomResourceDefinition{
				Status: apiextensionsv1.CustomResourceDefinitionStatus{
					Conditions: []apiextensionsv1.CustomResourceDefinitionCondition{
						{
							Type:               "Established",
							Status:             "True",
							Reason:             "TheReason",
							Message:            "a message.",
							LastTransitionTime: now,
						},
					},
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
