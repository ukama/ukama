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

require_cmd bash python3 mktemp readlink

INIT_BIN="${INIT_STARTER_BIN:-${1:-init-starter}}"
if ! command -v "$INIT_BIN" >/dev/null 2>&1 && [[ ! -x "$INIT_BIN" ]]; then
    echo "init-starter binary not found or not executable: $INIT_BIN" >&2
    exit 1
fi

TMP="$(mktemp -d /tmp/init-starter-smoke.XXXXXX)"
ROOT="$TMP/initroot"
SLOTS="$ROOT/slots"
READY_FILE="$TMP/ready"
LOG_FILE="$TMP/slot-runs.log"
INIT_LOG="$TMP/init-starter.log"
mkdir -p "$SLOTS/A" "$SLOTS/B"

cleanup() {
    set +e
    [[ -n "${INIT_PID:-}" ]] && kill "$INIT_PID" >/dev/null 2>&1 || true
    wait "${INIT_PID:-}" >/dev/null 2>&1 || true
    rm -rf "$TMP"
}
trap cleanup EXIT

cat > "$SLOTS/A/starter.d" <<'SH'
#!/usr/bin/env bash
set -euo pipefail
: "${STARTER_INIT_READY_FILE:?missing STARTER_INIT_READY_FILE}"
: "${STARTER_TEST_LOG_FILE:?missing STARTER_TEST_LOG_FILE}"
: "${STARTER_TEST_ROOT:?missing STARTER_TEST_ROOT}"
echo "A" >> "$STARTER_TEST_LOG_FILE"
touch "$STARTER_INIT_READY_FILE"
sleep 0.5
exit 77
SH
chmod +x "$SLOTS/A/starter.d"

cat > "$SLOTS/B/starter.d" <<'SH'
#!/usr/bin/env bash
set -euo pipefail
: "${STARTER_INIT_READY_FILE:?missing STARTER_INIT_READY_FILE}"
: "${STARTER_TEST_LOG_FILE:?missing STARTER_TEST_LOG_FILE}"
echo "B" >> "$STARTER_TEST_LOG_FILE"
touch "$STARTER_INIT_READY_FILE"
sleep 2
exit 0
SH
chmod +x "$SLOTS/B/starter.d"

ln -s slots/A "$ROOT/current"
ln -s slots/B "$ROOT/next"

export STARTER_INIT_ROOT="$ROOT"
export STARTER_INIT_READY_FILE="$READY_FILE"
export STARTER_INIT_READY_TIMEOUT_SEC=5
export STARTER_INIT_TERM_GRACE_SEC=1
export STARTER_TEST_LOG_FILE="$LOG_FILE"
export STARTER_TEST_ROOT="$ROOT"

"$INIT_BIN" >"$INIT_LOG" 2>&1 &
INIT_PID=$!

wait_for_file "$LOG_FILE" 10
wait_for_json_condition "file://unused" 'True' 0 >/dev/null 2>&1 || true

python3 - "$LOG_FILE" <<'PY'
import sys, time, pathlib
p = pathlib.Path(sys.argv[1])
end = time.time() + 10
while time.time() < end:
    txt = p.read_text() if p.exists() else ""
    if "A" in txt and "B" in txt:
        raise SystemExit(0)
    time.sleep(0.2)
raise SystemExit(1)
PY

[[ "$(readlink "$ROOT/current")" == "slots/B" ]]
[[ "$(readlink "$ROOT/prev")" == "slots/A" ]]

echo "[ok] init-starter observed exit 77 and switched A -> B"
echo "slot runs:"; cat "$LOG_FILE"
echo "smoke_init_switch: PASS"
