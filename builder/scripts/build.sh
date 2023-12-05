#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

mock_sysfs_for_noded() {

    id = $1

    chroot /mnt/${id} /bin/bash <<'EOL'
        set -e

        apt-get update
        apt-get install -y git

        git clone https://github.com/ukama/ukama.git
        cd ukama/nodes/ukamaOS/distro/system/noded

        ./utils/prepare_env.sh -u tnode -u anode
        ./build/genSchema --u ${id} --n com --m UK-SA9001-COM-A1-1103  \
            --f mfgdata/schema/com.json --n trx --m UK-SA9001-TRX-A1-1103 \
            --f mfgdata/schema/trx.json --n mask --m UK-SA9001-MSK-A1-1103 \
            --f mfgdata/schema/mask.json

        ./build/genInventory --n com --m UK-SA9001-COM-A1-1103 \
            --f mfgdata/schema/com.json -n trx --m UK-SA9001-TRX-A1-1103 \
            --f mfgdata/schema/trx.json --n mask -m UK-SA9001-MSK-A1-1103 \
            --f mfgdata/schema/mask.json
    EOL

    exit
}

if [ "$1" = "system" ]; then

    cd "$2" || exit 1

    find . -type f -name 'go.mod' | while read -r modfile; do
        dir=$(dirname "$modfile")
        cd "$dir" || exit 1
        go mod tidy
        make clean
        make
        cd - || exit 1
    done

    # cleanup and build with using cached layers
    docker-compose down -v || exit 1
    docker image prune -f  || exit 1
    docker-compose build --no-cache || exit 1

elif [ "$1" = "node" ]; then

    ukama_root = $2
    node_id    = $3

    # create bootable os image
    ./mkimage.sh ${node_id}        || exit 1

    # build all the capps
    ./build-capps.sh ${ukama_root} || exit 1

    # copy the apps and manifest.json into the os image
    mkdir -p /mnt/${node_id} || exit 1
    mount -o loop,offset=$((512*2048)) ${IMG_FILE} /mnt/${node_id} || exit 1

    cp -r ./pkgs /mnt/${node_id}/capps/
    cp ${ukama_root}/nodes/manifest.json /mnt/${node_id}

    # setup everything needed by node.d
    mock_sysfs_for_noded $node_id

    # update /etc/services to add ports

    # umount the image
    umount /mnt/${node_id}
    rmdir  /mnt/${node_id}

else
    echo "Invalid argument: $1. Use 'systems' or 'node'."
    exit 1
fi

exit 0
