#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
# Copyright (c) 2025-present, Ukama Inc.
#
# Creates a dual-partition device (USB or SD):
#   Partition 1: FAT32 “BOOT” with Alpine ISO
#   Partition 2: ext4 “DATA” with flash script + README

set -euo pipefail

: "${DEV:?Must set DEV (e.g. /dev/sdb or /dev/mmcblk0)}"
: "${FLASH_SCRIPT:?Must set FLASH_SCRIPT (e.g. flash-smarc.sh)}"
: "${BOARD_NAME:?Must set BOARD_NAME (e.g. SMARC or FEM-Control)}"
: "${ISO_FILE:=alpine.iso}"

# detect if partitions need a "p" (e.g. mmcblk0 → mmcblk0p1)
base=$(basename "$DEV")
if [[ "$base" =~ [0-9]$ ]]; then
    suffix="p"
else
    suffix=""
fi

BOOT_PART="${DEV}${suffix}1"
DATA_PART="${DEV}${suffix}2"

MNT_ISO="mnt-iso"
MNT_BOOT="mnt-boot"
MNT_DATA="mnt-data"
README_FILE="README.txt"

cleanup() {
    echo "Cleaning up mounts…"
    sudo umount "$MNT_ISO"  || true
    sudo umount "$MNT_BOOT" || true
    sudo umount "$MNT_DATA" || true
    rm -rf "$MNT_ISO" "$MNT_BOOT" "$MNT_DATA"
}
trap cleanup EXIT

echo "Wiping $DEV…"
sudo wipefs -a "$DEV"
sudo dd if=/dev/zero of="$DEV" bs=1M count=10 status=none

echo "Partitioning $DEV…"
sudo parted --script "$DEV" \
     mklabel msdos \
     mkpart primary fat32 1MiB 1024MiB \
     set 1 boot on \
     mkpart primary ext4 1024MiB 100%

echo "Formatting…"
sudo mkfs.vfat -F 32 -n BOOT "$BOOT_PART"
sudo mkfs.ext4 -L DATA "$DATA_PART"

echo "Mounting…"
mkdir -p "$MNT_ISO" "$MNT_BOOT" "$MNT_DATA"
sudo mount -o loop "$ISO_FILE" "$MNT_ISO"
sudo mount "$BOOT_PART" "$MNT_BOOT"
sudo mount "$DATA_PART" "$MNT_DATA"

echo "Copying Alpine-ISO → BOOT…"
sudo rsync -a --no-owner "$MNT_ISO"/ "$MNT_BOOT"/

echo "Installing flash script + README → DATA…"
sudo cp "$FLASH_SCRIPT" "$MNT_DATA/"
sudo chmod +x "$MNT_DATA/$FLASH_SCRIPT"

cat <<EOF | sudo tee "$MNT_DATA/$README_FILE" >/dev/null
Ukama ${BOARD_NAME} Flash Media
===============================

To flash the image onto the internal eMMC of your ${BOARD_NAME} board:

1. Boot the board from this media (Partition 1 is Alpine).
2. Log in as 'root' (no password).
3. Mount the DATA partition:
     mount -t ext4 /dev/$(basename "$DEV")${suffix}2 /mnt
4. Run the flashing script:
     /mnt/$(basename "$FLASH_SCRIPT")

NOTE: Requires network access to your build server.
EOF

sync
echo "Created dual-partition media for ${BOARD_NAME}: Alpine boot (FAT32) + flash script (ext4)."
