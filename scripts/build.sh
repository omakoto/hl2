#!/bin/bash

set -e

cd "${0%/*}/.."

out=bin
mkdir -p "$out"

go build -o "$out/hl" ./src/cmd/hl

if [[ "$1" == "-r" ]] ; then
    shift
    "$out/hl" "$@"
fi
