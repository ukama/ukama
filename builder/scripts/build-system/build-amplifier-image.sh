#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

set -euo pipefail

# all tools available?
for cmd in losetup parted kpartx mkfs.vfat mkfs.ext4 mkswap truncate make bc flex bison arm-linux-gnueabihf-gcc; do
    command -v $cmd >/dev/null || { echo "ERROR: '$cmd' not found in PATH"; exit 1; }
done

# Load helper utilities
PKG_UTILS="$(dirname "$0")/pkg-utils.sh"
if [[ ! -f "$PKG_UTILS" ]]; then
    echo "ERROR: Missing $PKG_UTILS"
    exit 1
fi
source "$PKG_UTILS"

# Mount points and image path
BOOT_MNT="/media/boot"
PRIMARY_MNT="/media/primary"
PASSIVE_MNT="/media/passive"
RAW_IMG="ukama-amplifier-image.img"

DIR="$(pwd)"
UKAMA_OS=$(realpath ../../../nodes/ukamaOS)
UKAMA_ROOT=$(realpath ../../../)
UKAMA_REPO_APP_PKG="${UKAMA_ROOT}/build/pkgs"
UKAMA_REPO_LIB_PKG="${UKAMA_ROOT}/build/libs"
COMMON_CONFIG_FILE="${UKAMA_ROOT}/builder/boards/common.config"
AMPLIFIER_CONFIG_FILE="${UKAMA_ROOT}/builder/boards/controller.config"
BOOT1_BIN="${UKAMA_OS}/firmware/build/boot/at91bootstrap/at91bootstrap.bin"
BOOT2_BIN="${UKAMA_OS}/firmware/build/boot/uboot/u-boot.bin"
FIRMWARE_DIR="${UKAMA_OS}/firmware"
KERNEL_SRC="${UKAMA_OS}/kernel/linux"
KERNEL_OUT=""
CROSS_COMPILE="arm-linux-gnueabihf-"
AT91_CONFIG="sama5d27_ukama_emmc_uboot_defconfig"
UBOOT_CONFIG="sama5d27_ukama_anode_emmc_defconfig"

ROOTFS_DIR=${UKAMA_ROOT}/builder/scripts/build-system/rootfs
MISC_FILES_DIR=${UKAMA_ROOT}/builder/scripts/build-system/amplifier
APP_NAMES=()

# Alpine parameters
ALPINE_URL="http://dl-cdn.alpinelinux.org/alpine"
ALPINE_VERSION="3.21"
ALPINE_ARCH="x86_64"

# Loop device handle
LOOPDEV=""

# Logging and status helpers
test -n "\${LOG_LEVEL-}" || LOG_LEVEL=INFO
log() {
    local type="$1" msg="$2"
    local color
    case "$type" in
        INFO)    color="\033[1;34m";;
        SUCCESS) color="\033[1;32m";;
        ERROR)   color="\033[1;31m";;
        *)       color="\033[1;37m";;
    esac
    echo -e "${color}${type}: ${msg}\033[0m"
}

check_status() {
    local code=$1 msg=$2 stage=$3
    if [[ $code -ne 0 ]]; then
        log ERROR "Stage '${stage}' failed"
        exit 1
    fi
    log SUCCESS "$msg"
}

check_pre_requisite() {
    if [ -d "${ROOTFS_DIR}" ] && [ "$(ls -A ${ROOTFS_DIR})" ]; then
        log "INFO" "ROOTFS exist."
    else
        log "INFO"  "Make sure you have ran build-env-setup and rootfs-env-setup.sh scripts"
        log "ERROR" "${ROOTFS_DIR} does not exist"
        exit 1
    fi

    if [ ! -f "${FIRMWARE_DIR}/Makefile" ] || [ ! -f "${KERNEL_SRC}/arch/arm/configs/ukama_anode_defconfig" ]; then
        log "ERROR" "firmware or kernel source missing, run: git submodule update --init nodes/ukamaOS/firmware nodes/ukamaOS/kernel"
        exit 1
    fi
}

build_emmc_firmware() {
    STAGE="build_emmc_firmware"
    log INFO "Building at91bootstrap (${AT91_CONFIG}) and u-boot (${UBOOT_CONFIG})"

    make -C "${FIRMWARE_DIR}" distclean TARGET=amplifier ROOTFSPATH="${FIRMWARE_DIR}/build" > /dev/null 2>&1 || true
    make -C "${FIRMWARE_DIR}" TARGET=amplifier ROOTFSPATH="${FIRMWARE_DIR}/build" \
         AT91CONFIG="${AT91_CONFIG}" UBOOTCONFIG="${UBOOT_CONFIG}"

    for f in "${BOOT1_BIN}" "${BOOT2_BIN}" "${FIRMWARE_DIR}/u-boot/tools/mkenvimage"; do
        [[ -f "$f" ]] || { log ERROR "Firmware build did not produce ${f}"; exit 1; }
    done
    log SUCCESS "eMMC boot firmware built"
}

