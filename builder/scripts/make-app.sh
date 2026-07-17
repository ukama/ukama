#!/bin/sh

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

set -eu

usage() {
    echo "Usage: $0 ACTION [ARG ...]" >&2
    exit 1
}

build_app() {
    src="$1"
    cmd="$2"
    old_pwd="$(pwd)"

    cd "$src"

    if [ -f Makefile ] || [ -f makefile ]; then
        BUILD_MODE="${BUILD_MODE:-debug}" make clean
    fi

    BUILD_MODE="${BUILD_MODE:-debug}" sh -c "$cmd"
    cd "$old_pwd"
}

copy_all_libs() {
    bin="$1"
    dest="$2"

    ldd "$bin" | awk '
        /=>/ && $3 ~ /^\// { print $3 }
        $1 ~ /^\// { print $1 }
    ' | while IFS= read -r lib; do
        [ -f "$lib" ] || continue

        case "$lib" in
            *libusys.so*)
                cp "$lib" "$dest/lib/"
                ;;
            *)
                cp --parents "$lib" "$dest"
                ;;
        esac
    done
}

[ "$#" -ge 1 ] || usage
action="$1"
shift

case "$action" in
    init)
        [ "$#" -eq 1 ] || usage
        rm -rf -- "$1"
        mkdir -p -- "$1/sbin" "$1/bin" "$1/lib" "$1/conf"
        ;;

    build)
        [ "$#" -eq 3 ] || usage
        build_app "$2" "$3"
        ;;

    cp)
        [ "$#" -eq 2 ] || usage
        cp -- "$1" "$2"
        ;;

    exec)
        [ "$#" -eq 1 ] || usage
        sh -c "$1"
        ;;

    patchelf)
        [ "$#" -eq 1 ] || usage
        patchelf --set-rpath /lib "$1"
        ;;

    mkdir)
        [ "$#" -eq 1 ] || usage
        mkdir -p -- "$1"
        ;;

    libs)
        [ "$#" -eq 2 ] || usage
        copy_all_libs "$1" "$2"
        ;;

    clean)
        [ "$#" -eq 1 ] || usage
        rm -rf -- "$1"
        ;;

    pack)
        [ "$#" -eq 4 ] || usage
        repo_root="$1"
        archive="$2"
        source_dir="$3"
        remove_source="$4"
        pkg_dir="${repo_root}/build/pkgs"
        latest_link="$(
            printf '%s\n' "$archive" |
                sed -E 's/_[^/]+\.tar\.gz$/_latest.tar.gz/'
        )"

        mkdir -p -- "$pkg_dir"
        tar -czf "${pkg_dir}/${archive}" -- "$source_dir"
        if [ "${UKAMA_SKIP_LATEST_LINK:-0}" -ne 1 ]; then
            rm -f -- "${pkg_dir}/${latest_link}"
            ln -s -- "$archive" "${pkg_dir}/${latest_link}"
        fi

        if [ "$remove_source" -eq 1 ]; then
            rm -rf -- "$source_dir"
        fi
        ;;

    *)
        usage
        ;;
esac
