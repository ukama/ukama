#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
# Copyright (c) 2025-present, Ukama Inc.

# This script creates a dual-partition USB:
# - Partition 1: bootable Alpine ISO (vanilla)
# - Partition 2: FAT32 with flash script + README

set -euo pipefail
set -x

: "${USB_DEV:?Must set USB_DEV}"
: "${FLASH_SCRIPT:?Must set FLASH_SCRIPT}"

BOOT_PART="${USB_DEV}1"
DATA_PART="${USB_DEV}2"
ISO_FILE="alpine.iso"

MNT_ISO="mnt-iso"
MNT_BOOT="mnt-boot"
MNT_DATA="mnt-data"

README_FILE="README.txt"

cleanup() {
    echo "Cleaning up..."
    sudo umount "$MNT_ISO"  || true
    sudo umount "$MNT_BOOT" || true
    sudo umount "$MNT_DATA" || true
    rm -rf "$MNT_ISO" "$MNT_BOOT" "$MNT_DATA"
}
trap cleanup EXIT

# Wipe and partition the USB
sudo wipefs -a "$USB_DEV"
sudo dd if=/dev/zero of="$USB_DEV" bs=1M count=10

echo "Partitioning USB..."
sudo parted --script "$USB_DEV" \
    mklabel msdos \
    mkpart primary fat32 1MiB 1024MiB \
    set 1 boot on \
    mkpart primary fat32 1024MiB 100%

sudo mkfs.vfat -F 32 -n BOOT "$BOOT_PART"
sudo mkfs.vfat -F 32 -n DATA "$DATA_PART"

# Mount everything
mkdir -p "$MNT_ISO" "$MNT_BOOT" "$MNT_DATA"

sudo mount -o loop "$ISO_FILE" "$MNT_ISO"
sudo mount "$BOOT_PART" "$MNT_BOOT"
sudo mount "$DATA_PART" "$MNT_DATA"

# Copy Alpine ISO to BOOT partition
sudo rsync -a "$MNT_ISO"/ "$MNT_BOOT"/

# Move flash script + README to DATA partition
sudo mv "$FLASH_SCRIPT" "$MNT_DATA/$FLASH_SCRIPT"
sudo chmod +x "$MNT_DATA/flash-smarc.sh"

cat <<EOF | sudo tee "$MNT_DATA/$README_FILE" >/dev/null
Ukama SMARC Flash USB
======================

To flash the image onto the internal eMMC:

1. Boot the SMARC board using this USB (Partition 1 boots Alpine).
2. Log into Alpine as 'root' (no password).
3. Mount the second partition:

   mount /dev/sda2 /mnt

4. Run the flashing script:

   /mnt/flash-smarc.sh

5. Wait for reboot after success.

NOTE: This method requires network connection to the build server.
EOF

sync

echo "Dual-partition USB created. Bootable Alpine + flash script in place."
