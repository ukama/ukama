/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef CONTROLLER_CLIENT_H_
#define CONTROLLER_CLIENT_H_

#include <stdbool.h>
#include "jansson.h"
#include "config.h"

typedef json_t JsonObj;

typedef struct {
    char path[AISGD_MAX_STR];
    int timeoutMs;
} ControllerClient;

typedef struct {
    bool ok;
    char code[AISGD_MAX_STR];
    char reason[AISGD_MAX_STR];
    JsonObj *payload;
} CtrlResponse;

void controller_client_init(ControllerClient *client, Config *config);
void ctrl_response_free(CtrlResponse *response);
JsonObj *ctrl_response_steal_payload(CtrlResponse *response);
bool controller_client_call(ControllerClient *client,
                            const char *type,
                            JsonObj *payload,
                            CtrlResponse *response);

#endif /* CONTROLLER_CLIENT_H_ */
