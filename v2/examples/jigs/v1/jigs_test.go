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

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	diebatchv1 "reconciler.io/dies/v2/apis/batch/v1"
	diecorev1 "reconciler.io/dies/v2/apis/core/v1"
	"reconciler.io/dies/v2/die"
	jigsv1 "reconciler.io/dies/v2/examples/jigs/v1"
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

// injectSidecar is written once against the jig, and works for every
// workload the jig knows.
func injectSidecar(ps *jigsv1.PodSpecableDie) {
	ps.Template().Spec().Containers().At("sidecar").Image().Set("proxy:1")
	ps.Template().Metadata().Annotations().Set("sidecar.example.dev/injected", "true")
}

func TestJigWritesThroughToEachSource(t *testing.T) {
	deploy := dieappsv1.DeploymentBlank.DieWith(func(d *dieappsv1.DeploymentDie) {
		d.Metadata().Name().Set("web")
		d.Spec().Template().Spec().Containers().At("app").Image().Set("app:1")
	}).DieCommit().DieImmutable(false)
	cron := diebatchv1.CronJobBlank.DieWith(func(d *diebatchv1.CronJobDie) {
		d.Metadata().Name().Set("nightly")
	}).DieCommit().DieImmutable(false)
	pod := diecorev1.PodBlank.DieWith(func(d *diecorev1.PodDie) {
		d.Metadata().Name().Set("debug")
	}).DieCommit().DieImmutable(false)

	for _, source := range []die.Die{deploy, cron, pod} {
		ps, ok := jigsv1.PodSpecableJig.Of(source)
		if !ok {
			t.Fatalf("no jig of %T", source)
		}
		injectSidecar(ps)
	}

	// each change landed at its source's own location, as intent
	got, err := deploy.DieReleaseApplyPatchE()
	assertJSON(t, got, err, `{
		"apiVersion": "apps/v1", "kind": "Deployment", "metadata": {"name": "web"},
		"spec": {"template": {
			"metadata": {"annotations": {"sidecar.example.dev/injected": "true"}},
			"spec": {"containers": [{"name": "sidecar", "image": "proxy:1"}]}
		}}
	}`)
	got, err = cron.DieReleaseApplyPatchE()
	assertJSON(t, got, err, `{
		"apiVersion": "batch/v1", "kind": "CronJob", "metadata": {"name": "nightly"},
		"spec": {"jobTemplate": {"spec": {"template": {
			"metadata": {"annotations": {"sidecar.example.dev/injected": "true"}},
			"spec": {"containers": [{"name": "sidecar", "image": "proxy:1"}]}
		}}}}
	}`)
	// a pod's template metadata is the pod's own metadata
	got, err = pod.DieReleaseApplyPatchE()
	assertJSON(t, got, err, `{
		"apiVersion": "v1", "kind": "Pod",
		"metadata": {"name": "debug", "annotations": {"sidecar.example.dev/injected": "true"}},
		"spec": {"containers": [{"name": "sidecar", "image": "proxy:1"}]}
	}`)

	// the existing container is untouched
	if got := deploy.Spec().Template().Spec().Containers().Keys(); !reflect.DeepEqual(got, []string{"app", "sidecar"}) {
		t.Errorf("containers %v", got)
	}
}

func TestJigReads(t *testing.T) {
	pod := diecorev1.PodBlank.DieWith(func(d *diecorev1.PodDie) {
		d.Metadata().Name().Set("debug")
		d.Spec().Containers().At("app").Image().Set("app:1")
	})
	ps, _ := jigsv1.PodSpecableJig.Of(pod)

	// reads through bound locations, including the overlapping metadata
	if ps.Metadata().Name().Get() != "debug" || ps.Template().Metadata().Name().Get() != "debug" {
		t.Errorf("names %q %q", ps.Metadata().Name().Get(), ps.Template().Metadata().Name().Get())
	}
	// the whole jig is assembled from the bindings
	got := ps.DieRelease()
	if got.Name != "debug" || got.Template.Name != "debug" || got.Template.Spec.Containers[0].Image != "app:1" {
		t.Errorf("released jig %+v", got)
	}
}

