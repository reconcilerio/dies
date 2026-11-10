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
	"errors"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	"reconciler.io/dies/v2/die"
	dietesting "reconciler.io/dies/v2/die/testing"
)

func unstructuredDeployment() *unstructured.Unstructured {
	return &unstructured.Unstructured{Object: map[string]any{
		"apiVersion": "apps/v1", "kind": "Deployment",
		"metadata": map[string]any{"namespace": "default", "name": "web"},
		"spec": map[string]any{
			"replicas":    int64(3),
			"paused":      false,
			"futureField": map[string]any{"from": "a newer API"},
			"template": map[string]any{"spec": map[string]any{
				"containers": []any{map[string]any{"name": "app", "image": "app:1"}},
			}},
		},
	}}
}

// Clients strip the type metadata of a typed object as they decode it, which
// isn't intent: patches don't carry it, and the type metadata is gone.
func TestSetGroupVersionKindIsNotIntent(t *testing.T) {
	d := dieappsv1.DeploymentBlank.DieMutable()
	dietesting.NoError(t, die.FeedUnstructured(d, unstructuredDeployment()))
	dietesting.NoError(t, d.Spec().Replicas().Set(4))
	d.GetObjectKind().SetGroupVersionKind(schema.GroupVersionKind{})

	if gvk := d.GetObjectKind().GroupVersionKind(); !gvk.Empty() {
		t.Errorf("gvk %v", gvk)
	}
	if got, err := die.ReleaseMergePatch(d); err != nil || string(got) != `{"spec":{"replicas":4}}` {
		t.Errorf("merge patch %s, %v", got, err)
	}
	data, err := d.MarshalJSON()
	if err != nil || strings.Contains(string(data), `"apiVersion"`) || strings.Contains(string(data), `"kind"`) {
		t.Errorf("encoded %s, %v", data, err)
	}
	// an apply always identifies its object
	if got, err := die.ReleaseApplyPatch(d); err != nil || !strings.Contains(string(got), `"apiVersion":"apps/v1","kind":"Deployment"`) {
		t.Errorf("apply patch %s, %v", got, err)
	}
}

// An object die feeds and releases unstructured objects through its codec, so
// presence and fields its type doesn't have survive both ways.
func TestDieFeedUnstructured(t *testing.T) {
	d := dieappsv1.DeploymentBlank.DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietesting.NoError(t, die.FeedUnstructured(d, unstructuredDeployment()))
	})
	if d.GetName() != "web" || d.Spec().Replicas().Get() != 3 || d.Spec().Template().Spec().Containers().At("app").Image().Get() != "app:1" {
		t.Errorf("fed %+v", *d.DieRelease())
	}
	if !d.Spec().Paused().Has() || d.Spec().MinReadySeconds().Has() {
		t.Errorf("paused %v, minReadySeconds %v", d.Spec().Paused().Has(), d.Spec().MinReadySeconds().Has())
	}

	d = d.DieStamp(func(d *dieappsv1.DeploymentDieMut) { dietesting.NoError(t, d.Spec().Replicas().Set(4)) })
	u := die.MustReleaseUnstructured(d)
	if got, _, _ := unstructured.NestedInt64(u.Object, "spec", "replicas"); got != 4 {
		t.Errorf("replicas %v", u.Object["spec"])
	}
	if got, _, _ := unstructured.NestedString(u.Object, "spec", "futureField", "from"); got != "a newer API" {
		t.Errorf("futureField %v", u.Object["spec"])
	}
	if got, ok, _ := unstructured.NestedBool(u.Object, "spec", "paused"); !ok || got {
		t.Errorf("paused %v, %v", got, ok)
	}
	if _, ok := u.Object["status"]; ok {
		t.Errorf("status %v, which was never present", u.Object["status"])
	}

	// as with DieFeed, a mutable root discards its intent
	m := dieappsv1.DeploymentBlank.DieMutable()
	dietesting.NoError(t, m.Spec().Replicas().Set(9))
	dietesting.NoError(t, die.FeedUnstructured(m, unstructuredDeployment()))
	if got, err := die.ReleaseMergePatch(m); err != nil || string(got) != "{}" {
		t.Errorf("patch %s, %v", got, err)
	}

	bad := unstructuredDeployment()
	bad.Object["spec"] = "not an object"
	if err := die.FeedUnstructured(dieappsv1.DeploymentBlank.DieMutable(), bad); err == nil || !strings.Contains(err.Error(), "feeding") {
		t.Errorf("error %v", err)
	}
}

// DieStampE runs an operation that fills in a die, such as a client's, on a
// mutable copy of an immutable die, and returns the result immutable, with
// the operation's error.
func TestDieStampE(t *testing.T) {
	blank := dieappsv1.DeploymentBlank
	got, err := blank.DieStampE(func(d *dieappsv1.DeploymentDieMut) error {
		dietesting.NoError(t, die.FeedUnstructured(d, unstructuredDeployment()))
		return nil
	})
	if err != nil || got.GetName() != "web" {
		t.Errorf("got %s, %v", got.GetName(), err)
	}
	if blank.GetName() != "" {
		t.Errorf("the blank changed")
	}

	failed := errors.New("failed")
	got, err = blank.DieStampE(func(d *dieappsv1.DeploymentDieMut) error {
		dietesting.NoError(t, d.Spec().Replicas().Set(2))
		return failed
	})
	if err != failed || got.Spec().Replicas().Get() != 2 {
		t.Errorf("got replicas %d, %v", got.Spec().Replicas().Get(), err)
	}

	// a client zeroes the object it decodes into, replacing the die's state
	got, err = blank.DieStampE(func(d *dieappsv1.DeploymentDieMut) error {
		*d = dieappsv1.DeploymentDieMut{}
		return json.Unmarshal([]byte(`{"metadata": {"name": "decoded"}}`), d)
	})
	if err != nil || got.GetName() != "decoded" {
		t.Errorf("got %q, %v", got.GetName(), err)
	}

	// a mutable die passes itself, and keeps what a client decoded into it
	m := blank.DieMutable()
	got2, err := m.DieStampE(func(d *dieappsv1.DeploymentDieMut) error {
		*d = dieappsv1.DeploymentDieMut{}
		return json.Unmarshal([]byte(`{"metadata": {"name": "decoded"}}`), d)
	})
	if err != nil || got2.GetName() != "decoded" || m.GetName() != "decoded" {
		t.Errorf("got %q, die %q, %v", got2.GetName(), m.GetName(), err)
	}
}
