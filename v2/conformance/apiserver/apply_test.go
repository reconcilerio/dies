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
	diecorev1 "reconciler.io/dies/v2/apis/core/v1"
	"reconciler.io/dies/v2/die/base"
)

// Server-side apply: what a die sends is what its manager owns.

// deployment is a minimal valid Deployment, with one container, app.
func deployment(namespace string) *dieappsv1.DeploymentDie {
	return dieappsv1.DeploymentBlank.DieWith(func(d *dieappsv1.DeploymentDieMut) {
		d.Metadata().Namespace().Set(namespace)
		d.Metadata().Name().Set("web")
		d.Spec().Selector().MatchLabels().Set("app", "web")
		d.Spec().Template().Metadata().Labels().Set("app", "web")
		d.Spec().Template().Spec().Containers().At("app").Image().Set("app:1")
	})
}

// identity is a Deployment with only its identity, to apply nothing else.
func identity(namespace string) *dieappsv1.DeploymentDie {
	return dieappsv1.DeploymentBlank.DieWith(func(d *dieappsv1.DeploymentDieMut) {
		d.Metadata().Namespace().Set(namespace)
		d.Metadata().Name().Set("web")
	})
}

// apply sends a die's apply patch as a field manager, and returns what the
// server has after.
func apply(t *testing.T, d *dieappsv1.DeploymentDie, manager string) *dieappsv1.DeploymentDie {
	t.Helper()
	body, err := d.DieReleaseApplyPatchE()
	if err != nil {
		t.Fatal(err)
	}
	got := d.DieMutable()
	if err := dies.Patch(context.Background(), got, client.RawPatch(types.ApplyPatchType, body), client.FieldOwner(manager), client.ForceOwnership); err != nil {
		t.Fatalf("apply %s: %v", body, err)
	}
	return got.DieImmutable()
}

func TestApplyOwnsOnlyWhatWasSet(t *testing.T) {
	ns := namespace(t)
	got := apply(t, deployment(ns), "conformance")

	assertOwned(t, got, "conformance", "",
		".spec.selector",
		".spec.template.metadata.labels.app",
		`.spec.template.spec.containers[name="app"].image`,
		`.spec.template.spec.containers[name="app"].name`,
	)
	// the server defaulted what wasn't set, and doesn't say the manager owns it
	if got.Spec().Replicas().Get() != 1 {
		t.Errorf("replicas %d, want the default", got.Spec().Replicas().Get())
	}
}

// A zero value that was set is sent, so it's owned and kept, rather than
// defaulted.
func TestApplySendsExplicitZeroValues(t *testing.T) {
	ns := namespace(t)
	got := apply(t, deployment(ns).DieWith(func(d *dieappsv1.DeploymentDieMut) {
		d.Spec().Replicas().Set(0)
		d.Spec().Paused().Set(false)
	}), "conformance")

	if !got.Spec().Replicas().Has() || got.Spec().Replicas().Get() != 0 {
		t.Errorf("replicas %v, want 0", got.Spec().Replicas().Ptr())
	}
	assertOwns(t, got, "conformance", "", ".spec.replicas", ".spec.paused")
}

// A list map merges by key: each manager's elements are its own, and an
// element its manager no longer applies is removed.
func TestApplyMergesListMapsByKey(t *testing.T) {
	ns := namespace(t)
	apply(t, deployment(ns), "app")
	sidecar := deployment(ns).DieWith(func(d *dieappsv1.DeploymentDieMut) {
		d.Spec().Template().Spec().Containers().Remove("app")
		d.Spec().Template().Spec().Containers().At("sidecar").Image().Set("sidecar:1")
	})
	got := apply(t, sidecar, "sidecar")
	if keys := got.Spec().Template().Spec().Containers().Keys(); len(keys) != 2 {
		t.Fatalf("containers %v, want app and sidecar", keys)
	}

	// app applies again, without sidecar, which it doesn't own: it stays
	got = apply(t, deployment(ns), "app")
	if !got.Spec().Template().Spec().Containers().Has("sidecar") {
		t.Errorf("app's apply removed sidecar, which it doesn't own")
	}

	// sidecar applies without its container: it's removed
	got = apply(t, identity(ns), "sidecar")
	if got.Spec().Template().Spec().Containers().Has("sidecar") {
		t.Errorf("sidecar's container stayed after sidecar stopped applying it")
	}
	if !got.Spec().Template().Spec().Containers().Has("app") {
		t.Errorf("app's container was removed")
	}
}

