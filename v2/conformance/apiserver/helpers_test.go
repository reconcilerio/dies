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

package apiserver_test

import (
	"testing"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"reconciler.io/dies/v2/conformance/internal/harness"
	dieunstructured "reconciler.io/dies/v2/die/unstructured"
)

func keyOf(namespace, name string) client.ObjectKey {
	return client.ObjectKey{Namespace: namespace, Name: name}
}

// namespace creates a namespace for a test, removed after it.
func namespace(t *testing.T) string {
	t.Helper()
	return harness.Namespace(t, dies)
}

var (
	owned       = harness.Owned
	assertOwned = harness.AssertOwned
	assertOwns  = harness.AssertOwns
)

func serverMolds(t *testing.T) *dieunstructured.Molds {
	t.Helper()
	return harness.Molds(t, config)
}

func moldedBlank(t *testing.T, molds *dieunstructured.Molds, gvk schema.GroupVersionKind) *dieunstructured.Object {
	t.Helper()
	return harness.MoldedBlank(t, molds, gvk)
}
