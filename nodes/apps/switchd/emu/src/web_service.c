/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "http_status.h"
#include "jserdes.h"
#include "model.h"
#include "scenario.h"
#include "utils.h"
#include "web_service.h"

#include "version.h"

static int port_id_from_path(const char *path) {

    const char *ptr;

    ptr = strstr(path, "/debug/ports/");
    if (ptr == NULL) {
        return -1;
    }

    ptr += strlen("/debug/ports/");
    return atoi(ptr);
}

int web_service_handle(EmuModel *model,
                       const char *method,
                       const char *path,
                       const char *body,
                       char *out,
                       size_t outLen,
                       int *status) {

    int id;
    char value[16];
    int reachable;

    id        = -1;
    reachable = 1;

    memset(value, 0, sizeof(value));

    if (!method || !path || !out || outLen == 0 || !status) {
        return STATUS_NOK;
    }

    *status = HttpStatus_OK;
    out[0] = '\0';

    /*
     * Ukama app lifecycle endpoints.
     *
     * starter.d requires these after launching an app:
     *   GET /v1/ping     -> 200 OK
     *   GET /v1/version  -> body must match manifest tag, usually "latest"
     */
    if (strcmp(method, "GET") == 0 && strcmp(path, "/v1/ping") == 0) {
        snprintf(out, outLen, "OK");
        return STATUS_OK;
    }

    if (strcmp(method, "GET") == 0 && strcmp(path, "/v1/version") == 0) {
        snprintf(out, outLen, "%s", VERSION);
        return STATUS_OK;
    }

    if (!model) {
        *status = HttpStatus_InternalServerError;
        json_serialize_error("modelNotAvailable", out, outLen);
        return STATUS_OK;
    }

    pthread_mutex_lock(&model->lock);

    if (strcmp(method, "GET") == 0 && strcmp(path, "/debug/state") == 0) {
        json_serialize_state(model, out, outLen);

    } else if (strcmp(method, "GET") == 0 &&
               strcmp(path, "/debug/ports") == 0) {
        json_serialize_ports(model, out, outLen);

    } else if (strcmp(method, "GET") == 0 &&
               strcmp(path, "/debug/firmware") == 0) {
        json_serialize_firmware(model, out, outLen);

    } else if (strcmp(method, "POST") == 0 &&
               strcmp(path, "/debug/scenario") == 0) {
        if (json_get_string_field(body, "name",
                                  value, sizeof(value)) != STATUS_OK) {
            *status = HttpStatus_BadRequest;
            json_serialize_error("invalidBody", out, outLen);
        } else {
            scenario_apply(model, value);
            json_serialize_result_ok(out, outLen);
        }

    } else if (strcmp(method, "POST") == 0 &&
               strcmp(path, "/debug/switch/reachable") == 0) {
        if (json_get_bool_field(body, "value", &reachable) != STATUS_OK) {
            *status = HttpStatus_BadRequest;
            json_serialize_error("invalidBody", out, outLen);
        } else {
            model_set_reachable(model, reachable);
            json_serialize_result_ok(out, outLen);
        }

    } else if (strcmp(method, "POST") == 0 &&
               strstr(path, "/debug/ports/") == path &&
               strstr(path, "/link") != NULL) {
        id = port_id_from_path(path);

        if (json_get_string_field(body, "value",
                                  value, sizeof(value)) != STATUS_OK ||
            model_set_port_link(model,
                                (unsigned int)id,
                                strcmp(value, "up") == 0) != STATUS_OK) {
            *status = HttpStatus_BadRequest;
            json_serialize_error("badRequest", out, outLen);
        } else {
            json_serialize_result_ok(out, outLen);
        }

    } else if (strcmp(method, "POST") == 0 &&
               strstr(path, "/debug/ports/") == path &&
               strstr(path, "/poe") != NULL) {
        id = port_id_from_path(path);

        if (json_get_string_field(body, "value",
                                  value, sizeof(value)) != STATUS_OK ||
            model_set_port_poe(model,
                               (unsigned int)id,
                               strcmp(value, "on") == 0) != STATUS_OK) {
            *status = HttpStatus_BadRequest;
            json_serialize_error("badRequest", out, outLen);
        } else {
            json_serialize_result_ok(out, outLen);
        }

    } else {
        *status = HttpStatus_NotFound;
        json_serialize_error("notFound", out, outLen);
    }

    pthread_mutex_unlock(&model->lock);

    return STATUS_OK;
}
