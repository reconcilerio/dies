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
	"context"
	"testing"

	diecorev1 "reconciler.io/dies/v2/apis/core/v1"
)

func TestControlPlaneServesDies(t *testing.T) {
	ns := diecorev1.NamespaceBlank.DieImmutable(false)
	if err := dies.Get(context.Background(), keyOf("", "default"), ns); err != nil {
		t.Fatal(err)
	}
	if ns.Status().Phase().Get() != diecorev1.NamespaceActive {
		t.Errorf("phase %q", ns.Status().Phase().Get())
	}
}
