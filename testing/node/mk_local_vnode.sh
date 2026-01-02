#!/usr/bin/env bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

# Script to create ukama's virtual node locally
#
# Behavior:
# - BOOTSTRAP_SERVER defaults to dev.bootstrap.ukama.com (env override allowed)
# - VNODE_METADATA is hard-coded (env override allowed)
# - VNODE_ID comes from:
#     1) --node-id <id> or --node-id=<id>
#     2) env VNODE_ID
#     3) default hard-coded value

set -euo pipefail

REPO_SERVER_URL="testing"
REPO_NAME="virtualnode"
REGISTRY="localhost:5000"

BOOTSTRAP_SERVER="dev.bootstrap.ukama.com"

# Always overwrite VNODE_METADATA with known-good JSON
VNODE_METADATA=$(cat <<'JSON'
{
  "nodeInfo": {
    "type": "hnode",
    "partNumber": "",
    "skew": "",
    "mac": "",
    "swVersion": "",
    "mfgSwVersion": "",
    "assemblyDate": "2022-05-09T14:08:02.985079028-07:00",
    "oem": "",
    "mfgTestStatus": "pending",
    "status": "LabelGenerated"
  },
  "nodeConfig": [
    {
      "moduleID": "ukma-sa2219-trx-m0-e479",
      "type": "trx",
      "partNumber": "",
      "hwVersion": "",
      "mac": "",
      "swVersion": "",
      "mfgSwVersion": "",
      "mfgDate": "2022-05-09T14:08:02.985112609-07:00",
      "mfgName": "",
      "status": "AssemblyCompleted"
    }
  ]
}
JSON
)
export VNODE_METADATA

DEFAULT_VNODE_ID="uk-sa2601-tnode-v0-62f1"

while [[ $# -gt 0 ]]; do
    case "$1" in
        --node-id)
            [[ $# -ge 2 ]] || { echo "ERROR: --node-id requires a value"; exit 1; }
            VNODE_ID="$2"
            shift 2
            ;;
        --node-id=*)
            VNODE_ID="${1#*=}"
            shift
            ;;
        *)
            echo "ERROR: Unknown argument: $1"
            exit 1
            ;;
    esac
done

: "${VNODE_ID:=$DEFAULT_VNODE_ID}"

export VNODE_ID VNODE_METADATA BOOTSTRAP_SERVER

VERSION="$(make -s print-version)"
LOCAL_IMAGE="localhost/${REPO_SERVER_URL}/${REPO_NAME}:${VERSION}"
REG_IMAGE="${REGISTRY}/${REPO_SERVER_URL}/${REPO_NAME}:${VERSION}"

RUN() { "$@"; }

# Build the container image locally (podman build via Makefile)
RUN make clean
RUN make container

# Start (or restart) local registry on port 5000
RUN docker rm -f local_registry >/dev/null 2>&1 || true
RUN docker run -d -p 5000:5000 --name local_registry registry:latest

# Push the base image to the local registry
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

# Retag pulled image into local short name so inspect works
OUT_LOCAL_IMAGE="${REPO_SERVER_URL}/${REPO_NAME}:${VNODE_ID}"
RUN podman tag "${OUT_REG_IMAGE}" "${OUT_LOCAL_IMAGE}"

# Stop registry
RUN docker rm -f local_registry

# cleanup
RUN make clean

echo "Done. Image available: ${OUT_LOCAL_IMAGE}"
