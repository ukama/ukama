#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

partition_name() {
    local dev="$1"
    local num="$2"
    if [[ "$dev" =~ [0-9]$ ]]; then
        echo "${dev}p${num}"
    else
        echo "${dev}${num}"
    fi
}

partition_unmount_all() {
    local dev="$1"
    local part
    for part in $(lsblk -nrpo NAME,MOUNTPOINT "$dev" 2>/dev/null | awk 'NF>1 && $2!=""{print $1}'); do
        sudo umount "$part" 2>/dev/null || true
    done
    sudo partprobe "$dev" 2>/dev/null || true
    sleep 1
}

partition_wipe() {
    local dev="$1"
    sudo sgdisk --clear "$dev" 2>/dev/null || true
    sudo dd if=/dev/zero of="$dev" bs=1M count=8 2>/dev/null || true
}

partition_apply() {
    local config="$1"
    local dev="$2"

    sudo parted -s "$dev" mklabel msdos

    local count
    count=$(yq_count "$config" partitions)

    local i
    for ((i=0; i<count; i++)); do
        local type fs start end
        type=$(yq_read "$config" "partitions[$i].type")
        fs=$(yq_read "$config" "partitions[$i].fs")
        start=$(yq_read "$config" "partitions[$i].start")
        end=$(yq_read "$config" "partitions[$i].end")

        if [ "$type" = "extended" ]; then
            sudo parted -s "$dev" mkpart extended "$start" "$end"
        elif [ -n "$fs" ] && [ "$fs" != "null" ]; then
            sudo parted -s "$dev" mkpart "$type" "$fs" "$start" "$end"
        else
            sudo parted -s "$dev" mkpart "$type" "$start" "$end"
        fi
    done

    sync
    sleep 2
    sudo partprobe "$dev" 2>/dev/null || sudo blockdev --rereadpt "$dev" || true
    sleep 2
}

partition_format() {
    local config="$1"
    local dev="$2"

    local count
    count=$(yq_count "$config" partitions)

    local i
    for ((i=0; i<count; i++)); do
        local num type fs label
        num=$(yq_read "$config" "partitions[$i].num")
        type=$(yq_read "$config" "partitions[$i].type")
        fs=$(yq_read "$config" "partitions[$i].fs")
        label=$(yq_read "$config" "partitions[$i].label")

        [ "$type" = "extended" ] && continue
        [ "$fs" = "null" ] || [ -z "$fs" ] && continue

        local part_name
        part_name=$(partition_name "$dev" "$num")

        case "$fs" in
            fat32)
                sudo mkfs.vfat -F 32 -n "$label" "$part_name" >/dev/null
                ;;
            ext4)
                sudo mkfs.ext4 -F -L "$label" "$part_name" >/dev/null
                ;;
            linux-swap)
                sudo mkswap -L "$label" "$part_name" >/dev/null
                ;;
        esac
    done
    sync
}
