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

package base_test

import (
	"encoding/json"
	"strings"
	"testing"

	corev1 "k8s.io/api/core/v1"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	diecorev1 "reconciler.io/dies/v2/apis/core/v1"
	"reconciler.io/dies/v2/die"
	dietesting "reconciler.io/dies/v2/die/testing"
)

func assertDuplicate(t *testing.T, name string, err error) {
	t.Helper()
	if err == nil || !strings.Contains(err.Error(), "has more than one element") {
		t.Errorf("%s: error %v, want a duplicate key refused", name, err)
	}
}

// A list map's keys are unique in every die: content with a key twice is
// refused where it enters, so that reads, which find elements by key, can't
// fail.
func TestDuplicateKeysAreRefusedWhereTheyEnter(t *testing.T) {
	dup := &corev1.PodSpec{Containers: []corev1.Container{{Name: "app", Image: "a"}, {Name: "app", Image: "b"}}}

	pod := diecorev1.PodSpecBlank.DieMutable()
	assertDuplicate(t, "feeding a root", pod.DieFeed(dup))
	assertDuplicate(t, "feeding JSON", die.FeedJSON(pod, []byte(`{"containers": [{"name": "a"}, {"name": "a"}]}`)))
	assertDuplicate(t, "decoding", json.Unmarshal([]byte(`{"containers": [{"name": "a"}, {"name": "a"}]}`), pod))
	assertDuplicate(t, "stamping the value", die.StampValue(pod, func(r *corev1.PodSpec) { r.Containers = dup.Containers }))
	if pod.Containers().Len() != 0 {
		t.Errorf("a refused feed changed the die: %+v", pod.DieRelease())
	}

	d := dieappsv1.DeploymentBlank.DieMutable()
	assertDuplicate(t, "feeding a die within another", d.Spec().Template().Spec().DieFeed(dup))

	// an element renamed to another's key, whole or by its key field
	dietesting.NoError(t, d.Spec().Template().Spec().Containers().At("app").Image().Set("app:1"))
	dietesting.NoError(t, d.Spec().Template().Spec().Containers().At("sidecar").Image().Set("proxy:1"))
	assertDuplicate(t, "renaming by the key field", d.Spec().Template().Spec().Containers().At("sidecar").Name().Set("app"))
	assertDuplicate(t, "feeding an element", d.Spec().Template().Spec().Containers().At("sidecar").DieFeed(&corev1.Container{Name: "app"}))
	if got := d.Spec().Template().Spec().Containers().Keys(); len(got) != 2 || got[0] != "app" || got[1] != "sidecar" {
		t.Errorf("a refused write changed the die: %v", got)
	}

	// a key left to its default and the default set are the same key
	c := diecorev1.ContainerBlank.DieMutable()
	assertDuplicate(t, "a defaulted key", c.DieFeed(&corev1.Container{Ports: []corev1.ContainerPort{
		{ContainerPort: 80}, {ContainerPort: 80, Protocol: corev1.ProtocolTCP},
	}}))

	// distinct keys are fine, and reads of a die can't fail
	dietesting.NoError(t, pod.DieFeed(&corev1.PodSpec{Containers: []corev1.Container{{Name: "app"}, {Name: "sidecar"}}}))
	if !pod.Containers().Has("app") || pod.Containers().At("sidecar").Image().Get() != "" {
		t.Errorf("containers %v", pod.Containers().Keys())
	}
}
