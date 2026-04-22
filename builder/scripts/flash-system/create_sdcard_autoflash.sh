#!/bin/bash
# Create SD card with auto-flash script for Microchip SOM controller
set -euo pipefail

SD_DEVICE="${SD_DEVICE:-}"
IMAGE_FILE="${IMAGE_FILE:-}"
TARGET_EMMC="${TARGET_EMMC:-/dev/mmcblk0}"

if [ -z "$SD_DEVICE" ] || [ -z "$IMAGE_FILE" ]; then
    echo "Usage: SD_DEVICE=/dev/sdX IMAGE_FILE=controller.img TARGET_EMMC=/dev/mmcblk0 $0"
    exit 1
fi

if [ ! -b "$SD_DEVICE" ]; then
    echo "Error: $SD_DEVICE is not a block device"
    exit 1
fi

if [ ! -f "$IMAGE_FILE" ]; then
    echo "Error: Image file $IMAGE_FILE not found"
    exit 1
fi

echo "WARNING: This will erase $SD_DEVICE completely!"
echo "Press Ctrl+C to cancel, or Enter to continue..."
read

# Unmount any mounted partitions
sudo umount ${SD_DEVICE}* 2>/dev/null || true

# Write image to SD card
echo "Writing image to SD card..."
sudo dd if="$IMAGE_FILE" of="$SD_DEVICE" bs=4M status=progress conv=fsync
sync

# Wait for partitions to appear
sleep 2
sudo partprobe "$SD_DEVICE" 2>/dev/null || true
sleep 2

# Find the root partition (usually partition 2)
if [[ "$SD_DEVICE" =~ mmcblk ]] || [[ "$SD_DEVICE" =~ loop ]]; then
    ROOT_PART="${SD_DEVICE}p2"
else
    ROOT_PART="${SD_DEVICE}2"
fi

# Mount root partition
MOUNT_POINT="/tmp/sdcard_root"
sudo mkdir -p "$MOUNT_POINT"
sudo mount "$ROOT_PART" "$MOUNT_POINT"

# Create auto-flash script
echo "Creating auto-flash script..."
sudo tee "$MOUNT_POINT/usr/local/bin/auto-flash-emmc.sh" > /dev/null <<'EOF'
#!/bin/bash
set -e

echo "========================================="
echo "Auto-flashing SD card to eMMC"
echo "========================================="

# Find SD card device (source)
SD_DEV=$(lsblk -ndo NAME,TYPE | grep disk | grep -v mmcblk | head -n1 | awk '{print $1}')
SD_DEV="/dev/${SD_DEV}"

# Target eMMC
EMMC_DEV="TARGET_EMMC_PLACEHOLDER"

echo "Source: $SD_DEV"
echo "Target: $EMMC_DEV"
echo ""
echo "Starting copy in 5 seconds... (Ctrl+C to cancel)"
sleep 5

echo "Copying $SD_DEV to $EMMC_DEV..."
dd if="$SD_DEV" of="$EMMC_DEV" bs=4M status=progress conv=fsync
sync

echo ""
echo "========================================="
echo "Flash complete! Rebooting in 5 seconds..."
echo "Remove SD card after shutdown."
echo "========================================="
sleep 5

# Disable this service so it doesn't run again
systemctl disable auto-flash-emmc.service || true

reboot
EOF

# Replace placeholder with actual target device
sudo sed -i "s|TARGET_EMMC_PLACEHOLDER|$TARGET_EMMC|g" "$MOUNT_POINT/usr/local/bin/auto-flash-emmc.sh"
sudo chmod +x "$MOUNT_POINT/usr/local/bin/auto-flash-emmc.sh"

# Create systemd service to run on boot
sudo tee "$MOUNT_POINT/etc/systemd/system/auto-flash-emmc.service" > /dev/null <<'EOF'
[Unit]
Description=Auto-flash SD card to eMMC
After=multi-user.target
ConditionPathExists=/usr/local/bin/auto-flash-emmc.sh

[Service]
Type=oneshot
ExecStart=/usr/local/bin/auto-flash-emmc.sh
StandardOutput=journal+console
StandardError=journal+console

[Install]
WantedBy=multi-user.target
EOF

# Enable the service
sudo chroot "$MOUNT_POINT" systemctl enable auto-flash-emmc.service 2>/dev/null || {
    # If chroot doesn't work, create symlink manually
    sudo mkdir -p "$MOUNT_POINT/etc/systemd/system/multi-user.target.wants"
    sudo ln -sf /etc/systemd/system/auto-flash-emmc.service \
        "$MOUNT_POINT/etc/systemd/system/multi-user.target.wants/auto-flash-emmc.service"
}

# Unmount
sudo umount "$MOUNT_POINT"
sudo rmdir "$MOUNT_POINT"

echo ""
echo "========================================="
echo "SD card is ready!"
echo "========================================="
echo "Next steps:"
echo "1. Remove SD card from your computer"
echo "2. Insert SD card into Microchip SOM controller board"
echo "3. Power on the board"
echo "4. Board will boot from SD and auto-copy to eMMC"
echo "5. Board will reboot automatically"
echo "6. Remove SD card after reboot"
echo "7. Board will boot from eMMC"
echo ""
