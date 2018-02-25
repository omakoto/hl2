#!/bin/bash

set -e

cd "${0%/*}/.."

out=bin
mkdir -p "$out"

./scripts/build.sh
time ./bin/hl2 "$@" -r ./samples/highlighter-logcat.toml <./samples/sample.log | wc -l
