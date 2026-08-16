/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <jansson.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <ulfius.h>

#include "http_status.h"
#include "web_service.h"

#include "usys_log.h"

#include "version.h"

static int reply_json(struct _u_response *response,
                      HttpStatus status,
                      json_t *json) {

    if (!json) {
        ulfius_set_string_body_response(response,
                                        HttpStatus_InternalServerError,
                                        "internal error");
        return U_CALLBACK_CONTINUE;
    }

    ulfius_set_json_body_response(response, status, json);
    json_decref(json);
    return U_CALLBACK_CONTINUE;
}

static json_t *load_body(const struct _u_request *request) {

    json_error_t error;

    if (!request || !request->binary_body ||
        request->binary_body_length == 0) {
        return NULL;
    }

    return json_loadb((const char *)request->binary_body,
                      request->binary_body_length,
                      0,
                      &error);
}

static const char *json_string_or_null(json_t *json, const char *key) {

    json_t *value;

    if (!json_is_object(json) || !key) return NULL;
    value = json_object_get(json, key);

    return json_is_string(value) ? json_string_value(value) : NULL;
}

static int ping_cb(const struct _u_request *request,
                   struct _u_response *response,
                   void *userData) {

    (void)request;
    (void)userData;

    ulfius_set_string_body_response(response, HttpStatus_OK, "OK");
    return U_CALLBACK_CONTINUE;
}

static int version_cb(const struct _u_request *request,
                      struct _u_response *response,
                      void *userData) {

    (void)request;
    (void)userData;

    ulfius_set_string_body_response(response, HttpStatus_OK, VERSION);
    return U_CALLBACK_CONTINUE;
}

static int ready_cb(const struct _u_request *request,
                    struct _u_response *response,
                    void *userData) {

    json_t *json;

    (void)request;
    (void)userData;

    json = json_object();
    json_object_set_new(json, "ready", json_true());
    json_object_set_new(json,
                        "reason",
                        json_string("lifecycle manager running"));

    return reply_json(response, HttpStatus_OK, json);
}

static json_t *fsm_status_json(const LifecycleFsm *fsm,
                               const StarterSnapshot *starter,
                               const char *bootId,
                               size_t pendingEvents) {

    json_t *root;
    json_t *checkIn;
    json_t *configuration;
    json_t *starterJson;
    int64_t nowMs;
    int64_t remainingMs;

    root = json_object();
    checkIn = json_object();
    configuration = json_object();
    starterJson = json_object();

    if (!root || !checkIn || !configuration || !starterJson) {
        json_decref(root);
        json_decref(checkIn);
        json_decref(configuration);
        json_decref(starterJson);
        return NULL;
    }

    nowMs = lifecycle_boottime_ms();

    json_object_set_new(root,
                        "state",
                        json_string(lifecycle_state_str(fsm->state)));
    json_object_set_new(root, "reason", json_string(fsm->reason));
    json_object_set_new(root,
                        "stateSince",
                        json_integer(fsm->stateSince));
    json_object_set_new(root, "bootId", json_string(bootId));
    json_object_set_new(root,
                        "sequence",
                        json_integer((json_int_t)fsm->sequence));

    remainingMs = fsm->checkInDeadlineMs - nowMs;
    if (remainingMs < 0) remainingMs = 0;
    json_object_set_new(checkIn,
                        "gateOpen",
                        json_boolean(fsm->gateOpen));
    json_object_set_new(checkIn,
                        "remainingSec",
                        json_integer((remainingMs + 999) / 1000));
    json_object_set_new(root, "checkIn", checkIn);

    remainingMs = fsm->configDeadlineMs - nowMs;
    if (remainingMs < 0) remainingMs = 0;
    json_object_set_new(configuration,
                        "requestId",
                        fsm->requestId[0] ?
                        json_string(fsm->requestId) : json_null());
    json_object_set_new(configuration,
                        "assignmentId",
                        fsm->assignmentId[0] ?
                        json_string(fsm->assignmentId) : json_null());
    json_object_set_new(configuration,
                        "received",
                        json_boolean(fsm->configurationSeen));
    json_object_set_new(configuration,
                        "applied",
                        json_boolean(fsm->configurationApplied));
    json_object_set_new(configuration,
                        "remainingSec",
                        json_integer((remainingMs + 999) / 1000));
    json_object_set_new(root, "configuration", configuration);

    json_object_set_new(starterJson,
                        "available",
                        json_boolean(starter->available));
    json_object_set_new(starterJson,
                        "readiness",
                        json_string(starter_aggregate_str(
                            starter->aggregate)));
    json_object_set_new(starterJson,
                        "reason",
                        json_string(starter->aggregateReason));
    json_object_set_new(starterJson,
                        "configPhase",
                        json_string(config_phase_str(
                            starter->configPhase)));
    json_object_set_new(starterJson,
                        "configReason",
                        json_string(starter->configReason));
    json_object_set_new(root, "starter", starterJson);

    json_object_set_new(root,
                        "notificationPending",
                        json_boolean(pendingEvents > 0));
    return root;
}

