#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -eu

if [ "$#" -ne 3 ] && { [ "$#" -ne 4 ] || [ "$4" != exhaust ]; }; then
    echo "usage: $0 <ue-id-or-ref> <amount-mb> <run-dir> [exhaust]" >&2
    exit 2
fi

UE_KEY="$1"
AMOUNT_MB="$2"
RUN_DIR="$3"
TRAFFIC_MODE="${4:-}"
STATE_FILE="$RUN_DIR/runtime-ues/$(printf "%s" "$UE_KEY" | tr -c 'A-Za-z0-9_.-' '-').env"
ATTEMPT_TIMEOUT_SEC="${ULAB_TRAFFIC_ATTEMPT_TIMEOUT_SEC:-60}"
EXHAUST_GRACE_SEC="${ULAB_TRAFFIC_EXHAUST_GRACE_SEC:-10}"

case "$ATTEMPT_TIMEOUT_SEC" in
    ''|*[!0-9]*)
        echo "invalid ULAB_TRAFFIC_ATTEMPT_TIMEOUT_SEC: $ATTEMPT_TIMEOUT_SEC" >&2
        exit 2
        ;;
esac

case "$EXHAUST_GRACE_SEC" in
    ''|*[!0-9]*)
        echo "invalid ULAB_TRAFFIC_EXHAUST_GRACE_SEC: $EXHAUST_GRACE_SEC" >&2
        exit 2
        ;;
esac

case "$AMOUNT_MB" in
    ''|*[!0-9]*)
        echo "invalid amount-mb: $AMOUNT_MB" >&2
        exit 2
        ;;
esac

if [ "$ATTEMPT_TIMEOUT_SEC" -eq 0 ]; then
    echo "ULAB_TRAFFIC_ATTEMPT_TIMEOUT_SEC must be greater than zero" >&2
    exit 2
fi

if [ ! -f "$STATE_FILE" ]; then
    echo "UE state not found: $STATE_FILE" >&2
    exit 1
fi

# shellcheck disable=SC1090
. "$STATE_FILE"

if [ ! -x "$UE_DIR/scripts/traffic-ue.sh" ]; then
    echo "missing $UE_DIR/scripts/traffic-ue.sh" >&2
    exit 1
fi

if [ -z "${UE_CONTAINER:-}" ]; then
    echo "UE_CONTAINER missing in $STATE_FILE" >&2
    exit 1
fi

if [ -z "${MEDIA_IP:-}" ]; then
    echo "MEDIA_IP missing in $STATE_FILE" >&2
    exit 1
fi

if ! podman inspect -f '{{.State.Running}}' "$UE_CONTAINER" 2>/dev/null | grep -q '^true$'; then
    echo "UE container is not running: $UE_CONTAINER" >&2
    podman ps -a --filter "name=$UE_CONTAINER" >&2 || true
    podman logs --tail 120 "$UE_CONTAINER" >&2 || true
    exit 1
fi

if ! podman exec "$UE_CONTAINER" test -d /sys/class/net/tun0; then
    echo "UE tun0 not found in $UE_CONTAINER" >&2
    podman exec "$UE_CONTAINER" ip addr >&2 || true
    podman logs --tail 120 "$UE_CONTAINER" >&2 || true
    exit 1
fi

# Force only the media target through the UE tunnel. Without this route,
# traffic may use the UE container's podman eth0 path and never enter
# EPCEMU/PCRF.
podman exec "$UE_CONTAINER" \
    ip route replace "$MEDIA_IP/32" dev tun0

echo "traffic route ue=$UE_KEY media=$MEDIA_IP via=tun0"
podman exec "$UE_CONTAINER" ip route get "$MEDIA_IP" || true

echo "traffic ue=$UE_KEY imsi=$IMSI mb=$AMOUNT_MB media=$MEDIA_IP"

# Sum the Transfer column of iperf3's interval lines, in bytes.
sent_bytes() {
    [ -f "$1" ] || { echo 0; return 0; }
    awk '/ sec / && !/sender|receiver/ {
        for (i = 1; i < NF; i++) if ($i == "sec") break
        v = $(i + 1); u = $(i + 2); m = 0
        if (u == "Bytes") m = 1
        if (u == "KBytes") m = 1024
        if (u == "MBytes") m = 1048576
        if (u == "GBytes") m = 1073741824
        t += v * m
    } END { printf "%.0f\n", t }' "$1"
}

media_rx_bytes() {
    podman exec "$MEDIA_CONTAINER" cat /sys/class/net/eth0/statistics/rx_bytes
}

