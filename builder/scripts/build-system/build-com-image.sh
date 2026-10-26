#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

set -euo pipefail

# all tools available?
for cmd in losetup parted kpartx mkfs.vfat mkfs.ext4 mkswap truncate debootstrap rsync docker; do
    command -v $cmd >/dev/null || { echo "ERROR: '$cmd' not found in PATH"; exit 1; }
done

# Load helper utilities
PKG_UTILS="$(dirname "$0")/pkg-utils.sh"
if [[ ! -f "$PKG_UTILS" ]]; then
    echo "ERROR: Missing $PKG_UTILS"
    exit 1
fi
source "$PKG_UTILS"

# Mount points and image path
BOOT_MNT="/media/boot"
PRIMARY_MNT="/media/primary"
PASSIVE_MNT="/media/passive"
RAW_IMG="ukama-com-image.img"

DIR="$(pwd)"
UKAMA_OS=$(realpath ../../../nodes/ukamaOS)
UKAMA_ROOT=$(realpath ../../../)
UKAMA_REPO_APP_PKG="${UKAMA_ROOT}/build/pkgs"
UKAMA_REPO_LIB_PKG="${UKAMA_ROOT}/build/libs"
COMMON_CONFIG_FILE="${UKAMA_ROOT}/builder/boards/common.config"
COM_CONFIG_FILE="${UKAMA_ROOT}/builder/boards/com.config"
APP_CONFIGS_DIR="${UKAMA_ROOT}/nodes/configs/apps"
BOOTSTRAP_SERVER="${BOOTSTRAP_SERVER:-dev.bootstrap.ukama.com}"
PLATFORM_LIB="${UKAMA_OS}/distro/platform/build/libusys.so"
APP_NAMES=()

# Ubuntu parameters
UBUNTU_IMAGE="${UBUNTU_IMAGE:-ubuntu:24.04}"
UBUNTU_RELEASE="${UBUNTU_RELEASE:-}"
UBUNTU_MIRROR="${UBUNTU_MIRROR:-http://archive.ubuntu.com/ubuntu}"

# Loop device handle
LOOPDEV=""

# Logging and status helpers
test -n "\${LOG_LEVEL-}" || LOG_LEVEL=INFO
log() {
    local type="$1" msg="$2"
    local color
    case "$type" in
        INFO)    color="\033[1;34m";;
        SUCCESS) color="\033[1;32m";;
        ERROR)   color="\033[1;31m";;
        *)       color="\033[1;37m";;
    esac
    echo -e "${color}${type}: ${msg}\033[0m"
}
check_status() {
    local code=$1 msg=$2 stage=$3
    if [[ $code -ne 0 ]]; then
        log ERROR "Stage '${stage}' failed"
        exit 1
    fi
    log SUCCESS "$msg"
}

# Cleanup function ensures mounts and loop are detached
cleanup() {
    log INFO "Cleaning up mounts and loop device"
    release_chroot "$PRIMARY_MNT"
    sudo umount "$BOOT_MNT"     2>/dev/null || true
    sudo umount "$PASSIVE_MNT"  2>/dev/null || true
    sudo umount "$PRIMARY_MNT"  2>/dev/null || true
    if [[ -n "$LOOPDEV" ]]; then
        sudo kpartx -dv "$LOOPDEV" 2>/dev/null || true
        sudo losetup -d "$LOOPDEV"  2>/dev/null || true
    fi
}
trap cleanup EXIT

release_chroot() {
    local rootfs="$1"
    for fs in dev/pts dev proc sys; do
        sudo umount "${rootfs}/${fs}" 2>/dev/null || true
    done
}

resolve_ubuntu_release() {
    if [[ -z "${UBUNTU_RELEASE}" ]]; then
        UBUNTU_RELEASE=$(docker run --rm "${UBUNTU_IMAGE}" sh -c '. /etc/os-release && echo "$VERSION_CODENAME"')
    fi
    [[ -n "${UBUNTU_RELEASE}" ]] || { log ERROR "Could not read the release of ${UBUNTU_IMAGE}"; exit 1; }
    log INFO "Ubuntu release: ${UBUNTU_RELEASE} (same as ${UBUNTU_IMAGE}, which builds the apps)"
}

