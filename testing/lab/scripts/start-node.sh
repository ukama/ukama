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
PROBE_IMAGE="${ULAB_NET_PROBE_IMAGE:-alpine:3.20}"
NET_STATE="$RUN_DIR/runtime-net/net.env"
LAB_NET="${LAB_NET:-}"
HOST_CONTROL="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)/node-host-control.sh"

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
need_cmd awk

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
    case "$NODE_ID" in
        *-tnode-*)
            PUBLISH_ARGS="-p 18001:18001 \
                -p 18026:18026 \
                -p 18028:18028 \
                -p 18029:18029/udp \
                -p 18030:18030 \
                -p 127.0.0.1:${ULAB_LIFECYCLE_TOWER_HOST_PORT:-18033}:${ULAB_LIFECYCLE_PORT:-18033}"
            ;;
        *-cnode-*)
            PUBLISH_ARGS="-p 127.0.0.1:${ULAB_LIFECYCLE_CONTROLLER_HOST_PORT:-18034}:${ULAB_LIFECYCLE_PORT:-18033}"
            ;;
        *-anode-*)
            PUBLISH_ARGS="-p 127.0.0.1:${ULAB_LIFECYCLE_AMPLIFIER_HOST_PORT:-18035}:${ULAB_LIFECYCLE_PORT:-18033}"
            ;;
    esac
fi

"$HOST_CONTROL" check "$CONTAINER_NAME"
"$HOST_CONTROL" remove "$CONTAINER_NAME"
echo "podman: removing existing container if present: $CONTAINER_NAME"
podman rm -f "$CONTAINER_NAME" >/dev/null 2>&1 || true

# Let Podman choose a free address, then store it as a static address on the
# node before its first boot. UE endpoints and media routes retain this IP
# while systemd stops and starts the same node container during a reboot.
if ! PROBE_ADDR="$(podman run --rm --network "$LAB_NET" \
    --entrypoint /bin/sh "$PROBE_IMAGE" -c 'ip -4 addr show dev eth0')"; then
    echo "podman: failed to allocate node IP on $LAB_NET" >&2
    exit 1
fi
NODE_IP="$(printf '%s\n' "$PROBE_ADDR" | \
    awk '$1 == "inet" { split($2, addr, "/"); print addr[1]; exit }')"
if [ -z "$NODE_IP" ]; then
    echo "podman: network probe has no IPv4 address on $LAB_NET" >&2
    exit 1
fi

echo "podman: starting $CONTAINER_NAME from $IMAGE network=$LAB_NET ip=$NODE_IP"

if [ -n "${ULAB_NODE_ENTRYPOINT:-}" ]; then
    # shellcheck disable=SC2086
    if ! podman create \
        --name "$CONTAINER_NAME" \
        --restart=no \
        --label io.ukama.lab.restart-manager=systemd \
        --privileged \
        --device /dev/net/tun \
        --network "$LAB_NET" \
        --ip "$NODE_IP" \
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
    if ! podman create \
        --name "$CONTAINER_NAME" \
        --restart=no \
        --label io.ukama.lab.restart-manager=systemd \
        --privileged \
        --device /dev/net/tun \
        --network "$LAB_NET" \
        --ip "$NODE_IP" \
        $PUBLISH_ARGS \
        "$IMAGE"; then
        echo "podman: failed to start $CONTAINER_NAME on $LAB_NET" >&2
        podman network inspect "$LAB_NET" >&2 || true
        exit 1
    fi
fi

if ! "$HOST_CONTROL" install "$CONTAINER_NAME" "$RUN_DIR"; then
    echo "podman: failed to start host service for $CONTAINER_NAME" >&2
    "$HOST_CONTROL" remove "$CONTAINER_NAME" || exit 1
    podman rm -f "$CONTAINER_NAME" >/dev/null 2>&1 || true
    exit 1
fi

CONTAINER_IP="$(container_ip_on_network "$CONTAINER_NAME" "$LAB_NET")"
if [ -z "$CONTAINER_IP" ]; then
    echo "podman: container has no IP on $LAB_NET: $CONTAINER_NAME" >&2
    podman inspect "$CONTAINER_NAME" >&2 || true
    exit 1
fi

echo "node-started node=$NODE_ID container=$CONTAINER_NAME ip=$CONTAINER_IP network=$LAB_NET run_dir=$RUN_DIR"
