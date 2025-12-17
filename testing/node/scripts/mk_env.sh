#!/usr/bin/env bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2022-present, Ukama Inc.

set -euo pipefail

BUILD_ENV=container

UKAMA_OS_TAR_GLOB="/ukama/ukama_*.tgz"
EXTRACT_ROOT="/tmp/virtnode"          # where we extract inside container
UKAMA_OS_PATH=""                      # will be discovered after extract

if_host() {
    # Your logic is fine; keep it as-is (quoted + safer)
    local val
    val="$(grep -i "pids" /proc/1/cgroup | awk -F":" 'NR==1{print $NF}')"
    if [[ "$val" == "/init.scope" || "$val" == "/" ]]; then
        BUILD_ENV=local
    fi
}

pick_tarball() {
    shopt -s nullglob
    local matches=( $UKAMA_OS_TAR_GLOB )
    shopt -u nullglob

    if (( ${#matches[@]} == 0 )); then
        echo "ERROR: No tarball found matching: $UKAMA_OS_TAR_GLOB" >&2
        exit 1
    fi
    if (( ${#matches[@]} > 1 )); then
        echo "ERROR: Multiple tarballs found:" >&2
        printf '  %s\n' "${matches[@]}" >&2
        echo "Please leave only one, or tighten the glob." >&2
        exit 1
    fi

    echo "${matches[0]}"
}

extract_source() {
    local tarball
    tarball="$(pick_tarball)"

    mkdir -p "$EXTRACT_ROOT"
    echo "Extracting: $tarball -> $EXTRACT_ROOT"
    tar -zxf "$tarball" -C "$EXTRACT_ROOT"

    # Find nodes/ukamaOS anywhere under EXTRACT_ROOT
    local found
    found="$(find "$EXTRACT_ROOT" -type d -path "*/nodes/ukamaOS" -print -quit || true)"
    if [[ -z "$found" ]]; then
        echo "ERROR: After extract, could not find */nodes/ukamaOS under $EXTRACT_ROOT" >&2
        echo "Tip: inspect tar contents with: tar -tzf \"$tarball\" | head" >&2
        exit 1
    fi

    UKAMA_OS_PATH="$found"
}

# main
if_host
echo "Build environment is $BUILD_ENV"

if [[ "$BUILD_ENV" == "local" ]]; then
    UKAMA_OS_PATH="$(realpath ../../nodes/ukamaOS)"
elif [[ "$BUILD_ENV" == "container" ]]; then
    extract_source
else
    echo "Unknown environment: $BUILD_ENV" >&2
    exit 1
fi

if [[ -d "$UKAMA_OS_PATH" ]]; then
    export UKAMA_OS="$UKAMA_OS_PATH"
    echo "Build environment is set for the Virtual Node on $BUILD_ENV."
    echo "UKAMA_OS=$UKAMA_OS"
    exit 0
else
    echo "UkamaOS not found at: $UKAMA_OS_PATH" >&2
    exit 1
fi
