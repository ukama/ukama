#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

set -e

source "$(dirname "$0")/chroot-utils.sh"

# Set the installation directory as "rootfs" inside the current directory
INSTALL_DIR="$(pwd)/rootfs"
ARCH="x86_64"
VERSION="latest-stable"
MIRROR="http://dl-cdn.alpinelinux.org/alpine"
MOUNT_SRC=""             # Source directory to mount
MOUNT_DEST="/ukamarepo"  # Destination inside chroot

trap 'if [ -x "${INSTALL_DIR}/destroy" ]; then sudo "${INSTALL_DIR}/destroy"; else unmount_chroot_binds "${INSTALL_DIR}" "${MOUNT_DEST}"; fi' EXIT

# Parse command-line arguments
while getopts "a:v:m:i:h" opt; do
    case $opt in
        a) ARCH="$OPTARG" ;;
        v) VERSION="$OPTARG" ;;
        m) MIRROR="$OPTARG" ;;
        i) MOUNT_SRC="$OPTARG" ;;
        h)
            echo "Usage: $0 [-a arch] [-v version] [-m mirror] [-i source]"
            echo "Example: $0 -a armhf -v v3.17 -m http://dl-cdn.alpinelinux.org/alpine -i /home/user/shared"
            exit 0
            ;;
        *) echo "Invalid option"; exit 1 ;;
    esac
done

# mount detination ignored by script
MOUNT_DEST="ukamarepo"  # Destination inside chroot

if [ -d "${INSTALL_DIR}" ]; then
    echo "Directory exists. Deleting ${INSTALL_DIR}"
    rm -rf "${INSTALL_DIR}"
fi

if ! command -v alpine-chroot-install &>/dev/null; then
    echo "Installing alpine-chroot-install..."
    wget -O alpine-chroot-install https://raw.githubusercontent.com/alpinelinux/alpine-chroot-install/master/alpine-chroot-install
    chmod +x alpine-chroot-install
    sudo mv alpine-chroot-install /usr/local/bin/
fi

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
    echo "rootfs chroot env setup completed "
else
    echo "rootfs chroot env setup failed."
    exit 1
fi

if [[ -n "${MOUNT_SRC}" ]]; then
    mount_chroot_binds "${INSTALL_DIR}" "${MOUNT_SRC}" "${MOUNT_DEST}"
fi

sleep 5;
sync;

# starting build
${INSTALL_DIR}/enter-chroot /bin/ash -c '/ukamarepo/builder/scripts/build-system/build-rootfs.sh "$@"' -- \
              "-p" "active" \
              "-r" "v3.17" \
              "-n" "starterd" \
              "-c" "/sbin/starter.d" \
              "-A" "${ARCH}" \
              "-V" "${VERSION}" \
              "-M" "${MIRROR}"

if [ $? -eq 0 ]; then
    echo "rootfs created successfully."
    #${INSTALL_DIR}/destroy
else
    echo "rootfs creation failed"
    #${INSTALL_DIR}/destroy
    exit 1
fi
