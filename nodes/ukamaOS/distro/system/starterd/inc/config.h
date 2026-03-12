/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#pragma once

#include <stdbool.h>
#include <stdint.h>

typedef struct {

    char *manifestPath;
    char *logPath;
    char *readyFile;

    char *appsRoot;
    char *pkgsDir;
    char *stateDir;

    char *httpAddr;
    int  httpPort;

    char *wimcHost;
    int  wimcPort;
    char *wimcPathTemplate;

    int commitTimeoutSec;
    int pingTimeoutSec;

    int termGraceSec;

    int restartMaxBackoffSec;
    int restartStableResetSec;

    char *bootSpace;

} Config;

bool config_load(Config *config);
void config_free(Config *config);
