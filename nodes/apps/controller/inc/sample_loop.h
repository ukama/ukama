/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef SAMPLE_LOOP_H
#define SAMPLE_LOOP_H

#include <pthread.h>
#include <stdbool.h>

#include "config.h"
#include "driver.h"
#include "metrics_store.h"

/*
 * Sample loop context
 */
typedef struct {
    pthread_t               thread;
    volatile bool           running;

    const Config            *config;
    const ControllerDriver  *driver;
    void                    *driver_ctx;
    MetricsStore            *store;

    /* Statistics */
    uint64_t                samples_taken;
    uint64_t                samples_failed;
} SampleLoop;

/* Lifecycle */
int  sample_loop_start(SampleLoop *loop, const Config *config,
                       const ControllerDriver *driver, void *driver_ctx,
                       MetricsStore *store);
void sample_loop_stop(SampleLoop *loop);

/* Get sample statistics */
void sample_loop_get_stats(SampleLoop *loop, uint64_t *taken, uint64_t *failed);

#endif /* SAMPLE_LOOP_H */
