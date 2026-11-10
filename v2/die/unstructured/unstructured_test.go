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

package unstructured_test

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"reconciler.io/dies/v2/die"
	dietesting "reconciler.io/dies/v2/die/testing"
	dieunstructured "reconciler.io/dies/v2/die/unstructured"
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

func assertPanicsWith(t *testing.T, want string, fn func()) {
	t.Helper()
	defer func() {
		t.Helper()
		if r := recover(); r == nil || !strings.Contains(r.(error).Error(), want) {
			t.Errorf("expected a panic containing %q, got %v", want, r)
		}
	}()
	fn()
}

// the state of an object of a kind with no Go types, as read from a server
func server() *unstructured.Unstructured {
	u := &unstructured.Unstructured{}
	_ = json.Unmarshal([]byte(`{
		"apiVersion": "acme.dev/v1", "kind": "Workload",
		"metadata": {"name": "web", "namespace": "default", "labels": {"team": "a"}},
		"spec": {
			"replicas": 3,
			"rollout": {"strategy": "canary"},
			"containers": [
				{"name": "app", "image": "app:1", "futureField": {"from": "a newer API"}},
				{"name": "proxy", "image": "proxy:1"}
			]
		}
	}`), &u.Object)
	return u
}

func TestUnstructuredDieTraversal(t *testing.T) {
	d := dieunstructured.BlankFor("acme.dev/v1", "Workload").DieStamp(func(d *dieunstructured.ObjectMut) {
		dietesting.NoError(t, d.Node().Field("metadata").Field("name").Set("web"))
		spec := d.Node().Field("spec")
		// an explicit zero is intent, as on a typed die
		dietesting.NoError(t, spec.Field("replicas").Set(0))
		dietesting.NoError(t, spec.Field("containers").Item("name", "app").Field("image").Set("app:1"))
		dietesting.NoError(t, spec.Field("volumes").EmptyList())
	})

	got, err := die.ReleaseApplyPatch(d)
	assertJSON(t, got, err, `{
		"apiVersion": "acme.dev/v1", "kind": "Workload",
		"metadata": {"name": "web"},
		"spec": {"replicas": 0, "containers": [{"name": "app", "image": "app:1"}], "volumes": []}
	}`)
	app := d.Node().Field("spec").Field("containers").Item("name", "app")
	if app.Field("image").Get() != "app:1" || dieunstructured.MustGetAs[int32](d.Node().Field("spec").Field("replicas")) != 0 {
		t.Errorf("image %v, replicas %v", app.Field("image").Get(), d.Node().Field("spec").Field("replicas").Get())
	}
	// reads create nothing
	if d.Node().Field("spec").Field("nothing").Has() || d.Node().Field("spec").Keys()[0] != "containers" {
		t.Errorf("keys %v", d.Node().Field("spec").Keys())
	}

}

func TestUnstructuredDieIsAnObject(t *testing.T) {
	var obj runtime.Object = dieunstructured.Blank.DieStamp(func(d *dieunstructured.ObjectMut) { dietesting.NoError(t, die.FeedUnstructured(d, server())) })
	if gvk := obj.GetObjectKind().GroupVersionKind(); gvk != (schema.GroupVersionKind{Group: "acme.dev", Version: "v1", Kind: "Workload"}) {
		t.Errorf("gvk %v", gvk)
	}
	d := dieunstructured.Blank.DieMutable()
	dietesting.NoError(t, die.FeedUnstructured(d, server()))
	accessor, err := meta.Accessor(d)
	if err != nil {
		t.Fatal(err)
	}
	if accessor.GetName() != "web" || accessor.GetLabels()["team"] != "a" {
		t.Errorf("name %q, labels %v", accessor.GetName(), accessor.GetLabels())
	}
	accessor.SetAnnotations(map[string]string{"note": "hi"})
	if got := d.Node().Field("metadata").Field("annotations").Field("note").Get(); got != "hi" {
		t.Errorf("annotation %v", got)
	}
}

