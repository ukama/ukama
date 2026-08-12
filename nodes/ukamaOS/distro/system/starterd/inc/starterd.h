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

#include "usys_services.h"

#define STARTERD_SERVICE_NAME          SERVICE_STARTER
#define STARTERD_DEFAULT_MANIFEST_FILE "/ukama/manifest.json"
#define STARTERD_DEFAULT_READY_FILE    "/ukama/init/starter/ready"
#define STARTERD_DEFAULT_RLOG_SOCKET \
    "/run/ukama/rlog-ingest.sock"
#define STARTERD_DEFAULT_LOG_SPOOL_DIR \
    "/ukama/state/starterd/log-spool"
#define STARTERD_DEFAULT_LOG_SPOOL_MAX_BYTES \
    (32 * 1024 * 1024)
#define STARTERD_DEFAULT_LOG_RECORD_MAX 32768
#define STARTERD_DEFAULT_LOG_RECONNECT_MS 1000
#define STARTERD_DEFAULT_READY_TIMEOUT_SEC 600
#define STARTERD_DEFAULT_READY_POLL_SEC 5
#define STARTERD_DEFAULT_LIFECYCLE_PORT 8097
#define STARTERD_READY_REASON_LEN 192
#define STARTERD_REQUEST_ID_LEN 96

typedef enum {
    APP_STATE_STOPPED = 0,
    APP_STATE_STARTING,
    APP_STATE_RUNNING,
    APP_STATE_STOPPING,
    APP_STATE_FAILED
} AppState;

typedef enum {
    APP_REASON_NONE = 0,
    APP_REASON_STARTED,
    APP_REASON_EXITED,
    APP_REASON_TERMINATED,
    APP_REASON_KILLED,
    APP_REASON_CRASHED,
    APP_REASON_UPDATE,
    APP_REASON_RESTART,
    APP_REASON_START_FAILED,
    APP_REASON_RESTART_FAILED,
    APP_REASON_PACKAGE_MISSING,
    APP_REASON_UNKNOWN
} AppReason;

typedef enum {
    INSTALL_STATE_NONE = 0,
    INSTALL_STATE_FETCHING,
    INSTALL_STATE_STAGING,
    INSTALL_STATE_SWITCHED,
    INSTALL_STATE_PENDING,
    INSTALL_STATE_FAILED
} InstallState;

typedef enum {
    APP_READINESS_IGNORED = 0,
    APP_READINESS_PENDING,
    APP_READINESS_READY,
    APP_READINESS_FAULTY
} AppReadinessState;
