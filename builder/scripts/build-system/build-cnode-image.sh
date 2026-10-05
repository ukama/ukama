#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

# Build image for Ukama controller node (cnode, CM4/CM5).

set -e

PKG_UTILS="$(dirname "$0")/pkg-utils.sh"
if [ ! -f "$PKG_UTILS" ]; then
    echo "ERROR: Missing $PKG_UTILS"
    exit 1
fi
source "$PKG_UTILS"

STAGE="init"
UKAMA_ROOT=$(realpath ../../../)
UKAMA_REPO_APP_PKG="${UKAMA_ROOT}/build/pkgs"
UKAMA_REPO_LIB_PKG="${UKAMA_ROOT}/build/libs"
APP_CONFIGS_DIR="${UKAMA_ROOT}/nodes/configs/apps"
COMMON_CONFIG_FILE="${UKAMA_ROOT}/builder/boards/common.config"
CNODE_CONFIG_FILE="${UKAMA_ROOT}/builder/boards/controller.config"

BOOT_MOUNT="/media/boot"
PRIMARY_MOUNT="/media/primary"
PASSIVE_MOUNT="/media/passive"

RAW_IMG="cnode-image.img"
IMG_SIZE="28800M"

ROOTFS_DIR=${UKAMA_ROOT}/builder/scripts/build-system/rootfs
FIRMWARE_ZIP_URL="https://github.com/raspberrypi/firmware/archive/refs/heads/master.zip"
FIRMWARE_DIR="/tmp/rpi-firmware"
BOOTSTRAP_SERVER="${BOOTSTRAP_SERVER:-dev.bootstrap.ukama.com}"
LOOPDISK=""

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
    for mount in ${BOOT_MOUNT} ${PRIMARY_MOUNT} ${PASSIVE_MOUNT}; do
        sudo umount "${mount}" 2>/dev/null || true
    done
    sudo kpartx -d "${LOOPDISK}" 2>/dev/null || true
    sudo losetup -d "${LOOPDISK}" 2>/dev/null || true
    LOOPDISK=""
}
trap cleanup EXIT

check_prerequisites() {
    STAGE="check_prerequisites"
    for cmd in sfdisk kpartx losetup mkfs.vfat mkfs.ext4 mkswap e2label rsync wget unzip truncate; do
        command -v "$cmd" >/dev/null || { log "ERROR" "'$cmd' not found"; exit 1; }
    done
    [ -x "${ROOTFS_DIR}/sbin/starter.d" ] || {
        log "ERROR" "No rootfs at ${ROOTFS_DIR}. Run rootfs-env-setup.sh -a aarch64 first"
        exit 1
    }
    [ -d "${UKAMA_REPO_APP_PKG}" ] && [ -f "${UKAMA_REPO_LIB_PKG}/vendor_libs.tgz" ] || {
        log "ERROR" "No apps in ${UKAMA_ROOT}/build. Run build-env-setup.sh -a aarch64 first"
        exit 1
    }
    log "SUCCESS" "Prerequisites OK"
}

download_firmware() {
    STAGE="download_firmware"
    if [ -d "${FIRMWARE_DIR}/firmware-master/boot" ]; then
        log "INFO" "Using cached firmware in ${FIRMWARE_DIR}"
        return
    fi
    log "INFO" "Downloading Raspberry Pi firmware (boot files, kernel, modules)"
    rm -rf "${FIRMWARE_DIR}"
    mkdir -p "${FIRMWARE_DIR}"
    wget -qO "${FIRMWARE_DIR}/master.zip" "${FIRMWARE_ZIP_URL}"
    unzip -q "${FIRMWARE_DIR}/master.zip" -d "${FIRMWARE_DIR}"
    rm -f "${FIRMWARE_DIR}/master.zip"
    [ -f "${FIRMWARE_DIR}/firmware-master/boot/kernel8.img" ]
    check_status $? "Firmware downloaded" ${STAGE}
}

create_disk_image() {
    STAGE="create_disk_image"
    log "INFO" "Creating raw image ${RAW_IMG} (${IMG_SIZE})"
    rm -f "${RAW_IMG}"
    truncate -s "${IMG_SIZE}" "${RAW_IMG}"
    check_status $? "Raw image created" ${STAGE}
}

