/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef CONFIG_H
#define CONFIG_H

#include <stdint.h>
#include <stdbool.h>

typedef struct {
    /* Web service */
    char        *listenAddr;
    uint16_t    listenPort;

    /* Sampling */
    uint32_t    sampleMs;

    /* Driver selection */
    char        *driverName;        /* "victron", "epever", etc. */

    /* Serial port configuration */
    char        *serialPort;        /* e.g., /dev/ttyUSB0 */
    int         baudRate;

    /* Notify.d integration */
    char        notifyHost[64];
    int         notifyPort;
    char        notifyPath[128];
    bool        enableNotify;

    /* Alarm thresholds */
    double      lowVoltageWarn;     /* V */
    double      lowVoltageCrit;     /* V */
    double      highTempWarn;       /* °C */
    double      highTempCrit;       /* °C */

    /* Node identification */
    char        *nodeId;
} Config;

int  config_load_from_env(Config *config);
void config_log(const Config *config);
void config_print_env_help(void);
void config_free(Config *config);

#endif /* CONFIG_H */
