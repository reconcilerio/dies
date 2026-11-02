#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

cd $(dirname "${BASH_SOURCE[0]}")/../v2


go generate ./...
for dir in $(find . -name "go.mod" -exec dirname {} \;); do
    echo "For $dir"
    (cd "$dir" && go mod tidy && go test ./...)
done

