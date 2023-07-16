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
 * Description of various env variables:
 *
 * ENV_SYSTEM_NAME - name of my system
 * ENV_SYSTEM_ADDR - address for the api gw of my system
 * ENV_SYSTEM_PORT - listen port for the api gw
 * ENV_SYSTEM_CERT - certificate.
 *
 * ENV_INIT_SYSTEM_ADDR - init system address (api-gw)
 * ENV_INIT_SYSTEM_PORT - init system port (api-gw)
 */

/*
 * read_config_from_env -- read configuration params from the env variables
 *
 */
int read_config_from_env(Config **config){

	char *port=NULL, *addr=NULL, *tempFile=NULL;
	char *systemName=NULL, *systemAddr=NULL, *systemPort=NULL;
	char *initSystemAddr=NULL, *initSystemPort=NULL;
	char *systemOrg=NULL, *systemCert=NULL, *apiVersion=NULL;

	if ((addr = getenv(ENV_INIT_CLIENT_ADDR)) == NULL ||
		(port = getenv(ENV_INIT_CLIENT_PORT)) == NULL ||
		(tempFile = getenv(ENV_INIT_CLIENT_TEMP_FILE)) == NULL) {
		log_error("Required env variables: %s %s %s missing",
				  ENV_INIT_CLIENT_ADDR,
				  ENV_INIT_CLIENT_PORT,
				  ENV_INIT_CLIENT_TEMP_FILE);
		return FALSE;
	}

	if ((systemName = getenv(ENV_SYSTEM_NAME)) == NULL ||
		(systemAddr = getenv(ENV_SYSTEM_ADDR)) == NULL ||
		(systemPort = getenv(ENV_SYSTEM_PORT)) == NULL ||
		(systemCert = getenv(ENV_SYSTEM_CERT)) == NULL ||
		(initSystemAddr = getenv(ENV_INIT_SYSTEM_ADDR)) == NULL ||
		(initSystemPort = getenv(ENV_INIT_SYSTEM_PORT)) == NULL ) {
	    log_error("Required env variables not defined");
		return FALSE;
	}

	if ((systemOrg = getenv(ENV_SYSTEM_ORG)) == NULL) {
		systemOrg = DEFAULT_SYSTEM_ORG;
	}

	if ((apiVersion = getenv(ENV_INIT_SYSTEM_API)) == NULL) {
		apiVersion = DEFAULT_API_VER;
	}

	*config = (Config *)calloc(1, sizeof(Config));
	if (*config == NULL) {
		log_error("Memory allocation failure: %d", sizeof(Config));
		return FALSE;
	}

	(*config)->logLevel   = getenv(ENV_INIT_CLIENT_LOG_LEVEL);

	(*config)->addr     = strdup(addr);
	(*config)->port     = strdup(port);
	(*config)->tempFile = strdup(tempFile);

	(*config)->systemOrg  = strdup(systemOrg);
	(*config)->systemName = strdup(systemName);
	(*config)->systemAddr = strdup(systemAddr);
	(*config)->systemPort = strdup(systemPort);
	(*config)->systemCert = strdup(systemCert);

	(*config)->initSystemAPIVer = strdup(apiVersion);
	(*config)->initSystemAddr   = strdup(initSystemAddr);
	(*config)->initSystemPort   = strdup(initSystemPort);

	if (!(*config)->logLevel) {
		log_debug("Log level not defined, setting to default: DEBUG");
		(*config)->logLevel = DEFAULT_LOG_LEVEL;
	}

	return TRUE;
}

/*
 * clear_config --
 *
 */
void clear_config(Config *config) {

	if (config == NULL) return;

	if (config->addr)             free(config->addr);
	if (config->port)             free(config->port);
	if (config->tempFile)         free(config->tempFile);

	if (config->systemOrg)        free(config->systemOrg);
	if (config->systemName)       free(config->systemName);
	if (config->systemAddr)       free(config->systemAddr);
	if (config->systemPort)       free(config->systemPort);
	if (config->systemCert)       free(config->systemCert);
	if (config->initSystemAddr)   free(config->initSystemAddr);
	if (config->initSystemPort)   free(config->initSystemPort);
	if (config->initSystemAPIVer) free(config->initSystemAPIVer);

	free(config);
}
