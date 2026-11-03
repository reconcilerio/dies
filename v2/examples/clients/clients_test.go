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

package clients_test

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"reconciler.io/runtime/duck"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	"reconciler.io/dies/v2/die"
	"reconciler.io/dies/v2/die/dietest"
	dieunstructured "reconciler.io/dies/v2/die/unstructured"
	jigsv1 "reconciler.io/dies/v2/examples/jigs/v1"
)

var workloadGVK = schema.GroupVersionKind{Group: "acme.dev", Version: "v1", Kind: "Workload"}

// a kind with no Go types, stored as unstructured
func stored(t *testing.T) *unstructured.Unstructured {
	t.Helper()
	u := &unstructured.Unstructured{}
	if err := json.Unmarshal([]byte(`{
		"apiVersion": "acme.dev/v1", "kind": "Workload",
		"metadata": {"name": "web", "namespace": "default"},
		"spec": {
			"rollout": {"strategy": "canary"},
			"template": {"spec": {"containers": [{"name": "app", "image": "app:1", "futureField": true}]}}
		}
	}`), &u.Object); err != nil {
		t.Fatal(err)
	}
	return u
}

func newClient(t *testing.T, objs ...client.Object) client.WithWatch {
	t.Helper()
	mapper := meta.NewDefaultRESTMapper(nil)
	mapper.Add(workloadGVK, meta.RESTScopeNamespace)
	mapper.Add(schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}, meta.RESTScopeNamespace)
	// a scheme that knows neither the workload kind nor the dies
	return fake.NewClientBuilder().WithScheme(runtime.NewScheme()).WithRESTMapper(mapper).WithObjects(objs...).Build()
}

func key() client.ObjectKey { return client.ObjectKey{Namespace: "default", Name: "web"} }

// controller-runtime's client sends any runtime.Unstructured through its
// unstructured client, so an unstructured die needs no registration to read.
// The fake client's writes depend on the concrete unstructured type, so they
// are tested with the die registered, below; the real client's writes route
// by the interface too.
func TestClientReadsAnUnregisteredUnstructuredDie(t *testing.T) {
	c := newClient(t, stored(t))
	d := dieunstructured.BlankFor("acme.dev/v1", "Workload").DieMutable()
	if err := c.Get(context.Background(), key(), d); err != nil {
		t.Fatal(err)
	}
	if got := d.Node().Field("spec").Field("template").Field("spec").Field("containers").Item("name", "app").Field("image").Get(); got != "app:1" {
		t.Errorf("image %v", got)
	}
}

// Registered in a scheme for a kind, an unstructured die is that kind's Go
// type: clients decode into it, and informers create it, so the cache can
// hand it out.
func TestClientTakesARegisteredUnstructuredDie(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	scheme.AddKnownTypeWithName(workloadGVK, &dieunstructured.ObjectMut{})
	mapper := meta.NewDefaultRESTMapper(nil)
	mapper.Add(workloadGVK, meta.RESTScopeNamespace)
	seed := dieunstructured.Blank.DieMutable()
	dietest.NoError(t, die.FeedUnstructured(seed, stored(t)))
	c := fake.NewClientBuilder().WithScheme(scheme).WithRESTMapper(mapper).WithObjects(seed).Build()

	d := dieunstructured.BlankFor("acme.dev/v1", "Workload").DieMutable()
	if err := c.Get(ctx, key(), d); err != nil {
		t.Fatal(err)
	}

	// change it through a typed jig, and write it back with an exact patch
	ps, _ := jigsv1.PodSpecableJig.OfMut(d)
	dietest.NoError(t, ps.Template().Spec().Containers().At("app").Image().Set("app:2"))
	patch, err := die.ReleaseMergePatchE(d)
	if err != nil {
		t.Fatal(err)
	}
	assertJSON(t, patch, `{"spec": {"template": {"spec": {"containers": [{"name": "app", "image": "app:2", "futureField": true}]}}}}`)
	if err := c.Patch(ctx, d, client.RawPatch(types.MergePatchType, patch)); err != nil {
		t.Fatal(err)
	}

	got := dieunstructured.BlankFor("acme.dev/v1", "Workload").DieMutable()
	if err := c.Get(ctx, key(), got); err != nil {
		t.Fatal(err)
	}
	app := got.Node().Field("spec").Field("template").Field("spec").Field("containers").Item("name", "app")
	if app.Field("image").Get() != "app:2" || app.Field("futureField").Get() != true {
		t.Errorf("container %v", app.Get())
	}
	if got.Node().Field("spec").Field("rollout").Field("strategy").Get() != "canary" {
		t.Errorf("rollout lost")
	}

	// update and create take it too
	dietest.NoError(t, got.Node().Field("metadata").Field("labels").Field("tier").Set("web"))
	if err := c.Update(ctx, got); err != nil {
		t.Fatal(err)
	}
	created := dieunstructured.BlankFor("acme.dev/v1", "Workload").DieStamp(func(d *dieunstructured.ObjectMut) {
		d.SetNamespace("default")
		d.SetName("api")
		dietest.NoError(t, d.Node().Field("spec").Field("replicas").Set(1))
	}).DieMutable()
	if err := c.Create(ctx, created); err != nil {
		t.Fatal(err)
	}
	if created.GetResourceVersion() == "" {
		t.Errorf("the client should decode the created object back into the die")
	}
}

