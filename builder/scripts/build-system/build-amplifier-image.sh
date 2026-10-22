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
APP_CONFIGS_DIR="${UKAMA_ROOT}/nodes/configs/apps"
BOOTSTRAP_SERVER="${BOOTSTRAP_SERVER:-dev.bootstrap.ukama.com}"
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
    log INFO "Creating sparse image ${RAW_IMG} (28 GiB)"
    rm -f "$RAW_IMG"
    sudo rm -f "${DIR}/boot.tgz" "${DIR}/rootfsA.tgz" "${DIR}/rootfsB.tgz"
    truncate -s 28674MiB "$RAW_IMG"
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
    log INFO "Partitioning $LOOPDEV: boot, recovery, primary, passive, data, swap"
    sudo parted -s "$LOOPDEV" mklabel msdos \
         mkpart primary fat32      1MiB     1025MiB  \
         mkpart primary ext4       1025MiB  5121MiB  \
         mkpart extended           5121MiB  28673MiB \
         mkpart logical ext4       5122MiB  11265MiB \
         mkpart logical ext4       11266MiB 17409MiB \
         mkpart logical ext4       17410MiB 27649MiB \
         mkpart logical linux-swap 27650MiB 28673MiB \
         set 1 boot on
    check_status $? "Partitions created" "partition_image"
}

map_partitions() {
    log INFO "Mapping partitions for $LOOPDEV"
    sudo kpartx -av "$LOOPDEV" | tee /dev/stderr
    sleep 1
    check_status $? "Partitions mapped" "map_partitions"
}

format_partitions() {
    log INFO "Formatting partitions"
    STAGE="format_partitions"
    local part="/dev/mapper/$(basename "$LOOPDEV")p"

    sudo mkfs.vfat -F32 -n boot "${part}1"
    check_status $? "boot formatted" "${STAGE}"

    for p in 2:recovery 5:primary 6:passive 7:data; do
        sudo mkfs.ext4 -L "${p#*:}" \
             -O ^64bit,^metadata_csum \
             "${part}${p%%:*}"
        check_status $? "${p#*:} formatted" "${STAGE}"
    done

    sudo mkswap -L swap "${part}8"
    check_status $? "swap created" "${STAGE}"
}

mount_partition() {
    local part=$1 mp=$2
    log INFO "Mounting $part → $mp"
    sudo mkdir -p "$mp"
    sudo mount "$part" "$mp"
    check_status $? "Mounted $part" "mount_partition"
}

write_uboot_env() {
    local dst="$1"
    local work env_txt
    STAGE="write_uboot_env"

    work=$(mktemp -d)
    env_txt="${work}/uboot-env.txt"
    (cd "$work" && CROSS_COMPILE="${CROSS_COMPILE}" \
        "${FIRMWARE_DIR}/u-boot/scripts/get_default_envs.sh" "${FIRMWARE_DIR}/u-boot") |
        awk -F= '!($1 in v) { k[++n] = $1 } { v[$1] = $0 } END { for (i = 1; i <= n; i++) print v[k[i]] }' \
        > "$env_txt"

    if ! grep -q '^bootcmd=run emmc_boot;' "$env_txt"; then
        log ERROR "u-boot default env has no A/B bootcmd (run emmc_boot)"
        rm -rf "$work"
        exit 1
    fi

    sudo "${FIRMWARE_DIR}/u-boot/tools/mkenvimage" -s 0x4000 -o "${dst}/uboot.env" "$env_txt"
    check_status $? "uboot.env written with $(grep '^bootcmd=' "$env_txt")" "${STAGE}"
    rm -rf "$work"
}

