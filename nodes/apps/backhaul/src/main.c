/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include <signal.h>
#include <unistd.h>

#include "backhauld.h"
#include "config.h"
#include "worker.h"
#include "scheduler.h"
#include "metrics_store.h"
#include "web_service.h"
#include "web_client.h"

#include "usys_log.h"
#include "usys_getopt.h"
#include "usys_string.h"
#include "usys_api.h"

#include "version.h"

static void handle_sigint(int signum) {
	usys_log_debug("Terminate signal.\n");
	usys_exit(0);
}

static UsysOption longOptions[] = {
	{ "logs",	required_argument, 0, 'l' },
	{ "help",	no_argument,	   0, 'h' },
	{ "version",	no_argument,	   0, 'v' },
	{ 0, 0, 0, 0 }
};

static void set_log_level(char *slevel) {

	int ilevel = USYS_LOG_TRACE;

	if (!strcmp(slevel, "TRACE")) {
		ilevel = USYS_LOG_TRACE;
	} else if (!strcmp(slevel, "DEBUG")) {
		ilevel = USYS_LOG_DEBUG;
	} else if (!strcmp(slevel, "INFO")) {
		ilevel = USYS_LOG_INFO;
	}

	usys_log_set_level(ilevel);
}

static void usage() {

	usys_puts("Usage: backhaul.d [options]");
	usys_puts("Options:");
	usys_puts("-h, --help                    Help menu");
	usys_puts("-l, --logs <TRACE|DEBUG|INFO> Log level for the process");
	usys_puts("-v, --version                 Software version");
}

int main(int argc, char **argv) {

	int opt, optIdx;
	char *debug = DEF_LOG_LEVEL;

	Config config = {0};
	MetricsStore store = {0};
	Worker worker = {0};
	Scheduler sched = {0};

	UInst serviceInst;
	EpCtx epCtx = {0};

	usys_log_set_service(SERVICE_NAME);
	usys_log_remote_init(SERVICE_NAME);

	while (1) {

		opt = 0; optIdx = 0;
		opt = usys_getopt_long(argc, argv, "vh:l:", longOptions, &optIdx);
		if (opt == -1) break;

		switch (opt) {
		case 'h':
			usage();
			usys_exit(0);
			break;
		case 'v':
			usys_puts(VERSION);
			usys_exit(0);
			break;
		case 'l':
			debug = optarg;
			set_log_level(debug);
			break;
		default:
			usage();
			usys_exit(0);
		}
	}

	if (!config_load_from_env(&config)) {
		usys_log_error("Failed to load config");
		usys_exit(1);
	}

	if (!metrics_store_init(&store, config.windowMicroSamples, config.windowMultiSamples, config.windowChgSamples)) {
		usys_log_error("Failed to init metrics store");
		usys_exit(1);
	}

	if (!wc_init()) {
		usys_log_error("Failed to init web client");
		usys_exit(1);
	}

	if (!worker_init(&worker, &config, &store, 64)) {
		usys_log_error("Failed to init worker");
		usys_exit(1);
	}

	if (!worker_start(&worker)) {
		usys_log_error("Failed to start worker");
		usys_exit(1);
	}

	if (!scheduler_start(&sched, &config, &worker)) {
		usys_log_error("Failed to start scheduler");
		usys_exit(1);
	}

	epCtx.config = &config;
	epCtx.worker = &worker;
	epCtx.store = &store;

	signal(SIGINT, handle_sigint);

	if (!start_web_service(&config, &serviceInst, &epCtx)) {
		usys_log_error("Web service failed to start");
		usys_exit(1);
	}

	pause();

	/* cleanup */
	scheduler_stop(&sched);
	worker_stop(&worker);

	ulfius_stop_framework(&serviceInst);
	ulfius_clean_instance(&serviceInst);

	worker_free(&worker);
	metrics_store_free(&store);
	wc_cleanup();
	config_free(&config);

	return USYS_TRUE;
}
