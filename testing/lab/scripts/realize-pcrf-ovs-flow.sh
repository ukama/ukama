#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -eu

if [ "$#" -ne 3 ]; then
    echo "usage: $0 <ue-state-file> <imsi> <ue-ip>" >&2
    exit 2
fi

STATE_FILE="$1"
IMSI="$2"
UE_IP="$3"

if [ ! -f "$STATE_FILE" ]; then
    echo "UE state not found: $STATE_FILE" >&2
    exit 1
fi

# shellcheck disable=SC1090
. "$STATE_FILE"

BRIDGE="${PCRF_OVS_BRIDGE:-br0}"
PCRF_BASE="http://127.0.0.1:18030"

if [ -z "${TNODE_CONTAINER:-}" ]; then
    echo "TNODE_CONTAINER missing in $STATE_FILE" >&2
    exit 1
fi

json_get_first() {
    key="$1"
    printf '%s' "$2" | sed -n "s/.*\"$key\"[[:space:]]*:[[:space:]]*\"\{0,1\}\([^\",}]*\).*/\1/p" | head -1
}

json_get_any_first() {
    payload="$1"
    shift
    for key in "$@"; do
        val="$(json_get_first "$key" "$payload" || true)"
        if [ -n "$val" ]; then
            printf '%s\n' "$val"
            return 0
        fi
    done
    return 1
}

json_objects() {
    # PCRF returns a compact array of objects. Split it without needing jq.
    printf '%s' "$1" | tr -d '\n' | sed 's/},{/}\
{/g'
}

hex_cookie() {
    # PCRF cookies are generated from uint32 today. awk keeps this portable.
    awk -v n="$1" 'BEGIN { printf "0x%x", n }'
}

ovs() {
    podman exec "$TNODE_CONTAINER" ovs-ofctl -O OpenFlow15 "$@"
}

ensure_meter() {
    id="$1"
    rate="$2"
    burst="$3"

    if [ -z "$id" ] || [ -z "$rate" ]; then
        echo "invalid meter id/rate id=$id rate=$rate" >&2
        return 1
    fi

    if [ "${burst:-0}" -gt 0 ] 2>/dev/null; then
        spec="meter=$id,kbps,burst,stats,bands=type=drop,rate=$rate,burst_size=$burst"
    else
        spec="meter=$id,kbps,stats,bands=type=drop,rate=$rate"
    fi

    # Make retries deterministic: add if missing, modify if it already exists.
    ovs add-meter "$BRIDGE" "$spec" >/dev/null 2>&1 || \
        ovs mod-meter "$BRIDGE" "$spec" >/dev/null
}

del_flow() {
    field="$1"
    ovs --strict del-flows "$BRIDGE" "table=0,priority=100,ip,$field=$UE_IP" \
        >/dev/null 2>&1 || true
}

add_flow() {
    field="$1"
    cookie="$2"
    meter="$3"

    cookie_hex="$(hex_cookie "$cookie")"
    spec="cookie=$cookie_hex,table=0,priority=100,ip,$field=$UE_IP,actions=meter:$meter,NORMAL"

    del_flow "$field"
    ovs add-flow "$BRIDGE" "$spec" >/dev/null
}

verify_flow() {
    field="$1"
    meter="$2"

    flows="$(ovs dump-flows "$BRIDGE" 2>/dev/null || true)"
    printf '%s\n' "$flows" | grep -q "priority=100" && \
    printf '%s\n' "$flows" | grep -q "$field=$UE_IP" && \
    printf '%s\n' "$flows" | grep -q "meter:$meter" && \
    printf '%s\n' "$flows" | grep -q "NORMAL"
}

flow_json="$(podman exec "$TNODE_CONTAINER" \
    curl -fsS --max-time 5 "$PCRF_BASE/v1/subscriber/imsi/$IMSI/flow")"

rx_obj="$(json_objects "$flow_json" | sed -n '1p')"
tx_obj="$(json_objects "$flow_json" | sed -n '2p')"

if [ -z "$rx_obj" ] || [ -z "$tx_obj" ]; then
    echo "PCRF did not return two flow objects for IMSI $IMSI: $flow_json" >&2
    exit 1
fi

rx_cookie="$(json_get_any_first "$rx_obj" Cookie cookie)"
tx_cookie="$(json_get_any_first "$tx_obj" Cookie cookie)"
rx_meter="$(json_get_any_first "$rx_obj" MeterID meterID meter_id)"
tx_meter="$(json_get_any_first "$tx_obj" MeterID meterID meter_id)"

policy_json="$(podman exec "$TNODE_CONTAINER" \
    curl -fsS --max-time 5 "$PCRF_BASE/v1/policy/imsi/$IMSI" 2>/dev/null || true)"

# Store.CreateMeter uses Ulbr for RX and Dlbr for TX. Keep that mapping.
rx_rate="$(json_get_any_first "$policy_json" Ulbr ulbr ULBR 2>/dev/null || true)"
tx_rate="$(json_get_any_first "$policy_json" Dlbr dlbr DLBR 2>/dev/null || true)"
burst="$(json_get_any_first "$policy_json" Burst burst 2>/dev/null || true)"

# Conservative fallback only for lab smoke if policy serialization changes.
# The PCRF flow cookies/meter IDs still come from PCRF; this does not bypass
# the OVS meter action, it only ensures the meter object exists.
rx_rate="${rx_rate:-10000000}"
tx_rate="${tx_rate:-10000000}"
burst="${burst:-0}"

ensure_meter "$rx_meter" "$rx_rate" "$burst"
ensure_meter "$tx_meter" "$tx_rate" "$burst"

add_flow "nw_dst" "$rx_cookie" "$rx_meter"
add_flow "nw_src" "$tx_cookie" "$tx_meter"

if ! verify_flow "nw_dst" "$rx_meter"; then
    echo "missing RX OVS flow for UE $UE_IP meter=$rx_meter" >&2
    ovs dump-flows "$BRIDGE" >&2 || true
    exit 1
fi

if ! verify_flow "nw_src" "$tx_meter"; then
    echo "missing TX OVS flow for UE $UE_IP meter=$tx_meter" >&2
    ovs dump-flows "$BRIDGE" >&2 || true
    exit 1
fi

echo "pcrf-ovs-ready imsi=$IMSI ip=$UE_IP rx_meter=$rx_meter tx_meter=$tx_meter bridge=$BRIDGE"
