#!/bin/bash

set -e

cd "${0%/*}/.."

go install ./src/cmd/hl2
