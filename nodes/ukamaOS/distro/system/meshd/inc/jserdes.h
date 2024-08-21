/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef MESH_JSERDES_H
#define MESH_JSERDES_H

#include <jansson.h>
#include <uuid/uuid.h>

#include "config.h"

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
#define JSON_CODE     "code"
#define JSON_UUID     "uuid"

#define JSON_KEY   "key"
#define JSON_VALUE "value"
#define JSON_LEN   "len"

int serialize_local_service_response(char **response, Message *message,
                                     int code, int len, char *data);
int serialize_websocket_message(char **str, URequest *request, char *uuid);
int deserialize_forward_request(MRequest **req, json_t *json);
int deserialize_response(MResponse **response, json_t *json);
int deserialize_websocket_message(Message **message, json_t *json);
int deserialize_request_info(URequest **request, char *str);

#endif /* MESH_JSERDES_H */
