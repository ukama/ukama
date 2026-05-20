#!/bin/bash
set -e

LOG_FILE="/var/log/ukama-autoflash.log"
exec > >(tee -a "$LOG_FILE")
exec 2>&1

echo "========================================"
echo "  Ukama Auto-Flash to eMMC"
echo "  $(date)"
echo "========================================"

partition_path() {
    local dev="$1"
    local num="$2"
    case "$dev" in
        *[0-9]) echo "${dev}p${num}" ;;
        *)      echo "${dev}${num}" ;;
    esac
}

ROOT_SOURCE=$(findmnt -n -o SOURCE /)
SD_PARENT=$(lsblk -nro PKNAME "$ROOT_SOURCE" 2>/dev/null | head -1)
if [ -n "$SD_PARENT" ]; then
    SD_DEV="/dev/$SD_PARENT"
else
    SD_DEV="$ROOT_SOURCE"
fi

if [ ! -b "$SD_DEV" ]; then
    echo "ERROR: could not determine SD device"
    exit 1
fi

EMMC_DEV=""
for dev in /dev/mmcblk*; do
    if [ -b "$dev" ] && [ -e "${dev}boot0" ] && [ -e "${dev}boot1" ]; then
        if ! lsblk -nrpo NAME,MOUNTPOINT "$dev" | awk 'NF>1 && $2!=""{m=1} END{exit !m}'; then
            EMMC_DEV="$dev"
            break
        fi
    fi
done

if [ -z "$EMMC_DEV" ]; then
    echo "ERROR: no eMMC found"
    exit 1
fi

echo "SD: $SD_DEV"
echo "eMMC: $EMMC_DEV"

SD_SIZE=$(blockdev --getsize64 "$SD_DEV")
EMMC_SIZE=$(blockdev --getsize64 "$EMMC_DEV")
echo "SD size: $(numfmt --to=iec $SD_SIZE)"
echo "eMMC size: $(numfmt --to=iec $EMMC_SIZE)"

if [ "$EMMC_SIZE" -lt "$SD_SIZE" ]; then
    echo "WARNING: eMMC ($EMMC_SIZE bytes) smaller than SD ($SD_SIZE bytes)"
    echo "WARNING: partitions that exceed eMMC will fail; continuing"
fi

echo "Zeroing eMMC start..."
dd if=/dev/zero of="$EMMC_DEV" bs=1M count=64

echo "Cloning partition table..."
sfdisk --dump "$SD_DEV" | sfdisk "$EMMC_DEV"

echo "Waiting for eMMC partitions..."
sleep 2
i=0
while [ $i -lt 10 ]; do
    if [ -b "$(partition_path "$EMMC_DEV" 1)" ]; then
        break
    fi
    sleep 1
    i=$((i + 1))
    partprobe "$EMMC_DEV" 2>/dev/null || true
done

for num in 1 2 5 6 7 8; do
    sd_part=$(partition_path "$SD_DEV" "$num")
    emmc_part=$(partition_path "$EMMC_DEV" "$num")
    if [ -b "$sd_part" ] && [ -b "$emmc_part" ]; then
        echo "Copying partition $num..."
        dd if="$sd_part" of="$emmc_part" bs=4M status=progress conv=fsync
    else
        echo "Skipping partition $num"
    fi
done

sync

mkdir -p /var/lib
touch /var/lib/ukama-autoflash-complete

for num in 5 6; do
    emmc_part=$(partition_path "$EMMC_DEV" "$num")
    if [ -b "$emmc_part" ]; then
        e2fsck -yf "$emmc_part" >/dev/null 2>&1 || true
        mnt=$(mktemp -d)
        if mount "$emmc_part" "$mnt" 2>/dev/null; then
            mkdir -p "$mnt/var/lib"
            touch "$mnt/var/lib/ukama-autoflash-complete"
            sync
            umount "$mnt"
        fi
        rm -rf "$mnt"
    fi
done

echo "Flash complete. Rebooting."
reboot
