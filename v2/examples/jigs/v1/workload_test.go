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
	"fmt"
	"reflect"
	"strings"
	"testing"

	corev1 "k8s.io/api/core/v1"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	"reconciler.io/dies/v2/die"
	"reconciler.io/dies/v2/die/jigs"
	dietesting "reconciler.io/dies/v2/die/testing"
	dieunstructured "reconciler.io/dies/v2/die/unstructured"
	jigsv1 "reconciler.io/dies/v2/examples/jigs/v1"
)

func assertPanicsWith(t *testing.T, want string, fn func()) {
	t.Helper()
	defer func() {
		t.Helper()
		r := recover()
		if r == nil || !strings.Contains(fmt.Sprint(r), want) {
			t.Errorf("expected a panic containing %q, got %v", want, r)
		}
	}()
	fn()
}

func assertErrorsWith(t *testing.T, want string, fn func() error) {
	t.Helper()
	err := fn()
	if err == nil || !strings.Contains(fmt.Sprint(err), want) {
		t.Errorf("expected an error containing %q, got %v", want, err)
	}
}

func deployment(t testing.TB) *dieappsv1.DeploymentDieMut {
	return dieappsv1.DeploymentBlank.DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietesting.NoError(t, d.Metadata().Name().Set("web"))
		ps := d.Spec().Template().Spec()
		dietesting.NoError(t, ps.Containers().At("app").Image().Set("app:1"))
		dietesting.NoError(t, ps.InitContainers().At("migrate").Image().Set("migrate:1"))
	}).DieTare().DieMutable()
}

func TestAggregateOfContainersAndInitContainers(t *testing.T) {
	deploy := deployment(t)
	w, _ := jigsv1.WorkloadJig.MustHoldMut(deploy)

	if got := w.Containers().Keys(); !reflect.DeepEqual(got, []string{"app", "migrate"}) {
		t.Errorf("containers %v", got)
	}
	if got := w.Containers().At("migrate").Image().Get(); got != "migrate:1" {
		t.Errorf("image %q", got)
	}

	// what ServiceBinding does: changes within each container go back to
	// the list that container is in
	for _, name := range w.Containers().Keys() {
		c := w.Containers().At(name)
		dietesting.NoError(t, c.Env().At("SERVICE_BINDING_ROOT").Value().Set("/bindings"))
		dietesting.NoError(t, c.VolumeMounts().At("/bindings/db").Name().Set("db"))
	}
	dietesting.NoError(t, w.Volumes().At("db").Secret().SecretName().Set("db-credentials"))

	got, err := die.ReleaseApplyPatch(deploy)
	assertJSON(t, got, err, `{
		"apiVersion": "apps/v1", "kind": "Deployment", "metadata": {"name": "web"},
		"spec": {"template": {"spec": {
			"containers": [{"name": "app",
				"env": [{"name": "SERVICE_BINDING_ROOT", "value": "/bindings"}],
				"volumeMounts": [{"mountPath": "/bindings/db", "name": "db"}]}],
			"initContainers": [{"name": "migrate",
				"env": [{"name": "SERVICE_BINDING_ROOT", "value": "/bindings"}],
				"volumeMounts": [{"mountPath": "/bindings/db", "name": "db"}]}],
			"volumes": [{"name": "db", "secret": {"secretName": "db-credentials"}}]
		}}}
	}`)

	// adding a container could go to either list, a fact of the jig, so
	// it panics
	assertPanicsWith(t, "it could go to .spec.template.spec.containers and .spec.template.spec.initContainers", func() {
		_ = w.Containers().At("new").Image().Set("new:1")
	})
	// so could writing the whole list
	assertPanicsWith(t, "write its elements", func() {
		_ = w.Containers().Empty()
	})

	// removing goes back to the list the container is in
	dietesting.NoError(t, w.Containers().Remove("migrate"))
	if deploy.Spec().Template().Spec().InitContainers().Has("migrate") {
		t.Errorf("init container not removed")
	}
}

