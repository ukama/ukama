#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

set -e

STAGE="init"
DIR="$(pwd)"

ARCH="armhf"
BOOT_MOUNT="/media/boot"
PRIMARY_MOUNT="/media/primary"
PASSIVE_MOUNT="/media/passive"
UNUSED_MOUNT="/media/unused"

ALPINE_VERSION="3.20"
ALPINE_MINOR_VERSION=".3"
ALPINE_URL_HOST="https://dl-cdn.alpinelinux.org"
ALPINE_URL_PATH="alpine/v${ALPINE_VERSION}/releases/${ARCH}"
ALPINE_URL=${ALPINE_URL_HOST}/${ALPINE_URL_PATH}
ALPINE_TAR="alpine-minirootfs-${ALPINE_VERSION}${ALPINE_MINOR_VERSION}-${ARCH}.tar.gz"

TMP_ROOTFS="/tmp/alpine-rootfs"

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

validate_partition() {
    local partition=$1
    if [ ! -b "${partition}" ]; then
        log "ERROR" "Partition ${partition} does not exist or is not a block device. Aborting."
        exit 1
    fi
}

download_alpine() {
    log "INFO" "Downloading Alpine Linux base..."
    if [[ ! -f "$ALPINE_TAR" ]]; then
        wget "${ALPINE_URL}/${ALPINE_TAR}" -O "$ALPINE_TAR"
    else
        log "INFO" "Using existing copy of Alpine Linux"
    fi
    log "SUCCESS" "Alpine Linux base downloaded."
}

extract_alpine_rootfs() {
    log "INFO" "Extracting Alpine Linux root filesystem..."
    sudo rm -rf "$TMP_ROOTFS"
    sudo mkdir -p "$TMP_ROOTFS"
    sudo tar -xzf "$ALPINE_TAR" -C "$TMP_ROOTFS"
    log "SUCCESS" "Alpine Linux root filesystem extracted."
}

setup_alpine_rootfs() {
    log "INFO" "Setting up Alpine Linux in primary partition..."

    # Copy Alpine to primary
    sudo rsync -aAX --progress "$TMP_ROOTFS/" "$PRIMARY_MOUNT/"

    # Configure basic Alpine settings
    echo "ukamaOS" | sudo tee "$PRIMARY_MOUNT/etc/hostname"
    echo "127.0.0.1 localhost" | sudo tee "$PRIMARY_MOUNT/etc/hosts"
    echo "nameserver 8.8.8.8" | sudo tee "$PRIMARY_MOUNT/etc/resolv.conf"

    # Install essential Alpine packages
    sudo chroot "$PRIMARY_MOUNT" /bin/sh <<'EOF'
apk update
apk add busybox dropbear dhcpcd openrc e2fsprogs
rc-update add networking boot
EOF

    log "SUCCESS" "Alpine Linux setup in primary partition."
}

build_firmware() {
    STAGE="build_firmware"
    local node=$1
    cwd=$(pwd)
    log "INFO" "Building firmware for Node: ${node}"
    cd "${UKAMA_ROOT}/nodes/ukamaOS/firmware"
    make TARGET="${node}"
    check_status $? "Linux kernel build successfull" ${STAGE}
    cd "${cwd}"
}

copy_bootloaders() {
    STAGE="copy_bootloaders"
    log "INFO" "Copying bootloaders to ${BOOT_MOUNT}"

    sudo cp -v ${BOOT1_BIN} ${BOOT_MOUNT}/boot.bin
    sudo cp -v ${BOOT2_BIN} ${BOOT_MOUNT}/

    check_status $? "Bootloaders copied" ${STAGE}
}

copy_to_passive_rootfs() {
    log "INFO" "Copying primary root filesystem to passive partition..."
    sudo rsync -aAX --progress "$PRIMARY_MOUNT/" "$PASSIVE_MOUNT/"
    log "SUCCESS" "Passive root filesystem set up."
}

create_disk_image() {
    STAGE="create_disk_image"
    log "INFO" "Creating a new raw image: ${RAW_IMG}"
    rm -f "${RAW_IMG}"
    dd if=/dev/zero of="${RAW_IMG}" bs=1M count=0 seek=8192
    check_status $? "Raw image created" ${STAGE}
}

clean_first_50MB() {
    STAGE="clean_first_50MB"
    log "INFO" "Cleaning the first 50MB of ${LOOPDISK}"
    sudo dd if=/dev/zero of="${LOOPDISK}" bs=1M count=50
    check_status $? "First 50MB cleaned" "${STAGE}"
}

