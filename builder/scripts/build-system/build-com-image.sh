#!/usr/bin/env bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
# Copyright (c) 2026-present, Ukama Inc.

# Build the Ubuntu COM disk image with the cnode Alpine app build pipeline.

set -euo pipefail
export LC_ALL=C

SCRIPT_DIR=$(dirname -- "${BASH_SOURCE[0]}")
SCRIPT_DIR=$(realpath -e "$SCRIPT_DIR")
UKAMA_ROOT=$(realpath -e "$SCRIPT_DIR/../../..")
OUTPUT="$PWD/ukama-com-image.img"
ALPINE_RELEASE=v3.17
UBUNTU_RELEASE=${UBUNTU_RELEASE:-noble}
UBUNTU_MIRROR=${UBUNTU_MIRROR:-http://archive.ubuntu.com/ubuntu}
BOOTSTRAP_SERVER=${BOOTSTRAP_SERVER:-dev.bootstrap.ukama.com}
WORK=""
HIDDEN_CONFIGS=""
APP_BUILD_STARTED=no
LOOPDEV=""
MOUNTS=()
APPS=()

usage() {
    cat <<'EOF'
Usage: ./build-com-image.sh [--output IMAGE] [--repo UKAMA_ROOT]

Creates a 16 GiB Ubuntu 24.04 amd64 image for the COM/SMARC board:
  P1 boot (UEFI ESP) 512 MiB; P2 passive 4 GiB;
  P3 primary 7.5 GiB; P4 swap 4095 MiB.

The cnode build-env-setup.sh pipeline builds COM apps on Alpine v3.17,
using x86_64 instead of aarch64. Existing app recipes are unchanged.
The companion validate-com-image.sh must pass before IMAGE is replaced.
Run as the build user with sudo privileges. No container engine is needed.
Host packages: debootstrap dosfstools e2fsprogs fdisk rsync jq binutils
               initramfs-tools-core zstd make wget git sudo.

Environment: UBUNTU_RELEASE=noble,
             UBUNTU_MIRROR=http://archive.ubuntu.com/ubuntu,
             BOOTSTRAP_SERVER=dev.bootstrap.ukama.com
Ubuntu remains the base OS. Alpine app libraries have a private runtime.
The Alpine branch is pinned; package updates follow the cnode build policy.
The original root/root and ukama/ukama debug credentials are retained.
Secure Boot is not supported by this unsigned standalone GRUB image.
EOF
}

log() {
    printf '%s: %s\n' "$1" "$2"
}

die() {
    log ERROR "$*" >&2
    exit 1
}

as_root() {
    if [ "$EUID" -eq 0 ]; then
        "$@"
    else
        sudo "$@"
    fi
}

release_mounts() {
    local index

    index=${#MOUNTS[@]}
    while [ "$index" -gt 0 ]; do
        index=$((index - 1))
        if ! as_root umount "${MOUNTS[index]}"; then
            return 1
        fi
        unset 'MOUNTS[index]'
    done

    if [ -n "$LOOPDEV" ]; then
        if ! as_root losetup --detach "$LOOPDEV"; then
            return 1
        fi
        LOOPDEV=""
    fi
}

cleanup() {
    local status=$?

    trap - EXIT
    trap '' HUP INT TERM
    if ! restore_configs; then
        status=1
    fi
    if release_mounts; then
        if ! cleanup_app_build; then
            status=1
        fi
        if [ -n "$WORK" ]; then
            if work_is_unmounted && [ -z "$HIDDEN_CONFIGS" ]; then
                if ! as_root rm -rf -- "$WORK"; then
                    status=1
                fi
            else
                log ERROR "Retained work directory: $WORK" >&2
                status=1
            fi
        fi
    else
        log ERROR "Cleanup incomplete; retained work directory: $WORK" >&2
        status=1
    fi
    if [ "$status" -ne 0 ]; then
        log ERROR "COM image build failed; see the error above" >&2
        if [ "$APP_BUILD_STARTED" = yes ]; then
            log ERROR "App build log: $UKAMA_ROOT/apps_build.log" >&2
        fi
    fi
    exit "$status"
}

mount_fs() {
    local target=$1

    shift
    as_root mkdir -p "$target"
    as_root mount "$@" "$target"
    MOUNTS+=("$target")
}

parse_args() {
    local output_dir
    local output_name

    while [ "$#" -gt 0 ]; do
        case "$1" in
            --output)
                if [ "$#" -lt 2 ]; then
                    die "Missing value for --output"
                fi
                OUTPUT=$2
                shift 2
                ;;
            --repo)
                if [ "$#" -lt 2 ]; then
                    die "Missing value for --repo"
                fi
                UKAMA_ROOT=$2
                shift 2
                ;;
            -h|--help)
                usage
                exit 0
                ;;
            *)
                die "Unknown argument: $1"
                ;;
        esac
    done

    UKAMA_ROOT=$(realpath -e "$UKAMA_ROOT")
    output_dir=$(dirname -- "$OUTPUT")
    output_dir=$(realpath -e "$output_dir")
    output_name=$(basename -- "$OUTPUT")
    OUTPUT="$output_dir/$output_name"

    # Resolve only the parent. Never follow an existing output symlink.
    if [ -L "$OUTPUT" ]; then
        die "Output must not be a symlink: $OUTPUT"
    fi
    if [ -d "$OUTPUT" ]; then
        die "Output must not be a directory: $OUTPUT"
    fi
    if [ -b "$OUTPUT" ]; then
        die "Output must not be a block device: $OUTPUT"
    fi
}

