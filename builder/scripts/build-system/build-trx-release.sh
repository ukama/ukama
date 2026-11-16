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

if [[ $# != 2 ]]; then
    die "Usage: $0 <built-input-dir> <release-dir>"
fi
input=$(realpath -e -- "$1")
out=$(realpath -m -- "$2")
if [[ $input == "$out" ]]; then
    die "Input and release directories must differ"
fi

need mkenvimage
need crc32
required_files=(
    u-boot-octeon_trx.bin lsm_os.gz lsm_rd.gz fw.bin
    jffs2.bin startup.cfg inputs.sha256 env.size build.info
    host-tools/oct-remote-boot
)
for file in "${required_files[@]}"; do
    if [[ ! -s $input/$file ]]; then
        die "Missing packaging input: $input/$file"
    fi
done

parent=$(dirname -- "$out")
mkdir -p -- "$parent"
temp=$(mktemp -d -- "$parent/.trx-release.XXXXXXXX")
trap 'rm -rf -- "$temp"' EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

release=$temp/release
mkdir "$release"
cd "$release"

# Construct the bootloader, environment and primary partitions.
check_uboot "$input/u-boot-octeon_trx.bin"
write_partition "$input/u-boot-octeon_trx.bin" flash_uboot.img "${capacity[uboot]}"
write_environment
write_partition "$input/lsm_os.gz" flash_os0.img "${capacity[os0]}"
write_partition "$input/lsm_rd.gz" flash_rd0.img "${capacity[rd0]}"

app_size=$(file_size "$input/jffs2.bin")
if (( app_size != capacity[app0] )); then
    die "Wrong JFFS2 size: $app_size bytes"
fi
cp -- "$input/jffs2.bin" flash_app0.img

# Populate slot 1, then assemble the combined image in physical NOR order.
for part in os rd app; do
    cp -- "flash_${part}0.img" "flash_${part}1.img"
done
for part in "${partitions[@]}"; do
    cat "flash_$part.img" >> flash-trx.img
done

# Include the unpadded bootstrap images, flash tools and build records.
cp -- "$input"/{startup.cfg,inputs.sha256,env.size,build.info,fw.bin} .
cp -- "$input/u-boot-octeon_trx.bin" .
cp -- "$input/lsm_os.gz" lsm_os_trx.gz
cp -- "$input/lsm_rd.gz" lsm_rd_trx.gz
cp -R -- "$input/host-tools" .
radio_profile startup.cfg > radio.conf
release_hashes > SHA256SUMS

"$scripts/validate-trx-image.sh" "$release"
publish_release
