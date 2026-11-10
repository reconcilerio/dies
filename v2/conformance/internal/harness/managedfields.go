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

// Package harness has what the conformance suites share: clients for dies,
// the CRD they install, namespaces for tests, and assertions about what the
// server records.
package harness

import (
	"bytes"
	"slices"
	"strings"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/structured-merge-diff/v6/fieldpath"
)

// Owned returns the leaf fields a manager owns, as the server records them,
// such as .spec.template.spec.containers[name="app"].image.
func Owned(t testing.TB, obj metav1.Object, manager, subresource string) []string {
	t.Helper()
	var leaves []string
	found := false
	for _, mf := range obj.GetManagedFields() {
		if mf.Manager != manager || mf.Subresource != subresource || mf.FieldsV1 == nil {
			continue
		}
		found = true
		set := &fieldpath.Set{}
		if err := set.FromJSON(bytes.NewReader(mf.FieldsV1.GetRawBytes())); err != nil {
			t.Fatal(err)
		}
		set.Leaves().Iterate(func(p fieldpath.Path) {
			leaves = append(leaves, p.String())
		})
	}
	if !found {
		t.Fatalf("%s owns nothing, managers: %v", manager, managers(obj))
	}
	slices.Sort(leaves)
	return slices.Compact(leaves)
}

func managers(obj metav1.Object) []string {
	var out []string
	for _, mf := range obj.GetManagedFields() {
		out = append(out, mf.Manager+"/"+string(mf.Operation)+"/"+mf.Subresource)
	}
	return out
}

// AssertOwned asserts the leaf fields a manager owns, exactly.
func AssertOwned(t testing.TB, obj metav1.Object, manager, subresource string, want ...string) {
	t.Helper()
	got := Owned(t, obj, manager, subresource)
	slices.Sort(want)
	if !slices.Equal(got, want) {
		t.Errorf("%s owns:\n  %s\nwant:\n  %s", manager, strings.Join(got, "\n  "), strings.Join(want, "\n  "))
	}
}

// AssertOwns asserts a manager owns at least the leaf fields.
func AssertOwns(t testing.TB, obj metav1.Object, manager, subresource string, want ...string) {
	t.Helper()
	got := Owned(t, obj, manager, subresource)
	for _, w := range want {
		if !slices.Contains(got, w) {
			t.Errorf("%s doesn't own %s; it owns:\n  %s", manager, w, strings.Join(got, "\n  "))
		}
	}
}
