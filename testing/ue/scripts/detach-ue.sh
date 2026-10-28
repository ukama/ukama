#!/usr/bin/env bash
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -euo pipefail

IMSI=""

while [[ $# -gt 0 ]]; do
    case "$1" in
        --imsi)
            IMSI="$2"
            shift 2
            ;;
        -h|--help)
            echo "Usage: $0 --imsi <imsi>"
            exit 0
            ;;
        *)
            echo "unknown arg $1" >&2
            exit 1
            ;;
    esac
done

: "${IMSI:?--imsi required}"

UE_CONTAINER="ue-${IMSI}"
status="$(podman exec "$UE_CONTAINER" cat /run/ukama-ue/status.json)"
if ! printf '%s' "$status" | jq -e --arg imsi "$IMSI" \
    '.source == "ue-agent" and .imsi == $imsi' >/dev/null; then
    echo "UE status does not match the requested IMSI" >&2
    exit 1
fi

podman kill --signal USR1 "$UE_CONTAINER" >/dev/null
for ((i=0; i<60; i++)); do
    status="$(podman exec "$UE_CONTAINER" \
        cat /run/ukama-ue/status.json)"
    epc="$(podman exec "$UE_CONTAINER" sh -c '
        curl -sS --max-time 5 -o /dev/null -w "%{http_code}" \
            "${EPCEMU_URL}/v1/ue/${UE_IMSI}"' || true)"
    pcrf="$(podman exec "$UE_CONTAINER" sh -c '
        curl -sS --max-time 5 -o /dev/null -w "%{http_code}" \
            "${PCRF_URL}/v1/pdn/imsi/${UE_IMSI}"' || true)"
    if printf '%s' "$status" | jq -e '.state == "detached"' >/dev/null &&
       [[ "$epc" == "404" && "$pcrf" == "404" ]]; then
        podman rm -f "$UE_CONTAINER" >/dev/null
        exit 0
    fi
    sleep 1
done

echo "PDN cleanup is pending; leaving $UE_CONTAINER for retry" >&2
exit 1
