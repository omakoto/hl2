#!/bin/bash

set -e
cd "${0%/*}/../"

name=hl2_docker

docker build $DOCKER_BUILD_OPTS -t $name .
