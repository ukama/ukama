#!/usr/bin/env bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

# Validate the release on the host. The builders also use the layout and format checks below.

set -euo pipefail
export LC_ALL=C

# Physical NOR order. Each size is in bytes.
partitions=(uboot os0 rd0 app0 os1 rd1 app1 env)
declare -A capacity=(
    [uboot]=2228224
    [os0]=5242880
    [rd0]=45088768
    [app0]=15597568
    [os1]=5242880
    [rd1]=45088768
    [app1]=15597568
    [env]=131072
)
payloads=(startup.cfg lsmD dsp.tgz LFMSOFT_OCT_D.tgz)

die() {
    echo "ERROR: $*" >&2
    exit 1
}

need() {
    if ! command -v "$1" >/dev/null; then
        die "Missing host command: $1"
    fi
}

file_size() {
    stat -Lc %s -- "$1"
}

read_region() {
    local image=$1
    local offset=$2
    local length=$3

    dd if="$image" bs=1M iflag=skip_bytes,count_bytes \
        skip="$offset" count="$length" status=none
}

hex_bytes() {
    local file=$1
    local offset=$2
    local length=$3

    od -An -v -tx1 -j "$offset" -N "$length" -- "$file" | tr -d ' \n'
}

file_crc32() {
    # crc32 can append a tab-separated comment when a filename contains a checksum.
    crc32 "$1" | cut -f1
}

check_erased_padding() {
    local image=$1
    local payload_size=$2
    local first_padding_byte=$((payload_size + 1))
    local unexpected_bytes

    # tail counts from one. Every byte after the payload must be erased (0xff).
    unexpected_bytes=$(tail -c "+$first_padding_byte" "$image" | tr -d '\377' | wc -c)
    if (( unexpected_bytes != 0 )); then
        die "Unexpected data after byte $payload_size in $image"
    fi
}

