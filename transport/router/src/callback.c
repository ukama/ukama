/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Callback functions for various endpoints and REST methods.
 */

#include <ulfius.h>
#include <curl/curl.h>
#include <uuid/uuid.h>
#include <string.h>

#include "router.h"
#include "callback.h"
#include "jserdes.h"
#include "log.h"
#include "httpStatus.h"
#include "pattern.h"
#include "forward.h"

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
 * log_json_params --
 *
 */
static void log_json_params(json_t *json, char *type) {

  char *str = NULL;

  str = json_dumps(json, 0);
  if (str) {
    log_debug("JSON %s str: %s", type, str);
    free(str);
  }
}

/*
 * add_service_entry --
 *
 */
static int add_service_entry(Router **router, Service *service) {

  Service *ptr=NULL;

  if ((*router)->services == NULL) {
    (*router)->services = (Service *)calloc(1, sizeof(Service));
    if ((*router)->services == NULL) {
      log_error("Error allocating memory of size: %lu", sizeof(Service));
      return FALSE;
    }
    ptr = (*router)->services;
  } else {
    for (ptr=(*router)->services; ptr->next; ptr=ptr->next);
    ptr->next = (Service *)calloc(1, sizeof(Service));
    if (ptr->next == NULL) {
      log_error("Error allocating memory of size: %lu", sizeof(Service));
      return FALSE;
    }
    ptr = ptr->next;
  }

  uuid_generate(service->uuid);

  uuid_copy(ptr->uuid, service->uuid);
  ptr->pattern = service->pattern;
  ptr->forward = service->forward;
  ptr->next    = NULL;

  return TRUE;
}

/*
 * parse_request_params --
 *
 */
static int parse_request_params(struct _u_map * map, Pattern **pattern) {

  const char **keys, *value;
  int count, i;
  Pattern *ptr=NULL;

  if (map == NULL) return FALSE;

  count = u_map_count(map);

  if (count == 0) {
    log_error("No key-value pair in the header");
    return FALSE;
  }

  if (*pattern == NULL) {
    *pattern = (Pattern *)calloc(1, sizeof(Pattern));
    if (*pattern == NULL) {
      log_error("Error allocating memory of size: %d", sizeof(Pattern));
      return FALSE;
    }
  }

  ptr = *pattern;

  keys = u_map_enum_keys(map);
  for (i=0; keys[i] != NULL; i++) {
    value = u_map_get(map, keys[i]);

    if (ptr == NULL) {
      ptr = (Pattern *)calloc(1, sizeof(Pattern));
      if (ptr == NULL) {
	log_error("Error allocating memory of size: %d", sizeof(Pattern));
	goto failure;
      }
    }

    ptr->key   = strdup(keys[i]);
    ptr->value = strdup(value);

    ptr = ptr->next;
  }

  return TRUE;

 failure:
  ptr = *pattern;
  while (ptr) {
    if (ptr->key)   free(ptr->key);
    if (ptr->value) free(ptr->value);

    ptr = ptr->next;
  }

  if (*pattern) free(*pattern);
  *pattern = NULL;

  return FALSE;
}

/*
 * free_service --
 *
 */
static void free_service(Service *service) {

  Pattern *ptr=NULL, *tmp=NULL;
  Forward *fPtr=NULL;

  if (service == NULL) return;

  ptr  = service->pattern;
  fPtr = service->forward;

  if (fPtr) {
    if (fPtr->ip)   free(fPtr->ip);
    if (fPtr->port) free(fPtr->port);
    free(fPtr);
  }

  while (ptr) {
    if (ptr->key)   free(ptr->key);
    if (ptr->value) free(ptr->value);
    tmp = ptr->next;
    free(ptr);
    ptr = tmp;
  }

  free(service);
}

/*
 * callback_get_route --
 *
 */
int callback_get_route(const struct _u_request *request,
		       struct _u_response *response,
		       void *user_data) {

  ulfius_set_string_body_response(response, 200, "");

  return U_CALLBACK_CONTINUE;
}

/*
 * callback_post_route --
 *
 */
