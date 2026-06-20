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
NET_STATE="$RUN_DIR/runtime-net/net.env"
LAB_NET="${LAB_NET:-}"

need_cmd() {
    if ! command -v "$1" >/dev/null 2>&1; then
        echo "missing required command: $1" >&2
        exit 1
    fi
}

container_ip_on_network() {
    container="$1"
    network="$2"

    podman inspect -f '{{range $name, $net := .NetworkSettings.Networks}}{{if eq $name "'"$network"'"}}{{$net.IPAddress}}{{end}}{{end}}' \
        "$container" 2>/dev/null
}

need_cmd podman

if [ -z "$LAB_NET" ]; then
    if [ ! -f "$NET_STATE" ]; then
        echo "lab network state not found: $NET_STATE" >&2
        exit 1
    fi
    # shellcheck disable=SC1090
    . "$NET_STATE"
fi

if [ -z "${LAB_NET:-}" ]; then
    echo "LAB_NET missing" >&2
    exit 1
fi

if ! podman network exists "$LAB_NET" >/dev/null 2>&1; then
    echo "podman network does not exist: $LAB_NET" >&2
    exit 1
fi

if ! podman image exists "$IMAGE"; then
    echo "missing node image: $IMAGE" >&2
    exit 1
fi

PUBLISH_ARGS=""
if [ "${ULAB_PUBLISH_NODE_PORTS:-0}" = "1" ]; then
    PUBLISH_ARGS="-p 18001:18001 -p 18026:18026 -p 18028:18028 -p 18029:18029/udp -p 18030:18030"
fi

echo "podman: removing existing container if present: $CONTAINER_NAME"
podman rm -f "$CONTAINER_NAME" >/dev/null 2>&1 || true

echo "podman: starting $CONTAINER_NAME from $IMAGE network=$LAB_NET"

if [ -n "${ULAB_NODE_ENTRYPOINT:-}" ]; then
    # shellcheck disable=SC2086
    if ! podman run -d \
        --name "$CONTAINER_NAME" \
        --privileged \
        --device /dev/net/tun \
        --network "$LAB_NET" \
        --entrypoint "$ULAB_NODE_ENTRYPOINT" \
        $PUBLISH_ARGS \
        "$IMAGE" \
        ${ULAB_NODE_CMD:-}; then
        echo "podman: failed to start $CONTAINER_NAME on $LAB_NET" >&2
        podman network inspect "$LAB_NET" >&2 || true
        exit 1
    fi
else
    # shellcheck disable=SC2086
    if ! podman run -d \
        --name "$CONTAINER_NAME" \
        --privileged \
        --device /dev/net/tun \
        --network "$LAB_NET" \
        $PUBLISH_ARGS \
        "$IMAGE"; then
        echo "podman: failed to start $CONTAINER_NAME on $LAB_NET" >&2
        podman network inspect "$LAB_NET" >&2 || true
        exit 1
    fi
fi

CONTAINER_IP="$(container_ip_on_network "$CONTAINER_NAME" "$LAB_NET")"
if [ -z "$CONTAINER_IP" ]; then
    echo "podman: container has no IP on $LAB_NET: $CONTAINER_NAME" >&2
    podman inspect "$CONTAINER_NAME" >&2 || true
    exit 1
fi

echo "node-started node=$NODE_ID container=$CONTAINER_NAME ip=$CONTAINER_IP network=$LAB_NET run_dir=$RUN_DIR"
