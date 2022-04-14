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
 * add_service_entry --
 *
 */
static int add_service_entry(Router **router, Service *service) {

  Service *ptr=NULL;

  if ((*router)->services == NULL) {
    ptr = (*router)->services;
  } else {
    for (ptr=(*router)->services; ptr->next; ptr=ptr->next);
  }

  ptr = (Service *)calloc(1, sizeof(Service));
  if (ptr == NULL) {
    log_error("Error allocating memory of size: %lu", sizeof(Service));
    return FALSE;
  }

  uuid_generate(ptr->uuid);
  ptr->pattern = service->pattern;
  ptr->forward = service->forward;
  ptr->next    = NULL;

  return TRUE;
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

  int retCode=400;
  json_t *jreq=NULL;
  json_error_t jerr;
  Service *service=NULL;
  Router *router=NULL;
  json_t *jResp=NULL;
  char *jRespStr=NULL;

  router = (Router *)user_data;

  log_request(request);
  
  /* get json body */
  jreq = ulfius_get_json_body_request(request, &jerr);
  if (!jreq) {
    log_error("json error for POST %s: %s", EP_ROUTE, jerr.text);
  } else {
    deserialize_post_route_request(&service, jreq);
  }

  if (service) {
    add_service_entry(&router, service);
    serialize_post_route_response(&jResp, UUID, (void *)&service->uuid, NULL);
    retCode=200;
  } else {
    serialize_post_route_response(&jResp, ERROR, NULL, "Invalid request");
    retCode=400;
  }

  /* response back */
  jRespStr = json_dumps(jResp, 0);
  ulfius_set_string_body_response(response, retCode, jRespStr);

  if (jRespStr) free(jRespStr);
  json_decref(jResp);

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
