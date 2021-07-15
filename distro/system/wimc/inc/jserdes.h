/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef WIMC_JSERDES_H
#define WIMC_JSERDES_J

#include "wimc.h"
#include "agent.h"
#include "err.h"

#define JSON_AGENT_CB     "agent-cbURL"
#define JSON_AGENT_CMD    "cmd"
#define JSON_AGENT_METHOD "method"

#define JSON_TYPE            "type"
#define JSON_TYPE_REGISTER   "type_register"
#define JSON_TYPE_UNREGISTER "type_unregister"
#define JSON_TYPE_UPDATE     "type_update"
#define JSON_TYPE_FETCH      "type_fetch"
#define JSON_TYPE_CANCEL     "type_cancel"
#define JSON_TYPE_PROCESS    "type_processing"
#define JSON_TYPE_RESULT     "type_result"
#define JSON_TYPE_STATUS     "type_status"

#define JSON_METHOD          "method"
#define JSON_URL             "url"
#define JSON_ID              "uuid"
#define JSON_CMD             "cmd"
#define JSON_ACTION          "action"
#define JSON_CONTENT         "content"
#define JSON_NAME            "name"
#define JSON_TAG             "tag"
#define JSON_AGENT_URL       "agent_url"
#define JSON_PROVIDER_URL    "provider_url"
#define JSON_CALLBACK_URL    "callback_url"
#define JSON_UPDATE_INTERVAL "update_interval"
#define JSON_CONTENT         "content"
#define JSON_TASK            "task"
#define JSON_STATUS          "status"

#define JSON_EVENT           "event_type"
#define JSON_EVENT_UPDATE    "update"
#define JSON_TOTAL_KBYTES    "total_kilobytes"
#define JSON_TRANSFER_KBYTES "transfer_kilobytes"
#define JSON_TRANSFER_STATE  "transfer_state"
#define JSON_VOID_STR        "void"

#define JSON_WIMC_REQUEST    "wimc_request"
#define JSON_AGENT_REQUEST   "agent_request"

#define JSON_PROVIDER_RESPONSE "provider_response"
#define JSON_WIMC_RESPONSE     "wimc_response"

/* function def. */
int serialize_wimc_request(WimcReq *request, json_t **json);
static int serialize_wimc_request_fetch(WimcReq *req, json_t **json);

#endif /* WIMC_JSERDES_H */
