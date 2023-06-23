/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef MESH_H
#define MESH_H

#include <getopt.h>
#include <ulfius.h>
#include <uuid/uuid.h>

#include "log.h"

#define DEF_FILENAME "cert.crt"
#define DEF_CA_FILE  ""
#define DEF_CRL_FILE ""
#define DEF_CA_PATH  ""
#define DEF_SERVER_NAME "localhost"
#define DEF_TLS_SERVER_PORT "4444"
#define DEF_LOG_LEVEL "TRACE"
#define DEF_CLOUD_SERVER_NAME "localhost"
#define DEF_CLOUD_SERVER_PORT "4444"
#define DEF_CLOUD_SERVER_CERT "certs/test.crt"

#define TRUE  1
#define FALSE 0

#define PROXY_NONE    0x01
#define PROXY_FORWARD 0x02
#define PROXY_REVERSE 0x04

#define PREFIX_WEBSOCKET "/websocket"
#define PREFIX_WEBSERVICE "*"

#define MESH_CLIENT_AGENT "Mesh-client"
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

#define MESH_LOCK_TIMEOUT 10 /* seconds */

typedef struct _u_instance UInst;
typedef struct _u_request  URequest;
typedef struct _u_response UResponse;
typedef struct _websocket_manager WSManager;
typedef struct _websocket_message WSMessage;
typedef struct _u_map UMap;

typedef struct {
    void      *config;
    WSManager *handler;
} ThreadArgs;

typedef struct {

    struct _websocket_client_handler *handler;
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
	int seqNo;     /* Sequence number of the request. */

	DeviceInfo  *deviceInfo;  /* Info. about originating device. */
	ServiceInfo *serviceInfo; /* Info. about origniating service. */
	URequest    *requestInfo; /* Actual request. */
} MRequest;

/* struct to define the response back from the service. */
typedef struct {

	char        *reqType;
	int         seqNo;
	int         size;
	void        *data;
	ServiceInfo *serviceInfo;
} MResponse;

typedef struct {

    char        *reqType;
    int         seqNo;
    NodeInfo    *nodeInfo;
    ServiceInfo *serviceInfo;
    int         dataSize;
    int         code;
    char        *data;   /* RequestInfo or actual response */
} Message;

#endif /* MESH_H */
