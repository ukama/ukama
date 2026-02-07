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

#include "config.h"
#include "metrics_store.h"

typedef struct {

	const Config	*cfg;
	MetricsStore	*store;
	pthread_t		thread;
	int			    stop;
} SampleLoop;

int sample_loop_start(SampleLoop *l, const Config *cfg, MetricsStore *store);
void sample_loop_stop(SampleLoop *l);

#endif
