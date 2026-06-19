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

. "$STATE_FILE"

start_ts="$(date +%s)"
while :; do
    if podman inspect -f '{{.State.Running}}' "$MEDIA_CONTAINER" 2>/dev/null | grep -q '^true$' && \
       podman exec "$MEDIA_CONTAINER" curl -fsS --max-time 2 http://127.0.0.1:8080/ >/dev/null 2>&1 && \
       podman exec "$MEDIA_CONTAINER" sh -lc 'pgrep iperf3 >/dev/null' >/dev/null 2>&1; then
        echo "media-ready container=$MEDIA_CONTAINER ip=$MEDIA_IP"
        exit 0
    fi

    now_ts="$(date +%s)"
    if [ $((now_ts - start_ts)) -ge "$TIMEOUT_SEC" ]; then
        echo "media not ready: $MEDIA_CONTAINER" >&2
        podman ps -a --filter "name=$MEDIA_CONTAINER" >&2 || true
        podman logs --tail 80 "$MEDIA_CONTAINER" >&2 || true
        exit 1
    fi

    sleep "$SLEEP_SEC"
done
