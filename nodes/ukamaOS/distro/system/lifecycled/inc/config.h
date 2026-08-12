/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#pragma once

#include <stdbool.h>

typedef struct {
    char *httpAddr;
    int httpPort;

    char *starterHost;
    int starterPort;

    char *notifyHost;
    int notifyPort;

    char *stateFile;

    int checkInTimeoutSec;
    int configTimeoutSec;
    int starterUnavailableTimeoutSec;
    int pollIntervalMs;
    int requestTimeoutSec;
} Config;

bool config_load(Config *config);
void config_free(Config *config);

