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
package dietest

import (
	"os"

	"reconciler.io/dies/v2/die/base"
	"reconciler.io/dies/v2/die/internal/core"
)

// FeedYAMLFile returns an immutable die of the YAML, or JSON, in a file, fed
// to a copy of d as [base.DieMutBase.DieFeedYAML] feeds it:
// FeedYAMLFile(DeploymentBlank, "testdata/web.yaml"). It panics if the file
// can't be read or doesn't decode, so that fixtures can be package variables.
func FeedYAMLFile[T any, D base.DieOf[T, D, M], M base.DieMutOf[T, D, M]](d D, name string) D {
	return core.Must(FeedYAMLFileE(d, name))
}

// FeedYAMLFileE is [FeedYAMLFile], returning an error rather than panicking.
func FeedYAMLFileE[T any, D base.DieOf[T, D, M], M base.DieMutOf[T, D, M]](d D, name string) (D, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		var zero D
		return zero, core.Errorf("dietest: %w", err)
	}
	return d.DieWithE(func(m M) error {
		_, err := m.DieFeedYAMLE(data)
		return err
	})
}
