#!/usr/bin/env bash
# SPDX-License-Identifier: MPL-2.0
# Copyright (c) 2026-present, Ukama Inc.
# Build the TRX release from the SDK and LTE source trees.

set -euo pipefail
umask 022
export TZ=UTC

host_env=${XDG_CACHE_HOME:-$HOME/.cache}/ukama/trx-host
export PATH="$host_env/bin:$PATH:/usr/sbin:/sbin"

script_path=$(realpath -- "$0")
scripts=$(dirname -- "$script_path")
source "$scripts/validate-trx-image.sh"

usage() {
    echo "Usage: $0 --sdk DIR --stack DIR --out DIR [--mode FDD|TDD --band 3|40|41] [--dl-earfcn N] [--profile FILE]"
    echo "       $0 --setup-host    # Ubuntu 22.04 build host"
    echo "Builds in the supplied SDK and stack working directories."
    echo "Default: FDD/B3, 20 MHz, using the channel in the stack's 64UE profile."
    echo "TDD uses the stack's 20 MHz INTRA profile. Supply --dl-earfcn when changing its band."
    echo "--profile FILE.profile reads mode, band, dl_earfcn, bandwidth_mhz and ue_target."
    echo "Radio flags override .profile defaults. Use a vendor .cfg to embed that file as written."
}

# Read KEY=VALUE as data. Never execute a profile as a shell script.
read_defaults() {
    local file=$1
    local line key
    local -A values=()

    while IFS= read -r line || [[ -n $line ]]; do
        line=${line%$'\r'}
        case "$line" in
            ''|'#'*) continue ;;
            mode=*|band=*|dl_earfcn=*|bandwidth_mhz=*|ue_target=*) ;;
            *) die "Unknown profile entry in $file: $line" ;;
        esac
        key=${line%%=*}
        if [[ -n ${values[$key]+set} ]]; then
            die "Duplicate $key in $file"
        fi
        values[$key]=${line#*=}
    done < "$file"

    for key in mode band dl_earfcn bandwidth_mhz ue_target; do
        if [[ -z ${values[$key]:-} ]]; then
            die "Missing $key in $file"
        fi
    done
    if [[ ${values[bandwidth_mhz]} != 20 || ${values[ue_target]} != 64 ]]; then
        die "Supported releases use bandwidth_mhz=20 and ue_target=64"
    fi
    requested_mode=${requested_mode:-${values[mode]^^}}
    requested_band=${requested_band:-${values[band]}}
    requested_dl=${requested_dl:-${values[dl_earfcn]}}
}

