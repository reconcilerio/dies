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
	k8sunstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	"reconciler.io/dies/v2/die/dietest"
	"reconciler.io/dies/v2/die/unstructured"
)

type umut = unstructured.ObjectMut

func widget(fn func(d *umut)) *unstructured.Object {
	return unstructured.BlankFor("example.dev/v1", "Widget").DieWith(fn)
}

func fedWidget(data string) *unstructured.Object {
	return unstructured.Blank.DieWith(func(d *umut) { d.DieFeedJSON([]byte(data)) })
}

func containers(d *umut) *unstructured.NodeMut {
	return d.Node().Field("spec").Field("containers")
}

// Unstructured dies have a diff exactly when they aren't Equal, as typed dies
// do, though their values have no Go types to normalize them.
func TestUnstructuredDiffIsEmptyExactlyWhenEqual(t *testing.T) {
	size := func(v any) *unstructured.Object {
		return widget(func(d *umut) { d.Node().Field("spec").Field("size").Set(v) })
	}
	fed := `{"apiVersion": "example.dev/v1", "kind": "Widget", "spec": {"size": 1, "ratio": 0.5, "big": 9007199254740993}}`
	pairs := []struct {
		name  string
		x, y  *unstructured.Object
		equal bool
	}{
		{"blanks", unstructured.Blank, unstructured.Blank, true},
		{"same writes", size(int64(1)), size(int64(1)), true},
		{"different writes", size(int64(1)), size(int64(2)), false},
		// Go's integer types encode alike
		{"int and int64", size(1), size(int64(1)), true},
		{"an integer and a float", size(int64(1)), size(1.5), false},
		{"a number and a string", size(int64(1)), size("1"), false},
		{"null and absent", size(nil), widget(func(d *umut) { d.Node().Field("spec").EmptyObject() }), false},
		{"empty object and empty list", widget(func(d *umut) { d.Node().Field("spec").Field("x").EmptyObject() }), widget(func(d *umut) { d.Node().Field("spec").Field("x").EmptyList() }), false},
		{"fed alike", fedWidget(fed), fedWidget(fed), true},
		{"large integers fed", fedWidget(`{"spec": {"big": 9007199254740993}}`), fedWidget(`{"spec": {"big": 9007199254740992}}`), false},
		{"fed, then the same value written", fedWidget(fed), fedWidget(fed).DieWith(func(d *umut) { d.Node().Field("spec").Field("size").Set(int64(1)) }), false},
		{"list items in a different order", widget(func(d *umut) {
			containers(d).Set([]any{map[string]any{"name": "a"}, map[string]any{"name": "b"}})
		}), widget(func(d *umut) {
			containers(d).Set([]any{map[string]any{"name": "b"}, map[string]any{"name": "a"}})
		}), false},
		{"nil and a blank", nil, unstructured.Blank, false},
	}
	for _, p := range pairs {
		equal := p.x == p.y
		if p.x != nil {
			equal = p.x.Equal(p.y)
		}
		check(t, p.name, p.x, p.y, p.equal, equal)
	}

	// as for mutable dies
	m := func(v int64) *umut {
		return widget(func(d *umut) { d.Node().Field("spec").Field("size").Set(v) }).DieMutable()
	}
	check(t, "mutable, same", m(1), m(1), true, m(1).Equal(m(1)))
	check(t, "mutable, different", m(1), m(2), false, m(1).Equal(m(2)))
}

// A diff names the fields and intent that differ, as for typed dies.
func TestUnstructuredDiffIsReadable(t *testing.T) {
	x := widget(func(d *umut) { containers(d).Item("name", "app").Field("image").Set("app:1") })
	y := widget(func(d *umut) { containers(d).Item("name", "app").Field("image").Set("app:2") })
	diff := dietest.Diff(x, y)
	for _, want := range []string{`"app:1"`, `"app:2"`, `set .spec.containers[name=app].image = "app:1"`} {
		if !strings.Contains(diff, want) {
			t.Errorf("diff missing %s:\n%s", want, diff)
		}
	}
	t.Logf("diff:\n%s", diff)
}

// Released, an unstructured die is its map, compared as go-cmp compares maps.
func TestUnstructuredDiffRelease(t *testing.T) {
	size := func(v any) *unstructured.Object {
		return widget(func(d *umut) { d.Node().Field("spec").Field("size").Set(v) })
	}
	if diff := dietest.DiffRelease(size(int64(1)), size(int64(1))); diff != "" {
		t.Errorf("unexpected diff:\n%s", diff)
	}
	if diff := dietest.DiffRelease(size(int64(1)), size(int64(2))); !strings.Contains(diff, "int64(2)") {
		t.Errorf("diff:\n%s", diff)
	}
	// an unstructured die holds integers as int64, as decoded JSON does
	if !size(1).Equal(size(int64(1))) || dietest.DiffRelease(size(1), size(int64(1))) != "" {
		t.Errorf("an int and an int64 should be the same value: %s", dietest.DiffRelease(size(1), size(int64(1))))
	}
}

