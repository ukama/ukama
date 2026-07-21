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
NETWORKS_TEMPLATE='{{range $name, $_ := .NetworkSettings.Networks}}{{println $name}}{{end}}'

if [ ! -f "$STATE_FILE" ]; then
    echo "reconnect-node: state not found $STATE_FILE" >&2
    exit 1
fi

# shellcheck disable=SC1090
. "$STATE_FILE"

if [ -z "${LAB_NET:-}" ] && [ -f "$RUN_DIR/runtime-net/net.env" ]; then
    # shellcheck disable=SC1090
    . "$RUN_DIR/runtime-net/net.env"
fi

if [ -z "${CONTAINER_NAME:-}" ] || [ -z "${LAB_NET:-}" ]; then
    echo "reconnect-node: container or network missing in $STATE_FILE" >&2
    exit 1
fi

if ! podman inspect -f '{{.State.Running}}' "$CONTAINER_NAME" 2>/dev/null |
    grep -q '^true$'; then
    echo "reconnect-node: container is not running: $CONTAINER_NAME" >&2
    exit 1
fi

if podman inspect -f "$NETWORKS_TEMPLATE" \
    "$CONTAINER_NAME" 2>/dev/null | grep -Fxq "$LAB_NET"; then
    echo "node-network-already-connected node=$NODE_KEY" \
        "container=$CONTAINER_NAME network=$LAB_NET"
else
    podman network connect "$LAB_NET" "$CONTAINER_NAME"
    echo "node-network-reconnected node=$NODE_KEY" \
        "container=$CONTAINER_NAME network=$LAB_NET"
fi
