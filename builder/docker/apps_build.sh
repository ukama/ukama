#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

# To run within specific container

set -e
set -x

export UKAMA_ROOT=/workspace
export ALPINE_BUILD=1

#default branch is main
APPS=$1

cd "${UKAMA_ROOT}/builder" && make clean && make ALPINE_BUILD=1 all

IFS=',' read -r -a array <<< "${APPS}"
for app in "${array[@]}"; do
    "${UKAMA_ROOT}/builder/app_builder" \
        --create \
        --config "${UKAMA_ROOT}/builder/configs/${app}.toml"
done

cd "${UKAMA_ROOT}/builder" && make clean

exit 0
