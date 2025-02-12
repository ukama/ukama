#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

# Script to build and package ukamaOS app

# Set the installation directory as "rootfs" inside the current directory
INSTALL_DIR="$(pwd)/rootfs"

# Default values
ARCH="x86_64"
VERSION="latest-stable"
MIRROR="http://dl-cdn.alpinelinux.org/alpine"
MOUNT_SRC=""   # Source directory to mount
MOUNT_DEST=""  # Destination inside chroot

# Parse command-line arguments
while getopts "a:v:m:M:h" opt; do
  case $opt in
    a) ARCH="$OPTARG" ;;
    v) VERSION="$OPTARG" ;;
    m) MIRROR="$OPTARG" ;;
    M) 
      IFS=":" read -r MOUNT_SRC MOUNT_DEST <<< "$OPTARG"
      ;;
    h) 
      echo "Usage: $0 [-a arch] [-v version] [-m mirror] [-M source:destination]"
      echo "Example: $0 -a x86_64 -v 3.18 -m http://dl-cdn.alpinelinux.org/alpine -M /home/user/shared:/mnt/shared"
      exit 0
      ;;
    *) echo "Invalid option"; exit 1 ;;
  esac
done

echo "Getting alpine-chroot-command."

wget -O alpine-chroot-install https://raw.githubusercontent.com/alpinelinux/alpine-chroot-install/master/alpine-chroot-install

chmod +x alpine-chroot-install
sudo mv alpine-chroot-install /usr/local/bin/

# Ensure alpine-chroot-install is available
if ! command -v alpine-chroot-install &>/dev/null; then
  echo "Error: alpine-chroot-install is not installed."
  echo "Install it from: https://github.com/alpinelinux/alpine-chroot-install"
  exit 1
fi


# Run the installation
echo "Installing AlpineLinux $VERSION in $INSTALL_DIR with architecture $ARCH using mirror $MIRROR..."
alpine-chroot-install -d "$INSTALL_DIR" -a "$ARCH" -m "${MIRROR}" -b ${VERSION} -i ${MOUNT_SRC}

# Check if installation was successful
if [ $? -eq 0 ]; then
  echo "Installation completed successfully."
else
  echo "Installation failed."
  exit 1
fi

