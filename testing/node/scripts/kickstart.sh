#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2022-present, Ukama Inc.

set -euo pipefail

CONF="/etc/supervisor.conf"
SUPERVISORCTL=(supervisorctl -c "$CONF")

# IMPORTANT: noded_latest is part of the group "on-boot"
NODED="on-boot:noded_latest"
BOOTSTRAP="bootstrap_latest"
MESHD="meshd_latest"

ctl() {
    "${SUPERVISORCTL[@]}" "$@"
}

status_line() {
    local name="$1"
    ctl status "$name" 2>/dev/null || true
}

is_state() {
    local name="$1"
    local want="$2"
    ctl status "$name" 2>/dev/null | awk '{print $2}' | grep -qx "$want"
}

wait_state() {
    local name="$1"
    local want="$2"
    local interval="${3:-1}"

    echo "Waiting for ${name} to be ${want}..."
    while ! is_state "$name" "$want"; do
        status_line "$name"
        sleep "$interval"
    done
}

start_prog() {
    local name="$1"
    echo "Starting ${name}..."
    ctl start "$name"
}

start_group() {
    local group="$1"
    echo "Starting ${group} group..."
    ctl start "${group}:*"
}

svc_port_from_etc_services() {
    # Usage: svc_port_from_etc_services bootstrap 18014
    local svc="$1"
    local def="$2"
    local port=""

    # Match lines like: "bootstrap       18014/tcp   bootstrap"
    port="$(awk -v s="$svc" '
        $1==s {
          split($2,a,"/");
          if (a[1] ~ /^[0-9]+$/) { print a[1]; exit }
        }' /etc/services 2>/dev/null || true)"

    if [[ -n "$port" ]]; then
        echo "$port"
    else
        echo "$def"
    fi
}

wait_http_ready() {
    local url="$1"
    local timeout="$2"
    local interval="$3"

    echo "Waiting for bootstrap readiness: ${url} (timeout ${timeout}s)..."

    local start now
    start="$(date +%s)"

    while true; do
        # If bootstrap crashes while we're waiting, fail fast and show status.
        if is_state "$BOOTSTRAP" "FATAL" || is_state "$BOOTSTRAP" "BACKOFF" || is_state "$BOOTSTRAP" "STOPPED"; then
            echo "ERROR: ${BOOTSTRAP} is not healthy while waiting for readiness:"
            status_line "$BOOTSTRAP"
            return 1
        fi

        if command -v curl >/dev/null 2>&1; then
            if curl -fsS "$url" >/dev/null 2>&1; then
                echo "Bootstrap is ready."
                return 0
            fi
        elif command -v wget >/dev/null 2>&1; then
            if wget -q -O /dev/null "$url" >/dev/null 2>&1; then
                echo "Bootstrap is ready."
                return 0
            fi
        else
            echo "ERROR: neither curl nor wget found; can't probe ${url}"
            return 1
        fi

        now="$(date +%s)"
        if (( now - start >= timeout )); then
            echo "ERROR: bootstrap readiness timed out after ${timeout}s"
            echo "Bootstrap status:"
            status_line "$BOOTSTRAP"
            return 1
        fi

        sleep "$interval"
    done
}

# ---- readiness config (override via env if needed) ----
BOOTSTRAP_HOST="${BOOTSTRAP_HOST:-127.0.0.1}"
BOOTSTRAP_PORT="${BOOTSTRAP_PORT:-$(svc_port_from_etc_services bootstrap 18014)}"
BOOTSTRAP_EP="${BOOTSTRAP_EP:-/v1/ping}"
BOOTSTRAP_URL="http://${BOOTSTRAP_HOST}:${BOOTSTRAP_PORT}${BOOTSTRAP_EP}"

READY_TIMEOUT_SECS="${READY_TIMEOUT_SECS:-120}"
READY_POLL_SECS="${READY_POLL_SECS:-1}"

echo "Kickstart using supervisorctl config: $CONF"
echo "Bootstrap readiness URL: ${BOOTSTRAP_URL}"
echo "Supervisor status (sanity):"
ctl status || true

echo "Starting on-boot group..."
start_group "on-boot"
wait_state "$NODED" "RUNNING" 1

echo "Starting bootstrap..."
start_prog "$BOOTSTRAP"
wait_state "$BOOTSTRAP" "RUNNING" 1

wait_http_ready "$BOOTSTRAP_URL" "$READY_TIMEOUT_SECS" "$READY_POLL_SECS"

echo "Starting meshd..."
start_prog "$MESHD"
wait_state "$MESHD" "RUNNING" 1

echo "Starting sys-service group..."
ctl start "sys-service:*" || true

echo "Kickstart complete."
