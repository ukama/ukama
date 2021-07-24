/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <ulfius.h>
#include <string.h>
#include <curl/curl.h>
#include <curl/easy.h>
#include <jansson.h>

#include "agent.h"
#include "wimc.h"
#include "log.h"
#include "err.h"
#include "agent/jserdes.h"
#include "common/utils.h"

/*
 * Callback functions for various endpoints.
 */

static void free_wimc_request(WimcReq *req);
static int valid_req_type(WimcReq *req, MethodType method);
static int validate_post_request(WimcReq *req, MethodType method);


/*
 * decode a u_map into a string
 */

static char *print_map(const struct _u_map * map) {

  char * line, * to_return = NULL;
  const char **keys, * value;
  
  int len, i;
  
  if (map != NULL) {
    keys = u_map_enum_keys(map);
    for (i=0; keys[i] != NULL; i++) {
      value = u_map_get(map, keys[i]);
      len = snprintf(NULL, 0, "key is %s, value is %s", keys[i], value);
      line = o_malloc((len+1)*sizeof(char));
      snprintf(line, (len+1), "key is %s, value is %s", keys[i], value);
      if (to_return != NULL) {
        len = o_strlen(to_return) + o_strlen(line) + 1;
        to_return = o_realloc(to_return, (len+1)*sizeof(char));
        if (o_strlen(to_return) > 0) {
          strcat(to_return, "\n");
        }
      } else {
        to_return = o_malloc((o_strlen(line) + 1)*sizeof(char));
        to_return[0] = 0;
      }
      strcat(to_return, line);
      o_free(line);
    }
    return to_return;
  } else {
    return NULL;
  }
}

/*
 * valid_req_type --
 *
 */

static int valid_req_type(WimcReq *req, MethodType method) {

  if (req->type == (WReqType)WREQ_FETCH)
    return TRUE;

  if (req->type == (WReqType)WREQ_UPDATE)
    return TRUE;

  if (req->type == (WReqType)WREQ_CANCEL)
    return TRUE;

  return FALSE;
}

/*
 * validate_post_request -- validate all parameters of POST are valid. 
 *
 */

static int validate_post_request(WimcReq *req, MethodType method) {

  WFetch *fetch=NULL;
  WContent *content=NULL;

  if (!valid_req_type(req, method)){
    return WIMC_ERROR_BAD_TYPE;
  }

  if (req->type == (WReqType)WREQ_FETCH) {
    fetch = req->fetch;
  } else {
    goto done;
  }

  /* Id and interval are always positive. */
  if (uuid_is_null(fetch->uuid)) {
    return WIMC_ERROR_BAD_ID;
  }

  if (fetch->interval <= 0) {
    return WIMC_ERROR_BAD_INTERVAL;
  }

  if (validate_url(fetch->cbURL) != WIMC_OK) {
    return WIMC_ERROR_BAD_URL;
  }
      
  if (!fetch->content) {
    return WIMC_ERROR_MISSING_CONTENT;
  } else {
    content = fetch->content;
  }
  
  if (!content->name || !content->tag) {
    return WIMC_ERROR_BAD_NAME;
  }

  if (method == (MethodType)CHUNK) {
    if (validate_url(content->indexURL) != WIMC_OK) {
      return WIMC_ERROR_BAD_URL;
    }

    if (validate_url(content->storeURL) != WIMC_OK) {
      return WIMC_ERROR_BAD_URL;
    }
  } else {
    if (validate_url(content->providerURL) != WIMC_OK) {
      return WIMC_ERROR_BAD_URL;
    }
  }

 done:
  return WIMC_OK;
}

/*
 * agent_callback_delete --
 *
 */

int agent_callback_delete(const struct _u_request *request,
			  struct _u_response *response,
			  void *user_data) {
  
  char *post_params, *response_body;

  post_params = print_map(request->map_post_body);
  response_body = msprintf("OK!\n%s", post_params);
  
  ulfius_set_string_body_response(response, 200, response_body);
  o_free(response_body);
  o_free(post_params);
  
  return U_CALLBACK_CONTINUE;
}

/*
 * agent_callback_get --
 *
 */

int agent_callback_get(const struct _u_request *request,
		       struct _u_response *response,
		       void *user_data) {
  
  char *post_params, *response_body;

  post_params = print_map(request->map_post_body);
  response_body = msprintf("OK!\n%s", post_params);
  
  ulfius_set_string_body_response(response, 200, response_body);
  o_free(response_body);
  o_free(post_params);

  return U_CALLBACK_CONTINUE;
}

/*
 * agent_callback_put --
 *
 */
int agent_callback_put(const struct _u_request *request,
		       struct _u_response *response,
		       void *user_data) {
  
  char *post_params, *response_body;

  post_params = print_map(request->map_post_body);
  response_body = msprintf("OK!\n%s", post_params);

  ulfius_set_string_body_response(response, 200, response_body);
  o_free(response_body);
  o_free(post_params);

  return U_CALLBACK_CONTINUE;
}

/*
 * agent_callback_post -- callback to handle WIMC request to fetch new content.
 *
 */
int agent_callback_post(const struct _u_request *request,
			struct _u_response *response, void *user_data) {

  int ret, retCode;
  char *resBody=NULL;
  json_t *jreq=NULL;
  json_error_t jerr;
  MethodType *method = (MethodType *)user_data;
  WimcReq *req=NULL;

  jreq = ulfius_get_json_body_request(request, &jerr);
  if (!jreq) {
    log_error("json error: %s", jerr.text);
    goto done;
  }

  req = (WimcReq *)calloc(1, sizeof(WimcReq));

  ret = deserialize_wimc_request(&req, jreq);

  if (ret) {
    char *jStr;

    jStr = json_dumps(jreq, 0);
    if (jStr) {
      log_debug("Wimc request received str: %s", jStr);
      free(jStr);
    }
  }

  ret = validate_post_request(req, *method);
  if (ret != WIMC_OK) {
    retCode = 400;
    resBody = msprintf("%s\n", error_to_str(ret));
    goto done;
  } else {
    retCode = 200;
    resBody = msprintf("OK");
  }

  request_handler(req->fetch);
  
 done:

  ulfius_set_string_body_response(response, retCode, resBody);
  free_wimc_request(req);
  o_free(resBody);
  json_decref(jreq);

  return U_CALLBACK_CONTINUE;
}

/*
 * free_wimc_request --
 *
 */

static void free_wimc_request(WimcReq *req) {

  if (!req)
    return;
  
  if (req->type == WREQ_FETCH) {
    WContent *content;
    
    free(req->fetch->cbURL);
    content = req->fetch->content;

    if (content) {
      free(content->name);
      free(content->tag);
      free(content->method);
      free(content->providerURL);
      free(content->indexURL);
      free(content->storeURL);
      free(content);
    }

    free(req->fetch);
  }
  
  free(req);
}  

/*
 * agent_callback_stats 
 *
 */
int agent_callback_stats(const struct _u_request *request,
			 struct _u_response *response,
			 void *user_data) {
  
  char *post_params = print_map(request->map_post_body);
  char *response_body = msprintf("OK!\n%s", post_params);
  
  ulfius_set_string_body_response(response, 200, response_body);
  o_free(response_body);
  o_free(post_params);
  
  return U_CALLBACK_CONTINUE;
}

/*
 * agent_callback_default -- default callback for no-match
 *
 */
int agent_callback_default(const struct _u_request *request,
			   struct _u_response *response, void *user_data) {

  ulfius_set_string_body_response(response, 404, "What's wrong?\n");
  return U_CALLBACK_CONTINUE;
}
