/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>

#include "status.h"

static void copy_str(char *dst, size_t size, const char *src) {

    if (dst == NULL || size == 0) return;

    if (src == NULL) {
        dst[0] = '\0';
        return;
    }

    snprintf(dst, size, "%s", src);
}

const char *status_state_str(InitState state) {

    switch (state) {
    case InitStateStarting:           return "starting";
    case InitStateCheckTools:         return "check-tools";
    case InitStateStartOvs:           return "start-ovs";
    case InitStateSetupEpcIf:         return "setup-epc-if";
    case InitStateSetupTun:           return "setup-tun";
    case InitStateSetupBridge:        return "setup-bridge";
    case InitStateSetupForwarding:    return "setup-forwarding";
    case InitStateSetupGateway:       return "setup-gateway";
    case InitStateSetupFlows:         return "setup-flows";
    case InitStateSetupPolicyRouting: return "setup-policy-routing";
    case InitStateReady:              return "ready";
    case InitStateFailed:             return "failed";
    default:                          return "unknown";
    }
}

void status_init(AppStatus *status) {

    if (status == NULL) return;

    memset(status, 0, sizeof(AppStatus));
    pthread_mutex_init(&status->mutex, NULL);

    status->state = InitStateStarting;
    status->ready = false;
    copy_str(status->reason, sizeof(status->reason), "starting");
}

void status_destroy(AppStatus *status) {

    if (status == NULL) return;

    pthread_mutex_destroy(&status->mutex);
}

void status_set(AppStatus *status, InitState state, const char *reason) {

    if (status == NULL) return;

    pthread_mutex_lock(&status->mutex);

    status->state = state;
    status->ready = (state == InitStateReady);

    if (reason != NULL) {
        copy_str(status->reason, sizeof(status->reason), reason);
    }

    pthread_mutex_unlock(&status->mutex);
}

void status_fail(AppStatus *status, const char *reason) {

    status_set(status, InitStateFailed, reason);
}

bool status_is_ready(AppStatus *status) {

    bool ready;

    if (status == NULL) return false;

    pthread_mutex_lock(&status->mutex);
    ready = status->ready;
    pthread_mutex_unlock(&status->mutex);

    return ready;
}