build_apps() {
    STAGE="build_apps"
    local apps
    apps=$(IFS=,; echo "${APPS[*]}")
    log INFO "Building ${apps} on ${UBUNTU_IMAGE} with builder/docker/apps_setup.sh"

    (cd "${UKAMA_ROOT}/builder/docker" && UBUNTU_IMAGE="${UBUNTU_IMAGE}" ./apps_setup.sh com ubuntu "${UKAMA_ROOT}" "${apps}") || {
        log ERROR "App build failed, see ${UKAMA_ROOT}/apps_build.log"
        exit 1
    }
    check_status 0 "Apps built into ${UKAMA_REPO_APP_PKG}" "${STAGE}"
}

create_disk_image() {
    log INFO "Creating sparse image ${RAW_IMG} (16 GiB)"
    rm -f "$RAW_IMG"
    truncate -s 16G "$RAW_IMG"
    # smaller image for testing
    #    truncate -s 5G "$RAW_IMG"
    check_status $? "Raw image ready" create_disk_image
}

attach_loop() {
    log INFO "Attaching ${RAW_IMG} to loop device"
    LOOPDEV=$(sudo losetup --show --find --partscan "$RAW_IMG")
    if [[ -z "$LOOPDEV" ]]; then
        log ERROR "Failed to attach loop device"
        exit 1
    fi
    check_status 0 "Loop device: $LOOPDEV" attach_loop
}

partition_image() {
    log INFO "Partitioning $LOOPDEV: boot, passive, primary, swap"
    sudo parted -s "$LOOPDEV" mklabel msdos \
         mkpart primary fat32      1MiB    513MiB   \
         mkpart primary ext4       513MiB  4609MiB  \
         mkpart primary ext4       4609MiB 12289MiB \
         mkpart primary linux-swap 12289MiB 100%     \
         set 1 boot on
    check_status $? "Partitions created" partition_image
}

# smaller image for testing (5GB)
#partition_image() {
#    log INFO "Partitioning $LOOPDEV: boot (512 MiB), primary (2 GiB), passive (2 GiB), swap (rest)"
#    sudo parted -s "$LOOPDEV" mklabel msdos \
#         mkpart primary fat32      1MiB     513MiB   \
#         mkpart primary ext4       513MiB   2561MiB  \
#         mkpart primary ext4       2561MiB  4609MiB  \
#         mkpart primary linux-swap 4609MiB  100%     \
#         set 1 boot on
#    check_status $? "Partitions created" partition_image
#}

map_partitions() {
    log INFO "Mapping partitions for $LOOPDEV"
    sudo kpartx -av "$LOOPDEV" | tee /dev/stderr
    sleep 1
    check_status $? "Partitions mapped" map_partitions
}

format_partitions() {
    log INFO "Formatting partitions"
    local b=$(basename "$LOOPDEV")
    local p1="/dev/mapper/${b}p1" p2="/dev/mapper/${b}p2"
    local p3="/dev/mapper/${b}p3" p4="/dev/mapper/${b}p4"

    sudo mkfs.vfat -F32 -n boot    "$p1"
    check_status $? "boot formatted" format_partitions

    sudo mkfs.ext4 -L passive \
         -O ^64bit,^metadata_csum \
         "$p2"
    check_status $? "passive formatted" format_partitions

    sudo mkfs.ext4 -L primary \
         -O ^64bit,^metadata_csum \
         "$p3"
    check_status $? "primary formatted" format_partitions

    sudo mkswap    -L swap         "$p4"
    check_status $? "swap created" format_partitions
}

mount_partition() {
    local part=$1 mp=$2
    log INFO "Mounting $part → $mp"
    sudo mkdir -p "$mp"
    sudo mount "$part" "$mp"
    check_status $? "Mounted $part" mount_partition
}

