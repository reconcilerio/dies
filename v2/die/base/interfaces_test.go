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

package base_test

import (
	"encoding/json"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	diecorev1 "reconciler.io/dies/v2/apis/core/v1"
	"reconciler.io/dies/v2/die"
	"reconciler.io/dies/v2/die/base"
	"reconciler.io/dies/v2/die/lists"
	dieunstructured "reconciler.io/dies/v2/die/unstructured"
)

// Every die is a Die and a DieOf its types, and every mutable die a Die
// and a DieMutOf; typed dies are DieTypeds, object dies DieObjects, and list
// dies DieLists, of mutable items. An unstructured die is a DieObject, but
// not a DieObjectOf: it has no typed accessors.
var (
	_ die.Die                                                                                                       = (*diecorev1.PodSpecDie)(nil)
	_ die.Die                                                                                                       = (*diecorev1.PodSpecDieMut)(nil)
	_ base.DieOf[corev1.PodSpec, *diecorev1.PodSpecDie, *diecorev1.PodSpecDieMut]                                   = (*diecorev1.PodSpecDie)(nil)
	_ base.DieMutOf[corev1.PodSpec, *diecorev1.PodSpecDie, *diecorev1.PodSpecDieMut]                                = (*diecorev1.PodSpecDieMut)(nil)
	_ die.DieTyped                                                                                                  = (*diecorev1.PodAttachOptionsDie)(nil)
	_ die.DieTypedOf[corev1.PodAttachOptions, *diecorev1.PodAttachOptionsDie, *diecorev1.PodAttachOptionsDieMut]    = (*diecorev1.PodAttachOptionsDie)(nil)
	_ die.DieTypedMutOf[corev1.PodAttachOptions, *diecorev1.PodAttachOptionsDie, *diecorev1.PodAttachOptionsDieMut] = (*diecorev1.PodAttachOptionsDieMut)(nil)
	_ die.DieObject                                                                                                 = (*dieappsv1.DeploymentDie)(nil)
	_ die.DieObject                                                                                                 = (*dieappsv1.DeploymentDieMut)(nil)
	_ die.DieObjectOf[appsv1.Deployment, *dieappsv1.DeploymentDie, *dieappsv1.DeploymentDieMut]                     = (*dieappsv1.DeploymentDie)(nil)
	_ die.DieObjectMutOf[appsv1.Deployment, *dieappsv1.DeploymentDie, *dieappsv1.DeploymentDieMut]                  = (*dieappsv1.DeploymentDieMut)(nil)
	_ die.DieObject                                                                                                 = (*dieunstructured.Object)(nil)
	_ die.DieObject                                                                                                 = (*dieunstructured.ObjectMut)(nil)
	_ die.DieList                                                                                                   = (*dieappsv1.DeploymentListDie)(nil)
	_ die.DieListOf[*dieappsv1.DeploymentDieMut]                                                                    = (*dieappsv1.DeploymentListDie)(nil)
	_ die.DieList                                                                                                   = (*dieunstructured.ObjectList)(nil)
	_ lists.DieListOf[*dieunstructured.ObjectMut]                                                                   = (*dieunstructured.ObjectList)(nil)
)

// Generic code over DieOf calls a die's methods directly, and writes
// through its mutable die.
func stampName[T any, D die.DieObjectOf[T, D, M], M die.DieObjectMutOf[T, D, M]](d D, name string) D {
	return d.DieWith(func(m M) {
		m.Metadata().Name().Set(name)
		m.APIVersion().Set("apps/v1")
	})
}

// An immutable die has DieMutable, and a mutable die DieImmutable, never
// both, so the two are told apart by type: ImmutableDie and MutableDie.
func TestMutabilityIsExclusive(t *testing.T) {
	var d, m any = dieappsv1.DeploymentBlank, dieappsv1.DeploymentBlank.DieMutable()
	if _, ok := d.(base.MutableDie); ok {
		t.Error("an immutable die is a MutableDie")
	}
	if _, ok := m.(base.ImmutableDie); ok {
		t.Error("a mutable die is an ImmutableDie")
	}
	if _, ok := d.(interface {
		DieImmutable() *dieappsv1.DeploymentDie
	}); ok {
		t.Error("an immutable die has DieImmutable")
	}
	if _, ok := m.(interface {
		DieMutable() *dieappsv1.DeploymentDieMut
	}); ok {
		t.Error("a mutable die has DieMutable")
	}
	if !base.IsMutable(m) || base.IsMutable(d) {
		t.Error("IsMutable")
	}
}

// Code over any die or list uses what every one has.
func describe(d die.Die) string {
	data, _ := json.Marshal(d)
	return string(data)
}

func TestDieInterfaces(t *testing.T) {
	d := stampName(dieappsv1.DeploymentBlank, "web")
	if d.GetName() != "web" {
		t.Errorf("name %q", d.GetName())
	}
	if got := describe(d.Metadata()); got != `{"name":"web"}` {
		t.Errorf("metadata %s", got)
	}

	var l die.DieList = dieappsv1.NewDeploymentList(d.DieMutable())
	items, err := meta.ExtractList(l)
	if err != nil || len(items) != 1 {
		t.Fatalf("items %v, %v", items, err)
	}
	var typed die.DieListOf[*dieappsv1.DeploymentDieMut] = dieappsv1.NewDeploymentList(d.DieMutable())
	if typed.DieItems()[0].GetName() != "web" {
		t.Errorf("items %v", typed.DieItems())
	}
	// a list's metadata is its own
	typed.Metadata().ResourceVersion = "42"
	if typed.GetResourceVersion() != "42" {
		t.Errorf("resource version %q", typed.GetResourceVersion())
	}
}
