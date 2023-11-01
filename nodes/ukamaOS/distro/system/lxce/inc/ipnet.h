/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef LXCE_IPNET_H
#define LXCE_IPNET_H

#define DEF_BRIDGE "br0"
#define DEF_IFACE  "eth0"
#define NET_EXEC   "setup_space_network.sh"
#define PATH       "/sbin"
#define IP_BIN     "ip"
#define PING_BIN   "/bin/ping"
#define TEST_IP    "8.8.8.8"

#define TRUE  1
#define FALSE 0

#define IPNET_DEV_TYPE_BRIDGE 1
#define IPNET_DEV_TYPE_CSPACE 2

#define IPNET_DEV_BRIDGE "bridge"
#define IPNET_DEV_CSPACE "cspace"

int ipnet_setup(int type, char *brName, char *iface, char *brIP, char *vethIP,
		char *spName, pid_t pid);
int ipnet_test(char *spName);

#endif /* LXCE_IPNET_H */
