#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

: "${DEV:?Must set DEV}"
: "${BOARD_NAME:?Must set BOARD_NAME}"
: "${BOOT_TARBALL:?Must set BOOT_TARBALL}"
: "${ROOTFSA_TARBALL:?Must set ROOTFSA_TARBALL}"
: "${ROOTFSB_TARBALL:?Must set ROOTFSB_TARBALL}"

MOUNT_POINTS=()

cleanup() {
    for mnt in "${MOUNT_POINTS[@]:-}"; do
        [ -z "$mnt" ] && continue
        if mountpoint -q "$mnt" 2>/dev/null; then
            sudo umount "$mnt" 2>/dev/null || true
        fi
        rm -rf "$mnt" 2>/dev/null || true
    done
}
trap cleanup EXIT

log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $*"
}

partition_name() {
    local dev="$1"
    local num="$2"
    if [[ "$dev" =~ [0-9]$ ]]; then
        echo "${dev}p${num}"
    else
        echo "${dev}${num}"
    fi
}

verify_device() {
    if [ ! -b "$DEV" ]; then
        log "ERROR: $DEV is not a block device"
        lsblk
        exit 1
    fi
    log "WARNING: this will erase $DEV"
    lsblk "$DEV"
    read -rp "Type 'yes' to continue: " confirm
    if [[ "$confirm" != "yes" ]]; then
        log "Aborted"
        exit 1
    fi
}

verify_tarballs() {
    local validate_lib="${SCRIPT_DIR}/lib/validate.sh"
    if [ ! -f "$validate_lib" ]; then
        log "ERROR: missing helper $validate_lib"
        exit 1
    fi
    source "$validate_lib"

    log "Validating tarballs (no disk writes yet)..."
    local fail=0
    validate_tarball "$BOOT_TARBALL"    "boot"    "u-boot.bin"      3    || fail=1
    validate_tarball "$ROOTFSA_TARBALL" "rootfsA" "bin,lib,etc,usr" 1000 || fail=1
    validate_tarball "$ROOTFSB_TARBALL" "rootfsB" "bin,lib,etc,usr" 1000 || fail=1

    if [ "$fail" -ne 0 ]; then
        log "ERROR: tarball validation failed — SD card NOT modified."
        log "Fix the tarballs above and re-run. They must be created with:"
        log "  cd <mounted_partition_dir> && sudo tar czf out.tgz ."
        exit 1
    fi
    log "Tarballs OK."
}

unmount_existing_partitions() {
    local part
    for part in $(lsblk -nrpo NAME,MOUNTPOINT "$DEV" | awk 'NF>1 && $2!=""{print $1}'); do
        log "Unmounting $part"
        sudo umount "$part" 2>/dev/null || true
    done
    sudo partprobe "$DEV" 2>/dev/null || true
    sleep 1
}

wipe_and_partition() {
    unmount_existing_partitions

    log "Wiping $DEV..."
    sudo sgdisk --clear "$DEV" 2>/dev/null || true
    sudo dd if=/dev/zero of="$DEV" bs=1M count=8 2>/dev/null || true

    log "Creating partitions..."
    sudo parted -s "$DEV" mklabel msdos
    sudo parted -s "$DEV" mkpart primary fat32 1MiB 1025MiB
    sudo parted -s "$DEV" mkpart primary ext4 1025MiB 5121MiB
    sudo parted -s "$DEV" mkpart extended 5121MiB 100%
    sudo parted -s "$DEV" mkpart logical ext4 5122MiB 11265MiB
    sudo parted -s "$DEV" mkpart logical ext4 11266MiB 17409MiB
    sudo parted -s "$DEV" mkpart logical ext4 17410MiB 28877MiB
    sudo parted -s "$DEV" mkpart logical linux-swap 28878MiB 100%

    sync
    sleep 2
    sudo partprobe "$DEV" 2>/dev/null || sudo blockdev --rereadpt "$DEV" || true
    sleep 2
}

format_partitions() {
    local boot_part recovery_part primary_part passive_part data_part swap_part
    boot_part=$(partition_name "$DEV" 1)
    recovery_part=$(partition_name "$DEV" 2)
    primary_part=$(partition_name "$DEV" 5)
    passive_part=$(partition_name "$DEV" 6)
    data_part=$(partition_name "$DEV" 7)
    swap_part=$(partition_name "$DEV" 8)

    log "Formatting partitions..."
    sudo mkfs.vfat -F 32 -n boot "$boot_part" >/dev/null
    sudo mkfs.ext4 -F -L recovery "$recovery_part" >/dev/null
    sudo mkfs.ext4 -F -L primary "$primary_part" >/dev/null
    sudo mkfs.ext4 -F -L passive "$passive_part" >/dev/null
    sudo mkfs.ext4 -F -L data "$data_part" >/dev/null
    sudo mkswap -L swap "$swap_part" >/dev/null
    sync
}

