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

// Package dietest has helpers for tests and fixtures that use dies.
package testing

import (
	"testing"
)

// NoError fails a test, as t.Fatal does, with the first of errs that isn't
// nil. Writes return only an error, so a test or fixture checks them inline:
//
//	dietesting.NoError(t, d.Metadata().Name().Set("web"), d.Spec().Replicas().Set(3))
//
// Each argument is evaluated before any is checked, so writes after one that
// fails are still made.
func NoError(t testing.TB, errs ...error) {
	t.Helper()
	for _, err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
}
