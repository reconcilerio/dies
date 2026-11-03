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
	"testing"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	diemetav1 "reconciler.io/dies/v2/apis/meta/v1"
	"reconciler.io/dies/v2/conformance/internal/harness"
	"reconciler.io/dies/v2/die"
	"reconciler.io/dies/v2/die/dietest"
	widgetv1 "reconciler.io/dies/v2/examples/widget/v1"
)

// Patches from dies, applied by a real cluster. The objects are ones no
// controller changes in ways the tests look at: Deployments scaled to zero
// and paused, whose specs only their appliers write, and Widgets, which have
// no controller.

func keyOf(namespace, name string) client.ObjectKey {
	return client.ObjectKey{Namespace: namespace, Name: name}
}

// identity is a Deployment with only its identity, to apply nothing else.
func identity(t testing.TB, ns string) *dieappsv1.DeploymentDie {
	return dieappsv1.DeploymentBlank.DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietest.NoError(t, d.Metadata().Namespace().Set(ns))
		dietest.NoError(t, d.Metadata().Name().Set("web"))
	})
}

// deployment is a valid Deployment that runs nothing.
func deployment(t testing.TB, ns string) *dieappsv1.DeploymentDie {
	return identity(t, ns).DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietest.NoError(t, d.Spec().Replicas().Set(0))
		dietest.NoError(t, d.Spec().Paused().Set(true))
		dietest.NoError(t, d.Spec().Selector().MatchLabels().Set("app", "web"))
		dietest.NoError(t, d.Spec().Template().Metadata().Labels().Set("app", "web"))
		dietest.NoError(t, d.Spec().Template().Spec().Containers().At("app").Image().Set("registry.example/app:1"))
	})
}

// apply sends a die's apply as a manager, and returns what the server has
// after: the die is tared first, so the result isn't overlaid with the
// apply's intent.
func apply[T any, D die.DieOf[T, D, M], M interface {
	die.DieMutOf[T, D, M]
	client.Object
}](t *testing.T, d D, manager string, status bool) D {
	t.Helper()
	var scopes []die.Scope
	if status {
		scopes = append(scopes, die.Subresource("status"))
	}
	body, err := die.ReleaseApplyPatchE(d, scopes...)
	if err != nil {
		t.Fatal(err)
	}
	got, err := d.DieTare().DieStampE(func(d M) error {
		patch := client.RawPatch(types.ApplyPatchType, body)
		if status {
			return dies.Status().Patch(context.Background(), d, patch, client.FieldOwner(manager), client.ForceOwnership)
		}
		return dies.Patch(context.Background(), d, patch, client.FieldOwner(manager), client.ForceOwnership)
	})
	if err != nil {
		t.Fatalf("apply %s: %v", body, err)
	}
	return got
}

func TestApplyOwnsWhatWasSet(t *testing.T) {
	ns := harness.Namespace(t, dies)
	got := apply(t, deployment(t, ns), "app", false)
	harness.AssertOwned(t, got, "app", "",
		".spec.paused",
		".spec.replicas",
		".spec.selector",
		".spec.template.metadata.labels.app",
		`.spec.template.spec.containers[name="app"].image`,
		`.spec.template.spec.containers[name="app"].name`,
	)
	if got.Spec().Replicas().Get() != 0 {
		t.Errorf("replicas %d, want the explicit 0", got.Spec().Replicas().Get())
	}
}

func TestApplyMergesListMapsByKey(t *testing.T) {
	ns := harness.Namespace(t, dies)
	apply(t, deployment(t, ns), "app", false)
	got := apply(t, identity(t, ns).DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietest.NoError(t, d.Spec().Template().Spec().Containers().At("sidecar").Image().Set("registry.example/sidecar:1"))
	}), "sidecar", false)
	if keys := got.Spec().Template().Spec().Containers().Keys(); len(keys) != 2 {
		t.Fatalf("containers %v, want app and sidecar", keys)
	}
	// sidecar stops applying its container, which it owned: it's removed
	got = apply(t, identity(t, ns), "sidecar", false)
	if containers := got.Spec().Template().Spec().Containers(); containers.Has("sidecar") || !containers.Has("app") {
		t.Errorf("containers %v, want only app", containers.Keys())
	}
}

// Clear sends null for a built-in type's field, which removes it, and the
// clearer owns it; removing a map entry is left out of an apply, so it's
// another manager's to remove.
func TestApplyClearsAndRemoves(t *testing.T) {
	ns := harness.Namespace(t, dies)
	apply(t, deployment(t, ns).DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietest.NoError(t, d.Spec().MinReadySeconds().Set(5))
		dietest.NoError(t, d.Metadata().Labels().Set("tier", "web"))
	}), "app", false)

	got := apply(t, identity(t, ns).DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietest.NoError(t, d.Spec().MinReadySeconds().Clear())
		dietest.NoError(t, d.Metadata().Labels().Remove("tier"))
	}), "clearer", false)
	if got.Spec().MinReadySeconds().Has() {
		t.Errorf("minReadySeconds %d, want it removed", got.Spec().MinReadySeconds().Get())
	}
	harness.AssertOwns(t, got, "clearer", "", ".spec.minReadySeconds")
	if got.GetLabels()["tier"] != "web" {
		t.Errorf("labels %v, want app's tier label kept", got.GetLabels())
	}
}

