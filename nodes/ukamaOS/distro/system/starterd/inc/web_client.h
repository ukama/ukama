/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#pragma once

#include <stdbool.h>

#include "ulfius.h"

#include "config.h"
#include "app.h"

typedef struct _u_instance  UInst;
typedef struct _u_instance  UInst;
typedef struct _u_request   URequest;
typedef struct _u_response  UResponse;
typedef json_t              JsonObj;
typedef json_error_t        JsonErrObj;

typedef struct {
    int status;
    bool ready;
    char reason[STARTERD_READY_REASON_LEN];
    char requestId[STARTERD_REQUEST_ID_LEN];
} AppReadyResponse;

typedef enum {
    LIFECYCLE_GATE_UNAVAILABLE = 0,
    LIFECYCLE_GATE_WAITING,
    LIFECYCLE_GATE_OPEN,
    LIFECYCLE_GATE_FAULTY
} LifecycleGateState;

bool wc_lifecycle_check_in(Config *config, bool bootHealthy);
LifecycleGateState wc_lifecycle_gate(Config *config);

bool wc_app_ping(Config *config, App *app);
bool wc_app_version_matches(Config *config,
                            App *app,
                            const char *tag);
bool wc_app_ready(Config *config, App *app, AppReadyResponse *result);
bool wc_mesh_status(Config *config,
                    App *app,
                    bool *connected,
                    char *reason,
                    size_t reasonSize);
bool wc_fetch_package(Config *config,
                      const char *appName,
                      const char *tag,
                      const char *hub,
                      char **pathOut,
                      char **versionOut);