// Ports are keyed by containerPort and protocol, which defaults to TCP. A
// port set without a protocol is the TCP one, on the die and on the server.
func TestApplyDefaultedListMapKeys(t *testing.T) {
	ns := namespace(t)
	withPort := deployment(ns).DieWith(func(d *dieappsv1.DeploymentDieMut) {
		d.Spec().Template().Spec().Containers().At("app").Ports().At(diecorev1.ContainerPortKey{ContainerPort: 8080}).Name().Set("http")
	})
	first := apply(t, withPort, "app")
	again := apply(t, withPort, "app")
	if again.GetResourceVersion() != first.GetResourceVersion() {
		t.Errorf("applying the same port again changed the object")
	}

	udp := deployment(ns).DieWith(func(d *dieappsv1.DeploymentDieMut) {
		d.Spec().Template().Spec().Containers().At("app").Ports().At(diecorev1.ContainerPortKey{ContainerPort: 8080, Protocol: diecorev1.ProtocolUDP}).Name().Set("dns")
	})
	got := apply(t, udp, "dns")
	ports := got.Spec().Template().Spec().Containers().At("app").Ports()
	if ports.Len() != 2 {
		t.Fatalf("ports %v, want TCP and UDP", ports.Keys())
	}
	if ports.At(diecorev1.ContainerPortKey{ContainerPort: 8080}).Name().Get() != "http" {
		t.Errorf("the TCP port is %+v", ports.At(diecorev1.ContainerPortKey{ContainerPort: 8080}).DieRelease())
	}
}

// The server replaces an atomic struct, such as a ConfigMapKeySelector,
// whole, so a die applies it whole: another manager changing its key, from
// what the server has, keeps its name.
func TestApplySendsAtomicStructsWhole(t *testing.T) {
	ns := namespace(t)
	withEnv := deployment(ns).DieWith(func(d *dieappsv1.DeploymentDieMut) {
		ref := d.Spec().Template().Spec().Containers().At("app").Env().At("COLOR").ValueFrom().ConfigMapKeyRef()
		ref.Name().Set("settings")
		ref.Key().Set("color")
		ref.Optional().Set(true)
	})
	actual := apply(t, withEnv, "app")

	// from the server's state, change only the key, as another manager: an
	// apply is all of a manager's intent, so app applying only this would
	// drop the rest of what it owns
	changed := dieappsv1.DeploymentBlank.DieWith(func(d *dieappsv1.DeploymentDieMut) {
		d.DieFeed(actual.DieRelease())
		d.Spec().Template().Spec().Containers().At("app").Env().At("COLOR").ValueFrom().ConfigMapKeyRef().Key().Set("colour")
	})
	got := apply(t, changed, "localizer")
	assertOwns(t, got, "localizer", "", `.spec.template.spec.containers[name="app"].env[name="COLOR"].valueFrom.configMapKeyRef`)
	ref := got.Spec().Template().Spec().Containers().At("app").Env().At("COLOR").ValueFrom().ConfigMapKeyRef()
	if ref.Name().Get() != "settings" || ref.Key().Get() != "colour" || !ref.Optional().Get() {
		t.Errorf("configMapKeyRef %+v, want settings/colour, optional", ref.DieRelease())
	}
}

// Clear on a field sends null: the manager owns it, and it's absent.
func TestApplyClearRemovesAField(t *testing.T) {
	ns := namespace(t)
	apply(t, deployment(ns).DieWith(func(d *dieappsv1.DeploymentDieMut) {
		d.Spec().MinReadySeconds().Set(5)
	}), "rollout")

	got := apply(t, identity(ns).DieWith(func(d *dieappsv1.DeploymentDieMut) {
		d.Spec().MinReadySeconds().Clear()
	}), "app")
	if got.Spec().MinReadySeconds().Has() {
		t.Errorf("minReadySeconds %d, want it removed", got.Spec().MinReadySeconds().Get())
	}
	assertOwns(t, got, "app", "", ".spec.minReadySeconds")
}

