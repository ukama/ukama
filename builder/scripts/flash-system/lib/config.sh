#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2025-present, Ukama Inc.

YQ_BIN="${YQ_BIN:-${FLASH_DIR}/.bin/yq}"

ensure_yq() {
    if [ ! -x "$YQ_BIN" ]; then
        mkdir -p "$(dirname "$YQ_BIN")"
        curl -sL https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 -o "$YQ_BIN"
        chmod +x "$YQ_BIN"
    fi
}

yq_read() {
    local file="$1"
    local path="$2"
    "$YQ_BIN" eval ".${path}" "$file"
}

yq_exists() {
    local file="$1"
    local path="$2"
    local val
    val=$("$YQ_BIN" eval ".${path} // \"__MISSING__\"" "$file")
    [ "$val" != "__MISSING__" ] && [ "$val" != "null" ]
}

yq_count() {
    local file="$1"
    local path="$2"
    "$YQ_BIN" eval ".${path} | length" "$file"
}

yq_keys() {
    local file="$1"
    local path="$2"
    "$YQ_BIN" eval ".${path} | keys | .[]" "$file"
}

yq_join() {
    local file="$1"
    local path="$2"
    local sep="${3:-,}"
    "$YQ_BIN" eval ".${path} | join(\"${sep}\")" "$file"
}
