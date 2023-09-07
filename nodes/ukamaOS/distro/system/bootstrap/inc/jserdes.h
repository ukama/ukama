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

#include "nodeInfo.h"
#include "server.h"

#define JSON_NODE_INFO      "nodeInfo"
#define JSON_UUID           "UUID"
#define JSON_NAME           "name"
#define JSON_TYPE           "type"
#define JSON_PART_NUMBER    "partNumber"
#define JSON_SKEW           "skew"
#define JSON_MAC            "mac"
#define JSON_OEM            "oemName"
#define JSON_ASSEMBLY_DATE  "assemblyDate"

/* serverInfo */
#define JSON_NODE        "nodeId"
#define JSON_ORG         "orgName"
#define JSON_IP          "ip"
#define JSON_CERTIFICATE "certificate"

int deserialize_node_info(NodeInfo **nodeInfo, json_t *json);
int deserialize_server_info(ServerInfo *serverInfo, json_t *json);
#endif /* JSERDES_H */
