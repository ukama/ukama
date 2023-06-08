/**
 * Copyright (c) 2022-present, Ukama Inc.
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
#include <rabbitmq-c/amqp.h>
#include <rabbitmq-c/tcp_socket.h>

#include "log.h"

#define PREFIX_WEBSOCKET "/websocket"
#define PREFIX_WEBSERVICE "*"

#define EP_PING "/ping"

#define MESH_CLIENT_AGENT "Mesh-client"
#define MESH_CLIENT_VERSION "0.0.1"

#define MESH_TYPE_FWD_REQ "forward_request"
#define MESH_TYPE_FWD_RESP "forward_response"

/* For MAP */
#define MESH_MAP_TYPE_URL  1
#define MESH_MAP_TYPE_HDR  2
#define MESH_MAP_TYPE_POST 3
#define MESH_MAP_TYPE_COOKIE 4

#define MESH_MAP_TYPE_URL_STR  "map_url"
#define MESH_MAP_TYPE_HDR_STR  "map_header"
#define MESH_MAP_TYPE_POST_STR "map_post"
#define MESH_MAP_TYPE_COOKIE_STR "map_cookie"

#define MESH_LOCK_TIMEOUT 10 /* seconds */
#define MAX_QUEUE_SIZE    100
#define MAX_BUFFER        256

#define TRUE  1
#define FALSE 0

#define ENV_WEBSOCKET_PORT   "ENV_WEBSOCKET_PORT"
#define ENV_SERVICES_PORT    "ENV_SERVICES_PORT"
#define ENV_AMQP_HOST        "ENV_AMQP_HOST"
#define ENV_AMQP_PORT        "ENV_AMQP_PORT"
#define ENV_INIT_CLIENT_HOST "ENV_INIT_CLIENT_HOST"
#define ENV_INIT_CLIENT_PORT "ENV_INIT_CLIENT_PORT"
#define ENV_MESH_CERT_FILE   "ENV_MESH_CERT_FILE"
#define ENV_MESH_KEY_FILE    "ENV_MESH_KEY_FILE"
#define ENV_MESH_LOG_LEVEL   "ENV_MESH_LOG_LEVEL"

#define DEFAULT_MESH_AMQP_EXCHANGE "amqp.direct"
#define DEFAULT_MESH_CERT_FILE     "certs/test.cert"
#define DEFAULT_MESH_KEY_FILE      "certs/server.key"

typedef struct _u_instance UInst;
typedef struct _u_request  URequest;
typedef struct _u_response UResponse;
typedef struct _websocket_manager WSManager;
typedef struct _websocket_message WSMessage;
typedef struct _u_map UMap;

typedef struct {

    int  connectionStatus; /* websocket connection status */
	char *nodeID;          /* recevied from the node */
} NodeInfo;

typedef struct {

	uuid_t uuid;
} ServiceInfo;

typedef struct {

	char *reqType; /* Type: forward_request, command, stats, etc. */
	int seqNo;     /* Sequence number of the request. */

	NodeInfo    *nodeInfo;  /* Info. about originating device. */
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

#endif /* MESH_H */