setup_loop_device() {
    STAGE="setup_loop_device"
    LOOPDISK=$(sudo losetup -f --show "${RAW_IMG}")
    [ -n "${LOOPDISK}" ]
    check_status $? "Loop device set up at ${LOOPDISK}" ${STAGE}
}

partition_image() {
    STAGE="partition_image"
    log "INFO" "Creating partitions on ${LOOPDISK}"
    sudo sfdisk "${LOOPDISK}" <<-__EOF__
label: dos
,1G,c
,4G,83
,,5
,6G,83
,6G,83
,10G,83
,1G,82
__EOF__
    check_status $? "Partitions created" ${STAGE}
}

map_partitions() {
    STAGE="map_partitions"
    sudo kpartx -v -a "${LOOPDISK}"
    check_status $? "Partitions mapped" ${STAGE}
    DISK="/dev/mapper/$(basename "${LOOPDISK}")p"
}

format_partitions() {
    STAGE="format_partitions"
    sudo mkfs.vfat -F 32 -n boot "${DISK}1"
    sudo mkfs.ext4 -F -L recovery "${DISK}2"
    sudo mkfs.ext4 -F -L primary  "${DISK}5"
    sudo mkfs.ext4 -F -L passive  "${DISK}6"
    sudo mkfs.ext4 -F -L data     "${DISK}7"
    sudo mkswap -L swap "${DISK}8"
    check_status $? "Partitions formatted" ${STAGE}

    [ "$(sudo e2label "${DISK}5")" = "primary" ] && \
        [ "$(sudo e2label "${DISK}6")" = "passive" ]
    check_status $? "Partition labels confirmed" ${STAGE}
}

mount_partition() {
    sudo mkdir -p "$2"
    sudo mount "$1" "$2"
    check_status $? "Mounted $1 on $2" "mount_partition"
}

copy_boot() {
    STAGE="copy_boot"
    local fw="${FIRMWARE_DIR}/firmware-master/boot"

    log "INFO" "Copying firmware boot files to boot partition"
    sudo cp -r "${fw}/." "${BOOT_MOUNT}/"
    sudo cp "${BOOT_MOUNT}/kernel8.img" "${BOOT_MOUNT}/kernel.img"

    sudo tee "${BOOT_MOUNT}/config.txt" > /dev/null <<EOF
enable_uart=1
arm_64bit=1
kernel=kernel.img
gpu_mem=64
boot_delay=1
disable_splash=1
dtoverlay=disable-bt
dtoverlay=disable-wifi
EOF

    # no initramfs, so root by device (p5 = primary)
    sudo tee "${BOOT_MOUNT}/cmdline.txt" > /dev/null <<EOF
console=serial0,115200 console=tty1 root=/dev/mmcblk0p5 rootfstype=ext4 fsck.repair=yes rootwait
EOF
    check_status $? "Boot partition ready" ${STAGE}
}

