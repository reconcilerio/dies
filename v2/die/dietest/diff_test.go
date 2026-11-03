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
	"github.com/google/go-cmp/cmp/cmpopts"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	"reconciler.io/dies/v2/die"
	"reconciler.io/dies/v2/die/base"
	"reconciler.io/dies/v2/die/dietest"
	"reconciler.io/dies/v2/die/unstructured"
)

type deploy = dieappsv1.DeploymentDieMut

func named(t testing.TB, name string) *dieappsv1.DeploymentDie {
	return dieappsv1.DeploymentBlank.DieStamp(func(d *deploy) { dietest.NoError(t, d.Metadata().Name().Set(name)) })
}

func at(t testing.TB, sec, nsec int64) *dieappsv1.DeploymentDie {
	return named(t, "web").DieStamp(func(d *deploy) {
		dietest.NoError(t, d.Metadata().CreationTimestamp().Set(metav1.Unix(sec, nsec)))
	})
}

func fedJSON(t testing.TB, data string) *dieappsv1.DeploymentDie {
	return dieappsv1.DeploymentBlank.DieStamp(func(d *deploy) { dietest.NoError(t, die.FeedJSON(d, []byte(data))) })
}

// Every pair of dies has a diff exactly when they aren't Equal.
func TestDiffIsEmptyExactlyWhenEqual(t *testing.T) {
	type pair[D any] struct {
		name  string
		x, y  D
		equal bool
	}
	immutable := []pair[*dieappsv1.DeploymentDie]{
		{"blanks", dieappsv1.DeploymentBlank, dieappsv1.DeploymentBlank, true},
		{"same writes", named(t, "web"), named(t, "web"), true},
		{"different values", named(t, "web"), named(t, "api"), false},
		{"fed and tared", named(t, "web").DieTare(), fedJSON(t, `{"apiVersion": "apps/v1", "kind": "Deployment", "metadata": {"name": "web"}}`), true},
		{"same value, different intent", named(t, "web"), fedJSON(t, `{"apiVersion": "apps/v1", "kind": "Deployment", "metadata": {"name": "web"}}`), false},
		{"intent of the same value, written again", named(t, "web"), named(t, "web").DieStamp(func(d *deploy) { dietest.NoError(t, d.Metadata().Name().Set("web")) }), true},
		{"a field the type doesn't have", fedJSON(t, `{"spec": {"futureField": 1}}`), fedJSON(t, `{"spec": {"futureField": 2}}`), false},
		{"the same field the type doesn't have", fedJSON(t, `{"spec": {"futureField": 1}}`), fedJSON(t, `{"spec": {"futureField": 1}}`), true},
		{"large integers", fedJSON(t, `{"spec": {"futureField": 9007199254740993}}`), fedJSON(t, `{"spec": {"futureField": 9007199254740992}}`), false},
		{"times within a second, as encoded", at(t, 10, 1), at(t, 10, 2), true},
		{"times a second apart", at(t, 10, 0), at(t, 11, 0), false},
		{"nil and a blank", nil, dieappsv1.DeploymentBlank, false},
		{"nils", nil, nil, true},
	}
	for _, p := range immutable {
		equal := p.x == p.y
		if p.x != nil {
			equal = p.x.Equal(p.y)
		}
		check(t, p.name, p.x, p.y, p.equal, equal)
	}

	mutable := []pair[*deploy]{
		{"mutable, same writes", named(t, "web").DieMutable(), named(t, "web").DieMutable(), true},
		{"mutable, different values", named(t, "web").DieMutable(), named(t, "api").DieMutable(), false},
	}
	for _, p := range mutable {
		check(t, p.name, p.x, p.y, p.equal, p.x.Equal(p.y))
	}

	// a die within another is compared by its own content and intent, at
	// paths within it, wherever its root is
	spec := func(d *dieappsv1.DeploymentDie) *dieappsv1.DeploymentSpecDie { return d.Spec() }
	replicas := func(n int32) *dieappsv1.DeploymentDie {
		return named(t, "web").DieStamp(func(d *deploy) { dietest.NoError(t, d.Spec().Replicas().Set(n)) })
	}
	children := []pair[*dieappsv1.DeploymentSpecDie]{
		{"children, same writes under different roots", spec(replicas(1)), spec(replicas(1).DieStamp(func(d *deploy) { dietest.NoError(t, d.Metadata().Name().Set("api")) })), true},
		{"children, different writes", spec(replicas(1)), spec(replicas(2)), false},
	}
	for _, p := range children {
		check(t, p.name, p.x, p.y, p.equal, p.x.Equal(p.y))
	}

	object := func(fn func(d *unstructured.ObjectMut)) *unstructured.Object {
		return unstructured.BlankFor("example.dev/v1", "Widget").DieStamp(fn)
	}
	objects := []pair[*unstructured.Object]{
		{"unstructured, same writes", object(func(d *unstructured.ObjectMut) {
			dietest.NoError(t, d.Node().Field("spec").Field("size").Set(int64(1)))
		}), object(func(d *unstructured.ObjectMut) {
			dietest.NoError(t, d.Node().Field("spec").Field("size").Set(int64(1)))
		}), true},
		{"unstructured, different writes", object(func(d *unstructured.ObjectMut) {
			dietest.NoError(t, d.Node().Field("spec").Field("size").Set(int64(1)))
		}), object(func(d *unstructured.ObjectMut) {
			dietest.NoError(t, d.Node().Field("spec").Field("size").Set(int64(2)))
		}), false},
	}
	for _, p := range objects {
		check(t, p.name, p.x, p.y, p.equal, p.x.Equal(p.y))
	}
}

