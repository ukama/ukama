/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include <string.h>
#include <unistd.h>
#include <time.h>
#include <stdlib.h>

#include "sample_loop.h"
#include "power_collector.h"

#include "usys_log.h"

#ifndef POWER_I2C_DEV_ENV
#define POWER_I2C_DEV_ENV   "POWER_I2C_DEV"
#endif

#ifndef POWER_LM75_ADDR_ENV
#define POWER_LM75_ADDR_ENV "POWER_LM75_ADDR"
#endif

static uint64_t now_ms(void) {

	struct timespec ts;
	clock_gettime(CLOCK_REALTIME, &ts);
	return ((uint64_t)ts.tv_sec * 1000ULL) +
	       ((uint64_t)ts.tv_nsec / 1000000ULL);
}

static int env_get_i2c_dev(char *out, size_t outlen) {

	const char *v = getenv(POWER_I2C_DEV_ENV);
	if (!v || !*v) v = "/dev/i2c-1";
	snprintf(out, outlen, "%s", v);
	return 0;
}

static int env_get_lm75_addr7(void) {

	const char *v = getenv(POWER_LM75_ADDR_ENV);
	long a;

	if (!v || !*v) return 0x48;

	/* accept "0x48" or "72" */
	a = strtol(v, NULL, 0);
	if (a < 0x03 || a > 0x77) return 0x48;
	return (int)a;
}

static void *loop_thread(void *arg) {

	SampleLoop *l = (SampleLoop *)arg;
	PowerCollectorCtx pc;

	if (!l || !l->store) return NULL;

	memset(&pc, 0, sizeof(pc));
	pc.store = l->store;
	pc.cal = NULL;
	pc.lm75_board = l->lm75_opened ? &l->lm75_board : NULL;

	if (l->period_ms == 0) l->period_ms = 1000;

	usys_log_info("sample_loop: start period=%ums", l->period_ms);

	while (!l->stop) {
		(void)power_collect_once(&pc, now_ms());
		usleep((useconds_t)l->period_ms * 1000);
	}

	usys_log_info("sample_loop: stop");
	return NULL;
}

int sample_loop_start(SampleLoop *l, Config *cfg, MetricsStore *store) {

	char dev[64];
	int addr7;

	if (!l || !cfg || !store) return USYS_FALSE;

	memset(l, 0, sizeof(*l));
	l->store = store;

	/*
	 * Sampling period: use config if you already have it, otherwise default.
	 * If your Config has a field (e.g., cfg->samplePeriodMs), wire it here.
	 */
	l->period_ms = 1000;

	/*
	 * Bring up LM75. If it fails, we still start the loop but it will report
	 * "lm75 not configured" in JSON.
	 */
	(void)env_get_i2c_dev(dev, sizeof(dev));
	addr7 = env_get_lm75_addr7();

	if (drv_lm75_open(&l->lm75_board, dev, addr7) == 0) {
		l->lm75_opened = 1;
		usys_log_info("sample_loop: lm75 opened dev=%s addr=0x%02x", dev, addr7);
	} else {
		l->lm75_opened = 0;
		usys_log_error("sample_loop: lm75 open failed dev=%s addr=0x%02x", dev, addr7);
	}

	l->stop = 0;

	if (pthread_create(&l->tid, NULL, loop_thread, l) != 0) {
		usys_log_error("sample_loop: pthread_create failed");
		if (l->lm75_opened) drv_lm75_close(&l->lm75_board);
		memset(l, 0, sizeof(*l));
		return USYS_FALSE;
	}

	return USYS_TRUE;
}

void sample_loop_stop(SampleLoop *l) {

	if (!l) return;

	l->stop = 1;

	if (l->tid) {
		pthread_join(l->tid, NULL);
		l->tid = 0;
	}

	if (l->lm75_opened) {
		drv_lm75_close(&l->lm75_board);
		l->lm75_opened = 0;
	}
}