# Each cell setting must occur once and end with a decimal value.
setting() {
    local profile=$1
    local key=$2
    local value

    value=$(awk -v key="$key" '
        {
            sub(/\r$/, "")
            value = $NF
            $NF = ""
            if ($0 == key " ") {
                print value
            }
        }
    ' "$profile")

    if [[ ! $value =~ ^(0|[1-9][0-9]*)$ ]]; then
        die "Expected one numeric '$key' in $profile"
    fi
    echo "$value"
}

# Write the settings carried by a supported 20 MHz release. Read only the serving cell.
radio_profile() {
    local profile=$1
    local mode band duplex srs_period
    local dl_width ul_width dl_channel ul_channel
    local first_channel last_channel
    local tdd_assignment tdd_special_subframe

    dl_width=$(setting "$profile" 'configure cell edlrscblock 0')
    ul_width=$(setting "$profile" 'configure cell freq eulbandwidth 0')
    if [[ $dl_width != 5 || $ul_width != 5 ]]; then
        die "Release profiles require 20 MHz in both directions (SDK bandwidth code 5)"
    fi

    band=$(setting "$profile" 'configure cell u8freqbandindicator 0')
    duplex=$(setting "$profile" 'configure cell etc eduplexmode 0')
    dl_channel=$(setting "$profile" 'configure cell u16dlcarrierfreq 0')
    ul_channel=$(setting "$profile" 'configure cell freq u16ulcarrierfreq 0')
    srs_period=$(setting "$profile" 'configure cell u16srsperiod 0')

    # These inclusive EARFCN limits leave 10 MHz on each side of the carrier.
    if [[ $duplex == 1 && $band == 3 ]]; then
        mode=FDD
        first_channel=1300
        last_channel=1850

        if (( ul_channel != dl_channel + 18000 )); then
            die "Invalid Band 3 uplink channel $ul_channel for downlink $dl_channel"
        fi
        if (( srs_period != 80 )); then
            die "Select the FDD 64UE profile (SRS period 80)"
        fi
    elif [[ $duplex == 0 && ( $band == 40 || $band == 41 ) ]]; then
        mode=TDD
        if [[ $band == 40 ]]; then
            first_channel=38750
            last_channel=39550
        else
            first_channel=39750
            last_channel=41490
        fi

        if (( ul_channel != dl_channel )); then
            die "TDD uplink and downlink channels differ in $profile"
        fi
        tdd_assignment=$(setting "$profile" 'configure tddspecific tddsubframeassignment 0')
        tdd_special_subframe=$(setting "$profile" 'configure tddspecific tddspecialsubframe 0')
        if (( tdd_assignment > 6 || tdd_special_subframe > 9 )); then
            die "Invalid TDD frame settings $tdd_assignment/$tdd_special_subframe"
        fi
    else
        die "Expected FDD/B3, TDD/B40 or TDD/B41 in $profile"
    fi

    if (( dl_channel < first_channel || dl_channel > last_channel )); then
        die "Channel $dl_channel is outside $first_channel..$last_channel for Band $band at 20 MHz"
    fi

    echo "mode=$mode"
    echo "band=$band"
    echo "bandwidth_mhz=20"
    echo "dl_earfcn=$dl_channel"
    echo "ul_earfcn=$ul_channel"
    echo "srs_period=$srs_period"
    if [[ $mode == TDD ]]; then
        echo "tdd_assignment=$tdd_assignment"
        echo "tdd_special_subframe=$tdd_special_subframe"
    fi
}

check_uboot() {
    local image=$1
    local header_size=192
    local image_size payload_size stored_crc computed_crc board_id

    if [[ $(hex_bytes "$image" 8 4) != 424f4f54 ]]; then
        die "Missing OCTEON BOOT header in $image"
    fi
    if [[ $(hex_bytes "$image" 16 2) != 00c0 ]]; then
        die "Unexpected U-Boot header size in $image"
    fi
    board_id=$(od -An --endian=big -tu2 -j 22 -N 2 "$image")
    if (( board_id != 59 )); then
        die "U-Boot is not TRX board 59"
    fi

    image_size=$(file_size "$image")
    payload_size=$(od -An --endian=big -tu4 -j 24 -N 4 "$image")
    if (( payload_size == 0 || header_size + payload_size != image_size )); then
        die "U-Boot payload length does not match the file size"
    fi
    if (( image_size > 0xc0000 )); then
        die "U-Boot exceeds the failsafe slot"
    fi

    # The OCTEON header checksum excludes its own four bytes at offset 12.
    head -c 12 "$image" > "$temp/uboot.header"
    read_region "$image" 16 "$((header_size - 16))" >> "$temp/uboot.header"
    stored_crc=$(hex_bytes "$image" 12 4)
    computed_crc=$(file_crc32 "$temp/uboot.header")
    if [[ $computed_crc != "$stored_crc" ]]; then
        die "U-Boot header CRC mismatch"
    fi

    read_region "$image" "$header_size" "$payload_size" > "$temp/uboot.payload"
    stored_crc=$(hex_bytes "$image" 28 4)
    computed_crc=$(file_crc32 "$temp/uboot.payload")
    if [[ $computed_crc != "$stored_crc" ]]; then
        die "U-Boot payload CRC mismatch"
    fi
}

check_env_size() {
    local record_size=$1

    if [[ ! $record_size =~ ^[1-9][0-9]{0,5}$ ]]; then
        die "Invalid U-Boot environment size: $record_size"
    fi
    if (( record_size <= 4 || record_size > capacity[env] )); then
        die "Environment record must fit in the ${capacity[env]}-byte partition"
    fi
}

check_environment() {
    local image=$1
    local record_size=$2
    local data_file="$temp/env.data"
    local settings_file="$temp/env.settings"
    local data_size stored_crc computed_crc

    check_env_size "$record_size"

    # The first four bytes hold the big-endian CRC32 of the rest of the record.
    data_size=$((record_size - 4))
    read_region "$image" 4 "$data_size" > "$data_file"
    stored_crc=$(hex_bytes "$image" 0 4)
    computed_crc=$(file_crc32 "$data_file")
    if [[ $computed_crc != "$stored_crc" ]]; then
        die "U-Boot environment CRC mismatch"
    fi

    check_erased_padding "$image" "$record_size"

    # A second NUL terminates the list of name=value strings.
    tr '\000' '\n' < "$data_file" | sed '/^$/,$d' > "$settings_file"
    if ! grep -axq 'bootby=flash' "$settings_file"; then
        die "U-Boot environment does not select flash boot"
    fi
    if ! grep -aq '^bootcmd=.' "$settings_file"; then
        die "U-Boot environment has no bootcmd"
    fi
}

check_padded_copy() {
    local payload=$1
    local partition=$2
    local size

    size=$(file_size "$payload")
    if ! cmp -s -n "$size" "$payload" "$partition"; then
        die "$partition differs from $payload"
    fi
    check_erased_padding "$partition" "$size"
}

check_firmware_member() {
    local payload=$1
    local partition=$2
    local member=$3
    local tag=$4
    local boot_limit=$5
    local size expected_tag stored_tag

    size=$(file_size "$payload")
    if (( size == 0 || size > boot_limit )); then
        die "$payload exceeds its boot-time size limit or is empty"
    fi
    if [[ $(hex_bytes "$payload" 0 4) != 1f8b0808 ]]; then
        die "Missing SDK gzip header in $payload"
    fi

    expected_tag=$(printf '[%s %08x]' "$tag" "$size")
    stored_tag=$(read_region "$payload" 10 "${#expected_tag}")
    if [[ $stored_tag != "$expected_tag" ]]; then
        die "Expected SDK tag '$expected_tag' in $payload"
    fi
    if ! gzip -t "$payload"; then
        die "Invalid gzip stream in $payload"
    fi

    check_padded_copy "$payload" "$partition"
    tar -xOf fw.bin "$member" > "$temp/firmware-member.gz"
    if ! cmp -s "$payload" "$temp/firmware-member.gz"; then
        die "$member in fw.bin differs from $payload"
    fi
}

extract_app() {
    local image=$1
    local dest=$2
    local file

    if [[ $(hex_bytes "$image" 0 2) != 1985 ]]; then
        die "Expected a big-endian JFFS2 image"
    fi
    if ! jffs2dump -b -c "$image" > "$temp/jffs2.log" 2>&1; then
        cat "$temp/jffs2.log" >&2
        die "Cannot read $image"
    fi
    if grep -Ei 'wrong|unknown node|error|truncat' "$temp/jffs2.log" >&2; then
        die "Invalid JFFS2 node"
    fi
    if ! jefferson -d "$dest" "$image" > "$temp/extract.log" 2>&1; then
        cat "$temp/extract.log" >&2
        die "Cannot extract $image"
    fi

    for file in .sign1985 startup.cfg lsmD.gz dsp.tgz LFMSOFT_OCT_D.tgz; do
        if [[ ! -f $dest/$file || -L $dest/$file ]]; then
            die "Missing regular application file: $file"
        fi
    done
    if [[ -e $dest/lsmD || -L $dest/lsmD ]]; then
        die "Unexpected lsmD in compressed stack image"
    fi

    gzip -dc "$dest/lsmD.gz" > "$dest/lsmD"
    if [[ ! -s $dest/lsmD ]]; then
        die "Empty stack executable"
    fi
    tar -tzf "$dest/dsp.tgz" > /dev/null
    tar -tzf "$dest/LFMSOFT_OCT_D.tgz" > /dev/null
}

release_hashes() {
    local part
    for part in "${partitions[@]}"; do
        sha256sum "flash_$part.img"
    done
    sha256sum flash-trx.img startup.cfg radio.conf inputs.sha256 env.size build.info
    sha256sum u-boot-octeon_trx.bin lsm_os_trx.gz lsm_rd_trx.gz fw.bin host-tools/oct-remote-*
}

check_partitions() {
    local offset=0
    local part image expected_size actual_size

    if [[ $(file_size flash-trx.img) != 134217728 ]]; then
        die "Expected a 128 MiB TRX image"
    fi
    for part in "${partitions[@]}"; do
        image="flash_$part.img"
        expected_size=${capacity[$part]}
        actual_size=$(file_size "$image")
        if (( actual_size != expected_size )); then
            die "$image has $actual_size bytes; expected $expected_size"
        fi
        if ! cmp -s -n "$expected_size" "$image" flash-trx.img 0 "$offset"; then
            die "Wrong $part region in flash-trx.img at offset $offset"
        fi
        offset=$((offset + expected_size))
    done

    for part in os rd app; do
        if ! cmp -s "flash_${part}0.img" "flash_${part}1.img"; then
            die "$part slots differ"
        fi
    done
}

check_application() {
    local mode

    extract_app flash_app0.img "$temp/app"
    if ! cmp -s startup.cfg "$temp/app/startup.cfg"; then
        die "Wrong startup.cfg in application image"
    fi

    radio_profile "$temp/app/startup.cfg" > "$temp/radio.conf"
    if ! cmp -s radio.conf "$temp/radio.conf"; then
        die "Radio settings differ from release"
    fi
    mode=$(sed -n 's/^mode=//p' "$temp/radio.conf")
    if ! grep -Fxq "stack_mode=$mode" build.info; then
        die "Stack build mode differs from startup.cfg"
    fi
    if ! grep -Fxq 'ue_target=64' build.info; then
        die "Missing 64-UE release designation"
    fi

    (
        cd "$temp/app"
        sha256sum -- "${payloads[@]}"
    ) > "$temp/inputs.sha256"
    if ! cmp -s inputs.sha256 "$temp/inputs.sha256"; then
        die "Embedded profile or stack payload differs from build inputs"
    fi
}

validate_release() {
    local release=$1
    local env_size

    cd -- "$release"
    release_hashes > "$temp/actual.sha256"
    if ! cmp -s SHA256SUMS "$temp/actual.sha256"; then
        die "Release checksum mismatch"
    fi
    check_partitions
    echo "PASS  Partition sizes, full-image offsets, A/B slots and checksums"

    check_uboot u-boot-octeon_trx.bin
    check_padded_copy u-boot-octeon_trx.bin flash_uboot.img
    env_size=$(<env.size)
    check_environment flash_env.img "$env_size"

    check_firmware_member lsm_os_trx.gz flash_os0.img 03_lsm_os.gz LSM_OS "${capacity[os0]}"
    # The TRX boot command copies at most 40 MiB of ramdisk into RAM.
    check_firmware_member lsm_rd_trx.gz flash_rd0.img 01_lsm_rd.gz LSM_RD 41943040
    if [[ ! -x host-tools/oct-remote-boot ]]; then
        die "Missing executable host-tools/oct-remote-boot"
    fi
    echo "PASS  TRX U-Boot, environment, kernel and ramdisk"

    check_application
    echo "PASS  Embedded radio profile, lsmD, DSP and LFM payloads"
}

if [[ ${BASH_SOURCE[0]} == "$0" ]]; then
    if [[ $# != 1 ]]; then
        die "Usage: $0 <release-dir>"
    fi

    host_env=${XDG_CACHE_HOME:-$HOME/.cache}/ukama/trx-host
    export PATH="$host_env/bin:$PATH:/usr/sbin:/sbin"
    need crc32
    need jffs2dump
    need jefferson

    temp=$(mktemp -d)
    trap 'rm -rf -- "$temp"' EXIT
    trap 'exit 130' INT
    trap 'exit 143' TERM
    validate_release "$1"
fi
