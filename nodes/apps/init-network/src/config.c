/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "toml.h"

#include "config.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_file.h"
#include "usys_services.h"

static char *dup_or_default(const char *value, const char *defValue) {

    if (value != NULL && value[0] != '\0') {
        return strdup(value);
    }

    if (defValue != NULL) {
        return strdup(defValue);
    }

    return NULL;
}

static char *toml_string_or(toml_table_t *table,
                            const char *key,
                            const char *defValue) {

    toml_datum_t datum;

    if (table == NULL) return dup_or_default(NULL, defValue);

    datum = toml_string_in(table, key);
    if (datum.ok) return datum.u.s;

    return dup_or_default(NULL, defValue);
}

static int toml_int_or(toml_table_t *table, const char *key, int defValue) {

    toml_datum_t datum;

    if (table == NULL) return defValue;

    datum = toml_int_in(table, key);
    if (datum.ok) return (int)datum.u.i;

    return defValue;
}

static bool toml_bool_or(toml_table_t *table,
                         const char *key,
                         bool defValue) {

    toml_datum_t datum;

    if (table == NULL) return defValue;

    datum = toml_bool_in(table, key);
    if (datum.ok) return datum.u.b;

    return defValue;
}

static void load_tun_extra(Config *config, toml_table_t *tun) {

    toml_array_t *arr;
    toml_datum_t datum;
    int i;

    if (config == NULL || tun == NULL) return;

    for (i = 0; i < config->tunExtraCount; i++) {
        usys_free(config->tunExtraCidrs[i]);
        config->tunExtraCidrs[i] = NULL;
    }
    config->tunExtraCount = 0;

    arr = toml_array_in(tun, "extra_cidrs");
    if (arr == NULL) return;

    for (i = 0; i < INIT_NETWORK_MAX_EXTRA_IPS; i++) {
        datum = toml_string_at(arr, i);
        if (!datum.ok) break;
        config->tunExtraCidrs[config->tunExtraCount++] = datum.u.s;
    }
}

static void resolve_service_port(Config *config, int fallbackPort) {

    int port;

    if (config == NULL || config->serviceName == NULL) {
        return;
    }

    port = usys_find_service_port(config->serviceName);
    if (port > 0) {
        config->servicePort = port;
        usys_log_debug("resolved service port from /etc/services: %s=%d",
                       config->serviceName,
                       config->servicePort);
        return;
    }

    config->servicePort = fallbackPort;
    usys_log_debug("service %s not found in /etc/services, using port %d",
                   config->serviceName,
                   config->servicePort);
}

void config_set_defaults(Config *config) {

    if (config == NULL) return;

    memset(config, 0, sizeof(Config));

    config->serviceName         = strdup(INIT_NETWORK_SERVICE_NAME);
    config->servicePort         = DEF_SERVICE_PORT;
    config->cmdTimeoutSec       = DEF_CMD_TIMEOUT_SEC;

    config->bridge              = strdup(DEF_OVS_BRIDGE);
    config->openflow            = strdup(DEF_OVS_OF_VERSION);
    config->mgmtDir             = strdup(DEF_OVS_MGMT_DIR);
    config->runDir              = strdup(DEF_OVS_RUN_DIR);
    config->dbDir               = strdup(DEF_OVS_DB_DIR);
    config->schema              = strdup(DEF_OVS_SCHEMA);

    config->bridgeAddr          = strdup(DEF_BRIDGE_ADDR);
    config->bridgeNetmask       = strdup(DEF_BRIDGE_NETMASK);
    config->bridgeCidr          = strdup(DEF_BRIDGE_CIDR);
    config->bridgeSubnet        = strdup(DEF_BRIDGE_SUBNET);

    config->ueCidr              = strdup(DEF_UE_CIDR);
    config->defaultDrop         = true;

    config->tunEnable           = DEF_TUN_ENABLE;
    config->tunIf               = strdup(DEF_TUN_IF);
    config->tunPrimaryCidr      = strdup(DEF_TUN_PRIMARY_CIDR);

    config->epcEnable           = DEF_EPC_ENABLE;
    config->epcIf               = strdup(DEF_EPC_SCTP_IF);
    config->epcSctpAddr         = strdup(DEF_EPC_SCTP_ADDR);
    config->epcGtpuAddr         = strdup(DEF_EPC_GTPU_ADDR);

    config->externalIf          = strdup(DEF_EXT_IF);
    config->enableIpForward     = DEF_FORWARD_ENABLE;
    config->enableNat           = DEF_NAT_ENABLE;
    config->enablePolicyRouting = DEF_POLICY_ROUTING_ENABLE;

    config->gatewayContainer    = strdup(DEF_GATEWAY_CONTAINER);
    config->gatewayAddr         = strdup(DEF_GATEWAY_ADDR);
    config->gatewayIp           = strdup(DEF_GATEWAY_IP);
}

