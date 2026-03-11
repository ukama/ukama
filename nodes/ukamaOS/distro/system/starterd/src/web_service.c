/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include <ulfius.h>
#include <jansson.h>

#include "usys_log.h"
#include "starterd.h"
#include "version.h"
#include "web_service.h"
#include "web_client.h"
#include "network.h"
#include "jserdes.h"
#include "supervisor.h"
#include "http_status.h"

static int ws_reply_text(struct _u_response *resp, int status, const char *text) {

    ulfius_set_string_body_response(resp, status, text ? text : "");
    return U_CALLBACK_CONTINUE;
}

static int ws_ping_cb(const struct _u_request *req,
                      struct _u_response *resp,
                      void *userData) {

    (void)req;
    (void)userData;

    return ws_reply_text(resp,
                         HttpStatus_OK,
                         HttpStatusStr(HttpStatus_OK));
}

static int ws_version_cb(const struct _u_request *req,
                         struct _u_response *resp,
                         void *userData) {

    (void)req;
    (void)userData;

    return ws_reply_text(resp,
                         HttpStatus_OK,
                         VERSION);
}

static int ws_status_cb(const struct _u_request *req,
                        struct _u_response *resp,
                        void *userData) {

    StarterContext *ctx;
    json_t *j;
    json_t *meta;
    char *body;

    (void)req;

    ctx = (StarterContext *)userData;
    if (!ctx || !ctx->spaceList) {
        return ws_reply_text(resp,
                             HttpStatus_InternalServerError,
                             HttpStatusStr(HttpStatus_InternalServerError));
    }

    j = jserdes_status_json(ctx->spaceList);
    if (!j) {
        return ws_reply_text(resp,
                             HttpStatus_InternalServerError,
                             HttpStatusStr(HttpStatus_InternalServerError));
    }

    meta = json_object();
    if (meta) {
        json_object_set_new(meta, "updateInProgress",
                            json_boolean(ctx->updateInProgress ? 1 : 0));
        json_object_set_new(meta, "switchRequested",
                            json_boolean(ctx->switchRequested ? 1 : 0));
        json_object_set_new(meta, "exitCode",
                            json_integer(ctx->exitCode));
        json_object_set_new(j, "starterd", meta);
    }

    body = json_dumps(j, JSON_INDENT(2) | JSON_SORT_KEYS);
    json_decref(j);

    ulfius_add_header_to_response(resp, "Content-Type", "application/json");
    ulfius_set_string_body_response(resp,
                                    HttpStatus_OK,
                                    body ? body : "{}");
    free(body);

    return U_CALLBACK_CONTINUE;
}

static bool ws_parse_update(json_t *j,
                            char **spaceOut,
                            char **nameOut,
                            char **tagOut,
                            char **hubOut) {

    json_t *v;
    const char *space;
    const char *name;
    const char *tag;
    const char *hub;

    if (spaceOut) *spaceOut = NULL;
    if (nameOut)  *nameOut  = NULL;
    if (tagOut)   *tagOut   = NULL;
    if (hubOut)   *hubOut   = NULL;

    if (!j || !json_is_object(j)) {
        return false;
    }

    v = json_object_get(j, "space");
    space = json_is_string(v) ? json_string_value(v) : NULL;

    v = json_object_get(j, "name");
    name = json_is_string(v) ? json_string_value(v) : NULL;

    v = json_object_get(j, "tag");
    tag = json_is_string(v) ? json_string_value(v) : NULL;

    v = json_object_get(j, "hub");
    hub = json_is_string(v) ? json_string_value(v) : NULL;

    if (!space || !name || !tag || !hub || !*hub) {
        return false;
    }

    if (spaceOut) *spaceOut = strdup(space);
    if (nameOut)  *nameOut  = strdup(name);
    if (tagOut)   *tagOut   = strdup(tag);
    if (hubOut)   *hubOut   = strdup(hub);

    return true;
}

static int ws_update_cb(const struct _u_request *req,
                        struct _u_response *resp,
                        void *userData) {

    StarterContext *ctx;
    json_error_t err;
    json_t *j;
    char *space;
    char *name;
    char *tag;
    Action *a;

    ctx = (StarterContext *)userData;
    if (!ctx || !ctx->queue || !ctx->supervisor) {
        return ws_reply_text(resp,
                             HttpStatus_InternalServerError,
                             HttpStatusStr(HttpStatus_InternalServerError));
    }

    if (ctx->switchRequested || ctx->terminateRequested) {
        return ws_reply_text(resp,
                             HttpStatus_Conflict,
                             HttpStatusStr(HttpStatus_Conflict));
    }

    if (ctx->updateInProgress) {
        return ws_reply_text(resp,
                             HttpStatus_Locked,
                             HttpStatusStr(HttpStatus_Locked));
    }

    j = json_loads(req->binary_body ? (const char *)req->binary_body : "{}",
                   0,
                   &err);
    if (!j) {
        return ws_reply_text(resp,
                             HttpStatus_BadRequest,
                             HttpStatusStr(HttpStatus_BadRequest));
    }

    space = NULL;
    name  = NULL;
    tag   = NULL;

    if (!ws_parse_update(j, &space, &name, &tag, &hub)) {
        json_decref(j);
        free(space);
        free(name);
        free(tag);
        free(hub);
        return ws_reply_text(resp,
                             HttpStatus_BadRequest,
                             HttpStatusStr(HttpStatus_BadRequest));
    }

    ctx->updateInProgress = 1;

    a = action_new(ACTION_UPDATE_APP, space, name, tag);
    free(space);
    free(name);
    free(tag);
    free(hub);
    json_decref(j);

    if (!a || !actions_enqueue(ctx->queue, a)) {
        if (a) {
            free(a);
        }
        ctx->updateInProgress = 0;
        return ws_reply_text(resp,
                             HttpStatus_InternalServerError,
                             HttpStatusStr(HttpStatus_InternalServerError));
    }

    supervisor_signal((Supervisor *)ctx->supervisor);

    return ws_reply_text(resp,
                         HttpStatus_Accepted,
                         HttpStatusStr(HttpStatus_Accepted));
}

