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

#include "config.h"

#define JSON_MESH_FORWARD    "mesh_forward"
#define JSON_MESH_CMD        "mesh_cmd"

#define JSON_TYPE            "type"
#define JSON_TYPE_REQUEST    "type_request"
#define JSON_TYPE_RESPONSE   "type_response"

#define JSON_DEVICE_INFO   "device_info"
#define JSON_SERVICE_INFO  "service_info"
#define JSON_REQUEST_INFO  "request_info"
#define JSON_RESPONSE_INFO "response_info"

#define JSON_ID       "uuid"
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

#define JSON_KEY   "key"
#define JSON_VALUE "value"
#define JSON_LEN   "len"

/* Function headers. */
int serialize_response(json_t **json, int size, void *data, uuid_t uuid);
int serialize_forward_request(URequest *request, json_t **json,
			      Config *config, uuid_t uuid);
int deserialize_forward_request(MRequest **req, json_t *json);
int deserialize_response(MResponse **response, json_t *json);

#endif /* MESH_JSERDES_H */