build_kernel() {
    STAGE="build_kernel"
    KERNEL_OUT=$(mktemp -d)
    local kmake=(make -C "${KERNEL_SRC}" O="${KERNEL_OUT}" ARCH=arm CROSS_COMPILE="${CROSS_COMPILE}")

    log INFO "Building the Ukama amplifier kernel in ${KERNEL_OUT}"
    "${kmake[@]}" ukama_anode_defconfig
    "${kmake[@]}" -j"$(nproc)" zImage modules dtbs
    "${kmake[@]}" INSTALL_MOD_PATH="${KERNEL_OUT}/modroot" modules_install

    log SUCCESS "Kernel $(ls "${KERNEL_OUT}/modroot/lib/modules") built"
}

# Cleanup function ensures mounts and loop are detached
cleanup() {
    log INFO "Cleaning up mounts and loop device"
    sudo umount "$BOOT_MNT"     2>/dev/null || true
    sudo umount "$PASSIVE_MNT"  2>/dev/null || true
    sudo umount "$PRIMARY_MNT"  2>/dev/null || true
    if [[ -n "$LOOPDEV" ]]; then
        sudo kpartx -dv "$LOOPDEV" 2>/dev/null || true
        sudo losetup -d "$LOOPDEV"  2>/dev/null || true
    fi
    if [[ -n "$KERNEL_OUT" ]]; then
        sudo rm -rf "$KERNEL_OUT"
    fi
}
trap cleanup EXIT

create_disk_image() {
    log INFO "Creating sparse image ${RAW_IMG} (16 GiB)"
    rm -f "$RAW_IMG"
    truncate -s 16G "$RAW_IMG"
    # smaller image for testing
    #    truncate -s 5G "$RAW_IMG"
    check_status $? "Raw image ready" "create_disk_image"
}

attach_loop() {
    log INFO "Attaching ${RAW_IMG} to loop device"
    LOOPDEV=$(sudo losetup --show --find --partscan "$RAW_IMG")
    if [[ -z "$LOOPDEV" ]]; then
        log ERROR "Failed to attach loop device"
        exit 1
    fi
    check_status 0 "Loop device: $LOOPDEV" "attach_loop"
}

partition_image() {
    log INFO "Partitioning $LOOPDEV: boot, passive, primary, swap"
    sudo parted -s "$LOOPDEV" mklabel msdos \
         mkpart primary fat32      1MiB    513MiB   \
         mkpart primary ext4       513MiB  4609MiB  \
         mkpart primary ext4       4609MiB 12289MiB \
         mkpart primary linux-swap 12289MiB 100%     \
         set 1 boot on
    check_status $? "Partitions created" "partition_image"
}

# smaller image for testing (5GB)
#partition_image() {
#    log INFO "Partitioning $LOOPDEV: boot (512 MiB), primary (2 GiB), passive (2 GiB), swap (rest)"
#    sudo parted -s "$LOOPDEV" mklabel msdos \
#         mkpart primary fat32      1MiB     513MiB   \
#         mkpart primary ext4       513MiB   2561MiB  \
#         mkpart primary ext4       2561MiB  4609MiB  \
#         mkpart primary linux-swap 4609MiB  100%     \
#         set 1 boot on
#    check_status $? "Partitions created" partition_image
#}

map_partitions() {
    log INFO "Mapping partitions for $LOOPDEV"
    sudo kpartx -av "$LOOPDEV" | tee /dev/stderr
    sleep 1
    check_status $? "Partitions mapped" "map_partitions"
}

format_partitions() {
    log INFO "Formatting partitions"
    STAGE="format_partitions"
    local b=$(basename "$LOOPDEV")
    local p1="/dev/mapper/${b}p1" p2="/dev/mapper/${b}p2"
    local p3="/dev/mapper/${b}p3" p4="/dev/mapper/${b}p4"

    sudo mkfs.vfat -F32 -n boot    "$p1"
    check_status $? "boot formatted" "${STAGE}"

    sudo mkfs.ext4 -L passive \
         -O ^64bit,^metadata_csum \
         "$p2"
    check_status $? "passive formatted" "${STAGE}"

    sudo mkfs.ext4 -L primary \
         -O ^64bit,^metadata_csum \
         "$p3"
    check_status $? "primary formatted" "${STAGE}"

    sudo mkswap    -L swap         "$p4"
    check_status $? "swap created" "${STAGE}"
}