# Prepare the startup.cfg that will be validated and embedded in the image.
prepare_profile() {
    local base base_band old_dl dl ul duplex

    if [[ -n $startup_cfg ]]; then
        install -m 0644 -- "$startup_cfg" "$app/startup.cfg"
        return
    fi

    if [[ -z $requested_mode ]]; then
        requested_mode=FDD
        requested_band=3
    fi
    case "$requested_mode/$requested_band" in
        FDD/3)
            base="$stack/stack/fdd/configs/OTA/64UE/startup10.102.81.75_20mhz.cfg"
            duplex=1
            ;;
        TDD/40|TDD/41)
            base="$stack/stack/tdd/configs/HO/Manual/INTRA/20Mhz/startup192.168.1.25_20MHz_UD2_SSF_8_PF4.cfg"
            duplex=0
            ;;
        *)
            die "Supported radio selections are FDD/3, TDD/40 and TDD/41"
            ;;
    esac

    if [[ ! -r $base ]]; then
        die "Missing base profile $base. Use --profile FILE to select another."
    fi
    if [[ $(setting "$base" 'configure cell etc eduplexmode 0') != "$duplex" ]]; then
        die "The duplex mode in $base does not match $requested_mode"
    fi
    base_band=$(setting "$base" 'configure cell u8freqbandindicator 0')
    old_dl=$(setting "$base" 'configure cell u16dlcarrierfreq 0')
    if [[ $requested_band != "$base_band" && -z $requested_dl ]]; then
        die "Specify --dl-earfcn for Band $requested_band. The base profile uses Band $base_band."
    fi
    dl=${requested_dl:-$old_dl}
    if [[ $requested_band == "$base_band" && $dl == "$old_dl" ]]; then
        install -m 0644 -- "$base" "$app/startup.cfg"
        return
    fi
    if [[ $requested_mode == FDD ]]; then
        ul=$((dl + 18000))
    else
        ul=$dl
    fi

    # Keep intra-frequency neighbors on the serving channel. Preserve other carriers.
    awk -v band="$requested_band" -v old_dl="$old_dl" -v dl="$dl" -v ul="$ul" '
        {
            original = $0
            sub(/\r$/, "")
            if (NF == 0) {
                print original
                next
            }
            value = $NF
            $NF = ""
            if ($0 == "configure cell u8freqbandindicator 0 ") {
                print $0 band
            } else if ($0 == "configure cell u16dlcarrierfreq 0 ") {
                print $0 dl
            } else if ($0 == "configure cell freq u16ulcarrierfreq 0 ") {
                print $0 ul
            } else if ($0 ~ /^configure cell eutranfreqrelation u16arfcnvalueeutrandl 0 [0-9]+ $/ && value == old_dl) {
                print $0 dl
            } else {
                print original
            }
        }
    ' "$base" > "$app/startup.cfg"
}

# Prepare a new Jammy host once. No SDK checkout or build products are needed.
setup_host() {
    local packages=(
        build-essential git autoconf automake libtool pkg-config flex bison
        gawk bc gettext autopoint texinfo gperf help2man patch patchutils
        zlib1g-dev libncurses-dev libssl-dev libelf-dev libsctp-dev libisl-dev
        libtool-bin libltdl-dev libgnutls28-dev libgcrypt20-dev libsystemd-dev
        libcurl4-gnutls-dev libmicrohttpd-dev libnftnl-dev libmnl-dev
        libc6-i386 libc6-dev-i386 cpio fakeroot faketime mtd-utils u-boot-tools libarchive-zip-perl
        rsync lzop rpm2cpio bsdextrautils sharutils subversion cvs wget curl
        cmake tcl patchelf gnat docbook-utils python3-venv
    )

    sudo apt-get update
    sudo apt-get install --no-remove "${packages[@]}"
    /usr/bin/python3 -m venv "$host_env"
    "$host_env/bin/pip" install 'jefferson==0.4.8'
}

# Older SDKs have a stray ']' in the compiler check. Correct the recipe for
# this invocation while leaving the caller's Makefile and local edits intact.
sdk_make() {
    local tree=$1
    local target=$2
    local makefile=$tree/Makefile

    if grep -Fq 'gcc -v ]' "$makefile"; then
        sed 's/gcc -v ]/gcc -v/' "$makefile" > "$temp/sdk.mk"
        makefile=$temp/sdk.mk
    fi
    make -C "$tree" -f "$makefile" SHELL=/bin/bash J=1 "$target"
}

advertises_trx() {
    local repo=$1
    local targets

    if ! targets=$(make -s --no-print-directory -C "$repo" available-uboot-target 2>/dev/null); then
        return 1
    fi
    tr '[:space:]' '\n' <<< "$targets" | grep -Fx 'u-boot-octeon_trx' > /dev/null
}

find_trx_board() {
    local tree table defconfig config_dir

    for tree in "$sdk" "$stack/distro"; do
        if [[ -d $tree/bootloader/u-boot ]] && advertises_trx "$tree"; then
            echo "$tree/bootloader/u-boot"
            return
        fi
    done

    # Older U-Boot versions list the board in boards.cfg.
    while IFS= read -r -d '' table; do
        if grep -Eq '^[[:space:]]*octeon_trx[[:space:]]' "$table"; then
            dirname -- "$table"
            return
        fi
    done < <(find "$sdk" "$stack" -type f -name boards.cfg -print0)

    defconfig=$(find "$sdk" "$stack" -type f -path '*/configs/octeon_trx_defconfig' -print -quit)
    if [[ -z $defconfig ]]; then
        return 1
    fi
    config_dir=$(dirname -- "$defconfig")
    dirname -- "$config_dir"
}

