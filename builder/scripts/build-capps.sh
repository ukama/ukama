#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

export UKAMA_ROOT=$1
CONFIG_DIR=${UKAMA_ROOT}/nodes/builder/configs/

CWD=`pwd`
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

# cleanup
make clean

cp -r ${UKAMA_ROOT}/nodes/builder/pkgs $CWD
cd $CWD