require_file() {
    local file=$1

    if [ ! -s "$file" ]; then
        die "Missing input file: $file"
    fi
}

check_host() {
    local tool
    local architecture
    local tools=(
        debootstrap sfdisk losetup mount umount
        mkfs.vfat mkfs.ext4 mkswap truncate rsync chroot
        jq tar gzip readelf sha256sum blkid blockdev
        fsck.fat e2fsck objdump diff find realpath awk
        sed grep stat du df install lsinitramfs zstd
        make wget git sudo findmnt cmp env
    )

    for tool in "${tools[@]}"; do
        if ! command -v "$tool" >/dev/null; then
            die "Missing host command: $tool"
        fi
    done

    architecture=$(uname -m)
    if [ "$architecture" != x86_64 ]; then
        die "Build on an x86_64 Linux host"
    fi

    if [ "$EUID" -ne 0 ]; then
        if ! command -v sudo >/dev/null; then
            die "sudo is required for image operations"
        fi
        sudo -v
    fi
}

check_inputs() {
    local builder="$UKAMA_ROOT/builder"
    local distro="$UKAMA_ROOT/nodes/ukamaOS/distro"

    require_file "$builder/scripts/build-system/pkg-utils.sh"
    require_file "$builder/scripts/build-system/build-env-setup.sh"
    require_file "$builder/scripts/build-system/build-distro.sh"
    require_file "$builder/scripts/build-system/chroot-utils.sh"
    require_file "$builder/scripts/build-all-apps.sh"
    require_file "$builder/boards/common.config"
    require_file "$builder/boards/com.config"
    require_file "$distro/scripts/files/services"
    require_file "$distro/system/noded/mfgdata/property/property.json"
    require_file "$SCRIPT_DIR/validate-com-image.sh"
    require_file /usr/share/debootstrap/scripts/gutsy

    if [ "$UBUNTU_RELEASE" != noble ]; then
        die "This COM image supports Ubuntu noble (24.04)"
    fi
    case "$BOOTSTRAP_SERVER" in
        ''|*[!A-Za-z0-9._:/-]*)
            die "Invalid BOOTSTRAP_SERVER"
            ;;
    esac
}

preflight() {
    local builder="$UKAMA_ROOT/builder"
    local output_dir

    check_host
    check_inputs

    # shellcheck source=/dev/null
    source "$builder/scripts/build-system/pkg-utils.sh"
    # Package helpers are shared with the standalone validator.
    # shellcheck source=validate-com-image.sh
    source "$SCRIPT_DIR/validate-com-image.sh"

    get_enabled_apps "$builder/boards/common.config" \
        "$builder/boards/com.config"
    if [ "${#APPS[@]}" -eq 0 ]; then
        die "No COM apps enabled"
    fi
    if ! _manifest_has_app starterd "${APPS[@]}"; then
        die "COM requires starterd"
    fi
    log INFO "Enabled apps: ${APPS[*]}"

    output_dir=$(dirname -- "$OUTPUT")
    WORK=$(mktemp -d "$output_dir/.com-build.XXXXXXXX")
    trap cleanup EXIT
    trap 'exit 129' HUP
    trap 'exit 130' INT
    trap 'exit 143' TERM
    mkdir -p "$WORK/overlay" "$WORK/packages" "$WORK/rootfs"
}

