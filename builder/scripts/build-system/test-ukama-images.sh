#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

# Test integrity of one or more Ukama disk images

set -euo pipefail

IMG="ukama-com-image.img"
MOUNT_DIR="/mnt/testimg"
EXPECTED_PARTITIONS=(1 2 3 4)  # [boot, passive, primary (root), swap]

REQUIRED_FILES_PRIMARY=(
  "/sbin/starter.d"
  "/manifest.json"
  "/boot/vmlinuz"
)

REQUIRED_FILES_BOOT=(
  "/boot/bootcode.bin"  # optional
  "/boot/boot.bin"      # optional
)

log() {
    local type="$1"
    local message="$2"
    local color
    case "$type" in
        "INFO") color="\033[1;34m";;
        "PASS") color="\033[1;32m";;
        "ERROR") color="\033[1;31m";;
        *) color="\033[1;37m";;
    esac
    echo -e "${color}${type}: ${message}\033[0m"
}

fail() {
    log "ERROR" "$*"
    exit 1
}

cleanup() {
    log "INFO" "Cleaning up mounts and mappings"
    for p in "${EXPECTED_PARTITIONS[@]}"; do
        umount "$MOUNT_DIR/p$p" 2>/dev/null || true
    done
    if [[ -n "${LOOP_DEV:-}" ]]; then
        kpartx -dv "$LOOP_DEV" 2>/dev/null || true
        losetup -d "$LOOP_DEV" 2>/dev/null || true
    fi
    rm -rf "$MOUNT_DIR"
}
trap cleanup EXIT

attach_loop() {
    log "INFO" "Attaching image to loop device"
    LOOP_DEV=$(sudo losetup -f --show --partscan "$IMG")
    log "INFO" "Loop device: $LOOP_DEV"
}

check_layout() {
    log "INFO" "Partition table for $IMG:"
    sudo fdisk -l "$IMG"
    log "INFO" "Filesystems & labels:"
    lsblk -o NAME,FSTYPE,LABEL "$LOOP_DEV"*
}

check_labels() {
    log "INFO" "Verifying partition labels"
    declare -A expected_labels=( ["1"]="boot" ["2"]="passive" ["3"]="primary" ["4"]="swap" )

    for part in "${!expected_labels[@]}"; do
        dev="${LOOP_DEV}p${part}"
        expected="${expected_labels[$part]}"
        fstype=$(blkid -o value -s TYPE "$dev")
        if [[ "$fstype" == "vfat" ]]; then
            actual=$(sudo fatlabel "$dev" 2>/dev/null | tail -n 1)
        elif [[ "$fstype" == "ext4" ]]; then
            actual=$(sudo e2label "$dev" 2>/dev/null || echo "")
        elif [[ "$fstype" == "swap" ]]; then
            actual=$(blkid -s LABEL -o value "$dev" || echo "")
        else
            actual="unknown"
        fi

        if [[ "$actual" != "$expected" ]]; then
            log "ERROR" "$dev label mismatch: got '$actual', expected '$expected'"
            exit 1
        else
            log "INFO" "$dev label OK: $actual"
        fi
    done
}

check_fstab() {
    log "INFO" "Checking /etc/fstab in primary"
    cat "$MOUNT_DIR/p3/etc/fstab"
}

check_swap() {
    log "INFO" "Checking swap partition"
    sudo file -s "${LOOP_DEV}p4"
}

test_image() {
    local img="$1"
    log "INFO" "Testing image: $img"

    [[ -f "$img" ]] || fail "Image not found: $img"

    attach_loop
    check_layout

    mkdir -p "$MOUNT_DIR"
    for p in "${EXPECTED_PARTITIONS[@]}"; do
        dev="${LOOP_DEV}p$p"
        mount_path="$MOUNT_DIR/p$p"
        mkdir -p "$mount_path"
        if [[ $p -eq 4 ]]; then
            log "INFO" "Skipping mount for swap"
            continue
        fi
        log "INFO" "Mounting $dev to $mount_path"
        mount "$dev" "$mount_path" || fail "Failed to mount $dev"
    done

    for file in "${REQUIRED_FILES_PRIMARY[@]}"; do
        if [[ ! -f "$MOUNT_DIR/p3${file}" ]]; then
            fail "Missing required file in primary: $file"
        else
            log "INFO" "Found primary file: $file"
        fi
    done

    for file in "${REQUIRED_FILES_BOOT[@]}"; do
        if compgen -G "$MOUNT_DIR/p1${file}" > /dev/null; then
            log "INFO" "Found boot file: $file"
        else
            log "INFO" "Boot file not found (optional): $file"
        fi
    done

    check_labels
    check_fstab
    check_swap

    log "PASS" "$img passed all tests"
}

main() {
  [[ $# -lt 1 ]] && fail "Usage: $0 <img1> [img2 ...]"
  for img in "$@"; do
    test_image "$img"
  done
}

main "$@"
