#!/bin/sh
# Copyright (c) 2021-present, Ukama Inc.
# All rights reserved.

# Script to setup networking for cSpaces and cApps

# base parameters
UKAMA_OS=../../
DEF_BRIDGE="ukama_bridge"
MIN_ARGS=2

usage() {
    echo 'Usage: setup_space_network.sh interface space_name'
    exit
}

msg_usage() {
    echo "Usage:"
    echo "      setup_space_network.sh [args]"
    echo ""
    echo "Args:"
    echo "   pid         # PID of the cspace process"
    echo "   interface   # Interface to use to connect with Internet"
    echo "   space_name  # name of the space"
    echo "   bridge_name # name to the master bridge (optional)"
    exit 99
}

valid_args() {

    # Test valid PID
    ps -p ${PID} > /dev/null
    if [ "$?" -eq "1" ]
    then
	exit 1
    fi
      
    # Test valid interface
    ip -o a show | cut -d ' ' -f 2,7 | \
	awk '{print $1}' | grep ${IFACE} > /dev/null
    if [ "$?" -eq "1" ]
    then
	exit 2
    fi
}

# Script main

if [ "$#" -lt ${MIN_ARGS} ]
then
    msg_usage
fi

PID=$1
IFACE=$2
SPACE=$3
BRIDGE=$4
NET_NS=ns_${SPACE}

if [ -z "${BRIDGE}"]
then
    BRIDGE=${DEF_BRIDGE}
fi

# Test PID and IFACE are valid
valid_args

# Basic flow is:
#
#  setup paired veth device (veth1_* for host and veth2_* for cspace)
#  setup bridge and connect the host veth to it as master
#  add network namespace and attach it to the cspace PID
#  move the veth2_* to the network namespace
#  bring up loopback in the network namespace
#  setup host machine to allow NATing
#  finally, setup veth2_* in the network namespace and setup routing
#

# Setup paired veth for each cspace on the host.
ip link   add dev veth1_${SPACE} type veth peer name veth2_${SPACE}

# Bring up the host iface
ip link   set dev veth1_${SPACE} up
ip tuntap add tap_${SPACE} mode tap
ip link   set dev tapm_${SPACE} up

# Setup bridge
ip link   add $(BRIDGE} type bridge

# Attach ifaces to the bridge
ip link set tap_${SPACE}   master ${BRIDGE}
ip link set veth1_${SPACE} master ${BRIDGE}

# Give address to the bridge
ip addr add 10.0.0.1/24 dev ${BRIDGE}

# setup named network namespace and attach to cspace PID
ip netns add ${NET_NS}
ip netns attach ${NET_NS} ${PID}

# Move the veth2 into network namespace
ip link set veth2_${SPACE} netns ${NET_NS}

# Enable loopback interface on the new namespace
ip netns exec ${NET_NS} ip link set dev lo up

# Setup host machine to allow NATing.
iptables -t nat -A POSTROUTING -o ${BRIDGE} -j MASQUERADE
iptables -t nat -A POSTROUTING -o ${IFACE}  -j MASQUERADE

# Setup the veth2 on the cspace.
ip netns ${NET_NS} addr add 10.0.0.2/24 dev veth2_${CSPACE}
ip netns ${NET_NS} link set dev veth2_${CSPACE} up
ip netns ${NET_NS} route add default via 10.0.0.1

echo "Done"

exit 0
