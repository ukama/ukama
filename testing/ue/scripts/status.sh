#!/usr/bin/env bash
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -euo pipefail

: "${TOWER_IP:?TOWER_IP required}"

EPCEMU_PORT="${EPCEMU_PORT:-18092}"
PCRF_PORT="${PCRF_PORT:-18090}"
INITNET_PORT="${INITNET_PORT:-18091}"

json_or_raw() {
    if command -v jq >/dev/null 2>&1; then
        jq .
    else
        cat
        echo
    fi
}

echo "== init-network =="
curl -fsS "http://${TOWER_IP}:${INITNET_PORT}/v1/status" | json_or_raw

echo "== pcrf =="
curl -fsS "http://${TOWER_IP}:${PCRF_PORT}/v1/status" | json_or_raw

echo "== epcemu =="
curl -fsS "http://${TOWER_IP}:${EPCEMU_PORT}/v1/status" | json_or_raw
