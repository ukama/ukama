/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdbool.h>
#include <stdio.h>
#include <string.h>

#include "epcemu.h"
#include "http_client.h"
#include "netutil.h"
#include "pcrf.h"

bool pcrf_is_ready(EpcemuConfig *config) {

    char url[EPCEMU_MAX_STR * 2];
    JsonObj *root;
    JsonObj *ready;
    long code;
    bool ret;

    if (config == NULL) return false;

    root = NULL;
    code = 0;
    ret  = false;

    snprintf(url, sizeof(url), "%s/v1/status", config->pcrfUrl);

    if (!http_get_json(url, &root, &code)) {
        config->pcrfReady = false;
        return false;
    }

    if (code == 200 && root != NULL) {
        ready = json_object_get(root, "ready");
        ret = json_is_true(ready);
    }

    if (root != NULL) json_decref(root);

    config->pcrfReady = ret;
    return ret;
}

int pcrf_probe(EpcemuConfig *config, EpcemuStatus *status) {

    if (config == NULL || status == NULL) return USYS_FALSE;

    status_set(status, EpcemuStateCheckingPcrf, "checking pcrf");

    if (!pcrf_is_ready(config)) {
        usys_log_error("pcrf is not ready");
        return USYS_FALSE;
    }

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

    if (!pcrf_is_ready(config)) {
        usys_log_error("PCRF not ready for session create imsi=%s ip=%s",
                       imsi, ip);
        return USYS_FALSE;
    }

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
        config->pcrfReady = false;
        usys_log_error("PCRF create session failed imsi=%s ip=%s http=%ld",
                       imsi, ip, code);
        return USYS_FALSE;
    }

    config->pcrfReady = true;
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

    if (!pcrf_is_ready(config)) {
        usys_log_error("PCRF not ready for session delete imsi=%s", imsi);
        return USYS_FALSE;
    }

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
        config->pcrfReady = false;
        usys_log_error("PCRF delete session failed imsi=%s http=%ld",
                       imsi, code);
        return USYS_FALSE;
    }

    config->pcrfReady = true;
    return USYS_TRUE;
}
