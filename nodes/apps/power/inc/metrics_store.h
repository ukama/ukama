/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef METRICS_STORE_H
#define METRICS_STORE_H

#include <pthread.h>

#include "json_types.h"

typedef struct {
	pthread_mutex_t	lock;
	PowerMetrics	latest;
} MetricsStore;

int metrics_store_init(MetricsStore *s, const char *boardName);
void metrics_store_free(MetricsStore *s);

void metrics_store_set(MetricsStore *s, const PowerMetrics *m);
void metrics_store_get(MetricsStore *s, PowerMetrics *out);

#endif
