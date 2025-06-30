#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

set -e

# Load helper functions for chroot bind/unbind
source "$(dirname "$0")/chroot-utils.sh"

# Installation parameters (can be overridden via flags)
INSTALL_DIR="$(pwd)/rootfs"
ARCH="x86_64"
VERSION="latest-stable"
MIRROR="http://dl-cdn.alpinelinux.org/alpine"
MOUNT_SRC=""             # optional host path to mount into chroot
MOUNT_DEST="/ukamarepo"  # mount point inside the chroot

# Clean up on exit: either run destroy if available, or unmount binds
trap 'if [ -x "${INSTALL_DIR}/destroy" ]; then
         sudo "${INSTALL_DIR}/destroy"
     else
         unmount_chroot_binds "${INSTALL_DIR}" "${MOUNT_DEST}"
     fi' EXIT

# Parse command-line options
while getopts "a:v:m:i:h" opt; do
    case $opt in
        a) ARCH="$OPTARG" ;;
        v) VERSION="$OPTARG" ;;
        m) MIRROR="$OPTARG" ;;
        i) MOUNT_SRC="$OPTARG" ;;
        h)
            echo "Usage: $0 [-a arch] [-v version] [-m mirror] [-i source]"
            exit 0
            ;;
        *) echo "Invalid option"; exit 1 ;;
    esac
done

# Recreate a fresh rootfs directory
if [ -d "${INSTALL_DIR}" ]; then
    echo "Removing existing ${INSTALL_DIR}..."
    rm -rf "${INSTALL_DIR}"
fi
mkdir -p "${INSTALL_DIR}"

# Ensure alpine-chroot-install utility is installed
if ! command -v alpine-chroot-install &>/dev/null; then
    echo "Installing alpine-chroot-install..."
    wget -O alpine-chroot-install \
         https://raw.githubusercontent.com/alpinelinux/alpine-chroot-install/master/alpine-chroot-install
    chmod +x alpine-chroot-install
    sudo mv alpine-chroot-install /usr/local/bin/
fi

# Bootstrap the base chroot with bash and OpenRC
echo "Bootstrapping Alpine ${VERSION} (${ARCH}) into ${INSTALL_DIR}..."
alpine-chroot-install \
    -d "${INSTALL_DIR}" \
    -a "${ARCH}" \
    -m "${MIRROR}" \
    -b "${VERSION}" \
    -p bash \
    -p openrc \
    -p curl

echo "Base chroot created"

# Mount kernel interfaces (only if not already mounted)
for fs in proc sys dev; do
    TARGET="${INSTALL_DIR}/${fs}"
    mkdir -p "${TARGET}"
    if ! mountpoint -q "${TARGET}"; then
        case "${fs}" in
            proc) mount -t proc   proc   "${TARGET}" ;;
            sys)  mount -t sysfs  sysfs  "${TARGET}" ;;
            dev)  mount --bind    /dev   "${TARGET}" ;;
        esac
    else
        echo "${TARGET} already a mountpoint; skipping"
    fi
done

# Register OpenRC services per the Alpine wiki
chroot "${INSTALL_DIR}" /bin/ash <<'EOF'
rc-update add devfs     sysinit
rc-update add dmesg     sysinit
rc-update add mdev      sysinit

rc-update add hwclock   boot
rc-update add modules   boot
rc-update add sysctl    boot
rc-update add hostname  boot
rc-update add bootmisc  boot
rc-update add syslog    boot

rc-update add mount-ro  shutdown
rc-update add killprocs shutdown
rc-update add savecache shutdown
EOF

echo "OpenRC runlevels configured"

# bind-mount a host directory into the chroot
if [[ -n "${MOUNT_SRC}" ]]; then
    mount_chroot_binds "${INSTALL_DIR}" "${MOUNT_SRC}" "${MOUNT_DEST}"
fi

# Let filesystem catch up
sleep 1
sync

# Invoke build-rootfs.sh inside the chroot
"${INSTALL_DIR}/enter-chroot" /bin/ash -x -c \
                              '/ukamarepo/builder/scripts/build-system/build-rootfs.sh "$@"' -- \
                              "-n" "starterd" \
                              "-c" "/sbin/starter.d" \
                              "-A" "${ARCH}" \
                              "-V" "${VERSION}" \
                              "-M" "${MIRROR}"

if [ $? -eq 0 ]; then
    echo "rootfs created successfully."
else
    echo "rootfs creation failed."
    exit 1
fi
