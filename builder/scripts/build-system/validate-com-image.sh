#!/usr/bin/env bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
# Copyright (c) 2026-present, Ukama Inc.

# Offline validation of a COM image. No repairs, app execution, backend calls
# or flashing. Rootfs mounts and the loop device are read-only. A passing result
# establishes image structure, not firmware, NIC, EPC or on-board boot behavior.
# The builder also sources the package helpers below.
# Literal loader variables below must stay unexpanded.
# shellcheck disable=SC2016

com_valid_name() {
    case "$1" in
        ''|*[!A-Za-z0-9_.-]*)
            return 1
            ;;
    esac
    return 0
}

com_app_space() {
    local manifest=$1
    local name=$2

    jq -er --arg name "$name" '
        .spaces[] |
        select(any(.apps[]; .name == $name)) |
        .name
    ' "$manifest"
}

com_config_dir() {
    local package=$1
    local name=$2
    local entry

    if [ -d "$package/ukama/configs/$name" ]; then
        printf 'ukama/configs/%s\n' "$name"
    elif [ -d "$package/conf" ]; then
        entry=$(find "$package/conf" -mindepth 1 -print -quit)
        if [ -n "$entry" ]; then
            echo conf
        fi
    fi
    return 0
}

com_resolve_file() {
    # Resolve image symlinks without following them into the host filesystem.
    local root=$1
    local path=${2#/}
    local current=""
    local depth=0
    local part
    local link

    while [ -n "$path" ]; do
        part=${path%%/*}
        if [[ "$path" == */* ]]; then
            path=${path#*/}
        else
            path=""
        fi
        case "$part" in
            ''|.)
                continue
                ;;
            ..)
                if [[ "$current" == */* ]]; then
                    current=${current%/*}
                else
                    current=""
                fi
                continue
                ;;
        esac

        if [ -L "$root/$current/$part" ]; then
            depth=$((depth + 1))
            if [ "$depth" -gt 40 ]; then
                return 1
            fi
            if ! link=$(readlink "$root/$current/$part"); then
                return 1
            fi
            if [[ "$link" == /* ]]; then
                current=""
                link=${link#/}
            fi
            if [ -n "$path" ]; then
                path="$link/$path"
            else
                path=$link
            fi
        elif [ -n "$current" ]; then
            current="$current/$part"
        else
            current=$part
        fi
    done
    printf '%s/%s\n' "$root" "$current"
}

com_safe_member() {
    local member=$1

    case "$member" in
        ''|/*|*[!A-Za-z0-9_./+@=-]*)
            return 1
            ;;
    esac
    case "/$member/" in
        */../*|*/./*)
            return 1
            ;;
    esac
    return 0
}

com_unpack_package() {
    local archive=$1
    local destination=$2
    local member
    local top=""
    local entries

    if [ ! -s "$archive" ]; then
        return 1
    fi
    if ! mkdir -p "$destination"; then
        return 1
    fi
    if ! gzip -t "$archive"; then
        return 1
    fi
    if ! tar -tzf "$archive" >"$destination.members"; then
        return 1
    fi
    while IFS= read -r member; do
        if ! com_safe_member "$member"; then
            return 1
        fi
        if [ -z "$top" ]; then
            top=${member%%/*}
        fi
        if [ "${member%%/*}" != "$top" ]; then
            return 1
        fi
    done <"$destination.members"
    if [ -z "$top" ]; then
        return 1
    fi
    if ! tar --no-same-owner --no-same-permissions \
        --delay-directory-restore -xzf "$archive" -C "$destination"; then
        return 1
    fi

    COM_TOP=$top
    COM_PAYLOAD="$destination/$top"
    if [ ! -d "$COM_PAYLOAD" ]; then
        return 1
    fi
    if [ -L "$COM_PAYLOAD" ]; then
        return 1
    fi
    if [ ! -s "$COM_PAYLOAD/VERSION" ]; then
        return 1
    fi
    if [ -L "$COM_PAYLOAD/VERSION" ]; then
        return 1
    fi
    COM_VERSION=$(cat "$COM_PAYLOAD/VERSION")
    if ! com_valid_name "$COM_VERSION"; then
        return 1
    fi
    case "$COM_VERSION" in
        latest|[!A-Za-z0-9]*)
            return 1
            ;;
    esac
    entries=$(find "$destination" -mindepth 1 -maxdepth 1 | wc -l)
    if [ "$entries" -ne 1 ]; then
        return 1
    fi
    return 0
}

com_executable() {
    local file=$1
    local match

    if [ ! -s "$file" ]; then
        return 1
    fi
    # Test stored mode bits; the filesystem may be mounted with noexec.
    match=$(find -L "$file" -maxdepth 0 -type f -perm /111 -print)
    if [ -z "$match" ]; then
        return 1
    fi
    return 0
}

com_elf() {
    local file=$1
    local header

    if ! header=$(readelf -hW "$file"); then
        return 1
    fi
    if [[ "$header" != *ELF64* ]]; then
        return 1
    fi
    if [[ "$header" != *'Advanced Micro Devices X86-64'* ]]; then
        return 1
    fi
    return 0
}

com_dynamic_values() {
    local file=$1
    local tag=$2

    readelf -dW "$file" | awk -F '[][]' -v tag="($tag)" '
        index($0, tag) { print $2 }
    '
}

com_library_directory() {
    local root=$1
    local origin=$2
    local directory=$3

    # ORIGIN is a literal ELF token, relative to the inspected file.
    case "$directory" in
        *'${ORIGIN}'*)
            directory=${directory//\$\{ORIGIN\}/$origin}
            ;;
        *'$ORIGIN'*)
            directory=${directory//\$ORIGIN/$origin}
            ;;
        /*)
            directory="$root$directory"
            ;;
        *)
            return 1
            ;;
    esac
    printf '%s\n' "$directory"
}

com_find_library() {
    local root=$1
    local package=$2
    local file=$3
    local needed=$4
    local origin
    local runpath
    local directory
    local candidate
    local relative
    local paths=()
    local entries=()

    case "$needed" in
        libc.musl-x86_64.so.1|ld-musl-x86_64.so.1)
            com_resolve_file "$root" /lib/ld-musl-x86_64.so.1
            return
            ;;
        */*)
            # Shared objects must be relocatable into the image runtime.
            return 1
            ;;
    esac

    origin=$(dirname -- "$file")
    runpath=$(com_dynamic_values "$file" RUNPATH)
    if [ -z "$runpath" ]; then
        runpath=$(com_dynamic_values "$file" RPATH)
    fi
    IFS=: read -r -a entries <<<"$runpath"
    for directory in "${entries[@]}"; do
        if directory=$(com_library_directory \
            "$root" "$origin" "$directory"); then
            paths+=("$directory")
        fi
    done
    paths+=("$root/opt/ukama/runtime/lib")

    for directory in "${paths[@]}"; do
        case "$directory" in
            "$package"/*)
                relative=${directory#"$package"}
                candidate=$(com_resolve_file \
                    "$package" "$relative/$needed")
                ;;
            "$root"/*)
                relative=${directory#"$root"}
                candidate=$(com_resolve_file "$root" "$relative/$needed")
                ;;
            *)
                continue
                ;;
        esac
        if [ -s "$candidate" ]; then
            printf '%s\n' "$candidate"
            return 0
        fi
    done
    printf 'Unresolved Alpine dependency %s in %s\n' "$needed" "$file" >&2
    return 1
}

com_contains() {
    local wanted=$1
    local item

    shift
    for item in "$@"; do
        if [ "$item" = "$wanted" ]; then
            return 0
        fi
    done
    return 1
}

com_interpreter() {
    readelf -lW "$1" | awk '
        /Requesting program interpreter:/ {
            sub(/^.*interpreter: /, "")
            sub(/]$/, "")
            print
        }
    '
}

com_dependencies() {
    local root=$1
    local package=$2
    local executable=$3
    local loader=/lib/ld-musl-x86_64.so.1
    local file
    local interpreter
    local libraries
    local needed
    local found
    local index=0
    local queue=("$executable")
    local checked=()

    while [ "$index" -lt "${#queue[@]}" ]; do
        file=${queue[index]}
        index=$((index + 1))
        if com_contains "$file" "${checked[@]}"; then
            continue
        fi
        checked+=("$file")
        if ! com_elf "$file"; then
            echo "Invalid x86_64 ELF: $file"
            return 1
        fi

        interpreter=$(com_interpreter "$file")
        if [ -n "$interpreter" ]; then
            if [ "$interpreter" != "$loader" ]; then
                echo "Expected an Alpine x86_64 app: $file"
                echo "Found interpreter: $interpreter"
                return 1
            fi
            if ! found=$(com_resolve_file "$root" "$loader"); then
                return 1
            fi
            if ! com_executable "$found"; then
                return 1
            fi
            queue+=("$found")
        fi

        libraries=$(com_dynamic_values "$file" NEEDED)
        while IFS= read -r needed; do
            if [ -z "$needed" ]; then
                continue
            fi
            case "$needed" in
                libc.so.6|ld-linux*)
                    echo "glibc dependency in Alpine app runtime: $file"
                    return 1
                    ;;
            esac
            if ! found=$(com_find_library \
                "$root" "$package" "$file" "$needed"); then
                return 1
            fi
            queue+=("$found")
        done <<<"$libraries"
    done
    return 0
}

com_runtime_path() {
    local value

    if ! value=$(cat "$1/etc/ld-musl-x86_64.path"); then
        return 1
    fi
    if [ "$value" != /opt/ukama/runtime/lib ]; then
        return 1
    fi
    return 0
}

validate_app_runtime() {
    local root=$1
    local label=$2
    local loader
    local library
    local runtime="$root/opt/ukama/runtime/lib"

    validate_check "$label private Alpine library path" \
        com_runtime_path "$root"
    loader=$(com_resolve_file "$root" /lib/ld-musl-x86_64.so.1)
    validate_check "$label musl loader executable" com_executable "$loader"
    validate_check "$label musl loader matches app runtime" \
        cmp "$loader" "$runtime/ld-musl-x86_64.so.1"
    validate_check "$label starter uses the system musl interpreter" \
        grep -Fx 'exec "$root/sbin/starter.d" "$@"' \
        "$root/usr/sbin/starter.d"
    validate_check "$label no global app library override" \
        grep -Fx 'unset LD_LIBRARY_PATH LD_PRELOAD' \
        "$root/usr/sbin/starter.d"

    for library in "$runtime"/*.so*; do
        library=$(com_resolve_file "$root" "${library#"$root"}")
        validate_check "$label Alpine runtime ${library##*/}" \
            com_dependencies "$root" "$root" "$library"
    done
}

