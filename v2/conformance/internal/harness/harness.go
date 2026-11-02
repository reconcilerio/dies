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

package harness

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"reconciler.io/dies/v2/apis"
	diecorev1 "reconciler.io/dies/v2/apis/core/v1"
	dieunstructured "reconciler.io/dies/v2/die/unstructured"
	widgetv1 "reconciler.io/dies/v2/examples/widget/v1"
)

// Clients returns a client whose scheme has the dies for every type, so that
// objects are read and written as dies, and one whose scheme has client-go's
// types, to see what the server has without dies in the way.
func Clients(config *rest.Config) (dies, upstream client.Client, err error) {
	scheme := runtime.NewScheme()
	for _, add := range []func(*runtime.Scheme) error{apis.AddToScheme, widgetv1.AddToScheme} {
		if err := add(scheme); err != nil {
			return nil, nil, err
		}
	}
	// controller-runtime sends the types client-go knows as protobuf, unless
	// told otherwise; dies are JSON
	diesConfig := rest.CopyConfig(config)
	diesConfig.ContentType = runtime.ContentTypeJSON
	if dies, err = client.New(diesConfig, client.Options{Scheme: scheme}); err != nil {
		return nil, nil, err
	}
	if upstream, err = client.New(config, client.Options{Scheme: clientgoscheme.Scheme}); err != nil {
		return nil, nil, err
	}
	return dies, upstream, nil
}

// Label marks what the suites create, to find what a run left behind.
const Label = "conformance.dies.reconciler.io/test"

// Namespace creates a namespace for a test, named for it, and deletes it
// after.
func Namespace(t testing.TB, c client.Client) string {
	t.Helper()
	b := make([]byte, 3)
	_, _ = rand.Read(b)
	suffix := "-" + hex.EncodeToString(b)
	name := "dies-" + strings.ToLower(strings.NewReplacer("/", "-", "_", "-").Replace(t.Name()))
	if len(name) > 63-len(suffix) {
		name = name[:63-len(suffix)]
	}
	name = strings.TrimRight(name, "-") + suffix
	ns := diecorev1.NamespaceBlank.DieWith(func(d *diecorev1.NamespaceDieMut) {
		d.Metadata().Name().Set(name)
		d.Metadata().Labels().Set(Label, "true")
	}).DieMutable()
	ctx := context.Background()
	if err := c.Create(ctx, ns); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = c.Delete(ctx, ns) })
	return name
}

// Molds reads a server's OpenAPI schemas, as molds for unstructured blanks.
func Molds(t testing.TB, config *rest.Config) *dieunstructured.Molds {
	t.Helper()
	dc, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	return dieunstructured.NewMolds(dc.OpenAPIV3())
}

// MoldedBlank makes an unstructured blank of a kind, in its mold, waiting
// for it, as a CRD's is published shortly after the CRD is established.
func MoldedBlank(t testing.TB, molds *dieunstructured.Molds, gvk schema.GroupVersionKind) *dieunstructured.Object {
	t.Helper()
	deadline := time.Now().Add(30 * time.Second)
	for {
		d, err := molds.BlankE(gvk)
		if err == nil {
			return d
		}
		if time.Now().After(deadline) || !strings.Contains(err.Error(), "no OpenAPI schema") {
			t.Fatal(err)
		}
		time.Sleep(200 * time.Millisecond)
	}
}
