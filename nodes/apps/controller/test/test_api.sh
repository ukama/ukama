#!/usr/bin/env bash
# REST API smoke test for controller.d
# Run this after starting the daemon with the simulator.
#
# Usage:
#   ./test_api.sh [HOST] [PORT]
#   ./test_api.sh 127.0.0.1 8095

set -u
set -o pipefail

HOST="${1:-${CONTROLLER_HOST:-127.0.0.1}}"
PORT="${2:-${CONTROLLER_PORT:-8095}}"
BASE="${BASE_URL:-http://${HOST}:${PORT}/v1/controller}"
TMP_BODY="$(mktemp /tmp/controllerd-api.XXXXXX.json)"

PASS=0
FAIL=0

cleanup() {
    rm -f "$TMP_BODY"
}

trap cleanup EXIT

require_cmd() {
    local cmd="$1"
    if ! command -v "$cmd" >/dev/null 2>&1; then
        echo "Missing required command: $cmd" >&2
        exit 1
    fi
}

check() {
    local desc="$1"
    local url="$2"
    local method="${3:-GET}"
    local body="${4-}"
    local expect_status="${5:-200}"
    local response

    if [ -n "$body" ]; then
        response=$(curl -sS -o "$TMP_BODY" -w "%{http_code}" \
            -X "$method" -H "Content-Type: application/json" \
            -d "$body" "$url" 2>/dev/null || printf "000")
    else
        response=$(curl -sS -o "$TMP_BODY" -w "%{http_code}" \
            -X "$method" "$url" 2>/dev/null || printf "000")
    fi

    if [ "$response" = "$expect_status" ]; then
        echo "  PASS  [$response] $desc"
        PASS=$((PASS + 1))
    else
        echo "  FAIL  [$response != $expect_status] $desc"
        if [ -s "$TMP_BODY" ]; then
            cat "$TMP_BODY"
        fi
        FAIL=$((FAIL + 1))
    fi
}

require_cmd curl
require_cmd python3

echo ""
echo "============================================"
echo "  controller.d API smoke test"
echo "  Target: $BASE"
echo "============================================"
echo ""

echo "--- Lifecycle ---"
check "ping"                   "$BASE/ping"
check "version"                "$BASE/version"

echo ""
echo "--- Data endpoints ---"
check "status"                 "$BASE/status"
check "metrics"                "$BASE/metrics"
check "alarms"                 "$BASE/alarms"

echo ""
echo "--- Status fields ---"
if curl -s "$BASE/status" | python3 -c "
import sys, json
d = json.load(sys.stdin)
assert 'charge_state' in d, 'missing charge_state'
assert 'comm_ok' in d, 'missing comm_ok'
assert 'solar' in d, 'missing solar'
assert 'battery' in d, 'missing battery'
assert 'voltage_v' in d['battery'], 'missing battery.voltage_v'
assert 'power_w' in d['solar'], 'missing solar.power_w'
print('  PASS  status JSON structure valid')
" 2>/dev/null; then
    PASS=$((PASS + 1))
else
    echo "  FAIL  status JSON structure invalid"
    FAIL=$((FAIL + 1))
fi

echo ""
echo "--- Metrics fields ---"
if curl -s "$BASE/metrics" | python3 -c "
import sys, json
d = json.load(sys.stdin)
assert 'metrics' in d, 'missing metrics array'
names = [m['name'] for m in d['metrics']]
for required in ['solar_panel_power', 'battery_voltage', 'battery_current', 'mppt_efficiency']:
    assert required in names, f'missing metric: {required}'
print('  PASS  metrics JSON structure valid')
" 2>/dev/null; then
    PASS=$((PASS + 1))
else
    echo "  FAIL  metrics JSON structure invalid"
    FAIL=$((FAIL + 1))
fi

echo ""
echo "--- Method not allowed ---"
check "POST /status is 405"    "$BASE/status"  "POST" "{}" "405"
check "DELETE /metrics is 405" "$BASE/metrics" "DELETE" "" "405"

echo ""
echo "--- Control endpoints ---"
check "PUT /absorption (not impl)" "$BASE/absorption" "PUT" '{"voltage_v": 57.6}' "501"
check "PUT /float (not impl)"      "$BASE/float"      "PUT" '{"voltage_v": 55.2}' "501"
check "POST /relay (not impl)"     "$BASE/relay"      "POST" '{"state": true}' "501"

echo ""
echo "--- 404 on unknown path ---"
check "unknown path is 404"    "$BASE/nonexistent" "GET" "" "404"

echo ""
echo "============================================"
echo "  Results: ${PASS} passed, ${FAIL} failed"
echo "============================================"
echo ""

[ "$FAIL" -eq 0 ]
