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
	"reflect"
	"strings"
	"testing"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	diecorev1 "reconciler.io/dies/v2/apis/core/v1"
	dieviews "reconciler.io/dies/v2/die/jigs"
	jigsv1 "reconciler.io/dies/v2/examples/jigs/v1"
)

func TestImageReferenceTransformIsLossless(t *testing.T) {
	sources := []string{
		"",
		"busybox",
		"busybox:1.36",
		"registry.example:5000/team/app",
		"registry.example:5000/team/app:v2",
		"registry.example/app@sha256:abc",
		"registry.example/app:v2@sha256:abc",
		// odd strings are presented whole, so they still round trip
		"app:",
		"app@",
	}
	jigs := []jigsv1.ImageReference{
		{Repository: "busybox"},
		{Repository: "registry.example:5000/app", Tag: "v3"},
		{Repository: "app", Digest: "sha256:def"},
	}
	if err := dieviews.CheckTransform(jigsv1.ImageReferenceTransform, sources, jigs); err != nil {
		t.Error(err)
	}
}

// A transform that drops information fails the checks: this one presents an
// image by its tag alone.
func TestLossyTransformIsRejected(t *testing.T) {
	tagOnly := dieviews.Transform[string, string]{
		To:   func(image string) string { _, tag, _ := strings.Cut(image, ":"); return tag },
		From: func(tag string, image string) string { repo, _, _ := strings.Cut(image, ":"); return repo + ":" + tag },
	}
	if err := dieviews.CheckTransform(tagOnly, []string{"busybox"}, nil); err == nil {
		t.Errorf("expected get-put to fail: writing back the tag of an untagged image adds a colon")
	}
}

// Jigs compose: the workload jig finds the containers, and the container
// jig presents each image in parts. Setting a tag rewrites only the image.
func TestRetagThroughJigs(t *testing.T) {
	deploy := dieappsv1.DeploymentBlank.DieWith(func(d *dieappsv1.DeploymentDie) {
		d.Metadata().Name().Set("web")
		d.Spec().Template().Spec().Containers().At("app").Image().Set("registry.example:5000/team/app:v1@sha256:abc")
		d.Spec().Template().Spec().Containers().At("proxy").Image().Set("proxy:1")
	}).DieCommit().DieImmutable(false)

	ps, _ := jigsv1.PodSpecableJig.Of(deploy)
	c, ok := jigsv1.ContainerImageJig.Of(ps.Template().Spec().Containers().At("app"))
	if !ok {
		t.Fatal("no container jig")
	}
	if got := c.Image().DieRelease(); !reflect.DeepEqual(got, jigsv1.ImageReference{Repository: "registry.example:5000/team/app", Tag: "v1", Digest: "sha256:abc"}) {
		t.Errorf("parts %+v", got)
	}
	c.Image().Tag().Set("v2")
	c.Image().Digest().Clear()

	if got := deploy.Spec().Template().Spec().Containers().At("app").Image().Get(); got != "registry.example:5000/team/app:v2" {
		t.Errorf("image %q", got)
	}
	// the change is recorded on the image field alone
	got, err := deploy.DieReleaseApplyPatchE()
	assertJSON(t, got, err, `{
		"apiVersion": "apps/v1", "kind": "Deployment", "metadata": {"name": "web"},
		"spec": {"template": {"spec": {"containers": [{"name": "app", "image": "registry.example:5000/team/app:v2"}]}}}
	}`)
}

// A jig value the source can't represent is refused, rather than written
// and read back as something else.
func TestTransformRefusesWhatTheSourceCannotHold(t *testing.T) {
	container := diecorev1.ContainerBlank.DieImmutable(false)
	c, _ := jigsv1.ContainerImageJig.Of(container)
	c.Image().Repository().Set("app")
	defer func() {
		if recover() == nil {
			t.Errorf("expected a tag containing a colon to be refused")
		}
		if got := container.Image().Get(); got != "app" {
			t.Errorf("image changed to %q", got)
		}
	}()
	// "app:a:b" reads back as repository "app:a", tag "b"
	c.Image().Tag().Set("a:b")
}
