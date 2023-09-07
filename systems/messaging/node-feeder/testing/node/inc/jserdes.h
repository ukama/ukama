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

#include "node.h"

#define JSON_STRING  1
#define JSON_INTEGER 2

#define JSON_NODE_INFO       "nodeInfo"
#define JSON_NODE_CONFIG     "nodeConfig"
#define JSON_TYPE            "type"
#define JSON_PART_NUMBER     "partNumber"
#define JSON_SKEW            "skew"
#define JSON_SW_VERSION      "swVersion"
#define JSON_MFG_SW_VERSION  "mfgSwVersion"
#define JSON_ASSEMBLY_DATE   "assemblyDate"
#define JSON_OEM             "oem"
#define JSON_MFG_TEST_STATUS "mfgTestStatus"
#define JSON_STATUS          "status"
#define JSON_MODULE_ID       "moduleID"
#define JSON_HW_VERSION      "hwVersion"
#define JSON_MAC             "mac"
#define JSON_MFG_DATE        "mfgDate"
#define JSON_MFG             "mfgName"

int deserialize_node(Node **node, json_t *json);

#endif /* JSERDES_H */
