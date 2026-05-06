#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

image_magic_hex() {
    local image_path="$1"
    local byte_count="${2:-6}"

    od -An -tx1 -N"$byte_count" "$image_path" 2>/dev/null | tr -d ' \n'
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

detect_image_format() {
    local image_path="$1"
    local magic=""

    magic=$(image_magic_hex "$image_path")
    case "$magic" in
        1f8b*) echo "gzip" ;;
        fd377a585a00*) echo "xz" ;;
        28b52ffd*) echo "zstd" ;;
        425a68*) echo "bzip2" ;;
        *)     echo "raw" ;;
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

validate_image_for_board() {
    local image_path="$1"
    local board_name="$2"

    if [ ! -f "$image_path" ]; then
        echo "Image not found: $image_path" >&2
        return 1
    fi

    if ! image_has_mbr_signature "$image_path"; then
        echo "Image missing MBR signature" >&2
        return 1
    fi

    return 0
}

image_partition_count() {
    local image_path="$1"
    local base

    base=$(basename "$image_path")
    fdisk -l "$image_path" 2>/dev/null | grep -c "^${base}[0-9]" || true
}
