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

. "$STATE_FILE"

if [ ! -x "$UE_DIR/scripts/traffic-ue.sh" ]; then
    echo "missing $UE_DIR/scripts/traffic-ue.sh" >&2
    exit 1
fi

echo "traffic ue=$UE_KEY imsi=$IMSI mb=$AMOUNT_MB media=$MEDIA_IP"
MEDIA_IP="$MEDIA_IP" HTTP_PORT=8080 IPERF_PORT=5201 \
    "$UE_DIR/scripts/traffic-ue.sh" --imsi "$IMSI" --mode iperf --mb "$AMOUNT_MB"
echo "traffic-complete ue=$UE_KEY imsi=$IMSI mb=$AMOUNT_MB"
