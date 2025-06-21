#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

set -euo pipefail

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
COM_CONFIG_FILE="${UKAMA_ROOT}/builder/boards/com.config"
MANIFEST_FILE="manifest.json"

BOOT_MOUNT="/media/boot"
PRIMARY_MOUNT="/media/primary"
PASSIVE_MOUNT="/media/passive"
#DATA_MOUNT="/media/data"

RAW_IMG="ukama-com-image.img"

ROOTFS_DIR=${UKAMA_ROOT}/builder/scripts/build-system/rootfs 
UKAMA_APP_PKG="${ROOTFS_DIR}/ukama/apps/pkgs"
APP_NAMES=()
ALPINE_URL="http://dl-cdn.alpinelinux.org/alpine"
ALPINE_VERSION="3.19"
ALPINE_ARCH="x86_64"

trap cleanup EXIT

log() {
    local type="$1"
    local message="$2"
    local color
    case "$type" in
        "INFO") color="\033[1;34m";;
        "SUCCESS") color="\033[1;32m";;
        "ERROR") color="\033[1;31m";;
        *) color="\033[1;37m";;
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
    local mounts=(${BOOT_MOUNT} ${PRIMARY_MOUNT} ${PASSIVE_MOUNT})
    for mount in "${mounts[@]}"; do
        sudo umount "${mount}" 2>/dev/null || true
    done
    sudo kpartx -dv "${RAW_IMG}" 2>/dev/null || true
    sudo losetup -d "${LOOPDISK}" 2>/dev/null || true
    log "INFO" "Cleanup completed."
}

check_label() {
    local dev="$1"
    local expected="$2"
    local actual

    local fstype
    fstype=$(blkid -o value -s TYPE "$dev")

    case "$fstype" in
        vfat)
            actual=$(sudo fatlabel "$dev" | tail -n1)
            ;;
        swap)
            actual=$(blkid -s LABEL -o value "$dev")
            ;;
        ext*)
            actual=$(sudo e2label "$dev")
            ;;
        *)
            log "ERROR" "Unsupported filesystem type: $fstype on $dev"
            exit 1
            ;;
    esac

    if [[ "$actual" != "$expected" ]]; then
        log "ERROR" "$dev label mismatch: got '$actual', expected '$expected'"
        exit 1
    else
        log "SUCCESS" "$dev label confirmed: $actual"
    fi
}

