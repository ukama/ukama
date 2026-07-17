#!/usr/bin/env bash

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

set -euo pipefail

usage() {
    echo "Usage: $0 UKAMA_ROOT" >&2
    exit 1
}

[ "$#" -eq 1 ] || usage

UKAMA_ROOT="$(realpath "$1")"
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

cd "${BUILDER_ROOT}"

cleanup() {
    make clean >/dev/null 2>&1 || true
}
trap cleanup EXIT

make clean
make app_builder

# app_builder writes packages here. Clean the actual output directory so
# failed or removed apps cannot leave stale *_latest.tar.gz packages behind.
rm -rf -- "${PKG_DIR}"
mkdir -p -- "${PKG_DIR}"

found=0
for config in "${CONFIG_DIR}"/*.toml; do
    [ -f "${config}" ] || continue
    found=1
    ./app_builder --create --config "${config}"
done

[ "${found}" -eq 1 ] || {
    echo "No app configs found in: ${CONFIG_DIR}" >&2
    exit 1
}
