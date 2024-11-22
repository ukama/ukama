#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

set -e

STAGE="init"
DIR="$(pwd)"

BOOT_MOUNT="/media/boot"
PRIMARY_MOUNT="/media/primary"
PASSIVE_MOUNT="/media/passive"
UNUSED_MOUNT="/media/unused"

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
    log "INFO" "Cleaning up resources..."
    local mounts=(${BOOT_MOUNT} ${PRIMARY_MOUNT} ${PASSIVE_MOUNT} ${UNUSED_MOUNT})
    local cwd=$(pwd)

    for mount in "${mounts[@]}"; do
        sudo umount "${mount}" 2>/dev/null || true
    done

    if [ -n "${LOOPDISK}" ]; then
        log "INFO" "Removing kpartx mappings for ${LOOPDISK}"
        sudo kpartx -dv "${LOOPDISK}" 2>/dev/null || true
        log "INFO" "Detaching loop device ${LOOPDISK}"
        sudo losetup -d "${LOOPDISK}" 2>/dev/null || true
    fi

    cd ${UKAMA_OS}/firmware; make distclean
    cd ${UKAMA_OS}/kernel;   make distclean
    cd "${cwd}"; sudo rm -rf _ukama_os_rootfs
     rm -rf ${UKAMA_ROOT}/builder/scripts/ukamaOS_initrd_${NODE}_${UKAMAOS_VERSION}.img
    log "INFO" "Cleanup completed."
}

build_firmware() {
    STAGE="build_linux_kernel"
    local node=$1
    cwd=$(pwd)
    log "INFO" "Building firmware for Node: ${node}"
    cd "${UKAMA_ROOT}/nodes/ukamaOS/firmware"
    make TARGET="${node}"
    check_status $? "Linux kernel build successfull" ${STAGE}
    cd "${cwd}"
}

build_linux_kernel() {
    STAGE="build_linux_kernel"
    local node=$1
    cwd=$(pwd)
    log "INFO" "Building linux kernel for Node: ${node}"
    cd "${UKAMA_ROOT}/nodes/ukamaOS/kernel"
    make distclean
    make TARGET="${node}"
    check_status $? "Linux kernel build successfull" ${STAGE}
    cd "${cwd}"
}

build_ukamaos() {
    STAGE="build_ukamaos"
    local node=$1
    cwd=$(pwd)
    log "INFO" "Building ukamaOS for Node: ${node}"
    cd ${UKAMA_ROOT}/builder/scripts/
    ./build-ukamaos.sh ${node}
    check_status $? "ukamaOS build successfull" ${STAGE}
    cd "${cwd}"
}

create_disk_image() {
    STAGE="create_disk_image"
    log "INFO" "Creating a new raw image: ${RAW_IMG}"
    rm -f "${RAW_IMG}"
    dd if=/dev/zero of="${RAW_IMG}" bs=1M count=0 seek=8096
    check_status $? "Raw image created" ${STAGE}
}

clean_first_50MB() {
    STAGE="clean_first_50MB"
    log "INFO" "Cleaning the first 50MB of ${LOOPDISK}"
    sudo dd if=/dev/zero of="${LOOPDISK}" bs=1M count=50
    check_status $? "First 50MB cleaned" "${STAGE}"
}