install_ubuntu() {
    local rootfs="$1"
    STAGE="install_ubuntu"
    log INFO "Installing Ubuntu ${UBUNTU_RELEASE} into ${rootfs}"

    sudo debootstrap --arch amd64 "${UBUNTU_RELEASE}" "${rootfs}" "${UBUNTU_MIRROR}" \
         /usr/share/debootstrap/scripts/gutsy
    check_status $? "Ubuntu ${UBUNTU_RELEASE} installed in ${rootfs}" "${STAGE}"

    sudo mount --bind /dev "${rootfs}/dev"
    sudo mount --bind /dev/pts "${rootfs}/dev/pts"
    sudo mount -t proc proc "${rootfs}/proc"
    sudo mount -t sysfs sysfs "${rootfs}/sys"
    sudo cp -L /etc/resolv.conf "${rootfs}/etc/resolv.conf"
    printf '#!/bin/sh\nexit 101\n' | sudo tee "${rootfs}/usr/sbin/policy-rc.d" > /dev/null
    sudo chmod 0755 "${rootfs}/usr/sbin/policy-rc.d"

    sudo chroot "${rootfs}" /usr/bin/env -i \
        PATH=/usr/sbin:/usr/bin:/sbin:/bin DEBIAN_FRONTEND=noninteractive \
        RELEASE="${UBUNTU_RELEASE}" MIRROR="${UBUNTU_MIRROR}" \
        /bin/bash -e <<'EOF'
cat > /etc/apt/sources.list <<APT
deb ${MIRROR} ${RELEASE} main universe
deb ${MIRROR} ${RELEASE}-updates main universe
deb ${MIRROR} ${RELEASE}-security main universe
APT
apt-get update
apt-get install -y --no-install-recommends \
    linux-image-generic initramfs-tools grub-efi-amd64-bin \
    systemd-sysv systemd-resolved systemd-timesyncd udev kmod dbus \
    openssh-server sudo locales tzdata ca-certificates curl \
    iproute2 iputils-ping ethtool tcpdump i2c-tools libjansson4 libulfius2.7t64 \
    openvswitch-switch openvpn iptables
apt-get clean

locale-gen en_US.UTF-8
ln -sf /usr/share/zoneinfo/UTC /etc/localtime
echo "ukama-linux" > /etc/hostname
printf '127.0.0.1 localhost\n127.0.1.1 ukama-linux\n' > /etc/hosts
echo "root:root" | chpasswd
id ukama > /dev/null 2>&1 || useradd -m -s /bin/bash -G sudo ukama
echo "ukama:ukama" | chpasswd
echo "ukama ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/ukama
chmod 0440 /etc/sudoers.d/ukama

cat > /etc/systemd/network/20-wired.network <<NET
[Match]
Name=en* eth*

[Network]
DHCP=yes
NET
systemctl enable systemd-networkd systemd-resolved systemd-timesyncd ssh serial-getty@ttyS0.service
EOF
    check_status $? "Kernel, systemd, SSH, Open vSwitch, OpenVPN and iptables in ${rootfs}" "${STAGE}"

    sudo rm -f "${rootfs}/usr/sbin/policy-rc.d"
    sudo ln -sf /run/systemd/resolve/stub-resolv.conf "${rootfs}/etc/resolv.conf"
}

setup_ukama() {
    local rootfs="$1"
    STAGE="setup_ukama"
    log INFO "Installing starter.d, the noded sysfs link and Ukama files in ${rootfs}"

    sudo mkdir -p "${rootfs}/ukama/configs" \
                  "${rootfs}/ukama/apps/pkgs" \
                  "${rootfs}/ukama/apps/rootfs" \
                  "${rootfs}/ukama/apps/registry" \
                  "${rootfs}/ukama/mocksysfs" \
                  "${rootfs}/ukama/logs" \
                  "${rootfs}/ukama/state/starterd/log-spool"

    tar -zxOf "${UKAMA_REPO_APP_PKG}/starterd_latest.tar.gz" --wildcards '*/sbin/starter.d' \
        | sudo tee "${rootfs}/sbin/starter.d" > /dev/null
    sudo chmod 0755 "${rootfs}/sbin/starter.d"

    sudo tee "${rootfs}/etc/systemd/system/starterd.service" > /dev/null <<'EOF'
[Unit]
Description=Ukama's capp starter.d
After=network.target ukama-sysfs.service

[Service]
Environment=STARTERD_MANIFEST=/ukama/manifest.json
Environment=STARTERD_RLOG_SOCKET=/run/ukama/rlog-ingest.sock
Environment=STARTERD_LOG_SPOOL_DIR=/ukama/state/starterd/log-spool
ExecStart=/sbin/starter.d
Type=simple

[Install]
WantedBy=multi-user.target
EOF

    sudo tee "${rootfs}/etc/systemd/system/ukama-sysfs.service" > /dev/null <<'EOF'
[Unit]
Description=Link noded sysfs data into /tmp/sys
After=local-fs.target
Before=starterd.service

[Service]
Type=oneshot
RemainAfterExit=yes
ExecStart=/bin/sh -c 'if [ -d /ukama/mocksysfs/sys ]; then rm -rf /tmp/sys; ln -sfn /ukama/mocksysfs/sys /tmp/sys; fi'

[Install]
WantedBy=multi-user.target
EOF
    sudo chroot "${rootfs}" systemctl enable starterd.service ukama-sysfs.service

    sudo cp -r "${UKAMA_OS}/distro/system/noded/mfgdata" "${rootfs}/ukama/mocksysfs/"
    sudo ln -sfn /sys "${rootfs}/ukama/mocksysfs/sys"
    sudo cp "${UKAMA_OS}/distro/scripts/files/services" "${rootfs}/etc/services"
    check_status $? "starter.d, ukama-sysfs and Ukama files in ${rootfs}" "${STAGE}"
}

