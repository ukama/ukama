/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#ifndef INIT_CLIENT_JSERDES_H
#define INIT_CLIENT_JSERDES_H

#include <jansson.h>

#include "initClient.h"

#define JSON_SYSTEM_NAME "systemName"
#define JSON_SYSTEM_ID   "systemId"
#define JSON_CERTIFICATE "certificate"
#define JSON_IP          "ip"
#define JSON_PORT        "port"
#define JSON_HEALTH      "health"
#define JSON_GLOBAL_UUID "global_uuid"
#define JSON_LOCAL_UUID  "local_uuid"

#define JSON_NODE_GW_IP   "nodeGWip"
#define JSON_NODE_GW_PORT "nodeGWPort"

int serialize_request(Request *request, json_t **json);
int deserialize_response(ReqType reqType, QueryResponse **queryResponse,
						 char *str);
int serialize_uuids_from_file(SystemRegistrationId *sysReg, json_t **json);
int deserialize_uuids_from_file(char* str, SystemRegistrationId** sysReg);

#endif /* INIT_CLIENT_JSERDES_H */
