#!/usr/bin/env bash
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -euo pipefail

if [[ $# -ne 4 ]]; then
    echo "usage: $0 <csv> <imsi> <tower-ip> <media-ip>" >&2
    exit 2
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/csv-lib.sh"

CSV="$1"
IMSI="$2"
TOWER_IP="$3"
MEDIA_IP="$4"
UE_IMAGE="ukama/ue:dev"
EPCEMU_PORT=18028
EPCEMU_DATA_PORT=18029
PCRF_PORT=18030
UE_BASE_PORT=41000

row="$(csv_row_by_imsi "$CSV" "$IMSI")"
if [[ -z "$row" ]]; then
    echo "IMSI not found: $IMSI" >&2
    exit 1
fi

ICCID="$(csv_field "$CSV" "$row" ICCID)"
UE_IP="$(csv_field "$CSV" "$row" UE_IP)"
APN="$(csv_field "$CSV" "$row" APN)"
[[ -n "$APN" ]] || APN="internet"

idx="${IMSI: -3}"
UE_DATA_PORT=$((UE_BASE_PORT + 10#$idx))
NAME="ue-$IMSI"

podman rm -f "$NAME" >/dev/null 2>&1 || true
podman run -d \
    --name "$NAME" \
    --cap-add NET_ADMIN \
    --device /dev/net/tun \
    -e UE_IMSI="$IMSI" \
    -e UE_ICCID="$ICCID" \
    -e UE_IP="$UE_IP/22" \
    -e UE_APN="$APN" \
    -e UE_TUN="tun0" \
    -e EPCEMU_URL="http://$TOWER_IP:$EPCEMU_PORT" \
    -e EPCEMU_DATA_HOST="$TOWER_IP" \
    -e EPCEMU_DATA_PORT="$EPCEMU_DATA_PORT" \
    -e UE_DATA_HOST="0.0.0.0" \
    -e UE_DATA_PORT="$UE_DATA_PORT" \
    -e PCRF_URL="http://$TOWER_IP:$PCRF_PORT" \
    -e MEDIA_IP="$MEDIA_IP" \
    -e UE_DETACH_ON_EXIT="1" \
    "$UE_IMAGE" /bin/sh -c '/opt/ukama/ue-agent/ue-agent || exit 1; tail -f /dev/null' >/dev/null

echo "$NAME imsi=$IMSI ip=$UE_IP tower=$TOWER_IP media=$MEDIA_IP port=$UE_DATA_PORT"
