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
STATE_FILE="$RUN_DIR/runtime-nodes/$(printf "%s" "$NODE_KEY" | tr -c 'A-Za-z0-9_.-' '-').env"

if [ ! -f "$STATE_FILE" ]; then
    echo "stop-node: state not found $STATE_FILE"
    exit 0
fi

# shellcheck disable=SC1090
. "$STATE_FILE"

if [ -n "${CONTAINER_NAME:-}" ]; then
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
    MESH_POD_PREFIX="ukama-mesh-node-${FACTORY_NODE_ID}-"

    if [ ! -x "$ULAB_KUBECTL" ] &&
       ! command -v "$ULAB_KUBECTL" >/dev/null 2>&1; then
        echo "stop-node: kubectl not found: $ULAB_KUBECTL" >&2
        exit 1
    fi

    if ! MESH_PODS="$("$ULAB_KUBECTL" get pods \
        -n "$MESH_NAMESPACE" \
        -o custom-columns=NAME:.metadata.name --no-headers)"; then
        echo "stop-node: failed to list mesh pods for $FACTORY_NODE_ID" >&2
        exit 1
    fi

    MESH_POD="$(printf '%s\n' "$MESH_PODS" |
        awk -v prefix="$MESH_POD_PREFIX" \
            'index($0, prefix) == 1 { print; exit }')"

    if [ -n "$MESH_POD" ]; then
        echo "stop-node: delete mesh pod $MESH_POD"
        if ! "$ULAB_KUBECTL" delete pod "$MESH_POD" \
            -n "$MESH_NAMESPACE" --ignore-not-found=true; then
            echo "stop-node: failed to delete mesh pod $MESH_POD" >&2
            exit 1
        fi
    fi
fi

exit 0
