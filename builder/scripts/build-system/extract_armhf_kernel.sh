#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

set -e

ALPINE_VERSION="v3.17"
MINOR_VERSION="3"
ARCH="armhf"
MIRROR="http://dl-cdn.alpinelinux.org/alpine"
ROOTFS_DIR="/tmp/alpine-rootfs-armhf"
KERNEL_OUT="/tmp/alpine-kernel-armhf"
KERNEL_PKG="linux-rpi2"

# Clean old rootfs
sudo rm -rf "$ROOTFS_DIR" "$KERNEL_OUT"
mkdir -p "$ROOTFS_DIR" "$KERNEL_OUT"

# Download and extract Alpine minirootfs
curl -sSL "${MIRROR}/${ALPINE_VERSION}/releases/${ARCH}/alpine-minirootfs-${ALPINE_VERSION#v}.${MINOR_VERSION}-${ARCH}.tar.gz" | sudo tar -xz -C "$ROOTFS_DIR"

# Copy QEMU static binary
sudo cp /usr/bin/qemu-arm-static "$ROOTFS_DIR/usr/bin/"

# Setup basic config
echo "${MIRROR}/${ALPINE_VERSION}/main" | sudo tee "$ROOTFS_DIR/etc/apk/repositories"
echo "nameserver 8.8.8.8" | sudo tee "$ROOTFS_DIR/etc/resolv.conf"

# Chroot and install kernel
sudo chroot "$ROOTFS_DIR" /bin/sh -c "
    apk update &&
    apk add $KERNEL_PKG"

# Locate the kernel image (vmlinuz or zImage)
KERNEL_SRC=$(find "$ROOTFS_DIR/boot" -type f \( -name 'vmlinuz-*' -o -name 'zImage' \) | head -n1)

if [[ -z "$KERNEL_SRC" ]]; then
    echo "❌ No kernel image found in /boot/"
    exit 1
fi

# Extract kernel, dtbs, modules
sudo mkdir -p "$KERNEL_OUT/boot" "$KERNEL_OUT/lib/modules"
sudo cp "$KERNEL_SRC" "$KERNEL_OUT/boot/kernel.img"
sudo find "$ROOTFS_DIR/boot" -name '*.dtb' -exec cp --parents {} "$KERNEL_OUT/boot/" \;
sudo cp -a "$ROOTFS_DIR/lib/modules/"* "$KERNEL_OUT/lib/modules/"

echo "Kernel extracted to $KERNEL_OUT"

