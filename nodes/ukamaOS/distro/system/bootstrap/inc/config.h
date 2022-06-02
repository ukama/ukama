/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * config.h
 */

#ifndef CONFIG_H
#define CONFIG_H

/* used in the config file and for parsing. */
#define CONFIG "config"

#define NODED_HOST       "noded-host"
#define NODED_PORT       "noded-port"
#define MESH_CONFIG      "mesh-config"
#define REMOTE_IP_FILE   "remote-ip-file"
#define BOOTSTRAP_SERVER "bootstrap-server"

/* Some defaults */
#define DEF_CONFIG_FILE      "config.toml"
#define DEF_NODED_HOST       "localhost"
#define DEF_NODED_PORT       "8095"
#define DEF_MESH_CONFIG      "/conf/mesh/config.toml"
#define DEF_REMOTE_IP_FILE   "/conf/mesh/ip_file"
#define DEF_BOOTSTRAP_SERVER "kickstart.ukama.com"

#define MAX_BUFFER 256

#define TRUE  1
#define FALSE 0

/* Struct to define the configuration. */
typedef struct {

	char *nodedHost;       /* Host where node.d is running (hostname) */
	char *nodedPort;       /* Port where node.d is listening */
	char *meshConfig;      /* Mesh.d configuration file */
	char *remoteIPFile;    /* file storing the remote server IP */
	char *bootstrapServer; /* Bootstrap server */
} Config;

int process_config_file(char *fileName, Config *config);
void clear_config(Config *config);
void print_config(Config *config);
#endif /* CONFIG_H */
