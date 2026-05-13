#!/bin/sh
set -eux

echo "[@@BOARD_NAME@@] Enabling eth0"
ip link set eth0 up

echo "[@@BOARD_NAME@@] Bringing up eth0 via DHCP (udhcpc)"
udhcpc -i eth0 -q

echo "[@@BOARD_NAME@@] Waiting a couple seconds for lease…"
sleep 2

echo "[@@BOARD_NAME@@] Detecting eMMC device..."
for dev in /dev/mmcblk*; do
    if [ -e "${dev}boot0" ] && [ -e "${dev}boot1" ]; then
        EMMC_DEV="$dev"
        break
    fi
done

if [ -z "${EMMC_DEV:-}" ]; then
    echo "[ERROR] No eMMC device found with boot0/boot1"
    exit 1
fi
echo "[@@BOARD_NAME@@] Detected eMMC device: $EMMC_DEV"

echo "[@@BOARD_NAME@@] Downloading image from @@HOST_IP@@:@@HTTP_PORT@@/@@IMG_NAME@@"
wget "http://@@HOST_IP@@:@@HTTP_PORT@@/@@IMG_NAME@@" -O "/mnt/@@IMG_NAME@@"

ls -lh "/mnt/@@IMG_NAME@@"

if [ ! -f "/mnt/@@IMG_NAME@@" ]; then
    echo "[@@BOARD_NAME@@] Image not found after wget!"
    exit 1
fi

echo "[@@BOARD_NAME@@] Zeroing first 64MB of $EMMC_DEV"
dd if=/dev/zero of="$EMMC_DEV" bs=1M count=64

echo "[@@BOARD_NAME@@] Flashing image to $EMMC_DEV"
dd if="/mnt/@@IMG_NAME@@" of="$EMMC_DEV" bs=4M
sync

echo "[@@BOARD_NAME@@] Flash complete. Rebooting."
reboot
