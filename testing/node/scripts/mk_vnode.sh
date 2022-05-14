#!/bin/sh
# Copyright (c) 2022-present, Ukama Inc.
# All rights reserved.

# Script to create ukama's virtual node.

# Base parameters
UKAMA_OS=`realpath ../../nodes/ukamaOS`
NODED_ROOT=${UKAMA_OS}/distro/system/noded/
DEF_BUILD_DIR=./build/

# default target is local machine (gcc)
DEF_TARGET="local"
TARGET=${DEF_TARGET}

# default rootfs location is ${DEF_BUILD_DIR}
BUILD_DIR=`realpath ${DEF_BUILD_DIR}`

#
# Build needed tools, e.g., genSchema, genInventory, if needed.
#
build_utils() {

	CWD=`pwd`

	mkdir -p ${BUILD_DIR}/utils

	# Build genSchema
    cd ${NODED_ROOT} && make genSchema
	if [ -f ${NODED_ROOT}/build/genSchema ]; then
		cp ${NODED_ROOT}/build/genSchema ${BUILD_DIR}/utils/
	else
		echo "Error building genSchema. Exiting"
		exit 1
	fi

	# Build genInventory - to create the EEPROM data
	cd ${NODED_ROOT} && make genInventory
	if [ -f ${NODED_ROOT}/build/genInventory ]; then
		cp ${NODED_ROOT}/build/genInventory ${BUILD_DIR}/utils/
	else
		echo "Error building genSchema. Exiting"
		exit 1
	fi

	cd $CWD
}

#
# Build /sys for the virtual node
# 1. prepare_env.sh
# 2. genSchema
# 3. genInventory

build_sysfs() {

	CWD=`pwd`
	NODE_TPYE=$1
	NODE_UUID=$2
	MODULES_METADATA=$3

	${NODED_ROOT}/utils/prepare_env.sh clean
	${NODED_ROOT}/utils/prepare_env.sh --unittype $1

	# genSchema --u UK-7001-HNODE-SA03-1102 \
	# --n ComV1 --m UK-7001-COM-1102 --f mfgdata/schema/com.json \
	# --n LTE   --m UK-7001-TRX-1102 --f mfgdata/schema/lte.json \
	# --n MASK  --m UK-7001-MSK-1102 --f mfgdata/schema/mask.json

	# copy the schemas file locally and run genSchema
	mkdir -p ${BUILD_DIR}/schemas
	cp ${NODED_ROOT}/mfgdata/schema/*.json  ${BUILD_DIR}/schemas

	${BUILD_DIR}/utils/genSchema -u $NODE_UUID $MODULE_METADATA

	# create EEPROM data using genInventory
	${BUILD_DIR}/utils/genInventory $MODULE_METADATA

	#copy the sysfs to build dir
	cp -rf /tmp/sys ${BUILD_DIR}/sys
	rm -rf /tmp/sys
}

#
# Build image using buildah
#

build_image() {

	FILE=$1
	NAME_TAG=$2

	buildah bud -f $1 -t $2
}

# main

ACTION=$1
CWD=`pwd`

case "$ACTION" in
	"init")
		build_utils
		;;
	"sysfs")
		build_sysfs $2 $3 $4
		;;
	"build")
		build_image $1 $2
		;;
    "cp")
		cp $2 ${BUILD_DIR}/$3
		;;
    "clean")
		buildah rmi -f localhost/$1
		rm -rf ${BUILD_DIR}
		cd ${NODED_ROOT} && make clean && cd ${CWD}
		;;
esac

exit
