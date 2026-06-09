/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdlib.h>

#include "web_service.h"
#include "http_status.h"
#include "ops.h"
#include "version.h"

#define METHOD_GET        "GET"
#define METHOD_POST       "POST"
#define JSON_FIELD_ERR    "error"
#define JSON_FIELD_TILT   "targetTiltDeg"
#define JSON_FIELD_PROF   "profile"
#define JSON_FIELD_CFG    "configPath"

typedef bool (*AisgdSimpleOp)(AisgdContext *ctx, JsonObj **response);

typedef struct {
    const char       *method;
    const char       *path;
    int (*cb)(const URequest *, UResponse *, void *);
} AisgdEndpoint;

static void set_json_response(UResponse *response,
                              int status,
                              JsonObj *json) {

    if (json == NULL) {
        ulfius_set_string_body_response(
            response,
            HttpStatus_InternalServerError,
            HttpStatusStr(HttpStatus_InternalServerError));
        return;
    }

    ulfius_set_json_body_response(response, status, json);
    json_decref(json);
}

static void set_error_response(UResponse *response,
                               int status,
                               const char *message) {

    JsonObj *json = NULL;

    json = json_object();
    if (json == NULL) {
        ulfius_set_string_body_response(response,
                                        status,
                                        HttpStatusStr(status));
        return;
    }

    json_object_set_new(json,
                        JSON_FIELD_ERR,
                        json_string(message ? message : HttpStatusStr(status)));

    set_json_response(response, status, json);
}

static bool parse_json_body(const URequest *request, JsonObj **json) {

    if (json == NULL) {
        return false;
    }

    if (request->binary_body == NULL || request->binary_body_length == 0) {
        *json = json_object();
        return *json != NULL;
    }

    *json = ulfius_get_json_body_request(request, NULL);

    return *json != NULL;
}

static const char *json_get_string(JsonObj *json, const char *key) {

    JsonObj *value = NULL;

    value = json_object_get(json, key);
    if (!json_is_string(value)) {
        return "";
    }

    return json_string_value(value);
}

static int run_simple_op(UResponse *response,
                         AisgdContext *ctx,
                         AisgdSimpleOp op) {

    JsonObj *json = NULL;

    if (ctx == NULL || op == NULL) {
        set_error_response(response,
                           HttpStatus_InternalServerError,
                           "invalid request context");
        return U_CALLBACK_CONTINUE;
    }

    if (!op(ctx, &json)) {
        set_error_response(response,
                           HttpStatus_ServiceUnavailable,
                           "operation failed");
        return U_CALLBACK_CONTINUE;
    }

    set_json_response(response, HttpStatus_OK, json);

    return U_CALLBACK_CONTINUE;
}

static bool add_endpoint(UInst *instance,
                         const AisgdEndpoint *endpoint,
                         AisgdContext *ctx) {

    int ret;

    ret = ulfius_add_endpoint_by_val(instance,
                                     endpoint->method,
                                     URL_PREFIX,
                                     endpoint->path,
                                     0,
                                     endpoint->cb,
                                     ctx);

    return ret == U_OK;
}

