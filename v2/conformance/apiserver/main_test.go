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
	"errors"
	"fmt"
	"os"
	"testing"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/client-go/rest"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	"reconciler.io/dies/v2/conformance/internal/harness"
)

var (
	// config reaches the control plane as a cluster admin.
	config *rest.Config
	// dies is a client whose scheme has the dies for every type, so that
	// objects are read and written as dies.
	dies client.Client
	// upstream is a client whose scheme has client-go's types, to see what
	// the server has without dies in the way.
	upstream client.Client
)

func TestMain(m *testing.M) {
	os.Exit(run(m))
}

func run(m *testing.M) int {
	ctx := context.Background()
	hooks, err := startWebhooks()
	if err != nil {
		fmt.Fprintf(os.Stderr, "conformance: %v\n", err)
		return 1
	}
	defer hooks.stop()
	cp, err := startControlPlane(ctx, hooks.port)
	if err != nil {
		// without Docker the suite is skipped, except in CI, where it must run
		if errors.Is(err, errNoDocker) && os.Getenv("CI") == "" {
			fmt.Fprintf(os.Stderr, "conformance: skipped: %v\n", err)
			return 0
		}
		fmt.Fprintf(os.Stderr, "conformance: %v\n", err)
		return 1
	}
	defer cp.stop()

	crds, err := harness.CRDs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "conformance: %v\n", err)
		return 1
	}
	// envtest runs against the control plane in Docker as an existing
	// cluster, installing the CRDs the suite uses and waiting for them
	env := &envtest.Environment{
		UseExistingCluster: ptr.To(true),
		Config:             cp.config,
		CRDs:               crds,
	}
	config, err = env.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "conformance: %v\n", err)
		return 1
	}
	defer func() { _ = env.Stop() }()

	// the API server calls the test process to convert sprockets. The CRD is
	// installed apart from the others, without webhook options, as envtest
	// removes the conversion of a CRD whose kind isn't convertible in its
	// scheme when it has them, and an environment always does.
	sprockets, err := harness.SprocketCRD(hooks.url("/convert/sprockets"), hooks.caBundle)
	if err != nil {
		fmt.Fprintf(os.Stderr, "conformance: %v\n", err)
		return 1
	}
	if _, err := envtest.InstallCRDs(config, envtest.CRDInstallOptions{CRDs: []*apiextensionsv1.CustomResourceDefinition{sprockets}}); err != nil {
		fmt.Fprintf(os.Stderr, "conformance: %v\n", err)
		return 1
	}

	if dies, upstream, err = harness.Clients(config); err != nil {
		fmt.Fprintf(os.Stderr, "conformance: %v\n", err)
		return 1
	}
	return m.Run()
}
