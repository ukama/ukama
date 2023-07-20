/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * config.h
 */

#ifndef MESH_CONFIG_H
#define MESH_CONFIG_H

#define DEFAULT_LOG_LEVEL "DEBUG"

typedef struct {

	char *websocketPort;   /* to accept nodes via websocket */
	char *servicesPort;    /* to accept services */

	char *amqpHost;       /* Host where AMQP exchange is running (IP) */
	char *amqpPort;       /* Port where AMQP exchange is listening */
	char *amqpExchange;   /* AMQP exchange name */

	char *initClientHost; /* Host where initClient is running (IP) */
	char *initClientPort; /* Port where initClient is listening */

	char *certFile;       /* CA Cert file name. */
	char *keyFile;        /* Key file name.*/

    char *logLevel;
} Config;

void clear_config(Config *config);
void print_config(Config *config);
int read_config_from_env(Config **config);

#endif /* MESH_CONFIG_H */
