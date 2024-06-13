#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

node_id=$1
mocksysfs_root=/ukama/mocksysfs

if [ -z "$node_id" ]; then
    if [ -s "/ukama/nodeid" ]; then
        node_id=$(cat /ukama/nodeid)
    else
        echo "Error: No argument passed and /ukama/nodeid is empty"
        exit 1
    fi
fi

# mock sysfs under /tmp/sys
cd ${mocksysfs_root}
${mocksysfs_root}/utils/prepare_env.sh -u tnode -u anode || exit 1

# Generate schema using dummy schema at '/ukama/mocksysfs/mfgdata/schema'
${mocksysfs_root}/build/genSchema --u ${node_id} --n com --m UK-SA9001-COM-A1-1103  \
                 --f ${mocksysfs_root}/mfgdata/schema/com.json \
                 --n trx --m UK-SA9001-TRX-A1-1103  \
                 --f ${mocksysfs_root}/mfgdata/schema/trx.json \
                 --n mask --m UK-SA9001-MSK-A1-1103 \
                 --f ${mocksysfs_root}/mfgdata/schema/mask.json || exit 1

${mocksysfs_root}/build/genInventory --n com --m UK-SA9001-COM-A1-1103 \
                 --f ${mocksysfs_root}/mfgdata/schema/com.json \
                 -n trx --m UK-SA9001-TRX-A1-1103  \
                 --f ${mocksysfs_root}/mfgdata/schema/trx.json \
                 --n mask -m UK-SA9001-MSK-A1-1103 \
                 --f ${mocksysfs_root}/mfgdata/schema/mask.json  || exit 1

exit 0

