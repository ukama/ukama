#!/bin/sh
#Host config

set -x

TUNIF=tun3
EXTIF=wlo1
BRIF=br0
ETHIF=enp60s0

## EPC SCTP Interface
ifconfig $ETHIF:0 10.102.81.3 > /dev/null 2>&1

## EPC GTPU interface
ifconfig $ETHIF:1 10.102.81.75 up > /dev/null 2>&1

## make tun and bridge
ip link delete ${TUNIF}

openvpn --mktun --dev ${TUNIF}
ip link set ${TUNIF} up
ip addr add 192.168.9.20/22 dev ${TUNIF}
ip addr add 192.168.9.21/22 dev ${TUNIF} > /dev/null 2>&1
ip addr add 192.168.9.22/22 dev ${TUNIF} > /dev/null 2>&1
ip addr add 192.168.9.23/22 dev ${TUNIF} > /dev/null 2>&1
ip addr add 192.168.9.24/22 dev ${TUNIF} > /dev/null 2>&1
ip addr add 192.168.9.25/22 dev ${TUNIF} > /dev/null 2>&1
ip addr add 192.168.9.26/22 dev ${TUNIF} > /dev/null 2>&1
ip addr add 192.168.9.27/22 dev ${TUNIF} > /dev/null 2>&1
ip addr add 192.168.9.28/22 dev ${TUNIF} > /dev/null 2>&1

sysctl net.ipv4.conf.${TUNIF}.rp_filter=2
sysctl net.ipv4.conf.${TUNIF}.accept_local=1
sysctl net.ipv4.conf.${TUNIF}.log_martians=1
echo 1 > /proc/sys/net/ipv4/ip_forward

ovs-vsctl del-br ${BRIF}

sleep 1

ovs-vsctl add-br ${BRIF}

ovs-vsctl show 

sleep 1

ifconfig ${BRIF} 10.10.10.1 netmask 255.255.255.0 up

#TODO check
#iptables -t nat -A POSTROUTING -s 192.168.8.0/22 -o ${EXTIF} -j MASQUERADE

iptables -t nat -A POSTROUTING -s 10.10.10.0/24 -o ${EXTIF} -j MASQUERADE
iptables -t nat -L -n -v

sleep 2

iptables -A FORWARD -i ${EXTIF} -o ${BRIF} -m state --state RELATED,ESTABLISHED -j ACCEPT
iptables -A FORWARD -i ${BRIF} -o ${EXTIF} -j ACCEPT

iptables -A OUTPUT -o ${EXTIF} -p tcp --dport 443 -m state --state NEW,ESTABLISHED -j ACCEPT
iptables -A INPUT -i ${EXTIF} -p tcp --sport 443 -m state --state ESTABLISHED -j ACCEPT

iptables -t filter -L -n -v

sleep 2

docker kill cont1 

sleep 2

## Start docker container
docker run -ti -d --rm --name cont1 --cap-add NET_ADMIN --net=none vthakur7f/epc_gateway:v0.0.1

sleep 2

## Attach docker ip
ovs-docker add-port ${BRIF} eth0 cont1 --ipaddress=10.10.10.11/24 --gateway=10.10.10.1

## List rules
ip rule list

sleep 2

## Add table 2000 for interface ${TUNIF}
ip route add default via 10.10.10.11 dev ${BRIF} table 2000
#ip route add 10.10.10.0/24 dev ${BRIF} table 2000
ip route add 192.168.8.0/22 via 192.168.9.20 dev ${TUNIF} table 2000

ip route show table 2000

sleep 1

## Add table 1000 for interface ${BRIF}
ip route add 10.10.10.0/24 via 10.10.10.1 dev ${BRIF} table 1000
ip route add 192.168.8.0/22 via 192.168.9.20 dev ${TUNIF} table 1000
ip route add default via 192.168.0.1 dev ${EXTIF} table 1000

ip route show table 1000

sleep 2

## Add route tables to interface
ip rule add iif ${BRIF} table 1000
ip rule add iif ${TUNIF} table 2000

ip route

sleep 2

## set iptables rules
iptables -t filter -I FORWARD -i ${TUNIF} -o ${BRIF} -j ACCEPT
iptables -t filter -I FORWARD -i ${BRIF} -o ${TUNIF} -m state --state NEW,RELATED,ESTABLISHED -j ACCEPT


OUT_IF=${EXTIF}
IM_IF=${BRIF}
echo 1 > /proc/sys/net/ipv4/ip_forward
#iptables -t nat -A POSTROUTING -s 192.168.8.0/22 -o $OUT_IF -j MASQUERADE
iptables -t nat -A POSTROUTING -s 10.10.10.0/24 -o $OUT_IF -j MASQUERADE
iptables -A FORWARD -i $OUT_IF -o $IM_IF -m state --state RELATED,ESTABLISHED -j ACCEPT
iptables -A FORWARD -i $IM_IF -o $OUT_IF -j ACCEPT
iptables -A OUTPUT -o $OUT_IF -p tcp --dport 443 -m state --state NEW,ESTABLISHED -j ACCEPT
iptables -A INPUT -i $OUT_IF -p tcp --sport 443 -m state --state ESTABLISHED -j ACCEPT



#Docker container config
##In docker file you to set masqerade

docker exec -it cont1 /sbin/ifconfig

docker exec -it cont1 /bin/ping -c2 google.com

docker exec -it cont1 /sbin/iptables -t nat -A POSTROUTING -s 192.168.8.0/22 -o eth0 -j MASQUERADE

docker exec -it cont1 /sbin/ip route

# remove the default flows
# this is commented as the vthakur7f/epc_gateway:v0.0.1 requires this flow. We might need to add in and out flows the docker IP 
#ovs-ofctl -O OpenFlow15 del-flows ${BRIF}

#in and Out flows for the docker image
#ovs-ofctl -O OpenFlow15 add-flow ${BRIF} "priority=100,ip,nw_dst=10.10.10.11, actions=normal"
#ovs-ofctl -O OpenFlow15 add-flow ${BRIF} "priority=100,ip,nw_src=10.10.10.11, actions=normal"