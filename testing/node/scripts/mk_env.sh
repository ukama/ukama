#!/usr/bin/env bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2022-present, Ukama Inc.

set -euo pipefail

UKAMA_OS_TAR_GLOB="/ukama/ukama.tgz"
EXTRACT_ROOT="/tmp/virtnode"
UKAMA_OS_PATH=""

detect_env() {
    # Respect explicit user setting: BUILD_ENV=local ./script.sh
    if [[ -n "${BUILD_ENV:-}" && ( "$BUILD_ENV" == "local" || "$BUILD_ENV" == "container" ) ]]; then
        return
    fi

    BUILD_ENV="local"
    if [[ -f "/.dockerenv" || -f "/run/.containerenv" ]]; then
        BUILD_ENV="container"
        return
    fi
    if [[ -n "${container:-}" || -n "${CONTAINER:-}" ]]; then
        BUILD_ENV="container"
        return
    fi
    if grep -qaE '(docker|podman|containerd|kubepods|lxc)' /proc/1/cgroup 2>/dev/null; then
        BUILD_ENV="container"
        return
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
        exit 1
    fi

    echo "${matches[0]}"
}

extract_source() {
    local tarball found
    tarball="$(pick_tarball)"

    mkdir -p "$EXTRACT_ROOT"
    echo "Extracting: $tarball -> $EXTRACT_ROOT"

    # Use parallel decompression when possible (tgz/tar.gz), otherwise plain tar
    case "$tarball" in
        *.tar.gz|*.tgz)
            if command -v pigz >/dev/null 2>&1; then
                pigz -dc "$tarball" | tar -xf - -C "$EXTRACT_ROOT"
            else
                tar -xzf "$tarball" -C "$EXTRACT_ROOT"
            fi
            ;;
        *.tar)
            tar -xf "$tarball" -C "$EXTRACT_ROOT"
            ;;
        *.tar.zst|*.tzst)
            if command -v zstd >/dev/null 2>&1; then
                zstd -dc "$tarball" | tar -xf - -C "$EXTRACT_ROOT"
            else
                echo "ERROR: $tarball requires zstd but 'zstd' not found" >&2
                exit 1
            fi
            ;;
        *)
            # Let tar try its built-in auto-detection as a last resort
            tar -xf "$tarball" -C "$EXTRACT_ROOT"
            ;;
    esac

    found="$(find "$EXTRACT_ROOT" -type d -path "*/nodes/ukamaOS" -print -quit || true)"
    if [[ -z "$found" ]]; then
        echo "ERROR: After extract, could not find */nodes/ukamaOS under $EXTRACT_ROOT" >&2
        echo "Tip: tar -tf \"$tarball\" | head" >&2
        exit 1
    fi

    UKAMA_OS_PATH="$found"
}

main() {
    detect_env
    echo "Build environment is $BUILD_ENV"

    # Allow explicit override for CI/debug:
    # UKAMA_OS=/some/path ./script.sh
    if [[ -n "${UKAMA_OS:-}" ]]; then
        UKAMA_OS_PATH="$UKAMA_OS"
    elif [[ "$BUILD_ENV" == "local" ]]; then
        UKAMA_OS_PATH="$(realpath ../../nodes/ukamaOS)"
    else
        extract_source
    fi

    if [[ ! -d "$UKAMA_OS_PATH" ]]; then
        echo "ERROR: UkamaOS not found at: $UKAMA_OS_PATH" >&2
        exit 1
    fi

    export UKAMA_OS="$UKAMA_OS_PATH"
    echo "UKAMA_OS=$UKAMA_OS"
}

main "$@"
