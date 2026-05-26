/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef CONFIG_H_
#define CONFIG_H_

#include <stdbool.h>

#include "usys_services.h"

#define INIT_NETWORK_SERVICE_NAME       SERVICE_INIT_NETWORK
#define INIT_NETWORK_APP_NAME           SERVICE_INIT_NETWORK

#define INIT_NETWORK_MAX_STR            256
#define INIT_NETWORK_MAX_EXTRA_IPS      16

#define DEF_CONFIG_FILE                 "/ukama/configs/init-network/config.toml"
#define DEF_LOG_LEVEL                   "TRACE"

#define DEF_SERVICE_PORT                18091
#define DEF_CMD_TIMEOUT_SEC             10

#define DEF_OVS_BRIDGE                  "br0"
#define DEF_OVS_OF_VERSION              "OpenFlow15"
#define DEF_OVS_RUN_DIR                 "/var/run/openvswitch"
#define DEF_OVS_MGMT_DIR                DEF_OVS_RUN_DIR
#define DEF_OVS_DB_DIR                  "/etc/openvswitch"
#define DEF_OVS_SCHEMA                  "/usr/share/openvswitch/vswitch.ovsschema"

#define DEF_BRIDGE_ADDR                 "10.10.10.1"
#define DEF_BRIDGE_NETMASK              "255.255.255.0"
#define DEF_BRIDGE_CIDR                 "10.10.10.1/24"
#define DEF_BRIDGE_SUBNET               "10.10.10.0/24"

#define DEF_UE_CIDR                     "192.168.8.0/22"

#define DEF_TUN_ENABLE                  false
#define DEF_TUN_IF                      "tun3"
#define DEF_TUN_PRIMARY_CIDR            "192.168.8.1/22"

#define DEF_EPC_ENABLE                  false
#define DEF_EPC_SCTP_IF                 "enp60s0"
#define DEF_EPC_SCTP_ADDR               "10.102.81.3"
#define DEF_EPC_GTPU_ADDR               "10.102.81.75"

#define DEF_EXTERNAL_IF                 "eth0"

#define DEF_GATEWAY_ENABLE              true
#define DEF_GATEWAY_MODE                "root"
#define DEF_GATEWAY_NAME                "ukama-gw"
#define DEF_GATEWAY_BRIDGE_IF           "gw-br"
#define DEF_GATEWAY_NAMESPACE_IF        "gw0"
#define DEF_GATEWAY_ADDR                "10.10.10.11/24"
#define DEF_GATEWAY_IP                  "10.10.10.11"

#define DEF_TUN_TABLE                   2000
#define DEF_BRIDGE_TABLE                1000
#define DEF_DEFAULT_DROP                true

typedef struct {

    char *configFile;
    char *logLevel;

    int servicePort;
    int cmdTimeoutSec;

    char *bridge;
    char *openflow;
    char *mgmtDir;
    char *runDir;
    char *dbDir;
    char *schema;

    char *bridgeAddr;
    char *bridgeNetmask;
    char *bridgeCidr;
    char *bridgeSubnet;

    char *ueCidr;

    bool tunEnable;
    char *tunIf;
    char *tunPrimaryCidr;
    char *tunExtraCidrs[INIT_NETWORK_MAX_EXTRA_IPS];
    int tunExtraCount;

    bool epcEnable;
    char *epcSctpIf;
    char *epcSctpAddr;
    char *epcGtpuAddr;

    char *externalIf;

    bool gatewayEnable;
    char *gatewayMode;
    char *gatewayName;
    char *gatewayBridgeIf;
    char *gatewayNamespaceIf;
    char *gatewayAddr;
    char *gatewayIp;

    int tunTable;
    int bridgeTable;

    bool defaultDrop;

} Config;

Config *config_init(void);
void config_free(Config *config);
bool config_load(Config *config, const char *path);

#endif /* CONFIG_H_ */