validate_usage() {
    cat <<'EOF'
Usage: sudo ./validate-com-image.sh IMAGE [--repo UKAMA_ROOT]

Checks the 16 GiB COM layout, filesystems, UEFI boot files, both Ubuntu
roots, kernel/modules, services, manifest, packaged Alpine x86_64 apps,
musl runtime, config links and space for first-boot app extraction.
Exit 0: static checks passed (review any WARN lines). Exit 1: a check failed.
Host packages: dosfstools e2fsprogs fdisk jq binutils initramfs-tools-core zstd.

The repository supplies the board app list, package recipes and /etc/services.
If build/pkgs outputs are present, their bytes must match the image packages.
Checks do not execute image code or prove hardware boot, symbol-version ABI
compatibility, node registration, physical interface wiring or EPC connectivity.
Secure Boot must be disabled for this unsigned standalone GRUB image.
EOF
}

validate_die() {
    printf 'ERROR: %s\n' "$*" >&2
    exit 1
}
validate_warn() {
    printf '\033[1;33mWARN\033[0m: %s\n' "$*"
    WARNINGS=$((WARNINGS + 1))
}

validate_check() {
    local description=$1

    shift
    if "$@" >"$VWORK/check.log" 2>&1; then
        printf '\033[1;32mPASS\033[0m: %s\n' "$description"
        PASSED=$((PASSED + 1))
        CHECK_PASSED=yes
    else
        printf '\033[1;31mFAIL\033[0m: %s\n' "$description" >&2
        sed 's/^/      /' "$VWORK/check.log" >&2
        FAILED=$((FAILED + 1))
        CHECK_PASSED=no
    fi

    # Keep checking after a failure. validate_main returns the final status.
    # CHECK_PASSED tells callers whether dependent work can proceed.
    return 0
}

