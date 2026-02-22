#!/usr/bin/env bash
set -euo pipefail

# Smoke test for device.d
#
# Modes:
#   amplifier: tests amplifier node device.d (radio on/off)
#   tower:     tests tower node server device.d AND client-mode device.d
#
# Usage:
#   ./smoke-deviced.sh amplifier --host 127.0.0.1 --port 8080
#
#   ./smoke-deviced.sh tower \
#       --server-host 127.0.0.1 --server-port 8080 \
#       --client-host 127.0.0.1 --client-port 8090
#
# Options:
#   --timeout <sec>        Curl timeout (default: 4)
#   --allow-reboot         Actually call /v1/restart (DANGEROUS)
#
# Notes:
# - device.d must already be running in the relevant modes.
# - Safe by default: restart calls are skipped unless --allow-reboot is provided.

#
# Sample runs:
#  To run amplifier node tests:
# ./smoke-deviced.sh amplifier --host 127.0.0.1 --port 8080
#
#  Tower node test (including client-mode):
# ./smoke-deviced.sh tower \
#  --server-host 127.0.0.1 --server-port 8080 \
#  --client-host 127.0.0.1 --client-port 8090
#

AllowReboot=0
TimeoutSec=4

Mode=""

Host="127.0.0.1"
Port="18003"

# Always run in debug mode (otherwise restart will actualy restart this machine)
export DEVICED_DEBUG_MODE=1

ServerHost="localhost"
ServerPort="10803"
ClientHost="localhost"
ClientPort="18004"

die() { echo "ERROR: $*" >&2; exit 1; }
log() { printf '[%s] %s\n' "$(date +%H:%M:%S)" "$*"; }

usage() {
    cat <<EOF
Usage:
  $0 amplifier --host <ip> --port <port> [--timeout <sec>] [--allow-reboot]

  $0 tower --server-host <ip> --server-port <port> --client-host <ip> --client-port <port> \\
           [--timeout <sec>] [--allow-reboot]

Examples:
  $0 amplifier --host 127.0.0.1 --port 8080

  $0 tower --server-host 127.0.0.1 --server-port 8080 --client-host 127.0.0.1 --client-port 8090

EOF
}

need_bin() { command -v "$1" >/dev/null 2>&1 || die "Missing required binary: $1"; }

curl_req() {
    # args: METHOD URL BODY(optional)
    local method="$1"; shift
    local url="$1"; shift
    local body="${1:-}"

    if [[ -n "$body" ]]; then
        curl -sS -i --max-time "${TimeoutSec}" \
            -X "$method" "$url" \
            -H "Content-Type: application/json" \
            --data "$body"
    else
        curl -sS -i --max-time "${TimeoutSec}" \
            -X "$method" "$url"
    fi
}

http_code() {
    # args: response
    printf '%s' "$1" | head -n 1 | awk '{print $2}'
}

assert_http() {
    # args: response expected_code label
    local resp="$1"
    local exp="$2"
    local label="$3"

    local got
    got="$(http_code "$resp")"
    if [[ "$got" != "$exp" ]]; then
        echo "$resp" >&2
        die "$label: expected HTTP $exp, got $got"
    fi
    log "OK: $label -> HTTP $got"
}

test_ping() {
    local host="$1"
    local port="$2"
    local label="$3"

    local url="http://${host}:${port}/v1/ping"
    log "PING ($label): $url"
    local resp
    resp="$(curl_req GET "$url")"
    assert_http "$resp" "200" "ping ($label)"
}

test_restart() {
    local host="$1"
    local port="$2"
    local label="$3"

    local url="http://${host}:${port}/v1/restart"
    log "RESTART ($label): $url"

    if [[ "$AllowReboot" -eq 0 ]]; then
        log "SKIP: reboot disabled (use --allow-reboot to enable)"
        return 0
    fi

    local resp
    resp="$(curl_req POST "$url")"
    # Most async restart implementations return 202.
    # If yours returns 200, change this expectation.
    assert_http "$resp" "202" "restart ($label)"
}

test_invalid_json_400() {
    local host="$1"
    local port="$2"
    local path="$3"
    local label="$4"

    local url="http://${host}:${port}${path}"
    log "INVALID JSON ($label): $url"
    local resp code
    resp="$(curl_req POST "$url" '{bad json' || true)"
    code="$(http_code "$resp")"
    if [[ "$code" != "400" ]]; then
        echo "$resp" >&2
        die "invalid-json ($label): expected HTTP 400, got $code"
    fi
    log "OK: invalid-json ($label) -> HTTP $code"
}

