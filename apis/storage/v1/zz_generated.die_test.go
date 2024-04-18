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

func TestCSIDriverDie_MissingMethods(t *testingx.T) {
	die := CSIDriverBlank
	ignore := []string{"TypeMeta", "ObjectMeta"}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for CSIDriverDie: %s", diff.List())
	}
}

func TestCSIDriverSpecDie_MissingMethods(t *testingx.T) {
	die := CSIDriverSpecBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for CSIDriverSpecDie: %s", diff.List())
	}
}

func TestTokenRequestDie_MissingMethods(t *testingx.T) {
	die := TokenRequestBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for TokenRequestDie: %s", diff.List())
	}
}

func TestCSINodeDie_MissingMethods(t *testingx.T) {
	die := CSINodeBlank
	ignore := []string{"TypeMeta", "ObjectMeta"}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for CSINodeDie: %s", diff.List())
	}
}

func TestCSINodeSpecDie_MissingMethods(t *testingx.T) {
	die := CSINodeSpecBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for CSINodeSpecDie: %s", diff.List())
	}
}

func TestCSINodeDriverDie_MissingMethods(t *testingx.T) {
	die := CSINodeDriverBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for CSINodeDriverDie: %s", diff.List())
	}
}

func TestVolumeNodeResourcesDie_MissingMethods(t *testingx.T) {
	die := VolumeNodeResourcesBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for VolumeNodeResourcesDie: %s", diff.List())
	}
}

func TestStorageClassDie_MissingMethods(t *testingx.T) {
	die := StorageClassBlank
	ignore := []string{"TypeMeta", "ObjectMeta"}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for StorageClassDie: %s", diff.List())
	}
}

func TestVolumeAttachmentDie_MissingMethods(t *testingx.T) {
	die := VolumeAttachmentBlank
	ignore := []string{"TypeMeta", "ObjectMeta"}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for VolumeAttachmentDie: %s", diff.List())
	}
}

func TestVolumeAttachmentSpecDie_MissingMethods(t *testingx.T) {
	die := VolumeAttachmentSpecBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for VolumeAttachmentSpecDie: %s", diff.List())
	}
}

func TestVolumeAttachmentSourceDie_MissingMethods(t *testingx.T) {
	die := VolumeAttachmentSourceBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for VolumeAttachmentSourceDie: %s", diff.List())
	}
}

func TestVolumeAttachmentStatusDie_MissingMethods(t *testingx.T) {
	die := VolumeAttachmentStatusBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for VolumeAttachmentStatusDie: %s", diff.List())
	}
}

func TestVolumeErrorDie_MissingMethods(t *testingx.T) {
	die := VolumeErrorBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for VolumeErrorDie: %s", diff.List())
	}
}