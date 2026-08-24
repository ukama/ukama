#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -eu

if [ "$#" -ne 3 ]; then
    echo "usage: $0 <target> <logical-node-id> <run-dir>" >&2
    exit 2
fi

TARGET="$1"
NODE_KEY="$2"
RUN_DIR="$3"
STATE_NAME="$(printf "%s" "$NODE_KEY" | tr -c 'A-Za-z0-9_.-' '-')"
STATE_FILE="$RUN_DIR/runtime-nodes/$STATE_NAME.env"
HOLD_DIR="$RUN_DIR/failure-controls/$TARGET"
HOLD_FILE="$HOLD_DIR/$STATE_NAME.env"

if [ ! -f "$STATE_FILE" ]; then
    echo "hold-node: state not found $STATE_FILE" >&2
    exit 1
fi

# shellcheck disable=SC1090
. "$STATE_FILE"

if [ -z "${CONTAINER_NAME:-}" ]; then
    echo "hold-node: container missing in $STATE_FILE" >&2
    exit 1
fi

mkdir -p "$HOLD_DIR"

if [ -f "$HOLD_FILE" ]; then
    echo "hold-node: already held node=$NODE_KEY target=$TARGET"
    exit 0
fi

{
    echo "LOGICAL_NODE_ID=$NODE_KEY"
    echo "CONTAINER_NAME=$CONTAINER_NAME"
} > "$HOLD_FILE"

if podman inspect "$CONTAINER_NAME" >/dev/null 2>&1; then
    podman stop -t 1 "$CONTAINER_NAME" >/dev/null
fi

echo "node-held node=$NODE_KEY container=$CONTAINER_NAME target=$TARGET"
