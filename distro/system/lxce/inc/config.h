/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * config.h
 */

#ifndef LXCE_CONFIG_H
#define LXCE_CONFIG_H

/* used in the config file and for parsing. */
#define CONFIG "config"

#define LOCAL_ACCEPT "local-accept"
#define LOCAL_EP     "local-ep"
#define WIMC_HOST    "wimc-host"
#define WIMC_PORT    "wimc-port"
#define MESH_PORT    "mesh-port"

/* Some defaults */
#define DEF_LOCAL_ACCEPT "4448"
#define DEF_WIMC_PORT    "4441"
#define DEF_MESH_PORT    "4444"
#define DEF_WIMC_HOST    "localhost"
#define DEF_LOCAL_EP     "/lxce/"

#define MAX_BUFFER 256

#define TRUE  1
#define FALSE 0

/* Struct to define the configuration. */
typedef struct {

  char *localAccept; /* Port on which to accept clients */
  char *localEP;     /* root for lxce.d REST interface */

  char *wimcHost;    /* Host where WIMC.d is running (hostname) */
  char *wimcPort;    /* Port where WIMC.d is listening */
  char *meshPort;    /* Port where MESH.d is listening */
} Config;

int process_config_file(char *fileName, Config *config);
void clear_config(Config *config);
void print_config(Config *config);

#endif /* LXCE_CONFIG_H */
