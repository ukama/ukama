#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

set -e

STAGE="init"
DIR="$(pwd)"
UKAMA_OS=`realpath ../../ukamaOS`

log_info() {
    echo -e "\033[1;34mINFO: $1\033[0m"
}

log_success() {
    echo -e "\033[1;32mSUCCESS: $1\033[0m"
}

log_error() {
    echo -e "\033[1;31mERROR: $1\033[0m" >&2
}

check_status() {
    if [ $1 -eq 0 ]; then
        log_success "$2"
    else
        log_error "Script failed at stage: $3"
        cleanup
        exit 1
    fi
}

cleanup() {

    log_info "Cleaning up resources..."

    sudo umount /media/boot       || true
    sudo umount /media/primary    || true
    sudo umount /media/passive    || true
    sudo umount /media/unused     || true
    sudo kpartx -dv "${RAW_IMG}"  || true
    sudo losetup -d "${LOOPDISK}" || true

    log_info "Cleanup completed."
}

# Check for arguments
if [ $# -lt 3 ]; then
    log_error "No arguments supplied. "
    log_error "Usage: $0 <output_image_file> <os_target> <os_version>"
    exit 1
fi

RAW_IMG="$1"
OS_TARGET="$2"
OS_VERSION="$3"

# Main execution
log_info "Starting SD card image creation script."

create_image() {

    STAGE="create_image"
    log_info "Creating a new raw image: ${RAW_IMG}"

    rm -f "${RAW_IMG}"
    dd if=/dev/zero of="${RAW_IMG}" bs=1M count=0 seek=8096
    check_status $? "Raw image created" ${STAGE}
}

setup_loop_device() {

    STAGE="setup_loop_device"
    log_info "Attaching ${RAW_IMG} to a loop device"

    LOOPDISK=$(sudo losetup -f --show "${RAW_IMG}")
    if [ -z "${LOOPDISK}" ]; then
        log_error "Failed to set up loop device for ${RAW_IMG}."
        cleanup
        exit 1
    fi
    log_success "Loop device set up at ${LOOPDISK}"
}

clean_first_50MB() {

    STAGE="clean_first_50MB"
    log_info "Cleaning the first 50MB of ${LOOPDISK}"

    sudo dd if=/dev/zero of="${LOOPDISK}" bs=1M count=50
    check_status $? "First 50MB cleaned" ${STAGE}
}

partition_image() {

    STAGE="partition_image"
    log_info "Creating partitions on ${LOOPDISK} using sfdisk"

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
    log_info "Mapping partitions using kpartx"

    sudo kpartx -v -a "${LOOPDISK}"
    check_status $? "Partitions mapped" ${STAGE}
}

format_partitions() {

    STAGE="format_partitions"
    log_info "Formatting partitions"

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

mount_partitions() {

    STAGE="mount_partitions"
    log_info "Mounting partitions"

    sudo mkdir -p /media/boot /media/primary /media/passive /media/unused

    sudo mount "${DISK}1" /media/boot
    check_status $? "BOOT partition mounted" ${STAGE}

    sudo mount "${DISK}2" /media/primary
    check_status $? "Primary partition mounted" ${STAGE}

    sudo mount "${DISK}3" /media/passive
    check_status $? "Passive partition mounted" ${STAGE}

    sudo mount "${DISK}4" /media/unused
    check_status $? "Unused partition mounted" ${STAGE}
}

copy_bootloaders() {

    STAGE="copy_bootloaders"
    log_info "Copying bootloaders to /media/boot"

    sudo cp -v ${UKAMA_OS}/firmware/_ukamafs/boot/at91bootstrap/at91bootstrap.bin \
         /media/boot/boot.bin
    sudo cp -v ${UKAMA_OS}/firmware/_ukamafs/boot/uboot/u-boot.bin \
         /media/boot/

    check_status $? "Bootloaders copied" ${STAGE}
}

copy_rootfs() {

    STAGE="copy_rootfs"
    log_info "Copying rootfs to primary and passive"

    IMG_PATH="${UKAMA_OS}/distro/scripts/ukamaOS_initrd_$1_$2.img"
    MOUNT_IMG="/mnt/img"
    MOUNT_PRIMARY="/media/primary/"
    MOUNT_PASSIVE="/media/passive/"

    # Step 0: Check if the .img file exists
    if [ ! -f "${IMG_PATH}" ]; then
        log_error "Image file ${IMG_PATH} does not exist" ${STAGE}
        return 1
    fi

    log_info "Image file ${IMG_PATH} found"

    # Create temporary mount point for the image
    sudo mkdir -p ${MOUNT_IMG}

    # Step 1: Set up loop device for the .img file
    LOOP_DEVICE=$(sudo losetup -fP --show ${IMG_PATH})
    if [ $? -ne 0 ]; then
        log_error "Failed to set up loop device for ${IMG_PATH}" ${STAGE}
        return 1
    fi

    log_info "Mounted ${IMG_PATH} as ${LOOP_DEVICE}"

    # Step 2: Find the partition containing rootfs (if needed)
    ROOTFS_PARTITION="${LOOP_DEVICE}p1"

    # Step 3: Mount the rootfs from the image
    sudo mount ${ROOTFS_PARTITION} ${MOUNT_IMG}
    if [ $? -ne 0 ]; then
        log_error "Failed to mount rootfs partition ${ROOTFS_PARTITION}" ${STAGE}
        sudo losetup -d ${LOOP_DEVICE}
        return 1
    fi

    log_info "Rootfs partition ${ROOTFS_PARTITION} mounted to ${MOUNT_IMG}"

    # Step 4: Copy rootfs from .img to both primary and passive
    sudo rsync -aAX ${MOUNT_IMG}/ ${MOUNT_PRIMARY}
    check_status $? "Rootfs copied to primary" ${STAGE}

    sudo rsync -aAX ${MOUNT_IMG}/ ${MOUNT_PASSIVE}
    check_status $? "Rootfs copied to passive" ${STAGE}

    # Step 5: Unmount and clean up
    sudo umount ${MOUNT_IMG}
    sudo losetup -d ${LOOP_DEVICE}
    log_info "Unmounted ${MOUNT_IMG} and detached loop device ${LOOP_DEVICE}"

    sudo rm -rf ${MOUNT_IMG}

    log_info "Rootfs copy operation completed"
}

set_permissions() {

    STAGE="set_permissions"
    log_info "Setting permissions for primary and passive partitions"

    sudo chown -R root:root /media/primary/
    sudo chmod -R 755 /media/primary/
    check_status $? "Permissions set for primary" ${STAGE}

    sudo chown -R root:root /media/passive/
    sudo chmod -R 755 /media/passive/

    check_status $? "Permissions set for passive" ${STAGE}
}

copy_kernel() {

    STAGE="copy_kernel"
    log_info "Copying kernel to primary and passive"

    sudo cp -v ${UKAMA_OS}/_ukamafs/boot/zImage /media/primary/boot/zImage
    check_status $? "Kernel copied to primary" ${STAGE}

    sudo cp -v ${UKAMA_OS}/_ukamafs/boot/zImage /media/passive/boot/zImage
    check_status $? "Kernel copied to passive" ${STAGE}
}

copy_dtbs() {

    STAGE="copy_dtbs"
    log_info "Copying DTBs to primary and passive"

    kernel_version=$(awk '{print $3}' "${UKAMA_OS}/kernel/linux/include/generated/utsrelease.h" \
                         | sed 's/\"//g')

    sudo mkdir -p /media/primary/dtbs/${kernel_version}/
    sudo cp -v ${UKAMA_OS}/_ukamafs/u-boot/arch/arm/dts/*.dtb /media/primary/dtbs/${kernel_version}/
    sudo cp -v ${UKAMA_OS}/_ukamafs/u-boot/arch/arm/dts/*.dtb /media/primary/boot/ 
    check_status $? "DTBs copied to primary" ${STAGE}

    sudo mkdir -p /media/passive/dtbs/${kernel_version}/
    sudo cp -v ${UKAMA_OS}/_ukamafs/u-boot/arch/arm/dts/*.dtb /media/passive/dtbs/${kernel_version}/
    sudo cp -v ${UKAMA_OS}/_ukamafs/u-boot/arch/arm/dts/*.dtb /media/passive/boot/ 
    check_status $? "DTBs copied to passive" ${STAGE}
}

copy_modules() {

    STAGE="copy_modules"
    log_info "Copying kernel modules to primary and passive"

    sudo cp -a ${UKAMA_OS}/_ukamafs/lib/modules /media/primary/lib/
    check_status $? "Modules copied to primary" ${STAGE}

    sudo cp -a ${UKAMA_OS}/_ukamafs/lib/modules /media/passive/lib/
    check_status $? "Modules copied to passive" ${STAGE}
}

setup_fstab() {

    STAGE="setup_fstab"
    log_info "Setting up fstab for primary and passive partitions"

    # Primary
    sudo bash -c "cat << EOF > /media/primary/etc/fstab
        /dev/mmcblk0p4  /unused  auto  ro,nodev,noexec  0  2
        /dev/mmcblk0p3  /fs3     auto  ro              0  2
        /dev/mmcblk0p2  /        auto  errors=remount-ro  0  1
        EOF"
    check_status $? "fstab setup for primary" ${STAGE}

    # Passive
    sudo bash -c "cat << EOF > /media/passive/etc/fstab
        /dev/mmcblk0p4  /unused  auto  ro,nodev,noexec  0  2
        /dev/mmcblk0p2  /fs2     auto  ro              0  1
        /dev/mmcblk0p3  /        auto  errors=remount-ro  0  1
        EOF"
    check_status $? "fstab setup for passive" ${STAGE}
}

unmount_partitions() {

    STAGE="unmount_partitions"
    log_info "Unmounting partitions"

    sudo umount /media/boot
    sudo umount /media/primary
    sudo umount /media/passive
    sudo umount /media/unused

    check_status $? "Partitions unmounted" ${STAGE}
}

detach_loop_device() {

    STAGE="detach_loop_device"
    log_info "Detaching loop device and cleaning up"

    sudo kpartx -dv "${LOOPDISK}"
    sudo losetup -d "${LOOPDISK}"
    check_status $? "Loop device detached" ${STAGE}
}

# Execute each step
create_image
setup_loop_device
clean_first_50MB
partition_image
map_partitions
format_partitions
mount_partitions
copy_bootloaders
copy_rootfs ${OS_TARGET} ${OS_VERSION}
set_permissions
copy_kernel
copy_dtbs
copy_modules
setup_fstab
unmount_partitions
detach_loop_device

log_success "SD card image creation completed successfully!"
