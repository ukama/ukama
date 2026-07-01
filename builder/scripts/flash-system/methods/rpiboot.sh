#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

source "${LIB_DIR}/serial.sh"

RPIBOOT_DIR="/tmp/rpiboot"
MOUNT_POINT="/mnt/rpi"
BLOCK_DEVICE=""

REQUIRED_TOOLS=(dd lsblk grep timeout sudo lsusb make gcc git wget)
HEADER_TO_PACKAGE=("/usr/include/libusb-1.0/libusb.h:libusb-1.0-0-dev")

method_validate() {
    local image_path
    image_path=$(yq_read "$BOARD_CONFIG" image.path)

    if [ ! -f "$image_path" ]; then
        echo "  [FAIL] image not found at $image_path"
        return 1
    fi

    if declare -F detect_image_format >/dev/null; then
        local fmt
        fmt=$(detect_image_format "$image_path" 2>/dev/null || true)
        if [ -z "$fmt" ]; then
            echo "  [FAIL] could not detect image format for $image_path"
            return 1
        fi
        echo "  [ OK ] image: $image_path ($fmt)"
    else
        echo "  [ OK ] image: $image_path"
    fi

    local missing=()
    local tool
    for tool in "${REQUIRED_TOOLS[@]}"; do
        command -v "$tool" >/dev/null 2>&1 || missing+=("$tool")
    done

    local mapping
    for mapping in "${HEADER_TO_PACKAGE[@]}"; do
        local header="${mapping%%:*}"
        local pkg="${mapping##*:}"
        [ -f "$header" ] || missing+=("$pkg")
    done

    if [ "${#missing[@]}" -gt 0 ]; then
        echo "  [WARN] missing prerequisites: ${missing[*]}"
        echo "         install with: sudo apt install ${missing[*]}"
    fi

    return 0
}

method_confirm() {
    echo ""
    echo "Plan:"
    echo "  - build rpiboot from source if needed"
    echo "  - wait for CM4 to appear in USB boot mode"
    echo "  - dd image to eMMC over USB"
    echo ""
    read -rp "Type 'yes' to continue: " confirm
    [ "$confirm" = "yes" ]
}

_build_rpiboot() {
    if [ -x "$RPIBOOT_DIR/rpiboot" ]; then
        return 0
    fi
    rm -rf "$RPIBOOT_DIR"
    git clone --depth=1 https://github.com/raspberrypi/usbboot "$RPIBOOT_DIR"
    pushd "$RPIBOOT_DIR" >/dev/null
    make
    popd >/dev/null
}

_wait_for_cm4() {
    echo "Waiting for CM4 in USB boot mode..."
    local elapsed=0
    while [ "$elapsed" -lt 60 ]; do
        if lsusb | grep -q "Broadcom.*BCM2711 Boot"; then
            return 0
        fi
        sleep 1
        elapsed=$((elapsed + 1))
    done
    echo "ERROR: CM4 not detected in USB boot mode after 60s"
    return 1
}

_detect_emmc() {
    echo "Running rpiboot..."
    sudo "$RPIBOOT_DIR/rpiboot" >/dev/null &
    wait

    local i dev
    for i in $(seq 1 30); do
        sleep 5
        dev=$(sudo dmesg | tac | grep -m1 -oE 'sd[b-z]' | head -n1 || true)
        if [ -n "$dev" ] && [ -b "/dev/$dev" ]; then
            BLOCK_DEVICE="/dev/$dev"
            echo "CM4 eMMC detected as $BLOCK_DEVICE"
            return 0
        fi
    done
    echo "ERROR: eMMC device not detected after rpiboot"
    sudo dmesg | tail -30
    return 1
}

method_apply() {
    local image_path img_name
    image_path=$(yq_read "$BOARD_CONFIG" image.path)
    img_name=$(yq_read "$BOARD_CONFIG" image.name)
    [ "$img_name" = "null" ] && img_name=$(basename "$image_path")

    _build_rpiboot
    _wait_for_cm4
    _detect_emmc

    echo "Flashing $image_path to $BLOCK_DEVICE..."
    if command -v pv >/dev/null 2>&1; then
        sudo pv "$image_path" | sudo dd of="$BLOCK_DEVICE" bs=4M conv=fsync
    else
        sudo dd if="$image_path" of="$BLOCK_DEVICE" bs=4M status=progress conv=fsync
    fi
    sync
    echo "Image flashed."
}

method_verify() {
    [ -z "$BLOCK_DEVICE" ] && return 0

    local boot_part="${BLOCK_DEVICE}1"
    [[ "$BLOCK_DEVICE" =~ mmcblk ]] && boot_part="${BLOCK_DEVICE}p1"

    sudo mkdir -p "$MOUNT_POINT"
    sudo mount "$boot_part" "$MOUNT_POINT"
    if [ ! -f "$MOUNT_POINT/config.txt" ]; then
        echo "ERROR: boot partition missing config.txt"
        sudo umount "$MOUNT_POINT"
        return 1
    fi
    echo "Boot partition verified."
    sudo umount "$MOUNT_POINT"
}

method_monitor() {
    local serial_dev boot_marker
    serial_dev=$(yq_read "$BOARD_CONFIG" serial.device)
    boot_marker=$(yq_read "$BOARD_CONFIG" serial.boot_marker)

    if [ ! -e "$serial_dev" ]; then
        echo "Serial $serial_dev not available — skipping monitor."
        return 0
    fi

    echo ""
    echo "Next steps:"
    echo "  1. Power off the CM4"
    echo "  2. Remove BOOT jumper / switch it off"
    echo "  3. Power it back on"
    echo ""
    read -rp "Press ENTER once the board is powered on, or 's' to skip: " resp
    [ "$resp" = "s" ] && return 0

    serial_wait_for_marker "$serial_dev" "${LOG_DIR}/uart.log" "$boot_marker" 180 || {
        echo "Did not see boot marker within timeout."
        return 1
    }
    echo "CM4 boot confirmed."
}
