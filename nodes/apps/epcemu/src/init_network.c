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

    status_set(status, EpcemuStateCheckingInitNetwork,
               "checking init-network");

    snprintf(url, sizeof(url), "%s/v1/status", config->initNetworkUrl);

    if (!http_get_json(url, &root, &code)) {
        status_fail(status, "failed to call init-network status");
        return USYS_FALSE;
    }

    if (code != 200 || root == NULL) {
        if (root != NULL) json_decref(root);
        status_fail(status, "init-network status failed");
        return USYS_FALSE;
    }

    ready = json_object_get(root, "ready");
    if (!json_is_true(ready)) {
        json_decref(root);
        status_fail(status, "init-network is not ready");
        return USYS_FALSE;
    }

    bridge = json_object_get(root, "bridge");
    ue = json_object_get(root, "ue");
    if (!json_is_object(bridge) || !json_is_object(ue)) {
        json_decref(root);
        status_fail(status, "init-network status missing bridge/ue");
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
        status_fail(status, "init-network status missing ue.cidr");
        return USYS_FALSE;
    }

    config->initNetworkReady = true;
    json_decref(root);

    return USYS_TRUE;
}
