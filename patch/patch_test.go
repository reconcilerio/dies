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

package patch_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/types"
	"reconciler.io/dies/patch"
)

func TestCreate(t *testing.T) {
	tests := map[string]struct {
		original      interface{}
		modified      interface{}
		patchType     types.PatchType
		expectedPatch string
		expectedErr   string
	}{
		"invalid patch type": {
			patchType:   types.PatchType("invalid"),
			expectedErr: fmt.Sprintf("unknown patch type %q", "invalid"),
		},
		"unloaded JSONPatch": {
			patchType:   types.JSONPatchType,
			expectedErr: fmt.Sprintf("unknown patch type %q, load support with `import _ %q`", types.JSONPatchType, "reconciler.io/dies/patch/json"),
		},
		"unloaded MergePatchType": {
			patchType:   types.MergePatchType,
			expectedErr: fmt.Sprintf("unknown patch type %q, load support with `import _ %q`", types.MergePatchType, "reconciler.io/dies/patch/merge"),
		},
		"unloaded StrategicMergePatchType": {
			patchType:   types.StrategicMergePatchType,
			expectedErr: fmt.Sprintf("unknown patch type %q, load support with `import _ %q`", types.StrategicMergePatchType, "reconciler.io/dies/patch/strategicmerge"),
		},
		"unloaded ApplyPatchType": {
			patchType:   types.ApplyPatchType,
			expectedErr: fmt.Sprintf("unknown patch type %q, load support with `import _ %q`", types.ApplyPatchType, "reconciler.io/dies/patch/apply"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := patch.Create(tc.original, tc.modified, tc.patchType)
			if err == nil && tc.expectedErr != "" {
				t.Errorf("expected to err")
				return
			}
			if err != nil {
				if diff := cmp.Diff(tc.expectedErr, err.Error()); diff != "" {
					t.Errorf("unexpected err (-expected, +actual): %s", diff)
				}
			}
		})
	}
}
