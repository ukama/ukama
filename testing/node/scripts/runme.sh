#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2022-present, Ukama Inc.

#/sbin/runme.sh noded "/sbin/noded --p ..."

PID_DIR=/tmp
NAME=$1
CMD=$2

# Execute the command in background
$CMD &
echo $! > ${PID_DIR}/${NAME}.pid

exit 0
