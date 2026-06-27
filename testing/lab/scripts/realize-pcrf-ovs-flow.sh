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

# Fallback only exists for lab debugging. Keep default off so we do not silently
# bypass PCRF metering. Set ULAB_OVS_ALLOW_UNMETERED_FALLBACK=1 only when you
# explicitly want to prove the remaining UE->media path even if meters fail.
ALLOW_UNMETERED_FALLBACK="${ULAB_OVS_ALLOW_UNMETERED_FALLBACK:-0}"

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
    printf '%s' "$1" | tr '\n' ' ' | \
        sed 's/^[[:space:]]*\[[[:space:]]*//' | \
        sed 's/[[:space:]]*\][[:space:]]*$//' | \
        sed 's/}[[:space:]]*,[[:space:]]*{/}\
{/g'
}

only_digits() {
    printf '%s' "$1" | tr -cd '0-9'
}

hex_cookie() {
    awk -v n="$1" 'BEGIN { printf "0x%x", n }'
}

ovs15() {
    podman exec "$TNODE_CONTAINER" ovs-ofctl -O OpenFlow15 "$@"
}

ovs_proto() {
    proto="$1"
    shift
    podman exec "$TNODE_CONTAINER" ovs-ofctl -O "$proto" "$@"
}

dump_ovs_state() {
    echo "---- OVS meters ----" >&2
    ovs15 dump-meters "$BRIDGE" >&2 || true
    echo "---- OVS flows ----" >&2
    ovs15 dump-flows "$BRIDGE" >&2 || true
}

meter_exists() {
    id="$1"
    ovs15 dump-meters "$BRIDGE" 2>/dev/null | grep -Eq "(meter=$id([^0-9]|$)|meter:$id([^0-9]|$)|^ *$id:)"
}

try_meter() {
    id="$1"
    rate="$2"
    proto="$3"
    spec="$4"

    echo "pcrf-ovs: try meter id=$id proto=$proto spec=$spec" >&2

    # Delete first so every retry is deterministic. del-meter returns failure
    # if the meter is absent; ignore that.
    ovs_proto "$proto" del-meter "$BRIDGE" "meter=$id" >/dev/null 2>&1 || true

    err="$(ovs_proto "$proto" add-meter "$BRIDGE" "$spec" 2>&1 >/dev/null)" && {
        meter_exists "$id" && return 0
        echo "pcrf-ovs: add-meter returned success but meter $id is not visible" >&2
        return 1
    }

    echo "pcrf-ovs: add-meter failed id=$id proto=$proto err=$err" >&2
    return 1
}

ensure_meter() {
    id="$1"
    rate="$2"

    id="$(only_digits "$id")"
    rate="$(only_digits "$rate")"

    if [ -z "$id" ]; then
        echo "invalid meter id: $id" >&2
        return 1
    fi

    # Keep rate sane and definitely non-empty. PCRF policy values have had
    # serialization differences in lab; this script is only realizing the meter
    # objects PCRF already assigned by ID.
    if [ -z "$rate" ] || [ "$rate" -le 0 ] 2>/dev/null; then
        rate=1000000
    fi

    # Try the forms used by common OVS versions. Official docs describe
    # add-meter switch meter, with bands after all other fields; examples in
    # the wild use both "band=" and "bands=" depending on version.
    for proto in OpenFlow13 OpenFlow15; do
        try_meter "$id" "$rate" "$proto" \
            "meter=$id,kbps,band=type=drop,rate=$rate" && return 0

        try_meter "$id" "$rate" "$proto" \
            "meter=$id,kbps,bands=type=drop,rate=$rate" && return 0

        try_meter "$id" "$rate" "$proto" \
            "meter=$id,kbps,stats,band=type=drop,rate=$rate" && return 0

        try_meter "$id" "$rate" "$proto" \
            "meter=$id,kbps,stats,bands=type=drop,rate=$rate" && return 0
    done

    echo "pcrf-ovs: failed to create OVS meter id=$id rate=$rate" >&2
    dump_ovs_state
    return 1
}

del_flow() {
    field="$1"
    ovs15 --strict del-flows "$BRIDGE" "table=0,priority=100,ip,$field=$UE_IP" \
        >/dev/null 2>&1 || true
}