static int status_cb(const struct _u_request *request,
                     struct _u_response *response,
                     void *userData) {

    LifecycleContext *ctx;
    LifecycleFsm fsm;
    StarterSnapshot starter;
    json_t *json;
    size_t pendingEvents;

    (void)request;

    ctx = (LifecycleContext *)userData;
    if (!ctx) {
        ulfius_set_string_body_response(response,
                                        HttpStatus_InternalServerError,
                                        "context unavailable");
        return U_CALLBACK_CONTINUE;
    }

    pthread_mutex_lock(&ctx->mutex);
    fsm = ctx->fsm;
    starter = ctx->starter;
    pendingEvents = ctx->eventCount;
    pthread_mutex_unlock(&ctx->mutex);

    json = fsm_status_json(&fsm,
                           &starter,
                           ctx->bootId,
                           pendingEvents);
    return reply_json(response, HttpStatus_OK, json);
}

static int check_in_cb(const struct _u_request *request,
                       struct _u_response *response,
                       void *userData) {

    LifecycleContext *ctx;
    LifecycleFsm fsm;
    json_t *body;
    json_t *json;
    const char *bootId;
    const char *bootResult;
    char error[LIFECYCLED_REASON_LEN];
    bool healthy;

    ctx = (LifecycleContext *)userData;
    body = load_body(request);

    if (!ctx || !body) {
        json_decref(body);
        return reply_json(response,
                          HttpStatus_BadRequest,
                          json_pack("{s:s}",
                                    "error",
                                    "valid JSON body required"));
    }

    bootId = json_string_or_null(body, "bootId");
    bootResult = json_string_or_null(body, "bootResult");

    if (!bootResult ||
        (strcmp(bootResult, "ready") != 0 &&
         strcmp(bootResult, "degraded") != 0)) {
        json_decref(body);
        return reply_json(response,
                          HttpStatus_BadRequest,
                          json_pack("{s:s}",
                                    "error",
                                    "bootResult must be ready or degraded"));
    }

    healthy = strcmp(bootResult, "ready") == 0;
    memset(error, 0, sizeof(error));

    if (!lifecycle_context_check_in(ctx,
                                    bootId,
                                    healthy,
                                    error,
                                    sizeof(error))) {
        json_decref(body);
        return reply_json(response,
                          HttpStatus_Conflict,
                          json_pack("{s:s}", "error", error));
    }

    json_decref(body);
    lifecycle_context_snapshot(ctx, &fsm, NULL);
    json = json_pack("{s:s,s:s}",
                     "state",
                     lifecycle_state_str(fsm.state),
                     "status",
                     "accepted");
    return reply_json(response, HttpStatus_Accepted, json);
}

static int gate_cb(const struct _u_request *request,
                   struct _u_response *response,
                   void *userData) {

    LifecycleContext *ctx;
    LifecycleFsm fsm;
    json_t *json;
    int remainingSec;
    HttpStatus status;

    (void)request;

    ctx = (LifecycleContext *)userData;
    lifecycle_context_snapshot(ctx, &fsm, NULL);
    status = lifecycle_fsm_gate_status(&fsm,
                                       lifecycle_boottime_ms(),
                                       &remainingSec);

    json = json_object();
    json_object_set_new(json,
                        "proceed",
                        json_boolean(status == HttpStatus_OK));
    json_object_set_new(json,
                        "remainingSec",
                        json_integer(remainingSec));
    json_object_set_new(json,
                        "state",
                        json_string(lifecycle_state_str(fsm.state)));
    if (status == HttpStatus_ServiceUnavailable) {
        json_object_set_new(json, "reason", json_string(fsm.reason));
    }

    return reply_json(response, status, json);
}

