/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <stdio.h>

#include "usys_api.h"
#include "usys_error.h"
#include "usys_log.h"
#include "usys_file.h"
#include "usys_log.h"
#include "usys_mem.h"

#include "mesh.h"
#include "config.h"
#include "toml.h"

#include "static.h"

STATIC int parse_config_entries(Config *config, toml_table_t *configData);
STATIC int read_line(char *buffer, int size, FILE *fp);

void print_config(Config *config) {

	usys_log_debug("Remote connect port: %s", config->remoteConnect);
    usys_log_debug("Org ename: %s",           config->orgName);
	usys_log_debug("Forward port: %d",        config->forwardPort);
    usys_log_debug("Servuce port: %d",        config->servicePort);
    usys_log_debug("Local hostname: %s",      config->localHostname);
    usys_log_debug("TLS/SSL key file: %s",    config->keyFile);
    usys_log_debug("TLS/SSL cert file: %s",   config->certFile);
}

STATIC int read_line(char *buffer, int size, FILE *fp) {

	char *tmp;

	memset(buffer, 0, size);

	if (fgets(buffer, size, fp) == NULL) {
		*buffer = '\0';
		return FALSE;
	} else {
		/* remove newline */
		if ((tmp = strrchr(buffer, '\n')) != NULL) {
			*tmp = '\0';
		}
	}
	return TRUE;
}

void split_strings(char *input, char **str1, char **str2,
                   char *delimiter) {

    char *token=NULL;

    token = strtok(input, delimiter);

    if (token != NULL) {
        *str1 = strdup(token);

        token = strtok(NULL, delimiter);
        if (token != NULL) {
            *str2 = strdup(token);
        }
    }
}

STATIC int read_nodeid(char **nodeID) {

    char buffer[MAX_BUFFER] = {0};
    FILE *fp=NULL;

#ifdef UNIT_TEST
    fp = fopen(".nodeid", "r");
#else
    fp = fopen("/ukama/nodeid", "r");
#endif
    if (fp == NULL) {
        usys_log_error("Unable to open /ukama/nodeid file. Error: %s",
                       strerror(errno));
        return FALSE;
    }

    if (read_line(buffer, MAX_BUFFER, fp) <= 0) {
        usys_log_error("[%s] Error reading file. Error: %s", "/ukama/nodeid",
                       strerror(errno));
        return FALSE;
	} else {
        *nodeID = strdup(buffer);
    }

    return TRUE;
}

STATIC int read_org_name(char **orgName) {

    char buffer[MAX_BUFFER] = {0};
    FILE *fp=NULL;

#ifdef UNIT_TEST
    fp = fopen(".orgName", "r");
#else
    fp = fopen("/ukama/org", "r");
#endif
    if (fp == NULL) {
        usys_log_error("Unable to open /ukama/org file. Error: %s",
                       strerror(errno));
        return FALSE;
    }

    if (read_line(buffer, MAX_BUFFER, fp) <= 0) {
        usys_log_error("[%s] Error reading file. Error: %s", "/ukama/org",
                       strerror(errno));
        return FALSE;
	} else {
        *orgName = strdup(buffer);
    }

    return TRUE;
}

STATIC int read_hostname_and_nodeid(char *fileName,
                                    char **hostname,
                                    char **subnetMask) {

    int ret=TRUE;
	FILE *fp=NULL;
	char *buffer=NULL, *CIDR=NULL, *nodeID=NULL;

	buffer = (char *)malloc(MAX_BUFFER);
	if (!buffer) {
		usys_log_error("Error allocating memory of size: %s", MAX_BUFFER);
		return FALSE;
	}

	fp = fopen(fileName, "r");
	if (fp == NULL) {
		usys_log_error("[%s] Error opening file. Error: %s", fileName,
				  strerror(errno));
		return FALSE;
	}

	/* Read the file content. */
	if (read_line(buffer, MAX_BUFFER, fp)<=0) {
		usys_log_error("[%s] Error reading file. Error: %s", fileName,
				  strerror(errno));
        ret = FALSE;
	} else {
        split_strings(buffer, &CIDR, &nodeID, ";");
        split_strings(CIDR, hostname, subnetMask, "/");
    }

	fclose(fp);
    free(buffer);
    free(CIDR);
    free(nodeID);

	return ret;
}

