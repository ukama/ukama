#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

set -e
set -x

# Variables
IMG_NAME="access-node.img"
RPI_BASE_URL="https://downloads.raspberrypi.com"
RPI_URL_PATH="raspios_oldstable_armhf/images/raspios_oldstable_armhf-2024-10-28"
RPI_IMG="2024-10-22-raspios-bullseye-armhf.img"
RPI_IMG_OFFSET="272629760"
PRIMARY_MOUNT="/mnt/access-node"

TMP_DIR="/tmp/access-node"
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

function cleanup() {
    log "INFO" "Cleaning up resources..."
    sudo umount -R "${PRIMARY_MOUNT}" 2>/dev/null || true
#    sudo rm -rf "$TMP_DIR"
    log "INFO" "Cleanup completed."
}

function build_linux_kernel() {
    log "INFO" "Building linux kernel..."

    if [ -d "${TMP_LINUX}" ]; then
        log "INFO" "Using existing linux kernel at: ${TMP_LINUX}"
    else 
        wget https://cdn.kernel.org/pub/linux/kernel/v6.x/linux-6.1.34.tar.xz
        tar xJf linux-6.1.34.tar.xz
        mv linux-6.1.34 "${TMP_LINUX}"
    fi

    cd "${TMP_LINUX}"

    if [ -f "${TMP_LINUX}/arch/arm64/boot/Image" ]; then
        log "INFO" "Kernel image already exists, skipping"
    else
        # build linux kernel suitable for qemu
        ARCH=arm64 CROSS_COMPILE=/bin/aarch64-linux-gnu- make defconfig
        ARCH=arm64 CROSS_COMPILE=/bin/aarch64-linux-gnu- make kvm_guest.config
        ARCH=arm64 CROSS_COMPILE=/bin/aarch64-linux-gnu- make -j8
    fi

    cd "${TMP_DIR}"
    log "SUCCESS" "Linux kernel build completed."
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

function download_rpi_rootfs() {
    log "INFO" "Checking for RPI rootfs image..."

    if [ -f "${RPI_IMG}" ]; then
        log "INFO" "Using existing extracted image: ${RPI_IMG}"
        return
    fi

    if [ -f "${RPI_IMG}.xz" ]; then
        log "INFO" "Using existing compressed image: ${RPI_IMG}.xz"
    else
        log "INFO" "Downloading RPI rootfs image..."
        wget "${RPI_BASE_URL}/${RPI_URL_PATH}/${RPI_IMG}.xz" \
            || { log "ERROR" "Unable to download ${RPI_IMG}.xz"; exit 1; }
    fi

    log "INFO" "Extracting ${RPI_IMG}.xz..."
    xz -d -f "${RPI_IMG}.xz" \
        || { log "ERROR" "Unable to extract ${RPI_IMG}.xz"; exit 1; }

    log "SUCCESS" "RPI rootfs image downloaded and extracted: ${RPI_IMG}"
}

function install_starter_app() {

    path=$1

    log "INFO" "Installing starter.d"

    sudo chroot "$path" /bin/sh <<'EOF'
cd /ukama/apps/pkgs/
tar zxvf starterd_latest.tar.gz
cp starterd_latest/sbin/starter.d /sbin/
rm -rf starterd_latest/
EOF
}

function copy_linux_kernel() {
    log "INFO" "Copying linux kernel..."
    cp "${TMP_LINUX}/arch/arm64/boot/Image" "${CWD}/build_access_node/kernel.img"
    log "SUCCESS" "Linux kernel copied"
}

function copy_all_apps() {
    local ukama_root=$1
    local apps=$2

    log "INFO" "Copying apps"

    sudo mkdir -p "${PRIMARY_MOUNT}/ukama/apps/pkgs"
    IFS=',' read -r -a array <<< "$apps"
    for app in "${array[@]}"; do
        sudo cp "${ukama_root}/build/pkgs/${app}_latest.tar.gz" \
             "${PRIMARY_MOUNT}/ukama/apps/pkgs"
    done

    sudo rm -rf "${ukama_root}/build/"
}

function copy_misc_files() {
    local ukama_root=$1
    local apps=$2

    log "INFO" "Copying various files to image"

    create_manifest_file $apps
    sudo cp ${MANIFEST_FILE} "${PRIMARY_MOUNT}/manifest.json"
    rm ${MANIFEST_FILE}

    # install the starter.d app
    install_starter_app "${PRIMARY_MOUNT}"

    log "INFO" "Copy Ukama sys lib to the image"
    sudo mkdir -p "${PRIMARY_MOUNT}/lib/x86_64-linux-gnu/"
    sudo cp "${ukama_root}/nodes/ukamaOS/distro/platform/build/libusys.so" \
         "${PRIMARY_MOUNT}/lib/x86_64-linux-gnu/"

    # update /etc/services to add ports
    log "INFO" "Adding all the apps to /etc/services"
    sudo mkdir -p "${PRIMARY_MOUNT}/etc"
    sudo cp "${ukama_root}/nodes/ukamaOS/distro/scripts/files/services" \
         "${PRIMARY_MOUNT}/etc/services"
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

function mount_partitions() {
    log "INFO" "Mounting partitions..."

    sudo mkdir -p "$PRIMARY_MOUNT"
    sudo mount -o loop,offset="${RPI_IMG_OFFSET}" "${RPI_IMG}" "${PRIMARY_MOUNT}" || \
        { log "ERROR" "Unable to mount rootfs partition"; exit 1; }
       
    log "SUCCESS" "Partitions mounted."
}

function setup_ukama_dirs() {
    log "INFO" "Creating Ukama directories..."

    mkdir -p "${PRIMARY_MOUNT}/ukama"
    mkdir -p "${PRIMARY_MOUNT}/ukama/configs"
    mkdir -p "${PRIMARY_MOUNT}/ukama/apps"
    mkdir -p "${PRIMARY_MOUNT}/ukama/apps/pkgs"
    mkdir -p "${PRIMARY_MOUNT}/ukama/apps/rootfs"
    mkdir -p "${PRIMARY_MOUNT}/ukama/apps/registry"

    echo "${NODE_ID}" > "${PRIMARY_MOUNT}/ukama/nodeid"
    echo "localhost"  > "${PRIMARY_MOUNT}/ukama/bootstrap"

    touch "${PRIMARY_MOUNT}/ukama/apps.log"

    log "SUCCESS" "Ukama directories created."
}

function configure_openrc_service() {
    local rootfs_path=$1

    log "INFO" "Setting up minimal Debian rootfs with debootstrap"

    # Step 1: Install debootstrap if not already installed
    if ! command -v debootstrap &>/dev/null; then
        log "INFO" "Installing debootstrap"
        sudo apt update && sudo apt install debootstrap -y
    fi

    # Step 2: Bootstrap a minimal Debian system
    if [ ! -d "$rootfs_path" ]; then
        log "INFO" "Bootstrapping Debian system at $rootfs_path"
        sudo debootstrap --variant=minbase stable "$rootfs_path" http://deb.debian.org/debian/
        log "SUCCESS" "Minimal Debian system installed at $rootfs_path"
    else
        log "INFO" "Rootfs path already exists at $rootfs_path, skipping debootstrap"
    fi

    # Step 3: Bind system directories for chroot
    log "INFO" "Binding system directories for chroot"
    sudo mount --bind /dev "$rootfs_path/dev"
    sudo mount --bind /proc "$rootfs_path/proc"
    sudo mount --bind /sys "$rootfs_path/sys"

    # Step 4: Install OpenRC and configure the service
    log "INFO" "Installing OpenRC and configuring service in Debian chroot"
    sudo chroot "$rootfs_path" /bin/bash <<'EOF'
# Update apt repositories
apt update

# Install OpenRC
apt install -y openrc

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

    # Step 5: Unmount system directories
    log "INFO" "Unmounting system directories"
    sudo umount "$rootfs_path/dev"
    sudo umount "$rootfs_path/proc"
    sudo umount "$rootfs_path/sys"

    log "SUCCESS" "OpenRC service for starter.d configured in Debian rootfs"
}

function unmount_partitions() {
    log "INFO" "Unmounting partitions..."
    sudo umount "$PRIMARY_MOUNT"
    log "SUCCESS" "Partitions unmounted."
}

function pre_cleanup_and_dir_setup() {

    local image_name=$1
    local tmp_dir=$2
    local build_dir=$3

    if [ -f "$image_name" ]; then
        rm "$image_name"
    fi

    if [ -d "$build_dir" ]; then
        rm -rf "$build_dir"
    fi

    mkdir -p "$build_dir"
    mkdir -p "$tmp_dir"
}

# Main Script Execution
OS_TYPE="alpine"
OS_VERSION="0.0.1"
MANIFEST_FILE="manifest.json"
export TARGET="access"

if [[ $# -ne 3 ]]; then
    log "ERROR" "Error: Exactly 3 arguments are required!"
    log "INFO"  "Usage: $0 <ukama_root> <node_apps> <node_id>"
    exit 1
fi

UKAMA_ROOT=$1
NODE_APPS=$2
NODE_ID=$3

check_sudo
check_requirements
pre_cleanup_and_dir_setup "$IMG_NAME" "$TMP_DIR" "${CWD}/build_access_node"

cd ${TMP_DIR}

# Build linux kernel and get rpi image (rootfs)
build_linux_kernel
download_rpi_rootfs

# Mount partition, create ukama dir, build apps
mount_partitions
setup_ukama_dirs
build_apps_using_container "${UKAMA_ROOT}" "${NODE_APPS}"
copy_all_apps              "${UKAMA_ROOT}" "${NODE_APPS}"
copy_misc_files            "${UKAMA_ROOT}" "${NODE_APPS}"
copy_linux_kernel

# setup openrc to run starter.d
configure_openrc_service "${PRIMARY_MOUNT}"
cp "${TMP_DIR}/${RPI_IMG}" "${CWD}/build_access_node/${IMG_NAME}"

# cleanup
unmount_partitions
cleanup

cd ${CWD}
log "SUCCESS" "Access node image built successfully: $IMG_NAME"
