#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.
#
# Creates SD card for SD-boot boards by writing raw image and adding auto-flash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
IMAGE_UTILS="${SCRIPT_DIR}/image_utils.sh"

: "${DEV:?Must set DEV (e.g., /dev/sdb or /dev/mmcblk0)}"
: "${IMAGE_PATH:?Must set IMAGE_PATH (path to OS image)}"
: "${BOARD_NAME:?Must set BOARD_NAME (e.g., anode)}"

if [ ! -f "$IMAGE_UTILS" ]; then
    echo "ERROR: Missing $IMAGE_UTILS"
    exit 1
fi

source "$IMAGE_UTILS"

log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $*"
}

verify_device() {
    if [ ! -b "$DEV" ]; then
        log "ERROR: Device '$DEV' not found or is not a block device"
        log "Available devices:"
        lsblk
        exit 1
    fi

    log "WARNING: This will erase all data on $DEV"
    log "Device info:"
    lsblk "$DEV"

    read -rp "Are you sure you want to continue? (yes/no): " confirm
    if [[ "$confirm" != "yes" ]]; then
        log "Aborted by user"
        exit 1
    fi
}

flash_raw_image() {
    local image_format=""

    log "Flashing raw image to $DEV..."
    log "Image: $IMAGE_PATH"
    log "Target: $DEV"

    if [ ! -f "$IMAGE_PATH" ]; then
        log "ERROR: Image not found: $IMAGE_PATH"
        exit 1
    fi

    if ! image_format=$(detect_image_format "$IMAGE_PATH"); then
        log "ERROR: Unsupported image format: $IMAGE_PATH"
        log "Expected a raw disk image or a gzip/xz/zstd/bzip2-compressed raw image"
        exit 1
    fi

    if [ "$image_format" != "raw" ]; then
        log "Detected $image_format-compressed image; flashing decompressed contents"
    fi

    # Get image size
    IMG_SIZE=$(stat -c%s "$IMAGE_PATH")
    log "Image size: $IMG_SIZE bytes ($(numfmt --to=iec $IMG_SIZE))"

    # Flash with progress
    log "Writing image... (this may take several minutes)"
    case "$image_format" in
        raw)
            if command -v pv &>/dev/null; then
                pv "$IMAGE_PATH" | sudo dd of="$DEV" bs=4M conv=fsync
            else
                sudo dd if="$IMAGE_PATH" of="$DEV" bs=4M status=progress conv=fsync
            fi
            ;;
        *)
            stream_image_to_stdout "$IMAGE_PATH" "$image_format" | sudo dd of="$DEV" bs=4M status=progress conv=fsync
            ;;
    esac

    sync
    log "Raw image flashed successfully"

    # Re-read partition table so new partitions are visible
    log "Re-reading partition table..."
    sudo partprobe "$DEV" 2>/dev/null || sudo blockdev --rereadpt "$DEV" || true
    sleep 2

    # Verify partitions exist
    if ! lsblk "$DEV" | grep -q "part"; then
        log "WARNING: No partitions found on $DEV after flashing"
        log "You may need to manually run: sudo partprobe $DEV"
    fi
}

