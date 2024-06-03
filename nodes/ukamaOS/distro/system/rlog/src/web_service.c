/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#include "http_status.h"
#include "usys_error.h"
#include "usys_log.h"

#include "rlogd.h"
#include "version.h"

extern ThreadData *gData;

int web_service_cb_ping(const URequest *request, UResponse *response,
                        void *data) {

    ulfius_set_string_body_response(response, HttpStatus_OK,
                                    HttpStatusStr(HttpStatus_OK));

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_version(const URequest *request, UResponse *response,
                           void *data) {

    ulfius_set_string_body_response(response, HttpStatus_OK, VERSION);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_get_level(const URequest *request, UResponse *response,
                             void *data) {

    char *level = NULL;

    if (gData->level == USYS_LOG_DEBUG) {
        level = "debug";
    } else if (gData->level == USYS_LOG_INFO) {
        level = "info";
    } else if (gData->level == USYS_LOG_ERROR) {
        level = "error";
    }

    ulfius_set_string_body_response(response, HttpStatus_OK, level);

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

int web_service_cb_default(const URequest *request, UResponse *response,
                           void *data) {

    ulfius_set_string_body_response(response, HttpStatus_NotFound,
                                    HttpStatusStr(HttpStatus_NotFound));

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_post_level(const URequest *request, UResponse *response,
                              void *data) {

    const char *level=NULL;

    level = u_map_get(request->map_url, "level");
    if (level == NULL) {
        ulfius_set_string_body_response(response,
                                        HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
        return U_CALLBACK_CONTINUE;
    }

    if (strcasecmp(level, "debug") == 0) {
        gData->level = USYS_LOG_DEBUG;
    } else if (strcasecmp(level, "info") == 0) {
        gData->level = USYS_LOG_INFO;
    } else if (strcasecmp(level, "error") == 0) {
        gData->level = USYS_LOG_ERROR;
    }
    
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