restore_configs() {
    local config
    local target

    if [ -z "$HIDDEN_CONFIGS" ]; then
        return 0
    fi
    for config in "$HIDDEN_CONFIGS"/*.toml; do
        if [ ! -f "$config" ]; then
            continue
        fi
        target="$UKAMA_ROOT/builder/configs/${config##*/}"
        if [ -e "$target" ]; then
            log ERROR "Cannot restore recipe; already exists: $target"
            return 1
        fi
        if ! mv -- "$config" "$target"; then
            return 1
        fi
    done
    if ! rmdir "$HIDDEN_CONFIGS"; then
        return 1
    fi
    HIDDEN_CONFIGS=""
}

clean_build_outputs() {
    local distro="$UKAMA_ROOT/nodes/ukamaOS/distro"
    local vendor="$distro/vendor"
    local recipe
    local failed=0
    local recipes=(
        "$UKAMA_ROOT"/nodes/apps/*/Makefile
        "$distro"/system/*/Makefile
        "$UKAMA_ROOT/builder/Makefile"
        "$distro/platform/Makefile"
        "$vendor/orcania/src/Makefile"
        "$vendor/yder/src/Makefile"
        "$vendor/ulfius/src/Makefile"
        "$vendor/tomlc/Makefile"
        "$vendor/libcap/Makefile"
        "$vendor/Makefile"
    )

    log INFO "Cleaning apps, builder, platform and vendor outputs"
    for recipe in "${recipes[@]}"; do
        if [ ! -f "$recipe" ]; then
            continue
        fi
        if ! as_root make -C "${recipe%/Makefile}" clean; then
            failed=1
        fi
    done
    if ! as_root rm -rf -- \
        "$vendor/jansson/tmpbuild" "$vendor/amqp/tmpbuild" \
        "$vendor/libuuid/tmpbuild" "$vendor/casync/tmpbuild" \
        "$vendor/prometheus-client/prom/build" \
        "$vendor/prometheus-client/promhttp/build" \
        "$UKAMA_ROOT/build/pkgs"; then
        failed=1
    fi
    if ! as_root rm -f -- "$UKAMA_ROOT/build/libs/vendor_libs.tgz"; then
        failed=1
    fi
    return "$failed"
}

cleanup_app_build() {
    local buildenv="$WORK/app-build/buildenv"
    local distro="$UKAMA_ROOT/nodes/ukamaOS/distro"
    local owner
    local dir
    local failed=0

    if [ "$APP_BUILD_STARTED" = no ]; then
        return 0
    fi
    if [ -x "$buildenv/destroy" ]; then
        if ! as_root "$buildenv/destroy" --remove; then
            return 1
        fi
    fi
    if ! work_is_unmounted; then
        return 1
    fi
    if ! clean_build_outputs; then
        failed=1
    fi
    if ! owner=$(stat -c '%u:%g' "$UKAMA_ROOT"); then
        return 1
    fi
    for dir in "$distro/vendor" "$distro/system" "$distro/platform" \
        "$UKAMA_ROOT/nodes/apps" "$UKAMA_ROOT/builder" \
        "$UKAMA_ROOT/build"; do
        if [ ! -d "$dir" ]; then
            continue
        fi
        if ! as_root find "$dir" -xdev -name .git -prune -o -uid 0 \
            -exec chown -h "$owner" {} +; then
            failed=1
        fi
    done
    if [ "$failed" -eq 0 ]; then
        APP_BUILD_STARTED=no
    fi
    return "$failed"
}