// An apply can't remove another manager's map entry, so removing one is left
// out of the apply, rather than sent as null, which the server would take
// as setting it to "". The manager that owns an entry removes it by no
// longer applying it.
func TestApplyRemovingAMapEntry(t *testing.T) {
	ns := namespace(t)
	apply(t, deployment(ns).DieWith(func(d *dieappsv1.DeploymentDieMut) {
		d.Metadata().Labels().Set("tier", "web")
	}), "labeler")

	got := apply(t, identity(ns).DieWith(func(d *dieappsv1.DeploymentDieMut) {
		d.Metadata().Labels().Remove("tier")
	}), "app")
	if tier := got.GetLabels()["tier"]; tier != "web" {
		t.Errorf("tier %q, want labeler's label untouched", tier)
	}

	got = apply(t, deployment(ns), "labeler")
	if got.Metadata().Labels().Has("tier") {
		t.Errorf("labels %v, want tier removed once labeler stopped applying it", got.GetLabels())
	}
}

// Empty sends []: the manager owns a list that's present and empty.
func TestApplyEmptyOwnsAnEmptyList(t *testing.T) {
	ns := namespace(t)
	got := apply(t, deployment(ns).DieWith(func(d *dieappsv1.DeploymentDieMut) {
		d.Spec().Template().Spec().Tolerations().Empty()
	}), "app")
	assertOwns(t, got, "app", "", ".spec.template.spec.tolerations")
}

// Status is applied to its subresource, from a blank, as reconcilers do:
// the spec is untouched, and the manager owns only the status it set.
func TestApplyStatusSubresource(t *testing.T) {
	ns := namespace(t)
	apply(t, deployment(ns), "app")

	status := identity(ns).DieWith(func(d *dieappsv1.DeploymentDieMut) {
		d.Status().ObservedGeneration().Set(1)
		d.Status().Conditions().At(dieappsv1.DeploymentAvailable).DieWith(func(c *dieappsv1.DeploymentConditionDieMut) {
			c.Status().Set(diecorev1.ConditionTrue)
			c.Reason().Set("Reconciled")
		})
	}).DieSubresource("status")
	body, err := status.DieReleaseApplyPatchE()
	if err != nil {
		t.Fatal(err)
	}
	got := status.DieMutable()
	if err := dies.Status().Patch(context.Background(), got, client.RawPatch(types.ApplyPatchType, body), client.FieldOwner("controller"), client.ForceOwnership); err != nil {
		t.Fatal(err)
	}
	if got.Status().Conditions().At(dieappsv1.DeploymentAvailable).Reason().Get() != "Reconciled" {
		t.Errorf("status %+v", got.Status().DieRelease())
	}
	if got.Spec().Template().Spec().Containers().At("app").Image().Get() != "app:1" {
		t.Errorf("spec %+v", got.Spec().DieRelease())
	}
	assertOwned(t, got, "controller", "status",
		".status.observedGeneration",
		`.status.conditions[type="Available"].reason`,
		`.status.conditions[type="Available"].status`,
		`.status.conditions[type="Available"].type`,
	)
}

// An object die is an apply configuration, as its generated
// DieReleaseApplyConfiguration and the generic adapter both make it, for
// controller-runtime's typed Apply.
func TestApplyConfigurations(t *testing.T) {
	ns := namespace(t)
	ctx := context.Background()

	if err := dies.Apply(ctx, deployment(ns).DieReleaseApplyConfiguration(), client.FieldOwner("typed")); err != nil {
		t.Fatal(err)
	}

	ac := base.ApplyConfigurationOf(deployment(ns).DieWith(func(d *dieappsv1.DeploymentDieMut) {
		d.Spec().Replicas().Set(3)
	}))
	if err := dies.Apply(ctx, ac, client.FieldOwner("adapted"), client.ForceOwnership); err != nil {
		t.Fatal(err)
	}
	got := dieappsv1.DeploymentBlank.DieWith(func(d *dieappsv1.DeploymentDieMut) { d.DieFeedJSON(ac.Result()) })
	if got.Spec().Replicas().Get() != 3 {
		t.Errorf("replicas %v", got.Spec().Replicas().Ptr())
	}
	assertOwned(t, got, "adapted", "",
		".spec.replicas",
		".spec.selector",
		".spec.template.metadata.labels.app",
		`.spec.template.spec.containers[name="app"].image`,
		`.spec.template.spec.containers[name="app"].name`,
	)
}
