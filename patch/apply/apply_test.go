/*
Copyright 2024 the original author or authors.

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

package apply_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	dieappsv1 "reconciler.io/dies/apis/apps/v1"
	"reconciler.io/dies/patch"
)

func TestCreate(t *testing.T) {
	original := dieappsv1.DeploymentBlank.
		SpecDie(func(d *dieappsv1.DeploymentSpecDie) {
			d.Replicas(ptr.To(int32(1)))
		})
	modified := original.
		SpecDie(func(d *dieappsv1.DeploymentSpecDie) {
			d.Replicas(ptr.To(int32(2)))
		})

	tests := map[string]struct {
		original      interface{}
		modified      interface{}
		patchType     types.PatchType
		expectedPatch string
		shouldErr     bool
	}{
		"identity ApplyPatchType": {
			original:      original.DieRelease(),
			modified:      original.DieRelease(),
			patchType:     types.ApplyPatchType,
			expectedPatch: `{"metadata":{},"spec":{"replicas":1,"selector":null,"template":{"metadata":{},"spec":{"containers":null}},"strategy":{}},"status":{}}`,
		},
		"diff ApplyPatchType": {
			original:      original.DieRelease(),
			modified:      modified.DieRelease(),
			patchType:     types.ApplyPatchType,
			expectedPatch: `{"metadata":{},"spec":{"replicas":2,"selector":null,"template":{"metadata":{},"spec":{"containers":null}},"strategy":{}},"status":{}}`,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			b, err := patch.Create(tc.original, tc.modified, tc.patchType)
			if err == nil && tc.shouldErr {
				t.Errorf("expected to err")
				return
			}
			if err != nil && !tc.shouldErr {
				t.Errorf("unexpected err: %s", err)
			}
			if diff := cmp.Diff(tc.expectedPatch, string(b)); diff != "" {
				t.Errorf("patch bytes (-expected, +actual): %s", diff)
			}
		})
	}
}
