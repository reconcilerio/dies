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
	"slices"
	"testing"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	diemetav1 "reconciler.io/dies/v2/apis/meta/v1"
	dieunstructured "reconciler.io/dies/v2/die/unstructured"
	widgetv1 "reconciler.io/dies/v2/examples/widget/v1"
)

// Dies for a custom resource, declared locally, and as unstructured dies,
// against its CRD, which declares the list maps the dies' markers do.

func widget(ns string) *widgetv1.WidgetDie {
	return widgetv1.WidgetBlank.DieWith(func(d *widgetv1.WidgetDieMut) {
		d.Metadata().Namespace().Set(ns)
		d.Metadata().Name().Set("web")
	})
}

// applyObject sends an apply patch, to the status subresource for status,
// and decodes the result into obj.
func applyObject(t *testing.T, obj client.Object, body []byte, manager string, status bool) {
	t.Helper()
	var err error
	if status {
		err = dies.Status().Patch(context.Background(), obj, client.RawPatch(types.ApplyPatchType, body), client.FieldOwner(manager), client.ForceOwnership)
	} else {
		err = dies.Patch(context.Background(), obj, client.RawPatch(types.ApplyPatchType, body), client.FieldOwner(manager), client.ForceOwnership)
	}
	if err != nil {
		t.Fatalf("apply %s: %v", body, err)
	}
}

func TestWidgetListMapsMergeByKey(t *testing.T) {
	ns := namespace(t)
	for _, port := range []struct {
		name    string
		port    int32
		manager string
	}{{"http", 80, "web"}, {"metrics", 9090, "monitor"}} {
		d := widget(ns).DieWith(func(d *widgetv1.WidgetDieMut) {
			d.Spec().Ports().At(port.name).Port().Set(port.port)
		})
		body, err := d.DieReleaseApplyPatchE()
		if err != nil {
			t.Fatal(err)
		}
		applyObject(t, d.DieMutable(), body, port.manager, false)
	}

	got := widgetv1.WidgetBlank.DieMutable()
	if err := dies.Get(context.Background(), keyOf(ns, "web"), got); err != nil {
		t.Fatal(err)
	}
	if keys := got.Spec().Ports().Keys(); !slices.Equal(keys, []string{"http", "metrics"}) {
		t.Errorf("ports %v", keys)
	}
	assertOwned(t, got, "monitor", "", `.spec.ports[name="metrics"].name`, `.spec.ports[name="metrics"].port`)
}

func TestWidgetStatusConditions(t *testing.T) {
	ns := namespace(t)
	d := widget(ns)
	body, _ := d.DieReleaseApplyPatchE()
	applyObject(t, d.DieMutable(), body, "web", false)

	for _, c := range []struct{ condition, manager string }{{"Ready", "controller"}, {"Healthy", "prober"}} {
		status := widget(ns).DieWith(func(d *widgetv1.WidgetDieMut) {
			d.Status().Conditions().At(c.condition).DieWith(func(cd *diemetav1.ConditionDieMut) {
				cd.Status().Set(diemetav1.ConditionTrue)
				cd.Reason().Set("Checked")
			})
		}).DieSubresource("status")
		body, err := status.DieReleaseApplyPatchE()
		if err != nil {
			t.Fatal(err)
		}
		applyObject(t, status.DieMutable(), body, c.manager, true)
	}

	got := widgetv1.WidgetBlank.DieMutable()
	if err := dies.Get(context.Background(), keyOf(ns, "web"), got); err != nil {
		t.Fatal(err)
	}
	conditions := got.Status().Conditions()
	if conditions.Len() != 2 || conditions.At("Ready").Reason().Get() != "Checked" || conditions.At("Healthy").Status().Get() != diemetav1.ConditionTrue {
		t.Errorf("conditions %+v", got.Status().DieRelease().Conditions)
	}
}

// An unstructured die needs no Go type, or registration: it's applied,
// read and listed as the CRD's kind.
func TestUnstructuredWidget(t *testing.T) {
	ns := namespace(t)
	ctx := context.Background()
	d := dieunstructured.BlankFor("example.dev/v1", "Widget").DieWith(func(d *dieunstructured.ObjectMut) {
		d.SetNamespace(ns)
		d.SetName("web")
		d.Node().Field("spec").Field("ports").Item("name", "http").Field("port").Set(int64(80))
	})
	body, err := d.DieReleaseApplyPatchE()
	if err != nil {
		t.Fatal(err)
	}
	applyObject(t, d.DieMutable(), body, "web", false)

	// another manager's port, by key, from what the server has
	read := dieunstructured.BlankFor("example.dev/v1", "Widget").DieMutable()
	if err := dies.Get(ctx, keyOf(ns, "web"), read); err != nil {
		t.Fatal(err)
	}
	more := read.DieImmutable().DieWith(func(d *dieunstructured.ObjectMut) {
		d.Node().Field("spec").Field("ports").Item("name", "metrics").Field("port").Set(int64(9090))
	})
	body, err = more.DieReleaseApplyPatchE()
	if err != nil {
		t.Fatal(err)
	}
	applyObject(t, more.DieMutable(), body, "monitor", false)

	list := dieunstructured.NewList("example.dev/v1", "WidgetList")
	if err := dies.List(ctx, list, client.InNamespace(ns)); err != nil {
		t.Fatal(err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("items %d", len(list.Items))
	}
	ports := list.Items[0].Node().Field("spec").Field("ports")
	if ports.Item("name", "http").Field("port").Get() != int64(80) || ports.Item("name", "metrics").Field("port").Get() != int64(9090) {
		t.Errorf("ports %v", ports.Get())
	}
	assertOwned(t, list.Items[0], "monitor", "", `.spec.ports[name="metrics"].name`, `.spec.ports[name="metrics"].port`)
}
