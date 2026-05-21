#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

TFTP_PID=""
TFTP_ROOT=""
TFTP_SYSTEMD_WAS_ACTIVE=0

tftp_serve() {
    local serve_dir="$1"

    if [ ! -d "$serve_dir" ]; then
        echo "tftp_serve: directory not found: $serve_dir" >&2
        return 1
    fi

    if ! command -v in.tftpd >/dev/null 2>&1 && [ ! -x /usr/sbin/in.tftpd ]; then
        sudo apt-get update -qq || true
        if ! sudo apt-get install -y tftpd-hpa; then
            echo "tftp_serve: failed to install tftpd-hpa." >&2
            echo "  Please install manually: sudo apt-get install tftpd-hpa" >&2
            return 1
        fi
    fi

    if systemctl is-active tftpd-hpa >/dev/null 2>&1; then
        echo "tftp_serve: stopping system tftpd-hpa so our server can bind port 69"
        sudo systemctl stop tftpd-hpa
        TFTP_SYSTEMD_WAS_ACTIVE=1
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
    if [ "$TFTP_SYSTEMD_WAS_ACTIVE" -eq 1 ]; then
        sudo systemctl start tftpd-hpa 2>/dev/null || true
        TFTP_SYSTEMD_WAS_ACTIVE=0
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
