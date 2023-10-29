/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

/*
 * capp config.json
 */

#ifndef CAPP_CONFIG_H
#define CAPP_CONFIG_H

#include "lxce_config.h"

/* for parsing config.json */
#define JSON_VERSION    "version"
#define JSON_SERIAL     "serial"
#define JSON_TARGET     "target"
#define JSON_ALL        "all"
#define JSON_PROCESS    "process"
#define JSON_EXEC       "exec"
#define JSON_ARGS       "args"
#define JSON_ENV        "env"
#define JSON_TYPE       "type"
#define JSON_HOSTNAME   "hostname"
#define JSON_NAMESPACES "namespaces"

/* some defaults */
#define CAPP_DEFAULT_HOSTNAME "localhost"

typedef struct capp_process_t {

  char *exec;  /* Executable name */
  char *argv; /* Arguments to the executable */
  char *env;  /* Environment variables setup for executable */
} CAppProc;

/* Store capp config.json values */
typedef struct capp_config_t {

  char *version;      /* capp version */
  char *serial;       /* serial of device, if applicable */
  int  target;        /* Target of this contained space (serial or general) */

  char *hostName;     /* host name associated with space */
  int nameSpaces;     /* linux namespaces enabled in this space */

  CAppProc *process;  /* Info related to the process/executable */
} CAppConfig;

int valid_path(char *path);
int process_capp_config_file(CAppConfig *config, char *fileName);

#endif /* CAPP_CONFIG_H */
