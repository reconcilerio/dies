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
	"testing"

	"reconciler.io/dies/v2/die"
	"reconciler.io/dies/v2/die/dietest"
	"reconciler.io/dies/v2/spike/corev1"
)

// Locations navigated from a blank are typed and can be declared at package
// level, with no die at hand.
var (
	statusConditions = DeploymentBlank.Status().Conditions()
	specReplicas     = DeploymentBlank.Spec().Replicas()
)

func TestBlankLocations(t *testing.T) {
	d := DeploymentBlank.DieStamp(func(d *DeploymentDieMut) {
		dietest.NoError(t, d.Metadata().Name().Set("my-app"))
		dietest.NoError(t, d.Spec().Replicas().Set(2))
		dietest.NoError(t, d.Spec().Paused().Set(true))
		dietest.NoError(t, d.Status().Conditions().At("Ready").True("Available"))
	})

	got, err := die.ReleaseApplyPatchE(d, die.Only(specReplicas, statusConditions))
	assertJSON(t, got, err, `{
		"apiVersion": "apps/v1",
		"kind": "Deployment",
		"metadata": {"name": "my-app"},
		"spec": {"replicas": 2},
		"status": {
			"conditions": [{"type": "Ready", "status": "True", "reason": "Available"}]
		}
	}`)

	// a location within another type is rejected instead of matching nothing
	if _, err := die.ReleaseApplyPatchE(d, die.Only(corev1.ContainerBlank.Image())); err == nil {
		t.Error("a scope of a location within another type was accepted")
	}
}

// bind stands in for a view binding: both locations must be the same die
// type, which the compiler checks.
func bind[CD die.Located](view, source CD) [2]string {
	return [2]string{view.DiePath().String(), source.DiePath().String()}
}

func TestTypedBinding(t *testing.T) {
	got := bind(corev1.PodTemplateSpecBlank, DeploymentBlank.Spec().Template())
	if want := [2]string{".", ".spec.template"}; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	// does not compile: *PodSpecDie and *PodTemplateSpecDie differ
	// bind(corev1.PodSpecBlank, DeploymentBlank.Spec().Template())
}
