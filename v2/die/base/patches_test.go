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
	"reflect"
	"strings"
	"testing"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	diecorev1 "reconciler.io/dies/v2/apis/core/v1"
)

// Patches fed as intent are released again as the same patch.

func sameJSON(t *testing.T, got, want []byte) {
	t.Helper()
	var g, w any
	if err := json.Unmarshal(got, &g); err != nil {
		t.Fatalf("%s: %v", got, err)
	}
	_ = json.Unmarshal(want, &w)
	if !reflect.DeepEqual(g, w) {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

// an apply with each kind of intent
func applied() *dieappsv1.DeploymentDie {
	return dieappsv1.DeploymentBlank.DieWith(func(d *dieappsv1.DeploymentDie) {
		d.Metadata().Namespace().Set("default")
		d.Metadata().Name().Set("web")
		d.Metadata().Labels().Set("tier", "web")
		d.Spec().Replicas().Set(0)
		d.Spec().Paused().Set(false)
		d.Spec().MinReadySeconds().Clear()
		d.Spec().Selector().MatchLabels().Set("app", "web")
		d.Spec().Template().Spec().Tolerations().Empty()
		d.Spec().Template().Spec().Containers().At("app").DieWith(func(c *diecorev1.ContainerDie) {
			c.Image().Set("app:1")
			c.Args().Set([]string{"--verbose"})
			c.Ports().At(diecorev1.ContainerPortKey{ContainerPort: 8080}).Name().Set("http")
			c.Env().At("COLOR").ValueFrom().ConfigMapKeyRef().Key().Set("color")
		})
	})
}

func TestDieFeedApplyPatchRoundTrips(t *testing.T) {
	want, err := applied().DieReleaseApplyPatchE()
	if err != nil {
		t.Fatal(err)
	}
	fed := dieappsv1.DeploymentBlank.DieFeedApplyPatch(want)
	got, err := fed.DieReleaseApplyPatchE()
	if err != nil {
		t.Fatal(err)
	}
	sameJSON(t, got, want)
	if !reflect.DeepEqual(fed.DieRelease(), applied().DieRelease()) {
		t.Errorf("fed value %+v", fed.DieRelease())
	}
}

// An apply patch is intent over the value, which is kept, and replaces the
// die's intent.
func TestDieFeedApplyPatchOverAValue(t *testing.T) {
	server := dieappsv1.DeploymentBlank.DieFeedUnstructured(unstructuredDeployment()).DieWith(func(d *dieappsv1.DeploymentDie) {
		d.Spec().Paused().Set(true)
	})
	fed := server.DieFeedApplyPatch([]byte(`{
		"apiVersion": "apps/v1", "kind": "Deployment",
		"metadata": {"name": "web", "namespace": "default"},
		"spec": {"replicas": 5, "minReadySeconds": null, "futureField": "ignored"}
	}`))
	if fed.Spec().Replicas().Get() != 5 || fed.Spec().Template().Spec().Containers().At("app").Image().Get() != "app:1" {
		t.Errorf("fed %+v", fed.DieRelease().Spec)
	}
	intent := fed.DieIntent()
	if strings.Contains(intent, "paused") || !strings.Contains(intent, "set .spec.replicas = 5") || !strings.Contains(intent, "clear .spec.minReadySeconds") {
		t.Errorf("intent:\n%s", intent)
	}
	// a field the type doesn't have isn't fed
	if strings.Contains(intent, "futureField") {
		t.Errorf("intent:\n%s", intent)
	}
}

// A merge patch fed to the state it was made against makes the same change,
// and releases as the same patch.
func TestDieFeedMergePatchRoundTrips(t *testing.T) {
	current := dieappsv1.DeploymentBlank.DieFeedUnstructured(unstructuredDeployment()).DieWith(func(d *dieappsv1.DeploymentDie) {
		d.Metadata().Labels().Set("canary", "true")
		d.Spec().MinReadySeconds().Set(5)
	}).DieCommit()
	update := current.DieWith(func(d *dieappsv1.DeploymentDie) {
		d.Spec().Replicas().Set(0)
		d.Spec().MinReadySeconds().Clear()
		d.Metadata().Labels().Remove("canary")
		d.Metadata().Labels().Set("tier", "web")
		d.Spec().Template().Spec().Containers().At("app").Image().Set("app:2")
		d.Spec().Template().Spec().Containers().At("sidecar").Image().Set("sidecar:1")
	})
	want, err := update.DieReleaseMergePatchE()
	if err != nil {
		t.Fatal(err)
	}
	fed := current.DieFeedMergePatch(want)
	got, err := fed.DieReleaseMergePatchE()
	if err != nil {
		t.Fatal(err)
	}
	sameJSON(t, got, want)
	if !reflect.DeepEqual(fed.DieRelease(), update.DieRelease()) {
		t.Errorf("fed:\n%+v\nupdate:\n%+v", fed.DieRelease(), update.DieRelease())
	}
}

func TestDieFeedPatchOnAChildPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil || !strings.Contains(r.(error).Error(), "only supported on a root die") {
			t.Errorf("recovered %v", r)
		}
	}()
	dieappsv1.DeploymentBlank.Spec().DieFeedMergePatch([]byte(`{}`))
}
