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

/*
 * Callback functions for various endpoints.
 */

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
 * log_request -- log various parameters for the incoming request. 
 *
 */

static void log_request(const struct _u_request *request) {

  log_debug("Recevied: %s %s %s", request->http_protocol, request->http_verb,
	    request->http_url);
}

/*
 * process_post_request --handle post request from WIMC.d
 *
 */

static int process_post_request(WimcReq *req) {

}

/*
 * validate_post_request -- validate all parameters of POST are valid. 
 */
static int validate_post_request(WimcReq *req, MethodType *method) {

  Content *content;

  if (req->type != (WReqType)AGENT) {
    return WIMC_ERROR_BAD_TYPE;
  }

  /* Only Fetch is via the POST. update and delete is via put and
   * delete respectively.
   */
  if (req->action != (ActionType)ACTION_FETCH) {
    return WIMC_ERROR_BAD_ACTION;
  }

  /* Id and interval are always positive. */
  if (req->id <= 0) { 
    return WIMC_ERROR_BAD_ID;
  }

  if (req->interval <= 0){
    return WIMC_ERROR_BAD_INTERVAL;
  }
      
  if (!req->content) {
    return WIMC_ERROR_MISSING_CONTENT;
  } else {
    content = req->content;
  }
  
  if (!content->name || !content->tag) {
    return WIMC_ERROR_BAD_NAME;
  }

  if (content->method != (MethodType)ACTION_FETCH) {
    log_error("Invalid method recevied for the POST. Ignoring");
    return WIMC_ERROR_BAD_METHOD;
  }

  if (!req->callbackURL || !content->providerURL) {
    return WIMC_ERROR_BAD_URL;
  }

  /* check if URL are reachable by the agent. */
  if (validate_url(req->callbackURL)) {
    return WIMC_ERROR_BAD_URL;
  }

  if (validate_url(content->providerURL)) {
    return WIMC_ERROR_BAD_URL;
  }
  
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
  char *params, *resBody;
  json_t *jreq=NULL;
  MethodType *method = (MethodType *)user_data;
  WimcReq req;

  jreq = ulfius_get_json_body_response(request, NULL);
  deserialize_wimc_request(&req, jreq);

  ret = validate_post_request(&req, method);

  if (ret == WIMC_OK) {
    //    ret = process_post_request(agents, &req);
    // XXX
  }
  
  if (ret == WIMC_OK) {
    retCode = 200;
  } else {
    retCode = 400;
  }

  params = print_map(request->map_post_body);
  resBody = msprintf("%s\n%s", error_to_str(ret), params);

  ulfius_set_string_body_response(response, retCode, resBody);
  o_free(resBody);
  o_free(params);

  return U_CALLBACK_CONTINUE;
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
