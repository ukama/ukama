#!/bin/sh

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

# Check if UKAMA_ROOT environment variable is set
if [ -z "$UKAMA_ROOT" ]; then
    echo "UKAMA_ROOT is not set. Please set and try again."
    exit 1
fi

CWD=`pwd`
CONFIG_DIR=${UKAMA_ROOT}/nodes/builder/configs/

cd ${UKAMA_ROOT}/nodes/builder

# builder
make clean; make
rm -rf ${UKAMA_ROOT}/nodes/builder/pkgs

# Loop through each file in the directory
for config in "$CONFIG_DIR"/*.toml; do
    if [ -f "$config" ]; then

        file=$(basename "$config")
        ./builder --create --config "$CONFIG_DIR/$file"
    fi
done

make clean

ls -l ${UKAMA_ROOT}/nodes/builder/pkgs
cd $CWD
