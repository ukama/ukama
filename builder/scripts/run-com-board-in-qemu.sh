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
    sudo mount "${LOOPDEV}p3" "$MOUNTPOINT"

    ROOT_UUID=$(sudo blkid -s UUID -o value "${LOOPDEV}p3" || true)

    if [[ -n "$ROOT_UUID" ]]; then
        KERNEL_ARGS="root=UUID=${ROOT_UUID} rootfstype=ext4 rw console=ttyS0"
    else
        echo "WARNING: Could not determine UUID, falling back to /dev/sda3"
        KERNEL_ARGS="root=/dev/vda3 rootfstype=ext4 rw console=ttyS0"
    fi

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
        -drive file="$IMG",format=raw,if=none,id=virtio \
        -device ide-hd,drive=hd0,bus=ide.0 \
        -enable-kvm \
        -nographic \
        -serial mon:stdio \
        -boot menu=on \
        -net nic -net user,hostfwd=tcp::2222-:22
}

setup_grub
run_qemu
