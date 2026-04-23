#!/bin/bash
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

sudo umount ${SD_DEVICE}* 2>/dev/null || true

echo "Writing image to SD card..."
sudo dd if="$IMAGE_FILE" of="$SD_DEVICE" bs=4M status=progress conv=fsync
sync

sleep 2
sudo partprobe "$SD_DEVICE" 2>/dev/null || true
sleep 2

echo "Checking for partitions..."
if ! lsblk "$SD_DEVICE" | grep -q part; then
    echo "ERROR: No partitions found. Image may be corrupted."
    exit 1
fi

if [[ "$SD_DEVICE" =~ mmcblk ]] || [[ "$SD_DEVICE" =~ loop ]]; then
    ROOT_PART="${SD_DEVICE}p2"
else
    ROOT_PART="${SD_DEVICE}2"
fi

if [ ! -b "$ROOT_PART" ]; then
    echo "ERROR: Root partition $ROOT_PART not found"
    lsblk "$SD_DEVICE"
    exit 1
fi

MOUNT_POINT="/tmp/sdcard_root"
sudo mkdir -p "$MOUNT_POINT"
sudo mount "$ROOT_PART" "$MOUNT_POINT"

sudo tee "$MOUNT_POINT/usr/local/bin/auto-flash-emmc.sh" > /dev/null <<'EOF'
#!/bin/bash
set -e

echo "Auto-flashing SD to eMMC..."

SD_DEV=$(lsblk -ndo NAME,TYPE | grep disk | grep -v mmcblk | head -n1 | awk '{print $1}')
SD_DEV="/dev/${SD_DEV}"
EMMC_DEV="TARGET_EMMC_PLACEHOLDER"

echo "Source: $SD_DEV -> Target: $EMMC_DEV"
sleep 5

dd if="$SD_DEV" of="$EMMC_DEV" bs=4M status=progress conv=fsync
sync

echo "Flash complete. Rebooting..."
sleep 3

systemctl disable auto-flash-emmc.service || true
reboot
EOF

sudo sed -i "s|TARGET_EMMC_PLACEHOLDER|$TARGET_EMMC|g" "$MOUNT_POINT/usr/local/bin/auto-flash-emmc.sh"
sudo chmod +x "$MOUNT_POINT/usr/local/bin/auto-flash-emmc.sh"

sudo tee "$MOUNT_POINT/etc/systemd/system/auto-flash-emmc.service" > /dev/null <<'EOF'
[Unit]
Description=Auto-flash SD to eMMC
After=multi-user.target

[Service]
Type=oneshot
ExecStart=/usr/local/bin/auto-flash-emmc.sh
StandardOutput=journal+console

[Install]
WantedBy=multi-user.target
EOF

sudo chroot "$MOUNT_POINT" systemctl enable auto-flash-emmc.service 2>/dev/null || {
    sudo mkdir -p "$MOUNT_POINT/etc/systemd/system/multi-user.target.wants"
    sudo ln -sf /etc/systemd/system/auto-flash-emmc.service \
        "$MOUNT_POINT/etc/systemd/system/multi-user.target.wants/auto-flash-emmc.service"
}

sudo umount "$MOUNT_POINT"
sudo rmdir "$MOUNT_POINT"

echo "SD card ready. Insert into board and power on."
