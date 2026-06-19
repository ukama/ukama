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
SITE_STATE="$RUN_DIR/runtime-sites/$(printf "%s" "$SITE_REF" | tr -c 'A-Za-z0-9_.-' '-').env"
MEDIA_STATE="$RUN_DIR/runtime-media/media.env"
UE_STATE_DIR="$RUN_DIR/runtime-ues"
UE_IMAGE="ukama/ue:dev"

mkdir -p "$UE_STATE_DIR"

if ! command -v podman >/dev/null 2>&1; then
    echo "podman is required" >&2
    exit 1
fi

if [ ! -x "$UE_DIR/scripts/run-ue.sh" ]; then
    echo "missing $UE_DIR/scripts/run-ue.sh" >&2
    exit 1
fi

if [ ! -f "$SITE_STATE" ]; then
    echo "site state not found: $SITE_STATE" >&2
    exit 1
fi

if [ ! -f "$MEDIA_STATE" ]; then
    echo "media state not found: $MEDIA_STATE" >&2
    exit 1
fi

. "$SITE_STATE"
. "$MEDIA_STATE"

if [ -z "${TNODE_CONTAINER:-}" ]; then
    echo "TNODE_CONTAINER missing in $SITE_STATE" >&2
    exit 1
fi

TOWER_IP="$(podman inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' "$TNODE_CONTAINER")"
if [ -z "$TOWER_IP" ]; then
    echo "tower container has no IP: $TNODE_CONTAINER" >&2
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

UE_IMAGE="$UE_IMAGE" \
UE_NETWORK_MODE=podman \
TOWER_IP="$TOWER_IP" \
MEDIA_IP="$MEDIA_IP" \
EPCEMU_PORT=18028 \
EPCEMU_DATA_PORT=18029 \
PCRF_PORT=18030 \
UE_DATA_HOST=0.0.0.0 \
ALLOW_LOCAL_MEDIA=true \
    "$UE_DIR/scripts/run-ue.sh" --csv "$CSV" --imsi "$IMSI" > "$UE_STATE_DIR/$UE_ID.start.out"

UE_CONTAINER="ue-$IMSI"
STATE_FILE="$UE_STATE_DIR/$UE_ID.env"
cat > "$STATE_FILE" <<STATE
UE_REF=$UE_REF
UE_ID=$UE_ID
UE_CONTAINER=$UE_CONTAINER
IMSI=$IMSI
ICCID=$ICCID
UE_IP=$UE_IP
SITE_REF=$SITE_REF
TNODE_CONTAINER=$TNODE_CONTAINER
TOWER_IP=$TOWER_IP
MEDIA_IP=$MEDIA_IP
HTTP_PORT=8080
IPERF_PORT=5201
CSV=$CSV
UE_DIR=$UE_DIR
STATE

cp "$STATE_FILE" "$UE_STATE_DIR/$UE_REF.env"

echo "ue-started ref=$UE_REF imsi=$IMSI ip=$UE_IP tower=$TNODE_CONTAINER media=$MEDIA_IP"