// Locations within unstructured dies ignore what they name, including list
// items by their keys, which an unstructured die has no types to know.
func TestUnstructuredIgnoreFields(t *testing.T) {
	blank := unstructured.Blank.Node().Field("spec").Field("containers")
	with := func(fns ...func(d *umut)) *unstructured.Object {
		return widget(func(d *umut) {
			for _, fn := range fns {
				fn(d)
			}
		})
	}
	image := func(name, image string) func(d *umut) {
		return func(d *umut) { containers(d).Item("name", name).Field("image").Set(image) }
	}
	ignoreSidecar := dietest.IgnoreFields(blank.Item("name", "sidecar"))
	assertNoDiffs(t, "a sidecar on one side", with(image("app", "app:1"), image("sidecar", "proxy:1")), with(image("app", "app:1")), ignoreSidecar)
	assertNoDiffs(t, "different sidecars", with(image("app", "app:1"), image("sidecar", "proxy:1")), with(image("app", "app:1"), image("sidecar", "proxy:2")), ignoreSidecar)
	assertDiffs(t, "different apps", with(image("app", "app:1"), image("sidecar", "proxy:1")), with(image("app", "app:2"), image("sidecar", "proxy:1")), "app:2", ignoreSidecar)

	// by its index, which ignores its value, but not its intent, which
	// locates it by its key
	byIndex := dietest.IgnoreFields(blank.Index(0))
	if diff := dietest.DiffRelease(with(image("a", "x:1"), image("b", "y:1")), with(image("a", "x:2"), image("b", "y:1")), byIndex); diff != "" {
		t.Errorf("unexpected release diff:\n%s", diff)
	}
	if diff := dietest.Diff(with(image("a", "x:1"), image("b", "y:1")), with(image("a", "x:2"), image("b", "y:1")), byIndex); !strings.Contains(diff, "Intent") || strings.Contains(diff, `string("x:2")`) {
		t.Errorf("diff should show only intent:\n%s", diff)
	}

	// a field below an item
	assertNoDiffs(t, "an item's field", with(image("app", "app:1")), with(image("app", "app:2")), dietest.IgnoreFields(blank.Item("name", "app").Field("image")))

	// fed values, whose intent sets them whole, aren't ignored by a location
	// within them, but are by one enclosing them
	fed := func(image string) *unstructured.Object {
		return fedWidget(`{"spec": {"containers": [{"name": "app", "image": "` + image + `"}]}}`)
	}
	if diff := dietest.Diff(fed("app:1"), fed("app:2"), dietest.IgnoreFields(blank)); diff != "" {
		t.Errorf("unexpected diff:\n%s", diff)
	}
}

// A typed view of an unstructured die compares as a typed die does, with
// the unstructured die's intent at its own paths.
func TestUnstructuredJigAsDiff(t *testing.T) {
	deployment := func(image string) *unstructured.Object {
		return unstructured.BlankFor("apps/v1", "Deployment").DieWith(func(d *umut) {
			d.Node().Field("metadata").Field("name").Set("web")
			d.Node().Field("spec").Field("template").Field("spec").Field("containers").Item("name", "app").Field("image").Set(image)
		})
	}
	x := unstructured.JigAs[*dieappsv1.DeploymentDie](deployment("app:1"))
	y := unstructured.JigAs[*dieappsv1.DeploymentDie](deployment("app:2"))
	check(t, "jigs, different", x, y, false, x.Equal(y))
	z := unstructured.JigAs[*dieappsv1.DeploymentDie](deployment("app:1"))
	check(t, "jigs, same", x, z, true, x.Equal(z))
	if diff := dietest.Diff(x, y); !strings.Contains(diff, `string("app:2")`) || !strings.Contains(diff, "set .spec.template.spec.containers[name=app].image") {
		t.Errorf("diff:\n%s", diff)
	}
	ignore := dietest.IgnoreFields(dieappsv1.DeploymentBlank.Spec().Template().Spec().Containers().At("app").Image())
	assertNoDiffs(t, "jigs, ignored", x, y, ignore)

	// and as the typed die written alike
	typed := dieappsv1.DeploymentBlank.DieWith(func(d *dieappsv1.DeploymentDieMut) {
		d.Metadata().Name().Set("web")
		d.Spec().Template().Spec().Containers().At("app").Image().Set("app:1")
	})
	check(t, "a jig and a typed die", x, typed, true, x.Equal(typed))
}

// Unstructured list dies compare by their items' content.
func TestUnstructuredListDiff(t *testing.T) {
	item := func(name string) *umut {
		return widget(func(d *umut) { d.Node().Field("metadata").Field("name").Set(name) }).DieMutable()
	}
	list := func(names ...string) *unstructured.ObjectList {
		var items []*umut
		for _, n := range names {
			items = append(items, item(n))
		}
		return unstructured.NewList("example.dev/v1", "WidgetList", items...)
	}
	if diff := dietest.Diff(list("a", "b"), list("a", "b")); diff != "" {
		t.Errorf("unexpected diff:\n%s", diff)
	}
	if diff := dietest.Diff(list("a", "b"), list("a", "c")); !strings.Contains(diff, `"c"`) {
		t.Errorf("diff:\n%s", diff)
	}
	// as from an unstructured list
	u := &k8sunstructured.UnstructuredList{Object: map[string]any{"apiVersion": "example.dev/v1", "kind": "WidgetList"}}
	u.Items = []k8sunstructured.Unstructured{{Object: map[string]any{"apiVersion": "example.dev/v1", "kind": "Widget", "metadata": map[string]any{"name": "a"}}}}
	if !cmp.Equal(unstructured.ListFrom(u), unstructured.ListFrom(u.DeepCopy()), dietest.Transform) {
		t.Errorf("expected lists from the same unstructured list to be equal")
	}
}
