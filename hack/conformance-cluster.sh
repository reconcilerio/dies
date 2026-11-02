#!/usr/bin/env bash

# Runs the cluster conformance suite against a kind cluster, dies-conformance,
# with Kubernetes' sample-apiserver installed behind an APIService, creating
# the cluster if it doesn't exist. The cluster is deleted after, unless
# --keep is passed. Arguments after -- are passed to go test. The cluster's
# kubeconfig is kept apart from the user's, which is left as it is.
#
# To run the suite against another cluster instead, nothing is created:
#
#   DIES_CLUSTER_CONTEXT=<kubeconfig context> go test -count=1 ./cluster/...
#
# from v2/conformance.

set -o errexit
set -o nounset
set -o pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."

name=dies-conformance
keep=false
while [[ $# -gt 0 ]]; do
  case "$1" in
    --keep) keep=true; shift ;;
    --) shift; break ;;
    *) echo "unknown argument: $1" >&2; exit 2 ;;
  esac
done

# a kubeconfig of its own, so the user's, and its current context, are left
# alone
export KUBECONFIG
KUBECONFIG="$(mktemp -t dies-conformance-kubeconfig.XXXXXX)"
cleanup() {
  rm -f "${KUBECONFIG}"
}
trap cleanup EXIT

if kind get clusters | grep -qx "${name}"; then
  kind export kubeconfig --name "${name}" --kubeconfig "${KUBECONFIG}"
else
  kind create cluster --name "${name}" --kubeconfig "${KUBECONFIG}" --wait 2m
  if [[ "${keep}" == false ]]; then
    trap 'kind delete cluster --name "${name}" --kubeconfig "${KUBECONFIG}"; cleanup' EXIT
  else
    echo "keeping the cluster; for its kubeconfig: kind export kubeconfig --name ${name}"
  fi
fi
context="kind-${name}"

kubectl --context "${context}" apply -f hack/conformance/sample-apiserver.yaml
kubectl --context "${context}" -n wardle rollout status deployment/sample-apiserver --timeout 3m
kubectl --context "${context}" wait apiservice/v1alpha1.wardle.example.com --for condition=Available --timeout 2m

cd v2/conformance
# never cached: the result is the cluster's
DIES_CLUSTER_CONTEXT="${context}" go test -count=1 ./cluster/... "$@"
