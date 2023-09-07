#!/bin/bash
# Copyright (c) 2022-present, Ukama Inc.
# All rights reserved.

#USAGE: waitfor.sh program-1 ... program-n 

PID_DIR=/tmp
MIN_ARG=3
ARG_ARRAY=( "$@" )
ARG_LEN=$#

if (( $# < $MIN_ARG )); then
	echo "Need atleast three arg: waitfor.sh program_name_towait exec_with_arg"
	exit 0
fi

for (( i=0; i<$ARG_LEN; i++));
do
	# wait for the program to exit.
	# it will return if the program is already done
	tail --pid=`cat ${PID_DIR}/${ARG_ARRAY[$i]}.pid` -f /dev/null
done

exit 0
