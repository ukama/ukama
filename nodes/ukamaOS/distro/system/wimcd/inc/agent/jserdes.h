/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef AGENT_JSERDES_H
#define AGENT_JSERDES_H

#include "usys_types.h"

int serialize_agent_request(AgentReq *request, json_t **json);
int serialize_agent_request_register(AgentReq *req, json_t **json);
int serialize_agent_request_update(AgentReq *req, json_t **json);
int serialize_agent_request_unregister(AgentReq *req, json_t **json);

bool deserialize_wimc_request(WimcReq **request, json_t *json);

#endif /* AGENT_JSERDES_H */
