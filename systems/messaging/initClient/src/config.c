/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "initClient.h"
#include "config.h"
#include "log.h"

/*
 * read_config_from_env -- read configuration params from the env variables
 *
 */
int read_config_from_env(Config **config){

	char *ip=NULL, *port=NULL;
	char *systemName=NULL, *systemAddr=NULL, *systemPort=NULL;
	char *initSystemAddr=NULL, *initSystemPort=NULL;
	char *orgName=NULL, *systemCert=NULL, *apiVersion=NULL;

	if ((ip = getenv(ENV_INIT_CLIENT_IP)) == NULL ||
		(port = getenv(ENV_INIT_CLIENT_PORT)) == NULL) {
		log_error("%s and/or %s env variables not defined",
				  ENV_INIT_CLIENT_IP, ENV_INIT_CLIENT_PORT);
		return FALSE;
	}

	if ((systemName = getenv(ENV_INIT_CLIENT_SYSTEM_NAME)) == NULL ||
		(systemAddr = getenv(ENV_INIT_CLIENT_SYSTEM_ADDR)) == NULL ||
		(systemPort = getenv(ENV_INIT_CLIENT_SYSTEM_PORT)) == NULL ||
		(systemCert = getenv(ENV_INIT_CLIENT_SYSTEM_CERT)) == NULL ||
		(initSystemAddr = getenv(ENV_INIT_SYSTEM_ADDR)) == NULL ||
		(initSystemPort = getenv(ENV_INIT_SYSTEM_PORT)) == NULL ) {
	    log_error("Required env variables not defined");
		return FALSE;
	}

	if ((orgName = getenv(ENV_INIT_CLIENT_SYSTEM_ORG)) == NULL) {
		orgName = DEFAULT_SYSTEM_ORG;
	}

	if ((apiVersion = getenv(ENV_INIT_SYSTEM_API_VER)) == NULL) {
		apiVersion = DEFAULT_API_VER;
	}

	*config = (Config *)calloc(1, sizeof(Config));
	if (*config == NULL) {
		log_error("Memory allocation failure: %d", sizeof(Config));
		return FALSE;
	}

	(*config)->logLevel   = getenv(ENV_INIT_CLIENT_LOG_LEVEL);
	(*config)->ip         = strdup(ip);
	(*config)->port       = strdup(port);
	(*config)->systemOrg  = strdup(orgName);
	(*config)->systemName = strdup(systemName);
	(*config)->systemAddr = strdup(systemAddr);
	(*config)->systemPort = strdup(systemPort);
	(*config)->systemCert = strdup(systemCert);

	(*config)->initSystemAPIVer = strdup(apiVersion);
	(*config)->initSystemAddr   = strdup(initSystemAddr);
	(*config)->initSystemPort   = strdup(initSystemPort);

	if (!(*config)->logLevel) {
		log_debug("Log level not defined, setting to default: DEBUG");
		(*config)->logLevel = DEF_LOG_LEVEL;
	}

	return TRUE;
}

/*
 * clear_config --
 *
 */
void clear_config(Config *config) {

	if (config == NULL) return;

	if (config->ip)             free(config->ip);
	if (config->port)           free(config->port);
	if (config->systemOrg)      free(config->systemOrg);
	if (config->systemName)     free(config->systemName);   
	if (config->systemAddr)     free(config->systemAddr);
	if (config->systemPort)     free(config->systemPort);
	if (config->systemCert)     free(config->systemCert);
	if (config->initSystemAddr) free(config->initSystemAddr);
	if (config->initSystemPort) free(config->initSystemPort);

	free(config);
}
