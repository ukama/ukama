#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

set -e

# Variables
ALPINE_VERSION="3.20.3"
ALPINE_ARCH="armv7"
IMG_SIZE="1G"
IMG_NAME="access-node.img"

ALPINE_URL_HOST="https://dl-cdn.alpinelinux.org"
ALPINE_URL_PATH="alpine/v${ALPINE_VERSION%.*}/releases/${ALPINE_ARCH}"
ALPINE_TAR="alpine-minirootfs-${ALPINE_VERSION}-${ALPINE_ARCH}.tar.gz"
ALPINE_URL=${ALPINE_URL_HOST}/${ALPINE_URL_PATH}/${ALPINE_TAR}

BOOT_MOUNT="/mnt/access-node/boot"
PRIMARY_MOUNT="//mnt/access-node/primary"
PASSIVE_MOUNT="/mnt/access-node/passive"
UNUSED_MOUNT="/mnt/access-node/unused"

TMP_DIR="/tmp/access-node/"
TMP_ROOTFS="${TMP_DIR}/alpine-rootfs"
TMP_LINUX="${TMP_DIR}/linux"

CWD=$(pwd)

trap cleanup EXIT

function log() {
    local type="$1"
    local message="$2"
    echo -e "\033[1;34m${type}:\033[0m ${message}"
}

function check_command() {
    command -v "$1" >/dev/null 2>&1 || {
        log "ERROR" "Command '$1' not found. Please install it."
        exit 1
    }
}

function check_requirements() {
    log "INFO" "Checking required commands..."
    for cmd in git make parted rsync wget tar dd losetup mkfs.vfat mkfs.ext4; do
        check_command "$cmd"
    done
    log "SUCCESS" "All required commands are available."
}

function check_status() {
    if [ $1 -ne 0 ]; then
        log "ERROR" "$3"
        exit 1
    fi
    log "SUCCESS" "$2"
}

function cleanup() {
    log "INFO" "Cleaning up resources..."
    sudo umount -R "$BOOT_MOUNT" \
         "$PRIMARY_MOUNT" \
         "$PASSIVE_MOUNT" \
         "$UNUSED_MOUNT" 2>/dev/null || true
    sudo rm -rf "$TMP_DIR"
    sudo losetup -d "$LOOP_DEVICE" 2>/dev/null || true
    log "INFO" "Cleanup completed."
}

function build_linux_kernel() {
    log "INFO" "Building linux kernel..."
    cd "$TMP_DIR"
    git clone --depth=1 https://github.com/raspberrypi/linux || true
    cd "$TMP_LINUX"
    make -j6 ARCH=arm64 CROSS_COMPILE=aarch64-linux-gnu- bcm2711_defconfig  || true
    make -j6 ARCH=arm64 CROSS_COMPILE=aarch64-linux-gnu- Image modules dtbs || true
    cd "$TMP_DIR"
    log "INFO" "Linux kernel build completed."
}

function download_alpine_rootfs() {
    log "INFO" "Downloading Alpine Linux Rootfs..."
    if [[ ! -f "$ALPINE_TAR" ]]; then
        wget "$ALPINE_URL" -O "$ALPINE_TAR"
    else
        log "INFO" "Using cached Alpine image."
    fi
    log "SUCCESS" "Alpine Linux downloaded."
}

function copy_rootfs() {
    log "INFO" "Extracting Alpine Linux root filesystem..."
    cd "$TMP_DIR"
    sudo mkdir -p "$TMP_ROOTFS"
    sudo tar -xzf "$ALPINE_TAR" -C "$TMP_ROOTFS"
    log "SUCCESS" "Alpine Linux extracted."
}

