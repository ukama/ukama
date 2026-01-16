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

    bool locked = false;

    (void)request;
    (void)epConfig;

    if (gData == NULL) {
        ulfius_set_string_body_response(response, HttpStatus_NotFound,
                                        HttpStatusStr(HttpStatus_NotFound));
        return U_CALLBACK_CONTINUE;
    }

    pthread_mutex_lock(&gData->mutex);
    locked = gData->gpsLock;
    pthread_mutex_unlock(&gData->mutex);

    if (!locked) {
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

    bool locked = false;
    const char *lat = NULL;
    const char *lon = NULL;
    char body[128] = {0};

    (void)request;
    (void)epConfig;

    if (gData == NULL) {
        ulfius_set_string_body_response(response, HttpStatus_NotFound,
                                        HttpStatusStr(HttpStatus_NotFound));
        return U_CALLBACK_CONTINUE;
    }

    pthread_mutex_lock(&gData->mutex);
    locked = gData->gpsLock;
    lat    = gData->latitude;
    lon    = gData->longitude;
    pthread_mutex_unlock(&gData->mutex);

    if (!locked || lat == NULL || lon == NULL || lat[0] == '\0' || lon[0] == '\0') {
        ulfius_set_string_body_response(response, HttpStatus_NotFound,
                                        HttpStatusStr(HttpStatus_NotFound));
        return U_CALLBACK_CONTINUE;
    }

    /* Return "longitude,latitude" */
    snprintf(body, sizeof(body), "%s,%s", lon, lat);
    ulfius_set_string_body_response(response, HttpStatus_OK, body);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_time(const URequest *request,
                        UResponse *response,
                        void *epConfig) {

    bool locked = false;
    const char *t = NULL;

    (void)request;
    (void)epConfig;

    if (gData == NULL) {
        ulfius_set_string_body_response(response, HttpStatus_NotFound,
                                        HttpStatusStr(HttpStatus_NotFound));
        return U_CALLBACK_CONTINUE;
    }

    pthread_mutex_lock(&gData->mutex);
    locked = gData->gpsLock;
    t      = gData->time;
    pthread_mutex_unlock(&gData->mutex);

    if (!locked || t == NULL || t[0] == '\0') {
        ulfius_set_string_body_response(response, HttpStatus_NotFound,
                                        HttpStatusStr(HttpStatus_NotFound));
    } else {
        ulfius_set_string_body_response(response, HttpStatus_OK, t);
    }

    return U_CALLBACK_CONTINUE;
}
