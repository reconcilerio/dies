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
	d := dieappsv1.DeploymentBlank.DieFeedUnstructured(unstructuredDeployment()).DieImmutable(false)
	d.Spec().Replicas().Set(4)
	d.GetObjectKind().SetGroupVersionKind(schema.GroupVersionKind{})

	if gvk := d.GetObjectKind().GroupVersionKind(); !gvk.Empty() {
		t.Errorf("gvk %v", gvk)
	}
	if got, err := d.DieReleaseMergePatchE(); err != nil || string(got) != `{"spec":{"replicas":4}}` {
		t.Errorf("merge patch %s, %v", got, err)
	}
	data, err := d.MarshalJSON()
	if err != nil || strings.Contains(string(data), `"apiVersion"`) || strings.Contains(string(data), `"kind"`) {
		t.Errorf("encoded %s, %v", data, err)
	}
	// an apply always identifies its object
	if got, err := d.DieReleaseApplyPatchE(); err != nil || !strings.Contains(string(got), `"apiVersion":"apps/v1","kind":"Deployment"`) {
		t.Errorf("apply patch %s, %v", got, err)
	}
}

// An object die feeds and releases unstructured objects through its codec, so
// presence and fields its type doesn't have survive both ways.
func TestDieFeedUnstructured(t *testing.T) {
	d := dieappsv1.DeploymentBlank.DieFeedUnstructured(unstructuredDeployment())
	if d.GetName() != "web" || d.Spec().Replicas().Get() != 3 || d.Spec().Template().Spec().Containers().At("app").Image().Get() != "app:1" {
		t.Errorf("fed %+v", d.DieRelease())
	}
	if !d.Spec().Paused().Has() || d.Spec().MinReadySeconds().Has() {
		t.Errorf("paused %v, minReadySeconds %v", d.Spec().Paused().Has(), d.Spec().MinReadySeconds().Has())
	}

	d = d.DieWith(func(d *dieappsv1.DeploymentDie) { d.Spec().Replicas().Set(4) })
	u := d.DieReleaseUnstructured()
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
	m := dieappsv1.DeploymentBlank.DieImmutable(false)
	m.Spec().Replicas().Set(9)
	m.DieFeedUnstructured(unstructuredDeployment())
	if got, err := m.DieReleaseMergePatchE(); err != nil || string(got) != "{}" {
		t.Errorf("patch %s, %v", got, err)
	}

	bad := unstructuredDeployment()
	bad.Object["spec"] = "not an object"
	defer func() {
		if r := recover(); r == nil || !strings.Contains(r.(error).Error(), "feeding") {
			t.Errorf("recovered %v", r)
		}
	}()
	dieappsv1.DeploymentBlank.DieFeedUnstructured(bad)
}

// DieWithE runs an operation that fills in a die, such as a client's, on a
// mutable copy of an immutable die, and returns the result immutable, with
// the operation's error.
func TestDieWithE(t *testing.T) {
	blank := dieappsv1.DeploymentBlank
	got, err := blank.DieWithE(func(d *dieappsv1.DeploymentDie) error {
		d.DieFeedUnstructured(unstructuredDeployment())
		return nil
	})
	if err != nil || got.GetName() != "web" || got.DieMutable() {
		t.Errorf("got %s, mutable %v, %v", got.GetName(), got.DieMutable(), err)
	}
	if blank.GetName() != "" {
		t.Errorf("the blank changed")
	}

	failed := errors.New("failed")
	got, err = blank.DieWithE(func(d *dieappsv1.DeploymentDie) error {
		d.Spec().Replicas().Set(2)
		return failed
	})
	if err != failed || got.Spec().Replicas().Get() != 2 {
		t.Errorf("got replicas %d, %v", got.Spec().Replicas().Get(), err)
	}

	// a client zeroes the object it decodes into, replacing the die's state
	got, err = blank.DieWithE(func(d *dieappsv1.DeploymentDie) error {
		*d = dieappsv1.DeploymentDie{}
		return json.Unmarshal([]byte(`{"metadata": {"name": "decoded"}}`), d)
	})
	if err != nil || got.GetName() != "decoded" {
		t.Errorf("got %q, %v", got.GetName(), err)
	}

	// a mutable die is passed through
	m := blank.DieImmutable(false)
	if got, _ := m.DieWithE(func(d *dieappsv1.DeploymentDie) error { return nil }); got.DieMutable() != true {
		t.Errorf("mutable die came back immutable")
	}
}
