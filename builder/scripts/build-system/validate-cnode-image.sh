#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

# Validate the cnode image without modifying it.

set -euo pipefail

export LC_ALL=C
export PATH="${PATH}:/usr/sbin:/sbin"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
STAGE="init"
RAW_IMG="cnode-image.img"
UKAMA_ROOT="${SCRIPT_DIR}/../../.."
MODEL="both"
COLOR="auto"
WORK=""
LOOPDISK=""
MOUNTS=()

PASS=0
FAIL=0
WARN=0
GREEN=""
RED=""
YELLOW=""
RESET=""

# Current application. load_app() resets these before checking each app.
APP_KEY=""
APP_NAME=""
APP_LABEL=""
APP_COMMAND=""
APP_TAG=""
APP_SOURCE=""
APP_PACKAGE=""
APP_DIRECTORY=""
APP_VERSION=""
APP_BINARY=""
APP_STATE=""

# Command line and console output.

usage() {
    echo "Usage: $0 [IMAGE] [cm4|cm5|both]"
    echo "       [--repo UKAMA_ROOT] [--color=auto|always|never]"
    echo "Run beside the build script, or supply --repo."
}

parse_args() {
    if [[ $# -gt 0 && $1 != --* ]]; then
        RAW_IMG=$1
        shift
    fi

    while [[ $# -gt 0 ]]; do
        case "$1" in
            cm4|cm5|both)
                MODEL=$1
                shift
                ;;
            --repo)
                if [[ $# -lt 2 ]]; then
                    usage
                    exit 2
                fi
                UKAMA_ROOT=$2
                shift 2
                ;;
            --color=auto|--color=always|--color=never)
                COLOR=${1#*=}
                shift
                ;;
            --help)
                usage
                exit 0
                ;;
            *)
                usage
                exit 2
                ;;
        esac
    done
}

setup_colors() {
    if [[ $COLOR == never ]]; then
        return
    fi

    if [[ $COLOR == auto ]]; then
        [[ -t 1 ]] || return 0
        [[ ${TERM:-} != dumb ]] || return 0
        [[ -z ${NO_COLOR:-} ]] || return 0
    fi

    GREEN=$'\033[1;32m'
    RED=$'\033[1;31m'
    YELLOW=$'\033[1;33m'
    RESET=$'\033[0m'
}

log() {
    local type=$1
    local message=$2
    local color=""

    case "$type" in
        PASS|SUCCESS)
            color=$GREEN
            PASS=$((PASS + 1))
            ;;
        FAIL|ERROR)
            color=$RED
            FAIL=$((FAIL + 1))
            ;;
        WARN)
            color=$YELLOW
            WARN=$((WARN + 1))
            ;;
    esac

    printf '%s%s%s  %s\n' "$color" "$type" "$RESET" "$message"
}

print_error_log() {
    cat "$1"
}

# Record a failed check and continue so the final report includes all failures.
check() {
    local message=$1
    local status
    shift

    if "$@" >"${WORK}/check.log" 2>&1; then
        log PASS "$message"
    else
        status=$?
        log FAIL "$message"
        if [[ -s ${WORK}/check.log ]]; then
            print_error_log "${WORK}/check.log"
        else
            printf '  Check returned %s:' "$status"
            printf ' %q' "$@"
            printf '\n'
        fi
    fi
}

die() {
    log ERROR "${STAGE}: $*"
    exit 2
}

# Setup and cleanup.