copy_boot() {
    STAGE="copy_boot_partition"
    log INFO "Copying the Ubuntu kernel, initrd and GRUB → ${BOOT_MNT}"

    local src="${PRIMARY_MNT}/boot"
    local dst="${BOOT_MNT}"
    local kernel initrd

    # 1) Clean out anything already on the boot partition
    sudo rm -rf "${dst:?}/"*
    check_status $? "Cleared old files on boot partition" "${STAGE}"

    kernel=$(ls "${src}"/vmlinuz-* | sort -V | tail -n 1)
    initrd="${src}/initrd.img-${kernel##*/vmlinuz-}"
    if [[ ! -f "${kernel}" || ! -f "${initrd}" ]]; then
        log ERROR "Missing kernel or initrd in ${src}"
        exit 1
    fi

    sudo mkdir -p "${dst}/boot/grub" "${dst}/EFI/BOOT"
    sudo cp "${kernel}" "${dst}/boot/vmlinuz"
    sudo cp "${initrd}" "${dst}/boot/initrd.img"
    check_status $? "Copied $(basename "${kernel}") and its initrd" "${STAGE}"

    sudo tee "${dst}/boot/grub/grub.cfg" > /dev/null <<EOF
insmod part_msdos
insmod fat

set default=0
set timeout=3

menuentry "ukama (linux)" {
    linux /boot/vmlinuz root=LABEL=primary rootfstype=ext4 rw \\
        console=ttyS0,115200 console=tty1
    initrd /boot/initrd.img
}
EOF
    check_status $? "grub.cfg written" "${STAGE}"

    printf 'insmod part_msdos\ninsmod fat\nsearch --no-floppy --label boot --set=root\nconfigfile ($root)/boot/grub/grub.cfg\n' \
        | sudo tee "${PRIMARY_MNT}/tmp/grub-embed.cfg" > /dev/null
    sudo chroot "${PRIMARY_MNT}" grub-mkstandalone -O x86_64-efi -o /tmp/BOOTX64.EFI \
        "boot/grub/grub.cfg=/tmp/grub-embed.cfg"
    sudo cp "${PRIMARY_MNT}/tmp/BOOTX64.EFI" "${dst}/EFI/BOOT/BOOTX64.EFI"
    sudo rm -f "${PRIMARY_MNT}/tmp/BOOTX64.EFI" "${PRIMARY_MNT}/tmp/grub-embed.cfg"
    check_status $? "UEFI BOOTX64.EFI installed" "${STAGE}"

    log SUCCESS "Boot partition populated under /boot" "${STAGE}"
}

