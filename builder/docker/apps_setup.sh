#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

set -e
set -x 

NODE=$1
TARGET=$2
UKAMA_ROOT=$3
APPS=$4
DOCKER_IMG=

if [ "$TARGET" = "Darwin" ]; then
    TARGETPLATFORM="amd64/ubuntu:latest"
    DOCKER_IMG="apps-builder-amd64"
elif [ "$TARGET" = "alpine" ]; then
    TARGETPLATFORM="alpine:latest"
    DOCKER_IMG="apps-builder-alpine"
elif [ "$TARGET" = "arm64" ]; then
    TARGETPLATFORM="ubuntu:latest"
    DOCKER_IMG="apps-builder-arm64"
else
    TARGETPLATFORM="ubuntu:latest"
    DOCKER_IMG="apps-builder-ubuntu"
fi

# Build the container
if [ "$TARGET" = "arm64" ]; then
    export DOCKER_BUILDKIT=1
    docker buildx build --platform linux/arm64 -t ${DOCKER_IMG} \
           --load --build-arg BUILDKIT_CPU_LIMIT=$(nproc) .
else
    docker build --build-arg TARGETPLATFORM=${TARGETPLATFORM} \
           -t ${DOCKER_IMG} .
fi

# Run the docker to build the apps
if [ "$TARGET" = "arm64" ]; then
    docker run --platform linux/arm64 --privileged \
           -v ${UKAMA_ROOT}:/workspace \
           ${DOCKER_IMG} \
           /bin/bash -c "cd /workspace/builder/scripts/ && /workspace/builder/docker/apps_build.sh ${NODE} ${APPS} > /workspace/apps_build.log 2>&1"
else
    docker run --privileged \
           -v ${UKAMA_ROOT}:/workspace \
           ${DOCKER_IMG} \
           /bin/bash -c "cd /workspace/builder/scripts/ && /workspace/builder/docker/apps_build.sh ${NODE} ${APPS} > /workspace/apps_build.log 2>&1"
fi

# clean up
docker image rm --force "${DOCKER_IMG}"

exit 0
