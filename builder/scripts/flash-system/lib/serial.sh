#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

serial_wait_for_marker() {
    local device="$1"
    local log_file="$2"
    local marker="$3"
    local timeout_secs="${4:-300}"

    if [ ! -e "$device" ]; then
        echo "Serial device $device not found"
        return 1
    fi

    touch "$log_file"

    local cat_pid
    cat "$device" >> "$log_file" &
    cat_pid=$!

    local elapsed=0
    while [ "$elapsed" -lt "$timeout_secs" ]; do
        if grep -q -- "$marker" "$log_file" 2>/dev/null; then
            kill "$cat_pid" 2>/dev/null || true
            return 0
        fi
        sleep 2
        elapsed=$((elapsed + 2))
    done

    kill "$cat_pid" 2>/dev/null || true
    return 1
}
