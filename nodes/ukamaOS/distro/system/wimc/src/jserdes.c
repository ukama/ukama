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
#include "common/utils.h"

#include "usys_error.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"
#include "usys_types.h"

void json_log(json_t *json) {

    char *str = NULL;

    str = json_dumps(json, 0);
    if (str) {
        log_debug("json str: %s", str);
        free(str);
    }
}

static bool get_json_entry(json_t *json, char *key, json_type type,
                           char **strValue, int *intValue,
                           double *doubleValue) {

    json_t *jEntry=NULL;

    if (json == NULL || key == NULL) return USYS_FALSE;

    jEntry = json_object_get(json, key);
    if (jEntry == NULL) {
        log_error("Missing %s key in json", key);
        return USYS_FALSE;
    }

    switch(type) {
    case (JSON_STRING):
        *strValue = strdup(json_string_value(jEntry));
        break;
    case (JSON_INTEGER):
        *intValue = json_integer_value(jEntry);
        break;
    case (JSON_REAL):
        *doubleValue = json_real_value(jEntry);
        break;
    default:
        log_error("Invalid type for json key-value: %d", type);
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

bool serialize_wimc_request_fetch(WimcReq *req, json_t **json) {

    json_t *jfetch=NULL, *jcontent=NULL;;
    WFetch *fetch=NULL;
    WContent *content=NULL;
    char idStr[36+1]; /* 36-bytes for UUID + trailing '\0' */

    if (req == NULL && req->fetch == NULL) {
        return USYS_FALSE;
    }

    fetch = req->fetch;
    content = fetch->content;

    if (content == NULL) return USYS_FALSE;

    json_object_set_new(*json, JSON_TYPE, json_string(WIMC_REQ_TYPE_FETCH));

    /* Add fetch object */
    json_object_set_new(*json, JSON_TYPE_FETCH, json_object());
    jfetch = json_object_get(*json, JSON_TYPE_FETCH);

    uuid_unparse(fetch->uuid, idStr);
    json_object_set_new(jfetch, JSON_ID, json_string(idStr));
    json_object_set_new(jfetch, JSON_UPDATE_INTERVAL,
                        json_integer(fetch->interval));

    json_object_set_new(jfetch, JSON_CONTENT, json_object());
    jcontent = json_object_get(jfetch, JSON_CONTENT);

    /* Add content object. */
    json_object_set_new(jcontent, JSON_NAME,
                        json_string(content->name));
    json_object_set_new(jcontent, JSON_TAG,
                        json_string(content->tag));
    json_object_set_new(jcontent, JSON_METHOD,
                        json_string(content->method));
    json_object_set_new(jcontent, JSON_INDEX_URL,
                        json_string(content->indexURL));
    json_object_set_new(jcontent, JSON_STORE_URL,
                        json_string(content->storeURL));

    return USYS_TRUE;
}

bool serialize_wimc_request(WimcReq *request, json_t **json) {

    json_t *req=NULL;

    *json = json_object();
    if (*json == NULL) return USYS_FALSE;

    json_object_set_new(*json, JSON_WIMC_REQUEST, json_object());
    req = json_object_get(*json, JSON_WIMC_REQUEST);
    if (req == NULL) return USYS_FALSE;

    if (request->type == (WReqType)WREQ_FETCH) {
        return serialize_wimc_request_fetch(request, &req);
    } else if (request->type == (WReqType)WREQ_UPDATE) {
        
    }


    return USYS_FALSE;
}

bool deserialize_agent_request_update(Update **update, json_t *json) {

    json_t *jupdate, *obj;

    jupdate = json_object_get(json, JSON_TYPE_UPDATE);

    if (jupdate == NULL) return USYS_FALSE;

    *update = (Update *)calloc(sizeof(Update), 1);
    if (*update == NULL) return USYS_FALSE;

    /* All updates must have the ID. */
    obj = json_object_get(jupdate, JSON_ID);
    if (obj == NULL) {
        return USYS_FALSE;
    } else {
        uuid_parse(json_string_value(obj), (*update)->uuid);
    }

    /* Total data to be transfered as part of this fetch (in KB) */
    obj = json_object_get(jupdate, JSON_TOTAL_KBYTES);
    if (obj) {
        (*update)->totalKB = json_integer_value(obj);
    }

    /* Activity so far. */
    obj = json_object_get(jupdate, JSON_TRANSFER_KBYTES);
    if (obj) {
        (*update)->transferKB = json_integer_value(obj);
    }

    /* Activity state. */
    obj = json_object_get(jupdate, JSON_TRANSFER_STATE);
    if (obj) {
        (*update)->transferState =
            convert_str_to_tx_state(json_string_value(obj));
        obj = json_object_get(jupdate, JSON_VOID_STR);
        if (obj) {
            (*update)->voidStr = strdup(json_string_value(obj));
        }
    }

    return USYS_TRUE;
}

int deserialize_provider_response(ServiceURL **urls, int *counter,
				  json_t *json) {

  int i=0, j=0, count=0;
  json_t *root=NULL, *jArray=NULL;
  json_t *elem=NULL, *method=NULL, *url=NULL;
  json_t *iURL=NULL, *sURL=NULL;

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
      iURL   = json_object_get(elem, JSON_INDEX_URL);
      sURL   = json_object_get(elem, JSON_STORE_URL);

      if (method && url && iURL && sURL) {
	urls[i]->method = strdup(json_string_value(method));
	urls[i]->url    = strdup(json_string_value(url));
	/* index and store URL will be strlen=0 if method!=chunk */
	urls[i]->iURL   = strdup(json_string_value(iURL));
	urls[i]->sURL   = strdup(json_string_value(sURL));
      }
    }
  }

  return TRUE;
  
 failure:
  for (j=0; j<i; j++) {
    free(urls[j]->method);
    free(urls[j]->url);
    free(urls[j]->iURL);
    free(urls[j]->sURL);
  }
  
  if (*urls)
    free(*urls);
  
  return FALSE;
}

