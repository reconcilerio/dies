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

	diecorev1 "dies.dev/apis/core/v1"
	diemetav1 "dies.dev/apis/meta/v1"
	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

func TestPod(t *testing.T) {
	now := metav1.Time{
		Time: time.Now().Round(time.Second),
	}

	tests := []struct {
		name     string
		die      *diecorev1.PodDie
		expected corev1.Pod
	}{
		{
			name:     "empty",
			die:      diecorev1.PodBlank,
			expected: corev1.Pod{},
		},
		{
			name: "object metadata",
			die: diecorev1.PodBlank.
				MetadataDie(func(d *diemetav1.ObjectMetaDie) {
					d.Namespace("my-namespace")
					d.Name("my-name")
				}),
			expected: corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "my-namespace",
					Name:      "my-name",
				},
			},
		},
		{
			name: "spec add container",
			die: diecorev1.PodBlank.
				SpecDie(func(d *diecorev1.PodSpecDie) {
					d.ContainerDie("workload", func(d *diecorev1.ContainerDie) {
						d.Image("ubuntu:bionic")
					})
				}),
			expected: corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "workload",
							Image: "ubuntu:bionic",
						},
					},
				},
			},
		},
		{
			name: "spec multiple containers",
			die: diecorev1.PodBlank.
				SpecDie(func(d *diecorev1.PodSpecDie) {
					d.ContainerDie("workload", func(d *diecorev1.ContainerDie) {
						d.Image("ubuntu:bionic")
					}).
						ContainerDie("sidecar", func(d *diecorev1.ContainerDie) {
							d.Image("gcr.io/kubebuilder/kube-rbac-proxy")
						})
				}),
			expected: corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "workload",
							Image: "ubuntu:bionic",
						},
						{
							Name:  "sidecar",
							Image: "gcr.io/kubebuilder/kube-rbac-proxy",
						},
					},
				},
			},
		},
		{
			name: "spec update containers",
			die: diecorev1.PodBlank.
				SpecDie(func(d *diecorev1.PodSpecDie) {
					d.ContainerDie("workload", func(d *diecorev1.ContainerDie) {
						d.Image("ubuntu:bionic")
					})
					d.ContainerDie("workload", func(d *diecorev1.ContainerDie) {
						d.Command("env")
					})
				}),
			expected: corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "workload",
							Image:   "ubuntu:bionic",
							Command: []string{"env"},
						},
					},
				},
			},
		},
		{
			name: "spec add init container",
			die: diecorev1.PodBlank.
				SpecDie(func(d *diecorev1.PodSpecDie) {
					d.InitContainerDie("workload", func(d *diecorev1.ContainerDie) {
						d.Image("ubuntu:bionic")
					})
				}),
			expected: corev1.Pod{
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						{
							Name:  "workload",
							Image: "ubuntu:bionic",
						},
					},
				},
			},
		},
		{
			name: "spec multiple init containers",
			die: diecorev1.PodBlank.
				SpecDie(func(d *diecorev1.PodSpecDie) {
					d.InitContainerDie("workload", func(d *diecorev1.ContainerDie) {
						d.Image("ubuntu:bionic")
					})
					d.InitContainerDie("sidecar", func(d *diecorev1.ContainerDie) {
						d.Image("gcr.io/kubebuilder/kube-rbac-proxy")
					})
				}),
			expected: corev1.Pod{
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						{
							Name:  "workload",
							Image: "ubuntu:bionic",
						},
						{
							Name:  "sidecar",
							Image: "gcr.io/kubebuilder/kube-rbac-proxy",
						},
					},
				},
			},
		},
		{
			name: "spec update init containers",
			die: diecorev1.PodBlank.
				SpecDie(func(d *diecorev1.PodSpecDie) {
					d.InitContainerDie("workload", func(d *diecorev1.ContainerDie) {
						d.Image("ubuntu:bionic")
					})
					d.InitContainerDie("workload", func(d *diecorev1.ContainerDie) {
						d.Command("env")
					})
				}),
			expected: corev1.Pod{
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						{
							Name:    "workload",
							Image:   "ubuntu:bionic",
							Command: []string{"env"},
						},
					},
				},
			},
		},
		{
			name: "spec add volume",
			die: diecorev1.PodBlank.
				SpecDie(func(d *diecorev1.PodSpecDie) {
					d.VolumeDie("config", func(d *diecorev1.VolumeDie) {
						d.ConfigMapDie(func(d *diecorev1.ConfigMapVolumeSourceDie) {
							d.LocalObjectReference(corev1.LocalObjectReference{
								Name: "my-config",
							})
						})
					})
				}),
			expected: corev1.Pod{
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "my-config",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "spec multiple volumes",
			die: diecorev1.PodBlank.
				SpecDie(func(d *diecorev1.PodSpecDie) {
					d.VolumeDie("config", func(d *diecorev1.VolumeDie) {
						d.ConfigMapDie(func(d *diecorev1.ConfigMapVolumeSourceDie) {
							d.LocalObjectReference(corev1.LocalObjectReference{
								Name: "my-config",
							})
						})
					})
					d.VolumeDie("scratch", func(d *diecorev1.VolumeDie) {
						d.EmptyDirDie(func(d *diecorev1.EmptyDirVolumeSourceDie) {
							d.Medium(corev1.StorageMediumMemory)
						})
					})
				}),
			expected: corev1.Pod{
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "my-config",
									},
								},
							},
						},
						{
							Name: "scratch",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{
									Medium: corev1.StorageMediumMemory,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "spec update volume",
			die: diecorev1.PodBlank.
				SpecDie(func(d *diecorev1.PodSpecDie) {
					d.VolumeDie("config", func(d *diecorev1.VolumeDie) {
						d.ConfigMapDie(func(d *diecorev1.ConfigMapVolumeSourceDie) {
							d.LocalObjectReference(corev1.LocalObjectReference{
								Name: "my-config",
							})
						})
					})
					d.VolumeDie("config", func(d *diecorev1.VolumeDie) {
						d.ConfigMapDie(func(d *diecorev1.ConfigMapVolumeSourceDie) {
							d.Optional(pointer.Bool(true))
						})
					})
				}),
			expected: corev1.Pod{
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "my-config",
									},
									Optional: pointer.Bool(true),
								},
							},
						},
					},
				},
			},
		},
		{
			name: "status condition",
			die: diecorev1.PodBlank.
				StatusDie(func(d *diecorev1.PodStatusDie) {
					d.ConditionsDie(
						diemetav1.ConditionBlank.
							Type("ContainersReady").
							Status(metav1.ConditionTrue).
							Reason("TheReason").
							Message("a message.").
							LastTransitionTime(now),
					)
				}),
			expected: corev1.Pod{
				Status: corev1.PodStatus{
					Conditions: []corev1.PodCondition{
						{
							Type:               corev1.ContainersReady,
							Status:             corev1.ConditionTrue,
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

func TestPodTemplateSpec(t *testing.T) {
	tests := []struct {
		name     string
		die      *diecorev1.PodTemplateSpecDie
		expected corev1.PodTemplateSpec
	}{
		{
			name:     "empty",
			die:      diecorev1.PodTemplateSpecBlank,
			expected: corev1.PodTemplateSpec{},
		},
		{
			name: "object metadata",
			die: diecorev1.PodTemplateSpecBlank.
				MetadataDie(func(d *diemetav1.ObjectMetaDie) {
					d.Namespace("my-namespace")
					d.Name("my-name")
				}),
			expected: corev1.PodTemplateSpec{
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
