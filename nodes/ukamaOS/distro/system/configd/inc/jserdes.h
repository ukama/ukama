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

#include "json_types.h"
#include "web_service.h"
#include "usys_types.h"
#include "session.h"

#define EMPTY_STRING  ""

#define JSON_OK                        STATUS_OK
#define JSON_FAILURE                   STATUS_NOTOK
#define JSON_ENCODING_OK               JSON_OK
#define JSON_DECODING_OK               JSON_OK

bool json_deserialize_config_data(JsonObj *json,ConfigData **cd);
bool json_deserialize_node_id(char **nodeID, json_t *json);
bool json_deserialize_running_config(char* file, ConfigData **cd);
void json_log(json_t *json);
void json_free(JsonObj** json);

#endif /* JSERDES_H_ */
