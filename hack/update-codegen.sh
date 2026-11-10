#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

cd $(dirname "${BASH_SOURCE[0]}")/../v2

# golangci-lint, with the linters it enables by default, as users' codebases
# run it; .golangci.yaml lints generated code too
GOLANGCI_LINT="go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0"

go generate ./...
for dir in $(find . -name "go.mod" -exec dirname {} \;); do
    echo "For $dir"
    (cd "$dir" && go mod tidy && $GOLANGCI_LINT run ./... && go test ./...)
done

