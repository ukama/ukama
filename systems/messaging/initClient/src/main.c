/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * initClient - client to register to init system.
 *
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <getopt.h>
#include <ulfius.h>
#include <signal.h>

#include "initClient.h"
#include "config.h"
#include "log.h"

#define VERSION "0.0.1"

typedef struct {

	struct _u_instance *webInst;
	Config             *config;
} State;

extern int start_web_services(Config *config, UInst *webtInst); /*network.c */
extern int send_request_to_init(ReqType reqType, Config *config); /* init.c */

/* Global */
State *state=NULL; 

/*
 * usage -- Usage options for initClient
 *
 */
void usage() {

	fprintf(stdout, "Usage: initClient [options] \n");
	fprintf(stdout, "Options:\n");
	fprintf(stdout, "--h, --help     this menu\n");
	fprintf(stdout, "--V, --version  Version\n");
	fprintf(stdout, "Environment variable used are: \n");
	fprintf(stdout, "\t %s \n\t %s \n\t %s \n\t %s \n\t %s \n\t %s\n\t %s \n",
			ENV_INIT_CLIENT_LOG_LEVEL,
			ENV_SYSTEM_ORG,
			ENV_SYSTEM_NAME,
			ENV_SYSTEM_ADDR,
			ENV_SYSTEM_PORT,
			ENV_INIT_SYSTEM_ADDR,
			ENV_INIT_SYSTEM_PORT);
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

	/* un-register the system */
	if (send_request_to_init(REQ_UNREGISTER, state->config) != TRUE) {
		log_error("Error registrating with the init system");
	}

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
 * Life of initClient:
 *
 * Look for environment variables
 * signal handling and graceful exit if SIGTERM
 * setup client webinstance for /ping
 * register the 'system' to the init system at INIT_SYSTEM_ADDR/PORT etc
 * send periodic health, config update, restart, de-reg
 * run GRPC server to:
 *   - handle queries from other services about particular system (via init)
 */
int main (int argc, char *argv[]) {

	int exitStatus=0;
	char *debug=DEFAULT_LOG_LEVEL;
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

	/* Step-3: system registration with init */
	if (send_request_to_init(REQ_REGISTER, config) != TRUE) {
		log_error("Error registrating with the init system");
		exitStatus = 1;
		goto exit_program;
	}

	/* Wait here for ever. XXX */

	log_debug("initClient running ...");

	getchar(); /* For now. */

	log_debug("Goodbye ... ");

	send_request_to_init(REQ_UNREGISTER, config);
	ulfius_stop_framework(&webInst);
	ulfius_clean_instance(&webInst);

	clear_config(config);

 exit_program:
	free(state);

	return exitStatus;
}
