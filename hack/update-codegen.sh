#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT=$(cd $(dirname "${BASH_SOURCE[0]}")/.. && pwd)

(cd "${REPO_ROOT}/diegen" ; go build)

# genreate apis
(cd "${REPO_ROOT}" ; ./diegen/diegen die:headerFile="hack/boilerplate.go.txt" paths="./apis/...")
(cd "${REPO_ROOT}" ; go mod tidy)
(cd "${REPO_ROOT}" ; go test -cover ./apis/...)

# genreate diegen testdata
(cd "${REPO_ROOT}/diegen" ; rm -f testdata/zz_generated.die*)
(cd "${REPO_ROOT}/diegen" ; ./diegen die:headerFile="../hack/boilerplate.go.txt" paths="./testdata")
