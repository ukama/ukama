#!/bin/bash
# Copyright (c) 2023-present, Ukama Inc.
# All rights reserved.

# Script to create ukama's virtual node locally 

RUN() {

	# Execute the command
	"$@"

	# Check the exit status of the command
	if [[ $? -eq 0 ]]; then
		echo -e "\033[32mCommand '$@' Ok\033[0m"
	else
		echo -e "\033[31mCommand '$@' Failed\033[0m"
		exit 1
	fi
}

if [[ -z "${VNODE_METADATA}" ]]; then
	echo "VNODE_METADATA environment variable is not set"
	exit 1
fi

if [[ -z "${VNODE_ID}" ]]; then
	echo "VNODE_ID environment variable is not set"
	exit 1
fi

# This will create the ukamaOS image and container within which the virtual
# node will be incubated
RUN make clean; make container

# Kill an existing running registry
RUN sudo docker rm --force registry

# Created local registry at port 5000
RUN sudo docker run -d -p 5000:5000 --name registry registry:latest

# start the virtual node incubator via podman
RUN podman run --network host --privileged  -it \
	-e VNODE_METADATA="$VNODE_METADATA" \
	-e VNODE_ID="$VNODE_ID" \
	-e VNODE_RUN_TARGET="local" \
	-e REPO_SERVER_URL="testing" \
	-e REPO_NAME="virtualnode" \
	localhost/testing/virtualnode:74ba00fc1-dirty

# Pull image from local registry and shut it down.
RUN podman pull --tls-verify=false \
	localhost:5000/${REPO_SERVER_URL}/${REPO_NAME}:${VNODE_ID}

RUN sudo docker rm --force registry

# Info on the newly created image.
RUN buildah info ${REPO_SERVER_URL}/${REPO_NAME}:${VNODE_ID}

echo "Done"
