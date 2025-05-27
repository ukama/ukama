#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.
# Run the generated image with compiled kernel in QEMU (access - rpi)

set -euo pipefail
set -x

DEFAULT_IMG_NAME="ukama-access-node.img"
KERNEL_IMAGE="kernel.img"
RAM_SIZE="4G"
BUILD_DIR="$(pwd)"
TMP_BOOT="/tmp/boot-kernel"
LOOP_DEV=""
MAPPED_DEVS=()

log() {
    local type="$1"
    local message="$2"
    local color
    case "$type" in
        "INFO")    color="\033[1;34m";;
        "SUCCESS") color="\033[1;32m";;
        "ERROR")   color="\033[1;31m";;
        *)         color="\033[1;37m";;
    esac

    echo -e "${color}${type}: ${message}\033[0m"
}

error_exit() {
    log "ERROR" "$1"
    cleanup
    exit 1
}

cleanup() {
    log "INFO" "Cleaning up loop device and mounts"
    sync
    sudo fuser -k "$TMP_BOOT" 2>/dev/null || true
    sudo umount "$TMP_BOOT" 2>/dev/null || true
    sudo rm -rf "$TMP_BOOT"

    if [[ -n "$LOOP_DEV" ]]; then
        for dev in "${MAPPED_DEVS[@]}"; do
            sudo umount "/dev/mapper/${dev}" 2>/dev/null || true
        done
        sudo kpartx -dv "$LOOP_DEV" || true
        sudo losetup -d "$LOOP_DEV" || true
    fi
}

trap cleanup EXIT

IMG_NAME="${1:-$DEFAULT_IMG_NAME}"
[ -f "${BUILD_DIR}/${IMG_NAME}" ] || error_exit "Disk image '${IMG_NAME}' not found in ${BUILD_DIR}"

log "INFO" "Attaching image to loop device"
LOOP_DEV=$(sudo losetup --find --show "${BUILD_DIR}/${IMG_NAME}")
sudo kpartx -av "$LOOP_DEV" | while read -r _ name _; do
    MAPPED_DEVS+=("$name")
done

BOOT_PART="/dev/mapper/$(basename "${LOOP_DEV}")p1"
log "INFO" "Boot partition is $BOOT_PART"

log "INFO" "Mounting boot partition to extract kernel"
sudo mkdir -p "$TMP_BOOT"
sudo mount "$BOOT_PART" "$TMP_BOOT"

KERNEL_PATH_FIRMWARE="${TMP_BOOT}/firmware/kernel.img"
KERNEL_PATH_ROOT="${TMP_BOOT}/kernel.img"

if [[ -f "$KERNEL_PATH_FIRMWARE" ]]; then
    cp "$KERNEL_PATH_FIRMWARE" "${BUILD_DIR}/${KERNEL_IMAGE}"
    log "INFO" "Kernel extracted from firmware directory"
elif [[ -f "$KERNEL_PATH_ROOT" ]]; then
    cp "$KERNEL_PATH_ROOT" "${BUILD_DIR}/${KERNEL_IMAGE}"
    log "INFO" "Kernel extracted from root of boot directory"
else
    error_exit "kernel.img not found in either /boot/firmware/ or /boot/"
fi

sync
sudo fuser -k "$TMP_BOOT" 2>/dev/null || true
sudo umount "$TMP_BOOT"

log "INFO" "Starting QEMU with image '${IMG_NAME}'..."
qemu-system-aarch64 \
  -machine virt \
  -cpu cortex-a72 \
  -smp 4 \
  -m "${RAM_SIZE}" \
  -kernel "${BUILD_DIR}/${KERNEL_IMAGE}" \
  -append "root=/dev/vda5 rootfstype=ext4 rw panic=0 console=ttyAMA0" \
  -drive format=raw,file="${BUILD_DIR}/${IMG_NAME}",if=none,id=hd0,cache=writeback \
  -device virtio-blk,drive=hd0,bootindex=0 \
  -netdev user,id=mynet,hostfwd=tcp::2222-:22 \
  -device virtio-net-pci,netdev=mynet \
  -serial mon:stdio \
  -nographic \
  -monitor telnet:127.0.0.1:5555,server,nowait

log "SUCCESS" "QEMU exited cleanly."
