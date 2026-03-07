/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#pragma once

#include <stdbool.h>
#include <stdint.h>

#define STARTERD_SERVICE_NAME          "starter"
#define STARTERD_DEFAULT_MANIFEST_FILE "/manifest.json"
#define STARTERD_DEFAULT_LOG_PATH      "/ukama/apps.log"
#define STARTERD_DEFAULT_READY_FILE    "/ukama/init/starter/ready"

typedef enum {
    APP_STATE_STOPPED = 0,
    APP_STATE_STARTING,
    APP_STATE_RUNNING,
    APP_STATE_STOPPING,
    APP_STATE_FAILED
} AppState;

typedef enum {
    INSTALL_STATE_NONE = 0,
    INSTALL_STATE_FETCHING,
    INSTALL_STATE_STAGING,
    INSTALL_STATE_SWITCHED,
    INSTALL_STATE_FAILED
} InstallState;
