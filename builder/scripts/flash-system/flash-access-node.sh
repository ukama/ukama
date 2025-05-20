#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

# Flash and verify Ukama image for the access node (rpi4)
set -x
set -euo pipefail

REQUIRED_TOOLS=(dd lsblk grep timeout screen expect sudo lsusb make gcc libusb-1.0-0-dev git)

IMAGE_NAME="ukama-access-node.img"
UART_PORT="${1:-/dev/ttyUSB0}"
EXPECTED_BOOT_STRING="${2:-Ukama OS boot complete}"
SSH_IP="${3:-}"
SSH_USER="${4:-}"
SSH_PASS="${5:-}"

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

usage() {
    echo -e "\nUsage:"
    echo "  $0 <UART_PORT> [EXPECTED_BOOT_STRING] [SSH_IP] [SSH_USER] [SSH_PASS]"
    echo
    echo "Arguments:"
    echo "  UART_PORT             UART port connected to CM4 (e.g. /dev/ttyUSB0)"
    echo "  EXPECTED_BOOT_STRING  (Optional) String to detect in UART output (default: 'Ukama OS boot complete')"
    echo "  SSH_IP                (Optional) IP address of Pi for SSH verification"
    echo "  SSH_USER              (Optional) SSH username"
    echo "  SSH_PASS              (Optional) SSH password"
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
    log INFO "Waiting for mass storage device..."
    for i in {1..20}; do
        sleep 2
        BLOCK_DEVICE=$(lsblk -d -o NAME,SIZE,MODEL | grep -E 'sd[b-z]|mmcblk[0-9]' | awk '{print $1}' | tail -n1)
        if [[ -n "$BLOCK_DEVICE" && -b "/dev/$BLOCK_DEVICE" ]]; then
            BLOCK_DEVICE="/dev/$BLOCK_DEVICE"
            log SUCCESS "CM4 eMMC detected as $BLOCK_DEVICE"
            return
        fi
    done
    log ERROR "eMMC device not detected after rpiboot."
    exit 1
}

flash_image() {
    log INFO "Flashing image to $BLOCK_DEVICE..."
    sudo dd if="$IMAGE_NAME" of="$BLOCK_DEVICE" bs=$BLOCK_SIZE status=progress conv=fsync
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

test_ssh_connectivity() {
    if [[ -z "$SSH_IP" || -z "$SSH_USER" || -z "$SSH_PASS" ]]; then
        log INFO "Skipping SSH check — no credentials provided."
        return
    fi

    log INFO "Attempting SSH connection to $SSH_USER@$SSH_IP"
    expect <<EOF
        spawn ssh -o StrictHostKeyChecking=no $SSH_USER@$SSH_IP "echo SSH OK"
        expect {
            "yes/no" { send "yes\r"; exp_continue }
            "assword:" { send "$SSH_PASS\r" }
        }
        expect "SSH OK"
EOF

    if [[ $? -eq 0 ]]; then
        log SUCCESS "SSH check passed."
    else
        log ERROR "SSH failed."
        return 1
    fi
}

if [[ $# -lt 1 ]]; then
    usage
    exit 1
fi

check_dependencies
build_rpiboot
wait_for_cm4_usb
start_rpiboot

if [[ ! -f "$IMAGE_NAME" ]]; then
    log ERROR "$IMAGE_NAME not found"
    exit 1
fi

flash_image
mount_and_verify_boot_partition
monitor_uart_boot_log
test_ssh_connectivity

log SUCCESS "All steps completed successfully. CM4 is flashed and verified!"
