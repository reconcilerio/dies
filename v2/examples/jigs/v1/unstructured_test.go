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
	"encoding/json"
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	diecorev1 "reconciler.io/dies/v2/apis/core/v1"
	"reconciler.io/dies/v2/die"
	"reconciler.io/dies/v2/die/dietest"
	dieunstructured "reconciler.io/dies/v2/die/unstructured"
	jigsv1 "reconciler.io/dies/v2/examples/jigs/v1"
)

func unstructuredFromJSON(t *testing.T, s string) *unstructured.Unstructured {
	t.Helper()
	u := &unstructured.Unstructured{}
	if err := json.Unmarshal([]byte(s), &u.Object); err != nil {
		t.Fatal(err)
	}
	return u
}

func assertObject(t *testing.T, u *unstructured.Unstructured, want string) {
	t.Helper()
	got, _ := json.Marshal(u.Object)
	assertJSON(t, got, nil, want)
}

// A kind with no Go types: the jig works on the part it knows, and the rest,
// including fields its types don't know inside the parts it changes, is kept.
func TestStructuredJigOfAnUnstructuredObject(t *testing.T) {
	u := unstructuredFromJSON(t, `{
		"apiVersion": "acme.dev/v1", "kind": "Workload",
		"metadata": {"name": "web"},
		"spec": {
			"replicas": 3,
			"rollout": {"strategy": "canary"},
			"template": {"spec": {"containers": [{
				"name": "app", "image": "app:1",
				"ports": [{"containerPort": 8080}],
				"futureField": {"from": "a newer API"}
			}]}}
		}
	}`)
	ps, ok := jigsv1.PodSpecableJig.OfUnstructured(u)
	if !ok {
		t.Fatal("no jig of an unknown kind")
	}

	// reads are typed
	app := ps.Template().Spec().Containers().At("app")
	if app.Image().Get() != "app:1" || ps.Metadata().Name().Get() != "web" {
		t.Errorf("read image %q, name %q", app.Image().Get(), ps.Metadata().Name().Get())
	}
	// the port has no protocol; its default still identifies it
	dietest.NoError(t, app.Ports().At(diecorev1.ContainerPortKey{ContainerPort: 8080}).Name().Set("http"))

	dietest.NoError(t, injectSidecar(ps))

	assertObject(t, u, `{
		"apiVersion": "acme.dev/v1", "kind": "Workload",
		"metadata": {"name": "web"},
		"spec": {
			"replicas": 3,
			"rollout": {"strategy": "canary"},
			"template": {
				"metadata": {"annotations": {"sidecar.example.dev/injected": "true"}},
				"spec": {"containers": [
					{
						"name": "app", "image": "app:1",
						"ports": [{"containerPort": 8080, "name": "http"}],
						"futureField": {"from": "a newer API"}
					},
					{"name": "sidecar", "image": "proxy:1"}
				]}
			}
		}
	}`)
}

func TestUnstructuredBindingsByKind(t *testing.T) {
	u := unstructuredFromJSON(t, `{"apiVersion": "batch/v1", "kind": "CronJob", "metadata": {"name": "nightly"}}`)
	ps, _ := jigsv1.PodSpecableJig.OfUnstructured(u)
	dietest.NoError(t, injectSidecar(ps))
	containers, found, _ := unstructured.NestedSlice(u.Object, "spec", "jobTemplate", "spec", "template", "spec", "containers")
	if !found || len(containers) != 1 {
		t.Errorf("containers %v", containers)
	}
	if _, found, _ := unstructured.NestedFieldNoCopy(u.Object, "spec", "template"); found {
		t.Errorf("a CronJob uses its own bindings, not the fallback")
	}
}