STATIC int parse_config_entries(Config *config, toml_table_t *configData) {

	int ret=TRUE;
    int remote=0;
	char *hostname=NULL, *subnetMask=NULL;
	toml_datum_t cert, key, localHostname, remoteIPFile;

	remoteIPFile  = toml_string_in(configData, REMOTE_IP_FILE);
    localHostname = toml_string_in(configData, LOCAL_HOSTNAME);
	cert          = toml_string_in(configData, CERT);
	key           = toml_string_in(configData, KEY);

	if (!remoteIPFile.ok) {
		usys_log_error("[%s] is missing but is mandatory", REMOTE_IP_FILE);
        ret=FALSE;
        goto done;
	} else {
		/* Read the content of the IP file. */
		if (read_hostname_and_nodeid(remoteIPFile.u.s,
                                     &hostname,
                                     &subnetMask) == FALSE) {
			goto done;
		}
	}

    remote = usys_find_service_port(SERVICE_REMOTE);
    if (remote == 0) {
        usys_log_error("Error getting remote mesh.d port from service db");
        ret = FALSE;
        goto done;
    }

	config->remoteConnect = (char *)calloc(1, MAX_BUFFER);
    sprintf(config->remoteConnect, "ws://%s:%d/%s", hostname,
            remote, PREFIX_WEBSOCKET);

	config->deviceInfo = (DeviceInfo *)malloc(sizeof(DeviceInfo));
	if (config->deviceInfo == NULL) {
		usys_log_error("Error allocating memory of size: %d", sizeof(DeviceInfo));
		goto done;
	}

    if (!read_nodeid(&config->deviceInfo->nodeID)) {
        usys_log_error("Unable to read nodeID from /ukama/nodeid");
        goto done;
    }

    if (!read_org_name(&config->orgName)) {
        usys_log_error("Unable to read orgName from /ukama/org");
        goto done;
    }

	if (!localHostname.ok) {
		usys_log_debug("[%s] is missing, setting to default: %s", LOCAL_HOSTNAME,
				  DEFAULT_LOCAL_HOSTNAME);
		config->localHostname = strdup(DEFAULT_LOCAL_HOSTNAME);
	} else {
		config->localHostname = strdup(localHostname.u.s);
	}

	if (cert.ok) {
		config->certFile = strdup(cert.u.s);
	} else {
		config->certFile = strdup(DEFAULT_CERT);
	}

	if (key.ok) {
		config->keyFile = strdup(key.u.s);
	} else {
		config->keyFile = strdup(DEFAULT_KEY);
	}

 done:
	/* clear up toml allocations. */
	if (key.ok)           free(key.u.s);
	if (cert.ok)          free(cert.u.s);
	if (remoteIPFile.ok)  free(remoteIPFile.u.s);
    if (hostname)         free(hostname);
    if (subnetMask)       free(subnetMask);

	return ret;
}

int process_config_file(Config *config, char *fileName) {

	int ret=TRUE;
	FILE *fp;
	toml_table_t *fileData=NULL, *localConfig=NULL;
	char errBuf[MAX_BUFFER];

	/* Sanity check. */
	if (fileName == NULL || config == NULL)
		return FALSE;

	if ((fp = fopen(fileName, "r")) == NULL) {
		usys_log_error("Error opening config file: %s: %s\n", fileName,
				  strerror(errno));
		return FALSE;
	}

	/* Parse the TOML file entries. */
	fileData = toml_parse_file(fp, errBuf, sizeof(errBuf));
  	fclose(fp);
 	if (!fileData) {
		usys_log_error("Error parsing the config file %s: %s\n", fileName, errBuf);
		return FALSE;
	}

	localConfig = toml_table_in(fileData, LOCAL_CONFIG);
	if (localConfig == NULL) {
		usys_log_error("[%s] section parsing error in file: %s\n", LOCAL_CONFIG,
				  fileName);
		ret = FALSE;
		goto done;
	}
	ret = parse_config_entries(config, localConfig);
	if (ret == FALSE) {
		goto done;
	}

done:
	toml_free(fileData);
	return ret;
}

void clear_config(Config *config) {

	if (!config) return;

	usys_free(config->remoteConnect);
    usys_free(config->localHostname);
	usys_free(config->certFile);
	usys_free(config->keyFile);
    usys_free(config->orgName);
}
