#!/bin/bash

# Copyright (c) 2010 The Chromium OS Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

# abort on error
set -e

# Load common constants and variables.
. "$(dirname "$0")/common.sh"

main() {
  if [[ $# -ne 1 ]]; then
    echo "Usage $0 <image>"
    exit 1
  fi

  local image="$1"

  local loopdev rootfs
  if [[ -d "${image}" ]]; then
    rootfs="${image}"
  else
    rootfs=$(make_temp_dir)
    loopdev=$(loopback_partscan "${image}")
    mount_loop_image_partition_ro "${loopdev}" 3 "${rootfs}"
  fi

  if ! no_chronos_password "${rootfs}"; then
    die "chronos password is set! Shouldn't be for release builds."
  fi
}
main "$@"