/*
 * deserialize_hub_response --
 *
 */
int deserialize_hub_response(Artifact ***artifacts, int *counter,
			     json_t *json) {

    int i=0, j=0, count=0, formatsCount=0;

    json_t *name=NULL;
    json_t *jArray=NULL, *formatsArray=NULL, *elem=NULL;
    json_t *url=NULL, *additionalInfo=NULL;
    json_t *version=NULL, *createdAt=NULL;
    json_t *size=NULL, *type=NULL, *extraInfo=NULL;
    json_t *formatElem=NULL, *value=NULL;
    
    ArtifactFormat *formats = NULL;

    /* sanity check */
    if (!json) return USYS_FALSE;

    usys_log_debug("Deserializing hub response ...");
    json_log(json);

    *counter=0;

    name = json_object_get(json, JSON_NAME);
    if (name == NULL) return USYS_FALSE;

    jArray = json_object_get(json, JSON_ARTIFACTS);
    if (!jArray) return USYS_FALSE;

    if (json_is_array(jArray)) {

        count = json_array_size(jArray);
        *counter = count;

        if (count == 0) return USYS_TRUE;

        *artifacts = (Artifact **)calloc(count, sizeof(Artifact *));
        if (*artifacts == NULL) {
            *counter = 0;
            return USYS_FALSE;
        }

        for (i=0; i<count; i++) {
            elem = json_array_get(jArray, i);
            
            if (elem == NULL) {
                goto failure;
            }

            (*artifacts)[i] = (Artifact *)calloc(1, sizeof(Artifact));

            /* name */
            (*artifacts[i])->name = strdup(json_string_value(name));

            /* version */
            version = json_object_get(elem, JSON_VERSION);
            if (version == NULL) goto failure;
            (*artifacts)[i]->version = strdup(json_string_value(version));

            /* formats */
            formatsArray = json_object_get(elem, JSON_FORMATS);
            if (!formatsArray)                goto failure;
            if (!json_is_array(formatsArray)) goto failure;


            formatsCount = json_array_size(formatsArray);
            (*artifacts)[i]->formatsCount = formatsCount;

            (*artifacts)[i]->formats = (ArtifactFormat **)
                calloc(formatsCount, sizeof(ArtifactFormat *));
            if ((*artifacts)[i]->formats == NULL) {
                *counter = 0;
                goto failure;
            }

            for (j=0; j<formatsCount; j++) {

                (*artifacts)[i]->formats[j] = (ArtifactFormat *)
                    calloc(1, sizeof(ArtifactFormat));
                formats = (*artifacts)[i]->formats[j];

                formatElem = json_array_get(formatsArray, j);

                type      = json_object_get(formatElem, JSON_TYPE);
                extraInfo = json_object_get(formatElem, JSON_EXTRA_INFO);
                url       = json_object_get(formatElem, JSON_URL);
                createdAt = json_object_get(formatElem, JSON_CREATED_AT);
                size      = json_object_get(formatElem, JSON_SIZE_BYTES);

                if (type && url && createdAt) {
                    formats->type = strdup(json_string_value(type));
                    formats->url  = strdup(json_string_value(url));
                    formats->createdAt = strdup(json_string_value(createdAt));
                } else {
                    goto failure;
                }

                if (strcmp(formats->type, WIMC_METHOD_CHUNK_STR)==0) {
                    formats->size = 0;
                    if (extraInfo) {
                        value = json_object_get(extraInfo, JSON_CHUNKS);
                        if (value==NULL) goto failure;
                        formats->extraInfo = strdup(json_string_value(value));
                    } else goto failure;
                } else {
                    if (size==NULL) goto failure;
                    else formats->size = (int)json_integer_value(size);
                }
            }
        }
    }
    
    return USYS_TRUE;

failure:
    usys_log_error("Error deserializing the hub response");
    *counter = 0;
    for (i=0; i<count; i++) {
        usys_free((*artifacts)[i]->name);
        usys_free((*artifacts)[i]->version);
        
        for (j=0; j<(*artifacts)[i]->formatsCount; j++) {
            
            formats = (*artifacts)[i]->formats[j];
            
            usys_free(formats->type);
            usys_free(formats->url);
            usys_free(formats->extraInfo);
            usys_free(formats);
        }
        usys_free((*artifacts)[i]->formats);
        usys_free((*artifacts)[i]);
    }
    
    usys_free(*artifacts);

    *artifacts = NULL;
    *counter = 0;

    return USYS_FALSE;
}

