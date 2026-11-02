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

package v1_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"reconciler.io/dies/v2/die/base"
	gadgetv1 "reconciler.io/dies/v2/examples/gadget/v1"
	widgetv1 "reconciler.io/dies/v2/examples/widget/v1"
)

func assertJSON(t *testing.T, got []byte, err error, want string) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var g, w any
	if err := json.Unmarshal(got, &g); err != nil {
		t.Fatalf("invalid json %s: %v", got, err)
	}
	if err := json.Unmarshal([]byte(want), &w); err != nil {
		t.Fatalf("invalid expected json %s: %v", want, err)
	}
	if !reflect.DeepEqual(g, w) {
		t.Errorf("json mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestLocalTypeDies(t *testing.T) {
	w := widgetv1.WidgetBlank.DieWith(func(d *widgetv1.WidgetDie) {
		d.Metadata().Name().Set("my-widget")
		// a list map from the local type's own markers
		d.Spec().Ports().At("http").Port().Set(8080)
		// an upstream type, through the core/v1 dies
		d.Spec().Template().Spec().Containers().At("app").Image().Set("busybox")
		// upstream metav1.Condition, through the meta/v1 dies
		d.Status().Conditions().At("Ready").Status().Set("True")
	})

	got, err := w.DieMain().DieReleaseApplyPatchE()
	assertJSON(t, got, err, `{
		"apiVersion": "example.dev/v1",
		"kind": "Widget",
		"metadata": {"name": "my-widget"},
		"spec": {
			"ports": [{"name": "http", "port": 8080}],
			"template": {"spec": {"containers": [{"name": "app", "image": "busybox"}]}}
		}
	}`)

	// +kubebuilder:subresource:status splits off the status
	got, err = w.DieSubresource("status").DieReleaseApplyPatchE()
	assertJSON(t, got, err, `{
		"apiVersion": "example.dev/v1",
		"kind": "Widget",
		"metadata": {"name": "my-widget"},
		"status": {"conditions": [{"type": "Ready", "status": "True"}]}
	}`)

	var released widgetv1.Widget = w.DieRelease()
	if released.Spec.Ports[0].Port != 8080 {
		t.Errorf("released %+v", released.Spec)
	}
}

func TestHandWrittenMethod(t *testing.T) {
	w := widgetv1.WidgetBlank.Ready("Available")
	got, err := w.DieSubresource("status").DieReleaseApplyPatchE()
	assertJSON(t, got, err, `{
		"apiVersion": "example.dev/v1",
		"kind": "Widget",
		"status": {"conditions": [{"type": "Ready", "status": "True", "reason": "Available"}]}
	}`)
}

// WidgetGadget embeds gadget's GadgetSpec, whose die package overrides parts
// to atomic. The widget's dies and the runtime both follow that override.
func TestEmbeddedTypeFromAnotherDiePackage(t *testing.T) {
	w := widgetv1.WidgetBlank.DieWith(func(d *widgetv1.WidgetDie) {
		d.Spec().Gadget().Label().Set("knob")
		d.Spec().Gadget().Parts().Set([]gadgetv1.GadgetPart{{Name: "a"}, {Name: "b"}})
	})
	for _, f := range base.ManagedFields[widgetv1.WidgetGadget]() {
		if f.Wire == "parts" && f.Kind != base.FieldValue {
			t.Errorf("parts is %s, want a whole value", f.Kind)
		}
	}
	got, err := w.DieReleaseApplyPatchE()
	assertJSON(t, got, err, `{
		"apiVersion": "example.dev/v1",
		"kind": "Widget",
		"spec": {"gadget": {"label": "knob", "parts": [{"name": "a"}, {"name": "b"}]}}
	}`)
}

func TestRenamedAccessor(t *testing.T) {
	w := widgetv1.WidgetBlank.DieWith(func(d *widgetv1.WidgetDie) {
		d.Spec().Match().Exact().Set(true)
	})
	if !w.Spec().Match().Exact().Get() || !w.DieRelease().Spec.Match.Equal {
		t.Errorf("exact match not set")
	}
	// Equal is still the die method
	if !w.Spec().Match().Equal(w.Spec().Match()) {
		t.Errorf("a die should equal itself")
	}
	got, err := w.DieReleaseApplyPatchE()
	assertJSON(t, got, err, `{
		"apiVersion": "example.dev/v1",
		"kind": "Widget",
		"spec": {"match": {"equal": true}}
	}`)
}
