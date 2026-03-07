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

wait_for_json_condition() {
    local url="$1"
    local pyexpr="$2"
    local timeout="${3:-20}"
    local start now body
    start="$(date +%s)"
    while true; do
        if body="$(curl -fsS "$url" 2>/dev/null || true)"; then
            if [[ -n "$body" ]]; then
                if python3 - "$pyexpr" <<'PY' <<<"$body" >/dev/null 2>&1
import json, sys
expr = sys.argv[1]
data = json.load(sys.stdin)
if eval(expr, {"__builtins__": {}}, {"data": data}):
    raise SystemExit(0)
raise SystemExit(1)
PY
                then
                    return 0
                fi
            fi
        fi
        now="$(date +%s)"
        if (( now - start >= timeout )); then
            echo "timeout waiting for JSON condition on: $url" >&2
            echo "last body: ${body:-<none>}" >&2
            return 1
        fi
        sleep 0.25
    done
}

assert_json_condition() {
    local url="$1"
    local pyexpr="$2"
    local body
    body="$(curl -fsS "$url")"
    python3 - "$pyexpr" <<'PY' <<<"$body"
import json, sys
expr = sys.argv[1]
data = json.load(sys.stdin)
assert eval(expr, {"__builtins__": {}}, {"data": data}), json.dumps(data, indent=2)
PY
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
