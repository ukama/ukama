#!/bin/bash -x

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

CWD=`pwd`

# Build docker image using local Dockerfile
docker build -t builder .

# Run the docker
docker run -v ${CWD}:/workspace builder \
       /bin/bash -c "/workspace/build.sh > /workspace/build.log 2>&1"

# clean up
docker image rm --force builder:latest
#sudo rm -rf ./ukama

exit 0
