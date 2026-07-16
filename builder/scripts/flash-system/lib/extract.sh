#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

extract_tarball_to_part() {
    local tarball="$1"
    local part="$2"
    local label="$3"

    local mnt
    mnt=$(mktemp -d "/tmp/ukama-${label}.XXXXXX")

    sudo mount "$part" "$mnt"
    sudo tar -xzf "$tarball" -C "$mnt" --no-same-owner --no-same-permissions --warning=no-timestamp --touch
    sync
    sudo umount "$mnt"
    rm -rf "$mnt"
}

install_payloads_on_part() {
    local payload_dir="$1"
    local part="$2"
    local label="$3"

    local mnt
    mnt=$(mktemp -d "/tmp/ukama-${label}-payload.XXXXXX")
    sudo mount "$part" "$mnt"

    if [ -f "$payload_dir/ukama-auto-flash.sh" ]; then
        sudo mkdir -p "$mnt/usr/local/sbin"
        sudo install -m 755 "$payload_dir/ukama-auto-flash.sh" "$mnt/usr/local/sbin/ukama-auto-flash.sh"
    fi

    if [ -f "$payload_dir/ukama-autoflash.service" ]; then
        sudo mkdir -p "$mnt/etc/systemd/system/multi-user.target.wants"
        sudo install -m 644 "$payload_dir/ukama-autoflash.service" "$mnt/etc/systemd/system/ukama-autoflash.service"
        sudo ln -sf /etc/systemd/system/ukama-autoflash.service \
            "$mnt/etc/systemd/system/multi-user.target.wants/ukama-autoflash.service"
    fi

    sync
    sudo umount "$mnt"
    rm -rf "$mnt"
}

render_template() {
    local template="$1"
    local output="$2"
    shift 2

    cp "$template" "$output"

    while [ $# -ge 2 ]; do
        local placeholder="$1"
        local value="$2"
        shift 2
        sed -i.bak "s|@@${placeholder}@@|${value}|g" "$output"
    done
    rm -f "${output}.bak"
}