restart_iperf_endpoints() {
    podman exec "$UE_CONTAINER" pkill -x iperf3 >/dev/null 2>&1 || true
    podman exec "$MEDIA_CONTAINER" sh -c \
        'pkill -x iperf3; sleep 1; iperf3 -s -D --forceflush' >/dev/null 2>&1 || {
        echo "traffic: could not restart iperf3 server in $MEDIA_CONTAINER" >&2
        return 1
    }
}

# The user spends the whole allowance, so the PCRF may close the PDN before
# iperf3 exchanges results. Pass on what the backend delivered to media: at
# least 95% of the allowance, since ACKs also count against the quota.
exhaust_traffic() {
    out="$RUN_DIR/runtime-ues/traffic-$IMSI.out"
    rc_file="$out.rc"
    want=$((AMOUNT_MB * 1048576 * 99 / 100))
    allowance=$((AMOUNT_MB * 1048576))
    need_delivered=$((allowance * 95 / 100))
    rm -f "$out" "$rc_file"
    rx_before="$(media_rx_bytes)" || {
        echo "traffic: cannot read rx_bytes in $MEDIA_CONTAINER" >&2
        return 1
    }
    (
        if timeout -k 5s "${ATTEMPT_TIMEOUT_SEC}s" \
            env MEDIA_IP="$MEDIA_IP" HTTP_PORT=8080 IPERF_PORT=5201 \
            "$UE_DIR/scripts/traffic-ue.sh" --imsi "$IMSI" --mode iperf \
            --mb "$AMOUNT_MB" >"$out" 2>&1; then
            echo 0 >"$rc_file"
        else
            echo "$?" >"$rc_file"
        fi
    ) &
    sent_at=""
    while [ ! -s "$rc_file" ]; do
        if [ -z "$sent_at" ] && [ "$(sent_bytes "$out")" -ge "$want" ]; then
            sent_at="$(date +%s)"
        fi
        if [ -n "$sent_at" ] && [ $(( $(date +%s) - sent_at )) -ge "$EXHAUST_GRACE_SEC" ]; then
            podman exec "$UE_CONTAINER" pkill -x iperf3 >/dev/null 2>&1 || true
        fi
        sleep 1
    done
    wait || true
    cat "$out"
    traffic_rc="$(cat "$rc_file")"
    sent="$(sent_bytes "$out")"
    rm -f "$out" "$rc_file"
    if [ "$traffic_rc" -eq 0 ]; then
        echo "traffic-complete ue=$UE_KEY imsi=$IMSI mb=$AMOUNT_MB"
        return 0
    fi
    rx_after="$(media_rx_bytes)" || {
        echo "traffic: cannot read rx_bytes in $MEDIA_CONTAINER" >&2
        return 1
    }
    delivered=$((rx_after - rx_before))
    if [ "$delivered" -lt "$need_delivered" ]; then
        echo "traffic cut before the allowance was used rc=$traffic_rc" \
             "delivered_bytes=$delivered need_bytes=$need_delivered sent_bytes=$sent" >&2
        return 1
    fi
    restart_iperf_endpoints || return 1
    echo "traffic-exhausted ue=$UE_KEY imsi=$IMSI mb=$AMOUNT_MB" \
         "delivered_bytes=$delivered allowance_bytes=$allowance sent_bytes=$sent iperf_rc=$traffic_rc"
}

if [ "$TRAFFIC_MODE" = exhaust ]; then
    if [ -z "${MEDIA_CONTAINER:-}" ]; then
        echo "MEDIA_CONTAINER missing in $STATE_FILE" >&2
        exit 1
    fi
    exhaust_traffic
    exit $?
fi

if timeout -k 5s "${ATTEMPT_TIMEOUT_SEC}s" \
    env MEDIA_IP="$MEDIA_IP" HTTP_PORT=8080 IPERF_PORT=5201 \
    "$UE_DIR/scripts/traffic-ue.sh" --imsi "$IMSI" --mode iperf \
    --mb "$AMOUNT_MB"; then
    :
else
    traffic_rc=$?
    if [ "$traffic_rc" -eq 124 ] || [ "$traffic_rc" -eq 137 ]; then
        echo "traffic transfer timed out after ${ATTEMPT_TIMEOUT_SEC}s" >&2
    else
        echo "traffic transfer failed rc=$traffic_rc" >&2
    fi
    exit 1
fi
echo "traffic-complete ue=$UE_KEY imsi=$IMSI mb=$AMOUNT_MB"
