/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <stdio.h>
#include <string.h>
#include <jansson.h>

#include "agent.h"
#include "wimc.h"
#include "jserdes.h"

#include "utils.h"

#include "agent/jserdes.h"

/* JSON (de)-serialization functions. */

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
 *   id: "some_id",
 *   cmd: "fetch",
 *          content: {name: "name",
 *                     tag: "tag",
 *                     provider_url: "http://www/www/www"},
 *          callback_url: "http://www.xyz.ccc/cc/cc/",
 *          update_interval: 10 }
 *
 * Unregister:
 * {type: "unregister",
 *  id: "same_id"}
 */

/*
 * serialize_agent_request --
 *
 */

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
  
  if (request->type == (ReqType)REQ_REG) {
    ret = serialize_agent_request_register(request, &req);
  } else if (request->type == (ReqType)REQ_UPDATE) {
    ret = serialize_agent_request_update(request, &req);
  } else if (request->type == (ReqType)REQ_UNREG) {
    ret = serialize_agent_request_unregister(request, &req);
  }

  if (ret) {
    
    char *str;
    str = json_dumps(*json, 0);

    if (str) {
      log_debug("Agent request str: %s", str);
      free(str);
    }
    ret = TRUE;
  }

  return ret;
}

/*
 * agent_request -> { type: "register",
 *              type_register: {
 *                      method: "ftp",
 *                      agent_url: "/some/url/path"
 *			}
 *              }
 */

/*
 * serialize_agent_request_register -- register agent into WIMC.d
 *
 */
int serialize_agent_request_register(AgentReq *req, json_t **json) {

  json_t *jreg;
  Register *reg;

  if (req==NULL && req->reg==NULL) {
    return FALSE;
  }

  reg = req->reg;

  json_object_set_new(*json, JSON_TYPE, json_string(AGENT_REQ_TYPE_REG));

  /* Add register object */
  json_object_set_new(*json, JSON_TYPE_REGISTER, json_object());
  jreg = json_object_get(*json, JSON_TYPE_REGISTER);

  json_object_set_new(jreg, JSON_METHOD, json_string(reg->method));
  json_object_set_new(jreg, JSON_AGENT_URL, json_string(reg->url));

  return TRUE;
}

/*
 * agent_request -> { type: "update",
 *                  type_update: {
 *                      id: "same_id",
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

  if (req==NULL || req->update==NULL) {
    return FALSE;
  }

  update = req->update;

  json_object_set_new(*json, JSON_TYPE, json_string(AGENT_REQ_TYPE_UPDATE));

  /* Add update object */
  json_object_set_new(*json, JSON_TYPE_UPDATE, json_object());
  jupdate = json_object_get(*json, JSON_TYPE_UPDATE);

  json_object_set_new(jupdate, JSON_ID, json_integer(update->id));
  json_object_set_new(jupdate, JSON_TOTAL_KBYTES,
		      json_integer(update->totalKB));
  json_object_set_new(jupdate, JSON_TRANSFER_KBYTES,
		      json_integer(update->transferKB));
  json_object_set_new(jupdate, JSON_TRANSFER_STATE,
		      json_string(convert_state_to_str(update->transferState)));
  
  /* void str is non-zero only if there was an error otherwise is empty */
  if (update->transferState == (TransferState)ERR) {
    json_object_set_new(jupdate, JSON_VOID_STR, json_string(update->voidStr));
  }
  
  return TRUE;
}

/*
 * agent_request -> { type: "unregister",
 *                    type_unregister: {
 *                          id: "same_id",
 *                    }
 *                 }
 */

int serialize_agent_request_unregister(AgentReq *req, json_t **json) {

  json_t *jUnReg=NULL;
  UnRegister *unReg;

  if (req==NULL || req->unReg==NULL) {
    return FALSE;
  }

  unReg = req->unReg;

  json_object_set_new(*json, JSON_TYPE, json_string(AGENT_REQ_TYPE_UNREG));

  /* Add un-register object */
  json_object_set_new(*json, JSON_TYPE_UNREGISTER, json_object());
  jUnReg = json_object_get(*json, JSON_TYPE_UNREGISTER);

  json_object_set_new(jUnReg, JSON_ID, json_integer(unReg->id));

  return TRUE;
}

/*
 * deserialize_wimc_request_to_agent --
 *
 */

