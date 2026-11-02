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

// Package conformance tests dies against real Kubernetes API servers, in two
// suites:
//
//   - apiserver checks the dies themselves, against a control plane it runs
//     in Docker: etcd and kube-apiserver, for the Kubernetes version the dies
//     are generated from, with every API version enabled and no controllers.
//     That every die agrees with the server's schema, that patches from dies
//     do what they say, and that objects round-trip through dies exactly.
//     Without Docker it's skipped, unless CI is set.
//   - cluster checks dies with a real cluster, of any version or
//     distribution: that server-side apply and merge patches work as they
//     should, and that unstructured dies find the schemas of the kinds it
//     serves, custom resources and aggregated APIs included. It only runs when
//     DIES_CLUSTER_CONTEXT names the kubeconfig context of the cluster to test.
//     hack/conformance-cluster.sh runs it against a kind cluster.
//
// It is a module of its own, so the dies don't depend on controller-runtime,
// envtest or testcontainers.
package conformance
