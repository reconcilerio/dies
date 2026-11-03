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
	"reconciler.io/dies/v2/die"
	"reconciler.io/dies/v2/die/dietest"
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
func applied(t testing.TB) *dieappsv1.DeploymentDie {
	return dieappsv1.DeploymentBlank.DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietest.NoError(t, d.Metadata().Namespace().Set("default"))
		dietest.NoError(t, d.Metadata().Name().Set("web"))
		dietest.NoError(t, d.Metadata().Labels().Set("tier", "web"))
		dietest.NoError(t, d.Spec().Replicas().Set(0))
		dietest.NoError(t, d.Spec().Paused().Set(false))
		dietest.NoError(t, d.Spec().MinReadySeconds().Clear())
		dietest.NoError(t, d.Spec().Selector().MatchLabels().Set("app", "web"))
		dietest.NoError(t, d.Spec().Template().Spec().Tolerations().Empty())
		d.Spec().Template().Spec().Containers().At("app").DieStamp(func(c *diecorev1.ContainerDieMut) {
			dietest.NoError(t, c.Image().Set("app:1"))
			dietest.NoError(t, c.Args().Set([]string{"--verbose"}))
			dietest.NoError(t, c.Ports().At(diecorev1.ContainerPortKey{ContainerPort: 8080}).Name().Set("http"))
			dietest.NoError(t, c.Env().At("COLOR").ValueFrom().ConfigMapKeyRef().Key().Set("color"))
		})
	})
}

func TestDieFeedApplyPatchRoundTrips(t *testing.T) {
	want, err := die.ReleaseApplyPatchE(applied(t))
	if err != nil {
		t.Fatal(err)
	}
	fed := dieappsv1.DeploymentBlank.DieStamp(func(d *dieappsv1.DeploymentDieMut) { dietest.NoError(t, die.FeedApplyPatch(d, want)) })
	got, err := die.ReleaseApplyPatchE(fed)
	if err != nil {
		t.Fatal(err)
	}
	sameJSON(t, got, want)
	if !reflect.DeepEqual(*fed.DieRelease(), *applied(t).DieRelease()) {
		t.Errorf("fed value %+v", *fed.DieRelease())
	}
}

// An apply patch is intent over the value, which is kept, and replaces the
// die's intent.
func TestDieFeedApplyPatchOverAValue(t *testing.T) {
	server := dieappsv1.DeploymentBlank.DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietest.NoError(t, die.FeedUnstructured(d, unstructuredDeployment()))
	}).DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietest.NoError(t, d.Spec().Paused().Set(true))
	})
	fed := server.DieMutable()
	dietest.NoError(t, die.FeedApplyPatch(fed, []byte(`{
		"apiVersion": "apps/v1", "kind": "Deployment",
		"metadata": {"name": "web", "namespace": "default"},
		"spec": {"replicas": 5, "minReadySeconds": null, "futureField": "ignored"}
	}`)))
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
	current := dieappsv1.DeploymentBlank.DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietest.NoError(t, die.FeedUnstructured(d, unstructuredDeployment()))
	}).DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietest.NoError(t, d.Metadata().Labels().Set("canary", "true"))
		dietest.NoError(t, d.Spec().MinReadySeconds().Set(5))
	}).DieTare()
	update := current.DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietest.NoError(t, d.Spec().Replicas().Set(0))
		dietest.NoError(t, d.Spec().MinReadySeconds().Clear())
		dietest.NoError(t, d.Metadata().Labels().Remove("canary"))
		dietest.NoError(t, d.Metadata().Labels().Set("tier", "web"))
		dietest.NoError(t, d.Spec().Template().Spec().Containers().At("app").Image().Set("app:2"))
		dietest.NoError(t, d.Spec().Template().Spec().Containers().At("sidecar").Image().Set("sidecar:1"))
	})
	want, err := die.ReleaseMergePatchE(update)
	if err != nil {
		t.Fatal(err)
	}
	fed := current.DieMutable()
	dietest.NoError(t, die.FeedMergePatch(fed, want))
	got, err := die.ReleaseMergePatchE(fed)
	if err != nil {
		t.Fatal(err)
	}
	sameJSON(t, got, want)
	if !reflect.DeepEqual(*fed.DieRelease(), *update.DieRelease()) {
		t.Errorf("fed:\n%+v\nupdate:\n%+v", *fed.DieRelease(), *update.DieRelease())
	}
}

func TestDieFeedPatchOnAChildFails(t *testing.T) {
	if err := die.FeedMergePatch(dieappsv1.DeploymentBlank.DieMutable().Spec(), []byte(`{}`)); err == nil || !strings.Contains(err.Error(), "only supported on a root die") {
		t.Errorf("error %v", err)
	}
}
