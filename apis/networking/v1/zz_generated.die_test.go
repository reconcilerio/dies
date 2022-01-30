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
	testing "dies.dev/testing"
	testingx "testing"
)

func TestIngressDie_MissingMethods(t *testingx.T) {
	die := IngressBlank
	ignore := []string{"TypeMeta", "ObjectMeta"}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for IngressDie: %s", diff.List())
	}
}

func TestIngressSpecDie_MissingMethods(t *testingx.T) {
	die := IngressSpecBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for IngressSpecDie: %s", diff.List())
	}
}

func TestIngressBackendDie_MissingMethods(t *testingx.T) {
	die := IngressBackendBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for IngressBackendDie: %s", diff.List())
	}
}

func TestIngressServiceBackendDie_MissingMethods(t *testingx.T) {
	die := IngressServiceBackendBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for IngressServiceBackendDie: %s", diff.List())
	}
}

func TestServiceBackendPortDie_MissingMethods(t *testingx.T) {
	die := ServiceBackendPortBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for ServiceBackendPortDie: %s", diff.List())
	}
}

func TestIngressTLSDie_MissingMethods(t *testingx.T) {
	die := IngressTLSBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for IngressTLSDie: %s", diff.List())
	}
}

func TestIngressRuleDie_MissingMethods(t *testingx.T) {
	die := IngressRuleBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for IngressRuleDie: %s", diff.List())
	}
}

func TestHTTPIngressRuleValueDie_MissingMethods(t *testingx.T) {
	die := HTTPIngressRuleValueBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for HTTPIngressRuleValueDie: %s", diff.List())
	}
}

func TestHTTPIngressPathDie_MissingMethods(t *testingx.T) {
	die := HTTPIngressPathBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for HTTPIngressPathDie: %s", diff.List())
	}
}

func TestIngressStatusDie_MissingMethods(t *testingx.T) {
	die := IngressStatusBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for IngressStatusDie: %s", diff.List())
	}
}

func TestIngressClassDie_MissingMethods(t *testingx.T) {
	die := IngressClassBlank
	ignore := []string{"TypeMeta", "ObjectMeta"}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for IngressClassDie: %s", diff.List())
	}
}

func TestIngressClassSpecDie_MissingMethods(t *testingx.T) {
	die := IngressClassSpecBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for IngressClassSpecDie: %s", diff.List())
	}
}

func TestIngressClassParametersReferenceDie_MissingMethods(t *testingx.T) {
	die := IngressClassParametersReferenceBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for IngressClassParametersReferenceDie: %s", diff.List())
	}
}

func TestNetworkPolicyDie_MissingMethods(t *testingx.T) {
	die := NetworkPolicyBlank
	ignore := []string{"TypeMeta", "ObjectMeta"}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for NetworkPolicyDie: %s", diff.List())
	}
}

func TestNetworkPolicySpecDie_MissingMethods(t *testingx.T) {
	die := NetworkPolicySpecBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for NetworkPolicySpecDie: %s", diff.List())
	}
}

func TestNetworkPolicyIngressRuleDie_MissingMethods(t *testingx.T) {
	die := NetworkPolicyIngressRuleBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for NetworkPolicyIngressRuleDie: %s", diff.List())
	}
}

func TestNetworkPolicyEgressRuleDie_MissingMethods(t *testingx.T) {
	die := NetworkPolicyEgressRuleBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for NetworkPolicyEgressRuleDie: %s", diff.List())
	}
}

func TestNetworkPolicyPortDie_MissingMethods(t *testingx.T) {
	die := NetworkPolicyPortBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for NetworkPolicyPortDie: %s", diff.List())
	}
}

func TestNetworkPolicyPeerDie_MissingMethods(t *testingx.T) {
	die := NetworkPolicyPeerBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for NetworkPolicyPeerDie: %s", diff.List())
	}
}

func TestIPBlockDie_MissingMethods(t *testingx.T) {
	die := IPBlockBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for IPBlockDie: %s", diff.List())
	}
}
