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
#include "usys_getopt.h"
#include "usys_api.h"

#include "version.h"

static volatile int gStop = 0;

static void handle_term(int signum) {
    (void)signum;
    usys_log_info("Terminate signal.");
    gStop = 1;
}

static UsysOption longOptions[] = {
    { "logs",    required_argument, 0, 'l' },
    { "help",    no_argument,       0, 'h' },
    { "version", no_argument,       0, 'v' },
    { 0, 0, 0, 0 }
};

static void set_log_level(char *slevel) {

    int ilevel = USYS_LOG_TRACE;

    if (!strcmp(slevel, "TRACE"))      ilevel = USYS_LOG_TRACE;
    else if (!strcmp(slevel, "DEBUG")) ilevel = USYS_LOG_DEBUG;
    else if (!strcmp(slevel, "INFO"))  ilevel = USYS_LOG_INFO;

    usys_log_set_level(ilevel);
}

static void usage(void) {

    usys_puts("Usage: power-monitor [options]");
    usys_puts("Options:");
    usys_puts("-h, --help                    Help menu");
    usys_puts("-l, --logs <TRACE|DEBUG|INFO> Log level for the process");
    usys_puts("-v, --version                 Software version");
}

int main(int argc, char **argv) {

    int opt, optIdx;
    char *debug = DEF_LOG_LEVEL;

    Config config = {0};
    SampleLoop sampler;
    MetricsStore store = {0};
    struct _u_instance inst;
    EpCtx ctx = {0};

    pthread_t probeThread = 0;

    usys_log_set_service(SERVICE_NAME);
    usys_log_info("starting %s", SERVICE_NAME);

    signal(SIGINT,  handle_term);
    signal(SIGTERM, handle_term);

    while (1) {
        opt = 0; optIdx = 0;
        opt = usys_getopt_long(argc, argv, "vh:l:", longOptions, &optIdx);
        if (opt == -1) break;

        switch (opt) {
        case 'h':
            usage();
            config_print_env_help();
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

	if (config_load_from_env(&config) != 0) {
		usys_log_error("Failed to load required config");
        usys_exit(1);
	}

	if (metrics_store_init(&store, config.boardName) != 0) {
		usys_log_error("Failed to initialize metrics store");
		config_free(&config);
        usys_exit(1);
	}

    ctx.config = &config;
	ctx.store  = &store;

    if (!start_web_service(&config, &inst, &ctx)) {
		usys_log_error("Failed to start the web services");
		metrics_store_free(&store);
		config_free(&config);
        usys_exit(1);
	}

	if (sample_loop_start(&sampler, &config, &store) != 0) {
		usys_log_error("failed to start the sampling loop");
        web_service_stop(&inst);
		metrics_store_free(&store);
		config_free(&config);
		usys_exit(1);
	}

    pause();

	usys_log_info("stopping %s", SERVICE_NAME);

	web_service_stop(&inst);
	sample_loop_stop(&sampler);

	metrics_store_free(&store);
	config_free(&config);

	return 0;
}
