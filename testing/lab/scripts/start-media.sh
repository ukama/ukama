#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -eu

if [ "$#" -ne 2 ]; then
    echo "usage: $0 <repo> <run-dir>" >&2
    exit 2
fi

REPO="$1"
RUN_DIR="$2"

UE_DIR="$REPO/testing/ue"
NET_STATE="$RUN_DIR/runtime-net/net.env"
SITE_STATE_DIR="$RUN_DIR/runtime-sites"
MEDIA_STATE_DIR="$RUN_DIR/runtime-media"
MEDIA_STATE_FILE="$MEDIA_STATE_DIR/media.env"

MEDIA_IMAGE="ukama/media:dev"
HTTP_PORT=8080
IPERF_PORT=5201

# Media is a user-plane host behind the tower OVS bridge.
MEDIA_BR="${ULAB_MEDIA_BR:-br0}"
MEDIA_IP="${ULAB_MEDIA_IP:-10.10.10.100}"
MEDIA_PREFIX="${ULAB_MEDIA_PREFIX:-24}"
MEDIA_GW="${ULAB_MEDIA_GW:-10.10.10.1}"
TOWER_IF="${ULAB_MEDIA_TOWER_IF:-ulabmed0}"
MEDIA_IF="${ULAB_MEDIA_IF:-eth0}"

# Only host network namespace plumbing needs privilege.
# Keep podman rootless. Example:
#   sudo -v
#   ULAB_NETNS_SUDO=sudo ./bin/ukama-lab validate ...
NETNS_SUDO="${ULAB_NETNS_SUDO:-}"

RUN_ID="$(basename "$RUN_DIR")"
SAFE_RUN_ID="$(printf "%s" "$RUN_ID" | tr -c 'A-Za-z0-9-' '-' | sed 's/^-*//;s/-*$//')"
MEDIA_CONTAINER="ukama-media-$SAFE_RUN_ID"

HOST_TOWER_IF="umt$$"
HOST_MEDIA_IF="umm$$"

run_priv() {
    if [ -n "$NETNS_SUDO" ]; then
        # shellcheck disable=SC2086
        $NETNS_SUDO "$@"
    else
        "$@"
    fi
}

host_ip() {
    run_priv ip "$@"
}

host_nsenter() {
    run_priv nsenter "$@"
}

need_cmd() {
    if ! command -v "$1" >/dev/null 2>&1; then
        echo "missing required command: $1" >&2
        exit 1
    fi
}

need_priv_cmd() {
    if [ -n "$NETNS_SUDO" ]; then
        cmd="${NETNS_SUDO%% *}"
        if ! command -v "$cmd" >/dev/null 2>&1; then
            echo "missing required command: $cmd" >&2
            exit 1
        fi
    fi
}

need_file() {
    if [ ! -f "$1" ]; then
        echo "missing $1" >&2
        exit 1
    fi
}

container_pid() {
    podman inspect -f '{{.State.Pid}}' "$1" 2>/dev/null
}

container_running() {
    podman inspect -f '{{.State.Running}}' "$1" 2>/dev/null | grep -q '^true$'
}

