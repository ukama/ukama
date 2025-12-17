#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2022-present, Ukama Inc.

#USAGE: prepare_env.sh

# Preapre the envirnment for building virtual node.

BUILD_ENV=container

UKAMA_OS_TAR=/ukama/ukama_*.tgz
UKAMA_OS_PATH=/tmp/virtnode/nodes/ukamaOS

# Check if building on local or in container
if_host() {
    val=`cat /proc/1/cgroup | grep -i "pids" |  awk -F":" 'NR==1{print $NF}'`
    if [ "${val}" == "/init.scope" ] || [ "${val}" == "/" ]; then
        BUILD_ENV=local
    fi
}

extract_source() {
	tar -zxvf $UKAMA_OS_TAR -C /
	if [ $? == 0 ]; then
		echo "Extraction for ukama source is success."
	else
		exit 1
	fi
}

set_ukama_os_env() {
	UKAMA_OS=$1
	export UKAMA_OS
}

# main

if_host

echo "Build envionment is $BUILD_ENV"

if [ $BUILD_ENV == "local" ]; then
	UKAMA_OS_PATH=`realpath ../../nodes/ukamaOS`
elif [ $BUILD_ENV == "container" ]; then
	extract_source
else
	echo "Unkown enviornment."
	exit 1
fi

if [ -d $UKAMA_OS_PATH ]; then
    echo "Build environment is set for the Virtual Node on $BUILD_ENV."
    exit 0;
else
    echo "UkamaOS not found."
    exit 1;
fi
