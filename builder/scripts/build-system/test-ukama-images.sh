#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

# Test integrity of one or more Ukama disk images

set -euo pipefail

MOUNT_DIR="/mnt/testimg"
TMP_DIR="/tmp/testimg"
EXPECTED_PARTITIONS=(1 2 5 6 7 8)  # boot, recovery, primary, passive, data, swap
REQUIRED_FILES_PRIMARY=(
  "/boot/kernel.img"
  "/sbin/starter.d"
  "/manifest.json"
)
REQUIRED_FILES_BOOT=(
  "/boot/bootcode.bin"  # for access node (optional)
  "/boot/boot.bin"      # for amplifier node
)

log() {
    local type="$1"; shift
    echo -e "\033[1;36m[$type]\033[0m $*"
}

fail() {
    echo -e "\033[1;31m[FAIL]\033[0m $*"
    exit 1
}

cleanup() {
    log "INFO" "Cleaning up mounts and mappings"
    for p in "${EXPECTED_PARTITIONS[@]}"; do
        sudo umount "$MOUNT_DIR/p$p" 2>/dev/null || true
    done
    [ -n "${LOOP_DEV:-}" ] && sudo kpartx -dv "$LOOP_DEV" && sudo losetup -d "$LOOP_DEV"
    sudo rm -rf "$MOUNT_DIR"
}
trap cleanup EXIT

test_image() {
    local img="$1"
    log "INFO" "Testing image: $img"

    [ -f "$img" ] || fail "Image not found: $img"
    [ $(stat -c %s "$img") -gt 1048576 ] || fail "Image too small: $img"

    LOOP_DEV=$(sudo losetup -f --show "$img")
    sudo kpartx -av "$LOOP_DEV" >/dev/null

    DEVICE=$(basename "$LOOP_DEV")
    mkdir -p "$MOUNT_DIR"

    for p in "${EXPECTED_PARTITIONS[@]}"; do
        local part_dev="/dev/mapper/${DEVICE}p${p}"
        local part_mount="$MOUNT_DIR/p${p}"
        mkdir -p "$part_mount"
        log "INFO" "Mounting $part_dev to $part_mount"
        if ! sudo mount "$part_dev" "$part_mount"; then
            [[ $p == 8 ]] && log "INFO" "Skipping mount test for swap partition" && continue
            fail "Failed to mount partition $p"
        fi
    done

    # Check boot partition files
    local bootdir="$MOUNT_DIR/p1"
    for file in "${REQUIRED_FILES_BOOT[@]}"; do
        if compgen -G "$bootdir$file" > /dev/null; then
            log "INFO" "Found boot file: $file"
        else
            log "WARN" "Boot file not found (optional): $file"
        fi
    done

    # Check primary rootfs files
    local primarydir="$MOUNT_DIR/p5"
    for file in "${REQUIRED_FILES_PRIMARY[@]}"; do
        if [[ ! -f "$primarydir$file" ]]; then
            fail "Missing required file in primary: $file"
        else
            log "INFO" "Found primary file: $file"
        fi
    done

    log "PASS" "$img passed all tests"
}

main() {
    [[ $# -lt 1 ]] && fail "Usage: $0 <img1> [img2 ...]"
    for img in "$@"; do
        test_image "$img"
        cleanup  # needed between runs
    done
}

main "$@"
