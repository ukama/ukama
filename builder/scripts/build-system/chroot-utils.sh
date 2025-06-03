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

# Unmount a bind mount from a chroot path
function unmount_chroot_binds() {
    local CHROOT="$1"
    local DST="$2"
    local TARGET="${CHROOT}/${DST}"

    if mountpoint -q "$TARGET"; then
        echo "Unmounting ${TARGET}"
        sudo umount -l "$TARGET"
    fi
}
