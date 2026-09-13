#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
# Copyright (c) 2026-present, Ukama Inc.

set -eu

if [ "$#" -gt 2 ] || [ "${1:-}" = "--help" ]; then
    echo "usage: $0 [base-url|all] [interval-seconds]"
    echo "default: monitor all three nodes every 1 second; Ctrl-C to stop"
    exit 0
fi

TARGET="${1:-all}"
INTERVAL="${2:-1}"
case "$INTERVAL" in
    ''|*[!0-9]*|0) echo "interval must be a positive integer" >&2; exit 2 ;;
esac
for cmd in curl jq; do
    command -v "$cmd" >/dev/null 2>&1 || {
        echo "missing required command: $cmd" >&2
        exit 1
    }
done
trap 'exit 0' INT TERM

poll_node() {
    label="$1"
    url="$2"
    if ! body=$(curl -fsS --connect-timeout 1 --max-time 2 "${url%/}/v1/status" 2>/dev/null); then
        printf '%-10s UNREACHABLE  %s\n' "$label" "$url"
        return
    fi
    if ! summary=$(printf '%s' "$body" | jq -er '
        if (.state | type) != "string" then error("missing state") else
        "state=\(.state) config=\(.configuration.phase // "-")" +
        " mode=\(.configuration.mode // "-")" +
        " starter=\(.starter.readiness // "-")" +
        " notifyPending=\(if has("notificationPending") then .notificationPending else "-" end)" +
        " request=\(.configuration.requestId // "-")" +
        " reason=\(.reason // "-")"
        end' 2>/dev/null); then
        printf '%-10s INVALID RESPONSE  %s\n' "$label" "$url"
        return
    fi
    printf '%-10s %s\n' "$label" "$summary"
}

while :; do
    date '+%Y-%m-%d %H:%M:%S %Z'
    if [ "$TARGET" = "all" ]; then
        poll_node tower "http://127.0.0.1:${ULAB_LIFECYCLE_TOWER_HOST_PORT:-18033}"
        poll_node controller "http://127.0.0.1:${ULAB_LIFECYCLE_CONTROLLER_HOST_PORT:-18034}"
        poll_node amplifier "http://127.0.0.1:${ULAB_LIFECYCLE_AMPLIFIER_HOST_PORT:-18035}"
    else
        poll_node node "$TARGET"
    fi
    printf '\n'
    sleep "$INTERVAL"
done
