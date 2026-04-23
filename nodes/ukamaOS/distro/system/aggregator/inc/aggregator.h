/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef AGGREGATOR_H_
#define AGGREGATOR_H_

#include <pthread.h>
#include <stdbool.h>
#include <stdint.h>
#include <time.h>

#include "config.h"

#include "usys_services.h"

#define RETURN_OK        0
#define RETURN_NOTOK    -1

#define DEFAULT_CONFIG_PATH "./config/config.toml"
#define DEFAULT_LOG_LEVEL   "DEBUG"

#define SERVICE_NAME       SERVICE_AGGREGATOR
#define SERVICE_NAME_ADMIN SERVICE_AGGREGATOR_ADMIN

#define AGG_MAX_NAME_LEN   64

typedef struct {
    char name[AGG_MAX_NAME_LEN];
    char *url;
    int required;
    int up;
    int lastHttpCode;
    int errorCount;
    time_t lastAttempt;
    time_t lastSuccess;
    char *body;
} SourceState;

typedef struct {
    pthread_mutex_t mutex;
    SourceState *sources;
    int sourceCount;
    int refreshIntervalSec;
    int requestTimeoutMs;
    int staleGraceSec;
    int running;
    pthread_t thread;
    char *snapshot;
    time_t snapshotAt;
} AppState;

int app_state_init(AppState *state, const Config *config);
void app_state_cleanup(AppState *state);
int app_state_start(AppState *state);
void app_state_stop(AppState *state);
char *app_state_dup_snapshot(AppState *state);
char *app_state_status_json(AppState *state);


#endif /* AGGREGATOR_H_ */