work_is_unmounted() {
    local mounts
    local target

    if ! mounts=$(findmnt -rn -o TARGET); then
        return 1
    fi
    while IFS= read -r target; do
        printf -v target '%b' "$target"
        case "$target" in
            "$WORK"|"$WORK"/*)
                log ERROR "Still mounted: $target"
                return 1
                ;;
        esac
    done <<<"$mounts"
    return 0
}

build_apps() {
    local config
    local key
    local script="$UKAMA_ROOT/builder/scripts/build-system"

    APP_BUILD_STARTED=yes
    clean_build_outputs
    mkdir -p "$WORK/app-build" "$WORK/alpine-tmp"
    HIDDEN_CONFIGS="$WORK/hidden-configs"
    mkdir -p "$HIDDEN_CONFIGS"

    # Match cnode: expose only this board's recipes to build-all-apps.sh.
    for config in "$UKAMA_ROOT"/builder/configs/*.toml; do
        key=$(basename "$config" .toml)
        if ! _manifest_has_app "$key" "${APPS[@]}"; then
            mv -- "$config" "$HIDDEN_CONFIGS/"
        fi
    done

    log INFO "Building COM apps with cnode's Alpine $ALPINE_RELEASE pipeline"
    log INFO "Build log: $UKAMA_ROOT/apps_build.log"
    (
        cd "$WORK/app-build"
        as_root env -u ALPINE_BUILD TEMP_DIR="$WORK/alpine-tmp" \
            bash "$script/build-env-setup.sh" \
            -a x86_64 -v "$ALPINE_RELEASE" -i "$UKAMA_ROOT"
    ) >"$UKAMA_ROOT/apps_build.log" 2>&1
    restore_configs
    log INFO "COM app packages and vendor_libs.tgz built"
}

stage_packages() {
    local overlay="$WORK/overlay"
    local config
    local key
    local name
    local tag
    local package
    local payload
    local space
    local target
    local config_dir

    create_manifest_file "$overlay/ukama/manifest.json" "${APPS[@]}"
    jq -e . "$overlay/ukama/manifest.json" >/dev/null

    for key in "${APPS[@]}"; do
        config="$UKAMA_ROOT/builder/configs/$key.toml"
        name=$(_manifest_toml_value "$config" capp-exec name)
        tag=$(_manifest_toml_value "$config" capp-exec version)
        if [ -z "$name" ]; then
            name=$key
        fi
        if [ -z "$tag" ]; then
            tag=latest
        fi
        if ! com_valid_name "$name"; then
            die "Invalid app name in $config"
        fi
        if ! com_valid_name "$tag"; then
            die "Invalid app tag in $config"
        fi

        package="$UKAMA_ROOT/build/pkgs/${name}_${tag}.tar.gz"
        if [ ! -s "$package" ]; then
            die "Missing built package: $package"
        fi
        if ! com_unpack_package "$package" "$WORK/packages/$name"; then
            die "Invalid package: $package"
        fi
        if [ "$COM_TOP" != "${name}_$COM_VERSION" ]; then
            die "Package name/VERSION mismatch: $package"
        fi
        payload=$COM_PAYLOAD
        cp -L -- "$package" "$overlay/ukama/apps/pkgs/"

        if [ "$key" = starterd ]; then
            cp -a "$payload" "$overlay/opt/ukama/starterd"
            continue
        fi

        space=$(com_app_space "$overlay/ukama/manifest.json" "$name")
        target="/ukama/apps/$space/$name/current"
        config_dir=$(com_config_dir "$payload" "$name")

        # starter.d creates current before launching the app. These links
        # then expose the packaged configs at the paths used by its argv.
        if [ -n "$config_dir" ]; then
            ln -s "$target/$config_dir" "$overlay/ukama/configs/$name"
        fi
    done
}

copy_runtime_libraries() {
    local source=$1
    local destination="$WORK/overlay/opt/ukama/runtime/lib"
    local entry
    local file
    local name

    find "$source" \( -type f -o -type l \) \
        -name '*.so*' -print0 >"$WORK/runtime-files.list"
    while IFS= read -r -d '' entry; do
        file=$(com_resolve_file "$source" "${entry#"$source"}")
        if ! com_elf "$file"; then
            die "Expected an x86_64 runtime library: $entry"
        fi
        name=${entry##*/}
        if [ -e "$destination/$name" ]; then
            if ! cmp -s "$file" "$destination/$name"; then
                die "Different app libraries share a name: $name"
            fi
        else
            cp -- "$file" "$destination/$name"
        fi
    done <"$WORK/runtime-files.list"
}

