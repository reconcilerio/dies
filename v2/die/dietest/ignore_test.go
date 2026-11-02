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

package dietest_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	diebatchv1 "reconciler.io/dies/v2/apis/batch/v1"
	"reconciler.io/dies/v2/die/dietest"
	"reconciler.io/dies/v2/die/unstructured"
)

// diffs reports the diff of x and y both ways: by content, and released.
func diffs(x, y any, opts ...cmp.Option) map[string]string {
	return map[string]string{
		"Diff":        dietest.Diff(x, y, opts...),
		"DiffRelease": dietest.DiffRelease(x, y, opts...),
	}
}

func assertNoDiffs(t *testing.T, what string, x, y any, opts ...cmp.Option) {
	t.Helper()
	for flavor, diff := range diffs(x, y, opts...) {
		if diff != "" {
			t.Errorf("%s: %s:\n%s", what, flavor, diff)
		}
	}
}

func assertDiffs(t *testing.T, what string, x, y any, contains string, opts ...cmp.Option) {
	t.Helper()
	for flavor, diff := range diffs(x, y, opts...) {
		if !strings.Contains(diff, contains) {
			t.Errorf("%s: %s should show %s:\n%s", what, flavor, contains, diff)
		}
	}
}

func TestIgnoreFields(t *testing.T) {
	blank := dieappsv1.DeploymentBlank
	with := func(fn func(d *deploy)) *dieappsv1.DeploymentDie { return named("web").DieWith(fn) }

	x := with(func(d *deploy) { d.Metadata().ResourceVersion().Set("1"); d.Spec().Replicas().Set(1) })
	y := with(func(d *deploy) { d.Metadata().ResourceVersion().Set("2"); d.Spec().Replicas().Set(1) })
	assertDiffs(t, "resource versions", x, y, "2")
	assertNoDiffs(t, "resource versions, ignored", x, y, dietest.IgnoreFields(blank.Metadata().ResourceVersion()))
	assertNoDiffs(t, "resource versions, metadata ignored", x, y, dietest.IgnoreFields(blank.Metadata()))

	// only what's ignored
	z := with(func(d *deploy) { d.Metadata().ResourceVersion().Set("2"); d.Spec().Replicas().Set(2) })
	assertDiffs(t, "replicas", x, z, "2", dietest.IgnoreFields(blank.Metadata().ResourceVersion()))

	// an inline field
	api := with(func(d *deploy) { d.APIVersion().Set("apps/v2") })
	assertDiffs(t, "apiVersion", named("web"), api, "apps/v2")
	assertNoDiffs(t, "apiVersion, ignored", named("web"), api, dietest.IgnoreFields(blank.APIVersion()))

	// a subresource
	ready := with(func(d *deploy) { d.Status().ReadyReplicas().Set(1) })
	assertNoDiffs(t, "status, ignored", named("web"), ready, dietest.IgnoreFields(blank.Status()))

	// a list-map element by its key, present on either side or both
	app := func(image string) func(d *deploy) {
		return func(d *deploy) { d.Spec().Template().Spec().Containers().At("app").Image().Set(image) }
	}
	sidecar := func(image string) func(d *deploy) {
		return func(d *deploy) { d.Spec().Template().Spec().Containers().At("sidecar").Image().Set(image) }
	}
	both := func(fns ...func(d *deploy)) *dieappsv1.DeploymentDie {
		return with(func(d *deploy) {
			for _, fn := range fns {
				fn(d)
			}
		})
	}
	ignoreSidecar := dietest.IgnoreFields(blank.Spec().Template().Spec().Containers().At("sidecar"))
	assertNoDiffs(t, "a sidecar on one side", both(app("app:1"), sidecar("proxy:1")), both(app("app:1")), ignoreSidecar)
	assertNoDiffs(t, "different sidecars", both(app("app:1"), sidecar("proxy:1")), both(app("app:1"), sidecar("proxy:2")), ignoreSidecar)
	assertDiffs(t, "different apps", both(app("app:1"), sidecar("proxy:1")), both(app("app:2"), sidecar("proxy:1")), "app:2", ignoreSidecar)

	// a field the type doesn't have is compared by content only, and
	// ignored by an enclosing location
	future := fedJSON(`{"spec": {"futureField": 1}}`)
	future2 := fedJSON(`{"spec": {"futureField": 2}}`)
	if diff := dietest.Diff(future, future2, dietest.IgnoreFields(blank.Spec())); diff != "" {
		t.Errorf("unexpected diff:\n%s", diff)
	}
}

// A location within a blank applies to dies of the blank's type only.
func TestIgnoreFieldsByType(t *testing.T) {
	cron := func(name string) *diebatchv1.CronJobDie {
		return diebatchv1.CronJobBlank.DieWith(func(d *diebatchv1.CronJobDieMut) { d.Metadata().Name().Set(name) })
	}
	type pair struct {
		Deployment *dieappsv1.DeploymentDie
		CronJob    *diebatchv1.CronJobDie
	}
	x, y := pair{named("web"), cron("web")}, pair{named("api"), cron("api")}
	ignore := dietest.IgnoreFields(dieappsv1.DeploymentBlank.Metadata().Name())
	assertNoDiffs(t, "deployments", x.Deployment, y.Deployment, ignore)
	assertDiffs(t, "cron jobs", x.CronJob, y.CronJob, `"api"`, ignore)
	assertDiffs(t, "pairs", x, y, `"api"`, ignore)
	// and within dies wherever they are, such as a list die's items
	list := func(names ...string) *dieappsv1.DeploymentListDie {
		var items []*deploy
		for _, n := range names {
			items = append(items, named("web").DieWith(func(d *deploy) { d.Metadata().ResourceVersion().Set(n) }).DieMutable())
		}
		return dieappsv1.NewDeploymentList(items...)
	}
	assertNoDiffs(t, "list items", list("1", "2"), list("3", "4"), dietest.IgnoreFields(dieappsv1.DeploymentBlank.Metadata().ResourceVersion()))
}

// Unstructured dies are ignored by the locations of unstructured blanks.
func TestIgnoreFieldsOfUnstructured(t *testing.T) {
	obj := func(size int64) *unstructured.Object {
		return unstructured.BlankFor("example.dev/v1", "Widget").DieWith(func(d *unstructured.ObjectMut) {
			d.Node().Field("spec").Field("size").Set(size)
			d.Node().Field("spec").Field("color").Set("blue")
		})
	}
	assertDiffs(t, "sizes", obj(1), obj(2), "2")
	assertNoDiffs(t, "sizes, ignored", obj(1), obj(2), dietest.IgnoreFields(unstructured.Blank.Node().Field("spec").Field("size")))
}
