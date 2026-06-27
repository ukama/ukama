#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -eu

if [ "$#" -ne 3 ]; then
    echo "usage: $0 <ue-id-or-ref> <amount-mb> <run-dir>" >&2
    exit 2
fi

UE_KEY="$1"
AMOUNT_MB="$2"
RUN_DIR="$3"
STATE_FILE="$RUN_DIR/runtime-ues/$(printf "%s" "$UE_KEY" | tr -c 'A-Za-z0-9_.-' '-').env"

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

# Media is now on the tower user-plane bridge (for example 10.10.10.100).
# That IP is outside the UE subscriber subnet, so force only the media target
# through the UE tunnel. Without this route, iperf may use the UE container's
# podman eth0 path and never enter EPCEMU/PCRF.
podman exec "$UE_CONTAINER" \
    ip route replace "$MEDIA_IP/32" dev tun0

echo "traffic route ue=$UE_KEY media=$MEDIA_IP via=tun0"
podman exec "$UE_CONTAINER" ip route get "$MEDIA_IP" || true

echo "traffic precheck http ue=$UE_KEY media=$MEDIA_IP"
if ! MEDIA_IP="$MEDIA_IP" HTTP_PORT=8080 IPERF_PORT=5201 \
    "$UE_DIR/scripts/traffic-ue.sh" --imsi "$IMSI" --mode http; then
    echo "traffic HTTP precheck failed; dumping datapath state" >&2
    podman exec "$TNODE_CONTAINER" sh -lc 'curl -sS http://127.0.0.1:18028/v1/status || true' >&2 || true
    podman exec "$TNODE_CONTAINER" sh -lc 'curl -sS http://127.0.0.1:18030/v1/status || true' >&2 || true
    podman exec "$TNODE_CONTAINER" sh -lc 'ip rule; ip route show table 2000; ip route show table 1000; iptables -S FORWARD' >&2 || true
    podman exec "$TNODE_CONTAINER" sh -lc 'ovs-ofctl -O OpenFlow15 dump-meters br0 || true; ovs-ofctl -O OpenFlow15 dump-flows br0 || true' >&2 || true
    exit 1
fi

echo "traffic ue=$UE_KEY imsi=$IMSI mb=$AMOUNT_MB media=$MEDIA_IP"
MEDIA_IP="$MEDIA_IP" HTTP_PORT=8080 IPERF_PORT=5201 \
    "$UE_DIR/scripts/traffic-ue.sh" --imsi "$IMSI" --mode iperf --mb "$AMOUNT_MB"
echo "traffic-complete ue=$UE_KEY imsi=$IMSI mb=$AMOUNT_MB"
