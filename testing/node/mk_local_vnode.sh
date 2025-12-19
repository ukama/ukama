#!/usr/bin/env bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

# Script to create ukama's virtual node locally

set -xeuo pipefail
REPO_SERVER_URL="testing"
REPO_NAME="virtualnode"

# Always use the Makefile's computed version (same as the image tag)
VERSION="$(make -s print-version)"
LOCAL_IMAGE="$(make -s print-local-image)"

REGISTRY="localhost:5000"
REG_IMAGE="${REGISTRY}/${LOCAL_IMAGE}"

RUN() {
    "$@"
}

if [[ -z "${VNODE_METADATA:-}" ]]; then
    echo "VNODE_METADATA environment variable is not set"
    exit 1
fi

if [[ -z "${VNODE_ID:-}" ]]; then
    echo "VNODE_ID environment variable is not set"
    exit 1
fi

# Build the container image locally (podman build via Makefile)
RUN make clean
RUN make container

# Start (or restart) local registry on port 5000
sudo docker rm -f local_registry >/dev/null 2>&1 || true
RUN sudo docker run -d -p 5000:5000 --name local_registry registry:latest

# ---- Key fix: push the base image to the local registry ----
RUN podman tag "${LOCAL_IMAGE}" "${REG_IMAGE}"
RUN podman push --tls-verify=false "${REG_IMAGE}"

# Run incubator from the registry ref (consistent path)
RUN podman run --network host --privileged -it \
    -v "${PWD}/ukama_${VERSION}.tgz:/ukama/ukama.tgz:ro" \
    --tls-verify=false \
    -e VNODE_METADATA="$VNODE_METADATA" \
    -e VNODE_ID="$VNODE_ID" \
    -e VNODE_RUN_TARGET="local" \
    -e REPO_SERVER_URL="$REPO_SERVER_URL" \
    -e REPO_NAME="$REPO_NAME" \
    "${REG_IMAGE}"

# Pull output image from registry
OUT_REG_IMAGE="${REGISTRY}/${REPO_SERVER_URL}/${REPO_NAME}:${VNODE_ID}"
RUN podman pull --tls-verify=false "${OUT_REG_IMAGE}"

# retag pulled image into local short name so inspect works
OUT_LOCAL_IMAGE="${REPO_SERVER_URL}/${REPO_NAME}:${VNODE_ID}"
RUN podman tag "${OUT_REG_IMAGE}" "${OUT_LOCAL_IMAGE}"

# Stop registry
RUN sudo docker rm -f local_registry

# Inspect the newly created image (local tag)
RUN buildah inspect "${OUT_LOCAL_IMAGE}"

# cleanup
RUN make clean

echo "Done"