add_flow_metered() {
    field="$1"
    cookie="$2"
    meter="$3"

    cookie_hex="$(hex_cookie "$cookie")"
    spec="cookie=$cookie_hex,table=0,priority=100,ip,$field=$UE_IP,actions=meter:$meter,NORMAL"

    echo "pcrf-ovs: add flow $spec" >&2
    del_flow "$field"
    ovs15 add-flow "$BRIDGE" "$spec" >/dev/null
}

add_flow_unmetered() {
    field="$1"
    cookie="$2"

    cookie_hex="$(hex_cookie "$cookie")"
    spec="cookie=$cookie_hex,table=0,priority=100,ip,$field=$UE_IP,actions=NORMAL"

    echo "pcrf-ovs: add UNMETERED LAB FALLBACK flow $spec" >&2
    del_flow "$field"
    ovs15 add-flow "$BRIDGE" "$spec" >/dev/null
}

verify_flow_metered() {
    field="$1"
    meter="$2"

    flows="$(ovs15 dump-flows "$BRIDGE" 2>/dev/null || true)"
    printf '%s\n' "$flows" | grep -q "priority=100" && \
    printf '%s\n' "$flows" | grep -q "$field=$UE_IP" && \
    printf '%s\n' "$flows" | grep -q "meter:$meter" && \
    printf '%s\n' "$flows" | grep -q "NORMAL"
}

verify_flow_unmetered() {
    field="$1"

    flows="$(ovs15 dump-flows "$BRIDGE" 2>/dev/null || true)"
    printf '%s\n' "$flows" | grep -q "priority=100" && \
    printf '%s\n' "$flows" | grep -q "$field=$UE_IP" && \
    printf '%s\n' "$flows" | grep -q "actions=NORMAL"
}

flow_json="$(podman exec "$TNODE_CONTAINER" \
    curl -fsS --max-time 5 "$PCRF_BASE/v1/subscriber/imsi/$IMSI/flow")"

rx_obj="$(json_objects "$flow_json" | sed -n '1p')"
tx_obj="$(json_objects "$flow_json" | sed -n '2p')"

if [ -z "$rx_obj" ] || [ -z "$tx_obj" ]; then
    echo "failed to parse two PCRF flow objects for IMSI $IMSI: $flow_json" >&2
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

rx_rate="${rx_rate:-1000000}"
tx_rate="${tx_rate:-1000000}"

echo "pcrf-ovs: flow records imsi=$IMSI ip=$UE_IP rx_cookie=$rx_cookie tx_cookie=$tx_cookie rx_meter=$rx_meter tx_meter=$tx_meter rx_rate=$rx_rate tx_rate=$tx_rate" >&2

metered=1
ensure_meter "$rx_meter" "$rx_rate" || metered=0
ensure_meter "$tx_meter" "$tx_rate" || metered=0

if [ "$metered" -eq 1 ]; then
    add_flow_metered "nw_dst" "$rx_cookie" "$rx_meter"
    add_flow_metered "nw_src" "$tx_cookie" "$tx_meter"

    if ! verify_flow_metered "nw_dst" "$rx_meter"; then
        echo "missing RX metered OVS flow for UE $UE_IP meter=$rx_meter" >&2
        dump_ovs_state
        exit 1
    fi

    if ! verify_flow_metered "nw_src" "$tx_meter"; then
        echo "missing TX metered OVS flow for UE $UE_IP meter=$tx_meter" >&2
        dump_ovs_state
        exit 1
    fi

    echo "pcrf-ovs-ready imsi=$IMSI ip=$UE_IP rx_meter=$rx_meter tx_meter=$tx_meter bridge=$BRIDGE"
    exit 0
fi

if [ "$ALLOW_UNMETERED_FALLBACK" = "1" ]; then
    echo "pcrf-ovs: WARNING using unmetered lab fallback because meter creation failed" >&2

    add_flow_unmetered "nw_dst" "$rx_cookie"
    add_flow_unmetered "nw_src" "$tx_cookie"

    verify_flow_unmetered "nw_dst" && verify_flow_unmetered "nw_src" || {
        echo "missing unmetered fallback OVS flow for UE $UE_IP" >&2
        dump_ovs_state
        exit 1
    }

    echo "pcrf-ovs-ready-unmetered imsi=$IMSI ip=$UE_IP bridge=$BRIDGE"
    exit 0
fi

echo "pcrf-ovs: meter creation failed and ULAB_OVS_ALLOW_UNMETERED_FALLBACK is not enabled" >&2
exit 1