extract_tarball() {
    local tarball="$1"
    local part="$2"
    local label="$3"
    local mnt

    log "Extracting $label..."
    mnt=$(mktemp -d "/tmp/ukama-${label}.XXXXXX")
    MOUNT_POINTS+=("$mnt")

    sudo mount "$part" "$mnt"
    log "Tarball $tarball: $(tar tzf "$tarball" | wc -l) entries, $(du -h "$tarball" | cut -f1) size"
    sudo tar -xzf "$tarball" -C "$mnt" --no-same-owner --no-same-permissions --warning=no-timestamp
    sync
    log "Extracted top-level: $(ls -A "$mnt" | head -5 | tr '\n' ' ')"
    sudo umount "$mnt"
    rm -rf "$mnt"
}

install_autoflash() {
    local root_part="$1"
    local label="$2"
    local mnt

    mnt=$(mktemp -d "/tmp/ukama-${label}-flash.XXXXXX")
    MOUNT_POINTS+=("$mnt")
    sudo mount "$root_part" "$mnt"

    sudo mkdir -p "$mnt/usr/local/sbin"
    sudo tee "$mnt/usr/local/sbin/ukama-auto-flash.sh" >/dev/null <<'FLASH_EOF'
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
    echo "ERROR: eMMC smaller than SD"
    exit 1
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

echo "Flash complete. Rebooting."
reboot
FLASH_EOF

    sudo chmod +x "$mnt/usr/local/sbin/ukama-auto-flash.sh"

    sudo mkdir -p "$mnt/etc/systemd/system"
    sudo tee "$mnt/etc/systemd/system/ukama-autoflash.service" >/dev/null <<'UNIT_EOF'
[Unit]
Description=Ukama Auto-Flash to eMMC
After=multi-user.target
ConditionPathExists=!/var/lib/ukama-autoflash-complete

[Service]
Type=oneshot
ExecStart=/usr/local/sbin/ukama-auto-flash.sh
RemainAfterExit=yes
StandardOutput=journal+console
StandardError=journal+console

[Install]
WantedBy=multi-user.target
UNIT_EOF

    sudo mkdir -p "$mnt/etc/systemd/system/multi-user.target.wants"
    sudo ln -sf /etc/systemd/system/ukama-autoflash.service "$mnt/etc/systemd/system/multi-user.target.wants/ukama-autoflash.service"

    sync
    sudo umount "$mnt"
    rm -rf "$mnt"
}

verify_layout() {
    local boot_part primary_part passive_part mnt
    boot_part=$(partition_name "$DEV" 1)
    primary_part=$(partition_name "$DEV" 5)
    passive_part=$(partition_name "$DEV" 6)

    log "Verifying layout..."

    for part in "$boot_part" "$primary_part" "$passive_part"; do
        if [ ! -b "$part" ]; then
            log "ERROR: missing partition $part"
            exit 1
        fi
    done

    mnt=$(mktemp -d "/tmp/ukama-verify.XXXXXX")
    MOUNT_POINTS+=("$mnt")

    sudo mount "$boot_part" "$mnt"
    if [ -z "$(ls -A "$mnt")" ]; then
        log "ERROR: boot partition is empty"
        sudo umount "$mnt"
        exit 1
    fi
    sudo umount "$mnt"

    for part in "$primary_part" "$passive_part"; do
        sudo mount "$part" "$mnt"
        if [ ! -f "$mnt/bin/sh" ] && [ ! -L "$mnt/bin/sh" ]; then
            log "ERROR: $part missing rootfs content (found: $(ls -A "$mnt" | head -5 | tr '\n' ' '))"
            sudo umount "$mnt"
            exit 1
        fi
        if [ ! -f "$mnt/usr/local/sbin/ukama-auto-flash.sh" ]; then
            log "ERROR: auto-flash missing in $part"
            sudo umount "$mnt"
            exit 1
        fi
        sudo umount "$mnt"
    done

    rm -rf "$mnt"
    log "Verification passed"
}

log "========================================"
log "Ukama Legacy SD Card Creator for $BOARD_NAME"
log "========================================"

verify_device
verify_tarballs
wipe_and_partition
format_partitions

boot_part=$(partition_name "$DEV" 1)
primary_part=$(partition_name "$DEV" 5)
passive_part=$(partition_name "$DEV" 6)

extract_tarball "$BOOT_TARBALL" "$boot_part" "boot"
extract_tarball "$ROOTFSA_TARBALL" "$primary_part" "primary"
extract_tarball "$ROOTFSB_TARBALL" "$passive_part" "passive"

install_autoflash "$primary_part" "primary"
install_autoflash "$passive_part" "passive"

verify_layout

log ""
log "========================================"
log "SD card ready"
log "========================================"
log "Next steps:"
log "1. Eject $DEV"
log "2. Insert into $BOARD_NAME"
log "3. Connect serial: screen /dev/ttyUSB0 115200"
log "4. Power on — board will flash to eMMC and reboot"

exit 0
