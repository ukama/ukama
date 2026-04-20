/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef __POWER_COLLECTOR_H__
#define __POWER_COLLECTOR_H__

#include <stdint.h>

#include "config.h"
#include "metrics_store.h"
#include "power_kpi.h"
#include "drv_lm75.h"
#include "drv_ads1015.h"
#include "drv_lm25066.h"

typedef struct {
    MetricsStore *store;
    const Config *config;
    PowerCal     *cal;

    Lm75         *lm75_board;
    Ads1015      *ads1015;
    Lm25066      *lm25066;

    int          mockMode;
} PowerCollectorCtx;

int power_collect_once(PowerCollectorCtx *c, uint64_t now_ms);

#endif
