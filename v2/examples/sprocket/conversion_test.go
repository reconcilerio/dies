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

package sprocket_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"reconciler.io/dies/v2/die/conversion"
	"reconciler.io/dies/v2/die/conversion/webhook"
	"reconciler.io/dies/v2/die/jigs"
	dietesting "reconciler.io/dies/v2/die/testing"
	"reconciler.io/dies/v2/examples/sprocket"
	sprocketv1 "reconciler.io/dies/v2/examples/sprocket/v1"
	sprocketv2 "reconciler.io/dies/v2/examples/sprocket/v2"
	widgetv1 "reconciler.io/dies/v2/examples/widget/v1"
)

func TestImageReferenceIsLossless(t *testing.T) {
	sources := []sprocketv2.SprocketImage{
		{},
		{Repository: "busybox"},
		{Repository: "busybox", Tag: "1.36"},
		{Repository: "registry.example:5000/team/app"},
		{Repository: "registry.example:5000/team/app", Tag: "v2"},
	}
	refs := []string{"", "busybox", "busybox:1.36", "registry.example:5000/team/app", "registry.example:5000/team/app:v2", "app:"}
	if err := jigs.CheckTransform(sprocket.ImageReference, sources, refs); err != nil {
		t.Fatal(err)
	}
}

// convert converts an object, as JSON, to a version.
func convert(t *testing.T, obj string, to schema.GroupVersion) string {
	t.Helper()
	out, err := sprocket.Conversion.Convert([]byte(obj), to)
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}

// equalJSON fails the test unless two JSON documents have the same value.
func equalJSON(t *testing.T, want, got string) {
	t.Helper()
	var w, g any
	if err := json.Unmarshal([]byte(want), &w); err != nil {
		t.Fatalf("want: %v", err)
	}
	if err := json.Unmarshal([]byte(got), &g); err != nil {
		t.Fatalf("got: %v", err)
	}
	if diff := cmp.Diff(w, g); diff != "" {
		t.Errorf("(-want, +got):\n%s", diff)
	}
}

const v1Sprocket = `{
	"apiVersion": "example.dev/v1",
	"kind": "Sprocket",
	"metadata": {"name": "chain", "namespace": "default", "labels": {"tier": "drive"}},
	"spec": {"size": 42, "image": "registry.example:5000/sprocket:v1", "lubricant": "oil"},
	"status": {"observedGeneration": 3}
}`

// A v1 sprocket converts to v2, which stashes the lubricant it has no place
// for, and back to v1 unchanged.
func TestConvertV1RoundTrips(t *testing.T) {
	v2 := convert(t, v1Sprocket, sprocketv2.SchemeGroupVersion)
	equalJSON(t, `{
		"apiVersion": "example.dev/v2",
		"kind": "Sprocket",
		"metadata": {
			"name": "chain", "namespace": "default", "labels": {"tier": "drive"},
			"annotations": {"conversion.dies.reconciler.io/v1": "{\"spec\":{\"lubricant\":\"oil\"}}"}
		},
		"spec": {"teeth": 42, "image": {"repository": "registry.example:5000/sprocket", "tag": "v1"}},
		"status": {"observedGeneration": 3}
	}`, v2)

	equalJSON(t, v1Sprocket, convert(t, v2, sprocketv1.SchemeGroupVersion))
}

const v2Sprocket = `{
	"apiVersion": "example.dev/v2",
	"kind": "Sprocket",
	"metadata": {"name": "chain", "namespace": "default", "annotations": {"owner": "team"}},
	"spec": {"teeth": 12, "image": {"repository": "busybox"}, "paused": true}
}`

// A v2 sprocket converts to v1, which stashes paused, and back unchanged.
func TestConvertV2RoundTrips(t *testing.T) {
	v1 := convert(t, v2Sprocket, sprocketv1.SchemeGroupVersion)
	equalJSON(t, `{
		"apiVersion": "example.dev/v1",
		"kind": "Sprocket",
		"metadata": {
			"name": "chain", "namespace": "default",
			"annotations": {"owner": "team", "conversion.dies.reconciler.io/v2": "{\"spec\":{\"paused\":true}}"}
		},
		"spec": {"size": 12, "image": "busybox"}
	}`, v1)

	equalJSON(t, v2Sprocket, convert(t, v1, sprocketv2.SchemeGroupVersion))
}

