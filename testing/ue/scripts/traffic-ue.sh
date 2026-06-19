#!/usr/bin/env bash
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -euo pipefail

if [[ $# -lt 4 ]]; then
    echo "usage: $0 <imsi> <ping|http|iperf> <mb> <media-ip>" >&2
    exit 2
fi

IMSI="$1"
MODE="$2"
MB="$3"
MEDIA_IP="$4"
HTTP_PORT=8080
IPERF_PORT=5201

case "$MODE" in
    ping)
        podman exec "ue-$IMSI" ping -c 5 "$MEDIA_IP"
        ;;
    http)
        podman exec "ue-$IMSI" curl -fsS "http://$MEDIA_IP:$HTTP_PORT/"
        ;;
    iperf)
        podman exec "ue-$IMSI" iperf3 -c "$MEDIA_IP" -p "$IPERF_PORT" -n "${MB}M"
        ;;
    *)
        echo "unknown mode: $MODE" >&2
        exit 1
        ;;
esac
