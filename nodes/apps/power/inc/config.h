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

typedef struct {
    char        *listenAddr;
    uint16_t    listenPort;
    uint32_t    sampleMs;
    char        *boardName;

    int         mockMode;

    /* LM25066 */
    char        *lm25066Dev;
    int         lm25066Addr;
    int         lm25066ClHigh;
    int         lm25066RsMohm;

    /* LM75 */
    char        *lm75Dev;
    int         lm75Addr;

    /* ADS1015 */
    char        *ads1015Dev;
    int         ads1015Addr;
    int         adsChVin;
    int         adsChVpa;
    int         adsChAux;
} Config;

int  config_load_from_env(Config *config);
int  config_validate_env(Config *config);
void config_log(Config *config);
void config_print_env_help(void);
void config_free(Config *config);

#endif
