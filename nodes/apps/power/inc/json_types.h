/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef JSON_TYPES_H
#define JSON_TYPES_H

#include <stdint.h>

typedef struct {
    uint64_t    sampleUnixMs;
    char        board[16];

    int         ok;
    char        err[128];

    char        severity[16];
    char        reason[128];

    int         haveLm25066;
    double      inVolts;
    double      outVolts;
    double      inAmps;
    double      inWatts;
    double      hsTempC;
    uint16_t    statusWord;
    uint16_t    diagnosticWord;
    int         assumedDirect;

    int         haveLm75;
    double      boardTempC;

    int         haveAds1015;
    double      adcVin;
    double      adcVpa;
    double      adcAux;

    double      totalWatts;
    double      energyWh;
} PowerMetrics;

#endif
