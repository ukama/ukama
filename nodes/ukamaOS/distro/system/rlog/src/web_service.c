/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#include <strings.h>

#include "http_status.h"
#include "rlogd.h"
#include "usys_log.h"
#include "version.h"

extern ThreadData *gData;

static int level_from_string(const char *value) {
    if (!value) return -1;

    if (strcasecmp(value, "trace") == 0) return USYS_LOG_TRACE;
    if (strcasecmp(value, "debug") == 0) return USYS_LOG_DEBUG;
    if (strcasecmp(value, "info") == 0) return USYS_LOG_INFO;
    if (strcasecmp(value, "warn") == 0) return USYS_LOG_WARN;
    if (strcasecmp(value, "error") == 0) return USYS_LOG_ERROR;
    if (strcasecmp(value, "critical") == 0 ||
        strcasecmp(value, "fatal") == 0) {
        return USYS_LOG_CRITICAL;
    }

    return -1;
}

static const char *level_to_string(int level) {
    switch (level) {
    case USYS_LOG_TRACE:
        return "trace";
    case USYS_LOG_DEBUG:
        return "debug";
    case USYS_LOG_INFO:
        return "info";
    case USYS_LOG_WARN:
        return "warn";
    case USYS_LOG_ERROR:
        return "error";
    case USYS_LOG_CRITICAL:
        return "critical";
    default:
        return "unknown";
    }
}

int web_service_cb_ping(const URequest *request,
                        UResponse *response,
                        void *data) {
    (void)request;
    (void)data;

    ulfius_set_string_body_response(response, HttpStatus_OK,
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
    json_t *status;

    (void)request;
    (void)data;

    if (!gData || !gData->store || !gData->ingest) {
        ulfius_set_string_body_response(
            response, HttpStatus_ServiceUnavailable,
            HttpStatusStr(HttpStatus_ServiceUnavailable));
        return U_CALLBACK_CONTINUE;
    }

    status = json_pack("{s:b,s:s,s:s,s:I,s:I,s:s,s:s}",
                       "ready", 1,
                       "nodeId", log_store_node_id(gData->store),
                       "bootId", log_store_boot_id(gData->store),
                       "currentSeq",
                       (json_int_t)log_store_current_seq(gData->store),
                       "activeBytes",
                       (json_int_t)log_store_active_bytes(gData->store),
                       "ingestSocket",
                       ingest_socket_path(gData->ingest),
                       "level", level_to_string(gData->level));
    if (!status) {
        ulfius_set_string_body_response(
            response, HttpStatus_InternalServerError,
            HttpStatusStr(HttpStatus_InternalServerError));
        return U_CALLBACK_CONTINUE;
    }

    ulfius_set_json_body_response(response, HttpStatus_OK, status);
    json_decref(status);
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_get_level(const URequest *request,
                             UResponse *response,
                             void *data) {
    (void)request;
    (void)data;

    if (!gData) {
        ulfius_set_string_body_response(
            response, HttpStatus_ServiceUnavailable,
            HttpStatusStr(HttpStatus_ServiceUnavailable));
        return U_CALLBACK_CONTINUE;
    }

    ulfius_set_string_body_response(response, HttpStatus_OK,
                                    level_to_string(gData->level));
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_post_level(const URequest *request,
                              UResponse *response,
                              void *data) {
    const char *value;
    int level;

    (void)data;

    value = u_map_get(request->map_url, "level");
    level = level_from_string(value);
    if (level < 0) {
        ulfius_set_string_body_response(
            response, HttpStatus_BadRequest,
            HttpStatusStr(HttpStatus_BadRequest));
        return U_CALLBACK_CONTINUE;
    }

    gData->level = level;
    usys_log_set_level(level);
    ulfius_set_string_body_response(response, HttpStatus_OK,
                                    level_to_string(level));
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_not_allowed(const URequest *request,
                               UResponse *response,
                               void *data) {
    (void)request;
    (void)data;

    ulfius_set_string_body_response(
        response, HttpStatus_MethodNotAllowed,
        HttpStatusStr(HttpStatus_MethodNotAllowed));
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_default(const URequest *request,
                           UResponse *response,
                           void *data) {
    (void)request;
    (void)data;

    ulfius_set_string_body_response(response, HttpStatus_NotFound,
                                    HttpStatusStr(HttpStatus_NotFound));
    return U_CALLBACK_CONTINUE;
}
