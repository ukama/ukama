#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -eu

if [ "$#" -ne 2 ]; then
    echo "usage: $0 <target> <run-dir>" >&2
    exit 2
fi

TARGET="$1"
RUN_DIR="$2"
HOLD_DIR="$RUN_DIR/failure-controls/$TARGET"

if [ ! -d "$HOLD_DIR" ]; then
    echo "release-held-nodes: no held nodes target=$TARGET"
    exit 0
fi

for HOLD_FILE in "$HOLD_DIR"/*.env; do
    [ -f "$HOLD_FILE" ] || continue

    LOGICAL_NODE_ID=""
    CONTAINER_NAME=""

    # shellcheck disable=SC1090
    . "$HOLD_FILE"

    if [ -n "$CONTAINER_NAME" ] &&
       podman inspect "$CONTAINER_NAME" >/dev/null 2>&1; then
        podman start "$CONTAINER_NAME" >/dev/null
    fi

    echo "node-released node=$LOGICAL_NODE_ID" \
        "container=$CONTAINER_NAME target=$TARGET"
    rm -f "$HOLD_FILE"
done

rmdir "$HOLD_DIR" 2>/dev/null || true
exit 0
