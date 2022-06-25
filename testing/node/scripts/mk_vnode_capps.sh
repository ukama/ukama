#!/bin/sh
# Copyright (c) 2022-present, Ukama Inc.
# All rights reserved.

# Script to generate capps for the virtual node

# Base parameters
UKAMA_OS=`realpath ../../nodes/ukamaOS`
SYS_ROOT=${UKAMA_OS}/distro
SCRIPTS_ROOT=${SYS_ROOT}/scripts/
DEF_BUILD_DIR=./build/capps

#Various network related parameters
HOSTNAME="localhost"

# default target is local machine (gcc)
DEF_TARGET="local"
TARGET=${DEF_TARGET}
BUILD_DIR=${DEF_BUILD_DIR}

#
# Build the app at given src path and cmd
#
build_app() {

    CWD=`pwd`
    SRC=${UKAMA_OS}$1
    CMD=$2

    cd ${SRC} && ${CMD} && cd ${CWD}
}

#
# copy all the required lib to rootfs
#
copy_all_libs() {

    BIN=$1
	CAPP=$2

	mkdir -p ${BUILD_DIR}/$2/lib

    for lib in $(ldd ${BIN} | cut -d '>' -f2 | awk '{print $1}')
    do
        if [ -f "${lib}" ]; then
            cp --parents "${lib}" ${BUILD_DIR}
            cp "${lib}" ${BUILD_DIR}/$2/lib
        fi
    done
}

# main

# Action can be 'build', 'cp' and 'mkdir'
ACTION=$1

case "$ACTION" in
    "build")
	if [ "$2" = "app" ]
	then
	    build_app $3 "$4"
	fi
	;;
    "cp")
	cp ${UKAMA_OS}/$2 ${BUILD_DIR}/$3
	;;
    "exec")
	$2
	;;
    "patchelf")
	patchelf --set-rpath /usys/lib ${UKAMA_OS}/$2
	;;
    "mkdir")
	mkdir -p ${BUILD_DIR}/$2
	;;
    "libs")
	copy_all_libs ${UKAMA_OS}/$2 $3
	;;
    "rename")
	mv ${BUILD_DIR} $2
	;;
    "clean")
	if [ "$2" = "" ]
	then
	    rm -rf ${BUILD_DIR}
	else
	    rm -rf ${BUILD_DIR}/$2
	fi
esac

exit
