/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
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