partition_disk_image() {
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

    IMG_PATH="${UKAMA_ROOT}/builder/scripts/ukamaOS_initrd_${NODE}_${UKAMAOS_VERSION}.img"
    PRIMARY_PARTITION="/dev/mapper/$(basename ${LOOPDISK})p2"
    PASSIVE_PARTITION="/dev/mapper/$(basename ${LOOPDISK})p3"
    TMP_DIR="/mnt/rootfs_extracted"
    MOUNT_IMG="/mnt/img"

    # Validate the image file
    if [ ! -f "${IMG_PATH}" ]; then
        log "ERROR" "Image file ${IMG_PATH} not found. Aborting."
        exit 1
    fi

    # Extract the initramfs contents
    log "INFO" "Extracting rootfs image ${IMG_PATH}"
    sudo mkdir -p "${TMP_DIR}"
    sudo rm -rf "${TMP_DIR}/*"  # Ensure the directory is clean
    sudo gunzip -c "${IMG_PATH}" | sudo cpio -idmv -D "${TMP_DIR}"
    check_status $? "Rootfs image extracted to ${TMP_DIR}" ${STAGE}

    # Copy to primary partition
    log "INFO" "Mounting primary partition ${PRIMARY_PARTITION}"
    sudo mkdir -p "${MOUNT_IMG}"
    sudo mount "${PRIMARY_PARTITION}" "${MOUNT_IMG}"
    check_status $? "Primary partition mounted to ${MOUNT_IMG}" ${STAGE}

    log "INFO" "Copying rootfs to primary partition"
    sudo rsync -aAX --progress "${TMP_DIR}/" "${MOUNT_IMG}/"
    check_status $? "Rootfs copied to primary partition" ${STAGE}
    sudo umount "${MOUNT_IMG}"
    log "INFO" "Unmounted ${MOUNT_IMG}"

    # Copy to passive partition
    log "INFO" "Mounting passive partition ${PASSIVE_PARTITION}"
    sudo mount "${PASSIVE_PARTITION}" "${MOUNT_IMG}"
    check_status $? "Passive partition mounted to ${MOUNT_IMG}" ${STAGE}

    log "INFO" "Copying rootfs to passive partition"
    sudo rsync -aAX --progress "${TMP_DIR}/" "${MOUNT_IMG}/"
    check_status $? "Rootfs copied to passive partition" ${STAGE}
    sudo umount "${MOUNT_IMG}"
    log "INFO" "Unmounted ${MOUNT_IMG}"

    # Cleanup
    sudo rm -rf "${TMP_DIR}"
    sudo rm -rf "${MOUNT_IMG}"
    log "INFO" "Rootfs copy operation completed"
}

copy_linux_kernel() {
    STAGE="copy_linux_kernel"
    log "INFO" "Copying kernel to primary and passive"

    sudo mkdir -p ${PRIMARY_MOUNT}/boot/
    sudo mkdir -p ${PASSIVE_MOUNT}/boot/

    sudo cp -v ${UKAMA_ROOT}/nodes/ukamaOS/kernel/_ukamafs/boot/zImage \
         ${PRIMARY_MOUNT}/boot/zImage
    check_status $? "Kernel copied to primary" ${STAGE}

    sudo cp -v ${UKAMA_ROOT}/nodes/ukamaOS/kernel/_ukamafs/boot/zImage \
         ${PASSIVE_MOUNT}/boot/zImage
    check_status $? "Kernel copied to passive" ${STAGE}
}

