#!/usr/bin/env bash
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.
#
set -euo pipefail

require_cmd() {
    local c
    for c in "$@"; do
        command -v "$c" >/dev/null 2>&1 || {
            echo "missing required command: $c" >&2
            exit 1
        }
    done
}

wait_for_http_ok() {
    local url="$1"
    local timeout="${2:-20}"
    local start now
    start="$(date +%s)"
    while true; do
        if curl -fsS "$url" >/dev/null 2>&1; then
            return 0
        fi
        now="$(date +%s)"
        if (( now - start >= timeout )); then
            echo "timeout waiting for HTTP OK: $url" >&2
            return 1
        fi
        sleep 0.2
    done
}

wait_for_http_down() {
    local url="$1"
    local timeout="${2:-20}"
    local start now
    start="$(date +%s)"
    while true; do
        if ! curl -fsS "$url" >/dev/null 2>&1; then
            return 0
        fi
        now="$(date +%s)"
        if (( now - start >= timeout )); then
            echo "timeout waiting for HTTP down: $url" >&2
            return 1
        fi
        sleep 0.2
    done
}

wait_for_file() {
    local path="$1"
    local timeout="${2:-20}"
    local start now
    start="$(date +%s)"
    while true; do
        if [[ -e "$path" ]]; then
            return 0
        fi
        now="$(date +%s)"
        if (( now - start >= timeout )); then
            echo "timeout waiting for file: $path" >&2
            return 1
        fi
        sleep 0.2
    done
}

wait_for_command_ok() {
    local timeout="${1:-20}"
    shift

    local start now
    start="$(date +%s)"
    while true; do
        if bash -o pipefail -c "$*" >/dev/null 2>&1; then
            return 0
        fi
        now="$(date +%s)"
        if (( now - start >= timeout )); then
            echo "timeout waiting for command: $*" >&2
            return 1
        fi
        sleep 0.2
    done
}

wait_for_json_condition() {
    local url="${1:?url required}"
    local pyexpr="${2:?python expression required}"
    local timeout_sec="${3:-20}"
    local start_ts
    local now
    local body

    start_ts="$(date +%s)"
    while true; do
        body="$(curl -fsS "$url" 2>/dev/null || true)"
        if [[ -n "$body" ]]; then
            if JSON_BODY="$body" python3 -c '
import json
import os
import sys

expr = sys.argv[1]
data = json.loads(os.environ["JSON_BODY"])

safe_builtins = {
    "any": any,
    "all": all,
    "len": len,
    "min": min,
    "max": max,
    "sum": sum,
    "sorted": sorted,
}

if eval(expr, {"__builtins__": safe_builtins}, {"data": data}):
    raise SystemExit(0)
raise SystemExit(1)
' "$pyexpr" >/dev/null 2>&1
            then
                return 0
            fi
        fi

        now="$(date +%s)"
        if (( now - start_ts >= timeout_sec )); then
            echo "timeout waiting for JSON condition on: $url" >&2
            echo "last body: ${body:-<empty>}" >&2
            return 1
        fi

        sleep 1
    done
}

assert_json_condition() {
    local url="${1:?url required}"
    local pyexpr="${2:?python expression required}"
    local body

    body="$(curl -fsS "$url")"

    JSON_BODY="$body" python3 -c '
import json
import os
import sys

expr = sys.argv[1]
data = json.loads(os.environ["JSON_BODY"])

safe_builtins = {
    "any": any,
    "all": all,
    "len": len,
    "min": min,
    "max": max,
    "sum": sum,
    "sorted": sorted,
}

assert eval(expr, {"__builtins__": safe_builtins}, {"data": data}), json.dumps(data, indent=2)
' "$pyexpr"
}

find_space_app() {
    local status_url="$1"
    local space="$2"
    local app="$3"
    curl -fsS "$status_url" | python3 - "$space" "$app" <<'PY'
import json, sys
space_name = sys.argv[1]
app_name = sys.argv[2]
data = json.load(sys.stdin)
for space in data.get("spaces", []):
    if space.get("name") != space_name:
        continue
    for app in space.get("apps", []):
        if app.get("name") == app_name:
            print(json.dumps(app))
            raise SystemExit(0)
raise SystemExit(1)
PY
}
