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
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	diecorev1 "reconciler.io/dies/v2/apis/core/v1"
	dieunstructured "reconciler.io/dies/v2/die/unstructured"
	jigsv1 "reconciler.io/dies/v2/examples/jigs/v1"
)

// A die within another returns its root die, sharing its state.
func TestDieRoot(t *testing.T) {
	d := dieappsv1.DeploymentBlank.DieWith(func(d *dieappsv1.DeploymentDieMut) {
		d.Metadata().Namespace().Set("default")
		d.Metadata().Name().Set("web")
	}).DieMutable()
	spec := d.Spec().Template().Spec()

	root, ok := spec.DieRoot().(*dieappsv1.DeploymentDieMut)
	if !ok || root.GetName() != "web" {
		t.Fatalf("root %T", spec.DieRoot())
	}
	// it shares the tree's state
	root.Spec().Replicas().Set(3)
	if d.Spec().Replicas().Get() != 3 {
		t.Errorf("a write through the root didn't reach the die")
	}
	// an object root is a metav1.Object, whatever its type
	if namespace := spec.DieRoot().(metav1.Object).GetNamespace(); namespace != "default" {
		t.Errorf("namespace %q", namespace)
	}
	// handles, and a root itself
	if _, ok := spec.Containers().At("app").Image().DieRoot().(*dieappsv1.DeploymentDieMut); !ok {
		t.Errorf("a handle's root is %T", spec.Containers().At("app").Image().DieRoot())
	}
	if _, ok := d.DieRoot().(*dieappsv1.DeploymentDieMut); !ok {
		t.Errorf("a root's root is %T", d.DieRoot())
	}
	// a copy of a die within another is unlinked: its own root
	copied := spec.DeepCopy()
	if _, ok := copied.Containers().DieRoot().(*diecorev1.PodSpecDieMut); !ok {
		t.Errorf("a copy's root is %T", copied.Containers().DieRoot())
	}
	if _, ok := spec.DieRoot().(*diecorev1.PodSpecDieMut); ok {
		t.Errorf("the root is a %T", spec.DieRoot())
	}
}

// An unstructured node's root is its object, and a jig's dies are within the
// jig.
func TestDieRootOfNodesAndJigs(t *testing.T) {
	u := dieunstructured.BlankFor("acme.dev/v1", "Workload").DieMutable()
	if _, ok := u.Node().Field("spec").Field("replicas").DieRoot().(*dieunstructured.ObjectMut); !ok {
		t.Errorf("a node's root is %T", u.Node().Field("spec").DieRoot())
	}

	d := dieappsv1.DeploymentBlank.DieMutable()
	ps, ok := jigsv1.PodSpecableJig.OfMut(d)
	if !ok {
		t.Fatal("no jig")
	}
	if _, ok := ps.Template().Spec().DieRoot().(*jigsv1.PodSpecableDieMut); !ok {
		t.Errorf("a jig's die's root is %T", ps.Template().Spec().DieRoot())
	}
}
