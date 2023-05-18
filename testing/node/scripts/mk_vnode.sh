#!/bin/bash
# Copyright (c) 2022-present, Ukama Inc.
# All rights reserved.

# Script to create ukama's virtual node.

# Base parameters
#UKAMA_OS=`realpath ../../nodes/ukamaOS`
NODED_ROOT=
DEF_BUILD_DIR=./build/
BUILD_ENV=container

# default target is local machine (gcc)
DEF_TARGET="local"
TARGET=${DEF_TARGET}

# default rootfs location is ${DEF_BUILD_DIR}
BUILD_DIR=`realpath ${DEF_BUILD_DIR}`

REGISTRY_URL=${REPO_SERVER_URL}

REGISTRY_NAME=${REPO_NAME}

#
# Check if building on local or in container
#
if_host() {
    val=`cat /proc/1/cgroup | grep -i "pids" |  awk -F":" 'NR==1{print $NF}'`
    if [ "${val}" == "/init.scope" ] || [ "${val}" == "/" ]; then
        BUILD_ENV=local
    fi
}

#
# Update UKAMA_OS
#

update_ukama_os_env() {
	if_host

	if [ "$BUILD_ENV" == "local" ]; then
		UKAMA_OS=`realpath ../../nodes/ukamaOS`
	elif [ "$BUILD_ENV" == "container" ]; then
		UKAMA_OS="/tmp/virtnode/ukamaOS"
		if [ -z $UAKMA_OS ]; then
			echo "UKAMA OS env set to $UKAMA_OS"
		else
			echo "Failed to find ukamaOS at $UAKMA_OS"
			exit 1
		fi
	else
		echo "Unkown enviornment."
		exit 1
	fi

	NODED_ROOT=${UKAMA_OS}/distro/system/noded

}

#
# Build needed tools, e.g., genSchema, genInventory, if needed.
#
build_utils() {

	CWD=`pwd`

	mkdir -p ${BUILD_DIR}/utils

    update_ukama_os_env

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

	update_ukama_os_env

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
	if [ $? != 0 ]; then
        echo "Failed to create schema for $NODE_UUID $VNODE_SCHEMA_ARGS."
        exit 1
	fi

	# create EEPROM data using genInventory
	${BUILD_DIR}/utils/genInventory $VNODE_SCHEMA_ARGS
	if [ $? != 0 ]; then
        echo "Failed to create inventory DB $VNODE_SCHEMA_ARGS."
        exit 1
	fi

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

	buildah bud -f $1 -t ${REGISTRY_URL}/${REGISTRY_NAME}:${NAME_TAG} .
	if [ $? == 0 ]; then
        echo "Buildah created image ${REGISTRY_URL}/${REGISTRY_NAME}:${NAME_TAG}"
	else
        echo "Buildah image creation failed"
        exit 1
	fi
}

#
# push image to repo
#
push_image() {

	UUID=$1
	TARGET=$2
	TAG=`echo ${UUID} | awk '{print tolower($0)}'`

	if [ ${TARGET} != "REMOTE" ]; then
		buildah push --tls-verify=false \
				${REGISTRY_URL}/${REGISTRY_NAME}:${TAG} \
				localhost:5000/${REGISTRY_URL}/${REGISTRY_NAME}:${TAG}
		echo "Image ${REGISTRY_URL}/${REGISTRY_NAME}:${TAG} pushed to ${TARGET}"
		return
	fi

	pass=`aws ecr get-login-password`

	buildah login --username ${DOCKER_USER} --password ${pass} ${REGISTRY_URL}
	if [ $? == 0 ]; then
		echo "Registry login success."
	else
		echo "Registry login failure."
		exit 1
	fi

	if [ ${DOCKER_USER} != "AWS" ]; then

		buildah push --tls-verify=false --creds ${DOCKER_USER}:${DOCKER_PASS} \
				 ${REGISTRY_URL}/${REGISTRY_NAME}:${TAG}
	else

		buildah push ${REGISTRY_URL}/${REGISTRY_NAME}:${TAG}
		if [ $? == 0 ]; then
			echo "Image ${REGISTRY_URL}/${REGISTRY_NAME}:${TAG} pushed to registry."
		else
			echo "Failure to push image ${REGISTRY_URL}/${REGISTRY_NAME}:${TAG} to registry"
			exit 1
		fi
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
		push_image $2 $3
		;;
	"cp")
		cp $2 ${BUILD_DIR}/$3
		;;
	"clean")
		update_ukama_os_env
		rm ContainerFile; rm supervisor.conf
		buildah rmi -f localhost/$1
		cd ${NODED_ROOT} && make clean && cd ${CWD}
		;;
esac

exit 0
