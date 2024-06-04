/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#ifndef CONFIG_H
#define CONFIG_H

/* used in the config file and for parsing. */
#define CONFIG "config"

#define NODED_HOST       "noded-host"
#define NODED_PORT       "noded-port"
#define MESH_CONFIG      "mesh-config"
#define REMOTE_IP_FILE   "remote-ip-file"

/* Some defaults */
#define DEF_CONFIG_FILE      "config.toml"
#define DEF_NODED_HOST       "localhost"
#define DEF_MESH_CONFIG      "/conf/mesh/config.toml"
#define DEF_REMOTE_IP_FILE   "/conf/mesh/ip_file"
#define DEF_BOOTSTRAP_SERVER "kickstart.ukama.com"
#define DEF_BOOTSTRAP_FILE   "/ukama/bootstrap"

#define MAX_BUFFER 256

#define TRUE  1
#define FALSE 0

/* Struct to define the configuration. */
typedef struct {

	char *nodedHost;       /* Host where node.d is running (hostname) */
	int   nodedPort;       /* Port where node.d is listening */
	char *meshConfig;      /* Mesh.d configuration file */
	char *remoteIPFile;    /* file storing the remote server IP */
	char *bootstrapRemoteServer; /* Bootstrap server */
    int   bootstrapRemotePort;   /* Bootstrap listening port */
} Config;

int process_config_file(char *fileName, Config *config);
void clear_config(Config *config);
void print_config(Config *config);
#endif /* CONFIG_H */
