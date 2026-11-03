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

package appsv1

import (
	"fmt"
	"testing"

	"reconciler.io/dies/v2/die/dietest"
	"reconciler.io/dies/v2/spike/corev1"
)

// largeDeployment has n containers with n env vars each.
func largeDeployment(n int) Deployment {
	d := serverDeployment()
	d.Spec.Template.Spec.Containers = nil
	for i := range n {
		c := corev1.Container{Name: fmt.Sprintf("c%d", i), Image: "busybox"}
		for j := range n {
			c.Env = append(c.Env, corev1.EnvVar{Name: fmt.Sprintf("E%d", j), Value: "v"})
		}
		d.Spec.Template.Spec.Containers = append(d.Spec.Template.Spec.Containers, c)
	}
	return d
}

// A write to an immutable root records intent without copying the base, so
// its cost should not depend on the size of the object.
func BenchmarkImmutableSet(b *testing.B) {
	for _, n := range []int{1, 10, 50} {
		d := DeploymentBlank.DieMutable()
		dietest.NoError(b, d.DieFeed(new(largeDeployment(n))))
		b.Run(fmt.Sprintf("size=%d", n), func(b *testing.B) {
			for b.Loop() {
				dietest.NoError(b, d.Kind().Set("Deployment"))
			}
		})
	}
}

// For contrast, the cost v1 dies pay on every stamp.
func BenchmarkDeepCopyBaseline(b *testing.B) {
	for _, n := range []int{1, 10, 50} {
		d := largeDeployment(n)
		b.Run(fmt.Sprintf("size=%d", n), func(b *testing.B) {
			for b.Loop() {
				dietest.NoError(b, DeploymentBlank.DieMutable().DieFeed(&d))
			}
		})
	}
}

// A read of a scalar walks the path; it does not materialize the object.
func BenchmarkReadLeaf(b *testing.B) {
	for _, n := range []int{1, 10, 50} {
		d := DeploymentBlank.DieStamp(func(d *DeploymentDieMut) { dietest.NoError(b, d.DieFeed(new(largeDeployment(n)))) }).DieStamp(func(d *DeploymentDieMut) {
			dietest.NoError(b, d.Spec().Template().Spec().Containers().At("c0").Image().Set("busybox:2"))
		})
		b.Run(fmt.Sprintf("size=%d", n), func(b *testing.B) {
			for b.Loop() {
				_ = d.Spec().Template().Spec().Containers().At("c0").Image().Get()
			}
		})
	}
}
