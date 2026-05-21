#!/usr/bin/env bash
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/csv-lib.sh"

CSV=""
IMSI=""

while [[ $# -gt 0 ]]; do
    case "$1" in
        --csv)
            CSV="$2"
            shift 2
            ;;
        --imsi)
            IMSI="$2"
            shift 2
            ;;
        *)
            echo "unknown arg $1" >&2
            exit 1
            ;;
    esac
done

: "${CSV:?--csv required}"
: "${IMSI:?--imsi required}"
: "${TOWER_IP:?TOWER_IP required}"
: "${MEDIA_IP:?MEDIA_IP required}"

UE_IMAGE="${UE_IMAGE:-ukama/ue:dev}"
EPCEMU_PORT="${EPCEMU_PORT:-18092}"
EPCEMU_DATA_PORT="${EPCEMU_DATA_PORT:-18110}"
UE_DATA_HOST="${UE_DATA_HOST:-$(hostname -I | awk '{print $1}')}"
BASE_PORT="${UE_BASE_PORT:-41000}"

row="$(csv_row_by_imsi "$CSV" "$IMSI")"
if [[ -z "$row" ]]; then
    echo "IMSI not found: $IMSI" >&2
    exit 1
fi

ICCID="$(csv_field "$CSV" "$row" ICCID)"
UE_IP="$(csv_field "$CSV" "$row" UE_IP)"
APN="$(csv_field "$CSV" "$row" APN)"
idx="${IMSI: -3}"
UE_DATA_PORT=$((BASE_PORT + 10#$idx))
NAME="ue-${IMSI}"

podman rm -f "$NAME" >/dev/null 2>&1 || true
podman run -d --name "$NAME" \
    --cap-add NET_ADMIN \
    --device /dev/net/tun \
    -p "${UE_DATA_PORT}:${UE_DATA_PORT}/udp" \
    -e UE_IMSI="$IMSI" \
    -e UE_ICCID="$ICCID" \
    -e UE_IP="${UE_IP}/22" \
    -e UE_APN="$APN" \
    -e EPCEMU_URL="http://${TOWER_IP}:${EPCEMU_PORT}" \
    -e EPCEMU_DATA_HOST="$TOWER_IP" \
    -e EPCEMU_DATA_PORT="$EPCEMU_DATA_PORT" \
    -e UE_DATA_HOST="$UE_DATA_HOST" \
    -e UE_DATA_PORT="$UE_DATA_PORT" \
    -e MEDIA_IP="$MEDIA_IP" \
    "$UE_IMAGE"

echo "$NAME $IMSI $ICCID $UE_IP $UE_DATA_HOST:$UE_DATA_PORT"
