/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#ifndef __SAMPLE_LOOP_H__
#define __SAMPLE_LOOP_H__

#include <stdint.h>
#include <pthread.h>

#include "config.h"
#include "metrics_store.h"
#include "drv_lm75.h"

/*
 * SampleLoop is what main.c expects.
 * It owns the thread + stop flag and the sensor handles used by the loop.
 */
typedef struct {
	pthread_t	tid;
	volatile int	stop;

	uint32_t	period_ms;

	MetricsStore	*store;

	/* sensors */
	Lm75		lm75_board;
	int		lm75_opened;
} SampleLoop;

int sample_loop_start(SampleLoop *l, Config *cfg, MetricsStore *store);
void sample_loop_stop(SampleLoop *l);

#endif /* __SAMPLE_LOOP_H__ */
