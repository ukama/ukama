/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
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
 * ENV_SYSTEM_NAME         - name of my system
 * ENV_SYSTEM_ADDR         - address for the api gw of my system
 * ENV_SYSTEM_PORT         - listen port for the api gw
 * ENV_SYSTEM_CERT         - certificate.
 * ENV_SYSTEM_NODE_GW_ADDR - address for the node gw of my system
 * ENV_SYSTEM_NODE_GW_PORT - listening port for the node gw
 *
 * (optional)
 * ENV_INIT_SYSTEM_ADDR - init system address (api-gw)
 * ENV_INIT_SYSTEM_PORT - init system port (api-gw)
 */

int read_config_from_env(Config **config){

	char *port=NULL, *addr=NULL, *tempFile=NULL;
	char *systemName=NULL, *systemAddr=NULL, *systemPort=NULL;
	char *initSystemAddr=NULL, *initSystemPort=NULL;
	char *globalInitSystemEnable=NULL, *globalInitSystemAddr=NULL, *globalInitSystemPort=NULL;
	char *systemOrg=NULL, *systemCert=NULL, *apiVersion=NULL;
	char *systemDNS=NULL, *timePeriod=NULL, *dnsServer=NULL, *nameServer=NULL;
	char *systemDNSNodeGw=NULL;
	char *systemNodeGwAddr=NULL, *systemNodeGwPort=NULL;

	int period = 0;
	int freeSystemAddr = 0;
	int freeSystemNodeGwAddr = 0;
	int nodegw_any = 0;

	addr     = getenv(ENV_INIT_CLIENT_ADDR);
	port     = getenv(ENV_INIT_CLIENT_PORT);
	tempFile = getenv(ENV_INIT_CLIENT_TEMP_FILE);
	if (!addr || !port || !tempFile) {
		log_error("Required env variables: %s %s %s missing",
		          ENV_INIT_CLIENT_ADDR,
		          ENV_INIT_CLIENT_PORT,
		          ENV_INIT_CLIENT_TEMP_FILE);
		return FALSE;
	}

	systemName     = getenv(ENV_SYSTEM_NAME);
	systemPort     = getenv(ENV_SYSTEM_PORT);
	systemCert     = getenv(ENV_SYSTEM_CERT);
	initSystemAddr = getenv(ENV_INIT_SYSTEM_ADDR);
	initSystemPort = getenv(ENV_INIT_SYSTEM_PORT);
	if (!systemName || !systemPort || !systemCert ||
	    !initSystemAddr || !initSystemPort) {
		log_error("Required env variables not defined");
		return FALSE;
	}

	systemDNSNodeGw  = getenv(ENV_SYSTEM_DNS_NODE_GW);
	systemNodeGwAddr = getenv(ENV_SYSTEM_NODE_GW_ADDR);
	systemNodeGwPort = getenv(ENV_SYSTEM_NODE_GW_PORT);
	nodegw_any = (systemDNSNodeGw != NULL) ||
	             (systemNodeGwAddr != NULL) ||
	             (systemNodeGwPort != NULL);

	if (nodegw_any) {
		if (!systemNodeGwPort) {
			log_error("NodeGW misconfigured: %s must be set",
			          ENV_SYSTEM_NODE_GW_PORT);
			return FALSE;
		}
		if (systemDNSNodeGw && systemNodeGwAddr) {
			log_error("NodeGW misconfigured: only one of %s or %s may be set",
			          ENV_SYSTEM_DNS_NODE_GW,
			          ENV_SYSTEM_NODE_GW_ADDR);
			return FALSE;
		}
		if (!systemDNSNodeGw && !systemNodeGwAddr) {
			log_error("NodeGW misconfigured: one of %s or %s must be set",
			          ENV_SYSTEM_DNS_NODE_GW,
			          ENV_SYSTEM_NODE_GW_ADDR);
			return FALSE;
		}
	}

	globalInitSystemEnable = getenv(ENV_GLOBAL_INIT_ENABLE);
	if (!globalInitSystemEnable) {
		globalInitSystemEnable = GLOBAL_INIT_SYSTEM_DISABLE_STR;
    }

	if (strcmp(globalInitSystemEnable, GLOBAL_INIT_SYSTEM_ENABLE_STR) == 0) {
		globalInitSystemAddr = getenv(ENV_GLOBAL_INIT_SYSTEM_ADDR);
		globalInitSystemPort = getenv(ENV_GLOBAL_INIT_SYSTEM_PORT);
		if (!globalInitSystemAddr || !globalInitSystemPort) {
			log_error("ENV_GLOBAL_INIT_SYSTEM_ADDR and PORT must be set");
			return FALSE;
		}
	}

	dnsServer = getenv(ENV_DNS_SERVER);
	if (dnsServer && strcmp(dnsServer, "true") == 0) {
		log_info("Resolving nameserver from /etc/resolv.conf");
		nameServer = parse_resolveconf(); /* heap */
	}

	systemDNS = getenv(ENV_SYSTEM_DNS);
	if (systemDNS) {
		systemAddr = nslookup(systemDNS, nameServer);
		if (!systemAddr) {
			log_error("Failed to resolve %s=%s", ENV_SYSTEM_DNS, systemDNS);
			goto fail;
		}
		freeSystemAddr = 1;
	} else {
		systemAddr = getenv(ENV_SYSTEM_ADDR);
	}

	if (!systemAddr) {
		log_error("ENV_SYSTEM_DNS or ENV_SYSTEM_ADDR must be set");
		goto fail;
	}

	if (nodegw_any && systemDNSNodeGw) {
		systemNodeGwAddr = nslookup(systemDNSNodeGw, nameServer);
		if (!systemNodeGwAddr) {
			log_error("Failed to resolve %s=%s",
			          ENV_SYSTEM_DNS_NODE_GW, systemDNSNodeGw);
			goto fail;
		}
		freeSystemNodeGwAddr = 1;
	}

	timePeriod = getenv(ENV_DNS_REFRESH_TIME_PERIOD);
	if (!timePeriod)
		period = DEFAULT_TIME_PERIOD;
	else
		period = atoi(timePeriod);

	systemOrg = getenv(ENV_SYSTEM_ORG);
	if (!systemOrg)
		systemOrg = DEFAULT_SYSTEM_ORG;

	apiVersion = getenv(ENV_INIT_SYSTEM_API);
	if (!apiVersion)
		apiVersion = DEFAULT_API_VER;

	*config = calloc(1, sizeof(Config));
	if (!*config) {
		log_error("Memory allocation failure");
		goto fail;
	}

	(*config)->logLevel = getenv(ENV_INIT_CLIENT_LOG_LEVEL);

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

	if (nodegw_any) {
		(*config)->systemNodeGwAddr = strdup(systemNodeGwAddr);
		(*config)->systemNodeGwPort = strdup(systemNodeGwPort);
	} else {
		(*config)->systemNodeGwAddr = strdup("0.0.0.0");
		(*config)->systemNodeGwPort = strdup("0");
	}

	if (nameServer) {
		(*config)->nameServer = strdup(nameServer);
    }

	if (dnsServer) {
		(*config)->dnsServer = strdup(dnsServer);
    }

	(*config)->timePeriod = period;

	(*config)->globalInitSystemEnable =
		(strcmp(globalInitSystemEnable, GLOBAL_INIT_SYSTEM_ENABLE_STR) == 0)
			? GLOBAL_INIT_SYSTEM_ENABLE
			: GLOBAL_INIT_SYSTEM_DISABLE;

	if ((*config)->globalInitSystemEnable) {
		(*config)->globalInitSystemAddr = strdup(globalInitSystemAddr);
		(*config)->globalInitSystemPort = strdup(globalInitSystemPort);
	}

	if (systemDNS) {
		(*config)->systemDNS = strdup(systemDNS);
    }

	if (!(*config)->logLevel) {
		(*config)->logLevel = DEFAULT_LOG_LEVEL;
    }

	if (freeSystemAddr)       free(systemAddr);
	if (freeSystemNodeGwAddr) free(systemNodeGwAddr);
	if (nameServer)           free(nameServer);

	return TRUE;

fail:
	if (freeSystemAddr && systemAddr)             free(systemAddr);
	if (freeSystemNodeGwAddr && systemNodeGwAddr) free(systemNodeGwAddr);
	if (nameServer)                               free(nameServer);
	return FALSE;
}

