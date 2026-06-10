/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>

#include "response.h"

void emu_response_init(EmuResponse *response, const char *id)
{
    if (response == NULL) {
        return;
    }

    memset(response, 0, sizeof(EmuResponse));
    snprintf(response->id, sizeof(response->id), "%s", id ? id : "req");
}

void emu_response_free(EmuResponse *response)
{
    if (response == NULL) {
        return;
    }

    json_decref(response->payload);
    memset(response, 0, sizeof(EmuResponse));
}

bool emu_response_set_ok(EmuResponse *response, JsonObj *payload)
{
    if (response == NULL) {
        json_decref(payload);
        return false;
    }

    json_decref(response->payload);

    response->ok = true;
    snprintf(response->code, sizeof(response->code), "OK");
    response->reason[0] = '\0';
    response->payload = payload ? payload : json_object();

    return response->payload != NULL;
}

bool emu_response_set_error(EmuResponse *response,
                            const char *code,
                            const char *reason)
{
    if (response == NULL) {
        return false;
    }

    json_decref(response->payload);

    response->ok = false;
    snprintf(response->code,
             sizeof(response->code),
             "%s",
             code ? code : "Error");
    snprintf(response->reason,
             sizeof(response->reason),
             "%s",
             reason ? reason : response->code);

    response->payload = json_object();

    return response->payload != NULL;
}

JsonObj *emu_response_to_json(EmuResponse *response)
{
    JsonObj *json = NULL;
    JsonObj *payload = NULL;

    if (response == NULL) {
        return NULL;
    }

    payload = response->payload ? json_deep_copy(response->payload)
                                : json_object();
    if (payload == NULL) {
        return NULL;
    }

    json = json_object();
    if (json == NULL) {
        json_decref(payload);
        return NULL;
    }

    json_object_set_new(json, "id", json_string(response->id));
    json_object_set_new(json, "ok", json_boolean(response->ok));
    json_object_set_new(json, "code", json_string(response->code));
    json_object_set_new(json, "reason", json_string(response->reason));
    json_object_set_new(json, "payload", payload);

    return json;
}
