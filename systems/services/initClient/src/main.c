/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
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
#include <errno.h>

#include "initClient.h"
#include "config.h"
#include "jserdes.h"
#include "log.h"

#define VERSION "0.0.1"

typedef struct {

	struct _u_instance *webInst;
	Config             *config;
} State;

extern int start_web_services(Config *config, UInst *webtInst); /*network.c */

/* Global */
State *state=NULL;
pthread_t child = 0;
int globalInit = 0;

#define ENV_STR(x) ((x) ? (x) : "(unset)")

void usage() {

	fprintf(stdout, "Usage: initClient [options] \n");
	fprintf(stdout, "Options:\n");
	fprintf(stdout, "--h, --help     this menu\n");
	fprintf(stdout, "--V, --version  Version\n");
	fprintf(stdout, "Environment variable used are: \n");
    fprintf(stdout,
            "\t ENV_INIT_CLIENT_LOG_LEVEL   = %s\n"
            "\t ENV_SYSTEM_ORG              = %s\n"
            "\t ENV_SYSTEM_NAME             = %s\n"
            "\t ENV_SYSTEM_DNS              = %s\n"
            "\t ENV_SYSTEM_ADDR             = %s\n"
            "\t ENV_SYSTEM_PORT             = %s\n"
            "\t ENV_SYSTEM_NODE_GW_ADDR     = %s\n"
            "\t ENV_SYSTEM_NODE_GW_PORT     = %s\n"
            "\t ENV_INIT_SYSTEM_ADDR        = %s\n"
            "\t ENV_INIT_SYSTEM_PORT        = %s\n"
            "\t ENV_GLOBAL_INIT_ENABLE      = %s\n"
            "\t ENV_GLOBAL_INIT_SYSTEM_ADDR = %s\n"
            "\t ENV_GLOBAL_INIT_SYSTEM_PORT = %s\n",
            ENV_STR(getenv(ENV_INIT_CLIENT_LOG_LEVEL)),
            ENV_STR(getenv(ENV_SYSTEM_ORG)),
            ENV_STR(getenv(ENV_SYSTEM_NAME)),
            ENV_STR(getenv(ENV_SYSTEM_DNS)),
            ENV_STR(getenv(ENV_SYSTEM_ADDR)),
            ENV_STR(getenv(ENV_SYSTEM_PORT)),
            ENV_STR(getenv(ENV_SYSTEM_NODE_GW_ADDR)),
            ENV_STR(getenv(ENV_SYSTEM_NODE_GW_PORT)),
            ENV_STR(getenv(ENV_INIT_SYSTEM_ADDR)),
            ENV_STR(getenv(ENV_INIT_SYSTEM_PORT)),
            ENV_STR(getenv(ENV_GLOBAL_INIT_ENABLE)),
            ENV_STR(getenv(ENV_GLOBAL_INIT_SYSTEM_ADDR)),
            ENV_STR(getenv(ENV_GLOBAL_INIT_SYSTEM_PORT))
        );
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

void signal_term_handler(void) {

    char *response=NULL;

	if (state == NULL) exit(1);

	/* un-register the system */
	if (send_request_to_init(REQ_UNREGISTER,
                             state->config,
                             state->config->systemOrg, NULL,
							 &response, REGISTER_TO_LOCAL_INIT) != TRUE) {
		log_error("Error registrating with the init system");
	}

	if (globalInit) {
		if (send_request_to_init(REQ_UNREGISTER,
                                 state->config,
                                 state->config->systemOrg, NULL,
                                 &response, REGISTER_TO_GLOBAL_INIT) != TRUE) {
			log_error("Error registrating with the init system");
		}
	}

	if (child)	{
		pthread_cancel(child);
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

void catch_sigterm(void) {

	static struct sigaction saction;

    memset(&saction, 0, sizeof(saction));

    saction.sa_sigaction = signal_term_handler;
	sigemptyset(&saction.sa_mask);
    saction.sa_flags     = 0;

    sigaction(SIGTERM, &saction, NULL);
}

int create_temp_file_and_store_uuid(char *fileName, SystemRegistrationId* sysReg) {

	json_t *json = NULL;
	FILE *fp=NULL;
	char* str = NULL;

	if ((fp = fopen(fileName, "w")) == NULL) {
		log_error("Unable to create cache temp file: %s Error: %s",
				fileName, strerror(errno));
		return FALSE;
	}

	if (!serialize_uuids_from_file(sysReg, &json)) {
		log_error("Error serializing registration status in file : %s Error :%s",
				fileName, strerror(errno));
		return REG_STATUS_NO_UUID;
	}

	str = json_dumps(json, 0);
	if (str) {
		fputs(str, fp);
		free(str);
	} else {
		log_error("Unable to create cache temp file: %s Error: %s",
						fileName, strerror(errno));
		return FALSE;
	}
	fclose(fp);

	return TRUE;
}

int store_cache_uuid(char *fileName, char* uuid, int global) {

	SystemRegistrationId *sysReg = NULL;

	if (!parse_cache_uuid(fileName, sysReg)) {
		/* Parsing Failed this means problem with file */
		sysReg = (SystemRegistrationId*)calloc(1, sizeof(SystemRegistrationId));
	}

	if (!sysReg) {
		return FALSE;
	}

	if (global) {
		if (sysReg->globalUUID) free(sysReg->globalUUID);
		sysReg->globalUUID = strdup(uuid);
	} else {
		if (sysReg->localUUID) free(sysReg->localUUID);
		sysReg->localUUID = strdup(uuid);
	}

	log_debug("Creating file %s", fileName);
	if (!create_temp_file_and_store_uuid(fileName, sysReg)) {
		return FALSE;
	}
	return TRUE;
}

int register_system(Config *config, int global){

	int regStatus=REG_STATUS_NONE;
	char *response=NULL;
	char *cacheUUID=NULL, *systemUUID=NULL;
	QueryResponse *queryResponse=NULL;

	/* Step-1: check current registration status */
	regStatus = existing_registration(config, &cacheUUID, &systemUUID, global);

	/* Step-2: take action(s) */
	switch(regStatus) {
	case REG_STATUS_MATCH | REG_STATUS_HAVE_UUID:
	log_debug("System already registerd with init.");
	break;

	case REG_STATUS_MATCH | REG_STATUS_NO_UUID:
	log_debug("Storing UUID %s to tempFile: %s", systemUUID,
			config->tempFile);
	store_cache_uuid(config->tempFile,
			systemUUID, global);

	break;

	case (REG_STATUS_NO_MATCH | REG_STATUS_HAVE_UUID):
		log_info("Sending registration request for system %s for org %s",
                 config->systemName,
                 config->systemOrg);
        if (send_request_to_init(REQ_UPDATE,
                                 config,config->systemOrg,
                                 NULL, &response, global) != TRUE) {
				log_error("Error updating with the init system");
				return FALSE;
			}
	break;

	case (REG_STATUS_NO_MATCH | REG_STATUS_NO_UUID):
	case REG_STATUS_NO_MATCH:
		/* first time registering */
		log_info("Sending registration request for system %s for org %s",
                 config->systemName,
                 config->systemOrg);
		if (send_request_to_init(REQ_REGISTER,
                                 config, config->systemOrg,
                                 NULL, &response, global) != TRUE) {
			log_error("Error registering with the init system");
			return FALSE;
		}

		/* read the UUID and log it into tempfile. */
		if (deserialize_response(REQ_REGISTER,
                                 &queryResponse,
                                 response) != TRUE) {
			log_error("Error deserialize the registration response. Str: %s",
                      response);
			return FALSE;
		}

		log_info("Storing registration response for system %s for org %s in %s",
                 config->systemName,
                 config->systemOrg,
                 config->tempFile);
        
		store_cache_uuid(config->tempFile, 
                         queryResponse->systemID,
                         global);
		break;

	default:
		break;
	}

	if (queryResponse) free_query_response(queryResponse);
	if (response)      free(response);

	return TRUE;
}

int register_to_inits(Config *config) {

	/* registration process for local Init */
	if (!register_system(config, REGISTER_TO_LOCAL_INIT)) {
		return 1;
	}

	/* registration process for global Init */
	if (config->globalInitSystemEnable) {
		/* registration process for global Init */
		if (!register_system(config, REGISTER_TO_GLOBAL_INIT)) {
			return 1;
		}
	}
	return 0;
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
	char *response=NULL;
	struct _u_instance webInst;
	Config *config=NULL;
	pthread_t child;
	int *childStatus;

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

	/* Step 2: register callback to update Inits */
	register_callback(&register_to_inits);

	/* Step-3: start webservice */
	if (start_web_services(config, &webInst) != TRUE) {
		log_error("Webservice failed to setup for clients. Exiting.");
		exitStatus = 1;
		goto exit_program;
	}

	/* Step-3: registration to init systems */
	exitStatus = register_to_inits(config);
	if (exitStatus) {
		goto exit_program;
	}

	log_debug("initClient running ...");

	if (config->systemDNS) {
		/* Start thread : Need a cleanup so that it's always 
         * dns no IP as arg for system */
		pthread_create(&child, NULL, refresh_lookup, config);
		pthread_join (child, (void **)&childStatus);
	} else {
        pause();
	}

	log_debug("Exiting initClient ... ");

	send_request_to_init(REQ_UNREGISTER,
                         config,
                         config->systemOrg,
                         NULL,
                         &response,
                         REGISTER_TO_LOCAL_INIT);
	if (config->globalInitSystemEnable) {
		send_request_to_init(REQ_UNREGISTER,
                             config,
                             config->systemOrg,
                             NULL,
                             &response,
                             REGISTER_TO_GLOBAL_INIT);
	}

	if (child) {
		pthread_cancel(child);
	}
	ulfius_stop_framework(&webInst);
	ulfius_clean_instance(&webInst);

	clear_config(config);

 exit_program:
	free(state);

	return exitStatus;
}
