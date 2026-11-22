#!/usr/bin/env bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

# Package the built inputs, verify the staged release, then publish it.

set -euo pipefail
umask 022

script_path=$(realpath -- "$0")
scripts=$(dirname -- "$script_path")
source "$scripts/validate-trx-image.sh"

write_partition() {
    local source=$1
    local image=$2
    local partition_size=$3
    local used padding

    used=$(file_size "$source")
    if (( used == 0 || used > partition_size )); then
        die "$source is empty or exceeds the $partition_size bytes available in $image"
    fi

    cp -- "$source" "$image"
    padding=$((partition_size - used))
    head -c "$padding" /dev/zero | tr '\000' '\377' >> "$image"
}

write_uboot_partition() {
    local binary=$1 size offset
    check_uboot "$binary"
    size=$(file_size "$binary")
    offset=$(normal_uboot_offset "$size")
    (( offset + size <= capacity[uboot] )) || die "Both U-Boot copies must fit in flash_uboot.img"
    write_partition "$binary" flash_uboot.img "${capacity[uboot]}"
    dd if="$binary" of=flash_uboot.img bs=65536 seek="$((offset / 65536))" conv=notrunc status=none
    printf 'TRX U-Boot: reset at 0x0, normal boot at 0x%x\n' "$offset"
}

write_environment() {
    local record_size
    record_size=$(<"$input/env.size")
    check_env_size "$record_size"

    # TRX bring-up commands and the SDK's flash boot command.
    cat > "$temp/env.txt" <<'ENV'
baudrate=115200
bootby=flash
cfgloadby=flash
swloadby=flash
i2cinit=i2c dev 0; i2c probe; i2c dev 1; i2c probe
bootcmd=run i2cinit; run namedalloc; run bootcby${bootby}
namedalloc=namedalloc dsp-dump 0x400000 0x7f4D0000; namedalloc cazac 0x630000 0x7f8D0000; namedalloc cpu-dsp-if 0x100000 0x7ff00000; namedalloc dsp-log-buf 0x4000000 0x80000000; namedalloc initrd 0x2800000 0x30800000;
bootcbyflash=gunzip 0x17e20000 0x20000000; cp.l 0x18320000 0x30800000 0xA00000; bootoctlinux 0x20000000 coremask=0x7 endbootargs rd_name=initrd mem=512M;
bootcbytftp=tftp 0x21000000 lsm_os.gz; gunzip 0x21000000 0x20000000 0x1000000; tftp 0x30800000 lsm_rd.gz; bootoctlinux 0x20000000 coremask=0x7 endbootargs rd_name=initrd mem=512M;
mk_ubootenv=1
ENV

    # mkenvimage writes the big-endian CRC and terminates the environment entries.
    mkenvimage -b -p 0xff -s "$record_size" -o "$temp/env.record" "$temp/env.txt"
    write_partition "$temp/env.record" "$release/flash_env.img" "${capacity[env]}"
}

