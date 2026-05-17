/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>
#include <jansson.h>

#include "agent.h"
#include "wimc.h"
#include "jserdes.h"
#include "utils.h"
#include "agent/jserdes.h"

#include "usys_mem.h"
#include "usys_log.h"

/* utils.c */
extern char *convert_tx_state_to_str(TransferState state);

static void json_log(json_t *json) {

    char *str = NULL;

    str = json_dumps(json, 0);
    if (str) {
        usys_log_debug("json str: %s", str);
        usys_free(str);
    }
}

/*
 * agent_request -> { type: "update",
 *                  type_update: {
 *                      uuid: "same_id",
 *                      total_kbytes: "1234"
 *                      transfer_kbytes:  "34"
 *                      transfer_state: "fetch"
 *                      void_str: "some_string_"
 *			}
 *              }
 *
 */

int serialize_agent_request_update(AgentReq *req, json_t **json) {

  json_t *jupdate=NULL;
  Update *update;
  char idStr[36+1];
  char *state=NULL;

  if (req==NULL || req->update==NULL) {
    return FALSE;
  }

  update = req->update;
  state = convert_tx_state_to_str(update->transferState);

  json_object_set_new(*json, JSON_TYPE, json_string(AGENT_REQ_TYPE_UPDATE));

  /* Add update object */
  json_object_set_new(*json, JSON_TYPE_UPDATE, json_object());
  jupdate = json_object_get(*json, JSON_TYPE_UPDATE);

  uuid_unparse(update->uuid, &idStr[0]);
  json_object_set_new(jupdate, JSON_ID, json_string(idStr));
  json_object_set_new(jupdate, JSON_TOTAL_KBYTES,
		      json_integer(update->totalKB));
  json_object_set_new(jupdate, JSON_TRANSFER_KBYTES,
		      json_integer(update->transferKB));
  json_object_set_new(jupdate, JSON_TRANSFER_STATE,
		      json_string(state));
  
  /* void str is non-zero only if there was an error or we are done (will
   * have final path. Otherwise is empty
   */
  if (update->transferState == (TransferState)ERR ||
      update->transferState == (TransferState)DONE) {
    json_object_set_new(jupdate, JSON_VOID_STR, json_string(update->voidStr));
  } else {
    json_object_set_new(jupdate, JSON_VOID_STR, json_string(""));
  }

  free(state);
  return TRUE;
}

int serialize_agent_request(AgentReq *request, json_t **json) {

    int ret=FALSE;
    json_t *req=NULL;

    *json = json_object();
    if (*json == NULL) {
        return ret;
    }
  
    json_object_set_new(*json, JSON_AGENT_REQUEST, json_object());
    req = json_object_get(*json, JSON_AGENT_REQUEST);

    if (req==NULL) {
        return ret;
    }
    ret = serialize_agent_request_update(request, &req);
    json_log(*json);

    return ret;
}

bool deserialize_wimc_request(WimcReq **request, json_t *json) {

    json_t *jUUID;
    json_t *jInterval;
    json_t *jName;
    json_t *jTag;
    json_t *jMethod;
    json_t *jIndexURL;
    json_t *jStoreURL;
    json_t *jExpected;
    WFetch *fetch;
    WContent *content;
    const char *uuidStr;
    uuid_t parsedUuid;

    jUUID     = NULL;
    jInterval = NULL;
    jName     = NULL;
    jTag      = NULL;
    jMethod   = NULL;
    jIndexURL = NULL;
    jStoreURL = NULL;
    jExpected = NULL;
    fetch     = NULL;
    content   = NULL;
    uuidStr   = NULL;
    uuid_clear(parsedUuid);

    if (request == NULL || json == NULL) {
        return USYS_FALSE;
    }

    if (*request == NULL) {
        *request = (WimcReq *)calloc(1, sizeof(WimcReq));
        if (*request == NULL) {
            return USYS_FALSE;
        }
    }

    jUUID     = json_object_get(json, JSON_ID);
    jInterval = json_object_get(json, JSON_UPDATE_INTERVAL);
    jName     = json_object_get(json, JSON_NAME);
    jTag      = json_object_get(json, JSON_TAG);
    jMethod   = json_object_get(json, JSON_METHOD);
    jIndexURL = json_object_get(json, JSON_INDEX_URL);
    jStoreURL = json_object_get(json, JSON_STORE_URL);
    jExpected = json_object_get(json, JSON_EXPECTED_SIZE);

    if (!jUUID || !jInterval || !jName || !jTag || !jMethod ||
        !jIndexURL || !jStoreURL) {
        usys_log_error("Invalid json received from WIMC");
        json_log(json);
        return USYS_FALSE;
    }

    uuidStr = json_string_value(jUUID);
    if (uuidStr == NULL ||
        uuid_parse(uuidStr, parsedUuid) != 0 ||
        uuid_is_null(parsedUuid)) {
        usys_log_error("Invalid or null WIMC request uuid: %s",
                       uuidStr ? uuidStr : "");
        return USYS_FALSE;
    }

    fetch = (WFetch *)calloc(1, sizeof(WFetch));
    if (fetch == NULL) {
        return USYS_FALSE;
    }

    content = (WContent *)calloc(1, sizeof(WContent));
    if (content == NULL) {
        usys_free(fetch);
        return USYS_FALSE;
    }

    uuid_copy(fetch->uuid, parsedUuid);
    fetch->interval = json_integer_value(jInterval);
    fetch->content = content;

    content->name     = strdup(json_string_value(jName));
    content->tag      = strdup(json_string_value(jTag));
    content->method   = strdup(json_string_value(jMethod));
    content->indexURL = strdup(json_string_value(jIndexURL));
    content->storeURL = strdup(json_string_value(jStoreURL));
    content->expectedSizeBytes = jExpected ?
                                 json_integer_value(jExpected) : 0;

    if (content->name == NULL ||
        content->tag == NULL ||
        content->method == NULL ||
        content->indexURL == NULL ||
        content->storeURL == NULL) {
        usys_free(content->name);
        usys_free(content->tag);
        usys_free(content->method);
        usys_free(content->indexURL);
        usys_free(content->storeURL);
        usys_free(content);
        usys_free(fetch);
        return USYS_FALSE;
    }

    (*request)->type = WREQ_FETCH;
    (*request)->fetch = fetch;

    return USYS_TRUE;
}