validate_cleanup() {
    local status=$?
    local safe=yes
    local index

    trap - EXIT
    index=${#VMOUNTS[@]}
    while [ "$index" -gt 0 ]; do
        index=$((index - 1))
        if ! umount "${VMOUNTS[index]}"; then
            safe=no
        fi
    done
    if [ "$safe" = yes ] && [ -n "$VLOOP" ]; then
        if ! losetup --detach "$VLOOP"; then
            safe=no
        fi
    fi
    if [ "$safe" = yes ]; then
        if [ -n "$VWORK" ]; then
            rm -rf -- "$VWORK"
        fi
    else
        printf 'ERROR: Cleanup incomplete; retained %s and %s\n' \
            "$VWORK" "$VLOOP" >&2
        status=1
    fi
    exit "$status"
}

validate_partition_table() {
    local size

    size=$(stat -Lc %s "$VIMAGE")
    if [ "$size" -ne 17179869184 ]; then
        validate_die "Image must be exactly 16 GiB"
    fi
    if ! sfdisk --verify "$VIMAGE"; then
        validate_die "Invalid partition table"
    fi
    sfdisk --json "$VIMAGE" >"$VWORK/partitions.json"
    if ! jq -e '
        .partitiontable |
        .label == "dos" and .sectorsize == 512 and
        ([.partitions[].start] == [2048,1050624,9439232,25167872]) and
        ([.partitions[].size] == [1048576,8388608,15728640,8386560]) and
        ([.partitions[].type] == ["ef","83","83","82"])
    ' "$VWORK/partitions.json" >/dev/null; then
        validate_die "Unexpected COM partition layout"
    fi
}

validate_partition() {
    local number=$1
    local label=$2
    local type=$3
    local device="${VLOOP}p$number"
    local attempt=0
    local actual

    while [ "$attempt" -lt 50 ]; do
        if [ -b "$device" ]; then
            break
        fi
        sleep 0.1
        attempt=$((attempt + 1))
    done
    if [ ! -b "$device" ]; then
        validate_die "Missing $device"
    fi

    actual=$(blkid -p -s TYPE -o value "$device")
    validate_check "$label filesystem type" test "$actual" = "$type"
    actual=$(blkid -p -s LABEL -o value "$device")
    validate_check "$label filesystem label" test "$actual" = "$label"
    case "$type" in
        vfat)
            validate_check "$label filesystem integrity" \
                fsck.fat -n "$device"
            ;;
        ext4)
            validate_check "$label filesystem integrity" \
                e2fsck -fn "$device"
            ;;
    esac
}

