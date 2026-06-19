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

UE_IMAGE="${UE_IMAGE:-ukama/ue:dev}"
UE_NETWORK_MODE="${UE_NETWORK_MODE:-container}"
VNODE_NAME="${VNODE_NAME:-ukama-vnode}"

TOWER_IP="${TOWER_IP:-127.0.0.1}"
MEDIA_IP="${MEDIA_IP:-127.0.0.1}"

EPCEMU_PORT="${EPCEMU_PORT:-18028}"
EPCEMU_DATA_PORT="${EPCEMU_DATA_PORT:-18029}"
PCRF_PORT="${PCRF_PORT:-18030}"

UE_DATA_HOST="${UE_DATA_HOST:-127.0.0.1}"
UE_BASE_PORT="${UE_BASE_PORT:-41000}"
UE_TUN="${UE_TUN:-tun0}"
UE_DETACH_ON_EXIT="${UE_DETACH_ON_EXIT:-1}"

ALLOW_LOCAL_MEDIA="${ALLOW_LOCAL_MEDIA:-true}"

usage() {
    cat <<USAGE
Usage: $0 --csv <csv-file> --imsi <imsi>

Environment:
  UE_IMAGE             UE image tag. Default: ukama/ue:dev
  UE_NETWORK_MODE      container|podman. Default: container
  VNODE_NAME           Virtual-node container name. Default: ukama-vnode

  TOWER_IP             Tower/API IP. Default: 127.0.0.1
  MEDIA_IP             Media target IP. Default: 127.0.0.1

  EPCEMU_PORT          Default: 18028
  EPCEMU_DATA_PORT     Default: 18029
  PCRF_PORT            Default: 18030

  UE_DATA_HOST         UE UDP return host. Default: 127.0.0.1
  UE_BASE_PORT         Default: 41000
  UE_TUN               Default: tun0
  UE_DETACH_ON_EXIT    Default: 1

  ALLOW_LOCAL_MEDIA    Default: true for lab mode

Examples:
  $0 --csv testing/ue/csv/SimPool.with-imsi.csv --imsi 001010000000001

  UE_IMAGE=localhost/ukama/ue:dev \\
  VNODE_NAME=ukama-vnode \\
  $0 --csv testing/ue/csv/SimPool.with-imsi.csv --imsi 001010000000001
USAGE
}

is_local_media_ip() {
    local ip="$1"

    case "$ip" in
        127.*|0.0.0.0|localhost)
            return 0
            ;;
        *)
            return 1
            ;;
    esac
}

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
        -h|--help)
            usage
            exit 0
            ;;
        *)
            echo "unknown arg $1" >&2
            usage >&2
            exit 1
            ;;
    esac
done

: "${CSV:?--csv required}"
: "${IMSI:?--imsi required}"

if is_local_media_ip "$MEDIA_IP" && [[ "$ALLOW_LOCAL_MEDIA" != "true" ]]; then
    echo "MEDIA_IP=$MEDIA_IP is local; real E2E requires external media" >&2
    echo "set ALLOW_LOCAL_MEDIA=true only for temporary lab testing" >&2
    exit 1
fi

row="$(csv_row_by_imsi "$CSV" "$IMSI")"
if [[ -z "$row" ]]; then
    echo "IMSI not found: $IMSI" >&2
    exit 1
fi

ICCID="$(csv_field "$CSV" "$row" ICCID)"
UE_IP="$(csv_field "$CSV" "$row" UE_IP)"
APN="$(csv_field "$CSV" "$row" APN)"

idx="${IMSI: -3}"
UE_DATA_PORT=$((UE_BASE_PORT + 10#$idx))
NAME="ue-${IMSI}"

podman rm -f "$NAME" >/dev/null 2>&1 || true

if [[ "$UE_NETWORK_MODE" == "container" ]]; then
    podman run -d \
        --name "$NAME" \
        --network "container:${VNODE_NAME}" \
        --privileged \
        --device /dev/net/tun \
        -e UE_IMSI="$IMSI" \
        -e UE_ICCID="$ICCID" \
        -e UE_IP="${UE_IP}/22" \
        -e UE_APN="$APN" \
        -e UE_TUN="$UE_TUN" \
        -e EPCEMU_URL="http://127.0.0.1:${EPCEMU_PORT}" \
        -e EPCEMU_DATA_HOST="127.0.0.1" \
        -e EPCEMU_DATA_PORT="$EPCEMU_DATA_PORT" \
        -e UE_DATA_HOST="127.0.0.1" \
        -e UE_DATA_PORT="$UE_DATA_PORT" \
        -e PCRF_URL="http://127.0.0.1:${PCRF_PORT}" \
        -e MEDIA_IP="$MEDIA_IP" \
        -e UE_DETACH_ON_EXIT="$UE_DETACH_ON_EXIT" \
        "$UE_IMAGE" /bin/sh -c '/opt/ukama/ue-agent/ue-agent || exit 1; tail -f /dev/null'
else
    podman run -d \
        --name "$NAME" \
        --cap-add NET_ADMIN \
        --device /dev/net/tun \
        -p "${UE_DATA_PORT}:${UE_DATA_PORT}/udp" \
        -e UE_IMSI="$IMSI" \
        -e UE_ICCID="$ICCID" \
        -e UE_IP="${UE_IP}/22" \
        -e UE_APN="$APN" \
        -e UE_TUN="$UE_TUN" \
        -e EPCEMU_URL="http://${TOWER_IP}:${EPCEMU_PORT}" \
        -e EPCEMU_DATA_HOST="$TOWER_IP" \
        -e EPCEMU_DATA_PORT="$EPCEMU_DATA_PORT" \
        -e UE_DATA_HOST="$UE_DATA_HOST" \
        -e UE_DATA_PORT="$UE_DATA_PORT" \
        -e PCRF_URL="http://${TOWER_IP}:${PCRF_PORT}" \
        -e MEDIA_IP="$MEDIA_IP" \
        -e UE_DETACH_ON_EXIT="$UE_DETACH_ON_EXIT" \
        "$UE_IMAGE" /bin/sh -c '/opt/ukama/ue-agent/ue-agent || exit 1; tail -f /dev/null'
fi

echo "$NAME $IMSI $ICCID $UE_IP port=$UE_DATA_PORT mode=$UE_NETWORK_MODE"