int web_service_cb_ping(const URequest *request,
                        UResponse *response,
                        void *data) {

    AisgdContext *ctx = NULL;

    (void)request;

    ctx = (AisgdContext *)data;
    if (ctx == NULL || ctx->status == NULL) {
        ulfius_set_string_body_response(response,
                                        HttpStatus_ServiceUnavailable,
                                        "Service Unavailable");
        return U_CALLBACK_CONTINUE;
    }

    if (!status_is_ready(ctx->status)) {
        ulfius_set_string_body_response(response,
                                        HttpStatus_ServiceUnavailable,
                                        "Service Unavailable");
        return U_CALLBACK_CONTINUE;
    }

    ulfius_set_string_body_response(response, HttpStatus_OK, "OK");

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

    AisgdContext *ctx = NULL;

    (void)request;

    ctx = (AisgdContext *)data;
    if (ctx == NULL || ctx->status == NULL) {
        set_error_response(response,
                           HttpStatus_InternalServerError,
                           "invalid service context");
        return U_CALLBACK_CONTINUE;
    }

    set_json_response(response, HttpStatus_OK, status_to_json(ctx->status));

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_reconcile(const URequest *request,
                             UResponse *response,
                             void *data) {

    JsonObj *json = NULL;

    (void)request;

    if (!aisgd_ops_reconcile((AisgdContext *)data, &json)) {
        set_error_response(response,
                           HttpStatus_ServiceUnavailable,
                           "reconcile failed");
        return U_CALLBACK_CONTINUE;
    }

    set_json_response(response, HttpStatus_OK, json);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_scan(const URequest *request,
                        UResponse *response,
                        void *data) {
    (void)request;

    return run_simple_op(response, data, aisgd_ops_scan);
}

int web_service_cb_device(const URequest *request,
                          UResponse *response,
                          void *data) {
    (void)request;

    return run_simple_op(response, data, aisgd_ops_get_device);
}

int web_service_cb_info(const URequest *request,
                        UResponse *response,
                        void *data) {
    (void)request;

    return run_simple_op(response, data, aisgd_ops_get_info);
}

int web_service_cb_alarms(const URequest *request,
                          UResponse *response,
                          void *data) {
    (void)request;

    return run_simple_op(response, data, aisgd_ops_get_alarms);
}

int web_service_cb_clear_alarms(const URequest *request,
                                UResponse *response,
                                void *data) {
    (void)request;

    return run_simple_op(response, data, aisgd_ops_clear_alarms);
}

int web_service_cb_subscribe_alarms(const URequest *request,
                                    UResponse *response,
                                    void *data) {
    (void)request;

    return run_simple_op(response, data, aisgd_ops_subscribe_alarms);
}

int web_service_cb_self_test(const URequest *request,
                             UResponse *response,
                             void *data) {
    (void)request;

    return run_simple_op(response, data, aisgd_ops_self_test);
}

int web_service_cb_calibrate(const URequest *request,
                             UResponse *response,
                             void *data) {
    (void)request;

    return run_simple_op(response, data, aisgd_ops_calibrate);
}

int web_service_cb_get_tilt(const URequest *request,
                            UResponse *response,
                            void *data) {
    (void)request;

    return run_simple_op(response, data, aisgd_ops_get_tilt);
}

int web_service_cb_reset(const URequest *request,
                         UResponse *response,
                         void *data) {
    (void)request;

    return run_simple_op(response, data, aisgd_ops_reset);
}

int web_service_cb_configure(const URequest *request,
                             UResponse *response,
                             void *data) {

    JsonObj *body       = NULL;
    JsonObj *json       = NULL;
    const char *profile = NULL;
    const char *path    = NULL;

    if (!parse_json_body(request, &body)) {
        set_error_response(response,
                           HttpStatus_BadRequest,
                           "invalid json body");
        return U_CALLBACK_CONTINUE;
    }

    profile = json_get_string(body, JSON_FIELD_PROF);
    path    = json_get_string(body, JSON_FIELD_CFG);

    if (!aisgd_ops_configure((AisgdContext *)data, profile, path, &json)) {
        json_decref(body);
        set_error_response(response,
                           HttpStatus_ServiceUnavailable,
                           "configure failed");
        return U_CALLBACK_CONTINUE;
    }

    json_decref(body);
    set_json_response(response, HttpStatus_OK, json);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_set_tilt(const URequest *request,
                            UResponse *response,
                            void *data) {

    JsonObj *body  = NULL;
    JsonObj *json  = NULL;
    JsonObj *value = NULL;

    if (!parse_json_body(request, &body)) {
        set_error_response(response,
                           HttpStatus_BadRequest,
                           "invalid json body");
        return U_CALLBACK_CONTINUE;
    }

    value = json_object_get(body, JSON_FIELD_TILT);
    if (!json_is_number(value)) {
        json_decref(body);
        set_error_response(response,
                           HttpStatus_BadRequest,
                           "missing targetTiltDeg");
        return U_CALLBACK_CONTINUE;
    }

    if (!aisgd_ops_set_tilt((AisgdContext *)data,
                            json_number_value(value),
                            &json)) {
        json_decref(body);
        set_error_response(response,
                           HttpStatus_ServiceUnavailable,
                           "set tilt failed");
        return U_CALLBACK_CONTINUE;
    }

    json_decref(body);
    set_json_response(response, HttpStatus_OK, json);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_get_device_data(const URequest *request,
                                   UResponse *response,
                                   void *data) {

    JsonObj *json     = NULL;
    const char *field = NULL;

    field = u_map_get(request->map_url, "field");
    if (field == NULL) {
        set_error_response(response, HttpStatus_BadRequest, "missing field");
        return U_CALLBACK_CONTINUE;
    }

    if (!aisgd_ops_get_device_data((AisgdContext *)data,
                                   atoi(field),
                                   &json)) {
        set_error_response(response,
                           HttpStatus_ServiceUnavailable,
                           "get device data failed");
        return U_CALLBACK_CONTINUE;
    }

    set_json_response(response, HttpStatus_OK, json);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_default(const URequest *request,
                           UResponse *response,
                           void *data) {

    (void)request;
    (void)data;

    set_error_response(response, HttpStatus_NotFound, "resource not found");

    return U_CALLBACK_CONTINUE;
}

int start_web_service(AisgdContext *ctx, UInst *serviceInst) {

    size_t i;

    static const AisgdEndpoint endpoints[] = {
        { METHOD_GET,  "ping",                    web_service_cb_ping },
        { METHOD_GET,  "version",                 web_service_cb_version },
        { METHOD_GET,  "status",                  web_service_cb_status },
        { METHOD_POST, "reconcile",               web_service_cb_reconcile },

        { METHOD_POST, "device/scan",             web_service_cb_scan },
        { METHOD_GET,  "device",                  web_service_cb_device },
        { METHOD_GET,  "device/info",             web_service_cb_info },
        { METHOD_GET,  "device/alarms",           web_service_cb_alarms },
        { METHOD_POST, "device/alarms/clear",     web_service_cb_clear_alarms },
        { METHOD_POST, "device/alarms/subscribe",
          web_service_cb_subscribe_alarms },

        { METHOD_POST, "device/self-test",        web_service_cb_self_test },
        { METHOD_POST, "device/config",           web_service_cb_configure },
        { METHOD_POST, "device/calibrate",        web_service_cb_calibrate },
        { METHOD_GET,  "device/tilt",             web_service_cb_get_tilt },
        { METHOD_POST, "device/tilt",             web_service_cb_set_tilt },
        { METHOD_GET,  "device/data/:field",
          web_service_cb_get_device_data },
        { METHOD_POST, "device/reset",            web_service_cb_reset },
    };

    if (ctx == NULL || serviceInst == NULL) {
        return USYS_FALSE;
    }

    if (ulfius_init_instance(serviceInst,
                             ctx->config->servicePort,
                             NULL,
                             NULL) != U_OK) {
        return USYS_FALSE;
    }

    for (i = 0; i < sizeof(endpoints) / sizeof(endpoints[0]); i++) {
        if (!add_endpoint(serviceInst, &endpoints[i], ctx)) {
            return USYS_FALSE;
        }
    }

    ulfius_set_default_endpoint(serviceInst, web_service_cb_default, ctx);

    if (ulfius_start_framework(serviceInst) != U_OK) {
        return USYS_FALSE;
    }

    return USYS_TRUE;
}
