#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -eu

if [ "$#" -ne 2 ]; then
    echo "usage: $0 <ue-id-or-ref> <run-dir>" >&2
    exit 2
fi

UE_KEY="$1"
RUN_DIR="$2"
STATE_FILE="$RUN_DIR/runtime-ues/$(printf "%s" "$UE_KEY" | tr -c 'A-Za-z0-9_.-' '-').env"

if [ ! -f "$STATE_FILE" ]; then
    echo "UE state not found: $STATE_FILE" >&2
    exit 1
fi

# shellcheck disable=SC1090
. "$STATE_FILE"

require_value() {
    name="$1"
    eval "value=\${$name:-}"
    if [ -z "$value" ]; then
        echo "$name missing in $STATE_FILE" >&2
        exit 1
    fi
}

require_value UE_CONTAINER
require_value TNODE_CONTAINER
require_value MEDIA_CONTAINER
require_value MEDIA_IP
require_value IMSI
require_value UE_IP

HTTP_PORT="${HTTP_PORT:-8080}"

if ! podman inspect -f '{{.State.Running}}' "$UE_CONTAINER" 2>/dev/null | grep -q '^true$'; then
    echo "UE container is not running: $UE_CONTAINER" >&2
    exit 1
fi

if ! podman inspect -f '{{.State.Running}}' "$TNODE_CONTAINER" 2>/dev/null | grep -q '^true$'; then
    echo "tower container is not running: $TNODE_CONTAINER" >&2
    exit 1
fi

if ! podman inspect -f '{{.State.Running}}' "$MEDIA_CONTAINER" 2>/dev/null | grep -q '^true$'; then
    echo "media container is not running: $MEDIA_CONTAINER" >&2
    exit 1
fi

if ! podman exec "$UE_CONTAINER" test -d /sys/class/net/tun0; then
    echo "UE tun0 not found in $UE_CONTAINER" >&2
    exit 1
fi

if ! podman exec "$TNODE_CONTAINER" \
    curl -fsS --max-time 3 "http://127.0.0.1:18028/v1/ue/$IMSI" 2>/dev/null | \
    grep -qi '"state"[[:space:]]*:[[:space:]]*"attached"'; then
    echo "EPC session is not attached for imsi=$IMSI tower=$TNODE_CONTAINER" >&2
    exit 1
fi

if ! podman exec "$TNODE_CONTAINER" \
    curl -fsS --max-time 3 "http://127.0.0.1:18030/v1/subscriber/imsi/$IMSI" \
    >/dev/null 2>&1; then
    echo "PCRF subscriber is unavailable for imsi=$IMSI tower=$TNODE_CONTAINER" >&2
    exit 1
fi

if ! podman exec "$TNODE_CONTAINER" \
    curl -fsS --max-time 3 "http://$MEDIA_IP:$HTTP_PORT/" >/dev/null 2>&1; then
    echo "media is unreachable from tower=$TNODE_CONTAINER ip=$MEDIA_IP" >&2
    exit 1
fi

echo "ue-session-healthy ue=$UE_KEY imsi=$IMSI ip=$UE_IP tower=$TNODE_CONTAINER media=$MEDIA_IP"
