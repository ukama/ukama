/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "web_service.h"
#include "http_status.h"
#include "version.h"

#define EMU_URL_PREFIX "/v1"

static int cb_ping(const URequest *request, UResponse *response, void *data)
{
    (void)request;
    (void)data;

    ulfius_set_string_body_response(response, HttpStatus_OK, "OK");

    return U_CALLBACK_CONTINUE;
}

static int cb_version(const URequest *request, UResponse *response, void *data)
{
    (void)request;
    (void)data;

    ulfius_set_string_body_response(response, HttpStatus_OK, VERSION);

    return U_CALLBACK_CONTINUE;
}

static int cb_status(const URequest *request, UResponse *response, void *data)
{
    JsonObj *json = NULL;

    (void)request;

    json = emu_model_status((EmuModel *)data);
    ulfius_set_json_body_response(response, HttpStatus_OK, json);
    json_decref(json);

    return U_CALLBACK_CONTINUE;
}

static bool add_endpoint(UInst *instance,
                         const char *path,
                         int (*cb)(const URequest *, UResponse *, void *),
                         EmuModel *model)
{
    int ret;

    ret = ulfius_add_endpoint_by_val(instance,
                                     "GET",
                                     EMU_URL_PREFIX,
                                     path,
                                     0,
                                     cb,
                                     model);

    return ret == U_OK;
}

bool start_web_service(UInst *instance, EmuConfig *config, EmuModel *model)
{
    if (instance == NULL || config == NULL || model == NULL) {
        return false;
    }

    if (ulfius_init_instance(instance,
                             config->servicePort,
                             NULL,
                             NULL) != U_OK) {
        return false;
    }

    if (!add_endpoint(instance, "/ping", cb_ping, model)) {
        return false;
    }

    if (!add_endpoint(instance, "/version", cb_version, model)) {
        return false;
    }

    if (!add_endpoint(instance, "/status", cb_status, model)) {
        return false;
    }

    if (ulfius_start_framework(instance) != U_OK) {
        return false;
    }

    return true;
}

void stop_web_service(UInst *instance)
{
    if (instance == NULL) {
        return;
    }

    ulfius_stop_framework(instance);
    ulfius_clean_instance(instance);
}
