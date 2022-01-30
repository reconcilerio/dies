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

package v1beta1

import (
	testing "dies.dev/testing"
	testingx "testing"
)

func TestPodSecurityPolicyDie_MissingMethods(t *testingx.T) {
	die := PodSecurityPolicyBlank
	ignore := []string{"TypeMeta", "ObjectMeta"}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for PodSecurityPolicyDie: %s", diff.List())
	}
}

func TestPodSecurityPolicySpecDie_MissingMethods(t *testingx.T) {
	die := PodSecurityPolicySpecBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for PodSecurityPolicySpecDie: %s", diff.List())
	}
}

func TestHostPortRangeDie_MissingMethods(t *testingx.T) {
	die := HostPortRangeBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for HostPortRangeDie: %s", diff.List())
	}
}

func TestSELinuxStrategyOptionsDie_MissingMethods(t *testingx.T) {
	die := SELinuxStrategyOptionsBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for SELinuxStrategyOptionsDie: %s", diff.List())
	}
}

func TestRunAsUserStrategyOptionsDie_MissingMethods(t *testingx.T) {
	die := RunAsUserStrategyOptionsBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for RunAsUserStrategyOptionsDie: %s", diff.List())
	}
}

func TestRunAsGroupStrategyOptionsDie_MissingMethods(t *testingx.T) {
	die := RunAsGroupStrategyOptionsBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for RunAsGroupStrategyOptionsDie: %s", diff.List())
	}
}

func TestSupplementalGroupsStrategyOptionsDie_MissingMethods(t *testingx.T) {
	die := SupplementalGroupsStrategyOptionsBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for SupplementalGroupsStrategyOptionsDie: %s", diff.List())
	}
}

func TestFSGroupStrategyOptionsDie_MissingMethods(t *testingx.T) {
	die := FSGroupStrategyOptionsBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for FSGroupStrategyOptionsDie: %s", diff.List())
	}
}

func TestAllowedHostPathDie_MissingMethods(t *testingx.T) {
	die := AllowedHostPathBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for AllowedHostPathDie: %s", diff.List())
	}
}

func TestAllowedFlexVolumeDie_MissingMethods(t *testingx.T) {
	die := AllowedFlexVolumeBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for AllowedFlexVolumeDie: %s", diff.List())
	}
}

func TestAllowedCSIDriverDie_MissingMethods(t *testingx.T) {
	die := AllowedCSIDriverBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for AllowedCSIDriverDie: %s", diff.List())
	}
}

func TestRuntimeClassStrategyOptionsDie_MissingMethods(t *testingx.T) {
	die := RuntimeClassStrategyOptionsBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for RuntimeClassStrategyOptionsDie: %s", diff.List())
	}
}

func TestIDRangeDie_MissingMethods(t *testingx.T) {
	die := IDRangeBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for IDRangeDie: %s", diff.List())
	}
}
