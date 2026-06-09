#!/usr/bin/env bash
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -euo pipefail

TOWER_IP="${TOWER_IP:-127.0.0.1}"
EPCEMU_PORT="${EPCEMU_PORT:-18028}"
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

curl -fsS -X DELETE "http://${TOWER_IP}:${EPCEMU_PORT}/v1/ue/detach" \
    -H 'Content-Type: application/json' \
    -d "{\"imsi\":\"${IMSI}\"}"

for i in $(seq 1 20); do
    epcAttached="$(curl -fsS "http://${TOWER_IP}:${EPCEMU_PORT}/v1/status" |
        jq -r '.ues.attached')"

    pcrfActive="$(curl -fsS "http://${TOWER_IP}:${PCRF_PORT:-18030}/v1/status" |
        jq -r '.sessions.active')"

    if [[ "$epcAttached" == "0" && "$pcrfActive" == "0" ]]; then
        break
    fi

    sleep 1
done

podman rm -f "ue-${IMSI}" >/dev/null 2>&1 || true
