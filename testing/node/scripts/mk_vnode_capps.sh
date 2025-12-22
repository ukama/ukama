#!/usr/bin/env bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2022-present, Ukama Inc.

set -euo pipefail

# Defaults (can be overridden by env vars)
: "${BUILD_ENV:=}"   # if empty, we'll auto-detect
: "${UKAMA_OS:=}"    # optional explicit override
: "${UKAMA_ROOT:=}"
UKAMA_ROOT_DEFAULT="/tmp/virtnode/ukama"
UKAMA_OS_PATH_DEFAULT="/tmp/virtnode/ukama/nodes/ukamaOS"

DEF_BUILD_DIR="./build/capps"
BUILD_DIR="${DEF_BUILD_DIR}"

# default target is local machine (gcc)
DEF_TARGET="local"
TARGET="${DEF_TARGET}"

detect_env() {
    # Respect explicit BUILD_ENV=local|container if set
    if [[ -n "${BUILD_ENV:-}" ]]; then
        if [[ "$BUILD_ENV" == "local" || "$BUILD_ENV" == "container" ]]; then
            return
        fi
        echo "ERROR: Invalid BUILD_ENV='$BUILD_ENV' (expected 'local' or 'container')" >&2
        exit 1
    fi

    BUILD_ENV="local"

    # Docker/Podman markers
    if [[ -f "/.dockerenv" || -f "/run/.containerenv" ]]; then
        BUILD_ENV="container"
        return
    fi

    # Sometimes present in container runtimes
    if [[ -n "${container:-}" || -n "${CONTAINER:-}" ]]; then
        BUILD_ENV="container"
        return
    fi

    # Fallback heuristic
    if grep -qaE '(docker|podman|containerd|kubepods|lxc)' /proc/1/cgroup 2>/dev/null; then
        BUILD_ENV="container"
        return
    fi
}

resolve_ukama_os() {
    detect_env
    echo "Build environment is: $BUILD_ENV"

    # Highest priority: explicit UKAMA_OS env var
    if [[ -n "${UKAMA_OS:-}" ]]; then
        :
    elif [[ "$BUILD_ENV" == "local" ]]; then
        UKAMA_OS="$(realpath ../../nodes/ukamaOS)"
        UKAMA_ROOT="$(realpath ../..)"
    else
        UKAMA_ROOT="$UKAMA_ROOT_DEFUALT"
        UKAMA_OS="$UKAMA_OS_PATH_DEFAULT"
    fi

    if [[ ! -d "$UKAMA_OS" ]]; then
        echo "ERROR: ukamaOS not found at: $UKAMA_OS (BUILD_ENV=$BUILD_ENV)" >&2
        exit 1
    fi

    # Optional extra validation (helps catch wrong mounts)
    if [[ ! -d "$UKAMA_OS/distro" ]]; then
        echo "ERROR: Expected '$UKAMA_OS/distro' to exist; looks like wrong ukamaOS root: $UKAMA_OS" >&2
        exit 1
    fi

    export UKAMA_ROOT
    export UKAMA_OS
    echo "UKAMA_ROOT=$UKAMA_ROOT"
    echo "UKAMA_OS=$UKAMA_OS"
}

build_app() {
    local cwd src cmd
    cwd="$(pwd)"
    src="${UKAMA_OS}$1"
    cmd="$2"

    cd "$src"
    eval "$cmd"
    echo "CApp build done for: ${cmd} (src=$src)"

    cd "$cwd"
}

copy_all_libs() {
    local bin capp lib
    bin="$1"
    capp="$2"

    mkdir -p "${BUILD_DIR}/${capp}/lib"

    # ldd output parsing is imperfect but OK for your use-case.
    # Read libs safely line-by-line to reduce word-splitting issues.
    while IFS= read -r lib; do
        [[ -z "$lib" ]] && continue
        if [[ -f "$lib" ]]; then
            cp --parents "$lib" "$BUILD_DIR"
            cp "$lib" "${BUILD_DIR}/${capp}/lib"
        fi
    done < <(ldd "$bin" | awk '/=>/ {print $3} /^[[:space:]]*\// {print $1}' | sort -u)
}

main() {
    # Action can be 'build', 'cp', 'mkdir', etc.
    local action="${1:-}"
    if [[ -z "$action" ]]; then
        echo "ERROR: missing ACTION. Usage: $0 <build|cp|exec|patchelf|mkdir|libs|rename|clean> ..." >&2
        exit 1
    fi

    resolve_ukama_os

    local SYS_ROOT="${UKAMA_OS}/distro"
    local SCRIPTS_ROOT="${SYS_ROOT}/scripts/"  # kept in case used later

    case "$action" in
        build)
            if [[ "${2:-}" == "app" ]]; then
                build_app "${3:-}" "${4:-}"
            fi
            ;;
        cp-config)
            mkdir -p "${3}"
            cp "${UKAMA_ROOT}/${2:?missing src path}" "${3:?missing dest path}"
            ;;
        cp)
            cp "${UKAMA_OS}/${2:?missing src path}" "${BUILD_DIR}/${3:?missing dest path}"
            ;;
        exec)
            "${2:?missing command}"
            ;;
        patchelf)
            patchelf --set-rpath /usys/lib "${UKAMA_OS}/${2:?missing path}"
            ;;
        mkdir)
            mkdir -p "${BUILD_DIR}/${2:?missing dir}"
            ;;
        libs)
            copy_all_libs "${UKAMA_OS}/${2:?missing bin}" "${3:?missing capp}"
            ;;
        rename)
            mv "${BUILD_DIR}" "${2:?missing new name}"
            ;;
        clean)
            # your original: if [ "$2" = "" ] then ... else ...
            if [[ -z "${2:-}" ]]; then
                rm -rf "${BUILD_DIR}"
            else
                rm -rf "${BUILD_DIR}/${2}"
            fi
            ;;
        *)
            echo "ERROR: Unknown ACTION '$action'" >&2
            exit 1
            ;;
    esac
}

main "$@"
