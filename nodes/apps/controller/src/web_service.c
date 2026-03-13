/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>

#include <ulfius.h>
#include <jansson.h>

#include "web_service.h"
#include "jserdes.h"
#include "http_status.h"
#include "usys_log.h"
#include "version.h"

static void setup_unsupported_methods(UInst *instance, const char *allowed,
                                      const char *prefix, const char *resource) {
    if (strcmp(allowed, "GET") != 0) {
        ulfius_add_endpoint_by_val(instance, "GET", prefix, resource, 0,
                                   &web_service_cb_not_allowed, (void *)allowed);
    }
    if (strcmp(allowed, "POST") != 0) {
        ulfius_add_endpoint_by_val(instance, "POST", prefix, resource, 0,
                                   &web_service_cb_not_allowed, (void *)allowed);
    }
    if (strcmp(allowed, "PUT") != 0) {
        ulfius_add_endpoint_by_val(instance, "PUT", prefix, resource, 0,
                                   &web_service_cb_not_allowed, (void *)allowed);
    }
    if (strcmp(allowed, "DELETE") != 0) {
        ulfius_add_endpoint_by_val(instance, "DELETE", prefix, resource, 0,
                                   &web_service_cb_not_allowed, (void *)allowed);
    }
}

static int respond_json(UResponse *response, int status, json_t *obj) {
    char *body;

    if (!obj) {
        ulfius_add_header_to_response(response, "Content-Type", "application/json");
        ulfius_set_string_body_response(response, status, "{}");
        return U_CALLBACK_CONTINUE;
    }

    body = json_dumps(obj, JSON_INDENT(2));
    json_decref(obj);

    if (!body) {
        ulfius_add_header_to_response(response, "Content-Type", "application/json");
        ulfius_set_string_body_response(response, status, "{}");
        return U_CALLBACK_CONTINUE;
    }

    ulfius_add_header_to_response(response, "Content-Type", "application/json");
    ulfius_set_string_body_response(response, status, body);
    free(body);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_get_ping(const URequest *request, UResponse *response,
                            void *user_data) {
    (void)request;
    (void)user_data;

    ulfius_set_string_body_response(response, HttpStatus_OK,
                                    HttpStatusStr(HttpStatus_OK));
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_get_version(const URequest *request, UResponse *response,
                               void *user_data) {
    (void)request;
    (void)user_data;

    ulfius_set_string_body_response(response, HttpStatus_OK, VERSION);
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_get_status(const URequest *request, UResponse *response,
                              void *user_data) {
    EpCtx *ctx = (EpCtx *)user_data;
    MetricsSnapshot snap;

    (void)request;

    if (!ctx || !ctx->store) {
        ulfius_set_string_body_response(response, HttpStatus_InternalServerError,
                                        HttpStatusStr(HttpStatus_InternalServerError));
        return U_CALLBACK_CONTINUE;
    }

    metrics_store_get(ctx->store, &snap);
    return respond_json(response, HttpStatus_OK, json_serialize_status(&snap));
}

int web_service_cb_get_metrics(const URequest *request, UResponse *response,
                               void *user_data) {
    EpCtx *ctx = (EpCtx *)user_data;
    MetricsSnapshot snap;
    const char *node_id;

    (void)request;

    if (!ctx || !ctx->store) {
        ulfius_set_string_body_response(response, HttpStatus_InternalServerError,
                                        HttpStatusStr(HttpStatus_InternalServerError));
        return U_CALLBACK_CONTINUE;
    }

    metrics_store_get(ctx->store, &snap);
    node_id = ctx->config ? ctx->config->nodeId : NULL;

    return respond_json(response, HttpStatus_OK,
                        json_serialize_metrics(&snap, node_id));
}

int web_service_cb_get_alarms(const URequest *request, UResponse *response,
                              void *user_data) {
    EpCtx *ctx = (EpCtx *)user_data;
    AlarmRecord history[64];
    MetricsSnapshot snap;
    json_t *obj, *arr;
    int count;

    (void)request;

    if (!ctx || !ctx->store) {
        ulfius_set_string_body_response(response, HttpStatus_InternalServerError,
                                        HttpStatusStr(HttpStatus_InternalServerError));
        return U_CALLBACK_CONTINUE;
    }

    metrics_store_get(ctx->store, &snap);
    count = metrics_store_get_alarm_history(ctx->store, history, 64);

    obj = json_object();
    json_object_set_new(obj, "active",  json_serialize_alarms(snap.alarms, ALARM_MAX));
    json_object_set_new(obj, "history", (arr = json_serialize_alarms(history, count), arr));

    return respond_json(response, HttpStatus_OK, obj);
}

int web_service_cb_put_absorption(const URequest *request, UResponse *response,
                                  void *user_data) {
    EpCtx *ctx = (EpCtx *)user_data;
    json_t *body;
    json_error_t err;
    double voltage_v;

    if (!ctx || !ctx->driver || !ctx->driver_ctx) {
        ulfius_set_string_body_response(response, HttpStatus_InternalServerError,
                                        HttpStatusStr(HttpStatus_InternalServerError));
        return U_CALLBACK_CONTINUE;
    }

    body = ulfius_get_json_body_request(request, &err);
    if (!body) {
        ulfius_set_string_body_response(response, HttpStatus_BadRequest, "Invalid JSON");
        return U_CALLBACK_CONTINUE;
    }

    if (json_deserialize_voltage_request(body, &voltage_v) != 0) {
        json_decref(body);
        ulfius_set_string_body_response(response, HttpStatus_BadRequest,
                                        "Missing 'voltage_v'");
        return U_CALLBACK_CONTINUE;
    }
    json_decref(body);

    if (!ctx->driver->set_absorption_voltage) {
        ulfius_set_string_body_response(response, HttpStatus_NotImplemented,
                                        HttpStatusStr(HttpStatus_NotImplemented));
        return U_CALLBACK_CONTINUE;
    }

    if (ctx->driver->set_absorption_voltage(ctx->driver_ctx, voltage_v) != 0) {
        ulfius_set_string_body_response(response, HttpStatus_ServiceUnavailable,
                                        "Driver returned error");
        return U_CALLBACK_CONTINUE;
    }

    ulfius_set_string_body_response(response, HttpStatus_OK,
                                    HttpStatusStr(HttpStatus_OK));
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_put_float(const URequest *request, UResponse *response,
                             void *user_data) {
    EpCtx *ctx = (EpCtx *)user_data;
    json_t *body;
    json_error_t err;
    double voltage_v;

    if (!ctx || !ctx->driver || !ctx->driver_ctx) {
        ulfius_set_string_body_response(response, HttpStatus_InternalServerError,
                                        HttpStatusStr(HttpStatus_InternalServerError));
        return U_CALLBACK_CONTINUE;
    }

    body = ulfius_get_json_body_request(request, &err);
    if (!body) {
        ulfius_set_string_body_response(response, HttpStatus_BadRequest, "Invalid JSON");
        return U_CALLBACK_CONTINUE;
    }

    if (json_deserialize_voltage_request(body, &voltage_v) != 0) {
        json_decref(body);
        ulfius_set_string_body_response(response, HttpStatus_BadRequest,
                                        "Missing 'voltage_v'");
        return U_CALLBACK_CONTINUE;
    }
    json_decref(body);

    if (!ctx->driver->set_float_voltage) {
        ulfius_set_string_body_response(response, HttpStatus_NotImplemented,
                                        HttpStatusStr(HttpStatus_NotImplemented));
        return U_CALLBACK_CONTINUE;
    }

    if (ctx->driver->set_float_voltage(ctx->driver_ctx, voltage_v) != 0) {
        ulfius_set_string_body_response(response, HttpStatus_ServiceUnavailable,
                                        "Driver returned error");
        return U_CALLBACK_CONTINUE;
    }

    ulfius_set_string_body_response(response, HttpStatus_OK,
                                    HttpStatusStr(HttpStatus_OK));
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_post_relay(const URequest *request, UResponse *response,
                              void *user_data) {
    EpCtx *ctx = (EpCtx *)user_data;
    json_t *body;
    json_error_t err;
    bool state;

    if (!ctx || !ctx->driver || !ctx->driver_ctx) {
        ulfius_set_string_body_response(response, HttpStatus_InternalServerError,
                                        HttpStatusStr(HttpStatus_InternalServerError));
        return U_CALLBACK_CONTINUE;
    }

    body = ulfius_get_json_body_request(request, &err);
    if (!body) {
        ulfius_set_string_body_response(response, HttpStatus_BadRequest, "Invalid JSON");
        return U_CALLBACK_CONTINUE;
    }

    if (json_deserialize_relay_request(body, &state) != 0) {
        json_decref(body);
        ulfius_set_string_body_response(response, HttpStatus_BadRequest,
                                        "Missing 'state'");
        return U_CALLBACK_CONTINUE;
    }
    json_decref(body);

    if (!ctx->driver->set_relay) {
        ulfius_set_string_body_response(response, HttpStatus_NotImplemented,
                                        HttpStatusStr(HttpStatus_NotImplemented));
        return U_CALLBACK_CONTINUE;
    }

    if (ctx->driver->set_relay(ctx->driver_ctx, state) != 0) {
        ulfius_set_string_body_response(response, HttpStatus_ServiceUnavailable,
                                        "Driver returned error");
        return U_CALLBACK_CONTINUE;
    }

    ulfius_set_string_body_response(response, HttpStatus_OK,
                                    HttpStatusStr(HttpStatus_OK));
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_default(const URequest *request, UResponse *response,
                           void *user_data) {
    (void)request;
    (void)user_data;

    ulfius_set_string_body_response(response, HttpStatus_NotFound,
                                    HttpStatusStr(HttpStatus_NotFound));
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_not_allowed(const URequest *request, UResponse *response,
                               void *user_data) {
    (void)request;
    (void)user_data;

    ulfius_set_string_body_response(response, HttpStatus_MethodNotAllowed,
                                    HttpStatusStr(HttpStatus_MethodNotAllowed));
    return U_CALLBACK_CONTINUE;
}

static void register_endpoints(UInst *inst, EpCtx *ctx) {
    ulfius_add_endpoint_by_val(inst, "GET", URL_PREFIX, API_RES_EP("ping"), 0,
                               &web_service_cb_get_ping, ctx);
    setup_unsupported_methods(inst, "GET", URL_PREFIX, API_RES_EP("ping"));

    ulfius_add_endpoint_by_val(inst, "GET", URL_PREFIX, API_RES_EP("version"), 0,
                               &web_service_cb_get_version, ctx);
    setup_unsupported_methods(inst, "GET", URL_PREFIX, API_RES_EP("version"));

    ulfius_add_endpoint_by_val(inst, "GET", URL_PREFIX, API_RES_EP("status"), 0,
                               &web_service_cb_get_status, ctx);
    setup_unsupported_methods(inst, "GET", URL_PREFIX, API_RES_EP("status"));

    ulfius_add_endpoint_by_val(inst, "GET", URL_PREFIX, API_RES_EP("metrics"), 0,
                               &web_service_cb_get_metrics, ctx);
    setup_unsupported_methods(inst, "GET", URL_PREFIX, API_RES_EP("metrics"));

    ulfius_add_endpoint_by_val(inst, "GET", URL_PREFIX, API_RES_EP("alarms"), 0,
                               &web_service_cb_get_alarms, ctx);
    setup_unsupported_methods(inst, "GET", URL_PREFIX, API_RES_EP("alarms"));

    ulfius_add_endpoint_by_val(inst, "PUT", URL_PREFIX, API_RES_EP("absorption"), 0,
                               &web_service_cb_put_absorption, ctx);
    setup_unsupported_methods(inst, "PUT", URL_PREFIX, API_RES_EP("absorption"));

    ulfius_add_endpoint_by_val(inst, "PUT", URL_PREFIX, API_RES_EP("float"), 0,
                               &web_service_cb_put_float, ctx);
    setup_unsupported_methods(inst, "PUT", URL_PREFIX, API_RES_EP("float"));

    ulfius_add_endpoint_by_val(inst, "POST", URL_PREFIX, API_RES_EP("relay"), 0,
                               &web_service_cb_post_relay, ctx);
    setup_unsupported_methods(inst, "POST", URL_PREFIX, API_RES_EP("relay"));

    ulfius_set_default_endpoint(inst, &web_service_cb_default, ctx);
}

int web_service_start(const Config *config, UInst *inst, EpCtx *ctx) {
    if (!config || !inst || !ctx) return -1;

    if (ulfius_init_instance(inst, config->listenPort, NULL, NULL) != U_OK) {
        usys_log_error("web_service: failed to init ulfius instance");
        return -1;
    }

    u_map_put(inst->default_headers, "Access-Control-Allow-Origin", "*");

    register_endpoints(inst, ctx);

    if (ulfius_start_framework(inst) != U_OK) {
        usys_log_error("web_service: failed to start framework");
        ulfius_stop_framework(inst);
        ulfius_clean_instance(inst);
        return -1;
    }

    usys_log_info("web_service: listening on %s:%d",
                  config->listenAddr, config->listenPort);
    return 0;
}

void web_service_stop(UInst *inst) {
    if (!inst) return;

    ulfius_stop_framework(inst);
    ulfius_clean_instance(inst);
}
