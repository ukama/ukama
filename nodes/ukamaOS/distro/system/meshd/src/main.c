/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <getopt.h>
#include <ulfius.h>
#include <unistd.h>
#include <signal.h>
#include <pthread.h>

#include "usys_log.h"
#include "usys_services.h"

#include "mesh.h"
#include "config.h"
#include "work.h"
#include "map.h"
#include "websocket.h"

#include "version.h"

/* Global */
State *state=NULL;

/* Global variables. */
WorkList *Transmit=NULL; /* Used by websocket to transmit packet between proxy*/
WorkList *Receive=NULL;
MapTable *ClientTable=NULL;
pthread_mutex_t websocketMutex, mutex;
pthread_cond_t  websocketFail, hasData;

void usage() {

	printf("Usage: mesh.d [options] \n");
	printf("Options:\n");
	printf("--h, --help                         Help menu.\n");
	printf("--c, --config                       Config file name\n");
	printf("--l, --level <ERROR | DEBUG | INFO> Log level for the process.\n");
	printf("--V, --version                      Version.\n");
}

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

WorkList **get_transmit(void) {
	return &Transmit;
}

WorkList **get_receive(void) {
	return &Receive;
}

void close_websocket(struct _websocket_client_handler *handler) {

    int ret;

    ret = ulfius_websocket_client_connection_close(handler);

    switch(ret) {
    case U_WEBSOCKET_STATUS_CLOSE:
        usys_log_debug("Websocket connection with server is closed");
        break;

    case U_WEBSOCKET_STATUS_OPEN:
    case U_WEBSOCKET_STATUS_ERROR:
        usys_log_error("Unable to close websocket connection with server");
        break;
    }
}

void signal_term_handler(int signal) {

    usys_log_debug("Received signal: %d (%s)\n", signal, strsignal(signal));

    if (state == NULL) exit(1);

    close_websocket(state->handler);

    if (state->webInst && state->webInst->port) {
        ulfius_stop_framework(state->webInst);
        ulfius_clean_instance(state->webInst);
    }

    if (state->fwdInst && state->fwdInst->port) {
        ulfius_stop_framework(state->fwdInst);
        ulfius_stop_framework(state->fwdInst);
    }

    if (state->config) {
        clear_config(state->config);
        free(state->config);
    }

    free(state);
    exit(1);
}

void catch_sigterm(void) {

    static struct sigaction saction;

    memset(&saction, 0, sizeof(saction));

    saction.sa_sigaction = signal_term_handler;
    sigemptyset(&saction.sa_mask);
    saction.sa_flags     = 0;

    sigaction(SIGINT, &saction, NULL);
    sigaction(SIGTERM, &saction, NULL);
}

int main (int argc, char *argv[]) {

	char *configFile=NULL;
	char *debug=DEF_LOG_LEVEL;
	Config *config=NULL;
    ThreadArgs threadArgs;
    pthread_t thread;
    
	struct _u_instance webInst;
    struct _u_instance fwdInst;
	struct _websocket_client_handler websocketHandler = {NULL, NULL};

    log_set_service(SERVICE_NAME);

    state = (State *)calloc(1, sizeof(State));
    if (state == NULL) {
        printf("Unable to allocate memory of size: %ld\n", sizeof(State));
        return 1;
    }
    state->fwdInst = &fwdInst;
    state->webInst = &webInst;
    state->handler = &websocketHandler;

    catch_sigterm();

	/* Prase command line args. */
	while (TRUE) {

		int opt = 0;
		int opdidx = 0;

		static struct option long_options[] = {
			{ "config",    required_argument, 0, 'c'},
			{ "level",     required_argument, 0, 'l'},
			{ "help",      no_argument,       0, 'h'},
			{ "version",   no_argument,       0, 'V'},
			{ 0,           0,                 0,  0}
		};

		opt = getopt_long(argc, argv, "l:c:sphV:", long_options, &opdidx);
		if (opt == -1) {
			break;
		}

		switch (opt) {
		case 'h':
			usage();
			exit(0);
			break;

		case 'c':
			configFile = optarg;
			break;

		case 'l':
			debug = optarg;
			set_log_level(debug);
			break;

		case 'V':
			fprintf(stdout, "Mesh.d - Version: %s\n", VERSION);
			exit(0);

		default:
			usage();
			exit(0);
		}
	} /* while */

	if (argc == 1 || configFile == NULL) {
		fprintf(stderr, "Missing required parameters\n");
		usage();
		exit(1);
	}

	config = (Config *)calloc(1, sizeof(Config));
	if (!config) {
		usys_log_error("Memory allocation failure: %d", sizeof(Config));
		exit(1);
	}

	/* Step-1: read config file. */
	if (process_config_file(config, configFile) != TRUE) {
		usys_log_error("Error parsing config file: %s. Exiting.",
                       configFile);
		exit(1);
	}

    config->forwardPort = usys_find_service_port(SERVICE_UKAMA);
    config->servicePort = usys_find_service_port(SERVICE_NAME);
    if (!config->forwardPort) {
        usys_log_error("Unable to find forward port in /etc/service");
        exit(1);
    }

    if (!config->servicePort) {
        usys_log_error("Unable to find %s port in /etc/service", SERVICE_NAME);
        exit(1);
    }
    state->config = config;
	print_config(config);

	/* Setup transmit and receiving queues for the websocket */
	Transmit = (WorkList *)malloc(sizeof(WorkList));
	Receive  = (WorkList *)malloc(sizeof(WorkList));

	if (Transmit == NULL && Receive == NULL) {
		usys_log_error("Memory allocation failure: %d", sizeof(WorkList));
		exit(1);
	}

	/* Initializa the transmit and receive list for the websocket. */
	init_work_list(&Transmit);
	init_work_list(&Receive);

    pthread_mutex_init(&websocketMutex, NULL);
    pthread_mutex_init(&mutex, NULL);
    pthread_cond_init(&websocketFail, NULL);
    pthread_cond_init(&hasData, NULL);

	/* Setup ip:port to UUID mapping table, if client. */
	ClientTable = (MapTable *)malloc(sizeof(MapTable));
	if (ClientTable == NULL) {
		usys_log_error("Memory allocation failure: %d", sizeof(MapTable));
		exit(1);
	}
	init_map_table(&ClientTable);

    if (start_web_services(config, &webInst) != TRUE) {
        usys_log_error("Web service failed to setup. Exiting.");
		exit(1);
	}

    while (start_websocket_client(config, &websocketHandler) != TRUE) {
		usys_log_error("Websocket failed to setup for client. Retrying in 5 seconds ...");
        sleep(5);
	}

	if (start_forward_services(config, &fwdInst) != TRUE) {
		usys_log_error("Forward service failed to setup. Exiting.");
		exit(1);
	}

    /* create websocket monitoring thread */
    threadArgs.config  = config;
    threadArgs.handler = &websocketHandler;
    if (pthread_create(&thread, NULL, monitor_websocket,
                       (void *)&threadArgs) != 0) {
        usys_log_error("Unable to create websocket monitoring thread.");
    }

    pthread_detach(thread);

	usys_log_debug("%s running ...", SERVICE_NAME);

    pause();

	ulfius_websocket_client_connection_close(&websocketHandler);
	ulfius_stop_framework(&fwdInst);
    ulfius_stop_framework(&webInst);
	ulfius_clean_instance(&fwdInst);
	ulfius_clean_instance(&webInst);

	clear_config(config);
	free(config);

	return 1;
}
