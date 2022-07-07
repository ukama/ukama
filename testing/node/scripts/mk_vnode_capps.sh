#!/bin/bash
# Copyright (c) 2022-present, Ukama Inc.
# All rights reserved.

# Script to generate capps for the virtual node

# Base parameters
#UKAMA_OS=`realpath ../../nodes/ukamaOS`
BUILD_ENV=container
UKAMA_OS_PATH=/tmp/virtnode/ukamaOS

DEF_BUILD_DIR=./build/capps

#Various network related parameters
HOSTNAME="localhost"

# default target is local machine (gcc)
DEF_TARGET="local"
TARGET=${DEF_TARGET}
BUILD_DIR=${DEF_BUILD_DIR}

# Check if building on local or in container
if_host() {
    val=`cat /proc/1/cgroup | grep -i "pids" |  awk -F":" 'NR==1{print $NF}'`
    if [ "${val}" == "/init.scope" ] || [ "${val}" == "/" ]; then
        BUILD_ENV=local
    fi
}

#
# Build the app at given src path and cmd
#
build_app() {

    CWD=`pwd`
    SRC=${UKAMA_OS}$1
    CMD=$2

    cd ${SRC} && ${CMD}
    if [ $? == 0 ]; then
       echo "CApp build done for ${CMD} ${SRC}"
    else
        echo "CApp build failed for ${CMD} ${SRC}"
        exit 1
    fi

	cd ${CWD}
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

if_host

echo "Build envionment is $BUILD_ENV"

if [ $BUILD_ENV == "local" ]; then
    UKAMA_OS=`realpath ../../nodes/ukamaOS`
elif [ $BUILD_ENV == "container" ]; then
    UKAMA_OS=$UKAMA_OS_PATH
else
   echo "Unkown enviornment."
   exit 1
fi

SYS_ROOT=${UKAMA_OS}/distro
SCRIPTS_ROOT=${SYS_ROOT}/scripts/

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

exit 0
