/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * mesh_config.h
 */

#ifndef MESH_CONFIG_H
#define MESH_CONFIG_H

#include <uuid/uuid.h>

#define CLIENT_CONFIG "client-config"

#define REMOTE_ACCEPT  "remote-accept"
#define LOCAL_ACCEPT   "local-accept"
#define REMOTE_CONNECT "remote-connect"

#define REMOTE_IP_FILE "remote-ip-file"
#define CFG_CERT "cert"
#define CFG_KEY  "key"

#define MAX_BUFFER 256

#define TRUE 1
#define FALSE 0

typedef struct {

	char *remoteIPFile; /* Remote server IP */
	char *certFile;     /* CA Cert file name. */
	char *keyFile;      /* Key file name.*/
} MeshConfig;

int read_mesh_config_file(char *fileName, MeshConfig *config);
void clear_mesh_config(MeshConfig *config);
#endif /* MESH_CONFIG_H */
