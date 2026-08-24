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
#include "app.h"
#include "readiness.h"

static int ws_reply_text(struct _u_response *resp,
                         int status,
                         const char *text) {

    ulfius_set_string_body_response(resp, status, text ? text : "");
    return U_CALLBACK_CONTINUE;
}

static json_t *ws_load_json_body(const struct _u_request *req) {

    json_error_t err;

    if (!req || !req->binary_body || req->binary_body_length <= 0) {
        usys_log_error("web: empty json body");
        return NULL;
    }

    return json_loadb((const char *)req->binary_body,
                      req->binary_body_length,
                      0,
                      &err);
}

static bool ws_app_exists(StarterContext *ctx,
                          const char *space,
                          const char *name) {

    if (!ctx || !ctx->spaceList || !space || !name) {
        return false;
    }

    if (!app_find(ctx->spaceList, space, name)) {
        usys_log_error("web: app not found: space=%s name=%s", space, name);
        return false;
    }

    return true;
}

static bool ws_update_target_exists(StarterContext *ctx,
                                    const char *space,
                                    const char *name) {

    Space *sp = NULL;

    if (!ctx || !ctx->spaceList || !name) {
        return false;
    }

    if (space) {
        return ws_app_exists(ctx, space, name);
    }

    sp = ctx->spaceList;
    while (sp) {
        if (sp->name && app_find(ctx->spaceList, sp->name, name)) {
            return true;
        }
        sp = sp->next;
    }

    usys_log_error("web: app not found in any space: name=%s", name);
    return false;
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

static int ws_ready_cb(const struct _u_request *req,
                       struct _u_response *resp,
                       void *userData) {

    StarterContext *ctx;
    NodeReadinessState state;
    json_t *json;
    char reason[STARTERD_READY_REASON_LEN];
    int status;

    (void)req;

    ctx = (StarterContext *)userData;
    state = readiness_get(ctx ? ctx->readiness : NULL,
                          reason,
                          sizeof(reason));

    status = HttpStatus_Accepted;
    if (state == NODE_READINESS_READY) {
        status = HttpStatus_OK;
    } else if (state == NODE_READINESS_FAULTY) {
        status = HttpStatus_ServiceUnavailable;
    }

    json = json_object();
    if (!json) {
        return ws_reply_text(resp,
                             HttpStatus_InternalServerError,
                             HttpStatusStr(HttpStatus_InternalServerError));
    }

    json_object_set_new(json,
                        "ready",
                        json_boolean(state == NODE_READINESS_READY));
    if (state != NODE_READINESS_READY) {
        json_object_set_new(json, "reason", json_string(reason));
    }

    ulfius_set_json_body_response(resp, status, json);
    json_decref(json);
    return U_CALLBACK_CONTINUE;
}

static int ws_status_cb(const struct _u_request *req,
                        struct _u_response *resp,
                        void *userData) {

    StarterContext *ctx;
    json_t *j;
    json_t *meta;
    json_t *readiness;
    json_t *connectivity;
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
        json_object_set_new(meta, "terminateRequested",
                            json_boolean(ctx->terminateRequested ? 1 : 0));
        json_object_set_new(meta, "restartRequested",
                            json_boolean(ctx->restartRequested ? 1 : 0));
        json_object_set_new(meta, "switchRequested",
                            json_boolean(ctx->switchRequested ? 1 : 0));
        json_object_set_new(meta, "exitCode",
                            json_integer(ctx->exitCode));
        readiness = readiness_status_json(ctx->readiness);
        connectivity = readiness_connectivity_json(ctx->readiness);
        if (readiness) {
            json_object_set_new(meta, "readiness", readiness);
        }
        if (connectivity) {
            json_object_set_new(meta, "connectivity", connectivity);
        }
        json_object_set_new(j, "starterd", meta);
    }

    body = json_dumps(j, JSON_INDENT(2));
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
    json_t *item;
    const char *space;
    const char *name;
    const char *tag;
    const char *hub;
    char *hubPayload;
    size_t i;
    size_t count;

    if (spaceOut) *spaceOut = NULL;
    if (nameOut)  *nameOut  = NULL;
    if (tagOut)   *tagOut   = NULL;
    if (hubOut)   *hubOut   = NULL;

    if (!j || !json_is_object(j)) {
        return false;
    }

    space = NULL;

    v = json_object_get(j, "space");
    if (json_is_string(v)) {
        space = json_string_value(v);
        if (!space || !*space) {
            return false;
        }
    } else if (v != NULL) {
        return false;
    }

    v = json_object_get(j, "name");
    name = json_is_string(v) ? json_string_value(v) : NULL;

    v = json_object_get(j, "tag");
    tag = json_is_string(v) ? json_string_value(v) : NULL;

    if (!name || !*name || !tag || !*tag) {
        return false;
    }

    hubPayload = NULL;

    v = json_object_get(j, "hub");
    if (json_is_string(v)) {
        hub = json_string_value(v);
        if (hub && *hub) {
            hubPayload = strdup(hub);
            if (!hubPayload) {
                return false;
            }
        }
    } else if (json_is_array(v)) {
        count = json_array_size(v);

        if (count == 0 || count > 8) {
            return false;
        }

        json_array_foreach(v, i, item) {
            hub = json_is_string(item) ? json_string_value(item) : NULL;
            if (!hub || !*hub) {
                return false;
            }
        }

        hubPayload = json_dumps(v, JSON_COMPACT);
        if (!hubPayload) {
            return false;
        }
    } else if (v != NULL) {
        return false;
    }

    if (spaceOut && space) {
        *spaceOut = strdup(space);
        if (!*spaceOut) {
            free(hubPayload);
            return false;
        }
    }

    if (nameOut) {
        *nameOut = strdup(name);
        if (!*nameOut) {
            free(hubPayload);
            free(spaceOut ? *spaceOut : NULL);
            if (spaceOut) *spaceOut = NULL;
            return false;
        }
    }

    if (tagOut) {
        *tagOut = strdup(tag);
        if (!*tagOut) {
            free(hubPayload);
            free(spaceOut ? *spaceOut : NULL);
            free(nameOut ?  *nameOut : NULL);

            if (spaceOut) *spaceOut = NULL;
            if (nameOut)  *nameOut  = NULL;

            return false;
        }
    }

    if (hubOut) {
        *hubOut = hubPayload;
    } else {
        free(hubPayload);
    }

    return true;
}

