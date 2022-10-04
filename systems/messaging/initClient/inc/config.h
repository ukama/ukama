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

#ifndef INIT_CLIENT_CONFIG_H
#define INIT_CLIENT_CONFIG_H

#define DEFAULT_SYSTEM_ORG "Ukama"
#define DEFAULT_API_VER    "v1"

#define ENV_INIT_CLIENT_LOG_LEVEL "INIT_CLIENT_LOG_LEVEL"
#define ENV_INIT_CLIENT_IP        "INIT_CLIENT_IP"
#define ENV_INIT_CLIENT_PORT      "INIT_CLIENT_PORT"

#define ENV_INIT_CLIENT_SYSTEM_ORG  "ENV_INIT_CLIENT_SYSTEM_ORG"
#define ENV_INIT_CLIENT_SYSTEM_NAME "ENV_INIT_CLIENT_SYSTEM_NAME"
#define ENV_INIT_CLIENT_SYSTEM_ADDR "ENV_INIT_CLIENT_SYSTEM_ADDR"
#define ENV_INIT_CLIENT_SYSTEM_PORT "ENV_INIT_CLIENT_SYSTEM_PORT"
#define ENV_INIT_CLIENT_SYSTEM_CERT "ENV_INIT_CLIENT_SYSTEM_CERT"

#define ENV_INIT_SYSTEM_API_VER     "ENV_INIT_SYSTEM_API_VER"
#define ENV_INIT_SYSTEM_ADDR        "ENV_INIT_SYSTEM_ADDR"
#define ENV_INIT_SYSTEM_PORT        "ENV_INIT_SYSTEM_PORT"

/* Struct to various env variables and runtime config parameters */
typedef struct {

	char *logLevel;   /* Log level */
	char *ip;         /* IP bind */
	char *port;       /* Port listen */

    char *systemOrg;  /* Name of the organization the system belongs */
	char *systemName; /* Name of the system */
	char *systemAddr; /* address where system can be reached at */
	char *systemPort; /* port where system can be reached at */
    char *systemCert; /* Certificate for the system */

	char *initSystemAPIVer; /* API version for init system */
	char *initSystemAddr;   /* address for init system */
	char *initSystemPort;   /* port for init system */
} Config;

void clear_config(Config *config);
int read_config_from_env(Config **config);

#endif /* INIT_CLIENT_CONFIG_H */