copy_boot() {
    STAGE="copy_boot"
    log INFO "Copying eMMC firmware and uboot.env → ${BOOT_MNT}"

    local dst="${BOOT_MNT}"

    # 1) Clean out anything already on the boot partition
    sudo rm -rf "${dst:?}/"*
    check_status $? "Cleared old files on boot partition" "${STAGE}"

    sudo cp "${BOOT1_BIN}" "${dst}/BOOT.BIN"
    sudo cp "${BOOT2_BIN}" "${dst}/u-boot.bin"
    check_status $? "Copied BOOT.BIN and u-boot.bin" "${STAGE}"

    write_uboot_env "${dst}"

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

install_kernel() {
    local rootfs="$1"
    sudo cp "${KERNEL_OUT}/arch/arm/boot/zImage" "${rootfs}/boot/zImage"
    sudo cp "${KERNEL_OUT}/arch/arm/boot/dts/microchip/ukama_anode.dtb" "${rootfs}/boot/ukama_anode.dtb"
    sudo mkdir -p "${rootfs}/lib/modules"
    sudo cp -a "${KERNEL_OUT}/modroot/lib/modules/." "${rootfs}/lib/modules/"
    check_status $? "zImage, ukama_anode.dtb and modules installed in ${rootfs}" "install_kernel"
}

update_fstab() {
    local rootfs_dir="$1"
    local root_part="$2"
    local fstab="${rootfs_dir}/etc/fstab"
    STAGE="update_fstab"

    log INFO "Writing fstab with root on mmcblk0p${root_part} in ${fstab}"

    sudo mkdir -p "${rootfs_dir}/media/boot" "${rootfs_dir}/data"
    sudo tee "${fstab}" > /dev/null <<FSTAB
# pseudo-file systems
devtmpfs        /dev        devtmpfs    defaults    0 0
proc            /proc       proc        defaults    0 0
sysfs           /sys        sysfs       defaults    0 0
devpts          /dev/pts    devpts      defaults    0 0
tmpfs           /tmp        tmpfs       defaults    0 0

# eMMC partitions
/dev/mmcblk0p${root_part}  /            ext4   defaults           0 1
/dev/mmcblk0p1  /media/boot  vfat   defaults,nofail    0 0
/dev/mmcblk0p7  /data        ext4   defaults,nofail    0 0
/dev/mmcblk0p8  none         swap   sw,nofail          0 0
FSTAB

    check_status $? "fstab written" "${STAGE}"
}

setup_network() {
    local rootfs="$1"
    sudo tee "${rootfs}/etc/dhcpcd.conf" > /dev/null <<EOF
hostname
clientid
persistent
option domain_name_servers, domain_name, domain_search, host_name
option classless_static_routes
option interface_mtu
require dhcp_server_identifier
EOF
    sudo ln -sfn /etc/init.d/loopback "${rootfs}/etc/runlevels/boot/loopback"
    sudo ln -sfn /etc/init.d/ntpd "${rootfs}/etc/runlevels/default/ntpd"
    check_status $? "DHCP, loopback and ntpd in ${rootfs}" "setup_network"
}

setup_serial_console() {
    local rootfs="$1"
    grep -q '^ttyS0:' "${rootfs}/etc/inittab" || \
        echo 'ttyS0::respawn:/sbin/getty -L 115200 ttyS0 vt100' | sudo tee -a "${rootfs}/etc/inittab" > /dev/null
    check_status $? "Serial login on ttyS0 in ${rootfs}" "setup_serial_console"
}

setup_noded_sysfs() {
    local rootfs="$1"
    sudo tee "${rootfs}/etc/init.d/ukama-sysfs" > /dev/null <<'EOF'
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
    sudo chmod 0755 "${rootfs}/etc/init.d/ukama-sysfs"
    sudo ln -sfn /etc/init.d/ukama-sysfs "${rootfs}/etc/runlevels/boot/ukama-sysfs"
    sudo mkdir -p "${rootfs}/ukama/mocksysfs"
    sudo ln -sfn /sys "${rootfs}/ukama/mocksysfs/sys"
    check_status $? "ukama-sysfs boot service in ${rootfs}" "setup_noded_sysfs"
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

    sudo mkdir -p "${rootfs}/ukama/configs/metricsd" \
                  "${rootfs}/ukama/configs/aggregator"
    for app in bootstrap meshd notifyd femd; do
        sudo cp -r "${APP_CONFIGS_DIR}/${app}" "${rootfs}/ukama/configs/"
    done
    sudo cp "${APP_CONFIGS_DIR}/metricsd/amplifier_config.toml" \
            "${rootfs}/ukama/configs/metricsd/config.toml"
    sudo cp -r "${APP_CONFIGS_DIR}/metricsd/anode" "${rootfs}/ukama/configs/metricsd/"
    sudo cp "${APP_CONFIGS_DIR}/aggregator/anode_config.toml" \
            "${rootfs}/ukama/configs/aggregator/config.toml"
    check_status $? "App configs copied to ${rootfs}" "deploy_to_rootfs"

    echo "${BOOTSTRAP_SERVER}" | sudo tee "${rootfs}/ukama/bootstrap" > /dev/null
    check_status $? "Bootstrap server ${BOOTSTRAP_SERVER} set in ${rootfs}" "deploy_to_rootfs"
}

pack_tarballs() {
    STAGE="pack_tarballs"
    log INFO "Packing boot.tgz, rootfsA.tgz and rootfsB.tgz into ${DIR}"
    sync

    sudo tar -czf "${DIR}/boot.tgz" -C "${BOOT_MNT}" .
    check_status $? "boot.tgz packed" "${STAGE}"

    sudo tar --xattrs --xattrs-include='*' --acls -czf "${DIR}/rootfsA.tgz" -C "${PRIMARY_MNT}" .
    check_status $? "rootfsA.tgz packed" "${STAGE}"

    sudo tar --xattrs --xattrs-include='*' --acls -czf "${DIR}/rootfsB.tgz" -C "${PASSIVE_MNT}" .
    check_status $? "rootfsB.tgz packed" "${STAGE}"
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
mount_partition "/dev/mapper/$(basename ${LOOPDEV})p5" "$PRIMARY_MNT"
mount_partition "/dev/mapper/$(basename ${LOOPDEV})p6" "$PASSIVE_MNT"

# Signal success; cleanup trap will unmount and detach
log SUCCESS "Disk image ${RAW_IMG} prepared successfully"

copy_boot
copy_rootfs
install_kernel "${PRIMARY_MNT}"
install_kernel "${PASSIVE_MNT}"
update_fstab "${PRIMARY_MNT}" 5
update_fstab "${PASSIVE_MNT}" 6
setup_serial_console "${PRIMARY_MNT}"
setup_serial_console "${PASSIVE_MNT}"
setup_network "${PRIMARY_MNT}"
setup_network "${PASSIVE_MNT}"
setup_noded_sysfs "${PRIMARY_MNT}"
setup_noded_sysfs "${PASSIVE_MNT}"

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
pack_tarballs

log "SUCCESS" "Disk image creation completed successfully! (bootstrap: ${BOOTSTRAP_SERVER})"
exit 0
