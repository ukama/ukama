#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

set -e

STAGE="init"
DIR="$(pwd)"
UKAMA_OS=$(realpath ../../ukamaOS)

BOOT_MOUNT="/media/boot"
PRIMARY_MOUNT="/media/primary"
PASSIVE_MOUNT="/media/passive"
UNUSED_MOUNT="/media/unused"

BOOT1_BIN=${UKAMA_OS}/firmware/_ukamafs/boot/at91bootstrap/at91bootstrap.bin
BOOT2_BIN=${UKAMA_OS}/firmware/_ukamafs/boot/uboot/u-boot.bin

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
    local mounts=(${BOOT_MOUNT} ${PRIMARY_MOUNT} ${PASSIVE_MOUNT} ${UNUSED_MOUNT})
    for mount in "${mounts[@]}"; do
        sudo umount "${mount}" 2>/dev/null || true
    done
    sudo kpartx -dv "${RAW_IMG}" 2>/dev/null || true
    sudo losetup -d "${LOOPDISK}" 2>/dev/null || true
    log "INFO" "Cleanup completed."
}

create_image() {
    STAGE="create_image"
    log "INFO" "Creating a new raw image: ${RAW_IMG}"
    rm -f "${RAW_IMG}"
    dd if=/dev/zero of="${RAW_IMG}" bs=1M count=0 seek=8096
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

partition_image() {
    STAGE="partition_image"
    log "INFO" "Creating partitions on ${LOOPDISK} using sfdisk"
    sudo sfdisk "${LOOPDISK}" <<-__EOF__
1M,48M,0xE,*
49M,2048M,,-
2097M,2048M,,-
4145M,,,-
__EOF__
    check_status $? "Partitions created" ${STAGE}
}

map_partitions() {
    STAGE="map_partitions"
    log "INFO" "Mapping partitions using kpartx"
    sudo kpartx -v -a "${LOOPDISK}"
    check_status $? "Partitions mapped" ${STAGE}
}

format_partitions() {
    STAGE="format_partitions"
    log "INFO" "Formatting partitions"
    DEVICE=$(basename "${LOOPDISK}")
    DISK="/dev/mapper/${DEVICE}p"

    sudo mkfs.vfat -F 16 -n boot "${DISK}1"
    check_status $? "boot partition formatted" ${STAGE}

    sudo mkfs.ext4 -L primary "${DISK}2"
    check_status $? "primary rootfs formatted" ${STAGE}

    sudo mkfs.ext4 -L passive "${DISK}3"
    check_status $? "passive rootfs formatted" ${STAGE}

    sudo mkfs.ext4 -L unused "${DISK}4"
    check_status $? "unused partition formatted" ${STAGE}
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

copy_bootloaders() {
    STAGE="copy_bootloaders"
    log "INFO" "Copying bootloaders to ${BOOT_MOUNT}"

    sudo cp -v ${BOOT1_BIN} ${BOOT_MOUNT}/boot.bin
    sudo cp -v ${BOOT2_BIN} ${BOOT_MOUNT}/

    check_status $? "Bootloaders copied" ${STAGE}
}

copy_rootfs() {
    STAGE="copy_rootfs"
    log "INFO" "Copying rootfs to primary and passive"

    IMG_PATH="${UKAMA_OS}/distro/scripts/ukamaOS_initrd_${OS_TARGET}_${OS_VERSION}.img"
    MOUNT_IMG="/mnt/img"

    log "INFO" "Image file ${IMG_PATH} found"
    sudo mkdir -p ${MOUNT_IMG}

    LOOP_DEVICE=$(sudo losetup -fP --show ${IMG_PATH})
    if [ $? -ne 0 ]; then
        log "ERROR" "Failed to set up loop device for ${IMG_PATH}"
        exit 1
    fi

    log "INFO" "Mounted ${IMG_PATH} as ${LOOP_DEVICE}"
    ROOTFS_PARTITION="${LOOP_DEVICE}p1"

    sudo mount ${ROOTFS_PARTITION} ${MOUNT_IMG}
    if [ $? -ne 0 ]; then
        log "ERROR" "Failed to mount rootfs partition ${ROOTFS_PARTITION}"
        sudo losetup -d ${LOOP_DEVICE}
        exit 1
    fi

    log "INFO" "Rootfs partition ${ROOTFS_PARTITION} mounted to ${MOUNT_IMG}"
    sudo rsync -aAX ${MOUNT_IMG}/ ${PRIMARY_MOUNT}
    check_status $? "Rootfs copied to primary" ${STAGE}

    sudo rsync -aAX ${MOUNT_IMG}/ ${PASSIVE_MOUNT}
    check_status $? "Rootfs copied to passive" ${STAGE}

    sudo umount ${MOUNT_IMG}
    sudo losetup -d ${LOOP_DEVICE}
    log "INFO" "Unmounted ${MOUNT_IMG} and detached loop device ${LOOP_DEVICE}"
    sudo rm -rf ${MOUNT_IMG}
    log "INFO" "Rootfs copy operation completed"
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

copy_kernel() {
    STAGE="copy_kernel"
    log "INFO" "Copying kernel to primary and passive"
    sudo cp -v ${UKAMA_OS}/_ukamafs/boot/zImage ${PRIMARY_MOUNT}/boot/zImage
    check_status $? "Kernel copied to primary" ${STAGE}

    sudo cp -v ${UKAMA_OS}/_ukamafs/boot/zImage ${PASSIVE_MOUNT}/boot/zImage
    check_status $? "Kernel copied to passive" ${STAGE}
}

copy_dtbs() {
    STAGE="copy_dtbs"
    log "INFO" "Copying DTBs to primary and passive"
    kernel_version=$(awk '{print $3}' "${UKAMA_OS}/kernel/linux/include/generated/utsrelease.h" | sed 's/\"//g')

    sudo mkdir -p ${PRIMARY_MOUNT}/dtbs/${kernel_version}/
    sudo cp -v ${UKAMA_OS}/_ukamafs/u-boot/arch/arm/dts/*.dtb \
         ${PRIMARY_MOUNT}/dtbs/${kernel_version}/
    sudo cp -v ${UKAMA_OS}/_ukamafs/u-boot/arch/arm/dts/*.dtb \
         ${PRIMARY_MOUNT}/boot/
    check_status $? "DTBs copied to primary" ${STAGE}

    sudo mkdir -p ${PASSIVE_MOUNT}/dtbs/${kernel_version}/
    sudo cp -v ${UKAMA_OS}/_ukamafs/u-boot/arch/arm/dts/*.dtb \
         ${PASSIVE_MOUNT}/dtbs/${kernel_version}/
    sudo cp -v ${UKAMA_OS}/_ukamafs/u-boot/arch/arm/dts/*.dtb \
         ${PASSIVE_MOUNT}/boot/
    check_status $? "DTBs copied to passive" ${STAGE}
}

copy_modules() {
    STAGE="copy_modules"
    log "INFO" "Copying kernel modules to primary and passive"
    sudo cp -a ${UKAMA_OS}/_ukamafs/lib/modules ${PRIMARY_MOUNT}/lib/
    check_status $? "Modules copied to primary" ${STAGE}

    sudo cp -a ${UKAMA_OS}/_ukamafs/lib/modules ${PASSIVE_MOUNT}/lib/
    check_status $? "Modules copied to passive" ${STAGE}
}

setup_fstab() {
    STAGE="setup_fstab"
    log "INFO" "Setting up fstab for primary and passive partitions"

    sudo bash -c "cat << EOF > ${PRIMARY_MOUNT}/etc/fstab
        /dev/mmcblk0p4  /unused  auto  ro,nodev,noexec  0  2
        /dev/mmcblk0p3  /fs3     auto  ro              0  2
        /dev/mmcblk0p2  /        auto  errors=remount-ro  0  1
EOF"
    check_status $? "fstab setup for primary" ${STAGE}

    sudo bash -c "cat << EOF > ${PASSIVE_MOUNT}/etc/fstab
        /dev/mmcblk0p4  /unused  auto  ro,nodev,noexec  0  2
        /dev/mmcblk0p2  /fs2     auto  ro              0  1
        /dev/mmcblk0p3  /        auto  errors=remount-ro  0  1
EOF"
    check_status $? "fstab setup for passive" ${STAGE}
}

detach_loop_device() {
    STAGE="detach_loop_device"
    log "INFO" "Detaching loop device and cleaning up"
    sudo kpartx -dv "${LOOPDISK}"
    sudo losetup -d "${LOOPDISK}"
    check_status $? "Loop device detached" ${STAGE}
}

# Execution entry
if [ $# -lt 3 ]; then
    log "ERROR" "Usage: $0 <output_image_file> <os_target> <os_version>"
    exit 1
fi

RAW_IMG="$1"
OS_TARGET="$2"
OS_VERSION="$3"
IMG_PATH="${UKAMA_OS}/distro/scripts/ukamaOS_initrd_${OS_TARGET}_${OS_VERSION}.img"

# Sanity check.
if [ ! -f "${BOOT1_BIN}" ]; then
    log "ERROR" "boot file ${BOOT1_BIN} does not exist"
    exit 1
fi

if [ ! -f "${BOOT2_BIN}" ]; then
    log "ERROR" "boot file ${BOOT2_BIN} does not exist"
    exit 1
fi

if [ ! -f "${IMG_PATH}" ]; then
    log "ERROR" "Image file ${IMG_PATH} does not exist"
    exit 1
fi

# Execute each step
create_image
setup_loop_device
clean_first_50MB
partition_image
map_partitions
format_partitions
mount_partition "${DISK}1" "${BOOT_MOUNT}"
mount_partition "${DISK}2" "${PRIMARY_MOUNT}"
mount_partition "${DISK}3" "${PASSIVE_MOUNT}"
mount_partition "${DISK}4" "${UNUSED_MOUNT}"
copy_bootloaders
copy_rootfs
set_permissions
copy_kernel
copy_dtbs
copy_modules
setup_fstab
unmount_partition "${BOOT_MOUNT}"
unmount_partition "${PRIMARY_MOUNT}"
unmount_partition "${PASSIVE_MOUNT}"
unmount_partition "${UNUSED_MOUNT}"
detach_loop_device

log "SUCCESS" "SD card image creation completed successfully!"
