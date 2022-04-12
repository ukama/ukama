#!/bin/bash

# Copyright (c) 2013 The Chromium OS Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

# Script to remove /boot directory from an image.

# Load common constants.  This should be the first executable line.
# The path to common.sh should be relative to your script's location.
. "$(dirname "$0")/common.sh"

load_shflags

DEFINE_string image "chromiumos_image.bin" \
  "Input file name of Chrome OS image to strip /boot from, or path to rootfs."

# Parse command line.
FLAGS "$@" || exit 1
eval set -- "${FLAGS_ARGV}"

# Abort on error.
set -e

# Swiped/modifed from $SRC/src/scripts/base_library/base_image_util.sh.
zero_free_space() {
  local rootfs="$1"

  info "Zeroing freespace in ${rootfs}"
  sudo fstrim -v "${rootfs}"
}


strip_boot() {
  local image=$1

  local rootfs_dir=$(make_temp_dir)
  if [[ -b "${image}" ]]; then
    enable_rw_mount "${image}"
    sudo mount "${image}" "${rootfs_dir}"
    tag_as_needs_to_be_resigned "${rootfs_dir}"
  else
    # Mount image so we can modify it.
    local loopdev=$(loopback_partscan "${image}")
    mount_loop_image_partition "${loopdev}" 3 "${rootfs_dir}"
  fi

  sudo rm -rf "${rootfs_dir}/boot" &&
    info "/boot directory was removed."

  # To prevent the files we just removed from the FS from remaining as non-
  # zero trash blocks that bloat payload sizes, need to zero them. This was
  # done when the image was built, but needs to be repeated now that we've
  # modified it in a non-trivial way.
  zero_free_space "${rootfs_dir}"
}

IMAGE=$(readlink -f "${FLAGS_image}")
if [[ ! -f "${IMAGE}" && ! -b "${IMAGE}" ]]; then
  IMAGE=
fi
if [[ -z "${IMAGE}" ]]; then
  die "Missing required argument: --from (image to update)"
fi

strip_boot "${IMAGE}"
