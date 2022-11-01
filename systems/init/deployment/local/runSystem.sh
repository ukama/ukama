#!/bin/bash
# Copyright (c) 2022-present, Ukama Inc.
# All rights reserved.

# compile the binaries

INIT_ROOT=`realpath ../../`

SUB_SYSTEMS="api-gateway lookup node-gateway"

#
# Docker stuff:
# docker stop init-postgres
# docker system prune

#
# check_openport
#
check_openport() {

	PORT=$1

	if ( nc -zv localhost $PORT 2>&1 >/dev/null ); then
		echo "Port in use: $PORT"
		exit 0
	fi
}

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
# clean all the residual of init system
#
clean_subsystems() {

	sudo docker stop init-postgres
	sudo docker rm init-postgres
	kill -9 `ps -ewf | grep "./bin/lookup" | awk '{print $2}' | head -1`
	kill -9 `ps -ewf | grep "./bin/api-gateway" | awk '{print $2}' | head -1`
	kill -9 `ps -ewf | grep "./bin/node-gateway" | awk '{print $2}' | head -1`

}

#
# run all subsystems
#
exec_subsystems() {

	# Steps are:
	# 1. Run postgres
	# 2. Run lookup sub-system
	# 3. Run api-gateway
	# 4. Run node-gateway

	sudo docker run --name init-postgres -p 5432:5432 \
		 -e POSTGRES_PASSWORD=Pass2020! -d postgres

	# Lookup
	export DB_HOST="localhost"
	export DB_USERNAME="postgres"
	export MSGCLIENT_HOST="msgclient:9091"
	./bin/lookup &

	# API gateway
	export SERVER_PORT=8081
	export GRPC_PORT=9090
	export METRICS_PORT=10240
	export METRICS_ENABLED="false"
	./bin/api-gateway &

	# Node gateway
	export SERVER_PORT=8082
	export METRICS_PORT=10241
	export METRICS_ENABLED="false"
	export SERVICES_LOOKUP="lookup:9090"
	./bin/node-gateway &
}

#
# print_usage
#
print_usage() {
	echo "usage: runSystem.sh help | build | exec | clean"
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
	"exec")
		if [ ! -d "./bin" ]; then
			echo "Binaries not found, build first"
			exit 0
		fi

		# Check if the postgres port is in use.
		check_openport 5432
		echo "Runing init system"
		echo "sudo is to run postgres docker"
		exec_subsystems
		;;
	"clean")
		echo "Cleaning binary dir: ./bin"
		rm -rf ./bin
		echo "Stop existing exec"
		clean_subsystems
		;;
	*)
		print_usage
		;;
esac

exit 0
