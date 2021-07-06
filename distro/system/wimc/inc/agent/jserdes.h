/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef AGENT_JSERDES_H
#define AGENT_JSERDES_H

int serialize_agent_request(AgentReq *request, json_t **json);
int serialize_agent_request_register(AgentReq *req, json_t **json);
int serialize_agent_request_update(AgentReq *req, json_t **json);
int serialize_agent_request_unregister(AgentReq *req, json_t **json);
int deserialize_wimc_request_to_agent(WimcReq *req, json_t *json);
int deserialize_wimc_request(WimcReq *req, json_t *json);

#endif /* AGENT_JSERDES_H */
