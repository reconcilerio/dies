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

	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	diemetav1 "reconciler.io/dies/apis/meta/v1"
)

func TestCondition(t *testing.T) {
	tests := []struct {
		name     string
		die      *diemetav1.ConditionDie
		expected metav1.Condition
	}{
		{
			name:     "empty",
			die:      diemetav1.ConditionBlank,
			expected: metav1.Condition{},
		},
		{
			name:     "true",
			die:      diemetav1.ConditionBlank.True(),
			expected: metav1.Condition{Status: metav1.ConditionTrue},
		},
		{
			name:     "false",
			die:      diemetav1.ConditionBlank.False(),
			expected: metav1.Condition{Status: metav1.ConditionFalse},
		},
		{
			name:     "unknowm",
			die:      diemetav1.ConditionBlank.Unknown(),
			expected: metav1.Condition{Status: metav1.ConditionUnknown},
		},
		{
			name:     "messagef",
			die:      diemetav1.ConditionBlank.Messagef("%d coconuts", 10),
			expected: metav1.Condition{Message: "10 coconuts"},
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
