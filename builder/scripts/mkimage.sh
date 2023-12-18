#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

# this script create a bootable image which can then run with QEMU
# and act as Ukama Node. Image is Ubuntu (20.04) with 5GB HD.

set -e  # Exit immediately if a command exits with a non-zero status.

NODE_ID=$1

UBUNTU_ISO_URL="https://releases.ubuntu.com/22.04/ubuntu-22.04.3-live-server-amd64.iso"
ISO_FILE="ubuntu.iso"
IMG_FILE="$NODE_ID.img"
IMG_SIZE="5G"

# Step 0: make sure have all the right packages
apt-get update
apt-get install -y qemu-kvm qemu virt-manager virt-viewer libvirt-daemon-system \
      libvirt-clients bridge-utils debootstrap \
      extlinux kpartx

# Step 1: Download Ubuntu ISO
echo "Downloading Ubuntu 22.04 (jammy) ISO..."
wget $UBUNTU_ISO_URL -O $ISO_FILE || { echo "Failed to download ISO"; exit 1; }

# Step 2: Create a Raw Disk Image, format, parition and mount 
echo "Creating and partitioning disk image..."
qemu-img create -f raw $IMG_FILE $IMG_SIZE 
LOOP_DEVICE=$(losetup -fP --show $IMG_FILE)
sleep 5
echo -e "o\nn\np\n1\n\n\nw" | fdisk $LOOP_DEVICE
partprobe $LOOP_DEVICE
mkfs.ext4 ${LOOP_DEVICE}p1

# Mount the partition
mkdir -p /mnt/image
mount ${LOOP_DEVICE}p1 /mnt/image || { echo "Unable to mount the partition"; exit 1;}

# Step 3: Install Ubuntu on the Disk Image
echo "Installing Ubuntu on the disk..."
debootstrap --arch amd64 jammy /mnt/image || { echo "Debootstrap failed"; exit 1; }

# Mounting necessary filesystems and setting up chroot
mount --bind /dev  /mnt/image/dev
mount --bind /proc /mnt/image/proc
mount --bind /sys  /mnt/image/sys

echo "Installing packages on the disk ..."
chroot /mnt/image /bin/bash <<'EOL'
    set -e	
    export DEBIAN_FRONTEND=noninteractive
    locale-gen en_US.UTF-8
    update-locale LANG=en_US.UTF-8
    debconf-set-selections <<< "grub-pc grub-pc/install_devices_empty boolean true"
    apt-get update
    apt-get install -y -o Dpkg::Options::="--force-confnew" linux-image-generic

    mkdir -p /capps/pkgs
    mkdir -p /capps/rootfs
    mkdir -p /capps/registry

    mkdir -p /ukama
    echo $NODE_ID > /ukama/nodeid

    # create systemd service for the starter.d program
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

    # Enable the service
    systemctl enable starterd.service

EOL

# Unmount filesystems
umount /mnt/image/dev /mnt/image/proc /mnt/image/sys

# Step 4: Set Up  EXTLINUX as Bootloader and create the cfg file
echo "Setting up EXTLINUX..."
extlinux --install /mnt/image/boot

cat <<EOF > /mnt/image/boot/extlinux.cfg
DEFAULT linux
LABEL linux
    KERNEL /boot/vmlinuz-$(ls /mnt/image/boot/ | grep vmlinuz | head -n 1)
    APPEND root=${LOOP_DEVICE}p1 ro quiet
EOF

# Unmount and detach loop device
umount /mnt/image
losetup -d $LOOP_DEVICE
rm -f $ISO_FILE
sleep 5

# Step 5: mount the image and extract kenerl and initramFS. This will
# be needed by QEMU to run the image.
#
# QEMU command would be:
# sudo qemu-system-x86_64 -hda ${IMG_FILE} -m 1024 -kernel ./vmlinuz-5.4.0-26-generic \
    # -initrd ./initrd.img-5.4.0-26-generic -append "root=/dev/sda1"
echo "Extracting kernel and initRAMfs from the OS image"
mkdir -p /mnt/${NODE_ID}
mount -o loop,offset=$((512*2048)) ${IMG_FILE} /mnt/${NODE_ID}

cp /mnt/${NODE_ID}/boot/vmlinuz-*    .
cp /mnt/${NODE_ID}/boot/initrd.img-* .

echo "Cleanup and done!"
# Unmount the image
umount /mnt/${NODE_ID}
rmdir /mnt/${NODE_ID}

exit 0
