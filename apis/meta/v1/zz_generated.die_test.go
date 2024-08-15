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

func TestConditionDie_MissingMethods(t *testingx.T) {
	die := ConditionBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for ConditionDie: %s", diff.List())
	}
}

func TestGroupResourceDie_MissingMethods(t *testingx.T) {
	die := GroupResourceBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for GroupResourceDie: %s", diff.List())
	}
}

func TestGroupVersionDie_MissingMethods(t *testingx.T) {
	die := GroupVersionBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for GroupVersionDie: %s", diff.List())
	}
}

func TestGroupVersionKindDie_MissingMethods(t *testingx.T) {
	die := GroupVersionKindBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for GroupVersionKindDie: %s", diff.List())
	}
}

func TestGroupVersionResourceDie_MissingMethods(t *testingx.T) {
	die := GroupVersionResourceBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for GroupVersionResourceDie: %s", diff.List())
	}
}

func TestGroupVersionForDiscoveryDie_MissingMethods(t *testingx.T) {
	die := GroupVersionForDiscoveryBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for GroupVersionForDiscoveryDie: %s", diff.List())
	}
}

func TestListMetaDie_MissingMethods(t *testingx.T) {
	die := ListMetaBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for ListMetaDie: %s", diff.List())
	}
}

func TestObjectMetaDie_MissingMethods(t *testingx.T) {
	die := ObjectMetaBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for ObjectMetaDie: %s", diff.List())
	}
}

func TestManagedFieldsEntryDie_MissingMethods(t *testingx.T) {
	die := ManagedFieldsEntryBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for ManagedFieldsEntryDie: %s", diff.List())
	}
}

func TestLabelSelectorDie_MissingMethods(t *testingx.T) {
	die := LabelSelectorBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for LabelSelectorDie: %s", diff.List())
	}
}

func TestLabelSelectorRequirementDie_MissingMethods(t *testingx.T) {
	die := LabelSelectorRequirementBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for LabelSelectorRequirementDie: %s", diff.List())
	}
}

func TestFieldSelectorRequirementDie_MissingMethods(t *testingx.T) {
	die := FieldSelectorRequirementBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for FieldSelectorRequirementDie: %s", diff.List())
	}
}

func TestStatusDie_MissingMethods(t *testingx.T) {
	die := StatusBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for StatusDie: %s", diff.List())
	}
}

func TestStatusDetailsDie_MissingMethods(t *testingx.T) {
	die := StatusDetailsBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for StatusDetailsDie: %s", diff.List())
	}
}

func TestStatusCauseDie_MissingMethods(t *testingx.T) {
	die := StatusCauseBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for StatusCauseDie: %s", diff.List())
	}
}
