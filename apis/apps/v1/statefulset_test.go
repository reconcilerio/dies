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

	dieappsv1 "dies.dev/apis/apps/v1"
	diecorev1 "dies.dev/apis/core/v1"
	diemetav1 "dies.dev/apis/meta/v1"
	"github.com/google/go-cmp/cmp"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestStatefulSet(t *testing.T) {
	now := metav1.Time{
		Time: time.Now().Round(time.Second),
	}

	tests := []struct {
		name     string
		die      dieappsv1.StatefulSetDie
		expected appsv1.StatefulSet
	}{
		{
			name:     "empty",
			die:      dieappsv1.StatefulSetBlank,
			expected: appsv1.StatefulSet{},
		},
		{
			name: "object metadata",
			die: dieappsv1.StatefulSetBlank.
				MetadataDie(func(d diemetav1.ObjectMetaDie) {
					d.Namespace("my-namespace")
					d.Name("my-name")
				}),
			expected: appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "my-namespace",
					Name:      "my-name",
				},
			},
		},
		{
			name: "spec template die",
			die: dieappsv1.StatefulSetBlank.
				SpecDie(func(d dieappsv1.StatefulSetSpecDie) {
					d.TemplateDie(func(d diecorev1.PodTemplateSpecDie) {
						d.MetadataDie(func(d diemetav1.ObjectMetaDie) {
							d.Name("my-name")
						})
						d.SpecDie(func(d diecorev1.PodSpecDie) {
							d.ContainerDie("app", func(d diecorev1.ContainerDie) {
								d.Command("/executable")
							})
						})
					})
				}),
			expected: appsv1.StatefulSet{
				Spec: appsv1.StatefulSetSpec{
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
			name: "spec volume claim templates die",
			die: dieappsv1.StatefulSetBlank.
				SpecDie(func(d dieappsv1.StatefulSetSpecDie) {
					d.VolumeClaimTemplatesDie(
						diecorev1.PersistentVolumeClaimBlank.
							SpecDie(func(d diecorev1.PersistentVolumeClaimSpecDie) {
								d.VolumeName("my-volume")
							}),
					)
				}),
			expected: appsv1.StatefulSet{
				Spec: appsv1.StatefulSetSpec{
					VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
						{
							Spec: corev1.PersistentVolumeClaimSpec{
								VolumeName: "my-volume",
							},
						},
					},
				},
			},
		},
		{
			name: "status conditions die",
			die: dieappsv1.StatefulSetBlank.
				StatusDie(func(d dieappsv1.StatefulSetStatusDie) {
					d.ConditionsDie(
						diemetav1.ConditionBlank.
							Type("Ready").
							Status(metav1.ConditionTrue).
							Reason("TheReason").
							Message("a message.").
							LastTransitionTime(now),
					)
				}),
			expected: appsv1.StatefulSet{
				Status: appsv1.StatefulSetStatus{
					Conditions: []appsv1.StatefulSetCondition{
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