test_concurrency_soft() {
    # Not super strict; just ensures server stays responsive and returns valid codes.
    local host="$1"
    local port="$2"
    local path="$3"
    local a="$4"
    local b="$5"
    local label="$6"

    local url="http://${host}:${port}${path}"
    log "CONCURRENCY ($label): $url"

    local r1 r2 c1 c2
    r1="$(curl_req POST "$url" "$a" || true)"; c1="$(http_code "$r1")"
    r2="$(curl_req POST "$url" "$b" || true)"; c2="$(http_code "$r2")"

    log "CONCURRENCY ($label) codes: $c1 , $c2"
    [[ "$c1" != "000" && "$c2" != "000" ]] || die "concurrency ($label): curl failed (000)."

    # Typical good patterns: 202+409 or 202+202 depending on implementation.
    # We won't enforce exact behavior here.
}

test_amplifier() {
    local host="$1"
    local port="$2"

    [[ -n "$port" ]] || die "amplifier: missing --port"

    log "=== Amplifier smoke test @ ${host}:${port} ==="

    test_ping "$host" "$port" "amplifier"

    local resp

    resp="$(curl_req POST "http://${host}:${port}/v1/radio" '{"state":"on"}')"
    assert_http "$resp" "202" "radio on"

    resp="$(curl_req POST "http://${host}:${port}/v1/radio" '{"state":"off"}')"
    assert_http "$resp" "202" "radio off"

    test_invalid_json_400 "$host" "$port" "/v1/radio" "amplifier/radio"
    test_concurrency_soft "$host" "$port" "/v1/radio" '{"state":"on"}' '{"state":"off"}' "amplifier/radio"

    test_restart "$host" "$port" "amplifier"

    log "=== Amplifier smoke test PASS ==="
}

test_tower_server() {
    local host="$1"
    local port="$2"

    log "--- Tower SERVER checks @ ${host}:${port} ---"

    test_ping "$host" "$port" "tower-server"

    local resp

    resp="$(curl_req POST "http://${host}:${port}/v1/service" '{"state":"on"}')"
    assert_http "$resp" "202" "service on"

    resp="$(curl_req POST "http://${host}:${port}/v1/service" '{"state":"off"}')"
    assert_http "$resp" "202" "service off"

    test_invalid_json_400 "$host" "$port" "/v1/service" "tower/service"
    test_concurrency_soft "$host" "$port" "/v1/service" '{"state":"on"}' '{"state":"off"}' "tower/service"

    test_restart "$host" "$port" "tower-server"
}

test_tower_client() {
    local host="$1"
    local port="$2"

    log "--- Tower CLIENT-MODE checks @ ${host}:${port} ---"

    test_ping "$host" "$port" "tower-client"

    # Optional: only call if allowed. If client-mode doesn't implement it, it may be 404/405.
    if [[ "$AllowReboot" -eq 1 ]]; then
        local resp code
        resp="$(curl_req POST "http://${host}:${port}/v1/restart" || true)"
        code="$(http_code "$resp")"
        if [[ "$code" == "202" ]]; then
            log "OK: restart (tower-client) -> HTTP 202"
        else
            log "WARN: restart (tower-client) returned HTTP $code (may be unimplemented)"
        fi
    else
        log "SKIP: tower-client restart disabled (use --allow-reboot)"
    fi
}

test_tower_combined() {
    local shost="$1" sport="$2"
    local chost="$3" cport="$4"

    [[ -n "$sport" ]] || die "tower: missing --server-port"
    [[ -n "$cport" ]] || die "tower: missing --client-port"

    log "=== Tower smoke test (SERVER + CLIENT) ==="
    log "Server: ${shost}:${sport}"
    log "Client: ${chost}:${cport}"

    test_tower_server "$shost" "$sport"
    test_tower_client "$chost" "$cport"

    log "=== Tower smoke test PASS ==="
}

# ---- arg parsing ----
[[ $# -ge 1 ]] || { usage; exit 1; }
Mode="$1"; shift

while [[ $# -gt 0 ]]; do
    case "$1" in
        --timeout) TimeoutSec="$2"; shift 2;;
        --allow-reboot) AllowReboot=1; shift;;

        # amplifier
        --host) Host="$2"; shift 2;;
        --port) Port="$2"; shift 2;;

        # tower combined
        --server-host) ServerHost="$2"; shift 2;;
        --server-port) ServerPort="$2"; shift 2;;
        --client-host) ClientHost="$2"; shift 2;;
        --client-port) ClientPort="$2"; shift 2;;

        -h|--help) usage; exit 0;;
        *) die "Unknown arg: $1";;
    esac
done

need_bin curl

case "$Mode" in
    amplifier)
        test_amplifier "$Host" "$Port"
        ;;
    tower)
        [[ -n "$ServerHost" ]] || ServerHost="127.0.0.1"
        [[ -n "$ClientHost" ]] || ClientHost="127.0.0.1"
        test_tower_combined "$ServerHost" "$ServerPort" "$ClientHost" "$ClientPort"
        ;;
    *) usage; die "Invalid mode: $Mode" ;;
esac
