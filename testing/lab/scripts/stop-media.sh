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

MEDIA_BR="${MEDIA_BR:-br0}"
TOWER_IF="${TOWER_IF:-ulabmed0}"

if [ -n "${TNODE_CONTAINER:-}" ] && [ -n "${TOWER_IF:-}" ]; then
    echo "stop-media: del-port $MEDIA_BR $TOWER_IF"
    podman exec "$TNODE_CONTAINER" \
        ovs-vsctl --if-exists del-port "$MEDIA_BR" "$TOWER_IF" \
        >/dev/null 2>&1 || true
fi

if [ -n "${MEDIA_CONTAINER:-}" ]; then
    echo "stop-media: rm $MEDIA_CONTAINER"
    podman rm -f "$MEDIA_CONTAINER" >/dev/null 2>&1 || true
fi

exit 0
