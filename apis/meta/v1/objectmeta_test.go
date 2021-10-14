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
	diecorev1 "github.com/scothis/dies/apis/core/v1"
	diemetav1 "github.com/scothis/dies/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/utils/pointer"
)

func TestObjectMeta(t *testing.T) {
	scheme := runtime.NewScheme()
	utilruntime.Must(corev1.AddToScheme(scheme))

	tests := []struct {
		name     string
		die      *diemetav1.ObjectMetaDie
		expected metav1.ObjectMeta
	}{
		{
			name:     "empty",
			die:      diemetav1.ObjectMetaBlank,
			expected: metav1.ObjectMeta{},
		},
		{
			name: "add label",
			die: diemetav1.ObjectMetaBlank.
				AddLabel("key", "value"),
			expected: metav1.ObjectMeta{
				Labels: map[string]string{
					"key": "value",
				},
			},
		},
		{
			name: "add annotation",
			die: diemetav1.ObjectMetaBlank.
				AddAnnotation("key", "value"),
			expected: metav1.ObjectMeta{
				Annotations: map[string]string{
					"key": "value",
				},
			},
		},
		{
			name: "controlled by",
			die: diemetav1.ObjectMetaBlank.
				ControlledBy(
					diecorev1.PodBlank.
						MetadataDie(func(d *diemetav1.ObjectMetaDie) {
							d.Name("my-name")
							d.UID("123e4567-e89b-12d3-a456-426614174000")
						}),
					scheme,
				),
			expected: metav1.ObjectMeta{
				OwnerReferences: []metav1.OwnerReference{
					{
						APIVersion:         "v1",
						Kind:               "Pod",
						Name:               "my-name",
						UID:                "123e4567-e89b-12d3-a456-426614174000",
						BlockOwnerDeletion: pointer.Bool(true),
						Controller:         pointer.Bool(true),
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

func TestFrozenObjectMeta(t *testing.T) {
	timestamp := metav1.Date(2021, 10, 10, 12, 0, 0, 0, time.UTC)
	metadata := metav1.ObjectMeta{
		Name:                       "my-name",
		GenerateName:               "my-name-",
		Namespace:                  "my-namespace",
		Generation:                 1,
		SelfLink:                   "http://localhost/resource",
		UID:                        "123e4567-e89b-12d3-a456-426614174000",
		ResourceVersion:            "1",
		CreationTimestamp:          timestamp,
		DeletionTimestamp:          &timestamp,
		DeletionGracePeriodSeconds: pointer.Int64(30),
		Labels: map[string]string{
			"key": "value",
		},
		Annotations: map[string]string{
			"key": "value",
		},
		Finalizers: []string{"finalizer1", "finalizer2"},
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: "v1",
				Kind:       "Foo",
				Name:       "my-name",
				UID:        "123e4567-e89b-12d3-a456-426614174000",
			},
		},
		ClusterName: "my-cluster",
		ManagedFields: []metav1.ManagedFieldsEntry{
			{
				Manager:    "kubectl",
				Operation:  "apply",
				APIVersion: "v1",
			},
		},
	}
	frozen := diemetav1.FreezeObjectMeta(metadata)

	if &frozen != frozen.GetObjectMeta() {
		t.Errorf("FrozenObjectMeta#GetObjectMeta")
	}
	if diff := cmp.Diff(metadata.Name, frozen.GetName()); diff != "" {
		t.Errorf("FrozenObjectMeta#GetName (-expected, +actual): %s", diff)
	}
	if diff := cmp.Diff(metadata.GenerateName, frozen.GetGenerateName()); diff != "" {
		t.Errorf("FrozenObjectMeta#GetGenerateName (-expected, +actual): %s", diff)
	}
	if diff := cmp.Diff(metadata.Namespace, frozen.GetNamespace()); diff != "" {
		t.Errorf("FrozenObjectMeta#GetNamespace (-expected, +actual): %s", diff)
	}
	if diff := cmp.Diff(metadata.Generation, frozen.GetGeneration()); diff != "" {
		t.Errorf("FrozenObjectMeta#GetGeneration (-expected, +actual): %s", diff)
	}
	if diff := cmp.Diff(metadata.SelfLink, frozen.GetSelfLink()); diff != "" {
		t.Errorf("FrozenObjectMeta#GetSelfLink (-expected, +actual): %s", diff)
	}
	if diff := cmp.Diff(metadata.UID, frozen.GetUID()); diff != "" {
		t.Errorf("FrozenObjectMeta#GetUID (-expected, +actual): %s", diff)
	}
	if diff := cmp.Diff(metadata.ResourceVersion, frozen.GetResourceVersion()); diff != "" {
		t.Errorf("FrozenObjectMeta#GetResourceVersion (-expected, +actual): %s", diff)
	}
	if diff := cmp.Diff(metadata.CreationTimestamp, frozen.GetCreationTimestamp()); diff != "" {
		t.Errorf("FrozenObjectMeta#GetCreationTimestamp (-expected, +actual): %s", diff)
	}
	if diff := cmp.Diff(metadata.DeletionTimestamp, frozen.GetDeletionTimestamp()); diff != "" {
		t.Errorf("FrozenObjectMeta#GetDeletionTimestamp (-expected, +actual): %s", diff)
	}
	if diff := cmp.Diff(metadata.DeletionGracePeriodSeconds, frozen.GetDeletionGracePeriodSeconds()); diff != "" {
		t.Errorf("FrozenObjectMeta#GetDeletionGracePeriodSeconds (-expected, +actual): %s", diff)
	}
	if diff := cmp.Diff(metadata.Labels, frozen.GetLabels()); diff != "" {
		t.Errorf("FrozenObjectMeta#GetLabels (-expected, +actual): %s", diff)
	}
	if diff := cmp.Diff(metadata.Annotations, frozen.GetAnnotations()); diff != "" {
		t.Errorf("FrozenObjectMeta#GetAnnotations (-expected, +actual): %s", diff)
	}
	if diff := cmp.Diff(metadata.Finalizers, frozen.GetFinalizers()); diff != "" {
		t.Errorf("FrozenObjectMeta#GetFinalizers (-expected, +actual): %s", diff)
	}
	if diff := cmp.Diff(metadata.OwnerReferences, frozen.GetOwnerReferences()); diff != "" {
		t.Errorf("FrozenObjectMeta#GetOwnerReferences (-expected, +actual): %s", diff)
	}
	if diff := cmp.Diff(metadata.ClusterName, frozen.GetClusterName()); diff != "" {
		t.Errorf("FrozenObjectMeta#GetClusterName (-expected, +actual): %s", diff)
	}
	if diff := cmp.Diff(metadata.ManagedFields, frozen.GetManagedFields()); diff != "" {
		t.Errorf("FrozenObjectMeta#GetManagedFields (-expected, +actual): %s", diff)
	}

	tests := []struct {
		name        string
		shouldPanic func()
	}{
		{
			name:        "set name",
			shouldPanic: func() { frozen.SetName(metadata.Name) },
		},
		{
			name:        "set generate name",
			shouldPanic: func() { frozen.SetGenerateName(metadata.GenerateName) },
		},
		{
			name:        "set namespace",
			shouldPanic: func() { frozen.SetNamespace(metadata.Namespace) },
		},
		{
			name:        "set generation",
			shouldPanic: func() { frozen.SetGeneration(metadata.Generation) },
		},
		{
			name:        "set self link",
			shouldPanic: func() { frozen.SetSelfLink(metadata.SelfLink) },
		},
		{
			name:        "set uid",
			shouldPanic: func() { frozen.SetUID(metadata.UID) },
		},
		{
			name:        "set resource version",
			shouldPanic: func() { frozen.SetResourceVersion(metadata.ResourceVersion) },
		},
		{
			name:        "set creation timestamp",
			shouldPanic: func() { frozen.SetCreationTimestamp(metadata.CreationTimestamp) },
		},
		{
			name:        "set deletion timestamp",
			shouldPanic: func() { frozen.SetDeletionTimestamp(metadata.DeletionTimestamp) },
		},
		{
			name:        "set deletion grace period seconds",
			shouldPanic: func() { frozen.SetDeletionGracePeriodSeconds(metadata.DeletionGracePeriodSeconds) },
		},
		{
			name:        "set labels",
			shouldPanic: func() { frozen.SetLabels(metadata.Labels) },
		},
		{
			name:        "set annotations",
			shouldPanic: func() { frozen.SetAnnotations(metadata.Annotations) },
		},
		{
			name:        "set finalizers",
			shouldPanic: func() { frozen.SetFinalizers(metadata.Finalizers) },
		},
		{
			name:        "set owner references",
			shouldPanic: func() { frozen.SetOwnerReferences(metadata.OwnerReferences) },
		},
		{
			name:        "set cluster name",
			shouldPanic: func() { frozen.SetClusterName(metadata.ClusterName) },
		},
		{
			name:        "set managed fields",
			shouldPanic: func() { frozen.SetManagedFields(metadata.ManagedFields) },
		},
	}
	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			defer func() {
				if recover() != nil {
					// we want this
					_ = ""
				}
			}()
			c.shouldPanic()
			t.Errorf("expected to panic")
		})
	}

}
