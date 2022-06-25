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

REGISTRY_URL=

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
#
build_sysfs() {

	CWD=`pwd`
	NODE_TPYE=$1
	NODE_UUID=$2

	${NODED_ROOT}/utils/prepare_env.sh --clean
	${NODED_ROOT}/utils/prepare_env.sh --unittype $1

	# genSchema --u UK-SA7001-HNODE-m0-1102 \
	# --n LTE   --m UK-SA7001-TRX-m0-1102 --f mfgdata/schema/hnode_trx.json

	# genSchema --u UK-SA7001-TNODE-m0-1102 \
	# --n ComV1 --m UK-SA7001-COM-m0-1102 --f mfgdata/schema/com.json \
	# --n LTE   --m UK-SA7001-TRX-m0-1102 --f mfgdata/schema/trx.json \
	# --n MASK  --m UK-SA7001-MSK-m0-1102 --f mfgdata/schema/mask.json

	# copy the mfgdata locally and run genSchema/genInventory
	mkdir -p ${BUILD_DIR}/schemas
	cp ${NODED_ROOT}/mfgdata/schema/*.json  ${BUILD_DIR}/schemas
	cp -rf ${NODED_ROOT}/mfgdata ${BUILD_DIR}

	cd ${BUILD_DIR}
	${BUILD_DIR}/utils/genSchema -u $NODE_UUID $VNODE_SCHEMA_ARGS

	# create EEPROM data using genInventory
	${BUILD_DIR}/utils/genInventory $VNODE_SCHEMA_ARGS

	#copy the sysfs to build dir
	cp -rf /tmp/sys ${BUILD_DIR}/sys
	rm -rf /tmp/sys
	cd ${CWD}
}

#
# Build image using buildah
#
build_image() {

	FILE=$1
	UUID=$2

	NAME_TAG=`echo ${UUID} | awk '{print tolower($0)}'`

	# copy capp's sbin, conf and lib to /sbin, /conf and /lib
	mkdir -p ${BUILD_DIR}/sbin ${BUILD_DIR}/lib ${BUILD_DIR}/conf
	mkdir -p ${BUILD_DIR}/tmp ${BUILD_DIR}/bin

	cp -rf ${BUILD_DIR}/capps/*/sbin ${BUILD_DIR}
	cp -rf ${BUILD_DIR}/capps/*/conf ${BUILD_DIR}
	cp -rf ${BUILD_DIR}/capps/*/lib  ${BUILD_DIR}

	cp ./scripts/runme.sh   ${BUILD_DIR}/bin/
	cp ./scripts/waitfor.sh ${BUILD_DIR}/bin/
	cp ./scripts/kickstart.sh ${BUILD_DIR}/bin/

	buildah bud -f $1 -t $NAME_TAG ${REGISTRY_URL}
}

#
# push image to repo
#
push_image() {

	UUID=$1
	TAG="latest"
	NAME=`echo ${UUID} | awk '{print tolower($0)}'`

	if [ ${DOCKER_USER} != "aws" ]; then
		REGISTRY_URL=${DOCKER_USER}
		buildah push --tls-verify=false --creds ${DOCKER_USER}:${DOCKER_PASS} \
				 ${NAME}:${TAG}
	else
		REGISTRY_URL=${REPO_SERVER_URL}
		echo "Docker login to AWS"
		mkdir -p ~/.docker/
		echo '{ "credsStore": "ecr-login" }' > ~/.docker/config.json
		buildah push --format docker://${REGISTRY_URL}/repo/${NAME}:${TAG}
	fi
}

# main

ACTION=$1
CWD=`pwd`

case "$ACTION" in
	"init")
		build_utils
		;;
	"sysfs")
		build_sysfs $2 $3
		;;
	"build")
		build_image $2 $3
		;;
	"push")
		push_image $2
		;;
	"cp")
		cp $2 ${BUILD_DIR}/$3
		;;
	"clean")
		rm ContainerFile; rm supervisor.conf
		buildah rmi -f localhost/$1
		cd ${NODED_ROOT} && make clean && cd ${CWD}
		;;
esac

exit