func TestAggregateOfAListASingleValueAndAJig(t *testing.T) {
	pipeline := jigsv1.PipelineBlank.DieStamp(func(d *jigsv1.PipelineDieMut) {
		dietesting.NoError(t, d.Steps().At("build").Image().Set("build:1"))
		dietesting.NoError(t, d.Sidecar().Name().Set("proxy"))
		dietesting.NoError(t, d.Builder().Name().Set("compile"))
	}).DieMutable()
	w, _ := jigsv1.WorkloadJig.MustHoldMut(pipeline)

	if got := w.Containers().Keys(); !reflect.DeepEqual(got, []string{"build", "proxy", "compile"}) {
		t.Errorf("containers %v", got)
	}

	// the builder is presented as a container by a jig, and changes go
	// through that jig back to the builder
	dietesting.NoError(t, w.Containers().At("compile").Env().At("GOFLAGS").Value().Set("-mod=mod"))
	dietesting.NoError(t, w.Containers().At("compile").Image().Set("golang:1.27"))
	if got := *pipeline.Builder().DieRelease(); got.Image != "golang:1.27" || got.Env[0].Name != "GOFLAGS" {
		t.Errorf("builder %+v", got)
	}
	// the sidecar is a single value
	dietesting.NoError(t, w.Containers().At("proxy").Image().Set("proxy:2"))
	if got := pipeline.Sidecar().Image().Get(); got != "proxy:2" {
		t.Errorf("sidecar image %q", got)
	}

	// the steps are the only list, so a new container is added to them
	dietesting.NoError(t, w.Containers().At("test").Image().Set("test:1"))
	if got := pipeline.Steps().Keys(); !reflect.DeepEqual(got, []string{"build", "test"}) {
		t.Errorf("steps %v", got)
	}

	// removing a single value clears it
	dietesting.NoError(t, w.Containers().Remove("proxy"))
	if pipeline.Sidecar().DieHas() {
		t.Errorf("sidecar not cleared")
	}

	// the same name in two sources is ambiguous, even to read
	dietesting.NoError(t, pipeline.Sidecar().Name().Set("build"))
	assertPanicsWith(t, "it is in .steps and .sidecar", func() {
		w.Containers().At("build").Image().Get()
	})
}

// Jigs nest: the workload jig's containers are ordinary container dies,
// so element changes can come from any code that takes one.
func TestAggregateElementsAreContainerDies(t *testing.T) {
	deploy := deployment(t)
	w, _ := jigsv1.WorkloadJig.MustHoldMut(deploy)
	c := w.Containers().At("migrate")
	dietesting.NoError(t, die.StampValue(c, func(r *corev1.Container) { r.Args = []string{"--dry-run"} }))
	if got := deploy.Spec().Template().Spec().InitContainers().At("migrate").Args().Get(); !reflect.DeepEqual(got, []string{"--dry-run"}) {
		t.Errorf("args %v", got)
	}
}

func TestAggregateOfAnAtomicList(t *testing.T) {
	batch := jigsv1.BatchBlank.DieStamp(func(d *jigsv1.BatchDieMut) {
		dietesting.NoError(t, d.DieFeed(&jigsv1.Batch{
			Containers: []corev1.Container{{Name: "extract", Image: "extract:1"}, {Name: "load", Image: "load:1"}},
			Init:       &corev1.Container{Name: "prepare"},
		}))
	}).DieMutable()
	w, _ := jigsv1.WorkloadJig.MustHoldMut(batch)

	if got := w.Containers().Keys(); !reflect.DeepEqual(got, []string{"extract", "load", "prepare"}) {
		t.Errorf("containers %v", got)
	}

	// an element of the atomic list is found by its index, and a change is
	// written back as the whole list. The list is encoded whole with
	// encoding/json, which writes empty structs such as resources; the codec
	// spike replaces that.
	dietesting.NoError(t, w.Containers().At("load").Env().At("MODE").Value().Set("full"))
	got, err := die.ReleaseApplyPatch(batch)
	assertJSON(t, got, err, `{"containers": [
		{"name": "extract", "image": "extract:1", "resources": {}},
		{"name": "load", "image": "load:1", "env": [{"name": "MODE", "value": "full"}], "resources": {}}
	]}`)

	// the atomic list is the only list, so a new container is appended
	dietesting.NoError(t, w.Containers().At("report").Image().Set("report:1"))
	dietesting.NoError(t, w.Containers().Remove("extract"))
	if got := batch.DieRelease().Containers; len(got) != 2 || got[0].Name != "load" || got[1].Name != "report" || got[1].Image != "report:1" {
		t.Errorf("containers %+v", got)
	}

	// the single value is still found by its key
	dietesting.NoError(t, w.Containers().At("prepare").Image().Set("prepare:1"))
	if got := batch.Init().Image().Get(); got != "prepare:1" {
		t.Errorf("init image %q", got)
	}
}

