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

#ifndef MSG_CLIENT_CONFIG_H
#define MSG_CLIENT_CONFIG_H

#define ENV_MSG_CLIENT_LOG_LEVEL "MSG_CLIENT_LOG_LEVEL"
#define ENV_MSG_CLIENT_IP        "MSG_CLIENT_IP"
#define ENV_MSG_CLIENT_PORT      "MSG_CLIENT_PORT"

#define ENV_MSG_CLIENT_AMQP_LOGIN  "MSG_CLIENT_AMQP_LOGIN"
#define ENV_MSG_CLIENT_AMQP_PASSWD "MSG_CLIENT_AMQP_PASSWD"

/* Struct to various env variables and runtime config parameters */
typedef struct {

	char *logLevel; /* Log level */
	char *ip;       /* IP bind */
	char *port;     /* Port listen */
	char *login;    /* AMQP login */
	char *passwd;   /* AMQP password */
} Config;

void print_config(Config *config);
void clear_config(Config *config);
int read_config_from_env(Config **config);

#endif /* MSG_CLIENT_CONFIG_H */
