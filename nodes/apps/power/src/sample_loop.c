/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <unistd.h>
#include <stdint.h>
#include <time.h>

#include "sample_loop.h"
#include "usys_log.h"

#include "power_collector.h"

static uint64_t now_ms(void) {

	struct timespec ts;
	clock_gettime(CLOCK_REALTIME, &ts);
	return ((uint64_t)ts.tv_sec * 1000ULL) + ((uint64_t)ts.tv_nsec / 1000000ULL);
}

static int env_get_int_default(const char *k, int defv) {

	char *v = getenv(k);
	if (!v || !*v) return defv;
	return atoi(v);
}

void *sample_loop(void *arg) {

	SampleCtx *ctx = (SampleCtx *)arg;
	PowerCollectorCtx pc;
	int sample_ms;

	if (!ctx || !ctx->store) return NULL;

	sample_ms = env_get_int_default("POWER_SAMPLE_MS", 1000);

	pc.store = ctx->store;
	pc.cal = ctx->cal;

	usys_log_info("sample_loop: start sample_ms=%d", sample_ms);

	while (!ctx->stop) {
		uint64_t t = now_ms();
		power_collect_once(&pc, t);
		usleep((useconds_t)sample_ms * 1000);
	}

	usys_log_info("sample_loop: stop");
	return NULL;
}
