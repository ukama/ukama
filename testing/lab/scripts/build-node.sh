#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -eu

if [ "$#" -lt 2 ]; then
    echo "usage: $0 <repo> <node-id> [runtime]" >&2
    exit 2
fi

REPO="$1"
NODE_ID="$2"
NODE_RUNTIME="${3:-${NODE_RUNTIME:-starter}}"

IMAGE_REPO="${IMAGE_REPO:-testing/virtualnode}"
BASE_IMAGE_REPO="${BASE_IMAGE_REPO:-testing/virtualnode-base}"
IMAGE="$IMAGE_REPO:$NODE_ID"

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
NODE_DIR="$REPO/testing/node"
MK_LOCAL_VNODE="$NODE_DIR/mk_local_vnode.sh"

need_cmd() {
    if ! command -v "$1" >/dev/null 2>&1; then
        echo "missing required command: $1" >&2
        exit 1
    fi
}

node_type() {
    case "$1" in
        *-tnode-*) echo "tnode" ;;
        *-cnode-*) echo "cnode" ;;
        *-anode-*) echo "anode" ;;
        *)
            echo "cannot determine node type from node id: $1" >&2
            exit 1
            ;;
    esac
}

need_cmd podman

if [ ! -x "$MK_LOCAL_VNODE" ]; then
    echo "missing executable $MK_LOCAL_VNODE" >&2
    exit 1
fi

if [ ! -x "$SCRIPT_DIR/stamp-node-image.sh" ]; then
    echo "missing executable $SCRIPT_DIR/stamp-node-image.sh" >&2
    exit 1
fi

NODE_TYPE="$(node_type "$NODE_ID")"
BASE_IMAGE="$BASE_IMAGE_REPO:${NODE_TYPE}-${NODE_RUNTIME}"

echo "build-node: $NODE_ID type=$NODE_TYPE runtime=$NODE_RUNTIME"

if podman image exists "$IMAGE"; then
    echo "build-node: image exists $IMAGE"

    if ! podman image exists "$BASE_IMAGE"; then
        echo "build-node: seed base image $BASE_IMAGE from existing $IMAGE"
        podman tag "$IMAGE" "$BASE_IMAGE"
    fi

    exit 0
fi

if podman image exists "$BASE_IMAGE"; then
    echo "build-node: stamp $IMAGE from $BASE_IMAGE"
    "$SCRIPT_DIR/stamp-node-image.sh" \
        "$BASE_IMAGE" \
        "$NODE_ID" \
        "$IMAGE" \
        "$NODE_RUNTIME"
    exit 0
fi

echo "build-node: no base image for $NODE_TYPE/$NODE_RUNTIME; full build $IMAGE"
(
    cd "$NODE_DIR"
    ./mk_local_vnode.sh --node-id "$NODE_ID" --runtime "$NODE_RUNTIME"
)

if ! podman image exists "$IMAGE"; then
    echo "build-node: expected image not found after build: $IMAGE" >&2
    exit 1
fi

echo "build-node: cache base image $BASE_IMAGE from $IMAGE"
podman tag "$IMAGE" "$BASE_IMAGE"
