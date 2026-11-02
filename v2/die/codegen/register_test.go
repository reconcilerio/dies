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

package codegen_test

import (
	"slices"
	"strings"
	"testing"

	upstreamappsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	"reconciler.io/dies/v2/die/base"
	"reconciler.io/dies/v2/die/codegen"
)

func TestCheckRegistered(t *testing.T) {
	gv := dieappsv1.SchemeGroupVersion
	if err := codegen.CheckRegistered(dieappsv1.AddToScheme, upstreamappsv1.AddToScheme, gv); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	for name, tc := range map[string]struct {
		dies func(*runtime.Scheme) error
		want string
	}{
		"a kind without a die": {
			dies: func(s *runtime.Scheme) error {
				s.AddKnownTypeWithName(gv.WithKind("Deployment"), &dieappsv1.DeploymentDie{})
				metav1.AddToGroupVersion(s, gv)
				return nil
			},
			want: "kinds without dies [ControllerRevision ControllerRevisionList DaemonSet",
		},
		"the upstream types": {
			dies: upstreamappsv1.AddToScheme,
			want: "kinds registered as the upstream type [ControllerRevision ControllerRevisionList DaemonSet",
		},
		"a kind upstream doesn't have": {
			dies: func(s *runtime.Scheme) error {
				if err := dieappsv1.AddToScheme(s); err != nil {
					return err
				}
				s.AddKnownTypeWithName(gv.WithKind("Widget"), &dieappsv1.DeploymentDie{})
				return nil
			},
			want: "kinds upstream doesn't have [Widget]",
		},
	} {
		t.Run(name, func(t *testing.T) {
			err := codegen.CheckRegistered(tc.dies, upstreamappsv1.AddToScheme, gv)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Errorf("error %v, want %q", err, tc.want)
			}
		})
	}
}

// Field tables are built on first use, and asking for the fields a die
// manages is a use.
func TestManagedFieldsOfAnUnusedDie(t *testing.T) {
	var wires []string
	for _, f := range base.ManagedFields[upstreamappsv1.ControllerRevision]() {
		wires = append(wires, f.Wire)
	}
	if want := []string{"kind", "apiVersion", "metadata", "data", "revision"}; !slices.Equal(wires, want) {
		t.Errorf("fields %v, want %v", wires, want)
	}
}
