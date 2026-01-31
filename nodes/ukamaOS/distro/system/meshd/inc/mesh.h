/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef MESH_H
#define MESH_H

#include <getopt.h>
#include <ulfius.h>
#include <uuid/uuid.h>

#include "usys_api.h"
#include "usys_file.h"
#include "usys_types.h"
#include "usys_log.h"
#include "usys_services.h"

#define SERVICE_NAME SERVICE_MESH

#define DEF_FILENAME          "cert.crt"
#define DEF_CA_FILE           ""
#define DEF_CRL_FILE          ""
#define DEF_CA_PATH           ""
#define DEF_SERVER_NAME       "localhost"
#define DEF_TLS_SERVER_PORT   "4444"
#define DEF_LOG_LEVEL         "TRACE"
#define DEF_CLOUD_SERVER_NAME "localhost"
#define DEF_CLOUD_SERVER_PORT "4444"
#define DEF_CLOUD_SERVER_CERT "certs/test.crt"

#define TRUE  1
#define FALSE 0

#define PROXY_NONE    0x01
#define PROXY_FORWARD 0x02
#define PROXY_REVERSE 0x04

#define PREFIX_WEBSOCKET "/websocket"
#define PREFIX_FWDSERVICE "*"

#define MESH_CLIENT_AGENT   "Mesh-client"
#define MESH_CLIENT_VERSION "0.0.1"

#define MESH_SERVICE_REQUEST  "service_request"
#define MESH_SERVICE_RESPONSE "service_response"
#define MESH_NODE_REQUEST     "node_request"
#define MESH_NODE_RESPONSE    "node_response"

/* For MAP */
#define MESH_MAP_TYPE_URL  1
#define MESH_MAP_TYPE_HDR  2
#define MESH_MAP_TYPE_POST 3
#define MESH_MAP_TYPE_COOKIE 4

#define MESH_MAP_TYPE_URL_STR    "map_url"
#define MESH_MAP_TYPE_HDR_STR    "map_header"
#define MESH_MAP_TYPE_POST_STR   "map_post"
#define MESH_MAP_TYPE_COOKIE_STR "map_cookie"

#define MESH_LOCK_TIMEOUT 1 /* seconds */

#ifndef SAFE_FREE
#define SAFE_FREE(p) do { if ((p) != NULL) { free(p); (p) = NULL; } } while (0)
#endif

typedef struct _u_instance UInst;
typedef struct _u_request  URequest;
typedef struct _u_response UResponse;
typedef struct _websocket_manager WSManager;
typedef struct _websocket_message WSMessage;
typedef struct _u_map UMap;

typedef struct {
    void   *config;
    struct _websocket_client_handler *handler;
} ThreadArgs;

typedef struct {

    struct _websocket_client_handler *handler;
    struct _u_instance *fwdInst;
    struct _u_instance *webInst;
    void    *config;
} State;

typedef struct {

	char *nodeID;
    char *port;
} NodeInfo;

typedef struct {

    char *nodeID;
} DeviceInfo;

typedef struct {

    char *name;
    char *port;
} ServiceInfo;

typedef struct {

	char *reqType; /* Type: forward_request, command, stats, etc. */
	char *seqNo;     /* Sequence number of the request. */

	DeviceInfo  *deviceInfo;  /* Info. about originating device. */
	ServiceInfo *serviceInfo; /* Info. about origniating service. */
	URequest    *requestInfo; /* Actual request. */
} MRequest;

/* struct to define the response back from the service. */
typedef struct {

	char        *reqType;
	char        *seqNo;
	int         size;
	void        *data;
	ServiceInfo *serviceInfo;
} MResponse;

typedef struct {

    char        *reqType;
    char        *seqNo;
    NodeInfo    *nodeInfo;
    ServiceInfo *serviceInfo;
    int         dataSize;
    int         code;
    char        *data;   /* RequestInfo or actual response */
} Message;

#endif /* MESH_H */
