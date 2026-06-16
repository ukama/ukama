#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -eu

if [ "$#" -lt 3 ]; then
    echo "usage: $0 <node-id> <container-name> <run-dir>" >&2
    exit 2
fi

NODE_ID="$1"
CONTAINER_NAME="$2"
RUN_DIR="$3"
IMAGE_REPO="${IMAGE_REPO:-testing/virtualnode}"
IMAGE="$IMAGE_REPO:$NODE_ID"

need_cmd() {
    if ! command -v "$1" >/dev/null 2>&1; then
        echo "missing required command: $1" >&2
        exit 1
    fi
}

need_cmd podman

PUBLISH_ARGS=""

if [ "${ULAB_PUBLISH_NODE_PORTS:-0}" = "1" ]; then
    PUBLISH_ARGS="-p 18001:18001 -p 18026:18026 -p 18028:18028 -p 18029:18029/udp -p 18030:18030"
fi

echo "podman: removing existing container if present: $CONTAINER_NAME"
podman rm -f "$CONTAINER_NAME" >/dev/null 2>&1 || true

echo "podman: starting $CONTAINER_NAME from $IMAGE"

if [ -n "${ULAB_NODE_ENTRYPOINT:-}" ]; then
    # shellcheck disable=SC2086
    podman run -d \
        --name "$CONTAINER_NAME" \
        --privileged \
        --device /dev/net/tun \
        --entrypoint "$ULAB_NODE_ENTRYPOINT" \
        $PUBLISH_ARGS \
        "$IMAGE" \
        ${ULAB_NODE_CMD:-}
else
    # shellcheck disable=SC2086
    podman run -d \
        --name "$CONTAINER_NAME" \
        --privileged \
        --device /dev/net/tun \
        $PUBLISH_ARGS \
        "$IMAGE"
fi

echo "node-started node=$NODE_ID container=$CONTAINER_NAME run_dir=$RUN_DIR"