# Source the vendor environment in the SDK directory that owns the toolchain.
setup_octeon() {
    local setup_status

    cd -- "$1"
    # The vendor probes for csh with setenv and inspects its failure under Bash.
    set +eu
    # shellcheck source=/dev/null
    source ./env-setup OCTEON_CNF71XX_PASS1_1
    setup_status=$?
    set -eu
    if (( setup_status != 0 )); then
        die "$1/env-setup failed (status $setup_status)"
    fi

    export PATH="$PWD/host/bin:$PATH"
    if ! command -v mips64-octeon-linux-gnu-gcc >/dev/null; then
        die "Missing OCTEON compiler after env-setup in $PWD"
    fi
}

write_env_size() {
    local boot_sdk=$1
    local env_size

    if [[ -f .config ]]; then
        env_size=$(sed -n 's/^CONFIG_ENV_SIZE=//p' .config)
        if [[ $env_size =~ ^(0x[[:xdigit:]]+|[1-9][0-9]*)$ ]]; then
            echo "$((env_size))" > "$inputs/env.size"
            return
        fi
    fi

    # Legacy U-Boot defines CONFIG_ENV_SIZE as a C expression. The cross
    # compiler resolves it using the board headers selected by the build.
    cat > "$temp/env-size.c" <<'ENV_SIZE'
#include <config.h>
const unsigned int size __attribute__((section(".trx_env_size"))) = CONFIG_ENV_SIZE;
ENV_SIZE
    mips64-octeon-linux-gnu-gcc -EB -D__KERNEL__ -Iinclude -Iarch/mips/include \
        -I"$boot_sdk/executive" -c "$temp/env-size.c" -o "$temp/env-size.o"
    mips64-octeon-linux-gnu-objcopy -O binary -j .trx_env_size "$temp/env-size.o" "$temp/env-size.bin"
    if [[ $(file_size "$temp/env-size.bin") != 4 ]]; then
        die "Cannot read CONFIG_ENV_SIZE for octeon_trx"
    fi
    od -An -N4 --endian=big -tu4 "$temp/env-size.bin" | tr -d ' ' > "$inputs/env.size"
}

build_bootloader() (
    local boot_sdk=$sdk
    local cross=mips64-octeon-linux-gnu-
    local config_target=octeon_trx_config
    local binary header

    if [[ $board == "$stack/distro/"* ]]; then
        boot_sdk=$stack/distro
    fi
    setup_octeon "$boot_sdk"
    cd "$board"

    if [[ -f configs/octeon_trx_defconfig ]]; then
        config_target=octeon_trx_defconfig
    elif [[ -f include/image.h && -f include/libfdt.h && -d lib/libfdt ]]; then
        # Make the host tools find the vendor FDT headers before system headers.
        mkdir "$temp/fdt-include"
        for header in fdt.h libfdt.h libfdt_env.h; do
            if [[ -f include/$header ]]; then
                ln -s "$board/include/$header" "$temp/fdt-include/$header"
            fi
        done
        export CPATH="$temp/fdt-include${CPATH:+:$CPATH}"
    fi

    if [[ $board == "$boot_sdk/bootloader/u-boot" ]] && advertises_trx "$boot_sdk"; then
        mkdir -p "$boot_sdk/target/bin"
        binary="$boot_sdk/output/u-boot-octeon_trx.bin"
        sdk_make "$boot_sdk" u-boot-octeon_trx
    else
        make ARCH=mips CROSS_COMPILE="$cross" distclean
        make ARCH=mips CROSS_COMPILE="$cross" "$config_target"
        # Both names are generated by this target. Select only a newly built image.
        rm -f -- u-boot.bin u-boot-octeon_trx.bin
        make -j1 ARCH=mips CROSS_COMPILE="$cross" all
        binary=u-boot-octeon_trx.bin
        if [[ ! -s $binary ]]; then
            binary=u-boot.bin
        fi
    fi

    install -m 0644 "$binary" "$inputs/u-boot-octeon_trx.bin"
    check_uboot "$inputs/u-boot-octeon_trx.bin"

    write_env_size "$boot_sdk"
)