// What a version doesn't set, its conversion doesn't either: a sprocket with
// no spec has none in the other version, not one of zero values.
func TestConvertKeepsUnsetUnset(t *testing.T) {
	bare := `{"apiVersion": "example.dev/v1", "kind": "Sprocket", "metadata": {"name": "bare"}}`
	v2 := convert(t, bare, sprocketv2.SchemeGroupVersion)
	equalJSON(t, `{"apiVersion": "example.dev/v2", "kind": "Sprocket", "metadata": {"name": "bare"}}`, v2)
	equalJSON(t, bare, convert(t, v2, sprocketv1.SchemeGroupVersion))
}

// A value set to its zero value is a value, unlike one that's unset: it
// converts as one.
func TestConvertKeepsZeroValues(t *testing.T) {
	zero := `{"apiVersion": "example.dev/v1", "kind": "Sprocket", "metadata": {"name": "zero"}, "spec": {"size": 0}}`
	v2 := convert(t, zero, sprocketv2.SchemeGroupVersion)
	equalJSON(t, `{"apiVersion": "example.dev/v2", "kind": "Sprocket", "metadata": {"name": "zero"}, "spec": {"teeth": 0}}`, v2)
	equalJSON(t, zero, convert(t, v2, sprocketv1.SchemeGroupVersion))
}

// A stash fills only what the conversion leaves unset: a value changed in
// the hub since the stash was made is kept.
func TestStashDoesNotOverrideConvertedValues(t *testing.T) {
	v2 := `{
		"apiVersion": "example.dev/v2", "kind": "Sprocket",
		"metadata": {"name": "chain", "annotations": {"conversion.dies.reconciler.io/v1": "{\"spec\":{\"size\":0,\"lubricant\":\"oil\"}}"}},
		"spec": {"teeth": 7}
	}`
	equalJSON(t, `{
		"apiVersion": "example.dev/v1", "kind": "Sprocket",
		"metadata": {"name": "chain"},
		"spec": {"size": 7, "lubricant": "oil"}
	}`, convert(t, v2, sprocketv1.SchemeGroupVersion))
}

func TestConvertRefusesOtherKinds(t *testing.T) {
	if _, err := sprocket.Conversion.Convert([]byte(`{"apiVersion": "example.dev/v1", "kind": "Widget"}`), sprocketv2.SchemeGroupVersion); err == nil {
		t.Error("expected an error converting a Widget")
	}
	if _, err := sprocket.Conversion.Convert([]byte(v1Sprocket), schema.GroupVersion{Group: "example.dev", Version: "v3"}); err == nil {
		t.Error("expected an error converting to a version without a conversion")
	}
}

func TestWebhook(t *testing.T) {
	server := httptest.NewServer(webhook.Handler(sprocket.Conversion))
	defer server.Close()

	review := func(t *testing.T, objs ...string) *apiextensionsv1.ConversionResponse {
		t.Helper()
		req := apiextensionsv1.ConversionReview{
			TypeMeta: metav1.TypeMeta{APIVersion: "apiextensions.k8s.io/v1", Kind: "ConversionReview"},
			Request:  &apiextensionsv1.ConversionRequest{UID: "1234", DesiredAPIVersion: "example.dev/v2"},
		}
		for _, o := range objs {
			req.Request.Objects = append(req.Request.Objects, runtime.RawExtension{Raw: []byte(o)})
		}
		body, err := json.Marshal(req)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := http.Post(server.URL, "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = resp.Body.Close() }()
		var out apiextensionsv1.ConversionReview
		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			t.Fatal(err)
		}
		if out.Response == nil || out.Response.UID != "1234" {
			t.Fatalf("the response is not for the request: %+v", out.Response)
		}
		return out.Response
	}

	t.Run("converts", func(t *testing.T) {
		resp := review(t, v1Sprocket, v2Sprocket)
		if resp.Result.Status != metav1.StatusSuccess || len(resp.ConvertedObjects) != 2 {
			t.Fatalf("unexpected response: %+v", resp)
		}
		equalJSON(t, convert(t, v1Sprocket, sprocketv2.SchemeGroupVersion), string(resp.ConvertedObjects[0].Raw))
		equalJSON(t, v2Sprocket, string(resp.ConvertedObjects[1].Raw))
	})
	t.Run("fails", func(t *testing.T) {
		resp := review(t, v1Sprocket, `{"apiVersion": "example.dev/v1", "kind": "Widget"}`)
		if resp.Result.Status != metav1.StatusFailure || resp.ConvertedObjects != nil {
			t.Fatalf("unexpected response: %+v", resp)
		}
	})
}