static int configure_cb(const struct _u_request *request,
                        struct _u_response *response,
                        void *userData) {

    LifecycleContext *ctx;
    LifecycleConfigureResult result;
    json_t *body;
    json_t *json;
    const char *requestId;
    const char *assignmentId;
    char error[LIFECYCLED_REASON_LEN];
    HttpStatus status;

    ctx = (LifecycleContext *)userData;
    body = load_body(request);

    if (!ctx || !body) {
        json_decref(body);
        return reply_json(response,
                          HttpStatus_BadRequest,
                          json_pack("{s:s}",
                                    "error",
                                    "valid JSON body required"));
    }

    requestId = json_string_or_null(body, "requestId");
    assignmentId = json_string_or_null(body, "assignmentId");
    memset(error, 0, sizeof(error));

    result = lifecycle_context_configure(ctx,
                                         requestId,
                                         assignmentId,
                                         error,
                                         sizeof(error));
    json_decref(body);

    if (result == LIFECYCLE_CONFIGURE_ACCEPTED) {
        status = HttpStatus_Accepted;
        json = json_pack("{s:s,s:s}",
                         "state",
                         "CONFIGURING",
                         "status",
                         "accepted");
    } else if (result == LIFECYCLE_CONFIGURE_DUPLICATE) {
        status = HttpStatus_OK;
        json = json_pack("{s:s}", "status", "already accepted");
    } else if (result == LIFECYCLE_CONFIGURE_INVALID_REQUEST) {
        status = HttpStatus_BadRequest;
        json = json_pack("{s:s}", "error", error);
    } else {
        status = HttpStatus_Conflict;
        json = json_pack("{s:s}", "error", error);
    }

    return reply_json(response, status, json);
}

bool web_service_start(LifecycleContext *ctx) {

    if (!ctx || !ctx->config) return false;

    ctx->uInstance = calloc(1, sizeof(struct _u_instance));
    if (!ctx->uInstance) return false;

    if (ulfius_init_instance(ctx->uInstance,
                             ctx->config->httpPort,
                             NULL,
                             NULL) != U_OK) {
        usys_log_error("web: init failed");
        free(ctx->uInstance);
        ctx->uInstance = NULL;
        return false;
    }

    ulfius_add_endpoint_by_val(ctx->uInstance,
                               "GET", "/v1", "/ping", 0,
                               &ping_cb, ctx);
    ulfius_add_endpoint_by_val(ctx->uInstance,
                               "GET", "/v1", "/version", 0,
                               &version_cb, ctx);
    ulfius_add_endpoint_by_val(ctx->uInstance,
                               "GET", "/v1", "/ready", 0,
                               &ready_cb, ctx);
    ulfius_add_endpoint_by_val(ctx->uInstance,
                               "GET", "/v1", "/status", 0,
                               &status_cb, ctx);
    ulfius_add_endpoint_by_val(ctx->uInstance,
                               "POST", "/v1", "/check-in", 0,
                               &check_in_cb, ctx);
    ulfius_add_endpoint_by_val(ctx->uInstance,
                               "GET", "/v1", "/gate", 0,
                               &gate_cb, ctx);
    ulfius_add_endpoint_by_val(ctx->uInstance,
                               "POST", "/v1", "/configure", 0,
                               &configure_cb, ctx);

    if (ulfius_start_framework(ctx->uInstance) != U_OK) {
        usys_log_error("web: start failed");
        ulfius_clean_instance(ctx->uInstance);
        free(ctx->uInstance);
        ctx->uInstance = NULL;
        return false;
    }

    usys_log_info("web: listening on %s:%d",
                  ctx->config->httpAddr,
                  ctx->config->httpPort);
    return true;
}

void web_service_stop(LifecycleContext *ctx) {

    if (!ctx || !ctx->uInstance) return;

    ulfius_stop_framework(ctx->uInstance);
    ulfius_clean_instance(ctx->uInstance);
    free(ctx->uInstance);
    ctx->uInstance = NULL;
}
