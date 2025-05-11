#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

set -e
set -x

# Initialize variables
PARTITION_TYPE=""
ROOTFS_VERSION=""
SERVICE_NAME=""
SERVICE_CMD=""
SERVICE_ARGS=""
ARCH=""
VERSION=""
MIRROR=""

UKAMA_ROOT="/ukamarepo"
UKAMA_REPO_APP_PKG="${UKAMA_ROOT}/build/pkgs"
UKAMA_REPO_LIB_PKG="${UKAMA_ROOT}/build/libs"

UKAMA_APP_PKG="/ukama/apps/pkgs"

LOG_FILE=/setup.log
NODE_ID="uk-sa12-4567-a1"

MANIFEST_FILE="manifest.json"

# Need to pass this as arg or read from file
APP_NAMES=("wimcd" "configd" "metricsd" "lookoutd" "deviced" "notifyd" "noded" "rlog")

# Logging function
log_message() {
    log "INFO" "$(date '+%Y-%m-%d %H:%M:%S') - [Partition: $PARTITION_TYPE] [RootFS: $ROOTFS_VERSION] $1"
}

# Function to show usage
usage() {
    echo "Usage: $0 -p <partition_type> -r <rootfs_version> -n <service_name> -c <service_command> -a <service_args>"
    echo "  -p   Partition type (active or passive)"
    echo "  -r   RootFS version"
    echo "  -n   Service name"
    echo "  -c   Service command"
    echo "  -a   Service arguments (optional)"
    exit 1
}

function log() {
    local type="$1"
    local message="$2"
    local timestamp
    local file_name
    local func_name
    local color
    local reset="\033[0m"

    timestamp=$(date +"%Y-%m-%d %H:%M:%S")
    file_name=$(basename "${BASH_SOURCE[1]}")
    func_name="${FUNCNAME[1]}"

    # Set color based on log type
    case "$type" in
        INFO)
            color="\033[1;34m" # Blue
            ;;
        SUCCESS)
            color="\033[1;32m" # Green
            ;;
        WARNING)
            color="\033[1;33m" # Yellow
            ;;
        ERROR)
            color="\033[1;31m" # Red
            ;;
        *)
            color="$reset" # Default (no color)
            ;;
    esac

    printf "%s %b%s%b %s:%s \"%s\"\n" "$timestamp" "$color" "$type" "$reset" "$file_name" "$func_name" "$message" | tee -a "$LOG_FILE"
}

function LOG_EXEC() {
    log "EXEC" "$*"
    "$@" >>"$LOG_FILE" 2>&1
    if [[ $? -ne 0 ]]; then
        log "ERROR" "Command failed: $*"
        exit 1
    fi
}

function check_command() {
    command -v "$1" >/dev/null 2>&1 || {
        log "ERROR" "Command '$1' not found. Please install it."
        exit 1
    }
}

function install_starter_app() {
    log "INFO" "Installing starter.d"
    cd ${UKAMA_REPO_APP_PKG}
    tar zxvf starterd_latest.tar.gz
    cp starterd_latest/sbin/starter.d /sbin/
    rm -rf starterd_latest/
}

function copy_linux_kernel() {
    log "INFO" "Setting up kernel for ARCH=$ARCH..."

    FINAL_ROOTFS="/"
    KERNEL_TMP_DIR="/tmp/alpine-kernel-${ARCH}"
    ROOTFS_TMP_DIR="/tmp/alpine-rootfs-${ARCH}"

    case "$ARCH" in
        x86_64)
            KERNEL_PKG="linux-lts"
            ;;
        aarch64)
            KERNEL_PKG="linux-lts"
            ;;
        armv7)
            KERNEL_PKG="linux-vanilla"
            ;;
        armhf)
            log "INFO" "Using QEMU-based method to extract ARMHF kernel"
            LOG_EXEC "${UKAMA_ROOT}/builder/scripts/build-system/extract_armhf_kernel.sh"

            cp -a "${KERNEL_TMP_DIR}/boot/"* "${FINAL_ROOTFS}/boot/"
            cp -a "${KERNEL_TMP_DIR}/lib/modules/"* "${FINAL_ROOTFS}/lib/modules/"

            log "INFO" "Cleaning up temp kernel rootfs"
            rm -rf "${KERNEL_TMP_DIR}" "${ROOTFS_TMP_DIR}"

            log "SUCCESS" "ARMHF kernel installed via QEMU"
            return
            ;;
        *)
            log "ERROR" "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac

    # apk --initdb method for all other supported archs
    APK_TMP="/tmp/alpine-kernel-$ARCH"
    mkdir -p "$APK_TMP"

    LOG_EXEC apk --root "$APK_TMP" --arch "$ARCH" --update-cache \
        --repositories-file /dev/null \
        --repository "$MIRROR/$VERSION/main" add --initdb "$KERNEL_PKG"

    mkdir -p "$FINAL_ROOTFS/boot"

    case "$ARCH" in
        x86_64)
            cp "$(find "$APK_TMP/boot" -name 'vmlinuz-*' | head -n1)" "$FINAL_ROOTFS/boot/kernel.img"
            ;;
        aarch64)
            cp "$APK_TMP/boot/Image" "$FINAL_ROOTFS/boot/kernel.img"
            ;;
        armv7)
            cp "$APK_TMP/boot/zImage" "$FINAL_ROOTFS/boot/kernel.img"
            ;;
    esac

    find "$APK_TMP" -type f -name '*.dtb' -exec cp --parents {} "$FINAL_ROOTFS/boot/" \; 2>/dev/null || true
    mkdir -p "$FINAL_ROOTFS/lib/modules"
    cp -a "$APK_TMP/lib/modules/"* "$FINAL_ROOTFS/lib/modules/" 2>/dev/null || true

    rm -rf "$APK_TMP"
    log "SUCCESS" "Kernel installed using apk --initdb method"
}

