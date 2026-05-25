#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

bdi_require_expect() {
    if ! command -v expect >/dev/null 2>&1; then
        sudo apt-get update -qq || true
        if ! sudo apt-get install -y expect; then
            echo "bdi_require_expect: failed to install 'expect'." >&2
            echo "  Please install manually: sudo apt-get install expect" >&2
            return 1
        fi
    fi
}

bdi_wait_prompt() {
    local ip="$1"
    local prompt="$2"
    local timeout_secs="${3:-30}"

    bdi_require_expect

    expect <<EOF
set timeout ${timeout_secs}
spawn telnet ${ip}
expect {
    "${prompt}" { send "quit\r"; expect eof; exit 0 }
    timeout      { exit 1 }
}
EOF
}

bdi_send_command() {
    local ip="$1"
    local prompt="$2"
    local command="$3"
    local timeout_secs="${4:-60}"

    bdi_require_expect

    expect <<EOF
set timeout ${timeout_secs}
spawn telnet ${ip}
expect {
    "${prompt}" {}
    timeout     { exit 1 }
}
send "${command}\r"
expect {
    "${prompt}" { send "quit\r"; expect eof; exit 0 }
    timeout     { exit 2 }
}
EOF
}

bdi_send_sequence() {
    local ip="$1"
    local prompt="$2"
    local timeout_secs="$3"
    shift 3

    bdi_require_expect

    local script
    script="set timeout ${timeout_secs}
spawn telnet ${ip}
expect {
    \"${prompt}\" {}
    timeout     { exit 1 }
}
"
    local cmd
    for cmd in "$@"; do
        script="${script}
send \"${cmd}\r\"
expect {
    \"${prompt}\" {}
    timeout     { exit 2 }
}
"
    done

    script="${script}
send \"quit\r\"
expect eof
exit 0
"
    expect -c "$script"
}
