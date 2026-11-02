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

	upstreamappsv1 "k8s.io/api/apps/v1"
	upstreamcorev1 "k8s.io/api/core/v1"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	diecorev1 "reconciler.io/dies/v2/apis/core/v1"
)

// A value set whole has the presence a Go value can have: nil and omitempty
// zero fields are absent, never null.
func TestCodecSetWhole(t *testing.T) {
	d := dieappsv1.DeploymentBlank.DieWith(func(d *dieappsv1.DeploymentDie) {
		d.Metadata().Name().Set("web")
		d.Spec().Template().DieFeed(upstreamcorev1.PodTemplateSpec{
			Spec: upstreamcorev1.PodSpec{ServiceAccountName: "sa"},
		})
	})
	got, err := d.DieReleaseApplyPatchE()
	assertJSON(t, got, err, `{
		"apiVersion": "apps/v1", "kind": "Deployment", "metadata": {"name": "web"},
		"spec": {"template": {"spec": {"serviceAccountName": "sa"}}}
	}`)
}

// Decoding records what was present on the wire, zero or not.
func TestCodecWirePresence(t *testing.T) {
	d := dieappsv1.DeploymentBlank.DieImmutable(false)
	if err := json.Unmarshal([]byte(`{"spec": {"paused": false, "minReadySeconds": 0}}`), d); err != nil {
		t.Fatal(err)
	}
	if !d.Spec().Paused().Has() || !d.Spec().MinReadySeconds().Has() {
		t.Errorf("paused %v, minReadySeconds %v", d.Spec().Paused().Has(), d.Spec().MinReadySeconds().Has())
	}
	if d.Spec().ProgressDeadlineSeconds().Has() {
		t.Error("progressDeadlineSeconds wasn't on the wire")
	}
}

// Fields the Go type doesn't have are kept, and encoded again.
func TestCodecUnknownFields(t *testing.T) {
	d := dieappsv1.DeploymentBlank.DieImmutable(false)
	if err := json.Unmarshal([]byte(`{
		"apiVersion": "apps/v1", "kind": "Deployment", "metadata": {"name": "web"},
		"spec": {"futureField": {"from": "a newer API"}, "template": {"spec": {"containers": [{"name": "app", "futureFlag": true}]}}}
	}`), d); err != nil {
		t.Fatal(err)
	}
	d.Spec().Replicas().Set(2)
	// only what was on the wire or written since, unknown fields included, and
	// none of the nulls and empty objects the Go types would add
	got, err := json.Marshal(d)
	assertJSON(t, got, err, `{
		"apiVersion": "apps/v1", "kind": "Deployment", "metadata": {"name": "web"},
		"spec": {"replicas": 2, "futureField": {"from": "a newer API"},
			"template": {"spec": {"containers": [{"name": "app", "futureFlag": true}]}}}
	}`)
}

func decoded(t *testing.T, data string) *dieappsv1.DeploymentDie {
	t.Helper()
	d := dieappsv1.DeploymentBlank.DieImmutable(false)
	if err := json.Unmarshal([]byte(data), d); err != nil {
		t.Fatal(err)
	}
	return d
}

// Presence follows writes and survives a commit.
func TestCodecPresenceThroughWrites(t *testing.T) {
	d := decoded(t, `{"spec": {"paused": false, "minReadySeconds": 0}}`)
	d.Spec().Paused().Clear()
	d.Spec().RevisionHistoryLimit().Set(0)
	if d.Spec().Paused().Has() || !d.Spec().RevisionHistoryLimit().Has() {
		t.Errorf("paused %v, revisionHistoryLimit %v", d.Spec().Paused().Has(), d.Spec().RevisionHistoryLimit().Has())
	}
	got, err := d.DieReleaseMergePatchE()
	assertJSON(t, got, err, `{"spec": {"paused": null, "revisionHistoryLimit": 0}}`)

	d.DieCommit()
	if d.Spec().Paused().Has() || !d.Spec().MinReadySeconds().Has() || !d.Spec().RevisionHistoryLimit().Has() {
		t.Errorf("after commit: paused %v, minReadySeconds %v, revisionHistoryLimit %v", d.Spec().Paused().Has(), d.Spec().MinReadySeconds().Has(), d.Spec().RevisionHistoryLimit().Has())
	}
	got, err = json.Marshal(d)
	assertJSON(t, got, err, `{"spec": {"minReadySeconds": 0, "revisionHistoryLimit": 0}}`)

	// a value fed from Go has the presence the Go rule gives it
	fed := dieappsv1.DeploymentBlank.DieFeed(d.DieRelease())
	if fed.Spec().MinReadySeconds().Has() || !fed.Spec().RevisionHistoryLimit().Has() {
		t.Errorf("fed: minReadySeconds %v, revisionHistoryLimit %v", fed.Spec().MinReadySeconds().Has(), fed.Spec().RevisionHistoryLimit().Has())
	}
}

// Unknown fields are kept within list maps, found by keys with defaults, and
// within atomic lists, by index.
func TestCodecUnknownFieldsInLists(t *testing.T) {
	d := decoded(t, `{"spec": {"template": {"spec": {
		"containers": [{"name": "app", "ports": [{"containerPort": 80, "futurePort": 1}]}],
		"tolerations": [{"key": "a", "futureToleration": 2}]
	}}}}`)
	d.Spec().Template().Spec().Containers().At("app").Ports().At(diecorev1.ContainerPortKey{ContainerPort: 80}).Name().Set("http")
	got, err := json.Marshal(d)
	assertJSON(t, got, err, `{"spec": {"template": {"spec": {
		"containers": [{"name": "app", "ports": [{"containerPort": 80, "name": "http", "futurePort": 1}]}],
		"tolerations": [{"key": "a", "futureToleration": 2}]
	}}}}`)
}

// Replacing a value from Go replaces its unknown fields too.
func TestCodecReplacingDropsUnknownFields(t *testing.T) {
	d := decoded(t, `{"spec": {"strategy": {"type": "Recreate", "futureStrategy": true}}}`)
	d.Spec().Strategy().DieFeed(upstreamappsv1.DeploymentStrategy{Type: upstreamappsv1.RollingUpdateDeploymentStrategyType})
	got, err := d.DieReleaseMergePatchE()
	assertJSON(t, got, err, `{"spec": {"strategy": {"type": "RollingUpdate", "futureStrategy": null}}}`)
}
