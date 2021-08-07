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

#include "log.h"
#include "config.h"

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

typedef struct _u_instance UInst;
typedef struct _u_request  URequest;
typedef struct _u_response UResponse;
typedef struct _websocket_manager WSManager;
typedef struct _websocket_message WSMessage;

#endif /* MESH_H */
