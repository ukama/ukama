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

while [[ $# -gt 0 ]]; do
    case "$1" in
        --imsi) IMSI="$2"; shift 2 ;;
        --mode) MODE="$2"; shift 2 ;;
        *) echo "unknown arg $1" >&2; exit 1 ;;
    esac
done

: "${IMSI:?--imsi required}"

case "$MODE" in
    ping)
        podman exec "ue-${IMSI}" ping -c 5 "$MEDIA_IP"
        ;;
    http)
        podman exec "ue-${IMSI}" curl -fsS "http://${MEDIA_IP}:8080/"
        ;;
    iperf)
        podman exec "ue-${IMSI}" iperf3 -c "$MEDIA_IP" -t "${IPERF_TIME:-10}"
        ;;
    *)
        echo "unknown mode $MODE" >&2
        exit 1
        ;;
esac
