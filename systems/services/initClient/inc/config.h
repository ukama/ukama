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

#define ENV_INIT_CLIENT_ADDR "ENV_INIT_CLIENT_ADDR"
#define ENV_INIT_CLIENT_PORT "ENV_INIT_CLIENT_PORT"

#define ENV_SYSTEM_NAME "ENV_SYSTEM_NAME"
#define ENV_SYSTEM_ADDR "ENV_SYSTEM_ADDR"
#define ENV_SYSTEM_PORT "ENV_SYSTEM_PORT"
#define ENV_SYSTEM_CERT "ENV_SYSTEM_CERT"
#define ENV_SYSTEM_ORG  "ENV_SYSTEM_ORG"

#define ENV_INIT_SYSTEM_ADDR "ENV_INIT_SYSTEM_ADDR"
#define ENV_INIT_SYSTEM_PORT "ENV_INIT_SYSTEM_PORT"
#define ENV_INIT_SYSTEM_API  "ENV_INIT_SYSTEM_API"

#define ENV_INIT_CLIENT_LOG_LEVEL "ENV_INIT_CLIENT_LOG_LEVEL"
#define ENV_INIT_CLIENT_TEMP_FILE "ENV_INIT_CLIENT_TEMP_FILE"

/* Struct to various env variables and runtime config parameters */
typedef struct {

	char *logLevel;   /* Log level */
	char *tempFile;   /* Temp file to log stuff */

	char *addr;       /* initClient bind address */
	char *port;       /* initClient listening port */

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
