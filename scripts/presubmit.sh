#!/bin/bash

set -e

. "${0%/*}"/build.sh

gofmt -s -d $(find . -type f -name '*.go') |& perl -pe 'END{exit($. > 0 ? 1 : 0)}'

go test -v -race ./...
./tests/testall.sh

go vet ./...
staticcheck ./...
golint $(go list ./...) |& grep -v '\(exported .* should have\|comment on exported\)' | perl -pe 'END{exit($. > 0 ? 1 : 0)}'

echo "All check passed."