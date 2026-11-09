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
	"reconciler.io/dies/v2/die/schemes"
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
				s.AddKnownTypeWithName(gv.WithKind("Deployment"), &dieappsv1.DeploymentDieMut{})
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
				s.AddKnownTypeWithName(gv.WithKind("Widget"), &dieappsv1.DeploymentDieMut{})
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

func TestCheckRegisteredCategories(t *testing.T) {
	gv := dieappsv1.SchemeGroupVersion
	if err := codegen.CheckRegisteredCategories(dieappsv1.SchemeBuilder, upstreamappsv1.AddToScheme, gv); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	// a builder that registers a category as another
	as := func(c schemes.Category) schemes.AddFunc {
		return func(s *runtime.Scheme, _ schemes.Category) error { return dieappsv1.SchemeBuilder.Add(s, c) }
	}
	for name, tc := range map[string]struct {
		b    schemes.Builder
		want []string
	}{
		"the mutable dies as each": {
			b:    schemes.NewBuilder(as(schemes.DieMuts)),
			want: []string{"Dies: the dies of apps/v1 register kinds as other than dies of the category: [ControllerRevision DaemonSet", "Types: the types of apps/v1 don't register as upstream does"},
		},
		"the types as each": {
			b:    schemes.NewBuilder(as(schemes.Types)),
			want: []string{"Dies: the dies of apps/v1 don't register as upstream does", "DieMuts: the dies of apps/v1 don't register as upstream does"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			err := codegen.CheckRegisteredCategories(tc.b, upstreamappsv1.AddToScheme, gv)
			for _, w := range tc.want {
				if err == nil || !strings.Contains(err.Error(), w) {
					t.Errorf("error %v, want it to contain %q", err, w)
				}
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
