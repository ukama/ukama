/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef MESH_JSERDES_H
#define MESH_JSERDES_H

#include <jansson.h>
#include <uuid/uuid.h>

#include "mesh.h"
#include "config.h"
#include "initClient.h"

#define JSON_MESH_FORWARD    "mesh_forward"
#define JSON_MESH_CMD        "mesh_cmd"

#define JSON_TYPE            "type"
#define JSON_TYPE_REQUEST    "type_request"
#define JSON_TYPE_RESPONSE   "type_response"

#define JSON_NODE_INFO     "node_info"
#define JSON_SERVICE_INFO  "service_info"
#define JSON_REQUEST_INFO  "request_info"
#define JSON_RESPONSE_INFO "response_info"

#define JSON_NODE_ID  "node_id"
#define JSON_NAME     "name"
#define JSON_PROTOCOL "protocol"
#define JSON_METHOD   "method"
#define JSON_URL      "url"
#define JSON_PATH     "path"
#define JSON_MAP      "map"
#define JSON_MAP_URL  "map_url"
#define JSON_MAP_HDR  "map_header"
#define JSON_MAP_POST "map_post"
#define JSON_RAW_DATA "raw_data"
#define JSON_LENGTH   "length"
#define JSON_DATA     "data"
#define JSON_SEQ      "seq"
#define JSON_PORT     "port"
#define JSON_MESSAGE  "message"

/* for system info*/
#define JSON_SYSTEM_NAME "systemName"
#define JSON_SYSTEM_ID   "systemId"
#define JSON_CERTIFICATE "certificate"
#define JSON_IP          "ip"
#define JSON_PORT        "port"
#define JSON_HEALTH      "health"

#define JSON_KEY   "key"
#define JSON_VALUE "value"
#define JSON_LEN   "len"

/* Function headers. */
int serialize_response(json_t **json, int size, void *data, char *nodeID);
int serialize_websocket_message(char **str, URequest *request, char *nodeID,
                                char *nodePort, char *agent);
int serialize_device_info(json_t **json, NodeInfo *device);
int deserialize_forward_request(MRequest **req, json_t *json);
int deserialize_response(MResponse **response, json_t *json);
int deserialize_system_info(SystemInfo **systemInfo, json_t *json);
int deserialize_websocket_message(Message **message, char *data);

#endif /* MESH_JSERDES_H */