build_sdk() (
    local version=develop-v0.0.1
    local compiler_version tool

    setup_octeon "$sdk"

    if [[ ! -f linux/embedded_rootfs/.config ]]; then
        cp linux/embedded_rootfs/default.config linux/embedded_rootfs/.config
    fi
    if [[ ! -e linux/kernel/linux/.config && -f linux/kernel/kernel.config ]]; then
        cp linux/kernel/kernel.config linux/kernel/linux/.config
    fi
    sdk_make "$sdk" basic_check

    # A vendor source distribution can include host objects built on another machine.
    sdk_make "$sdk" host-tools-clean
    sdk_make "$sdk" host-tools
    if [[ ! -x host/remote-utils/oct-remote-boot ]]; then
        die "SDK did not build oct-remote-boot"
    fi

    mkdir "$inputs/host-tools"
    for tool in host/remote-utils/oct-remote-*; do
        if [[ -f $tool && -x $tool && $tool != *.* ]]; then
            install -m 0755 "$tool" "$inputs/host-tools/"
        fi
    done

    sdk_make "$sdk" kernel-deb
    sdk_make "$sdk" rootfs

    # Normalize CPIO metadata while preserving device nodes, owners, modes and links.
    mkdir "$temp/rootfs"
    fakeroot -- bash -s -- "$temp" "$sdk/linux/embedded_rootfs/rootfs.cpio.gz" <<'RAMDISK'
set -euo pipefail
work=$1
archive=$2

cd "$work/rootfs"
gzip -dc "$archive" | cpio -idm --quiet --no-absolute-filenames
find . -exec touch -h -d "@$SOURCE_DATE_EPOCH" {} +

find . -print0 | sort -z > "$work/rootfs.list"
cpio --quiet --null -o -H newc --reproducible < "$work/rootfs.list" |
    gzip -9n > "$work/rd.gz"
RAMDISK
    gzip -9nc output/lsm_os > "$temp/os.gz"

    # gzipinfo writes the clock into the SDK tag. Use the same clock for every rebuild.
    if [[ -s .rev ]]; then
        version=$(<.rev)
    fi
    faketime -f "$KBUILD_BUILD_TIMESTAMP" gzipinfo -t LSM_OS -a "$version" "$temp/os.gz" > output/lsm_os.gz
    faketime -f "$KBUILD_BUILD_TIMESTAMP" gzipinfo -t LSM_RD -a "$version" "$temp/rd.gz" > output/lsm_rd.gz

    rm -f output/fw.bin
    TAR_OPTIONS="--mtime=@$SOURCE_DATE_EPOCH --sort=name --owner=0 --group=0 --numeric-owner" \
        sdk_make "$sdk" fw
    cp output/lsm_os.gz output/lsm_rd.gz output/fw.bin "$inputs/"

    printf 'sdk_version=%s\n' "$version" >> "$inputs/build.info"
    compiler_version=$(mips64-octeon-linux-gnu-gcc -dumpversion)
    echo "sdk_gcc=$compiler_version" >> "$inputs/build.info"
)

