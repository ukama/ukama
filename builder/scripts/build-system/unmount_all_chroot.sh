#/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

# Script to build and package ukamaOS app

SCRIPT_DIR=$1
cat /proc/mounts | cut -d' ' -f2 | grep "^$SCRIPT_DIR." | sort -r | while read path; do
			echo "Unmounting $path" >&2
			sudo umount -fn "$path" || exit 1
		done
