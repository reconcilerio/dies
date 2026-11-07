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
	"strings"
	"testing"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"reconciler.io/dies/v2/die"
	dietesting "reconciler.io/dies/v2/die/testing"
	"reconciler.io/dies/v2/examples/sprocket"
	sprocketv1 "reconciler.io/dies/v2/examples/sprocket/v1"
	sprocketv2 "reconciler.io/dies/v2/examples/sprocket/v2"
)

// Sprockets are stored as v2 and served as v1 too: the server converts them
// with the webhook the suite serves, sprocket.Conversion. What each version
// has no place for survives a trip through the other, in its stash.

func sprocketV1(t testing.TB, ns string) *sprocketv1.SprocketDieMut {
	return sprocketv1.SprocketBlank.DieStamp(func(d *sprocketv1.SprocketDieMut) {
		dietesting.NoError(t, d.Metadata().Namespace().Set(ns))
		dietesting.NoError(t, d.Metadata().Name().Set("chain"))
		dietesting.NoError(t, d.Spec().Size().Set(42))
		dietesting.NoError(t, d.Spec().Image().Set("registry.example:5000/sprocket:v1"))
		dietesting.NoError(t, d.Spec().Lubricant().Set("oil"))
	}).DieMutable()
}

func getV1(t *testing.T, ns string) *sprocketv1.SprocketDieMut {
	t.Helper()
	d := sprocketv1.SprocketBlank.DieMutable()
	if err := dies.Get(context.Background(), client.ObjectKey{Namespace: ns, Name: "chain"}, d); err != nil {
		t.Fatal(err)
	}
	return d
}

func getV2(t *testing.T, ns string) *sprocketv2.SprocketDieMut {
	t.Helper()
	d := sprocketv2.SprocketBlank.DieMutable()
	if err := dies.Get(context.Background(), client.ObjectKey{Namespace: ns, Name: "chain"}, d); err != nil {
		t.Fatal(err)
	}
	return d
}

func expect[T comparable](t *testing.T, what string, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %v, want %v", what, got, want)
	}
}

// A sprocket created as v1 is stored as v2, with v1's lubricant stashed, and
// reads back as v1 as it was created.
func TestSprocketCreatedAsV1(t *testing.T) {
	ns := namespace(t)
	if err := dies.Create(context.Background(), sprocketV1(t, ns)); err != nil {
		t.Fatal(err)
	}

	v2 := getV2(t, ns)
	expect(t, "v2 teeth", v2.Spec().Teeth().Get(), 42)
	expect(t, "v2 repository", v2.Spec().Image().Repository().Get(), "registry.example:5000/sprocket")
	expect(t, "v2 tag", v2.Spec().Image().Tag().Get(), "v1")
	expect(t, "v2 paused is unset", v2.Spec().Paused().Has(), false)
	stash := v2.Metadata().Annotations().At(sprocket.Conversion.StashAnnotation(sprocketv1.SchemeGroupVersion))
	expect(t, "v1 stash", stash, `{"spec":{"lubricant":"oil"}}`)

	v1 := getV1(t, ns)
	expect(t, "v1 size", v1.Spec().Size().Get(), 42)
	expect(t, "v1 image", v1.Spec().Image().Get(), "registry.example:5000/sprocket:v1")
	expect(t, "v1 lubricant", v1.Spec().Lubricant().Get(), "oil")
	expect(t, "v1 annotations", v1.Metadata().Annotations().Len(), 0)
}

// Updates through each version keep what the other has no place for.
func TestSprocketUpdatedThroughEachVersion(t *testing.T) {
	ns := namespace(t)
	ctx := context.Background()
	if err := dies.Create(ctx, sprocketV1(t, ns)); err != nil {
		t.Fatal(err)
	}

	// v2 pauses the sprocket, which v1 has no place for, and changes its
	// teeth
	v2 := getV2(t, ns)
	dietesting.NoError(t, v2.Spec().Paused().Set(true))
	dietesting.NoError(t, v2.Spec().Teeth().Set(50))
	if err := dies.Update(ctx, v2); err != nil {
		t.Fatal(err)
	}

	// v1 sees the teeth, keeps its lubricant, and stashes paused
	v1 := getV1(t, ns)
	expect(t, "v1 size", v1.Spec().Size().Get(), 50)
	expect(t, "v1 lubricant", v1.Spec().Lubricant().Get(), "oil")
	stash := v1.Metadata().Annotations().At(sprocket.Conversion.StashAnnotation(sprocketv2.SchemeGroupVersion))
	expect(t, "v2 stash", stash, `{"spec":{"paused":true}}`)

	// v1 changes the image and the lubricant, sending the stash back as it
	// read it
	dietesting.NoError(t, v1.Spec().Image().Set("busybox:1.36"))
	dietesting.NoError(t, v1.Spec().Lubricant().Set("grease"))
	if err := dies.Update(ctx, v1); err != nil {
		t.Fatal(err)
	}

	v2 = getV2(t, ns)
	expect(t, "v2 paused", v2.Spec().Paused().Get(), true)
	expect(t, "v2 teeth", v2.Spec().Teeth().Get(), 50)
	expect(t, "v2 repository", v2.Spec().Image().Repository().Get(), "busybox")
	expect(t, "v2 tag", v2.Spec().Image().Tag().Get(), "1.36")
	expect(t, "v2 annotations", strings.Join(v2.Metadata().Annotations().Keys(), ","), sprocket.Conversion.StashAnnotation(sprocketv1.SchemeGroupVersion))
	v1 = getV1(t, ns)
	expect(t, "v1 lubricant", v1.Spec().Lubricant().Get(), "grease")
}

// Server-side apply through v1, by two managers, merges as it does for one
// version: the server converts the apply, and the managed fields, through
// the webhook.
func TestSprocketAppliedAsV1(t *testing.T) {
	ns := namespace(t)
	sizer := sprocketv1.SprocketBlank.DieStamp(func(d *sprocketv1.SprocketDieMut) {
		dietesting.NoError(t, d.Metadata().Namespace().Set(ns))
		dietesting.NoError(t, d.Metadata().Name().Set("chain"))
		dietesting.NoError(t, d.Spec().Size().Set(42))
	})
	applyObject(t, sizer.DieMutable(), die.MustReleaseApplyPatch(sizer), "sizer", false)
	imager := sprocketv1.SprocketBlank.DieStamp(func(d *sprocketv1.SprocketDieMut) {
		dietesting.NoError(t, d.Metadata().Namespace().Set(ns))
		dietesting.NoError(t, d.Metadata().Name().Set("chain"))
		dietesting.NoError(t, d.Spec().Image().Set("busybox:1.36"))
	})
	applyObject(t, imager.DieMutable(), die.MustReleaseApplyPatch(imager), "imager", false)

	v1 := getV1(t, ns)
	expect(t, "v1 size", v1.Spec().Size().Get(), 42)
	expect(t, "v1 image", v1.Spec().Image().Get(), "busybox:1.36")
	v2 := getV2(t, ns)
	expect(t, "v2 teeth", v2.Spec().Teeth().Get(), 42)
	expect(t, "v2 repository", v2.Spec().Image().Repository().Get(), "busybox")
}
