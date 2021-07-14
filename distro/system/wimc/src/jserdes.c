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

/*
 * serialize_wimc_request --
 *
 */
int serialize_wimc_request(WimcReq *request, json_t **json) {

  int ret=FALSE;
  json_t *req=NULL;

  *json = json_object();
  if (*json == NULL)
    return ret;

  json_object_set_new(*json, JSON_WIMC_REQUEST, json_object());
  req = json_object_get(*json, JSON_WIMC_REQUEST);

  if (req==NULL) {
    return ret;
  }

  if (request->type == (WReqType)WREQ_FETCH) {
    ret = serialize_wimc_request_fetch(request, &req);
  } else if (request->type == (WReqType)WREQ_UPDATE) {

  }

  if (ret) {
    char *str;
    str = json_dumps(*json, 0);

    if (str) {
      log_debug("Wimc request str: %s", str);
      free(str);
    }
    ret = TRUE;
  }

  return ret;
}

/*
 * serialize_wimc_request_fetch --
 *
 */
static int serialize_wimc_request_fetch(WimcReq *req, json_t **json) {

  json_t *jfetch=NULL, *jcontent=NULL;;
  WFetch *fetch=NULL;
  WContent *content=NULL;
  char idStr[36+1]; /* 36-bytes for UUID + trailing '\0' */

  if (req==NULL && req->fetch==NULL) {
    return FALSE;
  }

  fetch = req->fetch;

  if (fetch->content==NULL) {
    return FALSE;
  }

  content = fetch->content;

  json_object_set_new(*json, JSON_TYPE, json_string(WIMC_REQ_TYPE_FETCH));

  /* Add fetch object */
  json_object_set_new(*json, JSON_TYPE_FETCH, json_object());
  jfetch = json_object_get(*json, JSON_TYPE_FETCH);

  uuid_unparse(fetch->uuid, idStr);
  json_object_set_new(jfetch, JSON_ID, json_string(idStr));
  json_object_set_new(jfetch, JSON_CALLBACK_URL, json_string(fetch->cbURL));
  json_object_set_new(jfetch, JSON_UPDATE_INTERVAL,
		      json_integer(fetch->interval));

  json_object_set_new(jfetch, JSON_CONTENT, json_object());
  jcontent = json_object_get(jfetch, JSON_CONTENT);

  /* Add content object. */
  json_object_set_new(jcontent, JSON_NAME, json_string(content->name));
  json_object_set_new(jcontent, JSON_TAG, json_string(content->tag));
  json_object_set_new(jcontent, JSON_METHOD, json_string(content->method));
  json_object_set_new(jcontent, JSON_PROVIDER_URL,
		      json_string(content->providerURL));

  return TRUE;
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
 * deserialize_agent_request_register --
 *
 */
static int deserialize_agent_request_register(Register **reg, json_t *json) {

  json_t *jreg, *jmethod, *jurl;

  jreg = json_object_get(json, JSON_TYPE_REGISTER);

  if (jreg == NULL) {
    return FALSE;
  }

  *reg = (Register *)calloc(sizeof(Register), 1);
  if (reg == NULL) {
    return FALSE;
  }

  jmethod = json_object_get(jreg, JSON_METHOD);
  jurl = json_object_get(jreg, JSON_AGENT_URL);
  if (jurl == NULL || jmethod == NULL) {
    return FALSE;
  }

  (*reg)->method = strdup(json_string_value(jmethod));
  (*reg)->url = strdup(json_string_value(jurl));

  return TRUE;
}

/*
 * deserialize_agent_request_update --
 */
static int deserialize_agent_request_update(Update *update, json_t *json) {

  json_t *jupdate, *obj;

  jupdate = json_object_get(json, JSON_TYPE_UPDATE);

  if (jupdate==NULL) {
    return FALSE;
  }

  update = (Update *)calloc(sizeof(Update), 1);
  if (update == NULL) {
    return FALSE;
  }

  /* All updates must have the ID. */
  obj = json_object_get(jupdate, JSON_ID);
  if (obj == NULL) {
    return FALSE;
  } else {
    uuid_unparse(json_string_value(obj), update->uuid);
  }

  /* Total data to be transfered as part of this fetch (in KB) */
  obj = json_object_get(jupdate, JSON_TOTAL_KBYTES);
  if (obj) {
    update->totalKB = json_integer_value(obj);
  }

  /* Activity so far. */
  obj = json_object_get(jupdate, JSON_TRANSFER_KBYTES);
  if (obj) {
    update->transferKB = json_integer_value(obj);
  }

  /* Activity state. */
  obj = json_object_get(jupdate, JSON_TRANSFER_STATE);
  if (obj) {
    update->transferState = convert_str_to_state(json_string_value(obj));

    if (update->transferState == (TransferState)ERR ||
	update->transferState == (TransferState)DONE) {
      obj = json_object_get(jupdate, JSON_VOID_STR);
      if (obj) {
	update->voidStr = json_string_value(obj);
      }
    }
  }

  return TRUE;
}

/*
 * deserialize_agent_request_unreg --
 */
static int deserialize_agent_request_unreg(UnRegister *unReg, json_t *json) {

  json_t *jreq, *obj;

  jreq = json_object_get(json, JSON_TYPE_UNREGISTER);

  if (jreq == NULL) {
    return FALSE;
  }

  unReg = (UnRegister *)calloc(sizeof(UnRegister), 1);
  if (unReg == NULL) {
    return FALSE;
  }

  obj = json_object_get(jreq, JSON_ID);
  if (obj) {
    uuid_unparse(json_string_value(obj), unReg->uuid);
  } else {
    return FALSE;
  }

  return TRUE;
}

/*
 * deserialize_agent_request --
 *
 */
int deserialize_agent_request(AgentReq **request, json_t *json) {

  int ret=FALSE;
  char *str;
  json_t *jreq, *jtype;

  AgentReq *req = *request;
  
  if (!json) {
    return FALSE;
  }
  
  jreq = json_object_get(json, JSON_AGENT_REQUEST);
  if (jreq == NULL) {
    return FALSE;
  }

  jtype = json_object_get(jreq, JSON_TYPE);

  if (jtype==NULL) {
    return FALSE;
  }

  req->type = convert_str_to_type(json_string_value(jtype));

  if (req->type == (ReqType)REQ_REG) {
    ret = deserialize_agent_request_register(&req->reg, jreq);
  } else if (req->type == (ReqType)REQ_UPDATE) {
    ret = deserialize_agent_request_update(req->update, jreq);
  } else if (req->type == (ReqType)REQ_UNREG) {
    ret = deserialize_agent_request_unreg(req->unReg, jreq);
  }

  return ret;
}

static int deserialize_wimc_request_to_provider(WimcReq *req, json_t *json) {

  return TRUE;
}

#if 0
// XXX Move this code to provider. 
/*
 * deserialize_wimc_request --
 *
 */
int deserialize_wimc_request(WimcReq *req, json_t *json) {

  int ret=FALSE;
  json_t *jreq, *jtype;
  WReqType type;

  jreq = json_object_get(json, JSON_WIMC_REQUEST);

  if (jreq==NULL) {
    return FALSE;
  }

  jtype = json_object_get(jreq, JSON_TYPE);
  if (jtype == NULL) {
    return FALSE;
  }

  type = convert_str_to_wType(json_string_value(jtype));
  req->type = type;

  if (type == (WReqType)AGENT) {
    ret = deserialize_wimc_request_to_agent(req, jreq);
  } else if (type == (WReqType)PROVIDER) {
    ret = deserialize_wimc_request_to_provider(req, jreq);
  }

  return ret;
}
#endif

/*
 * deserialize_provider_response --
 *
 */

int deserialize_provider_response(ServiceURL **urls, int *counter,
				  json_t *json) {

  int i=0, j=0, count=0;
  json_t *root=NULL, *jArray=NULL;
  json_t *elem=NULL, *method=NULL, *url=NULL;

  /* sanity check */
  if (!json) return FALSE;

  root   = json_object_get(json, JSON_PROVIDER_RESPONSE);
  jArray = json_object_get(root, JSON_AGENT_CB);

  if (json_is_array(jArray)) {
    
    count = json_array_size(jArray);
    *counter = count;

    if (count==0) { /* No match found. */
      return TRUE;
    }
    
    *urls = (ServiceURL *)calloc(sizeof(ServiceURL), count);
    
    for (i=0; i<count; i++) {
      elem = json_array_get(jArray, i);

      if (elem == NULL) {
	goto failure;
      }
      
      method = json_object_get(elem, JSON_METHOD);
      url    = json_object_get(elem, JSON_URL);

      if (method && url) {
	urls[i]->method = strdup(json_string_value(method));
	urls[i]->url    = strdup(json_string_value(url));
      }
    }
  }

  return TRUE;
  
 failure:
  for (j=0; j<i; j++) {
    free(urls[j]->method);
    free(urls[j]->url);
  }
  
  if (*urls)
    free(*urls);
  
  return FALSE;
}

/*
 * serialize_task -- Serialize task struct.
 *
 */
int serialize_task(WTasks *task, json_t **json) {

  int ret=FALSE;
  json_t *jtask=NULL;
  WContent *content=NULL;
  Update *update=NULL;
  char idStr[36+1];

  /* Sanity check. */
  if (task==NULL) {
    return ret;
  }

  content = task->content;
  update  = task->update;
  if (uuid_is_null(task->uuid) && content==NULL && update == NULL){
    return ret;
  }

  /* Go ahead, serialize them all objects. */
  *json = json_object();
  if (*json == NULL) {
    return ret;
  }

  json_object_set_new(*json, JSON_TASK, json_object());
  jtask = json_object_get(*json, JSON_TASK);

  if (jtask==NULL) {
    return ret;
  }

  uuid_unparse(task->uuid, &idStr[0]);
  json_object_set_new(jtask, JSON_ID, json_string(idStr));
  json_object_set_new(jtask, JSON_NAME, json_string(content->name));
  json_object_set_new(jtask, JSON_TAG, json_string(content->tag));
  json_object_set_new(jtask, JSON_METHOD,
		      json_string(convert_method_to_str(content->method)));

  json_object_set_new(jtask, JSON_TOTAL_KBYTES, json_integer(update->totalKB));
  json_object_set_new(jtask, JSON_TRANSFER_KBYTES,
		      json_integer(update->transferKB));
  json_object_set_new(jtask, JSON_TRANSFER_STATE,
		      json_string(convert_state_to_str(update->transferState)));

  if (update->voidStr) {
    json_object_set_new(jtask, JSON_VOID_STR, json_string(update->voidStr));
  } else {
    json_object_set_new(jtask, JSON_VOID_STR, json_string(""));
  }

  return TRUE;
}
