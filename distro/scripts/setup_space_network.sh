#!/bin/sh
# Copyright (c) 2021-present, Ukama Inc.
# All rights reserved.

# Script to setup networking for cSpaces and cApps

# base parameters
UKAMA_OS=../../
DEF_BRIDGE="ukama_bridge"
MIN_ARGS=4

CMD="NaN"
DEV="NaN"
PID="NaN"
IFACE="NaN"
SPACE="NaN"
BRIDGE="NaN"

usage() {
    echo 'Usage: setup_space_network.sh interface space_name'
    exit
}

msg_usage() {
    echo "Usage:"
    echo "      setup_space_network.sh --add [bridge | cspace] [args]"
    echo ""
    echo "Options:"
    echo "   Dev=bridge:          # args when dev is 'bridge'"
    echo "       interface        # Interface to use to connect with Internet"
    echo "       bridge_name      # name to the master bridge"
    echo "   Dev=cspace:          # args when dev is 'cspace'"
    echo "       pid              # PID of the cspace process"
    echo "       space_name       # name of space"
    echo "       bridge_name      # name to valid master bridge"
    exit 100
}

valid_args() {

    # --add bridge iface bridge_name
    if [ "${DEV}" = "cspace" ]; then
	# Test valid PID
	ps -p ${PID} > /dev/null
	if [ "$?" -eq "1" ]
	then
	    exit 101
	fi
    fi

    # --add cspace pid space_name bridge_name
    if [ "${DEV}" = "bridge" ]; then
	# Test valid interface
	ip -o a show | cut -d ' ' -f 2,7 | \
	    awk '{print $1}' | grep ${IFACE} > /dev/null
	if [ "$?" -eq "1" ]
	then
	    exit 102
	fi
    fi
}

add_bridge() {

    IF=$1
    BR=$2

    # Setup bridge
    ip link add ${BR} type bridge

    # Setup host machine to allow NATing.
    iptables -t nat -A POSTROUTING -o ${BR} -j MASQUERADE
    iptables -t nat -A POSTROUTING -o ${IF} -j MASQUERADE
}

add_cspace() {

    ID=$1
    SP=$2
    BR=$3

    NS=${SP}

    # Setup paired veth for each cspace on the host
    ip link   add dev veth1_${SP} type veth peer name veth2_${SP}

    # Bring up the host iface
    ip link   set dev veth1_${SP} up
    ip tuntap add tap_${SP} mode tap
    ip link   set dev tap_${SP} up

    # Attach iface to the bridge
    ip link set tap_${SP}   master ${BR}
    ip link set veth1_${SP} master ${BR}

    # Give address to the bridge
    ip addr add 10.0.0.1/24 dev ${BR}
    ip link set dev ${BR} up

    # setup named network namespace and attach to cspace PID
    ip netns add ${NS}
    ip netns attach ${NS} ${ID}

    # Move the veth2 into network namespace
    ip link set veth2_${SP} netns ${NS}

    # Enable loopback interface on the new namespace
    ip netns exec ${NS} ip link set dev lo up

    # Setup the veth2 on the cspace
    ip netns exec ${NS} ip addr add 10.0.0.2/24 dev veth2_${SP}
    ip netns exec ${NS} ip link set dev veth2_${SP} up
    ip netns exec ${NS} ip route add default via 10.0.0.1

}

# Script main

if [ "$#" -lt ${MIN_ARGS} ]
then
    msg_usage
fi

if [ "$#" -gt 0 ]; then

    case $1 in
	-a|--add)
	    DEV=$2
	    CMD="add"
	    if [ "${DEV}" = "bridge" ]; then
		IFACE=$3
		BRIDGE=$4
	    elif [ "${DEV}" = "cspace" ]; then
		PID=$3
		SPACE=$4
		BRIDGE=$5
	    else
		exit 100
	    fi
	    shift
	    ;;
	*)
	    exit 100
    esac
fi

# Test PID and IFACE are valid
valid_args

if [ "$CMD" = "add" ]; then
    if [ "${DEV}" = "bridge" ]; then
	add_bridge $IFACE $BRIDGE
    elif [ "${DEV}" = "cspace" ]; then
	add_cspace $PID $SPACE $BRIDGE
    else
	exit 100
    fi
fi

exit 0
