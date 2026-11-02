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

package cluster_test

import (
	"fmt"
	"os"
	"testing"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlconfig "sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	"reconciler.io/dies/v2/conformance/internal/harness"
)

var (
	config   *rest.Config
	dies     client.Client
	upstream client.Client
)

func TestMain(m *testing.M) {
	os.Exit(run(m))
}

func run(m *testing.M) int {
	// a cluster is only ever tested when it's named
	context := os.Getenv("DIES_CLUSTER_CONTEXT")
	if context == "" {
		fmt.Fprintln(os.Stderr, "cluster conformance: skipped: set DIES_CLUSTER_CONTEXT to the kubeconfig context of a cluster to test, or run hack/conformance-cluster.sh for a kind cluster")
		return 0
	}
	cfg, err := ctrlconfig.GetConfigWithContext(context)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cluster conformance: %v\n", err)
		return 1
	}
	if dc, err := discovery.NewDiscoveryClientForConfig(cfg); err == nil {
		if v, err := dc.ServerVersion(); err == nil {
			fmt.Fprintf(os.Stderr, "cluster conformance: testing %s, Kubernetes %s\n", context, v.GitVersion)
		}
	}

	// envtest installs the CRD the suite uses, waits for it, and removes it
	// after
	env := &envtest.Environment{
		UseExistingCluster: ptr.To(true),
		Config:             cfg,
		CRDInstallOptions: envtest.CRDInstallOptions{
			CRDs:            harness.CRDs(),
			CleanUpAfterUse: true,
		},
	}
	if config, err = env.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "cluster conformance: %v\n", err)
		return 1
	}
	defer func() { _ = env.Stop() }()
	if dies, upstream, err = harness.Clients(config); err != nil {
		fmt.Fprintf(os.Stderr, "cluster conformance: %v\n", err)
		return 1
	}
	return m.Run()
}