// The ServiceBinding aggregate works on unstructured workloads too.
func TestAggregateOfAnUnstructuredObject(t *testing.T) {
	u := unstructuredFromJSON(t, `{
		"apiVersion": "acme.dev/v1", "kind": "Workload",
		"spec": {"template": {"spec": {
			"containers": [{"name": "app"}],
			"initContainers": [{"name": "migrate", "futureField": true}]
		}}}
	}`)
	w, _ := jigsv1.WorkloadJig.OfUnstructured(u)
	if got := w.Containers().Keys(); !reflect.DeepEqual(got, []string{"app", "migrate"}) {
		t.Errorf("containers %v", got)
	}
	for _, name := range w.Containers().Keys() {
		dietest.NoError(t, w.Containers().At(name).VolumeMounts().At("/bindings/db").Name().Set("db"))
	}
	dietest.NoError(t, w.Volumes().At("db").Secret().SecretName().Set("db-credentials"))
	assertObject(t, u, `{
		"apiVersion": "acme.dev/v1", "kind": "Workload",
		"spec": {"template": {"spec": {
			"containers": [{"name": "app", "volumeMounts": [{"mountPath": "/bindings/db", "name": "db"}]}],
			"initContainers": [{"name": "migrate", "futureField": true, "volumeMounts": [{"mountPath": "/bindings/db", "name": "db"}]}],
			"volumes": [{"name": "db", "secret": {"secretName": "db-credentials"}}]
		}}}
	}`)
	assertErrorsWith(t, "ambiguous", func() error {
		return w.Containers().At("new").Image().Set("x")
	})
}

func TestUnstructuredKeysMustBeUnique(t *testing.T) {
	u := unstructuredFromJSON(t, `{
		"apiVersion": "acme.dev/v1", "kind": "Workload",
		"spec": {"template": {"spec": {"containers": [{"name": "app"}, {"name": "app"}]}}}
	}`)
	ps, _ := jigsv1.PodSpecableJig.OfUnstructured(u)
	assertErrorsWith(t, "the list has it more than once", func() error {
		return ps.Template().Spec().Containers().At("app").Image().Set("x")
	})
}

// A typed jig over an unstructured die records what it writes as intent on
// the die, so it has exact patches, for a kind with no Go types.
func TestStructuredJigOfAnUnstructuredDie(t *testing.T) {
	d := dieunstructured.Blank.DieStamp(func(d *dieunstructured.ObjectMut) {
		dietest.NoError(t, die.FeedUnstructured(d, unstructuredFromJSON(t, `{
		"apiVersion": "acme.dev/v1", "kind": "Workload",
		"metadata": {"name": "web"},
		"spec": {
			"replicas": 3,
			"template": {"spec": {"containers": [{
				"name": "app", "image": "app:1",
				"ports": [{"containerPort": 8080}],
				"futureField": {"from": "a newer API"}
			}]}}
		}
	}`)))
	}).DieMutable()

	ps, ok := jigsv1.PodSpecableJig.OfMut(d)
	if !ok {
		t.Fatal("no jig of an unstructured die")
	}
	if got := ps.Template().Spec().Containers().At("app").Image().Get(); got != "app:1" {
		t.Errorf("image %q", got)
	}
	dietest.NoError(t, ps.Template().Spec().Containers().At("app").Ports().At(diecorev1.ContainerPortKey{ContainerPort: 8080}).Name().Set("http"))
	dietest.NoError(t, injectSidecar(ps))

	got, err := die.ReleaseApplyPatchE(d)
	assertJSON(t, got, err, `{
		"apiVersion": "acme.dev/v1", "kind": "Workload", "metadata": {"name": "web"},
		"spec": {"template": {
			"metadata": {"annotations": {"sidecar.example.dev/injected": "true"}},
			"spec": {"containers": [
				{"name": "app", "ports": [{"containerPort": 8080, "protocol": "TCP", "name": "http"}]},
				{"name": "sidecar", "image": "proxy:1"}
			]}
		}}
	}`)
	// the whole object keeps what the jig doesn't know
	u := die.ReleaseUnstructured(d)
	if future, _, _ := unstructured.NestedMap(u.Object, "spec", "template", "spec"); len(future) == 0 {
		t.Errorf("spec lost")
	}
	containers, _, _ := unstructured.NestedSlice(u.Object, "spec", "template", "spec", "containers")
	if containers[0].(map[string]any)["futureField"] == nil || len(containers) != 2 {
		t.Errorf("containers %v", containers)
	}

	// and the same jig of an immutable die reads
	ro, _ := jigsv1.PodSpecableJig.Of(dieunstructured.Blank.DieStamp(func(d *dieunstructured.ObjectMut) { dietest.NoError(t, die.FeedUnstructured(d, u)) }))
	if ro.Template().Spec().Containers().Len() != 2 {
		t.Errorf("containers %d", ro.Template().Spec().Containers().Len())
	}
}
