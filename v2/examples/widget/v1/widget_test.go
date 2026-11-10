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
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"reconciler.io/dies/v2/die"
	"reconciler.io/dies/v2/die/base"
	dietesting "reconciler.io/dies/v2/die/testing"
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
	w := widgetv1.WidgetBlank.DieStamp(func(d *widgetv1.WidgetDieMut) {
		dietesting.NoError(t, d.Metadata().Name().Set("my-widget"))
		// a list map from the local type's own markers
		dietesting.NoError(t, d.Spec().Ports().At("http").Port().Set(8080))
		// an upstream type, through the core/v1 dies
		dietesting.NoError(t, d.Spec().Template().Spec().Containers().At("app").Image().Set("busybox"))
		// upstream metav1.Condition, through the meta/v1 dies
		dietesting.NoError(t, d.Status().Conditions().At("Ready").Status().Set("True"))
	})

	got, err := die.ReleaseApplyPatch(w, die.MainResource())
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
	got, err = die.ReleaseApplyPatch(w, die.Subresource("status"))
	assertJSON(t, got, err, `{
		"apiVersion": "example.dev/v1",
		"kind": "Widget",
		"metadata": {"name": "my-widget"},
		"status": {"conditions": [{"type": "Ready", "status": "True"}]}
	}`)

	released := *w.DieRelease()
	if released.Spec.Ports[0].Port != 8080 {
		t.Errorf("released %+v", released.Spec)
	}
}

func TestHandWrittenMethod(t *testing.T) {
	w := widgetv1.WidgetBlank.DieMutable()
	dietesting.NoError(t, w.Ready("Available"))
	got, err := die.ReleaseApplyPatch(w, die.Subresource("status"))
	assertJSON(t, got, err, `{
		"apiVersion": "example.dev/v1",
		"kind": "Widget",
		"status": {"conditions": [{"type": "Ready", "status": "True", "reason": "Available"}]}
	}`)
}

// WidgetGadget embeds gadget's GadgetSpec, whose die package overrides parts
// to atomic. The widget's dies and the runtime both follow that override.
func TestEmbeddedTypeFromAnotherDiePackage(t *testing.T) {
	w := widgetv1.WidgetBlank.DieStamp(func(d *widgetv1.WidgetDieMut) {
		dietesting.NoError(t, d.Spec().Gadget().Label().Set("knob"))
		dietesting.NoError(t, d.Spec().Gadget().Parts().SetValue([]gadgetv1.GadgetPart{{Name: "a"}, {Name: "b"}}...))
	})
	for _, f := range base.ManagedFields[widgetv1.WidgetGadget]() {
		if f.Wire == "parts" && f.Kind != base.FieldValue {
			t.Errorf("parts is %s, want a whole value", f.Kind)
		}
	}
	got, err := die.ReleaseApplyPatch(w)
	assertJSON(t, got, err, `{
		"apiVersion": "example.dev/v1",
		"kind": "Widget",
		"spec": {"gadget": {"label": "knob", "parts": [{"name": "a"}, {"name": "b"}]}}
	}`)
}

func TestRenamedAccessor(t *testing.T) {
	w := widgetv1.WidgetBlank.DieStamp(func(d *widgetv1.WidgetDieMut) {
		dietesting.NoError(t, d.Spec().Match().Exact().Set(true))
	})
	if !w.Spec().Match().Exact().Get() || !w.DieRelease().Spec.Match.Equal {
		t.Errorf("exact match not set")
	}
	// Equal is still the die method
	if !w.Spec().Match().Equal(w.Spec().Match()) {
		t.Errorf("a die should equal itself")
	}
	got, err := die.ReleaseApplyPatch(w)
	assertJSON(t, got, err, `{
		"apiVersion": "example.dev/v1",
		"kind": "Widget",
		"spec": {"match": {"equal": true}}
	}`)
}

func TestDuckIsNotRegistered(t *testing.T) {
	scheme := runtime.NewScheme()
	if err := widgetv1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	if _, ok := scheme.KnownTypes(widgetv1.SchemeGroupVersion)["Widget"]; !ok {
		t.Errorf("Widget is not registered")
	}
	// a duck type has no kind of its own
	if gvks, _, err := scheme.ObjectKinds(&widgetv1.TemplatedDieMut{}); err == nil {
		t.Errorf("Templated is registered as %v", gvks)
	}
}