mount_partition() {
    local part=$1 mp=$2
    log INFO "Mounting $part → $mp"
    sudo mkdir -p "$mp"
    sudo mount "$part" "$mp"
    check_status $? "Mounted $part" "mount_partition"
}

copy_boot() {
    STAGE="copy_boot"
    log INFO "Copying eMMC firmware, kernel and uboot.env → ${BOOT_MNT} (which will become /boot)"

    local dst="${BOOT_MNT}"
    local env_txt

    # 1) Clean out anything already on the boot partition
    sudo rm -rf "${dst:?}/"*
    check_status $? "Cleared old files on boot partition" "${STAGE}"

    sudo cp "${BOOT1_BIN}" "${dst}/BOOT.BIN"
    sudo cp "${BOOT2_BIN}" "${dst}/u-boot.bin"
    sudo cp "${KERNEL_OUT}/arch/arm/boot/zImage" "${dst}/zImage"
    sudo cp "${KERNEL_OUT}/arch/arm/boot/dts/microchip/ukama_anode.dtb" "${dst}/ukama_anode.dtb"

    env_txt=$(mktemp)
    cat > "$env_txt" <<EOF
baudrate=115200
bootdelay=3
bootargs=console=ttyS0,115200 earlyprintk root=/dev/mmcblk0p3 rootfstype=ext4 rw rootwait
bootcmd=fatload mmc 0 0x22000000 ukama_anode.dtb; fatload mmc 0 0x23000000 zImage; bootz 0x23000000 - 0x22000000
EOF
    sudo "${FIRMWARE_DIR}/u-boot/tools/mkenvimage" -s 0x4000 -o "${dst}/uboot.env" "$env_txt"
    rm -f "$env_txt"
    check_status $? "Copied BOOT.BIN, u-boot.bin, uboot.env, zImage and ukama_anode.dtb" "${STAGE}"

    log SUCCESS "Boot partition populated"
}

