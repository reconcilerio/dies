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

package apiserver_test

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	diecorev1 "reconciler.io/dies/v2/apis/core/v1"
	"reconciler.io/dies/v2/die"
	"reconciler.io/dies/v2/die/dietest"
	dieunstructured "reconciler.io/dies/v2/die/unstructured"
	widgetv1 "reconciler.io/dies/v2/examples/widget/v1"
)

// Unstructured dies with the server's molds apply as typed dies do.

var (
	deploymentGVK = schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}
	widgetGVK     = schema.GroupVersionKind{Group: "example.dev", Version: "v1", Kind: "Widget"}
)

// For the same changes, an unstructured Deployment with the server's schema
// sends the same apply as the typed die. The changes are made through a typed
// jig of the unstructured die, so they're the same intent.
func TestUnstructuredAppliesAsTyped(t *testing.T) {
	blank := moldedBlank(t, serverMolds(t), deploymentGVK)
	server := die.ReleaseJSON(deployment(t, "default").DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietest.NoError(t, d.Metadata().Labels().Set("tier", "web"))
		dietest.NoError(t, d.Spec().MinReadySeconds().Set(5))
		ref := d.Spec().Template().Spec().Containers().At("app").Env().At("COLOR").ValueFrom().ConfigMapKeyRef()
		dietest.NoError(t, ref.Name().Set("settings"))
		dietest.NoError(t, ref.Key().Set("color"))
	}))

	for _, c := range []struct {
		name   string
		server bool
		change func(d *dieappsv1.DeploymentDieMut)
	}{
		{"from a blank", false, func(d *dieappsv1.DeploymentDieMut) {
			dietest.NoError(t, d.Metadata().Namespace().Set("default"))
			dietest.NoError(t, d.Metadata().Name().Set("web"))
			dietest.NoError(t, d.Spec().Replicas().Set(0))
			dietest.NoError(t, d.Spec().Selector().MatchLabels().Set("app", "web"))
			d.Spec().Template().Spec().Containers().At("app").DieStamp(func(c *diecorev1.ContainerDieMut) {
				dietest.NoError(t, c.Image().Set("app:1"))
				dietest.NoError(t, c.Args().Set([]string{"--verbose"}))
				dietest.NoError(t, c.Ports().At(diecorev1.ContainerPortKey{ContainerPort: 8080, Protocol: diecorev1.ProtocolTCP}).Name().Set("http"))
			})
			dietest.NoError(t, d.Spec().Template().Spec().Tolerations().Empty())
		}},
		{"clearing a field", true, func(d *dieappsv1.DeploymentDieMut) {
			dietest.NoError(t, d.Spec().MinReadySeconds().Clear())
		}},
		{"removing a map entry", true, func(d *dieappsv1.DeploymentDieMut) {
			dietest.NoError(t, d.Metadata().Labels().Remove("tier"))
		}},
		{"within an atomic struct", true, func(d *dieappsv1.DeploymentDieMut) {
			dietest.NoError(t, d.Spec().Template().Spec().Containers().At("app").Env().At("COLOR").ValueFrom().ConfigMapKeyRef().Key().Set("colour"))
		}},
		{"removing a list-map element", true, func(d *dieappsv1.DeploymentDieMut) {
			dietest.NoError(t, d.Spec().Template().Spec().Containers().Remove("app"))
			dietest.NoError(t, d.Spec().Template().Spec().Containers().At("sidecar").Image().Set("sidecar:1"))
		}},
	} {
		t.Run(c.name, func(t *testing.T) {
			typed, u := dieappsv1.DeploymentBlank, blank
			if c.server {
				typed = typed.DieStamp(func(d *dieappsv1.DeploymentDieMut) { dietest.NoError(t, die.FeedJSON(d, server)) })
				u = u.DieStamp(func(d *dieunstructured.ObjectMut) { dietest.NoError(t, die.FeedJSON(d, server)) })
			}
			typed = typed.DieStamp(c.change)
			mu := u.DieMutable()
			c.change(dieunstructured.JigAsMut[*dieappsv1.DeploymentDieMut](mu.Node()))

			want, err := die.ReleaseApplyPatchE(typed)
			if err != nil {
				t.Fatal(err)
			}
			got, err := die.ReleaseApplyPatchE(mu)
			if err != nil {
				t.Fatal(err)
			}
			var g, w any
			_ = json.Unmarshal(got, &g)
			_ = json.Unmarshal(want, &w)
			if !reflect.DeepEqual(g, w) {
				t.Errorf("unstructured applies:\n%s\ntyped applies:\n%s", got, want)
			}
		})
	}
}