func TestDefaultAndValidate(t *testing.T) {
	ctx := context.Background()
	w := widgetv1.WidgetBlank.DieStamp(func(d *widgetv1.WidgetDieMut) {
		dietesting.NoError(t, d.Spec().Ports().At("http").Port().Set(8080))
		dietesting.NoError(t, d.Spec().Ports().At("admin").Port().Set(70000))
		dietesting.NoError(t, d.Spec().Ports().At("web").Name().Set("web"))
		dietesting.NoError(t, d.Spec().Ports().At("SSH").Port().Set(22))
		dietesting.NoError(t, d.Spec().Template().Spec().Containers().At("app").Image().Set("busybox"))
	})

	// a port's markers and its die validate it, through the widget and its
	// spec, which have no methods of their own
	assertErrors(t, validated(t)(die.ValidateCreate(ctx, w)),
		"spec.ports[1].port: Invalid value: 70000: must be less than or equal to 65535",
		`spec.ports[3].name: Invalid value: "SSH": a lowercase RFC 1123 label must consist of lower case alphanumeric characters or '-', and must start and end with an alphanumeric character (e.g. 'my-name',  or '123-abc', regex used for validation is '[a-z0-9]([-a-z0-9]*[a-z0-9])?')`,
		"spec.ports[3].port: Invalid value: 22: is reserved for ssh",
	)

	w = w.DieStamp(func(d *widgetv1.WidgetDieMut) {
		dietesting.NoError(t, die.Default(ctx, d))
	})
	if got := w.Spec().Ports().At("web").Port().Get(); got != 80 {
		t.Errorf("defaulted port %d, want 80", got)
	}
	if got := w.Spec().Ports().At("http").Port().Get(); got != 8080 {
		t.Errorf("set port %d was defaulted", got)
	}

	// by kubebuilder's convention, a field without omitempty is required:
	// the package has kubebuilder's markers of custom resources
	m := widgetv1.WidgetBlank.DieMutable()
	dietesting.NoError(t, die.FeedYAML(m, []byte(`
spec:
  ports:
  - port: 8080
`)))
	assertErrors(t, validated(t)(die.ValidateCreate(ctx, m.DieImmutable())),
		// a die's errors, then those of the dies within it
		"spec.template: Required value",
		"spec.ports[0].name: Required value",
	)

	// a die's error stops defaulting, at its location
	err := die.Default(ctx, widgetv1.WidgetBlank.DieStamp(func(d *widgetv1.WidgetDieMut) {
		dietesting.NoError(t, d.Spec().Match().Note().Set("fail"))
	}).DieMutable())
	if err == nil || err.Error() != "defaulting .spec.match: a note of fail" {
		t.Errorf("got %v", err)
	}
}

func assertErrors(t *testing.T, errs field.ErrorList, want ...string) {
	t.Helper()
	if len(errs) != len(want) {
		t.Fatalf("got %v, want %q", errs, want)
	}
	for i, err := range errs {
		if err.Error() != want[i] {
			t.Errorf("error %d: got %q, want %q", i, err.Error(), want[i])
		}
	}
}

func TestValidateUpdate(t *testing.T) {
	ctx := context.Background()
	old := widgetv1.WidgetBlank.DieStamp(func(d *widgetv1.WidgetDieMut) {
		dietesting.NoError(t, d.Spec().Ports().At("http").Port().Set(8080))
		dietesting.NoError(t, d.Spec().Ports().At("admin").Port().Set(70000))
		dietesting.NoError(t, d.Spec().Ports().At("metrics").Port().Set(9090))
		dietesting.NoError(t, d.Spec().Template().Spec().Containers().At("app").Image().Set("busybox"))
	})
	// the old widget's admin port is out of range, as a create reports
	if errs := validated(t)(die.ValidateCreate(ctx, old)); len(errs) != 1 {
		t.Errorf("create: %v", errs)
	}

	updated := old.DieStamp(func(d *widgetv1.WidgetDieMut) {
		// reordered, and http's number changed
		dietesting.NoError(t, d.Spec().Ports().Clear())
		dietesting.NoError(t, d.Spec().Ports().At("admin").Port().Set(70000))
		dietesting.NoError(t, d.Spec().Ports().At("http").Port().Set(8081))
		dietesting.NoError(t, d.Spec().Match().Exact().Set(true))
	})
	assertErrors(t, validated(t)(die.ValidateUpdate(ctx, updated, old)),
		// +k8s:update=NoRemoveItem: metrics was removed
		"spec.ports: Forbidden: item may not be removed",
		// the admin port is unchanged, so its range isn't checked again, and
		// the http port's die has its old value, by its name
		"spec.ports[1].port: Forbidden: a named port's number may not change",
		// +k8s:immutable
		"spec.match.equal: Invalid value: null: field is immutable",
	)

	// a widget being deleted has no schema to meet, nor does its die check
	// its ports
	if errs := validated(t)(die.ValidateDelete(ctx, old.DieStamp(func(d *widgetv1.WidgetDieMut) {
		dietesting.NoError(t, d.Spec().Ports().At("ssh").Port().Set(22))
	}))); len(errs) > 0 {
		t.Errorf("delete: %v", errs)
	}
}

func TestValidateError(t *testing.T) {
	w := widgetv1.WidgetBlank.DieStamp(func(d *widgetv1.WidgetDieMut) {
		dietesting.NoError(t, d.Spec().Ports().At("ssh").Port().Set(22))
		dietesting.NoError(t, d.Spec().Match().Note().Set("unknown"))
	})
	// a validator's error stops validation, at its location, with the errors
	// found until then
	errs, err := die.ValidateCreate(context.Background(), w)
	if err == nil || err.Error() != "validating .spec.match: a note of unknown can't be validated" {
		t.Errorf("got %v", err)
	}
	if len(errs) == 0 || errs[len(errs)-1].Field != "spec.ports[0].port" {
		t.Errorf("got %v", errs)
	}
}

// validated fails the test if validation fails, and returns its errors.
func validated(t *testing.T) func(field.ErrorList, error) field.ErrorList {
	t.Helper()
	return func(errs field.ErrorList, err error) field.ErrorList {
		t.Helper()
		if err != nil {
			t.Fatalf("validating: %v", err)
		}
		return errs
	}
}
