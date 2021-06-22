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

#include "jserdes.h"
#include "agent.h"


/*
 * deserialize_provider_response --
 *
 */

int deserialize_provider_response(json_t *resp, AgentCB **agent) {

  int i=0, j=0, count=0;
  json_t *array=NULL;
  json_t *elem=NULL, *method=NULL, *url=NULL;
  
  if (!resp) return FALSE;
  
  array = json_object_get(resp, JSON_AGENT_URL);

  if (json_is_array(array)) {
    count = json_array_size(array);

    agent = (AgentCB **)calloc(sizeof(AgentCB), count);
    
    for (i=0; i<count; i++) {
      elem = json_array_get(array, i);

      if (elem == NULL) {
	goto failure;
      }
      
      method = json_object_get(elem, JSON_AGENT_METHOD);
      url  = json_object_get(elem, JSON_AGENT_URL);

      if (method && url) {
	agent[i]->method = strdup(json_string_value(method));
	agent[i]->url = strdup(json_string_value(url));
      }
    }
  }

  return TRUE;
  
 failure:
  for (j=0; j<i; j++) {
    free(agent[j]->method);
    free(agent[j]->url);
  }

  if (&agent[0]) free(&agent[0]);
  return FALSE;
}

/*
 * serialize_wimc_request -- Serialize WIMC.d request.
 *
 * WIMC requests are of two types:
 * 1. Provider: asking provider for content (URL for Agent).
 * 2. Agent: asking agent to initiate content transfer using provider CB url.
 *
 * wimc_request -> { type:   "agent",
 *                   action: "fetch",
 *                   id:     "some_id",
 *                   callback_url: "/a/b/c/",
 *                   interval: 10,
 *                   content { name:   "name",
 *                             tag:    "tag",
 *                             method: "raw",
 *                             provider_url: "/a/b/c" }
 *                 }
 *
 * wimc_request -> {type:   "agent",
 *                  action: "cancel",
 *                  id:     "some_id"}
 *
 * wimc_request -> {type:     "agent",
 *                  action:   "update",
 *                  id:       "some_id",
 *                  interval: "100",
 *                  callback_url: "/a/b/d"}
 *
 * wimc_request -> {type:   "provider",
 *                  action: "request",
 *                  content { name : "name",
 *                            tag: "tag"}
 *
 */

json_t *serialize_wimc_request(WimcReq *request) {

  int ret=FALSE;
  json_t *json, *req;

  json = json_object();
  if (json == NULL) {
    return NULL;
  }

  json_object_set_new(json, JSON_WIMC_REQUEST, json_object());

  req = json_object_get(json, JSON_WIMC_REQUEST);

  if (req==NULL) {
    return NULL;
  }

  if (request->type == (WReqType)AGENT) {
    ret = serialize_wimc_request_to_agent(request, json);
  } else if (request->type == (WReqType)PROVIDER) {
    ret = serialize_wimc_request_to_provider(request, json);
  }

  if (ret) {
    return json;
  }

  return NULL;
}

/*
 * serialize_wimc_request_to_provider --
 *
 */
int serialize_wimc_request_to_provider(WimcReq *req, json_t *json) {
  /* Currently we are using simple GET with URL specifying the container
   * name and tag. Probably good idea to use JSON so we can expand
   * this interface as needed.
   */
}

/*
 * serialize_wimc_request_to_agent --
 *
 */
int serialize_wimc_request_to_agent(WimcReq *req, json_t *json) {

  json_t *jreq, *jcont;

  jreq = json_object_get(json, JSON_WIMC_REQUEST);

  if (req==NULL) {
    return FALSE;
  }
  json_object_set_new(jreq, JSON_TYPE,
		      json_string(convert_type_to_str(req->type)));
  json_object_set_new(jreq, JSON_ID, json_integer(req->id));
  json_object_set_new(jreq, JSON_ACTION,
		      json_string(convert_action_to_str(req->action)));

  /* Three types of action: 1. Fetch, 2. Update and 3. Cancel */
  if (req->action == (ActionType)ACTION_FETCH) {

    Content *content = req->content;

    json_object_set_new(req, JSON_CALLBACK_URL, json_string(req->callbackURL));
    json_object_set_new(req, JSON_UPDATE_INTERVAL, json_integer(req->interval));

    /* Add Content object */
    json_object_set_new(jreq, JSON_CONTENT, json_object());
    jcont = json_object_get(req, JSON_CONTENT);

    json_object_set_new(jcont, JSON_NAME, json_string(content->name));
    json_object_set_new(jcont, JSON_TAG, json_string(content->tag));
    json_object_set_new(jcont, JSON_METHOD,
			json_string(convert_method_to_str(content->method)));
    json_object_set_new(jcont, JSON_PROVIDER_URL,
			json_string(content->providerURL));
  }

  if (req->action == (ActionType)ACTION_UPDATE) {

    if (req->callbackURL) {
      json_object_set_new(req, JSON_CALLBACK_URL,
			  json_string(req->callbackURL));
    }

    if (req->interval>0) {
      json_object_set_new(req, JSON_UPDATE_INTERVAL,
			  json_integer(req->interval));
    }
  }

  if (req->action == (ActionType)ACTION_CANCEL) {
    /* Do nothing. */
  }

  return TRUE;
}

