#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
EMU_BIN="${AISG_EMU_BIN:-$ROOT/emu/aisg-emu}"
WORKDIR="${AISG_TEST_DIR:-$(mktemp -d /tmp/aisg-ret-emu-neg.XXXXXX)}"
PTY="$WORKDIR/aisg-ret0"
EMU_LOG="$WORKDIR/aisg-emu.log"
KEEP_LOGS="${AISG_KEEP_LOGS:-0}"
EMU_PID=""

cleanup() {
    set +e
    if [[ -n "$EMU_PID" ]]; then kill "$EMU_PID" 2>/dev/null || true; fi
    wait "$EMU_PID" 2>/dev/null || true
    if [[ "$KEEP_LOGS" != "1" ]]; then
        rm -rf "$WORKDIR"
    else
        echo "logs kept in $WORKDIR"
    fi
}
trap cleanup EXIT

wait_for_path() {
    local path="$1"
    local i
    for i in $(seq 1 80); do
        [[ -e "$path" ]] && return 0
        sleep 0.1
    done
    echo "timeout waiting for PTY: $path" >&2
    exit 1
}

[[ -x "$EMU_BIN" ]] || { echo "missing or non-executable: $EMU_BIN" >&2; exit 1; }
command -v python3 >/dev/null 2>&1 || { echo "missing python3" >&2; exit 1; }

mkdir -p "$WORKDIR"
"$EMU_BIN" --mode ret --pty "$PTY" --requires-config true -l TRACE >"$EMU_LOG" 2>&1 &
EMU_PID="$!"
wait_for_path "$PTY"

python3 - "$PTY" <<'PY'
import os
import select
import sys
import termios
import time

def fcs16(data):
    fcs = 0xffff
    for b in data:
        fcs ^= b
        for _ in range(8):
            if fcs & 1:
                fcs = (fcs >> 1) ^ 0x8408
            else:
                fcs >>= 1
            fcs &= 0xffff
    return (~fcs) & 0xffff

def hdlc_frame(payload):
    fcs = fcs16(payload)
    data = bytes(payload) + bytes([fcs & 0xff, (fcs >> 8) & 0xff])
    out = bytearray([0x7e])
    for b in data:
        if b in (0x7e, 0x7d):
            out.extend([0x7d, b ^ 0x20])
        else:
            out.append(b)
    out.append(0x7e)
    return bytes(out)

def no_response(fd, frame, name):
    os.write(fd, frame)
    ready, _, _ = select.select([fd], [], [], 0.75)
    if ready:
        data = os.read(fd, 4096)
        raise SystemExit(f"{name}: expected no response, got {data.hex(' ')}")
    print(f"ok - {name} produced no response")

path = sys.argv[1]
fd = os.open(path, os.O_RDWR | os.O_NOCTTY | os.O_NONBLOCK)
attrs = termios.tcgetattr(fd)
attrs[0] = 0
attrs[1] = 0
attrs[2] |= termios.CLOCAL | termios.CREAD
attrs[3] = 0
termios.tcsetattr(fd, termios.TCSANOW, attrs)
time.sleep(0.1)
try:
    # Old invalid scan payload: FF BF 81 F0 00.
    # It has FI/GI with GL=0 and no PI=1/PI=3 scan parameters.
    no_response(fd, hdlc_frame([0xff, 0xbf, 0x81, 0xf0, 0x00]), "old fake scan")

    # RETAP GetInformation I-frame before address assignment/SNRM must be ignored.
    no_response(fd, hdlc_frame([0x01, 0x10, 0x05, 0x00, 0x00]), "RETAP before address assignment")
finally:
    os.close(fd)
PY

echo "AISG ret emulator negative tests passed"
