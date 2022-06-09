#!/bin/sh
# Copyright (c) 2022-present, Ukama Inc.
# All rights reserved.

# Script to setup and test bootstrap

# Base parameters
UKAMA_OS=`realpath ../../../../`
NODED_ROOT=${UKAMA_OS}/distro/system/noded/
BOOTSTRAP_ROOT=${UKAMA_OS}/distro/system/bootstrap/
DEF_BUILD_DIR=./run/

# default target is local machine (gcc)
DEF_TARGET="local"
TARGET=${DEF_TARGET}

# default rootfs location is ${DEF_BUILD_DIR}
BUILD_DIR=`realpath ${DEF_BUILD_DIR}`

NODE_UUID="ukma-sa3333-tnode-m0-e465"

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
# Build /sys for the node.d
# 1. prepare_env.sh
# 2. genSchema
# 3. genInventory
#
build_sysfs() {

	CWD=`pwd`

	${NODED_ROOT}/utils/prepare_env.sh --clean
	${NODED_ROOT}/utils/prepare_env.sh --unittype "tnode"

	# genSchema --u UK-7001-HNODE-SA03-1102 \
	# --n ComV1 --m UK-7001-COM-1102 --f mfgdata/schema/com.json \
	# --n LTE   --m UK-7001-TRX-1102 --f mfgdata/schema/lte.json \
	# --n MASK  --m UK-7001-MSK-1102 --f mfgdata/schema/mask.json

	# copy the mfgdata locally and run genSchema/genInventory
	mkdir -p ${BUILD_DIR}/schema
	cp ${NODED_ROOT}/mfgdata/schema/*.json  ${BUILD_DIR}/schema
	cp -rf ${NODED_ROOT}/mfgdata ${BUILD_DIR}

	cd ${BUILD_DIR}
	${BUILD_DIR}/utils/genSchema -u ${NODE_UUID} \
			--n com  --m ukma-7001-com-1102 --f ${BUILD_DIR}/schema/com.json \
			--n trx  --m ukma-7001-trx-1102 --f ${BUILD_DIR}/schema/trx.json \
			--n mask --m ukma-7001-mask-1102 --f ${BUILD_DIR}/schema/mask.json

	# create EEPROM data using genInventory
	${BUILD_DIR}/utils/genInventory \
			--n com  --m ukma-7001-com-1102 --f ${BUILD_DIR}/schema/com.json \
			--n trx  --m ukma-7001-trx-1102 --f ${BUILD_DIR}/schema/trx.json \
			--n mask --m ukma-7001-mask-1102 --f ${BUILD_DIR}/schema/mask.json

	# copy the sysfs to build dir
	cp -rf /tmp/sys ${BUILD_DIR}/sys

	rm ${BUILD_DIR}/sys/tnode_inventory_db
	ln -s ${BUILD_DIR}/sys/bus/i2c/devices/i2c-0/0-0050/eeprom \
	   ${BUILD_DIR}/sys/tnode_inventory_db

	cd ${CWD}
}

#
# Build noded and bootstrap
#
build_bins() {

	mkdir -p ${BUILD_DIR}/bin

	CWD=`pwd`
	cd ${BOOTSTRAP_ROOT} && make clean && make
	cp ${BOOTSTRAP_ROOT}/bootstrap ${BUILD_DIR}/bin
	cd ${CWD}

	cd ${BOOTSTRAP_ROOT}/test && make clean && make
	cp ${BOOTSTRAP_ROOT}/test/build/bootstrap_server ${BUILD_DIR}/bin
	cd ${CWD}

	cd ${NODED_ROOT} && make clean && make
	cp ${NODED_ROOT}/build/noded ${BUILD_DIR}/bin
	cd ${CWD}
}

#
# create config.toml for mesh.d
#
create_mesh_config() {

	MESH_CONFIG=${BUILD_DIR}/config/mesh_config.toml
	rm -f ${MESH_CONFIG}

	echo "[client-config]"                                    >> ${MESH_CONFIG}
	echo " remote-ip-file = \"${BUILD_DIR}/config/ip_file\" " >> ${MESH_CONFIG}
	echo " cert           = \"${BUILD_DIR}/cert/test.crt\" "  >> ${MESH_CONFIG}
	echo " key            = \"${BUILD_DIR}/cert/test.crt\" "  >> ${MESH_CONFIG}
	
}

#
# create config.toml for bootstrap
#
create_bootstrap_config() {

	BOOT_CONFIG=${BUILD_DIR}/config/bootstrap_config.toml
	rm -f ${BOOT_CONFIG}

	echo "[config]"                                 >> ${BOOT_CONFIG}
	echo " noded-host = \"localhost\" "             >> ${BOOT_CONFIG}
	echo " noded-port = \"8095\" "                  >> ${BOOT_CONFIG}
	echo " remote-ip-file = \"file\" "              >> ${BOOT_CONFIG}
	echo " bootstrap-server  = \"localhost:4444\" " >> ${BOOT_CONFIG}
	echo " mesh-config = \"${BUILD_DIR}/config/mesh_config.toml\" " \
		 >> ${BOOT_CONFIG}
}

#
# Copy config etc.
#
setup_build_dir() {

	mkdir -p ${BUILD_DIR}/config
	mkdir -p ${BUILD_DIR}/cert
	mkdir -p ${BUILD_DIR}/property

	cp ${BOOTSTRAP_ROOT}/test/test.crt              ${BUILD_DIR}/cert
	cp ${NODED_ROOT}/mfgdata/property/property.json ${BUILD_DIR}/property

	create_mesh_config
	create_bootstrap_config
}

# main

ACTION=$1
CWD=`pwd`

echo "Target directory is: ${BUILD_DIR}"

case "$ACTION" in
	"help")
		echo "options are: help clean setup run"
		;;
	"clean")
		rm -rf ${BUILD_DIR}
		rm -rf /tmp/sys
		
		echo "${BUILD_DIR} deleted"
		echo "/tmp/sys deleted"
		;;

	"setup")
		mkdir -p ${BUILD_DIR}

		setup_build_dir
		build_bins
		build_utils
		build_sysfs

		echo "Done setup. Exiting."
		;;

	"run")
		echo "Build directory is: ${BUILD_DIR}"
		echo "Run following commands; each in different terminal:"
		echo "${BUILD_DIR}/bin/noded --p ${BUILD_DIR}/property/property.json --i /tmp/sys/tnode_inventory_db"
		echo "${BUILD_DIR}/bin/bootstrap_server 4444 192.168.0.1 ${BUILD_DIR}/cert/test.crt testOrg"
		echo "${BUILD_DIR}/bin/bootstrap --config ${BUILD_DIR}/config/bootstrap_config.toml"
		;;
esac

exit
