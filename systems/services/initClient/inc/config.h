/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#ifndef INIT_CLIENT_CONFIG_H
#define INIT_CLIENT_CONFIG_H

#define ENV_INIT_CLIENT_ADDR "ENV_INIT_CLIENT_ADDR"
#define ENV_INIT_CLIENT_PORT "ENV_INIT_CLIENT_PORT"

#define ENV_DNS_SERVER	"ENV_DNS_SERVER"
#define ENV_SYSTEM_NAME "ENV_SYSTEM_NAME"
#define ENV_SYSTEM_ADDR "ENV_SYSTEM_ADDR"
#define ENV_SYSTEM_DNS  "ENV_SYSTEM_DNS"
#define ENV_SYSTEM_PORT "ENV_SYSTEM_PORT"
#define ENV_SYSTEM_CERT "ENV_SYSTEM_CERT"
#define ENV_SYSTEM_ORG  "ENV_SYSTEM_ORG"

#define ENV_SYSTEM_NODE_GW_ADDR "ENV_SYSTEM_NODE_GW_ADDR"
#define ENV_SYSTEM_NODE_GW_PORT "ENV_SYSTEM_NODE_GW_PORT"

#define ENV_INIT_SYSTEM_ADDR "ENV_INIT_SYSTEM_ADDR"
#define ENV_INIT_SYSTEM_PORT "ENV_INIT_SYSTEM_PORT"

#define ENV_DNS_REFRESH_TIME_PERIOD "ENV_DNS_REFRESH_TIME_PERIOD"
#define ENV_GLOBAL_INIT_ENABLE "ENV_GLOBAL_INIT_ENABLE"
#define ENV_GLOBAL_INIT_SYSTEM_ADDR "ENV_GLOBAL_INIT_SYSTEM_ADDR"
#define ENV_GLOBAL_INIT_SYSTEM_PORT "ENV_GLOBAL_INIT_SYSTEM_PORT"

#define ENV_INIT_SYSTEM_API  "ENV_INIT_SYSTEM_API"

#define ENV_INIT_CLIENT_LOG_LEVEL "ENV_INIT_CLIENT_LOG_LEVEL"
#define ENV_INIT_CLIENT_TEMP_FILE "ENV_INIT_CLIENT_TEMP_FILE"

#define GLOBAL_INIT_SYSTEM_ENABLE_STR "true"
#define GLOBAL_INIT_SYSTEM_DISABLE_STR "false"

#define GLOBAL_INIT_SYSTEM_ENABLE   1
#define GLOBAL_INIT_SYSTEM_DISABLE	0

#define DEFAULT_TIME_PERIOD 10
#define DEFAULT_DNS_SERVER "127.0.0.1"

/* Struct to various env variables and runtime config parameters */
typedef struct {

	char *logLevel;   /* Log level */
	char *tempFile;   /* Temp file to log stuff */
	char *dnsServer;  /* DNS server */
	int  timePeriod;  /* time period in seconds.*/

	char *addr;       /* initClient bind address */
	char *port;       /* initClient listening port */

	char *systemOrg;  /* Name of the organization the system belongs */
	char *systemName; /* Name of the system */
	char *systemDNS;  /* DNS for system */
	char *systemAddr; /* address where system can be reached at */
	char *systemPort; /* port where system can be reached at */
	char *systemCert; /* Certificate for the system */

    char *systemNodeGWAddr; /* address where system's node-gw is available */
    char *systemNodeGWPort; /* poer where system's node-gw port is */

    char *nameServer;
	char *initSystemAPIVer; /* API version for init system */
	char *initSystemAddr;   /* address for init system */
	char *initSystemPort;   /* port for init system */
	int  globalInitSystemEnable; /* 1: if global init system is enabled else 0 */
	char *globalInitSystemAddr;   /* address for global init system */
	char *globalInitSystemPort;   /* port for global init system */
} Config;

typedef int (*UpdateIpCallback)(Config *config);
void clear_config(Config *config);
int read_config_from_env(Config **config);
char* nslookup(char* name, char *server);
void* refresh_lookup(void* args);
void register_callback(UpdateIpCallback cb);
char* parse_resolveconf();
#endif /* INIT_CLIENT_CONFIG_H */