bool config_load_from_file(Config *config, const char *path) {

    FILE *fp;
    char errbuf[256];
    toml_table_t *root;
    toml_table_t *server;
    toml_table_t *ovs;
    toml_table_t *bridge;
    toml_table_t *ue;
    toml_table_t *tun;
    toml_table_t *epc;
    toml_table_t *routing;
    toml_table_t *gateway;
    int fallbackPort;

    if (config == NULL || path == NULL) return false;

    fallbackPort = config->servicePort;

    fp = fopen(path, "r");
    if (fp == NULL) {
        usys_log_debug("config file not found: %s, using defaults", path);
        resolve_service_port(config, fallbackPort);
        return true;
    }

    root = toml_parse_file(fp, errbuf, sizeof(errbuf));
    fclose(fp);

    if (root == NULL) {
        usys_log_error("failed to parse config %s: %s", path, errbuf);
        return false;
    }

    server  = toml_table_in(root, "server");
    ovs     = toml_table_in(root, "ovs");
    bridge  = toml_table_in(root, "bridge");
    ue      = toml_table_in(root, "ue");
    tun     = toml_table_in(root, "tun");
    epc     = toml_table_in(root, "epc");
    routing = toml_table_in(root, "routing");
    gateway = toml_table_in(root, "gateway");

    fallbackPort = toml_int_or(server, "port", config->servicePort);
    config->cmdTimeoutSec = toml_int_or(server, "cmd_timeout_sec",
                                        config->cmdTimeoutSec);

    usys_free(config->bridge);
    usys_free(config->openflow);
    usys_free(config->mgmtDir);
    usys_free(config->runDir);
    usys_free(config->dbDir);
    usys_free(config->schema);

    config->bridge   = toml_string_or(ovs, "bridge", DEF_OVS_BRIDGE);
    config->openflow = toml_string_or(ovs, "openflow", DEF_OVS_OF_VERSION);
    config->mgmtDir  = toml_string_or(ovs, "management_dir", DEF_OVS_MGMT_DIR);
    config->runDir   = toml_string_or(ovs, "run_dir", DEF_OVS_RUN_DIR);
    config->dbDir    = toml_string_or(ovs, "db_dir", DEF_OVS_DB_DIR);
    config->schema   = toml_string_or(ovs, "schema", DEF_OVS_SCHEMA);

    usys_free(config->bridgeAddr);
    usys_free(config->bridgeNetmask);
    usys_free(config->bridgeCidr);
    usys_free(config->bridgeSubnet);

    config->bridgeAddr    = toml_string_or(bridge, "address", DEF_BRIDGE_ADDR);
    config->bridgeNetmask = toml_string_or(bridge, "netmask",
                                           DEF_BRIDGE_NETMASK);
    config->bridgeCidr    = toml_string_or(bridge, "cidr", DEF_BRIDGE_CIDR);
    config->bridgeSubnet  = toml_string_or(bridge, "subnet", DEF_BRIDGE_SUBNET);

    usys_free(config->ueCidr);
    config->ueCidr      = toml_string_or(ue, "cidr", DEF_UE_CIDR);
    config->defaultDrop = toml_bool_or(ue, "default_drop", true);

    usys_free(config->tunIf);
    usys_free(config->tunPrimaryCidr);
    config->tunEnable      = toml_bool_or(tun, "enable", DEF_TUN_ENABLE);
    config->tunIf          = toml_string_or(tun, "interface", DEF_TUN_IF);
    config->tunPrimaryCidr = toml_string_or(tun, "primary_cidr",
                                            DEF_TUN_PRIMARY_CIDR);
    load_tun_extra(config, tun);

    usys_free(config->epcIf);
    usys_free(config->epcSctpAddr);
    usys_free(config->epcGtpuAddr);
    config->epcEnable   = toml_bool_or(epc, "enable", DEF_EPC_ENABLE);
    config->epcIf       = toml_string_or(epc, "interface", DEF_EPC_SCTP_IF);
    config->epcSctpAddr = toml_string_or(epc, "sctp_address", DEF_EPC_SCTP_ADDR);
    config->epcGtpuAddr = toml_string_or(epc, "gtpu_address", DEF_EPC_GTPU_ADDR);

    usys_free(config->externalIf);
    config->externalIf          = toml_string_or(routing, "external_if",
                                                 DEF_EXT_IF);
    config->enableIpForward     = toml_bool_or(routing, "enable_ip_forward",
                                               DEF_FORWARD_ENABLE);
    config->enableNat           = toml_bool_or(routing, "enable_nat",
                                               DEF_NAT_ENABLE);
    config->enablePolicyRouting = toml_bool_or(routing, "enable_policy_routing",
                                               DEF_POLICY_ROUTING_ENABLE);

    usys_free(config->gatewayContainer);
    usys_free(config->gatewayAddr);
    usys_free(config->gatewayIp);
    config->gatewayContainer = toml_string_or(gateway, "container",
                                              DEF_GATEWAY_CONTAINER);
    config->gatewayAddr      = toml_string_or(gateway, "address",
                                              DEF_GATEWAY_ADDR);
    config->gatewayIp        = toml_string_or(gateway, "ip", DEF_GATEWAY_IP);

    resolve_service_port(config, fallbackPort);

    toml_free(root);
    return true;
}

void config_free(Config *config) {

    int i;

    if (config == NULL) return;

    usys_free(config->serviceName);
    usys_free(config->bridge);
    usys_free(config->openflow);
    usys_free(config->mgmtDir);
    usys_free(config->runDir);
    usys_free(config->dbDir);
    usys_free(config->schema);
    usys_free(config->bridgeAddr);
    usys_free(config->bridgeNetmask);
    usys_free(config->bridgeCidr);
    usys_free(config->bridgeSubnet);
    usys_free(config->ueCidr);
    usys_free(config->tunIf);
    usys_free(config->tunPrimaryCidr);
    usys_free(config->epcIf);
    usys_free(config->epcSctpAddr);
    usys_free(config->epcGtpuAddr);
    usys_free(config->externalIf);
    usys_free(config->gatewayContainer);
    usys_free(config->gatewayAddr);
    usys_free(config->gatewayIp);

    for (i = 0; i < config->tunExtraCount; i++) {
        usys_free(config->tunExtraCidrs[i]);
    }

    memset(config, 0, sizeof(Config));
}
