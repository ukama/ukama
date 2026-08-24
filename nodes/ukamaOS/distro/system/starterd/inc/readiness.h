/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#pragma once

#include <stddef.h>

#include <jansson.h>

#include "config.h"
#include "space.h"
#include "web_service.h"

typedef enum {
    NODE_READINESS_PENDING = 0,
    NODE_READINESS_READY,
    NODE_READINESS_FAULTY
} NodeReadinessState;

ReadinessMonitor *readiness_start(Config *config,
                                  Space *spaceList,
                                  StarterContext *ctx);
void readiness_stop(ReadinessMonitor *monitor);

NodeReadinessState readiness_get(ReadinessMonitor *monitor,
                                 char *reason,
                                 size_t reasonSize);
NodeReadinessState readiness_get_boot(ReadinessMonitor *monitor,
                                      char *reason,
                                      size_t reasonSize);
void readiness_boot_started(ReadinessMonitor *monitor);
void readiness_app_started(ReadinessMonitor *monitor, App *app);
void readiness_app_exited(ReadinessMonitor *monitor, App *app);
json_t *readiness_status_json(ReadinessMonitor *monitor);
json_t *readiness_connectivity_json(ReadinessMonitor *monitor);
