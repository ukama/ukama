/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
#ifndef BACKHAULD_H
#define BACKHAULD_H

#include <stddef.h>
#include <time.h>

#include "usys_api.h"

#define SERVICE_NAME      "backhaul.d"
#define DEF_LOG_LEVEL     "INFO"

#define URL_PREFIX        "/v1"
#define API_RES_EP(x)     "/" x

/* Common OK/NOK */
#ifndef STATUS_OK
#define STATUS_OK   0
#endif
#ifndef STATUS_NOK
#define STATUS_NOK  1
#endif

typedef enum {
    BACKHAUL_STATE_UNKNOWN = 0,
    BACKHAUL_STATE_GOOD,
    BACKHAUL_STATE_DEGRADED,
    BACKHAUL_STATE_DOWN,
    BACKHAUL_STATE_CAPPED
} BackhaulState;

typedef enum {
    BACKHAUL_LINK_UNKNOWN = 0,
    BACKHAUL_LINK_TERRESTRIAL_LIKE,
    BACKHAUL_LINK_SAT_LEO_LIKE,
    BACKHAUL_LINK_SAT_GEO_LIKE,
    BACKHAUL_LINK_CELLULAR_LIKE
} BackhaulLinkGuess;

typedef struct {
    int     ok;
    long    httpCode;
    double  ttfbMs;
    double  totalMs;
    int     stalled;
} ProbeResult;

typedef struct {
    int     ok;
    long    httpCode;
    double  seconds;
    double  mbps;
} TransferResult;

typedef struct {
    long    ts;
    int     ok;
    double  ttfbMs;
    int     stalled;
} MicroSample;

typedef struct {
    long    ts;
    int     ok;
    double  dlMbps;
    double  ulMbps;
    double  dlSec;
    double  ulSec;
} ChgSample;

typedef struct {

    BackhaulState      backhaulState;
    BackhaulLinkGuess  linkGuess;
    double             confidence;

    /* classifier hysteresis (internal state, useful for debugging) */
    int                consecFails;
    int                consecOk;

    double             dlGoodputMbps;
    double             ulGoodputMbps;

    double             bufferbloatInflationFactor;
    double             capDetectedMbps;

    double             nearTtfbMedianMs;
    double             nearTtfbP95Ms;
    double             nearTtfbP99Ms;

    double             farTtfbMedianMs;
    double             farTtfbP95Ms;
    double             farTtfbP99Ms;

    double             probeSuccessRatePct;
    double             stallRatePct;

    time_t             lastMicroTs;
    time_t             lastMultiTs;
    time_t             lastChgTs;
    time_t             lastClassifyTs;

    time_t             lastDiagTs;
    char               lastDiagName[32];

} BackhaulMetrics;

typedef struct {
    char nearUrl[256];
    char farUrl[256];
    long ts;
} ReflectorSet;

/* Forward decls */
typedef struct Config Config;
typedef struct MetricsStore MetricsStore;

typedef struct {
    Config       *config;
    MetricsStore *store;
} EpCtx;

#endif /* BACKHAULD_H */
