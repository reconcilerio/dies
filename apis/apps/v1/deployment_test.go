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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	dieappsv1 "reconciler.io/dies/apis/apps/v1"
	diecorev1 "reconciler.io/dies/apis/core/v1"
	diemetav1 "reconciler.io/dies/apis/meta/v1"
)

func TestDeployment(t *testing.T) {
	now := metav1.Time{
		Time: time.Now().Round(time.Second),
	}

	tests := []struct {
		name     string
		die      *dieappsv1.DeploymentDie
		expected appsv1.Deployment
	}{
		{
			name:     "empty",
			die:      dieappsv1.DeploymentBlank,
			expected: appsv1.Deployment{},
		},
		{
			name: "object metadata",
			die: dieappsv1.DeploymentBlank.
				MetadataDie(func(d *diemetav1.ObjectMetaDie) {
					d.Namespace("my-namespace")
					d.Name("my-name")
				}),
			expected: appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "my-namespace",
					Name:      "my-name",
				},
			},
		},
		{
			name: "spec template die",
			die: dieappsv1.DeploymentBlank.
				SpecDie(func(d *dieappsv1.DeploymentSpecDie) {
					d.TemplateDie(func(d *diecorev1.PodTemplateSpecDie) {
						d.MetadataDie(func(d *diemetav1.ObjectMetaDie) {
							d.Name("my-name")
						})
						d.SpecDie(func(d *diecorev1.PodSpecDie) {
							d.ContainerDie("app", func(d *diecorev1.ContainerDie) {
								d.Image("registry.example/image:latest")
								d.Command("/executable")
							})
						})
					})
				}),
			expected: appsv1.Deployment{
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Name: "my-name",
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:    "app",
									Image:   "registry.example/image:latest",
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
			die: dieappsv1.DeploymentBlank.
				StatusDie(func(d *dieappsv1.DeploymentStatusDie) {
					d.ConditionsDie(
						diemetav1.ConditionBlank.
							Type(string(appsv1.DeploymentAvailable)).
							Status(metav1.ConditionTrue).
							Reason("MinimumReplicasAvailable").
							Message("Deployment has minimum availability.").
							LastTransitionTime(now),
						diemetav1.ConditionBlank.
							Type(string(appsv1.DeploymentProgressing)).
							Status(metav1.ConditionTrue).
							Reason("NewReplicaSetAvailable").
							Message(`ReplicaSet "test" has successfully progressed.`).
							LastTransitionTime(now),
					)
				}),
			expected: appsv1.Deployment{
				Status: appsv1.DeploymentStatus{
					Conditions: []appsv1.DeploymentCondition{
						{
							Type:               appsv1.DeploymentAvailable,
							Status:             corev1.ConditionTrue,
							Reason:             "MinimumReplicasAvailable",
							Message:            "Deployment has minimum availability.",
							LastTransitionTime: now,
						},
						{
							Type:               appsv1.DeploymentProgressing,
							Status:             corev1.ConditionTrue,
							Reason:             "NewReplicaSetAvailable",
							Message:            `ReplicaSet "test" has successfully progressed.`,
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
