/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#ifndef JSERDES_H
#define JSERDES_H

#include <jansson.h>

#include "nodeInfo.h"

#define JTAG_NODE_INFO      "nodeInfo"
#define JTAG_UUID           "UUID"
#define JTAG_NAME           "name"
#define JTAG_TYPE           "type"
#define JTAG_PART_NUMBER    "partNumber"
#define JTAG_SKEW           "skew"
#define JTAG_MAC            "mac"
#define JTAG_OEM            "oemName"
#define JTAG_ASSEMBLY_DATE  "assemblyDate"

#define JTAG_LOGS      "logs"
#define JTAG_APP_NAME  "app_name"
#define JTAG_TIME      "time"
#define JTAG_LEVEL     "level"
#define JTAG_MESSAGE   "message"

int deserialize_node_info(NodeInfo **nodeInfo, json_t *json);
#endif /* JSERDES_H */
