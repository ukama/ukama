/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "web_service.h"
#include "http_status.h"
#include "version.h"

#include "usys_log.h"

int web_service_cb_ping(const URequest *request,
                        UResponse *response,
                        void *data) {

    ServiceContext *ctx;

    (void)request;

    ctx = (ServiceContext *)data;
    if (ctx == NULL || ctx->status == NULL) {
        ulfius_set_string_body_response(response,
                                        HttpStatus_ServiceUnavailable,
                                        HttpStatusStr(HttpStatus_ServiceUnavailable));
        return U_CALLBACK_CONTINUE;
    }

    if (status_is_ready(ctx->status)) {
        ulfius_set_string_body_response(response,
                                        HttpStatus_OK,
                                        HttpStatusStr(HttpStatus_OK));
    } else {
        ulfius_set_string_body_response(response,
                                        HttpStatus_ServiceUnavailable,
                                        HttpStatusStr(HttpStatus_ServiceUnavailable));
    }

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
