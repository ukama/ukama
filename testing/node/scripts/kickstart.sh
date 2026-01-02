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
    ctl start "$name" || true
}

start_group() {
    local group="$1"
    echo "Starting ${group} group..."
    ctl start "${group}:*" || true
}

svc_port_from_etc_services() {
    local svc="$1"
    local def="$2"

    awk -v s="$svc" '
        $1==s {
            split($2,a,"/");
            if (a[1] ~ /^[0-9]+$/) { print a[1]; exit }
        }
    ' /etc/services 2>/dev/null || echo "$def"
}

http_ping() {
    local url="$1"

    if command -v curl >/dev/null 2>&1; then
        curl -fsS "$url" >/dev/null 2>&1
    elif command -v wget >/dev/null 2>&1; then
        wget -q -O /dev/null "$url" >/dev/null 2>&1
    else
        return 2
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
        if is_state "$BOOTSTRAP" "FATAL" || \
           is_state "$BOOTSTRAP" "BACKOFF" || \
           is_state "$BOOTSTRAP" "STOPPED"; then
            echo "ERROR: ${BOOTSTRAP} unhealthy:"
            status_line "$BOOTSTRAP"
            return 1
        fi

        if http_ping "$url"; then
            echo "Bootstrap is ready."
            return 0
        fi

        now="$(date +%s)"
        if (( now - start >= timeout )); then
            echo "ERROR: bootstrap readiness timed out"
            status_line "$BOOTSTRAP"
            return 1
        fi

        sleep "$interval"
    done
}

# ---- bootstrap readiness ----
BOOTSTRAP_HOST="127.0.0.1"
BOOTSTRAP_PORT="$(svc_port_from_etc_services bootstrap 18014)"
BOOTSTRAP_EP="/v1/ping"
BOOTSTRAP_URL="http://${BOOTSTRAP_HOST}:${BOOTSTRAP_PORT}${BOOTSTRAP_EP}"

READY_TIMEOUT_SECS=120
READY_POLL_SECS=1

# ---- post-boot apps: start + single ping ----
APP_HOST="127.0.0.1"
APP_EP="/v1/ping"

# Format per entry:
#   "supervisor_program|service_name_in_/etc/services|default_port"
#
# supervisor_program: what you pass to supervisorctl (e.g., "inventory_latest")
# service_name: first column in /etc/services (e.g., "inventory")
# default_port: fallback if not found in /etc/services
POST_BOOT_APPS=(
    "rlog_latest|rlog|0"
    "deviced_latest|device|0"
    "wimcd_latest|wimc|0"
    "lookoutd_latest|lookout|0"
    "configd_latest|config|0"
    "metricsd_latest|metrics-admin|0"
    "gpsd_latest|gps|0"
    "notifyd_latest|notify-admin|0"
)

start_and_check_post_boot_apps() {
    if (( ${#POST_BOOT_APPS[@]} == 0 )); then
        echo "No post-boot apps configured."
        return 0
    fi

    echo "Starting post-boot apps (not in on-boot) and probing once..."

    local item prog svc def_port port url
    for item in "${POST_BOOT_APPS[@]}"; do
        prog="${item%%|*}"
        svc="${item#*|}"; svc="${svc%%|*}"
        def_port="${item##*|}"

        echo "Starting ${prog}..."
        ctl start "$prog" || true

        port="$(svc_port_from_etc_services "$svc" "$def_port")"
        url="http://${APP_HOST}:${port}${APP_EP}"

        if http_ping "$url"; then
            echo "APP ${prog} (${url}) : RUNNING"
        else
            echo "APP ${prog} (${url}) : ERROR"
        fi
    done
}

echo "Kickstart using supervisorctl config: $CONF"
ctl status || true

echo "Starting on-boot group..."
start_group "on-boot"
wait_state "$NODED" "RUNNING"

echo "Starting bootstrap..."
start_prog "$BOOTSTRAP"
wait_state "$BOOTSTRAP" "RUNNING"

wait_http_ready "$BOOTSTRAP_URL" "$READY_TIMEOUT_SECS" "$READY_POLL_SECS"

echo "Starting meshd..."
start_prog "$MESHD"
wait_state "$MESHD" "RUNNING"

echo "Starting sys-service group..."
ctl start "sys-service:*" || true

start_and_check_post_boot_apps

echo "Kickstart complete."
