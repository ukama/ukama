#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

set -e

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
    TARGETPLATFORM="arm64v8/ubuntu:20.04"
    DOCKER_IMG="apps-builder-arm64"
else
    TARGETPLATFORM="ubuntu:latest"
    DOCKER_IMG="apps-builder-ubuntu"
fi

# Build the container
docker build --build-arg TARGETPLATRFORM=${TARGETPLATFORM} -t ${DOCKER_IMG} .

# Run the docker to build the apps 
docker run --privileged \
       -v ${UKAMA_ROOT}:/workspace \
       ${DOCKER_IMG} \
       /bin/bash -c "cd /workspace/builder/scripts/ && /workspace/builder/docker/apps_build.sh ${NODE} ${APPS} > /workspace/apps_build.log 2>&1"

# clean up
docker image rm --force apps-builder-${TARGETPLATFORM}

exit 0
