/**
 * Copyright (c) 2021-present, Ukama Inc.
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
  toml_datum_t localAccept, localEP, wimcHost, wimcPort, meshPort;

  /* sanity check */
  if (config == NULL) return FALSE;
  if (configData == NULL) return FALSE;

  /* Read the config data from the config.toml and load into Config. */
  /* local-accept */
  localAccept = toml_string_in(configData, LOCAL_ACCEPT);
  if (!localAccept.ok) {
    log_debug("[%s] is missing, setting to default: %s", LOCAL_ACCEPT,
	      DEF_LOCAL_ACCEPT);
    config->localAccept = strdup(DEF_LOCAL_ACCEPT);
  } else {
    config->localAccept = strdup(localAccept.u.s);
  }

  /* local-ep */
  localEP = toml_string_in(configData, LOCAL_EP);
  if (!localEP.ok) {
    log_debug("[%s] is missing, setting to default: %s", LOCAL_EP,
	      DEF_LOCAL_EP);
    config->localEP = strdup(DEF_LOCAL_EP);
  } else {
    config->localEP = strdup(localEP.u.s);
  }

  /* wimc-host */
  wimcHost = toml_string_in(configData, WIMC_HOST);
  if (!wimcHost.ok) {
    log_debug("[%s] is missing, setting to default: %s", WIMC_HOST,
	      DEF_WIMC_HOST);
    config->wimcHost = strdup(DEF_WIMC_HOST);
  } else {
    config->wimcHost = strdup(wimcHost.u.s);
  }

  /* wimc-port */
  wimcPort = toml_string_in(configData, WIMC_PORT);
  if (!wimcPort.ok) {
    log_debug("[%s] is missing, setting to default: %s", WIMC_PORT,
	      DEF_WIMC_PORT);
    config->wimcPort = strdup(DEF_WIMC_PORT);
  } else {
    config->wimcPort = strdup(wimcPort.u.s);
  }

  /* mesh-host */
  meshPort = toml_string_in(configData, MESH_PORT);
  if (!meshPort.ok) {
    log_debug("[%s] is missing, setting to default: %s", MESH_PORT,
	      DEF_MESH_PORT);
    config->meshPort = strdup(DEF_MESH_PORT);
  } else {
    config->meshPort = strdup(meshPort.u.s);
  }

  if (localAccept.ok) free(localAccept.u.s);
  if (localEP.ok)     free(localEP.u.s);
  if (wimcHost.ok)    free(wimcHost.u.s);
  if (wimcPort.ok)    free(wimcPort.u.s);
  if (meshPort.ok)    free(meshPort.u.s);

  return ret;
}

/*
 * process_config_file -- read and parse the config file.
 *
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

  if (config->localAccept) {
    log_debug("Local-Accept Port: %s", config->localAccept);
  }

  if (config->localEP) {
    log_debug("Local-EP: %s", config->localEP);
  }

  if (config->wimcHost) {
    log_debug("wimcHost: %s", config->wimcHost);
  }

  if (config->wimcPort) {
    log_debug("wimcPort: %s", config->wimcPort);
  }

  if (config->meshPort) {
    log_debug("meshPort: %s", config->meshPort);
  }
}

/*
 * clear_config --
 *
 */
void clear_config(Config *config) {

  if (!config) return;

  free(config->localAccept);
  free(config->localEP);
  free(config->wimcHost);
  free(config->wimcPort);
  free(config->meshPort);
}
