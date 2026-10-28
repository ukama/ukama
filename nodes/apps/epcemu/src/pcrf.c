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

static const char *read_string(JsonObj *root, const char *name, size_t size) {

    JsonObj *value;
    const char *text;
    size_t len;

    value = json_object_get(root, name);
    if (!json_is_string(value)) return NULL;
    text = json_string_value(value);
    len = json_string_length(value);
    if (len >= size || strlen(text) != len) return NULL;
    return text;
}

static int read_epc(JsonObj *root, PcrfEpc *epc) {

    const char *id;
    const char *state;

    id = read_string(root, "epc_instance_id", sizeof(epc->instanceId));
    state = read_string(root, "state", sizeof(epc->state));
    if (id == NULL || state == NULL ||
        (*id != '\0' && !identity_valid(id))) return USYS_FALSE;
    if (strcmp(state, "ready") && strcmp(state, "restart_required")) {
        return USYS_FALSE;
    }
    snprintf(epc->instanceId, sizeof(epc->instanceId), "%s", id);
    snprintf(epc->state, sizeof(epc->state), "%s", state);
    return USYS_TRUE;
}

int pcrf_get_epc(EpcemuConfig *config, PcrfEpc *epc) {

    char url[EPCEMU_MAX_STR * 2];
    JsonObj *reply;
    long code;
    int ret;

    reply = NULL;
    code = 0;
    snprintf(url, sizeof(url), "%s/v1/epc", config->pcrfUrl);
    ret = http_get_json(url, &reply, &code);
    if (!ret || (code == 200 && !read_epc(reply, epc))) code = 0;
    if (reply != NULL) json_decref(reply);
    return (int)code;
}

int pcrf_register_epc(EpcemuConfig *config, const char *id,
                      const char *expected) {

    char url[EPCEMU_MAX_STR * 2];
    JsonObj *body;
    JsonObj *reply;
    PcrfEpc epc;
    long code;
    int ret;

    body = json_pack("{s:s,s:s}", "epc_instance_id", id,
                     "expected_instance_id", expected);
    if (body == NULL) return 0;
    reply = NULL;
    code = 0;
    snprintf(url, sizeof(url), "%s/v1/epc", config->pcrfUrl);
    ret = http_send_json("POST", url, body, &reply, &code);
    if (!ret || (code == 202 &&
        (!read_epc(reply, &epc) || strcmp(epc.instanceId, id) ||
         strcmp(epc.state, "ready")))) code = 0;
    if (reply != NULL) json_decref(reply);
    json_decref(body);
    return (int)code;
}

static JsonObj *session_body(const UeEntry *ue, bool create) {

    JsonObj *body;
    JsonObj *imsi;
    uint32_t ip;

    imsi = imsi_to_json_array(ue->imsi);
    if (imsi == NULL) return NULL;
    body = json_pack("{s:s,s:s,s:o}",
                     "epc_instance_id", ue->instanceId,
                     "pdn_id", ue->pdnId, "imsi", imsi);
    if (body == NULL) return NULL;
    if (create) {
        if (!ip_to_uint32(ue->ip, &ip) ||
            json_object_set_new(body, "apn_name", json_string(ue->apn)) ||
            json_object_set_new(body, "pdn_address", json_integer(ip))) {
            json_decref(body);
            return NULL;
        }
    }
    return body;
}

static int matches(JsonObj *root, const char *name, const char *expected) {

    const char *value;

    value = read_string(root, name, EPCEMU_MAX_STR);
    return value != NULL && !strcmp(value, expected);
}

static int read_pdn(JsonObj *root, const UeEntry *ue, PcrfPdn *pdn) {

    const char *state;
    const char *data;
    const char *reason;

    if (!matches(root, "epc_instance_id", ue->instanceId) ||
        !matches(root, "pdn_id", ue->pdnId) ||
        !matches(root, "imsi", ue->imsi) ||
        !matches(root, "ip", ue->ip) ||
        !matches(root, "apn_name", ue->apn)) return USYS_FALSE;

    state = read_string(root, "state", sizeof(pdn->state));
    data = read_string(root, "data_state", sizeof(pdn->dataState));
    reason = read_string(root, "reason", sizeof(pdn->reason));
    if (state == NULL || data == NULL || reason == NULL || *reason == '\0') {
        return USYS_FALSE;
    }
    if (strcmp(state, "live") && strcmp(state, "closing")) {
        return USYS_FALSE;
    }
    if (strcmp(data, "allowed") && strcmp(data, "blocked") &&
        strcmp(data, "unknown")) return USYS_FALSE;
    snprintf(pdn->state, sizeof(pdn->state), "%s", state);
    snprintf(pdn->dataState, sizeof(pdn->dataState), "%s", data);
    snprintf(pdn->reason, sizeof(pdn->reason), "%s", reason);
    return USYS_TRUE;
}

static int pdn_request(EpcemuConfig *config, const UeEntry *ue,
                       bool create, PcrfPdn *pdn) {

    char url[EPCEMU_MAX_STR * 2];
    JsonObj *body;
    JsonObj *reply;
    const char *error;
    long code;
    int ret;

    memset(pdn, 0, sizeof(*pdn));
    snprintf(pdn->reason, sizeof(pdn->reason), "pcrf_unavailable");
    body = NULL;
    reply = NULL;
    code = 0;
    if (create) {
        body = session_body(ue, true);
        if (body == NULL) return 0;
        snprintf(url, sizeof(url), "%s/v1/session", config->pcrfUrl);
    } else {
        snprintf(url, sizeof(url), "%s/v1/pdn/imsi/%s",
                 config->pcrfUrl, ue->imsi);
    }
    ret = http_send_json(create ? "POST" : "GET", url, body, &reply,
                          &code);
    if (!ret) {
        code = 0;
    } else if (code == (create ? 202 : 200)) {
        if (!read_pdn(reply, ue, pdn)) {
            code = 0;
            snprintf(pdn->reason, sizeof(pdn->reason),
                     "invalid_pcrf_response");
        }
    } else {
        error = read_string(reply, "error", sizeof(pdn->reason));
        if (error != NULL && *error != '\0') {
            snprintf(pdn->reason, sizeof(pdn->reason), "%s", error);
        }
    }
    if (reply != NULL) json_decref(reply);
    if (body != NULL) json_decref(body);
    return (int)code;
}

int pcrf_create_session(EpcemuConfig *config, const UeEntry *ue,
                        PcrfPdn *pdn) {

    return pdn_request(config, ue, true, pdn);
}

int pcrf_get_pdn(EpcemuConfig *config, const UeEntry *ue, PcrfPdn *pdn) {

    return pdn_request(config, ue, false, pdn);
}

int pcrf_delete_session(EpcemuConfig *config, const UeEntry *ue) {

    char url[EPCEMU_MAX_STR * 2];
    JsonObj *body;
    JsonObj *reply;
    long code;
    int ret;

    body = session_body(ue, false);
    if (body == NULL) return 0;
    reply = NULL;
    code = 0;
    snprintf(url, sizeof(url), "%s/v1/session", config->pcrfUrl);
    ret = http_send_json("DELETE", url, body, &reply, &code);
    if (!ret) code = 0;
    if (reply != NULL) json_decref(reply);
    json_decref(body);
    return (int)code;
}
