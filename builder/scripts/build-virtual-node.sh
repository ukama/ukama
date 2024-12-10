#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

MANIFEST_FILE=manifest.json

log() {
    local type="$1"
    local message="$2"
    local color

    case "$type" in
        "INFO") color="\033[1;34m";;
        "SUCCESS") color="\033[1;32m";;
        "ERROR") color="\033[1;31m";;
        *) color="\033[1;37m";;
    esac
    echo -e "${color}${type}: ${message}\033[0m"
}

check_status() {

    if [ $1 -ne 0 ]; then
        log "ERROR" "Script failed at stage: $3"
        exit 1
    fi
    log "SUCCESS" "$2"
}

mock_sysfs_for_noded() {

    repo=$1
    node_id=$2

    mkdir /mnt/${node_id}/ukama/mocksysfs/
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
        ExecStart=/ukama/mocksysfs/mocksysfs.sh

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

install_starter_app() {

    path=$1

    chroot $path /bin/bash <<EOF

        cd /ukama/apps/pkgs/
        tar zxvf starterd_latest.tar.gz starterd_latest/sbin/starter.d .
        mv starterd_latest/sbin/starter.d /sbin/
        rm -rf starterd_latest/
EOF
}

generate_manifest_file() {

    local apps_to_include="$1"

    # Create an array from the comma-separated list
    IFS=',' read -r -a apps_array <<< "$apps_to_include"

    cat <<EOF > ${MANIFEST_FILE}
{   
    "version": "0.1",
    
    "spaces" : [
        { "name" : "boot" },   
        { "name" : "services" },
        { "name" : "reboot" }  
    ],

    "capps": [
        {
            "name"   : "noded",
            "tag"    : "latest",
            "restart": "yes",
            "space"  : "boot"
        },
        {
            "name"   : "bootstrap",
            "tag"    : "latest",
            "restart": "yes",
            "space"  : "boot",
	        "depends_on" : [
                {
                    "capp"  : "noded",
			        "state" : "active"
		        }
	        ]
        },
        {
            "name"   : "meshd",
            "tag"    : "latest",
            "restart": "yes",
            "space"  : "boot",
	        "depends_on" : [
                {
                    "capp"  : "bootstrap",
			        "state" : "done"
		        }
	        ]
        }
EOF

  echo '        ,' >> ${MANIFEST_FILE}
  echo '        {"name" : "services", "capps" : [' >> ${MANIFEST_FILE}

  for app in "${apps_array[@]}"; do
    case "$app" in
      "wimcd"|"configd"|"metricsd"|"lookoutd"|"deviced"|"notifyd")
        cat <<EOF >> ${MANIFSET_FILE}
        {
            "name"   : "$app",
            "tag"    : "latest",
            "restart": "yes",
            "space"  : "services"
        },
EOF
        ;;
    esac
  done

  # Remove the last comma and close the JSON array
  sed -i '$ s/,$//' ${MANIFEST_FILE}
  echo '    ]}'  >> ${MANIFEST_FILE}
  echo '}'       >> ${MANIFEST_FILE}
}

build_base_image() {

    ukama_root=$1
    base_id=$2

    # create bootable os image (based on Ubuntu)
    echo "Creating bootsable OS image for virtual env."
    ./mk_virtual_os_image.sh ${base_id} ${ukama_root} || \
        { echo "Unable to make OS image"; exit 1; }

}

build_base_apps() {

    local ukama_root=$1
    local base_id=$2
    local apps=$3
    
    # build all apps
    echo "Building all apps"
    ./build-capps.sh ${ukama_root} || exit 1

    # copy the apps and manifest.json into the os image
    echo "Copying apps, installing starter.d and manifesto to the OS image"
    mkdir -p /mnt/${base_id} || exit 1
    mount -o loop,offset=$((512*2048)) ${base_id}.img /mnt/${base_id} || exit 1

    # Only copy the apps needed and generate the manifest file.
    IFS=',' read -r -a array <<< "$apps"
    for app in "${array[@]}"; do
        mv ./pkgs/${app}_latest.tar.gz /mnt/${base_id}/ukama/apps/
    done

    generate_manifest_file $apps
    
    cp ${MANIFEST_FILE} /mnt/${base_id}/manifest.json

    # install the starter.d app
    install_starter_app /mnt/${base_id}/

    echo "Copy Ukama sys and vendor libs to the OS image"
    cp ${ukama_root}/nodes/ukamaOS/distro/platform/build/libusys.so \
       /mnt/${base_id}/lib/x86_64-linux-gnu/
    cp -rf ${ukama_root}/nodes/ukamaOS/distro/vendor/build/lib/* \
       /mnt/${base_id}/lib/x86_64-linux-gnu/

    # setup everything needed by node.d
    echo "mocking FS for node.d"
    mock_sysfs_for_noded $ukama_root $base_id

    # update /etc/services to add ports
    echo "Adding all the apps to /etc/services"
    cp ${ukama_root}/nodes/ukamaOS/distro/scripts/files/services \
       /mnt/${base_id}/etc/services
}

build_node_from_base_image() {

    ukama_root=$1
    node_id=$2
    based_id=$3

    cp ${base_id}.img ${node_id}.img

    mkdir -p /mnt/${node_id} || exit 1
    mount -o loop,offset=$((512*2048)) ${node_id}.img /mnt/${node_id} || exit 1

    rm /mnt/${node_id}/ukama/nodeid
    echo $node_id > /mnt/${node_id}/ukama/nodeid

    rm -rf /mnt/${node_id}/ukama/mocksysfs/
    mock_sysfs_for_noded $ukama_root $node_id

    umount /mnt/${node_id}
    rmdir  /mnt/${node_id}
    echo "All done"
}

# Main entry point for the script

if [ "$1" = "base-image" ]; then

    ukama_root=$2
    base_id=$3
    apps=$4

    # If already mounted, clean up
    if [ -d /mnt/$base_id ]; then
        umount /mnt/${base_id}
        rmdir /mnt/${base_id} || { echo "Unable to remove /mnt/$base_id"; exit 1; }
    fi
    
    echo "Building base image with base_id: $base_id"
    build_base_image $ukama_root $base_id
    build_base_apps $apps

    # umount the image
    umount /mnt/${base_id}
    rmdir  /mnt/${base_id}

    echo "All done"
    
elif [ "$1" = "create-node" ]; then

    ukama_root=$2
    node_id=$3
    base_id=$4

    echo "Building node image with base: $base_id and node: $node_id"
    build_node_from_base_image $ukama_root $node_id $base_id

else
    echo "Invalid argument: $0. Use 'base-image' or 'create-node'"
    exit 1
fi

exit 0
