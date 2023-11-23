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

static void json_log(json_t *json) {

    char *str = NULL;

    str = json_dumps(json, 0);
    if (str) {
        usys_log_debug("json str: %s", str);
        usys_free(str);
    }
}

/*
 * serialize_agent_request -- Serialize Agent request.
 *
 * Agent request are of three types (ReqType):
 * 1. Register
 * 2. Update
 * 3. Unregister
 *
 * Register:
 * { type: "register",
 *   method: "ca-sync",
 *   agent_url: "/a/b/c" }
 *
 * Update:
 * { type: "update",
 *   uuid: "some_id",
 *   cmd: "fetch",
 *          content: {name: "name",
 *                     tag: "tag",
 *                     provider_url: "http://www/www/www"},
 *          callback_url: "http://www.xyz.ccc/cc/cc/",
 *          update_interval: 10 }
 *
 * Unregister:
 * {type: "unregister",
 *  uuid: "same_id"}
 */

/*
 * agent_request -> { type: "register",
 *              type_register: {
 *                      method: "ftp",
 *                      agent_url: "/some/url/path"
 *			}
 *              }
 */

bool serialize_agent_register_request(char *method,
                                      char *url,
                                      json_t **json) {

    json_t   *jreg;
    Register *reg;

    *json = json_object();
    if (*json == NULL) {
        usys_log_error("Unable to initialize json object");
        return USYS_FALSE;
    }
  
    json_object_set_new(*json, JSON_METHOD, json_string(method));
    json_object_set_new(*json, JSON_URL,    json_string(url));

    return USYS_TRUE;
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

/*
 * agent_request -> { type: "unregister",
 *                    type_unregister: {
 *                          uuid: "same_id",
 *                    }
 *                 }
 */

int serialize_agent_request_unregister(AgentReq *req, json_t **json) {

  json_t *jUnReg=NULL;
  UnRegister *unReg;
  char idStr[36+1];

  if (req==NULL || req->unReg==NULL) {
    return FALSE;
  }

  unReg = req->unReg;

  json_object_set_new(*json, JSON_TYPE, json_string(AGENT_REQ_TYPE_UNREG));

  /* Add un-register object */
  json_object_set_new(*json, JSON_TYPE_UNREGISTER, json_object());
  jUnReg = json_object_get(*json, JSON_TYPE_UNREGISTER);

  uuid_unparse(unReg->uuid, &idStr[0]);
  json_object_set_new(jUnReg, JSON_ID, json_string(idStr));

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
#if 0  
    if (request->type == (ReqType)REQ_REG) {
        ret = serialize_agent_request_register(request, &req);
    } else if (request->type == (ReqType)REQ_UPDATE) {
        ret = serialize_agent_request_update(request, &req);
    } else if (request->type == (ReqType)REQ_UNREG) {
        ret = serialize_agent_request_unregister(request, &req);
    }
#endif 
    json_log(*json);

    return ret;
}

static bool deserialize_wimc_request_fetch(WFetch **fetch, json_t *json) {

    json_t *jfetch   = NULL;
    json_t *jcontent = NULL;
    json_t *jObj     = NULL;

    jfetch = json_object_get(json, JSON_TYPE_FETCH);
    if (jfetch == NULL) return USYS_FALSE;

    *fetch = (WFetch *)calloc(1, sizeof(WFetch));
    if (*fetch == NULL) return USYS_FALSE;

    jObj = json_object_get(jfetch, JSON_ID);
    uuid_parse(json_string_value(jObj), (*fetch)->uuid);

    jObj = json_object_get(jfetch, JSON_UPDATE_INTERVAL);
    (*fetch)->interval = json_integer_value(jObj);

    jcontent = json_object_get(jfetch, JSON_CONTENT);
    if (jcontent == NULL) return USYS_FALSE;

    (*fetch)->content = (WContent *)calloc(1, sizeof(WContent));
  
    jObj = json_object_get(jcontent, JSON_NAME);
    (*fetch)->content->name = strdup(json_string_value(jObj));
  
    jObj = json_object_get(jcontent, JSON_TAG);
    (*fetch)->content->tag = strdup(json_string_value(jObj));
  
    jObj = json_object_get(jcontent, JSON_METHOD);
    (*fetch)->content->method = strdup(json_string_value(jObj));

    jObj = json_object_get(jcontent, JSON_INDEX_URL);
    (*fetch)->content->indexURL = strdup(json_string_value(jObj));

    jObj = json_object_get(jcontent, JSON_STORE_URL);
    (*fetch)->content->storeURL = strdup(json_string_value(jObj));

    return USYS_TRUE;
}

bool deserialize_wimc_request(WimcReq **request, json_t *json) {

    json_t *jreq  = NULL;
    json_t *jtype = NULL;

    if (!json) return USYS_FALSE;

    jreq = json_object_get(json, JSON_WIMC_REQUEST);
    if (jreq == NULL) return USYS_FALSE;

    jtype = json_object_get(jreq, JSON_TYPE);
    if (jtype == NULL) return USYS_FALSE;

    if (strcmp(json_string_value(jtype), "fetch") == 0) {
        return deserialize_wimc_request_fetch(&(*request)->fetch, jreq);
    }

    return USYS_FALSE;
}
