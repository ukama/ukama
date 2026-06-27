#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -eu

if [ "$#" -ne 1 ]; then
    echo "usage: $0 <run-dir>" >&2
    exit 2
fi

RUN_DIR="$1"
STATE_FILE="$RUN_DIR/runtime-media/media.env"
TIMEOUT_SEC=60
SLEEP_SEC=2

if [ ! -f "$STATE_FILE" ]; then
    echo "media state not found: $STATE_FILE" >&2
    exit 1
fi

# shellcheck disable=SC1090
. "$STATE_FILE"

HTTP_PORT="${HTTP_PORT:-8080}"
IPERF_PORT="${IPERF_PORT:-5201}"
MEDIA_IF="${MEDIA_IF:-eth0}"
MEDIA_BR="${MEDIA_BR:-br0}"
TOWER_IF="${TOWER_IF:-ulabmed0}"
TUN_TABLE="${TUN_TABLE:-2000}"

media_local_ready() {
    podman inspect -f '{{.State.Running}}' "$MEDIA_CONTAINER" 2>/dev/null | grep -q '^true$' && \
    podman exec "$MEDIA_CONTAINER" curl -fsS --max-time 2 \
        "http://127.0.0.1:$HTTP_PORT/" >/dev/null 2>&1 && \
    podman exec "$MEDIA_CONTAINER" sh -lc 'pgrep iperf3 >/dev/null' \
        >/dev/null 2>&1 && \
    podman exec "$MEDIA_CONTAINER" sh -lc \
        "ip addr show '$MEDIA_IF' | grep -q '$MEDIA_IP/'" >/dev/null 2>&1 && \
    podman exec "$MEDIA_CONTAINER" sh -lc \
        "ip route | grep -q 'default via $MEDIA_GW'" >/dev/null 2>&1
}

tower_path_ready() {
    podman inspect -f '{{.State.Running}}' "$TNODE_CONTAINER" 2>/dev/null | grep -q '^true$' && \
    podman exec "$TNODE_CONTAINER" ovs-vsctl port-to-br "$TOWER_IF" \
        >/dev/null 2>&1 && \
    podman exec "$TNODE_CONTAINER" curl -fsS --max-time 2 \
        "http://$MEDIA_IP:$HTTP_PORT/" >/dev/null 2>&1 && \
    podman exec "$TNODE_CONTAINER" sh -lc \
        "ip route show table '$TUN_TABLE' | grep -q '^$MEDIA_IP dev $MEDIA_BR'" \
        >/dev/null 2>&1
}

start_ts="$(date +%s)"
while :; do
    if media_local_ready && tower_path_ready; then
        echo "media-ready container=$MEDIA_CONTAINER ip=$MEDIA_IP bridge=$MEDIA_BR if=$TOWER_IF"
        exit 0
    fi

    now_ts="$(date +%s)"
    if [ $((now_ts - start_ts)) -ge "$TIMEOUT_SEC" ]; then
        echo "media not ready: $MEDIA_CONTAINER" >&2
        podman ps -a --filter "name=$MEDIA_CONTAINER" >&2 || true
        podman logs --tail 80 "$MEDIA_CONTAINER" >&2 || true
        echo "---- media net ----" >&2
        podman exec "$MEDIA_CONTAINER" ip addr >&2 || true
        podman exec "$MEDIA_CONTAINER" ip route >&2 || true
        echo "---- tower bridge ----" >&2
        podman exec "$TNODE_CONTAINER" ovs-vsctl show >&2 || true
        podman exec "$TNODE_CONTAINER" ip addr show "$TOWER_IF" >&2 || true
        podman exec "$TNODE_CONTAINER" ip route show table "$TUN_TABLE" >&2 || true
        podman exec "$TNODE_CONTAINER" curl -v --max-time 3 \
            "http://$MEDIA_IP:$HTTP_PORT/" >&2 || true
        exit 1
    fi

    sleep "$SLEEP_SEC"
done