copy_dtbs() {
    STAGE="copy_dtbs"
    log "INFO" "Copying DTBs to primary and passive"
    kernel_version=$(awk '{print $3}' "${UKAMA_OS}/kernel/linux/include/generated/utsrelease.h" | sed 's/\"//g')

    sudo mkdir -p ${PRIMARY_MOUNT}/dtbs/${kernel_version}/

    sudo cp -v ${UKAMA_OS}/firmware/u-boot/arch/arm/dts/*.dtb \
         ${PRIMARY_MOUNT}/dtbs/${kernel_version}/
    sudo cp -v ${UKAMA_OS}/firmware/u-boot/arch/arm/dts/*.dtb \
         ${PRIMARY_MOUNT}/boot/
    check_status $? "DTBs copied to primary" ${STAGE}

    sudo mkdir -p ${PASSIVE_MOUNT}/dtbs/${kernel_version}/
    sudo cp -v ${UKAMA_OS}/firmware/u-boot/arch/arm/dts/*.dtb \
         ${PASSIVE_MOUNT}/dtbs/${kernel_version}/
    sudo cp -v ${UKAMA_OS}/firmware/u-boot/arch/arm/dts/*.dtb \
         ${PASSIVE_MOUNT}/boot/
    check_status $? "DTBs copied to passive" ${STAGE}
}

copy_modules() {
    STAGE="copy_modules"
    log "INFO" "Copying kernel modules to primary and passive"

    sudo mkdir -p ${PRIMARY_MOUNT}/lib/
    sudo mkdir -p ${PASSIVE_MOUNT}/lib/

    sudo cp -a ${UKAMA_OS}/kernel/_ukamafs/lib/modules ${PRIMARY_MOUNT}/lib/
    check_status $? "Modules copied to primary" ${STAGE}

    sudo cp -a ${UKAMA_OS}/kernel/_ukamafs/lib/modules ${PASSIVE_MOUNT}/lib/
    check_status $? "Modules copied to passive" ${STAGE}
}

setup_loop_device() {
    STAGE="setup_loop_device"
    log "INFO" "Attaching ${RAW_IMG} to a loop device"

    # Check if the image is already attached to a loop device
    EXISTING_LOOP=$(losetup -a | grep "${RAW_IMG}" | grep -v "(deleted)" | cut -d: -f1)
    if [ -n "${EXISTING_LOOP}" ]; then
        LOOPDISK="${EXISTING_LOOP}"
        log "INFO" "Loop device already attached: ${LOOPDISK}"
    else
        # Attach a new loop device
        LOOPDISK=$(sudo losetup -f --show "${RAW_IMG}")
        if [ -z "${LOOPDISK}" ]; then
            log "ERROR" "Failed to set up loop device for ${RAW_IMG}."
            exit 1
        fi
        log "SUCCESS" "Loop device set up at ${LOOPDISK}"
    fi
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

setup_fstab() {
    STAGE="setup_fstab"
    log "INFO" "Setting up fstab for primary and passive partitions"

    sudo mkdir -p "${PRIMARY_MOUNT}/etc/"
    sudo mkdir -p "${PASSIVE_MOUNT}/etc/"
   
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
    sudo kpartx -dv "${LOOPDISK}" 2>/dev/null || true
    sudo losetup -d "${LOOPDISK}" 2>/dev/null || true
    check_status $? "Loop device detached" "${STAGE}"
}

# Execution entry
if [ $# -lt 3 ]; then
    log "ERROR" "Usage: $0 <node> <ukama-root> <ukamaos_version>"
    exit 1
fi

NODE="$1"
UKAMA_ROOT="$2"
UKAMAOS_VERSION="$3"
RAW_IMG="$4"

LOOPDISK=
UKAMA_OS=$(realpath ${UKAMA_ROOT}/nodes/ukamaOS)
IMG_PATH="${UKAMA_ROOT}/builder/scripts/ukamaOS_initrd_${NODE}_${UKAMAOS_VERSION}.img"
if [ "${NODE}" = "anode" ]; then
    BOOT1_BIN="${UKAMA_OS}/firmware/_ukamafs/boot/at91-bootstrap/at91bootstrap.bin"
    BOOT2_BIN="${UKAMA_OS}/firmware/_ukamafs/boot/u-boot/u-boot.bin"
fi

if [ -d "_ukama_os_rootfs" ]
then
    sudo rm -rf ./_ukama_os_rootfs/
    log "Removed existing copy of _ukama_os_rootfs/"
fi

# build images
build_linux_kernel ${NODE}
build_firmware     ${NODE}
build_ukamaos      ${NODE}

# and rest
create_disk_image
setup_loop_device
clean_first_50MB
partition_disk_image
map_partitions
format_partitions
mount_partition "${DISK}1" "${BOOT_MOUNT}"
mount_partition "${DISK}2" "${PRIMARY_MOUNT}"
mount_partition "${DISK}3" "${PASSIVE_MOUNT}"
mount_partition "${DISK}4" "${UNUSED_MOUNT}"
copy_bootloaders
copy_rootfs
set_permissions
copy_linux_kernel
if [ "${NODE}" = "anode" ]; then
    copy_dtbs
    copy_modules
fi
setup_fstab

# cleanp
unmount_partition "${BOOT_MOUNT}"
unmount_partition "${PRIMARY_MOUNT}"
unmount_partition "${PASSIVE_MOUNT}"
unmount_partition "${UNUSED_MOUNT}"
detach_loop_device

log "SUCCESS" "SD card image creation completed successfully!"

exit 0
