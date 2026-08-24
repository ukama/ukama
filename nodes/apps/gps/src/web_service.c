/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#include <string.h>

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

int web_service_cb_ready(const URequest *request,
                         UResponse *response,
                         void *epConfig) {

    JsonObj *json;
    bool locked;
    time_t lockLostAt;
    time_t lastLockAt;
    time_t now;
    int lockTimeoutSec;
    int status;
    const char *reason;

    (void)request;
    (void)epConfig;

    if (gData == NULL) {
        ulfius_set_string_body_response(
            response,
            HttpStatus_ServiceUnavailable,
            HttpStatusStr(HttpStatus_ServiceUnavailable));
        return U_CALLBACK_CONTINUE;
    }

    pthread_mutex_lock(&gData->mutex);
    locked = gData->gpsLock;
    lockLostAt = gData->lockLostAt;
    lastLockAt = gData->lastLockAt;
    lockTimeoutSec = gData->lockTimeoutSec;
    pthread_mutex_unlock(&gData->mutex);

    now = time(NULL);
    status = HttpStatus_OK;
    reason = "ready";

    if (locked &&
        (lastLockAt == 0 ||
         now - lastLockAt > (GPS_WAIT_TIME * 3))) {
        locked = false;
        lockLostAt = lastLockAt;
        reason = "GPS collection is stale";
    }

    if (!locked &&
        lockLostAt > 0 &&
        now - lockLostAt >= lockTimeoutSec) {
        status = HttpStatus_ServiceUnavailable;
        reason = "GPS lock timeout";
    } else if (!locked) {
        status = HttpStatus_Accepted;
        if (strcmp(reason, "ready") == 0) {
            reason = "waiting for GPS lock";
        }
    }

    json = json_object();
    if (!json) {
        ulfius_set_string_body_response(
            response,
            HttpStatus_InternalServerError,
            HttpStatusStr(HttpStatus_InternalServerError));
        return U_CALLBACK_CONTINUE;
    }

    json_object_set_new(json, "ready", json_boolean(locked));
    if (!locked) {
        json_object_set_new(json, "reason", json_string(reason));
    }

    ulfius_set_json_body_response(response, status, json);
    json_decref(json);
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

int web_service_cb_status(const URequest *request,
                          UResponse *response,
                          void *epConfig) {

    JsonObj *json = NULL;
    bool locked = false;
    time_t lastLockAt = 0;
    time_t now = 0;
    char *lat = NULL;
    char *lon = NULL;
    char *gpsTime = NULL;
    char coordinates[MAX_BUFFER] = {0};

    (void)request;
    (void)epConfig;

    if (gData == NULL) {
        ulfius_set_string_body_response(
            response,
            HttpStatus_ServiceUnavailable,
            HttpStatusStr(HttpStatus_ServiceUnavailable));
        return U_CALLBACK_CONTINUE;
    }

    pthread_mutex_lock(&gData->mutex);
    locked = gData->gpsLock;
    lastLockAt = gData->lastLockAt;
    lat = gData->latitude ? usys_strdup(gData->latitude) : NULL;
    lon = gData->longitude ? usys_strdup(gData->longitude) : NULL;
    gpsTime = gData->time ? usys_strdup(gData->time) : NULL;
    pthread_mutex_unlock(&gData->mutex);

    now = time(NULL);
    if (locked &&
        (lastLockAt == 0 || now - lastLockAt > (GPS_WAIT_TIME * 3))) {
        locked = false;
    }

    json = json_object();
    if (json == NULL) {
        ulfius_set_string_body_response(
            response,
            HttpStatus_InternalServerError,
            HttpStatusStr(HttpStatus_InternalServerError));
        goto cleanup;
    }

    json_object_set_new(json, "lock", json_boolean(locked));

    if (lat && lat[0] != '\0' && lon && lon[0] != '\0') {
        snprintf(coordinates, sizeof(coordinates), "%s,%s", lat, lon);
        json_object_set_new(json,
                            "coordinates",
                            json_string(coordinates));
    } else {
        json_object_set_new(json, "coordinates", json_null());
    }

    if (gpsTime && gpsTime[0] != '\0') {
        json_object_set_new(json, "time", json_string(gpsTime));
    } else {
        json_object_set_new(json, "time", json_null());
    }

    ulfius_set_json_body_response(response, HttpStatus_OK, json);

cleanup:
    json_decref(json);
    usys_free(lat);
    usys_free(lon);
    usys_free(gpsTime);

    return U_CALLBACK_CONTINUE;
}
