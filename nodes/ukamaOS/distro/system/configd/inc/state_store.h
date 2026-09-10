/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef CONFIG_STATE_STORE_H
#define CONFIG_STATE_STORE_H

#include <pthread.h>
#include <stdint.h>
#include <limits.h>

#define CONFIG_STATE_ID_SIZE   96
#define CONFIG_STATE_NODE_SIZE 128
#define CONFIG_STATE_APPS      32
#define CONFIG_STATE_APP_SIZE  128
#define DEF_CONFIG_STATE_FILE "/ukama/configs/configd/state/configuration.status"

typedef enum {
    CONFIG_MODE_NONE = 0,
    CONFIG_MODE_NOCONFIG,
    CONFIG_MODE_CONFIG
} ConfigMode;

typedef enum {
    CONFIG_PHASE_AWAITING = 0,
    CONFIG_PHASE_PENDING,
    CONFIG_PHASE_COMPLETED,
    CONFIG_PHASE_FAILED
} ConfigPhase;

typedef struct {
    ConfigMode  mode;
    ConfigPhase phase;
    uint64_t    generation;
    char requestId[CONFIG_STATE_ID_SIZE];
    int  revision;
    int  appCount;
    char apps[CONFIG_STATE_APPS][CONFIG_STATE_APP_SIZE];
    char error[128];
} ConfigRecord;

typedef struct {
    pthread_mutex_t mutex;
    int  lockFd;
    char path[PATH_MAX];
    char directory[PATH_MAX];
    char configRoot[PATH_MAX];
    char nodeId[CONFIG_STATE_NODE_SIZE];
    ConfigRecord record;
} ConfigStateStore;

/* Return HTTP-compatible outcomes: 200, 400, 409 or 503. */
int config_store_open(ConfigStateStore *store, const char *path,
                      const char *configRoot, const char *nodeId);
void config_store_close(ConfigStateStore *store);
void config_store_snapshot(ConfigStateStore *store, ConfigRecord *record);
int config_store_noconfig(ConfigStateStore *store, const char *requestId);
int config_store_begin(ConfigStateStore *store, const char *requestId,
                       int revision);
int config_store_track_app(ConfigStateStore *store, const char *app);
int config_store_finish(ConfigStateStore *store, int success);
int config_store_valid_id(const char *value);

#endif