validate_mount() {
    local number=$1
    local label=$2
    local options=ro,nodev,nosuid,noexec

    if [ "$number" -ne 1 ]; then
        options="$options,noload"
    fi
    mkdir -p "$VWORK/$label"
    mount -o "$options" "${VLOOP}p$number" "$VWORK/$label"
    VMOUNTS+=("$VWORK/$label")
}

validate_layout() {
    local read_only

    validate_partition_table
    VLOOP=$(losetup --find --show --read-only --partscan "$VIMAGE")
    read_only=$(blockdev --getro "$VLOOP")
    if [ "$read_only" -ne 1 ]; then
        validate_die "Loop device is not read-only"
    fi

    validate_partition 1 boot vfat
    validate_partition 2 passive ext4
    validate_partition 3 primary ext4
    validate_partition 4 swap swap
    if [ "$FAILED" -ne 0 ]; then
        validate_die "Partition/filesystem checks failed; refusing to mount"
    fi
    validate_mount 1 boot
    validate_mount 2 passive
    validate_mount 3 primary
}

validate_json_equal() {
    local actual=$1
    local expected=$2

    if ! jq -S . "$actual" >"$VWORK/actual.json"; then
        return 1
    fi
    if ! jq -S . "$expected" >"$VWORK/expected.json"; then
        return 1
    fi
    diff -u "$VWORK/actual.json" "$VWORK/expected.json"
}

validate_link() {
    local link=$1
    local expected=$2
    local actual

    if ! actual=$(readlink "$link"); then
        return 1
    fi
    if [ "$actual" != "$expected" ]; then
        return 1
    fi
    return 0
}

