/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Functions related to config
 */

#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <stdio.h>

#include "config.h"
#include "toml.h"
#include "log.h"

/*
 * parse_config -- process [config] stuff
 *
 */
int parse_config(Config *config, toml_table_t *configData) {

	int ret=FALSE;
	toml_datum_t nodedHost, nodedPort;
	toml_datum_t meshConfig, remoteIPFile;
	toml_datum_t bootstrapServer;

	/* sanity check */
	if (config == NULL) return FALSE;
	if (configData == NULL) return FALSE;

	/* Read the config data from the config.toml and load into Config. */
	/* noded-host */
	nodedHost = toml_string_in(configData, NODED_HOST);
	if (!nodedHost.ok) {
		log_debug("[%s] is missing, setting to default: %s", NODED_HOST,
				  DEF_NODED_HOST);
		config->nodedHost = strdup(DEF_NODED_HOST);
	} else {
		config->nodedHost = strdup(nodedHost.u.s);
	}

	/* noded-port */
	nodedPort = toml_string_in(configData, NODED_PORT);
	if (!nodedPort.ok) {
		log_debug("[%s] is missing, setting to default: %s", NODED_PORT,
				  DEF_NODED_PORT);
		config->nodedPort = strdup(DEF_NODED_PORT);
	} else {
		config->nodedPort = strdup(nodedPort.u.s);
	}

	/* mesh-config */
	meshConfig = toml_string_in(configData, MESH_CONFIG);
	if (!meshConfig.ok) {
		log_debug("[%s] is missing, setting to default: %s", MESH_CONFIG,
				  DEF_MESH_CONFIG);
		config->meshConfig = strdup(DEF_MESH_CONFIG);
	} else {
		config->meshConfig = strdup(meshConfig.u.s);
	}

	/* remote-IP-file */
	remoteIPFile = toml_string_in(configData, REMOTE_IP_FILE);
	if (!remoteIPFile.ok) {
		log_debug("[%s] is missing, setting to default: %s", REMOTE_IP_FILE,
				  DEF_REMOTE_IP_FILE);
		config->remoteIPFile = strdup(DEF_REMOTE_IP_FILE);
	} else {
		config->remoteIPFile = strdup(remoteIPFile.u.s);
	}

	/* bootstrap-server */
	bootstrapServer = toml_string_in(configData, BOOTSTRAP_SERVER);
	if (!bootstrapServer.ok) {
		log_debug("[%s] is missing, setting to default: %s", BOOTSTRAP_SERVER,
				  DEF_BOOTSTRAP_SERVER);
		config->bootstrapServer = strdup(DEF_BOOTSTRAP_SERVER);
	} else {
		config->bootstrapServer = strdup(bootstrapServer.u.s);
	}

	if (nodedHost.ok)       free(nodedHost.u.s);
	if (nodedPort.ok)       free(nodedPort.u.s);
	if (meshConfig.ok)      free(meshConfig.u.s);
	if (remoteIPFile.ok)    free(remoteIPFile.u.s);
	if (bootstrapServer.ok) free(bootstrapServer.u.s);

	return ret;
}

/*
 * process_config_file -- read and parse the config file
 *
 */
int process_config_file(char *fileName, Config *config) {

	FILE *fp;
	toml_table_t *fileData, *configData;
	char errBuf[MAX_BUFFER];

	if ((fp = fopen(fileName, "r")) == NULL) {
		log_error("Error opening config file: %s: %s\n", fileName,
				  strerror(errno));
		return FALSE;
	}

	/* Prase the TOML file entries. */
	fileData = toml_parse_file(fp, errBuf, sizeof(errBuf));
	fclose(fp);
	if (!fileData) {
		log_error("Error parsing the config file %s: %s\n", fileName, errBuf);
		return FALSE;
	}

	/* Parse the config. */
	configData = toml_table_in(fileData, CONFIG);

	if (configData == NULL) {
		log_error("[Config] section parsing error in file: %s\n", fileName);
		toml_free(fileData);
		return FALSE;
	}

	parse_config(config, configData);

	toml_free(fileData);
	return TRUE;
}

/*
 * print_config -- print the config
 *
 */
void print_config(Config *config) {

	if (config == NULL) return;

	if (config->nodedHost) {
		log_debug("noded host: %s", config->nodedHost);
	}

	if (config->nodedPort) {
		log_debug("noded port: %s", config->nodedPort);
	}

	if (config->meshConfig) {
		log_debug("mesh config: %s", config->meshConfig);
	}

	if (config->remoteIPFile) {
		log_debug("remote IP file: %s", config->remoteIPFile);
	}

	if (config->bootstrapServer) {
	    log_debug("bootstrap server: %s", config->bootstrapServer);
	}
}

/*
 * clear_config --
 *
 */
void clear_config(Config *config) {

	if (!config) return;

	free(config->nodedHost);
	free(config->nodedPort);
	free(config->meshConfig);
	free(config->remoteIPFile);
	free(config->bootstrapServer);
}
