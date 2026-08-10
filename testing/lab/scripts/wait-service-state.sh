#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

# Wait until EPCEMU and PCRF both hold the requested service state long
# enough to be safe for the next scenario step. This script is read-only: it
# verifies controller convergence without forcing either node application.

set -eu

if [ "$#" -ne 3 ]; then
    echo "usage: $0 <logical-tower-node-id> <on|off> <run-dir>" >&2
    exit 2
fi

LOGICAL_NODE_ID="$1"
EXPECTED_STATE="$2"
RUN_DIR="$3"

STATE_DIR="$RUN_DIR/runtime-nodes"
TIMEOUT_SEC="${ULAB_SERVICE_STATE_TIMEOUT_SEC:-120}"
POLL_SEC="${ULAB_SERVICE_STATE_POLL_SEC:-2}"
STABLE_POLLS="${ULAB_SERVICE_STATE_STABLE_POLLS:-3}"

case "$EXPECTED_STATE" in
    on)
        EXPECTED_ADMISSION="enabled"
        ;;
    off)
        EXPECTED_ADMISSION="disabled"
        ;;
    *)
        echo "invalid service state: $EXPECTED_STATE" >&2
        exit 2
        ;;
esac

safe_name() {
    printf "%s" "$1" | tr -c 'A-Za-z0-9_.-' '-'
}

container_running() {
    podman inspect -f '{{.State.Running}}' "$1" 2>/dev/null | grep -q '^true$'
}

container_get() {
    port="$1"
    path="$2"

    podman exec "$CONTAINER_NAME" sh -lc \
        "curl -sS --max-time 5 'http://127.0.0.1:${port}${path}'" \
        2>/dev/null || true
}

matches_state() {
    body="$1"

    printf "%s" "$body" | \
        grep -Eq '"state"[[:space:]]*:[[:space:]]*"'"$EXPECTED_STATE"'"' && \
    printf "%s" "$body" | \
        grep -Eq '"admission"[[:space:]]*:[[:space:]]*"'"$EXPECTED_ADMISSION"'"'
}

state_file="$STATE_DIR/$(safe_name "$LOGICAL_NODE_ID").env"

if [ ! -f "$state_file" ]; then
    echo "service-state: node state not found: $state_file" >&2
    exit 1
fi

# shellcheck disable=SC1090
. "$state_file"

if [ -z "${CONTAINER_NAME:-}" ] || ! container_running "$CONTAINER_NAME"; then
    echo "service-state: tower container is not running: ${CONTAINER_NAME:-}" >&2
    exit 1
fi

start_ts="$(date +%s)"
stable=0
epc_body=""
pcrf_body=""

while :; do
    epc_body="$(container_get 18028 /v1/status)"
    pcrf_body="$(container_get 18030 /v1/service)"

    if matches_state "$epc_body" && matches_state "$pcrf_body"; then
        stable=$((stable + 1))
        if [ "$stable" -ge "$STABLE_POLLS" ]; then
            echo "service-state: ok state=$EXPECTED_STATE stable_polls=$stable"
            exit 0
        fi
    else
        stable=0
    fi

    now_ts="$(date +%s)"
    if [ $((now_ts - start_ts)) -ge "$TIMEOUT_SEC" ]; then
        echo "service-state: failed state=$EXPECTED_STATE after ${TIMEOUT_SEC}s" >&2
        echo "EPCEMU: $epc_body" >&2
        echo "PCRF: $pcrf_body" >&2
        exit 1
    fi

    sleep "$POLL_SEC"
done
