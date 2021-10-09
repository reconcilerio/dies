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
	dieappsv1 "github.com/scothis/dies/apis/apps/v1"
	diecorev1 "github.com/scothis/dies/apis/core/v1"
	diemetav1 "github.com/scothis/dies/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestReplicaSet(t *testing.T) {
	tests := []struct {
		name     string
		die      *dieappsv1.ReplicaSetDie
		expected appsv1.ReplicaSet
	}{
		{
			name:     "empty",
			die:      dieappsv1.ReplicaSetBlank,
			expected: appsv1.ReplicaSet{},
		},
		{
			name: "object metadata",
			die: dieappsv1.ReplicaSetBlank.
				MetadataDie(func(d *diemetav1.ObjectMetaDie) {
					d.Namespace("my-namespace")
					d.Name("my-name")
				}),
			expected: appsv1.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "my-namespace",
					Name:      "my-name",
				},
			},
		},
		{
			name: "spec template die",
			die: dieappsv1.ReplicaSetBlank.
				SpecDie(func(d *dieappsv1.ReplicaSetSpecDie) {
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
			expected: appsv1.ReplicaSet{
				Spec: appsv1.ReplicaSetSpec{
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
			die: dieappsv1.ReplicaSetBlank.
				StatusDie(func(d *dieappsv1.ReplicaSetStatusDie) {
					d.ConditionsDie(
						diemetav1.ConditionBlank.
							Type("Ready").
							Status(metav1.ConditionTrue).
							Reason("TheReason").
							Message("a message.").
							LastTransitionTime(metav1.Time{Time: time.Unix(1, 0)}),
					)
				}),
			expected: appsv1.ReplicaSet{
				Status: appsv1.ReplicaSetStatus{
					Conditions: []appsv1.ReplicaSetCondition{
						{
							Type:               "Ready",
							Status:             "True",
							Reason:             "TheReason",
							Message:            "a message.",
							LastTransitionTime: metav1.Time{Time: time.Unix(1, 0)},
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
