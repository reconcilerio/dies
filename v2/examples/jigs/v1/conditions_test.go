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
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	"reconciler.io/dies/v2/die"
	"reconciler.io/dies/v2/die/dietest"
	"reconciler.io/dies/v2/die/jigs"
	jigsv1 "reconciler.io/dies/v2/examples/jigs/v1"
)

var updated = metav1.Unix(1700000000, 0)

func legacyConditions() []appsv1.DeploymentCondition {
	return []appsv1.DeploymentCondition{
		{Type: appsv1.DeploymentAvailable, Status: corev1.ConditionFalse, Reason: "Starting", LastUpdateTime: updated},
		{Type: appsv1.DeploymentProgressing, Status: corev1.ConditionTrue, Reason: "NewReplicaSetAvailable", LastUpdateTime: updated},
	}
}

func TestDeploymentConditionsAreLossless(t *testing.T) {
	err := jigs.CheckTransform(jigsv1.DeploymentConditions,
		[][]appsv1.DeploymentCondition{nil, {}, legacyConditions()},
		[][]metav1.Condition{nil, {{Type: "Available", Status: metav1.ConditionTrue, Reason: "Ready"}}},
	)
	if err != nil {
		t.Error(err)
	}
}

func TestHarmonizedConditionsKeepWhatTheJigDoesNotShow(t *testing.T) {
	deploy := dieappsv1.DeploymentBlank.DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietest.NoError(t, d.Metadata().Name().Set("web"))
		dietest.NoError(t, die.StampValue(d.Status(), func(r *appsv1.DeploymentStatus) { r.Conditions = legacyConditions() }))
	}).DieTare().DieMutable()

	c, _ := jigsv1.ConditionedJig.OfMut(deploy.Status())
	dietest.NoError(t, c.Conditions().At("Available").Status().Set(metav1.ConditionTrue))
	dietest.NoError(t, c.Conditions().At("Available").Reason().Set("MinimumReplicasAvailable"))

	got := deploy.Status().Conditions().At(appsv1.DeploymentAvailable)
	if got.Status().Get() != corev1.ConditionTrue || got.Reason().Get() != "MinimumReplicasAvailable" {
		t.Errorf("condition %+v", *got.DieRelease())
	}
	// lastUpdateTime isn't in the jig, and is kept
	if lut := got.LastUpdateTime().Get(); !lut.Equal(&updated) {
		t.Errorf("lastUpdateTime %v, want it kept", got.LastUpdateTime().Get())
	}
	// only the changed fields of the changed condition are intent
	patch, err := die.ReleaseApplyPatchE(deploy, die.Subresource("status"))
	assertJSON(t, patch, err, `{
		"apiVersion": "apps/v1", "kind": "Deployment", "metadata": {"name": "web"},
		"status": {"conditions": [{"type": "Available", "status": "True", "reason": "MinimumReplicasAvailable"}]}
	}`)

	// a deployment condition has nowhere to keep an observed generation
	assertErrorsWith(t, "the transform for .conditions cannot represent", func() error {
		return c.Conditions().At("Available").ObservedGeneration().Set(3)
	})
}