validate_enabled() {
    local root=$1
    local unit=$2
    local target=${3:-}
    local link
    local file
    local relative

    if [ -n "$target" ]; then
        link="$root/etc/systemd/system/$target/$unit"
        if [ ! -L "$link" ]; then
            return 1
        fi
        relative=${link#"$root"}
        if ! file=$(com_resolve_file "$root" "$relative"); then
            return 1
        fi
        if [ ! -s "$file" ]; then
            return 1
        fi
        return 0
    fi

    for link in "$root/etc/systemd/system/"*.wants/"$unit"; do
        if [ ! -L "$link" ]; then
            continue
        fi
        relative=${link#"$root"}
        if ! file=$(com_resolve_file "$root" "$relative"); then
            continue
        fi
        if [ -s "$file" ]; then
            return 0
        fi
    done
    return 1
}

validate_fstab() {
    local root=$1
    local label=$2

    awk -v label="$label" '
        /^[[:space:]]*#/ { next }
        $2 == "/" {
            count++
            correct_root = ($1 == "LABEL=" label && $3 == "ext4")
        }
        $1 == "LABEL=boot" && $2 == "/boot" && $3 == "vfat" { boot++ }
        $1 == "LABEL=swap" && $3 == "swap" { swap++ }
        END { exit !(count == 1 && correct_root && boot == 1 && swap == 1) }
    ' "$root/etc/fstab"
}

validate_no_ssh_keys() {
    local root=$1
    local key

    key=$(find "$root/etc/ssh" -maxdepth 1 -name 'ssh_host_*' -print -quit)
    if [ -n "$key" ]; then
        return 1
    fi
    return 0
}

validate_services() {
    local root=$1
    local label=$2
    local unit
    local units=(
        starterd.service
        ukama-sysfs.service
        systemd-networkd.service
        systemd-resolved.service
        systemd-timesyncd.service
        ssh.service
    )

    for unit in "${units[@]}"; do
        validate_check "$label enabled $unit" \
            validate_enabled "$root" "$unit"
    done
    validate_check "$label serial debugging enabled" \
        validate_enabled "$root" serial-getty@ttyS0.service getty.target.wants
    validate_check "$label per-board SSH keys enabled" \
        validate_enabled "$root" ssh-keygen.service ssh.service.requires

    units=(
        openvswitch-switch.service
        ovsdb-server.service
        ovs-vswitchd.service
    )
    for unit in "${units[@]}"; do
        validate_check "$label init-network owns $unit" \
            validate_link "$root/etc/systemd/system/$unit" /dev/null
    done
    unit="$root/etc/systemd/system/starterd.service"
    validate_check "$label starter manifest environment" \
        grep -qx 'Environment=STARTERD_MANIFEST=/ukama/manifest.json' "$unit"
    validate_check "$label stable per-boot lifecycle identity" \
        grep -Fq '/proc/sys/kernel/random/boot_id' "$unit"
}

validate_kernel_files() {
    local root=$1
    local label=$2
    local modules="$root/lib/modules/$KERNEL"
    local config="$root/boot/config-$KERNEL"
    local list="$VWORK/$label.modules.list"

    validate_check "$label kernel matches ESP" \
        cmp "$root/boot/vmlinuz-$KERNEL" "$VWORK/boot/boot/vmlinuz"
    validate_check "$label initramfs matches ESP" \
        cmp "$root/boot/initrd.img-$KERNEL" "$VWORK/boot/boot/initrd.img"
    validate_check "$label kernel module index" test -s "$modules/modules.dep"
    validate_check "$label kernel configuration" test -s "$config"
    if [ "$label" = passive ]; then
        validate_check "$label kernel configuration matches primary" \
            cmp "$config" "$VWORK/primary/boot/config-$KERNEL"
    fi

    validate_check "$label kernel module files readable" \
        find "$modules" -type f -name '*.ko*' -print
    if [ "$CHECK_PASSED" = yes ]; then
        cp "$VWORK/check.log" "$list"
    else
        : >"$list"
    fi
    validate_check "$label kernel OVS support" validate_driver \
        "$root" openvswitch CONFIG_OPENVSWITCH "$list"
    validate_check "$label kernel connection tracking support" \
        validate_driver "$root" nf_conntrack CONFIG_NF_CONNTRACK "$list"
    validate_check "$label kernel TUN/TAP support (required)" \
        validate_driver "$root" tun CONFIG_TUN "$list"
}

validate_os() {
    local root=$1
    local label=$2
    local services="$VREPO/nodes/ukamaOS/distro/scripts/files/services"
    local property="$root/ukama/mocksysfs/mfgdata/property/property.json"
    local path
    local file
    local executables=(
        /sbin/init
        /bin/sh
        /sbin/starter.d
        /usr/bin/ssh-keygen
        /usr/sbin/sshd
        /usr/sbin/ovsdb-server
        /usr/sbin/ovs-vswitchd
        /usr/bin/ovs-vsctl
        /usr/bin/ovs-ofctl
        /usr/sbin/ip
        /usr/sbin/iptables
    )

    validate_check "$label Ubuntu noble" \
        grep -qx 'VERSION_CODENAME=noble' "$root/etc/os-release"
    validate_check "$label independent fstab" \
        validate_fstab "$root" "$label"
    validate_check "$label board manifest" validate_json_equal \
        "$root/ukama/manifest.json" "$VWORK/manifest.json"
    validate_check "$label Ukama service registry" \
        cmp "$root/etc/services" "$services"
    validate_check "$label bootstrap server" test -s "$root/ukama/bootstrap"
    validate_check "$label noded property config" \
        jq -e 'type == "object" and length > 0' "$property"
    validate_check "$label physical sysfs link" \
        validate_link "$root/ukama/mocksysfs/sys" /sys

    for path in "${executables[@]}"; do
        if ! file=$(com_resolve_file "$root" "$path"); then
            file=""
        fi
        validate_check "$label executable $path" com_executable "$file"
    done
    validate_check "$label OVS schema" jq -e '.name == "Open_vSwitch"' \
        "$root/usr/share/openvswitch/vswitch.ovsschema"
    validate_services "$root" "$label"
    validate_check "$label DHCP enabled" grep -qx 'DHCP=yes' \
        "$root/etc/systemd/network/20-wired.network"
    validate_check "$label empty machine identity" \
        test ! -s "$root/etc/machine-id"
    validate_check "$label no pre-generated SSH keys" \
        validate_no_ssh_keys "$root"
    validate_kernel_files "$root" "$label"
}

validate_config_links() {
    local root=$1
    local name=$2
    local package=$3
    local space=$4
    local directory
    local expected
    local arguments
    local argument
    local relative

    directory=$(com_config_dir "$package" "$name")
    if [ -n "$directory" ]; then
        expected="/ukama/apps/$space/$name/current/$directory"
        if ! validate_link "$root/ukama/configs/$name" "$expected"; then
            return 1
        fi
    fi

    # Check absolute config arguments against what starter.d will unpack.
    arguments=$(jq -r --arg name "$name" '
        .spaces[].apps[] | select(.name == $name) | .argv[]
    ' "$root/ukama/manifest.json")
    while IFS= read -r argument; do
        case "$argument" in
            /ukama/configs/"$name"/*)
                if [ -z "$directory" ]; then
                    return 1
                fi
                relative=${argument#"/ukama/configs/$name/"}
                if [ ! -s "$package/$directory/$relative" ]; then
                    return 1
                fi
                ;;
            /ukama/configs/*)
                return 1
                ;;
        esac
    done <<<"$arguments"
    return 0
}

validate_profile_warnings() {
    local name=$1
    local package=$2
    local config

    case "$name" in
        metricsd)
            config="$package/conf/config.toml"
            if grep -Eq 'ukamnode[[:space:]]*=[[:space:]]*"TRX"' "$config"; then
                validate_warn \
                    "metricsd identifies TRX; confirm the COM profile"
            fi
            ;;
        aggregator)
            config="$package/conf/config.toml"
            if ! grep -q '127.0.0.10' "$config"; then
                validate_warn \
                    "Aggregator has no 127.0.0.10 TRX source; review coverage"
            fi
            ;;
        init-network)
            config="$package/ukama/configs/init-network/config.toml"
            if grep -Eq 'virtual tower|EPCEMU' "$config"; then
                validate_warn \
                    "init-network has virtual-node settings; review COM NICs"
            fi
            ;;
    esac
}

validate_helpers() {
    local root=$1
    local package=$2
    local executable=$3
    local label=$4
    local extra

    # Include helper binaries such as wimc-agent. Libraries are checked
    # through each executable's DT_NEEDED entries in com_dependencies.
    find "$package" -type f \
        \( -path "$package/sbin/*" -o -path "$package/bin/*" \) \
        -print0 >"$VWORK/helpers.list"
    while IFS= read -r -d '' extra; do
        if [ "$extra" = "$executable" ]; then
            continue
        fi
        if readelf -hW "$extra" >/dev/null 2>&1; then
            validate_check "$label helper ${extra##*/}" \
                com_dependencies "$root" "$package" "$extra"
        fi
    done <"$VWORK/helpers.list"
}

validate_apps() {
    local key
    local config
    local name
    local tag
    local binary
    local path
    local label
    local root
    local archive
    local package
    local reference
    local space
    local file
    local size
    local bytes=0
    local available
    local required

    for label in primary passive; do
        validate_app_runtime "$VWORK/$label" "$label"
    done

    for key in "${APPS[@]}"; do
        config="$VREPO/builder/configs/$key.toml"
        name=$(_manifest_toml_value "$config" capp-exec name)
        tag=$(_manifest_toml_value "$config" capp-exec version)
        binary=$(_manifest_toml_value "$config" capp-exec bin)
        path=$(_manifest_toml_value "$config" capp-exec path)
        if [ -z "$name" ]; then
            name=$key
        fi
        if [ -z "$tag" ]; then
            tag=latest
        fi
        if [ -z "$path" ]; then
            path=/sbin
        fi
        if ! com_valid_name "$name"; then
            validate_die "Invalid app name in $config"
        fi
        if ! com_valid_name "$tag"; then
            validate_die "Invalid app tag in $config"
        fi

        archive="$VWORK/primary/ukama/apps/pkgs/${name}_${tag}.tar.gz"
        validate_check "$name package integrity" \
            com_unpack_package "$archive" "$VWORK/unpack/$name"
        if [ "$CHECK_PASSED" = no ]; then
            continue
        fi
        package=$COM_PAYLOAD
        validate_check "$name package name and VERSION" \
            test "$COM_TOP" = "${name}_$COM_VERSION"
        if [ "$tag" != latest ]; then
            validate_check "$name requested version" \
                test "$tag" = "$COM_VERSION"
        fi

        reference="$VREPO/build/pkgs/${name}_${tag}.tar.gz"
        if [ -s "$reference" ]; then
            validate_check "$name matches build output" \
                cmp "$archive" "$reference"
        fi
        size=$(du -s --block-size=1 "$package" | awk '{print $1}')
        bytes=$((bytes + size))
        validate_profile_warnings "$name" "$package"

        for label in primary passive; do
            root="$VWORK/$label"
            reference="$root/ukama/apps/pkgs/${name}_${tag}.tar.gz"
            validate_check "$label $name package bytes" \
                cmp "$archive" "$reference"
            if ! file=$(com_resolve_file "$package" "${path%/}/$binary"); then
                file=""
            fi
            validate_check "$label $name executable" com_executable "$file"
            validate_check "$label $name ELF and shared libraries" \
                com_dependencies "$root" "$package" "$file"
            validate_helpers "$root" "$package" "$file" "$label $name"

            if [ "$key" = starterd ]; then
                validate_check "$label complete starter.d runtime" \
                    diff -qr "$package" "$root/opt/ukama/starterd"
            else
                space=$(com_app_space "$root/ukama/manifest.json" "$name")
                validate_check "$label $name packaged config paths" \
                    validate_config_links \
                    "$root" "$name" "$package" "$space"
            fi
        done
    done

    # Archives remain after unpacking. Allow unpacked size, 20% allocation
    # overhead and another 512 MiB on each root partition.
    required=$((bytes + bytes / 5 + 512 * 1024 * 1024))
    for label in primary passive; do
        available=$(df -B1 --output=avail "$VWORK/$label" | awk 'NR == 2')
        validate_check "$label space for app extraction and 512 MiB reserve" \
            test "$available" -ge "$required"
    done
}

validate_efi_loader() {
    local loader=$1

    if ! objdump -p "$loader" >"$VWORK/efi-header.txt"; then
        return 1
    fi
    if ! objdump -f "$loader" >"$VWORK/efi-format.txt"; then
        return 1
    fi
    if ! grep -q 'pei-x86-64' "$VWORK/efi-format.txt"; then
        return 1
    fi
    if ! awk '
        $1 == "Subsystem" && $2 == "0000000a" { found = 1 }
        END { exit !found }
    ' "$VWORK/efi-header.txt"; then
        return 1
    fi
    return 0
}

validate_module_list() {
    local list=$1
    local module=$2
    local version=${3:-}

    awk -v version="$version" -v module="$module" '
        BEGIN { gsub(/-/, "_", module) }
        {
            path = $0
            sub(/^\.\//, "", path)
            sub(/^\//, "", path)
            if (version != "") {
                prefix = "lib/modules/" version "/"
                if (index(path, prefix) != 1 &&
                    index(path, "usr/" prefix) != 1) next
            }
            count = split(path, parts, "/")
            name = parts[count]
            sub(/\.(xz|zst|gz)$/, "", name)
            if (name !~ /\.ko$/) next
            sub(/\.ko$/, "", name)
            gsub(/-/, "_", name)
            if (name == module) found = 1
        }
        END { exit !found }
    ' "$list"
}

validate_driver() {
    local root=$1
    local module=$2
    local option=$3
    local list=$4
    local version=${5:-}
    local config="$root/boot/config-$KERNEL"

    if [ ! -s "$config" ]; then
        echo "Missing kernel configuration: $config"
        return 1
    fi
    if grep -qx "$option=y" "$config"; then
        echo "$module is built into kernel $KERNEL ($option=y)"
        return 0
    fi
    if ! grep -qx "$option=m" "$config"; then
        echo "$module is not enabled in kernel $KERNEL ($option)"
        return 1
    fi
    if ! validate_module_list "$list" "$module" "$version"; then
        echo "$option=m, but $module.ko is missing from $list"
        return 1
    fi
    echo "$module is available as a loadable module ($option=m)"
    return 0
}

validate_initramfs() {
    local initramfs=$1
    local root="$VWORK/primary"
    local list="$VWORK/initramfs.list"

    validate_check "Initramfs can be read" lsinitramfs "$initramfs"
    if [ "$CHECK_PASSED" = no ]; then
        return
    fi
    cp "$VWORK/check.log" "$VWORK/initramfs.list"
    validate_check "Initramfs init program" \
        grep -Ex 'init|/?usr/bin/init|/?sbin/init' "$VWORK/initramfs.list"
    # Built-in drivers need no .ko in the initramfs. Loadable storage
    # drivers must be present there so the kernel can mount the rootfs.
    validate_check "Boot storage support mmc_block" validate_driver \
        "$root" mmc_block CONFIG_MMC_BLOCK "$list" "$KERNEL"
    validate_check "Boot storage support sdhci" validate_driver \
        "$root" sdhci CONFIG_MMC_SDHCI "$list" "$KERNEL"
    validate_check "Boot storage support sdhci_acpi" validate_driver \
        "$root" sdhci_acpi CONFIG_MMC_SDHCI_ACPI "$list" "$KERNEL"
    validate_check "Boot storage support sdhci_pci" validate_driver \
        "$root" sdhci_pci CONFIG_MMC_SDHCI_PCI "$list" "$KERNEL"
    validate_check "Boot filesystem support ext4" validate_driver \
        "$root" ext4 CONFIG_EXT4_FS "$list" "$KERNEL"
}

validate_boot() {
    local boot="$VWORK/boot"
    local config="$boot/boot/grub/grub.cfg"

    validate_check "x86_64 UEFI fallback loader" \
        validate_efi_loader "$boot/EFI/BOOT/BOOTX64.EFI"
    validate_check "Kernel present" test -s "$boot/boot/vmlinuz"
    validate_check "Initramfs present" test -s "$boot/boot/initrd.img"
    validate_check "Primary GRUB entry" \
        grep -Fq 'root=LABEL=primary' "$config"
    validate_check "Passive GRUB entry" \
        grep -Fq 'root=LABEL=passive' "$config"
    validate_check "Serial kernel debugging" \
        grep -Fq 'console=ttyS0,115200n8' "$config"
    validate_check "Autonomous boot timeout" grep -qx 'set timeout=3' "$config"

    KERNEL=$(cat "$boot/boot/kernel-version")
    case "$KERNEL" in
        ''|*[!A-Za-z0-9_.+-]*)
            validate_die "Invalid kernel version marker"
            ;;
    esac
    validate_initramfs "$boot/boot/initrd.img"
}

validate_parse_args() {
    while [ "$#" -gt 0 ]; do
        case "$1" in
            --repo)
                if [ "$#" -lt 2 ]; then
                    validate_die "Missing --repo value"
                fi
                VREPO=$2
                shift 2
                ;;
            -*)
                validate_die "Unknown argument: $1"
                ;;
            *)
                if [ -n "$VIMAGE" ]; then
                    validate_die "Only one image may be supplied"
                fi
                VIMAGE=$1
                shift
                ;;
        esac
    done
    if [ -z "$VIMAGE" ]; then
        validate_usage >&2
        exit 1
    fi
}

