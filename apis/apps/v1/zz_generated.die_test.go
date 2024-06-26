//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2021-2022 the original author or authors.

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

// Code generated by diegen. DO NOT EDIT.

package v1

import (
	testing "reconciler.io/dies/testing"
	testingx "testing"
)

func TestControllerRevisionDie_MissingMethods(t *testingx.T) {
	die := ControllerRevisionBlank
	ignore := []string{"TypeMeta", "ObjectMeta"}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for ControllerRevisionDie: %s", diff.List())
	}
}

func TestDaemonSetDie_MissingMethods(t *testingx.T) {
	die := DaemonSetBlank
	ignore := []string{"TypeMeta", "ObjectMeta"}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for DaemonSetDie: %s", diff.List())
	}
}

func TestDaemonSetSpecDie_MissingMethods(t *testingx.T) {
	die := DaemonSetSpecBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for DaemonSetSpecDie: %s", diff.List())
	}
}

func TestDaemonSetUpdateStrategyDie_MissingMethods(t *testingx.T) {
	die := DaemonSetUpdateStrategyBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for DaemonSetUpdateStrategyDie: %s", diff.List())
	}
}

func TestRollingUpdateDaemonSetDie_MissingMethods(t *testingx.T) {
	die := RollingUpdateDaemonSetBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for RollingUpdateDaemonSetDie: %s", diff.List())
	}
}

func TestDaemonSetStatusDie_MissingMethods(t *testingx.T) {
	die := DaemonSetStatusBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for DaemonSetStatusDie: %s", diff.List())
	}
}

func TestDeploymentDie_MissingMethods(t *testingx.T) {
	die := DeploymentBlank
	ignore := []string{"TypeMeta", "ObjectMeta"}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for DeploymentDie: %s", diff.List())
	}
}

func TestDeploymentSpecDie_MissingMethods(t *testingx.T) {
	die := DeploymentSpecBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for DeploymentSpecDie: %s", diff.List())
	}
}

func TestDeploymentStrategyDie_MissingMethods(t *testingx.T) {
	die := DeploymentStrategyBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for DeploymentStrategyDie: %s", diff.List())
	}
}

func TestRollingUpdateDeploymentDie_MissingMethods(t *testingx.T) {
	die := RollingUpdateDeploymentBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for RollingUpdateDeploymentDie: %s", diff.List())
	}
}

func TestDeploymentStatusDie_MissingMethods(t *testingx.T) {
	die := DeploymentStatusBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for DeploymentStatusDie: %s", diff.List())
	}
}

func TestReplicaSetDie_MissingMethods(t *testingx.T) {
	die := ReplicaSetBlank
	ignore := []string{"TypeMeta", "ObjectMeta"}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for ReplicaSetDie: %s", diff.List())
	}
}

func TestReplicaSetSpecDie_MissingMethods(t *testingx.T) {
	die := ReplicaSetSpecBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for ReplicaSetSpecDie: %s", diff.List())
	}
}

func TestReplicaSetStatusDie_MissingMethods(t *testingx.T) {
	die := ReplicaSetStatusBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for ReplicaSetStatusDie: %s", diff.List())
	}
}

func TestStatefulSetDie_MissingMethods(t *testingx.T) {
	die := StatefulSetBlank
	ignore := []string{"TypeMeta", "ObjectMeta"}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for StatefulSetDie: %s", diff.List())
	}
}

func TestStatefulSetSpecDie_MissingMethods(t *testingx.T) {
	die := StatefulSetSpecBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for StatefulSetSpecDie: %s", diff.List())
	}
}

func TestStatefulSetUpdateStrategyDie_MissingMethods(t *testingx.T) {
	die := StatefulSetUpdateStrategyBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for StatefulSetUpdateStrategyDie: %s", diff.List())
	}
}

func TestRollingUpdateStatefulSetStrategyDie_MissingMethods(t *testingx.T) {
	die := RollingUpdateStatefulSetStrategyBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for RollingUpdateStatefulSetStrategyDie: %s", diff.List())
	}
}

func TestStatefulSetPersistentVolumeClaimRetentionPolicyDie_MissingMethods(t *testingx.T) {
	die := StatefulSetPersistentVolumeClaimRetentionPolicyBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for StatefulSetPersistentVolumeClaimRetentionPolicyDie: %s", diff.List())
	}
}

func TestStatefulSetOrdinalsDie_MissingMethods(t *testingx.T) {
	die := StatefulSetOrdinalsBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for StatefulSetOrdinalsDie: %s", diff.List())
	}
}

func TestStatefulSetStatusDie_MissingMethods(t *testingx.T) {
	die := StatefulSetStatusBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for StatefulSetStatusDie: %s", diff.List())
	}
}
