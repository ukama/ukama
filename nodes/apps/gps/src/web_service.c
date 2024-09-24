/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#include "web_service.h"
#include "web_client.h"
#include "http_status.h"
#include "config.h"

#include "usys_error.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

#include "version.h"

/* global */
extern GPSData *gData; 

int web_service_cb_ping(const URequest *request,
                        UResponse *response,
                        void *epConfig) {

    ulfius_set_string_body_response(response, HttpStatus_OK,
                                    HttpStatusStr(HttpStatus_OK));

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_version(const URequest *request,
                           UResponse *response,
                           void *epConfig) {

    ulfius_set_string_body_response(response,
                                    HttpStatus_OK,
                                    VERSION);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_default(const URequest *request,
                           UResponse *response,
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

int web_service_cb_lock(const URequest *request,
                                UResponse *response,
                                void *epConfig) {

    if (gData->gpsLock == USYS_FALSE) {
        ulfius_set_string_body_response(response, HttpStatus_NotFound,
                                        HttpStatusStr(HttpStatus_NotFound));
    } else {
        ulfius_set_empty_body_response(response, HttpStatus_OK);
    }

    return U_CALLBACK_CONTINUE;
}


int web_service_cb_coordinates(const URequest *request,
                               UResponse *response,
                               void *epConfig) {

    if (gData->gpsLock == USYS_FALSE) {
        ulfius_set_string_body_response(response, HttpStatus_NotFound,
                                        HttpStatusStr(HttpStatus_NotFound));
    } else {
        ulfius_set_empty_body_response(response, HttpStatus_Accepted);
    }

    return U_CALLBACK_CONTINUE;
}


int web_service_cb_time(const URequest *request,
                        UResponse *response,
                        void *epConfig) {

    if (gData->gpsLock == USYS_FALSE) {
        ulfius_set_string_body_response(response, HttpStatus_NotFound,
                                        HttpStatusStr(HttpStatus_NotFound));
    } else {
        ulfius_set_empty_body_response(response, HttpStatus_Accepted);
    }
    
    ulfius_set_empty_body_response(response, HttpStatus_Accepted);

    return U_CALLBACK_CONTINUE;
}








