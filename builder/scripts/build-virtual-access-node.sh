#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

set -e
set -x

# Variables
LOG_FILE="build_run.log"
IMG_NAME="access-node.img"
RPI_BASE_URL="https://downloads.raspberrypi.com"
RPI_URL_PATH="raspios_lite_arm64/images/raspios_lite_arm64-2024-11-19"
RPI_IMG="2024-11-19-raspios-bookworm-arm64-lite.img"
RPI_IMG_OFFSET="541065216"
PRIMARY_MOUNT="/mnt/access-node"

# for 32-bit ARM
#RPI_URL_PATH="raspios_oldstable_armhf/images/raspios_oldstable_armhf-2024-10-28"
#RPI_IMG="2024-10-22-raspios-bullseye-armhf.img"
#RPI_IMG_OFFSET="272629248"

TMP_DIR="/tmp/access-node"
TMP_ROOTFS="${TMP_DIR}/alpine-rootfs"
TMP_LINUX="${TMP_DIR}/linux"

USER_NAME="ukama"

CWD=$(pwd)

trap cleanup EXIT

function check_sudo() {
    if ! sudo -v; then
        echo "You do not have sudo privileges or sudo is not configured correctly."
        exit 1
    fi
}

function log() {
    local type="$1"
    local message="$2"
    local timestamp
    local file_name
    local func_name
    local color
    local reset="\033[0m"

    timestamp=$(date +"%Y-%m-%d %H:%M:%S")
    file_name=$(basename "${BASH_SOURCE[1]}")
    func_name="${FUNCNAME[1]}"

    # Set color based on log type
    case "$type" in
        INFO)
            color="\033[1;34m" # Blue
            ;;
        SUCCESS)
            color="\033[1;32m" # Green
            ;;
        WARNING)
            color="\033[1;33m" # Yellow
            ;;
        ERROR)
            color="\033[1;31m" # Red
            ;;
        *)
            color="$reset" # Default (no color)
            ;;
    esac

    printf "%s %b%s%b %s:%s \"%s\"\n" "$timestamp" "$color" "$type" "$reset" "$file_name" "$func_name" "$message" | tee -a "$LOG_FILE"
}

function LOG_EXEC() {
    log "EXEC" "$*"
    "$@" >>"$LOG_FILE" 2>&1
    if [[ $? -ne 0 ]]; then
        log "ERROR" "Command failed: $*"
        exit 1
    fi
}

function check_command() {
    command -v "$1" >/dev/null 2>&1 || {
        log "ERROR" "Command '$1' not found. Please install it."
        exit 1
    }
}

function check_requirements() {
    log "INFO" "Checking required commands..."
    for cmd in git make parted rsync wget tar dd losetup mkfs.vfat mkfs.ext4; do
        check_command "$cmd"
    done
    log "SUCCESS" "All required commands are available."
}

function cleanup() {
    log "INFO" "Cleaning up resources..."
    sudo umount -R "${PRIMARY_MOUNT}" 2>/dev/null || true
#    sudo rm -rf "$TMP_DIR"
    log "INFO" "Cleanup completed."
}

function build_linux_kernel() {
    log "INFO" "Building linux kernel..."

    if [ -d "${TMP_LINUX}" ]; then
        log "INFO" "Using existing linux kernel at: ${TMP_LINUX}"
    else
        wget https://cdn.kernel.org/pub/linux/kernel/v6.x/linux-6.1.34.tar.xz
        tar xJf linux-6.1.34.tar.xz
        mv linux-6.1.34 "${TMP_LINUX}"
    fi

    cd "${TMP_LINUX}"

    if [ -f "${TMP_LINUX}/arch/arm64/boot/Image" ]; then
        log "INFO" "Kernel image already exists, skipping"
    else
        # build linux kernel suitable for qemu
        ARCH=arm64 CROSS_COMPILE=/bin/aarch64-linux-gnu- make defconfig
        ARCH=arm64 CROSS_COMPILE=/bin/aarch64-linux-gnu- make kvm_guest.config
        ARCH=arm64 CROSS_COMPILE=/bin/aarch64-linux-gnu- make -j8
    fi

    cd "${TMP_DIR}"
    log "SUCCESS" "Linux kernel build completed."
}