func check(t *testing.T, name string, x, y any, want, equal bool) {
	t.Helper()
	diff := dietest.Diff(x, y)
	if equal != want {
		t.Errorf("%s: Equal %v, want %v", name, equal, want)
	}
	if (diff == "") != want {
		t.Errorf("%s: diff %q, want equal %v", name, diff, want)
	}
	if want != cmp.Equal(x, y, dietest.Transform) {
		t.Errorf("%s: cmp.Equal %v, want %v", name, !want, want)
	}
}

// A diff names the fields and intent that differ, not a die's internals.
func TestDiffIsReadable(t *testing.T) {
	x := named(t, "web").DieStamp(func(d *deploy) { dietest.NoError(t, d.Spec().Replicas().Set(1)) })
	y := named(t, "web").DieStamp(func(d *deploy) { dietest.NoError(t, d.Spec().Replicas().Set(2)) })
	diff := dietest.Diff(x, y)
	for _, want := range []string{`"replicas"`, "int64(1)", "int64(2)", `"set .spec.replicas = 1"`, `"set .spec.replicas = 2"`} {
		if !strings.Contains(diff, want) {
			t.Errorf("diff missing %s:\n%s", want, diff)
		}
	}
	for _, internal := range []string{"Holder", "state", "Onode"} {
		if strings.Contains(diff, internal) {
			t.Errorf("diff shows %s:\n%s", internal, diff)
		}
	}
	t.Logf("diff:\n%s", diff)
}

// Dies are compared by content wherever they are, such as the items of a
// list die or the fields of a struct.
func TestDiffWithinValues(t *testing.T) {
	type tc struct {
		Name   string
		Expect *dieappsv1.DeploymentDie
	}
	if diff := dietest.Diff(tc{"a", named(t, "web")}, tc{"a", named(t, "web")}); diff != "" {
		t.Errorf("unexpected diff:\n%s", diff)
	}
	if diff := dietest.Diff(tc{"a", named(t, "web")}, tc{"a", named(t, "api")}); !strings.Contains(diff, `"api"`) {
		t.Errorf("diff:\n%s", diff)
	}
	list := func(names ...string) *dieappsv1.DeploymentListDie {
		var items []*deploy
		for _, n := range names {
			items = append(items, named(t, n).DieMutable())
		}
		return dieappsv1.NewDeploymentList(items...)
	}
	if diff := dietest.Diff(list("web", "api"), list("web", "api")); diff != "" {
		t.Errorf("unexpected diff:\n%s", diff)
	}
	if diff := dietest.Diff(list("web", "api"), list("web", "db")); !strings.Contains(diff, `"db"`) {
		t.Errorf("diff:\n%s", diff)
	}
}

// Options for Go types don't apply to a die's content, but options for the
// content do, and with TransformRelease options for Go types do.
func TestDiffOptions(t *testing.T) {
	x := named(t, "web").DieStamp(func(d *deploy) { dietest.NoError(t, d.Metadata().ResourceVersion().Set("1")) })
	y := named(t, "web").DieStamp(func(d *deploy) { dietest.NoError(t, d.Metadata().ResourceVersion().Set("2")) })
	ignoreRV := cmpopts.IgnoreFields(metav1.ObjectMeta{}, "ResourceVersion")

	if diff := dietest.Diff(x, y, ignoreRV); diff == "" {
		t.Errorf("an option for a Go type applied to content")
	}
	byName := cmpopts.IgnoreMapEntries(func(k string, _ any) bool { return k == "resourceVersion" })
	if diff := dietest.Diff(x, y, byName); !strings.Contains(diff, "set .metadata.resourceVersion") {
		t.Errorf("intent should still differ:\n%s", diff)
	}
	if diff := dietest.Diff(x, y, byName, cmpopts.IgnoreFields(base.Content{}, "Intent")); diff != "" {
		t.Errorf("unexpected diff:\n%s", diff)
	}

	if diff := dietest.DiffRelease(x, y); !strings.Contains(diff, "ResourceVersion") {
		t.Errorf("diff:\n%s", diff)
	}
	if diff := dietest.DiffRelease(x, y, ignoreRV); diff != "" {
		t.Errorf("unexpected diff:\n%s", diff)
	}
	// options for the type's own struct apply, as to the type itself
	if diff := dietest.DiffRelease(x, y, cmpopts.IgnoreFields(appsv1.Deployment{}, "ObjectMeta.ResourceVersion")); diff != "" {
		t.Errorf("unexpected diff:\n%s", diff)
	}
	// intent isn't compared
	fed := fedJSON(t, `{"apiVersion": "apps/v1", "kind": "Deployment", "metadata": {"name": "web"}}`)
	if diff := dietest.DiffRelease(named(t, "web"), fed); diff != "" || named(t, "web").Equal(fed) {
		t.Errorf("diff %q, equal %v", diff, named(t, "web").Equal(fed))
	}
	// nor are nil dies a problem
	if diff := dietest.DiffRelease((*dieappsv1.DeploymentDie)(nil), dieappsv1.DeploymentBlank); diff == "" {
		t.Errorf("nil and a blank should differ")
	}
	t.Logf("release diff:\n%s", dietest.DiffRelease(x, y))
}
