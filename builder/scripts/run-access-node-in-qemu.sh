#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

# Run the generated image with compiled kernel in QEMU (rpi)

set -e
set -x

BUILD_DIR="./build_access_node"
IMG_NAME="access-node.img"
KERNEL_IMAGE="kernel.img"
RAM_SIZE="4G"

function log() {
    local type="$1"
    local message="$2"
    echo -e "\033[1;34m${type}:\033[0m ${message}"
}

function error_exit() {
    log "ERROR" "$1"
    exit 1
}

[ -f "${BUILD_DIR}/${IMG_NAME}" ]     || error_exit "Disk image '$IMG_NAME' not found"
[ -f "${BUILD_DIR}/${KERNEL_IMAGE}" ] || error_exit "Kernel '$KERNEL_IMAGE' not found"

log "INFO" "Starting QEMU for access node (rpi cm4) ..." 

if ! qemu-system-aarch64 \
     -machine virt \
     -cpu cortex-a72 \
     -smp 4 \
     -m "${RAM_SIZE}" \
     -kernel "${BUILD_DIR}/${KERNEL_IMAGE}" \
     -append "root=/dev/vda2 rootfstype=ext4 rw panic=0 console=ttyAMA0" \
     -drive format=raw,file="${BUILD_DIR}/${IMG_NAME}",if=none,id=hd0,cache=writeback \
     -device virtio-blk,drive=hd0,bootindex=0 \
     -netdev user,id=mynet,hostfwd=tcp::2222-:22 \
     -device virtio-net-pci,netdev=mynet \
     -monitor telnet:127.0.0.1:5555,server,nowait; then
    error_exit "QEMU execution failed."
fi

log "SUCCESS" "QEMU started successfully."
