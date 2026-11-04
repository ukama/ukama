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
#include "service_lease.h"
#include "http_status.h"
#include "netutil.h"
#include "pdn.h"
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

static JsonObj *status_json(ServiceContext *ctx) {

    JsonObj *root;
    JsonObj *pcrf;
    JsonObj *initNetwork;
    JsonObj *service;
    JsonObj *ues;
    JsonObj *userPlane;
    EpcemuState state;
    bool ready;
    bool serviceOn;
    bool dataPlaneStarted;
    char reason[EPCEMU_MAX_REASON];

    if (ctx == NULL || ctx->config == NULL || ctx->status == NULL) return NULL;

    pthread_mutex_lock(&ctx->status->mutex);
    state = ctx->status->state;
    ready = ctx->status->ready;
    snprintf(reason, sizeof(reason), "%s", ctx->status->reason);
    pthread_mutex_unlock(&ctx->status->mutex);

    pthread_mutex_lock(&ctx->serviceMutex);
    serviceOn = ctx->serviceOn;
    dataPlaneStarted = ctx->dataPlaneStarted;

    root = json_object();
    if (root == NULL) {
        pthread_mutex_unlock(&ctx->serviceMutex);
        return NULL;
    }

    pcrf = json_object();
    initNetwork = json_object();
    service = json_object();
    ues = ue_summary_json();
    userPlane = dataPlaneStarted ?
        data_plane_json(&gDataPlane, ctx->config) : json_object();

    if (userPlane != NULL && !dataPlaneStarted) {
        json_object_set_new(userPlane, "enabled", json_false());
    }

    json_object_set_new(root, "ready",  json_boolean(ready));
    json_object_set_new(root, "state",  json_string(status_state_str(state)));
    json_object_set_new(root, "reason", json_string(reason));

    json_object_set_new(service, "state",
                        json_string(serviceOn ? "on" : "off"));
    json_object_set_new(service, "admission",
                        json_string(!serviceOn || !service_lease_valid() ?
                                    "disabled" :
                                    ready && ctx->config->pcrfReady ?
                                    "enabled" : "pending"));
    json_object_set_new(service, "pdn_count", json_integer(ue_count_total()));
    json_object_set_new(service, "cleanup_pending",
                        json_boolean(!serviceOn && ue_count_total() != 0));
    json_object_set_new(service, "confirmed",
                        json_boolean(serviceOn ? service_lease_valid()
                                               : ue_count_total() == 0));
    json_object_set_new(root, "service", service);

    json_object_set_new(pcrf, "url",   json_string(ctx->config->pcrfUrl));
    json_object_set_new(pcrf, "ready", json_boolean(ctx->config->pcrfReady));
    json_object_set_new(pcrf, "registered",
                        json_boolean(ctx->pcrfRegistered));
    json_object_set_new(pcrf, "epc_instance_id",
                        json_string(ctx->instanceId));
    json_object_set_new(pcrf, "restartRequired",
                        json_boolean(ctx->ownerLost));
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

    pthread_mutex_unlock(&ctx->serviceMutex);
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

int web_service_cb_ready(const URequest *request,
                         UResponse *response,
                         void *data) {

    ServiceContext *ctx;
    JsonObj *json;
    EpcemuState state;
    char reason[EPCEMU_MAX_REASON];
    int httpStatus;

    (void)request;

    ctx = (ServiceContext *)data;
    if (ctx == NULL || ctx->status == NULL) {
        ulfius_set_string_body_response(
            response,
            HttpStatus_InternalServerError,
            HttpStatusStr(HttpStatus_InternalServerError));
        return U_CALLBACK_CONTINUE;
    }

    pthread_mutex_lock(&ctx->status->mutex);
    state = ctx->status->state;
    snprintf(reason, sizeof(reason), "%s", ctx->status->reason);
    pthread_mutex_unlock(&ctx->status->mutex);

    httpStatus = HttpStatus_Accepted;
    if (state == EpcemuStateReady) {
        httpStatus = HttpStatus_OK;
    } else if (state == EpcemuStateFailed) {
        httpStatus = HttpStatus_ServiceUnavailable;
    }

    json = json_object();
    if (json == NULL) {
        ulfius_set_string_body_response(
            response,
            HttpStatus_InternalServerError,
            HttpStatusStr(HttpStatus_InternalServerError));
        return U_CALLBACK_CONTINUE;
    }

    json_object_set_new(json,
                        "ready",
                        json_boolean(state == EpcemuStateReady));
    if (state != EpcemuStateReady) {
        json_object_set_new(json, "reason", json_string(reason));
    }

    ulfius_set_json_body_response(response, httpStatus, json);
    json_decref(json);
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

int web_service_cb_service(const URequest *request,
                           UResponse *response,
                           void *data) {

    ServiceContext *ctx = data;
    JsonObj *body;
    JsonObj *reply;
    const char *state;
    int code;

    if (ctx == NULL || ctx->config == NULL || ctx->status == NULL) {
        set_error(response, HttpStatus_ServiceUnavailable,
                  "epcemu unavailable");
        return U_CALLBACK_CONTINUE;
    }

    body = ulfius_get_json_body_request(request, NULL);
    state = json_string_default(body, "state", NULL);
    if (state == NULL ||
        (strcmp(state, "on") != 0 && strcmp(state, "off") != 0)) {
        json_decref(body);
        set_error(response, HttpStatus_BadRequest, "state must be on or off");
        return U_CALLBACK_CONTINUE;
    }

    code = pdn_set_service(ctx, strcmp(state, "on") == 0);

    json_decref(body);

    reply = status_json(ctx);
    if (reply == NULL) {
        set_error(response, HttpStatus_InternalServerError,
                  "failed to build service response");
        return U_CALLBACK_CONTINUE;
    }

    ulfius_set_json_body_response(response, code,
                                  json_object_get(reply, "service"));
    json_decref(reply);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_attach(const URequest *request,
                          UResponse *response,
                          void *data) {

    ServiceContext *ctx;
    JsonObj *body;
    JsonObj *reply;
    char reason[EPCEMU_MAX_REASON];
    int code;

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
    snprintf(reason, sizeof(reason), "pdn_registration_failed");
    reply = NULL;
    code = pdn_attach(ctx, body, &reply, reason, sizeof(reason));
    json_decref(body);
    if (reply != NULL) {
        ulfius_set_json_body_response(response, code, reply);
        json_decref(reply);
    } else {
        set_error(response, code, reason);
    }
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_detach(const URequest *request,
                          UResponse *response,
                          void *data) {

    ServiceContext *ctx;
    JsonObj *body;
    char reason[EPCEMU_MAX_REASON];
    int code;

    ctx = (ServiceContext *)data;
    if (ctx == NULL || ctx->config == NULL) {
        set_error(response, HttpStatus_ServiceUnavailable, "epcemu not ready");
        return U_CALLBACK_CONTINUE;
    }
    body = ulfius_get_json_body_request(request, NULL);
    if (body == NULL) {
        set_error(response, HttpStatus_BadRequest, "invalid JSON body");
        return U_CALLBACK_CONTINUE;
    }
    code = pdn_detach(ctx, body, reason, sizeof(reason));
    json_decref(body);
    if (code == HttpStatus_OK) {
        ulfius_set_string_body_response(response, code, "OK");
    } else {
        set_error(response, code, reason);
    }
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
