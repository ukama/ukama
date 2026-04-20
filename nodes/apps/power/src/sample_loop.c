/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>
#include <unistd.h>
#include <time.h>

#include "sample_loop.h"
#include "power_collector.h"
#include "usys_log.h"

static uint64_t now_ms(void) {

    struct timespec ts;

    clock_gettime(CLOCK_REALTIME, &ts);

    return ((uint64_t)ts.tv_sec * 1000ULL) +
           ((uint64_t)ts.tv_nsec / 1000000ULL);
}

static void *loop_thread(void *arg) {

    SampleLoop *l = (SampleLoop *)arg;
    PowerCollectorCtx pc;

    if (!l || !l->store || !l->config) return NULL;

    memset(&pc, 0, sizeof(pc));
    pc.store      = l->store;
    pc.config     = l->config;
    pc.cal        = &l->cal;
    pc.lm75_board = l->lm75_opened ? &l->lm75_board : NULL;
    pc.lm25066    = l->lm25066_opened ? &l->lm25066 : NULL;
    pc.ads1015    = l->ads1015_opened ? &l->ads1015 : NULL;
    pc.mockMode   = l->mockMode;

    if (l->period_ms == 0) l->period_ms = 1000;

    usys_log_info("sample_loop: start period=%u mock=%d",
                  l->period_ms, l->mockMode);

    while (!l->stop) {
        (void)power_collect_once(&pc, now_ms());
        usleep((useconds_t)l->period_ms * 1000);
    }

    usys_log_info("sample_loop: stop");
    return NULL;
}

int sample_loop_start(SampleLoop *l, Config *cfg, MetricsStore *store) {

    if (!l || !cfg || !store) return USYS_FALSE;

    memset(l, 0, sizeof(*l));

    l->store = store;
    l->config = cfg;
    l->period_ms = cfg->sampleMs;
    l->mockMode = cfg->mockMode;

    memset(&l->cal, 0, sizeof(l->cal));

    if (!l->mockMode) {
        if (cfg->lm75Dev && cfg->lm75Addr > 0) {
            if (drv_lm75_open(&l->lm75_board, cfg->lm75Dev, cfg->lm75Addr) == 0) {
                l->lm75_opened = 1;
                usys_log_info("sample_loop: lm75 opened dev=%s addr=0x%02x",
                              cfg->lm75Dev, cfg->lm75Addr);
            } else {
                usys_log_error("sample_loop: lm75 open failed dev=%s addr=0x%02x",
                               cfg->lm75Dev, cfg->lm75Addr);
            }
        }

        if (cfg->lm25066Dev && cfg->lm25066Addr > 0) {
            if (drv_lm25066_open(&l->lm25066,
                                 cfg->lm25066Dev,
                                 cfg->lm25066Addr,
                                 cfg->lm25066ClHigh,
                                 cfg->lm25066RsMohm) == 0) {
                l->lm25066_opened = 1;
                usys_log_info("sample_loop: lm25066 opened dev=%s addr=0x%02x",
                              cfg->lm25066Dev, cfg->lm25066Addr);
            } else {
                usys_log_error("sample_loop: lm25066 open failed dev=%s addr=0x%02x",
                               cfg->lm25066Dev, cfg->lm25066Addr);
            }
        }

        if (cfg->ads1015Dev && cfg->ads1015Addr > 0) {
            if (drv_ads1015_open(&l->ads1015,
                                 cfg->ads1015Dev,
                                 cfg->ads1015Addr) == 0) {
                l->ads1015_opened = 1;
                usys_log_info("sample_loop: ads1015 opened dev=%s addr=0x%02x",
                              cfg->ads1015Dev, cfg->ads1015Addr);
            } else {
                usys_log_error("sample_loop: ads1015 open failed dev=%s addr=0x%02x",
                               cfg->ads1015Dev, cfg->ads1015Addr);
            }
        }
    } else {
        usys_log_info("sample_loop: mock mode enabled");
    }

    l->stop = 0;

    if (pthread_create(&l->tid, NULL, loop_thread, l) != 0) {
        usys_log_error("sample_loop: pthread_create failed");
        sample_loop_stop(l);
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

void sample_loop_stop(SampleLoop *l) {

    if (!l) return;

    l->stop = 1;

    if (l->tid) {
        pthread_join(l->tid, NULL);
        l->tid = 0;
    }

    if (l->lm75_opened) {
        drv_lm75_close(&l->lm75_board);
        l->lm75_opened = 0;
    }

    if (l->lm25066_opened) {
        drv_lm25066_close(&l->lm25066);
        l->lm25066_opened = 0;
    }

    if (l->ads1015_opened) {
        drv_ads1015_close(&l->ads1015);
        l->ads1015_opened = 0;
    }
}