stage_app_runtime() {
    local overlay="$WORK/overlay"
    local archive="$UKAMA_ROOT/build/libs/vendor_libs.tgz"
    local vendor="$WORK/vendor-libs"
    local runtime="$overlay/opt/ukama/runtime/lib"
    local member
    local payload
    local loader=ld-musl-x86_64.so.1

    require_file "$archive"
    mkdir -p "$vendor" "$runtime" "$overlay/usr/lib"
    tar -tzf "$archive" >"$WORK/vendor-members.list"
    while IFS= read -r member; do
        # cnode's build-distro.sh produces a flat library archive.
        if ! com_safe_member "$member"; then
            die "Invalid vendor library archive member: $member"
        fi
        case "$member" in
            */*)
                die "Expected a flat vendor library archive: $member"
                ;;
        esac
    done <"$WORK/vendor-members.list"
    tar --no-same-owner -xzf "$archive" -C "$vendor"
    copy_runtime_libraries "$vendor"
    for payload in "$WORK"/packages/*/*; do
        if [ -d "$payload" ]; then
            copy_runtime_libraries "$payload"
        fi
    done

    # Use the loader and dependencies from the same Alpine build as the apps.
    # Ubuntu's glibc libraries and OVS binaries keep their normal search path.
    require_file "$runtime/$loader"
    chmod 0755 "$runtime/$loader"
    cp "$runtime/$loader" "$overlay/usr/lib/$loader"
    if [ ! -e "$runtime/libc.musl-x86_64.so.1" ]; then
        ln -s "$loader" "$runtime/libc.musl-x86_64.so.1"
    fi
    printf '/opt/ukama/runtime/lib\n' \
        >"$overlay/etc/ld-musl-x86_64.path"
    ln -s usr/lib "$overlay/lib"
}

stage_ukama() {
    local overlay="$WORK/overlay"
    local distro="$UKAMA_ROOT/nodes/ukamaOS/distro"

    mkdir -p "$overlay/ukama/apps/pkgs"
    mkdir -p "$overlay/ukama/apps/rootfs"
    mkdir -p "$overlay/ukama/apps/registry"
    mkdir -p "$overlay/ukama/configs"
    mkdir -p "$overlay/ukama/logs"
    mkdir -p "$overlay/ukama/mocksysfs"
    mkdir -p "$overlay/ukama/state/starterd/log-spool"
    mkdir -p "$overlay/opt/ukama"
    mkdir -p "$overlay/usr/sbin"
    mkdir -p "$overlay/etc/systemd/system"
    mkdir -p "$overlay/etc/tmpfiles.d"
    mkdir -p "$overlay/etc/systemd/journald.conf.d"
    stage_packages
    stage_app_runtime

    # starter.d must run before it can unpack other packages.
    # Keep its complete runtime together under /opt/ukama/starterd.
    cat >"$overlay/usr/sbin/starter.d" <<'EOF'
#!/bin/sh
set -eu
root=/opt/ukama/starterd
cd "$root"
unset LD_LIBRARY_PATH LD_PRELOAD
exec "$root/sbin/starter.d" "$@"
EOF
    chmod 0755 "$overlay/usr/sbin/starter.d"
    cat >"$overlay/etc/systemd/system/starterd.service" <<'EOF'
[Unit]
Description=Ukama app supervisor
After=network.target ukama-sysfs.service
Requires=ukama-sysfs.service

[Service]
Type=simple
RuntimeDirectory=ukama
RuntimeDirectoryPreserve=yes
Environment=STARTERD_MANIFEST=/ukama/manifest.json
Environment=STARTERD_RLOG_SOCKET=/run/ukama/rlog-ingest.sock
Environment=STARTERD_LOG_SPOOL_DIR=/ukama/state/starterd/log-spool
Environment=LIFECYCLED_BOOT_ID_FILE=/run/ukama/boot-id
ExecStartPre=/bin/sh -ec \
    'cat /proc/sys/kernel/random/boot_id > /run/ukama/boot-id'
ExecStart=/sbin/starter.d
Restart=on-failure
RestartSec=5
KillMode=control-group

[Install]
WantedBy=multi-user.target
EOF
    cat >"$overlay/etc/systemd/system/ukama-sysfs.service" <<'EOF'
[Unit]
Description=Expose COM sysfs to noded
After=local-fs.target systemd-tmpfiles-setup.service
Before=starterd.service

[Service]
Type=oneshot
RemainAfterExit=yes
ExecStart=/bin/ln -sfnT /ukama/mocksysfs/sys /tmp/sys

[Install]
WantedBy=multi-user.target
EOF
    cat >"$overlay/etc/tmpfiles.d/ukama.conf" <<'EOF'
# Keep noded's sysfs alias through periodic /tmp cleanup.
x /tmp/sys - - - -
d /run/openvswitch 0755 root root -
EOF
    cat >"$overlay/etc/systemd/journald.conf.d/ukama.conf" <<'EOF'
[Journal]
SystemMaxUse=128M
RuntimeMaxUse=32M
EOF
    cp -a "$distro/system/noded/mfgdata" "$overlay/ukama/mocksysfs/"
    ln -s /sys "$overlay/ukama/mocksysfs/sys"
    # Service discovery is needed before starter.d can unpack any app.
    cp "$distro/scripts/files/services" "$overlay/etc/services"
    printf '%s\n' "$BOOTSTRAP_SERVER" >"$overlay/ukama/bootstrap"
}

