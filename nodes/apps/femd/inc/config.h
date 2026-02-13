/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

#ifndef CONFIG_H
#define CONFIG_H

#include <stdint.h>
#include <stdbool.h>

#include "femd.h"

typedef struct {
    int  servicePort;
    int  nodedPort;
    char serviceName[64];

    char nodeID[64];
    char nodeType[32];
    
    char gpioBasePath[256];

    int  i2cBusFem1;
    int  i2cBusFem2;
    int  i2cBusCtrl;

    char safetyConfigPath[256];

    char notifyHost[64];
    int  notifyPort;
    char notifyPath[128];

    uint32_t samplePeriodMs;
    uint32_t safetyPeriodMs;

    bool enableWeb;
    bool enableSafety;
    bool enableNotify;
} Config;

typedef struct {
    Config *config;
} ServerConfig;

int  config_set_defaults(Config *cfg, const char *path);
int  config_load(Config *cfg, const char *path);
int  config_validate(const Config *cfg);
void config_print(const Config *cfg);

#endif /* CONFIG_H */
