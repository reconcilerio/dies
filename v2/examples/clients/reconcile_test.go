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
	"fmt"
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	appsv1ac "k8s.io/client-go/applyconfigurations/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	diecorev1 "reconciler.io/dies/v2/apis/core/v1"
	"reconciler.io/dies/v2/die/base"
	dieunstructured "reconciler.io/dies/v2/die/unstructured"
)

// reconcile is a controller's loop over a Deployment: get it, work out the
// status this controller owns, and apply only that, from a blank, to the
// status subresource.
func reconcile(ctx context.Context, c client.Client, key client.ObjectKey) error {
	actual, err := dieappsv1.DeploymentBlank.DieWithE(func(d *dieappsv1.DeploymentDieMut) error {
		return c.Get(ctx, key, d)
	})
	if err != nil {
		return err
	}

	status := dieappsv1.DeploymentBlank.DieWith(func(d *dieappsv1.DeploymentDieMut) {
		d.Metadata().Namespace().Set(actual.GetNamespace())
		d.Metadata().Name().Set(actual.GetName())
		d.Status().ObservedGeneration().Set(actual.GetGeneration())
		d.Status().Replicas().Set(actual.Spec().Replicas().Get())
		d.Status().Conditions().At(dieappsv1.DeploymentAvailable).DieWith(func(c *dieappsv1.DeploymentConditionDieMut) {
			c.Status().Set(diecorev1.ConditionTrue)
			c.Reason().Set("Reconciled")
		})
	}).DieSubresource("status")
	body, err := status.DieReleaseApplyPatchE()
	if err != nil {
		return err
	}
	// the client decodes the server's response into the die it's given
	applied, err := status.DieWithE(func(d *dieappsv1.DeploymentDieMut) error {
		return c.Status().Patch(ctx, d, client.RawPatch(types.ApplyPatchType, body), client.FieldOwner("litmus"), client.ForceOwnership)
	})
	if err != nil {
		return err
	}
	if applied.Status().ObservedGeneration().Get() != actual.GetGeneration() {
		return fmt.Errorf("the server has observed generation %d, applied %d", applied.Status().ObservedGeneration().Get(), actual.GetGeneration())
	}
	return nil
}

func TestReconcileAppliesOnlyWhatItOwns(t *testing.T) {
	gvk := schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}
	scheme := runtime.NewScheme()
	if err := dieappsv1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	mapper := meta.NewDefaultRESTMapper(nil)
	mapper.Add(gvk, meta.RESTScopeNamespace)

	seed := dieappsv1.DeploymentBlank.DieWith(func(d *dieappsv1.DeploymentDieMut) {
		d.Metadata().Namespace().Set("default")
		d.Metadata().Name().Set("web")
		d.Spec().Replicas().Set(2)
		d.Spec().Template().Spec().Containers().At("app").Image().Set("app:1")
	}).DieCommit().DieMutable()
	c := fake.NewClientBuilder().WithScheme(scheme).WithRESTMapper(mapper).
		WithObjects(seed).WithStatusSubresource(&dieappsv1.DeploymentDieMut{}).WithReturnManagedFields().Build()
	ctx := context.Background()

	// another controller owns a condition of its own
	other := dieappsv1.DeploymentBlank.DieWith(func(d *dieappsv1.DeploymentDieMut) {
		d.Metadata().Namespace().Set("default")
		d.Metadata().Name().Set("web")
		d.Status().Conditions().At(dieappsv1.DeploymentProgressing).Status().Set(diecorev1.ConditionTrue)
	}).DieSubresource("status")
	body, err := other.DieReleaseApplyPatchE()
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Status().Patch(ctx, other.DieMutable(), client.RawPatch(types.ApplyPatchType, body), client.FieldOwner("other")); err != nil {
		t.Fatal(err)
	}

	for range 2 {
		if err := reconcile(ctx, c, key()); err != nil {
			t.Fatal(err)
		}
	}

	got := dieappsv1.DeploymentBlank.DieMutable()
	if err := c.Get(ctx, key(), got); err != nil {
		t.Fatal(err)
	}
	status := got.Status()
	if status.Replicas().Get() != 2 {
		t.Errorf("replicas %d", status.Replicas().Get())
	}
	conditions := status.Conditions()
	if conditions.At(dieappsv1.DeploymentAvailable).Reason().Get() != "Reconciled" {
		t.Errorf("available %+v", conditions.At(dieappsv1.DeploymentAvailable).DieRelease())
	}
	// the other owner's condition, and the spec, are as they were
	if conditions.At(dieappsv1.DeploymentProgressing).Status().Get() != diecorev1.ConditionTrue {
		t.Errorf("progressing %+v", conditions.At(dieappsv1.DeploymentProgressing).DieRelease())
	}
	if got.Spec().Template().Spec().Containers().At("app").Image().Get() != "app:1" {
		t.Errorf("spec %+v", got.Spec().DieRelease())
	}
	// each controller applied as its own field owner (the fake client does not
	// record the subresource an apply was made to)
	owners := map[string]bool{}
	for _, mf := range got.GetManagedFields() {
		if mf.Operation == "Apply" {
			owners[mf.Manager] = true
		}
	}
	if !owners["litmus"] || !owners["other"] {
		t.Errorf("owners %v", owners)
	}
}