install_ubuntu() {
    local root=$WORK/rootfs

    log INFO "Installing Ubuntu noble amd64 rootfs"
    as_root debootstrap --arch=amd64 \
        "$UBUNTU_RELEASE" "$root" "$UBUNTU_MIRROR" \
        /usr/share/debootstrap/scripts/gutsy
    mount_fs "$root/dev" --bind /dev
    mount_fs "$root/dev/pts" --bind /dev/pts
    mount_fs "$root/proc" -t proc proc
    mount_fs "$root/sys" -t sysfs sysfs
    as_root cp -L /etc/resolv.conf "$root/etc/resolv.conf"
    printf '#!/bin/sh\nexit 101\n' |
        as_root tee "$root/usr/sbin/policy-rc.d" >/dev/null
    as_root chmod 0755 "$root/usr/sbin/policy-rc.d"
    as_root chroot "$root" /usr/bin/env -i \
        PATH=/usr/sbin:/usr/bin:/sbin:/bin DEBIAN_FRONTEND=noninteractive \
        RELEASE="$UBUNTU_RELEASE" MIRROR="$UBUNTU_MIRROR" /bin/bash -eu <<'EOF'
cat >/etc/apt/sources.list <<APT
deb ${MIRROR} ${RELEASE} main universe restricted multiverse
deb ${MIRROR} ${RELEASE}-updates main universe restricted multiverse
deb ${MIRROR} ${RELEASE}-security main universe restricted multiverse
APT
apt-get update
apt-get install -y --no-install-recommends \
    linux-image-generic linux-firmware initramfs-tools \
    grub-efi-amd64-bin grub-common \
    systemd-sysv systemd-resolved systemd-timesyncd udev kmod dbus \
    openssh-server sudo locales tzdata ca-certificates curl \
    iproute2 iputils-ping ethtool tcpdump i2c-tools \
    openvswitch-switch openvpn iptables net-tools
apt-get clean
rm -rf /var/lib/apt/lists/*
sed -i 's/^# *en_US.UTF-8 UTF-8/en_US.UTF-8 UTF-8/' /etc/locale.gen
locale-gen
echo 'LANG=en_US.UTF-8' >/etc/default/locale
ln -sf /usr/share/zoneinfo/UTC /etc/localtime
echo ukama-linux >/etc/hostname
cat >/etc/hosts <<HOSTS
127.0.0.1 localhost
127.0.1.1 ukama-linux
::1 localhost ip6-localhost
HOSTS
echo 'root:root' | chpasswd
if ! id ukama >/dev/null 2>&1; then
    useradd -m -s /bin/bash -G sudo ukama
fi
echo 'ukama:ukama' | chpasswd
echo 'ukama ALL=(ALL) NOPASSWD: ALL' >/etc/sudoers.d/ukama
chmod 0440 /etc/sudoers.d/ukama

mkdir -p /etc/systemd/network
cat >/etc/systemd/network/20-wired.network <<NET
[Match]
Name=en* eth*

[Network]
DHCP=yes
NET
systemctl enable systemd-networkd systemd-resolved systemd-timesyncd \
    ssh serial-getty@ttyS0.service
# init-network owns ovsdb-server and ovs-vswitchd, using the packaged schema.
# Prevent a second OVS supervisor from racing it for sockets and DB locks.
systemctl mask openvswitch-switch.service \
    ovsdb-server.service ovs-vswitchd.service
EOF
    as_root rsync -a --chown=0:0 "$WORK/overlay/" "$root/"
    as_root chroot "$root" systemctl enable starterd.service ukama-sysfs.service
    # Generate fresh machine and SSH identities on each board's first boot.
    cat >"$WORK/ssh-keygen.service" <<'EOF'
[Unit]
Description=Create per-board SSH host keys
Before=ssh.service

[Service]
Type=oneshot
ExecStart=/usr/bin/ssh-keygen -A
RemainAfterExit=yes

[Install]
RequiredBy=ssh.service
EOF
    as_root install -m 0644 "$WORK/ssh-keygen.service" \
        "$root/etc/systemd/system/ssh-keygen.service"
    as_root chroot "$root" systemctl enable ssh-keygen.service
    as_root chroot "$root" /bin/bash -eu <<'EOF'
rm -f /etc/ssh/ssh_host_* /var/lib/dbus/machine-id /var/lib/systemd/random-seed
truncate -s 0 /etc/machine-id
ln -s /etc/machine-id /var/lib/dbus/machine-id
rm -f /usr/sbin/policy-rc.d /etc/resolv.conf
ln -s /run/systemd/resolve/stub-resolv.conf /etc/resolv.conf
# Do not capture the build host's swap UUID in the target initramfs.
mkdir -p /etc/initramfs-tools/conf.d
echo RESUME=none >/etc/initramfs-tools/conf.d/resume
# Include COM storage support even when the build host boots from NVMe/SATA.
cat >>/etc/initramfs-tools/modules <<MODULES
mmc_block
sdhci
sdhci_acpi
sdhci_pci
ext4
MODULES
update-initramfs -u -k all
EOF
    release_mounts
}

prepare_boot() {
    local root="$WORK/rootfs"
    local boot="$WORK/boot"
    local kernel
    local version
    local modules

    log INFO "Preparing UEFI boot, primary/passive menus and serial console"
    kernel=$(find "$root/boot" -maxdepth 1 -type f -name 'vmlinuz-*' |
        sort -V | tail -n 1)
    if [ -z "$kernel" ]; then
        die "Ubuntu kernel is missing"
    fi
    version=${kernel##*/vmlinuz-}
    if [ ! -s "$root/boot/initrd.img-$version" ]; then
        die "Missing initramfs for $version"
    fi
    mkdir -p "$boot/boot/grub" "$boot/EFI/BOOT"
    as_root cp "$kernel" "$boot/boot/vmlinuz"
    as_root cp "$root/boot/initrd.img-$version" "$boot/boot/initrd.img"
    printf '%s\n' "$version" >"$boot/boot/kernel-version"
    cat >"$boot/boot/grub/grub.cfg" <<'EOF'
