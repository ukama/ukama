#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

set -e
set -x

TARGET=$1
UKAMA_ROOT=$2
APPS=$3

if [ "$TARGET" = "Darwin" ]; then
    TARGETPLATFORM="amd64/ubuntu:latest"
elif [ "$TARGET" = "alpine" ]; then
    TARGETPLATFORM="alpine:latest"
else
    TARGETPLATFORM="ubuntu:latest"
fi

# Check if the Docker image already exists
IMAGE_EXISTS=$(docker images -q apps-builder-${TARGETPLATFORM})

if [ -z "$IMAGE_EXISTS" ]; then
    # Build docker image using local Dockerfile
    docker build --build-arg TARGETPLATFORM=${TARGETPLATFORM} \
           -t apps-builder-${TARGETPLATFORM} .
else
    echo "Docker image 'apps-builder:latest' already exists. Skipping build."
fi

# Run the docker to build the apps 
docker run --privileged \
       -v ${UKAMA_ROOT}:/workspace \
       apps-builder-${TARGETPLATFORM} \
       /bin/bash -c "cd /workspace/builder/scripts/ && /workspace/builder/docker/apps_build.sh ${APPS} > /workspace/apps_build.log 2>&1"

# clean up
docker image rm --force apps-builder-${TARGETPLATFORM}

exit 0
