#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -eu

if [ "$#" -ne 8 ]; then
    echo "usage: $0 <repo> <ue-ref> <ue-id> <imsi> <iccid> <ue-ip> <site-ref> <run-dir>" >&2
    exit 2
fi

REPO="$1"
UE_REF="$2"
UE_ID="$3"
IMSI="$4"
ICCID="$5"
UE_IP="$6"
SITE_REF="$7"
RUN_DIR="$8"

UE_DIR="$REPO/testing/ue"
NET_STATE="$RUN_DIR/runtime-net/net.env"
SITE_STATE="$RUN_DIR/runtime-sites/$(printf "%s" "$SITE_REF" | tr -c 'A-Za-z0-9_.-' '-').env"
MEDIA_STATE="$RUN_DIR/runtime-media/media.env"
UE_STATE_DIR="$RUN_DIR/runtime-ues"
UE_IMAGE="ukama/ue:dev"
UE_CONTAINER="ue-$IMSI"

mkdir -p "$UE_STATE_DIR"

need_cmd() {
    if ! command -v "$1" >/dev/null 2>&1; then
        echo "missing required command: $1" >&2
        exit 1
    fi
}

need_file() {
    if [ ! -f "$1" ]; then
        echo "missing $1" >&2
        exit 1
    fi
}

container_ip_on_network() {
    container="$1"
    network="$2"

    podman inspect -f '{{range $name, $net := .NetworkSettings.Networks}}{{if eq $name "'"$network"'"}}{{$net.IPAddress}}{{end}}{{end}}' \
        "$container" 2>/dev/null
}

ue_data_port() {
    last3="$(printf "%s" "$IMSI" | sed 's/.*\(...\)$/\1/')"
    awk -v x="$last3" 'BEGIN { printf "%d", 41000 + x }'
}

need_cmd podman
need_cmd awk
need_file "$UE_DIR/ue/Containerfile"
need_file "$NET_STATE"
need_file "$SITE_STATE"
need_file "$MEDIA_STATE"

# shellcheck disable=SC1090
. "$NET_STATE"
# shellcheck disable=SC1090
. "$SITE_STATE"
# shellcheck disable=SC1090
. "$MEDIA_STATE"

if [ -z "${LAB_NET:-}" ]; then
    echo "LAB_NET missing in $NET_STATE" >&2
    exit 1
fi

if [ -z "${TNODE_CONTAINER:-}" ]; then
    echo "TNODE_CONTAINER missing in $SITE_STATE" >&2
    exit 1
fi

TOWER_IP="$(container_ip_on_network "$TNODE_CONTAINER" "$LAB_NET")"
if [ -z "$TOWER_IP" ]; then
    echo "tower container has no IP on $LAB_NET: $TNODE_CONTAINER" >&2
    podman inspect "$TNODE_CONTAINER" >&2 || true
    exit 1
fi

if [ -z "${MEDIA_IP:-}" ]; then
    echo "MEDIA_IP missing in $MEDIA_STATE" >&2
    exit 1
fi

if [ ! -f "$UE_STATE_DIR/.ue-image-built" ]; then
    echo "ue: build $UE_IMAGE"
    podman build -t "$UE_IMAGE" -f "$UE_DIR/ue/Containerfile" "$UE_DIR"
    touch "$UE_STATE_DIR/.ue-image-built"
fi

CSV="$UE_STATE_DIR/$UE_ID.csv"
cat > "$CSV" <<CSVEOF
IMSI,ICCID,MSISDN,SmDpAddress,ActivationCode,IsPhysical,QrCode,UE_IP,APN,Enabled
$IMSI,$ICCID,10000000000,lab,lab,TRUE,lab,$UE_IP,internet,TRUE
CSVEOF

UE_DATA_PORT="$(ue_data_port)"

echo "ue: start $UE_CONTAINER network=$LAB_NET tower=$TOWER_IP media=$MEDIA_IP ip=$UE_IP"
podman rm -f "$UE_CONTAINER" >/dev/null 2>&1 || true

podman run -d \
    --name "$UE_CONTAINER" \
    --network "$LAB_NET" \
    --cap-add NET_ADMIN \
    --device /dev/net/tun \
    -e UE_IMSI="$IMSI" \
    -e UE_ICCID="$ICCID" \
    -e UE_IP="$UE_IP/22" \
    -e UE_APN="internet" \
    -e UE_TUN="tun0" \
    -e EPCEMU_URL="http://$TOWER_IP:18028" \
    -e EPCEMU_DATA_HOST="$TOWER_IP" \
    -e EPCEMU_DATA_PORT="18029" \
    -e UE_DATA_HOST="0.0.0.0" \
    -e UE_DATA_PORT="$UE_DATA_PORT" \
    -e PCRF_URL="http://$TOWER_IP:18030" \
    -e MEDIA_IP="$MEDIA_IP" \
    -e UE_DETACH_ON_EXIT="1" \
    "$UE_IMAGE" \
    /bin/sh -c '/opt/ukama/ue-agent/ue-agent || exit 1; tail -f /dev/null' \
    > "$UE_STATE_DIR/$UE_ID.start.out"

STATE_FILE="$UE_STATE_DIR/$UE_ID.env"
cat > "$STATE_FILE" <<STATE
UE_REF=$UE_REF
UE_ID=$UE_ID
UE_CONTAINER=$UE_CONTAINER
IMSI=$IMSI
ICCID=$ICCID
UE_IP=$UE_IP
UE_DATA_PORT=$UE_DATA_PORT
SITE_REF=$SITE_REF
TNODE_CONTAINER=$TNODE_CONTAINER
TOWER_IP=$TOWER_IP
MEDIA_IP=$MEDIA_IP
HTTP_PORT=8080
IPERF_PORT=5201
LAB_NET=$LAB_NET
CSV=$CSV
UE_DIR=$UE_DIR
STATE

cp "$STATE_FILE" "$UE_STATE_DIR/$UE_REF.env"

echo "ue-started ref=$UE_REF imsi=$IMSI ip=$UE_IP tower=$TOWER_IP media=$MEDIA_IP network=$LAB_NET"