insmod part_msdos
insmod fat
search --no-floppy --label boot --set=root
serial --unit=0 --speed=115200 --word=8 --parity=no --stop=1
terminal_input console serial
terminal_output console serial
set default=0
set timeout=3

menuentry 'Ukama COM - primary' {
    linux /boot/vmlinuz root=LABEL=primary rootfstype=ext4 \
        rootwait rw console=tty1 console=ttyS0,115200n8
    initrd /boot/initrd.img
}
menuentry 'Ukama COM - passive' {
    linux /boot/vmlinuz root=LABEL=passive rootfstype=ext4 \
        rootwait rw console=tty1 console=ttyS0,115200n8
    initrd /boot/initrd.img
}
EOF
    cat >"$WORK/grub-embed.cfg" <<'EOF'
insmod part_msdos
insmod fat
search --no-floppy --label boot --set=root
configfile ($root)/boot/grub/grub.cfg
EOF
    as_root cp "$WORK/grub-embed.cfg" "$root/tmp/grub-embed.cfg"
    as_root cp "$boot/boot/grub/grub.cfg" "$root/tmp/grub.cfg"
    as_root chroot "$root" grub-script-check /tmp/grub.cfg
    modules="part_msdos fat normal linux search search_label"
    modules="$modules configfile serial terminal"
    as_root chroot "$root" grub-mkstandalone -O x86_64-efi \
        --modules="$modules" \
        -o /tmp/BOOTX64.EFI 'boot/grub/grub.cfg=/tmp/grub-embed.cfg'
    as_root cp "$root/tmp/BOOTX64.EFI" "$boot/EFI/BOOT/BOOTX64.EFI"
    as_root rm -f "$root/tmp/BOOTX64.EFI" \
        "$root/tmp/grub-embed.cfg" "$root/tmp/grub.cfg"
}