validate_inputs() {
    local tool
    local file
    local tools=(
        sfdisk losetup mount umount blockdev blkid fsck.fat e2fsck
        jq tar gzip readelf objdump diff cmp awk sed grep find
        sha256sum stat du df realpath lsinitramfs zstd
    )
    local inputs=(
        builder/scripts/build-system/pkg-utils.sh
        builder/boards/common.config
        builder/boards/com.config
        nodes/ukamaOS/distro/scripts/files/services
    )

    if [ ! -f "$VIMAGE" ]; then
        validate_die "Expected an image file: $VIMAGE"
    fi
    if [ "$EUID" -ne 0 ]; then
        validate_die "Run with sudo; read-only loop mounts require root"
    fi
    for tool in "${tools[@]}"; do
        if ! command -v "$tool" >/dev/null; then
            validate_die "Missing host command: $tool"
        fi
    done
    VREPO=$(realpath -e "$VREPO")
    VIMAGE=$(realpath -e "$VIMAGE")
    for file in "${inputs[@]}"; do
        if [ ! -s "$VREPO/$file" ]; then
            validate_die "Missing repository input: $file"
        fi
    done
}

validate_manifest() {
    local builder="$VREPO/builder"

    # Use the existing board selection and manifest logic unchanged.
    # shellcheck source=/dev/null
    source "$builder/scripts/build-system/pkg-utils.sh"
    get_enabled_apps "$builder/boards/common.config" \
        "$builder/boards/com.config"
    if [ "${#APPS[@]}" -eq 0 ]; then
        validate_die "No COM apps enabled"
    fi
    # The existing manifest helper calls log().
    log() {
        printf '%s: %s\n' "$1" "$2"
    }
    create_manifest_file "$VWORK/manifest.json" "${APPS[@]}"
    jq -e . "$VWORK/manifest.json" >/dev/null
}

