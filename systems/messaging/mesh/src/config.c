/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */
#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <stdio.h>

#include "mesh.h"
#include "config.h"
#include "log.h"

/*
 * read_config_from_env -- read configuration params from the env variables
 *
 */
int read_config_from_env(Config **config) {

    char *websocketPort=NULL, *servicesPort=NULL;
    char *amqpHost=NULL, *amqpPort=NULL, *amqpUser=NULL, *amqpPassword=NULL;
	char *initClientHost=NULL, *initClientPort=NULL;
	char *certFile=NULL, *keyFile=NULL;
    char *orgName=NULL, *orgID=NULL;
    char *bindingIP=NULL;

    if ((bindingIP = getenv(ENV_BINDING_IP)) == NULL ||
        (websocketPort = getenv(ENV_WEBSOCKET_PORT)) == NULL ||
        (servicesPort = getenv(ENV_SERVICES_PORT)) == NULL ||
        (amqpHost = getenv(ENV_AMQP_HOST)) == NULL ||
        (amqpPort = getenv(ENV_AMQP_PORT)) == NULL ||
        (amqpUser = getenv(ENV_AMQP_USER)) == NULL ||
        (amqpPassword = getenv(ENV_AMQP_PASSWORD)) == NULL ||
        (initClientHost = getenv(ENV_INIT_SYSTEM_ADDR)) == NULL ||
        (initClientPort = getenv(ENV_INIT_SYSTEM_PORT)) == NULL ||
        (orgName = getenv(ENV_SYSTEM_ORG)) == NULL ||
        (orgID   = getenv(ENV_SYSTEM_ORG_ID)) == NULL) {
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

    (*config)->bindingIP      = strdup(bindingIP);
	(*config)->websocketPort  = strdup(websocketPort);
    (*config)->servicesPort   = strdup(servicesPort);
    (*config)->amqpHost       = strdup(amqpHost);
    (*config)->amqpPort       = strdup(amqpPort);
    (*config)->amqpUser       = strdup(amqpUser);
    (*config)->amqpPassword   = strdup(amqpPassword);
    (*config)->amqpExchange   = strdup(DEFAULT_MESH_AMQP_EXCHANGE);
    (*config)->initClientHost = strdup(initClientHost);
    (*config)->initClientPort = strdup(initClientPort);
    (*config)->orgName        = strdup(orgName);
    (*config)->orgID          = strdup(orgID);

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

    log_debug("Ukama org name: %s",  config->orgName);
    log_debug("Ukama org ID:   %s",  config->orgID);
    log_debug("Binding IP:     %s",  config->bindingIP);
	log_debug("Websocket port: %s",  config->websocketPort);
    log_debug("Services port:  %s",  config->servicesPort);
	log_debug("AMQP: %s:***@%s:%s",   config->amqpUser, config->amqpHost, config->amqpPort);
	log_debug("initClient: %s:%s",   config->initClientHost,
                                     config->initClientPort);
    log_debug("Cert file: %s",       config->certFile);
    log_debug("Key file:  %s",       config->keyFile);
}

void clear_config(Config *config) {

	if (!config) return;

    free(config->bindingIP);
    free(config->websocketPort);
    free(config->servicesPort);
	free(config->amqpHost);
	free(config->amqpPort);
	free(config->amqpUser);
	free(config->amqpPassword);
	free(config->amqpExchange);
    free(config->initClientHost);
    free(config->initClientPort);
	free(config->certFile);
	free(config->keyFile);
    free(config->orgName);
    free(config->orgID);
}
