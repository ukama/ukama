#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

# Flash and verify Ukama image for the access node (rpi4)

set -euo pipefail

REQUIRED_TOOLS=(dd lsblk grep timeout screen expect sudo)

IMAGE_NAME="ukama-access-node.img"
DEVICE=${1:-}
UART_PORT=${2:-/dev/ttyUSB0}
EXPECTED_BOOT_STRING=${3:-"Ukama OS boot complete"}
BAUD_RATE=115200
UART_LOG="/tmp/rpi_uart.log"
SSH_IP=${4:-}
SSH_USER=${5:-}
SSH_PASS=${6:-}
MOUNT_POINT="/mnt/rpi"
BLOCK_SIZE=4M

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
    echo "  $0 <DEVICE> <UART_PORT> (EXPECTED_BOOT_STRING) [SSH_IP] [SSH_USER] [SSH_PASS]"
    echo
    echo "Arguments:"
    echo "  DEVICE                Target block device (e.g. /dev/sdX or /dev/mmcblk0)"
    echo "  UART_PORT             UART port connected to CM4 (e.g. /dev/ttyUSB0)"
    echo "  EXPECTED_BOOT_STRING  (Optional) String to detect in UART output (e.g. 'UkamaOS boot complete')"
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

validate_device() {
    if [[ -z "$DEVICE" || ! -b "$DEVICE" ]]; then
        log ERROR "Block device $DEVICE is not valid or not specified."
        exit 1
    fi
    log INFO "Flashing to $DEVICE"
}

flash_image() {
    log INFO "Flashing image to $DEVICE ..."
    sudo dd if="$IMAGE_NAME" of="$DEVICE" bs=$BLOCK_SIZE status=progress conv=fsync
    sync
    log SUCCESS "Image flashed to $DEVICE"
}

mount_and_verify_boot() {
    log INFO "Verifying boot partition structure ..."
    BOOT_PART="${DEVICE}1"
    [[ "$DEVICE" =~ mmcblk ]] && BOOT_PART="${DEVICE}p1"

    mkdir -p "$MOUNT_POINT"
    sudo mount "$BOOT_PART" "$MOUNT_POINT"
    if [ -f "$MOUNT_POINT/config.txt" ]; then
        log SUCCESS "Boot partition looks good."
    else
        log ERROR "Boot partition seems invalid or missing files."
        sudo umount "$MOUNT_POINT"
        exit 1
    fi
    sudo umount "$MOUNT_POINT"
}

monitor_uart_boot_log() {
    log INFO "Waiting for RPi boot via UART log on $UART_PORT ..."
    sudo timeout 120s cat "$UART_PORT" > "$UART_LOG" &
    sleep 10  # Give it time to boot

    for i in {1..24}; do
        if grep -q "$EXPECTED_BOOT_STRING" "$UART_LOG"; then
            log SUCCESS "Detected expected boot message: \"$EXPECTED_BOOT_STRING\""
            return 0
        fi
        sleep 5
    done

    log ERROR "Boot log does not contain expected string after timeout."
    tail "$UART_LOG"
    return 1
}

test_ssh_connectivity() {
    if [[ -z "$SSH_IP" || -z "$SSH_USER" || -z "$SSH_PASS" ]]; then
        log INFO "SSH test skipped — IP or credentials not provided."
        return
    fi

    log INFO "Testing SSH connection to $SSH_USER@$SSH_IP"

    expect <<EOF
        spawn ssh -o StrictHostKeyChecking=no $SSH_USER@$SSH_IP "echo 'SSH OK'"
        expect {
            "yes/no" { send "yes\r"; exp_continue }
            "assword:" { send "$SSH_PASS\r" }
        }
        expect "SSH OK"
EOF

    if [[ $? -eq 0 ]]; then
        log SUCCESS "SSH test passed."
    else
        log ERROR "SSH test failed."
        return 1
    fi
}

# Main
if [[ $# -lt 2 ]]; then
    log ERROR "Not enough arguments."
    usage
fi

check_and_install_dependencies
validate_device

if [[ ! -f "$IMAGE_NAME" ]]; then
    log ERROR "Image file $IMAGE_NAME not found."
    exit 1
fi

flash_image
mount_and_verify_boot
monitor_uart_boot_log
test_ssh_connectivity

log SUCCESS "All steps completed successfully."
