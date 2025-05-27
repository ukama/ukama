#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

# Build image for Ukama access node.

set -e

PKG_UTILS="$(dirname "$0")/pkg-utils.sh"
if [ ! -f "$PKG_UTILS" ]; then
    echo "ERROR: Missing $PKG_UTILS"
    exit 1
fi
source "$PKG_UTILS"

STAGE="init"
DIR="$(pwd)"
UKAMA_OS=$(realpath ../../../nodes/ukamaOS)
UKAMA_ROOT=$(realpath ../../../)
UKAMA_REPO_APP_PKG="${UKAMA_ROOT}/build/pkgs"
UKAMA_REPO_LIB_PKG="${UKAMA_ROOT}/build/libs"
COMMON_CONFIG_FILE="${UKAMA_ROOT}/builder/boards/common.config"
ACCESS_CONFIG_FILE="${UKAMA_ROOT}/builder/boards/access.config"
MANIFEST_FILE="manifest.json"

BOOT_MOUNT="/media/boot"
RECOVERY_MOUNT="/media/recovery"
PRIMARY_MOUNT="/media/primary"
PASSIVE_MOUNT="/media/passive"
DATA_MOUNT="/media/data"

RAW_IMG="ukama-access-node.img"

ROOTFS_DIR=${UKAMA_ROOT}/builder/scripts/build-system/rootfs/
UKAMA_APP_PKG="${ROOTFS_DIR}/ukama/apps/pkgs"
APP_NAMES=()
FIRMWARE_REPO="https://github.com/raspberrypi/firmware"
FIRMWARE_ZIP_URL="${FIRMWARE_REPO}/archive/refs/heads/master.zip"

trap cleanup EXIT

log() {
    local type="$1"
    local message="$2"
    local color
    case "$type" in
        "INFO")    color="\033[1;34m";;
        "SUCCESS") color="\033[1;32m";;
        "ERROR")   color="\033[1;31m";;
        *)         color="\033[1;37m";;
    esac
    echo -e "${color}${type}: ${message}\033[0m"
}

check_status() {
    if [ $1 -ne 0 ]; then
        log "ERROR" "Script failed at stage: $3"
        exit 1
    fi
    log "SUCCESS" "$2"
}

cleanup() {
    if [ -z "$LOOPDISK" ]; then
        return
    fi
    log "INFO" "Cleaning up resources..."
    local mounts=(${BOOT_MOUNT} ${PRIMARY_MOUNT} ${PASSIVE_MOUNT} ${UNUSED_MOUNT})
    for mount in "${mounts[@]}"; do
        sudo umount "${mount}" 2>/dev/null || true
    done
    sudo kpartx -dv "${RAW_IMG}" 2>/dev/null || true
    sudo losetup -d "${LOOPDISK}" 2>/dev/null || true
    log "INFO" "Cleanup completed."
}

create_disk_image() {
    STAGE="create_disk_image"
    log "INFO" "Creating a new raw image: ${RAW_IMG}"
    rm -f "${RAW_IMG}"
    dd if=/dev/zero of="${RAW_IMG}" bs=512 count=0 seek=61120512
    check_status $? "Raw image created" ${STAGE}
}

