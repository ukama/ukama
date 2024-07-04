#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

CWD=$(pwd)
BRANCH=$1

set -e

# Determine the platform
PLATFORM=$(uname -s)

if [ "$PLATFORM" = "Darwin" ]; then
    TARGETPLATFORM="amd64/ubuntu:latest"
else
    TARGETPLATFORM="ubuntu:latest"
fi

# Check if the Docker image already exists
IMAGE_EXISTS=$(docker images -q builder:latest)

if [ -z "$IMAGE_EXISTS" ]; then
    # Build docker image using local Dockerfile
    docker build --build-arg TARGETPLATFORM=${TARGETPLATFORM} -t builder .
else
    echo "Docker image 'builder:latest' already exists. Skipping build."
fi

# Run the docker to build the UkamaOS
docker run --privileged \
       -v ${CWD}:/workspace \
       -v /dev:/dev \
       builder:latest \
       /bin/bash -c "/workspace/build.sh ${BRANCH} > /workspace/build.log 2>&1"

# clean up
docker image rm --force builder:latest
sudo rm -rf ./ukama

exit 0