cleanup() {
    local status=$?
    local busy=0
    local i

    trap - EXIT

    for ((i=${#MOUNTS[@]}-1; i>=0; i--)); do
        if ! umount "${MOUNTS[i]}"; then
            log ERROR "Cannot unmount ${MOUNTS[i]}"
            busy=1
            status=2
        fi
    done

    if [[ $busy -eq 0 ]]; then
        if [[ -n $LOOPDISK ]]; then
            losetup -d "$LOOPDISK" || status=2
        fi
        if [[ -n $WORK ]]; then
            rm -rf -- "$WORK"
        fi
    fi

    exit "$status"
}

check_tools() {
    local cmd
    local commands=(
        sfdisk losetup blockdev blkid fsck.fat e2fsck mount umount
        jq tar gzip readelf sha256sum cmp diff find awk realpath stat
        readlink head cut tr grep
    )

    for cmd in "${commands[@]}"; do
        command -v "$cmd" >/dev/null || die "Missing tool: $cmd"
    done

    [[ $EUID -eq 0 ]] || die "Run with sudo"
}

check_prerequisites() {
    local common_config
    local cnode_config
    local app

    STAGE="check_prerequisites"
    RAW_IMG=$(realpath "$RAW_IMG")
    UKAMA_ROOT=$(realpath "$UKAMA_ROOT")

    [[ -f $RAW_IMG ]] || die "Image missing: $RAW_IMG"
    if [[ -n $(losetup -j "$RAW_IMG") ]]; then
        die "Image is already attached; finish that use first"
    fi

    PKG_UTILS="${UKAMA_ROOT}/builder/scripts/build-system/pkg-utils.sh"
    APP_PKG="${UKAMA_ROOT}/build/pkgs"
    LIB_PKG="${UKAMA_ROOT}/build/libs/vendor_libs.tgz"
    common_config="${UKAMA_ROOT}/builder/boards/common.config"
    cnode_config="${UKAMA_ROOT}/builder/boards/cnode.config"

    [[ -f $PKG_UTILS ]] || die "Missing pkg-utils.sh"
    [[ -d $APP_PKG ]] || die "Missing build/pkgs"
    [[ -s $LIB_PKG ]] || die "Missing vendor_libs.tgz"
    [[ -f $common_config ]] || die "Missing common.config"
    [[ -f $cnode_config ]] || die "Missing cnode.config"

    source "$PKG_UTILS"
    declare -F get_enabled_apps create_manifest_file \
        _manifest_toml_value >/dev/null || die "Missing package helpers"

    get_enabled_apps "$common_config" "$cnode_config"
    [[ ${#APPS[@]} -gt 0 ]] || die "No cnode apps enabled"

    for app in "${APPS[@]}"; do
        [[ -f ${UKAMA_ROOT}/builder/configs/${app}.toml ]] ||
            die "Missing config for $app"
    done

    WORK=$(mktemp -d /tmp/cnode-check.XXXXXXXX)
    BOOT_MOUNT="${WORK}/boot"
    PRIMARY_MOUNT="${WORK}/primary"
    PASSIVE_MOUNT="${WORK}/passive"

    create_manifest_file "${WORK}/expected-manifest.json" "${APPS[@]}"
    log INFO "Apps: ${APPS[*]}"
}

wait_for_partitions() {
    local partition
    local device
    local attempt

    for partition in 1 2 5 6 7 8; do
        device="${LOOPDISK}p${partition}"

        for ((attempt=0; attempt<50; attempt++)); do
            [[ -b $device ]] && break
            sleep 0.1
        done

        [[ -b $device ]] || die "Missing partition $partition"
    done
}

# Read-only partition inspection and mounting.

setup_loop_device() {
    STAGE="setup_loop_device"

    LOOPDISK=$(losetup --find --show --read-only --partscan "$RAW_IMG")
    [[ $(blockdev --getro "$LOOPDISK") == 1 ]] ||
        die "Loop device is not read-only"

    if command -v udevadm >/dev/null; then
        udevadm settle --timeout=10 || true
    fi

    wait_for_partitions
}

check_partition() {
    local partition=$1
    local label=$2
    local type=$3
    local gib=$4
    local device="${LOOPDISK}p${partition}"
    local actual
    local bytes

    actual=$(blkid -p -s LABEL -o value "$device" || true)
    check "Partition $partition label: $label (found: $actual)" \
        test "$actual" = "$label"

    actual=$(blkid -p -s TYPE -o value "$device" || true)
    check "Partition $partition filesystem: $type (found: $actual)" \
        test "$actual" = "$type"

    bytes=$(blockdev --getsize64 "$device")
    check "Partition $partition size: $gib GiB" \
        test "$bytes" -eq "$((gib * 1024 * 1024 * 1024))"

    case "$type" in
        vfat)
            check "$label filesystem integrity" fsck.fat -n "$device"
            ;;
        ext4)
            check "$label filesystem integrity" e2fsck -fn "$device"
            ;;
    esac
}

check_partitions() {
    local partition
    local label
    local type
    local gib

    STAGE="check_partitions"

    check "Image size: 28,800 MiB" \
        test "$(stat -c %s "$RAW_IMG")" -eq 30198988800
    check "Partition table consistency" sfdisk --verify "$LOOPDISK"

    sfdisk --json "$LOOPDISK" >"${WORK}/partitions.json"
    check "DOS partition layout" jq -e '
        .partitiontable |
        .label == "dos" and
        ([.partitions[].type] == ["c","83","5","83","83","83","82"])
    ' "${WORK}/partitions.json"

    while read -r partition label type gib; do
        check_partition "$partition" "$label" "$type" "$gib"
    done <<'LAYOUT'
1 boot vfat 1
2 recovery ext4 4
5 primary ext4 6
6 passive ext4 6
7 data ext4 10
8 swap swap 1
LAYOUT
}

mount_partition() {
    local partition=$1
    local target=$2
    local options=ro,nodev,nosuid,noexec

    if [[ $partition != 1 ]]; then
        options+=,noload
    fi

    mkdir -p "$target"
    mount -o "$options" "${LOOPDISK}p${partition}" "$target" ||
        die "Cannot mount partition $partition"
    MOUNTS+=("$target")
}

# Resolve image/package symlinks relative to their root, never the host root.
resolve_file() {
    local root=$1
    local path=${2#/}
    local current=""
    local depth=0
    local part
    local link

    while [[ -n $path ]]; do
        part=${path%%/*}
        if [[ $path == */* ]]; then
            path=${path#*/}
        else
            path=""
        fi

        case "$part" in
            ''|.)
                continue
                ;;
            ..)
                if [[ $current == */* ]]; then
                    current=${current%/*}
                else
                    current=""
                fi
                continue
                ;;
        esac

        if [[ -L $root/$current/$part ]]; then
            depth=$((depth + 1))
            [[ $depth -le 40 ]] || return 1

            link=$(readlink "$root/$current/$part")
            if [[ $link == /* ]]; then
                current=""
                link=${link#/}
            fi
            path="${link}${path:+/$path}"
        else
            current="${current:+$current/}${part}"
        fi
    done

    printf '%s/%s\n' "$root" "$current"
}

find_nonempty() {
    local found

    found=$(find "$1" -type f -name "$2" -size +0c -print -quit)
    [[ -n $found ]]
}

is_executable() {
    local file=$1
    local mode

    if [[ ! -f $file ]]; then
        echo "Missing regular file: ${file:-unresolved image path}"
        return 1
    fi
    if [[ ! -s $file ]]; then
        echo "Empty executable: $file"
        return 1
    fi

    # Inspect stored mode bits; -x also tests the host's noexec mount.
    mode=$(stat -Lc %a -- "$file") || return 1
    if (( (8#$mode & 0111) == 0 )); then
        echo "No execute bits: $file (mode $mode)"
        echo 'Expected at least one execute bit in the stored file mode.'
        return 1
    fi
}

has_kernel_modules() {
    local found

    found=$(find "$1/lib/modules" -type f -path '*-v8+/*' \
        -name '*.ko*' -size +0c -print -quit)
    [[ -n $found ]]
}

# Boot files and root filesystem configuration.

check_boot() {
    local file
    local boot_files=(
        kernel.img kernel8.img config.txt cmdline.txt
        overlays/disable-bt.dtbo overlays/disable-wifi.dtbo
    )

    STAGE="check_boot"

    for file in "${boot_files[@]}"; do
        check "Boot file: $file" test -s "${BOOT_MOUNT}/$file"
    done

    check "Selected kernel matches kernel8.img" \
        cmp -s "${BOOT_MOUNT}/kernel.img" "${BOOT_MOUNT}/kernel8.img"

    check "64-bit kernel selection" awk '
        /^arm_64bit=1$/ { arm64=1 }
        /^kernel=kernel.img$/ { kernel=1 }
        END { exit !(arm64 && kernel) }
    ' "${BOOT_MOUNT}/config.txt"

    check "Boot command line selects primary" awk '
        {
            for (i=1; i<=NF; i++) {
                if ($i=="root=/dev/mmcblk0p5") root=1
                if ($i=="rootfstype=ext4") ext4=1
                if ($i=="rootwait") wait=1
            }
        }
        END { exit !(NR==1 && root && ext4 && wait) }
    ' "${BOOT_MOUNT}/cmdline.txt"

    if [[ $MODEL != cm5 ]]; then
        for file in start4.elf fixup4.dat bcm2711-rpi-cm4.dtb; do
            check "CM4: $file" test -s "${BOOT_MOUNT}/$file"
        done
    fi

    if [[ $MODEL != cm4 ]]; then
        check "CM5 device tree" \
            find_nonempty "$BOOT_MOUNT" 'bcm2712-rpi-cm5*.dtb'
    fi
}

starter_enabled() {
    local root=$1
    local link
    local file

    for link in "$root"/etc/runlevels/*/starterd; do
        [[ -L $link ]] || continue
        file=$(resolve_file "$root" "${link#"$root"}") || return 1
        is_executable "$file" && return 0
    done

    echo 'No usable starterd link under /etc/runlevels/*/starterd.'
    return 1
}

compare_json() {
    jq -S . "$1" >"$WORK/actual.json" || return 1
    jq -S . "$2" >"$WORK/reference.json" || return 1
    cmp -s "$WORK/actual.json" "$WORK/reference.json"
}

check_rootfs_executables() {
    local root=$1
    local label=${root##*/}
    local path
    local file
    local executables=(
        /sbin/init /bin/sh /sbin/starter.d /sbin/openrc-run
        /etc/runlevels/boot/loopback
        /etc/runlevels/default/ntpd
        /etc/runlevels/boot/ukama-sysfs
    )

    for path in "${executables[@]}"; do
        file=$(resolve_file "$root" "$path") || file=""
        check "$label executable: $path" is_executable "$file"
    done

    check "$label starterd enabled" starter_enabled "$root"

    file=$(resolve_file "$root" /lib/ld-musl-aarch64.so.1) || file=""
    check "$label ARM64 musl loader" valid_elf "$file"
}

check_fstab() {
    local root=$1
    local label=$2

    check "$label root fstab" awk -v label="$label" '
        $1=="LABEL="label && $2=="/" && $3=="ext4" { ok=1 }
        END { exit !ok }
    ' "$root/etc/fstab"

    check "$label boot fstab" awk '
        $1=="LABEL=boot" && $2=="/boot/firmware" &&
            $3=="vfat" && $4=="ro" { ok=1 }
        END { exit !ok }
    ' "$root/etc/fstab"
}

check_app_configs() {
    local root=$1
    local label=${root##*/}
    local app

    for app in bootstrap meshd notifyd switchd; do
        check "$label $app config" \
            find_nonempty "$root/ukama/configs/$app" '*'
    done

    for app in metricsd aggregator; do
        check "$label $app config.toml" \
            test -s "$root/ukama/configs/$app/config.toml"
    done
}

has_mock_sysfs() {
    local root=$1
    local path="/ukama/mocksysfs/sys"
    local directory

    directory=$(resolve_file "$root" "$path") || return 1
    if [[ ! -d $directory ]]; then
        echo "Missing directory in image: $path"
        echo 'ukama-sysfs needs this directory to link /tmp/sys at boot.'
        echo 'Populate it with generated sysfs data during the image build.'
        return 1
    fi
}

check_rootfs() {
    local root=$1
    local label=$2

    STAGE="check_rootfs:$label"

    check_rootfs_executables "$root"
    check "$label manifest matches board apps" compare_json \
        "$root/ukama/manifest.json" "$WORK/expected-manifest.json"
    check "$label cnode schema" \
        jq -e 'type == "object" and length > 0' \
        "$root/ukama/mocksysfs/mfgdata/schema/cnode.json"
    check "$label mock sysfs" has_mock_sysfs "$root"
    check "$label DHCP config" test -s "$root/etc/dhcpcd.conf"
    check "$label bootstrap server" test -s "$root/ukama/bootstrap"

    check_fstab "$root" "$label"
    check "$label v8 kernel modules" has_kernel_modules "$root"
    check_app_configs "$root"
}

# Package integrity, executable format, and dependency helpers.

same_file() {
    local actual=$1
    local reference=$2
    local actual_hash
    local reference_hash

    [[ -s $actual && -s $reference ]] || return 1
    [[ $(stat -Lc %s "$actual") == "$(stat -Lc %s "$reference")" ]] ||
        return 1

    actual_hash=$(sha256sum "$actual" | cut -d' ' -f1)
    reference_hash=$(sha256sum "$reference" | cut -d' ' -f1)
    [[ $actual_hash == "$reference_hash" ]]
}

unpack_package() {
    local package=$1
    local destination=$2

    gzip -t "$package" || return 1
    tar -tzf "$package" >"${WORK}/members.txt" || return 1
    awk '
        /^\// || /(^|\/)\.\.(\/|$)/ { bad=1 }
        END { exit bad }
    ' "${WORK}/members.txt" || return 1

    mkdir -p "$destination"
    tar --no-same-owner --no-same-permissions \
        -xzf "$package" -C "$destination"
}

valid_elf() {
    readelf -aW "$1" >"${WORK}/elf.txt" \
        2>"${WORK}/elf-errors.txt" || return 1

    if [[ -s ${WORK}/elf-errors.txt ]]; then
        cat "${WORK}/elf-errors.txt"
        return 1
    fi

    awk '
        /Class:/ { elf64=($2=="ELF64") }
        /Machine:/ { arm64=($2=="AArch64") }
        END { exit !(elf64 && arm64) }
    ' "${WORK}/elf.txt"
}

lookup_file() {
    local package=$1
    local root=$2
    local path=$3
    local file

    file=$(resolve_file "$package" "$path") || return 1
    if [[ -s $file ]]; then
        printf '%s\n' "$file"
        return
    fi

    file=$(resolve_file "$root" "$path") || return 1
    [[ -s $file ]] && printf '%s\n' "$file"
}

dependency_search_paths() {
    local package=$1
    local root=$2
    local file=$3
    local env_paths=$4
    local defaults="/lib:/usr/local/lib:/usr/lib"
    local paths
    local origin

    if [[ -s $root/etc/ld-musl-aarch64.path ]]; then
        defaults=$(tr '\n' ':' <"$root/etc/ld-musl-aarch64.path")
    fi

    paths=$(awk -F'[][]' '
        /\((RUNPATH|RPATH)\)/ { printf "%s:", $2 }
    ' "$WORK/elf.txt")

    if [[ $file == "$package/"* ]]; then
        origin="/${file#"$package/"}"
    else
        origin="/${file#"$root/"}"
    fi

    origin=${origin%/*}
    paths="${env_paths}:${paths}${defaults}"
    paths=${paths//\$\{ORIGIN\}/$origin}
    paths=${paths//\$ORIGIN/$origin}

    printf '%s\n' "$paths"
}

find_library() {
    local package=$1
    local root=$2
    local library=$3
    local paths=$4
    local directory
    local found

    if [[ $library == /* ]]; then
        lookup_file "$package" "$root" "$library"
        return
    fi

    while IFS= read -r directory; do
        [[ $directory == /* && $directory != *'$'* ]] || continue

        if found=$(lookup_file "$package" "$root" "$directory/$library")
        then
            printf '%s\n' "$found"
            return 0
        fi
    done < <(tr ':' '\n' <<<"$paths")

    return 1
}

find_script_interpreter() {
    local package=$1
    local root=$2
    local file=$3
    local interpreter
    local found

    interpreter=$(head -n1 "$file" | awk '{sub(/^#!/, ""); print $1}')
    if [[ $interpreter != /* || $interpreter == */env ]]; then
        echo "Unsupported script interpreter: $file" >&2
        return 1
    fi

    if ! found=$(lookup_file "$package" "$root" "$interpreter"); then
        echo "Missing interpreter: $interpreter" >&2
        return 1
    fi

    printf '%s\n' "$found"
}

check_dependencies() {
    local package=$1
    local root=$2
    local binary=$3
    local env_paths=$4
    local file
    local library
    local found
    local paths
    local interpreter
    local queue=("$binary")
    local needed=()
    local -A seen=()

    while [[ ${#queue[@]} -gt 0 ]]; do
        file=${queue[0]}
        queue=("${queue[@]:1}")
        [[ -z ${seen[$file]:-} ]] || continue
        seen[$file]=1

        if [[ $(head -c2 "$file") == '#!' ]]; then
            found=$(find_script_interpreter "$package" "$root" "$file") ||
                return 1
            queue+=("$found")
            continue
        fi

        if ! valid_elf "$file"; then
            echo "Invalid ARM64 ELF: $file"
            return 1
        fi

        interpreter=$(awk '
            /Requesting program interpreter:/ {
                sub(/.*interpreter: /, "")
                sub(/].*/, "")
                print
            }
        ' "$WORK/elf.txt")

        if [[ -n $interpreter ]]; then
            if ! found=$(lookup_file "$package" "$root" "$interpreter")
            then
                echo "Missing loader: $interpreter"
                return 1
            fi
            is_executable "$found" || return 1
            queue+=("$found")
        fi

        mapfile -t needed < <(
            awk -F'[][]' '/\(NEEDED\)/ {print $2}' "$WORK/elf.txt"
        )
        paths=$(dependency_search_paths "$package" "$root" \
            "$file" "$env_paths")

        for library in "${needed[@]}"; do
            if ! found=$(find_library "$package" "$root" "$library" \
                "$paths")
            then
                echo "Missing library: $library (needed by $file)"
                return 1
            fi
            queue+=("$found")
        done
    done
}

# Application checks: load settings, inspect the package, then check rootfs.

load_app() {
    local root=$1
    local key=$2
    local config="${UKAMA_ROOT}/builder/configs/${key}.toml"
    local directory
    local binary

    APP_KEY=$key
    APP_NAME=$(_manifest_toml_value "$config" capp-exec name)
    APP_NAME=${APP_NAME:-$key}
    APP_LABEL="${root##*/} $APP_NAME"

    directory=$(_manifest_toml_value "$config" capp-exec path)
    binary=$(_manifest_toml_value "$config" capp-exec bin)
    directory=${directory:-/sbin}
    APP_COMMAND="${directory%/}/$binary"
    APP_TAG=$(_manifest_toml_value "$config" capp-exec version)
    APP_TAG=${APP_TAG:-latest}

    APP_SOURCE="${APP_PKG}/${key}_latest.tar.gz"
    if [[ ! -e $APP_SOURCE ]]; then
        APP_SOURCE="${APP_PKG}/${APP_NAME}_latest.tar.gz"
    fi
    APP_PACKAGE="$root/ukama/apps/pkgs/${APP_SOURCE##*/}"

    APP_DIRECTORY=""
    APP_VERSION="-"
    APP_BINARY=""
    APP_STATE="STAGED"

    printf '%s\n' "${APP_SOURCE##*/}" >>"$WORK/expected-packages.txt"
}

manifest_entry_matches() {
    local manifest=$1
    local entries

    entries=$(jq -c --arg name "$APP_NAME" '
        [.spaces[].apps[] | select(.name==$name)]
    ' "$manifest") || return 1

    if jq -e \
        --arg cmd "${APP_COMMAND#/}" \
        --arg tag "$APP_TAG" '
        length==1 and .[0].cmd==$cmd and .[0].tag==$tag
    ' <<<"$entries" >/dev/null; then
        return 0
    fi

    echo "Expected one manifest entry for $APP_NAME:"
    echo "  cmd=${APP_COMMAND#/} tag=$APP_TAG"
    jq -r '
        "Found \(length) matching entries:",
        (.[] | "  cmd=\(.cmd) tag=\(.tag)")
    ' <<<"$entries"
    return 1
}

check_app_manifest() {
    local root=$1

    # starterd starts from OpenRC, so it has no manifest entry.
    [[ $APP_KEY != starterd ]] || return 0

    check "$APP_LABEL: manifest entry" \
        manifest_entry_matches "$root/ukama/manifest.json"
}

open_app_package() {
    local extracted
    local roots=()

    if [[ ! -s $APP_PACKAGE ]]; then
        log FAIL "$APP_LABEL: missing package $APP_PACKAGE"
        APP_STATE="MISSING"
        return 1
    fi

    check "$APP_LABEL: package bytes/SHA-256 match build" \
        same_file "$APP_PACKAGE" "$APP_SOURCE"

    extracted=$(mktemp -d "${WORK}/app.XXXXXXXX") ||
        die "Cannot create package inspection directory"

    if ! unpack_package "$APP_PACKAGE" "$extracted" \
        >"$WORK/archive.log" 2>&1
    then
        log FAIL "$APP_LABEL: invalid archive"
        print_error_log "$WORK/archive.log"
        APP_STATE="INVALID"
        return 1
    fi

    mapfile -t roots < <(
        find "$extracted" -mindepth 1 -maxdepth 1 -printf '%p\n'
    )
    if [[ ${#roots[@]} -eq 1 ]]; then
        APP_DIRECTORY=${roots[0]}
    fi

    if [[ -z $APP_DIRECTORY ||
          ! -d $APP_DIRECTORY/sbin ||
          ! -s $APP_DIRECTORY/VERSION ]]; then
        log FAIL "$APP_LABEL: expected one app root with VERSION and sbin"
        APP_STATE="INVALID"
        return 1
    fi

    APP_VERSION=$(cat "$APP_DIRECTORY/VERSION")
    APP_BINARY=$(resolve_file "$APP_DIRECTORY" "$APP_COMMAND") ||
        APP_BINARY=""
}

check_app_version() {
    check "$APP_LABEL: package name/VERSION" \
        test "${APP_DIRECTORY##*/}" = "${APP_NAME}_${APP_VERSION}"

    if [[ $APP_TAG != latest ]]; then
        check "$APP_LABEL: manifest version" \
            test "$APP_VERSION" = "$APP_TAG"
    fi
}

check_app_binary() {
    local file
    local relative_path

    check "$APP_LABEL: executable present and nonempty" \
        is_executable "$APP_BINARY"

    while IFS= read -r -d '' file; do
        if [[ $(head -c4 "$file") == $'\177ELF' ]]; then
            relative_path=${file#"$APP_DIRECTORY/"}
            check "$APP_LABEL: ARM64 $relative_path" valid_elf "$file"
        fi
    done < <(find "$APP_DIRECTORY" -type f -print0)
}

check_app_dependencies() {
    local root=$1
    local env_paths

    # A missing binary is already reported by check_app_binary().
    [[ -s $APP_BINARY ]] || return 0

    env_paths=$(jq -r --arg name "$APP_NAME" '
        [.spaces[].apps[] |
            select(.name==$name) |
            .env.LD_LIBRARY_PATH // ""] |
        first // ""
    ' "$root/ukama/manifest.json") || env_paths=""

    check "$APP_LABEL: loader and shared libraries" check_dependencies \
        "$APP_DIRECTORY" "$root" "$APP_BINARY" "$env_paths"
}

check_installed_app() {
    local root=$1
    local path
    local file
    local candidates=(
        "$APP_COMMAND"
        "/ukama/apps/rootfs/${APP_DIRECTORY##*/}${APP_COMMAND}"
    )

    # The builder stages most apps; compare installed copies when present.
    for path in "${candidates[@]}"; do
        file=$(resolve_file "$root" "$path") || continue
        [[ -e $file ]] || continue

        APP_STATE="INSTALLED"
        check "$APP_LABEL: installed binary bytes/SHA-256" \
            same_file "$file" "$APP_BINARY"
        check "$APP_LABEL: installed execute permission" \
            is_executable "$file"
    done
}

check_installed_starter() {
    local root=$1
    local label=${root##*/}
    local file

    [[ $APP_KEY == starterd ]] || return 0

    file=$(resolve_file "$root" /sbin/starter.d) || file=""
    check "$label starter: installed bytes/SHA-256" \
        same_file "$file" "$APP_BINARY"
    check "$label starter: installed dependencies" check_dependencies \
        "$WORK/empty" "$root" "$file" ""
}

file_size_mb() {
    local file=$1
    local bytes

    if [[ ! -s $file ]]; then
        printf '%s' '-'
        return
    fi

    bytes=$(stat -Lc %s "$file")
    awk -v bytes="$bytes" 'BEGIN { printf "%.3f", bytes / 1000000 }'
}

show_app() {
    local failures_before=$1
    local package_mb
    local binary_mb
    local color=$GREEN
    local result="PASS"

    package_mb=$(file_size_mb "$APP_PACKAGE")
    binary_mb=$(file_size_mb "$APP_BINARY")

    if [[ $FAIL -gt $failures_before ]]; then
        color=$RED
        result="FAIL"
    fi

    printf '%-20s %11s %11s %-25s %-10s %s%s%s\n' \
        "$APP_NAME" "$package_mb" "$binary_mb" "$APP_VERSION" \
        "$APP_STATE" "$color" "$result" "$RESET" \
        >>"${WORK}/apps.txt"
}

check_app() {
    local root=$1
    local key=$2
    local failures_before=$FAIL

    load_app "$root" "$key"
    check_app_manifest "$root"

    if open_app_package; then
        check_app_version
        check_app_binary
        check_app_dependencies "$root"
        check_installed_app "$root"
        check_installed_starter "$root"
    fi

    show_app "$failures_before"
}

check_extra_packages() {
    local root=$1
    local file

    for file in "$root"/ukama/apps/pkgs/*.tar.gz; do
        [[ -e $file ]] || continue

        if ! grep -Fxq "${file##*/}" "$WORK/expected-packages.txt"; then
            log WARN "${root##*/}: extra package ${file##*/}"
        fi
    done
}

check_apps() {
    local root=$1
    local app

    STAGE="check_apps:${root##*/}"
    : >"${WORK}/apps.txt"
    : >"${WORK}/expected-packages.txt"

    for app in "${APPS[@]}"; do
        check_app "$root" "$app"
    done
    check_extra_packages "$root"

    printf '\n%-20s %11s %11s %-25s %-10s %s\n' \
        APP 'PACKAGE MB' 'BINARY MB' VERSION STATE RESULT
    cat "${WORK}/apps.txt"
    echo 'MB = 1,000,000 bytes. STAGED = in package.'
    echo 'INSTALLED = executable also checked in rootfs; see RESULT.'
}

# Shared libraries and primary/passive consistency.

check_libraries() {
    local root=$1
    local label=${root##*/}
    local file
    local target

    STAGE="check_libraries:$label"

    for file in "$WORK"/vendor/*; do
        target=$(resolve_file "$root" "/lib/${file##*/}") || target=""
        check "$label: vendor library ${file##*/} bytes/SHA-256" \
            same_file "$file" "$target"

        if [[ -s $target && $(head -c4 "$target") == $'\177ELF' ]]; then
            check "$label: ARM64 library ${file##*/}" valid_elf "$target"
        fi
    done
}

package_names() {
    local root=$1

    cd "$root/ukama/apps/pkgs" || return 1
    printf '%s\n' *.tar.gz
}

check_copies() {
    local file

    for file in manifest.json bootstrap; do
        check "Primary/passive $file" same_file \
            "$PRIMARY_MOUNT/ukama/$file" "$PASSIVE_MOUNT/ukama/$file"
    done

    check "Primary/passive package lists" diff -q \
        <(package_names "$PRIMARY_MOUNT") \
        <(package_names "$PASSIVE_MOUNT")

    for file in "$PRIMARY_MOUNT"/ukama/apps/pkgs/*.tar.gz; do
        check "Primary/passive ${file##*/}" same_file \
            "$file" "$PASSIVE_MOUNT/ukama/apps/pkgs/${file##*/}"
    done
}

print_summary() {
    printf '\n%sPASS%s=%s  %sFAIL%s=%s  %sWARN%s=%s\n' \
        "$GREEN" "$RESET" "$PASS" \
        "$RED" "$RESET" "$FAIL" \
        "$YELLOW" "$RESET" "$WARN"

    echo 'Static checks only: boot the cnode to confirm application health'
    echo 'and runtime library compatibility.'

    if [[ $FAIL -eq 0 ]]; then
        printf 'RESULT: %sPASS%s\n' "$GREEN" "$RESET"
    else
        printf 'RESULT: %sFAIL%s\n' "$RED" "$RESET"
    fi
}

# Main
parse_args "$@"
setup_colors

trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

check_tools
check_prerequisites
setup_loop_device
check_partitions

mount_partition 1 "$BOOT_MOUNT"
mount_partition 5 "$PRIMARY_MOUNT"
mount_partition 6 "$PASSIVE_MOUNT"

check_boot
mkdir -p "$WORK/empty"
unpack_package "$LIB_PKG" "$WORK/vendor" || die "Invalid vendor_libs.tgz"

for target in "$PRIMARY_MOUNT" "$PASSIVE_MOUNT"; do
    check_rootfs "$target" "${target##*/}"
    check_apps "$target"
    check_libraries "$target"
done

check_copies
log INFO 'Recovery is empty; data and swap require runtime setup.'
log INFO 'Boot selects primary.'
print_summary

[[ $FAIL -eq 0 ]]
