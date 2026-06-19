#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -eu

if [ "$#" -ne 2 ]; then
    echo "usage: $0 <ue-id-or-ref> <run-dir>" >&2
    exit 2
fi

UE_KEY="$1"
RUN_DIR="$2"
STATE_FILE="$RUN_DIR/runtime-ues/$(printf "%s" "$UE_KEY" | tr -c 'A-Za-z0-9_.-' '-').env"
TIMEOUT_SEC=90
SLEEP_SEC=3

if [ ! -f "$STATE_FILE" ]; then
    echo "UE state not found: $STATE_FILE" >&2
    exit 1
fi

. "$STATE_FILE"

start_ts="$(date +%s)"
while :; do
    ue_running=0
    epc_attached=0
    pcrf_ready=0

    podman inspect -f '{{.State.Running}}' "$UE_CONTAINER" 2>/dev/null | grep -q '^true$' && ue_running=1

    podman exec "$TNODE_CONTAINER" \
        curl -fsS --max-time 2 "http://127.0.0.1:18028/v1/ue/$IMSI" 2>/dev/null | \
        grep -qi '"state"[[:space:]]*:[[:space:]]*"attached"' && epc_attached=1

    podman exec "$TNODE_CONTAINER" \
        curl -fsS --max-time 2 "http://127.0.0.1:18030/v1/subscriber/imsi/$IMSI/flow" 2>/dev/null | \
        grep -q '[{}\[]' && pcrf_ready=1

    if [ "$ue_running" -eq 1 ] && [ "$epc_attached" -eq 1 ] && [ "$pcrf_ready" -eq 1 ]; then
        echo "ue-attached ue=$UE_KEY imsi=$IMSI ip=$UE_IP"
        exit 0
    fi

    now_ts="$(date +%s)"
    if [ $((now_ts - start_ts)) -ge "$TIMEOUT_SEC" ]; then
        echo "UE not attached: ue=$UE_KEY imsi=$IMSI" >&2
        echo "ue_running=$ue_running epc_attached=$epc_attached pcrf_ready=$pcrf_ready" >&2
        podman ps -a --filter "name=$UE_CONTAINER" >&2 || true
        podman logs --tail 120 "$UE_CONTAINER" >&2 || true
        podman exec "$TNODE_CONTAINER" curl -fsS "http://127.0.0.1:18028/v1/ue/$IMSI" >&2 || true
        echo >&2
        podman exec "$TNODE_CONTAINER" curl -fsS "http://127.0.0.1:18030/v1/subscriber/imsi/$IMSI/flow" >&2 || true
        echo >&2
        exit 1
    fi

    sleep "$SLEEP_SEC"
done
