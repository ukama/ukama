#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -u

if [ "$#" -ne 2 ]; then
    echo "usage: $0 <logical-node-id> <run-dir>" >&2
    exit 2
fi

NODE_KEY="$1"
RUN_DIR="$2"
HOST_CONTROL="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)/node-host-control.sh"
STATE_FILE="$RUN_DIR/runtime-nodes/$(printf "%s" "$NODE_KEY" | tr -c 'A-Za-z0-9_.-' '-').env"

if [ ! -f "$STATE_FILE" ]; then
    echo "stop-node: state not found $STATE_FILE"
    exit 0
fi

# shellcheck disable=SC1090
. "$STATE_FILE"

if [ -n "${CONTAINER_NAME:-}" ]; then
    "$HOST_CONTROL" remove "$CONTAINER_NAME" || exit 1
    echo "stop-node: rm $CONTAINER_NAME"
    podman rm -f "$CONTAINER_NAME" >/dev/null 2>&1 || true
fi

# Remove only node-specific image tags. Do not remove by image id because the
# same image id may also carry virtualnode-base:* or git-sha tags we want kept.
if [ -n "${FACTORY_NODE_ID:-}" ]; then
    for img in \
        "testing/virtualnode:${FACTORY_NODE_ID}" \
        "localhost/testing/virtualnode:${FACTORY_NODE_ID}" \
        "localhost:5000/testing/virtualnode:${FACTORY_NODE_ID}"
    do
        if podman image exists "$img" >/dev/null 2>&1; then
            echo "stop-node: rmi tag $img"
            podman rmi "$img" >/dev/null 2>&1 || true
        fi
    done
fi

if [ -n "${ULAB_KUBECTL:-}" ] && [ -n "${FACTORY_NODE_ID:-}" ]; then
    MESH_NAMESPACE="${ULAB_MESH_NAMESPACE:-ukama-messaging}"
    MESH_NODE_NAME="ukama-mesh-node-${FACTORY_NODE_ID}"

    if [ ! -x "$ULAB_KUBECTL" ] &&
       ! command -v "$ULAB_KUBECTL" >/dev/null 2>&1; then
        echo "stop-node: kubectl not found: $ULAB_KUBECTL" >&2
        exit 1
    fi

    if ! MESH_RESOURCES="$("$ULAB_KUBECTL" get pods,services \
        -n "$MESH_NAMESPACE" -o name)"; then
        echo "stop-node: failed to list mesh pods/services for $FACTORY_NODE_ID" >&2
        exit 1
    fi

    MESH_RESOURCES="$(printf '%s\n' "$MESH_RESOURCES" |
        awk -F/ -v node="$MESH_NODE_NAME" \
            '$2 == node || index($2, node "-") == 1 { print }')"

    MESH_DELETE_FAILED=0
    for MESH_RESOURCE in $MESH_RESOURCES; do
        echo "stop-node: delete mesh resource $MESH_RESOURCE"
        if ! "$ULAB_KUBECTL" delete "$MESH_RESOURCE" \
            -n "$MESH_NAMESPACE" --ignore-not-found=true; then
            echo "stop-node: failed to delete mesh resource $MESH_RESOURCE" >&2
            MESH_DELETE_FAILED=1
        fi
    done
    if [ "$MESH_DELETE_FAILED" -ne 0 ]; then
        exit 1
    fi
fi

exit 0