function build_apps_using_target() {

    # Idea is to run QEMU so we can build everything on the target
    # rather than doing docker. This avoid lib mismatch.
    #
    # Steps:
    # 1. copy repo to workspace.
    # 2. mount the workspace.
    # 3. update system.d to run the following script:
    # 3. after running QEMU:
    #    install all pkgs (look into dockerfile)
    #    compile all pkgs (run the scripts under docker)
    #    copy the pkgs to right location (copy_all_apps)
    #    exit the qemu
    # 4. remove the above script from the system.d
    # 5. rest of stuff as before
   

}

function build_apps_using_container() {
    local ukama_root="$1"
    local apps="$2"

    log "INFO" "Packaging applications via container build ..."
    cwd=$(pwd)

    cd "${ukama_root}/builder/docker"
    ./apps_setup.sh "access" "arm64" "${ukama_root}" "${apps}"
    cd ${cwd}
}

function download_rpi_rootfs() {
    log "INFO" "Checking for RPI rootfs image..."

    if [ -f "${RPI_IMG}" ]; then
        log "INFO" "Using existing extracted image: ${RPI_IMG}"
        return
    fi

    if [ -f "${RPI_IMG}.xz" ]; then
        log "INFO" "Using existing compressed image: ${RPI_IMG}.xz"
    else
        log "INFO" "Downloading RPI rootfs image..."
        wget "${RPI_BASE_URL}/${RPI_URL_PATH}/${RPI_IMG}.xz" \
            || { log "ERROR" "Unable to download ${RPI_IMG}.xz"; exit 1; }
    fi

    log "INFO" "Extracting ${RPI_IMG}.xz..."
    xz -d -f "${RPI_IMG}.xz" \
        || { log "ERROR" "Unable to extract ${RPI_IMG}.xz"; exit 1; }

    log "SUCCESS" "RPI rootfs image downloaded and extracted: ${RPI_IMG}"
}

function install_starter_app() {

    path=$1

    log "INFO" "Installing starter.d"

    sudo chroot "$path" /bin/sh <<'EOF'
cd /ukama/apps/pkgs/
tar zxvf starterd_latest.tar.gz
cp starterd_latest/sbin/starter.d /sbin/
rm -rf starterd_latest/
EOF
}

function copy_linux_kernel() {
    log "INFO" "Copying linux kernel..."
    cp "${TMP_LINUX}/arch/arm64/boot/Image" "${CWD}/build_access_node/kernel.img"
    log "SUCCESS" "Linux kernel copied"
}

function copy_all_apps() {
    local ukama_root=$1
    local apps=$2

    log "INFO" "Copying apps"

    sudo mkdir -p "${PRIMARY_MOUNT}/ukama/apps/pkgs"
    IFS=',' read -r -a array <<< "$apps"
    for app in "${array[@]}"; do
        sudo cp "${ukama_root}/build/pkgs/${app}_latest.tar.gz" \
             "${PRIMARY_MOUNT}/ukama/apps/pkgs"
    done

    sudo rm -rf "${ukama_root}/build/"
}

function copy_misc_files() {
    local ukama_root=$1
    local apps=$2

    log "INFO" "Copying various files to image"

    create_manifest_file $apps
    sudo cp ${MANIFEST_FILE} "${PRIMARY_MOUNT}/manifest.json"
    rm ${MANIFEST_FILE}

    # install the starter.d app
    install_starter_app "${PRIMARY_MOUNT}"

    log "INFO" "Copy Ukama sys lib to the image"
    sudo mkdir -p "${PRIMARY_MOUNT}/lib/aarch64-linux-gnu/"
    sudo cp "${ukama_root}/nodes/ukamaOS/distro/platform/build/libusys.so" \
         "${PRIMARY_MOUNT}/lib/aarch64-linux-gnu/"

    # update /etc/services to add ports
    log "INFO" "Adding all the apps to /etc/services"
    sudo mkdir -p "${PRIMARY_MOUNT}/etc"
    sudo cp "${ukama_root}/nodes/ukamaOS/distro/scripts/files/services" \
         "${PRIMARY_MOUNT}/etc/services"
}

