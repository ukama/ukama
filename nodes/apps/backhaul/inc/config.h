/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef CONFIG_H
#define CONFIG_H

#include "usys_api.h"

typedef struct Config {
    char *serviceName;
    int   servicePort;

    char *bootstrapHost;
    char *bootstrapScheme;
    char *bootstrapEp;

    char *reflectorNearUrl;
    char *reflectorFarUrl;

    int reflectorRefreshSec; /* kept but unused in simplified design */

    int microPeriodMs;
    int multiPeriodMs;
    int chgPeriodSec;
    int classifyPeriodSec;

    int connectTimeoutMs;
    int totalTimeoutMs;
    int maxRetries;

    int pingBytes;
    int stallThresholdMs;

    int chgTargetSec;
    int chgWarmupBytes;
    int chgMinBytes;
    int chgMaxBytes;
    int chgSamples;

    int parallelStreams;
    int parallelMaxBytesTotal;

    int downConsecFails;
    int recoverConsecOk;
    int capStabilityPct;

    int windowMicroSamples;
    int windowMultiSamples;
    int windowChgSamples;
} Config;

int  config_load_from_env(Config *config);
int  config_validate_env(Config *config);
void config_log(Config *config);
void config_print_env_help(void);
void config_free(Config *config);

#endif /* CONFIG_H_ */