function copy_all_apps() {
    log "INFO" "Copying apps"
    cp -rvf ${UKAMA_REPO_APP_PKG} ${UKAMA_APP_PKG}
}

function copy_required_libs() {
    log "INFO" "Installing required libs"
    pushd ${UKAMA_REPO_LIB_PKG}
    tar zxvf vendor_libs.tgz 
    cp -vrf ./build/* /usr/
	popd
}

function copy_misc_files() {
    
	log "INFO" "Copying various files to image"
    rm -f ${MANIFEST_FILE}
    create_manifest_file
    
    #sudo cp ${MANIFEST_FILE} "/manifest.json"

    # install the starter.d app
    install_starter_app "/"

    # update /etc/services to add ports
    log "INFO" "Adding all the apps to /etc/services"
    sudo mkdir -p "/etc"
    sudo cp "${UKAMA_ROOT}/nodes/ukamaOS/distro/scripts/files/services" \
         "/etc/services"
}


# Update /etc/fstab based on partition type
update_fstab() {
    log_message "Updating /etc/fstab for partition type: $PARTITION_TYPE"

    if [[ "$PARTITION_TYPE" == "active" ]]; then
        cat <<FSTAB > /etc/fstab
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
        cat <<FSTAB > /etc/fstab
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

    log_message "/etc/fstab updated successfully."
}

# Configure network interface eth0
configure_network() {
    log_message "Configuring network for eth0"

    cat <<NETWORK > /etc/network/interfaces
auto eth0
iface eth0 inet static
    address 10.102.81.10
    netmask 255.255.255.0
    gateway 10.102.81.1
NETWORK

    # Apply network changes
    #ifdown eth0 && ifup eth0
    log_message "Network configuration updated for eth0"
}

# Create a custom OpenRC service
create_openrc_service() {
    log_message "Creating OpenRC service: $SERVICE_NAME"

    mkdir -p /etc/init.d

    cat <<SERVICE > /etc/init.d/$SERVICE_NAME
#!/sbin/openrc-run

description="OpenRC Service: $SERVICE_NAME"
command="$SERVICE_CMD"
command_args="$SERVICE_ARGS"

depend() {
    need net
}

start() {
    ebegin "Starting $SERVICE_NAME"
    start-stop-daemon --start --background --exec \$command -- \$command_args
    eend \$?
}

stop() {
    ebegin "Stopping $SERVICE_NAME"
    start-stop-daemon --stop --exec \$command
    eend \$?
}
SERVICE

    chmod +x /etc/init.d/$SERVICE_NAME
    rc-update add $SERVICE_NAME default
    log_message "OpenRC service $SERVICE_NAME created and added to startup."
}

# Function to set up the root filesystem
setup_rootfs() {
    log_message "Setting up root filesystem"

    # Set up DNS
    echo "nameserver 8.8.8.8" > /etc/resolv.conf

    # Set up package repositories
    echo "https://dl-cdn.alpinelinux.org/alpine/${ROOTFS_VERSION}/main" > /etc/apk/repositories
    echo "https://dl-cdn.alpinelinux.org/alpine/${ROOTFS_VERSION}/community" >> /etc/apk/repositories

    # Update packages
    apk update
    apk upgrade

    # Install essential packages
    apk add alpine-base openrc busybox bash sudo shadow tzdata
    apk add acpid busybox-openrc busybox-extras busybox-mdev-openrc
    apk add readline bash autoconf automake libmicrohttpd-dev gnutls-dev openssl-dev \
        iptables libuuid sqlite dhcpcd protobuf iproute2 zlib curl-dev nettle libcap \
        libidn2 libmicrohttpd gnutls openssl-dev curl-dev linux-headers bsd-compat-headers \
        tree libtool sqlite-dev openssl-dev readline cmake autoconf automake alpine-sdk \
        build-base git tcpdump ethtool iperf3 htop vim doas

    # Set timezone
    ln -sf /usr/share/zoneinfo/UTC /etc/localtime

    # Configure networking
    apk add dhcpcd iproute2 iputils
    rc-update add dhcpcd default
#    rc-service dhcpcd start

    # Set hostname
    echo "ukama-linux" > /etc/hostname

    # Set up root user
    echo "root:root" | chpasswd

    # Create 'ukama' user only if it doesn't already exist
    if ! id "ukama" >/dev/null 2>&1; then
        adduser -D -s /bin/bash -G wheel ukama
        echo "ukama:ukama" | chpasswd
    fi
    echo "%wheel ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/wheel

    # Configure doas (instead of sudo)
    apk add doas
    echo "permit persist ukama as root" > /etc/doas.d/doas.conf
    chmod 600 /etc/doas.d/doas.conf

    # Enable SSH access
    apk add openssh
    rc-update add sshd default
#    rc-service sshd start

    # Enable system services
    rc-update add networking default
    rc-update add sshd default
    rc-update add dhcpcd default
    rc-update add acpid default

    log_message "INFO" "root filesystem setup completed."
}

function create_manifest_file() {
    log "INFO" "Creating manifest file"

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
  echo "Adding manifest for ${APP_NAMES[@]}"
  for app in "${APP_NAMES[@]}"; do
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

  for app in "${APP_NAMES[@]}"; do
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

function setup_ukama_dirs() {
    log "INFO" "Creating Ukama directories..."

    mkdir -p "/ukama"
    mkdir -p "/ukama/configs"
    mkdir -p "/ukama/apps"
    mkdir -p "/ukama/apps/pkgs"
    mkdir -p "/ukama/apps/rootfs"
    mkdir -p "/ukama/apps/registry"
    mkdir -p "/passive"
    mkdir -p "/boot/firmware"
    mkdir -p "/data"
    mkdir -p "/recovery"

    echo "${NODE_ID}" > "/ukama/nodeid"
    echo "localhost"  > "/ukama/bootstrap"

    touch "/ukama/apps.log"

    log "SUCCESS" "Ukama directories created."
}

function get_apps_name() {

	# Loop through each .tar.gz file in the directory
	for file in "${UKAMA_REPO_APP_PKG}"/*.tar.gz; do
    	[ -e "$file" ] || continue  # Skip if no .tar.gz files exist
    	filename=$(basename "$file")  # Get filename without path
    	prefix=${filename%%_*}  # Extract prefix before first underscore
    	APP_NAMES+=("$prefix")  # Store in array
	done

	# Print the array elements
	echo "Extracted app prefixes: ${APP_NAMES[@]}"
}

#Main 
setup_ukama_dirs

log "INFO" "Script ${0} called with args $#"

index=0
for arg in "$@"; do
  log "INFO" "arg[${index}]: ${arg}"
  index=$((index + 1))
done

# Parse options using getopts
while getopts "p:r:n:c:a:A:V:M:" opt; do
    case "${opt}" in
        p) PARTITION_TYPE="${OPTARG}" ;;
        r) ROOTFS_VERSION="${OPTARG}" ;;
        n) SERVICE_NAME="${OPTARG}" ;;
        c) SERVICE_CMD="${OPTARG}" ;;
        a) SERVICE_ARGS="${OPTARG}" ;;
        A) ARCH="${OPTARG}" ;;
        V) VERSION="${OPTARG}" ;;
        M) MIRROR="${OPTARG}" ;;
        *) usage ;;
    esac
done

# Validate required arguments
if [[ -z "$PARTITION_TYPE" || -z "$ROOTFS_VERSION" || -z "$SERVICE_NAME" || -z "$SERVICE_CMD" ]]; then
    usage
fi

# Validate partition type
if [[ "$PARTITION_TYPE" != "active" && "$PARTITION_TYPE" != "passive" ]]; then
    echo "Error: Partition type must be 'active' or 'passive'."
    exit 1
fi

setup_rootfs  # Set up root filesystem

log "INFO" "Preparing mount for ${PARTITION_TYPE}"
update_fstab  ${PARTITION_TYPE} #Update fstab

log "INFO" "Network configuration"
configure_network  # Configure network

log "INFO" "OpenRC service steup for  ${SERVICE_NAME}  ${SERVICE_CMD}"
create_openrc_service ${SERVICE_NAME}  ${SERVICE_CMD}# Create OpenRC service

log "INFO" "Copying libs"
copy_required_libs

log "INFO" "Copying apps"
copy_all_apps
#get_apps_name

log "INFO" "Create manifest."
copy_misc_files 

log "INFO" "Copy kernel"
copy_linux_kernel

echo "Roootfs build success."
exit 0
~                                                                     