/*
 * serialize_request_to_provider --
 *
 */

/* wimc_request -> {type:   "provider",
 *                  action: "request",
 *                  content { name : "name",
 *                            tag: "tag"}
 */
static int serialize_request_to_provider(WimcReq *req, json_t *json) {

  json_t *jreq, *jcont;
  Content *content;

  jreq = json_object_get(json, JSON_WIMC_REQUEST);

  if (req==NULL) {
    return FALSE;
  }

  content = req->content;

  json_object_set_new(jreq, JSON_TYPE,
		      json_string(convert_type_to_str(req->type)));
  json_object_set_new(jreq, JSON_ID, json_integer(req->id));
  json_object_set_new(jreq, JSON_ACTION,
		      json_string(convert_action_to_str(req->action)));

  /* Add Content object */
  json_object_set_new(jreq, JSON_CONTENT, json_object());
  jcont = json_object_get(req, JSON_CONTENT);

  json_object_set_new(jcont, JSON_NAME, json_string(content->name));
  json_object_set_new(jcont, JSON_TAG, json_string(content->tag));

  return TRUE;
}

/* agent_request -> { type: "update",
 *                  type_update: {
 *                      id: "same_id",
 *                      total_kbytes: "1234"
 *                      transfer_kbytes:  "34"
 *                      transfer_state: "fetch"
 *                      void_str: "some_string_"
 *			}

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

json_t *serialize_agent_request(AgentReq *request) {

  int ret=FALSE;
  json_t *json, *req, *content;

  json = json_object();
  if (json == NULL) {
    return NULL;
  }

  json_object_set_new(json, JSON_AGENT_REQUEST, json_object());

  req = json_object_get(json, JSON_AGENT_REQUEST);

  if (req==NULL) {
    return NULL;
  }

  if (request->type == (ReqType)REQ_REG) {
    ret = serialize_agent_request_register(request, json);
  } else if (request->type == (ReqType)REQ_UPDATE) {
    ret = serialize_agent_request_update(request, json);
  } else if (request->type == (ReqType)REQ_UNREG) {
    ret = serialize_agent_request_unregister(request, json);
  }

  if (ret) {
    return json;
  }

  return NULL;
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
int serialize_agent_request_register(AgentReq *req, json_t *json) {

  json_t *jreq, *jreg;
  Register *reg;

  jreq = json_object_get(json, JSON_AGENT_REQUEST);

  if (req==NULL || req->reg==NULL) {
    return FALSE;
  }

  reg = req->reg;

  json_object_set_new(req, JSON_TYPE, json_string(AGENT_REQ_TYPE_REG));

  /* Add register object */
  json_object_set_new(req, JSON_TYPE_REGISTER, json_object());
  jreg = json_object_get(req, JSON_TYPE_REGISTER);

  json_object_set_new(jreg, JSON_METHOD, json_string(reg->method));
  json_object_set_new(jreg, JSON_AGENT_URL, json_integer(reg->url));

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

int serialize_agent_request_update(AgentReq *req, json_t *json) {

  json_t *jupdate, *jreg;
  Update *update;

  jreg = json_object_get(json, JSON_AGENT_REQUEST);

  if (req==NULL || req->update==NULL) {
    return FALSE;
  }

  update = req->update;

  json_object_set_new(jreg, JSON_TYPE, json_string(AGENT_REQ_TYPE_UPDATE));

  /* Add update object */
  json_object_set_new(req, JSON_TYPE_UPDATE, json_object());
  jupdate = json_object_get(req, JSON_TYPE_UPDATE);

  json_object_set_new(jupdate, JSON_ID, json_integer(update->id));
  json_object_set_new(jupdate, JSON_TOTAL_KBYTES,
		      json_integer(update->totalKB));
  json_object_set_new(jupdate, JSON_TRANSFER_KBYTES,
		      json_integer_value(update->transferKB));
  json_object_set_new(jupdate, JSON_TRANSFER_STATE,
	       json_string_value(convert_state_to_str(update->transferState)));

  /* void str is non-zero only if there was an error otherwise is empty */
  if (update->transferState == (TransferState)ERR) {
    json_object_set_new(jupdate, JSON_VOID_STR,
			json_string_value(update->voidStr));
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

int serialize_agent_request_unregister(AgentReq *req, json_t *json) {

  json_t *jUnReg, *jreg;
  UnRegister *unReg;

  jreg = json_object_get(json, JSON_AGENT_REQUEST);

  if (req==NULL || req->unReg==NULL) {
    return FALSE;
  }

  unReg = req->unReg;

  json_object_set_new(jreg, JSON_TYPE, json_string(AGENT_REQ_TYPE_UNREG));

  /* Add un-register object */
  json_object_set_new(jreg, JSON_TYPE_UNREGISTER, json_object());
  jUnReg = json_object_get(jreg, JSON_TYPE_UNREGISTER);

  json_object_set_new(jUnReg, JSON_ID, json_integer(unReg->id));

  return TRUE;
}

/*
 * deserialize_agent_request --
 *
 */
int deserialize_agent_request(AgentReq *re, json_t *json) {


}
