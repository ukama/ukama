#!/bin/bash

# Copyright (c) 2010 The Chromium OS Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

# Remove the test label from lsb-release to prepare an image for
# signing using the official keys.

# Load common constants and variables.
. "$(dirname "$0")/common.sh"

set -e
image=$1

loopdev=$(loopback_partscan "${image}")
rootfs=$(make_temp_dir)
mount_loop_image_partition "${loopdev}" 3 "${rootfs}"
sed -i 's/test//' "${rootfs}/etc/lsb-release"
restore_lsb_selinux "${rootfs}/etc/lsb-release"