build_stack() (
    local tools=$stack/tools
    local source="$stack/stack/${mode,,}/source"
    local bundle name

    cd "$tools/crosstool-ng"
    ./bootstrap
    ./configure --prefix="$tools/build"
    make -j1
    make install
    export PATH="$tools/build/bin:$PATH"

    setup_octeon "$stack/distro"
    cd "$stack"

    if [[ $mode == FDD ]]; then
        bundle="$source/LTE_Stack_FDD_Bin_startup_NMMdisable_IPV6disable_158476"
    else
        bundle="$source/LTE_Stack_TDD_Bin_startup_NMMdisable_IPV6disable_158239"
    fi

    # The vendor target owns its build outputs in the selected stack tree.
    export TAR_OPTIONS="--sort=name --mtime=@$SOURCE_DATE_EPOCH --owner=0 --group=0 --numeric-owner"
    sudo -E env PATH="$PATH" make SHELL=/bin/bash stack TYPE="$mode"

    for name in lsmD dsp.tgz LFMSOFT_OCT_D.tgz; do
        if [[ ! -s $bundle/$name ]]; then
            die "Stack build did not produce $bundle/$name"
        fi
        install -m 0644 "$bundle/$name" "$app/$name"
    done
)

cleanup() {
    cd /
    if ! rm -rf -- "$temp" 2>/dev/null; then
        # The staged rootfs can contain directories with restrictive permissions.
        sudo rm -rf -- "$temp"
    fi
}

if [[ $# == 1 && $1 == --setup-host ]]; then
    setup_host
    exit 0
fi

sdk=
stack=
profile=
startup_cfg=
requested_mode=
requested_band=
requested_dl=
out=

while (( $# > 0 )); do
    case "$1" in
        --help|-h)
            usage
            exit 0
            ;;
        --sdk)
            sdk=${2:-}
            ;;
        --stack)
            stack=${2:-}
            ;;
        --profile)
            profile=${2:-}
            ;;
        --mode)
            requested_mode=${2:-}
            requested_mode=${requested_mode^^}
            ;;
        --band)
            requested_band=${2:-}
            ;;
        --dl-earfcn)
            requested_dl=${2:-}
            ;;
        --out)
            out=${2:-}
            ;;
        *)
            die "Unknown argument $1. Use --help."
            ;;
    esac

    if [[ -z ${2:-} || $2 == --* ]]; then
        die "Missing value for $1"
    fi
    shift 2
done

if [[ -z $sdk || -z $stack || -z $out ]]; then
    usage >&2
    exit 2
fi
sdk=$(realpath -e -- "$sdk")
stack=$(realpath -e -- "$stack")
out=$(realpath -m -- "$out")

if [[ ! -f $sdk/Makefile || ! -f $sdk/env-setup ]]; then
    die "--sdk must name the SDK root containing Makefile and env-setup"
fi
if [[ ! -f $stack/Makefile ]]; then
    die "--stack must name the LTE source repository"
fi
if [[ -n $requested_mode && -z $requested_band || -z $requested_mode && -n $requested_band ]]; then
    die "Specify --mode and --band together"
fi
if [[ -n $profile ]]; then
    if [[ ! -f $profile || ! -r $profile ]]; then
        die "Cannot read profile $profile"
    fi
    if [[ $profile == *.profile ]]; then
        read_defaults "$profile"
    else
        startup_cfg=$(realpath -e -- "$profile")
    fi
fi
if [[ -n $requested_dl && ! $requested_dl =~ ^[1-9][0-9]{0,4}$ ]]; then
    die "--dl-earfcn must be a positive decimal LTE channel number (up to five digits)"
fi

# 1. Build in the supplied working trees. Serialize writes to each tree.
exec 7< "$sdk/Makefile"
if ! flock -n -x 7; then
    die "Another build is using $sdk"
fi
if [[ ! $sdk/Makefile -ef $stack/Makefile ]]; then
    exec 8< "$stack/Makefile"
    if ! flock -n -x 8; then
        die "Another build is using $stack"
    fi
fi
printf 'SDK:   %s\n' "$sdk"
printf 'Stack: %s\n' "$stack"

# Stage packaging inputs separately so cleanup only removes files we created.
work_root=${XDG_CACHE_HOME:-$HOME/.cache}/ukama/trx-build
mkdir -p -- "$work_root"
work_root=$(realpath -- "$work_root")
temp=$(mktemp -d -- "$work_root/work.XXXXXXXX")
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