copy_rootfs() {
    STAGE="copy_rootfs"
    log INFO "Copying rootfs → primary (${PRIMARY_MNT}) and passive (${PASSIVE_MNT}) partitions"

    # Common rsync exclude list
    local excludes=(
        --exclude=/dev/*
        --exclude=/proc/*
        --exclude=/sys/* 
        --exclude=/run/*
        --exclude=/tmp/*
        --exclude=/boot
        --exclude=/lib/modules
        --exclude=/efi
        --exclude=/ukamarepo
        --exclude=/passive
        --exclude=/recovery
        --exclude=/data
        --exclude=/destroy
        --exclude=/enter-chroot
        --exclude=env.sh
        --exclude=setup.log
    )

    # ---- Primary partition ----
    # 1) Clear old data
    sudo rm -rf "${PRIMARY_MNT:?}/"*
    check_status $? "Cleared ${PRIMARY_MNT}" "${STAGE}"

    # 2) Copy rootfs into PRIMARY
    sudo rsync -aAX --delete "${excludes[@]}" \
         "${ROOTFS_DIR}/" "${PRIMARY_MNT}/"
    check_status $? "Rootfs copied to ${PRIMARY_MNT}" "${STAGE}"

    # 3) Recreate mount-point dirs
    sudo mkdir -p \
         "${PRIMARY_MNT}/boot" \
         "${PRIMARY_MNT}/dev" \
         "${PRIMARY_MNT}/proc" \
         "${PRIMARY_MNT}/sys" \
         "${PRIMARY_MNT}/run" \
         "${PRIMARY_MNT}/tmp"

    # ---- Passive partition ----
    sudo rm -rf "${PASSIVE_MNT:?}/"*
    check_status $? "Cleared ${PASSIVE_MNT}" "${STAGE}"

    sudo rsync -aAX --delete "${excludes[@]}" \
         "${ROOTFS_DIR}/" "${PASSIVE_MNT}/"
    check_status $? "Rootfs copied to ${PASSIVE_MNT}" "${STAGE}"

    sudo mkdir -p \
         "${PASSIVE_MNT}/boot" \
         "${PASSIVE_MNT}/dev" \
         "${PASSIVE_MNT}/proc" \
         "${PASSIVE_MNT}/sys" \
         "${PASSIVE_MNT}/run" \
         "${PASSIVE_MNT}/tmp"

    sync
    log SUCCESS "Rootfs deployed to primary & passive" "${STAGE}"
}

install_kernel_modules() {
    local rootfs="$1"
    sudo mkdir -p "${rootfs}/lib/modules"
    sudo cp -a "${KERNEL_OUT}/modroot/lib/modules/." "${rootfs}/lib/modules/"
    check_status $? "Kernel modules installed in ${rootfs}" "install_kernel_modules"
}

update_fstab() {
    local rootfs_dir="$1"
    local fstab="${rootfs_dir}/etc/fstab"
    STAGE="update_fstab"

    log INFO "Writing label-based fstab in ${fstab}"

    sudo tee "${fstab}" > /dev/null <<'FSTAB'
# pseudo-file systems
devtmpfs        /dev        devtmpfs    defaults    0 0
proc            /proc       proc        defaults    0 0
sysfs           /sys        sysfs       defaults    0 0
devpts          /dev/pts    devpts      defaults    0 0
tmpfs           /tmp        tmpfs       defaults    0 0

# real partitions by LABEL
LABEL=primary   /           ext4      defaults         0 1
LABEL=boot      /boot       auto      rw,defaults      0 2
LABEL=swap      none        swap      sw               0 0
FSTAB

    check_status $? "fstab written" "${STAGE}"
}

deploy_to_rootfs() {
    local rootfs="$1"

    log INFO "Deploying apps/libs + manifest into ${rootfs}"

    # ensure dirs exist
    sudo mkdir -p "${rootfs}/ukama/apps/pkgs" \
                     "${rootfs}/ukama/logs" \
                     "${rootfs}/ukama/state/starterd/log-spool" \
                     "${rootfs}/lib"

    # copy apps
    copy_all_apps      "$UKAMA_REPO_APP_PKG" "${rootfs}/ukama/apps/pkgs"
    check_status $? "copy_all_apps to ${rootfs}" "deploy_to_rootfs"

    # copy libs
    copy_required_libs "$UKAMA_REPO_LIB_PKG" "${rootfs}/lib"
    check_status $? "copy_required_libs to ${rootfs}" "deploy_to_rootfs"

    # create manifest
    create_manifest_file "${rootfs}/ukama/manifest.json" "${APPS[@]}"
    check_status $? "create_manifest_file in ${rootfs}" "deploy_to_rootfs"
}

copy_misc_files() {
	DEST_DIR=$1
	# MISC FILES Copy
	if [ -d "${MISC_FILES_DIR}" ] && [ "$(ls -A ${MISC_FILES_DIR})" ]; then
    	log "INFO" "${MISC_FILES_DIR} exist."
		rsync -aAXv ${MISC_FILES_DIR}/* ${DEST_DIR}/
	else
    	log "Nothing to copy" "${MISC_FILES_DIR} does not exist"
	fi
	
}

# Main
check_pre_requisite
build_emmc_firmware
build_kernel

mkdir -p "$BOOT_MNT" "$PASSIVE_MNT" "$PRIMARY_MNT"

create_disk_image
attach_loop
partition_image
map_partitions
format_partitions
mount_partition "/dev/mapper/$(basename ${LOOPDEV})p1" "$BOOT_MNT"
mount_partition "/dev/mapper/$(basename ${LOOPDEV})p2" "$PASSIVE_MNT"
mount_partition "/dev/mapper/$(basename ${LOOPDEV})p3" "$PRIMARY_MNT"

# Signal success; cleanup trap will unmount and detach
log SUCCESS "Disk image ${RAW_IMG} prepared successfully"

copy_boot
copy_rootfs
install_kernel_modules "${PRIMARY_MNT}"
install_kernel_modules "${PASSIVE_MNT}"
update_fstab "${PRIMARY_MNT}"
update_fstab "${PASSIVE_MNT}"

# Gather the list of enabled apps
get_enabled_apps "$COMMON_CONFIG_FILE" "$AMPLIFIER_CONFIG_FILE"
if (( ${#APPS[@]} > 0 )); then
    log INFO "Enabled apps: ${APPS[*]}"
else
    log ERROR "No apps enabled. Aborting."
    exit 1
fi

deploy_to_rootfs "${PRIMARY_MNT}"
deploy_to_rootfs "${PASSIVE_MNT}"
copy_misc_files  "${PRIMARY_MNT}"
copy_misc_files  "${PASSIVE_MNT}"

log "SUCCESS" "Disk image creation completed successfully!"
exit 0
