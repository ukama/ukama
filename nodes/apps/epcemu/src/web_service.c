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

#include "data_plane.h"
#include "epcemu.h"
#include "http_status.h"
#include "netutil.h"
#include "pcrf.h"
#include "ue.h"
#include "web_service.h"
#include "version.h"

extern DataPlane gDataPlane;

static void set_error(UResponse *response, int code, const char *message) {

    JsonObj *obj;

    obj = json_object();
    if (obj == NULL) {
        ulfius_set_string_body_response(response, code,
                                        message ? message : HttpStatusStr(code));
        return;
    }

    json_object_set_new(obj, "error", json_string(message ? message : "error"));
    ulfius_set_json_body_response(response, code, obj);
    json_decref(obj);
}

static const char *json_string_default(JsonObj *obj,
                                       const char *key,
                                       const char *defValue) {

    JsonObj *value;

    if (obj == NULL || key == NULL) return defValue;

    value = json_object_get(obj, key);
    if (json_is_string(value)) return json_string_value(value);

    return defValue;
}

static int json_int_default(JsonObj *obj, const char *key, int defValue) {

    JsonObj *value;

    if (obj == NULL || key == NULL) return defValue;

    value = json_object_get(obj, key);
    if (json_is_integer(value)) return (int)json_integer_value(value);

    return defValue;
}

static JsonObj *status_json(ServiceContext *ctx) {

    JsonObj *root;
    JsonObj *pcrf;
    JsonObj *initNetwork;
    JsonObj *ues;
    JsonObj *userPlane;
    EpcemuState state;
    bool ready;
    char reason[EPCEMU_MAX_REASON];

    if (ctx == NULL || ctx->config == NULL || ctx->status == NULL) return NULL;

    pthread_mutex_lock(&ctx->status->mutex);
    state = ctx->status->state;
    ready = ctx->status->ready;
    snprintf(reason, sizeof(reason), "%s", ctx->status->reason);
    pthread_mutex_unlock(&ctx->status->mutex);

    root = json_object();
    if (root == NULL) return NULL;

    (void)pcrf_is_ready(ctx->config);

    pcrf = json_object();
    initNetwork = json_object();
    ues = ue_summary_json();
    userPlane = data_plane_json(&gDataPlane, ctx->config);

    json_object_set_new(root, "ready",  json_boolean(ready));
    json_object_set_new(root, "state",  json_string(status_state_str(state)));
    json_object_set_new(root, "reason", json_string(reason));

    json_object_set_new(pcrf, "url",   json_string(ctx->config->pcrfUrl));
    json_object_set_new(pcrf, "ready", json_boolean(ctx->config->pcrfReady));
    json_object_set_new(root, "pcrf",  pcrf);

    json_object_set_new(initNetwork, "url",
                        json_string(ctx->config->initNetworkUrl));
    json_object_set_new(initNetwork, "ready",
                        json_boolean(ctx->config->initNetworkReady));
    json_object_set_new(initNetwork, "routed",
                        json_boolean(ctx->config->initNetworkRouted));
    json_object_set_new(initNetwork, "bridge",
                        json_string(ctx->config->bridge));
    json_object_set_new(initNetwork, "bridgeCidr",
                        json_string(ctx->config->bridgeCidr));
    json_object_set_new(initNetwork, "ueCidr",
                        json_string(ctx->config->ueCidr));
    json_object_set_new(root, "initNetwork", initNetwork);

    if (ues != NULL) {
        json_object_set_new(root, "ues", ues);
    }

    if (userPlane != NULL) {
        json_object_set_new(root, "userPlane", userPlane);
    }

    return root;
}

