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

package cluster_test

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/runtime/schema"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	"reconciler.io/dies/v2/conformance/internal/harness"
	"reconciler.io/dies/v2/die"
	"reconciler.io/dies/v2/die/dietest"
	dieunstructured "reconciler.io/dies/v2/die/unstructured"
)

// Unstructured dies find the molds of the kinds the cluster serves: built
// in, custom resources, and an API served through an APIService.

var (
	deploymentGVK = schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}
	widgetGVK     = schema.GroupVersionKind{Group: "example.dev", Version: "v1", Kind: "Widget"}
	flunderGVK    = schema.GroupVersionKind{Group: "wardle.example.com", Version: "v1alpha1", Kind: "Flunder"}
	fischerGVK    = schema.GroupVersionKind{Group: "wardle.example.com", Version: "v1alpha1", Kind: "Fischer"}
)

// For the same changes, an unstructured Deployment with the cluster's schema
// sends the apply the typed die does, and the cluster takes it.
func TestUnstructuredBuiltInAppliesAsTyped(t *testing.T) {
	ns := harness.Namespace(t, dies)
	blank := harness.MoldedBlank(t, harness.Molds(t, config), deploymentGVK)
	change := func(d *dieappsv1.DeploymentDieMut) {
		dietest.NoError(t, d.Metadata().Namespace().Set(ns))
		dietest.NoError(t, d.Metadata().Name().Set("web"))
		dietest.NoError(t, d.Spec().Replicas().Set(0))
		dietest.NoError(t, d.Spec().Paused().Set(true))
		dietest.NoError(t, d.Spec().MinReadySeconds().Clear())
		dietest.NoError(t, d.Spec().Selector().MatchLabels().Set("app", "web"))
		dietest.NoError(t, d.Spec().Template().Metadata().Labels().Set("app", "web"))
		dietest.NoError(t, d.Spec().Template().Spec().Containers().At("app").Image().Set("registry.example/app:1"))
	}
	typed := dieappsv1.DeploymentBlank.DieStamp(change)
	u := blank.DieStamp(func(u *dieunstructured.ObjectMut) {
		change(dieunstructured.JigAsMut[*dieappsv1.DeploymentDieMut](u.Node()))
	})

	want, _ := die.ReleaseApplyPatchE(typed)
	got, err := die.ReleaseApplyPatchE(u)
	if err != nil {
		t.Fatal(err)
	}
	assertSameJSON(t, got, want)
	applied := apply(t, u, "app", false)
	harness.AssertOwns(t, applied, "app", "", ".spec.replicas", ".spec.minReadySeconds", ".spec.selector")
}

// A CRD's schema is found once the cluster publishes it. Clearing a field
// that isn't nullable is left out, as the CRD's validation would reject
// null, so it stays another manager's; a nullable one is sent as null, and
// removed.
func TestUnstructuredCustomResource(t *testing.T) {
	ns := harness.Namespace(t, dies)
	blank := harness.MoldedBlank(t, harness.Molds(t, config), widgetGVK).DieStamp(func(d *dieunstructured.ObjectMut) {
		d.SetNamespace(ns)
		d.SetName("web")
	})
	apply(t, blank.DieStamp(func(d *dieunstructured.ObjectMut) {
		match := d.Node().Field("spec").Field("match")
		dietest.NoError(t, match.Field("equal").Set(true))
		dietest.NoError(t, match.Field("note").Set("set by matcher"))
		dietest.NoError(t, d.Node().Field("spec").Field("ports").Item("name", "http").Field("port").Set(int64(80)))
	}), "matcher", false)

	got := apply(t, blank.DieStamp(func(d *dieunstructured.ObjectMut) {
		match := d.Node().Field("spec").Field("match")
		dietest.NoError(t, match.Field("equal").Clear())
		dietest.NoError(t, match.Field("note").Clear())
		dietest.NoError(t, d.Node().Field("spec").Field("ports").Item("name", "metrics").Field("port").Set(int64(9090)))
	}), "clearer", false)
	spec := got.Node().Field("spec")
	if spec.Field("match").Field("note").Has() || !spec.Field("match").Field("equal").Has() {
		t.Errorf("match %v, want note removed and equal kept", spec.Field("match").Get())
	}
	if spec.Field("ports").Len() != 2 {
		t.Errorf("ports %v, want http and metrics", spec.Field("ports").Get())
	}
}

