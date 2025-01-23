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

function check_sudo() {
    if ! sudo -v; then
        echo "You do not have sudo privileges or sudo is not configured correctly."
        exit 1
    fi
}

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

    wget https://cdn.kernel.org/pub/linux/kernel/v6.x/linux-6.1.34.tar.xz
    tar xvJf linux-6.1.34.tar.xz
    mv linux-6.1.34 "${TMP_LINUX}"

    cd "${TMP_LINUX}"

    # build linux kernel suitable for qemu
    ARCH=arm64 CROSS_COMPILE=/bin/aarch64-linux-gnu- make defconfig        || true
    ARCH=arm64 CROSS_COMPILE=/bin/aarch64-linux-gnu- make kvm_guest.config || true
    ARCH=arm64 CROSS_COMPILE=/bin/aarch64-linux-gnu- make -j8              || true

    cd "$TMP_DIR"
    log "INFO" "Linux kernel build completed."
}

function build_apps_using_container() {

    local ukama_root="$1"
    local apps="$2"

    log "INFO" "Packaging applications via container build ..."
    cwd=$(pwd)

    cd "${ukama_root}/builder/docker"
    ./apps_setup.sh "access" "alpine" "${ukama_root}" "${apps}"
    cd ${cwd}
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

function download_rpi_rootfs() {

    log "INFO" "Downloading rpi image..."
    
    wget https://downloads.raspberrypi.com/raspios_arm64/images/raspios_arm64-2024-11-19/2024-11-19-raspios-bookworm-arm64.img.xz || true
    xz -d 2024-11-19-raspios-bookworm-arm64.img.xz || true

    log "SUCCESS" "Rpi rootfs image downloaded."
}

function install_starter_app() {

    path=$1

    log "INFO" "Installing starter.d"

    sudo chroot "$path" /bin/sh <<'EOF'
cd /ukama/apps/pkgs/
tar zxvf starterd_latest.tar.gz starterd_latest/sbin/starter.d .
mv starterd_latest/sbin/starter.d /sbin/
rm -rf starterd_latest/
EOF
}

function copy_rootfs() {
    log "INFO" "Extracting Alpine Linux root filesystem..."
    cd "$TMP_DIR"
    sudo mkdir -p "$TMP_ROOTFS"
    sudo tar -xzf "$ALPINE_TAR" -C "$TMP_ROOTFS"
    log "SUCCESS" "Alpine Linux extracted."
}

function copy_linux_kernel() {
    log "INFO" "Copying linux kernel..."
    cp "${TMP_LINUX}/arch/arm64/boot/Image" "${CWD}/build_access_node/" || true
    log "SUCCESS" "Linux kernel copied"
}

function copy_all_apps() {
    local ukama_root=$1
    local apps=$2

    log "INFO" "Copying apps"

    sudo mkdir -p "${PRIMARY_MOUNT}/ukama/apps/pkgs"
    sudo mkdir -p "${PASSIVE_MOUNT}/ukama/apps/pkgs"

    IFS=',' read -r -a array <<< "$apps"
    for app in "${array[@]}"; do
        sudo cp "${ukama_root}/build/pkgs/${app}_latest.tar.gz" \
             "${PRIMARY_MOUNT}/ukama/apps/pkgs"
        sudo cp "${ukama_root}/build/pkgs/${app}_latest.tar.gz" \
             "${PASSIVE_MOUNT}/ukama/apps/pkgs"
    done

    # cleanup
    sudo rm -rf "${ukama_root}/build/"
}

function copy_misc_files() {
    local ukama_root=$1
    local apps=$2

    log "INFO" "Copying various files to image"

    create_manifest_file $apps
    sudo cp ${MANIFEST_FILE} "${PRIMARY_MOUNT}/manifest.json"
    sudo cp ${MANIFEST_FILE} "${PASSIVE_MOUNT}/manifest.json"
    rm ${MANIFEST_FILE}

    # install the starter.d app
    install_starter_app "${PRIMARY_MOUNT}"
    install_starter_app "${PASSIVE_MOUNT}"

    echo "Copy Ukama sys lib to the OS image"
    sudo mkdir -p "${PRIMARY_MOUNT}/lib/x86_64-linux-gnu/"
    sudo mkdir -p "${PASSIVE_MOUNT}/lib/x86_64-linux-gnu/"

    sudo cp "${ukama_root}/nodes/ukamaOS/distro/platform/build/libusys.so" \
         "${PRIMARY_MOUNT}/lib/x86_64-linux-gnu/"
    sudo cp "${ukama_root}/nodes/ukamaOS/distro/platform/build/libusys.so" \
         "${PASSIVE_MOUNT}/lib/x86_64-linux-gnu/"

    # update /etc/services to add ports
    echo "Adding all the apps to /etc/services"
    sudo mkdir -p "${PRIMARY_MOUNT}/etc"
    sudo mkdir -p "${PASSIVE_MOOUNT}/etc"

    sudo cp "${ukama_root}/nodes/ukamaOS/distro/scripts/files/services" \
         "${PRIMARY_MOUNT}/etc/services"
    sudo cp "${ukama_root}/nodes/ukamaOS/distro/scripts/files/services" \
         "${PASSIVE_MOUNT}/etc/services"
}

function create_manifest_file() {
    local apps_to_include="$1"
    log "INFO" "Creating manifest file"

    # Create an array from the comma-separated list
    IFS=',' read -r -a apps_array <<< "$apps_to_include"

   cat <<EOF > ${MANIFEST_FILE}
{
    "version": "0.1",

    "spaces" : [
        { "name" : "boot" },
        { "name" : "services" },
        { "name" : "reboot" }
    ],

    "capps": [
        {
            "name"   : "noded",
            "tag"    : "latest",
            "restart": "yes",
            "space"  : "boot"
        },
        {
            "name"   : "bootstrap",
            "tag"    : "latest",
            "restart": "yes",
            "space"  : "boot",
                "depends_on" : [
                {
                    "capp"  : "noded",
                                "state" : "active"
                        }
                ]
        },
        {
            "name"   : "meshd",
            "tag"    : "latest",
            "restart": "yes",
            "space"  : "boot",
                "depends_on" : [
                {
                    "capp"  : "bootstrap",
                                "state" : "done"
                        }
                ]
        }
EOF

  echo '        ,' >> ${MANIFEST_FILE}
  echo '        {"name" : "services", "capps" : [' >> ${MANIFEST_FILE}

  for app in "${apps_array[@]}"; do
    case "$app" in
      "wimcd"|"configd"|"metricsd"|"lookoutd"|"deviced"|"notifyd")
        cat <<EOF >> ${MANIFEST_FILE}
        {
            "name"   : "$app",
            "tag"    : "latest",
            "restart": "yes",
            "space"  : "services"
        },
EOF
        ;;
    esac
  done

  echo '        ,' >> ${MANIFEST_FILE}
  echo '        {"name" : "services", "capps" : [' >> ${MANIFEST_FILE}

  for app in "${apps_array[@]}"; do
    case "$app" in
      "wimcd"|"configd"|"metricsd"|"lookoutd"|"deviced"|"notifyd")
        cat <<EOF >> ${MANIFEST_FILE}
        {
            "name"   : "$app",
            "tag"    : "latest",
            "restart": "yes",
            "space"  : "services"
        },
EOF
        ;;
    esac
  done

  # Remove the last comma and close the JSON array
  sed -i '$ s/,$//' ${MANIFEST_FILE}
  echo '    ]}'  >> ${MANIFEST_FILE}
  echo '}'       >> ${MANIFEST_FILE}
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

function configure_openrc_service() {
    local rootfs_path=$1

    log "INFO" "Configuring OpenRC service for starter.d"

    sudo chroot "$rootfs_path" /bin/sh <<'EOF'
# Create OpenRC service script
cat <<'SERVICE' > /etc/init.d/starter
#!/sbin/openrc-run
description="Starter service for running starter.d"

command="/sbin/starter.d"
command_args=""
command_user="root"
pidfile="/var/run/starter.pid"
SERVICE

# Make the service script executable
chmod +x /etc/init.d/starter

# Add the service to the default runlevel
rc-update add starter default
EOF

    log "SUCCESS" "OpenRC service for starter.d configured"
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
UKAMA_ROOT=$1
NODE_APPS=$2

OS_TYPE="alpine"
OS_VERSION="0.0.1"
MANIFEST_FILE="manifest.json"
export TARGET="access"

check_sudo
check_requirements
pre_cleanup_and_dir_setup "$IMG_NAME" "$TMP_DIR" "${CWD}/build_access_node"

cd ${TMP_DIR}

build_linux_kernel
download_rpi_rootfs

create_disk_image
partition_disk_image
format_partitions
mount_partitions

copy_linux_kernel
copy_rootfs
setup_rootfs
setup_fstab

build_apps_using_container "${UKAMA_ROOT}" "${NODE_APPS}"
copy_all_apps              "${UKAMA_ROOT}" "${NODE_APPS}"
copy_misc_files            "${UKAMA_ROOT}" "${NODE_APPS}"

configure_openrc_service "${PRIMARY_MOUNT}"
configure_openrc_service "${PASSIVE_MOUNT}"

cp "${TMP_DIR}/${IMG_NAME}" ${CWD}

#unmount_partitions
#cleanup

cd ${CWD}
log "SUCCESS" "Access node image built successfully: $IMG_NAME"
