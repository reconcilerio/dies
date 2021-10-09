#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT=$(cd $(dirname "${BASH_SOURCE[0]}")/.. && pwd)

(cd "${REPO_ROOT}/cmd/diegen" ; go build)
(cd "${REPO_ROOT}" ; "${REPO_ROOT}/cmd/diegen/diegen" die:headerFile="hack/boilerplate.go.txt" paths="./apis/...")
go mod tidy

go test -cover ./apis/...
