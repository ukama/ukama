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
STATE_FILE="$RUN_DIR/runtime-ues/$(printf "%s" "$UE_KEY" | \
    tr -c 'A-Za-z0-9_.-' '-').env"
TIMEOUT_SEC="${ULAB_UE_DETACH_TIMEOUT:-30}"
SLEEP_SEC=1

if [ ! -f "$STATE_FILE" ]; then
    echo "UE state not found: $STATE_FILE" >&2
    exit 1
fi

# shellcheck disable=SC1090
. "$STATE_FILE"

ue_is_running() {
    podman inspect -f '{{.State.Running}}' "$UE_CONTAINER" 2>/dev/null | \
        grep -q '^true$'
}

epc_ue_absent() {
    code="$(podman exec "$TNODE_CONTAINER" \
        curl -sS -o /dev/null -w '%{http_code}' --max-time 2 \
        "http://127.0.0.1:18028/v1/ue/$IMSI" 2>/dev/null || true)"
    [ "$code" = "404" ]
}

pcrf_session_absent() {
    podman exec "$TNODE_CONTAINER" \
        curl -fsS --max-time 2 "http://127.0.0.1:18030/v1/status" \
        2>/dev/null | grep -q '"active"[[:space:]]*:[[:space:]]*0'
}

ovs_flow_absent() {
    ! podman exec "$TNODE_CONTAINER" \
        ovs-ofctl -O OpenFlow15 dump-flows br0 2>/dev/null | \
        grep -q "priority=100.*nw_src=$UE_IP.*NORMAL" && \
    ! podman exec "$TNODE_CONTAINER" \
        ovs-ofctl -O OpenFlow15 dump-flows br0 2>/dev/null | \
        grep -q "priority=100.*nw_dst=$UE_IP.*NORMAL"
}

start_ts="$(date +%s)"
while :; do
    running=0
    epc_absent=0
    pcrf_absent=0
    ovs_absent=0

    ue_is_running && running=1
    epc_ue_absent && epc_absent=1
    pcrf_session_absent && pcrf_absent=1
    ovs_flow_absent && ovs_absent=1

    if [ "$running" -eq 1 ] && \
       [ "$epc_absent" -eq 1 ] && \
       [ "$pcrf_absent" -eq 1 ] && \
       [ "$ovs_absent" -eq 1 ]; then
        echo "ue-detached ue=$UE_KEY imsi=$IMSI ip=$UE_IP"
        exit 0
    fi

    now_ts="$(date +%s)"
    if [ $((now_ts - start_ts)) -ge "$TIMEOUT_SEC" ]; then
        echo "UE not detached: ue=$UE_KEY imsi=$IMSI" >&2
        echo "running=$running epc_absent=$epc_absent "\
"pcrf_absent=$pcrf_absent ovs_absent=$ovs_absent" >&2
        exit 1
    fi

    sleep "$SLEEP_SEC"
done
