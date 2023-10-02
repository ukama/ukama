/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Config.c
 *
 */
#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <stdio.h>

#include "mesh.h"
#include "config.h"
#include "toml.h"
#include "log.h"

static int parse_config_entries(Config *config, toml_table_t *configData);
static int read_line(char *buffer, int size, FILE *fp);

/*
 * print_config --
 *
 */
void print_config(Config *config) {

	log_debug("Remote connect port: %s", config->remoteConnect);
	log_debug("Local accept port: %s", config->localAccept);
    log_debug("Local hostname: %s", config->localHostname);
    log_debug("TLS/SSL key file: %s", config->keyFile);
    log_debug("TLS/SSL cert file: %s", config->certFile);
}

/*
 * read_line -- read a line from file pointer.
 *
 */
static int read_line(char *buffer, int size, FILE *fp) {

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

/*
 * split_strings --
 *
 */
static void split_strings(char *input, char **str1, char **str2,
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

/*
 * read_hostname_and_nodeid -- read hostname (ip:port) and nodeID from the
 *                             passed file
 *
 */
static int read_hostname_and_nodeid(char *fileName, char **hostname,
                                    char **subnetMask, char **nodeID) {

    int ret=TRUE;
	FILE *fp=NULL;
	char *buffer=NULL, *CIDR=NULL;

	buffer = (char *)malloc(MAX_BUFFER);
	if (!buffer) {
		log_error("Error allocating memory of size: %s", MAX_BUFFER);
		return FALSE;
	}

	fp = fopen(fileName, "r");
	if (fp == NULL) {
		log_error("[%s] Error opening file. Error: %s", fileName,
				  strerror(errno));
		return FALSE;
	}

	/* Read the file content. */
	if (read_line(buffer, MAX_BUFFER, fp)<=0) {
		log_error("[%s] Error reading file. Error: %s", fileName,
				  strerror(errno));
        ret = FALSE;
	} else {
        split_strings(buffer, &CIDR, nodeID, ";");
        split_strings(CIDR, hostname, subnetMask, "/");
    }

	fclose(fp);
    free(buffer);
    free(CIDR);

	return ret;
}

/*
 * parse_config_entries -- Server/client stuff.
 *
 */
static int parse_config_entries(Config *config, toml_table_t *configData) {

	int ret=TRUE;
	char *hostname=NULL, *nodeID=NULL, *subnetMask=NULL;
	toml_datum_t localAccept, cert, key, localHostname, remoteIPFile;

	remoteIPFile  = toml_string_in(configData, REMOTE_IP_FILE);
	localAccept   = toml_string_in(configData, LOCAL_ACCEPT);
    localHostname = toml_string_in(configData, LOCAL_HOSTNAME);
	cert          = toml_string_in(configData, CERT);
	key           = toml_string_in(configData, KEY);

	if (!remoteIPFile.ok) {
		log_error("[%s] is missing but is mandatory", REMOTE_IP_FILE);
        ret=FALSE;
        goto done;
	} else {
		/* Read the content of the IP file. */
		if (read_hostname_and_nodeid(remoteIPFile.u.s, &hostname,
                                     &subnetMask, &nodeID) == FALSE) {
			goto done;
		}
	}

	config->remoteConnect = (char *)calloc(1, MAX_BUFFER);
    sprintf(config->remoteConnect, "ws://%s:%s/%s", hostname,
            DEFAULT_REMOTE_PORT, PREFIX_WEBSOCKET);

	config->deviceInfo = (DeviceInfo *)malloc(sizeof(DeviceInfo));
	if (config->deviceInfo == NULL) {
		log_error("Error allocating memory of size: %d", sizeof(DeviceInfo));
		goto done;
	}
    config->deviceInfo->nodeID = strdup(nodeID);

	if (!localAccept.ok) {
		log_debug("[%s] is missing, setting to default: %s", LOCAL_ACCEPT,
				  DEFAULT_LOCAL_ACCEPT);
		config->localAccept = strdup(DEFAULT_LOCAL_ACCEPT);
	} else {
		config->localAccept = strdup(localAccept.u.s);
	}

	if (!localHostname.ok) {
		log_debug("[%s] is missing, setting to default: %s", LOCAL_HOSTNAME,
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
	if (localAccept.ok)   free(localAccept.u.s);
	if (remoteIPFile.ok)  free(remoteIPFile.u.s);
    if (hostname)         free(hostname);
    if (subnetMask)       free(subnetMask);
    if (nodeID)           free(nodeID);

	return ret;
}

/*
 * process_config_file -- read and parse the config file. 
 *                       
 *
 */
int process_config_file(Config *config, char *fileName) {

	int ret=TRUE;
	FILE *fp;
	toml_table_t *fileData=NULL, *localConfig=NULL;
	char errBuf[MAX_BUFFER];

	/* Sanity check. */
	if (fileName == NULL || config == NULL)
		return FALSE;

	if ((fp = fopen(fileName, "r")) == NULL) {
		log_error("Error opening config file: %s: %s\n", fileName,
				  strerror(errno));
		return FALSE;
	}

	/* Parse the TOML file entries. */
	fileData = toml_parse_file(fp, errBuf, sizeof(errBuf));
  	fclose(fp);
 	if (!fileData) {
		log_error("Error parsing the config file %s: %s\n", fileName, errBuf);
		return FALSE;
	}

	localConfig = toml_table_in(fileData, LOCAL_CONFIG);
	if (localConfig == NULL) {
		log_error("[%s] section parsing error in file: %s\n", LOCAL_CONFIG,
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

/*
 * clear_config --
 */
void clear_config(Config *config) {

	if (!config) return;

	free(config->remoteConnect);
	free(config->localAccept);
    free(config->localHostname);
	free(config->certFile);
	free(config->keyFile);
}