add_auto_flash_script() {
    log "Adding auto-flash functionality..."

    # Detect partition suffix
    local part_suffix=""
    if [[ "$(basename "$DEV")" =~ [0-9]$ ]]; then
        part_suffix="p"
    fi

    # Find boot partition (usually first partition)
    local boot_part="${DEV}${part_suffix}1"
    local mnt_point="/tmp/sd-boot-$$"

    if [ ! -b "$boot_part" ]; then
        log "WARNING: Boot partition $boot_part not found"
        log "Partitions on $DEV:"
        lsblk "$DEV"
        return 1
    fi

    log "Mounting boot partition $boot_part..."
    mkdir -p "$mnt_point"
    sudo mount "$boot_part" "$mnt_point"

    # Create auto-flash script that runs on first boot
    log "Creating auto-flash script..."
    sudo tee "$mnt_point/ukama-auto-flash.sh" > /dev/null << 'AUTOFLASH_EOF'
#!/bin/bash
# Ukama Auto-Flash Script
# This script runs on first boot to copy SD card contents to eMMC

set -e

LOG_FILE="/var/log/ukama-autoflash.log"
exec > >(tee -a "$LOG_FILE")
exec 2>&1

echo "========================================"
echo "  Ukama Auto-Flash to eMMC"
echo "  $(date)"
echo "========================================"

# Find eMMC device
echo "Detecting eMMC device..."
EMMC_DEV=""
for dev in /dev/mmcblk*; do
    if [ -e "${dev}boot0" ] && [ -e "${dev}boot1" ]; then
        # Skip the device we're booted from
        if ! mount | grep -q "^$dev"; then
            EMMC_DEV="$dev"
            break
        fi
    fi
done

if [ -z "$EMMC_DEV" ]; then
    echo "ERROR: No eMMC device found"
    exit 1
fi

echo "Found eMMC device: $EMMC_DEV"
echo "This will copy SD card contents to eMMC"
echo "Starting in 5 seconds... (Ctrl+C to cancel)"
sleep 5

# Find SD card device (current root)
SD_DEV=$(findmnt -n -o SOURCE / | sed 's/p[0-9]*$//')
echo "SD card device: $SD_DEV"
echo "eMMC device: $EMMC_DEV"

# Check disk sizes
SD_SIZE=$(blockdev --getsize64 "$SD_DEV")
EMMC_SIZE=$(blockdev --getsize64 "$EMMC_DEV")
echo "SD card size: $(numfmt --to=iec $SD_SIZE)"
echo "eMMC size: $(numfmt --to=iec $EMMC_SIZE)"

if [ "$EMMC_SIZE" -lt "$SD_SIZE" ]; then
    echo "WARNING: eMMC is smaller than SD card"
    echo "Continuing anyway..."
fi

echo "Copying partition table..."
sudo sfdisk --dump "$SD_DEV" | sudo sfdisk "$EMMC_DEV"

echo "Copying partitions..."
for part_num in 1 2 3 4; do
    SD_PART="${SD_DEV}p${part_num}"
    EMMC_PART="${EMMC_DEV}p${part_num}"

    if [ -b "$SD_PART" ] && [ -b "$EMMC_PART" ]; then
        echo "Copying partition $part_num..."
        sudo dd if="$SD_PART" of="$EMMC_PART" bs=4M status=progress conv=fsync
    else
        echo "Partition $part_num not found, skipping"
    fi
done

echo "========================================"
echo "  Flash Complete!"
echo "========================================"
echo ""
echo "Next steps:"
echo "1. Power off the board"
echo "2. Remove SD card"
echo "3. Power on - system will boot from eMMC"
echo ""
echo "To disable this script from running again:"
echo "  systemctl disable ukama-autoflash"

# Mark as complete
touch /var/lib/ukama-autoflash-complete
AUTOFLASH_EOF

    sudo chmod +x "$mnt_point/ukama-auto-flash.sh"

    # Create systemd service to run on first boot
    log "Creating systemd service..."
    sudo tee "$mnt_point/ukama-autoflash.service" > /dev/null << 'SERVICE_EOF'
[Unit]
Description=Ukama Auto-Flash to eMMC (First Boot Only)
After=multi-user.target
ConditionPathExists=!/var/lib/ukama-autoflash-complete

[Service]
Type=oneshot
ExecStart=/boot/ukama-auto-flash.sh
RemainAfterExit=yes
StandardOutput=journal+console
StandardError=journal+console

[Install]
WantedBy=multi-user.target
SERVICE_EOF

    # Enable the service (create symlink in /etc/systemd/system)
    if [ -d "$mnt_point/etc/systemd/system/multi-user.target.wants" ]; then
        sudo ln -sf /boot/ukama-autoflash.service "$mnt_point/etc/systemd/system/multi-user.target.wants/"
    fi

    # Create README
    sudo tee "$mnt_point/README-ukama.txt" > /dev/null << EOF
Ukama ${BOARD_NAME} SD Card
===========================

This SD card contains the OS image and will auto-flash to eMMC on first boot.

What happens on first boot:
1. System boots from SD card
2. Auto-flash script copies SD contents to eMMC
3. Flash completes, ready for SD card removal

Manual Steps:
1. Insert SD card into ${BOARD_NAME}
2. Power on
3. Wait for "Flash Complete" message
4. Power off
5. Remove SD card
6. Power on - system boots from eMMC

Image: $(basename "$IMAGE_PATH")
Created: $(date)
EOF

    sync
    sudo umount "$mnt_point"
    rm -rf "$mnt_point"

    log "Auto-flash script added to boot partition"
}

create_flash_only_sd() {
    log "Creating flash-only SD card..."
    log "This SD card will boot and flash eMMC, then halt"

    # This is a simpler approach - the SD card just boots
    # and the user manually runs the flash script
    log "Raw image flashed. SD card ready to boot."
    log ""
    log "To flash manually:"
    log "  1. Boot from SD card"
    log "  2. Login as root"
    log "  3. Run: /boot/ukama-auto-flash.sh"
    log ""
}

# Main
log "========================================"
log "Ukama SD Card Creator for $BOARD_NAME"
log "========================================"
log ""

verify_device
flash_raw_image

# Try to add auto-flash script, but don't fail if it doesn't work
if add_auto_flash_script; then
    log "Auto-flash script added successfully"
else
    log "WARNING: Could not add auto-flash script"
    log "You will need to flash manually after booting"
fi

log ""
log "========================================"
log "SD card created successfully!"
log "========================================"
log ""
log "Next steps:"
log "1. Remove SD card from Ubuntu"
log "2. Insert into $BOARD_NAME board"
log "3. Power on"
log "4. System will boot from SD card"
log "5. If auto-flash enabled: wait for completion"
log "6. Power off, remove SD, power on"
log ""
log "To monitor: screen /dev/ttyUSB0 115200"

exit 0
