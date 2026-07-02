#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
CTRL_BIN="${AISG_CTRL_BIN:-$ROOT/ctrl/aisg-ctrl}"
TTY="${AISG_TTY:-${1:-/dev/ttyUSB0}}"
WORKDIR="${AISG_TEST_DIR:-$(mktemp -d /tmp/aisg-real.XXXXXX)}"
SOCK="$WORKDIR/aisg-ctrl.sock"
CTRL_CFG="$WORKDIR/aisg-ctrl.toml"
CTRL_LOG="$WORKDIR/aisg-ctrl.log"
CONFIG_BLOB="${AISG_CONFIG_BLOB:-}"
MOVE="${AISG_MOVE:-0}"
KEEP_LOGS="${AISG_KEEP_LOGS:-1}"
CTRL_PID=""

cleanup() {
    set +e
    if [[ -n "$CTRL_PID" ]]; then kill "$CTRL_PID" 2>/dev/null || true; fi
    wait "$CTRL_PID" 2>/dev/null || true
    if [[ "$KEEP_LOGS" != "1" ]]; then
        rm -rf "$WORKDIR"
    else
        echo "logs kept in $WORKDIR"
    fi
}
trap cleanup EXIT

need() {
    command -v "$1" >/dev/null 2>&1 || {
        echo "missing required command: $1" >&2
        exit 1
    }
}

wait_for_socket() {
    local path="$1"
    local i
    for i in $(seq 1 80); do
        [[ -S "$path" ]] && return 0
        sleep 0.1
    done
    echo "timeout waiting for socket: $path" >&2
    exit 1
}

ctrl_call() {
    local id="$1"
    local type="$2"
    local payload="${3:-{}}"
    printf '{"id":"%s","type":"%s","payload":%s}\n' "$id" "$type" "$payload" \
        | socat - "UNIX-CONNECT:$SOCK"
}

expect_ok() {
    local name="$1"
    local json="$2"
    echo "$json" | jq .
    if [[ "$(echo "$json" | jq -r '.ok')" != "true" ]]; then
        echo "FAILED: $name" >&2
        exit 1
    fi
    echo "OK: $name"
}

need jq
need socat
[[ -x "$CTRL_BIN" ]] || { echo "missing or non-executable: $CTRL_BIN" >&2; exit 1; }
[[ -e "$TTY" ]] || { echo "serial device not found: $TTY" >&2; exit 1; }

mkdir -p "$WORKDIR"
cat > "$CTRL_CFG" <<EOF_CTRL
[service]
socket = "$SOCK"

[backend]
type = "raw-rs485"

[raw_rs485]
device = "$TTY"
baud = 9600
EOF_CTRL

"$CTRL_BIN" -c "$CTRL_CFG" -l TRACE >"$CTRL_LOG" 2>&1 &
CTRL_PID="$!"
wait_for_socket "$SOCK"

resp="$(ctrl_call status-1 get_status)"
expect_ok "get_status" "$resp"

resp="$(ctrl_call scan-1 scan)"
expect_ok "scan/connect" "$resp"

echo "$resp" | jq -e '.payload.linkState == "CONNECTED"' >/dev/null || {
    echo "scan did not reach CONNECTED state" >&2
    exit 1
}

resp="$(ctrl_call info-1 get_info)"
expect_ok "get_info" "$resp"

resp="$(ctrl_call err-1 get_alarm_status)"
expect_ok "get_alarm_status" "$resp"

# Read-only smoke ends here by default.
# To run movement tests, set AISG_MOVE=1 and optionally AISG_CONFIG_BLOB=/path/to/vendor.cfg.
if [[ "$MOVE" == "1" ]]; then
    if [[ -n "$CONFIG_BLOB" ]]; then
        [[ -f "$CONFIG_BLOB" ]] || { echo "config blob not found: $CONFIG_BLOB" >&2; exit 1; }
        resp="$(ctrl_call cfg-1 send_configuration_data "{\"configPath\":\"$CONFIG_BLOB\"}")"
        expect_ok "send_configuration_data" "$resp"
    else
        echo "AISG_MOVE=1 without AISG_CONFIG_BLOB; skipping config upload"
    fi

    resp="$(ctrl_call cal-1 calibrate)"
    expect_ok "calibrate" "$resp"

    resp="$(ctrl_call tilt-1 get_tilt)"
    expect_ok "get_tilt" "$resp"

    resp="$(ctrl_call tilt-set set_tilt '{"targetTiltDeg":0.5}')"
    expect_ok "set_tilt 0.5" "$resp"

    resp="$(ctrl_call tilt-2 get_tilt)"
    expect_ok "get_tilt verify" "$resp"
else
    echo "read-only hardware ladder passed; set AISG_MOVE=1 to run config/calibrate/set-tilt"
fi
