#!/usr/bin/env bash
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -euo pipefail

: "${TOWER_IP:?TOWER_IP required}"
: "${VNODE_NAME:?VNODE_NAME required, e.g. ukama-vnode}"

EPCEMU_PORT="${EPCEMU_PORT:-18092}"
PCRF_PORT="${PCRF_PORT:-18090}"
INITNET_PORT="${INITNET_PORT:-18091}"
TUN_IF="${TUN_IF:-tun3}"
BRIDGE_IF="${BRIDGE_IF:-br0}"
UE_CIDR="${UE_CIDR:-192.168.8.0/22}"
TUN_TABLE="${TUN_TABLE:-2000}"
BRIDGE_TABLE="${BRIDGE_TABLE:-1000}"
GATEWAY_IP="${GATEWAY_IP:-10.10.10.11}"
IMSI="${IMSI:-}"
UE_IP="${UE_IP:-}"

json_or_raw() {
    if command -v jq >/dev/null 2>&1; then
        jq .
    else
        cat
        echo
    fi
}

vexec() {
    podman exec "$VNODE_NAME" "$@"
}

require_match() {
    local title="$1"
    local pattern="$2"
    shift 2

    echo "== check: $title =="
    if ! "$@" | tee /tmp/ukama-validate.$$ | grep -F "$pattern" >/dev/null; then
        echo "missing: $pattern" >&2
        rm -f /tmp/ukama-validate.$$
        exit 1
    fi
    rm -f /tmp/ukama-validate.$$
}

show_api() {
    local name="$1"
    local url="$2"

    echo "== $name =="
    curl -fsS "$url" | json_or_raw
}

show_api "init-network" "http://${TOWER_IP}:${INITNET_PORT}/v1/status"
show_api "pcrf" "http://${TOWER_IP}:${PCRF_PORT}/v1/status"
show_api "epcemu" "http://${TOWER_IP}:${EPCEMU_PORT}/v1/status"

echo "== reconcile init-network =="
curl -fsS -X POST "http://${TOWER_IP}:${INITNET_PORT}/v1/reconcile" | json_or_raw

require_match \
    "ip rule ${TUN_IF} -> table ${TUN_TABLE}" \
    "iif ${TUN_IF} lookup ${TUN_TABLE}" \
    vexec ip rule show

require_match \
    "ip rule ${BRIDGE_IF} -> table ${BRIDGE_TABLE}" \
    "iif ${BRIDGE_IF} lookup ${BRIDGE_TABLE}" \
    vexec ip rule show

require_match \
    "table ${TUN_TABLE} default via gateway" \
    "default via ${GATEWAY_IP} dev ${BRIDGE_IF}" \
    vexec ip route show table "$TUN_TABLE"

require_match \
    "table ${BRIDGE_TABLE} UE CIDR to ${TUN_IF}" \
    "${UE_CIDR} dev ${TUN_IF}" \
    vexec ip route show table "$BRIDGE_TABLE"

require_match \
    "OVS default drop UE source" \
    "nw_src=${UE_CIDR} actions=drop" \
    vexec ovs-ofctl -O OpenFlow15 dump-flows "$BRIDGE_IF"

require_match \
    "OVS default drop UE destination" \
    "nw_dst=${UE_CIDR} actions=drop" \
    vexec ovs-ofctl -O OpenFlow15 dump-flows "$BRIDGE_IF"

if [[ -n "$UE_IP" ]]; then
    require_match \
        "OVS UE uplink flow ${UE_IP}" \
        "nw_src=${UE_IP}" \
        vexec ovs-ofctl -O OpenFlow15 dump-flows "$BRIDGE_IF"

    require_match \
        "OVS UE downlink flow ${UE_IP}" \
        "nw_dst=${UE_IP}" \
        vexec ovs-ofctl -O OpenFlow15 dump-flows "$BRIDGE_IF"
fi

if [[ -n "$IMSI" ]]; then
    echo "== epcemu UE list =="
    curl -fsS "http://${TOWER_IP}:${EPCEMU_PORT}/v1/ues" | json_or_raw || true

    echo "== pcrf flow for IMSI ${IMSI} =="
    curl -fsS "http://${TOWER_IP}:${PCRF_PORT}/v1/subscriber/imsi/${IMSI}/flow" \
        | json_or_raw || true
fi

echo "datapath wiring validation passed"
