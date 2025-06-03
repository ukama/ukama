#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

# Flash and verify Ukama image for the access node (rpi4)
set -euo pipefail

trap cleanup EXIT

REQUIRED_TOOLS=(dd lsblk grep timeout screen expect sudo lsusb make gcc libusb-1.0-0-dev git)

IMAGE_NAME="ukama-access-node.img"
UART_PORT="/dev/ttyUSB0}"
EXPECTED_BOOT_STRING="Ukama OS boot complete}"

BAUD_RATE=115200
UART_LOG="/tmp/rpi_uart.log"
MOUNT_POINT="/mnt/rpi"
BLOCK_DEVICE=""
BLOCK_SIZE=4M
RPIBOOT_DIR="/tmp/rpiboot"

log() {
    local type="$1"; shift
    local msg="$*"
    case "$type" in
        INFO) color="\033[1;34m";;
        SUCCESS) color="\033[1;32m";;
        ERROR) color="\033[1;31m";;
        *) color="\033[0m";;
    esac
    echo -e "${color}${type}: ${msg}\033[0m"
}

cleanup() {
    log INFO "Cleaning up..."
    sudo pkill -f "cat $UART_PORT" || true
    sudo umount "$MOUNT_POINT" 2>/dev/null || true
}

usage() {
    echo -e "\nUsage:"
    echo "  $0 [--verify]"
    echo
    echo "Modes:"
    echo "  (no args)       Flash image to CM4 eMMC via rpiboot"
    echo "  --verify        After power cycle, verify UART boot log"
    echo
    exit 1
}

detect_package_manager() {
    for cmd in apt dnf yum pacman apk; do
        if command -v "$cmd" &>/dev/null; then
            echo "$cmd"
            return
        fi
    done
    log ERROR "No supported package manager found."
    exit 1
}

check_and_install_dependencies() {
    local missing=()
    for tool in "${REQUIRED_TOOLS[@]}"; do
        if ! command -v "$tool" &>/dev/null; then
            missing+=("$tool")
        fi
    done

    if [ ${#missing[@]} -eq 0 ]; then
        log SUCCESS "All required tools are installed."
        return
    fi

    log ERROR "Missing tools: ${missing[*]}"

    read -rp "Do you want to install the missing tools? (yes/[no]): " confirm
    if [[ "$confirm" != "yes" ]]; then
        log ERROR "Cannot proceed without required tools."
        exit 1
    fi

    local pkgmgr
    pkgmgr=$(detect_package_manager)
    log INFO "Using package manager: $pkgmgr"

    case "$pkgmgr" in
        apt)
            sudo apt update
            sudo apt install -y "${missing[@]}"
            ;;
        dnf|yum)
            sudo "$pkgmgr" install -y "${missing[@]}"
            ;;
        apk)
            sudo apk add "${missing[@]}"
            ;;
        pacman)
            sudo pacman -Sy --noconfirm "${missing[@]}"
            ;;
        *)
            log ERROR "Unsupported package manager: $pkgmgr"
            exit 1
            ;;
    esac

    log SUCCESS "Installed missing tools."
}

build_rpiboot() {
    if [[ ! -x "$RPIBOOT_DIR/rpiboot" ]]; then
        log INFO "Building rpiboot..."
        rm -rf "$RPIBOOT_DIR"
        git clone --depth=1 https://github.com/raspberrypi/usbboot "$RPIBOOT_DIR"
        pushd "$RPIBOOT_DIR" >/dev/null
        make
        popd >/dev/null
        log SUCCESS "rpiboot built successfully."
    fi
}

wait_for_cm4_usb() {
    log INFO "Waiting for CM4 in USB boot mode..."
    until lsusb | grep -q "Broadcom.*BCM2711 Boot"; do
        sleep 1
    done
    log SUCCESS "CM4 detected in USB boot mode."
}

start_rpiboot() {
    log INFO "Uploading bootloader to CM4 (rpiboot)..."
    sudo "$RPIBOOT_DIR/rpiboot" >/dev/null &

    log INFO "Waiting for rpiboot to complete..."
    wait

    log INFO "Waiting for eMMC device to appear via dmesg..."
    for i in {1..30}; do
        log INFO "Sleeping for 5 seconds ....."
        sleep 5

        dev=$(sudo dmesg | tac | grep -m1 -oE 'sd[b-z]' | head -n1 || true)
        log INFO "DEBUG: dev=$dev"

        if [[ -b "/dev/$dev" ]]; then
            log INFO "/dev/$dev exists as block"
        else
            log ERROR "/dev/$dev does NOT exist yet"
        fi

        if [[ -n "$dev" && -b "/dev/$dev" ]]; then
            BLOCK_DEVICE="/dev/$dev"
            log SUCCESS "CM4 eMMC detected as $BLOCK_DEVICE"
            return
        fi
    done

    log ERROR "eMMC device not detected after rpiboot."
    sudo dmesg | tail -30
    exit 1
}

flash_image() {
    log INFO "Flashing image to $BLOCK_DEVICE..."
    if command -v pv &>/dev/null; then
        log INFO "Using pv for flashing progress..."
        sudo pv "$IMAGE_NAME" | sudo dd of="$BLOCK_DEVICE" bs=$BLOCK_SIZE conv=fsync
    else
        sudo dd if="$IMAGE_NAME" of="$BLOCK_DEVICE" bs=$BLOCK_SIZE status=progress conv=fsync
    fi
    sync
    log SUCCESS "Image flashed to $BLOCK_DEVICE"
}

mount_and_verify_boot_partition() {
    log INFO "Verifying boot partition..."
    BOOT_PART="${BLOCK_DEVICE}1"
    [[ "$BLOCK_DEVICE" =~ mmcblk ]] && BOOT_PART="${BLOCK_DEVICE}p1"

    sudo mkdir -p "$MOUNT_POINT"
    sudo mount "$BOOT_PART" "$MOUNT_POINT"
    if [ -f "$MOUNT_POINT/config.txt" ]; then
        log SUCCESS "Boot partition verified."
    else
        log ERROR "Boot partition is missing config.txt!"
        sudo umount "$MOUNT_POINT"
        exit 1
    fi
    sudo umount "$MOUNT_POINT"
}

monitor_uart_boot_log() {
    log INFO "Monitoring UART on $UART_PORT for boot log..."
    sudo timeout 120s cat "$UART_PORT" > "$UART_LOG" &
    sleep 10
    for i in {1..24}; do
        if grep -q "$EXPECTED_BOOT_STRING" "$UART_LOG"; then
            log SUCCESS "Detected boot message: '$EXPECTED_BOOT_STRING'"
            return 0
        fi
        sleep 5
    done
    log ERROR "Expected boot string not found in UART log."
    tail "$UART_LOG"
    return 1
}

if [[ $# -ge 1 && "$1" == "--verify" ]]; then
    monitor_uart_boot_log
    log SUCCESS "Verification complete. CM4 boot confirmed."
    exit 0
fi

rm -rf "$UART_LOG"

check_and_install_dependencies
build_rpiboot
wait_for_cm4_usb
start_rpiboot

if [[ ! -f "$IMAGE_NAME" ]]; then
    log ERROR "$IMAGE_NAME not found"
    exit 1
fi

flash_image
mount_and_verify_boot_partition

log SUCCESS "All steps completed successfully. CM4 is flashed and verified!"
log INFO "\nNext steps:"
echo "  1. Power off the CM4"
echo "  2. Remove BOOT jumper or switch it off"
echo "  3. Power it back on"
echo "  4. Once booted, run this script with \"--verify\" to test UART boot log."

