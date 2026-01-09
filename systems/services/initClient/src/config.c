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
	char *systemDNS = NULL, *timePeriod = NULL, *dnsServer = NULL, *nameServer = NULL;
    char *systemNodeGwAddr = NULL, *systemNodeGwPort = NULL;
	int period = 0 ;

	if ((addr = getenv(ENV_INIT_CLIENT_ADDR)) == NULL ||
		(port = getenv(ENV_INIT_CLIENT_PORT)) == NULL ||
		(tempFile = getenv(ENV_INIT_CLIENT_TEMP_FILE)) == NULL) {
		log_error("Required env variables: %s %s %s missing",
				  ENV_INIT_CLIENT_ADDR,
				  ENV_INIT_CLIENT_PORT,
				  ENV_INIT_CLIENT_TEMP_FILE);
		return FALSE;
	}

	if ((systemName = getenv(ENV_SYSTEM_NAME)) == NULL ||
			(systemPort = getenv(ENV_SYSTEM_PORT)) == NULL ||
			(systemCert = getenv(ENV_SYSTEM_CERT)) == NULL ||
			(initSystemAddr = getenv(ENV_INIT_SYSTEM_ADDR)) == NULL ||
			(initSystemPort = getenv(ENV_INIT_SYSTEM_PORT)) == NULL ){
		log_error("Required env variables not defined");
		return FALSE;
	}

    systemNodeGwAddr = getenv(ENV_SYSTEM_NODE_GW_ADDR);
    systemNodeGwPort = getenv(ENV_SYSTEM_NODE_GW_PORT);
    if ((systemNodeGwAddr && !systemNodeGwPort) ||
        (!systemNodeGwAddr && systemNodeGwPort)) {
        log_error("Error: Both %s and %s must be set together\n",
                  ENV_SYSTEM_NODE_GW_ADDR,
                  ENV_SYSTEM_NODE_GW_PORT);
        return FALSE;
    }

	if ((globalInitSystemEnable = getenv(ENV_GLOBAL_INIT_ENABLE)) == NULL) {
		globalInitSystemEnable = GLOBAL_INIT_SYSTEM_DISABLE_STR;
	}

	if (strcmp(globalInitSystemEnable, GLOBAL_INIT_SYSTEM_ENABLE_STR) == 0){
		if ((globalInitSystemAddr = getenv(ENV_GLOBAL_INIT_SYSTEM_ADDR)) == NULL ||
				(globalInitSystemPort = getenv(ENV_GLOBAL_INIT_SYSTEM_PORT)) == NULL ){
			log_error("Required env variables system ENV_GLOBAL_INIT_SYSTEM_ADDR and ENV_GLOBAL_INIT_SYSTEM_PORT not defined");
		}
	}

	if ((dnsServer = getenv(ENV_DNS_SERVER)) != NULL) {
		if (strcmp(dnsServer, "true") == 0) {
            /* Fetching from /etc/resolv.conf */
		    log_info("Resolving nameserver from  /etc/resolv.conf");
		    nameServer = parse_resolveconf();
		} else {
			nameServer = NULL;
		}
	}

	if ((systemDNS = getenv(ENV_SYSTEM_DNS)) != NULL) {
		if (nameServer == NULL) {
			systemAddr       = nslookup(systemDNS, NULL);
            systemNodeGwAddr = nslookup(systemDNS, NULL);
		} else {
			systemAddr       = nslookup(systemDNS, nameServer);
            systemNodeGwAddr = nslookup(systemDNS, nameServer);
		}
	} else {
		systemAddr       = getenv(ENV_SYSTEM_ADDR);
        systemNodeGwAddr = getenv(ENV_SYSTEM_NODE_GW_ADDR);
	}

	if (!systemAddr) {
		log_error("Required one of env variable ENV_SYSTEM_DNS or ENV_SYSTEM_ADDR to be valid");
		return FALSE;
	}

	if ((timePeriod = getenv(ENV_DNS_REFRESH_TIME_PERIOD)) == NULL) {
			period = DEFAULT_TIME_PERIOD;
	}

	if (timePeriod) {
		period = atoi(timePeriod);
	}
	if ((systemOrg = getenv(ENV_SYSTEM_ORG)) == NULL) {
		systemOrg = DEFAULT_SYSTEM_ORG;
	}

	if ((apiVersion = getenv(ENV_INIT_SYSTEM_API)) == NULL) {
		apiVersion = DEFAULT_API_VER;
	}

	*config = (Config *)calloc(1, sizeof(Config));
	if (*config == NULL) {
		log_error("Memory allocation failure: %d", sizeof(Config));
		return FALSE;
	}

	(*config)->logLevel   = getenv(ENV_INIT_CLIENT_LOG_LEVEL);

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
    /* set 0.0.0.0:0 for system's node gw */
    if (systemNodeGwAddr == NULL && systemNodeGwPort == NULL) {
           (*config)->systemNodeGwAddr = strdup("0.0.0.0");
           (*config)->systemNodeGwPort = strdup("0");
    } else {
        if (systemNodeGwAddr) {
            (*config)->systemNodeGwAddr = strdup(systemNodeGwAddr);
        }
        if (systemNodeGwPort) {
            (*config)->systemNodeGwPort = strdup(systemNodeGwPort);
        } else {
            (*config)->systemNodeGwPort = strdup("0");
        }
    }

	if (nameServer) {
    	(*config)->nameServer = strdup(nameServer);
	}
	
	(*config)->timePeriod = period;
	
	if (dnsServer) {
		(*config)->dnsServer = dnsServer;
	}
	(*config)->globalInitSystemEnable = (strcmp(globalInitSystemEnable, GLOBAL_INIT_SYSTEM_ENABLE_STR) == 0) ? GLOBAL_INIT_SYSTEM_ENABLE : GLOBAL_INIT_SYSTEM_DISABLE ;

	if ((*config)->globalInitSystemEnable) {
		(*config)->globalInitSystemAddr   = strdup(globalInitSystemAddr);
		(*config)->globalInitSystemPort   = strdup(globalInitSystemPort);
	}

	if(systemDNS) {
		(*config)->systemDNS = strdup(systemDNS);
	}

	if (!(*config)->logLevel) {
		log_debug("Log level not defined, setting to default: DEBUG");
		(*config)->logLevel = DEFAULT_LOG_LEVEL;
	}

	return TRUE;
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
