#!/bin/bash
# Copyright (c) 2022-present, Ukama Inc.
# All rights reserved.

# compile the binaries

INIT_ROOT=`realpath ../../`

SUB_SYSTEMS="api-gateway lookup node-gateway"

#
# build subsystems
#
build_subsystems() {

	SS=$1
	CWD=`pwd`
	cd $INIT_ROOT/$SS
	make
	cp -rf ./bin/$SS $CWD/bin
	make clean
	cd $CWD
}

#
# print_usage
#
print_usage() {
	echo "usage: runSystem.sh help | build | clean"
}

# main

ACTION=$1
CWD=`pwd`

case "$ACTION" in
	"help")
		print_usage
		;;
	"build")
		mkdir -p ./bin
		for ss in `echo $SUB_SYSTEMS`
		do
			build_subsystems $ss
		done
		;;
	"clean")
		rm -rf ./bin
		;;
	*)
		print_usage
		;;
esac

exit 0
