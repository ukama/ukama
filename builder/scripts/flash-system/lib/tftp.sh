#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

TFTP_PID=""
TFTP_ROOT=""

tftp_serve() {
    local serve_dir="$1"

    if [ ! -d "$serve_dir" ]; then
        echo "tftp_serve: directory not found: $serve_dir" >&2
        return 1
    fi

    if ! command -v in.tftpd >/dev/null 2>&1; then
        sudo apt-get update -qq
        sudo apt-get install -y tftpd-hpa
    fi

    TFTP_ROOT="$serve_dir"
    sudo /usr/sbin/in.tftpd -L --secure --user root "$serve_dir" &
    TFTP_PID=$!
    sleep 1
}

tftp_stop() {
    if [ -n "$TFTP_PID" ]; then
        sudo kill "$TFTP_PID" 2>/dev/null || true
        TFTP_PID=""
    fi
}

tftp_stage_file() {
    local src="$1"
    local name
    name=$(basename "$src")

    if [ -z "$TFTP_ROOT" ]; then
        echo "tftp_stage_file: TFTP not started" >&2
        return 1
    fi

    sudo cp "$src" "${TFTP_ROOT}/${name}"
    sudo chmod 644 "${TFTP_ROOT}/${name}"
    echo "$name"
}
