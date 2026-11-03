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
	"testing"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	"reconciler.io/dies/v2/die"
	dietesting "reconciler.io/dies/v2/die/testing"
)

// JSON merge patches: what a die sends is its intent against what it read,
// so it doesn't undo what changed on the server since.

func get(t *testing.T, ns string) *dieappsv1.DeploymentDie {
	t.Helper()
	d := dieappsv1.DeploymentBlank.DieMutable()
	if err := dies.Get(context.Background(), keyOf(ns, "web"), d); err != nil {
		t.Fatal(err)
	}
	return d.DieImmutable()
}

func mergePatch(t *testing.T, d *dieappsv1.DeploymentDie) *dieappsv1.DeploymentDie {
	t.Helper()
	body, err := die.ReleaseMergePatch(d)
	if err != nil {
		t.Fatal(err)
	}
	got := d.DieMutable()
	if err := dies.Patch(context.Background(), got, client.RawPatch(types.MergePatchType, body)); err != nil {
		t.Fatalf("merge patch %s: %v", body, err)
	}
	return got.DieImmutable()
}

func TestMergePatchCarriesOnlyIntent(t *testing.T) {
	ns := namespace(t)
	apply(t, deployment(t, ns).DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietesting.NoError(t, d.Metadata().Labels().Set("tier", "web"))
		dietesting.NoError(t, d.Spec().Replicas().Set(2))
		dietesting.NoError(t, d.Spec().Paused().Set(true))
		dietesting.NoError(t, d.Spec().MinReadySeconds().Set(5))
	}), "app")
	read := get(t, ns)

	// another writer changes the object after it was read
	other := identity(t, ns).DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietesting.NoError(t, d.Metadata().Labels().Set("team", "a"))
	})
	apply(t, other, "other")

	got := mergePatch(t, read.DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietesting.NoError(t, d.Spec().Replicas().Set(0))
		dietesting.NoError(t, d.Spec().Paused().Set(false))
		dietesting.NoError(t, d.Spec().MinReadySeconds().Clear())
		dietesting.NoError(t, d.Metadata().Labels().Remove("tier"))
	}))

	if r := got.Spec().Replicas().Ptr(); r == nil || *r != 0 {
		t.Errorf("replicas %v, want 0", r)
	}
	if got.Spec().Paused().Get() {
		t.Errorf("paused, want not")
	}
	if got.Spec().MinReadySeconds().Has() {
		t.Errorf("minReadySeconds %d, want it removed", got.Spec().MinReadySeconds().Get())
	}
	if got.Metadata().Labels().Has("tier") {
		t.Errorf("labels %v, want tier removed", got.GetLabels())
	}
	// what the other writer changed is kept
	if got.GetLabels()["team"] != "a" {
		t.Errorf("labels %v, want the other writer's team label", got.GetLabels())
	}
}

// A merge patch can only replace a list, so it sends the list as the die
// has it, from what it read: the other elements are kept, as they were.
func TestMergePatchSendsListsWhole(t *testing.T) {
	ns := namespace(t)
	apply(t, deployment(t, ns).DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietesting.NoError(t, d.Spec().Template().Spec().Containers().At("sidecar").Image().Set("sidecar:1"))
	}), "app")

	got := mergePatch(t, get(t, ns).DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietesting.NoError(t, d.Spec().Template().Spec().Containers().At("app").Image().Set("app:2"))
	}))
	containers := got.Spec().Template().Spec().Containers()
	if containers.At("app").Image().Get() != "app:2" || containers.At("sidecar").Image().Get() != "sidecar:1" {
		t.Errorf("containers %+v", got.Spec().Template().Spec().DieRelease().Containers)
	}
	// defaults the server set on the container are sent back as they were
	if containers.At("app").TerminationMessagePath().Get() != "/dev/termination-log" {
		t.Errorf("terminationMessagePath %q", containers.At("app").TerminationMessagePath().Get())
	}
}

// The steps of a merge patch, as the README has them: the current state,
// read; the update, a change to it; the patch, from the update; and the
// patched state, decoded into the current die, which has no intent.
func TestMergePatchSteps(t *testing.T) {
	ns := namespace(t)
	ctx := context.Background()
	apply(t, deployment(t, ns).DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietesting.NoError(t, d.Metadata().Labels().Set("canary", "true"))
		dietesting.NoError(t, d.Spec().Replicas().Set(2))
	}), "app")

	current, err := dieappsv1.DeploymentBlank.DieStampE(func(d *dieappsv1.DeploymentDieMut) error {
		return dies.Get(ctx, keyOf(ns, "web"), d)
	})
	if err != nil {
		t.Fatal(err)
	}
	update := current.DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietesting.NoError(t, d.Spec().Replicas().Set(0))
		dietesting.NoError(t, d.Metadata().Labels().Remove("canary"))
	})
	body, err := die.ReleaseMergePatch(update)
	if err != nil {
		t.Fatal(err)
	}
	patched, err := current.DieStampE(func(d *dieappsv1.DeploymentDieMut) error {
		return dies.Patch(ctx, d, client.RawPatch(types.MergePatchType, body))
	})
	if err != nil {
		t.Fatal(err)
	}

	if patched.DieIntent() != "" {
		t.Errorf("patched has intent %q", patched.DieIntent())
	}
	if patched.Spec().Replicas().Get() != 0 || patched.Metadata().Labels().Has("canary") {
		t.Errorf("patched replicas %d, labels %v", patched.Spec().Replicas().Get(), patched.GetLabels())
	}
	if patched.GetResourceVersion() == current.GetResourceVersion() {
		t.Errorf("patched has the current resource version")
	}
	// the steps before are as they were
	if current.Spec().Replicas().Get() != 2 || update.Spec().Replicas().Get() != 0 {
		t.Errorf("current replicas %d, update replicas %d", current.Spec().Replicas().Get(), update.Spec().Replicas().Get())
	}
}
