#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -u

if [ "$#" -ne 1 ]; then
    echo "usage: $0 <run-dir>" >&2
    exit 2
fi

RUN_DIR="$1"
STATE_FILE="$RUN_DIR/runtime-media/media.env"

if [ ! -f "$STATE_FILE" ]; then
    echo "stop-media: state not found $STATE_FILE"
    exit 0
fi

# shellcheck disable=SC1090
. "$STATE_FILE"

TUN_IF="${TUN_IF:-tun3}"
TUN_TABLE="${TUN_TABLE:-2000}"
UE_CIDR="${UE_CIDR:-192.168.8.0/22}"
TOWER_LAB_IF="${TOWER_LAB_IF:-eth0}"

if [ -n "${TNODE_CONTAINER:-}" ] && [ -n "${MEDIA_IP:-}" ]; then
    echo "stop-media: delete media routes/rules"
    podman exec "$TNODE_CONTAINER" sh -lc "
        iptables -D FORWARD -i '$TUN_IF' -d '$MEDIA_IP/32' -j ACCEPT 2>/dev/null || true
        iptables -D FORWARD -i '$TOWER_LAB_IF' -d '$UE_CIDR' -j ACCEPT 2>/dev/null || true
        ip route del '$MEDIA_IP/32' table '$TUN_TABLE' 2>/dev/null || true
        ip route flush cache >/dev/null 2>&1 || true
    " >/dev/null 2>&1 || true
fi

if [ -n "${MEDIA_CONTAINER:-}" ]; then
    echo "stop-media: rm $MEDIA_CONTAINER"
    podman rm -f "$MEDIA_CONTAINER" >/dev/null 2>&1 || true
fi

exit 0