function copy_linux_kernel() {
    log "INFO" "Building linux kernel..."
    local build_dir="${CWD}/build_access_node/"

    cd "$TMP_LINUX"
    sudo mkdir -p ${BOOT_MOUNT}/overlays/
    env PATH=$PATH make -j6 ARCH=arm64 \
        CROSS_COMPILE=aarch64-linux-gnu- INSTALL_MOD_PATH=${PRIMARY} modules_install
    cp arch/arm64/boot/Image               ${BOOT_MOUNT}/kernel8.img
    cp arch/arm64/boot/dts/broadcom/*.dtb  ${BOOT_MOUNT}/
    cp arch/arm64/boot/dts/overlays/*.dtb* ${BOOT_MOUNT}/overlays/
    cp arch/arm64/boot/dts/overlays/README ${BOOT_MOUNT}/overlays/

    # also copy kernel and dtb - needed to run in QEMU
    cp arch/arm64/boot/Image               "${build_dir}/kernel8.img"
    cp arch/arm64/boot/dts/broadcom/*.dtb  "${build_dir}/"
    cp arch/arm64/boot/dts/overlays/*.dtb* "${build_dir}/overlays/"
    cp arch/arm64/boot/dts/overlays/README "${build_dir}/overlays/"

    log "SUCCESS" "Linux kernel copied"
}

function create_disk_image() {
    log "INFO" "Creating disk image..."
    dd if=/dev/zero of="$IMG_NAME" bs=1M count=0 seek=$((1024 * ${IMG_SIZE%G}))
    LOOP_DEVICE=$(sudo losetup -fP --show "$IMG_NAME")
    log "SUCCESS" "Disk image created: $LOOP_DEVICE"
}

function partition_disk_image() {
    log "INFO" "Partitioning disk image..."
    sudo parted --script "$LOOP_DEVICE" mklabel msdos
    sudo parted --script "$LOOP_DEVICE" mkpart primary fat32  1MiB  48MiB
    sudo parted --script "$LOOP_DEVICE" mkpart primary ext4  49MiB 300MiB
    sudo parted --script "$LOOP_DEVICE" mkpart primary ext4 301MiB 600MiB
    sudo parted --script "$LOOP_DEVICE" mkpart primary ext4 601MiB   100%
    log "SUCCESS" "Disk image partitioned."
}

function format_partitions() {
    log "INFO" "Formatting partitions..."
    sudo mkfs.vfat -F 16 -n boot ${LOOP_DEVICE}p1
    sudo mkfs.ext4 -L primary    ${LOOP_DEVICE}p2
    sudo mkfs.ext4 -L passive    ${LOOP_DEVICE}p3
    sudo mkfs.ext4 -L unused     ${LOOP_DEVICE}p4
    log "SUCCESS" "Partitions formatted."
}

function mount_partitions() {
    log "INFO" "Mounting partitions..."
    sudo mkdir -p "$BOOT_MOUNT" "$PRIMARY_MOUNT" "$PASSIVE_MOUNT" "$UNUSED_MOUNT"
    sudo mount ${LOOP_DEVICE}p1 "$BOOT_MOUNT"
    sudo mount ${LOOP_DEVICE}p2 "$PRIMARY_MOUNT"
    sudo mount ${LOOP_DEVICE}p3 "$PASSIVE_MOUNT"
    sudo mount ${LOOP_DEVICE}p4 "$UNUSED_MOUNT"
    log "SUCCESS" "Partitions mounted."
}

function setup_rootfs() {
    log "INFO" "Setting up root filesystem..."
    sudo rsync -aAX "$TMP_ROOTFS/" "$PRIMARY_MOUNT/"
    echo "127.0.0.1 localhost" | sudo tee "$PRIMARY_MOUNT/etc/hosts"
    echo "nameserver 8.8.8.8" | sudo tee "$PRIMARY_MOUNT/etc/resolv.conf"
    sudo chroot "$PRIMARY_MOUNT" /bin/sh <<EOF
apk update
apk add busybox dropbear dhcpcd openrc e2fsprogs
rc-update add networking boot
EOF
    sudo rsync -aAX "$PRIMARY_MOUNT/" "$PASSIVE_MOUNT/"
    log "SUCCESS" "Root filesystem setup complete."
}

function setup_fstab() {
    log "INFO" "Configuring fstab..."
    echo "/dev/mmcblk0p1  /boot  vfat defaults 0 1" | sudo tee "$PRIMARY_MOUNT/etc/fstab" > /dev/null
    echo "/dev/mmcblk0p2  /      ext4 defaults 0 1" | sudo tee -a "$PRIMARY_MOUNT/etc/fstab" > /dev/null
    echo "/dev/mmcblk0p1  /boot  vfat defaults 0 1" | sudo tee "$PASSIVE_MOUNT/etc/fstab" > /dev/null
    echo "/dev/mmcblk0p3  /      ext4 defaults 0 1" | sudo tee -a "$PASSIVE_MOUNT/etc/fstab" > /dev/null
    log "SUCCESS" "fstab configured."
}

function unmount_partitions() {
    log "INFO" "Unmounting partitions..."
    sudo umount -R "$BOOT_MOUNT" "$PRIMARY_MOUNT" "$PASSIVE_MOUNT" "$UNUSED_MOUNT"
    log "SUCCESS" "Partitions unmounted."
}

function pre_cleanup_and_dir_setup() {

    local image_name=$1
    local tmp_dir=$2
    local build_dir=$3

    if [ -f "$image_name" ]; then
        rm "$image_name"
    fi

    if [ -d "$tmp_dir" ]; then
        rm -rf "$tmp_dir"
    fi
    mkdir -p "$tmp_dir"

    if [ -d "$build_dir" ]; then
        rm -rf "$build_dir"
    fi
    mkdir "$build_dir"
}

# Main Script Execution
check_requirements
pre_cleanup_and_dir_setup "$IMG_NAME" "$TMP_DIR" "${CWD}/build_access_node"

cd ${TMP_DIR}

build_linux_kernel
download_alpine_rootfs

create_disk_image
partition_disk_image
format_partitions
mount_partitions

copy_linux_kernel
copy_rootfs
setup_rootfs
setup_fstab
unmount_partitions
cleanup

cd ${CWD}
log "SUCCESS" "Access node image built successfully: $IMG_NAME"
