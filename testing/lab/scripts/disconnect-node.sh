#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -eu

if [ "$#" -ne 2 ]; then
    echo "usage: $0 <logical-node-id> <run-dir>" >&2
    exit 2
fi

NODE_KEY="$1"
RUN_DIR="$2"
STATE_NAME="$(printf "%s" "$NODE_KEY" | tr -c 'A-Za-z0-9_.-' '-')"
STATE_FILE="$RUN_DIR/runtime-nodes/$STATE_NAME.env"

if [ ! -f "$STATE_FILE" ]; then
    echo "disconnect-node: state not found $STATE_FILE" >&2
    exit 1
fi

# shellcheck disable=SC1090
. "$STATE_FILE"

if [ -z "${CONTAINER_NAME:-}" ]; then
    echo "disconnect-node: container missing in $STATE_FILE" >&2
    exit 1
fi

if ! podman container exists "$CONTAINER_NAME" >/dev/null 2>&1; then
    echo "disconnect-node: container not found: $CONTAINER_NAME" >&2
    exit 1
fi

if ! podman inspect -f '{{.State.Running}}' "$CONTAINER_NAME" 2>/dev/null |
    grep -q '^true$'; then
    echo "node-already-disconnected node=$NODE_KEY container=$CONTAINER_NAME"
    exit 0
fi

# Stop the virtual node process rather than only removing its Podman network
# interface.  This deterministically closes the backend WebSocket so the
# production offline path is exercised immediately.  An explicit podman stop
# suppresses the container's --restart policy until reconnect-node.sh starts it.
podman stop "$CONTAINER_NAME" >/dev/null

if podman inspect -f '{{.State.Running}}' "$CONTAINER_NAME" 2>/dev/null |
    grep -q '^true$'; then
    echo "disconnect-node: container is still running: $CONTAINER_NAME" >&2
    exit 1
fi

echo "node-disconnected node=$NODE_KEY container=$CONTAINER_NAME"
