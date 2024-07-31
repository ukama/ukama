/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

/*
 * server.h
 */

#ifndef SERVER_H
#define SERVER_H

#define TRUE  1
#define FALSE 0

#define API_VERSION "v1"
#define EP_NODES    "nodes"

#define MAX_BACKOFF     30
#define MAX_GET_URL_LEN 2048

typedef struct _response {

	char *buffer;
	size_t size;
} Response;

/* Struct to define the server */
typedef struct {

	char *IP;   /* Server's IPv4 for Mesh.d */
	char *cert; /* Cert for connection with Server */
	char *org;  /* Organization this Node belong's */
} ServerInfo;

int send_request_to_init_with_exponential_backoff(char *bootstrapServer,
                                                  int bootstrapPort,
                                                  char *uuid,
                                                  ServerInfo *server);
void free_server_info(ServerInfo *server);
void log_debug_server(ServerInfo *server);

#endif /* SERVER_H */
