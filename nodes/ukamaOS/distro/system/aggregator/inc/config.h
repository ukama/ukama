/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef AGG_CONFIG_H_
#define AGG_CONFIG_H_

#define TAG_REFRESH_INTERVAL_SEC "refresh_interval_sec"
#define TAG_REQUEST_TIMEOUT_MS   "request_timeout_ms"
#define TAG_STALE_GRACE_SEC      "stale_grace_sec"
#define TAG_SOURCE               "source"
#define TAG_NAME                 "name"
#define TAG_URL                  "url"
#define TAG_REQUIRED             "required"

typedef struct SourceConfig {
    char *name;
    char *url;
    int required;

    struct SourceConfig *next;
} SourceConfig;

typedef struct {
    int refreshIntervalSec;
    int requestTimeoutMs;
    int staleGraceSec;
    int sourceCount;

    SourceConfig *sources;
} Config;

int config_load(const char *fileName, Config *config);
void config_free(Config *config);

#endif /* AGG_CONFIG_H_ */