// wardle is the molds of the cluster, if it serves sample-apiserver's
// wardle.example.com through an APIService; the test is skipped if not.
func wardle(t *testing.T) *dieunstructured.Molds {
	t.Helper()
	molds := harness.Molds(t, config)
	if _, err := molds.BlankE(flunderGVK); err != nil {
		if strings.Contains(err.Error(), "no OpenAPI schema") {
			t.Skip("the cluster doesn't serve wardle.example.com/v1alpha1; hack/conformance-cluster.sh installs Kubernetes' sample-apiserver, which does, behind an APIService")
		}
		t.Fatal(err)
	}
	return molds
}

// An aggregated API's schema is found through the cluster's OpenAPI, as the
// aggregator publishes it, and its server converts objects through Go
// types, so clearing any field sends null, which removes it.
func TestAggregatedAPI(t *testing.T) {
	molds := wardle(t)
	ns := harness.Namespace(t, dies)
	flunder := harness.MoldedBlank(t, molds, flunderGVK).DieStamp(func(d *dieunstructured.ObjectMut) {
		d.SetNamespace(ns)
		d.SetName("web")
	})
	apply(t, flunder.DieStamp(func(d *dieunstructured.ObjectMut) {
		dietest.NoError(t, d.Node().Field("spec").Field("reference").Set("other"))
	}), "referrer", false)
	got := apply(t, flunder.DieStamp(func(d *dieunstructured.ObjectMut) {
		dietest.NoError(t, d.Node().Field("spec").Field("referenceType").Set("Flunder"))
	}), "typer", false)
	if spec := got.Node().Field("spec"); spec.Field("reference").Get() != "other" || spec.Field("referenceType").Get() != "Flunder" {
		t.Errorf("spec %v, want each manager's field", spec.Get())
	}
	harness.AssertOwned(t, got, "typer", "", ".spec.referenceType")

	got = apply(t, flunder.DieStamp(func(d *dieunstructured.ObjectMut) {
		dietest.NoError(t, d.Node().Field("spec").Field("reference").Clear())
	}), "clearer", false)
	if got.Node().Field("spec").Field("reference").Has() {
		t.Errorf("spec %v, want reference removed", got.Node().Field("spec").Get())
	}
	harness.AssertOwns(t, got, "clearer", "", ".spec.reference")
}

// An aggregated API's atomic list is applied whole: a change to one element,
// from what the server has, keeps the others.
func TestAggregatedAPIAtomicList(t *testing.T) {
	molds := wardle(t)
	name := "dies-" + strings.ToLower(strings.ReplaceAll(t.Name(), "/", "-"))
	// cluster scoped, so it's labeled and deleted after, as namespaces are
	fischer := harness.MoldedBlank(t, molds, fischerGVK).DieStamp(func(d *dieunstructured.ObjectMut) {
		d.SetName(name)
		d.SetLabels(map[string]string{harness.Label: "true"})
	})
	actual := apply(t, fischer.DieStamp(func(d *dieunstructured.ObjectMut) {
		dietest.NoError(t, d.Node().Field("disallowedFlunders").Set([]any{"a", "b"}))
	}), "app", false)
	// t.Context is done by the time cleanups run
	t.Cleanup(func() { _ = dies.Delete(context.Background(), actual.DieMutable()) })

	got := apply(t, actual.DieStamp(func(d *dieunstructured.ObjectMut) {
		dietest.NoError(t, d.Node().Field("disallowedFlunders").Index(1).Set("c"))
	}), "app", false)
	if list := got.Node().Field("disallowedFlunders").Get(); !reflect.DeepEqual(list, []any{"a", "c"}) {
		t.Errorf("disallowedFlunders %v, want [a c]", list)
	}
}

func assertSameJSON(t *testing.T, got, want []byte) {
	t.Helper()
	var g, w any
	_ = json.Unmarshal(got, &g)
	_ = json.Unmarshal(want, &w)
	if !reflect.DeepEqual(g, w) {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}