func TestJigOfAnImmutableDieIsReadOnly(t *testing.T) {
	ps, _ := jigsv1.PodSpecableJig.Of(dieappsv1.DeploymentBlank)
	if ps.Template().Spec().Containers().Len() != 0 {
		t.Errorf("expected no containers")
	}
	defer func() {
		if recover() == nil {
			t.Errorf("expected a write to a read-only jig to panic")
		}
	}()
	ps.Template().Spec().Containers().At("app").Image().Set("x")
}

func TestJigOfAnUnknownType(t *testing.T) {
	if _, ok := jigsv1.PodSpecableJig.Of(diecorev1.ServiceBlank); ok {
		t.Errorf("a service is not PodSpecable")
	}
	// Of only takes dies, so a value that isn't one doesn't compile
}

func TestStampThroughAJig(t *testing.T) {
	deploy := dieappsv1.DeploymentBlank.DieImmutable(false)
	ps, _ := jigsv1.PodSpecableJig.Of(deploy)
	ps.Template().DieStamp(func(r *diecorev1.PodTemplateSpec) {
		r.Labels = map[string]string{"app": "web"}
	})
	if got := deploy.Spec().Template().Metadata().Labels().At("app"); got != "web" {
		t.Errorf("label %q", got)
	}
}

// A jig holds no state: two jig locations bound to the same source
// location read the same value, so a write through one is visible through
// the other at once, including through dies and handles obtained earlier.
func TestOverlappingBindingsShareTheSource(t *testing.T) {
	pod := diecorev1.PodBlank.DieImmutable(false)
	ps, _ := jigsv1.PodSpecableJig.Of(pod)

	// obtained before any write
	templateLabels := ps.Template().Metadata().Labels()
	templateName := ps.Template().Metadata().Name()

	ps.Metadata().Labels().Set("app", "web")
	ps.Metadata().Name().Set("debug")
	if templateLabels.At("app") != "web" || templateName.Get() != "debug" {
		t.Errorf("template metadata %v %q, want the write through metadata", templateLabels.Keys(), templateName.Get())
	}

	// and the other way
	ps.Template().Metadata().Annotations().Set("note", "hi")
	if got := ps.Metadata().Annotations().At("note"); got != "hi" {
		t.Errorf("annotation %q, want the write through template metadata", got)
	}

	// a value already read is a copy, and does not change
	snapshot := ps.Metadata().DieRelease()
	ps.Metadata().Labels().Set("app", "api")
	if snapshot.Labels["app"] != "web" || ps.Template().Metadata().Labels().At("app") != "api" {
		t.Errorf("snapshot %v, live %v", snapshot.Labels, ps.Template().Metadata().Labels().Keys())
	}
}

// A copy of a jig is unlinked from its source: it reads what the jig did,
// and writes to it stay in the copy.
func TestJigCopyIsUnlinked(t *testing.T) {
	deploy := dieappsv1.DeploymentBlank.DieWith(func(d *dieappsv1.DeploymentDie) {
		d.Spec().Template().Spec().Containers().At("app").Image().Set("app:1")
	}).DieCommit().DieImmutable(false)
	ps, _ := jigsv1.PodSpecableJig.Of(deploy)

	copied := ps.DeepCopy()
	copied.Template().Spec().Containers().At("app").Image().Set("app:2")
	if got := copied.Template().Spec().Containers().At("app").Image().Get(); got != "app:2" {
		t.Errorf("copy image %q", got)
	}
	if got := deploy.Spec().Template().Spec().Containers().At("app").Image().Get(); got != "app:1" {
		t.Errorf("source image %q", got)
	}
	// the copy is a root of its own, whose intent is only what was written to it
	got, err := copied.DieReleaseApplyPatchE()
	assertJSON(t, got, err, `{"template": {"spec": {"containers": [{"name": "app", "image": "app:2"}]}}}`)
}
