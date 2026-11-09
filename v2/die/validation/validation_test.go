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

package validation_test

import (
	"testing"

	"k8s.io/apimachinery/pkg/util/validation/field"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	"reconciler.io/dies/v2/die/validation"
	dietesting "reconciler.io/dies/v2/die/testing"
)

// A map entry's handle is the location of its value, so an error about the
// value is at the key, with the value, rather than at the map.
func TestErrorAtAMapEntry(t *testing.T) {
	d := dieappsv1.DeploymentBlank.DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietesting.NoError(t, d.Metadata().Labels().Set("app", "not valid!"))
	})
	labels := d.Metadata().Labels()
	err := validation.Invalid(labels.Entry("app"), "must be a valid label value")
	want := field.Invalid(field.NewPath("metadata", "labels").Key("app"), "not valid!", "must be a valid label value")
	if err.Error() != want.Error() {
		t.Errorf("got %v\nwant %v", err, want)
	}
	if got := validation.Required(labels.Entry("tier"), "").Field; got != "metadata.labels[tier]" {
		t.Errorf("path of a missing key %q", got)
	}
}
