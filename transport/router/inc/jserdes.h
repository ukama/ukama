/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef JSERDES_H
#define JSERDES_H

#include <jansson.h>
#include <uuid/uuid.h>

#include "router.h"

#define JSON_NAME         "name"
#define JSON_PATTERNS     "patterns"
#define JSON_PATH         "path"
#define JSON_FORWARD      "forward"
#define JSON_IP           "ip"
#define JSON_PORT         "port"
#define JSON_DEFAULT_PATH "default_path"
#define JSON_UUID         "uuid"
#define JSON_ERROR        "error"

int deserialize_delete_route_request(char **uuidStr, json_t *json);
int deserialize_post_route_request(Service **service, json_t *json);
int serialize_post_route_response(json_t **json, int respCode, uuid_t uuid,
				  char *errStr);

#endif /* JSERDES_H */