// An atomic struct is applied whole, so another manager changing part of it,
// from what the server has, keeps the rest.
func TestApplySendsAtomicStructsWhole(t *testing.T) {
	ns := harness.Namespace(t, dies)
	actual := apply(t, deployment(t, ns).DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		ref := d.Spec().Template().Spec().Containers().At("app").Env().At("COLOR").ValueFrom().ConfigMapKeyRef()
		dietest.NoError(t, ref.Name().Set("settings"))
		dietest.NoError(t, ref.Key().Set("color"))
	}), "app", false)

	got := apply(t, actual.DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietest.NoError(t, d.Spec().Template().Spec().Containers().At("app").Env().At("COLOR").ValueFrom().ConfigMapKeyRef().Key().Set("colour"))
	}), "localizer", false)
	ref := got.Spec().Template().Spec().Containers().At("app").Env().At("COLOR").ValueFrom().ConfigMapKeyRef()
	if ref.Name().Get() != "settings" || ref.Key().Get() != "colour" {
		t.Errorf("configMapKeyRef %+v", *ref.DieRelease())
	}
}

// A custom resource's status, applied to its subresource by two managers,
// each owning its own condition.
func TestApplyCustomResourceStatus(t *testing.T) {
	ns := harness.Namespace(t, dies)
	widget := widgetv1.WidgetBlank.DieStamp(func(d *widgetv1.WidgetDieMut) {
		dietest.NoError(t, d.Metadata().Namespace().Set(ns))
		dietest.NoError(t, d.Metadata().Name().Set("web"))
	})
	apply(t, widget.DieStamp(func(d *widgetv1.WidgetDieMut) {
		dietest.NoError(t, d.Spec().Ports().At("http").Port().Set(80))
	}), "web", false)

	var got *widgetv1.WidgetDie
	for _, manager := range []string{"Ready", "Healthy"} {
		got = apply(t, widget.DieStamp(func(d *widgetv1.WidgetDieMut) {
			d.Status().Conditions().At(manager).DieStamp(func(c *diemetav1.ConditionDieMut) {
				dietest.NoError(t, c.Status().Set(diemetav1.ConditionTrue))
				dietest.NoError(t, c.Reason().Set("Checked"))
			})
		}), manager, true)
	}
	if conditions := got.Status().Conditions(); conditions.Len() != 2 {
		t.Errorf("conditions %v, want Ready and Healthy", conditions.Keys())
	}
	if got.Spec().Ports().At("http").Port().Get() != 80 {
		t.Errorf("spec %+v", *got.Spec().DieRelease())
	}
}

// The steps of a merge patch: the current state, read; the update, a change
// to it; the patch, from the update; and the patched state. What another
// writer changed between the read and the patch is kept.
func TestMergePatchCarriesOnlyIntent(t *testing.T) {
	ns := harness.Namespace(t, dies)
	ctx := context.Background()
	apply(t, deployment(t, ns).DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietest.NoError(t, d.Metadata().Labels().Set("canary", "true"))
	}), "app", false)

	current, err := dieappsv1.DeploymentBlank.DieStampE(func(d *dieappsv1.DeploymentDieMut) error {
		return dies.Get(ctx, keyOf(ns, "web"), d)
	})
	if err != nil {
		t.Fatal(err)
	}
	apply(t, identity(t, ns).DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietest.NoError(t, d.Metadata().Labels().Set("team", "a"))
	}), "other", false)

	update := current.DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietest.NoError(t, d.Metadata().Labels().Remove("canary"))
		dietest.NoError(t, d.Spec().Template().Spec().Containers().At("app").Image().Set("registry.example/app:2"))
	})
	body, err := die.ReleaseMergePatchE(update)
	if err != nil {
		t.Fatal(err)
	}
	patched, err := current.DieStampE(func(d *dieappsv1.DeploymentDieMut) error {
		return dies.Patch(ctx, d, client.RawPatch(types.MergePatchType, body))
	})
	if err != nil {
		t.Fatalf("merge patch %s: %v", body, err)
	}
	if patched.Metadata().Labels().Has("canary") || patched.GetLabels()["team"] != "a" {
		t.Errorf("labels %v, want canary removed and team kept", patched.GetLabels())
	}
	if image := patched.Spec().Template().Spec().Containers().At("app").Image().Get(); image != "registry.example/app:2" {
		t.Errorf("image %q", image)
	}
}
