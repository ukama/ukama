/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "service_api.h"
#include "web_service.h"
#include "http_status.h"
#include "version.h"
#include "ovs.h"

#include "usys_log.h"

int web_service_cb_ping(const URequest *request,
                        UResponse *response,
                        void *data) {

    (void)request;
    (void)data;

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
    InitState state;
    char reason[STATUS_REASON_LEN];
    int status;

    (void)request;

    ctx = (ServiceContext *)data;
    if (ctx == NULL || ctx->status == NULL) {
        ulfius_set_string_body_response(
            response,
            HttpStatus_InternalServerError,
            HttpStatusStr(HttpStatus_InternalServerError));
        return U_CALLBACK_CONTINUE;
    }

    state = status_get(ctx->status, reason, sizeof(reason));
    status = HttpStatus_Accepted;
    if (state == InitStateReady) {
        status = HttpStatus_OK;
    } else if (state == InitStateFailed) {
        status = HttpStatus_ServiceUnavailable;
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
                        json_boolean(state == InitStateReady));
    if (state != InitStateReady) {
        json_object_set_new(json, "reason", json_string(reason));
    }

    ulfius_set_json_body_response(response, status, json);
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
    if (ctx == NULL || ctx->config == NULL || ctx->status == NULL) {
        ulfius_set_string_body_response(response,
                                        HttpStatus_ServiceUnavailable,
                                        HttpStatusStr(HttpStatus_ServiceUnavailable));
        return U_CALLBACK_CONTINUE;
    }

    json = status_to_json(ctx->status, ctx->config);
    if (json == NULL) {
        ulfius_set_string_body_response(response,
                                        HttpStatus_InternalServerError,
                                        HttpStatusStr(HttpStatus_InternalServerError));
        return U_CALLBACK_CONTINUE;
    }

    ulfius_set_json_body_response(response, HttpStatus_OK, json);
    json_decref(json);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_reconcile(const URequest *request,
                             UResponse *response,
                             void *data) {

    ServiceContext *ctx;
    JsonObj *json;

    (void)request;

    ctx = (ServiceContext *)data;
    if (ctx == NULL || ctx->config == NULL || ctx->status == NULL) {
        ulfius_set_string_body_response(response,
                                        HttpStatus_ServiceUnavailable,
                                        HttpStatusStr(HttpStatus_ServiceUnavailable));
        return U_CALLBACK_CONTINUE;
    }

    if (!ovs_reconcile(ctx->config, ctx->status)) {
        json = status_to_json(ctx->status, ctx->config);
        if (json != NULL) {
            ulfius_set_json_body_response(response,
                                         HttpStatus_ServiceUnavailable,
                                         json);
            json_decref(json);
        } else {
            ulfius_set_string_body_response(response,
                                            HttpStatus_ServiceUnavailable,
                                            HttpStatusStr(HttpStatus_ServiceUnavailable));
        }
        return U_CALLBACK_CONTINUE;
    }

    json = status_to_json(ctx->status, ctx->config);
    if (json == NULL) {
        ulfius_set_string_body_response(response,
                                        HttpStatus_InternalServerError,
                                        HttpStatusStr(HttpStatus_InternalServerError));
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

    ulfius_set_string_body_response(response,
                                    HttpStatus_NotFound,
                                    HttpStatusStr(HttpStatus_NotFound));

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_not_allowed(const URequest *request,
                               UResponse *response,
                               void *data) {

    (void)request;
    (void)data;

    ulfius_set_string_body_response(response,
                                    HttpStatus_MethodNotAllowed,
                                    HttpStatusStr(HttpStatus_MethodNotAllowed));

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_service(const URequest *request, UResponse *response,
                           void *data) {

    ServiceContext *ctx = data;
    JsonObj *body;
    JsonObj *reply;
    const char *state;
    int code = HttpStatus_OK;

    if (!strcmp(request->http_verb, HTTP_POST)) {
        body = ulfius_get_json_body_request(request, NULL);
        state = json_string_value(json_object_get(body, "state"));

        if (state == NULL || (strcmp(state, "on") && strcmp(state, "off"))) {
            json_decref(body);
            ulfius_set_empty_body_response(response, HttpStatus_BadRequest);

            return U_CALLBACK_CONTINUE;
        }

        if (!ovs_service_set(ctx->config, !strcmp(state, "on"))) {
            code = HttpStatus_ServiceUnavailable;
        }

        json_decref(body);
    }

    reply = ovs_service_status(ctx->config);
    ulfius_set_json_body_response(response, code, reply);
    json_decref(reply);

    return U_CALLBACK_CONTINUE;
}
