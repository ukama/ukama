/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef __SAMPLE_LOOP_H__
#define __SAMPLE_LOOP_H__

#include <stdint.h>
#include <pthread.h>

#include "config.h"
#include "metrics_store.h"
#include "power_kpi.h"
#include "drv_lm75.h"
#include "drv_lm25066.h"
#include "drv_ads1015.h"

typedef struct {
    pthread_t   tid;
    volatile int stop;

    uint32_t    period_ms;
    MetricsStore *store;
    const Config *config;

    PowerCal    cal;

    Lm75        lm75_board;
    int         lm75_opened;

    Lm25066     lm25066;
    int         lm25066_opened;

    Ads1015     ads1015;
    int         ads1015_opened;

    int         mockMode;
} SampleLoop;

int sample_loop_start(SampleLoop *l, Config *cfg, MetricsStore *store);
void sample_loop_stop(SampleLoop *l);

#endif
