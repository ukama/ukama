/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#include "config.h"
#include "configd.h"
#include "config_macros.h"
#include "service.h"
#include "web.h"
#include "web_client.h"
#include "web_service.h"

#include "usys_api.h"
#include "usys_file.h"
#include "usys_getopt.h"
#include "usys_mem.h"
#include "usys_log.h"
#include "usys_string.h"
#include "usys_types.h"
#include "usys_services.h"

#include "version.h"

/* network.c */
int start_web_services(Config *config, UInst *serviceInst);

void handle_sigint(int signum) {
	usys_log_debug("Caught terminate signal.\n");
	usys_log_debug("Cleanup complete.\n");
	usys_exit(0);
}

static UsysOption longOptions[] = {

    { "port",         required_argument, 0, 'p' },
    { "logs",         required_argument, 0, 'l' },
    { "noded-host",   required_argument, 0, 'n' },
    { "noded-port",   required_argument, 0, 'o' },
    { "noded-ep",     required_argument, 0, 'e' },
    { "starter-host", required_argument, 0, 's' },
    { "starter-port", required_argument, 0, 't' },
    { "help",         no_argument,       0, 'h' },
    { "version",      no_argument,       0, 'v' },
    { 0,              0,                 0,  0 }
};

void set_log_level(char *slevel) {

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

void usage() {

	usys_puts("Usage: configd [options] \n");
	usys_puts("Options:\n");
	usys_puts("-h, --help                          Help\n");
	usys_puts("-l, --logs <TRACE> <DEBUG> <INFO>   Log level for the process\n");
	usys_puts("-n, --noded-host <host>             Host at which node.d listen\n");
	usys_puts("-s, --starter-host <host>           Host at which starter.d listen\n");
	usys_puts("-v, --version                       Version.\n");
}

void free_config(Config *config) {

    usys_free(config->serviceName);
	usys_free(config->nodeId);
	usys_free(config->nodedEP);
	usys_free(config->nodedHost);
	usys_free(config->starterEP);
	usys_free(config->starterHost);
	free_config_data(config->runningConfig);
}

int main(int argc, char **argv) {

	char *debug        = DEF_LOG_LEVEL;
	char *nodedHost    = DEF_NODED_HOST;
	char *starterHost  = DEF_STARTER_HOST;

	UInst serviceInst;

	Config serviceConfig = {0};

	/* Parsing command line args. */
	while (true) {
		int opt = 0;
		int opdIdx = 0;

		opt = getopt_long(argc, argv, "s:l:n:hv", longOptions, &opdIdx);
		if (opt == -1) {
			break;
		}

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

		case 's':
			starterHost = optarg;
			if (!starterHost) {
				usage();
				usys_exit(0);
			}
			break;

		case 'n':
			nodedHost = optarg;
			if (!nodedHost) {
				usage();
				usys_exit(0);
			}
			break;

		default:
			usage();
			usys_exit(0);
		}
	}

	/* Service config update */
	serviceConfig.serviceName  = usys_strdup(SERVICE_CONFIG);
	serviceConfig.servicePort  = usys_find_service_port(SERVICE_CONFIG);
	serviceConfig.nodedEP      = usys_strdup(DEF_NODED_EP);
	serviceConfig.nodedHost    = usys_strdup(nodedHost);
	serviceConfig.nodedPort    = usys_find_service_port(SERVICE_NODE);
	serviceConfig.starterEP    = usys_strdup(DEF_STARTER_EP);
    serviceConfig.starterHost  = usys_strdup(starterHost);
	serviceConfig.starterPort  = usys_find_service_port(SERVICE_STARTER);

    if (!serviceConfig.servicePort ||
        !serviceConfig.nodedPort   ||
        !serviceConfig.starterPort) {

        usys_log_error("Unable to determine port for services");
        free_config(&serviceConfig);
        usys_exit(1);
    }

	usys_log_debug("Starting config.d ...");

	/* Signal handler */
	signal(SIGINT, handle_sigint);

	/* Read Node Info from noded */
	if (getenv(ENV_CONFIG_DEBUG_MODE)) {
		serviceConfig.nodeId = usys_strdup(DEF_NODE_ID);
		usys_log_debug("configd: Using default Node ID: %s", DEF_NODE_ID);
	} else {
		if (get_nodeid_from_noded(&serviceConfig) == STATUS_NOK) {
			usys_log_error("configd: Unable to connect with node.d");
            free_config(&serviceConfig);
            usys_exit(1);
		}
	}

	if (configd_read_running_config((ConfigData**)&serviceConfig.runningConfig)) {
		usys_log_error("Failed to read last running config.");
        free_config(&serviceConfig);
		usys_exit(1);
	}

	if (start_web_services(&serviceConfig, &serviceInst) != USYS_TRUE) {
		usys_log_error("Webservice failed to setup for clients. Exiting.");
        free_config(&serviceConfig);
        usys_exit(1);
	}

	pause();

	ulfius_stop_framework(&serviceInst);
	ulfius_clean_instance(&serviceInst);

    free_config(&serviceConfig);
	usys_log_debug("Exiting config.d ...");

	return 0;
}
