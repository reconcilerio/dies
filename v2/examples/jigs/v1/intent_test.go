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
	"testing"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	diebatchv1 "reconciler.io/dies/v2/apis/batch/v1"
	diecorev1 "reconciler.io/dies/v2/apis/core/v1"
	"reconciler.io/dies/v2/die"
	"reconciler.io/dies/v2/die/dietest"
	jigsv1 "reconciler.io/dies/v2/examples/jigs/v1"
)

func assertIntent(t *testing.T, what, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("%s intent\n got: %s\nwant: %s", what, got, want)
	}
}

// A jig's intent is its source's, through its bindings: what is written
// through the jig, or to the source, at the locations the jig binds, at the
// jig's own locations.
func TestJigIntentIsTheSources(t *testing.T) {
	deploy := dieappsv1.DeploymentBlank.DieMutable()
	ps, _ := jigsv1.PodSpecableJig.OfMut(deploy)
	dietest.NoError(t, ps.Template().Spec().Containers().At("app").Image().Set("app:1"))
	// not bound by the jig, so not its intent
	dietest.NoError(t, deploy.Spec().Replicas().Set(3))

	assertIntent(t, "jig", ps.DieIntent(), `set .template.spec.containers[name=app].image = "app:1"`)
	assertIntent(t, "jig's template", ps.Template().DieIntent(), `set .template.spec.containers[name=app].image = "app:1"`)

	// as for a jig of an immutable die, whose intent was written to it
	ro, _ := jigsv1.PodSpecableJig.Of(deploy.DieImmutable())
	assertIntent(t, "immutable jig", ro.DieIntent(), `set .template.spec.containers[name=app].image = "app:1"`)
}

// Two jig locations bound to the same source location both have its intent.
func TestOverlappingBindingsShareIntent(t *testing.T) {
	pod := diecorev1.PodBlank.DieStamp(func(d *diecorev1.PodDieMut) {
		dietest.NoError(t, d.Metadata().Labels().Set("app", "web"))
	})
	ps, _ := jigsv1.PodSpecableJig.Of(pod)
	assertIntent(t, "jig", ps.DieIntent(), "set .metadata.labels[\"app\"] = \"web\"\nset .template.metadata.labels[\"app\"] = \"web\"")
}

// Where a transform presents the source, its operations can't be mapped:
// any intent below the binding sets the binding's whole value.
func TestTransformedIntentSetsTheBinding(t *testing.T) {
	container := diecorev1.ContainerBlank.DieStamp(func(d *diecorev1.ContainerDieMut) {
		dietest.NoError(t, d.Name().Set("app"))
		dietest.NoError(t, d.Image().Set("registry.example/app:1"))
	})
	ci, _ := jigsv1.ContainerImageJig.Of(container)
	assertIntent(t, "jig", ci.DieIntent(), "set .name = \"app\"\nset .image = {\"repository\":\"registry.example/app\",\"tag\":\"1\"}")

	// the binding's value, not the source's: none of it was written
	unwritten, _ := jigsv1.ContainerImageJig.Of(diecorev1.ContainerBlank)
	assertIntent(t, "unwritten jig", unwritten.DieIntent(), "")
}

// As for an aggregate of several source locations.
func TestAggregatedIntentSetsTheBinding(t *testing.T) {
	deploy := dieappsv1.DeploymentBlank.DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietest.NoError(t, d.Spec().Template().Spec().Containers().At("app").Image().Set("app:1"))
	})
	w, _ := jigsv1.WorkloadJig.Of(deploy)
	assertIntent(t, "jig", w.DieIntent(), "set .containers = [{\"name\":\"app\",\"image\":\"app:1\",\"resources\":{}}]")
}

// Jigs of different types compare by what the jig binds, value and intent.
func TestJigsOfDifferentTypesCompare(t *testing.T) {
	deploy := func(image string) *jigsv1.PodSpecableDie {
		ps, _ := jigsv1.PodSpecableJig.Of(dieappsv1.DeploymentBlank.DieStamp(func(d *dieappsv1.DeploymentDieMut) {
			dietest.NoError(t, d.Spec().Replicas().Set(3))
			dietest.NoError(t, d.Spec().Template().Spec().Containers().At("app").Image().Set(image))
		}))
		return ps
	}
	cron := func(image string) *jigsv1.PodSpecableDie {
		ps, _ := jigsv1.PodSpecableJig.Of(diebatchv1.CronJobBlank.DieStamp(func(d *diebatchv1.CronJobDieMut) {
			dietest.NoError(t, d.Spec().Schedule().Set("@daily"))
			dietest.NoError(t, d.Spec().JobTemplate().Spec().Template().Spec().Containers().At("app").Image().Set(image))
		}))
		return ps
	}
	if x, y := deploy("app:1"), cron("app:1"); !x.Equal(y) || dietest.Diff(x, y) != "" {
		t.Errorf("equal %v, diff:\n%s", x.Equal(y), dietest.Diff(x, y))
	}
	if x, y := deploy("app:1"), cron("app:2"); x.Equal(y) || dietest.Diff(x, y) == "" {
		t.Errorf("expected a difference")
	}

	// the same value, with different intent, differs
	fed := diebatchv1.CronJobBlank.DieStamp(func(d *diebatchv1.CronJobDieMut) {
		dietest.NoError(t, die.FeedJSON(d, []byte(`{"spec": {"jobTemplate": {"spec": {"template": {"spec": {"containers": [{"name": "app", "image": "app:1"}]}}}}}}`)))
	})
	y, _ := jigsv1.PodSpecableJig.Of(fed)
	x := deploy("app:1")
	if x.DieRelease().Template.Spec.Containers[0].Image != y.DieRelease().Template.Spec.Containers[0].Image {
		t.Fatalf("expected the same image")
	}
	if x.Equal(y) || dietest.Diff(x, y) == "" {
		t.Errorf("equal %v, diff %q", x.Equal(y), dietest.Diff(x, y))
	}
	t.Logf("diff:\n%s", dietest.Diff(x, y))
}

// Jigs of different types compare by what the jig binds, ignoring what the
// comparison doesn't care about by the jig's own locations.
func TestJigsCompareIgnoringFields(t *testing.T) {
	deploy, _ := jigsv1.PodSpecableJig.Of(dieappsv1.DeploymentBlank.DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietest.NoError(t, d.Metadata().Name().Set("web"))
		dietest.NoError(t, d.Spec().Template().Spec().Containers().At("app").Image().Set("app:1"))
	}))
	cron, _ := jigsv1.PodSpecableJig.Of(diebatchv1.CronJobBlank.DieStamp(func(d *diebatchv1.CronJobDieMut) {
		dietest.NoError(t, d.Metadata().Name().Set("nightly"))
		dietest.NoError(t, d.Spec().JobTemplate().Spec().Template().Spec().Containers().At("app").Image().Set("app:1"))
	}))
	ignore := dietest.IgnoreFields(jigsv1.PodSpecableBlank.Metadata())
	if diff := dietest.Diff(deploy, cron); diff == "" {
		t.Errorf("expected the names to differ")
	}
	if diff := dietest.Diff(deploy, cron, ignore); diff != "" {
		t.Errorf("unexpected diff:\n%s", diff)
	}
	if diff := dietest.DiffRelease(deploy, cron, ignore); diff != "" {
		t.Errorf("unexpected release diff:\n%s", diff)
	}
}
