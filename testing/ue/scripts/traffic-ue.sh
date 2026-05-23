#!/usr/bin/env bash
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -euo pipefail

: "${MEDIA_IP:?MEDIA_IP required}"

IMSI=""
MODE="ping"
ALLOW_LOCAL_MEDIA="${ALLOW_LOCAL_MEDIA:-false}"
HTTP_PORT="${HTTP_PORT:-8080}"
IPERF_PORT="${IPERF_PORT:-5201}"

usage() {
    cat <<USAGE
Usage: $0 --imsi <imsi> [--mode ping|http|iperf]

Environment:
  MEDIA_IP      external media/sink server IP
  HTTP_PORT     default: 8080
  IPERF_PORT    default: 5201
  IPERF_TIME    default: 10
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
        --imsi) IMSI="$2"; shift 2 ;;
        --mode) MODE="$2"; shift 2 ;;
        -h|--help) usage; exit 0 ;;
        *) echo "unknown arg $1" >&2; usage >&2; exit 1 ;;
    esac
done

: "${IMSI:?--imsi required}"

if is_local_media_ip "$MEDIA_IP" && [[ "$ALLOW_LOCAL_MEDIA" != "true" ]]; then
    echo "MEDIA_IP=$MEDIA_IP is local; real E2E requires external media" >&2
    echo "set ALLOW_LOCAL_MEDIA=true only for temporary lab testing" >&2
    exit 1
fi

case "$MODE" in
    ping)
        podman exec "ue-${IMSI}" ping -c "${PING_COUNT:-5}" "$MEDIA_IP"
        ;;
    http)
        podman exec "ue-${IMSI}" curl -fsS "http://${MEDIA_IP}:${HTTP_PORT}/"
        ;;
    iperf)
        podman exec "ue-${IMSI}" iperf3 -c "$MEDIA_IP" \
            -p "$IPERF_PORT" -t "${IPERF_TIME:-10}"
        ;;
    *)
        echo "unknown mode $MODE" >&2
        exit 1
        ;;
esac
