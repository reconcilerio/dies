/*
Copyright 2026 the original author or authors.

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

	corev1 "reconciler.io/dies/v2/apis/core/v1"
	dietesting "reconciler.io/dies/v2/die/testing"
)

// The hand-written VolumeSource alias declares a union, which applies to the
// members flattened into Volume.
func TestVolumeSourceUnionOverride(t *testing.T) {
	v := corev1.VolumeBlank.DieStamp(func(d *corev1.VolumeDieMut) {
		dietesting.NoError(t, d.DieFeed(&corev1.Volume{
			Name:         "creds",
			VolumeSource: corev1.VolumeSource{Secret: &corev1.SecretVolumeSource{SecretName: "creds"}},
		}))
	}).DieStamp(func(d *corev1.VolumeDieMut) {
		dietesting.NoError(t, d.EmptyDir().Medium().Set(corev1.StorageMediumMemory))
	})

	got := *v.DieRelease()
	if got.Secret != nil {
		t.Errorf("setting emptyDir should clear secret, got %+v", got.Secret)
	}
	if got.EmptyDir == nil || got.EmptyDir.Medium != corev1.StorageMediumMemory {
		t.Errorf("emptyDir %+v", got.EmptyDir)
	}
	if want := "clear .secret\nset .emptyDir.medium = \"Memory\""; v.DieIntent() != want {
		t.Errorf("intent\n%s\nwant\n%s", v.DieIntent(), want)
	}
}
