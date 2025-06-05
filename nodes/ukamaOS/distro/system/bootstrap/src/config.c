/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <stdio.h>

#include "usys_log.h"

#include "config.h"
#include "toml.h"

int parse_config(Config *config, toml_table_t *configData) {

	toml_datum_t nodedHost;
	toml_datum_t meshConfig, remoteIPFile;

	/* sanity check */
	if (config == NULL) return FALSE;
	if (configData == NULL) return FALSE;

	/* Read the config data from the config.toml and load into Config. */
	/* noded-host */
	nodedHost = toml_string_in(configData, NODED_HOST);
	if (!nodedHost.ok) {
		usys_log_debug("[%s] is missing, setting to default: %s", NODED_HOST,
				  DEF_NODED_HOST);
		config->nodedHost = strdup(DEF_NODED_HOST);
	} else {
		config->nodedHost = strdup(nodedHost.u.s);
	}

	/* mesh-config */
	meshConfig = toml_string_in(configData, MESH_CONFIG);
	if (!meshConfig.ok) {
		usys_log_debug("[%s] is missing, setting to default: %s", MESH_CONFIG,
				  DEF_MESH_CONFIG);
		config->meshConfig = strdup(DEF_MESH_CONFIG);
	} else {
		config->meshConfig = strdup(meshConfig.u.s);
	}

	/* remote-IP-file */
	remoteIPFile = toml_string_in(configData, REMOTE_IP_FILE);
	if (!remoteIPFile.ok) {
		usys_log_debug("[%s] is missing, setting to default: %s", REMOTE_IP_FILE,
				  DEF_REMOTE_IP_FILE);
		config->remoteIPFile = strdup(DEF_REMOTE_IP_FILE);
	} else {
		config->remoteIPFile = strdup(remoteIPFile.u.s);
	}

	if (nodedHost.ok)       free(nodedHost.u.s);
	if (meshConfig.ok)      free(meshConfig.u.s);
	if (remoteIPFile.ok)    free(remoteIPFile.u.s);

    return TRUE;
}

bool read_bootstrap_server_info(char **buffer) {

    FILE *file = NULL;
    long length = 0;

    if ((file = fopen(DEF_BOOTSTRAP_FILE, "r")) == NULL) {
        usys_log_error("Error opening bootstrap file: %s %s",
                  DEF_BOOTSTRAP_FILE, strerror(errno));
        return FALSE;
    }

    fseek(file, 0, SEEK_END);
    length = ftell(file);
    rewind(file);

    *buffer = (char *)calloc((length + 1), sizeof(char));
    if (*buffer == NULL) {
        usys_log_error("Memory allocation failed: %s",
                  (length + 1) * sizeof(char));
        fclose(file);
        return FALSE;
    }

    fread(*buffer, sizeof(char), length-1, file);
    (*buffer)[length] = '\0';

    fclose(file);

    return TRUE;
}

int process_config_file(char *fileName, Config *config) {

	FILE *fp;
	toml_table_t *fileData, *configData;
	char errBuf[MAX_BUFFER];

	if ((fp = fopen(fileName, "r")) == NULL) {
		usys_log_error("Error opening config file: %s: %s", fileName,
				  strerror(errno));
		return FALSE;
	}

	/* Prase the TOML file entries. */
	fileData = toml_parse_file(fp, errBuf, sizeof(errBuf));
	fclose(fp);
	if (!fileData) {
		usys_log_error("Error parsing the config file %s: %s", fileName, errBuf);
		return FALSE;
	}

	/* Parse the config. */
	configData = toml_table_in(fileData, CONFIG);

	if (configData == NULL) {
		usys_log_error("[Config] section parsing error in file: %s", fileName);
		toml_free(fileData);
		return FALSE;
	}

	if (!parse_config(config, configData)) {
        usys_log_error("Unable to prase config file: %s", fileName);
        toml_free(fileData);
        return FALSE;
    }

    if (!read_bootstrap_server_info(&config->bootstrapRemoteServer)) {
        usys_log_debug("Unable to read bootstrap server info");
        config->bootstrapRemoteServer = strdup(DEF_BOOTSTRAP_SERVER);
    }

	toml_free(fileData);
	return TRUE;
}

void print_config(Config *config) {

	if (config == NULL) return;

	if (config->nodedHost) {
		usys_log_debug("noded host: %s", config->nodedHost);
	}

    usys_log_debug("noded port: %d", config->nodedPort);

	if (config->meshConfig) {
		usys_log_debug("mesh config: %s", config->meshConfig);
	}

	if (config->remoteIPFile) {
		usys_log_debug("remote IP file: %s", config->remoteIPFile);
	}

	if (config->bootstrapRemoteServer) {
	    usys_log_debug("bootstrap remote server: %s", config->bootstrapRemoteServer);
	}

    if (config->bootstrapRemotePort) {
        usys_log_debug("bootstrap remote port: %d", config->bootstrapRemotePort);
    }
}

void clear_config(Config *config) {

	if (!config) return;

	free(config->nodedHost);
	free(config->meshConfig);
	free(config->remoteIPFile);
	free(config->bootstrapRemoteServer);
}
