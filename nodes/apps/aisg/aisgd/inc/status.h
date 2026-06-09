/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef STATUS_H_
#define STATUS_H_

#include <stdbool.h>
#include <pthread.h>
#include "jansson.h"

#define STATUS_REASON_LEN                  256
#define AISGD_DEVICE_ID                    "ret-0"
#define STATUS_MAX_STR                     128

typedef json_t JsonObj;

typedef enum {
    AisgdStateStarting = 0,
    AisgdStateConnectController,
    AisgdStateScanDevice,
    AisgdStateSubscribeAlarms,
    AisgdStateVerifyConfig,
    AisgdStateReady,
    AisgdStateOperationRunning,
    AisgdStateDegraded,
    AisgdStateFailed
} AisgdState;

typedef struct {
    pthread_mutex_t mutex;
    AisgdState state;
    bool ready;
    char reason[STATUS_REASON_LEN];
    bool controllerConnected;
    char backend[STATUS_MAX_STR];
    char mode[STATUS_MAX_STR];
    bool powerManaged;
    bool devicePresent;
    bool configured;
    bool calibrated;
    bool busy;
    char model[STATUS_MAX_STR];
    bool operationActive;
    char operationType[STATUS_MAX_STR];
    char operationId[STATUS_MAX_STR];
} AppStatus;

void status_init(AppStatus *status);
void status_destroy(AppStatus *status);
void status_set(AppStatus *status, AisgdState state, const char *reason);
void status_set_operation(AppStatus *status, const char *type, const char *id);
void status_clear_operation(AppStatus *status);
void status_update_from_controller(AppStatus *status, JsonObj *payload);
void status_mark_controller_down(AppStatus *status, const char *reason);
void status_set_ready_if_idle(AppStatus *status, const char *reason);
bool status_is_ready(AppStatus *status);
JsonObj *status_to_json(AppStatus *status);

#endif /* STATUS_H_ */