int web_service_cb_ping(const URequest *request,
                        UResponse *response,
                        void *data) {

    ServiceContext *ctx;

    (void)request;

    ctx = (ServiceContext *)data;
    if (ctx == NULL || ctx->status == NULL || !status_is_ready(ctx->status)) {
        ulfius_set_string_body_response(response,
                                        HttpStatus_ServiceUnavailable,
                                        "not ready");
        return U_CALLBACK_CONTINUE;
    }

    ulfius_set_string_body_response(response,
                                    HttpStatus_OK,
                                    HttpStatusStr(HttpStatus_OK));
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_version(const URequest *request,
                           UResponse *response,
                           void *data) {

    (void)request;
    (void)data;

    ulfius_set_string_body_response(response, HttpStatus_OK, VERSION);
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_status(const URequest *request,
                          UResponse *response,
                          void *data) {

    ServiceContext *ctx;
    JsonObj *json;

    (void)request;

    ctx = (ServiceContext *)data;
    json = status_json(ctx);
    if (json == NULL) {
        set_error(response, HttpStatus_InternalServerError,
                  "failed to build status");
        return U_CALLBACK_CONTINUE;
    }

    ulfius_set_json_body_response(response, HttpStatus_OK, json);
    json_decref(json);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_attach(const URequest *request,
                          UResponse *response,
                          void *data) {

    ServiceContext *ctx;
    JsonObj *body;
    JsonObj *reply;
    const char *imsi;
    const char *iccid;
    const char *ip;
    const char *apn;
    int userPlanePort;
    char imsiBuf[UE_IMSI_LEN];
    char iccidBuf[UE_ICCID_LEN];
    char ipBuf[UE_IP_LEN];
    char apnBuf[UE_APN_LEN];
    char reason[EPCEMU_MAX_REASON];
    UeEntry existing;

    ctx = (ServiceContext *)data;
    body = NULL;
    reply = NULL;

    memset(reason, 0, sizeof(reason));
    memset(imsiBuf, 0, sizeof(imsiBuf));
    memset(iccidBuf, 0, sizeof(iccidBuf));
    memset(ipBuf, 0, sizeof(ipBuf));
    memset(apnBuf, 0, sizeof(apnBuf));
    memset(&existing, 0, sizeof(existing));

    if (ctx == NULL || ctx->config == NULL || !status_is_ready(ctx->status)) {
        set_error(response, HttpStatus_ServiceUnavailable, "epcemu not ready");
        return U_CALLBACK_CONTINUE;
    }

    body = ulfius_get_json_body_request(request, NULL);
    if (body == NULL) {
        set_error(response, HttpStatus_BadRequest, "invalid JSON body");
        return U_CALLBACK_CONTINUE;
    }

    imsi = json_string_default(body, "imsi", NULL);
    iccid = json_string_default(body, "iccid", "");
    ip = json_string_default(body, "ip", NULL);
    apn = json_string_default(body, "apn", EPCEMU_DEF_APN);
    userPlanePort = json_int_default(body, "userPlanePort", 0);

    (void)userPlanePort;

    if (!imsi_valid(imsi)) {
        json_decref(body);
        set_error(response, HttpStatus_BadRequest, "invalid imsi");
        return U_CALLBACK_CONTINUE;
    }

    if (ip == NULL || !ip_in_cidr(ip, ctx->config->ueCidr)) {
        json_decref(body);
        set_error(response, HttpStatus_BadRequest,
                  "ue ip outside configured ue cidr");
        return U_CALLBACK_CONTINUE;
    }

    snprintf(imsiBuf, sizeof(imsiBuf), "%s", imsi);
    snprintf(iccidBuf, sizeof(iccidBuf), "%s", iccid ? iccid : "");
    snprintf(ipBuf, sizeof(ipBuf), "%s", ip);
    snprintf(apnBuf, sizeof(apnBuf), "%s", apn ? apn : EPCEMU_DEF_APN);

    json_decref(body);
    body = NULL;

    /*
     * Attach is idempotent for the same IMSI/IP pair.
     *
     * This is important for lab scripts because the UE agent/container may
     * retry after a previous request committed UE/PCRF state but failed while
     * building the HTTP response.
     */
    if (ue_get(imsiBuf, &existing)) {
        if (existing.state == UeStateAttached &&
            strcmp(existing.ip, ipBuf) == 0) {

            reply = ue_get_json(imsiBuf);
            if (reply == NULL) {
                set_error(response, HttpStatus_InternalServerError,
                          "failed to read attached ue");
                return U_CALLBACK_CONTINUE;
            }

            ulfius_set_json_body_response(response, HttpStatus_OK, reply);
            json_decref(reply);
            return U_CALLBACK_CONTINUE;
        }

        set_error(response, HttpStatus_Conflict, "imsi already attached");
        return U_CALLBACK_CONTINUE;
    }

    if (!ue_attach_start(imsiBuf,
                         iccidBuf,
                         ipBuf,
                         apnBuf,
                         reason,
                         sizeof(reason))) {
        set_error(response, HttpStatus_Conflict, reason);
        return U_CALLBACK_CONTINUE;
    }

    if (!pcrf_is_ready(ctx->config)) {
        ue_attach_fail(imsiBuf, "pcrf not ready");
        set_error(response, HttpStatus_ServiceUnavailable,
                  "pcrf not ready");
        return U_CALLBACK_CONTINUE;
    }

    if (!pcrf_create_session(ctx->config, imsiBuf, ipBuf, apnBuf)) {
        ue_attach_fail(imsiBuf, "pcrf session create failed");
        set_error(response, HttpStatus_ServiceUnavailable,
                  "pcrf session create failed");
        return U_CALLBACK_CONTINUE;
    }

    ue_attach_complete(imsiBuf);

    reply = ue_get_json(imsiBuf);
    if (reply == NULL) {
        set_error(response, HttpStatus_InternalServerError,
                  "failed to read attached ue");
        return U_CALLBACK_CONTINUE;
    }

    ulfius_set_json_body_response(response, HttpStatus_Created, reply);
    json_decref(reply);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_detach(const URequest *request,
                          UResponse *response,
                          void *data) {

    ServiceContext *ctx;
    JsonObj *body;
    const char *imsi;
    UeEntry ue;

    ctx = (ServiceContext *)data;

    if (ctx == NULL || ctx->config == NULL || !status_is_ready(ctx->status)) {
        set_error(response, HttpStatus_ServiceUnavailable, "epcemu not ready");
        return U_CALLBACK_CONTINUE;
    }

    body = ulfius_get_json_body_request(request, NULL);
    if (body == NULL) {
        set_error(response, HttpStatus_BadRequest, "invalid JSON body");
        return U_CALLBACK_CONTINUE;
    }

    imsi = json_string_default(body, "imsi", NULL);
    if (!imsi_valid(imsi)) {
        json_decref(body);
        set_error(response, HttpStatus_BadRequest, "invalid imsi");
        return U_CALLBACK_CONTINUE;
    }

    if (!ue_detach_start(imsi, &ue)) {
        json_decref(body);
        ulfius_set_string_body_response(response, HttpStatus_OK, "OK");
        return U_CALLBACK_CONTINUE;
    }

    if (!pcrf_delete_session(ctx->config, imsi)) {
        usys_log_error("PCRF session delete failed imsi=%s", imsi);
    }

    ue_detach_complete(imsi);
    json_decref(body);

    ulfius_set_string_body_response(response, HttpStatus_OK, "OK");
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_get_ue(const URequest *request,
                          UResponse *response,
                          void *data) {

    const char *imsi;
    JsonObj *json;

    (void)data;

    imsi = u_map_get(request->map_url, "imsi");
    if (!imsi_valid(imsi)) {
        set_error(response, HttpStatus_BadRequest, "invalid imsi");
        return U_CALLBACK_CONTINUE;
    }

    json = ue_get_json(imsi);
    if (json == NULL) {
        set_error(response, HttpStatus_NotFound, "ue not found");
        return U_CALLBACK_CONTINUE;
    }

    ulfius_set_json_body_response(response, HttpStatus_OK, json);
    json_decref(json);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_list_ues(const URequest *request,
                            UResponse *response,
                            void *data) {

    JsonObj *json;

    (void)request;
    (void)data;

    json = ue_list_json();
    if (json == NULL) {
        set_error(response, HttpStatus_InternalServerError,
                  "failed to build ue list");
        return U_CALLBACK_CONTINUE;
    }

    ulfius_set_json_body_response(response, HttpStatus_OK, json);
    json_decref(json);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_default(const URequest *request,
                           UResponse *response,
                           void *data) {

    (void)request;
    (void)data;

    set_error(response, HttpStatus_NotFound, HttpStatusStr(HttpStatus_NotFound));
    return U_CALLBACK_CONTINUE;
}
