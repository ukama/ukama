#!/bin/bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

set -e

trap "cleanup; exit 1" ERR

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

cleanup() {
    unmount_all_partitions || \
        log "WARNING" "Some partitions might not have been unmounted"
    log "INFO" "Cleanup complete. Exiting."
    exit 1
}

validate_inputs() {
    if [[ -z "${UKAMA_ROOT}" || -z "${OS_TYPE}" || \
              -z "${OS_VERSION}" || -z "${NODE_ID}" ]]; then
        echo "Usage: $0 <UKAMA_ROOT> <OS_TYPE> <OS_VERSION> <NODE_ID> [NODE_APPS]"
        exit 1
    fi
}

cleanup_mount_dir() {
    local dir="$1"

    if [[ -d "${dir}" ]]; then
        if mountpoint -q "${dir}"; then
            log "INFO" "Unmounting existing mount point: ${dir}"
            sudo umount "${dir}" || {
                log "ERROR" "Failed to unmount ${dir}"
                exit 1
            }
        fi
        log "INFO" "Removing directory: ${dir}"
        sudo rm -rf "${dir}" || {
            log "ERROR" "Failed to remove ${dir}"
            exit 1
        }
    fi
}

mount_partition() {

    local name="$1"
    local start_sector="$2"
    local mount_point="${MOUNT_DIR}/${name}"
    local offset=$((start_sector * 512))

    log "INFO" "Mounting ${name} partition at ${mount_point}..."

    # Create the mount point directory
    sudo mkdir -p "${mount_point}"

    # Set up a loop device with offset
    loop_device=$(sudo losetup -f --show -o "${offset}" "${NODE_IMAGE}")
    if [[ -z "${loop_device}" ]]; then
        log "ERROR" "Failed to set up loop device for ${name} partition."
        exit 1
    fi

    # Mount the partition
    sudo mount "${loop_device}" "${mount_point}" || {
        log "ERROR" "Failed to mount ${name} partition."
        sudo losetup -d "${loop_device}" # Clean up loop device
        exit 1
    }

    # Store the loop device for cleanup
    LOOP_DEVICES+=("${loop_device}")
}

unmount_all_partitions() {

    log "INFO" "Unmounting all partitions and detaching loop devices..."

    for part in "${PARTITIONS[@]}"; do
        set -- $part
        local name="$1"
        local mount_point="${MOUNT_DIR}/${name}"

        # Unmount partition
        if mountpoint -q "${mount_point}"; then
            sudo umount "${mount_point}" && log "INFO" "Unmounted ${mount_point}" \
                || log "WARNING" "Failed to unmount ${mount_point}"
        fi
        sudo rmdir "${mount_point}" 2>/dev/null \
            || log "WARNING" "Failed to remove ${mount_point}"
    done

    # Detach all loop devices
    for loop_device in "${LOOP_DEVICES[@]}"; do
        sudo losetup -d "${loop_device}" && log "INFO" "Detached ${loop_device}" \
            || log "WARNING" "Failed to detach ${loop_device}"
    done

    LOOP_DEVICES=() # Reset loop device array
}