copy_rootfs() {
    STAGE="copy_rootfs"
    local target

    for target in "${PRIMARY_MOUNT}" "${PASSIVE_MOUNT}"; do
        log "INFO" "Copying rootfs to ${target}"
        sudo rsync -aAX \
             --exclude=/dev/* --exclude=/proc/* --exclude=/sys/* \
             --exclude=/run/* --exclude=/tmp/* --exclude=/ukamarepo \
             --exclude=/destroy --exclude=/enter-chroot --exclude=/env.sh \
             --exclude=/setup.log \
             "${ROOTFS_DIR}/" "${target}/"
        sudo mkdir -p "${target}/dev" "${target}/proc" "${target}/sys" \
                      "${target}/run" "${target}/tmp" "${target}/boot/firmware"
        sudo chmod 1777 "${target}/tmp"
    done
    check_status $? "Rootfs copied to primary and passive" ${STAGE}
}

copy_kernel_modules() {
    STAGE="copy_kernel_modules"
    local src="${FIRMWARE_DIR}/firmware-master/modules"
    local target dir found=0

    for target in "${PRIMARY_MOUNT}" "${PASSIVE_MOUNT}"; do
        sudo mkdir -p "${target}/lib/modules"
        for dir in "${src}"/*-v8+; do
            [ -d "$dir" ] || continue
            sudo cp -a "$dir" "${target}/lib/modules/"
            found=1
        done
    done
    [ "$found" -eq 1 ]
    check_status $? "Kernel modules (kernel8) copied" ${STAGE}
}

update_fstab() {
    local target="$1"
    local root_label="$2"

    sudo tee "${target}/etc/fstab" > /dev/null <<FSTAB
proc                 /proc           proc    defaults              0 0
sysfs                /sys            sysfs   defaults              0 0
devpts               /dev/pts        devpts  defaults              0 0
tmpfs                /tmp            tmpfs   defaults              0 0
LABEL=${root_label}  /               ext4    errors=remount-ro     0 1
LABEL=boot           /boot/firmware  vfat    ro                    0 2
FSTAB
}

setup_network() {
    local target="$1"

    # DHCP on every interface instead of the fixed IP from build-rootfs.sh
    sudo tee "${target}/etc/dhcpcd.conf" > /dev/null <<EOF
hostname
clientid
persistent
option domain_name_servers, domain_name, domain_search, host_name
option classless_static_routes
option interface_mtu
require dhcp_server_identifier
EOF
}

setup_noded_sysfs() {
    local target="$1"

    # noded reads /tmp/sys, and /tmp is a tmpfs at boot
    sudo tee "${target}/etc/init.d/ukama-sysfs" > /dev/null <<'EOF'
#!/sbin/openrc-run

description="Link noded sysfs data into /tmp/sys"

depend() {
    need localmount
    after bootmisc
    before starterd
}

start() {
    ebegin "Linking /tmp/sys to /ukama/mocksysfs/sys"
    if [ -d /ukama/mocksysfs/sys ]; then
        rm -rf /tmp/sys
        ln -sfn /ukama/mocksysfs/sys /tmp/sys
    fi
    eend 0
}
EOF
    sudo chmod 0755 "${target}/etc/init.d/ukama-sysfs"
    sudo mkdir -p "${target}/etc/runlevels/boot"
    sudo ln -sfn /etc/init.d/ukama-sysfs "${target}/etc/runlevels/boot/ukama-sysfs"
}

deploy_to_rootfs() {
    STAGE="deploy_to_rootfs"
    local target="$1"

    log "INFO" "Deploying apps, libs, configs and manifest into ${target}"
    sudo mkdir -p "${target}/ukama/apps/pkgs" \
                  "${target}/ukama/configs" \
                  "${target}/ukama/logs" \
                  "${target}/ukama/state/starterd/log-spool" \
                  "${target}/lib"

    copy_all_apps      "$UKAMA_REPO_APP_PKG" "${target}/ukama/apps/pkgs"
    copy_required_libs "$UKAMA_REPO_LIB_PKG" "${target}/lib"
    create_manifest_file "${target}/ukama/manifest.json" "${APPS[@]}"

    # configs needed to come online (bootstrap exits without its config)
    for app in bootstrap meshd; do
        sudo cp -r "${APP_CONFIGS_DIR}/${app}" "${target}/ukama/configs/"
    done

    echo "${BOOTSTRAP_SERVER}" | sudo tee "${target}/ukama/bootstrap" > /dev/null
    check_status $? "Deployed into ${target}" ${STAGE}
}

# Main
check_prerequisites
download_firmware

get_enabled_apps "$COMMON_CONFIG_FILE" "$CNODE_CONFIG_FILE"
if [[ ${#APPS[@]} -eq 0 ]]; then
    log "ERROR" "No apps enabled"
    exit 1
fi
log "INFO" "Apps: ${APPS[*]}"

create_disk_image
setup_loop_device
partition_image
map_partitions
format_partitions

mount_partition "${DISK}1" "${BOOT_MOUNT}"
mount_partition "${DISK}5" "${PRIMARY_MOUNT}"
mount_partition "${DISK}6" "${PASSIVE_MOUNT}"

copy_boot
copy_rootfs
copy_kernel_modules

update_fstab "${PRIMARY_MOUNT}" primary
update_fstab "${PASSIVE_MOUNT}" passive
for target in "${PRIMARY_MOUNT}" "${PASSIVE_MOUNT}"; do
    setup_network "${target}"
    setup_noded_sysfs "${target}"
    deploy_to_rootfs "${target}"
done

sync
cleanup
log "SUCCESS" "Image ${RAW_IMG} created (bootstrap: ${BOOTSTRAP_SERVER})"
exit 0
