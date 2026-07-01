#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

UBOOT_CAT_PID=""
UBOOT_LOG=""

uboot_open() {
    local dev="$1"
    local baud="$2"
    local log="$3"

    if [ ! -e "$dev" ]; then
        echo "uboot_open: serial device not found: $dev" >&2
        return 1
    fi

    sudo stty -F "$dev" "$baud" cs8 -cstopb -parenb -ixon -echo raw

    UBOOT_LOG="$log"
    : > "$log"

    sudo cat "$dev" >> "$log" &
    UBOOT_CAT_PID=$!
    sleep 1
}

uboot_close() {
    if [ -n "$UBOOT_CAT_PID" ]; then
        sudo kill "$UBOOT_CAT_PID" 2>/dev/null || true
        UBOOT_CAT_PID=""
    fi
}

uboot_drain() {
    local quiet_secs="${1:-3}"
    local last="" cur stable=0
    while [ "$stable" -lt "$quiet_secs" ]; do
        cur=$(wc -c < "$UBOOT_LOG" 2>/dev/null || echo 0)
        if [ "$cur" = "$last" ]; then
            stable=$((stable + 1))
        else
            stable=0
            last="$cur"
        fi
        sleep 1
    done
}

uboot_wait_for() {
    local pattern="$1"
    local timeout_secs="${2:-30}"

    if [ -z "$UBOOT_LOG" ]; then
        echo "uboot_wait_for: serial not opened" >&2
        return 1
    fi

    local elapsed=0
    while [ "$elapsed" -lt "$timeout_secs" ]; do
        if grep -qF -- "$pattern" "$UBOOT_LOG" 2>/dev/null; then
            return 0
        fi
        sleep 1
        elapsed=$((elapsed + 1))
    done
    return 1
}

uboot_send() {
    local dev="$1"
    local command="$2"
    local i
    exec 3>"$dev"
    for (( i=0; i<${#command}; i++ )); do
        printf '%s' "${command:$i:1}" >&3
        sleep 0.02
    done
    printf '\r\n' >&3
    exec 3>&-
}

uboot_send_and_wait() {
    local dev="$1"
    local command="$2"
    local prompt="$3"
    local timeout_secs="${4:-30}"

    uboot_send "$dev" "$command"
    sleep 1

    local plen="${#prompt}"
    local waited=1 last="" cur
    while [ "$waited" -lt "$timeout_secs" ]; do
        cur=$(wc -c < "$UBOOT_LOG" 2>/dev/null || echo 0)
        if [ "$cur" = "$last" ] && tail -c "$plen" "$UBOOT_LOG" 2>/dev/null | grep -qF -- "$prompt"; then
            return 0
        fi
        last="$cur"
        sleep 1
        waited=$((waited + 1))
    done
    return 1
}
