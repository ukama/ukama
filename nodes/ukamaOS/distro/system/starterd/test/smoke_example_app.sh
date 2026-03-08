#!/usr/bin/env bash
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.
#
set -euo pipefail

SCRIPT_DIR="$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)"
# shellcheck source=common.sh
source "$SCRIPT_DIR/common.sh"

require_cmd bash curl python3 tar mktemp

STARTERD_BIN="${STARTERD_BIN:-${1:-starter.d}}"
if ! command -v "$STARTERD_BIN" >/dev/null 2>&1 && [[ ! -x "$STARTERD_BIN" ]]; then
    echo "starter.d binary not found or not executable: $STARTERD_BIN" >&2
    exit 1
fi

TMP="$(mktemp -d /tmp/starterd-smoke.XXXXXX)"
export TMP
TEST_ROOT="$TMP/root"
PKG_REPO="$TMP/pkgrepo"
LOG_DIR="$TMP/logs"
mkdir -p "$TEST_ROOT" "$PKG_REPO" "$LOG_DIR"

STARTER_HOST="127.0.0.1"
STARTER_PORT="18001"
WIMC_HOST="127.0.0.1"
WIMC_PORT="18006"
APP_PORT="18110"
MANIFEST="$TMP/manifest.json"
READY_FILE="$TMP/starter.ready"
STATUS_URL="http://${STARTER_HOST}:${STARTER_PORT}/v1/status"
APP_URL="http://127.0.0.1:${APP_PORT}"
WIMC_LOG="$LOG_DIR/mock_wimc.log"
STARTER_LOG="$LOG_DIR/starterd.log"

show_debug() {
    set +e

    echo "===== starter.d /v1/status =====" >&2
    curl -fsS "$STATUS_URL" >&2 || true
    echo >&2

    echo "===== starter.d stdout.log =====" >&2
    [[ -f "$LOG_DIR/stdout.log" ]] && tail -n 200 "$LOG_DIR/stdout.log" >&2 || true

    echo "===== starterd.log =====" >&2
    [[ -f "$STARTER_LOG" ]] && tail -n 200 "$STARTER_LOG" >&2 || true

    echo "===== mock_wimc.log =====" >&2
    [[ -f "$WIMC_LOG" ]] && tail -n 200 "$WIMC_LOG" >&2 || true
}

cleanup() {
    local rc=$?

    if (( rc != 0 )); then
        show_debug
    fi

    set +e
    [[ -n "${STARTER_PID:-}" ]] && kill "$STARTER_PID" >/dev/null 2>&1 || true
    [[ -n "${WIMC_PID:-}" ]] && kill "$WIMC_PID" >/dev/null 2>&1 || true
    wait "${STARTER_PID:-}" >/dev/null 2>&1 || true
    wait "${WIMC_PID:-}" >/dev/null 2>&1 || true
    rm -rf "$TMP"

    exit "$rc"
}
trap cleanup EXIT

make_example_pkg() {
    local version="$1"
    local pkgdir="$TMP/pkg-$version"

    mkdir -p "$pkgdir"
    cat > "$pkgdir/example_app.py" <<PY
#!/usr/bin/env python3
import json
import os
import signal
from http.server import BaseHTTPRequestHandler, HTTPServer

VERSION = ${version@Q}
PORT = int(os.environ.get("APP_PORT", "${APP_PORT}"))
httpd = None

class Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path == "/v1/ping":
            self.send_response(200)
            self.end_headers()
            self.wfile.write(b"OK")
            return
        if self.path == "/v1/version":
            self.send_response(200)
            self.end_headers()
            self.wfile.write(VERSION.encode())
            return
        if self.path == "/v1/status":
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            self.end_headers()
            self.wfile.write(json.dumps({"version": VERSION, "port": PORT}).encode())
            return
        self.send_response(404)
        self.end_headers()

    def log_message(self, fmt, *args):
        pass

def on_term(signum, frame):
    global httpd
    if httpd is not None:
        httpd.shutdown()

signal.signal(signal.SIGTERM, on_term)
signal.signal(signal.SIGINT, on_term)
httpd = HTTPServer(("127.0.0.1", PORT), Handler)
httpd.serve_forever()
PY
    chmod +x "$pkgdir/example_app.py"
    tar -czf "$PKG_REPO/example_app-${version}.tar.gz" \
        -C "$pkgdir" \
        example_app.py
}

