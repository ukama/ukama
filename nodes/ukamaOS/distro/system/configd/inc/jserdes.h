/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#ifndef JSERDES_H
#define JSERDES_H

#include <jansson.h>

#include "json_types.h"
#include "web_service.h"
#include "usys_types.h"
#include "session.h"

bool json_deserialize_config_data(JsonObj *json, SessionData **sd);
bool json_deserialize_node_id(char **nodeID, json_t *json);
void json_log(json_t *json);
void json_free(JsonObj** json);

#endif /* JSERDES_H_ */
