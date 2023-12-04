#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

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
    mkdir -p /mnt/${NODE_ID} || exit 1
    mount -o loop,offset=$((512*2048)) ${IMG_FILE} /mnt/${NODE_ID} || exit 1

    cp -r ./pkgs /mnt/${NODE_ID}/capps/
    cp ${ukama_root}/nodes/manifest.json /mnt/${NODE_ID}

    # modify systemd config to start starter.d

    # umount the image
    umount /mnt/${NODE_ID}
    rmdir  /mnt/${NODE_ID}

else
    echo "Invalid argument: $1. Use 'systems' or 'node'."
    exit 1
fi

exit 0
