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
	"reconciler.io/dies/v2/die"
	dietesting "reconciler.io/dies/v2/die/testing"
)

const webJSON = `{"apiVersion": "apps/v1", "kind": "Deployment", "metadata": {"name": "web"},
	"spec": {"paused": false, "futureField": {"from": "a newer API"}}}`

// Encodings feed and release through the codec: presence and fields the type
// doesn't have survive.
func TestDieFeedAndReleaseEncodings(t *testing.T) {
	for name, d := range map[string]*dieappsv1.DeploymentDie{
		"JSON": dieappsv1.DeploymentBlank.DieStamp(func(d *dieappsv1.DeploymentDieMut) { dietesting.NoError(t, die.FeedJSON(d, []byte(webJSON))) }),
		"YAML": dieappsv1.DeploymentBlank.DieStamp(func(d *dieappsv1.DeploymentDieMut) {
			dietesting.NoError(t, die.FeedYAML(d, []byte(`
apiVersion: apps/v1
kind: Deployment
metadata: {name: web}
spec:
  paused: false
  futureField: {from: a newer API}
`)))
		}),
		"a raw extension's raw JSON": dieappsv1.DeploymentBlank.DieStamp(func(d *dieappsv1.DeploymentDieMut) {
			dietesting.NoError(t, die.FeedRawExtension(d, runtime.RawExtension{Raw: []byte(webJSON)}))
		}),
		"a raw extension's object": dieappsv1.DeploymentBlank.DieStamp(func(d *dieappsv1.DeploymentDieMut) {
			web := dieappsv1.DeploymentBlank.DieMutable()
			dietesting.NoError(t, die.FeedJSON(web, []byte(webJSON)), die.FeedRawExtension(d, runtime.RawExtension{Object: web}))
		}),
	} {
		t.Run(name, func(t *testing.T) {
			if d.GetName() != "web" || !d.Spec().Paused().Has() {
				t.Errorf("fed %s", die.MustReleaseJSON(d))
			}
			assertJSON(t, die.MustReleaseJSON(d), nil, webJSON)
			assertJSON(t, die.MustReleaseRawExtension(d).Raw, nil, webJSON)
			if y := string(die.MustReleaseYAML(d)); !strings.Contains(y, "futureField:\n    from: a newer API") || !strings.Contains(y, "paused: false") {
				t.Errorf("YAML %s", y)
			}
		})
	}
}

func TestDieFeedJSONIsDieFeed(t *testing.T) {
	// a mutable root discards its intent
	m := dieappsv1.DeploymentBlank.DieMutable()
	dietesting.NoError(t, m.Spec().Replicas().Set(9))
	dietesting.NoError(t, die.FeedJSON(m, []byte(webJSON)))
	if got, err := die.ReleaseMergePatch(m); err != nil || string(got) != "{}" {
		t.Errorf("patch %s, %v", got, err)
	}
	// a die within another sets the value as intent
	d := dieappsv1.DeploymentBlank.DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietesting.NoError(t, die.FeedJSON(d.Spec().Template().Spec(), []byte(`{"serviceAccountName": "sa"}`)))
	})
	got, err := die.ReleaseMergePatch(d)
	assertJSON(t, got, err, `{"spec": {"template": {"spec": {"serviceAccountName": "sa"}}}}`)

	if err := die.FeedJSON(dieappsv1.DeploymentBlank.DieMutable(), []byte(`{"spec": "not an object"}`)); err == nil || !strings.Contains(err.Error(), "feeding v1.Deployment from JSON") {
		t.Errorf("error %v", err)
	}
}
