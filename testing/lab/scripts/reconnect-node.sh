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
    echo "reconnect-node: state not found $STATE_FILE" >&2
    exit 1
fi

# shellcheck disable=SC1090
. "$STATE_FILE"

if [ -z "${CONTAINER_NAME:-}" ]; then
    echo "reconnect-node: container missing in $STATE_FILE" >&2
    exit 1
fi

if ! podman container exists "$CONTAINER_NAME" >/dev/null 2>&1; then
    echo "reconnect-node: container not found: $CONTAINER_NAME" >&2
    exit 1
fi

if podman inspect -f '{{.State.Running}}' "$CONTAINER_NAME" 2>/dev/null |
    grep -q '^true$'; then
    echo "node-already-connected node=$NODE_KEY container=$CONTAINER_NAME"
    exit 0
fi

# Start the same stopped container.  Podman preserves its network attachment,
# filesystem and node identity, so the normal backend WebSocket reconnect path
# is exercised without rebuilding or replacing the virtual node.
podman start "$CONTAINER_NAME" >/dev/null

if ! podman inspect -f '{{.State.Running}}' "$CONTAINER_NAME" 2>/dev/null |
    grep -q '^true$'; then
    echo "reconnect-node: container failed to start: $CONTAINER_NAME" >&2
    exit 1
fi

echo "node-reconnected node=$NODE_KEY container=$CONTAINER_NAME"
