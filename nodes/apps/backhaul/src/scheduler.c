/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include <unistd.h>
#include <time.h>

#include "backhauld.h"
#include "scheduler.h"
#include "web_client.h"
#include "usys_log.h"

static long now_ms(void) {

	struct timespec ts;
	clock_gettime(CLOCK_MONOTONIC, &ts);
	return (ts.tv_sec * 1000L) + (ts.tv_nsec / 1000000L);
}

static void ms_sleep(long ms) {
	if (ms <= 0) return;
	usleep((useconds_t)(ms * 1000));
}

static void* scheduler_thread(void *arg) {

	Scheduler *s = (Scheduler *)arg;

	long lastMicro = 0;
	long lastMulti = 0;
	long lastClassify = 0;
	long lastChg = 0;

	long lastRefRefresh = 0;

	while (!s->stop) {

		long now = now_ms();

		/* reflector refresh */
		if (lastRefRefresh == 0 || (time(NULL) - lastRefRefresh) >= s->config->reflectorRefreshSec) {
			ReflectorSet set;
			if (wc_fetch_reflectors(s->config, &set) == STATUS_OK && *set.nearUrl && *set.farUrl) {
				worker_set_reflectors(s->worker, set.nearUrl, set.farUrl, set.ts);
			}
			lastRefRefresh = time(NULL);
		}

		if (lastMicro == 0 || (now - lastMicro) >= s->config->microPeriodMs) {
			worker_enqueue(s->worker, JOB_MICRO_PROBE);
			lastMicro = now;
		}

		if (lastMulti == 0 || (now - lastMulti) >= s->config->multiPeriodMs) {
			worker_enqueue(s->worker, JOB_MULTI_REFLECTOR);
			lastMulti = now;
		}

		if (lastClassify == 0 || (time(NULL) - (lastClassify/1000)) >= s->config->classifyPeriodSec) {
			worker_enqueue(s->worker, JOB_CLASSIFY);
			lastClassify = now;
		}

		if (lastChg == 0 || (time(NULL) - (lastChg/1000)) >= s->config->chgPeriodSec) {
			worker_enqueue(s->worker, JOB_CHG);
			lastChg = now;
		}

		ms_sleep(200); /* keep cheap */
	}

	return NULL;
}

int scheduler_start(Scheduler *s, Config *config, Worker *worker) {

	if (!s || !config || !worker) return USYS_FALSE;

	s->stop = 0;
	s->config = config;
	s->worker = worker;

	if (pthread_create(&s->thread, NULL, scheduler_thread, s) != 0) {
		return USYS_FALSE;
	}

	return USYS_TRUE;
}

void scheduler_stop(Scheduler *s) {

	if (!s) return;

	s->stop = 1;
	if (s->thread) pthread_join(s->thread, NULL);
	s->thread = 0;
}
