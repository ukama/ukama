#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

# Run the generated disk image in QEMU

set -e

BUILD_DIR="./build_access_node"
IMG_NAME="access-node.img"
KERNEL_IMAGE="kernel8.img"
DTB_FILE="bcm2711-rpi-cm4.dtb"
RAM_SIZE="2G"

function log() {
    local type="$1"
    local message="$2"
    echo -e "\033[1;34m${type}:\033[0m ${message}"
}

function error_exit() {
    log "ERROR" "$1"
    exit 1
}

[ -f "$IMG_NAME" ]     || error_exit "Disk image '$IMG_NAME' not found"
[ -f "$KERNEL_IMAGE" ] || error_exit "Kernel '$KERNEL_IMAGE' not found"
[ -f "$DTB_FILE" ]     || error_exit "DTB '$DTB_FILE' not found"

log "INFO" "Starting QEMU for access node (rpi cm4) ..." 

if ! qemu-system-aarch64 \
    -M raspi3 \
    -cpu cortex-a72 \
    -m "$RAM_SIZE" \
    -kernel "${BUILD_DIR}/$KERNEL_IMAGE" \
    -dtb "${BUILD_DIR}/$DTB_FILE" \
    -drive file="$IMG_NAME",if=sd,format=raw \
    -append "rw root=/dev/mmcblk0p2 console=ttyAMA0" \
    -serial mon:stdio \
    -nographic \
    -netdev user,id=net0,hostfwd=tcp::2222-:22 \
    -device usb-net,netdev=net0; then
    error_exit "QEMU execution failed."
fi

log "SUCCESS" "QEMU started successfully."