JsonObj *status_to_json(AppStatus *status, Config *config) {

    JsonObj *root;
    JsonObj *bridge;
    JsonObj *ovs;
    JsonObj *ue;
    JsonObj *tun;
    JsonObj *routing;
    JsonObj *gateway;
    JsonObj *steps;
    InitState state;
    bool ready;
    bool toolsOk;
    bool ovsdbRunning;
    bool vswitchdRunning;
    bool epcIfReady;
    bool tunReady;
    bool bridgeReady;
    bool forwardingReady;
    bool gatewayReady;
    bool flowsReady;
    bool policyRoutingReady;
    char reason[STATUS_REASON_LEN];
    char mgmtSocket[INIT_NETWORK_MAX_STR * 2];

    if (status == NULL || config == NULL) return NULL;

    pthread_mutex_lock(&status->mutex);

    state              = status->state;
    ready              = status->ready;
    toolsOk            = status->toolsOk;
    ovsdbRunning       = status->ovsdbRunning;
    vswitchdRunning    = status->vswitchdRunning;
    epcIfReady         = status->epcIfReady;
    tunReady           = status->tunReady;
    bridgeReady        = status->bridgeReady;
    forwardingReady    = status->forwardingReady;
    gatewayReady       = status->gatewayReady;
    flowsReady         = status->flowsReady;
    policyRoutingReady = status->policyRoutingReady;
    copy_str(reason, sizeof(reason), status->reason);

    pthread_mutex_unlock(&status->mutex);

    snprintf(mgmtSocket, sizeof(mgmtSocket), "%s/%s.mgmt",
             config->mgmtDir, config->bridge);

    root = json_object();
    if (root == NULL) return NULL;

    bridge  = json_object();
    ovs     = json_object();
    ue      = json_object();
    tun     = json_object();
    routing = json_object();
    gateway = json_object();
    steps   = json_object();

    json_object_set_new(root, "ready", json_boolean(ready));
    json_object_set_new(root, "state", json_string(status_state_str(state)));
    json_object_set_new(root, "reason", json_string(reason));

    json_object_set_new(bridge, "name", json_string(config->bridge));
    json_object_set_new(bridge, "address", json_string(config->bridgeAddr));
    json_object_set_new(bridge, "netmask", json_string(config->bridgeNetmask));
    json_object_set_new(bridge, "cidr", json_string(config->bridgeCidr));
    json_object_set_new(bridge, "managementSocket", json_string(mgmtSocket));
    json_object_set_new(bridge, "openflow", json_string(config->openflow));
    json_object_set_new(root,   "bridge", bridge);

    json_object_set_new(ovs, "ovsdb",
                        json_string(ovsdbRunning ? "running" : "unknown"));
    json_object_set_new(ovs, "vswitchd",
                        json_string(vswitchdRunning ? "running" : "unknown"));
    json_object_set_new(ovs, "runDir", json_string(config->runDir));
    json_object_set_new(ovs, "dbDir", json_string(config->dbDir));
    json_object_set_new(root, "ovs", ovs);

    json_object_set_new(ue,   "cidr", json_string(config->ueCidr));
    json_object_set_new(ue,   "defaultDrop", json_boolean(config->defaultDrop));
    json_object_set_new(root, "ue", ue);

    json_object_set_new(tun,  "enabled", json_boolean(config->tunEnable));
    json_object_set_new(tun,  "interface", json_string(config->tunIf));
    json_object_set_new(tun,  "primaryCidr",
                        json_string(config->tunPrimaryCidr));
    json_object_set_new(root, "tun", tun);

    json_object_set_new(routing, "externalIf",
                        json_string(config->externalIf));
    json_object_set_new(routing, "ipForward",
                        json_boolean(config->enableIpForward));
    json_object_set_new(routing, "nat", json_boolean(config->enableNat));
    json_object_set_new(routing, "policyRouting",
                        json_boolean(config->enablePolicyRouting));
    json_object_set_new(routing, "tunTable", json_integer(config->tunTable));
    json_object_set_new(routing, "bridgeTable",
                        json_integer(config->bridgeTable));
    json_object_set_new(root,    "routing", routing);

    json_object_set_new(gateway, "enabled",
                        json_boolean(config->gatewayEnable));
    json_object_set_new(gateway, "mode", json_string(config->gatewayMode));
    json_object_set_new(gateway, "name", json_string(config->gatewayName));
    json_object_set_new(gateway, "bridgeIf",
                        json_string(config->gatewayBridgeIf));
    json_object_set_new(gateway, "namespaceIf",
                        json_string(config->gatewayNamespaceIf));
    json_object_set_new(gateway, "address", json_string(config->gatewayAddr));
    json_object_set_new(gateway, "ip", json_string(config->gatewayIp));
    json_object_set_new(root,    "gateway", gateway);

    json_object_set_new(steps, "toolsOk", json_boolean(toolsOk));
    json_object_set_new(steps, "ovsdbRunning", json_boolean(ovsdbRunning));
    json_object_set_new(steps, "vswitchdRunning",
                        json_boolean(vswitchdRunning));
    json_object_set_new(steps, "epcIfReady", json_boolean(epcIfReady));
    json_object_set_new(steps, "tunReady", json_boolean(tunReady));
    json_object_set_new(steps, "bridgeReady", json_boolean(bridgeReady));
    json_object_set_new(steps, "forwardingReady",
                        json_boolean(forwardingReady));
    json_object_set_new(steps, "gatewayReady", json_boolean(gatewayReady));
    json_object_set_new(steps, "flowsReady", json_boolean(flowsReady));
    json_object_set_new(steps, "policyRoutingReady",
                        json_boolean(policyRoutingReady));
    json_object_set_new(root,  "steps", steps);

    return root;
}
