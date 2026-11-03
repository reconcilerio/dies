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
	dietesting "reconciler.io/dies/v2/die/testing"
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
func injectSidecar(ps *jigsv1.PodSpecableDieMut) error {
	if err := ps.Template().Spec().Containers().At("sidecar").Image().Set("proxy:1"); err != nil {
		return err
	}
	return ps.Template().Metadata().Annotations().Set("sidecar.example.dev/injected", "true")
}

func TestJigWritesThroughToEachSource(t *testing.T) {
	deploy := dieappsv1.DeploymentBlank.DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietesting.NoError(t, d.Metadata().Name().Set("web"))
		dietesting.NoError(t, d.Spec().Template().Spec().Containers().At("app").Image().Set("app:1"))
	}).DieTare().DieMutable()
	cron := diebatchv1.CronJobBlank.DieStamp(func(d *diebatchv1.CronJobDieMut) {
		dietesting.NoError(t, d.Metadata().Name().Set("nightly"))
	}).DieTare().DieMutable()
	pod := diecorev1.PodBlank.DieStamp(func(d *diecorev1.PodDieMut) {
		dietesting.NoError(t, d.Metadata().Name().Set("debug"))
	}).DieTare().DieMutable()

	for _, source := range []die.MutableDie{deploy, cron, pod} {
		ps, ok := jigsv1.PodSpecableJig.MustHoldMut(source)
		if !ok {
			t.Fatalf("no jig of %T", source)
		}
		dietesting.NoError(t, injectSidecar(ps))
	}

	// each change landed at its source's own location, as intent
	got, err := die.ReleaseApplyPatch(deploy)
	assertJSON(t, got, err, `{
		"apiVersion": "apps/v1", "kind": "Deployment", "metadata": {"name": "web"},
		"spec": {"template": {
			"metadata": {"annotations": {"sidecar.example.dev/injected": "true"}},
			"spec": {"containers": [{"name": "sidecar", "image": "proxy:1"}]}
		}}
	}`)
	got, err = die.ReleaseApplyPatch(cron)
	assertJSON(t, got, err, `{
		"apiVersion": "batch/v1", "kind": "CronJob", "metadata": {"name": "nightly"},
		"spec": {"jobTemplate": {"spec": {"template": {
			"metadata": {"annotations": {"sidecar.example.dev/injected": "true"}},
			"spec": {"containers": [{"name": "sidecar", "image": "proxy:1"}]}
		}}}}
	}`)
	// a pod's template metadata is the pod's own metadata
	got, err = die.ReleaseApplyPatch(pod)
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
	pod := diecorev1.PodBlank.DieStamp(func(d *diecorev1.PodDieMut) {
		dietesting.NoError(t, d.Metadata().Name().Set("debug"))
		dietesting.NoError(t, d.Spec().Containers().At("app").Image().Set("app:1"))
	})
	ps, _ := jigsv1.PodSpecableJig.MustHold(pod)

	// reads through bound locations, including the overlapping metadata
	if ps.Metadata().Name().Get() != "debug" || ps.Template().Metadata().Name().Get() != "debug" {
		t.Errorf("names %q %q", ps.Metadata().Name().Get(), ps.Template().Metadata().Name().Get())
	}
	// the whole jig is assembled from the bindings
	got := *ps.DieRelease()
	if got.Name != "debug" || got.Template.Name != "debug" || got.Template.Spec.Containers[0].Image != "app:1" {
		t.Errorf("released jig %+v", got)
	}
}

// The jig of an immutable die is immutable, and of a mutable die mutable:
// a read-only jig of a mutable die would see its writes.
func TestJigOfAnImmutableDieIsImmutable(t *testing.T) {
	ps, _ := jigsv1.PodSpecableJig.MustHold(dieappsv1.DeploymentBlank)
	if ps.Template().Spec().Containers().Len() != 0 {
		t.Errorf("expected no containers")
	}
	// DieStamp on the jig works on a copy, which is not linked to the source
	written := ps.DieStamp(func(ps *jigsv1.PodSpecableDieMut) {
		dietesting.NoError(t, ps.Template().Spec().Containers().At("app").Image().Set("x"))
	})
	if written.Template().Spec().Containers().At("app").Image().Get() != "x" || dieappsv1.DeploymentBlank.Spec().Template().Spec().Containers().Len() != 0 {
		t.Errorf("written %s", die.MustReleaseJSON(written))
	}
}