// A converter can stash under another annotation, followed by the version,
// as the default is. A stash under the default isn't restored by it.
func TestConvertWithStashAnnotation(t *testing.T) {
	c := sprocket.Conversion.WithStashAnnotation("sprockets.example.dev/dropped-")
	if got, want := c.StashAnnotation(sprocketv1.SchemeGroupVersion), "sprockets.example.dev/dropped-v1"; got != want {
		t.Errorf("got %s, want %s", got, want)
	}
	if got, want := sprocket.Conversion.StashAnnotation(sprocketv1.SchemeGroupVersion), conversion.StashAnnotation+"v1"; got != want {
		t.Errorf("the default: got %s, want %s", got, want)
	}

	v2, err := c.Convert([]byte(v1Sprocket), sprocketv2.SchemeGroupVersion)
	if err != nil {
		t.Fatal(err)
	}
	equalJSON(t, `{
		"apiVersion": "example.dev/v2",
		"kind": "Sprocket",
		"metadata": {
			"name": "chain", "namespace": "default", "labels": {"tier": "drive"},
			"annotations": {"sprockets.example.dev/dropped-v1": "{\"spec\":{\"lubricant\":\"oil\"}}"}
		},
		"spec": {"teeth": 42, "image": {"repository": "registry.example:5000/sprocket", "tag": "v1"}},
		"status": {"observedGeneration": 3}
	}`, string(v2))
	v1, err := c.Convert(v2, sprocketv1.SchemeGroupVersion)
	if err != nil {
		t.Fatal(err)
	}
	equalJSON(t, v1Sprocket, string(v1))

	// the default converter leaves the other's stash as an annotation
	v1, err = sprocket.Conversion.Convert(v2, sprocketv1.SchemeGroupVersion)
	if err != nil {
		t.Fatal(err)
	}
	var obj map[string]any
	if err := json.Unmarshal(v1, &obj); err != nil {
		t.Fatal(err)
	}
	spec := obj["spec"].(map[string]any)
	if _, ok := spec["lubricant"]; ok {
		t.Errorf("the default converter restored another prefix's stash: %s", v1)
	}
	ann := obj["metadata"].(map[string]any)["annotations"].(map[string]any)
	if _, ok := ann["sprockets.example.dev/dropped-v1"]; !ok {
		t.Errorf("the other prefix's stash is gone: %s", v1)
	}
}

// panics fails the test unless fn panics.
func panics(t *testing.T, what string, fn func()) {
	t.Helper()
	defer func() {
		t.Helper()
		if recover() == nil {
			t.Errorf("expected a panic for %s", what)
		}
	}()
	fn()
}

func TestStashAnnotationMustBeValid(t *testing.T) {
	for _, prefix := range []string{"not valid/", "a/b/", "-dropped-", strings.Repeat("x", 62)} {
		panics(t, prefix, func() { sprocket.Conversion.WithStashAnnotation(prefix) })
	}
}