create_disk_image() {
    STAGE="create_disk_image"
    log "INFO" "Creating a new raw image: ${RAW_IMG}"
    rm -f "${RAW_IMG}"
    # Allocate 4GB image
    dd if=/dev/zero of="${RAW_IMG}" bs=512 count=0 seek=8388608
    check_status $? "Raw image created" ${STAGE}
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

# partitions:
# 512 MB -> Boot
#   1 GB -> Passive
#   2 GB -> Primary (root)
# 512 MB -> Swap
partition_image() {
    STAGE="partition_image"
    log "INFO" "Creating 4 aligned primary partitions on ${LOOPDISK} using sfdisk"

    sudo sfdisk "${LOOPDISK}" <<-__EOF__
label: dos
unit: sectors

start=2048,    size=524288,  type=ef
start=526336,  size=2097152, type=83
start=2623488, size=4194304, type=83
start=6817792, size=524288,  type=82
__EOF__

    check_status $? "Partitions created" ${STAGE}
}

map_partitions() {
    STAGE="map_partitions"
    log "INFO" "Mapping partitions using kpartx"
    sudo kpartx -v -a "${LOOPDISK}"
    sudo partprobe "${LOOPDISK}"
    sleep 2
    check_status $? "Partitions mapped" ${STAGE}
}

format_partitions() {
    STAGE="format_partitions"
    log "INFO" "Formatting partitions"
    sudo partprobe "${LOOPDISK}"
    sleep 1

    MAPPED_LOOP_NAME=$(basename "$LOOPDISK")
    DISK="/dev/mapper/${MAPPED_LOOP_NAME}p"
    
    mkfs.vfat -F 32 -n boot "${DISK}1"
    check_status $? "boot partition formatted" ${STAGE}

    mkfs.ext4 -L passive "${DISK}2"
    check_status $? "passive partition formatted" ${STAGE}

    mkfs.ext4 -L primary "${DISK}3"
    check_status $? "primary partition formatted" ${STAGE}

    mkswap -L swap "${DISK}4"
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

copy_boot_partition() {
    STAGE="copy_boot_partition"
    log "INFO" "Copying /boot from rootfs into boot partition"

    if [ ! -d "${ROOTFS_DIR}/boot" ]; then
        log "ERROR" "Missing /boot in rootfs (${ROOTFS_DIR}/boot)"
        exit 1
    fi

    rsync -aAX "${ROOTFS_DIR}/boot/" "${BOOT_MOUNT}/"

    if [ -d "${ROOTFS_DIR}/efi" ]; then
        rsync -aAX "${ROOTFS_DIR}/efi/" "${BOOT_MOUNT}/efi/"
        log "INFO" "/efi directory copied to boot partition"
    else
        log "INFO" "No /efi directory found in rootfs — UEFI boot might fail"
    fi

    check_status $? "/boot copied to boot partition" ${STAGE}
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

update_fstab() {
    PARTITION_TYPE=$1
    log "INFO" "Updating /etc/fstab for partition type: ${PARTITION_TYPE}"

    # Detect QEMU
    if grep -qi qemu /proc/cpuinfo 2>/dev/null || [[ "$(systemd-detect-virt 2>/dev/null || true)" == "qemu" ]]; then
        ROOT_DEV="/dev/sda3"
        BOOT_DEV="/dev/sda1"
        SWAP_DEV="/dev/sda4"
    else
        ROOT_DEV="/dev/mmcblk1p3"
        BOOT_DEV="/dev/mmcblk1p1"
        SWAP_DEV="/dev/mmcblk1p4"
    fi

    # Clean fstab without redundant mounts
    cat <<FSTAB > ${PARTITION_TYPE}/etc/fstab
proc            /proc        proc    defaults    0 0
sysfs           /sys         sysfs   defaults    0 0
devpts          /dev/pts     devpts  defaults    0 0
tmpfs           /tmp         tmpfs   defaults    0 0
${ROOT_DEV}     /            ext4    defaults    0 1
# ${BOOT_DEV}   /boot/firmware vfat  ro          0 2
${SWAP_DEV}     none         swap    sw          0 0
FSTAB

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
    log "INFO" "Make sure you have ran build-env-setup and rootfs-env-setup.sh scripts"
    exit 1
fi

create_disk_image
setup_loop_device
clean_first_50MB
partition_image
map_partitions
format_partitions
mount_partition "${DISK}1" "${BOOT_MOUNT}"
mount_partition "${DISK}2" "${PASSIVE_MOUNT}"
mount_partition "${DISK}3" "${PRIMARY_MOUNT}"

check_label "/dev/mapper/$(basename ${LOOPDISK})p1" "boot"
check_label "/dev/mapper/$(basename ${LOOPDISK})p2" "passive"
check_label "/dev/mapper/$(basename ${LOOPDISK})p3" "primary"
check_label "/dev/mapper/$(basename ${LOOPDISK})p4" "swap"

# create board specific manifest and cp its pkds/libs
get_enabled_apps "$COMMON_CONFIG_FILE" "$COM_CONFIG_FILE"
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
copy_all_apps        "$UKAMA_REPO_APP_PKG" "$UKAMA_APP_PKG"
copy_required_libs   "$UKAMA_REPO_LIB_PKG" "$ROOTFS_DIR/lib"
create_manifest_file "$MANIFEST_FILE"      "${APPS[@]}"

copy_rootfs
copy_boot_partition
cp -rf "${MANIFEST_FILE}" "${PRIMARY_MOUNT}"
cp -rf "${MANIFEST_FILE}" "${PASSIVE_MOUNT}"
rm -rf "${MANIFEST_FILE}"
set_permissions
update_fstab "${PRIMARY_MOUNT}"
#update_fstab "${PASSIVE_MOUNT}"

unmount_partition "${BOOT_MOUNT}"
unmount_partition "${PRIMARY_MOUNT}"
unmount_partition "${PASSIVE_MOUNT}"
detach_loop_device

log "SUCCESS" "Disk image creation completed successfully!"
