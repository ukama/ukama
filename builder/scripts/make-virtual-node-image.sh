#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

# this script create a bootable image which can then run with QEMU
# and act as Ukama Node. Image is Ubuntu (20.04) with 5GB HD.

set -e  # Exit immediately if a command exits with a non-zero status.

UBUNTU_ISO_URL="https://releases.ubuntu.com/22.04/ubuntu-22.04.4-live-server-amd64.iso"
ISO_FILE="ubuntu.iso"
IMG_SIZE="5G"
BOOTSTRAP_PORT=0
NODE_ID=$1
UKAMA_REPO=$2
IMG_FILE="$NODE_ID.img"

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

# Function to check if the script is run as root
check_root() {
    if [ "$EUID" -ne 0 ]; then
        log "ERROR" "This script must be run as root."
        exit 1
    fi
    log "INFO" "Running as root, proceeding..."
}

# Function to install necessary packages
install_prerequisites() {
    log "INFO" "Installing necessary packages..."
    apt-get update
    apt-get install -y qemu-system virt-manager virt-viewer libvirt-daemon-system \
        libvirt-clients bridge-utils debootstrap extlinux kpartx
    check_status $? "Packages installed successfully" "install_prerequisites"
}

# Function to download the Ubuntu ISO
download_ubuntu_iso() {
    log "INFO" "Downloading Ubuntu 22.04 (jammy) ISO..."
    wget --no-check-certificate $UBUNTU_ISO_URL -O $ISO_FILE
    check_status $? "Ubuntu ISO downloaded" "download_ubuntu_iso"
}

# Function to create and partition the disk image
create_disk_image() {
    log "INFO" "Creating and partitioning disk image..."
    qemu-img create -f raw $IMG_FILE $IMG_SIZE 
    check_status $? "Raw disk image created" "create_disk_image"
    
    LOOP_DEVICE=$(losetup -fP --show $IMG_FILE)
    sleep 5
    echo -e "o\nn\np\n1\n\n\nw" | fdisk $LOOP_DEVICE
    check_status $? "Disk partitioned" "create_disk_image"
    
    partprobe $LOOP_DEVICE
    mkfs.ext4 ${LOOP_DEVICE}p1
    check_status $? "Filesystem created" "create_disk_image"
}

# Function to mount the disk image
mount_disk_image() {
    log "INFO" "Mounting the disk image..."
    mkdir -p /mnt/image
    mount ${LOOP_DEVICE}p1 /mnt/image
    check_status $? "Disk image mounted" "mount_disk_image"
}

# Function to install Ubuntu on the image
install_ubuntu() {
    log "INFO" "Installing Ubuntu on the disk..."
    debootstrap --arch amd64 jammy /mnt/image
    check_status $? "Ubuntu installed" "install_ubuntu"
}

# Function to set up the chroot environment and install additional packages
setup_chroot_environment() {
    log "INFO" "Setting up chroot environment..."
    mount --bind /dev  /mnt/image/dev
    mount --bind /proc /mnt/image/proc
    mount --bind /sys  /mnt/image/sys

    log "INFO" "Installing packages inside the chroot environment..."
    get_bootstrap_port $UKAMA_REPO
    export BOOTSTRAP_PORT
    export NODE_ID

    chroot /mnt/image /bin/bash <<'EOL'
        set -e	
        export DEBIAN_FRONTEND=noninteractive
        locale-gen en_US.UTF-8
        update-locale LANG=en_US.UTF-8
        debconf-set-selections <<< "grub-pc grub-pc/install_devices_empty boolean true"
        apt-get update
        apt-get install -y -o Dpkg::Options::="--force-confnew" linux-image-generic

        mkdir -p /ukama /ukama/configs /ukama/apps/pkgs /ukama/apps/rootfs /ukama/apps/registry
        echo $NODE_ID > /ukama/nodeid
        echo "localhost" > /ukama/bootstrap
        touch /ukama/apps.log

        cat > /etc/systemd/system/starterd.service << EOF
        [Unit]
        Description=Ukama's capp starter.d
        After=network.target

        [Service]
        ExecStart=/sbin/starter.d --manifest-file /manifest.json
        Type=simple

        [Install]
        WantedBy=multi-user.target
        EOF

        systemctl enable starterd.service
EOL

    check_status $? "Chroot environment set up" "setup_chroot_environment"
}

# Function to unmount the filesystems
unmount_filesystems() {
    log "INFO" "Unmounting filesystems..."
    umount /mnt/image/dev /mnt/image/proc /mnt/image/sys
    check_status $? "Filesystems unmounted" "unmount_filesystems"
}

# Function to set up EXTLINUX bootloader
setup_bootloader() {
    log "INFO" "Setting up EXTLINUX..."
    extlinux --install /mnt/image/boot
    check_status $? "EXTLINUX bootloader installed" "setup_bootloader"

    cat <<EOF > /mnt/image/boot/extlinux.cfg
DEFAULT linux
LABEL linux
    KERNEL /boot/vmlinuz-$(ls /mnt/image/boot/ | grep vmlinuz | head -n 1)
    APPEND root=${LOOP_DEVICE}p1 ro quiet
EOF
    check_status $? "EXTLINUX bootloader configured" "setup_bootloader"
}

# Function to extract kernel and initRAMFS
extract_kernel_and_initramfs() {
    log "INFO" "Extracting kernel and initRAMFS from the OS image..."
    mkdir -p /mnt/${NODE_ID}
    mount -o loop,offset=$((512*2048)) ${IMG_FILE} /mnt/${NODE_ID}
    check_status $? "OS image mounted for extraction" "extract_kernel_and_initramfs"

    cp /mnt/${NODE_ID}/boot/vmlinuz-* .
    cp /mnt/${NODE_ID}/boot/initrd.img-* .
    check_status $? "Kernel and initRAMFS extracted" "extract_kernel_and_initramfs"

    umount /mnt/${NODE_ID}
    rmdir /mnt/${NODE_ID}
    check_status $? "Unmounted and cleaned up image extraction" "extract_kernel_and_initramfs"
}

# Function to get the bootstrap port
get_bootstrap_port() {
    log "INFO" "Getting bootstrap port..."
    local repo=$1
    local service_name="node-gateway-init"
    local compose_file="${repo}/systems/init/docker-compose.yml"

    # Extract the port
    local port=$(grep -A 10 "${service_name}:" "${compose_file}" | \
                     grep -A 2 ports | awk -F"'" '{print $2}' | cut -d ':' -f 1)

    if [ -z "$port" ]; then
        log "ERROR" "Port not found for service ${service_name}"
        exit 1
    else
        BOOTSTRAP_PORT=$port
    fi
    check_status $? "Bootstrap port retrieved" "get_bootstrap_port"
}

# Main
check_root

log "INFO" "Starting the script to create Ukama Node Image"

install_prerequisites
download_ubuntu_iso
create_disk_image
mount_disk_image
install_ubuntu
setup_chroot_environment
unmount_filesystems
setup_bootloader
extract_kernel_and_initramfs

log "SUCCESS" "Node image creation completed successfully!"

exit 0
