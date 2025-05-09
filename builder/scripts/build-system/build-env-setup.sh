#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

set -e
set -x

# Set the installation directory as "rootfs" inside the current directory
INSTALL_DIR="$(pwd)/buildenv"

# Default values
ARCH="x86_64"
VERSION="latest-stable"
MIRROR="http://dl-cdn.alpinelinux.org/alpine"
MOUNT_SRC=""   # Source directory to mount
MOUNT_DEST="ukamarepo"  # Destination inside chroot

# Parse command-line arguments
while getopts "a:v:m:i:h" opt; do
  case $opt in
    a) ARCH="$OPTARG" ;;
    v) VERSION="$OPTARG" ;;
    m) MIRROR="$OPTARG" ;;
    i) MOUNT_SRC="$OPTARG" ;;
    h) 
      echo "Usage: $0 [-a arch] [-v version] [-m mirror] [-i source]"
      echo "Example: $0 -a armhf -v v3.21 -m http://dl-cdn.alpinelinux.org/alpine -i /home/user/shared"
      exit 0
      ;;
    *) echo "Invalid option"; exit 1 ;;
  esac
done

echo "Getting alpine-chroot-command."
wget -O alpine-chroot-install https://raw.githubusercontent.com/alpinelinux/alpine-chroot-install/master/alpine-chroot-install

chmod +x alpine-chroot-install
sudo mv alpine-chroot-install /usr/local/bin/

mkdir -p ${INSTALL_DIR}

# Ensure alpine-chroot-install is available
if ! command -v alpine-chroot-install &>/dev/null; then
  echo "Error: alpine-chroot-install is not installed."
  echo "Install it from: https://github.com/alpinelinux/alpine-chroot-install"
  exit 1
fi

# Run the installation
echo "Installing Alpine Linux ${VERSION} in ${INSTALL_DIR} with architecture ${ARCH} using mirror ${MIRROR}."
alpine-chroot-install -d "${INSTALL_DIR}" -a "${ARCH}" -m "${MIRROR}" -b "${VERSION}" -p "bash" 

# Check if installation was successful
if [ $? -eq 0 ]; then
  echo "Installation for build env completed successfully."
else
  echo "Installation for build env failed."
  exit 1
fi

# mount dir
mkdir -p ${INSTALL_DIR}/${MOUNT_DEST}
if [[ -n "{$MOUNT_SRC}" && -n "${MOUNT_DEST}" ]]; then
  echo "Mounting ${MOUNT_SRC} to ${INSTALL_DIR}/${MOUNT_DEST}"
  mount --bind "${MOUNT_SRC}" "${INSTALL_DIR}/${MOUNT_DEST}"
  if [ $? -eq 0 ]; then
  	echo "Mount success"
  else
  	echo "Mount failed"
#	${INSTALL_DIR}/destroy
  	exit 1
  fi
fi

sleep 2;
sync;

# starting build
${INSTALL_DIR}/enter-chroot /bin/ash -c '/ukamarepo/builder/scripts/build-system/build-distro.sh "$@"' "v3.21 /ukamarepo"
if [ $? -eq 0 ]; then
  echo "Build completed successfully."
#  ${INSTALL_DIR}/destroy
else
  echo "Build failed."
#  ${INSTALL_DIR}/destroy
  exit 1
fi

echo "Closing ${0}"
