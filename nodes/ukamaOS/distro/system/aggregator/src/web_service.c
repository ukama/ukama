/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdlib.h>

#include "aggregator.h"
#include "httpStatus.h"
#include "web_service.h"

#include "usys_log.h"

#include "version.h"

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

    ulfius_set_string_body_response(response,
                                    HttpStatus_OK,
                                    VERSION);
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_status(const URequest *request,
                          UResponse *response,
                          void *data) {

    char *json = NULL;

    (void)request;

    json = app_state_status_json((AppState *)data);
    if (json == NULL) {
        ulfius_set_string_body_response(response,
                                        HttpStatus_InternalServerError,
                                        HttpStatusStr(HttpStatus_InternalServerError));
        return U_CALLBACK_CONTINUE;
    }

    u_map_put(response->map_header, "Content-Type", "application/json");
    ulfius_set_string_body_response(response,
                                    HttpStatus_OK,
                                    json);
    free(json);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_metrics(const URequest *request,
                           UResponse *response,
                           void *data) {

    char *snapshot = NULL;

    (void)request;

    snapshot = app_state_dup_snapshot((AppState *)data);
    if (snapshot == NULL) {
        ulfius_set_string_body_response(response,
                                        HttpStatus_InternalServerError,
                                        HttpStatusStr(HttpStatus_InternalServerError));
        return U_CALLBACK_CONTINUE;
    }

    u_map_put(response->map_header, "Content-Type", "text/plain; version=0.0.4");
    ulfius_set_string_body_response(response,
                                    HttpStatus_OK,
                                    snapshot);
    free(snapshot);

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
