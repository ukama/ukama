/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#ifndef MESH_CONFIG_H
#define MESH_CONFIG_H

#define DEFAULT_LOG_LEVEL "DEBUG"

typedef struct {

    char *bindingIP;      /* binding IP for websocket */
	char *websocketPort;  /* to accept nodes via websocket */
	char *servicesPort;   /* to accept services */
    char *adminPort;      /* to accept admin services */

	char *amqpHost;       /* Host where AMQP exchange is running (IP) */
	char *amqpPort;       /* Port where AMQP exchange is listening */
	char *amqpUser;       /* User for AMQP connection */
	char *amqpPassword;   /* Password for AMQP connection */
	char *amqpExchange;   /* AMQP exchange name */

	char *initClientHost; /* Host where initClient is running (IP) */
	char *initClientPort; /* Port where initClient is listening */

	char *certFile;       /* CA Cert file name. */
	char *keyFile;        /* Key file name.*/

    char *orgName;        /* Ukama organization name */
    char *orgID;          /* and its ID */

    char *logLevel;
} Config;

void clear_config(Config *config);
void print_config(Config *config);
int read_config_from_env(Config **config);

#endif /* MESH_CONFIG_H */
