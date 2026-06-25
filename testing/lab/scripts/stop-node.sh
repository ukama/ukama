#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -eu

if [ "$#" -lt 2 ]; then
    echo "usage: $0 <logical-node-id> <run-dir>" >&2
    exit 2
fi

LOGICAL_NODE_ID="$1"
RUN_DIR="$2"
STATE_DIR="$RUN_DIR/runtime-nodes"

safe_name() {
    printf "%s" "$1" | tr -c 'A-Za-z0-9_.-' '-'
}

state_file="$STATE_DIR/$(safe_name "$LOGICAL_NODE_ID").env"

if [ ! -f "$state_file" ]; then
    echo "stop-node: no state for $LOGICAL_NODE_ID"
    exit 0
fi

# shellcheck disable=SC1090
. "$state_file"

echo "stop-node: logical=$LOGICAL_NODE_ID factory=${FACTORY_NODE_ID:-} container=${CONTAINER_NAME:-}"

if [ -n "${CONTAINER_NAME:-}" ]; then
    podman rm -f "$CONTAINER_NAME" >/dev/null 2>&1 || true
fi

# Remove only node-specific image tags. Do not remove by IMAGE ID, because the
# same image id may also have the virtualnode-base:* tag.
if [ -n "${FACTORY_NODE_ID:-}" ]; then
    for ref in \
        "localhost/testing/virtualnode:${FACTORY_NODE_ID}" \
        "testing/virtualnode:${FACTORY_NODE_ID}" \
        "localhost:5000/testing/virtualnode:${FACTORY_NODE_ID}"
    do
        if podman image exists "$ref" 2>/dev/null; then
            echo "remove-node-image-tag: $ref"
            podman rmi "$ref" >/dev/null 2>&1 || true
        fi
    done
fi

exit 0