// A version's group, version and kind are its type's, and converted objects start
// empty: what's stamped on the blanks a conversion is made from isn't part of
// what it converts to.
func TestConversionIgnoresTheBlanksValues(t *testing.T) {
	v1 := sprocketv1.SprocketBlank.DieStamp(func(d *sprocketv1.SprocketDieMut) {
		d.GetObjectKind().SetGroupVersionKind(schema.GroupVersionKind{Group: "wrong.dev", Version: "v9", Kind: "Gear"})
		dietesting.NoError(t, d.Spec().Lubricant().Set("stamped"))
	})
	v2 := sprocketv2.SprocketBlank.DieStamp(func(d *sprocketv2.SprocketDieMut) {
		dietesting.NoError(t, d.Spec().Paused().Set(true))
	})
	c := conversion.New(v2).
		AddSpoke(v1, jigs.Bind(v1.Spec().Size(), v2.Spec().Teeth()))

	bare := `{"apiVersion": "example.dev/v1", "kind": "Sprocket", "metadata": {"name": "bare"}, "spec": {"size": 3}}`
	out, err := c.Convert([]byte(bare), sprocketv2.SchemeGroupVersion)
	if err != nil {
		t.Fatal(err)
	}
	equalJSON(t, `{"apiVersion": "example.dev/v2", "kind": "Sprocket", "metadata": {"name": "bare"}, "spec": {"teeth": 3}}`, string(out))
	if out, err = c.Convert(out, sprocketv1.SchemeGroupVersion); err != nil {
		t.Fatal(err)
	}
	equalJSON(t, bare, string(out))
}

func TestDefinitionsThatCantConvertPanic(t *testing.T) {
	hub := conversion.New(sprocketv2.SprocketBlank)
	// Templated is an object, but has no kind
	panics(t, "a hub without a group, version and kind", func() { conversion.New(widgetv1.TemplatedBlank) })
	panics(t, "a spoke without a group, version and kind", func() { hub.AddSpoke(widgetv1.TemplatedBlank) })
	panics(t, "a spoke of another kind", func() { hub.AddSpoke(widgetv1.WidgetBlank) })
	panics(t, "a spoke of the hub's own version", func() { hub.AddSpoke(sprocketv2.SprocketBlank) })
	panics(t, "a second conversion of a version", func() { sprocket.Conversion.AddSpoke(sprocketv1.SprocketBlank) })
	panics(t, "a binding to another hub", func() {
		hub.AddSpoke(sprocketv1.SprocketBlank, jigs.Bind(sprocketv1.SprocketBlank.Metadata().Name(), widgetv1.WidgetBlank.Metadata().Name()))
	})
}

// Each method returns a new converter, leaving the one it's called on as it
// was, so a converter can be derived from another that's shared.
func TestDefinitionsDontChangeTheirConverter(t *testing.T) {
	hub := conversion.New(sprocketv2.SprocketBlank)
	withV1 := hub.AddSpoke(sprocketv1.SprocketBlank)
	if _, err := hub.Convert([]byte(v1Sprocket), sprocketv2.SchemeGroupVersion); err == nil {
		t.Error("adding a spoke changed the converter it was added to")
	}
	if _, err := withV1.Convert([]byte(v1Sprocket), sprocketv2.SchemeGroupVersion); err != nil {
		t.Error(err)
	}
	withV1.WithStashAnnotation("sprockets.example.dev/dropped-")
	if got, want := withV1.StashAnnotation(sprocketv1.SchemeGroupVersion), conversion.StashAnnotation+"v1"; got != want {
		t.Errorf("setting the prefix changed the converter: got %s, want %s", got, want)
	}
}

// A version the converter doesn't know is refused, even when it's the
// version the object has.
func TestConvertRefusesUnknownVersions(t *testing.T) {
	other := `{"apiVersion": "other.dev/v1", "kind": "Sprocket"}`
	if _, err := sprocket.Conversion.Convert([]byte(other), schema.GroupVersion{Group: "other.dev", Version: "v1"}); err == nil {
		t.Error("expected an error for a version of another group")
	}
}
