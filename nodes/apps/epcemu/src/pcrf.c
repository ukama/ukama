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
#include "netutil.h"
#include "pcrf.h"

int pcrf_probe(EpcemuConfig *config, EpcemuStatus *status) {

    char url[EPCEMU_MAX_STR * 2];
    JsonObj *root;
    JsonObj *ready;
    long code;

    if (config == NULL || status == NULL) return USYS_FALSE;

    root = NULL;
    code = 0;

    status_set(status, EpcemuStateCheckingPcrf, "checking pcrf");

    snprintf(url, sizeof(url), "%s/v1/status", config->pcrfUrl);

    if (!http_get_json(url, &root, &code)) {
        status_fail(status, "failed to call pcrf status");
        return USYS_FALSE;
    }

    if (code != 200 || root == NULL) {
        if (root != NULL) json_decref(root);
        status_fail(status, "pcrf status failed");
        return USYS_FALSE;
    }

    ready = json_object_get(root, "ready");
    if (!json_is_true(ready)) {
        json_decref(root);
        status_fail(status, "pcrf is not ready");
        return USYS_FALSE;
    }

    config->pcrfReady = true;
    json_decref(root);

    return USYS_TRUE;
}

int pcrf_create_session(EpcemuConfig *config, const char *imsi,
                        const char *ip, const char *apn) {

    char url[EPCEMU_MAX_STR * 2];
    JsonObj *body;
    JsonObj *imsiArray;
    JsonObj *reply;
    long code;
    uint32_t ipValue;
    int ret;

    if (config == NULL || imsi == NULL || ip == NULL) return USYS_FALSE;

    if (!ip_to_uint32(ip, &ipValue)) return USYS_FALSE;

    imsiArray = imsi_to_json_array(imsi);
    if (imsiArray == NULL) return USYS_FALSE;

    body = json_object();
    if (body == NULL) {
        json_decref(imsiArray);
        return USYS_FALSE;
    }

    json_object_set_new(body, "imsi",        imsiArray);
    json_object_set_new(body, "pdn_address", json_integer(ipValue));
    json_object_set_new(body, "apn_name",    json_string(apn ? apn : EPCEMU_DEF_APN));

    snprintf(url, sizeof(url), "%s/v1/session", config->pcrfUrl);

    reply = NULL;
    code  = 0;
    ret   = http_send_json("POST", url, body, &reply, &code);

    if (reply != NULL) json_decref(reply);
    json_decref(body);

    if (!ret || (code < 200 || code >= 300)) {
        usys_log_error("PCRF create session failed imsi=%s ip=%s http=%ld",
                       imsi, ip, code);
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

int pcrf_delete_session(EpcemuConfig *config, const char *imsi) {

    char url[EPCEMU_MAX_STR * 2];
    JsonObj *body;
    JsonObj *imsiArray;
    JsonObj *reply;
    long code;
    int ret;

    if (config == NULL || imsi == NULL) return USYS_FALSE;

    imsiArray = imsi_to_json_array(imsi);
    if (imsiArray == NULL) return USYS_FALSE;

    body = json_object();
    if (body == NULL) {
        json_decref(imsiArray);
        return USYS_FALSE;
    }

    json_object_set_new(body, "imsi", imsiArray);

    snprintf(url, sizeof(url), "%s/v1/session", config->pcrfUrl);

    reply = NULL;
    code  = 0;
    ret   = http_send_json("DELETE", url, body, &reply, &code);

    if (reply != NULL) json_decref(reply);
    json_decref(body);

    if (!ret || (code < 200 || code >= 300)) {
        usys_log_error("PCRF delete session failed imsi=%s http=%ld",
                       imsi, code);
        return USYS_FALSE;
    }

    return USYS_TRUE;
}
