/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
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
bool json_serialize_health_report(JsonObj **json,
                                  char *nodeID,
                                  CappList *list,
                                  GPSClientData *gps);

#endif /* JSERDES_H_ */