install_starter_app() {

    path=$1

    sudo chroot $path /bin/bash <<EOF

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

package_all_apps_via_container_build() {

    local ukama_root="$1"
    local apps="$2"

    log "INFO" "Packaging applications via container build ..."
    cwd=$(pwd)

    cd "${ukama_root}/builder/docker"
    ./apps_setup.sh "alpine" "${ukama_root}" "${apps}"
    cd ${cwd}
}

package_all_apps() {

    local ukama_root="$1"
    local apps="$2"

    log "INFO" "Packaging applications..."
    cwd=$(pwd)
    cd "${ukama_root}/builder"; make clean; make all
    cd "${cwd}"

    IFS=',' read -r -a array <<< "$apps"
    for app in "${array[@]}"; do
        "${ukama_root}/builder/app_builder" \
                     --create \
                     --config "${ukama_root}/builder/configs/${app}.toml"
    done

    cd "${ukama_root}/builder" && make clean
    cd "${cwd}"
}

copy_all_apps_to_image() {

    local ukama_root=$1
    local node_id=$2
    local apps=$3

    echo "Copying apps"

    sudo mkdir -p "${PRIMARY}/ukama/apps/"
    sudo mkdir -p "${PASSIVE}/ukama/apps/"

    IFS=',' read -r -a array <<< "$apps"
    for app in "${array[@]}"; do
        sudo cp "${ukama_root}/build/pkgs/${app}_latest.tar.gz" "${PRIMARY}/ukama/apps/"
        sudo cp "${ukama_root}/build/pkgs/${app}_latest.tar.gz" "${PASSIVE}/ukama/apps/"
    done
}

copy_misc_files_to_image() {

    local os_type=$1
    local ukama_root=$2
    local node_id=$3
    local apps=$4

    generate_manifest_file $apps
    sudo cp ${MANIFEST_FILE} "${PRIMARY}/manifest.json"
    sudo cp ${MANIFEST_FILE} "${PASSIVE}/manifest.json"
    rm ${MANIFEST_FILE}

    # install the starter.d app
    install_starter_app $PRIMARY
    install_starter_app $PASSIVE

    echo "Copy Ukama sys and vendor libs to the OS image"
    sudo mkdir -p "${PRIMARY}/lib/x86_64-linux-gnu/"
    sudo mkdir -p "${PASSIVE}/lib/x86_64-linux-gnu/"

    sudo cp "${ukama_root}/nodes/ukamaOS/distro/platform/build/libusys.so" \
         "${PRIMARY}/lib/x86_64-linux-gnu/"
    sudo cp "${ukama_root}/nodes/ukamaOS/distro/platform/build/libusys.so" \
         "${PASSIVE}/lib/x86_64-linux-gnu/"

    if [[ "${os_type}" = "alpine" ]]; then
        echo "Skipping copying library for alpine based image"
    else
        sudo cp -rf "${ukama_root}/nodes/ukamaOS/distro/vendor/build/lib/"* \
             "${PRIMARY}/lib/x86_64-linux-gnu/"
        sudo cp -rf "${ukama_root}/nodes/ukamaOS/distro/vendor/build/lib/"* \
             "${PASSIVE}/lib/x86_64-linux-gnu/"
    fi

    # update /etc/services to add ports
    echo "Adding all the apps to /etc/services"
    sudo mkdir -p "${PRIMARY}/etc"
    sudo mkdir -p "${PASSIVE}/etc"

    sudo cp "${ukama_root}/nodes/ukamaOS/distro/scripts/files/services" \
         "${PRIMARY}/etc/services"
    sudo cp "${ukama_root}/nodes/ukamaOS/distro/scripts/files/services" \
         "${PASSIVE}/etc/services"
}

# build all the apps within Alpine Docker
# copy everything over to Alpine image
# install all the required pkgs (needed by apps) into the image.
# update alpine so that it can start the starter.d automatically

# Main entry point for the script to build an image for the amplifier node
UKAMA_ROOT="$1"
OS_TYPE="$2"
OS_VERSION="$3"
NODE_ID="$4"
NODE_APPS="$5"
NODE_IMAGE="amplifier_node_sdcard.img"
NODE_TYPE="anode"
MOUNT_DIR="/mnt/${NODE_ID}"
PRIMARY="${MOUNT_DIR}/primary"
PASSIVE="${MOUNT_DIR}/passive"
MANIFEST_FILE=manifest.json
LOOP_DEVICES=()
PARTITIONS=(
    "boot 2048"        # Boot partition
    "primary 100352"   # Primary partition
    "passive 4294656"  # Passive partition
    "unused  8488960"  # Unused partition
)

validate_inputs
cleanup_mount_dir "${PRIMARY}"
cleanup_mount_dir "${PASSIVE}"

# Build base image
log "INFO" "Building image (OS: ${OS_TYPE}) with Node ID: ${NODE_ID}"
if [[ "${OS_TYPE}" = "alpine" ]]; then
    ${UKAMA_ROOT}/builder/scripts/make-alpine-image.sh \
             "${NODE_TYPE}" \
             "${UKAMA_ROOT}" \
             "${OS_VERSION}" \
             "${NODE_IMAGE}"
else
    ${UKAMA_ROOT}/builder/scripts/make-node-ukamaos-image.sh \
                 "${NODE_TYPE}" \
                 "${UKAMA_ROOT}" \
                 "${OS_VERSION}" \
                 "${NODE_IMAGE}"
fi

# Mount paritions of the image.
for part in "${PARTITIONS[@]}"; do
    set -- $part
    mount_partition "$1" "$2"
done
log "INFO" "All partitions mounted successfully under ${MOUNT_BASE_DIR}."

# Compile and package the needed apps for this node
log "INFO" "Packaging and copying applications"
export UKAMA_ROOT

if [[ "${OS_TYPE}" = "alpine" ]]; then
    package_all_apps_via_container_build \
        "${UKAMA_ROOT}" \
        "${NODE_APPS}" || { log "ERROR" "Failed to package apps"; exit 1; }
else
    package_all_apps "${UKAMA_ROOT}" "${NODE_APPS}" \
        || { log "ERROR" "Failed to package apps"; exit 1; }
fi

copy_all_apps_to_image "${UKAMA_ROOT}" "${NODE_ID}" "${NODE_APPS}" \
    || { log "ERROR" "Failed to copy apps to image"; exit 1; }

copy_misc_files_to_image "${OS_TYPE}" "${UKAMA_ROOT}" "${NODE_ID}" "${NODE_APPS}" \
    || { log "ERROR" "Failed to copy misc files"; exit 1; }

# Add Node ID to the image
log "INFO" "Adding Node ID to the image"
echo "${NODE_ID}" | sudo tee "${PRIMARY}/ukama/nodeid" > /dev/null
echo "${NODE_ID}" | sudo tee "${PASSIVE}/ukama/nodeid" > /dev/null

# Unmount the image
log "INFO" "Unmounting image: ${NODE_IMAGE}"
unmount_all_partitions

log "SUCCESS" "Amplifier node image creation completed successfully."
log "INFO"    "Node ID: ${NODE_ID}, Image: ${NODE_IMAGE}"

exit 0
