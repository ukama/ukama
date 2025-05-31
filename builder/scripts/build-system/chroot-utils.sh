#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

# Ukama Chroot Helper Utilities

# Mount a source directory into a chroot bind mount
function mount_chroot_binds() {
    local CHROOT="$1"
    local SRC="$2"
    local DST="$3"

    mkdir -p "${CHROOT}/${DST}"
    echo "Mounting ${SRC} -> ${CHROOT}/${DST}"
    sudo mount --bind "${SRC}" "${CHROOT}/${DST}"
}

# Unmount chroot mounts safely (custom and system)
function unmount_chroot_binds() {
    local CHROOT="$1"
    local CUSTOM_DST="$2"

    echo "Unmounting chroot mounts under ${CHROOT}..."

    # Unmount in reverse order
    for mnt in \
        "$CHROOT/dev/pts" \
        "$CHROOT/dev/shm" \
        "$CHROOT/dev/mqueue" \
        "$CHROOT/dev/hugepages" \
        "$CHROOT/dev" \
        "$CHROOT/proc" \
        "$CHROOT/sys/kernel/config" \
        "$CHROOT/sys/fs/fuse/connections" \
        "$CHROOT/sys" \
        "$CHROOT/run" \
        "$CHROOT/$CUSTOM_DST"
    do
        if mountpoint -q "$mnt"; then
            echo "Unmounting $mnt"
            sudo umount -lf "$mnt"
        fi
    done
}