int deserialize_wimc_request_to_agent(WimcReq *req, json_t *json) {

  json_t *obj, *jcont;

#if 0
  if (json==NULL) {
    return FALSE;
  }

  obj = json_object_get(json, JSON_TYPE);
  if (obj) {
     req->type = convert_str_to_wType(json_string_value(obj));
  } else {
    return FALSE;
  }

  obj = json_object_get(json, JSON_ACTION);
  if (obj) {
     req->action = convert_str_to_action(json_string_value(obj));
  } else {
    return FALSE;
  }

  obj = json_object_get(json, JSON_ID);
  if (obj) {
    req->id = json_integer_value(obj);
  } else {
    return FALSE;
  }

  if (req->action == (ActionType)ACTION_FETCH ||
      req->action == (ActionType)ACTION_UPDATE) {

    obj = json_object_get(json, JSON_CALLBACK_URL);
    if (obj) {
      req->callbackURL = json_string_value(obj);
    } else {
      return FALSE;
    }

    obj = json_object_get(json, JSON_UPDATE_INTERVAL);
    if (obj) {
      req->interval = json_integer_value(obj);
    } else {
      return FALSE;
    }
  }

  if (req->action == (ActionType)ACTION_FETCH) {

    /* Get content object. */
    jcont = json_object_get(json, JSON_CONTENT);
    if (jcont==NULL) {
      return FALSE;
    }

    req->content = (Content *)calloc(sizeof(Content), 1);
    if (req->content == NULL) {
      return FALSE;
    }

    obj = json_object_get(jcont, JSON_NAME);
    if (obj) {
      req->content->name = json_string_value(obj);
    } else {
      free(req->content);
      return FALSE;
    }

    obj = json_object_get(jcont, JSON_TAG);
    if (obj) {
      req->content->tag = json_string_value(obj);
    } else {
      free(req->content);
      return FALSE;
    }

    obj = json_object_get(jcont, JSON_METHOD);
    if (obj) {
      req->content->method = convert_str_to_method(obj);
    } else {
      free(req->content);
      return FALSE;
    }

    obj = json_object_get(jcont, JSON_PROVIDER_URL);
    if (obj) {
      req->content->providerURL = json_string_value(obj);
    } else {
      free(req->content);
      return FALSE;
    }
  }

  if (req->action == (ActionType)ACTION_CANCEL) {
    /* Do nothing. */
  }
#endif
  return TRUE;
}

/*
 * deserialize_wimc_request --
 *
 */
int deserialize_wimc_request(WimcReq **request, json_t *json) {

  int ret=FALSE;
  char *str;
  json_t *jreq=NULL, *jtype=NULL;

  WimcReq *req = *request;

  /* sanity check. */
  if (!json) {
    return FALSE;
  }

  if (json) {
    char *str;

    str = json_dumps(json, 0);
    if (str) {
      log_debug("Deserializeing JSON: %s", str);
      free(str);
    }
  }
  
  jreq = json_object_get(json, JSON_WIMC_REQUEST);
  if (jreq == NULL) {
    return FALSE;
  }
    
  jtype = json_object_get(jreq, JSON_TYPE);
  if (jtype==NULL) {
    return FALSE;
  }

  req->type = convert_str_to_type(json_string_value(jtype));

  if (req->type == (WReqType)WREQ_FETCH) {
    ret = deserialize_wimc_request_fetch(&req->fetch, jreq);
  } else if (req->type == (WReqType)WREQ_UPDATE) {

  }

  return ret;
}

/*
 * deserialize_wimc_request_fetch --
 *
 */
static int deserialize_wimc_request_fetch(WFetch **fetch, json_t *json) {

  json_t *jfetch=NULL, *jcontent=NULL, *jObj=NULL;

  jfetch = json_object_get(json, JSON_TYPE_FETCH);
  if (jfetch == NULL) {
    return FALSE;
  }

  *fetch = (WFetch *)calloc(1, sizeof(WFetch));
  if (*fetch == NULL) {
    return FALSE;
  }

  jObj = json_object_get(jfetch, JSON_ID);
  (*fetch)->id = json_integer_value(jObj);

  jObj = json_object_get(jfetch, JSON_UPDATE_INTERVAL);
  (*fetch)->interval = json_integer_value(jObj);

  jObj = json_object_get(jfetch, JSON_CALLBACK_URL);
  (*fetch)->cbURL = strdup(json_string_value(jObj));

  jcontent = json_object_get(jfetch, JSON_CONTENT);
  if (jcontent == NULL) {
    return FALSE;
  }

  (*fetch)->content = (WContent *)calloc(1, sizeof(WContent));
  
  jObj = json_object_get(jcontent, JSON_NAME);
  (*fetch)->content->name = strdup(json_string_value(jObj));
  
  jObj = json_object_get(jcontent, JSON_TAG);
  (*fetch)->content->tag = strdup(json_string_value(jObj));
  
  jObj = json_object_get(jcontent, JSON_PROVIDER_URL);
  (*fetch)->content->providerURL = strdup(json_string_value(jObj));

  jObj = json_object_get(jcontent, JSON_METHOD);
  (*fetch)->content->method = strdup(json_string_value(jObj));

  return TRUE;
}
