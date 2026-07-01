#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

source "${LIB_DIR}/partition.sh"
source "${LIB_DIR}/extract.sh"
source "${LIB_DIR}/serial.sh"

method_validate() {
    local fail=0

    local keys
    keys=$(yq_keys "$BOARD_CONFIG" tarballs)

    local key
    for key in $keys; do
        local path must_contain min_entries
        path=$(yq_read "$BOARD_CONFIG" "tarballs.${key}.path")
        must_contain=$(yq_join "$BOARD_CONFIG" "tarballs.${key}.must_contain" ",")
        min_entries=$(yq_read "$BOARD_CONFIG" "tarballs.${key}.min_entries")
        [ "$min_entries" = "null" ] && min_entries=100
        validate_tarball "$path" "$key" "$must_contain" "$min_entries" || fail=1
    done

    return $fail
}

method_confirm() {
    local target
    target=$(yq_read "$BOARD_CONFIG" target)

    if [ ! -b "$target" ]; then
        echo "ERROR: $target is not a block device"
        lsblk
        return 1
    fi

    echo ""
    echo "Plan:"
    echo "  target : $target"
    lsblk "$target" || true
    echo ""
    echo "This will ERASE $target."
    read -rp "Type 'yes' to continue: " confirm
    [ "$confirm" = "yes" ]
}

method_apply() {
    local target
    target=$(yq_read "$BOARD_CONFIG" target)

    echo "Unmounting any auto-mounted partitions on $target..."
    partition_unmount_all "$target"

    echo "Wiping $target..."
    partition_wipe "$target"

    echo "Creating partitions..."
    partition_apply "$BOARD_CONFIG" "$target"

    echo "Formatting partitions..."
    partition_format "$BOARD_CONFIG" "$target"

    local count i
    count=$(yq_count "$BOARD_CONFIG" partitions)

    for ((i=0; i<count; i++)); do
        local num tarball_key install_autoflash
        num=$(yq_read "$BOARD_CONFIG" "partitions[$i].num")
        tarball_key=$(yq_read "$BOARD_CONFIG" "partitions[$i].tarball")
        install_autoflash=$(yq_read "$BOARD_CONFIG" "partitions[$i].install_autoflash")

        local part_name
        part_name=$(partition_name "$target" "$num")

        if [ -n "$tarball_key" ] && [ "$tarball_key" != "null" ]; then
            local tarball_path
            tarball_path=$(yq_read "$BOARD_CONFIG" "tarballs.${tarball_key}.path")
            echo "Extracting ${tarball_key} -> ${part_name}..."
            extract_tarball_to_part "$tarball_path" "$part_name" "$tarball_key"
        fi

        if [ "$install_autoflash" = "true" ]; then
            echo "Installing autoflash payload on ${part_name}..."
            install_payloads_on_part "${BOARD_DIR}/payloads" "$part_name" "p${num}"
        fi
    done
}

method_verify() {
    local target
    target=$(yq_read "$BOARD_CONFIG" target)

    local count i
    count=$(yq_count "$BOARD_CONFIG" partitions)

    for ((i=0; i<count; i++)); do
        local num install_autoflash
        num=$(yq_read "$BOARD_CONFIG" "partitions[$i].num")
        install_autoflash=$(yq_read "$BOARD_CONFIG" "partitions[$i].install_autoflash")

        [ "$install_autoflash" != "true" ] && continue

        local part_name mnt
        part_name=$(partition_name "$target" "$num")
        mnt=$(mktemp -d "/tmp/ukama-verify.XXXXXX")
        sudo mount "$part_name" "$mnt"

        if [ ! -e "$mnt/bin/sh" ] && [ ! -L "$mnt/bin/sh" ]; then
            echo "ERROR: $part_name missing rootfs content"
            sudo umount "$mnt"
            rm -rf "$mnt"
            return 1
        fi
        if [ ! -f "$mnt/usr/local/sbin/ukama-auto-flash.sh" ]; then
            echo "ERROR: autoflash payload missing on $part_name"
            sudo umount "$mnt"
            rm -rf "$mnt"
            return 1
        fi

        sudo umount "$mnt"
        rm -rf "$mnt"
    done
}

method_monitor() {
    local serial_dev success_marker boot_marker
    serial_dev=$(yq_read "$BOARD_CONFIG" serial.device)
    success_marker=$(yq_read "$BOARD_CONFIG" serial.success_marker)
    boot_marker=$(yq_read "$BOARD_CONFIG" serial.boot_marker)

    if [ ! -e "$serial_dev" ]; then
        echo "Serial $serial_dev not available — skipping monitor."
        echo "To monitor manually: screen $serial_dev $(yq_read "$BOARD_CONFIG" serial.baud)"
        return 0
    fi

    echo ""
    echo "Next steps:"
    echo "  1. Eject the SD card"
    echo "  2. Insert it into the ${BOARD} board"
    echo "  3. Power on the board"
    echo ""
    read -rp "Press ENTER once the board is powered on, or 's' to skip serial monitor: " resp
    [ "$resp" = "s" ] && return 0

    echo "Watching $serial_dev for '$success_marker'..."
    serial_wait_for_marker "$serial_dev" "${LOG_DIR}/serial.log" "$success_marker" 300 || {
        echo "Did not see success marker within timeout."
        return 1
    }
    echo "Flash completed on target."

    if [ -n "$boot_marker" ] && [ "$boot_marker" != "null" ]; then
        echo "Watching for boot marker '$boot_marker'..."
        serial_wait_for_marker "$serial_dev" "${LOG_DIR}/serial.log" "$boot_marker" 120 || {
            echo "Did not see boot marker within timeout."
            return 1
        }
        echo "Board booted."
    fi
}
