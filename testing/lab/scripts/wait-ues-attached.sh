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
TIMEOUT_SEC="${ULAB_UE_ATTACH_TIMEOUT:-90}"
SLEEP_SEC=3

if [ ! -f "$STATE_FILE" ]; then
    echo "UE state not found: $STATE_FILE" >&2
    exit 1
fi

# shellcheck disable=SC1090
. "$STATE_FILE"

pcrf_subscriber_exists() {
    podman exec "$TNODE_CONTAINER" \
        curl -fsS --max-time 2 "http://127.0.0.1:18030/v1/subscriber/imsi/$IMSI" \
        >/dev/null 2>&1
}

ue_is_running() {
    podman inspect -f '{{.State.Running}}' "$UE_CONTAINER" 2>/dev/null | grep -q '^true$'
}

epc_ue_attached() {
    podman exec "$TNODE_CONTAINER" \
        curl -fsS --max-time 2 "http://127.0.0.1:18028/v1/ue/$IMSI" 2>/dev/null | \
        grep -qi '"state"[[:space:]]*:[[:space:]]*"attached"'
}

pcrf_flow_ready() {
    podman exec "$TNODE_CONTAINER" \
        curl -fsS --max-time 2 "http://127.0.0.1:18030/v1/subscriber/imsi/$IMSI/flow" 2>/dev/null | \
        grep -q '[{}\[]'
}

realize_pcrf_ovs_flow() {
    script_dir="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
    "$script_dir/realize-pcrf-ovs-flow.sh" "$STATE_FILE" "$IMSI" "$UE_IP"
}

ovs_flow_ready() {
    podman exec "$TNODE_CONTAINER" \
        ovs-ofctl -O OpenFlow15 dump-flows br0 2>/dev/null | \
        grep -q "priority=100.*nw_src=$UE_IP.*NORMAL" && \
    podman exec "$TNODE_CONTAINER" \
        ovs-ofctl -O OpenFlow15 dump-flows br0 2>/dev/null | \
        grep -q "priority=100.*nw_dst=$UE_IP.*NORMAL"
}

start_ts="$(date +%s)"
while :; do
    ue_running=0
    epc_attached=0
    pcrf_subscriber=0
    pcrf_flow=0
    ovs_flow=0

    ue_is_running && ue_running=1
    pcrf_subscriber_exists && pcrf_subscriber=1
    epc_ue_attached && epc_attached=1
    pcrf_flow_ready && pcrf_flow=1

    if [ "$ue_running" -eq 1 ] && \
       [ "$pcrf_subscriber" -eq 1 ] && \
       [ "$epc_attached" -eq 1 ] && \
       [ "$pcrf_flow" -eq 1 ]; then
        if ovs_flow_ready; then
            ovs_flow=1
        else
            realize_pcrf_ovs_flow && ovs_flow=1 || ovs_flow=0
        fi
    fi

    if [ "$ue_running" -eq 1 ] && \
       [ "$pcrf_subscriber" -eq 1 ] && \
       [ "$epc_attached" -eq 1 ] && \
       [ "$pcrf_flow" -eq 1 ] && \
       [ "$ovs_flow" -eq 1 ]; then
        echo "ue-attached ue=$UE_KEY imsi=$IMSI ip=$UE_IP ovs_flow=ready"
        exit 0
    fi

    now_ts="$(date +%s)"
    if [ $((now_ts - start_ts)) -ge "$TIMEOUT_SEC" ]; then
        echo "UE not attached: ue=$UE_KEY imsi=$IMSI" >&2
        echo "ue_running=$ue_running pcrf_subscriber=$pcrf_subscriber epc_attached=$epc_attached pcrf_flow=$pcrf_flow ovs_flow=$ovs_flow" >&2
        echo "---- UE container ----" >&2
        podman ps -a --filter "name=$UE_CONTAINER" >&2 || true
        echo "---- UE logs ----" >&2
        podman logs --tail 120 "$UE_CONTAINER" >&2 || true
        echo "---- EPCEMU UE ----" >&2
        podman exec "$TNODE_CONTAINER" \
            curl -i --max-time 5 "http://127.0.0.1:18028/v1/ue/$IMSI" >&2 || true
        echo >&2
        echo "---- PCRF subscriber ----" >&2
        podman exec "$TNODE_CONTAINER" \
            curl -i --max-time 5 "http://127.0.0.1:18030/v1/subscriber/imsi/$IMSI" >&2 || true
        echo >&2
        echo "---- PCRF subscriber flow ----" >&2
        podman exec "$TNODE_CONTAINER" \
            curl -i --max-time 5 "http://127.0.0.1:18030/v1/subscriber/imsi/$IMSI/flow" >&2 || true
        echo >&2
        echo "---- EPCEMU status ----" >&2
        podman exec "$TNODE_CONTAINER" \
            curl -fsS --max-time 5 "http://127.0.0.1:18028/v1/status" >&2 || true
        echo >&2
        echo "---- PCRF status ----" >&2
        podman exec "$TNODE_CONTAINER" \
            curl -fsS --max-time 5 "http://127.0.0.1:18030/v1/status" >&2 || true
        echo >&2
        echo "---- OVS meters ----" >&2
        podman exec "$TNODE_CONTAINER" \
            ovs-ofctl -O OpenFlow15 dump-meters br0 >&2 || true
        echo >&2
        echo "---- OVS flows ----" >&2
        podman exec "$TNODE_CONTAINER" \
            ovs-ofctl -O OpenFlow15 dump-flows br0 >&2 || true
        echo >&2
        exit 1
    fi

    sleep "$SLEEP_SEC"
done