func assertJSON(t *testing.T, got []byte, want string) {
	t.Helper()
	var g, w any
	_ = json.Unmarshal(got, &g)
	_ = json.Unmarshal([]byte(want), &w)
	if !reflect.DeepEqual(g, w) {
		t.Errorf("json mismatch\n got: %s\nwant: %s", got, want)
	}
}

// reconciler.io/runtime's duck wrapper converts an object the scheme doesn't
// know, with its kind set, through unstructured, so a typed die needs no
// registration either.
func TestDuckClientTakesDies(t *testing.T) {
	ctx := context.Background()
	deploy := &unstructured.Unstructured{}
	deploy.SetGroupVersionKind(schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"})
	deploy.SetNamespace("default")
	deploy.SetName("web")
	_ = unstructured.SetNestedField(deploy.Object, int64(3), "spec", "replicas")
	c := newClient(t, stored(t), deploy)
	reader := duck.NewDuckAwareAPIReaderWrapper(c, c)

	// a typed die for a type the scheme doesn't have
	typed := dieappsv1.DeploymentBlank.DieMutable()
	if !duck.IsDuck(typed, c.Scheme()) {
		t.Fatal("expected a duck")
	}
	if err := reader.Get(ctx, key(), typed); err != nil {
		t.Fatal(err)
	}
	if got := typed.Spec().Replicas().Get(); got != 3 {
		t.Errorf("replicas %d", got)
	}

	// an unstructured die isn't a duck: schemes know any runtime.Unstructured
	// by its kind, so the wrapper passes it straight to the client
	unstructuredDie := dieunstructured.BlankFor("acme.dev/v1", "Workload").DieMutable()
	if duck.IsDuck(unstructuredDie, c.Scheme()) {
		t.Fatal("did not expect a duck")
	}
	if err := reader.Get(ctx, key(), unstructuredDie); err != nil {
		t.Fatal(err)
	}
	if got := unstructuredDie.GetName(); got != "web" {
		t.Errorf("name %q", got)
	}
}

// List dies are lists of dies, as Kubernetes lists are, so clients that fill
// a list by reflecting on its items, as the fake client and the informer
// cache do through apimachinery's meta.SetList, list into them, once a kind's
// dies are registered as its types.
func TestClientListsIntoRegisteredListDies(t *testing.T) {
	deployGVK := schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}
	scheme := runtime.NewScheme()
	scheme.AddKnownTypeWithName(workloadGVK, &dieunstructured.ObjectMut{})
	scheme.AddKnownTypeWithName(workloadGVK.GroupVersion().WithKind("WorkloadList"), &dieunstructured.ObjectList{})
	scheme.AddKnownTypeWithName(deployGVK, &dieappsv1.DeploymentDieMut{})
	scheme.AddKnownTypeWithName(deployGVK.GroupVersion().WithKind("DeploymentList"), &dieappsv1.DeploymentListDie{})
	mapper := meta.NewDefaultRESTMapper(nil)
	mapper.Add(workloadGVK, meta.RESTScopeNamespace)
	mapper.Add(deployGVK, meta.RESTScopeNamespace)

	// the fake client keeps objects as they are given; seeded as unstructured
	// objects, they would list as unstructured objects, which aren't dies
	web := dieunstructured.Blank.DieMutable()
	dietest.NoError(t, die.FeedUnstructured(web, stored(t)))
	api := web.DeepCopy()
	api.SetName("api")
	deploy := dieappsv1.DeploymentBlank.DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietest.NoError(t, d.Metadata().Namespace().Set("default"))
		dietest.NoError(t, d.Metadata().Name().Set("web"))
		dietest.NoError(t, d.Spec().Template().Spec().Containers().At("app").Image().Set("app:1"))
	}).DieTare().DieMutable()
	c := fake.NewClientBuilder().WithScheme(scheme).WithRESTMapper(mapper).WithObjects(web, api, deploy).Build()
	ctx := context.Background()

	workloads := dieunstructured.NewList("acme.dev/v1", "WorkloadList")
	if err := c.List(ctx, workloads, client.InNamespace("default")); err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, item := range workloads.Items {
		names = append(names, item.GetName())
	}
	if !reflect.DeepEqual(names, []string{"api", "web"}) {
		t.Errorf("names %v", names)
	}
	// items keep fields no Go type describes
	if got := workloads.Items[1].Node().Field("spec").Field("template").Field("spec").Field("containers").Index(0).Field("futureField").Get(); got != true {
		t.Errorf("futureField %v", got)
	}

	deployments := dieappsv1.NewDeploymentList()
	if err := c.List(ctx, deployments, client.InNamespace("default")); err != nil {
		t.Fatal(err)
	}
	if len(deployments.Items) != 1 || deployments.Items[0].Spec().Template().Spec().Containers().At("app").Image().Get() != "app:1" {
		t.Errorf("deployments %v", deployments.Items)
	}
	// items are mutable dies, ready to change and write back
	item := deployments.Items[0]
	dietest.NoError(t, item.Spec().Replicas().Set(3))
	if err := c.Update(ctx, item); err != nil {
		t.Fatal(err)
	}
	reloaded := dieappsv1.DeploymentBlank.DieMutable()
	if err := c.Get(ctx, key(), reloaded); err != nil {
		t.Fatal(err)
	}
	if got := reloaded.Spec().Replicas().Get(); got != 3 {
		t.Errorf("replicas %v", got)
	}
}

