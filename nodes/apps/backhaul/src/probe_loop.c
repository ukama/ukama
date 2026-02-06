/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <pthread.h>
#include <unistd.h>
#include <time.h>
#include <string.h>

#include "probe_loop.h"

#include "algo_micro_probe.h"
#include "algo_multi_reflector.h"
#include "algo_chg.h"
#include "algo_diag.h"
#include "classifier.h"

#include "usys_log.h"

typedef struct {
    Config        *config;
    MetricsStore  *store;
    volatile int  *stop;
} ProbeLoopCtx;

static long now_ms(void) {
    struct timespec ts;
    clock_gettime(CLOCK_MONOTONIC, &ts);
    return (ts.tv_sec * 1000L) + (ts.tv_nsec / 1000000L);
}

static void ms_sleep(long ms) {
    if (ms <= 0) return;
    usleep((useconds_t)(ms * 1000));
}

static void run_requested_diag(Config *config, MetricsStore *store) {

    BackhaulDiagRequest req = metrics_store_take_diag_request(store);
    if (req == BACKHAUL_DIAG_NONE) return;

    switch (req) {
    case BACKHAUL_DIAG_CHG:
        metrics_store_set_diag(store, "chg");
        usys_log_info("diag requested: chg");
        (void)algo_chg_run(config, store, NULL);
        break;

    case BACKHAUL_DIAG_PARALLEL:
        metrics_store_set_diag(store, "parallel");
        usys_log_info("diag requested: parallel");
        (void)algo_diag_parallel_run(config, store, NULL);
        break;

    case BACKHAUL_DIAG_BUFFERBLOAT:
        metrics_store_set_diag(store, "bufferbloat");
        usys_log_info("diag requested: bufferbloat");
        (void)algo_diag_bufferbloat_run(config, store, NULL);
        break;

    default:
        break;
    }
}

static void* probe_loop_thread(void *arg) {

    ProbeLoopCtx *ctx = (ProbeLoopCtx *)arg;
    Config *config = ctx->config;
    MetricsStore *store = ctx->store;

    const long microPeriodMs    = (config->microPeriodMs > 0) ?
        config->microPeriodMs : 10000;
    const long multiPeriodMs    = (config->multiPeriodMs > 0) ?
        config->multiPeriodMs : 30000;
    const long classifyPeriodMs = (config->classifyPeriodSec > 0) ?
        (config->classifyPeriodSec * 1000L) : 60000;
    const long chgPeriodMs      = (config->chgPeriodSec > 0) ?
        (config->chgPeriodSec * 1000L) : (1800L * 1000L);

    long lastMicro = 0, lastMulti = 0, lastClassify = 0, lastChg = 0;

    usys_log_info("probe-loop started: micro=%ldms multi=%ldms classify=%ldms chg=%ldms",
                  microPeriodMs, multiPeriodMs, classifyPeriodMs, chgPeriodMs);

    while (!(*ctx->stop)) {

        long now = now_ms();

        run_requested_diag(config, store);

        if (lastMicro == 0 || (now - lastMicro) >= microPeriodMs) {
            (void)algo_micro_probe_run(config, store, NULL);
            lastMicro = now;
        }

        if (lastMulti == 0 || (now - lastMulti) >= multiPeriodMs) {
            (void)algo_multi_reflector_run(config, store, NULL);
            lastMulti = now;
        }

        if (lastClassify == 0 || (now - lastClassify) >= classifyPeriodMs) {
            classifier_run(config, store);
            lastClassify = now;
        }

        if (lastChg == 0 || (now - lastChg) >= chgPeriodMs) {
            (void)algo_chg_run(config, store, NULL);
            lastChg = now;
        }

        ms_sleep(150);
    }

    usys_log_info("probe-loop stopped");
    return NULL;
}

int probe_loop_start(pthread_t *thread,
                     Config *config,
                     MetricsStore *store,
                     volatile int *stopFlag) {

    static ProbeLoopCtx ctx; /* one instance is enough */
    memset(&ctx, 0, sizeof(ctx));
    ctx.config = config;
    ctx.store = store;
    ctx.stop = stopFlag;

    if (!thread || !config || !store || !stopFlag)                  return 0;
    if (pthread_create(thread, NULL, probe_loop_thread, &ctx) != 0) return 0;

    return 1;
}
