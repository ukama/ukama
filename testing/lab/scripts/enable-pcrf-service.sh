#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -eu

if [ "$#" -ne 2 ]; then
    echo "usage: $0 <logical-tower-node-id> <run-dir>" >&2
    exit 2
fi

LOGICAL_NODE_ID="$1"
RUN_DIR="$2"

STATE_DIR="$RUN_DIR/runtime-nodes"
NET_STATE="$RUN_DIR/runtime-net/net.env"
PCRF_PORT="${ULAB_PCRF_PORT:-18030}"
PCRF_PATH="${ULAB_PCRF_SERVICE_PATH:-/v1/service}"
TIMEOUT_SEC="${ULAB_PCRF_ENABLE_TIMEOUT_SEC:-60}"
SLEEP_SEC="${ULAB_PCRF_ENABLE_SLEEP_SEC:-2}"

safe_name() {
    printf "%s" "$1" | tr -c 'A-Za-z0-9_.-' '-'
}

container_ip_on_network() {
    container="$1"
    network="$2"

    podman inspect -f '{{range $name, $net := .NetworkSettings.Networks}}{{if eq $name "'"$network"'"}}{{$net.IPAddress}}{{end}}{{end}}' \
        "$container" 2>/dev/null
}

need_cmd() {
    command -v "$1" >/dev/null 2>&1 || {
        echo "missing required command: $1" >&2
        exit 1
    }
}

need_cmd podman
need_cmd curl

state_file="$STATE_DIR/$(safe_name "$LOGICAL_NODE_ID").env"

if [ ! -f "$state_file" ]; then
    echo "pcrf-enable: node state not found: $state_file" >&2
    exit 1
fi

if [ ! -f "$NET_STATE" ]; then
    echo "pcrf-enable: network state not found: $NET_STATE" >&2
    exit 1
fi

# shellcheck disable=SC1090
. "$state_file"
# shellcheck disable=SC1090
. "$NET_STATE"

if [ -z "${CONTAINER_NAME:-}" ]; then
    echo "pcrf-enable: CONTAINER_NAME missing in $state_file" >&2
    exit 1
fi

if [ -z "${LAB_NET:-}" ]; then
    echo "pcrf-enable: LAB_NET missing in $NET_STATE" >&2
    exit 1
fi

node_ip="$(container_ip_on_network "$CONTAINER_NAME" "$LAB_NET")"
if [ -z "$node_ip" ]; then
    echo "pcrf-enable: unable to find IP for $CONTAINER_NAME on $LAB_NET" >&2
    podman inspect "$CONTAINER_NAME" >&2 || true
    exit 1
fi

url="http://${node_ip}:${PCRF_PORT}${PCRF_PATH}"

printf 'pcrf-enable: logical=%s factory=%s container=%s url=%s\n' \
    "$LOGICAL_NODE_ID" "${FACTORY_NODE_ID:-}" "$CONTAINER_NAME" "$url"

start_ts="$(date +%s)"

while :; do
    body_file="/tmp/ukama-pcrf-enable-body.$$"
    code="$(
        curl -sS \
            --max-time 5 \
            -o "$body_file" \
            -w '%{http_code}' \
            -X POST "$url" \
            -H 'Content-Type: application/json' \
            -d '{"state":"on"}' \
            2>/dev/null || true
    )"

    if [ "$code" -ge 200 ] 2>/dev/null && [ "$code" -lt 300 ] 2>/dev/null; then
        echo "pcrf-enable: ok http=$code"
        rm -f "$body_file"
        exit 0
    fi

    now_ts="$(date +%s)"
    elapsed=$((now_ts - start_ts))

    if [ "$elapsed" -ge "$TIMEOUT_SEC" ]; then
        echo "pcrf-enable: failed after ${TIMEOUT_SEC}s http=${code:-none}" >&2
        echo "---- response ----" >&2
        cat "$body_file" >&2 2>/dev/null || true
        echo >&2
        echo "---- pcrf status ----" >&2
        curl -sS --max-time 5 "http://${node_ip}:${PCRF_PORT}/v1/status" >&2 || true
        echo >&2
        rm -f "$body_file"
        exit 1
    fi

    rm -f "$body_file"
    sleep "$SLEEP_SEC"
done