static int ws_update_cb(const struct _u_request *req,
                        struct _u_response *resp,
                        void *userData) {

    StarterContext *ctx;
    json_t *j;
    char *space;
    char *name;
    char *tag;
    char *hub;
    Action *a;

    ctx = (StarterContext *)userData;
    if (!ctx || !ctx->queue || !ctx->supervisor || !ctx->spaceList) {
        return ws_reply_text(resp,
                             HttpStatus_InternalServerError,
                             HttpStatusStr(HttpStatus_InternalServerError));
    }

    if (ctx->switchRequested ||
        ctx->terminateRequested ||
        ctx->restartRequested) {
        return ws_reply_text(resp,
                             HttpStatus_Conflict,
                             HttpStatusStr(HttpStatus_Conflict));
    }

    if (ctx->updateInProgress) {
        return ws_reply_text(resp,
                             HttpStatus_Locked,
                             HttpStatusStr(HttpStatus_Locked));
    }

    j = ws_load_json_body(req);
    if (!j) {
        return ws_reply_text(resp,
                             HttpStatus_BadRequest,
                             HttpStatusStr(HttpStatus_BadRequest));
    }

    space = NULL;
    name  = NULL;
    tag   = NULL;
    hub   = NULL;

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

    if (!ws_update_target_exists(ctx, space, name)) {
        json_decref(j);
        free(space);
        free(name);
        free(tag);
        free(hub);
        return ws_reply_text(resp,
                             HttpStatus_NotFound,
                             HttpStatusStr(HttpStatus_NotFound));
    }

    ctx->updateInProgress = 1;

    a = action_new(ACTION_UPDATE_APP, space, name, tag, hub);
    free(space);
    free(name);
    free(tag);
    free(hub);
    json_decref(j);

    if (!a || !actions_enqueue(ctx->queue, a)) {
        if (a) {
            free(a->space);
            free(a->name);
            free(a->tag);
            free(a->hub);
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
    json_t *j;
    json_t *v;
    const char *space;
    const char *name;
    Action *a;

    ctx = (StarterContext *)userData;
    if (!ctx || !ctx->queue || !ctx->supervisor || !ctx->spaceList) {
        return ws_reply_text(resp,
                             HttpStatus_InternalServerError,
                             HttpStatusStr(HttpStatus_InternalServerError));
    }

    if (ctx->switchRequested ||
        ctx->updateInProgress ||
        ctx->terminateRequested ||
        ctx->restartRequested) {
        return ws_reply_text(resp,
                             HttpStatus_Conflict,
                             HttpStatusStr(HttpStatus_Conflict));
    }

    j = ws_load_json_body(req);
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

    if (!ws_app_exists(ctx, space, name)) {
        json_decref(j);
        return ws_reply_text(resp,
                             HttpStatus_NotFound,
                             HttpStatusStr(HttpStatus_NotFound));
    }

    ctx->terminateRequested = 1;

    a = action_new(ACTION_TERMINATE_APP, space, name, NULL, NULL);
    json_decref(j);

    if (!a || !actions_enqueue(ctx->queue, a)) {
        if (a) {
            free(a->space);
            free(a->name);
            free(a->tag);
            free(a->hub);
            free(a);
        }
        ctx->terminateRequested = 0;
        return ws_reply_text(resp,
                             HttpStatus_InternalServerError,
                             HttpStatusStr(HttpStatus_InternalServerError));
    }

    supervisor_signal((Supervisor *)ctx->supervisor);

    return ws_reply_text(resp,
                         HttpStatus_Accepted,
                         HttpStatusStr(HttpStatus_Accepted));
}

static int ws_restart_cb(const struct _u_request *req,
                         struct _u_response *resp,
                         void *userData) {

    StarterContext *ctx;
    json_t *json;
    json_t *value;
    const char *space;
    const char *name;
    Action *action;

    ctx = (StarterContext *)userData;
    if (!ctx || !ctx->queue || !ctx->supervisor || !ctx->spaceList) {
        return ws_reply_text(resp,
                             HttpStatus_InternalServerError,
                             HttpStatusStr(HttpStatus_InternalServerError));
    }

    if (ctx->switchRequested ||
        ctx->updateInProgress ||
        ctx->terminateRequested ||
        ctx->restartRequested) {
        return ws_reply_text(resp,
                             HttpStatus_Conflict,
                             HttpStatusStr(HttpStatus_Conflict));
    }

    json = ws_load_json_body(req);
    if (!json) {
        return ws_reply_text(resp,
                             HttpStatus_BadRequest,
                             HttpStatusStr(HttpStatus_BadRequest));
    }

    value = json_object_get(json, "space");
    space = json_is_string(value) ? json_string_value(value) : NULL;
    value = json_object_get(json, "name");
    name = json_is_string(value) ? json_string_value(value) : NULL;

    if (!space || !name) {
        json_decref(json);
        return ws_reply_text(resp,
                             HttpStatus_BadRequest,
                             HttpStatusStr(HttpStatus_BadRequest));
    }

    if (!ws_app_exists(ctx, space, name)) {
        json_decref(json);
        return ws_reply_text(resp,
                             HttpStatus_NotFound,
                             HttpStatusStr(HttpStatus_NotFound));
    }

    ctx->restartRequested = 1;
    action = action_new(ACTION_RESTART_APP, space, name, NULL, NULL);
    json_decref(json);

    if (!action || !actions_enqueue(ctx->queue, action)) {
        if (action) {
            free(action->space);
            free(action->name);
            free(action);
        }
        ctx->restartRequested = 0;
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
                               "/v1", "/ready", 0,
                               &ws_ready_cb, ctx);
    setup_unsupported_methods(ctx->uInstance, "GET",
                              "/v1", "/ready");

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

    ulfius_add_endpoint_by_val(ctx->uInstance, "POST",
                               "/v1", "/restart", 0,
                               &ws_restart_cb, ctx);
    setup_unsupported_methods(ctx->uInstance, "POST",
                              "/v1", "/restart");

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
