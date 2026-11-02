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

// Package cluster tests dies against a real cluster, of any distribution: that
// server-side apply and merge patches from dies do what they say, and that
// unstructured dies find the schemas of the kinds the cluster serves,
// including custom resources and an API served through an APIService.
//
// It doesn't care which APIs or fields the cluster serves, beyond those the
// tests use, so a cluster of another Kubernetes version, or a distribution
// with its own admission, is fair game; failures there are worth reporting.
//
// It only runs when DIES_CLUSTER_CONTEXT names the kubeconfig context of the
// cluster to test, run with -count=1, as Go's test cache can't see the
// cluster:
//
//	DIES_CLUSTER_CONTEXT=my-cluster go test -count=1 ./cluster/...
//
// It changes the cluster: it installs a CRD, which it
// removes after, and creates a namespace for each test, labeled
// conformance.dies.reconciler.io/test, which it deletes after. It needs to be
// able to do both, as a cluster admin can.
//
// The tests of an aggregated API need Kubernetes' sample-apiserver, which
// serves wardle.example.com behind an APIService, and are skipped without
// it. hack/conformance-cluster.sh makes a kind cluster with it installed, and
// runs the suite.
package cluster
