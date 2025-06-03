#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

set -e

export UKAMA_ROOT=$1
CONFIG_DIR=${UKAMA_ROOT}/builder/configs

CWD=`pwd`
cd ${UKAMA_ROOT}/builder

# make the app_builder and cleanup
make clean; make app_builder
rm -rf ${UKAMA_ROOT}/builder/pkgs

# Loop through each file in the directory
for config in "$CONFIG_DIR"/*.toml; do
    if [ -f "$config" ]; then
        file=$(basename "$config")
        ./app_builder --create --config "$CONFIG_DIR/$file"
    fi
done

# cleanup
make clean

cp -r ${UKAMA_ROOT}/build/pkgs $CWD
cd $CWD
