#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
UNIMGC_DIR="${SCRIPT_DIR}/unimgc"
UNIMGC_BIN="${SCRIPT_DIR}/.bin/unimgc"
HDD_RAW_COPY_IMGC_MAGIC="114844442052617720436f707920546f6f6c"

image_magic_hex() {
    local image_path="$1"
    local byte_count="${2:-6}"

    od -An -tx1 -N"$byte_count" "$image_path" 2>/dev/null | tr -d ' \n'
}

is_hdd_raw_copy_imgc() {
    local image_path="$1"

    [[ "$(image_magic_hex "$image_path" 18)" == "$HDD_RAW_COPY_IMGC_MAGIC" ]]
}

image_has_mbr_signature() {
    local sig
    sig=$(dd if="$1" bs=1 skip=510 count=2 2>/dev/null | od -An -tx1 | tr -d ' \n')
    [[ "$sig" == "55aa" ]]
}

is_probably_raw_disk_image() {
    local image_path="$1"

    if image_has_mbr_signature "$image_path"; then
        return 0
    fi

    fdisk -l "$image_path" >/dev/null 2>&1
}

find_c_compiler() {
    local candidate=""

    for candidate in cc gcc clang; do
        if command -v "$candidate" >/dev/null 2>&1; then
            printf '%s\n' "$candidate"
            return 0
        fi
    done

    return 1
}

ensure_unimgc() {
    local compiler=""
    local src=""
    local -a sources=(
        "${UNIMGC_DIR}/unimgc.c"
        "${UNIMGC_DIR}/image.c"
        "${UNIMGC_DIR}/lzo.c"
    )

    if [ -x "$UNIMGC_BIN" ]; then
        return 0
    fi

    if ! compiler=$(find_c_compiler); then
        echo "HDD Raw Copy .imgc support requires a C compiler (cc, gcc, or clang)" >&2
        return 1
    fi

    for src in "${sources[@]}" "${UNIMGC_DIR}/image.h" "${UNIMGC_DIR}/lzo.h" "${UNIMGC_DIR}/endian.h"; do
        if [ ! -f "$src" ]; then
            echo "Missing unimgc source file: $src" >&2
            return 1
        fi
    done

    mkdir -p "$(dirname "$UNIMGC_BIN")"
    "$compiler" -O2 -std=c99 -D_FILE_OFFSET_BITS=64 \
        -o "$UNIMGC_BIN" \
        "${sources[@]}"
}

detect_image_format() {
    local image_path="$1"
    local magic=""

    if is_hdd_raw_copy_imgc "$image_path"; then
        echo "imgc"
        return 0
    fi

    magic=$(image_magic_hex "$image_path")
    case "$magic" in
        1f8b*) echo "gzip" ;;
        fd377a585a00*) echo "xz" ;;
        28b52ffd*) echo "zstd" ;;
        425a68*) echo "bzip2" ;;
        *)
            if [[ "$image_path" == *.imgc ]]; then
                if is_probably_raw_disk_image "$image_path"; then
                    echo "raw"
                else
                    echo "unknown"
                    return 1
                fi
            else
                echo "raw"
            fi
            ;;
    esac
}

require_image_tool() {
    local tool="$1"

    if ! command -v "$tool" >/dev/null 2>&1; then
        echo "Required tool '$tool' is not installed" >&2
        return 1
    fi
}

stream_image_to_stdout() {
    local image_path="$1"
    local image_format="$2"

    case "$image_format" in
        raw)
            cat -- "$image_path"
            ;;
        imgc)
            ensure_unimgc
            "$UNIMGC_BIN" "$image_path"
            ;;
        gzip)
            require_image_tool gzip
            gzip -dc -- "$image_path"
            ;;
        xz)
            require_image_tool xz
            xz -dc -- "$image_path"
            ;;
        zstd)
            require_image_tool zstd
            zstd -dc -- "$image_path"
            ;;
        bzip2)
            require_image_tool bzip2
            bzip2 -dc -- "$image_path"
            ;;
        *)
            echo "Unsupported image format for '$image_path'" >&2
            return 1
            ;;
    esac
}

prepare_image_for_host_use() {
    local image_path="$1"
    local output_dir="$2"
    local image_format=""
    local prepared_path=""

    image_format=$(detect_image_format "$image_path")
    if [[ "$image_format" == "raw" ]]; then
        echo "$image_path"
        return 0
    fi

    mkdir -p "$output_dir"
    prepared_path=$(mktemp "${output_dir}/prepared-image.XXXXXX.img")
    stream_image_to_stdout "$image_path" "$image_format" > "$prepared_path"
    echo "$prepared_path"
}