find_tower_container() {
    state=""

    if [ ! -d "$SITE_STATE_DIR" ]; then
        return 1
    fi

    for state in "$SITE_STATE_DIR"/*.env; do
        [ -f "$state" ] || continue

        TNODE_CONTAINER=""
        TNODE_IP=""
        # shellcheck disable=SC1090
        . "$state"

        if [ -n "${TNODE_CONTAINER:-}" ]; then
            printf "%s\n" "$TNODE_CONTAINER"
            return 0
        fi
    done

    return 1
}

cleanup_partial_veth() {
    host_ip link del "$HOST_TOWER_IF" >/dev/null 2>&1 || true
    host_ip link del "$HOST_MEDIA_IF" >/dev/null 2>&1 || true

    if [ -n "${TNODE_CONTAINER:-}" ]; then
        podman exec "$TNODE_CONTAINER" \
            ovs-vsctl --if-exists del-port "$MEDIA_BR" "$TOWER_IF" >/dev/null 2>&1 || true
        podman exec "$TNODE_CONTAINER" \
            ip link del "$TOWER_IF" >/dev/null 2>&1 || true
    fi

    if [ -n "${MEDIA_CONTAINER:-}" ]; then
        podman exec "$MEDIA_CONTAINER" \
            ip link del "$MEDIA_IF" >/dev/null 2>&1 || true
    fi
}

attach_media_to_tower_bridge() {
    tower_container="$1"
    media_container="$2"
    tower_pid="$3"
    media_pid="$4"

    echo "media: attach user-plane ip=$MEDIA_IP/$MEDIA_PREFIX gw=$MEDIA_GW bridge=$MEDIA_BR"

    cleanup_partial_veth

    podman exec "$tower_container" \
        ovs-vsctl --if-exists del-port "$MEDIA_BR" "$TOWER_IF" >/dev/null 2>&1 || true

    host_ip link add "$HOST_TOWER_IF" type veth peer name "$HOST_MEDIA_IF"

    host_ip link set "$HOST_TOWER_IF" netns "$tower_pid"
    host_ip link set "$HOST_MEDIA_IF" netns "$media_pid"

    host_nsenter -t "$tower_pid" -n ip link set "$HOST_TOWER_IF" name "$TOWER_IF"
    host_nsenter -t "$tower_pid" -n ip link set "$TOWER_IF" up

    podman exec "$tower_container" \
        ovs-vsctl --may-exist add-port "$MEDIA_BR" "$TOWER_IF"
    podman exec "$tower_container" ip link set "$TOWER_IF" up

    host_nsenter -t "$media_pid" -n ip link set lo up
    host_nsenter -t "$media_pid" -n ip link set "$HOST_MEDIA_IF" name "$MEDIA_IF"
    host_nsenter -t "$media_pid" -n ip addr replace "$MEDIA_IP/$MEDIA_PREFIX" dev "$MEDIA_IF"
    host_nsenter -t "$media_pid" -n ip link set "$MEDIA_IF" up
    host_nsenter -t "$media_pid" -n ip route replace default via "$MEDIA_GW"

    if ! podman exec "$tower_container" curl -fsS --max-time 3 \
        "http://$MEDIA_IP:$HTTP_PORT/" >/dev/null 2>&1; then
        echo "media user-plane HTTP check failed from tower: http://$MEDIA_IP:$HTTP_PORT/" >&2
        podman exec "$tower_container" ip addr show "$TOWER_IF" >&2 || true
        podman exec "$tower_container" ovs-vsctl show >&2 || true
        podman exec "$media_container" ip addr >&2 || true
        podman exec "$media_container" ip route >&2 || true
        return 1
    fi

    return 0
}

need_cmd podman
need_cmd ip
need_cmd nsenter
need_priv_cmd
need_file "$UE_DIR/media/Containerfile"
need_file "$NET_STATE"

# shellcheck disable=SC1090
. "$NET_STATE"

if [ -z "${LAB_NET:-}" ]; then
    echo "LAB_NET missing in $NET_STATE" >&2
    exit 1
fi

if ! podman network exists "$LAB_NET" >/dev/null 2>&1; then
    echo "podman network does not exist: $LAB_NET" >&2
    exit 1
fi

TNODE_CONTAINER="$(find_tower_container || true)"
if [ -z "$TNODE_CONTAINER" ]; then
    echo "could not find tower container in $SITE_STATE_DIR" >&2
    find "$SITE_STATE_DIR" -maxdepth 1 -type f -name '*.env' -print >&2 2>/dev/null || true
    exit 1
fi

if ! container_running "$TNODE_CONTAINER"; then
    echo "tower container is not running: $TNODE_CONTAINER" >&2
    podman ps -a --filter "name=$TNODE_CONTAINER" >&2 || true
    exit 1
fi

mkdir -p "$MEDIA_STATE_DIR"

# Remove previous media and stale OVS port for this tower.
podman rm -f "$MEDIA_CONTAINER" >/dev/null 2>&1 || true
podman exec "$TNODE_CONTAINER" \
    ovs-vsctl --if-exists del-port "$MEDIA_BR" "$TOWER_IF" >/dev/null 2>&1 || true
podman exec "$TNODE_CONTAINER" \
    ip link del "$TOWER_IF" >/dev/null 2>&1 || true

trap cleanup_partial_veth EXIT HUP INT TERM

echo "media: build $MEDIA_IMAGE"
podman build -t "$MEDIA_IMAGE" -f "$UE_DIR/media/Containerfile" "$UE_DIR"

echo "media: start $MEDIA_CONTAINER user-plane=$MEDIA_BR"
podman run -d \
    --name "$MEDIA_CONTAINER" \
    --privileged \
    --network none \
    "$MEDIA_IMAGE" >/dev/null

if ! container_running "$MEDIA_CONTAINER"; then
    echo "media container is not running after start: $MEDIA_CONTAINER" >&2
    podman ps -a --filter "name=$MEDIA_CONTAINER" >&2 || true
    podman logs --tail 80 "$MEDIA_CONTAINER" >&2 || true
    exit 1
fi

TOWER_PID="$(container_pid "$TNODE_CONTAINER")"
MEDIA_PID="$(container_pid "$MEDIA_CONTAINER")"

if [ -z "$TOWER_PID" ] || [ "$TOWER_PID" = "0" ]; then
    echo "invalid tower pid for $TNODE_CONTAINER: $TOWER_PID" >&2
    exit 1
fi

if [ -z "$MEDIA_PID" ] || [ "$MEDIA_PID" = "0" ]; then
    echo "invalid media pid for $MEDIA_CONTAINER: $MEDIA_PID" >&2
    exit 1
fi

if ! podman exec "$TNODE_CONTAINER" ovs-vsctl br-exists "$MEDIA_BR"; then
    echo "tower bridge not found: $MEDIA_BR" >&2
    podman exec "$TNODE_CONTAINER" ovs-vsctl show >&2 || true
    exit 1
fi

attach_media_to_tower_bridge "$TNODE_CONTAINER" "$MEDIA_CONTAINER" \
    "$TOWER_PID" "$MEDIA_PID"

trap - EXIT HUP INT TERM

cat > "$MEDIA_STATE_FILE" <<STATE
MEDIA_CONTAINER=$MEDIA_CONTAINER
MEDIA_IMAGE=$MEDIA_IMAGE
MEDIA_IP=$MEDIA_IP
MEDIA_PREFIX=$MEDIA_PREFIX
MEDIA_GW=$MEDIA_GW
MEDIA_BR=$MEDIA_BR
MEDIA_IF=$MEDIA_IF
TOWER_IF=$TOWER_IF
TNODE_CONTAINER=$TNODE_CONTAINER
HTTP_PORT=$HTTP_PORT
IPERF_PORT=$IPERF_PORT
LAB_NET=$LAB_NET
NETNS_SUDO=$NETNS_SUDO
STATE

echo "media-ready container=$MEDIA_CONTAINER ip=$MEDIA_IP bridge=$MEDIA_BR tower=$TNODE_CONTAINER"
