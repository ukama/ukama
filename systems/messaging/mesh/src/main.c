/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Mesh.d - L7-websocket based forward/reversed proxy
 *
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <getopt.h>
#include <ulfius.h>
#include <signal.h>
#include <unistd.h>

#include "mesh.h"
#include "config.h"
#include "work.h"
#include "map.h"
#include "u_amqp.h"

#define VERSION "0.0.1"

typedef struct {

	struct _u_instance *websocketInst;
	struct _u_instance *servicesInst;
	Config             *config;
} ProcessState;

/* Defined in network.c */
extern int start_web_services(Config *config, UInst *servicesInst);
extern int start_websocket_server(Config *config, UInst *websocketInst);

/* Global variables. */
MapTable *IDsTable=NULL;
ProcessState *processState=NULL;

/* usage -- Usage options for the Mesh */
void usage(void) {

	printf("Usage: mesh [options] \n");
	printf("Options:\n");
	printf("--h, --help    Help menu.\n");
	printf("--V, --version Version.\n");
    printf("Environment variable needed are: \n");
    printf("\t %s \n\t %s \n\t %s \n\t %s \n\t %s \n\t %s\n\t %s \n\t %s \n",
           ENV_WEBSOCKET_PORT,
           ENV_SERVICES_PORT,
           ENV_AMQP_HOST,
           ENV_AMQP_PORT,
           ENV_INIT_CLIENT_HOST,
           ENV_INIT_CLIENT_PORT,
           ENV_MESH_CERT_FILE,
           ENV_MESH_KEY_FILE);
}

/* Set the verbosity level for logs. */
void set_log_level(char *slevel) {

	int ilevel = LOG_TRACE;

	if (!strcmp(slevel, "DEBUG")) {
		ilevel = LOG_DEBUG;
	} else if (!strcmp(slevel, "INFO")) {
		ilevel = LOG_INFO;
	} else if (!strcmp(slevel, "ERROR")) {
		ilevel = LOG_ERROR;
	}

	log_set_level(ilevel);
}

/* SIGTERM handling routine. Gracefully exit the process
 */
void signal_term_handler(void) {

	if (processState == NULL) exit(1);

	if (processState->websocketInst) {
		ulfius_stop_framework(processState->websocketInst);
		ulfius_clean_instance(processState->websocketInst);
	}

	if (processState->servicesInst) {
		ulfius_stop_framework(processState->servicesInst);
		ulfius_clean_instance(processState->servicesInst);
	}

	if (processState->config) {
		clear_config(processState->config);
		free(processState->config);
	}

    free(processState);

	exit(1);
}

/* setup SIGTERM catch
 */
void catch_sigterm(void) {

	static struct sigaction saction;

    memset(&saction, 0, sizeof(saction));

    saction.sa_sigaction = signal_term_handler;
	sigemptyset(&saction.sa_mask);
    saction.sa_flags     = 0;

    sigaction(SIGTERM, &saction, NULL);
}

int main (int argc, char *argv[]) {

	int exitStatus=0;
	char *debug=DEFAULT_LOG_LEVEL;
	Config *config=NULL;
	struct _u_instance websocketInst;
	struct _u_instance servicesInst;

	processState = (ProcessState *)calloc(1, sizeof(ProcessState));
	if (processState == NULL) return 1;
	processState->websocketInst = &websocketInst;
	processState->servicesInst  = &servicesInst;

	catch_sigterm();

    /* Parse command line args. */
    while (TRUE) {

        int opt = 0;
        int opdidx = 0;

        static struct option long_options[] = {
            { "help",      no_argument,       0, 'h'},
            { "version",   no_argument,       0, 'V'},
            { 0,           0,                 0,  0}
        };

        opt = getopt_long(argc, argv, "hV:", long_options, &opdidx);
        if (opt == -1) {
            break;
        }

        switch (opt) {
        case 'h':
            usage();
            goto exit_program;
            break;

        case 'l':
            debug = optarg;
            set_log_level(debug);
            break;

        default:
            usage();
            goto exit_program;
        }
    } /* while */

	config = (Config *)calloc(1, sizeof(Config));
	if (!config) {
		log_error("Memory allocation failure: %d", sizeof(Config));
		exit(1);
	}

	/* Step-1: read config file. */
    if (!read_config_from_env(&config)) {
        goto exit_program;
    }
	print_config(config);

	IDsTable = (MapTable *)malloc(sizeof(MapTable));
	if (IDsTable == NULL) {
		log_error("Memory allocation failure: %d", sizeof(MapTable));
        exitStatus=1;
        goto exit_program;
	}
	init_map_table(&IDsTable);

	/* Step-2a: setup all endpoints, cb and run websocket. Wait. */
	if (start_websocket(config, &websocketInst) != TRUE) {
		log_error("Websocket failed to setup for server. Exiting...");
		exit(1);
	}

    /* Step-2b: start webservice for the services */
	if (start_web_services(config, &servicesInst) != TRUE) {
		log_error("Webservice failed to setup for clients. Exiting.");
		exit(1);
	}

	log_debug("Mesh running ...");

    pause();

	ulfius_stop_framework(&websocketInst);
	ulfius_stop_framework(&servicesInst);

	ulfius_clean_instance(&websocketInst);
	ulfius_clean_instance(&servicesInst);

	clear_config(config);
	free(config);

exit_program:
	return exitStatus;
}
