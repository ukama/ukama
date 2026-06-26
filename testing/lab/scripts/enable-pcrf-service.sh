#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

# Enable PCRF service on the tower node for ukama-lab smoke runs.

set -eu

if [ "$#" -lt 2 ]; then
    echo "usage: $0 <logical-tower-node-id> <run-dir>" >&2
    exit 2
fi

LOGICAL_NODE_ID="$1"
RUN_DIR="$2"

STATE_DIR="$RUN_DIR/runtime-nodes"
PCRF_PORT="${ULAB_PCRF_PORT:-18030}"
PCRF_PATH="${ULAB_PCRF_SERVICE_PATH:-/v1/service}"
TIMEOUT_SEC="${ULAB_PCRF_ENABLE_TIMEOUT_SEC:-60}"
SLEEP_SEC="${ULAB_PCRF_ENABLE_SLEEP_SEC:-2}"

safe_name() {
    printf "%s" "$1" | tr -c 'A-Za-z0-9_.-' '-'
}

state_file="$STATE_DIR/$(safe_name "$LOGICAL_NODE_ID").env"

if [ ! -f "$state_file" ]; then
    echo "pcrf-enable: node state not found: $state_file" >&2
    exit 1
fi

# shellcheck disable=SC1090
. "$state_file"

if [ -z "${CONTAINER_NAME:-}" ]; then
    echo "pcrf-enable: missing CONTAINER_NAME in $state_file" >&2
    exit 1
fi

url="http://127.0.0.1:${PCRF_PORT}${PCRF_PATH}"

echo "pcrf-enable: logical=$LOGICAL_NODE_ID factory=${FACTORY_NODE_ID:-} container=$CONTAINER_NAME url=$url"

start_ts="$(date +%s)"

while :; do
    body_file="/tmp/ukama-pcrf-enable-body.$$"

    code="$(
        podman exec "$CONTAINER_NAME" sh -lc \
            "curl -sS --max-time 5 -o '$body_file' -w '%{http_code}' \
             -X POST '$url' \
             -H 'Content-Type: application/json' \
             -d '{\"state\":\"on\"}'" \
            2>/dev/null || true
    )"

    if [ "$code" -ge 200 ] 2>/dev/null && [ "$code" -lt 300 ] 2>/dev/null; then
        echo "pcrf-enable: ok http=$code"
        podman exec "$CONTAINER_NAME" rm -f "$body_file" >/dev/null 2>&1 || true
        exit 0
    fi

    now_ts="$(date +%s)"
    elapsed=$((now_ts - start_ts))

    if [ "$elapsed" -ge "$TIMEOUT_SEC" ]; then
        echo "pcrf-enable: failed after ${TIMEOUT_SEC}s http=${code:-none}" >&2

        echo "---- pcrf enable response ----" >&2
        podman exec "$CONTAINER_NAME" sh -lc "cat '$body_file' 2>/dev/null || true" >&2 || true
        echo >&2

        echo "---- pcrf status ----" >&2
        podman exec "$CONTAINER_NAME" sh -lc \
            "curl -sS --max-time 5 'http://127.0.0.1:${PCRF_PORT}/v1/status' || true" >&2 || true
        echo >&2

        echo "---- pcrf ping ----" >&2
        podman exec "$CONTAINER_NAME" sh -lc \
            "curl -sS --max-time 5 'http://127.0.0.1:${PCRF_PORT}/v1/ping' || true" >&2 || true
        echo >&2

        podman exec "$CONTAINER_NAME" rm -f "$body_file" >/dev/null 2>&1 || true
        exit 1
    fi

    podman exec "$CONTAINER_NAME" rm -f "$body_file" >/dev/null 2>&1 || true
    sleep "$SLEEP_SEC"
done