// Operations by key never pick the first match: a key that is not unique is
// ambiguous, whether it repeats within one source or across sources. A
// source with such a key is refused when the jig is made, so that the jig's
// reads can't fail.
func TestAggregateKeysMustBeUnique(t *testing.T) {
	// an atomic list may repeat a key, which the jig can't tell apart
	batch := jigsv1.BatchBlank.DieStamp(func(d *jigsv1.BatchDieMut) {
		dietesting.NoError(t, d.DieFeed(&jigsv1.Batch{
			Containers: []corev1.Container{{Name: "work"}, {Name: "work"}},
		}))
	}).DieMutable()
	if _, _, err := jigsv1.WorkloadJig.HoldMut(batch); err == nil || !strings.Contains(err.Error(), ".containers has more than one element [name=work]") {
		t.Errorf("error %v", err)
	}

	batch = jigsv1.BatchBlank.DieStamp(func(d *jigsv1.BatchDieMut) {
		dietesting.NoError(t, d.DieFeed(&jigsv1.Batch{
			Containers: []corev1.Container{{Name: "work"}},
			Init:       &corev1.Container{Name: "work"},
		}))
	}).DieMutable()
	if _, _, err := jigsv1.WorkloadJig.HoldMut(batch); err == nil || !strings.Contains(err.Error(), ".containers has more than one element [name=work]") {
		t.Errorf("error %v", err)
	}

	// made ambiguous through the source after the jig is made, a read
	// through the jig panics
	batch = jigsv1.BatchBlank.DieMutable()
	dietesting.NoError(t, batch.Init().Name().Set("work"))
	w, _ := jigsv1.WorkloadJig.MustHoldMut(batch)
	dietesting.NoError(t, batch.DieFeed(&jigsv1.Batch{Containers: []corev1.Container{{Name: "work"}}, Init: &corev1.Container{Name: "work"}}))
	assertPanicsWith(t, "it is in .containers and .init", func() {
		w.Containers().At("work").Image().Get()
	})
}

// A jig needn't bind every location for every source. Reading one it doesn't
// bind reads as absent, and writing one panics, so code that writes through
// a jig of any source asks first.
func TestUnboundLocations(t *testing.T) {
	batch := jigsv1.BatchBlank.DieMutable()
	w, _ := jigsv1.WorkloadJig.MustHoldMut(batch)
	if jigs.IsBound(w.Volumes()) || !jigs.IsBound(w.Containers()) {
		t.Errorf("a batch binds containers, not volumes")
	}
	if w.Volumes().Len() != 0 {
		t.Errorf("unbound volumes read %v", w.Volumes().Keys())
	}
	assertPanicsWith(t, "is not bound", func() { _ = w.Volumes().At("db").Secret().SecretName().Set("db") })

	deploy := deployment(t)
	dw, _ := jigsv1.WorkloadJig.MustHoldMut(deploy)
	if !jigs.IsBound(dw.Volumes().At("db")) || !jigs.IsBound(deploy.Spec()) {
		t.Errorf("a deployment binds volumes, and a die that isn't a jig binds everything")
	}
}

// A source the jig can't hold is refused when the jig is made, so that the
// jig's reads and writes can't fail on its content.
func TestMalformedSourcesAreRefused(t *testing.T) {
	for name, content := range map[string]string{
		"a template that isn't an object": `{"apiVersion": "apps/v1", "kind": "Deployment", "spec": {"template": "nope"}}`,
		"a path through a number":         `{"apiVersion": "batch/v1", "kind": "CronJob", "spec": {"jobTemplate": 7}}`,
		"a container repeated in a list":  `{"apiVersion": "apps/v1", "kind": "Deployment", "spec": {"template": {"spec": {"containers": [{"name": "a"}, {"name": "a"}]}}}}`,
	} {
		u := unstructuredFromJSON(t, content)
		if _, ok, err := jigsv1.PodSpecableJig.HoldUnstructured(u); !ok || err == nil {
			t.Errorf("%s: ok %v, error %v", name, ok, err)
		}
		d := dieunstructured.Blank.DieMutable()
		dietesting.NoError(t, die.FeedUnstructured(d, u))
		if _, ok, err := jigsv1.PodSpecableJig.HoldMut(d); !ok || err == nil {
			t.Errorf("%s, as a die: ok %v, error %v", name, ok, err)
		}
	}
}
