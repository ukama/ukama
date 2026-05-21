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
#include <stddef.h>
#include <pthread.h>

#include "jansson.h"
#include "config.h"

#define STATUS_REASON_LEN 256

typedef json_t JsonObj;

typedef enum {
    InitStateStarting = 0,
    InitStateCheckTools,
    InitStateStartOvs,
    InitStateSetupEpcIf,
    InitStateSetupTun,
    InitStateSetupBridge,
    InitStateSetupForwarding,
    InitStateSetupFlows,
    InitStateReady,
    InitStateFailed
} InitState;

typedef struct {

    InitState state;
    bool ready;

    bool toolsOk;
    bool ovsdbRunning;
    bool vswitchdRunning;
    bool epcIfReady;
    bool tunReady;
    bool bridgeReady;
    bool forwardingReady;
    bool flowsReady;

    char reason[STATUS_REASON_LEN];

    pthread_mutex_t mutex;

} AppStatus;

void status_init(AppStatus *status);
void status_destroy(AppStatus *status);
void status_set(AppStatus *status, InitState state, const char *reason);
void status_fail(AppStatus *status, const char *reason);
bool status_is_ready(AppStatus *status);
const char *status_state_str(InitState state);
JsonObj *status_to_json(AppStatus *status, Config *config);

#endif /* STATUS_H_ */