static int ws_terminate_cb(const struct _u_request *req,
                           struct _u_response *resp,
                           void *userData) {

    StarterContext *ctx;
    json_error_t err;
    json_t *j;
    json_t *v;
    const char *space;
    const char *name;
    Action *a;

    ctx = (StarterContext *)userData;
    if (!ctx || !ctx->queue || !ctx->supervisor) {
        return ws_reply_text(resp,
                             HttpStatus_InternalServerError,
                             HttpStatusStr(HttpStatus_InternalServerError));
    }

    if (ctx->switchRequested) {
        return ws_reply_text(resp,
                             HttpStatus_Conflict,
                             HttpStatusStr(HttpStatus_Conflict));
    }

    j = json_loads(req->binary_body ? (const char *)req->binary_body : "{}",
                   0,
                   &err);
    if (!j) {
        return ws_reply_text(resp,
                             HttpStatus_BadRequest,
                             HttpStatusStr(HttpStatus_BadRequest));
    }

    v = json_object_get(j, "space");
    space = json_is_string(v) ? json_string_value(v) : NULL;

    v = json_object_get(j, "name");
    name = json_is_string(v) ? json_string_value(v) : NULL;

    if (!space || !name) {
        json_decref(j);
        return ws_reply_text(resp,
                             HttpStatus_BadRequest,
                             HttpStatusStr(HttpStatus_BadRequest));
    }

    a = action_new(ACTION_TERMINATE_APP, space, name, NULL);
    json_decref(j);

    if (!a || !actions_enqueue(ctx->queue, a)) {
        if (a) {
            free(a);
        }
        return ws_reply_text(resp,
                             HttpStatus_InternalServerError,
                             HttpStatusStr(HttpStatus_InternalServerError));
    }

    supervisor_signal((Supervisor *)ctx->supervisor);

    return ws_reply_text(resp,
                         HttpStatus_Accepted,
                         HttpStatusStr(HttpStatus_Accepted));
}

static int ws_cb_not_allowed(const struct _u_request *request,
                             struct _u_response *response,
                             void *user_data) {

    const char *allowedMethod = (const char *)user_data;

    (void)request;

    u_map_put(response->map_header, "Allow", allowedMethod);
    ulfius_set_string_body_response(response,
                                    HttpStatus_MethodNotAllowed,
                                    HttpStatusStr(HttpStatus_MethodNotAllowed));
    return U_CALLBACK_CONTINUE;
}

static void setup_unsupported_methods(UInst *instance,
                                      const char *allowedMethod,
                                      const char *prefix,
                                      const char *resource) {

    if (strcmp(allowedMethod, "GET") != 0) {
        ulfius_add_endpoint_by_val(instance, "GET",
                                   prefix, resource, 0,
                                   &ws_cb_not_allowed,
                                   (void *)allowedMethod);
    }

    if (strcmp(allowedMethod, "POST") != 0) {
        ulfius_add_endpoint_by_val(instance, "POST",
                                   prefix, resource, 0,
                                   &ws_cb_not_allowed,
                                   (void *)allowedMethod);
    }

    if (strcmp(allowedMethod, "PUT") != 0) {
        ulfius_add_endpoint_by_val(instance, "PUT",
                                   prefix, resource, 0,
                                   &ws_cb_not_allowed,
                                   (void *)allowedMethod);
    }

    if (strcmp(allowedMethod, "DELETE") != 0) {
        ulfius_add_endpoint_by_val(instance, "DELETE",
                                   prefix, resource, 0,
                                   &ws_cb_not_allowed,
                                   (void *)allowedMethod);
    }
}

bool web_service_start(StarterContext *ctx) {

    if (!ctx || !ctx->uInstance) {
        return false;
    }

    ulfius_add_endpoint_by_val(ctx->uInstance, "GET",
                               "/v1", "/ping", 0,
                               &ws_ping_cb, ctx);
    setup_unsupported_methods(ctx->uInstance, "GET",
                              "/v1", "/ping");

    ulfius_add_endpoint_by_val(ctx->uInstance, "GET",
                               "/v1", "/version", 0,
                               &ws_version_cb, ctx);
    setup_unsupported_methods(ctx->uInstance, "GET",
                              "/v1", "/version");

    ulfius_add_endpoint_by_val(ctx->uInstance, "GET",
                               "/v1", "/status", 0,
                               &ws_status_cb, ctx);
    setup_unsupported_methods(ctx->uInstance, "GET",
                              "/v1", "/status");

    ulfius_add_endpoint_by_val(ctx->uInstance, "POST",
                               "/v1", "/update", 0,
                               &ws_update_cb, ctx);
    setup_unsupported_methods(ctx->uInstance, "POST",
                              "/v1", "/update");

    ulfius_add_endpoint_by_val(ctx->uInstance, "POST",
                               "/v1", "/terminate", 0,
                               &ws_terminate_cb, ctx);
    setup_unsupported_methods(ctx->uInstance, "POST",
                              "/v1", "/terminate");

    if (ulfius_start_framework(ctx->uInstance) != U_OK) {
        usys_log_error("web: start failed");
        return false;
    }

    return true;
}

void web_service_stop(StarterContext *ctx) {

    if (!ctx || !ctx->uInstance) {
        return;
    }

    ulfius_stop_framework(ctx->uInstance);
}
