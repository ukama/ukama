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
#include "network.h"
#include "jserdes.h"
#include "supervisor.h"
#include "http_status.h"

#include "version.h"

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
    char *body;

    (void)req;

    ctx = (StarterContext *)userData;
    if (!ctx || !ctx->spaceList) {
        return ws_reply_text(resp,
                             HttpStatus_InternalServerError,
                             HttpStatusStr(HttpStatus_InternalServerError));
    }

    j = jserdes_status_json(ctx->spaceList);
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
                            char **tagOut) {

    json_t *v;
    const char *space;
    const char *name;
    const char *tag;

    if (spaceOut) *spaceOut = NULL;
    if (nameOut)  *nameOut = NULL;
    if (tagOut)   *tagOut = NULL;

    if (!j || !json_is_object(j)) {
        return false;
    }

    v = json_object_get(j, "space");
    space = json_is_string(v) ? json_string_value(v) : NULL;

    v = json_object_get(j, "name");
    name = json_is_string(v) ? json_string_value(v) : NULL;

    v = json_object_get(j, "tag");
    tag = json_is_string(v) ? json_string_value(v) : NULL;

    if (!space || !name || !tag) return false;

    if (spaceOut) *spaceOut = strdup(space);
    if (nameOut)  *nameOut  = strdup(name);
    if (tagOut)   *tagOut   = strdup(tag);

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

    j = json_loads(req->binary_body ? (const char*)req->binary_body : "{}", 0, &err);
    if (!j) {
        return ws_reply_text(resp,
                             HttpStatus_BadRequest,
                             HttpStatusStr(HttpStatus_BadRequest));
    }

    space = NULL;
    name  = NULL;
    tag   = NULL;

    if (!ws_parse_update(j, &space, &name, &tag)) {
        json_decref(j);
        free(space);
        free(name);
        free(tag);
        return ws_reply_text(resp,
                             HttpStatus_BadRequest,
                             HttpStatusStr(HttpStatus_BadRequest));
    }

    a = action_new(ACTION_UPDATE_APP, space, name, tag);
    free(space);
    free(name);
    free(tag);
    json_decref(j);

    if (!a || !actions_enqueue(ctx->queue, a)) {
        if (a) free(a);
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

    j = json_loads(req->binary_body ? (const char*)req->binary_body : "{}", 0, &err);
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
        if (a) free(a);
        return ws_reply_text(resp,
                             HttpStatus_InternalServerError,
                             HttpStatusStr(HttpStatus_InternalServerError));
    }

    supervisor_signal((Supervisor *)ctx->supervisor);
    return ws_reply_text(resp,
                         HttpStatus_Accepted,
                         HttpStatusStr(HttpStatus_Accepted));
}

bool web_service_start(StarterContext *ctx) {

    if (!ctx || !ctx->uInstance) return false;

    ulfius_add_endpoint_by_val(ctx->uInstance, "GET",  "/v1", "/ping",     0, &ws_ping_cb,      ctx);
    ulfius_add_endpoint_by_val(ctx->uInstance, "GET",  "/v1", "/version",  0, &ws_version_cb,   ctx);
    ulfius_add_endpoint_by_val(ctx->uInstance, "GET",  "/v1", "/status",   0, &ws_status_cb,    ctx);
    ulfius_add_endpoint_by_val(ctx->uInstance, "POST", "/v1", "/update",   0, &ws_update_cb,    ctx);
    ulfius_add_endpoint_by_val(ctx->uInstance, "POST", "/v1", "/terminate",0, &ws_terminate_cb, ctx);

    if (ulfius_start_framework(ctx->uInstance) != U_OK) {
        usys_log_error("web: start failed");
        return false;
    }

    return true;
}

void web_service_stop(StarterContext *ctx) {

    (void)ctx;
}