validate_main() {
    set -euo pipefail
    export LC_ALL=C

    local script_dir
    local argument

    for argument in "$@"; do
        case "$argument" in
            -h|--help)
                validate_usage
                return
                ;;
        esac
    done

    script_dir=$(dirname -- "${BASH_SOURCE[0]}")
    script_dir=$(realpath -e "$script_dir")
    VREPO=$(realpath -e "$script_dir/../../..")
    VIMAGE=""
    VWORK=""
    VLOOP=""
    VMOUNTS=()
    APPS=()
    PASSED=0
    FAILED=0
    WARNINGS=0
    CHECK_PASSED=yes

    validate_parse_args "$@"
    validate_inputs
    VWORK=$(mktemp -d -t com-validate.XXXXXXXX)
    trap validate_cleanup EXIT
    trap 'exit 130' INT
    trap 'exit 143' TERM

    validate_manifest
    validate_layout
    validate_boot
    validate_os "$VWORK/primary" primary
    validate_os "$VWORK/passive" passive
    validate_apps

    printf '\nCOM validation: %s passed, %s failed, %s warnings.\n' \
        "$PASSED" "$FAILED" "$WARNINGS"
    if [ "$FAILED" -ne 0 ]; then
        return 1
    fi
    echo 'Static validation passed.'
    echo 'Hardware boot, networking and provisioning need separate testing.'
}

if [ "${BASH_SOURCE[0]}" = "$0" ]; then
    validate_main "$@"
fi
