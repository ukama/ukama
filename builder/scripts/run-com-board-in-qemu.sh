#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

set -euo pipefail
set -x

IMG="ukama-com-image.img"
MOUNTPOINT="/mnt/ukama-test"
GRUB_CFG="${MOUNTPOINT}/boot/grub/grub.cfg"

setup_grub() {
    echo "Setting up GRUB bootloader..."

    LOOPDEV=$(sudo losetup --find --partscan --show "$IMG")
    sudo mkdir -p "$MOUNTPOINT"

    ROOT_PART=""
    for i in {1..4}; do
        PART_DEV="${LOOPDEV}p${i}"
        sudo mount "$PART_DEV" "$MOUNTPOINT" || continue

        if [ -f "$MOUNTPOINT/sbin/init" ] || [ -f "$MOUNTPOINT/init" ]; then
            ROOT_PART="/dev/sda${i}"
            echo "Found root partition: $ROOT_PART"
            break
        fi

        sudo umount "$MOUNTPOINT"
    done

    if [ -z "$ROOT_PART" ]; then
        echo "ERROR: Could not detect root partition with /sbin/init"
        sudo losetup -d "$LOOPDEV"
        exit 1
    fi

    KERNEL_ARGS="root=${ROOT_PART} rootfstype=ext4 rw console=ttyS0 debug ignore_loglevel"

    sudo mkdir -p "${MOUNTPOINT}/boot/grub"
    sudo grub-install \
        --target=i386-pc \
        --boot-directory="${MOUNTPOINT}/boot" \
        --modules="normal part_msdos ext2 multiboot" \
        "$LOOPDEV"

    cat <<EOF | sudo tee "${GRUB_CFG}"
set timeout=5
set default=0

menuentry "Ukama Alpine x86_64" {
    linux /boot/vmlinuz ${KERNEL_ARGS}
}
EOF

    sudo umount "$MOUNTPOINT"
    sudo losetup -d "$LOOPDEV"
}

run_qemu() {
    echo "Booting image in QEMU..."

    qemu-system-x86_64 \
        -m 2048 \
        -kernel /boot/vmlinuz \
        -append "root=/dev/sda2 rootfstype=ext4 rw console=ttyS0" \
        -drive file="$IMG",format=raw,if=none,id=hd0 \
        -device ide-hd,drive=hd0,bus=ide.0 \
        -nographic \
        -serial mon:stdio \
        -enable-kvm
}

#setup_grub
run_qemu

