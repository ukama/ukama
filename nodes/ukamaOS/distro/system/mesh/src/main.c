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
#include <unistd.h>
#include <signal.h>
#include <pthread.h>

#include "mesh.h"
#include "config.h"
#include "work.h"
#include "map.h"
#include "websocket.h"

#define VERSION "0.0.1"

/* Global */
State *state=NULL;

/* Defined in network.c */
extern int start_web_services(Config *config, UInst *webInst);
extern int start_websocket_client(Config *config,
								  struct _websocket_client_handler *handler);

/* Global variables. */
WorkList *Transmit=NULL; /* Used by websocket to transmit packet between proxy*/
WorkList *Receive=NULL;
MapTable *ClientTable=NULL;
pthread_mutex_t websocketMutex, mutex;
pthread_cond_t  websocketFail, hasData;

/*
 * usage -- Usage options for the Mesh.d
 *
 */
void usage() {

	printf("Usage: mesh.d [options] \n");
	printf("Options:\n");
	printf("--h, --help                         Help menu.\n");
	printf("--c, --config                       Config file name\n");
	printf("--l, --level <ERROR | DEBUG | INFO> Log level for the process.\n");
	printf("--V, --version                      Version.\n");
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

WorkList **get_transmit(void) {
	return &Transmit;
}

WorkList **get_receive(void) {
	return &Receive;
}

/*
 * close_websocket -- close websocket connection with server with timeout
 *
 */
void close_websocket(struct _websocket_client_handler *handler) {

    int ret;

    ret = ulfius_websocket_client_connection_close(handler);

    switch(ret) {
    case U_WEBSOCKET_STATUS_CLOSE:
        log_debug("Websocket connection with server is closed");
        break;

    case U_WEBSOCKET_STATUS_OPEN:
    case U_WEBSOCKET_STATUS_ERROR:
        log_error("Unable to close websocket connection with server");
        break;
    }
}

/*
 * signal_term_handler -- SIGTERM handling routine. Gracefully exit the process
 *
 */
void signal_term_handler(int signal) {

    log_debug("Received signal: %d (%s)\n", signal, strsignal(signal));

    if (state == NULL) exit(1);

    close_websocket(state->handler);

    if (state->webInst) {
        ulfius_stop_framework(state->webInst);
        ulfius_clean_instance(state->webInst);
    }

    if (state->config) {
        clear_config(state->config);
        free(state->config);
    }

    free(state);
    exit(1);
}

/*
 *  catch_sigterm -- setup SIGTERM catch
 *
 */
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
	struct _websocket_client_handler websocketHandler = {NULL, NULL};

    state = (State *)calloc(1, sizeof(State));
    if (state == NULL) {
        printf("Unable to allocate memory of size: %ld\n", sizeof(State));
        return 1;
    }
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
		log_error("Memory allocation failure: %d", sizeof(Config));
		exit(1);
	}

	/* Step-1: read config file. */
	if (process_config_file(config, configFile) != TRUE) {
		fprintf(stderr, "Error parsing config file: %s. Exiting ... \n",
				configFile);
		exit(1);
	}
    state->config = config;
	print_config(config);

	/* Setup transmit and receiving queues for the websocket */
	Transmit = (WorkList *)malloc(sizeof(WorkList));
	Receive  = (WorkList *)malloc(sizeof(WorkList));

	if (Transmit == NULL && Receive == NULL) {
		log_error("Memory allocation failure: %d", sizeof(WorkList));
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
		log_error("Memory allocation failure: %d", sizeof(MapTable));
		exit(1);
	}
	init_map_table(&ClientTable);

	/* start webservice for local client. */
	if (start_web_services(config, &webInst) != TRUE) {
		log_error("Webservice failed to setup for clients. Exiting.");
		exit(1);
	}

	if (start_websocket_client(config, &websocketHandler) != TRUE) {
		log_error("Websocket failed to setup for client. Retrying soon ...");
	}

    /* create websocket monitoring thread */
    threadArgs.config  = config;
    threadArgs.handler = &websocketHandler;
    if (pthread_create(&thread, NULL, monitor_websocket,
                       (void *)&threadArgs) != 0) {
        log_error("Unable to create websocket monitoring thread.");
    }

    pthread_detach(thread);

	log_debug("Mesh.d running ...");

    pause();

	ulfius_websocket_client_connection_close(&websocketHandler);
	ulfius_stop_framework(&webInst);
	ulfius_clean_instance(&webInst);

	clear_config(config);
	free(config);

	return 1;
}
