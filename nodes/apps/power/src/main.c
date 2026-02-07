/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <signal.h>
#include <stdlib.h>

#include "powerd.h"
#include "config.h"
#include "metrics_store.h"
#include "sample_loop.h"
#include "web_service.h"

#include "usys_log.h"

static volatile int gStop = 0;

static void on_signal(int sig) {

	(void)sig;
	gStop = 1;
}

int main(void) {

	Config cfg;
	MetricsStore store;
	SampleLoop sampler;
	struct _u_instance inst;
	EpCtx ctx;

	signal(SIGINT, on_signal);
	signal(SIGTERM, on_signal);

	usys_log_init(POWERD_NAME);
	usys_log_info("starting %s", POWERD_NAME);

	if (config_load(&cfg) != 0) {
		usys_log_error("config_load failed");
		return 1;
	}

	if (metrics_store_init(&store, cfg.boardName) != 0) {
		usys_log_error("metrics_store_init failed");
		config_free(&cfg);
		return 1;
	}

	if (sample_loop_start(&sampler, &cfg, &store) != 0) {
		usys_log_error("sample_loop_start failed");
		metrics_store_free(&store);
		config_free(&cfg);
		return 1;
	}

	ctx.cfg = &cfg;
	ctx.store = &store;

	if (web_service_start(&inst, &ctx) != 0) {
		usys_log_error("web_service_start failed");
		sample_loop_stop(&sampler);
		metrics_store_free(&store);
		config_free(&cfg);
		return 1;
	}

	while (!gStop) {
		usleep(200 * 1000);
	}

	usys_log_info("stopping %s", POWERD_NAME);

	web_service_stop(&inst);
	sample_loop_stop(&sampler);

	metrics_store_free(&store);
	config_free(&cfg);

	return 0;
}