// A typed die keeps the fields its Go type doesn't have, such as those of a
// newer API, so an update doesn't remove them.
func TestUpdateKeepsUnknownFields(t *testing.T) {
	scheme := runtime.NewScheme()
	if err := dieappsv1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	gvk := schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}
	mapper := meta.NewDefaultRESTMapper(nil)
	mapper.Add(gvk, meta.RESTScopeNamespace)
	seed := dieappsv1.DeploymentBlank.DieMutable()
	if err := json.Unmarshal([]byte(`{
		"apiVersion": "apps/v1", "kind": "Deployment", "metadata": {"namespace": "default", "name": "web"},
		"spec": {"futureField": "from a newer API", "template": {"spec": {"containers": [{"name": "app", "image": "app:1"}]}}}
	}`), seed); err != nil {
		t.Fatal(err)
	}
	c := fake.NewClientBuilder().WithScheme(scheme).WithRESTMapper(mapper).WithObjects(seed).Build()
	ctx := context.Background()

	d := dieappsv1.DeploymentBlank.DieMutable()
	if err := c.Get(ctx, key(), d); err != nil {
		t.Fatal(err)
	}
	d.Spec().Replicas().Set(3)
	if err := c.Update(ctx, d); err != nil {
		t.Fatal(err)
	}

	var stored unstructured.Unstructured
	stored.SetGroupVersionKind(gvk)
	if err := c.Get(ctx, key(), &stored); err != nil {
		t.Fatal(err)
	}
	if got, _, _ := unstructured.NestedString(stored.Object, "spec", "futureField"); got != "from a newer API" {
		t.Errorf("futureField %q in %v", got, stored.Object["spec"])
	}
	if got, _, _ := unstructured.NestedInt64(stored.Object, "spec", "replicas"); got != 3 {
		t.Errorf("replicas %d", got)
	}
}

