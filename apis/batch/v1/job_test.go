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

	"github.com/google/go-cmp/cmp"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	diebatchv1 "reconciler.io/dies/apis/batch/v1"
	diecorev1 "reconciler.io/dies/apis/core/v1"
	diemetav1 "reconciler.io/dies/apis/meta/v1"
)

func TestJob(t *testing.T) {
	now := metav1.Time{
		Time: time.Now().Round(time.Second),
	}

	tests := []struct {
		name     string
		die      *diebatchv1.JobDie
		expected batchv1.Job
	}{
		{
			name:     "empty",
			die:      diebatchv1.JobBlank,
			expected: batchv1.Job{},
		},
		{
			name: "object metadata",
			die: diebatchv1.JobBlank.
				MetadataDie(func(d *diemetav1.ObjectMetaDie) {
					d.Namespace("my-namespace")
					d.Name("my-name")
				}),
			expected: batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "my-namespace",
					Name:      "my-name",
				},
			},
		},
		{
			name: "spec template die",
			die: diebatchv1.JobBlank.
				SpecDie(func(d *diebatchv1.JobSpecDie) {
					d.TemplateDie(func(d *diecorev1.PodTemplateSpecDie) {
						d.MetadataDie(func(d *diemetav1.ObjectMetaDie) {
							d.Name("my-name")
						})
						d.SpecDie(func(d *diecorev1.PodSpecDie) {
							d.ContainerDie("app", func(d *diecorev1.ContainerDie) {
								d.Command("/executable")
							})
						})
					})
				}),
			expected: batchv1.Job{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Name: "my-name",
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:    "app",
									Command: []string{"/executable"},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "status conditions die",
			die: diebatchv1.JobBlank.
				StatusDie(func(d *diebatchv1.JobStatusDie) {
					d.ConditionsDie(
						diemetav1.ConditionBlank.
							Type("Ready").
							Status(metav1.ConditionTrue).
							Reason("TheReason").
							Message("a message.").
							LastTransitionTime(now),
					)
				}),
			expected: batchv1.Job{
				Status: batchv1.JobStatus{
					Conditions: []batchv1.JobCondition{
						{
							Type:               "Ready",
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
