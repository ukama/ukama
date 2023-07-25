/**
 * Copyright (c) 2023-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef CAPP_CONFIG_H
#define CAPP_CONFIG_H

#include "usys_types.h"

/* for parsing config.json of capp */
#define JTAG_VERSION    "version"
#define JTAG_SERIAL     "serial"
#define JTAG_TARGET     "target"
#define JTAG_ALL        "all"
#define JTAG_PROCESS    "process"
#define JTAG_EXEC       "exec"
#define JTAG_ARGS       "args"
#define JTAG_ENV        "env"
#define JTAG_TYPE       "type"
#define JTAG_HOSTNAME   "hostname"
#define JTAG_NAMESPACES "namespaces"

/* some defaults */
#define CAPP_DEFAULT_HOSTNAME   "localhost"
#define CAPP_CONFIG_MAX_SIZE    1000000

typedef struct capp_process_t {

    char *exec; /* Executable name */
    char *argv; /* Arguments to the executable */
    char *env;  /* Environment variables setup for executable */
} CappProc;

/* Store capp config.json values */
typedef struct capp_config_t {

    char     *version;  /* capp version */
    char     *hostName; /* host name associated with space */
    CappProc *process;  /* Info related to the process/executable */
} CappConfig;

bool process_capp_config_file(CappConfig **config, char *fileName);

#endif /* CAPP_CONFIG_H */
