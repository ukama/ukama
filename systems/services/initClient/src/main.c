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
pthread_t reconciler = 0;
int globalInit = 0;
static volatile sig_atomic_t terminating = 0;

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
            "\t ENV_INIT_RECONCILE_PERIOD   = %s\n"
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
            ENV_STR(getenv(ENV_INIT_RECONCILE_PERIOD)),
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

/* Termination is only flagged here; the unregister and the teardown run on the
 * main thread, where they are allowed to take locks and talk to init. */
void signal_term_handler(int signum) {

	(void)signum;

	terminating = 1;
}

void catch_sigterm(void) {

	static struct sigaction saction;

    memset(&saction, 0, sizeof(saction));

    saction.sa_handler = signal_term_handler;
	sigemptyset(&saction.sa_mask);
    saction.sa_flags   = 0;

    sigaction(SIGTERM, &saction, NULL);
    sigaction(SIGINT,  &saction, NULL);
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
	json_decref(json);
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
	int ret = TRUE;

	if (uuid == NULL) return FALSE;

	/* The file holds both the local and the global registration id, so the one
	 * not being written here has to be carried over rather than dropped. */
	if (!parse_cache_uuid(fileName, &sysReg) || sysReg == NULL) {
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
		ret = FALSE;
	}

	free_system_registration(sysReg);

	return ret;
}

/* Registers this system with init and records the returned registration id as
 * this process's proof of ownership. Re-registering refreshes every field and
 * yields a new id, which is what stops a replaced instance from removing the
 * record of the instance that took over from it. */
static int claim_registration(Config *config, int global) {

	int ret=TRUE;
	char *response=NULL;
	QueryResponse *queryResponse=NULL;

	log_info("Sending registration request for system %s for org %s",
	         config->systemName, config->systemOrg);

	if (send_request_to_init(REQ_REGISTER, config, config->systemOrg,
	                         NULL, &response, global) != TRUE) {
		log_error("Error registering with the init system");
		return FALSE;
	}

	if (deserialize_response(REQ_REGISTER, &queryResponse, response) != TRUE) {
		log_error("Error deserialize the registration response. Str: %s",
		          response);
		ret = FALSE;
		goto return_function;
	}

	log_info("Storing registration id %s for system %s for org %s in %s",
	         queryResponse->systemID, config->systemName, config->systemOrg,
	         config->tempFile);

	if (!store_cache_uuid(config->tempFile, queryResponse->systemID, global)) {
		log_error("Unable to store registration id in %s", config->tempFile);
		ret = FALSE;
	}

 return_function:
	free_query_response(queryResponse);
	if (response) free(response);

	return ret;
}

int register_system(Config *config, int global){

	int regStatus=REG_STATUS_NONE, ret=TRUE;
	char *cacheUUID=NULL, *systemUUID=NULL;

	pthread_mutex_lock(&registrationLock);

	regStatus = existing_registration(config, &cacheUUID, &systemUUID, global);

	/* At start-up the record is left alone only when it is already current and
	 * carries the id this process holds. Every other state is claimed, so a
	 * freshly started instance always owns what it is serving. */
	if ((regStatus & REG_STATUS_MATCH) &&
	    (regStatus & REG_STATUS_HAVE_UUID) &&
	    !(regStatus & REG_STATUS_UUID_MISMATCH)) {
		log_debug("System %s already registered with init as %s",
		          config->systemName, cacheUUID);
	} else {
		ret = claim_registration(config, global);
	}

	pthread_mutex_unlock(&registrationLock);

	if (cacheUUID)  free(cacheUUID);
	if (systemUUID) free(systemUUID);

	return ret;
}

/* One reconcile pass: restores the record if it has gone missing and repairs it
 * if it has drifted, but never takes a record another instance owns. */
static int reconcile_registration(Config *config, int global) {

	int regStatus=REG_STATUS_NONE, ret=TRUE;
	char *cacheUUID=NULL, *systemUUID=NULL;

	pthread_mutex_lock(&registrationLock);

	regStatus = existing_registration(config, &cacheUUID, &systemUUID, global);

	if (regStatus & REG_STATUS_PARSING_FAILURE) {
		log_error("Unable to read the init record for system %s. "
		          "Skipping this reconcile pass", config->systemName);
		ret = FALSE;
	} else if (regStatus & REG_STATUS_NO_RECORD) {
		log_info("System %s is missing from init. Registering it again",
		         config->systemName);
		ret = claim_registration(config, global);
	} else if (regStatus & REG_STATUS_UUID_MISMATCH) {
		log_info("Init record for system %s belongs to another instance "
		         "(init holds %s, this instance registered %s). Skipping",
		         config->systemName,
		         (systemUUID) ? systemUUID : "null", cacheUUID);
	} else if (regStatus & REG_STATUS_NO_MATCH) {
		log_info("Init record for system %s no longer matches its "
		         "configuration. Registering it again", config->systemName);
		ret = claim_registration(config, global);
	}

	pthread_mutex_unlock(&registrationLock);

	if (cacheUUID)  free(cacheUUID);
	if (systemUUID) free(systemUUID);

	return ret;
}

static void* reconcile_loop(void *args) {

	Config *config = (Config *)args;

	block_termination_signals();

	while (TRUE) {

		sleep(config->reconcilePeriod);

		if (!reconcile_registration(config, REGISTER_TO_LOCAL_INIT)) {
			log_error("Reconcile with local init failed for system %s",
			          config->systemName);
		}

		if (config->globalInitSystemEnable) {
			if (!reconcile_registration(config, REGISTER_TO_GLOBAL_INIT)) {
				log_error("Reconcile with global init failed for system %s",
				          config->systemName);
			}
		}
	}

	return NULL;
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
	globalInit    = config->globalInitSystemEnable;

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

	/* Registration is re-checked on a timer as well as on an address change, so
	 * a record that disappears from init is restored within one period instead
	 * of staying gone until this container is replaced. */
	if (pthread_create(&reconciler, NULL, reconcile_loop, config) != 0) {
		log_error("Unable to start the registration reconcile thread");
		exitStatus = 1;
		goto exit_program;
	}

	if (config->systemDNS) {
		pthread_create(&child, NULL, refresh_lookup, config);
	}

	while (!terminating) {
		pause();
	}

	log_debug("Exiting initClient ... ");

	unregister_system(config, REGISTER_TO_LOCAL_INIT);
	if (config->globalInitSystemEnable) {
		unregister_system(config, REGISTER_TO_GLOBAL_INIT);
	}

	if (reconciler) {
		pthread_cancel(reconciler);
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