// What a manager applied comes back from the server as an apply
// configuration, through client-go's Extract, and is fed to a die as its
// intent: the die's apply is what was applied, whatever else has changed.
func TestApplyConfigurationRoundTrip(t *testing.T) {
	scheme := runtime.NewScheme()
	if err := dieappsv1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	gvk := schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}
	mapper := meta.NewDefaultRESTMapper(nil)
	mapper.Add(gvk, meta.RESTScopeNamespace)
	c := fake.NewClientBuilder().WithScheme(scheme).WithRESTMapper(mapper).WithReturnManagedFields().Build()
	ctx := context.Background()

	applied := dieappsv1.DeploymentBlank.DieWith(func(d *dieappsv1.DeploymentDieMut) {
		d.Metadata().Namespace().Set("default")
		d.Metadata().Name().Set("web")
		d.Spec().Replicas().Set(2)
		d.Spec().Template().Spec().Containers().At("app").Image().Set("app:1")
	})
	if err := c.Apply(ctx, applied.DieReleaseApplyConfiguration(), client.FieldOwner("litmus")); err != nil {
		t.Fatal(err)
	}
	// another manager changes something else
	if err := c.Apply(ctx, appsv1ac.Deployment("web", "default").WithLabels(map[string]string{"team": "a"}), client.FieldOwner("other")); err != nil {
		t.Fatal(err)
	}

	live := dieappsv1.DeploymentBlank.DieMutable()
	if err := c.Get(ctx, key(), live); err != nil {
		t.Fatal(err)
	}
	obj := live.DieRelease()
	owned, err := appsv1ac.ExtractDeployment(&obj, "litmus")
	if err != nil {
		t.Fatal(err)
	}
	mine := dieappsv1.DeploymentBlank.DieWith(func(d *dieappsv1.DeploymentDieMut) { d.DieFeedApplyConfiguration(owned) })
	got, err := mine.DieReleaseApplyPatchE()
	if err != nil {
		t.Fatal(err)
	}
	// exactly what Extract found, which is what was applied, but for the
	// empty status Extract adds, which claims nothing, and so isn't sent
	want, _ := json.Marshal(owned)
	var g, w map[string]any
	_ = json.Unmarshal(got, &g)
	_ = json.Unmarshal(want, &w)
	if status, ok := w["status"].(map[string]any); ok && len(status) == 0 {
		delete(w, "status")
	}
	if !reflect.DeepEqual(g, w) {
		t.Errorf("owned\n got %s\nwant %s", got, want)
	}
	if !reflect.DeepEqual(mine.DieRelease().Spec, applied.DieRelease().Spec) || mine.GetLabels() != nil {
		t.Errorf("owned spec %+v, labels %v", mine.DieRelease().Spec, mine.GetLabels())
	}
	if live.GetLabels()["team"] != "a" {
		t.Errorf("labels %v", live.GetLabels())
	}
}

// Any root die is an apply configuration for clients, with or without an
// apply configuration type: here a typed die, and an unstructured die for a
// kind with no Go type at all.
func TestApplyingAnyDie(t *testing.T) {
	scheme := runtime.NewScheme()
	if err := dieappsv1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	deployGVK := schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}
	mapper := meta.NewDefaultRESTMapper(nil)
	mapper.Add(deployGVK, meta.RESTScopeNamespace)
	mapper.Add(workloadGVK, meta.RESTScopeNamespace)
	c := fake.NewClientBuilder().WithScheme(scheme).WithRESTMapper(mapper).Build()
	ctx := context.Background()

	deploy := dieappsv1.DeploymentBlank.DieWith(func(d *dieappsv1.DeploymentDieMut) {
		d.Metadata().Namespace().Set("default")
		d.Metadata().Name().Set("web")
		d.Spec().Template().Spec().Containers().At("app").Image().Set("app:1")
	})
	ac := base.ApplyConfigurationOf(deploy)
	if err := c.Apply(ctx, ac, client.FieldOwner("litmus")); err != nil {
		t.Fatal(err)
	}
	// the object as the server returned it
	applied := dieappsv1.DeploymentBlank.DieWith(func(d *dieappsv1.DeploymentDieMut) { d.DieFeedJSON(ac.Result()) })
	if applied.GetResourceVersion() == "" || applied.Spec().Template().Spec().Containers().At("app").Image().Get() != "app:1" {
		t.Errorf("result %s", ac.Result())
	}

	workload := dieunstructured.BlankFor("acme.dev/v1", "Workload").DieWith(func(d *dieunstructured.ObjectMut) {
		d.SetNamespace("default")
		d.SetName("web")
		d.Node().Field("spec").Field("replicas").Set(2)
	})
	if err := c.Apply(ctx, base.ApplyConfigurationOf(workload), client.FieldOwner("litmus")); err != nil {
		t.Fatal(err)
	}
	var got unstructured.Unstructured
	got.SetGroupVersionKind(workloadGVK)
	if err := c.Get(ctx, key(), &got); err != nil {
		t.Fatal(err)
	}
	if replicas, _, _ := unstructured.NestedInt64(got.Object, "spec", "replicas"); replicas != 2 {
		t.Errorf("workload %v", got.Object)
	}
}
