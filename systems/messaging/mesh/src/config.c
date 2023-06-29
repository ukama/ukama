/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Config.c
 *
 */

#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <stdio.h>

#include "mesh.h"
#include "config.h"
#include "log.h"

/*
 * Various environment variables:
 * ENV_WEBSOCKET_PORT
 * ENV_SERVICE_PORT   
 * ENV_AMQP_HOST
 * ENV_AMQP_PORT 
 * ENV_INIT_CLIENT_HOST
 * ENV_INIT_CLIENT_PORT
 * ENV_MESH_CERT_FILE 
 * ENV_MESH_KEY_FILE
 *
 */

/*
 * read_config_from_env -- read configuration params from the env variables
 *
 */
int read_config_from_env(Config **config) {

    char *websocketPort=NULL, *servicesPort=NULL;
    char *amqpHost=NULL, *amqpPort=NULL;
	char *initClientHost=NULL, *initClientPort=NULL;
	char *certFile=NULL, *keyFile=NULL;

    if ((websocketPort = getenv(ENV_WEBSOCKET_PORT)) == NULL ||
        (servicesPort = getenv(ENV_SERVICES_PORT)) == NULL ||
        (amqpHost = getenv(ENV_AMQP_HOST)) == NULL ||
        (amqpPort = getenv(ENV_AMQP_PORT)) == NULL ||
        (initClientHost = getenv(ENV_INIT_CLIENT_HOST)) == NULL ||
        (initClientPort = getenv(ENV_INIT_CLIENT_PORT)) == NULL) {

        log_error("Required env variable not defined");
        return FALSE;
    }

    if ((certFile = getenv(ENV_MESH_CERT_FILE)) == NULL) {
        certFile = DEFAULT_MESH_CERT_FILE;
    }

    if ((keyFile = getenv(ENV_MESH_KEY_FILE)) == NULL) {
        keyFile = DEFAULT_MESH_KEY_FILE;
    }

    *config = (Config *)calloc(1, sizeof(Config));
    if (*config == NULL) {
        log_error("Memory allocation failure: %d", sizeof(Config));
        return FALSE;
    }

    (*config)->logLevel = getenv(ENV_MESH_LOG_LEVEL);

	(*config)->websocketPort  = strdup(websocketPort);
    (*config)->servicesPort   = strdup(servicesPort);
    (*config)->amqpHost       = strdup(amqpHost);
    (*config)->amqpPort       = strdup(amqpPort);
    (*config)->amqpExchange   = strdup(DEFAULT_MESH_AMQP_EXCHANGE);
    (*config)->initClientHost = strdup(initClientHost);
    (*config)->initClientPort = strdup(initClientPort);

    if (!(*config)->logLevel) {
        log_debug("Log level not defined, setting to default: DEBUG");
        (*config)->logLevel = DEFAULT_LOG_LEVEL;
    }

    return TRUE;
}

/*
 * print_config_data --
 *
 */
void print_config(Config *config) {

	log_debug("Websocket port: %s", config->websocketPort);
    log_debug("Services port:   %s", config->servicesPort);
	log_debug("AMQP: %s:%s", config->amqpHost, config->amqpPort);
	log_debug("initClient: %s:%s", config->initClientHost,
              config->initClientPort);
    log_debug("Cert file: %s", config->certFile);
    log_debug("Key file:  %s", config->keyFile);
}

/*
 * clear_config --
 */
void clear_config(Config *config) {

	if (!config) return;

    free(config->websocketPort);
    free(config->servicesPort);
	free(config->amqpHost);
	free(config->amqpPort);
	free(config->amqpExchange);
    free(config->initClientHost);
    free(config->initClientPort);
	free(config->certFile);
	free(config->keyFile);
}
