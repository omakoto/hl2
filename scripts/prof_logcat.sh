#!/bin/bash

set -e

cd "${0%/*}/.."

out=bin
mkdir -p "$out"

./scripts/build.sh
time ./bin/hl --cpuprofile prof/hl.prof "$@" -r ./samples/highlighter-logcat.toml <./samples/sample.log | wc -l
echo "web"| go tool pprof prof/hl.prof