// Through reconciler.io/runtime's duck wrapper, list dies need no
// registration: a typed list die the scheme doesn't know is a duck, listed as
// an unstructured list and converted, items and all, through the dies' JSON
// methods. As the wrapped reader only sees unstructured objects, a cache
// behind it serves typed dies too. An unstructured die isn't a duck, so a
// cache still sees it as itself, and can't fill it from the unstructured
// objects its informers hold.
func TestDuckClientListsIntoListDies(t *testing.T) {
	seed := &unstructured.Unstructured{}
	seed.SetGroupVersionKind(schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"})
	seed.SetNamespace("default")
	seed.SetName("web")
	_ = unstructured.SetNestedSlice(seed.Object, []any{map[string]any{"name": "app", "image": "app:1"}}, "spec", "template", "spec", "containers")
	c := newClient(t, seed, stored(t))
	reader := duck.NewDuckAwareAPIReaderWrapper(c, c)
	ctx := context.Background()

	deployments := dieappsv1.NewDeploymentList()
	if !duck.IsDuck(deployments, c.Scheme()) {
		t.Fatal("expected a duck")
	}
	if err := reader.List(ctx, deployments, client.InNamespace("default")); err != nil {
		t.Fatal(err)
	}
	if len(deployments.Items) != 1 || deployments.Items[0].Spec().Template().Spec().Containers().At("app").Image().Get() != "app:1" {
		t.Errorf("deployments %v", deployments.Items)
	}

	// an unstructured list die isn't a duck, and the client lists into it
	// unregistered
	workloads := dieunstructured.NewList("acme.dev/v1", "WorkloadList")
	if duck.IsDuck(workloads, c.Scheme()) {
		t.Error("expected an unstructured list die not to be a duck")
	}
	if err := reader.List(ctx, workloads, client.InNamespace("default")); err != nil {
		t.Fatal(err)
	}
	if len(workloads.Items) != 1 || workloads.Items[0].GetName() != "web" {
		t.Errorf("workloads %v", workloads.Items)
	}
}
