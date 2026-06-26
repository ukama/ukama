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
UE_CIDR="${ULAB_UE_CIDR:-192.168.8.0/22}"
UE_ROUTE_PROBE_IP="${ULAB_UE_ROUTE_PROBE_IP:-192.168.8.2}"

RUN_ID="$(basename "$RUN_DIR")"
SAFE_RUN_ID="$(printf "%s" "$RUN_ID" | tr -c 'A-Za-z0-9-' '-' | sed 's/^-*//;s/-*$//')"
MEDIA_CONTAINER="ukama-media-$SAFE_RUN_ID"

need_cmd() {
    if ! command -v "$1" >/dev/null 2>&1; then
        echo "missing required command: $1" >&2
        exit 1
    fi
}

need_file() {
    if [ ! -f "$1" ]; then
        echo "missing $1" >&2
        exit 1
    fi
}

container_ip_on_network() {
    container="$1"
    network="$2"

    podman inspect -f '{{range $name, $net := .NetworkSettings.Networks}}{{if eq $name "'"$network"'"}}{{$net.IPAddress}}{{end}}{{end}}' \
        "$container" 2>/dev/null
}

find_tower_ip() {
    state=""

    if [ -n "${ULAB_MEDIA_TOWER_IP:-}" ]; then
        printf "%s\n" "$ULAB_MEDIA_TOWER_IP"
        return 0
    fi

    if [ ! -d "$SITE_STATE_DIR" ]; then
        return 1
    fi

    for state in "$SITE_STATE_DIR"/*.env; do
        [ -f "$state" ] || continue

        TNODE_IP=""
        # shellcheck disable=SC1090
        . "$state"

        if [ -n "${TNODE_IP:-}" ]; then
            printf "%s\n" "$TNODE_IP"
            return 0
        fi
    done

    return 1
}

configure_media_route() {
    container="$1"
    tower_ip="$2"

    echo "media: route ue-cidr=$UE_CIDR via tower=$tower_ip"

    if ! podman exec "$container" ip route replace "$UE_CIDR" via "$tower_ip"; then
        echo "failed to install media return route: $UE_CIDR via $tower_ip" >&2
        podman exec "$container" ip route >&2 || true
        return 1
    fi

    if ! podman exec "$container" \
        sh -lc "ip route get '$UE_ROUTE_PROBE_IP' | grep -q 'via $tower_ip'"; then
        echo "media return route verification failed: probe=$UE_ROUTE_PROBE_IP via=$tower_ip" >&2
        podman exec "$container" ip route >&2 || true
        podman exec "$container" ip route get "$UE_ROUTE_PROBE_IP" >&2 || true
        return 1
    fi

    return 0
}

need_cmd podman
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

TOWER_IP="$(find_tower_ip || true)"
if [ -z "$TOWER_IP" ]; then
    echo "could not find tower IP in $SITE_STATE_DIR; media needs a return route to $UE_CIDR" >&2
    find "$SITE_STATE_DIR" -maxdepth 1 -type f -name '*.env' -print >&2 2>/dev/null || true
    exit 1
fi

mkdir -p "$MEDIA_STATE_DIR"

echo "media: build $MEDIA_IMAGE"
podman build -t "$MEDIA_IMAGE" -f "$UE_DIR/media/Containerfile" "$UE_DIR"

echo "media: start $MEDIA_CONTAINER network=$LAB_NET"
podman rm -f "$MEDIA_CONTAINER" >/dev/null 2>&1 || true

podman run -d \
    --name "$MEDIA_CONTAINER" \
    --cap-add NET_ADMIN \
    --network "$LAB_NET" \
    "$MEDIA_IMAGE" >/dev/null

MEDIA_IP="$(container_ip_on_network "$MEDIA_CONTAINER" "$LAB_NET")"
if [ -z "$MEDIA_IP" ]; then
    echo "media container has no IP on $LAB_NET: $MEDIA_CONTAINER" >&2
    podman inspect "$MEDIA_CONTAINER" >&2 || true
    exit 1
fi

configure_media_route "$MEDIA_CONTAINER" "$TOWER_IP"

cat > "$MEDIA_STATE_FILE" <<STATE
MEDIA_CONTAINER=$MEDIA_CONTAINER
MEDIA_IMAGE=$MEDIA_IMAGE
MEDIA_IP=$MEDIA_IP
HTTP_PORT=$HTTP_PORT
IPERF_PORT=$IPERF_PORT
LAB_NET=$LAB_NET
TOWER_IP=$TOWER_IP
UE_CIDR=$UE_CIDR
UE_ROUTE_PROBE_IP=$UE_ROUTE_PROBE_IP
STATE

echo "media-ready container=$MEDIA_CONTAINER ip=$MEDIA_IP network=$LAB_NET route=$UE_CIDR via=$TOWER_IP"
