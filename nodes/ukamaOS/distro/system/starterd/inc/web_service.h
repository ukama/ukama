/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#pragma once

#include <stdbool.h>
#include <signal.h>

#include "config.h"
#include "space.h"
#include "actions.h"

struct _u_instance;

typedef struct StarterContext {
    Config *config;
    Space *spaceList;
    ActionQueue *queue;
    void *supervisor;
    struct _u_instance *uInstance;

    volatile sig_atomic_t terminateRequested;
    volatile sig_atomic_t switchRequested;
    volatile sig_atomic_t updateInProgress;

    int exitCode;
} StarterContext;

bool web_service_start(StarterContext *ctx);
void web_service_stop(StarterContext *ctx);
