/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdlib.h>
#include <ulfius.h>

#include "httpStatus.h"

#include "usys_error.h"
#include "usys_types.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_file.h"
#include "usys_string.h"
#include "usys_services.h"

typedef struct _u_instance UInst;
typedef struct _u_request  URequest;
typedef struct _u_response UResponse;

static int callback_web_service_ping(const URequest *request,
                                     UResponse *response,
                                     void *data) {

    ulfius_set_string_body_response(response, HttpStatus_OK,
                                    HttpStatusStr(HttpStatus_OK));

    return U_CALLBACK_CONTINUE;
}

static int callback_web_service_default(const URequest *request,
                                        UResponse *response,
                                        void *epConfig) {

    ulfius_set_string_body_response(response,
                                    HttpStatus_NotFound,
                                    HttpStatusStr(HttpStatus_NotFound));

    return U_CALLBACK_CONTINUE;
}

int start_web_services(UInst *inst) {

    int port;

    port = usys_find_service_port(SERVICE_BOOTSTRAP);
    if (!port) {
        usys_log_error("Unable to find bootstrap service port in /etc/service");
        return USYS_FALSE;
    }

    if (ulfius_init_instance(inst, port, NULL, NULL) != U_OK) {
        usys_log_error("Error initializing instance for websocket remote port %d", port);
        return USYS_FALSE;
    }

    ulfius_add_endpoint_by_val(inst, "GET", "/v1/", "ping", 0,
                               &callback_web_service_ping, NULL);
    ulfius_set_default_endpoint(inst, &callback_web_service_default, NULL);

    if (ulfius_start_framework(inst) != U_OK) {
        usys_log_error("Error starting the webservice on port %d", port);
        ulfius_stop_framework(inst);
        ulfius_clean_instance(inst);

        return USYS_FALSE;
    }

    usys_log_debug("Web service on port: %d started.", port);

    return USYS_TRUE;
}