int callback_post_route(const struct _u_request *request,
			struct _u_response *response,
			void *user_data) {

  int retCode;
  json_t *jreq=NULL;
  json_error_t jerr;
  Service *service=NULL;
  Router *router=NULL;
  json_t *jResp=NULL;
  char *jRespStr=NULL;
  const char *statusStr=NULL;

  router = (Router *)user_data;

  log_request(request);
  
  /* get json body */
  jreq = ulfius_get_json_body_request(request, &jerr);
  if (!jreq) {
    log_error("json error for POST %s: %s", EP_ROUTE, jerr.text);
    retCode   = HttpStatus_BadRequest;
    statusStr = HttpStatusStr(retCode);
    log_error("%d: %s", retCode, statusStr);
    goto reply;
  } else {
    deserialize_post_route_request(&service, jreq);
  }

  log_json_params(jreq, "request");

  if (!service) {
    retCode   = HttpStatus_BadRequest;
    statusStr = HttpStatusStr(retCode);
    goto reply;
  }

  /* Validate the connection with forward service */
  if (valid_forward_route(service->forward->ip,
			  service->forward->port) != TRUE) {
    retCode   = HttpStatus_ServiceUnavailable;
    statusStr = HttpStatusStr(retCode);
    log_error("Matching forward service unavailable. %d: %s", retCode,
	      statusStr);
    goto reply;
  }

  /* Add to internal structure. UUID is assigned. */
  add_service_entry(&router, service);

  /* Reply back with uuid */
  serialize_post_route_response(&jResp, UUID, (void *)&(service->uuid), NULL);
  retCode   = HttpStatus_OK;
  jRespStr  = json_dumps(jResp, 0);
  statusStr = jRespStr;

 reply:
  ulfius_set_string_body_response(response, retCode, statusStr);

  log_debug("Registration response: %d %s", retCode, statusStr);

  if (retCode == HttpStatus_OK) free(jRespStr);

  if (retCode != HttpStatus_OK) {
    free_service(service);
  }

  json_decref(jResp);
  json_decref(jreq);

  return U_CALLBACK_CONTINUE;
}

/*
 * callback_get_stats --
 *
 */
int callback_get_stats(const struct _u_request *request,
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
 * callback_post_service --
 *
 */
int callback_post_service(const struct _u_request *request,
			  struct _u_response *response,
			  void *user_data) {

  int retCode, serviceResp;
  const char *statusStr=NULL;
  Router *router=NULL;
  Pattern *requestPattern=NULL;
  Forward *requestForward=NULL;
  struct _u_request  *fRequest=NULL;
  struct _u_response *fResponse=NULL;

  router = (Router *)user_data;

  log_request(request);

  /* Steps are:
   * 1. parse the header for requested param
   * 2. pattern matches against the existing registered services
   * 3. setup request to forward
   * 4. open connection to the service, forward request and recv respone
   * 5. forward service response to client
   */

  /* Step-1: Parse the requested pattern from the header */
  if (!parse_request_params(request->map_url, &requestPattern)) {
    retCode   = HttpStatus_BadRequest;
    statusStr = HttpStatusStr(retCode);
    log_error("%d: %s", retCode, statusStr);
    goto reply;
  } else {
    log_debug("Recevied forward request: %s", print_map(request->map_url));
  }

  /* Step-2: Pattern match to a service (if any)*/
  if (!find_matching_service(router, requestPattern, &requestForward)) {
    retCode   = HttpStatus_ServiceUnavailable;
    statusStr = HttpStatusStr(retCode);
    log_error("No matching forward service found. %d: %s", retCode, statusStr);
    goto reply;
  } else {
    log_debug("Matching service found at IP: %s port: %s",
	      requestForward->ip, requestForward->port);
  }

  /* Quick test connection */
  if (valid_forward_route(requestForward->ip, requestForward->port) != TRUE) {
    retCode   = HttpStatus_ServiceUnavailable;
    statusStr = HttpStatusStr(retCode);
    log_error("Matching forward service unavailable. %d: %s", retCode,
	      statusStr);
    goto reply;
  } else {
    log_debug("Connection Test OK. Service available at ip: %s port: %s",
	      requestForward->ip, requestForward->port);
  }

  /* Step-3: setup request to forward */
  fRequest = create_forward_request(requestForward, requestPattern, request);
  if (fRequest == NULL) {
    retCode   = HttpStatus_InternalServerError;
    statusStr = HttpStatusStr(retCode);
    log_error("Internal routing error. %d: %s", retCode, statusStr);
    goto reply;
  } else {
    log_debug("Forward request sucessfully created");
  }

  /* Step-4: setup connection to the service */
  fResponse = (struct _u_response *)malloc(sizeof(struct _u_response));
  ulfius_init_response(fResponse);
  serviceResp = ulfius_send_http_request(fRequest, fResponse);
  if (serviceResp != U_OK) {
    retCode   = fResponse->status;
    statusStr = HttpStatusStr(retCode);
    log_error("Service response error: %d retCode: %d", retCode, statusStr);
  } else {
    log_debug("Request Forward to the service");
  }

  retCode = fResponse->status;

 reply:
  /* Step-5: response back to client */
  if (statusStr) {
    ulfius_set_string_body_response(response, retCode, statusStr);
  } else {
    ulfius_set_binary_body_response(response, retCode,
				  (void *)fResponse->binary_body,
				  fResponse->binary_body_length);
  }

  ulfius_clean_response(fResponse);
  free(fResponse);

  return U_CALLBACK_CONTINUE;
}

/*
 * callback_not_allowed -- 
 *
 */
int callback_not_allowed(const struct _u_request *request,
			 struct _u_response *response, void *user_data) {

  ulfius_set_string_body_response(response, 403, "Operation not allowed\n");
  return U_CALLBACK_CONTINUE;
}

/*
 * callback_default -- default callback for no-match
 *
 */
int callback_default(const struct _u_request *request,
                     struct _u_response *response, void *user_data) {

  ulfius_set_string_body_response(response, 404, "You are clearly high!\n");
  return U_CALLBACK_CONTINUE;
}
