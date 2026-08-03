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
STATE_DIR="$RUN_DIR/runtime-media"
TIMEOUT_SEC="${ULAB_MEDIA_READY_TIMEOUT:-60}"
SLEEP_SEC=2

check_state() {
    STATE_FILE="$1"
    SITE_REF=""
    MEDIA_CONTAINER=""
    MEDIA_IP=""
    TNODE_CONTAINER=""
    TNODE_IP=""
    # shellcheck disable=SC1090
    . "$STATE_FILE"

    HTTP_PORT="${HTTP_PORT:-8080}"
    UE_CIDR="${UE_CIDR:-192.168.8.0/22}"
    TUN_TABLE="${TUN_TABLE:-2000}"

    start_ts="$(date +%s)"
    while :; do
        if podman inspect -f '{{.State.Running}}' "$MEDIA_CONTAINER" 2>/dev/null | grep -q '^true$' && \
           podman exec "$MEDIA_CONTAINER" curl -fsS --max-time 2 \
               "http://127.0.0.1:$HTTP_PORT/" >/dev/null 2>&1 && \
           podman exec "$MEDIA_CONTAINER" sh -lc 'pgrep iperf3 >/dev/null' \
               >/dev/null 2>&1 && \
           podman exec "$MEDIA_CONTAINER" sh -lc \
               "ip route show '$UE_CIDR' | grep -q 'via $TNODE_IP'" >/dev/null 2>&1 && \
           podman inspect -f '{{.State.Running}}' "$TNODE_CONTAINER" 2>/dev/null | grep -q '^true$' && \
           podman exec "$TNODE_CONTAINER" curl -fsS --max-time 2 \
               "http://$MEDIA_IP:$HTTP_PORT/" >/dev/null 2>&1; then
            echo "media-ready site=$SITE_REF container=$MEDIA_CONTAINER ip=$MEDIA_IP"
            return 0
        fi

        now_ts="$(date +%s)"
        if [ $((now_ts - start_ts)) -ge "$TIMEOUT_SEC" ]; then
            echo "media not ready: site=$SITE_REF container=$MEDIA_CONTAINER" >&2
            podman ps -a --filter "name=$MEDIA_CONTAINER" >&2 || true
            podman logs --tail 80 "$MEDIA_CONTAINER" >&2 || true
            podman exec "$MEDIA_CONTAINER" ip route >&2 || true
            podman exec "$TNODE_CONTAINER" ip route show table "$TUN_TABLE" >&2 || true
            exit 1
        fi
        sleep "$SLEEP_SEC"
    done
}

if [ ! -d "$STATE_DIR" ]; then
    echo "media state directory not found: $STATE_DIR" >&2
    exit 1
fi

found=0
for state in "$STATE_DIR"/*.env; do
    [ -f "$state" ] || continue
    [ "$(basename "$state")" = "media.env" ] && continue
    found=1
    check_state "$state"
done

if [ "$found" -eq 0 ] && [ -f "$STATE_DIR/media.env" ]; then
    check_state "$STATE_DIR/media.env"
    found=1
fi

if [ "$found" -eq 0 ]; then
    echo "no media state files found in $STATE_DIR" >&2
    exit 1
fi
