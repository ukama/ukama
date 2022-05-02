#!/bin/sh
# Copyright (c) 2022-present, Ukama Inc.
# All rights reserved.

# Script to generate capps for the virtual node

# Base parameters
UKAMA_OS=`realpath ../../../nodes/ukamaOS`
SYS_ROOT=${UKAMA_OS}/distro/
SCRIPTS_ROOT=${SYS_ROOT}/scripts/
BB_ROOT=${UKAMA_OS}/distro/system/busybox
BB_CONFIG=ukama_minimal_defconfig

DEF_BUILD_DIR=./build/

#Various network related parameters
HOSTNAME="localhost"

# default target is local machine (gcc)
DEF_TARGET="local"
TARGET=${DEF_TARGET}

# default rootfs location is ${DEF_BUILD_DIR}
BUILD_DIR=${DEF_BUILD_DIR}

#
# Build the app at given src path and cmd
#
build_app() {

    CWD=`pwd`
    SRC=$1
    CMD=$2

    cd ${SRC} && ${CMD} && cd ${CWD}
}

# main

mkdir -p ${BUILD_DIR}

# Action can be 'build', 'cp' and 'mkdir'
ACTION=$1

case "$ACTION" in
    "build")
	if [ "$2" = "app" ]
	then
	    build_app $3 "$4"
	elif [ "$2" = "busybox" ]
	then
	     build_busybox
	     build_rootfs_dirs
	     setup_etc
	fi
	;;
    "cp")
	cp $2 ${BUILD_DIR}/$3
	;;
    "exec")
	$2
	;;
    "patchelf")
	patchelf --set-rpath /lib $2
	;;
    "mkdir")
	mkdir ${BUILD_DIR}/$2
	;;
    "libs")
	copy_all_libs $2
	;;
    "rename")
	mv ${BUILD_DIR} $2
	;;
    "clean")
	if [ "$2" = "" ]
	then
	    rm -rf ${BUILD_DIR}
	else
	    rm -rf $2
	fi
esac

exit
