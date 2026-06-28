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

# No-sudo mode: media is a normal podman container on LAB_NET.
# UE still sends traffic through tun0 -> EPCEMU -> tower. The tower routes
# UE->media through its existing init-network uplink path. Media routes replies
# for the UE CIDR back to the tower's LAB_NET IP.
MEDIA_MODE="podman-net"
TUN_IF="${ULAB_MEDIA_TUN_IF:-tun3}"
TUN_TABLE="${ULAB_MEDIA_TUN_TABLE:-2000}"
BR_TABLE="${ULAB_MEDIA_BR_TABLE:-1000}"
MEDIA_BR="${ULAB_MEDIA_BR:-br0}"
UE_CIDR="${ULAB_UE_CIDR:-192.168.8.0/22}"
TOWER_LAB_IF="${ULAB_TOWER_LAB_IF:-eth0}"

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

container_running() {
    podman inspect -f '{{.State.Running}}' "$1" 2>/dev/null | grep -q '^true$'
}

container_ip_on_net() {
    c="$1"
    n="$2"
    podman inspect -f "{{with index .NetworkSettings.Networks \"$n\"}}{{.IPAddress}}{{end}}" "$c" 2>/dev/null
}

find_tower_state() {
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
            printf "%s\n" "$state"
            return 0
        fi
    done

    return 1
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

TOWER_STATE="$(find_tower_state || true)"
if [ -z "$TOWER_STATE" ]; then
    echo "could not find tower container in $SITE_STATE_DIR" >&2
    find "$SITE_STATE_DIR" -maxdepth 1 -type f -name '*.env' -print >&2 2>/dev/null || true
    exit 1
fi

# shellcheck disable=SC1090
. "$TOWER_STATE"

if [ -z "${TNODE_CONTAINER:-}" ]; then
    echo "TNODE_CONTAINER missing in $TOWER_STATE" >&2
    exit 1
fi

if ! container_running "$TNODE_CONTAINER"; then
    echo "tower container is not running: $TNODE_CONTAINER" >&2
    podman ps -a --filter "name=$TNODE_CONTAINER" >&2 || true
    exit 1
fi

if [ -z "${TNODE_IP:-}" ]; then
    TNODE_IP="$(container_ip_on_net "$TNODE_CONTAINER" "$LAB_NET")"
fi

if [ -z "$TNODE_IP" ]; then
    echo "could not determine tower IP on $LAB_NET for $TNODE_CONTAINER" >&2
    podman inspect "$TNODE_CONTAINER" >&2 || true
    exit 1
fi

mkdir -p "$MEDIA_STATE_DIR"

podman rm -f "$MEDIA_CONTAINER" >/dev/null 2>&1 || true

# Remove stale route/rules from earlier no-sudo media runs.
podman exec "$TNODE_CONTAINER" sh -lc "
    ip route del '${ULAB_MEDIA_OLD_IP:-0.0.0.0}/32' table '$TUN_TABLE' 2>/dev/null || true
" >/dev/null 2>&1 || true

echo "media: build $MEDIA_IMAGE"
podman build -t "$MEDIA_IMAGE" -f "$UE_DIR/media/Containerfile" "$UE_DIR"

echo "media: start $MEDIA_CONTAINER mode=$MEDIA_MODE network=$LAB_NET"
podman run -d \
    --name "$MEDIA_CONTAINER" \
    --privileged \
    --network "$LAB_NET" \
    "$MEDIA_IMAGE" >/dev/null

if ! container_running "$MEDIA_CONTAINER"; then
    echo "media container is not running after start: $MEDIA_CONTAINER" >&2
    podman ps -a --filter "name=$MEDIA_CONTAINER" >&2 || true
    podman logs --tail 80 "$MEDIA_CONTAINER" >&2 || true
    exit 1
fi

MEDIA_IP="$(container_ip_on_net "$MEDIA_CONTAINER" "$LAB_NET")"
if [ -z "$MEDIA_IP" ]; then
    echo "could not determine media IP on $LAB_NET for $MEDIA_CONTAINER" >&2
    podman inspect "$MEDIA_CONTAINER" >&2 || true
    exit 1
fi

# Return path: media must send UE subnet replies back to the tower container,
# not to the podman bridge gateway.
podman exec "$MEDIA_CONTAINER" \
    ip route replace "$UE_CIDR" via "$TNODE_IP"

# Uplink path: packets arriving from tun3 use table 2000. Use the existing
# init-network default gateway in that table, but add a specific host route so
# media traffic is deterministic.
TUN_GW="$(podman exec "$TNODE_CONTAINER" sh -lc "ip route show table '$TUN_TABLE' default 2>/dev/null | awk '{for (i=1; i<=NF; i++) if (\$i == \"via\") {print \$(i+1); exit}}'" || true)"

if [ -n "$TUN_GW" ]; then
    podman exec "$TNODE_CONTAINER" \
        ip route replace "$MEDIA_IP/32" via "$TUN_GW" dev "$MEDIA_BR" table "$TUN_TABLE"
else
    echo "media: warning no default gateway in table $TUN_TABLE; leaving table route unchanged" >&2
fi

podman exec "$TNODE_CONTAINER" sh -lc "
    ip route flush cache >/dev/null 2>&1 || true
    sysctl -w net.ipv4.ip_forward=1 >/dev/null 2>&1 || true
    sysctl -w net.ipv4.conf.'$TUN_IF'.rp_filter=0 >/dev/null 2>&1 || true
    sysctl -w net.ipv4.conf.'$TOWER_LAB_IF'.rp_filter=0 >/dev/null 2>&1 || true

    iptables -C FORWARD -i '$TUN_IF' -d '$MEDIA_IP/32' -j ACCEPT 2>/dev/null ||
        iptables -I FORWARD 1 -i '$TUN_IF' -d '$MEDIA_IP/32' -j ACCEPT

    iptables -C FORWARD -i '$TOWER_LAB_IF' -d '$UE_CIDR' -j ACCEPT 2>/dev/null ||
        iptables -I FORWARD 1 -i '$TOWER_LAB_IF' -d '$UE_CIDR' -j ACCEPT
" >/dev/null

# Basic checks.
if ! podman exec "$MEDIA_CONTAINER" curl -fsS --max-time 3 \
    "http://127.0.0.1:$HTTP_PORT/" >/dev/null 2>&1; then
    echo "media local HTTP check failed" >&2
    podman logs --tail 80 "$MEDIA_CONTAINER" >&2 || true
    exit 1
fi

if ! podman exec "$TNODE_CONTAINER" curl -fsS --max-time 3 \
    "http://$MEDIA_IP:$HTTP_PORT/" >/dev/null 2>&1; then
    echo "media HTTP check failed from tower: http://$MEDIA_IP:$HTTP_PORT/" >&2
    podman exec "$TNODE_CONTAINER" ip route >&2 || true
    podman exec "$TNODE_CONTAINER" ip route show table "$TUN_TABLE" >&2 || true
    podman exec "$MEDIA_CONTAINER" ip route >&2 || true
    exit 1
fi

cat > "$MEDIA_STATE_FILE" <<STATE
MEDIA_MODE=$MEDIA_MODE
MEDIA_CONTAINER=$MEDIA_CONTAINER
MEDIA_IMAGE=$MEDIA_IMAGE
MEDIA_IP=$MEDIA_IP
MEDIA_BR=$MEDIA_BR
MEDIA_IF=eth0
TNODE_CONTAINER=$TNODE_CONTAINER
TNODE_IP=$TNODE_IP
TUN_IF=$TUN_IF
TUN_TABLE=$TUN_TABLE
BR_TABLE=$BR_TABLE
UE_CIDR=$UE_CIDR
TOWER_LAB_IF=$TOWER_LAB_IF
TUN_GW=$TUN_GW
HTTP_PORT=$HTTP_PORT
IPERF_PORT=$IPERF_PORT
LAB_NET=$LAB_NET
STATE

echo "media-ready container=$MEDIA_CONTAINER ip=$MEDIA_IP mode=$MEDIA_MODE tower=$TNODE_CONTAINER"