download_and_copy_rpi_firmware() {
    STAGE="download_and_copy_rpi_firmware"
    BOOT_DIR="$1"
    TMP_DIR="/tmp/rpi-bootloader"

    log "INFO" "Downloading Raspberry Pi bootloader from $FIRMWARE_REPO ..."

    rm -rf "$TMP_DIR"
    mkdir -p "$TMP_DIR"

    wget -qO "$TMP_DIR/master.zip" "$FIRMWARE_ZIP_URL"
    unzip -q "$TMP_DIR/master.zip" -d "$TMP_DIR"

    FIRMWARE_BOOT="$TMP_DIR/firmware-master/boot"
    if [ ! -d "$FIRMWARE_BOOT" ]; then
        log "ERROR" "Failed to extract firmware boot folder"
        return 1
    fi

    mkdir -p "$BOOT_DIR/firmware"
    cp "$FIRMWARE_BOOT"/{start*.elf,fixup*.dat,bootcode.bin} "$BOOT_DIR/firmware/" 2>/dev/null || true
    cp "$FIRMWARE_BOOT"/*.dtb "$BOOT_DIR/firmware/" 2>/dev/null || true
    cp -r "$FIRMWARE_BOOT/overlays" "$BOOT_DIR/firmware/" 2>/dev/null || true

    # Copy actual kernel image
    if [ -f "${ROOTFS_DIR}/boot/kernel.img" ]; then
        cp "${ROOTFS_DIR}/boot/kernel.img" "$BOOT_DIR/firmware/"
        log "INFO" "Copied kernel.img from rootfs"
    elif [ -f "${ROOTFS_DIR}/boot/vmlinuz-rpi" ]; then
        cp "${ROOTFS_DIR}/boot/vmlinuz-rpi" "$BOOT_DIR/firmware/kernel.img"
        log "INFO" "Renamed vmlinuz-rpi to kernel.img and copied to firmware/"
    elif [ -f "${ROOTFS_DIR}/boot/Image" ]; then
        cp "${ROOTFS_DIR}/boot/Image" "$BOOT_DIR/firmware/kernel.img"
        log "INFO" "Renamed Image to kernel.img and copied to firmware/"
    else
        log "ERROR" "No kernel image found in ${ROOTFS_DIR}/boot"
        exit 1
    fi

    # Write a clean config.txt
    cat <<EOF > "$BOOT_DIR/firmware/config.txt"
enable_uart=1
arm_64bit=1
kernel=kernel.img
gpu_mem=64
boot_delay=1
disable_splash=1
dtoverlay=disable-bt
dtoverlay=disable-wifi
EOF

    # Write cmdline.txt
    cat <<EOF > "$BOOT_DIR/firmware/cmdline.txt"
console=serial0,115200 console=tty1 root=LABEL=primary rootfstype=ext4 fsck.repair=yes rootwait
EOF

    log "SUCCESS" "Bootloader files and boot config copied to $BOOT_DIR/firmware"

    rm -rf "$TMP_DIR"
}

setup_loop_device() {
    STAGE="setup_loop_device"
    log "INFO" "Attaching ${RAW_IMG} to a loop device"
    LOOPDISK=$(sudo losetup -f --show "${RAW_IMG}")
    if [ -z "${LOOPDISK}" ]; then
        log "ERROR" "Failed to set up loop device for ${RAW_IMG}."
        exit 1
    fi
    log "SUCCESS" "Loop device set up at ${LOOPDISK}"
}

clean_first_50MB() {
    STAGE="clean_first_50MB"
    log "INFO" "Cleaning the first 50MB of ${LOOPDISK}"
    sudo dd if=/dev/zero of="${LOOPDISK}" bs=1M count=50
    check_status $? "First 50MB cleaned" ${STAGE}
}

partition_image() {
    STAGE="partition_image"
    log "INFO" "Creating partitions on ${LOOPDISK} using sfdisk"
    sudo sfdisk "${LOOPDISK}" <<-__EOF__
label: dos
,1G,83
,4G,83
,,5
,6G,83
,6G,83
,11G,83
,1G,82
;
__EOF__
    check_status $? "Partitions created" ${STAGE}
}

map_partitions() {
    STAGE="map_partitions"
    log "INFO" "Mapping partitions using kpartx"
    sudo kpartx -v -a "${LOOPDISK}"
    check_status $? "Partitions mapped" ${STAGE}

    DEVICE=$(basename "${LOOPDISK}")
    DISK="/dev/mapper/${DEVICE}p"
}

format_partitions() {
    STAGE="format_partitions"
    log "INFO" "Formatting partitions"
    DEVICE=$(basename "${LOOPDISK}")
    DISK="/dev/mapper/${DEVICE}p"

    mkfs.vfat -F 32 -n boot "${DISK}1"
    check_status $? "boot partition formatted" ${STAGE}

    mkfs.ext4 -L recovery "${DISK}2"
    check_status $? "recovery rootfs formatted" ${STAGE}

    mkfs.ext4 -L primary "${DISK}5"
    check_status $? "primary rootfs formatted" ${STAGE}

    mkfs.ext4 -L passive "${DISK}6"
    check_status $? "passive rootfs formatted" ${STAGE}

    mkfs.ext4 -L data "${DISK}7"
    check_status $? "data partition formatted" ${STAGE}

    mkswap "${DISK}8"
    check_status $? "swap partition created" ${STAGE}
}

mount_partition() {
    local partition=$1
    local mount_point=$2
    log "INFO" "Mounting ${partition} to ${mount_point}"
    sudo mkdir -p "${mount_point}"
    sudo mount "${partition}" "${mount_point}"
    check_status $? "Partition ${partition} mounted to ${mount_point}" "${STAGE}"
}

unmount_partition() {
    local mount_point=$1
    log "INFO" "Unmounting ${mount_point}"
    sudo umount "${mount_point}"
    check_status $? "${mount_point} unmounted" "${STAGE}"
}

copy_rootfs() {
    STAGE="copy_rootfs"

    log "INFO" "Copying rootfs to primary and passive"
	rsync -aAX --exclude={"/dev","/sys","/proc"} ${ROOTFS_DIR}/* ${PRIMARY_MOUNT}/
	mkdir -p ${PRIMARY_MOUNT}/dev ${PRIMARY_MOUNT}/sys ${PRIMARY_MOUNT}/proc

	log "INFO" "Copying rootfs to primary and passive"
	rsync -aAX --exclude={"/dev","/sys","/proc"} ${ROOTFS_DIR}/* ${PASSIVE_MOUNT}/
    mkdir -p ${PASSIVE_MOUNT}/dev ${PASSIVE_MOUNT}/sys ${PASSIVE_MOUNT}/proc

	sync
}

set_permissions() {
    STAGE="set_permissions"
    log "INFO" "Setting permissions for primary and passive partitions"
    sudo chown -R root:root ${PRIMARY_MOUNT}
    sudo chmod -R 755 ${PRIMARY_MOUNT}
    check_status $? "Permissions set for primary" ${STAGE}

    sudo chown -R root:root ${PASSIVE_MOUNT}
    sudo chmod -R 755 ${PASSIVE_MOUNT}
    check_status $? "Permissions set for passive" ${STAGE}
}

# Update /etc/fstab based on partition type
update_fstab() {
    PARTITION_TYPE=$1
    log "INFO" "Updating /etc/fstab for partition type: ${PARTITION_TYPE}"

    if [[ "${PARTITION_TYPE}" == "${PRIMARY_MOUNT}" ]]; then
        cat <<FSTAB > ${PRIMARY_MOUNT}/etc/fstab
proc              /proc           proc    defaults              0 0
sysfs             /sys            sysfs   defaults              0 0
devpts            /dev/pts        devpts  defaults              0 0
tmpfs             /tmp            tmpfs   defaults              0 0
LABEL=recovery    /recovery       ext4    ro                    0 2
LABEL=data        /data           ext4    ro                    0 2
LABEL=passive     /passive        ext4    ro                    0 2
LABEL=primary     /               ext4    errors=remount-ro     0 1
LABEL=boot        /boot/firmware  vfat    ro                    0 2
FSTAB
    else
        cat <<FSTAB > ${PASSIVE_MOUNT}/etc/fstab
proc              /proc           proc    defaults              0 0
sysfs             /sys            sysfs   defaults              0 0
devpts            /dev/pts        devpts  defaults              0 0
tmpfs             /tmp            tmpfs   defaults              0 0
LABEL=recovery    /recovery       ext4    ro                    0 2
LABEL=data        /data           ext4    ro                    0 2
LABEL=primary     /passive        ext4    ro                    0 2
LABEL=passive     /               ext4    errors=remount-ro     0 1
LABEL=boot        /boot/firmware  vfat    ro                    0 2
FSTAB
    fi

    log "INFO" "${PARTITION_TYPE}/etc/fstab updated successfully."
}

detach_loop_device() {
    STAGE="detach_loop_device"
    log "INFO" "Detaching loop device and cleaning up"
    sudo kpartx -dv "${LOOPDISK}"
    sudo losetup -d "${LOOPDISK}"
    check_status $? "Loop device detached" ${STAGE}
}

# Main
if [ -d "${ROOTFS_DIR}" ] && [ "$(ls -A ${ROOTFS_DIR})" ]; then
    log "INFO" "ROOTFS exist."
else
    log "ERROR" "${ROOTFS_DIR} does not exist"
    log "INFO" "Make sure you have run build-env-setup and rootfs-env-setup.sh scripts"
    exit 1
fi

create_disk_image
setup_loop_device
clean_first_50MB
partition_image
map_partitions
format_partitions
mount_partition "${DISK}1" "${BOOT_MOUNT}"
mount_partition "${DISK}5" "${PRIMARY_MOUNT}"
mount_partition "${DISK}6" "${PASSIVE_MOUNT}"
download_and_copy_rpi_firmware "${BOOT_MOUNT}"

# create board specific manifest and cp its pkds/libs
get_enabled_apps "$COMMON_CONFIG_FILE" "$ACCESS_CONFIG_FILE"
if [[ ${#APPS[@]} -gt 0 ]]; then
    log "INFO" "APPS are: ${APPS[@]}"
else
    log "ERROR" "APPS not assigned. Exit"
    unmount_partition "${BOOT_MOUNT}"
    unmount_partition "${PRIMARY_MOUNT}"
    unmount_partition "${PASSIVE_MOUNT}"
    detach_loop_device

    log "ERROR" "Disk image creation unsuccessful"
    exit 1
fi
copy_all_apps "$UKAMA_REPO_APP_PKG" "$UKAMA_APP_PKG"
copy_required_libs "$UKAMA_REPO_LIB_PKG" "$ROOTFS_DIR"
create_manifest_file "$MANIFEST_FILE" "${APPS[@]}"

copy_rootfs
cp -rf "${MANIFEST_FILE}" "${PRIMARY_MOUNT}"
cp -rf "${MANIFEST_FILE}" "${PASSIVE_MOUNT}"
rm -rf "${MANIFEST_FILE}"
set_permissions
update_fstab "${PRIMARY_MOUNT}"
update_fstab "${PASSIVE_MOUNT}"

unmount_partition "${BOOT_MOUNT}"
unmount_partition "${PRIMARY_MOUNT}"
unmount_partition "${PASSIVE_MOUNT}"
detach_loop_device

log "SUCCESS" "Disk image creation completed successfully!"