int serialize_task(WTasks *task, json_t **json) {

  int ret=FALSE;
  json_t *jtask=NULL, *jresp=NULL;
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

  json_object_set_new(*json, JSON_WIMC_RESPONSE, json_object());
  jresp = json_object_get(*json, JSON_WIMC_RESPONSE);
  if (jresp==NULL) {
    return ret;
  }

  json_object_set_new(jresp, JSON_TYPE, json_string(WIMC_RESP_TYPE_STATUS));
  json_object_set_new(*json, JSON_STATUS, json_object());
  jtask = json_object_get(*json, JSON_STATUS);

  if (jtask==NULL) {
    return ret;
  }

  uuid_unparse(task->uuid, &idStr[0]);
  json_object_set_new(jtask, JSON_ID, json_string(idStr));
  json_object_set_new(jtask, JSON_NAME, json_string(content->name));
  json_object_set_new(jtask, JSON_TAG, json_string(content->tag));
  json_object_set_new(jtask, JSON_METHOD, json_string(content->method));

  json_object_set_new(jtask, JSON_TOTAL_KBYTES, json_integer(update->totalKB));
  json_object_set_new(jtask, JSON_TRANSFER_KBYTES,
		      json_integer(update->transferKB));
  json_object_set_new(jtask, JSON_TRANSFER_STATE,
		      json_string(convert_tx_state_to_str(update->transferState)));

  if (update->voidStr) {
    json_object_set_new(jtask, JSON_VOID_STR, json_string(update->voidStr));
  } else {
    json_object_set_new(jtask, JSON_VOID_STR, json_string(""));
  }

  return TRUE;
}

/*
 * serialize_result -- Serialize result string as wimc_response{}.
 *
 */
int serialize_result(WRespType type, char *str, json_t **json) {

  int ret=FALSE;
  json_t *jresp=NULL;

  /* Go ahead, serialize them all objects. */
  *json = json_object();
  if (*json == NULL) {
    return ret;
  }

  json_object_set_new(*json, JSON_WIMC_RESPONSE, json_object());
  jresp = json_object_get(*json, JSON_WIMC_RESPONSE);
  if (jresp==NULL) {
    return ret;
  }

  if (type == (WRespType)WRESP_RESULT) {
    json_object_set_new(jresp, JSON_TYPE, json_string(WIMC_RESP_TYPE_RESULT));
  } else if (type == (WRespType)WRESP_PROCESSING) {
    json_object_set_new(jresp, JSON_TYPE,
			json_string(WIMC_RESP_TYPE_PROCESSING));
  } else if (type == (WRespType)WRESP_ERROR) {
    json_object_set_new(jresp, JSON_TYPE, json_string(WIMC_RESP_TYPE_ERROR));
  }

  if (str != NULL) {
    json_object_set_new(jresp, JSON_VOID_STR, json_string(str));
  } else {
    json_object_set_new(jresp, JSON_VOID_STR, json_string(""));
  }

  return TRUE;
}

static void log_json(json_t *json) {

  char *str = NULL;

  str = json_dumps(json, 0);
  if (str) {
    log_debug("json str: %s", str);
    free(str);
  }
}
