#!/bin/bash
# Copyright (c) 2023-present, Ukama Inc.
# All rights reserved.

# Script to create ukama's virtual node locally 

set -x

# This will create the ukamaOS image and container within which the virtual
# node will be incubated
make clean; make container

# Kill an existing running registry
sudo docker rm --force registry

# Created local registry at port 5000
sudo docker run -d -p 5000:5000 --name registry registry:latest

# start the virtual node incubator via podman
podman run --network host --privileged  -it \
	   -e VNODE_METADATA="$VNODE_METADATA" \
	   -e VNODE_ID="$VNODE_ID" \
	   -e VNODE_RUN_TARGET="local" \
	   -e REPO_SERVER_URL="testing" \
	   -e REPO_NAME="virtualnode" \
	   localhost/testing/virtualnode:74ba00fc1-dirty

# Info on the newly created image.
buildah info ${REPO_SERVER_URL}/${REPO_NAME}:${VNODE_ID}

echo "Done"