// Fed from a server, only what is written is intent, and everything else,
// known or not, is kept.
func TestUnstructuredDiePatches(t *testing.T) {
	d := dieunstructured.Blank.DieStamp(func(d *dieunstructured.ObjectMut) { dietesting.NoError(t, die.FeedUnstructured(d, server())) }).DieStamp(func(d *dieunstructured.ObjectMut) {
		app := d.Node().Field("spec").Field("containers").Item("name", "app")
		dietesting.NoError(t, app.Field("env").Item("name", "MODE").Field("value").Set("full"))
		dietesting.NoError(t, d.Node().Field("metadata").Field("labels").Field("team").Clear())
	})

	// an object's field is a map entry to an unstructured die, which an apply
	// can't remove: it's left out
	got, err := die.ReleaseApplyPatch(d)
	assertJSON(t, got, err, `{
		"apiVersion": "acme.dev/v1", "kind": "Workload",
		"metadata": {"name": "web", "namespace": "default"},
		"spec": {"containers": [{"name": "app", "env": [{"name": "MODE", "value": "full"}]}]}
	}`)
	got, err = die.ReleaseMergePatch(d)
	assertJSON(t, got, err, `{
		"metadata": {"labels": {"team": null}},
		"spec": {"containers": [
			{"name": "app", "image": "app:1", "futureField": {"from": "a newer API"}, "env": [{"name": "MODE", "value": "full"}]},
			{"name": "proxy", "image": "proxy:1"}
		]}
	}`)
	u := die.MustReleaseUnstructured(d)
	if strategy, _, _ := unstructured.NestedString(u.Object, "spec", "rollout", "strategy"); strategy != "canary" {
		t.Errorf("unknown fields lost: %v", u.Object["spec"])
	}
	// a released object is an ordinary unstructured object
	_ = u.DeepCopy()
}

func TestUnstructuredDieKeysMustBeUnique(t *testing.T) {
	u := server()
	containers, _, _ := unstructured.NestedSlice(u.Object, "spec", "containers")
	_ = unstructured.SetNestedSlice(u.Object, append(containers, map[string]any{"name": "app"}), "spec", "containers")
	d := dieunstructured.Blank.DieStamp(func(d *dieunstructured.ObjectMut) { dietesting.NoError(t, die.FeedUnstructured(d, u)) })
	assertPanicsWith(t, "more than once", func() {
		d.Node().Field("spec").Field("containers").Item("name", "app").Get()
	})
}

