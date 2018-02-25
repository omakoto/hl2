#!/bin/bash

set -e

cd "${0%/*}/.."

out=bin
mkdir -p "$out"

go build -o "$out/hl2" ./src/cmd/hl2

if [[ "$1" == "-r" ]] ; then
    shift
    "$out/hl2" "$@"
fi
