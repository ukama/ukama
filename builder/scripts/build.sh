#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

mock_sysfs_for_noded() {

    repo=$1
    node_id=$2

    mkdir -p /mnt/${node_id}/ukama/mocksysfs/
    cp -p ./mocksysfs.sh /mnt/${node_id}/ukama/mocksysfs/

    cd ${repo}/nodes/ukamaOS/distro/system/noded; make
    cp -rfp * /mnt/${node_id}/ukama/mocksysfs/
    make clean
    cd -

    chroot /mnt/${node_id} /bin/bash <<EOF

    # Create systemd service file
    cat > /etc/systemd/system/mocksysfs.service <<EOL
        [Unit]
        Description=Mock sysfs for Ukama Node
        Before=starterd.service

        [Service]
        Type=oneshot
        ExecStart=/ukama/mocksysfs.sh

        [Install]
        WantedBy=multi-user.target
EOL

    # Check if the service file creation was successful
    if [ ! -f /etc/systemd/system/mocksysfs.service ]; then
        echo "Error: Failed to create systemd service file."
        exit 1
    fi

    # Enable the service
    systemctl enable mocksysfs.service || {
        echo "Error: Failed to enable mocksysfs.service."
        exit 1
    }
EOF
}

install_starter_capp() {

    path=$1

    chroot $path /bin/bash <<EOF

        cd /ukama/apps/pkgs/
        tar zxvf starterd_latest.tar.gz starterd_latest/sbin/starter.d .
        mv starterd_latest/sbin/starter.d /sbin/
        rm -rf starterd_latest/
EOF
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

    ukama_root=$2
    node_id=$3

    if [ -d /mnt/$node_id ]; then
        umount /mnt/${node_id}
        rmdir /mnt/${node_id} || { echo "Unable to remove /mnt/$node_id"; exit 1; }
    fi

    # create bootable os image
    echo "Creating bootsable OS image"
    ./mkimage.sh ${node_id} ${ukama_root} || { echo "Unable to make OS image"; exit 1; }

    # build all apps
    echo "Building all apps"
    ./build-capps.sh ${ukama_root} || exit 1

    # copy the apps and manifest.json into the os image
    echo "Copying apps, installing starter.d and manifesto to the OS image"
    mkdir -p /mnt/${node_id} || exit 1
    mkdir -p /mnt/${node_id}/ukama/ || exit 1
    mount -o loop,offset=$((512*2048)) ${node_id}.img /mnt/${node_id} || exit 1

    cp -r ./pkgs /mnt/${node_id}/ukama/apps/
    cp ${ukama_root}/nodes/manifest.json /mnt/${node_id}/

    # install the starter.d app
    install_starter_app /mnt/${node_id}/

    echo "Copy Ukama sys and vendor libs to the OS image"
    cp ${ukama_root}/nodes/ukamaOS/distro/platform/build/libusys.so \
       /mnt/${node_id}/lib/x86_64-linux-gnu/
    cp -rf ${ukama_root}/nodes/ukamaOS/distro/vendor/build/lib/* \
       /mnt/${node_id}/lib/x86_64-linux-gnu/

    # setup everything needed by node.d
    echo "mocking FS for node.d"
    mock_sysfs_for_noded $ukama_root $node_id

    # update /etc/services to add ports
    echo "Adding all the apps to /etc/services"
    cp ${ukama_root}/nodes/ukamaOS/distro/scripts/files/services \
       /mnt/${node_id}/etc/services

    # umount the image
    umount /mnt/${node_id}
    rmdir  /mnt/${node_id}
    echo "All done"

else
    echo "Invalid argument: $1. Use 'systems' or 'node'."
    exit 1
fi

exit 0