copy_rootfs() {
    STAGE="copy_rootfs"
    log INFO "Copying primary (${PRIMARY_MNT}) → passive (${PASSIVE_MNT})"

    sudo rsync -aAXH --delete \
         --exclude=/dev/* \
         --exclude=/proc/* \
         --exclude=/sys/* \
         --exclude=/run/* \
         --exclude=/tmp/* \
         "${PRIMARY_MNT}/" "${PASSIVE_MNT}/"
    check_status $? "Rootfs copied to ${PASSIVE_MNT}" "${STAGE}"
}

update_fstab() {
    local rootfs_dir="$1"
    local fstab="${rootfs_dir}/etc/fstab"
    STAGE="update_fstab"

    log INFO "Writing label-based fstab in ${fstab}"

    sudo tee "${fstab}" > /dev/null <<'FSTAB'
# pseudo-file systems
devtmpfs        /dev        devtmpfs    defaults    0 0
proc            /proc       proc        defaults    0 0
sysfs           /sys        sysfs       defaults    0 0
devpts          /dev/pts    devpts      defaults    0 0
tmpfs           /tmp        tmpfs       defaults    0 0

# real partitions by LABEL
LABEL=primary   /           ext4      defaults         0 1
LABEL=boot      /boot       auto      rw,defaults      0 2
LABEL=swap      none        swap      sw               0 0
FSTAB

    check_status $? "fstab written" "${STAGE}"
}

deploy_to_rootfs() {
    local rootfs="$1"

    log INFO "Deploying apps/libs + manifest into ${rootfs}"

    # ensure dirs exist
    sudo mkdir -p "${rootfs}/ukama/apps/pkgs" \
                     "${rootfs}/ukama/logs" \
                     "${rootfs}/ukama/state/starterd/log-spool" \
                     "${rootfs}/lib"

    # copy apps
    copy_all_apps      "$UKAMA_REPO_APP_PKG" "${rootfs}/ukama/apps/pkgs"
    check_status $? "copy_all_apps to ${rootfs}" "deploy_to_rootfs"

    # copy libs
    sudo cp "${PLATFORM_LIB}" "${rootfs}/lib/x86_64-linux-gnu/"
    check_status $? "libusys.so copied to ${rootfs}" "deploy_to_rootfs"
    if sudo chroot "${rootfs}" ldd /sbin/starter.d | grep "not found"; then
        log ERROR "starter.d is missing libraries in ${rootfs}"
        exit 1
    fi

    # create manifest
    create_manifest_file "${rootfs}/ukama/manifest.json" "${APPS[@]}"
    check_status $? "create_manifest_file in ${rootfs}" "deploy_to_rootfs"

    sudo mkdir -p "${rootfs}/ukama/configs/metricsd" \
                  "${rootfs}/ukama/configs/aggregator"
    for app in bootstrap meshd notifyd init-network; do
        sudo cp -r "${APP_CONFIGS_DIR}/${app}" "${rootfs}/ukama/configs/"
    done
    sudo cp "${APP_CONFIGS_DIR}/metricsd/com_config.toml" \
            "${rootfs}/ukama/configs/metricsd/config.toml"
    sudo cp "${APP_CONFIGS_DIR}/aggregator/tnode_config.toml" \
            "${rootfs}/ukama/configs/aggregator/config.toml"
    check_status $? "App configs copied to ${rootfs}" "deploy_to_rootfs"

    echo "${BOOTSTRAP_SERVER}" | sudo tee "${rootfs}/ukama/bootstrap" > /dev/null
    check_status $? "Bootstrap server ${BOOTSTRAP_SERVER} set in ${rootfs}" "deploy_to_rootfs"
}

# Main

# Gather the list of enabled apps
get_enabled_apps "$COMMON_CONFIG_FILE" "$COM_CONFIG_FILE"
if (( ${#APPS[@]} > 0 )); then
    log INFO "Enabled apps: ${APPS[*]}"
else
    log ERROR "No apps enabled. Aborting."
    exit 1
fi

resolve_ubuntu_release
build_apps

mkdir -p "$BOOT_MNT" "$PASSIVE_MNT" "$PRIMARY_MNT"

create_disk_image
attach_loop
partition_image
map_partitions
format_partitions
mount_partition "/dev/mapper/$(basename ${LOOPDEV})p1" "$BOOT_MNT"
mount_partition "/dev/mapper/$(basename ${LOOPDEV})p2" "$PASSIVE_MNT"
mount_partition "/dev/mapper/$(basename ${LOOPDEV})p3" "$PRIMARY_MNT"

# Signal success; cleanup trap will unmount and detach
log SUCCESS "Disk image ${RAW_IMG} prepared successfully"

install_ubuntu "${PRIMARY_MNT}"
setup_ukama "${PRIMARY_MNT}"
deploy_to_rootfs "${PRIMARY_MNT}"
copy_boot
release_chroot "${PRIMARY_MNT}"
update_fstab "${PRIMARY_MNT}"
copy_rootfs

log "SUCCESS" "Disk image creation completed successfully! (Ubuntu ${UBUNTU_RELEASE}, bootstrap: ${BOOTSTRAP_SERVER})"
exit 0
