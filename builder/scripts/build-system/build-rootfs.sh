#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

set -e

# Initialize variables
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
    log "INFO" "$(date '+%Y-%m-%d %H:%M:%S') - [RootFS: $VERSION] $1"
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

log() {
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

LOG_EXEC() {
    log "EXEC" "$*"
    "$@" >>"$LOG_FILE" 2>&1
    if [[ $? -ne 0 ]]; then
        log "ERROR" "Command failed: $*"
        exit 1
    fi
}

check_command() {
    command -v "$1" >/dev/null 2>&1 || {
        log "ERROR" "Command '$1' not found. Please install it."
        exit 1
    }
}

install_starter_app() {
    log "INFO" "Installing starter.d"

    cd ${UKAMA_REPO_APP_PKG}
    tar zxvf starterd_latest.tar.gz

    cp -rf starterd_latest/lib/*      /ukama/apps/lib
    cp -rf starterd_latest/usr/lib/*  /ukama/apps/lib
    cp starterd_latest/sbin/starter.d /sbin/

    rm -rf starterd_latest/
}

install_rpi4_kernel_from_tarball() {
    log "INFO" "Installing RPi4 kernel and boot files via Alpine RPi tarball"

    ALPINE_VERSION="${VERSION#v}"
    ALPINE_RPI_URL="https://dl-cdn.alpinelinux.org/alpine/v${ALPINE_VERSION}/releases/aarch64/alpine-rpi-${ALPINE_VERSION}.0-aarch64.tar.gz"
    TMP_RPI_DIR="/tmp/alpine-rpi"
    FINAL_BOOT="/boot"

    mkdir -p "$TMP_RPI_DIR/rootfs" "$FINAL_BOOT"

    log "INFO" "Downloading: $ALPINE_RPI_URL"
    wget -qO "$TMP_RPI_DIR/rpi.tar.gz" "$ALPINE_RPI_URL" || {
        log "ERROR" "Failed to download $ALPINE_RPI_URL"
        exit 1
    }

    log "INFO" "Extracting Alpine RPi image"
    tar -xzf "$TMP_RPI_DIR/rpi.tar.gz" -C "$TMP_RPI_DIR/rootfs"

    log "INFO" "Copying kernel to /boot/kernel.img"
    cp "$TMP_RPI_DIR/rootfs/boot/vmlinuz-rpi" "$FINAL_BOOT/kernel.img" || {
        log "ERROR" "Missing vmlinuz-rpi in tarball"
        exit 1
    }

    log "INFO" "Copying bootloader firmware and configs"
    cp "$TMP_RPI_DIR/rootfs"/bootcode.bin "$FINAL_BOOT/" 2>/dev/null || true
    cp "$TMP_RPI_DIR/rootfs"/start*.elf   "$FINAL_BOOT/" 2>/dev/null || true
    cp "$TMP_RPI_DIR/rootfs"/fixup*.dat   "$FINAL_BOOT/" 2>/dev/null || true
    cp "$TMP_RPI_DIR/rootfs"/config.txt   "$FINAL_BOOT/" 2>/dev/null || true
    cp "$TMP_RPI_DIR/rootfs"/cmdline.txt  "$FINAL_BOOT/" 2>/dev/null || true
    cp "$TMP_RPI_DIR/rootfs"/*.dtb        "$FINAL_BOOT/" 2>/dev/null || true

    log "INFO" "Copying overlays"
    mkdir -p "$FINAL_BOOT/overlays"
    cp -a "$TMP_RPI_DIR/rootfs/overlays/"* "$FINAL_BOOT/overlays/" 2>/dev/null || true

    if [ -d "$TMP_RPI_DIR/rootfs/lib/modules" ]; then
        log "INFO" "Copying kernel modules"
        mkdir -p "/lib/modules"
        cp -a "$TMP_RPI_DIR/rootfs/lib/modules/"* "/lib/modules/"
    else
        log "WARNING" "No /lib/modules found in RPi tarball"
    fi

    rm -rf "$TMP_RPI_DIR"
    log "SUCCESS" "RPi4 kernel, firmware, DTBs, and modules installed"
}

install_amplifier_toolchain() {
    log_message "INFO: Ensuring ARM cross-toolchain for amplifier is installed"

    # If the expected prefix already exists, nothing to do
    if command -v arm-linux-gnueabihf-gcc >/dev/null 2>&1; then
        log_message "INFO: arm-linux-gnueabihf-gcc already available, skipping toolchain install"
        return 0
    fi

    log_message "INFO: Installing gcc-arm-none-eabi (provides arm-none-eabi-* tools)"
    apk update
    apk add --no-cache gcc-arm-none-eabi || {
        log_message "ERROR: Failed to install gcc-arm-none-eabi toolchain via apk"
        exit 1
    }

    # Create arm-linux-gnueabihf-* symlinks pointing to arm-none-eabi-*
    local bindir="/usr/bin"
    # include ld.bfd because some build systems use it explicitly
    local tools=(gcc g++ cpp objcopy objdump ar as ld ld.bfd nm ranlib strip)

    for exe in "${tools[@]}"; do
        if [ -x "${bindir}/arm-none-eabi-${exe}" ] && [ ! -e "${bindir}/arm-linux-gnueabihf-${exe}" ]; then
            ln -s "arm-none-eabi-${exe}" "${bindir}/arm-linux-gnueabihf-${exe}"
            log_message "INFO: Linked arm-linux-gnueabihf-${exe} -> arm-none-eabi-${exe}"
        fi
    done

    if ! command -v arm-linux-gnueabihf-gcc >/dev/null 2>&1; then
        log_message "ERROR: arm-linux-gnueabihf-gcc still not found after installing toolchain"
        exit 1
    fi

    log_message "SUCCESS: ARM cross-toolchain for amplifier is ready"
}

build_armv7_boot() {
    local node=$1
    local path="${UKAMA_ROOT}/nodes/ukamaOS/firmware"
    local boot1="${path}/build/boot/at91bootstrap/at91bootstrap.bin"
    local boot2="${path}/build/boot/uboot/u-boot.bin"

    cwd=$(pwd)
    log_message "INFO: Building firmware for Node: ${node}"

    # 1) Ensure the cross-toolchain is installed and on PATH
    install_amplifier_toolchain

    # 2) Remove any stale host-built Kconfig 'conf' binary so it rebuilds
    if [ -f "${path}/at91-bootstrap/config/conf" ]; then
        log_message "INFO: Removing stale at91-bootstrap/config/conf to force rebuild inside rootfs"
        rm -f "${path}/at91-bootstrap/config/conf"
    fi

    cd "${path}"
    make clean TARGET="${node}" ROOTFSPATH="${path}/build"
    make TARGET="${node}" ROOTFSPATH="${path}/build"

    if [ ! -f "$boot1" ]; then
        log_message "ERROR: Firmware build failure. boot file does not exist: $boot1"
        exit 1
    fi

    if [ ! -f "$boot2" ]; then
        log_message "ERROR: Firmware build failure. boot file does not exist: $boot2"
        exit 1
    fi

    log_message "INFO: Firmware build OK"
    cd "${cwd}"
}

copy_armv7_boot() {

    local path="${UKAMA_ROOT}/nodes/ukamaOS/firmware"
    local boot1="${path}/build/boot/at91bootstrap/at91bootstrap.bin"
    local boot2="${path}/build/boot/uboot/u-boot.bin"

    cp "$boot1" /boot/BOOT.BIN
    cp "$boot2" /boot/

    # XXX also need to copy the ukama_anode.dtb and uboot.env file XXX
    # tree boot
    #boot
    #├── BOOT.BIN
    #├── boot.tar
    #├── firmware
    #├── u-boot.bin
    #├── uboot.env
    #└── ukama_anode.dtb
}

copy_armv7_kernel() {
    log_message "Installing linux-lts kernel and modules"

    # Update the APK index and install the package
    apk update
    apk add --no-cache linux-lts

    log_message "SUCCESS: linux-lts package and modules installed"
}

copy_x86_64_boot() {
    log_message "Extracting boot files from Alpine ISO"

    local ver="${VERSION#v}"
    local iso_url="${MIRROR}/${VERSION}/releases/x86_64/alpine-standard-${ver}.0-x86_64.iso"
    local tmpdir mnt
    tmpdir=$(mktemp -d)
    mnt="$tmpdir/mnt"

    mkdir -p "$mnt" /boot /boot/efi
    trap 'umount "$mnt" 2>/dev/null; rm -rf "$tmpdir"' EXIT

    log_message "Downloading ISO: $iso_url"
    curl -fsSL "$iso_url" -o "$tmpdir/alpine.iso"

    log_message "Mounting ISO"
    mount -o loop "$tmpdir/alpine.iso" "$mnt"

    log_message "Copying kernel, initramfs, and bootloader"
    cp -a "$mnt/boot/." /boot/
    cp -a "$mnt/efi/."  /boot/efi/

    log_message "Unmounting ISO"
    umount "$mnt"
    trap - EXIT

    log_message "Boot files extracted to /boot and /boot/efi"
}

copy_x86_64_kernel() {
    log_message "Installing linux-lts kernel and modules"

    # Update the APK index and install the package
    apk update
    apk add --no-cache linux-lts

    log_message "SUCCESS: linux-lts package and modules installed"
}

copy_misc_files() {
	log "INFO" "Copying various files to image"

    # install the starter.d app
    install_starter_app "/"

    # update /etc/services to add ports
    log "INFO" "Adding all the apps to /etc/services"
    cp "${UKAMA_ROOT}/nodes/ukamaOS/distro/scripts/files/services" \
       "/etc/services"

    # copy mocksysfs related files (not needed for actual HW) - XXX
    mkdir -p "/tmp/sys"
    cp -rf ${UKAMA_ROOT}/builder/scripts/build-system/mocksysfs/* \
       "/tmp/sys/"
    cp -rf "${UKAMA_ROOT}/nodes/ukamaOS/distro/system/noded/mfgdata" \
       "/ukama/mocksysfs/"
}

# Function to create and register a custom OpenRC service
setup_openrc_service() {
    log_message "Creating OpenRC service: $SERVICE_NAME"

    # Ensure init.d directory exists
    mkdir -p /etc/init.d

    # Write the service script
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
    start-stop-daemon --start --background --exec "$SERVICE_CMD" -- $SERVICE_ARGS
    eend $?
}

stop() {
    ebegin "Stopping $SERVICE_NAME"
    start-stop-daemon --stop --exec "$SERVICE_CMD"
    eend $?
}
SERVICE

    # Make it executable and add to default runlevel
    chmod +x /etc/init.d/$SERVICE_NAME
    rc-update add $SERVICE_NAME default

    log_message "INFO" "OpenRC service $SERVICE_NAME created and added to startup."
}

setup_rootfs() {
    log_message "Setting up root filesystem"

    cat > /etc/resolv.conf <<EOF
nameserver 8.8.8.8
EOF
    cat > /etc/apk/repositories <<EOF
https://dl-cdn.alpinelinux.org/alpine/${VERSION}/main
https://dl-cdn.alpinelinux.org/alpine/${VERSION}/community
EOF

    COMMON_DEV_PACKAGES="alpine-base bash sudo shadow tzdata openrc
            eudev eudev-openrc
            kmod dosfstools
            acpid dhcpcd iproute2 iputils
            openssh
            readline autoconf automake cmake
            alpine-sdk build-base libtool
            openssl-dev gnutls-dev curl curl-dev
            sqlite-dev zlib libuuid libcap libidn2 libmicrohttpd-dev
            protobuf e2fsprogs util-linux rsync jansson tree
            git tcpdump ethtool iperf3 htop vim doas
            kbd bison flex"

    ARM_PACKAGES="dtc coreutils"

    # Install base
    apk update && apk upgrade
    apk add --no-cache $COMMON_DEV_PACKAGES

    # Conditionally install extras
    if [[ "$ARCH" == "armv7" || "$ARCH" == "armv7l" ]]; then
        apk add --no-cache $ARM_PACKAGES
    fi

    cat > /etc/dhcpcd.conf <<'EOF'
interface eth0
static ip_address=10.102.81.10/24
static routers=10.102.81.1
static domain_name_servers=8.8.8.8 8.8.4.4
EOF

    # Enable agetty on /dev/tty1
    ln -sf /etc/init.d/agetty /etc/init.d/agetty.tty1

    # sysinit: early, before any daemons
    rc-update add devfs           sysinit    # mounts /dev,/proc,/sys
    rc-update add modules         sysinit    # modprobe usbcore, ehci_hcd, usbhid, vfat, etc.
    rc-update add loadkmap        sysinit    # apply your US keymap
    rc-update add udev-trigger    sysinit    # now sees the USB HID devices
    rc-update add udev-settle     sysinit
    rc-update add udev            boot

    # boot: filesystem mounts, hostname, syslog
    rc-update add sysctl          boot
    rc-update add bootmisc        boot
    rc-update add hostname        boot
    rc-update add syslog          boot

    # default: network, console, SSH
    rc-update add dhcpcd          default
    rc-update add sshd            default
    rc-update add acpid           default
    rc-update add agetty.tty1     default

    # vfat support
    mkdir -p /etc/modules-load.d
    cat > /etc/modules-load.d/vfat.conf <<EOF
vfat
fat
EOF

    # USB keyboard
    cat > /etc/modules-load.d/keyboard.conf <<EOF
usbcore
ehci_hcd
xhci_hcd
usbhid
hid_generic
EOF

    ln -sf /usr/share/zoneinfo/UTC /etc/localtime
    echo "ukama-linux" > /etc/hostname

    echo "root:root" | chpasswd
    if ! id ukama &>/dev/null; then
        adduser -D -s /bin/bash -G wheel ukama
        echo "ukama:ukama" | chpasswd
    fi
    echo "%wheel ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/wheel

    cat > /etc/doas.d/doas.conf <<EOF
permit persist ukama as root
EOF
    chmod 600 /etc/doas.d/doas.conf

    log_message "INFO" "Root filesystem setup completed."
}

setup_ukama_dirs() {
    log "INFO" "Creating Ukama directories..."

    mkdir -p "/ukama/configs"
    mkdir -p "/ukama/apps/lib"
    mkdir -p "/ukama/apps/pkgs"
    mkdir -p "/ukama/apps/rootfs"
    mkdir -p "/ukama/apps/registry"
    mkdir -p "/ukama/mocksysfs"

    echo "${NODE_ID}" > "/ukama/nodeid"
    echo "localhost"  > "/ukama/bootstrap"

    touch "/ukama/apps.log"

    log "SUCCESS" "Ukama directories created."
}

# Main
log "INFO" "Script ${0} called with args $#"

index=0
for arg in "$@"; do
  log "INFO" "arg[${index}]: ${arg}"
  index=$((index + 1))
done

# Parse options using getopts
while getopts "r:n:c:a:A:V:M:" opt; do
    case "${opt}" in
        n) SERVICE_NAME="${OPTARG}" ;;
        c) SERVICE_CMD="${OPTARG}" ;;
        a) SERVICE_ARGS="${OPTARG}" ;;
        A) ARCH="${OPTARG}" ;;
        V) VERSION="${OPTARG}" ;;
        M) MIRROR="${OPTARG}" ;;
        *) usage ;;
    esac
done

setup_rootfs
setup_ukama_dirs
setup_openrc_service "${SERVICE_NAME}" "${SERVICE_CMD}"
copy_misc_files

case "${ARCH}" in
    x86_64)
        copy_x86_64_kernel
        copy_x86_64_boot
        ;;
    armv7|armv7l)
        build_armv7_boot "amplifier"
        copy_armv7_kernel
        copy_armv7_boot
        ;;
    *)
        echo "Unsupported architecture: ${ARCH}"
        exit 1
        ;;
esac

echo "Rootfs build success."
exit 0