void clear_config(Config *config) {

	if (config == NULL) return;

	if (config->addr)             		free(config->addr);
	if (config->port)            	 	free(config->port);
	if (config->tempFile)         		free(config->tempFile);
	if (config->systemOrg)        		free(config->systemOrg);
	if (config->systemName)       		free(config->systemName);
	if (config->systemAddr)       		free(config->systemAddr);
	if (config->systemPort)       		free(config->systemPort);
	if (config->systemCert)       		free(config->systemCert);
	if (config->initSystemAddr)   		free(config->initSystemAddr);
	if (config->initSystemPort)   		free(config->initSystemPort);
	if (config->initSystemAPIVer) 		free(config->initSystemAPIVer);
	if (config->globalInitSystemPort)   free(config->globalInitSystemPort);
	if (config->globalInitSystemAddr)   free(config->globalInitSystemAddr);
	if (config->systemDNS)              free(config->systemDNS);
	if (config->dnsServer)              free(config->dnsServer);
	if (config->nameServer)             free(config->nameServer);
    if (config->systemNodeGwAddr)       free(config->systemNodeGwAddr);
    if (config->systemNodeGwPort)       free(config->systemNodeGwPort);

	free(config);
}

char* parse_resolveconf() {
	char * nameServer = NULL;
	FILE *resolv;
	resolv = fopen("/etc/resolv.conf", "r");
	if (resolv) {
		char line[512];	/* "search" is defined to be up to 256 chars */
		while (fgets(line, sizeof(line), resolv)) {
			char *p, *arg;
			p = strtok(line, " \t\n");
			if (!p)
				continue;
			log_debug("resolv_key:'%s'\n", p);
			arg = strtok(NULL, "\n");
			log_debug("resolv_arg:'%s'\n", arg);
			if (!arg)
				continue;
			/* May be parse them if required. Skipping for now */
			if (strcmp(p, "domain") == 0) {
				continue;
			}
			if (strcmp(p, "search") == 0) {
				continue;
			}

			if (strcmp(p, "nameserver") != 0)
				continue;
			/* only first for now nameserver DNS. We can have upto three in file*/
			nameServer =  strdup(arg);
			break;
		}
		fclose(resolv);
	}
	return nameServer;
}
