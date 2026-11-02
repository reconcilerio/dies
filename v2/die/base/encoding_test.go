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
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
)

const webJSON = `{"apiVersion": "apps/v1", "kind": "Deployment", "metadata": {"name": "web"},
	"spec": {"paused": false, "futureField": {"from": "a newer API"}}}`

// Encodings feed and release through the codec: presence and fields the type
// doesn't have survive.
func TestDieFeedAndReleaseEncodings(t *testing.T) {
	for name, d := range map[string]*dieappsv1.DeploymentDie{
		"JSON": dieappsv1.DeploymentBlank.DieFeedJSON([]byte(webJSON)),
		"YAML": dieappsv1.DeploymentBlank.DieFeedYAML([]byte(`
apiVersion: apps/v1
kind: Deployment
metadata: {name: web}
spec:
  paused: false
  futureField: {from: a newer API}
`)),
		"a raw extension's raw JSON": dieappsv1.DeploymentBlank.DieFeedRawExtension(runtime.RawExtension{Raw: []byte(webJSON)}),
		"a raw extension's object":   dieappsv1.DeploymentBlank.DieFeedRawExtension(runtime.RawExtension{Object: dieappsv1.DeploymentBlank.DieFeedJSON([]byte(webJSON))}),
	} {
		t.Run(name, func(t *testing.T) {
			if d.DieMutable() || d.GetName() != "web" || !d.Spec().Paused().Has() {
				t.Errorf("fed %s", d.DieReleaseJSON())
			}
			assertJSON(t, d.DieReleaseJSON(), nil, webJSON)
			assertJSON(t, d.DieReleaseRawExtension().Raw, nil, webJSON)
			if y := string(d.DieReleaseYAML()); !strings.Contains(y, "futureField:\n    from: a newer API") || !strings.Contains(y, "paused: false") {
				t.Errorf("YAML %s", y)
			}
		})
	}
}

func TestDieFeedJSONIsDieFeed(t *testing.T) {
	// a mutable root discards its intent
	m := dieappsv1.DeploymentBlank.DieImmutable(false)
	m.Spec().Replicas().Set(9)
	m.DieFeedJSON([]byte(webJSON))
	if got, err := m.DieReleaseMergePatchE(); err != nil || string(got) != "{}" {
		t.Errorf("patch %s, %v", got, err)
	}
	// a die within another sets the value as intent
	d := dieappsv1.DeploymentBlank.DieWith(func(d *dieappsv1.DeploymentDie) {
		d.Spec().Template().Spec().DieFeedJSON([]byte(`{"serviceAccountName": "sa"}`))
	})
	got, err := d.DieReleaseMergePatchE()
	assertJSON(t, got, err, `{"spec": {"template": {"spec": {"serviceAccountName": "sa"}}}}`)

	defer func() {
		if r := recover(); r == nil || !strings.Contains(r.(error).Error(), "feeding v1.Deployment from JSON") {
			t.Errorf("recovered %v", r)
		}
	}()
	dieappsv1.DeploymentBlank.DieFeedJSON([]byte(`{"spec": "not an object"}`))
}
