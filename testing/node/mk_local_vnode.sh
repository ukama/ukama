#!/usr/bin/env bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

# Script to create ukama's virtual node locally

# Setup variables before executing this script:
# VNODE_ID=uk-sa2549-tnode-v0-655b
# BOOTSTRAP_SERVER=dev.bootstrap.ukama.com
# VNODE_METADATA={ "nodeInfo": { "type": "hnode", "partNumber": "", "skew": "", "mac": "", "swVersion": "", "mfgSwVersion": "", "assemblyDate": "2022-05-09T14:08:02.985079028-07:00", "oem": "", "mfgTestStatus": "pending", "status": "LabelGenerated" }, "nodeConfig": [ { "moduleID": "ukma-sa2219-trx-m0-e479", "type": "trx", "partNumber": "", "hwVersion": "", "mac": "", "swVersion": "", "mfgSwVersion": "", "mfgDate": "2022-05-09T14:08:02.985112609-07:00", "mfgName": "", "status": "AssemblyCompleted" }] }
# export VNODE_ID BOOTSTRAP_SERVER VNODE_METADATA

set -euo pipefail
REPO_SERVER_URL="testing"
REPO_NAME="virtualnode"
REGISTRY="localhost:5000"

VERSION="$(make -s print-version)"
LOCAL_IMAGE="localhost/${REPO_SERVER_URL}/${REPO_NAME}:${VERSION}"
REG_IMAGE="${REGISTRY}/${REPO_SERVER_URL}/${REPO_NAME}:${VERSION}"

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
RUN docker rm -f local_registry >/dev/null 2>&1 || true
RUN docker run -d -p 5000:5000 --name local_registry registry:latest

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
    -e BOOTSTRAP_SERVER="$BOOTSTRAP_SERVER" \
    "${REG_IMAGE}"

# Pull output image from registry
OUT_REG_IMAGE="${REGISTRY}/${REPO_SERVER_URL}/${REPO_NAME}:${VNODE_ID}"
RUN podman pull --tls-verify=false "${OUT_REG_IMAGE}"

# retag pulled image into local short name so inspect works
OUT_LOCAL_IMAGE="${REPO_SERVER_URL}/${REPO_NAME}:${VNODE_ID}"
RUN podman tag "${OUT_REG_IMAGE}" "${OUT_LOCAL_IMAGE}"

# Stop registry
RUN docker rm -f local_registry

# Inspect the newly created image (local tag)
#RUN buildah inspect "${OUT_LOCAL_IMAGE}"

# cleanup
RUN make clean

echo "Done. Image available: ${OUT_LOCAL_IMAGE}"