// A custom resource's validation rejects null for a field that isn't
// nullable, so clearing one is left out of an apply, and the field stays
// another manager's; a nullable field is sent as null, and removed. A typed
// die knows from its markers, and an unstructured one from the CRD's schema.
func TestClearingCustomResourceFields(t *testing.T) {
	ns := namespace(t)
	ctx := context.Background()
	unstructuredBlank := moldedBlank(t, serverMolds(t), widgetGVK)

	applyU := func(t *testing.T, d *dieunstructured.Object, manager string) *dieunstructured.Object {
		t.Helper()
		body, err := die.ReleaseApplyPatchE(d)
		if err != nil {
			t.Fatal(err)
		}
		got, err := d.DieStampE(func(d *dieunstructured.ObjectMut) error {
			return dies.Patch(ctx, d, client.RawPatch(types.ApplyPatchType, body), client.FieldOwner(manager), client.ForceOwnership)
		})
		if err != nil {
			t.Fatalf("apply %s: %v", body, err)
		}
		return got
	}
	for _, c := range []struct {
		name  string
		clear func(manager string)
	}{
		{"typed", func(manager string) {
			d := widget(t, ns).DieStamp(func(d *widgetv1.WidgetDieMut) {
				dietest.NoError(t, d.Spec().Match().Exact().Clear())
				dietest.NoError(t, d.Spec().Match().Note().Clear())
			})
			body, err := die.ReleaseApplyPatchE(d)
			if err != nil {
				t.Fatal(err)
			}
			applyObject(t, d.DieMutable(), body, manager, false)
		}},
		{"unstructured", func(manager string) {
			applyU(t, unstructuredBlank.DieStamp(func(d *dieunstructured.ObjectMut) {
				d.SetNamespace(ns)
				d.SetName("web")
				match := d.Node().Field("spec").Field("match")
				dietest.NoError(t, match.Field("equal").Clear())
				dietest.NoError(t, match.Field("note").Clear())
			}), manager)
		}},
	} {
		t.Run(c.name, func(t *testing.T) {
			applyU(t, unstructuredBlank.DieStamp(func(d *dieunstructured.ObjectMut) {
				d.SetNamespace(ns)
				d.SetName("web")
				match := d.Node().Field("spec").Field("match")
				dietest.NoError(t, match.Field("equal").Set(true))
				dietest.NoError(t, match.Field("note").Set("set by matcher"))
			}), "matcher")

			c.clear("clearer-" + c.name)
			// what the server has, rather than a die with its intent on top
			got, err := dieunstructured.BlankFor("example.dev/v1", "Widget").DieStampE(func(d *dieunstructured.ObjectMut) error {
				return dies.Get(ctx, keyOf(ns, "web"), d)
			})
			if err != nil {
				t.Fatal(err)
			}
			match := got.Node().Field("spec").Field("match")
			if match.Field("note").Has() {
				t.Errorf("match %v, want note removed", match.Get())
			}
			if !match.Field("equal").Has() {
				t.Errorf("match %v, want equal left to matcher", match.Get())
			}
			assertOwns(t, got, "clearer-"+c.name, "", ".spec.match.note")
		})
	}
}
