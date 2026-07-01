#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

set -euo pipefail

image_magic_hex() {
    local image_path="$1"
    local byte_count="${2:-6}"

    od -An -tx1 -N"$byte_count" "$image_path" 2>/dev/null | tr -d ' \n'
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
