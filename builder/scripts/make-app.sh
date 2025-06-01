#!/bin/sh

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

# Script to build and package ukamaOS app

set -e

# Base parameters
UKAMA_OS=`realpath ../ukamaOS`
SYS_ROOT=${UKAMA_OS}/distro
SCRIPTS_ROOT=${SYS_ROOT}/scripts

# Build the app at given src path and cmd
build_app() {

    CWD=`pwd`
    SRC=$1
    CMD=$2

    cd ${SRC} && ${CMD} && cd ${CWD}
}

# copy all the required lib to rootfs
copy_all_libs() {

    BIN=$1
    DEST=$2

    for lib in $(ldd ${BIN} | cut -d '>' -f2 | awk '{print $1}')
    do
        if [ -f "${lib}" ]; then
            # Use case statement for substring match
            case "${lib}" in
                *libusys.so*)
                    # Copy libusys.so directly to /lib
                    cp "${lib}" "${DEST}/lib"
                    ;;
                *)
                    # Copy other libraries to the destination with their parents
                    cp --parents "${lib}" "${DEST}"
                    ;;
            esac
        fi
    done
}

# main
ACTION=$1
case "$ACTION" in
    "init")
        rm -rf $2
        mkdir $2
        mkdir $2/sbin
        mkdir $2/bin
        mkdir $2/lib
        mkdir $2/conf
        ;;
    "build")
	    build_app $3 "$4"
	    ;;
    "cp")
	    cp $2 $3
	    ;;
    "exec")
	    $2
	    ;;
    "patchelf")
#	    patchelf --set-rpath /lib $2
	    ;;
    "mkdir")
	    mkdir -p $2
	    ;;
    "libs")
	    copy_all_libs $2 $3
	    ;;
    "clean")
	    rm -rf $2
	    ;;
    "pack")
	    mkdir -p $2/build/pkgs/
	    tar -czf $2/build/pkgs/$3 $4
	    if [ $5 -eq 1 ]
	    then
	        rm -rf $4 ${ROOTFS}
	    fi
	    ;;
esac

exit 0
