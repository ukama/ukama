#!/bin/sh
# Copyright 2017 The Chromium OS Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

# This is a template file which provides settings for firmware update of a
# particular model. The pack_firmware.py script uses this to create a working
# setvars-model.sh script.

# Version information for model whitetip
TARGET_RO_FWID="Google_Coral.10068.45.0"
TARGET_FWID="Google_Coral.10068.45.0"
TARGET_ECID="coral_v1.1.7272-0b44fba22"
TARGET_PDID=""
TARGET_PLATFORM="Google_Coral"

# Image and key files for model whitetip
IMAGE_MAIN="images/bios_coral.bin"
IMAGE_EC=""
IMAGE_PD=""
SIGNATURE_ID="sig-id-in-customization-id"