inputs=$temp/inputs
app=$temp/app
mkdir "$inputs" "$app"

# 2. Establish that a TRX loader can be built before spending time on Linux or LTE.
if ! board=$(find_trx_board); then
    die "Cannot find octeon_trx U-Boot source in either repository. TRX requires board 59."
fi
for tool in git make sudo mkfs.jffs2 jffs2dump jefferson faketime fakeroot cpio crc32 mkenvimage; do
    need "$tool"
done

# 3. Validate the effective profile before compiling and use its mode to build the stack.
prepare_profile
radio_profile "$app/startup.cfg" > "$temp/radio.conf"
mode=$(sed -n 's/^mode=//p' "$temp/radio.conf")
band=$(sed -n 's/^band=//p' "$temp/radio.conf")
dl=$(sed -n 's/^dl_earfcn=//p' "$temp/radio.conf")
if [[ -n $requested_mode && $requested_mode/$requested_band != "$mode/$band" ]]; then
    die "The profile selects $mode/B$band. Requested $requested_mode/B$requested_band."
fi
if [[ -n $requested_dl && $requested_dl != "$dl" ]]; then
    die "The profile selects downlink EARFCN $dl. Requested $requested_dl."
fi
echo "Radio: $mode, Band $band, 20 MHz, downlink EARFCN $dl"

# 4. Pin the clock and record the base revisions of the supplied working trees.
if ! stack_rev=$(git -C "$stack" rev-parse HEAD 2>/dev/null); then
    stack_rev="source-tree"
fi
if ! sdk_rev=$(git -C "$sdk" rev-parse HEAD 2>/dev/null); then
    sdk_rev="source-tree"
fi

if [[ -z ${SOURCE_DATE_EPOCH:-} ]]; then
    if [[ $stack_rev == source-tree ]]; then
        SOURCE_DATE_EPOCH=0
    else
        SOURCE_DATE_EPOCH=$(git -C "$stack" log -1 --format=%ct)
    fi
fi
if [[ ! $SOURCE_DATE_EPOCH =~ ^[0-9]+$ ]]; then
    die "SOURCE_DATE_EPOCH must be a Unix timestamp"
fi
export SOURCE_DATE_EPOCH

KBUILD_BUILD_TIMESTAMP=$(date -u -d "@$SOURCE_DATE_EPOCH" '+%Y-%m-%d %H:%M:%S')
export KBUILD_BUILD_TIMESTAMP
export KBUILD_BUILD_USER=ukama
export KBUILD_BUILD_HOST=builder
export KBUILD_BUILD_VERSION=1

cat > "$inputs/build.info" <<BUILD_INFO
sdk_commit=$sdk_rev
stack_commit=$stack_rev
source_date_epoch=$SOURCE_DATE_EPOCH
stack_mode=$mode
ue_target=64
BUILD_INFO

# 5. Compile TRX U-Boot and determine the environment format from that same target.
build_bootloader

# 6. Build the SDK host utilities, kernel, root filesystem and firmware archive.
build_sdk

# 7. Build the matching FDD or TDD stack in its own toolchain environment.
build_stack

# 8. Record original stack hashes, then build the SDK's five-file JFFS2 layout.
(
    cd "$app"
    sha256sum -- "${payloads[@]}"
) > "$inputs/inputs.sha256"

cp "$app/startup.cfg" "$inputs/startup.cfg"
gzip -9n "$app/lsmD"
touch "$app/.sign1985"
mkfs.jffs2 --big-endian --pagesize=0x1000 --eraseblock=0x20000 \
    --pad=0xee0000 --faketime --squash-uids --root="$app" --output="$inputs/jffs2.bin"

# 9. Package, validate and publish. The exit trap then removes intermediate files.
"$scripts/build-trx-release.sh" "$inputs" "$out"
