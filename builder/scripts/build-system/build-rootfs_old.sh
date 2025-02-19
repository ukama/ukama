#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

# Script to build and package ukamaOS app

set -e

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

# Initialize variables
PARTITION_TYPE=""
ROOTFS_VERSION=""
SERVICE_NAME=""
SERVICE_CMD=""
SERVICE_ARGS=""
MAJOR_VERSION="v3.21"

# Parse options using getopts
while getopts "p:n:c:a:" opt; do
    case "${opt}" in
        p) PARTITION_TYPE="${OPTARG}" ;;
	r) ROOTFS_VERSION="${OPTARG}" ;;
        n) SERVICE_NAME="${OPTARG}" ;;
        c) SERVICE_CMD="${OPTARG}" ;;
        a) SERVICE_ARGS="${OPTARG}" ;;
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

# Logging function
log_message() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - [Partition: $PARTITION_TYPE] [RootFS: $ROOTFS_VERSION] $1"
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
    ifdown eth0 && ifup eth0
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
    apk add readline bash autoconf automake libmicrohttpd-dev gnutls-dev openssl-dev iptables libuuid sqlite dhcpcd protobuf iproute2 zlib curl-dev nettle libcap libidn2 libmicrohttpd gnutls openssl-dev curl-dev linux-headers bsd-compat-headers tree libtool sqlite-dev openssl-dev readline cmake autoconf automake alpine-sdk build-base git tcpdump ethtool iperf3 htop vim doas

    # Set timezone
    ln -sf /usr/share/zoneinfo/UTC /etc/localtime

    # Configure networking
    apk add dhcpcd iproute2 iputils
    rc-update add dhcpcd default
    rc-service dhcpcd start

    # Set hostname
    echo "ukama-linux" > /etc/hostname

    # Set up root user
    echo "root:root" | chpasswd

    # Create a new user
    adduser -D -s /bin/bash -G wheel ukama
    echo "ukama:ukama" | chpasswd
    echo "%wheel ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/wheel

    # Configure doas (instead of sudo)
    apk add doas
    echo "permit persist ukama as root" > /etc/doas.d/doas.conf
    chmod 600 /etc/doas.d/doas.conf

    # Enable SSH access
    apk add openssh
    rc-update add sshd default
    rc-service sshd start

    # Enable system services
    rc-update add networking default
    rc-update add sshd default
    rc-update add dhcpcd default
    rc-update add acpid default

    # Create necessary directories
    mkdir -p /recovery /data /passive /boot/firmware

    log_message "Root filesystem setup completed."
}


# Add packages


# Add vendor libs


#copy anyother scripts required

# Main execution
setup_rootfs  # Set up root filesystem
update_fstab  # Update fstab
configure_network  # Configure network
create_openrc_service  # Create OpenRC service


