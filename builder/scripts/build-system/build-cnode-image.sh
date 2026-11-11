#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

# Build image for Ukama controller node (cnode, CM4/CM5).

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "${SCRIPT_DIR}"

PKG_UTILS="${SCRIPT_DIR}/pkg-utils.sh"
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
CNODE_CONFIG_FILE="${UKAMA_ROOT}/builder/boards/cnode.config"

IMAGE_NAME="cnode-image.img"
IMG_SIZE="28800M"

ROOTFS_DIR=${UKAMA_ROOT}/builder/scripts/build-system/rootfs
BUILDENV_DIR=${UKAMA_ROOT}/builder/scripts/build-system/buildenv
# raspberrypi/firmware release 1.20260915 (CM4 and CM5).
RPI_FIRMWARE_COMMIT="6f0881cba8bea8ec24956a5718714730e9936ec5"
RPI_FIRMWARE_SHA256="848bdbe2870e627a0b35c1de772cb66c24b7cb1069494630f863caf51b2bf7b0"
BOOTSTRAP_SERVER="${BOOTSTRAP_SERVER:-dev.bootstrap.ukama.com}"
LOOPDISK=""
HIDDEN_CONFIGS=""
ALPINE_TMP=""
WORK_DIR=""

if [ $# -ne 0 ]; then
    echo "Usage: $0"
    exit 1
fi

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

restore_configs() {
    local config
    [ -n "$HIDDEN_CONFIGS" ] || return 0
    for config in "$HIDDEN_CONFIGS"/*.toml; do
        [ -f "$config" ] || continue
        mv -- "$config" "${UKAMA_ROOT}/builder/configs/" || return 1
    done
    rmdir "$HIDDEN_CONFIGS" || return 1
    HIDDEN_CONFIGS=""
}

clean_build_outputs() {
    local distro="${UKAMA_ROOT}/nodes/ukamaOS/distro"
    local vendor="${distro}/vendor"
    local recipe failed=0

    log INFO "Cleaning Ukama apps, builder, platform and vendor libraries"
    for recipe in "${UKAMA_ROOT}"/nodes/apps/*/Makefile "$distro"/system/*/Makefile \
                  "${UKAMA_ROOT}/builder/Makefile" "$distro/platform/Makefile" \
                  "$vendor"/{orcania/src,yder/src,ulfius/src,tomlc,libcap}/Makefile \
                  "$vendor/Makefile"; do
        [ -f "$recipe" ] || continue
        make -C "${recipe%/Makefile}" clean || failed=1
    done

    # The vendor clean target removes the installed libs. Interrupted builds
    # can also leave CMake or configure outputs inside individual libraries.
    rm -rf -- "$vendor"/{jansson,amqp,libuuid,casync}/tmpbuild \
              "$vendor"/prometheus-client/{prom,promhttp}/build \
              "$UKAMA_REPO_APP_PKG" || failed=1
    rm -f -- "${UKAMA_REPO_LIB_PKG}/vendor_libs.tgz" || failed=1
    return "$failed"
}

remove_build_dirs() {
    local dir mounts target failed=0
    for dir in "$BUILDENV_DIR" "$ROOTFS_DIR"; do
        [ -d "$dir" ] || continue
        if [ -x "$dir/destroy" ]; then
            "$dir/destroy" --remove || failed=1
        else
            # A partially installed chroot may still have the checkout bind-mounted.
            # --one-file-system alone cannot distinguish a bind on the same filesystem.
            mounts=$(findmnt -rn -o TARGET) || return 1
            while IFS= read -r target; do
                printf -v target '%b' "$target"
                if [[ "$target" == "$dir" || "$target" == "$dir/"* ]]; then
                    log ERROR "Still mounted: $target"
                    failed=1
                    continue 2
                fi
            done <<< "$mounts"
            rm -rf --one-file-system -- "$dir" || failed=1
        fi
    done
    rm -rf -- "${SCRIPT_DIR}/mocksysfs" || failed=1
    return "$failed"
}

release_image() {
    local dir
    [ -n "$LOOPDISK" ] || return 0
    for dir in "$PASSIVE_MOUNT" "$PRIMARY_MOUNT" "$BOOT_MOUNT"; do
        if mountpoint -q "$dir"; then
            sudo umount "$dir" || return 1
        fi
    done
    sudo kpartx -d "$LOOPDISK" || return 1
    sudo losetup -d "$LOOPDISK" || return 1
    LOOPDISK=""
}

cleanup() {
    local status=$1 failed=0 owner dir
    trap - EXIT
    trap '' HUP INT TERM
    restore_configs || failed=1
    release_image || failed=1
    remove_build_dirs || failed=1
    clean_build_outputs || failed=1
    if [ -n "$ALPINE_TMP" ]; then
        rm -rf -- "$ALPINE_TMP" || failed=1
    fi

    owner=$(stat -c '%u:%g' "$UKAMA_ROOT") || failed=1
    # Vendor configure scripts can regenerate source-tree files as root.
    find "${UKAMA_ROOT}/nodes/ukamaOS/distro/vendor" "${UKAMA_ROOT}/nodes/ukamaOS/distro/system" \
         "${UKAMA_ROOT}/nodes/ukamaOS/distro/platform" "${UKAMA_ROOT}/nodes/apps" \
         -xdev -name .git -prune -o -uid 0 -exec chown -h "$owner" {} + || failed=1

    for dir in "${UKAMA_ROOT}/build" "$UKAMA_REPO_LIB_PKG"; do
        [ -d "$dir" ] || continue
        chown --from=0 "$owner" "$dir" || failed=1
    done

    if [ "$status" -eq 0 ] && [ "$failed" -eq 0 ]; then
        chown "$owner" "$RAW_IMG" && mv -fT -- "$RAW_IMG" "$IMAGE_NAME" || failed=1
    fi
    if [ -z "$LOOPDISK" ]; then
        rm -rf -- "$WORK_DIR" || failed=1
    fi

    if [ "$failed" -ne 0 ]; then
        log ERROR "Cleanup or publication of the image failed"
        [ "$status" -ne 0 ] || status=1
    fi
    if [ "$status" -eq 0 ]; then
        log SUCCESS "Image ${IMAGE_NAME} created (bootstrap: ${BOOTSTRAP_SERVER})"
    else
        log ERROR "Script failed at stage: ${STAGE}"
    fi
    exit "$status"
}

check_tools() {
    for cmd in sfdisk kpartx losetup mkfs.vfat mkfs.ext4 mkswap rsync wget unzip sha256sum truncate jq readelf make flock findmnt mountpoint; do
        command -v "$cmd" >/dev/null || { log "ERROR" "'$cmd' not found"; exit 1; }
    done
}

check_prerequisites() {
    STAGE="check_prerequisites"
    [ -x "${ROOTFS_DIR}/sbin/starter.d" ] || {
        log "ERROR" "Rootfs build failed: no starter.d in ${ROOTFS_DIR}"
        exit 1
    }
    # x86 image builds (com, amplifier) use the same rootfs folder
    [ -e "${ROOTFS_DIR}/lib/ld-musl-aarch64.so.1" ] || {
        log "ERROR" "Rootfs build failed: ${ROOTFS_DIR} is not aarch64"
        exit 1
    }
    [ -f "${ROOTFS_DIR}/ukama/mocksysfs/mfgdata/schema/cnode.json" ] || {
        log "ERROR" "Rootfs build failed: no cnode schema in ${ROOTFS_DIR}"
        exit 1
    }
    log "SUCCESS" "Prerequisites OK"
}

build_apps() {
    STAGE="build_apps"
    local config
    HIDDEN_CONFIGS=$(mktemp -d)
    for config in "${UKAMA_ROOT}"/builder/configs/*.toml; do
        if [[ " ${APPS[*]} " != *" $(basename "${config}" .toml) "* ]]; then
            mv "${config}" "${HIDDEN_CONFIGS}/"
        fi
    done
    log "INFO" "Building the apps the cnode uses: ${APPS[*]}"
    TEMP_DIR="${ALPINE_TMP}" ./build-env-setup.sh -a aarch64 -i "${UKAMA_ROOT}"
    restore_configs
    log "SUCCESS" "Apps built"
}

build_rootfs() {
    STAGE="build_rootfs"
    mkdir -p mocksysfs
    [ -n "$(ls -A mocksysfs)" ] || touch mocksysfs/README
    log "INFO" "Building the rootfs"
    TEMP_DIR="${ALPINE_TMP}" ./rootfs-env-setup.sh -a aarch64 -i "${UKAMA_ROOT}" || \
        log "INFO" "rootfs-env-setup.sh always exits 1 on aarch64, checking the rootfs instead"
}

download_firmware() {
    STAGE="download_firmware"
    local zip="${WORK_DIR}/firmware.zip"
    local url="https://codeload.github.com/raspberrypi/firmware/zip/${RPI_FIRMWARE_COMMIT}"

    log INFO "Downloading Raspberry Pi firmware ${RPI_FIRMWARE_COMMIT}"
    wget -O "$zip" "$url"
    printf '%s  %s\n' "$RPI_FIRMWARE_SHA256" "$zip" | sha256sum --check --status || {
        log ERROR "SHA-256 mismatch for Raspberry Pi firmware ${RPI_FIRMWARE_COMMIT}"
        return 1
    }
    unzip -q "$zip" -d "$WORK_DIR"
    rm -f -- "$zip"
    log SUCCESS "Raspberry Pi firmware verified"
}

create_disk_image() {
    STAGE="create_disk_image"
    log "INFO" "Creating raw image ${RAW_IMG} (${IMG_SIZE})"
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
}

mount_partition() {
    sudo mkdir -p "$2"
    sudo mount "$1" "$2"
    check_status $? "Mounted $1 on $2" "mount_partition"
}

copy_boot() {
    STAGE="copy_boot"
    local fw="${FIRMWARE_DIR}/boot"

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
        # -x: never copy anything still mounted inside the chroot
        sudo rsync -aAXx \
             --exclude=/dev/* --exclude=/proc/* --exclude=/sys/* \
             --exclude=/run/* --exclude=/tmp/* --exclude=/ukamarepo \
             --exclude=/destroy --exclude=/enter-chroot --exclude=/env.sh \
             --exclude=/setup.log --exclude=/usr/bin/qemu-*-static \
             "${ROOTFS_DIR}/" "${target}/"
        sudo mkdir -p "${target}/boot/firmware"
    done
    check_status $? "Rootfs copied to primary and passive" ${STAGE}
}

copy_kernel_modules() {
    STAGE="copy_kernel_modules"
    local src="${FIRMWARE_DIR}/modules"
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

    # apps talk over localhost, and lo is down unless OpenRC brings it up
    sudo ln -sfn /etc/init.d/loopback "${target}/etc/runlevels/boot/loopback"

    # the clock starts at 1970, so set it from the network
    sudo ln -sfn /etc/init.d/ntpd "${target}/etc/runlevels/default/ntpd"
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
    sudo ln -sfn /etc/init.d/ukama-sysfs "${target}/etc/runlevels/boot/ukama-sysfs"
    sudo ln -sfn /sys "${target}/ukama/mocksysfs/sys"
}

deploy_to_rootfs() {
    STAGE="deploy_to_rootfs"
    local target="$1"

    log "INFO" "Deploying apps, libs, configs and manifest into ${target}"
    copy_all_apps      "$UKAMA_REPO_APP_PKG" "${target}/ukama/apps/pkgs"
    copy_required_libs "$UKAMA_REPO_LIB_PKG" "${target}/lib"
    create_manifest_file "${target}/ukama/manifest.json" "${APPS[@]}"

    # configs the apps read from /ukama/configs (bootstrap exits without its config)
    for app in bootstrap meshd notifyd switchd; do
        sudo cp -r "${APP_CONFIGS_DIR}/${app}" "${target}/ukama/configs/"
    done

    # metricsd and aggregator read config.toml; the repo has one per node type
    for app in metricsd aggregator; do
        sudo mkdir -p "${target}/ukama/configs/${app}"
        sudo cp "${APP_CONFIGS_DIR}/${app}/cnode_config.toml" \
                "${target}/ukama/configs/${app}/config.toml"
    done

    echo "${BOOTSTRAP_SERVER}" | sudo tee "${target}/ukama/bootstrap" > /dev/null
    check_status $? "Deployed into ${target}" ${STAGE}
}

# Main
[ "$EUID" -eq 0 ] || { log ERROR "Run the cnode image build with sudo"; exit 1; }
check_tools

get_enabled_apps "$COMMON_CONFIG_FILE" "$CNODE_CONFIG_FILE"
if [[ ${#APPS[@]} -eq 0 ]]; then
    log "ERROR" "No apps enabled"
    exit 1
fi
log "INFO" "Apps: ${APPS[*]}"

# A second cnode build in this checkout would clean the first one's objects.
exec 9<"$SCRIPT_DIR"
flock -n 9 || { log ERROR "Another cnode build is using this checkout"; exit 1; }
WORK_DIR=$(mktemp -d "${SCRIPT_DIR}/.cnode.XXXXXXXX")
FIRMWARE_DIR="${WORK_DIR}/firmware-${RPI_FIRMWARE_COMMIT}"
RAW_IMG="${WORK_DIR}/${IMAGE_NAME}"
BOOT_MOUNT="${WORK_DIR}/boot"
PRIMARY_MOUNT="${WORK_DIR}/primary"
PASSIVE_MOUNT="${WORK_DIR}/passive"
trap 'cleanup $?' EXIT
trap 'exit 129' HUP
trap 'exit 130' INT
trap 'exit 143' TERM

STAGE="clean_build_outputs"
remove_build_dirs
clean_build_outputs
download_firmware
ALPINE_TMP=$(mktemp -d)
build_apps
build_rootfs

check_prerequisites

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
release_image
STAGE="validate_image"
./validate-cnode-image.sh "$RAW_IMG"
