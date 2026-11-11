#!/usr/bin/env bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

set -euo pipefail

usage() {
    echo "Usage: $0 UKAMA_ROOT [BOARD]" >&2
    exit 1
}

[[ $# -ge 1 && $# -le 2 ]] || usage

UKAMA_ROOT="$(realpath "$1")"
BOARD="${2:-}"
BUILDER_ROOT="${UKAMA_ROOT}/builder"
CONFIG_DIR="${BUILDER_ROOT}/configs"
PKG_DIR="${UKAMA_ROOT}/build/pkgs"

export BUILD_MODE="${BUILD_MODE:-debug}"
export UKAMA_ROOT

[ -d "${BUILDER_ROOT}" ] || {
    echo "Builder directory not found: ${BUILDER_ROOT}" >&2
    exit 1
}

[ -d "${CONFIG_DIR}" ] || {
    echo "Builder config directory not found: ${CONFIG_DIR}" >&2
    exit 1
}

if [[ -n "$BOARD" ]]; then
    common_config="${BUILDER_ROOT}/boards/common.config"
    board_config="${BUILDER_ROOT}/boards/${BOARD}.config"

    for config in "$common_config" "$board_config"; do
        if [[ ! -f "$config" ]]; then
            echo "Board config not found: $config" >&2
            exit 1
        fi
    done

    source "${BUILDER_ROOT}/scripts/build-system/pkg-utils.sh"
    get_enabled_apps "$common_config" "$board_config"

    if [[ ${#APPS[@]} -eq 0 ]]; then
        echo "No apps enabled for board: $BOARD" >&2
        exit 1
    fi

    configs=()
    for app in "${APPS[@]}"; do
        config="${CONFIG_DIR}/${app}.toml"
        if [[ ! -f "$config" ]]; then
            echo "Missing app config for board $BOARD: $config" >&2
            exit 1
        fi
        configs+=("$config")
    done
else
    configs=("${CONFIG_DIR}"/*.toml)
fi

cd "${BUILDER_ROOT}"
package_work=$(mktemp -d)

cleanup() {
    rm -rf -- "$package_work"
    make -C "$BUILDER_ROOT" clean >/dev/null 2>&1 || true
}
trap cleanup EXIT

make clean
make app_builder

# app_builder writes packages here. Clean the actual output directory so
# failed or removed apps cannot leave stale *_latest.tar.gz packages behind.
rm -rf -- "${PKG_DIR}"
mkdir -p -- "${PKG_DIR}"

found=0
for config in "${configs[@]}"; do
    [ -f "${config}" ] || continue
    found=1
    (cd "$package_work" && "${BUILDER_ROOT}/app_builder" --create --config "${config}")
done

[ "${found}" -eq 1 ] || {
    echo "No app configs found in: ${CONFIG_DIR}" >&2
    exit 1
}
