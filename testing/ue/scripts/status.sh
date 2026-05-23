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
VNODE_NAME="${VNODE_NAME:-}"

json_or_raw() {
    if command -v jq >/dev/null 2>&1; then
        jq .
    else
        cat
        echo
    fi
}

show_container_cmd() {
    local title="$1"
    shift

    if [[ -z "$VNODE_NAME" ]]; then
        return 0
    fi

    echo "== $title =="
    if ! podman exec "$VNODE_NAME" "$@"; then
        echo "failed: podman exec $VNODE_NAME $*" >&2
        return 1
    fi
}

echo "== init-network =="
curl -fsS "http://${TOWER_IP}:${INITNET_PORT}/v1/status" | json_or_raw

echo "== pcrf =="
curl -fsS "http://${TOWER_IP}:${PCRF_PORT}/v1/status" | json_or_raw

echo "== epcemu =="
curl -fsS "http://${TOWER_IP}:${EPCEMU_PORT}/v1/status" | json_or_raw

if [[ -n "$VNODE_NAME" ]]; then
    show_container_cmd "ip rule" ip rule show || true
    show_container_cmd "route table 1000" ip route show table 1000 || true
    show_container_cmd "route table 2000" ip route show table 2000 || true
    show_container_cmd "ovs flows" ovs-ofctl -O OpenFlow15 dump-flows br0 || true
    show_container_cmd "ovs meters" ovs-ofctl -O OpenFlow15 dump-meters br0 || true
fi