wait_for_example_version() {
    local expected_version="$1"

    wait_for_http_ok "$APP_URL/v1/ping" 20

    wait_for_json_condition \
        "$STATUS_URL" \
        "any(
            space.get('name') == 'boot' and
            any(
                app.get('name') == 'example_app' and
                (
                    app.get('tag') == '${expected_version}' or
                    app.get('lastGoodTag') == '${expected_version}'
                )
                for app in space.get('apps', [])
            )
            for space in data.get('spaces', [])
        )" \
        20

    wait_for_command_ok 20 curl -fsS "$APP_URL/v1/version" \
        '|' grep -qx "$expected_version"
}

make_example_pkg v1
make_example_pkg v2

cat > "$MANIFEST" <<JSON
{
  "spaces": [
    {
      "name": "boot",
      "apps": [
        {
          "name": "example_app",
          "tag": "v1",
          "cmd": "example_app.py",
          "port": ${APP_PORT},
          "env": {
            "APP_PORT": "${APP_PORT}"
          }
        }
      ]
    }
  ]
}
JSON

python3 "$SCRIPT_DIR/mock_wimc.py" \
    --host "$WIMC_HOST" \
    --port "$WIMC_PORT" \
    --repo "$PKG_REPO" \
    >"$WIMC_LOG" 2>&1 &
WIMC_PID=$!
wait_for_http_ok "http://${WIMC_HOST}:${WIMC_PORT}/does-not-exist" 1 >/dev/null 2>&1 || true
sleep 0.2

export STARTERD_MANIFEST="$MANIFEST"
export STARTERD_LOG_PATH="$STARTER_LOG"
export STARTERD_READY_FILE="$READY_FILE"
export STARTERD_APPS_ROOT="$TEST_ROOT/apps"
export STARTERD_PKGS_DIR="$TEST_ROOT/pkgs"
export STARTERD_STATE_DIR="$TEST_ROOT/state"
export STARTERD_HTTP_ADDR="$STARTER_HOST"
export STARTERD_HTTP_PORT="$STARTER_PORT"
export STARTERD_WIMC_HOST="$WIMC_HOST"
export STARTERD_WIMC_PORT="$WIMC_PORT"
export STARTERD_WIMC_PATH_TEMPLATE="/v1/apps/%s/%s/pkg"
export STARTERD_COMMIT_TIMEOUT_SEC=10
export STARTERD_PING_TIMEOUT_SEC=2
export STARTERD_TERM_GRACE_SEC=3
export STARTERD_RESTART_MAX_BACKOFF_SEC=2
export STARTERD_RESTART_STABLE_RESET_SEC=5
export STARTERD_BOOT_SPACE=boot
export STARTERD_LOG_LEVEL=debug

"$STARTERD_BIN" >"$LOG_DIR/stdout.log" 2>&1 &
STARTER_PID=$!

wait_for_http_ok "http://${STARTER_HOST}:${STARTER_PORT}/v1/ping" 20
wait_for_file "$READY_FILE" 20
wait_for_example_version v1

echo "[ok] example_app booted with v1"

curl -fsS \
    -X POST \
    "http://${STARTER_HOST}:${STARTER_PORT}/v1/update" \
    -H 'Content-Type: application/json' \
    -d '{
        "space":"boot",
        "name":"example_app",
        "tag":"v2"
    }' \
    >/dev/null

wait_for_example_version v2
echo "[ok] example_app updated to v2 and is reachable"

sleep 5
curl -X GET \
     "http://${STARTER_HOST}:${STARTER_PORT}/v1/status"

echo "smoke_example_app: PASS"
