#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

set -e

STAGE="init"
DIR="$(pwd)"
UKAMA_OS=$(realpath ../../../nodes/ukamaOS)
UKAMA_ROOT=$(realpath ../../../)
BOOT_MOUNT="/media/boot"
RECOVERY_MOUNT="/media/recovery"
PRIMARY_MOUNT="/media/primary"
PASSIVE_MOUNT="/media/passive"
DATA_MOUNT="/media/data"

RAW_IMG="ukama-amplifier-node.img"

BOOT1_BIN=${UKAMA_OS}/firmware/build/boot/at91bootstrap/at91bootstrap.bin
BOOT2_BIN=${UKAMA_OS}/firmware/build/boot/uboot/u-boot.bin

ROOTFS_DIR=${UKAMA_ROOT}/builder/scripts/build-system/rootfs 
MISC_FILES_DIR=${UKAMA_ROOT}/builder/scripts/build-system/amplifier

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

create_disk_image() {
    STAGE="create_disk_image"
    log "INFO" "Creating a new raw image: ${RAW_IMG}"
    rm -f "${RAW_IMG}"
    dd if=/dev/zero of="${RAW_IMG}" bs=512 count=0 seek=61120512
    check_status $? "Raw image created" ${STAGE}
}

build_firmware() {
    STAGE="build_firmware"
    local node=$1
    local path="${UKAMA_ROOT}/nodes/ukamaOS/firmware"
    cwd=$(pwd)
    log "INFO" "Building firmware for Node: ${node}"

    cd "${path}"
    make clean TARGET="${node}" ROOTFSPATH="${path}/build"
    make TARGET="${node}" ROOTFSPATH="${path}/build"
    check_status $? "Firmware (at91 and uboot) build successful" ${STAGE}
    cd "${cwd}"
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
	rsync -aAXv --exclude={"/dev","/sys","/proc"} ${ROOTFS_DIR}/* ${PRIMARY_MOUNT}/
	mkdir -p ${PRIMARY_MOUNT}/dev ${PRIMARY_MOUNT}/sys ${PRIMARY_MOUNT}/proc	

	log "INFO" "Copying rootfs to primary and passive"
	rsync -aAXv --exclude={"/dev","/sys","/proc"} ${ROOTFS_DIR}/* ${PASSIVE_MOUNT}/
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

    if [[ "${PARTITION_TYPE}" == ${PRIMARY_MOUNT} ]]; then
        cat <<FSTAB > ${PRIMARY_MOUNT}/etc/fstab
proc            /proc        proc    defaults    0 0
sysfs           /sys         sysfs   defaults    0 0
devpts          /dev/pts     devpts  defaults    0 0
tmpfs           /tmp         tmpfs   defaults    0 0
/dev/mmcblk1p2  /recovery    auto    ro          0 2
/dev/mmcblk1p7  /data        auto    ro          0 2
/dev/mmcblk1p6  /passive     auto    ro          0 2
/dev/mmcblk1p5  /            auto    errors=remount-ro  0 1
/dev/mmcblk1p1  /boot/firmware auto  ro          0 2
FSTAB
    else
        cat <<FSTAB > ${PASSIVE_MOUNT}/etc/fstab
proc            /proc        proc    defaults    0 0
sysfs           /sys         sysfs   defaults    0 0
devpts          /dev/pts     devpts  defaults    0 0
tmpfs           /tmp         tmpfs   defaults    0 0
/dev/mmcblk1p2  /recovery    auto    ro          0 2
/dev/mmcblk1p7  /data        auto    ro          0 2
/dev/mmcblk1p5  /passive     auto    ro          0 2
/dev/mmcblk1p6  /            auto    errors=remount-ro  0 1
/dev/mmcblk1p1  /boot/firmware auto  ro          0 2
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
if [ -d "${ROOTFS_DIR}" ] && [ "$(ls -A ${ROOTFS_DIR})" ]; then
    log "INFO" "ROOTFS exist."
else
    log "ERROR" "${ROOTFS_DIR} does not exist"
    log "INFO" "Make sure you have ran build-env-setup and rootfs-env-setup.sh scripts"
    exit 1
fi

build_firmware "amplifier"
if [ ! -f "${BOOT1_BIN}" ]; then
    log "ERROR" "boot file ${BOOT1_BIN} does not exist"
    exit 1
fi

if [ ! -f "${BOOT2_BIN}" ]; then
    log "ERROR" "boot file ${BOOT2_BIN} does not exist"
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
copy_bootloaders
copy_rootfs
set_permissions
update_fstab "${PRIMARY_MOUNT}"
update_fstab "${PASSIVE_MOUNT}"

#these files are secific to amplifier nodes
copy_misc_files "${PRIMARY_MOUNT}"
copy_misc_files "${PASSIVE_MOUNT}"

unmount_partition "${BOOT_MOUNT}"
unmount_partition "${PRIMARY_MOUNT}"
unmount_partition "${PASSIVE_MOUNT}"
detach_loop_device

log "SUCCESS" "Disk image creation completed successfully!"