write_fstab() {
    local root=$1
    local label=$2

    as_root tee "$root/etc/fstab" >/dev/null <<EOF
LABEL=$label /     ext4 defaults 0 1
LABEL=boot   /boot vfat defaults 0 2
LABEL=swap   none  swap sw       0 0
tmpfs       /tmp  tmpfs defaults,nosuid,nodev,size=256M 0 0
EOF
}

assemble_image() {
    local part
    local attempt
    local used

    # Leave at least 512 MiB for starter.d to unpack apps on the smaller root.
    # The validator also budgets each package's uncompressed size on both roots.
    used=$(as_root du -sx --block-size=1 "$WORK/rootfs" | awk '{print $1}')
    if [ "$used" -ge $((3584 * 1024 * 1024)) ]; then
        die "Rootfs exceeds passive partition budget (3.5 GiB)"
    fi
    log INFO "Assembling the 16 GiB image"
    truncate -s 16G "$WORK/image.img"
    sfdisk "$WORK/image.img" <<'EOF'
label: dos
unit: sectors

start=2048,     size=1048576,  type=ef, bootable
start=1050624,  size=8388608,  type=83
start=9439232,  size=15728640, type=83
start=25167872, size=8386560,  type=82
EOF
    LOOPDEV=$(as_root losetup --find --show --partscan "$WORK/image.img")
    for part in 1 2 3 4; do
        attempt=0
        while [ "$attempt" -lt 50 ]; do
            if [ -b "${LOOPDEV}p$part" ]; then
                break
            fi
            sleep 0.1
            attempt=$((attempt + 1))
        done
        if [ ! -b "${LOOPDEV}p$part" ]; then
            die "Missing loop partition ${LOOPDEV}p$part"
        fi
    done
    as_root mkfs.vfat -F 32 -n boot "${LOOPDEV}p1"
    as_root mkfs.ext4 -F -L passive -O '^64bit,^metadata_csum' "${LOOPDEV}p2"
    as_root mkfs.ext4 -F -L primary -O '^64bit,^metadata_csum' "${LOOPDEV}p3"
    as_root mkswap -L swap "${LOOPDEV}p4"
    mount_fs "$WORK/mnt-boot" "${LOOPDEV}p1"
    mount_fs "$WORK/mnt-passive" "${LOOPDEV}p2"
    mount_fs "$WORK/mnt-primary" "${LOOPDEV}p3"
    as_root cp -r "$WORK/boot/." "$WORK/mnt-boot/"
    for part in primary passive; do
        as_root rsync -aHAX --numeric-ids "$WORK/rootfs/" "$WORK/mnt-$part/"
        write_fstab "$WORK/mnt-$part" "$part"
    done
    sync
    release_mounts
}

main() {
    parse_args "$@"
    preflight
    build_apps
    stage_ukama
    install_ubuntu
    prepare_boot
    assemble_image
    log INFO "Validating the finished image before publishing it"
    as_root bash "$SCRIPT_DIR/validate-com-image.sh" \
        "$WORK/image.img" --repo "$UKAMA_ROOT"
    cleanup_app_build
    mv -f -- "$WORK/image.img" "$OUTPUT"
    log SUCCESS "COM image built and statically validated: $OUTPUT"
}

if [ "${BASH_SOURCE[0]}" = "$0" ]; then
    main "$@"
fi
