#!/bin/sh
# Copyright 2017 The Chromium OS Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

# This is a template file which provides settings for firmware update of a
# particular model. The pack_firmware.py script uses this to create a working
# setvars-model.sh script.

# Image and key files for model peppy
IMAGE_MAIN="images/bios_peppy.bin"
IMAGE_EC="images/ec_peppy.bin"
IMAGE_PD=""
SIGNATURE_ID="peppy"