# S06network takes eth0 down before it applies the release IP and netmask.
# Complete older builder records too. Keep explicit site values and boot commands.
complete_network_environment() {
    local record_size entry key value used=0 terminated=no end extra crc
    local -A defaults=([ipaddr]=10.102.81.61 [netmask]=255.255.255.0 [gatewayip]=10.102.81.1)
    local -A seen=()
    local entries="$temp/network.entries" data="$temp/network.data"
    record_size=$(<env.size)
    check_environment flash_env.img "$record_size"

    while IFS= read -r -d '' entry; do
        used=$((used + ${#entry} + 1))
        if [[ -z $entry ]]; then terminated=yes; break; fi
        [[ $entry == ?*=* ]] || die "Malformed U-Boot environment entry"
        key=${entry%%=*}
        case $key in
            ipaddr|netmask|gatewayip)
                [[ -z ${seen[$key]:-} ]] || die "Duplicate U-Boot setting $key"
                seen[$key]=yes
                value=${entry#*=}
                entry="$key=${value:-${defaults[$key]}}"
                ;;
        esac
        printf '%s\0' "$entry"
    done < "$temp/env.data" > "$entries"
    [[ $terminated == yes ]] || die "Unterminated U-Boot environment"
    for key in ipaddr netmask gatewayip; do
        if [[ -z ${seen[$key]:-} ]]; then printf '%s=%s\0' "$key" "${defaults[$key]}" >> "$entries"; fi
    done
    printf '\0' >> "$entries"

    end=$(file_size "$entries")
    (( end <= record_size - 4 )) || die "Network settings exceed the U-Boot environment record"
    if cmp -s -n "$end" "$entries" "$temp/env.data"; then return 0; fi
    extra=$((end - used))
    if (( extra > 0 )) && [[ $(read_region "$temp/env.data" "$used" "$extra" | tr -d '\000\377' | wc -c) -ne 0 ]]; then
        die "Network settings would overwrite occupied environment bytes"
    fi
    { cat "$entries"; tail -c "+$((end + 1))" "$temp/env.data"; } > "$data"
    crc=$(file_crc32 "$data")
    [[ $crc =~ ^[[:xdigit:]]{8}$ ]] || die "Invalid environment CRC result"
    {
        printf '%b' "\\x${crc:0:2}\\x${crc:2:2}\\x${crc:4:2}\\x${crc:6:2}"
        cat "$data"
        tail -c "+$((record_size + 1))" flash_env.img
    } > "$temp/network.img"
    mv -f -- "$temp/network.img" flash_env.img
}

# Reuse the builder's JFFS2 layout. Reproducing the old slot first protects
# file contents and metadata that an extractor might otherwise omit.
pack_application() {
    mkfs.jffs2 --big-endian --pagesize=0x1000 --eraseblock=0x20000 \
        --pad="${capacity[app0]}" --faketime --squash-uids --root="$1" --output="$2"
}

install_ethernet_hook() {
    local tree="$temp/app-tree" image="$temp/app-slot.img"
    local hook="$tree/rc_post.local" source_hook="$scripts/trx-rc-post.local"

    extract_app flash_app0.img "$tree"
    # extract_app also expands this file for validation. NOR stores lsmD.gz.
    rm -- "$tree/lsmD"
    if [[ -e $hook || -L $hook ]]; then
        [[ -f $hook && ! -L $hook ]] && cmp -s -- "$source_hook" "$hook" ||
            die "Existing rc_post.local differs. Preserve its logic when merging the TRX Ethernet hook."
        # Some host extractors strip execute bits. Repack with the intended
        # mode and compare the image itself before accepting an existing hook.
        chmod 0755 "$hook"
        pack_application "$tree" "$image"
        if cmp -s -- "$image" flash_app0.img; then return 0; fi
        die "Cannot reproduce the application slot with its executable Ethernet hook"
    fi

    pack_application "$tree" "$image"
    cmp -s -- "$image" flash_app0.img ||
        die "Cannot reproduce the original application slot with this host's mkfs.jffs2"
    install -m 0755 -- "$source_hook" "$hook"
    pack_application "$tree" "$image"
    [[ $(file_size "$image") == "${capacity[app0]}" ]] || die "Ethernet hook exceeds the application slot"
    # Replace both private staging files (the published inputs may be read-only).
    rm -- flash_app0.img flash_app1.img
    mv -- "$image" flash_app0.img
    cp -- flash_app0.img flash_app1.img
}

publish_release() {
    local first_entry

    # Hold the lock across both the comparison and the rename.
    exec 9< "$parent"
    flock -x 9

    if [[ -d $out ]]; then
        first_entry=$(find "$out" -mindepth 1 -maxdepth 1 -print -quit)
        if [[ -n $first_entry ]]; then
            if ! diff -qr -- "$out" "$release"; then
                die "Existing release differs. Use a new output directory."
            fi
            echo "Reproduced $out"
            return
        fi
    elif [[ -e $out ]]; then
        die "Release path is not a directory: $out"
    fi

    # Staging lives on the same filesystem. rename also accepts an empty destination.
    mv -T -- "$release" "$out"
    echo "Built and verified $out"
}

finish_release() {
    local part
    complete_network_environment
    install_ethernet_hook
    for part in "${partitions[@]}"; do
        cat "flash_$part.img"
    done > flash-trx.img
    release_hashes > SHA256SUMS
    "$scripts/validate-trx-image.sh" "$release"
    publish_release
}

from_release=no
if [[ ${1:-} == --from-release ]]; then
    from_release=yes
    shift
fi
if [[ $# != 2 ]]; then
    die "Usage: $0 <built-input-dir> <release-dir> or $0 --from-release <existing-builder-release> <new-release-dir>"
fi
input=$(realpath -e -- "$1")
out=$(realpath -m -- "$2")
if [[ $input == "$out" || $from_release == yes && ( $out == "$input/"* || $input == "$out/"* ) ]]; then
    die "Use a separate output directory for the new release"
fi

host_env=${XDG_CACHE_HOME:-$HOME/.cache}/ukama/trx-host
export PATH="$host_env/bin:$PATH:/usr/sbin:/sbin"
need crc32
need jefferson
need jffs2dump
need mkfs.jffs2

parent=$(dirname -- "$out")
mkdir -p -- "$parent"
temp=$(mktemp -d -- "$parent/.trx-release.XXXXXXXX")
trap 'rm -rf -- "$temp"' EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

release=$temp/release
mkdir "$release"
cd "$release"

if [[ $from_release == yes ]]; then
    # Work on a snapshot. Validate the old partition and combined image
    # before changing either, then run the full validator on the result.
    cp -R --reflink=auto -- "$input/." .
    [[ -z $(find . ! -type d ! -type f -print -quit) ]] || die "Expected only regular release files and directories"
    release_hashes > "$temp/source.sha256"
    cmp -s SHA256SUMS "$temp/source.sha256" || die "Source release checksum mismatch"
    check_partitions
    # Accept the original single-copy layout or the verified V26 layout.
    check_uboot u-boot-octeon_trx.bin
    normal_offset=$(normal_uboot_offset "$(file_size u-boot-octeon_trx.bin)")
    if [[ $(hex_bytes flash_uboot.img "$normal_offset" 4) == ffffffff ]]; then
        check_padded_copy u-boot-octeon_trx.bin flash_uboot.img
    else
        check_uboot_partition u-boot-octeon_trx.bin flash_uboot.img
    fi
    # Published inputs may be read-only. Replace files in the snapshot.
    chmod u+w .
    rm -- flash_uboot.img flash-trx.img SHA256SUMS
    write_uboot_partition u-boot-octeon_trx.bin
    finish_release
    exit 0
fi

need mkenvimage
required_files=(
    u-boot-octeon_trx.bin lsm_os.gz lsm_rd.gz fw.bin
    jffs2.bin startup.cfg inputs.sha256 env.size build.info
    host-tools/oct-remote-boot
)
for file in "${required_files[@]}"; do
    [[ -s $input/$file ]] || die "Missing packaging input: $input/$file"
done

# Construct the bootloader, environment and primary partitions.
write_uboot_partition "$input/u-boot-octeon_trx.bin"
write_environment
write_partition "$input/lsm_os.gz" flash_os0.img "${capacity[os0]}"
write_partition "$input/lsm_rd.gz" flash_rd0.img "${capacity[rd0]}"

app_size=$(file_size "$input/jffs2.bin")
if (( app_size != capacity[app0] )); then
    die "Wrong JFFS2 size: $app_size bytes"
fi
cp -- "$input/jffs2.bin" flash_app0.img

# Populate slot 1.
for part in os rd app; do
    cp -- "flash_${part}0.img" "flash_${part}1.img"
done

# Include the unpadded bootstrap images, flash tools and build records.
cp -- "$input"/{startup.cfg,inputs.sha256,env.size,build.info,fw.bin} .
cp -- "$input/u-boot-octeon_trx.bin" .
cp -- "$input/lsm_os.gz" lsm_os_trx.gz
cp -- "$input/lsm_rd.gz" lsm_rd_trx.gz
cp -R -- "$input/host-tools" .
radio_profile startup.cfg > radio.conf
finish_release