resize_disk_image() {
    log "INFO" "Resizing the disk image: ${RPI_IMG}"

    truncate -s "8G" "${RPI_IMG}"

    # Use fdisk to modify the partition table
    echo -e "d\n2\nn\np\n2\n1056768\n\nw\n" | fdisk "${RPI_IMG}"

    loop_device=$(sudo losetup -Pf --show "${RPI_IMG}")
    if [[ -z "$loop_device" ]]; then
        log "ERROR" "Failed to map the image to a loop device."
    fi
    sudo resize2fs "${loop_device}p2"
    sudo losetup -d "$loop_device"

    log "SUCCESS" "disk image '${RPI_IMG}' has been resized and parition extended."
}

function install_required_libs() {
    log "INFO" "Installing library"

        # Bind mount necessary system directories
    log "INFO" "Binding system directories for chroot"
    sudo mount --bind /dev "$PRIMARY_MOUNT/dev"
    sudo mount --bind /dev/pts "$PRIMARY_MOUNT/dev/pts"
    sudo mount --bind /proc "$PRIMARY_MOUNT/proc"
    sudo mount --bind /sys "$PRIMARY_MOUNT/sys"
    sudo mount --bind /run "$PRIMARY_MOUNT/run"

    # Chroot into the target and configure the systemd service
    log "INFO" "Chrooting into target rootfs to install libs and pkgs"
    sudo chroot "$PRIMARY_MOUNT" /bin/bash <<'EOF'
set -e

# Configure locales
echo "en_US.UTF-8 UTF-8" > /etc/locale.gen
locale-gen
dpkg-reconfigure --frontend=noninteractive locales

apt-get clean && rm -rf /var/lib/apt/lists/*

apt-get update && apt-get install -y \
    software-properties-common \
    build-essential \
    git \
    wget \
    autoconf \
    automake \
    libtool \
    pkg-config \
    libssl-dev \
    texinfo \
    cmake \
    tcl \
    zlib1g-dev \
    texlive \
    texlive-latex-extra \
    ghostscript \
    gperf \
    gtk-doc-tools \
    libev-dev \
    bison \
    jq \
    util-linux \
    libgnutls28-dev \
    libmicrohttpd-dev \
    libcurl4-openssl-dev \
    libjansson-dev \
    gcc-aarch64-linux-gnu \
    g++-aarch64-linux-gnu \
    uuid-dev \
    libsqlite3-dev

apt-get clean && rm -rf /var/lib/apt/lists/*

#build and install libs from source
git clone https://github.com/babelouest/orcania.git && \
    cd orcania && \
    mkdir build && \
    cd build && \
    cmake .. && \
    make && \
    make install && \
    cd ../../ && rm -rf orcania

# Clone and build Ulfius with Yder disabled
git clone https://github.com/babelouest/ulfius.git && \
    cd ulfius && \
    mkdir build && \
    cd build && \
    cmake -DWITH_YDER=off -DWITH_JOURNALD=off -DCMAKE_C_FLAGS="-Wno-stringop-overflow" .. && \
    make YDERFLAG=1 && \
    make install && \
    cd ../../ && rm -rf ulfius

# Install tomlc99
git clone https://github.com/cktan/tomlc99.git && \
    cd tomlc99 && \
    make && \
    cp toml.h /usr/include/ && \
    cp libtoml.* /usr/lib/ && \
    cd .. && rm -rf tomlc99

# Clone Prometheus client repository and build it
git clone https://github.com/ukama/prometheus-client.git && \
    cd prometheus-client && \
    mkdir -p prom/build && \
    cd prom/build && \
    cmake ../ -DCMAKE_C_COMPILER=gcc -DCMAKE_C_FLAGS="-I/usr/include" \
        -DCMAKE_LD_FLAGS="-L/usr/lib" && \
    make install DESTDIR=/usr && \
    mkdir -p ../../promhttp/build && \
    cd ../../promhttp/build && \
    cmake ../ -DCMAKE_C_COMPILER=gcc -DCMAKE_C_FLAGS="-I/usr/include" \
        -DCMAKE_LD_FLAGS="-L/usr/lib" && \
    make install DESTDIR=/usr && \
    cd ../../.. && \
    rm -rf prometheus-client

# Copy libs to default location to avoid seeting LD_LIBRARY_PATH
cp -rf /usr/local/lib/* /lib/aarch64-linux-gnu/

EOF
    
    # Unmount system directories
    log "INFO" "Unmounting system directories"
    sudo umount "$PRIMARY_MOUNT/dev/pts"
    sudo umount "$PRIMARY_MOUNT/dev"
    sudo umount "$PRIMARY_MOUNT/proc"
    sudo umount "$PRIMARY_MOUNT/sys"
    sudo umount "$PRIMARY_MOUNT/run"

    log "SUCCESS" "Requried libraries and pkgs installed"
}

function create_manifest_file() {
    local apps_to_include="$1"
    log "INFO" "Creating manifest file"

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
        cat <<EOF >> ${MANIFEST_FILE}
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

  echo '        ,' >> ${MANIFEST_FILE}
  echo '        {"name" : "services", "capps" : [' >> ${MANIFEST_FILE}

  for app in "${apps_array[@]}"; do
    case "$app" in
      "wimcd"|"configd"|"metricsd"|"lookoutd"|"deviced"|"notifyd")
        cat <<EOF >> ${MANIFEST_FILE}
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

function mount_partitions() {
    log "INFO" "Mounting partitions..."

    sudo mkdir -p "$PRIMARY_MOUNT"
    sudo mount -o loop,offset="${RPI_IMG_OFFSET}" "${RPI_IMG}" "${PRIMARY_MOUNT}" || \
        { log "ERROR" "Unable to mount rootfs partition"; exit 1; }
       
    log "SUCCESS" "Partitions mounted."
}

function setup_ukama_dirs() {
    log "INFO" "Creating Ukama directories..."

    mkdir -p "${PRIMARY_MOUNT}/ukama"
    mkdir -p "${PRIMARY_MOUNT}/ukama/configs"
    mkdir -p "${PRIMARY_MOUNT}/ukama/apps"
    mkdir -p "${PRIMARY_MOUNT}/ukama/apps/pkgs"
    mkdir -p "${PRIMARY_MOUNT}/ukama/apps/rootfs"
    mkdir -p "${PRIMARY_MOUNT}/ukama/apps/registry"

    echo "${NODE_ID}" > "${PRIMARY_MOUNT}/ukama/nodeid"
    echo "localhost"  > "${PRIMARY_MOUNT}/ukama/bootstrap"

    touch "${PRIMARY_MOUNT}/ukama/apps.log"

    log "SUCCESS" "Ukama directories created."
}

function configure_systemd_service() {
    log "INFO" "Configuring systemd service for starter.d in target rootfs"

    # Ensure PRIMARY_MOUNT path exists
    if [ ! -d "$PRIMARY_MOUNT" ]; then
        log "ERROR" "PRIMARY_MOUNT path does not exist: $PRIMARY_MOUNT"
        exit 1
    fi

    # Create the systemd service file in the target root filesystem
    SERVICE_FILE="$PRIMARY_MOUNT/etc/systemd/system/starter.service"
    log "INFO" "Creating systemd service file at $SERVICE_FILE"

    sudo bash -c "cat <<'EOF' > $SERVICE_FILE
[Unit]
Description=Starter service for running starter.d
After=network.target

[Service]
ExecStart=/sbin/starter.d
Restart=always
User=root
PIDFile=/var/run/starter.pid

[Install]
WantedBy=multi-user.target
EOF"

    # Bind mount necessary system directories
    log "INFO" "Binding system directories for chroot"
    sudo mount --bind /dev "$PRIMARY_MOUNT/dev"
    sudo mount --bind /proc "$PRIMARY_MOUNT/proc"
    sudo mount --bind /sys "$PRIMARY_MOUNT/sys"

    # Chroot into the target and configure the systemd service
    log "INFO" "Chrooting into target rootfs to enable the systemd service"
    sudo chroot "$PRIMARY_MOUNT" /bin/bash <<'EOF'
# Set locale to avoid warnings
export LANGUAGE=en_US.UTF-8
export LC_ALL=en_US.UTF-8
export LANG=en_US.UTF-8
locale-gen en_US.UTF-8 || true

# Reload systemd and enable the service
systemctl daemon-reload
systemctl enable starter.service

# Remove any SysV init script association (if necessary)
rm -f /etc/init.d/starter || true
EOF

    # Unmount system directories
    log "INFO" "Unmounting system directories"
    sudo umount "$PRIMARY_MOUNT/dev"
    sudo umount "$PRIMARY_MOUNT/proc"
    sudo umount "$PRIMARY_MOUNT/sys"

    log "SUCCESS" "Systemd service for starter.d configured in target rootfs"
}

function unmount_partitions() {
    log "INFO" "Unmounting partitions..."
    sudo umount "$PRIMARY_MOUNT"
    log "SUCCESS" "Partitions unmounted."
}

function create_ssh_user() {
    log "INFO" "Adding ssh user..."

    echo "${USER_NAME}:x:1001:1001::/home/${USER_NAME}:/bin/bash" \
        | sudo tee -a "${PRIMARY_MOUNT}/etc/passwd"
    echo "${USER_NAME}:x:1001:" \
        | sudo tee -a "${PRIMARY_MOUNT}/etc/group"
    echo "${USER_NAME}::19000:0:99999:7:::" \
        | sudo tee -a "${PRIMARY_MOUNT}/etc/shadow"

    # Create home directory
    mkdir -p        "${PRIMARY_MOUNT}/home/${USER_NAME}"
    chown 1001:1001 "${PRIMARY_MOUNT}/home/${USER_NAME}"
    chmod 700       "${PRIMARY_MOUNT}/home/${USER_NAME}"

    # Add to sudoer
    echo "ukama ALL=(ALL:ALL) ALL" | sudo tee "${PRIMARY_MOUNT}/etc/sudoers.d/ukama"
    sudo chmod 0440 "${PRIMARY_MOUNT}/etc/sudoers.d/ukama"

    log "SUCCESS" "User $USER_NAME added with no password."
}

function pre_cleanup_and_dir_setup() {

    local image_name=$1
    local tmp_dir=$2
    local build_dir=$3

    if [ -f "$image_name" ]; then
        rm "$image_name"
    fi

    if [ -d "$build_dir" ]; then
        rm -rf "$build_dir"
    fi

    mkdir -p "$build_dir"
    mkdir -p "$tmp_dir"
}

# Main Script Execution
OS_TYPE="alpine"
OS_VERSION="0.0.1"
MANIFEST_FILE="manifest.json"
export TARGET="access"

rm -rf ${LOG_FILE}

if [[ $# -ne 3 ]]; then
    log "ERROR" "Error: Exactly 3 arguments are required!"
    log "INFO"  "Usage: $0 <ukama_root> <node_apps> <node_id>"
    exit 1
fi

UKAMA_ROOT=$1
NODE_APPS=$2
NODE_ID=$3

check_sudo
check_requirements
pre_cleanup_and_dir_setup "$IMG_NAME" "$TMP_DIR" "${CWD}/build_access_node"

cd ${TMP_DIR}

# Build linux kernel and get rpi image (rootfs)
build_linux_kernel
download_rpi_rootfs
resize_disk_image

# Mount partition, create ukama dir, build apps
mount_partitions
setup_ukama_dirs
build_apps_using_container "${UKAMA_ROOT}" "${NODE_APPS}"
copy_all_apps              "${UKAMA_ROOT}" "${NODE_APPS}"
copy_misc_files            "${UKAMA_ROOT}" "${NODE_APPS}"
copy_linux_kernel

install_required_libs
configure_systemd_service
create_ssh_user

# cleanup
unmount_partitions
cp "${TMP_DIR}/${RPI_IMG}" "${CWD}/build_access_node/${IMG_NAME}"
cleanup

cd ${CWD}
log "SUCCESS" "Access node image built successfully: $IMG_NAME"
