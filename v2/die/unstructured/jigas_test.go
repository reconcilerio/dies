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
	"testing"

	appsv1 "k8s.io/api/apps/v1"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	diecorev1 "reconciler.io/dies/v2/apis/core/v1"
	"reconciler.io/dies/v2/die"
	"reconciler.io/dies/v2/die/dietest"
	"reconciler.io/dies/v2/die/unstructured"
)

func TestJigAs(t *testing.T) {
	d := unstructured.Blank.DieStamp(func(d *unstructured.ObjectMut) { dietest.NoError(t, die.FeedUnstructured(d, server())) })

	app := unstructured.JigAs[*diecorev1.ContainerDie](d.Node().Field("spec").Field("containers").Item("name", "app"))
	if got := app.Image().Get(); got != "app:1" {
		t.Errorf("image %q", got)
	}
	// the jig of an immutable die is immutable, and of a mutable die,
	// JigAsMut, mutable: JigAs takes only immutable dies and nodes, and
	// JigAsMut only mutable ones

	d = d.DieStamp(func(d *unstructured.ObjectMut) {
		containers := d.Node().Field("spec").Field("containers")
		app := unstructured.JigAsMut[*diecorev1.ContainerDieMut](containers.Item("name", "app"))
		dietest.NoError(t, app.Image().Set("app:2"))
		dietest.NoError(t, app.Ports().At(diecorev1.ContainerPortKey{ContainerPort: 8080}).Name().Set("http"))
		dietest.NoError(t, app.Env().At("MODE").Value().Set("full"))
		// reads see writes to the node, through either
		dietest.NoError(t, containers.Item("name", "app").Field("workingDir").Set("/srv"))
		if got := app.WorkingDir().Get(); got != "/srv" {
			t.Errorf("workingDir %q", got)
		}
		// a jig of a list element by its index
		proxy := unstructured.JigAsMut[*diecorev1.ContainerDieMut](containers.Index(1))
		dietest.NoError(t, proxy.Args().Set([]string{"--verbose"}))
	})

	// only what was written is intent, and unknown fields are kept
	got, err := die.ReleaseMergePatchE(d)
	assertJSON(t, got, err, `{"spec": {"containers": [
		{"name": "app", "image": "app:2", "futureField": {"from": "a newer API"}, "workingDir": "/srv",
		 "ports": [{"containerPort": 8080, "protocol": "TCP", "name": "http"}],
		 "env": [{"name": "MODE", "value": "full"}]},
		{"name": "proxy", "image": "proxy:1", "args": ["--verbose"]}
	]}}`)

	// a jig's copy is detached, and keeps only the fields its type knows
	copied := unstructured.JigAs[*diecorev1.ContainerDie](d.Node().Field("spec").Field("containers").Item("name", "app")).DeepCopy()
	if got := copied.Image().Get(); got != "app:2" {
		t.Errorf("image %q", got)
	}
}

// The whole of an unstructured die can be read and held as a typed die.
func TestJigAsOfAnUnstructuredObject(t *testing.T) {
	d := unstructured.BlankFor("apps/v1", "Deployment").DieStamp(func(d *unstructured.ObjectMut) {
		dietest.NoError(t, d.Node().Field("metadata").Field("name").Set("web"))
		dietest.NoError(t, d.Node().Field("spec").Field("futureField").Set(true))
	}).DieTare()

	if got := unstructured.GetAs[appsv1.Deployment](d).Name; got != "web" {
		t.Errorf("name %q", got)
	}
	d = d.DieStamp(func(d *unstructured.ObjectMut) {
		deploy := unstructured.JigAsMut[*dieappsv1.DeploymentDieMut](d)
		if got := deploy.Metadata().Name().Get(); got != "web" {
			t.Errorf("name %q", got)
		}
		dietest.NoError(t, deploy.Spec().Template().Spec().Containers().At("app").Image().Set("app:1"))
	})
	got, err := die.ReleaseMergePatchE(d)
	assertJSON(t, got, err, `{"spec": {"template": {"spec": {"containers": [{"name": "app", "image": "app:1"}]}}}}`)
	if d.Node().Field("spec").Field("futureField").Get() != true {
		t.Error("unknown field lost")
	}
}
