/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>

#include "epcemu.h"
#include "http_client.h"
#include "init_network.h"

#define INIT_NETWORK_RECONCILE_TIMEOUT_SEC 30

static int json_bool_path(JsonObj *root,
                          const char *parent,
                          const char *key,
                          int *value) {

    JsonObj *obj;
    JsonObj *leaf;

    if (root == NULL || parent == NULL || key == NULL || value == NULL) {
        return USYS_FALSE;
    }

    obj = json_object_get(root, parent);
    if (!json_is_object(obj)) return USYS_FALSE;

    leaf = json_object_get(obj, key);
    if (json_is_boolean(leaf)) {
        *value = json_is_true(leaf) ? USYS_TRUE : USYS_FALSE;
        return USYS_TRUE;
    }

    return USYS_FALSE;
}

static int reconcile_response_ok(JsonObj *root) {

    int tunReady;
    int routingReady;

    tunReady = USYS_FALSE;
    routingReady = USYS_FALSE;

    return root != NULL &&
           json_is_true(json_object_get(root, "ready")) &&
           json_bool_path(root, "steps", "tunReady", &tunReady) &&
           json_bool_path(root, "steps", "policyRoutingReady",
                           &routingReady) && tunReady && routingReady;
}

static int read_tun_config(JsonObj *root, EpcemuConfig *config) {

    JsonObj *tun;
    JsonObj *enabled;
    const char *name;
    const char *addr;

    tun = json_object_get(root, "tun");
    if (!json_is_object(tun)) return USYS_FALSE;

    name = json_string_value(json_object_get(tun, "interface"));
    addr = json_string_value(json_object_get(tun, "primaryCidr"));
    enabled = json_object_get(tun, "enabled");
    if (name == NULL || addr == NULL || !json_is_boolean(enabled) ||
        name[0] == '\0' || addr[0] == '\0' ||
        strlen(name) >= sizeof(config->tunName) ||
        strlen(addr) >= sizeof(config->tunAddr)) return USYS_FALSE;

    snprintf(config->tunName, sizeof(config->tunName), "%s", name);
    snprintf(config->tunAddr, sizeof(config->tunAddr), "%s", addr);
    config->initNetworkRouted = reconcile_response_ok(root);

    return json_is_true(enabled) && config->initNetworkRouted;
}

int init_network_probe(EpcemuConfig *config, EpcemuStatus *status) {

    char url[EPCEMU_MAX_STR * 2];
    JsonObj *root;
    JsonObj *bridge;
    JsonObj *ue;
    JsonObj *ready;
    JsonObj *value;
    long code;

    if (config == NULL || status == NULL) return USYS_FALSE;

    root = NULL;
    code = 0;
    config->initNetworkReady = false;
    config->initNetworkRouted = false;

    status_set(status, EpcemuStateCheckingInitNetwork,
               "checking init-network");

    snprintf(url, sizeof(url), "%s/v1/status", config->initNetworkUrl);

    if (!http_get_json(url, &root, &code)) {
        status_set(status,
                   EpcemuStateCheckingInitNetwork,
                   "waiting for init-network status");
        return USYS_FALSE;
    }

    if (code != 200 || root == NULL) {
        if (root != NULL) json_decref(root);
        status_set(status,
                   EpcemuStateCheckingInitNetwork,
                   "waiting for init-network status");
        return USYS_FALSE;
    }

    /*
     * In this manifest, init-network runs in boot before EPCEMU.
     * Therefore this should already be ready by the time EPCEMU starts.
     */
    ready = json_object_get(root, "ready");
    if (!json_is_true(ready)) {
        json_decref(root);
        status_set(status,
                   EpcemuStateCheckingInitNetwork,
                   "waiting for init-network readiness");
        return USYS_FALSE;
    }

    bridge = json_object_get(root, "bridge");
    ue = json_object_get(root, "ue");
    if (!json_is_object(bridge) || !json_is_object(ue)) {
        json_decref(root);
        status_set(status,
                   EpcemuStateCheckingInitNetwork,
                   "waiting for init-network bridge and UE configuration");
        return USYS_FALSE;
    }

    value = json_object_get(bridge, "name");
    if (json_is_string(value)) {
        snprintf(config->bridge, sizeof(config->bridge), "%s",
                 json_string_value(value));
    }

    value = json_object_get(bridge, "cidr");
    if (json_is_string(value)) {
        snprintf(config->bridgeCidr, sizeof(config->bridgeCidr), "%s",
                 json_string_value(value));
    }

    value = json_object_get(ue, "cidr");
    if (json_is_string(value)) {
        snprintf(config->ueCidr, sizeof(config->ueCidr), "%s",
                 json_string_value(value));
    }

    if (config->ueCidr[0] == '\0') {
        json_decref(root);
        status_set(status,
                   EpcemuStateCheckingInitNetwork,
                   "waiting for init-network UE CIDR");
        return USYS_FALSE;
    }

    if (!read_tun_config(root, config)) {
        json_decref(root);
        status_set(status, EpcemuStateCheckingInitNetwork,
                   "waiting for init-network TUN configuration");
        return USYS_FALSE;
    }

    config->initNetworkReady = true;
    json_decref(root);

    return USYS_TRUE;
}

int init_network_reconcile(EpcemuConfig *config, EpcemuStatus *status) {

    char url[EPCEMU_MAX_STR * 2];
    JsonObj *root;
    long code;

    if (config == NULL || status == NULL) return USYS_FALSE;

    root = NULL;
    code = 0;
    config->initNetworkRouted = false;

    status_set(status, EpcemuStateReconcilingInitNetwork,
               "reconciling EPC tun routing through OVS");

    snprintf(url, sizeof(url), "%s/v1/reconcile", config->initNetworkUrl);

    if (!http_send_json_timeout("POST",
                                url,
                                NULL,
                                &root,
                                &code,
                                INIT_NETWORK_RECONCILE_TIMEOUT_SEC)) {
        status_set(status,
                   EpcemuStateReconcilingInitNetwork,
                   "waiting for init-network reconcile");
        return USYS_FALSE;
    }

    if (code != 200 || root == NULL) {
        if (root != NULL) json_decref(root);
        status_set(status,
                   EpcemuStateReconcilingInitNetwork,
                   "waiting for init-network reconcile");
        return USYS_FALSE;
    }

    if (!reconcile_response_ok(root)) {
        json_decref(root);
        status_set(status,
                   EpcemuStateReconcilingInitNetwork,
                   "waiting for init-network policy routing");
        return USYS_FALSE;
    }

    config->initNetworkRouted = true;
    json_decref(root);

    usys_log_info("init-network reconciled tun=%s bridge=%s ue=%s",
                  config->tunName,
                  config->bridge,
                  config->ueCidr);

    return USYS_TRUE;
}