func TestUnstructuredDieLists(t *testing.T) {
	list := dieunstructured.BlankFor("acme.dev/v1", "WorkloadList").DieStamp(func(d *dieunstructured.ObjectMut) {
		dietesting.NoError(t, d.Node().Field("items").Set([]any{
			map[string]any{"apiVersion": "acme.dev/v1", "kind": "Workload", "metadata": map[string]any{"name": "a"}},
			map[string]any{"apiVersion": "acme.dev/v1", "kind": "Workload", "metadata": map[string]any{"name": "b"}},
		}))
	})
	if !list.IsList() {
		t.Fatal("expected a list")
	}
	names := func(d *dieunstructured.Object) []string {
		var names []string
		_ = d.EachListItemWithAlloc(func(obj runtime.Object) error {
			names = append(names, obj.(*dieunstructured.ObjectMut).GetName())
			return nil
		})
		return names
	}

	// items are presented by the list, so a change to one changes the list
	list = list.DieStamp(func(d *dieunstructured.ObjectMut) {
		err := meta.EachListItem(d, func(obj runtime.Object) error {
			item := obj.(*dieunstructured.ObjectMut)
			item.SetName(item.GetName() + "-changed")
			dietesting.NoError(t, item.Node().Field("spec").Field("replicas").Set(2))
			dietesting.NoError(t, die.StampValue(item, func(m *map[string]any) { (*m)["stamped"] = true }))
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	})
	if got := names(list); !reflect.DeepEqual(got, []string{"a-changed", "b-changed"}) {
		t.Errorf("names %v", got)
	}
	if got := list.Node().Field("items").Get().([]any)[1].(map[string]any)["spec"]; !reflect.DeepEqual(got, map[string]any{"replicas": int64(2)}) {
		t.Errorf("spec %v", got)
	}
	if got := list.Node().Field("items").Get().([]any)[0].(map[string]any)["stamped"]; got != true {
		t.Errorf("stamped %v", got)
	}

	// the items of an immutable list are immutable
	_ = list.EachListItem(func(obj runtime.Object) error {
		item := obj.(*dieunstructured.Object)
		if !strings.HasSuffix(item.GetName(), "-changed") {
			t.Errorf("name %q", item.GetName())
		}
		assertPanicsWith(t, "immutable", func() { item.SetName("x") })
		return nil
	})

	// allocated items are copies, so a change to one leaves the list as it was
	list = list.DieStamp(func(d *dieunstructured.ObjectMut) {
		_ = meta.EachListItemWithAlloc(d, func(obj runtime.Object) error {
			obj.(*dieunstructured.ObjectMut).SetName("copy")
			return nil
		})
	})
	if got := names(list); !reflect.DeepEqual(got, []string{"a-changed", "b-changed"}) {
		t.Errorf("names %v", got)
	}

	// a copy of an item is unlinked from the list
	list = list.DieStamp(func(d *dieunstructured.ObjectMut) {
		_ = d.EachListItem(func(obj runtime.Object) error {
			obj.(*dieunstructured.ObjectMut).DeepCopy().SetName("copy")
			return nil
		})
	})
	if got := names(list); !reflect.DeepEqual(got, []string{"a-changed", "b-changed"}) {
		t.Errorf("names %v", got)
	}

	object := dieunstructured.Blank.DieStamp(func(d *dieunstructured.ObjectMut) { dietesting.NoError(t, die.FeedUnstructured(d, server())) })
	if object.IsList() {
		t.Error("an object is not a list")
	}
	if err := object.EachListItem(func(runtime.Object) error { return nil }); err == nil || !strings.Contains(err.Error(), "not a list") {
		t.Errorf("error %v", err)
	}
	bad := dieunstructured.BlankFor("acme.dev/v1", "WorkloadList").DieStamp(func(d *dieunstructured.ObjectMut) {
		dietesting.NoError(t, d.Node().Field("items").Set([]any{"not an object"}))
	})
	if err := bad.EachListItem(func(runtime.Object) error { return nil }); err == nil || !strings.Contains(err.Error(), "not an object") {
		t.Errorf("error %v", err)
	}
}

func TestUnstructuredDieIndex(t *testing.T) {
	d := dieunstructured.BlankFor("acme.dev/v1", "Workload").DieStamp(func(d *dieunstructured.ObjectMut) {
		dietesting.NoError(t, d.Node().Field("spec").Field("args").Set([]any{"serve", "--port", "8080"}))
		dietesting.NoError(t, d.Node().Field("spec").Field("rules").Set([]any{
			map[string]any{"hosts": []any{"a.dev"}},
			map[string]any{"hosts": []any{"b.dev"}},
		}))
	}).DieTare()

	args := d.Node().Field("spec").Field("args")
	if got := args.Index(1).Get(); got != "--port" {
		t.Errorf("arg %v", got)
	}
	if args.Index(3).Has() || args.Index(3).Get() != nil {
		t.Error("expected no element past the end")
	}
	if got := d.Node().Field("spec").Field("rules").Index(1).Field("hosts").Index(0).Get(); got != "b.dev" {
		t.Errorf("host %v", got)
	}

	d = d.DieStamp(func(d *dieunstructured.ObjectMut) {
		dietesting.NoError(t, d.Node().Field("spec").Field("args").Index(2).Set("9090"))
		dietesting.NoError(t, d.Node().Field("spec").Field("rules").Index(0).Field("hosts").Index(0).Set("c.dev"))
		dietesting.NoError(t, d.Node().Field("spec").Field("rules").Index(1).Field("paths").Set([]any{"/"}))
		dietesting.NoError(t, d.Node().Field("spec").Field("args").Index(0).Clear())
		// past the end of the list, as reflect's Index does
		assertPanicsWith(t, "out of range", func() { _ = d.Node().Field("spec").Field("args").Index(5).Set("x") })
	})

	// a write through an index sets the whole list
	got, err := die.ReleaseMergePatch(d)
	assertJSON(t, got, err, `{"spec": {
		"args": ["--port", "9090"],
		"rules": [{"hosts": ["c.dev"]}, {"hosts": ["b.dev"], "paths": ["/"]}]
	}}`)

	// an index sees what was written by key, and the reverse
	keyed := dieunstructured.Blank.DieStamp(func(d *dieunstructured.ObjectMut) { dietesting.NoError(t, die.FeedUnstructured(d, server())) }).DieStamp(func(d *dieunstructured.ObjectMut) {
		containers := d.Node().Field("spec").Field("containers")
		dietesting.NoError(t, containers.Item("name", "proxy").Field("image").Set("proxy:2"))
		if got := containers.Index(1).Field("image").Get(); got != "proxy:2" {
			t.Errorf("image %v", got)
		}
		dietesting.NoError(t, containers.Index(0).Field("image").Set("app:2"))
		if got := containers.Item("name", "app").Field("image").Get(); got != "app:2" {
			t.Errorf("image %v", got)
		}
	})
	if got := keyed.Node().Field("spec").Field("containers").Index(1).Field("image").Get(); got != "proxy:2" {
		t.Errorf("image %v", got)
	}
}

func TestUnstructuredNodeType(t *testing.T) {
	d := dieunstructured.BlankFor("acme.dev/v1", "Workload").DieStamp(func(d *dieunstructured.ObjectMut) {
		spec := d.Node().Field("spec")
		dietesting.NoError(t, spec.Field("args").Set([]string{"serve"}))
		dietesting.NoError(t, spec.Field("replicas").Set(3))
		dietesting.NoError(t, spec.Field("ratio").Set(0.5))
		dietesting.NoError(t, spec.Field("paused").Set(false))
		dietesting.NoError(t, spec.Field("cleared").Clear())
	})
	spec := d.Node().Field("spec")
	for node, want := range map[*dieunstructured.Node]dieunstructured.Type{
		spec:                        dieunstructured.TypeObject,
		spec.Field("args"):          dieunstructured.TypeList,
		spec.Field("args").Index(0): dieunstructured.TypeString,
		spec.Field("replicas"):      dieunstructured.TypeInt,
		spec.Field("ratio"):         dieunstructured.TypeFloat,
		spec.Field("paused"):        dieunstructured.TypeBool,
		spec.Field("cleared"):       dieunstructured.TypeAbsent,
		spec.Field("missing"):       dieunstructured.TypeAbsent,
		spec.Field("args").Index(1): dieunstructured.TypeAbsent,
	} {
		if got := node.Type(); got != want {
			t.Errorf("%s is %s, want %s", node.DiePath(), got, want)
		}
	}
}

func TestUnstructuredListDie(t *testing.T) {
	workload := func(name string) *dieunstructured.ObjectMut {
		return dieunstructured.BlankFor("acme.dev/v1", "Workload").DieStamp(func(d *dieunstructured.ObjectMut) {
			dietesting.NoError(t, d.Node().Field("metadata").Field("name").Set(name))
		}).DieMutable()
	}
	list := dieunstructured.NewList("acme.dev/v1", "WorkloadList", workload("a"), workload("b"))

	var _ metav1.ListInterface = list
	if !list.IsList() || len(list.Items) != 2 || list.Items[1].GetName() != "b" {
		t.Fatalf("items %v", list.Items)
	}
	if got := list.GetObjectKind().GroupVersionKind().Kind; got != "WorkloadList" {
		t.Errorf("kind %q", got)
	}
	list.SetResourceVersion("42")
	list.SetContinue("next")
	for _, item := range list.Items {
		item.SetLabels(map[string]string{"seen": "true"})
	}

	// apimachinery reads it as an unstructured list, and by its items
	var labeled []string
	for _, each := range []func(runtime.Object, func(runtime.Object) error) error{meta.EachListItem, meta.EachListItemWithAlloc} {
		labeled = nil
		_ = each(list, func(obj runtime.Object) error {
			o, _ := meta.Accessor(obj)
			labeled = append(labeled, o.GetName()+"="+o.GetLabels()["seen"])
			return nil
		})
		if !reflect.DeepEqual(labeled, []string{"a=true", "b=true"}) {
			t.Errorf("labeled %v", labeled)
		}
	}
	if objs, err := meta.ExtractList(list); err != nil || len(objs) != 2 {
		t.Errorf("extracted %v, %v", objs, err)
	}

	// round trips through an unstructured list
	u := list.UnstructuredList()
	if len(u.Items) != 2 || u.Items[0].GetLabels()["seen"] != "true" || u.GetResourceVersion() != "42" || u.GetContinue() != "next" {
		t.Errorf("released %v", u.UnstructuredContent())
	}
	back := dieunstructured.MustListFrom(u)
	if back.Kind != "WorkloadList" || back.GetResourceVersion() != "42" || back.Items[0].GetName() != "a" {
		t.Errorf("fed %v", back.UnstructuredContent())
	}
	// copies are independent
	copied := list.DeepCopy()
	copied.Items[0].SetName("copy")
	if list.Items[0].GetName() != "a" {
		t.Errorf("name %q", list.Items[0].GetName())
	}
}

func TestUnstructuredDieToList(t *testing.T) {
	if _, ok := dieunstructured.Blank.DieStamp(func(d *dieunstructured.ObjectMut) { dietesting.NoError(t, die.FeedUnstructured(d, server())) }).ToList(); ok {
		t.Error("an object is not a list")
	}

	d := dieunstructured.BlankFor("acme.dev/v1", "WorkloadList").DieStamp(func(d *dieunstructured.ObjectMut) {
		dietesting.NoError(t, d.Node().Field("metadata").Field("resourceVersion").Set("7"))
		dietesting.NoError(t, d.Node().Field("items").Set([]any{map[string]any{"metadata": map[string]any{"name": "a"}}}))
	})
	list, ok := d.ToList()
	if !ok {
		t.Fatal("expected a list")
	}
	if list.Kind != "WorkloadList" || list.GetResourceVersion() != "7" || list.Items[0].GetName() != "a" {
		t.Errorf("list %v", list.UnstructuredContent())
	}
	// the list is a copy, with mutable items
	list.Items[0].SetName("b")
	if got := d.Node().Field("items").Index(0).Field("metadata").Field("name").Get(); got != "a" {
		t.Errorf("name %v", got)
	}
}

// Clients decode into an unstructured die with UnmarshalJSON, which holds
// numbers as unstructured objects do.
func TestUnmarshalJSONNumbers(t *testing.T) {
	d := dieunstructured.Blank.DieMutable()
	if err := json.Unmarshal([]byte(`{"spec": {"replicas": 3, "ratio": 0.5}}`), d); err != nil {
		t.Fatal(err)
	}
	spec := d.Node().Field("spec")
	if got := spec.Field("replicas").Get(); got != int64(3) {
		t.Errorf("replicas %v (%T)", got, got)
	}
	if got := spec.Field("ratio").Get(); got != 0.5 {
		t.Errorf("ratio %v (%T)", got, got)
	}
}

// An unstructured die's content has no schema, so it is traversed as reflect
// traverses values: a read below a value of another kind reads as absent,
// and a write that would step into one panics, without changing the die.
// Type and Len say what a node holds, so careful code never panics.
func TestUnstructuredTraversalByKind(t *testing.T) {
	d := dieunstructured.Blank.DieMutable()
	dietesting.NoError(t, die.FeedJSON(d, []byte(`{"spec": {"name": "web", "args": ["a"], "ports": [{"name": "http"}, {"name": "http"}]}}`)))
	spec := d.Node().Field("spec")

	if got := spec.Field("name").Field("deeper"); got.Get() != nil || got.Has() || got.Type() != dieunstructured.TypeAbsent {
		t.Errorf("a field of a string reads %v", got.Get())
	}
	if got := spec.Field("args").Field("x").Get(); got != nil {
		t.Errorf("a field of a list reads %v", got)
	}
	if got := spec.Field("name").Index(0).Get(); got != nil {
		t.Errorf("an element of a string reads %v", got)
	}
	if got := spec.Field("args").Index(5).Get(); got != nil {
		t.Errorf("an element past the end reads %v", got)
	}

	before := string(die.MustReleaseJSON(d))
	assertPanicsWith(t, "is a string", func() { _ = spec.Field("name").Field("deeper").Set("v") })
	assertPanicsWith(t, "is a list", func() { _ = spec.Field("args").Field("x").Set("v") })
	assertPanicsWith(t, "out of range", func() { _ = spec.Field("args").Index(5).Set("v") })
	if after := string(die.MustReleaseJSON(d)); after != before {
		t.Errorf("a refused write changed the die:\n%s\n%s", before, after)
	}
	// so the die is still whole
	if spec.Field("name").Type() != dieunstructured.TypeString || spec.Field("args").Len() != 1 {
		t.Errorf("spec %v", spec.Get())
	}

	// a key with no schema to refuse it is ambiguous where it's used; a
	// filter finds every element with it
	assertPanicsWith(t, "more than once", func() { spec.Field("ports").Item("name", "http").Get() })
	if got := spec.Field("ports").MustFind(`[?(@.name=="http")]`); len(got) != 2 {
		t.Errorf("found %d elements", len(got))
	}
}
