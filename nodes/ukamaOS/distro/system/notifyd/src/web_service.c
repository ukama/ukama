/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#include "web_service.h"
#include "web_client.h"
#include "httpStatus.h"
#include "jserdes.h"
#include "notification.h"
#include "service.h"

#include "usys_error.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

#include "version.h"

extern ThreadData *gData;

int web_service_cb_ping(const URequest *request, UResponse *response,
                        void *epConfig) {

    ulfius_set_string_body_response(response, HttpStatus_OK,
                                    HttpStatusStr(HttpStatus_OK));

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_version(const URequest *request, UResponse *response,
                           void *epConfig) {

    ulfius_set_string_body_response(response, HttpStatus_OK, VERSION);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_default(const URequest *request, UResponse *response,
                           void *epConfig) {
    
    ulfius_set_string_body_response(response, HttpStatus_NotFound,
                                    HttpStatusStr(HttpStatus_NotFound));

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_not_allowed(const URequest *request,
                               UResponse *response,
                               void *user_data) {

    ulfius_set_string_body_response(response,
                                    HttpStatus_MethodNotAllowed,
                                    HttpStatusStr(HttpStatus_MethodNotAllowed));
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_post_event(const URequest *request,
                              UResponse *response, void *epConfig) {

    int ret = STATUS_NOK;
    const char *service=NULL;
    JsonObj *json=NULL;

    service = u_map_get(request->map_url, "service");
    json = ulfius_get_json_body_request(request, NULL);
    if (service == NULL || json == NULL) {
        ulfius_set_string_body_response(response,
                                        HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
        return U_CALLBACK_CONTINUE;
    }
    usys_log_trace("notify.d:: Received POST for an event from %s.", service);

    ret = process_incoming_notification(service,
                                        NOTIFICATION_EVENT,
                                        json,
                                        (Config *)epConfig);
    if (ret == STATUS_OK) {
        ulfius_set_empty_body_response(response, HttpStatus_Accepted);
    } else {
        ulfius_set_empty_body_response(response,
                                       HttpStatus_InternalServerError);
    }

    json_free(&json);
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_post_alert(const URequest *request,
                              UResponse *response, void *epConfig) {

    int ret = STATUS_NOK;
    const char *service=NULL;
    JsonObj *json=NULL;

    service = u_map_get(request->map_url, "service");
    json = ulfius_get_json_body_request(request, NULL);
    if (service == NULL || json == NULL) {
        ulfius_set_string_body_response(response,
                                        HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
        return U_CALLBACK_CONTINUE;
    }
    usys_log_trace("notify.d:: Received POST for an event from %s.", service);

    ret = process_incoming_notification(service,
                                        NOTIFICATION_ALERT,
                                        json,
                                        (Config *)epConfig);
    if (ret == STATUS_OK) {
        ulfius_set_empty_body_response(response, HttpStatus_Accepted);
    } else {
        ulfius_set_empty_body_response(response,
                                       HttpStatus_InternalServerError);
    }

    json_free(&json);
    return U_CALLBACK_CONTINUE;
}

int web_service_cb_get_output(const URequest *request, UResponse *response,
                              void *data) {

    char *output = NULL;

    if (gData->output == STDOUT) {
        output = "stdout";
    } else if (gData->output == STDERR) {
        output = "stderr";
    } else if (gData->output == LOG_FILE) {
        output = "file";
    } else if (gData->output == UKAMA_SERVICE) {
        output = "ukama";
    }

    ulfius_set_string_body_response(response, HttpStatus_OK, output);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_post_output(const URequest *request, UResponse *response,
                               void *data) {

    const char *output=NULL;

    output = u_map_get(request->map_url, "output");
    if (output == NULL) {
        ulfius_set_string_body_response(response,
                                        HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
        return U_CALLBACK_CONTINUE;
    }

    if (strcasecmp(output, "stdout") == 0) {
        gData->output = STDOUT;
    } else if (strcasecmp(output, "stderr") == 0) {
        gData->output = STDERR;
    } else if (strcasecmp(output, "file") == 0) {
        gData->output = LOG_FILE;
    } else if (strcasecmp(output, "ukama") == 0) {
        gData->output = UKAMA_SERVICE;
    }

    return U_CALLBACK_CONTINUE;
}
