/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#pragma once

#include <pthread.h>
#include <signal.h>
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

#include "config.h"
#include "fsm.h"

#define LIFECYCLED_SERVICE_NAME "lifecycle"
#define LIFECYCLED_REASON_LEN   192
#define LIFECYCLED_ID_LEN       96
#define LIFECYCLED_EVENT_QUEUE  32

typedef struct {
    LifecycleState state;
    uint64_t sequence;
    int64_t occurredAt;
    char reason[LIFECYCLED_REASON_LEN];
} LifecycleEvent;

typedef struct {
    Config *config;
    LifecycleFsm fsm;
    StarterSnapshot starter;

    pthread_mutex_t mutex;
    pthread_cond_t condition;
    pthread_t worker;
    bool workerStarted;
    bool running;

    LifecycleEvent events[LIFECYCLED_EVENT_QUEUE];
    size_t eventHead;
    size_t eventCount;

    char bootId[LIFECYCLED_ID_LEN];
    struct _u_instance *uInstance;
} LifecycleContext;

int64_t lifecycle_boottime_ms(void);
int64_t lifecycle_epoch_sec(void);
bool lifecycle_read_boot_id(char *buffer, size_t size);

bool lifecycle_context_init(LifecycleContext *ctx, Config *config);
void lifecycle_context_free(LifecycleContext *ctx);
bool lifecycle_context_start(LifecycleContext *ctx);
void lifecycle_context_stop(LifecycleContext *ctx);

bool lifecycle_context_check_in(LifecycleContext *ctx,
                                const char *bootId,
                                bool bootHealthy,
                                char *error,
                                size_t errorSize);

LifecycleConfigureResult lifecycle_context_configure(
    LifecycleContext *ctx,
    const char *requestId,
    const char *assignmentId,
    char *error,
    size_t errorSize);

void lifecycle_context_snapshot(LifecycleContext *ctx,
                                LifecycleFsm *fsm,
                                StarterSnapshot *starter);

