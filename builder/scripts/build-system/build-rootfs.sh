#!/bin/bash

set -e
set -x  

# Initialize variables
PARTITION_TYPE=""
ROOTFS_VERSION=""
SERVICE_NAME=""
SERVICE_CMD=""
SERVICE_ARGS=""
MAJOR_VERSION="v3.21"

UKAMA_ROOT="/ukama_repo"
UKAMA_REPO_APP_PKG="${UKAMA_ROOT}/build/pkg"
UKAMA_REPO_LIB_PKG="${UKAMA_ROOT}/build/lib"

UKAMA_REPO_APP_PKG="${UKAMA_ROOT}/build/pkg"

LOG_FILE=/setup.log
NODE_ID="uk-sa12-4567-a1"

# Logging function
log_message() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - [Partition: $PARTITION_TYPE] [RootFS: $ROOTFS_VERSION] $1"
}

# Function to show usage
usage() {
    echo "Usage: $0 -p <partition_type> -r <rootfs_version> -n <service_name> -c <service_command> -a <service_args>"
    echo "  -p   Partition type (active or passive)"
    echo "  -r   RootFS version"
    echo "  -n   Service name"
    echo "  -c   Service command"
    echo "  -a   Service arguments (optional)"
    exit 1
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

function install_starter_app() {
    log "INFO" "Installing starter.d"
    cd ${UKAMA_REPO_APP_PKG}
    tar zxvf starterd_latest.tar.gz
    cp starterd_latest/sbin/starter.d /sbin/
    rm -rf starterd_latest/
}

function copy_linux_kernel() {
    log "INFO" "Copying linux kernel..."
    cp "${TMP_LINUX}/arch/arm64/boot/Image" "${CWD}/build_access_node/kernel.img"
    log "SUCCESS" "Linux kernel copied"
}

function copy_all_apps() {
    log "INFO" "Copying apps"
    cp -rvf ${UKAMA_REPO_APP_PKG} ${UAKAMA_APP_PKG}
}

function copy_required_libs() {
    log "INFO" "Installing required libs"
    cd ${UKAMA_REPO_LIBS_PKG}
    tar zxvf vendor_libs.tgz -C /usr
}

function copy_misc_files() {
    log "INFO" "Copying various files to image"
    create_manifest_file $apps
    sudo cp ${MANIFEST_FILE} "/manifest.json"
    rm ${MANIFEST_FILE}

    # install the starter.d app
    install_starter_app "/"

    # update /etc/services to add ports
    log "INFO" "Adding all the apps to /etc/services"
    sudo mkdir -p "/etc"
    sudo cp "${UKAMA_ROOT}/nodes/ukamaOS/distro/scripts/files/services" \
         "/etc/services"
}


# Update /etc/fstab based on partition type
update_fstab() {
    log_message "Updating /etc/fstab for partition type: $PARTITION_TYPE"

    if [[ "$PARTITION_TYPE" == "active" ]]; then
        cat <<FSTAB > /etc/fstab
proc            /proc        proc    defaults    0 0
sysfs           /sys         sysfs   defaults    0 0
devpts          /dev/pts     devpts  defaults    0 0
tmpfs           /tmp         tmpfs   defaults    0 0
/dev/mmcblk1p2  /recovery    auto    ro          0 2
/dev/mmcblk1p7  /data        auto    ro          0 2
/dev/mmcblk1p6  /passive     auto    ro          0 2
/dev/mmcblk1p5  /            auto    errors=remount-ro  0 1
/dev/mmcblk1p1  /boot/firmware auto  ro          0 2
FSTAB
    else
        cat <<FSTAB > /etc/fstab
proc            /proc        proc    defaults    0 0
sysfs           /sys         sysfs   defaults    0 0
devpts          /dev/pts     devpts  defaults    0 0
tmpfs           /tmp         tmpfs   defaults    0 0
/dev/mmcblk1p2  /recovery    auto    ro          0 2
/dev/mmcblk1p7  /data        auto    ro          0 2
/dev/mmcblk1p5  /passive     auto    ro          0 2
/dev/mmcblk1p6  /            auto    errors=remount-ro  0 1
/dev/mmcblk1p1  /boot/firmware auto  ro          0 2
FSTAB
    fi

    log_message "/etc/fstab updated successfully."
}

# Configure network interface eth0
configure_network() {
    log_message "Configuring network for eth0"

    cat <<NETWORK > /etc/network/interfaces
auto eth0
iface eth0 inet static
    address 10.102.81.10
    netmask 255.255.255.0
    gateway 10.102.81.1
NETWORK

    # Apply network changes
    ifdown eth0 && ifup eth0
    log_message "Network configuration updated for eth0"
}



# Create a custom OpenRC service
create_openrc_service() {
    log_message "Creating OpenRC service: $SERVICE_NAME"

    mkdir -p /etc/init.d

    cat <<SERVICE > /etc/init.d/$SERVICE_NAME
#!/sbin/openrc-run

description="OpenRC Service: $SERVICE_NAME"
command="$SERVICE_CMD"
command_args="$SERVICE_ARGS"

depend() {
    need net
}

start() {
    ebegin "Starting $SERVICE_NAME"
    start-stop-daemon --start --background --exec \$command -- \$command_args
    eend \$?
}

stop() {
    ebegin "Stopping $SERVICE_NAME"
    start-stop-daemon --stop --exec \$command
    eend \$?
}
SERVICE

    chmod +x /etc/init.d/$SERVICE_NAME
    rc-update add $SERVICE_NAME default
    log_message "OpenRC service $SERVICE_NAME created and added to startup."
}

# Function to set up the root filesystem
setup_rootfs() {
    log_message "Setting up root filesystem"

    # Set up DNS
    echo "nameserver 8.8.8.8" > /etc/resolv.conf

    # Set up package repositories
    echo "https://dl-cdn.alpinelinux.org/alpine/${ROOTFS_VERSION}/main" > /etc/apk/repositories
    echo "https://dl-cdn.alpinelinux.org/alpine/${ROOTFS_VERSION}/community" >> /etc/apk/repositories

    # Update packages
    apk update
    apk upgrade

    # Install essential packages
    apk add alpine-base openrc busybox bash sudo shadow tzdata
    apk add acpid busybox-openrc busybox-extras busybox-mdev-openrc
    apk add readline bash autoconf automake libmicrohttpd-dev gnutls-dev openssl-dev iptables libuuid sqlite dhcpcd protobuf iproute2 zlib curl-dev nettle libcap libidn2 libmicrohttpd gnutls openssl-dev curl-dev linux-headers bsd-compat-headers tree libtool sqlite-dev openssl-dev readline cmake autoconf automake alpine-sdk build-base git tcpdump ethtool iperf3 htop vim doas

    # Set timezone
    ln -sf /usr/share/zoneinfo/UTC /etc/localtime

    # Configure networking
    apk add dhcpcd iproute2 iputils
    rc-update add dhcpcd default
    rc-service dhcpcd start

    # Set hostname
    echo "ukama-linux" > /etc/hostname

    # Set up root user
    echo "root:root" | chpasswd

    # Create a new user
    adduser -D -s /bin/bash -G wheel ukama
    echo "ukama:ukama" | chpasswd
    echo "%wheel ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/wheel

    # Configure doas (instead of sudo)
    apk add doas
    echo "permit persist ukama as root" > /etc/doas.d/doas.conf
    chmod 600 /etc/doas.d/doas.conf

    # Enable SSH access
    apk add openssh
    rc-update add sshd default
    rc-service sshd start

    # Enable system services
    rc-update add networking default
    rc-update add sshd default
    rc-update add dhcpcd default
    rc-update add acpid default

    # Create necessary directories
    mkdir -p /recovery /data /passive /boot/firmware

    log_message "Root filesystem setup completed."
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


function setup_ukama_dirs() {
    log "INFO" "Creating Ukama directories..."

    mkdir -p "/ukama"
    mkdir -p "/ukama/configs"
    mkdir -p "/ukama/apps"
    mkdir -p "/ukama/apps/pkgs"
    mkdir -p "/ukama/apps/rootfs"
    mkdir -p "/ukama/apps/registry"

    echo "${NODE_ID}" > "/ukama/nodeid"
    echo "localhost"  > "/ukama/bootstrap"

    touch "/ukama/apps.log"

    log "SUCCESS" "Ukama directories created."
}

setup_ukama_dirs

log "INFO" "Script ${0} called with args $#"
index=0
for arg in "$@"; do
  log "INFO" "arg[${index}]: ${arg}"
  index=$((index + 1))
done

# Parse options using getopts
while getopts "p:r:n:c:a:" opt; do
    case "${opt}" in
        p) PARTITION_TYPE="${OPTARG}" ;;
        r) ROOTFS_VERSION="${OPTARG}" ;;
        n) SERVICE_NAME="${OPTARG}" ;;
        c) SERVICE_CMD="${OPTARG}" ;;
        a) SERVICE_ARGS="${OPTARG}" ;;
        *) usage ;;
    esac
done

# Validate required arguments
if [[ -z "$PARTITION_TYPE" || -z "$ROOTFS_VERSION" || -z "$SERVICE_NAME" || -z "$SERVICE_CMD" ]]; then
    usage
fi

# Validate partition type
if [[ "$PARTITION_TYPE" != "active" && "$PARTITION_TYPE" != "passive" ]]; then
    echo "Error: Partition type must be 'active' or 'passive'."
    exit 1
fi

#copy anyother scripts required

# Main execution
setup_rootfs  # Set up root filesystem
update_fstab  # Update fstab
configure_network  # Configure network
create_openrc_service  # Create OpenRC service

copy_required_libs
copy_all_apps              "${UKAMA_ROOT}" "${NODE_APPS}"
copy_misc_files            "${UKAMA_ROOT}" "${NODE_APPS}"

install_starter_app

copy_linux_kernel




~                                                                     
