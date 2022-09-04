/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * msgClient - GRPC to AMQP service
 *
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <getopt.h>
#include <ulfius.h>
#include <signal.h>

#include "msgClient.h"
#include "config.h"
#include "log.h"

#define VERSION "0.0.1"

typedef struct {

	struct _u_instance *webInst;
	Config             *config;
} State;

extern int start_web_services(Config *config, UInst *webtInst);

/* Global */
State *state=NULL; 

/*
 * usage -- Usage options for msgClient
 *
 */
void usage() {

	fprintf(stdout, "Usage: msgClient [options] \n");
	fprintf(stdout, "Options:\n");
	fprintf(stdout, "--h, --help     this menu\n");
	fprintf(stdout, "--V, --version  Version\n");
	fprintf(stdout, "Environment variables are:\n");
	fprintf(stdout, "\t MSG_CLIENT_LOG_LEVEL\n");
	fprintf(stdout, "\t MSG_CLIENT_IP\n");
	fprintf(stdout, "\t MSG_CLIENT_PORT\n");
}

/*
 * set_log_level --  set the verbosity level for logs
 *
 */
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

/* 
 * signal_term_handler -- SIGTERM handling routine. Gracefully exit the process
 *
 */
void signal_term_handler(void) {

	if (state == NULL) exit(1);

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

    sigaction(SIGTERM, &saction, NULL);
}


/*
 * Look for environment variables.
 * signal handling and graceful exit (AMQP discon, mem clear, web close)
 * setup client webinstance for /ping
 * run GRPC server
 * listen for messaging system and establie to AMQP
 * forward all events to AMQP
 *
 */

int main (int argc, char *argv[]) {

	int exitStatus=0;
	char *debug=DEF_LOG_LEVEL;
	char address[MAX_BUFFER_SIZE] = {0};
	struct _u_instance webInst;
	Config *config=NULL;
	
	state = (State *)calloc(1, sizeof(State));
	if (state == NULL) {
		printf("Unable to allocate memory of size: %ld\n", sizeof(State));
		return 1;
	}
	state->webInst = &webInst;

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

	/* Step-1: read config params */
	if (!read_config_from_env(&config)) {
		goto exit_program;
	}
	state->config = config;

	/* Step-2: start webservice */
	if (start_web_services(config, &webInst) != TRUE) {
		log_error("Webservice failed to setup for clients. Exiting.");
		exitStatus = 1;
		goto exit_program;
	}

	/* Step-3: start GRPC server */
	sprintf(address, "%s:%s", config->ip, config->port);
	run_grpc_server(address);
	
	/* Wait here for ever. XXX */

	log_debug("Mesh.d running ...");

	getchar(); /* For now. */

	log_debug("Goodbye ... ");

	ulfius_stop_framework(&webInst);
	ulfius_clean_instance(&webInst);

	clear_config(config);

 exit_program:
	free(state);

	return exitStatus;
}
