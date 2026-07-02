#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
EMU_BIN="${AISG_EMU_BIN:-$ROOT/emu/aisg-emu}"
CTRL_BIN="${AISG_CTRL_BIN:-$ROOT/ctrl/aisg-ctrl}"
WORKDIR="${AISG_TEST_DIR:-$(mktemp -d /tmp/aisg-ret-emu.XXXXXX)}"
PTY="$WORKDIR/aisg-ret0"
SOCK="$WORKDIR/aisg-ctrl.sock"
CTRL_CFG="$WORKDIR/aisg-ctrl.toml"
CONFIG_BLOB="$WORKDIR/antenna.cfg"
EMU_LOG="$WORKDIR/aisg-emu.log"
CTRL_LOG="$WORKDIR/aisg-ctrl.log"
KEEP_LOGS="${AISG_KEEP_LOGS:-0}"

EMU_PID=""
CTRL_PID=""

cleanup() {
    set +e
    if [[ -n "$CTRL_PID" ]]; then kill "$CTRL_PID" 2>/dev/null || true; fi
    if [[ -n "$EMU_PID" ]]; then kill "$EMU_PID" 2>/dev/null || true; fi
    wait "$CTRL_PID" 2>/dev/null || true
    wait "$EMU_PID" 2>/dev/null || true
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

wait_for_path() {
    local path="$1"
    local name="$2"
    local i
    for i in $(seq 1 80); do
        [[ -e "$path" ]] && return 0
        sleep 0.1
    done
    echo "timeout waiting for $name: $path" >&2
    exit 1
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

expect_fail_code() {
    local name="$1"
    local expected="$2"
    local json="$3"
    echo "$json" | jq .
    if [[ "$(echo "$json" | jq -r '.ok')" != "false" ]]; then
        echo "FAILED: $name expected failure" >&2
        exit 1
    fi
    if [[ "$(echo "$json" | jq -r '.code')" != "$expected" ]]; then
        echo "FAILED: $name expected code $expected" >&2
        exit 1
    fi
    echo "OK: $name failed with $expected"
}

need jq
need socat
need timeout

[[ -x "$EMU_BIN" ]] || { echo "missing or non-executable: $EMU_BIN" >&2; exit 1; }
[[ -x "$CTRL_BIN" ]] || { echo "missing or non-executable: $CTRL_BIN" >&2; exit 1; }

mkdir -p "$WORKDIR"
printf 'ukama-aisg-test-config\n' > "$CONFIG_BLOB"

cat > "$CTRL_CFG" <<EOF_CTRL
[service]
socket = "$SOCK"

[backend]
type = "raw-rs485"

[raw_rs485]
device = "$PTY"
baud = 9600
EOF_CTRL

"$EMU_BIN" --mode ret --pty "$PTY" --requires-config true \
    --min-tilt 0 --max-tilt 10 --initial-tilt 3.0 \
    -l TRACE >"$EMU_LOG" 2>&1 &
EMU_PID="$!"
wait_for_path "$PTY" "RET PTY"

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
expect_ok "get_alarm_status before config" "$resp"
echo "$resp" | jq -e '.payload.active[]? | select(.name == "NotScaled")' >/dev/null || {
    echo "expected NotScaled before config" >&2
    exit 1
}

resp="$(ctrl_call tilt-pre get_tilt)"
expect_fail_code "get_tilt before calibration" "NotCalibrated" "$resp"

resp="$(ctrl_call cfg-1 send_configuration_data "{\"configPath\":\"$CONFIG_BLOB\"}")"
expect_ok "send_configuration_data" "$resp"
echo "$resp" | jq -e '.payload.chunks == 1' >/dev/null || {
    echo "expected one config chunk" >&2
    exit 1
}

resp="$(ctrl_call cal-1 calibrate)"
expect_ok "calibrate" "$resp"

resp="$(ctrl_call tilt-1 get_tilt)"
expect_ok "get_tilt" "$resp"

resp="$(ctrl_call tilt-set set_tilt '{"targetTiltDeg":0.5}')"
expect_ok "set_tilt" "$resp"

resp="$(ctrl_call tilt-2 get_tilt)"
expect_ok "get_tilt verify" "$resp"
echo "$resp" | jq -e '(.payload.currentTiltDeg >= 0.49) and (.payload.currentTiltDeg <= 0.51)' >/dev/null || {
    echo "expected final tilt near 0.5 deg" >&2
    exit 1
}

resp="$(ctrl_call status-2 get_status)"
expect_ok "final get_status" "$resp"

echo "AISG ret emulator ladder passed"
echo "workdir: $WORKDIR"
