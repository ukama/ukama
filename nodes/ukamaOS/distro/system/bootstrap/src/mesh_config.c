/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * mesh_config.c
 *
 */

#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <stdio.h>

#include "mesh_config.h"
#include "toml.h"
#include "log.h"

/*
 * read_mesh_config_file -- read and parse the mesh config file for
 *                          IP and cert/key filename
 *
 */
int read_mesh_config_file(char *fileName, MeshConfig *meshConfig) {

	int ret=FALSE;
	FILE *fp=NULL;
	toml_table_t *fileData=NULL;
	toml_table_t *clientConfig=NULL;
  	toml_datum_t cert, key, remoteIPFile;

	char errBuf[MAX_BUFFER];

	/* Sanity check. */
	if (fileName == NULL || meshConfig == NULL) return FALSE;

	if ((fp = fopen(fileName, "r")) == NULL) {
		log_error("Error opening config file: %s: %s\n", fileName,
				  strerror(errno));
		return ret;
	}

	/* Parse the TOML file entries. */
	fileData = toml_parse_file(fp, errBuf, sizeof(errBuf));
  	fclose(fp);

	if (!fileData) {
		log_error("Error parsing the config file %s: %s\n", fileName, errBuf);
		return ret;
	}

	/* client-mode entries only */
	clientConfig = toml_table_in(fileData, CLIENT_CONFIG);
	if (clientConfig == NULL) {
		log_error("[%s] section parsing error in file: %s\n", CLIENT_CONFIG,
				  fileName);
		goto done;
	}

	remoteIPFile  = toml_string_in(clientConfig, REMOTE_IP_FILE);
	cert          = toml_string_in(clientConfig, CFG_CERT);
	key           = toml_string_in(clientConfig, CFG_KEY);

	if (!remoteIPFile.ok) {
		log_error("[%s] is missing from %s", REMOTE_IP_FILE, fileName);
		goto done;
	}

	if (!cert.ok || !key.ok) {
		log_error("[%s] or [%s] is missing from %s", CFG_CERT, CFG_KEY,
				  fileName);
		goto done;
	}

	meshConfig->remoteIPFile = strdup(remoteIPFile.u.s);
	meshConfig->certFile     = strdup(cert.u.s);
	meshConfig->keyFile      = strdup(key.u.s);

	ret = TRUE;

 done:
	toml_free(fileData);

	if (remoteIPFile.u.s) free(remoteIPFile.u.s);
	if (cert.u.s)         free(cert.u.s);
	if (key.u.s)          free(key.u.s);

	return ret;
}

/*
 * clear_mesh_config --
 */
void clear_mesh_config(MeshConfig *config) {

	if (!config) return;

	if (config->remoteIPFile) free(config->remoteIPFile);
	if (config->certFile)     free(config->certFile);
	if (config->keyFile)      free(config->keyFile);
}