func TestJigOfAnUnknownType(t *testing.T) {
	if _, ok := jigsv1.PodSpecableJig.MustHold(diecorev1.ServiceBlank); ok {
		t.Errorf("a service is not PodSpecable")
	}
	// Of only takes dies, so a value that isn't one doesn't compile
}

func TestStampThroughAJig(t *testing.T) {
	deploy := dieappsv1.DeploymentBlank.DieMutable()
	ps, _ := jigsv1.PodSpecableJig.MustHoldMut(deploy)
	dietesting.NoError(t, die.StampValue(ps.Template(), func(r *diecorev1.PodTemplateSpec) {
		r.Labels = map[string]string{"app": "web"}
	}))
	if got := deploy.Spec().Template().Metadata().Labels().At("app"); got != "web" {
		t.Errorf("label %q", got)
	}
}

// A jig holds no state: two jig locations bound to the same source
// location read the same value, so a write through one is visible through
// the other at once, including through dies and handles obtained earlier.
func TestOverlappingBindingsShareTheSource(t *testing.T) {
	pod := diecorev1.PodBlank.DieMutable()
	ps, _ := jigsv1.PodSpecableJig.MustHoldMut(pod)

	// obtained before any write
	templateLabels := ps.Template().Metadata().Labels()
	templateName := ps.Template().Metadata().Name()

	dietesting.NoError(t, ps.Metadata().Labels().Set("app", "web"))
	dietesting.NoError(t, ps.Metadata().Name().Set("debug"))
	if templateLabels.At("app") != "web" || templateName.Get() != "debug" {
		t.Errorf("template metadata %v %q, want the write through metadata", templateLabels.Keys(), templateName.Get())
	}

	// and the other way
	dietesting.NoError(t, ps.Template().Metadata().Annotations().Set("note", "hi"))
	if got := ps.Metadata().Annotations().At("note"); got != "hi" {
		t.Errorf("annotation %q, want the write through template metadata", got)
	}

	// a value already read is a copy, and does not change
	snapshot := *ps.Metadata().DieRelease()
	dietesting.NoError(t, ps.Metadata().Labels().Set("app", "api"))
	if snapshot.Labels["app"] != "web" || ps.Template().Metadata().Labels().At("app") != "api" {
		t.Errorf("snapshot %v, live %v", snapshot.Labels, ps.Template().Metadata().Labels().Keys())
	}
}

// A copy of a jig is unlinked from its source: it reads what the jig did,
// and writes to it stay in the copy.
func TestJigCopyIsUnlinked(t *testing.T) {
	deploy := dieappsv1.DeploymentBlank.DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietesting.NoError(t, d.Spec().Template().Spec().Containers().At("app").Image().Set("app:1"))
	}).DieTare().DieMutable()
	ps, _ := jigsv1.PodSpecableJig.MustHoldMut(deploy)

	copied := ps.DeepCopy()
	dietesting.NoError(t, copied.Template().Spec().Containers().At("app").Image().Set("app:2"))
	if got := copied.Template().Spec().Containers().At("app").Image().Get(); got != "app:2" {
		t.Errorf("copy image %q", got)
	}
	if got := deploy.Spec().Template().Spec().Containers().At("app").Image().Get(); got != "app:1" {
		t.Errorf("source image %q", got)
	}
	// the copy is a root of its own, whose intent is only what was written to it
	got, err := die.ReleaseApplyPatch(copied)
	assertJSON(t, got, err, `{"template": {"spec": {"containers": [{"name": "app", "image": "app:2"}]}}}`)
}
