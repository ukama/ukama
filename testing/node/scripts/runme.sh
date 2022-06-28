#!/bin/bash
# Copyright (c) 2022-present, Ukama Inc.
# All rights reserved.

#/sbin/runme.sh noded "/sbin/noded --p ..."

PID_DIR=/tmp
NAME=$1
CMD=$2

# Execute the command in background
$CMD &
echo $! > ${PID_DIR}/${NAME}.pid

exit 0
