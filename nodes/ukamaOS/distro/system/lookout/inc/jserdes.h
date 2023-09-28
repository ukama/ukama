/**
 * Copyright (c) 2023-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef JSERDES_H_
#define JSERDES_H_

#include <jansson.h>

#include "lookout.h"
#include "json_types.h"

#include "usys_types.h"

#define EMPTY_STRING  ""

#define JSON_OK           STATUS_OK
#define JSON_FAILURE      STATUS_NOTOK
#define JSON_ENCODING_OK  JSON_OK
#define JSON_DECODING_OK  JSON_OK

void json_log(json_t *json);
void json_free(JsonObj** json);
bool json_deserialize_node_id(char **nodeID, json_t *json);
bool json_deserialize_capps(CappList **cappList, JsonObj *json);
bool json_serialize_health_report(JsonObj **json, CappList *list);

#endif /* JSERDES_H_ */
