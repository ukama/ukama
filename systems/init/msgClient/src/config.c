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

#include "msgClient.h"
#include "config.h"
#include "log.h"

/*
 * read_config_from_env -- read configuration params from the env variables
 *
 */
int read_config_from_env(Config **config){

	char *ip=NULL, *port=NULL;

	if ((ip = getenv(ENV_MSG_CLIENT_IP)) == NULL ||
		(port = getenv(ENV_MSG_CLIENT_PORT)) == NULL) {
		log_error("%s and/or %s env variables not defined",
				  ENV_MSG_CLIENT_IP, ENV_MSG_CLIENT_PORT);
		return FALSE;
	}
	
	*config = (Config *)calloc(1, sizeof(Config));
	if (*config == NULL) {
		log_error("Memory allocation failure: %d", sizeof(Config));
		return FALSE;
	}

	(*config)->logLevel = getenv(ENV_MSG_CLIENT_LOG_LEVEL);
	(*config)->ip       = strdup(ip);
	(*config)->port     = strdup(port);

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

	if (config->ip)   free(config->ip);
	if (config->port) free(config->port);

	free(config);
}
