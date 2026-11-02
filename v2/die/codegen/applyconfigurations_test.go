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

package codegen_test

import (
	"encoding/json"
	"reflect"
	"testing"

	appsv1ac "k8s.io/client-go/applyconfigurations/apps/v1"
	corev1ac "k8s.io/client-go/applyconfigurations/core/v1"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
)

func assertSameJSON(t *testing.T, got, want any) {
	t.Helper()
	g, _ := json.Marshal(got)
	w, _ := json.Marshal(want)
	var gv, wv any
	_ = json.Unmarshal(g, &gv)
	_ = json.Unmarshal(w, &wv)
	if !reflect.DeepEqual(gv, wv) {
		t.Errorf("got %s\nwant %s", g, w)
	}
}

func desired() *appsv1ac.DeploymentApplyConfiguration {
	return appsv1ac.Deployment("web", "default").
		WithSpec(appsv1ac.DeploymentSpec().
			WithReplicas(2).
			WithTemplate(corev1ac.PodTemplateSpec().
				WithSpec(corev1ac.PodSpec().
					WithContainers(corev1ac.Container().WithName("app").WithImage("app:2")))))
}

// A die's intent is an apply configuration, and back.
func TestApplyConfigurations(t *testing.T) {
	live := dieappsv1.DeploymentBlank.DieWith(func(d *dieappsv1.DeploymentDie) {
		d.Metadata().Namespace().Set("default")
		d.Metadata().Name().Set("web")
		d.Spec().Replicas().Set(3)
		d.Spec().Template().Spec().ServiceAccountName().Set("sa")
		d.Spec().Template().Spec().Containers().At("app").Image().Set("app:1")
	}).DieCommit()

	// fed over a value, an apply configuration is the intent, exactly, and
	// the value is the base with it applied
	d := live.DieWith(func(d *dieappsv1.DeploymentDie) { d.Spec().Paused().Set(true) }).
		DieFeedApplyConfiguration(desired())
	assertSameJSON(t, d.DieReleaseApplyConfiguration(), desired())
	got, err := d.DieReleaseApplyPatchE()
	if err != nil {
		t.Fatal(err)
	}
	assertSameJSON(t, json.RawMessage(got), desired())
	if d.Spec().Replicas().Get() != 2 || d.Spec().Paused().Get() || d.Spec().Template().Spec().ServiceAccountName().Get() != "sa" || d.Spec().Template().Spec().Containers().At("app").Image().Get() != "app:2" {
		t.Errorf("value %+v", d.DieRelease().Spec)
	}
	// feeding an immutable die returns a new one
	if live.Spec().Replicas().Get() != 3 {
		t.Error("the fed die changed")
	}

	// a die within another releases and feeds its own intent
	app := d.Spec().Template().Spec().Containers().At("app")
	assertSameJSON(t, app.DieReleaseApplyConfiguration(), corev1ac.Container().WithName("app").WithImage("app:2"))
	d = d.DieWith(func(d *dieappsv1.DeploymentDie) {
		d.Spec().Template().Spec().Containers().At("app").DieFeedApplyConfiguration(corev1ac.Container().WithName("app").WithArgs("--verbose"))
	})
	assertSameJSON(t, d.Spec().Template().Spec().DieReleaseApplyConfiguration(), corev1ac.PodSpec().
		WithContainers(corev1ac.Container().WithName("app").WithArgs("--verbose")))

	// a stamp edits the intent as an apply configuration
	d = d.DieStampApplyConfiguration(func(ac *appsv1ac.DeploymentApplyConfiguration) {
		ac.Spec.WithMinReadySeconds(5)
		ac.WithLabels(map[string]string{"app": "web"})
	})
	if got := d.DieReleaseApplyConfiguration(); *got.Spec.MinReadySeconds != 5 || got.Labels["app"] != "web" || *got.Spec.Replicas != 2 {
		t.Errorf("stamped %+v", got)
	}

	// an element of a list map with no other fields is still there
	keys := dieappsv1.DeploymentBlank.DieFeedApplyConfiguration(appsv1ac.Deployment("web", "default").
		WithSpec(appsv1ac.DeploymentSpec().WithTemplate(corev1ac.PodTemplateSpec().
			WithSpec(corev1ac.PodSpec().WithContainers(corev1ac.Container().WithName("sidecar"))))))
	if !keys.Spec().Template().Spec().Containers().Has("sidecar") {
		t.Error("the sidecar isn't there")
	}
}