setup_loop_device() {
    STAGE="setup_loop_device"
    log "INFO" "Attaching ${RAW_IMG} to a loop device"

    EXISTING_LOOP=$(losetup -a | grep "${RAW_IMG}" | grep -v "(deleted)" | cut -d: -f1)
    if [ -n "${EXISTING_LOOP}" ]; then
        LOOPDISK="${EXISTING_LOOP}"
        log "INFO" "Loop device already attached: ${LOOPDISK}"
    else
        LOOPDISK=$(sudo losetup -f --show "${RAW_IMG}")
        if [ -z "${LOOPDISK}" ]; then
            log "ERROR" "Failed to set up loop device for ${RAW_IMG}."
            exit 1
        fi
        log "SUCCESS" "Loop device set up at ${LOOPDISK}"
    fi
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

mount_partition() {
    local partition=$1
    local mount_point=$2
    validate_partition "${partition}"
    log "INFO" "Mounting ${partition} to ${mount_point}"

    if mountpoint -q "$mount_point"; then
        log "INFO" "${mount_point} is already mounted. Skipping."
        return 0
    fi

    sudo mkdir -p "${mount_point}"
    sudo mount "${partition}" "${mount_point}"
    check_status $? "Partition ${partition} mounted to ${mount_point}" "${STAGE}"
}

unmount_partition() {
    local mount_point=$1
    log "INFO" "Unmounting ${mount_point}"

    if ! mountpoint -q "$mount_point"; then
        log "INFO" "${mount_point} is not mounted. Skipping."
        return 0
    fi

    sudo umount "${mount_point}"
    check_status $? "${mount_point} unmounted" "${STAGE}"
}

cleanup() {
    log "INFO" "Cleaning up resources..."
    local mounts=(${BOOT_MOUNT} ${PRIMARY_MOUNT} ${PASSIVE_MOUNT} ${UNUSED_MOUNT})

    for mount in "${mounts[@]}"; do
        if mountpoint -q "${mount}"; then
            sudo umount -R "${mount}" 2>/dev/null || true
        fi
    done

    if [ -n "${LOOPDISK}" ]; then
        log "INFO" "Removing kpartx mappings for ${LOOPDISK}"
        sudo kpartx -dv "${LOOPDISK}" 2>/dev/null || true
        log "INFO" "Detaching loop device ${LOOPDISK}"
        sudo losetup -d "${LOOPDISK}" 2>/dev/null || true
    fi

    sudo rm -rf "$TMP_ROOTFS"
    log "INFO" "Cleanup completed."
}

# Main Execution
if [ $# -lt 4 ]; then
    log "ERROR" "Usage: $0 <node> <ukama-root> <ukamaos_version> <raw_img_name>"
    exit 1
fi

NODE="$1"
UKAMA_ROOT="$2"
UKAMAOS_VERSION="$3"
RAW_IMG="$4"

UKAMA_OS=$(realpath ${UKAMA_ROOT}/nodes/ukamaOS)

if [ "${NODE}" = "anode" ]; then
    BOOT1_BIN="${UKAMA_OS}/firmware/_ukamafs/boot/at91-bootstrap/at91bootstrap.bin"
    BOOT2_BIN="${UKAMA_OS}/firmware/_ukamafs/boot/u-boot/u-boot.bin"
fi

LOOPDISK=

# build images
download_alpine
extract_alpine_rootfs
build_firmware ${NODE}

create_disk_image
setup_loop_device
clean_first_50MB
partition_disk_image
map_partitions
format_partitions
mount_partition "/dev/mapper/$(basename ${LOOPDISK})p1" "$BOOT_MOUNT"
mount_partition "/dev/mapper/$(basename ${LOOPDISK})p2" "$PRIMARY_MOUNT"
mount_partition "/dev/mapper/$(basename ${LOOPDISK})p3" "$PASSIVE_MOUNT"
mount_partition "/dev/mapper/$(basename ${LOOPDISK})p4" "$UNUSED_MOUNT"
copy_bootloaders
setup_alpine_rootfs
copy_to_passive_rootfs
set_permissions
setup_fstab

#cleanup
unmount_partition "$BOOT_MOUNT"
unmount_partition "$PRIMARY_MOUNT"
unmount_partition "$PASSIVE_MOUNT"
unmount_partition "$UNUSED_MOUNT"
detach_loop_device

log "SUCCESS" "Alpine-based SD card image created successfully!"
exit 0

